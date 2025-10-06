import { promises as fs } from 'fs';
import * as crypto from 'crypto';
import { LockFile, FileEntry, LegacyLockFile } from '../types/lock';

/**
 * LockManager manages lock file operations with backward compatibility.
 *
 * Responsibilities:
 * - Read lock files (v1 legacy format or v2 with checksums)
 * - Write lock files (always v2 format with SHA-256 checksums)
 * - Auto-migrate v1 lock files to v2 format on read
 * - Generate SHA-256 checksums for files
 * - Compare checksums to detect file modifications
 * - Determine which files need reinstalling based on checksum comparison
 *
 * Backward Compatibility:
 * - Supports v1 lock files (string[] format without checksums)
 * - Auto-detects format by checking if files.agents[0] is string
 * - Migrates v1 to normalized v2 structure with undefined checksums
 * - Always writes v2 format for new installations
 */
export class LockManager {
  constructor(private readonly lockFilePath: string) {}

  /**
   * Reads lock file from disk with backward compatibility.
   *
   * Auto-detects and migrates v1 legacy format to v2 format.
   * V1 files have string arrays, v2 files have FileEntry arrays.
   *
   * @returns LockFile object or null if file doesn't exist
   * @throws Error if file exists but cannot be parsed
   */
  async readLockFile(): Promise<LockFile | null> {
    try {
      const content = await fs.readFile(this.lockFilePath, 'utf-8');
      const parsed = JSON.parse(content);

      // Detect format version
      if (this.isV2Format(parsed)) {
        return parsed as LockFile;
      }

      // Migrate v1 to v2 format
      return this.migrateV1ToV2(parsed as LegacyLockFile);
    } catch (error) {
      if ((error as NodeJS.ErrnoException).code === 'ENOENT') {
        return null;
      }
      throw error;
    }
  }

  /**
   * Writes lock file to disk in v2 format with checksums.
   *
   * Always writes v2 format regardless of what was read.
   * Categorizes files by path patterns into agents, commands, templates, etc.
   *
   * @param fileEntries - Array of file entries with paths and checksums
   * @param version - Package version for the lock file
   * @throws Error if write operation fails
   */
  async writeLockFile(
    fileEntries: FileEntry[],
    version: string
  ): Promise<void> {
    const categorized = this.categorizeFiles(fileEntries);

    const lockFile: LockFile = {
      version,
      installedAt: new Date().toISOString(),
      lockFileVersion: 2,
      files: categorized,
    };

    await fs.writeFile(
      this.lockFilePath,
      JSON.stringify(lockFile, null, 2),
      'utf-8'
    );
  }

  /**
   * Generates SHA-256 checksum for a file.
   *
   * @param filePath - Absolute path to file
   * @returns SHA-256 hash as hex string (64 characters)
   * @throws Error if file cannot be read
   */
  async generateChecksum(filePath: string): Promise<string> {
    const content = await fs.readFile(filePath);
    return crypto.createHash('sha256').update(content).digest('hex');
  }

  /**
   * Compares file's current checksum with stored checksum.
   *
   * @param fileEntry - File entry with path and optional checksum
   * @returns true if checksums match, false if different or no checksum
   * @throws Error if file cannot be read
   */
  async compareChecksums(fileEntry: FileEntry): Promise<boolean> {
    if (!fileEntry.checksum) {
      return false;
    }

    const currentChecksum = await this.generateChecksum(fileEntry.path);
    return currentChecksum === fileEntry.checksum;
  }

  /**
   * Determines which files need reinstalling based on checksum comparison.
   *
   * Rules:
   * - If no lock file exists: reinstall all files
   * - If file not in lock: mark for install (new file)
   * - If file in lock with matching checksum: skip (unchanged)
   * - If file in lock with different checksum: mark for reinstall (modified)
   * - If file in lock without checksum (legacy): mark for reinstall
   *
   * @param targetFiles - Files to check for reinstall
   * @returns Array of files that need installing/reinstalling
   */
  async getFilesNeedingReinstall(
    targetFiles: FileEntry[]
  ): Promise<FileEntry[]> {
    const existingLock = await this.readLockFile();

    if (!existingLock) {
      return targetFiles;
    }

    const existingFilesMap = this.buildFileMap(existingLock);
    const needReinstall: FileEntry[] = [];

    for (const targetFile of targetFiles) {
      if (this.shouldReinstall(targetFile, existingFilesMap)) {
        needReinstall.push(targetFile);
      }
    }

    return needReinstall;
  }

