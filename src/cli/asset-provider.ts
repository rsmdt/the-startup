import { readdir } from 'fs/promises';
import { join } from 'path';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

/**
 * Asset file representation
 */
export interface AssetFile {
  category: 'agents' | 'commands' | 'templates' | 'rules' | 'outputStyles';
  sourcePath: string;
  relativePath: string;  // Path relative to category root (e.g., "the-analyst/requirements-analysis.md")
}

/**
 * Asset Provider Interface
 */
export interface AssetProvider {
  getAssetFiles(): Promise<AssetFile[]>;
  getSettingsTemplate(): any;
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
    ];

    for (const { category, path } of categories) {
      const categoryPath = join(this.assetsRoot, path);

      try {
        const files = await this.scanDirectory(categoryPath, categoryPath);
        for (const file of files) {
          // Compute relative path from category root
          // E.g., /abs/path/assets/claude/agents/the-analyst/file.md -> the-analyst/file.md
          const relativePath = file.replace(categoryPath + '/', '');

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
   * Get the settings template with hooks configuration
   */
  getSettingsTemplate(): any {
    return {
      hooks: {
        'user-prompt-submit': {
          command: '{{STARTUP_PATH}}/bin/statusline.sh',
        },
      },
    };
  }
}

/**
 * Create default asset provider instance
 */
export function createAssetProvider(): AssetProvider {
  return new FileSystemAssetProvider();
}
