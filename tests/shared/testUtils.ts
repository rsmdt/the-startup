import { vi } from 'vitest';
import { promises as fs } from 'fs';
import * as path from 'path';
import * as os from 'os';

/**
 * Test Utilities
 *
 * Provides helper functions for test setup, teardown, and mock data generation.
 * These utilities are used across all test suites to ensure consistency.
 */

/**
 * Creates a temporary directory for test isolation
 *
 * @returns Promise resolving to absolute path of created temp directory
 * @throws Error if directory creation fails
 *
 * @example
 * const tempDir = await createTempDir();
 * // Use tempDir for test operations
 * await cleanupTempDir(tempDir);
 */
export async function createTempDir(): Promise<string> {
  const tempPrefix = path.join(os.tmpdir(), 'the-startup-test-');
  const tempDir = await fs.mkdtemp(tempPrefix);
  return tempDir;
}

/**
 * Removes a temporary directory and all its contents
 *
 * @param dirPath - Absolute path to directory to remove
 * @throws Error if directory removal fails
 *
 * @example
 * const tempDir = await createTempDir();
 * // ... test operations
 * await cleanupTempDir(tempDir);
 */
export async function cleanupTempDir(dirPath: string): Promise<void> {
  await fs.rm(dirPath, { recursive: true, force: true });
}

/**
 * Creates an in-memory file system mock using Vitest
 *
 * Returns a mock object that can be used to track fs operations without
 * touching the real file system. All operations are tracked and can be
 * verified in tests.
 *
 * @returns Mock file system object with common fs operations
 *
 * @example
 * const mockFs = mockFileSystem();
 * mockFs.readFile.mockResolvedValue('file content');
 * // Use mockFs in your test
 * expect(mockFs.readFile).toHaveBeenCalledWith('/path/to/file', 'utf-8');
 */
export function mockFileSystem() {
  return {
    readFile: vi.fn(),
    writeFile: vi.fn(),
    mkdir: vi.fn(),
    rm: vi.fn(),
    access: vi.fn(),
    stat: vi.fn(),
    readdir: vi.fn(),
    copyFile: vi.fn(),
  };
}

/**
 * Creates a mock lock file for testing
 *
 * @param version - Package version (e.g., '1.0.0')
 * @param lockFileVersion - Lock file format version (1 for legacy string[], 2 for FileEntry[])
 * @returns Mock lock file object matching LockFile or LegacyLockFile interface
 *
 * @example
 * // Create v2 lock file with checksums
 * const lockV2 = createMockLockFile('1.0.0', 2);
 *
 * // Create v1 legacy lock file
 * const lockV1 = createMockLockFile('0.9.0', 1);
 */
export function createMockLockFile(
  version: string,
  lockFileVersion: 1 | 2 = 2
): LockFileV1 | LockFileV2 {
  const installedAt = new Date('2025-10-06T12:00:00.000Z').toISOString();

  if (lockFileVersion === 1) {
    // Legacy format with string arrays
    return {
      version,
      installedAt,
      files: {
        agents: [
          '.claude/agents/specify.md',
          '.claude/agents/implement.md',
        ],
        commands: [
          '.claude/commands/s-specify.md',
          '.claude/commands/s-implement.md',
        ],
        templates: [
          '.the-startup/templates/SPEC.md',
          '.the-startup/templates/TASK-DOD.md',
        ],
        rules: [],
        outputStyles: [],
        binary: '.the-startup/bin/the-startup',
      },
    };
  }

  // V2 format with FileEntry[] and checksums
  return {
    version,
    installedAt,
    lockFileVersion: 2,
    files: {
      agents: [
        {
          path: '.claude/agents/specify.md',
          checksum: 'abc123def456',
        },
        {
          path: '.claude/agents/implement.md',
          checksum: 'def789ghi012',
        },
      ],
      commands: [
        {
          path: '.claude/commands/s-specify.md',
          checksum: 'ghi345jkl678',
        },
        {
          path: '.claude/commands/s-implement.md',
          checksum: 'jkl901mno234',
        },
      ],
      templates: [
        {
          path: '.the-startup/templates/SPEC.md',
          checksum: 'mno567pqr890',
        },
        {
          path: '.the-startup/templates/TASK-DOD.md',
          checksum: 'pqr123stu456',
        },
      ],
      rules: [],
      outputStyles: [],
      binary: {
        path: '.the-startup/bin/the-startup',
        checksum: 'stu789vwx012',
      },
    },
  };
}

