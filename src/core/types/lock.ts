/**
 * Lock file structure for tracking installed files and enabling idempotent reinstalls.
 *
 * The lock file supports two formats for backward compatibility:
 * - Version 1 (legacy): Simple string arrays without checksums
 * - Version 2 (current): FileEntry arrays with SHA-256 checksums
 *
 * Checksums enable:
 * - Detection of user-modified files during reinstall
 * - Idempotent reinstalls (skip unchanged files)
 * - Update command to detect changes and reinstall only modified files
 *
 * Migration strategy:
 * - New installations use v2 format with checksums
 * - Legacy lock files are auto-detected and migrated on read
 * - Migration is transparent to the user
 */

/**
 * File entry with optional checksum for backward compatibility.
 *
 * Checksum is SHA-256 hash of file contents, used to detect modifications.
 * Optional to support legacy lock files that don't have checksums.
 */
export interface FileEntry {
  /** File path relative to installation directory */
  path: string;

  /** SHA-256 hash of file contents (optional for backward compatibility) */
  checksum?: string;
}

/**
 * Lock file structure supporting both old (string[]) and new (FileEntry[]) formats.
 *
 * Version 2 format includes checksums for all installed files.
 * This enables idempotent reinstalls by comparing checksums to detect changes.
 */
export interface LockFile {
  /** Package version that installed these files */
  version: string;

  /** ISO timestamp of installation */
  installedAt: string;

  /**
   * Lock file format version (1 = old string format, 2 = with checksums)
   * If undefined, assume version 1 for backward compatibility
   */
  lockFileVersion?: number;

  /** Installed files organized by category */
  files: {
    /** Installed agent files with checksums */
    agents: FileEntry[];

    /** Installed command files with checksums */
    commands: FileEntry[];

    /** Installed template files with checksums */
    templates: FileEntry[];

    /** Installed rule files with checksums */
    rules: FileEntry[];

    /** Installed output style files with checksums */
    outputStyles: FileEntry[];

    /** Binary file with checksum */
    binary: FileEntry;
  };
}

/**
 * Legacy lock file format (backward compatibility).
 *
 * This format was used before checksums were added.
 * Files are stored as simple string arrays without checksum tracking.
 *
 * When reading a legacy lock file:
 * 1. Auto-detect format by checking if files.agents[0] is a string
 * 2. Convert to FileEntry[] format with checksum: undefined
 * 3. Mark as lockFileVersion: 1
 */
export interface LegacyLockFile {
  /** Package version that installed these files */
  version: string;

  /** ISO timestamp of installation */
  installedAt: string;

  /** Installed files organized by category (string arrays) */
  files: {
    /** Installed agent file paths */
    agents: string[];

    /** Installed command file paths */
    commands: string[];

    /** Installed template file paths */
    templates: string[];

    /** Installed rule file paths */
    rules: string[];

    /** Installed output style file paths */
    outputStyles: string[];

    /** Binary file path */
    binary: string;
  };
}