  /**
   * Determines if a single file needs reinstalling.
   *
   * @param targetFile - File to check
   * @param existingFilesMap - Map of existing files from lock
   * @returns true if file needs reinstalling
   */
  private shouldReinstall(
    targetFile: FileEntry,
    existingFilesMap: Map<string, FileEntry>
  ): boolean {
    const existingEntry = existingFilesMap.get(targetFile.path);

    if (!existingEntry) {
      return true; // New file not in lock
    }

    if (!existingEntry.checksum || !targetFile.checksum) {
      return true; // Legacy entry or target without checksum
    }

    return existingEntry.checksum !== targetFile.checksum; // Modified file
  }

  /**
   * Detects if lock file is v2 format.
   *
   * Detection logic:
   * 1. If lockFileVersion === 2, it's v2
   * 2. If files.agents[0] is object with 'path' field, it's v2
   * 3. Otherwise, it's v1
   */
  private isV2Format(lock: unknown): boolean {
    const lockObj = lock as Record<string, unknown>;

    if (lockObj.lockFileVersion === 2) {
      return true;
    }

    const files = lockObj.files as Record<string, unknown>;
    if (!files) {
      return false;
    }

    const agents = files.agents as unknown[];
    if (!agents || agents.length === 0) {
      return true; // Empty arrays are v2
    }

    const firstAgent = agents[0];
    return typeof firstAgent === 'object' && firstAgent !== null;
  }

  /**
   * Migrates v1 legacy lock file to v2 format.
   *
   * Converts string arrays to FileEntry arrays with undefined checksums.
   * Marks as lockFileVersion: 1 to indicate migrated legacy format.
   */
  private migrateV1ToV2(legacyLock: LegacyLockFile): LockFile {
    const migrateArray = (paths: string[]): FileEntry[] => {
      return paths.map((path) => ({ path, checksum: undefined }));
    };

    return {
      version: legacyLock.version,
      installedAt: legacyLock.installedAt,
      lockFileVersion: 1,
      files: {
        agents: migrateArray(legacyLock.files.agents),
        commands: migrateArray(legacyLock.files.commands),
        templates: migrateArray(legacyLock.files.templates),
        rules: migrateArray(legacyLock.files.rules),
        outputStyles: migrateArray(legacyLock.files.outputStyles),
        binary: {
          path: legacyLock.files.binary,
          checksum: undefined,
        },
      },
    };
  }

  /**
   * Categorizes file entries into appropriate categories based on path.
   *
   * Path patterns:
   * - .claude/agents/ -> agents
   * - .claude/commands/ -> commands
   * - .the-startup/templates/ -> templates
   * - .the-startup/rules/ -> rules
   * - .the-startup/outputStyles/ -> outputStyles
   * - .the-startup/bin/ -> binary
   */
  private categorizeFiles(fileEntries: FileEntry[]): LockFile['files'] {
    const categorized: LockFile['files'] = {
      agents: [],
      commands: [],
      templates: [],
      rules: [],
      outputStyles: [],
      binary: { path: '', checksum: undefined },
    };

    for (const entry of fileEntries) {
      const normalizedPath = entry.path.replace(/\\/g, '/');

      if (normalizedPath.includes('/agents/')) {
        categorized.agents.push(entry);
      } else if (normalizedPath.includes('/commands/')) {
        categorized.commands.push(entry);
      } else if (normalizedPath.includes('/templates/')) {
        categorized.templates.push(entry);
      } else if (normalizedPath.includes('/rules/')) {
        categorized.rules.push(entry);
      } else if (normalizedPath.includes('/outputStyles/')) {
        categorized.outputStyles.push(entry);
      } else if (normalizedPath.includes('/bin/')) {
        categorized.binary = entry;
      }
    }

    return categorized;
  }

  /**
   * Builds a map of file paths to FileEntry for quick lookup.
   */
  private buildFileMap(lockFile: LockFile): Map<string, FileEntry> {
    const fileMap = new Map<string, FileEntry>();

    const addEntries = (entries: FileEntry[]) => {
      for (const entry of entries) {
        fileMap.set(entry.path, entry);
      }
    };

    addEntries(lockFile.files.agents);
    addEntries(lockFile.files.commands);
    addEntries(lockFile.files.templates);
    addEntries(lockFile.files.rules);
    addEntries(lockFile.files.outputStyles);

    if (lockFile.files.binary.path) {
      fileMap.set(lockFile.files.binary.path, lockFile.files.binary);
    }

    return fileMap;
  }
}
