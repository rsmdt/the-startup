/**
 * Configuration types for command options and core business logic.
 *
 * This file defines:
 * - CLI command option flags (Commander.js flag definitions)
 * - Core installer options and results
 * - Init command options and results
 * - Spec command options and results
 * - Spec numbering interface
 */

/**
 * Install command flags (Commander.js options).
 *
 * Flags control installation behavior:
 * - --local: Non-interactive, uses default paths (./.the-startup, ~/.claude)
 * - --yes: Auto-confirms all prompts with recommended paths
 *
 * Flags can be combined: install --local --yes for fully non-interactive installation.
 */
export interface InstallCommandOptions {
  /** Skip prompts, use local paths (./.the-startup, ~/.claude) */
  local?: boolean;

  /** Auto-confirm all prompts with recommended paths */
  yes?: boolean;
}

/**
 * Uninstall command flags (Commander.js options).
 *
 * Flags control what gets removed during uninstall:
 * - --keep-logs: Preserve .the-startup/logs directory
 * - --keep-settings: Don't modify ~/.claude/settings.json
 *
 * Flags can be combined: uninstall --keep-logs --keep-settings for minimal cleanup.
 */
export interface UninstallCommandOptions {
  /** Preserve .the-startup/logs directory */
  keepLogs?: boolean;

  /** Don't modify ~/.claude/settings.json */
  keepSettings?: boolean;
}

/**
 * Init command flags (Commander.js options).
 *
 * Flags control template initialization behavior:
 * - --dry-run: Preview changes without writing files
 * - --force: Overwrite existing files without prompting
 *
 * Flags are mutually exclusive: cannot use --dry-run --force together.
 */
export interface InitCommandOptions {
  /** Preview changes without writing files */
  dryRun?: boolean;

  /** Overwrite existing files without prompting */
  force?: boolean;
}

/**
 * Spec command flags (Commander.js options).
 *
 * Flags control spec directory and template generation:
 * - --add: Generate template file in spec directory
 * - --read: Output spec state in TOML format
 */
export interface SpecCommandOptions {
  /** Generate template file in spec directory */
  add?: 'product-requirements' | 'solution-design' | 'implementation-plan' | 'business-requirements';

  /** Output spec state in TOML format */
  read?: boolean;
}

/**
 * Core installer options (framework-agnostic business logic).
 *
 * These options are passed to the Installer core after CLI/UI processing.
 */
export interface InstallerOptions {
  /** Installation directory (e.g., .the-startup) */
  startupPath: string;

  /** Claude config directory (e.g., ~/.claude) */
  claudePath: string;

  /** Selected file categories to install */
  selectedFiles: {
    /** Install agent files */
    agents: boolean;

    /** Install command files */
    commands: boolean;

    /** Install template files */
    templates: boolean;

    /** Install rule files */
    rules: boolean;

    /** Install output style files */
    outputStyles: boolean;
  };
}

/**
 * Installation result returned by Installer core.
 *
 * Provides success status and details of installed files or errors.
 */
export interface InstallResult {
  /** Whether installation completed successfully */
  success: boolean;

  /** List of installed file paths (absolute paths) */
  installedFiles: string[];

  /** Error messages if installation failed */
  errors?: string[];
}

/**
 * Init command options (core business logic).
 *
 * These options are passed to the Init core after prompt processing.
 */
export interface InitOptions {
  /** Templates to initialize */
  templates: ('DOR' | 'DOD' | 'TASK-DOD')[];

  /** User-provided custom values for template placeholders */
  customValues?: Record<string, string>;

  /** Target directory (default: docs/) */
  targetDirectory?: string;

  /** Preview mode (from --dry-run flag) */
  dryRun?: boolean;

  /** Overwrite mode (from --force flag) */
  force?: boolean;
}

/**
 * Init command result returned by Init core.
 *
 * Provides success status and details of created/previewed files.
 */
export interface InitResult {
  /** Whether initialization completed successfully */
  success: boolean;

  /** Paths of created files (empty if dry-run) */
  filesCreated?: string[];

  /** Paths that would be created (dry-run only) */
  filesPreview?: string[];

  /** Files skipped due to existing (no --force) */
  skipped?: string[];

  /** Error message if failed */
  error?: string;
}

/**
 * Spec command options (core business logic).
 *
 * These options are passed to the Spec core after argument parsing.
 */
export interface SpecOptions {
  /** Feature name (e.g., "user-authentication") */
  name: string;

  /** Template to generate (from --add flag) */
  template?: 'product-requirements' | 'solution-design' | 'implementation-plan' | 'business-requirements';

  /** Existing spec ID (e.g., "004") */
  specId?: string;
}

/**
 * Spec command result returned by Spec core.
 *
 * Provides success status and details of created spec directory.
 */
export interface SpecResult {
  /** Whether spec creation completed successfully */
  success: boolean;

  /** Generated or provided ID (e.g., "004") */
  specId: string;

  /** Spec directory path */
  directory: string;

  /** TOML output (for --read flag) */
  toml?: string;

  /** Template file path (for --add flag) */
  templateGenerated?: string;

  /** Error message if failed */
  error?: string;
}

/**
 * Spec numbering algorithm interface.
 *
 * Calculates next spec ID by reading existing directories.
 * IDs are 3-digit zero-padded numbers (001, 002, 003, etc.).
 */
export interface SpecNumbering {
  /**
   * Calculate next spec ID by reading existing spec directories.
   *
   * Returns "001", "002", etc. (3-digit zero-padded).
   */
  getNextSpecId(): Promise<string>;

  /**
   * Extract number from spec directory name.
   *
   * Examples:
   * - "004-feature-name" -> 4
   * - "001-test" -> 1
   * - "invalid" -> null
   *
   * @param dirname - Directory name to parse
   * @returns Extracted number or null if invalid format
   */
  parseSpecId(dirname: string): number | null;
}
