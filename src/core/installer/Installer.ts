import { join, resolve, dirname } from 'path';
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
 * Simplified asset file metadata
 */
interface AssetFile {
  sourcePath: string;       // Absolute path in assets/
  relativePath: string;     // Path relative to category root
  targetCategory: 'claude' | 'startup';  // Install to .claude/ or .the-startup/
  isJson: boolean;          // Merge (true) or copy (false)
}

/**
 * Asset Provider Interface for accessing embedded assets
 */
interface AssetProvider {
  getAssetFiles(): Promise<AssetFile[]>;
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
      const totalSteps = assetFiles.length + 1; // +1 for lock file (assets include settings.json)
      let currentStep = 0;

      // 3. Create directories
      this.reportProgress('Creating directories', currentStep++, totalSteps);
      await this.createDirectories(startupPath, claudePath, options.selectedFiles);

      // 4. Copy/merge asset files
      for (const asset of assetFiles) {
        const category = this.extractCategory(asset.relativePath) || 'settings';
        this.reportProgress(`Processing ${category}`, currentStep++, totalSteps);
        await this.copyAsset(asset, startupPath, claudePath);
      }

      // 5. Create lock file with checksums
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
   * Converts absolute paths to tilde notation when under home directory.
   * This makes paths more portable across different user environments.
   *
   * @param absolutePath - Absolute path (e.g., /Users/john/.the-startup)
   * @returns Path with tilde notation if under home (e.g., ~/.the-startup), otherwise unchanged
   *
   * @example
   * // On macOS/Linux with HOME=/Users/john
   * toTildePath('/Users/john/.the-startup')  // → '~/.the-startup'
   * toTildePath('/Users/john/.claude')       // → '~/.claude'
   * toTildePath('/opt/the-startup')          // → '/opt/the-startup' (unchanged)
   */
  private toTildePath(absolutePath: string): string {
    // Only convert if path is under home directory
    if (!this.homeDir || !absolutePath.startsWith(this.homeDir)) {
      return absolutePath;
    }

    // Replace home directory with ~
    const relativePath = absolutePath.substring(this.homeDir.length);
    return `~${relativePath}`;
  }

  /**
   * Gets selected asset files based on user selections.
   *
   * Extracts category from relativePath:
   * - "agents/file.md" → category "agents"
   * - "settings.json" → no category, always included
   */
  private async getSelectedAssets(
    selectedFiles: InstallerOptions['selectedFiles']
  ): Promise<AssetFile[]> {
    const allAssets = await this.assetProvider.getAssetFiles();

    return allAssets.filter((asset) => {
      // Extract category from first path component
      const category = this.extractCategory(asset.relativePath);

      // Files without category (like settings.json) are always included
      if (!category) {
        return true;
      }

      // Filter based on user selection
      return selectedFiles[category] === true;
    });
  }

  /**
   * Extracts category from relative path.
   *
   * @param relativePath - Path like "agents/the-chief.md" or "settings.json"
   * @returns Category name or null if no category
   */
  private extractCategory(relativePath: string): keyof InstallerOptions['selectedFiles'] | null {
    const firstComponent = relativePath.split('/')[0];

    // Map directory names to category keys (handles both kebab-case and camelCase)
    const categoryMap: Record<string, keyof InstallerOptions['selectedFiles']> = {
      'agents': 'agents',
      'commands': 'commands',
      'templates': 'templates',
      'rules': 'rules',
      'output-styles': 'outputStyles',  // kebab-case directory name
      'outputStyles': 'outputStyles',    // camelCase fallback
    };

    // Return mapped category or null if not found
    return categoryMap[firstComponent] || null;
  }

  /**
   * Creates necessary installation directories.
   *
   * Just creates base directories - subdirectories are created
   * as needed during file copying.
   */
  private async createDirectories(
    startupPath: string,
    claudePath: string,
    _selectedFiles: InstallerOptions['selectedFiles']
  ): Promise<void> {
    // Create base directories
    await this.fs.mkdir(startupPath, { recursive: true });
    await this.fs.mkdir(claudePath, { recursive: true });

    // Subdirectories (agents, commands, templates, etc.) are created
    // automatically during copyAsset based on file paths
  }

  /**
   * Copies or merges a single asset file to its destination.
   *
   * .json files are merged with existing files.
   * All other files are copied (overwritten if exist).
   */
  private async copyAsset(
    asset: AssetFile,
    startupPath: string,
    claudePath: string
  ): Promise<void> {
    const destPath = this.getDestinationPath(asset, startupPath, claudePath);

    // Ensure parent directory exists
    const destDir = dirname(destPath);
    await this.fs.mkdir(destDir, { recursive: true });

    // Handle .json files specially (merge instead of overwrite)
    if (asset.isJson) {
      await this.mergeJsonFile(asset.sourcePath, destPath, startupPath, claudePath);
    } else {
      // Copy file (with placeholder replacement)
      await this.copyFileWithPlaceholders(asset.sourcePath, destPath, startupPath, claudePath);
    }

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
    const basePath = asset.targetCategory === 'claude' ? claudePath : startupPath;
    return join(basePath, asset.relativePath);
  }

  /**
   * Copies a file with placeholder replacement.
   *
   * Reads source file, replaces placeholders, writes to destination.
   */
  private async copyFileWithPlaceholders(
    sourcePath: string,
    destPath: string,
    startupPath: string,
    claudePath: string
  ): Promise<void> {
    const content = await this.fs.readFile(sourcePath, 'utf-8');
    const replaced = this.replacePlaceholders(content, startupPath, claudePath);
    await this.fs.writeFile(destPath, replaced, 'utf-8');
  }

  /**
   * Merges a JSON file with existing file at destination.
   *
   * Reads source JSON, replaces placeholders, merges with existing, writes result.
   */
  private async mergeJsonFile(
    sourcePath: string,
    destPath: string,
    startupPath: string,
    claudePath: string
  ): Promise<void> {
    // Read and parse source JSON
    const sourceContent = await this.fs.readFile(sourcePath, 'utf-8');
    const sourceJson = JSON.parse(sourceContent);

    // Replace placeholders in source
    const isWindows = process.platform === 'win32';
    const shellScriptExtension = isWindows ? '.ps1' : '.sh';

    const placeholders: PlaceholderMap = {
      STARTUP_PATH: this.toTildePath(startupPath),
      CLAUDE_PATH: this.toTildePath(claudePath),
      SHELL_SCRIPT_EXTENSION: shellScriptExtension,
    };

    // Merge with existing file
    await this.settingsMerger.mergeFullSettings(destPath, sourceJson, placeholders);
  }

  /**
   * Replaces placeholders in a string.
   */
  private replacePlaceholders(
    content: string,
    startupPath: string,
    claudePath: string
  ): string {
    const isWindows = process.platform === 'win32';
    const shellScriptExtension = isWindows ? '.ps1' : '.sh';

    let result = content;
    result = result.replace(/\{\{STARTUP_PATH\}\}/g, this.toTildePath(startupPath));
    result = result.replace(/\{\{CLAUDE_PATH\}\}/g, this.toTildePath(claudePath));
    result = result.replace(/\{\{SHELL_SCRIPT_EXTENSION\}\}/g, shellScriptExtension);

    return result;
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
