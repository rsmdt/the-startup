import { join } from 'path';
import type { InitOptions, InitResult } from '../types/config';

/**
 * File System Interface for dependency injection and testing
 */
interface FileSystem {
  mkdir(path: string, options: { recursive: boolean }): Promise<void>;
  readFile(path: string, encoding: string): Promise<string>;
  writeFile(path: string, content: string, encoding: string): Promise<void>;
  access(path: string): Promise<void>;
  stat(path: string): Promise<{ isDirectory(): boolean }>;
}

/**
 * Asset Provider Interface for accessing embedded template assets
 */
interface AssetProvider {
  getTemplateContent(template: string): string;
  getTemplateFileName(template: string): string;
}

/**
 * Initializer - Core template initialization engine
 *
 * Responsibilities:
 * - Copy template files (DOR, DOD, TASK-DOD) to docs/ directory
 * - Replace custom values in template content
 * - Support dry-run mode (preview without creating files)
 * - Support force mode (overwrite existing files)
 * - Create target directory if missing
 *
 * SDD References:
 * - Init interfaces (lines 614-629)
 * - Init flow (lines 822-868)
 */
export class Initializer {
  constructor(
    private readonly fs: FileSystem,
    private readonly assetProvider: AssetProvider,
    private readonly cwd: string = process.cwd()
  ) {}

  /**
   * Initialize templates by copying them to target directory
   *
   * @param options - Initialization options
   * @returns Result with created/previewed files or error
   */
  async initialize(options: InitOptions): Promise<InitResult> {
    try {
      const targetDir = this.resolveTargetDir(options);

      // Create target directory
      await this.ensureDirectory(targetDir);

      // Process each template
      const results = await Promise.all(
        options.templates.map(template =>
          this.processTemplate(template, targetDir, options)
        )
      );

      // Aggregate results
      const filesCreated: string[] = [];
      const filesPreview: string[] = [];
      const skipped: string[] = [];

      for (const result of results) {
        if (result.created) filesCreated.push(result.path);
        if (result.preview) filesPreview.push(result.path);
        if (result.skipped) skipped.push(result.path);
      }

      return {
        success: true,
        ...(options.dryRun
          ? { filesPreview }
          : { filesCreated, ...(skipped.length > 0 && { skipped }) }
        ),
      };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : String(error),
      };
    }
  }

  /**
   * Process a single template
   */
  private async processTemplate(
    template: string,
    targetDir: string,
    options: InitOptions
  ): Promise<{ path: string; created?: boolean; preview?: boolean; skipped?: boolean }> {
    const fileName = this.assetProvider.getTemplateFileName(template);
    const filePath = join(targetDir, fileName);

    // Check if file exists (unless force mode)
    if (!options.force && !options.dryRun) {
      const exists = await this.fileExists(filePath);
      if (exists) {
        return { path: filePath, skipped: true };
      }
    }

    // Get and process template content
    const content = this.assetProvider.getTemplateContent(template);
    const processedContent = this.replaceValues(content, options.customValues || {});

    // Dry-run: preview only
    if (options.dryRun) {
      return { path: filePath, preview: true };
    }

    // Write file
    await this.fs.writeFile(filePath, processedContent, 'utf-8');
    return { path: filePath, created: true };
  }

  /**
   * Replace placeholder values in template content
   */
  private replaceValues(content: string, values: Record<string, string>): string {
    let result = content;
    for (const [key, value] of Object.entries(values)) {
      const placeholder = `{{${key}}}`;
      result = result.replace(new RegExp(placeholder, 'g'), value);
    }
    return result;
  }

  /**
   * Resolve target directory (use provided or default to docs/)
   */
  private resolveTargetDir(options: InitOptions): string {
    return options.targetDirectory || join(this.cwd, 'docs');
  }

  /**
   * Ensure directory exists, create if missing
   */
  private async ensureDirectory(path: string): Promise<void> {
    await this.fs.mkdir(path, { recursive: true });
  }

  /**
   * Check if file exists
   */
  private async fileExists(path: string): Promise<boolean> {
    try {
      await this.fs.access(path);
      return true;
    } catch {
      return false;
    }
  }
}