/**
 * Creates a mock Claude settings.json file for testing
 *
 * @param includeHooks - Whether to include existing user hooks in the settings
 * @returns Mock settings object matching ClaudeSettings interface
 *
 * @example
 * // Empty settings
 * const emptySettings = createMockSettings(false);
 *
 * // Settings with existing hooks
 * const settingsWithHooks = createMockSettings(true);
 */
export function createMockSettings(includeHooks: boolean = false): ClaudeSettings {
  const baseSettings: ClaudeSettings = {
    mcpServers: {},
    modelSettings: {
      defaultModel: 'claude-sonnet-4-5',
    },
  };

  if (includeHooks) {
    baseSettings.hooks = {
      'user-custom-hook': {
        command: 'echo "User custom hook"',
        description: 'User custom validation hook',
        continueOnError: false,
      },
      'another-user-hook': {
        command: 'npm run lint',
        description: 'Lint before commit',
        continueOnError: true,
      },
    };
  }

  return baseSettings;
}

/**
 * Creates a mock asset directory structure for testing
 *
 * @returns Map of file paths to their mock content
 *
 * @example
 * const assets = createMockAssetStructure();
 * const specifyAgentContent = assets.get('agents/specify.md');
 */
export function createMockAssetStructure(): Map<string, string> {
  const assets = new Map<string, string>();

  // Agent files
  assets.set('agents/specify.md', '# Specify Agent\nCreates specifications from requirements.');
  assets.set('agents/implement.md', '# Implement Agent\nImplements code from specifications.');
  assets.set('agents/refactor.md', '# Refactor Agent\nRefactors code for maintainability.');

  // Command files
  assets.set('commands/s-specify.md', '# /s:specify command\nCreate specification');
  assets.set('commands/s-implement.md', '# /s:implement command\nImplement from spec');
  assets.set('commands/s-refactor.md', '# /s:refactor command\nRefactor code');

  // Template files
  assets.set('templates/SPEC.md', '# Specification Template\n{{PLACEHOLDER}}');
  assets.set('templates/TASK-DOD.md', '# Task Definition of Done\n{{PLACEHOLDER}}');
  assets.set('templates/DOR.md', '# Definition of Ready\n{{PLACEHOLDER}}');
  assets.set('templates/DOD.md', '# Definition of Done\n{{PLACEHOLDER}}');

  // Settings template
  assets.set('settings.json', JSON.stringify({
    hooks: {
      'startup-validate-dor': {
        command: '{{STARTUP_PATH}}/bin/the-startup validate dor',
        description: 'Validate Definition of Ready',
        continueOnError: false,
      },
    },
  }, null, 2));

  return assets;
}

/**
 * Waits for a specified number of milliseconds
 * Useful for testing async operations and timeouts
 *
 * @param ms - Milliseconds to wait
 *
 * @example
 * await sleep(100); // Wait 100ms
 */
export async function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

/**
 * Creates a mock error object with specified properties
 *
 * @param message - Error message
 * @param code - Optional error code (e.g., 'ENOENT', 'EACCES')
 * @returns Error object with additional properties
 *
 * @example
 * const error = createMockError('File not found', 'ENOENT');
 * mockFs.readFile.mockRejectedValue(error);
 */
export function createMockError(message: string, code?: string): NodeJS.ErrnoException {
  const error = new Error(message) as NodeJS.ErrnoException;
  if (code) {
    error.code = code;
  }
  return error;
}

// Type definitions to match the SDD

interface FileEntry {
  path: string;
  checksum?: string;
}

interface LockFileV2 {
  version: string;
  installedAt: string;
  lockFileVersion: 2;
  files: {
    agents: FileEntry[];
    commands: FileEntry[];
    templates: FileEntry[];
    rules: FileEntry[];
    outputStyles: FileEntry[];
    binary: FileEntry;
  };
}

interface LockFileV1 {
  version: string;
  installedAt: string;
  files: {
    agents: string[];
    commands: string[];
    templates: string[];
    rules: string[];
    outputStyles: string[];
    binary: string;
  };
}

interface ClaudeSettings {
  hooks?: {
    [hookName: string]: {
      command: string;
      description?: string;
      continueOnError?: boolean;
    };
  };
  mcpServers?: Record<string, unknown>;
  modelSettings?: Record<string, unknown>;
  [key: string]: unknown;
}
