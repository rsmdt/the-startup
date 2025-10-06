import { join } from 'path';
import type { SpecOptions, SpecResult, SpecNumbering } from '../types/config';

/**
 * File System Interface for dependency injection and testing
 */
interface FileSystem {
  mkdir(path: string, options: { recursive: boolean }): Promise<void>;
  readdir(path: string): Promise<string[]>;
  writeFile(path: string, content: string, encoding: string): Promise<void>;
  readFile(path: string, encoding: string): Promise<string>;
  stat(path: string): Promise<{ isDirectory(): boolean }>;
}

/**
 * SpecGenerator - Core specification directory management
 *
 * Implements spec flow from SDD (lines 870-895):
 * - Create spec directories with auto-incrementing IDs (001, 002, 003...)
 * - Parse spec IDs from directory names (e.g., "004-typescript-npm-package-migration" → "004")
 * - Generate TOML output for --read flag (PRD line 206)
 * - Create template files for --add flag (PRD line 205)
 *
 * Key Responsibilities:
 * - Auto-increment spec ID based on existing directories
 * - Create spec directory structure (docs/specs/[id]-[name]/)
 * - Generate template files (PRD, SDD, PLAN, BRD)
 * - Output spec metadata in TOML format
 * - Parse directory names to extract spec IDs
 *
 * @example
 * const generator = new SpecGenerator(fs, 'docs/specs');
 * const result = await generator.createSpec({ name: 'user-authentication' });
 * // Creates: docs/specs/001-user-authentication/
 *
 * @example
 * const result = await generator.createSpec({ name: 'api-integration', template: 'PRD' });
 * // Creates: docs/specs/002-api-integration/PRD.md
 *
 * @example
 * const result = await generator.readSpec('001');
 * // Returns TOML: id = "001", name = "user-authentication", dir = "...", files = [...]
 */
export class SpecGenerator implements SpecNumbering {
  constructor(
    private fs: FileSystem,
    private specsDir: string = 'docs/specs'
  ) {}

  /**
   * Create a new spec directory with auto-incremented ID
   *
   * Implements SDD flow (lines 870-895):
   * 1. Read existing spec directories
   * 2. Calculate next ID (highest + 1)
   * 3. Create directory with format: [id]-[name]
   * 4. Optionally generate template file
   *
   * @param options - Spec creation options
   * @returns SpecResult with success status and metadata
   */
  async createSpec(options: SpecOptions): Promise<SpecResult> {
    try {
      const specId = await this.getNextSpecId();
      const sanitizedName = this.sanitizeName(options.name);
      const directory = join(this.specsDir, `${specId}-${sanitizedName}`);

      await this.fs.mkdir(directory, { recursive: true });

      const templateGenerated = options.template
        ? await this.generateTemplate(directory, options.template)
        : undefined;

      return this.createSuccessResult(specId, directory, { templateGenerated });
    } catch (error: any) {
      return this.createErrorResult(error);
    }
  }

  /**
   * Read spec metadata and output in TOML format
   *
   * Implements PRD requirement (line 206):
   * - Output spec state in TOML format
   * - Include: id, name, dir, files
   *
   * @param id - Spec ID to read (e.g., "001")
   * @returns SpecResult with TOML output
   */
  async readSpec(id: string): Promise<SpecResult> {
    try {
      const specDir = await this.findSpecDirectory(id);

      if (!specDir) {
        return this.createErrorResult(new Error(`Spec ${id} not found`));
      }

      const name = specDir.replace(/^\d{3}-/, '');
      const directory = join(this.specsDir, specDir);
      const files = await this.fs.readdir(directory);
      const toml = this.generateToml(id, name, directory, files);

      return this.createSuccessResult(id, directory, { toml });
    } catch (error: any) {
      return this.createErrorResult(error);
    }
  }

  /**
   * Calculate next spec ID by reading existing directories
   *
   * Implements SpecNumbering interface (SDD lines 648-652):
   * - Returns "001", "002", etc. (3-digit zero-padded)
   * - Finds highest existing ID and adds 1
   *
   * @returns Next spec ID (3-digit zero-padded)
   */
  async getNextSpecId(): Promise<string> {
    try {
      const dirs = await this.fs.readdir(this.specsDir);

      // Parse all directory names to extract IDs
      const ids = dirs
        .map((dir) => this.parseSpecId(dir))
        .filter((id): id is number => id !== null);

      // Find highest ID or default to 0
      const maxId = ids.length > 0 ? Math.max(...ids) : 0;

      // Return next ID with zero-padding
      return (maxId + 1).toString().padStart(3, '0');
    } catch (error) {
      // If directory doesn't exist or error reading, start at 001
      return '001';
    }
  }

