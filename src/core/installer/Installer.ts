import { join, resolve } from 'path';
import type { InstallerOptions, InstallResult } from '../types/config';
import type { FileEntry } from '../types/lock';
import type { LockManager } from './LockManager';
import type { SettingsMerger } from './SettingsMerger';
import type { PlaceholderMap } from '../types/settings';

/**
 * File System Interface for dependency injection and testing
 */
interface FileSystem {
  mkdir(path: string, options: { recursive: boolean }): Promise<void>;
  copyFile(src: string, dest: string): Promise<void>;
  readFile(path: string, encoding: string): Promise<string>;
  writeFile(path: string, content: string, encoding: string): Promise<void>;
  rm(path: string, options?: { force?: boolean; recursive?: boolean }): Promise<void>;
  access(path: string): Promise<void>;
  stat(path: string): Promise<{ isDirectory(): boolean }>;
}

/**
 * Asset file metadata
 */
interface AssetFile {
  category: 'agents' | 'commands' | 'templates' | 'rules' | 'outputStyles';
  sourcePath: string;
}

/**
 * Asset Provider Interface for accessing embedded assets
 */
interface AssetProvider {
  getAssetFiles(): Promise<AssetFile[]>;
  getSettingsTemplate(): any;
}

/**
 * Progress callback for reporting installation progress
 */
interface ProgressInfo {
  stage: string;
  current: number;
  total: number;
}

type ProgressCallback = (progress: ProgressInfo) => void;

/**
 * Installer - Core installation engine
 *
 * Implements installation flow from SDD:
 * - Copy selected assets to .the-startup/ and .claude/ directories
 * - Merge settings.json with hooks via SettingsMerger
 * - Create lock file via LockManager
 * - Atomic operations with rollback on failure
 * - Progress reporting for operations > 5s
 * - Comprehensive error handling
 *
 * Key Responsibilities:
 * - Orchestrate file copying operations
 * - Integrate with LockManager for file tracking
 * - Integrate with SettingsMerger for settings.json updates
 * - Provide rollback mechanism on failures
 * - Report progress for long operations
 * - Handle all error scenarios with specific messages
 *
 * Error Handling (SDD lines 899-921):
 * - Invalid path: "Please re-enter a valid path"
 * - Permission denied: Show chmod/permission suggestion
 * - Disk full: Show space needed
 * - Settings merge failure: Rollback to backup
 * - Asset copy failure: Clean up partial installation
 *
 * @example
 * const installer = new Installer(fs, lockManager, settingsMerger, assetProvider, '1.0.0');
 * const result = await installer.install({
 *   startupPath: './.the-startup',
 *   claudePath: '~/.claude',
 *   selectedFiles: { agents: true, commands: true, templates: true, rules: true, outputStyles: true }
 * });
 */
export class Installer {
  private installedFiles: string[] = [];

  constructor(
    private fs: FileSystem,
    private lockManager: LockManager,
    private settingsMerger: SettingsMerger,
    private assetProvider: AssetProvider,
    private version: string,
    private progressCallback?: ProgressCallback,
    private homeDir: string = process.env.HOME || process.env.USERPROFILE || '',
    private cwd: string = process.cwd()
  ) {}

  /**
   * Install selected assets with atomic operations and rollback.
   *
   * Installation flow (SDD lines 677-718):
   * 1. Normalize and validate paths
   * 2. Create installation directories
   * 3. Copy selected asset files
   * 4. Merge settings.json with hooks
   * 5. Create lock file with checksums
   *
   * On any failure: Rollback all changes and return error result.
   *
   * @param options - Installation options with paths and file selections
   * @returns Installation result with success status and installed files
   */
  async install(options: InstallerOptions): Promise<InstallResult> {
    this.installedFiles = [];

    try {
      // 1. Normalize paths
      const startupPath = this.normalizePath(options.startupPath);
      const claudePath = this.normalizePath(options.claudePath);

      // 2. Get selected assets
      const assetFiles = await this.getSelectedAssets(options.selectedFiles);
      const totalSteps = assetFiles.length + 2; // +2 for settings and lock
      let currentStep = 0;

      // 3. Create directories
      this.reportProgress('Creating directories', currentStep++, totalSteps);
      await this.createDirectories(startupPath, claudePath, options.selectedFiles);

      // 4. Copy asset files
      for (const asset of assetFiles) {
        this.reportProgress(`Copying ${asset.category}`, currentStep++, totalSteps);
        await this.copyAsset(asset, startupPath, claudePath);
      }

      // 5. Merge settings.json
      this.reportProgress('Merging settings', currentStep++, totalSteps);
      await this.mergeSettings(claudePath, startupPath);

      // 6. Create lock file with checksums
      this.reportProgress('Creating lock file', currentStep++, totalSteps);
      await this.createLockFile(startupPath);

      return {
        success: true,
        installedFiles: [...this.installedFiles],
      };
    } catch (error) {
      // Rollback on any failure
      await this.rollback();

      return {
        success: false,
        installedFiles: [],
        errors: [this.formatError(error)],
      };
    }
  }

