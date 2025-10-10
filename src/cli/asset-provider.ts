import { readdir, stat } from 'fs/promises';
import { join, relative } from 'path';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

/**
 * Simplified asset file representation
 */
export interface AssetFile {
  sourcePath: string;       // Absolute path in assets/ directory
  relativePath: string;     // Path relative to category root (preserves structure)
  targetCategory: 'claude' | 'startup';  // Whether to install to .claude/ or .the-startup/
  isJson: boolean;          // Whether to merge (true) or overwrite (false)
}

/**
 * Asset Provider Interface
 */
export interface AssetProvider {
  getAssetFiles(): Promise<AssetFile[]>;
  getAssetsRoot(): string;
}

/**
 * Simplified FileSystemAssetProvider - Automatically scans assets/ directory
 *
 * Architecture:
 * - Scans assets/claude/ and assets/the-startup/ recursively
 * - Auto-detects .json files for merging
 * - Filters OS-specific files (.sh on Unix, .ps1 on Windows)
 * - No hardcoded categories or special template methods!
 *
 * File handling:
 * - .json files: Marked for merge (isJson: true)
 * - .sh files: Only included on non-Windows (OS filtering)
 * - .ps1 files: Only included on Windows (OS filtering)
 * - Everything else: Copied as-is
 */
export class FileSystemAssetProvider implements AssetProvider {
  private assetsRoot: string;

  constructor() {
    // Get the directory of this file
    const currentFile = fileURLToPath(import.meta.url);
    const currentDir = dirname(currentFile);

    // Resolve to package root, then assets/
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
   * Scan and return all asset files with automatic discovery
   *
   * Scans:
   * - assets/claude/ → Install to ~/.claude/
   * - assets/the-startup/ → Install to .the-startup/
   *
   * Filters:
   * - .sh files: Only on Unix/Mac
   * - .ps1 files: Only on Windows
   * - .json files: Marked for merging
   */
  async getAssetFiles(): Promise<AssetFile[]> {
    const assets: AssetFile[] = [];
    const isWindows = process.platform === 'win32';

    // Scan claude assets
    const claudePath = join(this.assetsRoot, 'claude');
    const claudeFiles = await this.scanDirectory(claudePath);
    for (const file of claudeFiles) {
      // Skip OS-specific files for wrong OS
      if (!this.shouldIncludeFile(file, isWindows)) {
        continue;
      }

      assets.push({
        sourcePath: file,
        relativePath: relative(claudePath, file),
        targetCategory: 'claude',
        isJson: file.endsWith('.json'),
      });
    }

    // Scan the-startup assets
    const startupPath = join(this.assetsRoot, 'the-startup');
    const startupFiles = await this.scanDirectory(startupPath);
    for (const file of startupFiles) {
      // Skip OS-specific files for wrong OS
      if (!this.shouldIncludeFile(file, isWindows)) {
        continue;
      }

      assets.push({
        sourcePath: file,
        relativePath: relative(startupPath, file),
        targetCategory: 'startup',
        isJson: file.endsWith('.json'),
      });
    }

    return assets;
  }

  /**
   * Determines if a file should be included based on OS
   *
   * @param filePath - Path to check
   * @param isWindows - Whether running on Windows
   * @returns True if file should be included
   */
  private shouldIncludeFile(filePath: string, isWindows: boolean): boolean {
    // .sh files: Only on Unix/Mac
    if (filePath.endsWith('.sh')) {
      return !isWindows;
    }

    // .ps1 files: Only on Windows
    if (filePath.endsWith('.ps1')) {
      return isWindows;
    }

    // All other files: Always include
    return true;
  }

  /**
   * Recursively scan a directory for all files
   *
   * @param dir - Directory to scan
   * @returns Array of absolute file paths
   */
  private async scanDirectory(dir: string): Promise<string[]> {
    const files: string[]= [];

    try {
      const entries = await readdir(dir, { withFileTypes: true });

      for (const entry of entries) {
        const fullPath = join(dir, entry.name);

        if (entry.isDirectory()) {
          // Recursively scan subdirectories
          const subFiles = await this.scanDirectory(fullPath);
          files.push(...subFiles);
        } else if (entry.isFile()) {
          // Add file
          files.push(fullPath);
        }
      }
    } catch (error) {
      // If directory doesn't exist or can't be read, return empty array
      console.warn(`Warning: Could not read directory ${dir}:`, error instanceof Error ? error.message : error);
    }

    return files;
  }
}

/**
 * Create default asset provider instance
 */
export function createAssetProvider(): AssetProvider {
  return new FileSystemAssetProvider();
}