  /**
   * Extract number from spec directory name
   *
   * Implements SpecNumbering interface (SDD lines 225-235):
   * - Examples: "004-feature-name" -> 4, "001-test" -> 1
   * - Returns null for invalid formats
   *
   * @param dirname - Directory name to parse
   * @returns Extracted number or null if invalid format
   */
  parseSpecId(dirname: string): number | null {
    // Match pattern: 3 digits optionally followed by hyphen (e.g., "001-" or "001")
    const match = dirname.match(/^(\d{3})(-|$)/);
    if (!match) {
      return null;
    }

    return parseInt(match[1], 10);
  }

  /**
   * Sanitize feature name for directory creation
   *
   * Converts to lowercase and replaces spaces/special chars with hyphens.
   *
   * @param name - Feature name to sanitize
   * @returns Sanitized name safe for directories
   */
  private sanitizeName(name: string): string {
    return name
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '');
  }

  /**
   * Generate template file in spec directory
   *
   * Creates template file with boilerplate content based on type.
   *
   * @param directory - Spec directory path
   * @param template - Template type (PRD, SDD, PLAN, BRD)
   * @returns Path to generated template file
   */
  private async generateTemplate(
    directory: string,
    template: 'PRD' | 'SDD' | 'PLAN' | 'BRD'
  ): Promise<string> {
    const templatePath = join(directory, `${template}.md`);
    const content = this.getTemplateContent(template);

    await this.fs.writeFile(templatePath, content, 'utf-8');

    return templatePath;
  }

  /**
   * Get template boilerplate content
   *
   * Returns minimal template structure for each document type.
   *
   * @param template - Template type
   * @returns Template content string
   */
  private getTemplateContent(template: 'PRD' | 'SDD' | 'PLAN' | 'BRD'): string {
    switch (template) {
      case 'PRD':
        return '# Product Requirements Document\n\n';
      case 'SDD':
        return '# System Design Document\n\n';
      case 'PLAN':
        return '# Implementation Plan\n\n';
      case 'BRD':
        return '# Business Requirements Document\n\n';
    }
  }

  /**
   * Generate TOML output for spec metadata
   *
   * Implements PRD requirement (line 206):
   * - Format: id = "001", name = "...", dir = "...", files = [...]
   *
   * @param id - Spec ID
   * @param name - Spec name
   * @param directory - Spec directory path
   * @param files - Files in spec directory
   * @returns TOML formatted string
   */
  private generateToml(
    id: string,
    name: string,
    directory: string,
    files: string[]
  ): string {
    const filesArray = files.map((f) => `"${f}"`).join(', ');
    return `id = "${id}"
name = "${name}"
dir = "${directory}"
files = [${filesArray}]`;
  }

  /**
   * Find spec directory matching given ID
   *
   * Searches existing spec directories for one matching the given ID.
   *
   * @param id - Spec ID to find (e.g., "001")
   * @returns Directory name or null if not found
   */
  private async findSpecDirectory(id: string): Promise<string | null> {
    const dirs = await this.fs.readdir(this.specsDir);
    const specDir = dirs.find((dir) => {
      const parsed = this.parseSpecId(dir);
      return parsed !== null && parsed.toString().padStart(3, '0') === id;
    });
    return specDir || null;
  }

  /**
   * Create success result with spec metadata
   *
   * @param specId - Spec ID
   * @param directory - Spec directory path
   * @param options - Optional additional fields (toml, templateGenerated)
   * @returns Success SpecResult
   */
  private createSuccessResult(
    specId: string,
    directory: string,
    options?: { toml?: string; templateGenerated?: string }
  ): SpecResult {
    return {
      success: true,
      specId,
      directory,
      ...options,
    };
  }

  /**
   * Create error result with error message
   *
   * @param error - Error object
   * @returns Error SpecResult
   */
  private createErrorResult(error: any): SpecResult {
    return {
      success: false,
      specId: '',
      directory: '',
      error: error.message,
    };
  }
}