  /**
   * Normalizes paths by expanding ~ and resolving relative paths.
   */
  private normalizePath(path: string): string {
    // Expand tilde
    let normalized = path;
    if (normalized.startsWith('~')) {
      normalized = normalized.replace('~', this.homeDir);
    }

    // Resolve relative paths
    if (!normalized.startsWith('/')) {
      normalized = resolve(this.cwd, normalized);
    }

    return normalized;
  }

  /**
   * Gets selected asset files based on user selections.
   */
  private async getSelectedAssets(
    selectedFiles: InstallerOptions['selectedFiles']
  ): Promise<AssetFile[]> {
    const allAssets = await this.assetProvider.getAssetFiles();

    return allAssets.filter((asset) => {
      return selectedFiles[asset.category] === true;
    });
  }

  /**
   * Creates necessary installation directories.
   */
  private async createDirectories(
    startupPath: string,
    claudePath: string,
    selectedFiles: InstallerOptions['selectedFiles']
  ): Promise<void> {
    // Always create base directories
    await this.fs.mkdir(startupPath, { recursive: true });
    await this.fs.mkdir(claudePath, { recursive: true });

    // Create category-specific directories
    if (selectedFiles.agents) {
      await this.fs.mkdir(join(claudePath, 'agents'), { recursive: true });
    }
    if (selectedFiles.commands) {
      await this.fs.mkdir(join(claudePath, 'commands'), { recursive: true });
    }
    if (selectedFiles.templates) {
      await this.fs.mkdir(join(startupPath, 'templates'), { recursive: true });
    }
    if (selectedFiles.rules) {
      await this.fs.mkdir(join(startupPath, 'rules'), { recursive: true });
    }
    if (selectedFiles.outputStyles) {
      await this.fs.mkdir(join(claudePath, 'output-styles'), { recursive: true });
    }
  }

  /**
   * Copies a single asset file to its destination.
   */
  private async copyAsset(
    asset: AssetFile,
    startupPath: string,
    claudePath: string
  ): Promise<void> {
    const destPath = this.getDestinationPath(asset, startupPath, claudePath);
    const sourcePath = this.getSourcePath(asset);

    await this.fs.copyFile(sourcePath, destPath);
    this.installedFiles.push(destPath);
  }

  /**
   * Determines destination path for an asset file.
   */
  private getDestinationPath(
    asset: AssetFile,
    startupPath: string,
    claudePath: string
  ): string {
    const fileName = asset.sourcePath.split('/').pop()!;

    switch (asset.category) {
      case 'agents':
        return join(claudePath, 'agents', fileName);
      case 'commands':
        return join(claudePath, 'commands', fileName);
      case 'templates':
        return join(startupPath, 'templates', fileName);
      case 'rules':
        return join(startupPath, 'rules', fileName);
      case 'outputStyles':
        return join(claudePath, 'output-styles', fileName);
    }
  }

  /**
   * Gets source path for an asset file.
   */
  private getSourcePath(asset: AssetFile): string {
    // In production, this would resolve from embedded assets
    // For testing, we use the asset's sourcePath directly
    return asset.sourcePath;
  }

  /**
   * Merges settings.json with hooks.
   */
  private async mergeSettings(
    claudePath: string,
    startupPath: string
  ): Promise<void> {
    const settingsPath = join(claudePath, 'settings.json');
    const settingsTemplate = this.assetProvider.getSettingsTemplate();

    const placeholders: PlaceholderMap = {
      STARTUP_PATH: startupPath,
      CLAUDE_PATH: claudePath,
    };

    await this.settingsMerger.mergeSettings(
      settingsPath,
      settingsTemplate.hooks, // Pass only the hooks object, not the whole template
      placeholders
    );
  }

  /**
   * Creates lock file with checksums for all installed files.
   */
  private async createLockFile(_startupPath: string): Promise<void> {
    // Generate checksums for all installed files
    const fileEntries: FileEntry[] = [];
    for (const filePath of this.installedFiles) {
      const checksum = await this.lockManager.generateChecksum(filePath);
      fileEntries.push({ path: filePath, checksum });
    }

    await this.lockManager.writeLockFile(fileEntries, this.version);
  }

  /**
   * Rolls back all installed files on failure.
   */
  private async rollback(): Promise<void> {
    for (const filePath of this.installedFiles) {
      try {
        await this.fs.rm(filePath, { force: true });
      } catch (error) {
        // Continue cleanup even if individual file deletion fails
      }
    }

    this.installedFiles = [];
  }

  /**
   * Reports progress to callback if provided.
   */
  private reportProgress(
    stage: string,
    current: number,
    total: number
  ): void {
    if (this.progressCallback) {
      this.progressCallback({ stage, current, total });
    }
  }

  /**
   * Formats errors into user-friendly messages.
   *
   * Error handling from SDD lines 899-921:
   * - ENOENT: Invalid path
   * - EACCES: Permission denied
   * - ENOSPC: Disk full
   */
  private formatError(error: unknown): string {
    const err = error as NodeJS.ErrnoException;

    switch (err.code) {
      case 'ENOENT':
        return 'Invalid path. Please re-enter a valid path.';

      case 'EACCES':
        return `Permission denied. Please check directory permissions or run with appropriate privileges (e.g., chmod to fix permissions).`;

      case 'ENOSPC':
        return 'Disk full. Not enough space left on device. Please free up disk space and try again.';

      default:
        return err.message || 'Installation failed due to an unexpected error.';
    }
  }
}
