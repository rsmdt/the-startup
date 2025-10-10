import { readdir } from 'fs/promises';
import { join } from 'path';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

/**
 * Asset file representation
 */
export interface AssetFile {
  category: 'agents' | 'commands' | 'templates' | 'rules' | 'outputStyles' | 'bin';
  sourcePath: string;
  relativePath: string;  // Path relative to category root (e.g., "the-analyst/requirements-analysis.md")
}

/**
 * Asset Provider Interface
 */
export interface AssetProvider {
  getAssetFiles(): Promise<AssetFile[]>;
  getSettingsTemplate(): Promise<any>;
  getSettingsLocalTemplate(): Promise<any>;
  getAssetsRoot(): string;
}

/**
 * FileSystemAssetProvider - Reads assets from the file system
 *
 * This provider scans the assets/ directory to find all agent definitions,
 * commands, templates, rules, and output styles.
 *
 * Works in both development (npm run dev) and production (npx) by resolving
 * the package root from the current module location.
 */
export class FileSystemAssetProvider implements AssetProvider {
  private assetsRoot: string;

  constructor() {
    // Get the directory of this file
    const currentFile = fileURLToPath(import.meta.url);
    const currentDir = dirname(currentFile);

    // Resolve to package root, then assets/
    // In development (npm run dev): currentDir is src/cli/ -> go up 2 levels
    // In tests: currentDir is src/cli/ -> go up 2 levels
    // In production (bundled): currentDir is dist/ -> go up 1 level
    //
    // We detect which context we're in by checking if we're in 'src/' or 'dist/'
    const isSourceContext = currentDir.includes('/src/');
    const levelsUp = isSourceContext ? 2 : 1;

    const packageRoot = join(currentDir, ...Array(levelsUp).fill('..'));
    this.assetsRoot = join(packageRoot, 'assets');
  }

  /**
   * Get the assets root directory path
   */
  getAssetsRoot(): string {
    return this.assetsRoot;
  }

  /**
   * Scan and return all asset files organized by category
   * For bin category, only includes the script matching current OS
   */
  async getAssetFiles(): Promise<AssetFile[]> {
    const assets: AssetFile[] = [];

    // Scan each asset category
    const categories = [
      { category: 'agents' as const, path: 'claude/agents' },
      { category: 'commands' as const, path: 'claude/commands' },
      { category: 'templates' as const, path: 'the-startup/templates' },
      { category: 'rules' as const, path: 'the-startup/rules' },
      { category: 'outputStyles' as const, path: 'claude/output-styles' },
      { category: 'bin' as const, path: 'the-startup/bin' },
    ];

    for (const { category, path } of categories) {
      const categoryPath = join(this.assetsRoot, path);

      try {
        const files = await this.scanDirectory(categoryPath, categoryPath);
        for (const file of files) {
          // Compute relative path from category root
          // E.g., /abs/path/assets/claude/agents/the-analyst/file.md -> the-analyst/file.md
          const relativePath = file.replace(categoryPath + '/', '');

          // For bin category, only include the script matching current OS
          if (category === 'bin') {
            const isWindows = process.platform === 'win32';
            const expectedExt = isWindows ? '.ps1' : '.sh';
            if (!file.endsWith(expectedExt)) {
              continue; // Skip script for different OS
            }
          }

          assets.push({
            category,
            sourcePath: file,
            relativePath,
          });
        }
      } catch (error) {
        // Category directory might not exist, skip it
        console.warn(`Warning: Could not scan ${path}:`, error instanceof Error ? error.message : error);
      }
    }

    return assets;
  }

  /**
   * Recursively scan a directory for all files
   *
   * @param dir - Directory to scan
   * @param baseDir - Base directory for relative path calculation
   * @returns Array of file paths relative to assets root
   */
  private async scanDirectory(dir: string, baseDir: string): Promise<string[]> {
    const files: string[] = [];

    try {
      const entries = await readdir(dir, { withFileTypes: true });

      for (const entry of entries) {
        const fullPath = join(dir, entry.name);

        if (entry.isDirectory()) {
          // Recursively scan subdirectories
          const subFiles = await this.scanDirectory(fullPath, baseDir);
          files.push(...subFiles);
        } else if (entry.isFile()) {
          // Add file with path relative to base directory
          files.push(fullPath);
        }
      }
    } catch (error) {
      // If directory doesn't exist or can't be read, return empty array
      console.warn(`Warning: Could not read directory ${dir}:`, error instanceof Error ? error.message : error);
    }

    return files;
  }

  /**
   * Get the settings template with complete configuration (permissions, statusLine, hooks).
   *
   * Reads from assets/claude/settings.json and adds hooks section programmatically
   * to match the statusLine command.
   *
   * Uses {{SHELL_SCRIPT_EXTENSION}} placeholder that gets replaced with .sh or .ps1
   * during installation based on the user's operating system.
   */
  async getSettingsTemplate(): Promise<any> {
    const { readFile } = await import('fs/promises');

    // Read base settings from assets/claude/settings.json
    const settingsPath = join(this.assetsRoot, 'claude', 'settings.json');
    const settingsContent = await readFile(settingsPath, 'utf-8');
    const baseSettings = JSON.parse(settingsContent);

    // Add hooks section programmatically (matches statusLine command)
    if (!baseSettings.hooks) {
      baseSettings.hooks = {};
    }

    // Use the same command as statusLine for the user-prompt-submit hook
    if (baseSettings.statusLine?.command) {
      baseSettings.hooks['user-prompt-submit'] = {
        command: baseSettings.statusLine.command,
      };
    }

    return baseSettings;
  }

  /**
   * Get the settings.local.json template.
   *
   * Reads from assets/claude/settings.local.json which contains local
   * configuration like outputStyle that should be merged into the user's
   * settings.local.json file.
   *
   * @returns Settings local template object
   */
  async getSettingsLocalTemplate(): Promise<any> {
    const { readFile } = await import('fs/promises');

    try {
      const settingsLocalPath = join(this.assetsRoot, 'claude', 'settings.local.json');
      const settingsLocalContent = await readFile(settingsLocalPath, 'utf-8');
      return JSON.parse(settingsLocalContent);
    } catch (error) {
      // If settings.local.json doesn't exist, return empty object
      if ((error as NodeJS.ErrnoException).code === 'ENOENT') {
        return {};
      }
      throw error;
    }
  }
}

/**
 * Create default asset provider instance
 */
export function createAssetProvider(): AssetProvider {
  return new FileSystemAssetProvider();
}
