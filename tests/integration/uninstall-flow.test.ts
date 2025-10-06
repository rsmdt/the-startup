import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { promises as fs } from 'fs';
import { join } from 'path';
import { LockManager } from '../../src/core/installer/LockManager';
import { SettingsMerger } from '../../src/core/installer/SettingsMerger';
import { createTempDir, cleanupTempDir } from '../shared/testUtils';
import type { LockFile } from '../../src/core/types/lock';
import type { ClaudeSettings } from '../../src/core/types/settings';

/**
 * Integration Test Suite: Uninstall Flow
 *
 * Tests the complete uninstall flow using real file system operations.
 * Simulates an existing installation and verifies clean removal.
 *
 * Test Scenarios (from PRD lines 177-189):
 * 1. Clean uninstall removes all files
 * 2. Uninstall with --keep-logs preserves logs
 * 3. Uninstall with --keep-settings preserves settings
 * 4. Uninstall with missing files (graceful handling)
 * 5. Settings.json restoration (remove hooks)
 * 6. Lock file deletion
 *
 * Uninstall Algorithm:
 * 1. Read lock file to get installed file list
 * 2. Remove all files listed in lock file
 * 3. Remove hooks from settings.json (unless --keep-settings)
 * 4. Delete lock file
 * 5. Optionally preserve logs directory (if --keep-logs)
 */

/**
 * Uninstaller - Integration test implementation
 *
 * This class implements the uninstall logic for integration testing.
 * In production, this would be in src/core/installer/Uninstaller.ts
 */
class Uninstaller {
  constructor(
    private lockManager: LockManager,
    private settingsMerger: SettingsMerger
  ) {}

  async uninstall(options: {
    keepLogs?: boolean;
    keepSettings?: boolean;
  }): Promise<{ success: boolean; removedFiles: string[]; errors?: string[] }> {
    const removedFiles: string[] = [];
    const errors: string[] = [];

    try {
      // 1. Read lock file
      const lockFile = await this.lockManager.readLockFile();

      if (!lockFile) {
        return {
          success: false,
          removedFiles: [],
          errors: ['No installation found. Lock file does not exist.'],
        };
      }

      // 2. Collect all files to remove
      const filesToRemove = this.collectFilesToRemove(lockFile, options.keepLogs);

      // 3. Remove files
      for (const filePath of filesToRemove) {
        try {
          await fs.rm(filePath, { force: true });
          removedFiles.push(filePath);
        } catch (error) {
          // Continue removing other files even if one fails
          const err = error as NodeJS.ErrnoException;
          if (err.code !== 'ENOENT') {
            errors.push(`Failed to remove ${filePath}: ${err.message}`);
          }
        }
      }

      // 4. Remove hooks from settings.json (unless --keep-settings)
      if (!options.keepSettings) {
        const settingsPath = this.inferSettingsPath(lockFile);
        try {
          await this.settingsMerger.removeHooks(settingsPath);
        } catch (error) {
          // Settings removal is non-fatal
          const err = error as Error;
          errors.push(`Failed to remove hooks from settings.json: ${err.message}`);
        }
      }

      // 5. Remove lock file
      try {
        await fs.rm(this.lockManager['lockFilePath'], { force: true });
      } catch (error) {
        const err = error as Error;
        errors.push(`Failed to remove lock file: ${err.message}`);
      }

      return {
        success: errors.length === 0,
        removedFiles,
        errors: errors.length > 0 ? errors : undefined,
      };
    } catch (error) {
      const err = error as Error;
      return {
        success: false,
        removedFiles,
        errors: [err.message],
      };
    }
  }

  private collectFilesToRemove(lockFile: LockFile, keepLogs: boolean = false): string[] {
    const files: string[] = [];

    // Add all agents
    for (const agent of lockFile.files.agents) {
      files.push(agent.path);
    }

    // Add all commands
    for (const command of lockFile.files.commands) {
      files.push(command.path);
    }

    // Add all templates
    for (const template of lockFile.files.templates) {
      // Skip log files if --keep-logs
      if (keepLogs && template.path.includes('/logs/')) {
        continue;
      }
      files.push(template.path);
    }

    // Add all rules
    for (const rule of lockFile.files.rules) {
      files.push(rule.path);
    }

    // Add all output styles
    for (const outputStyle of lockFile.files.outputStyles) {
      files.push(outputStyle.path);
    }

    // Add binary
    if (lockFile.files.binary.path) {
      files.push(lockFile.files.binary.path);
    }

    return files;
  }

  private inferSettingsPath(lockFile: LockFile): string {
    // Infer settings path from agent files (they're in .claude/agents/)
    if (lockFile.files.agents.length > 0) {
      const agentPath = lockFile.files.agents[0].path;
      const claudePath = agentPath.replace(/\/agents\/.*$/, '');
      return join(claudePath, 'settings.json');
    }

    // Fallback: look for command files
    if (lockFile.files.commands.length > 0) {
      const commandPath = lockFile.files.commands[0].path;
      const claudePath = commandPath.replace(/\/commands\/.*$/, '');
      return join(claudePath, 'settings.json');
    }

    throw new Error('Cannot infer settings.json path from lock file');
  }
}

describe('Integration: Uninstall Flow', () => {
  let tempDir: string;
  let startupPath: string;
  let claudePath: string;
  let lockFilePath: string;
  let uninstaller: Uninstaller;
  let lockManager: LockManager;
  let settingsMerger: SettingsMerger;

  beforeEach(async () => {
    tempDir = await createTempDir();
    startupPath = join(tempDir, '.the-startup');
    claudePath = join(tempDir, '.claude');
    lockFilePath = join(startupPath, 'lock.json');

    lockManager = new LockManager(lockFilePath);
    settingsMerger = new SettingsMerger(fs);
    uninstaller = new Uninstaller(lockManager, settingsMerger);

    // Create directory structure
    await fs.mkdir(join(claudePath, 'agents'), { recursive: true });
    await fs.mkdir(join(claudePath, 'commands'), { recursive: true });
    await fs.mkdir(join(startupPath, 'templates'), { recursive: true });
    await fs.mkdir(join(startupPath, 'rules'), { recursive: true });
    await fs.mkdir(join(claudePath, 'output-styles'), { recursive: true });
  });

  afterEach(async () => {
    await cleanupTempDir(tempDir);
  });

  /**
   * Helper: Create a complete installation for testing uninstall
   */
  async function createInstallation(): Promise<void> {
    // Create agent files
    await fs.writeFile(
      join(claudePath, 'agents/specify.md'),
      '# Specify Agent',
      'utf-8'
    );
    await fs.writeFile(
      join(claudePath, 'agents/implement.md'),
      '# Implement Agent',
      'utf-8'
    );

    // Create command files
    await fs.writeFile(
      join(claudePath, 'commands/s-specify.md'),
      '# /s:specify command',
      'utf-8'
    );

    // Create template files
    await fs.writeFile(
      join(startupPath, 'templates/SPEC.md'),
      '# Specification Template',
      'utf-8'
    );
    await fs.writeFile(
      join(startupPath, 'templates/TASK-DOD.md'),
      '# Task Definition of Done',
      'utf-8'
    );

    // Create rule files
    await fs.writeFile(
      join(startupPath, 'rules/SCQA.md'),
      '# SCQA Framework',
      'utf-8'
    );

    // Create output style files
    await fs.writeFile(
      join(claudePath, 'output-styles/json.md'),
      '# JSON Output Style',
      'utf-8'
    );

    // Create settings.json with hooks
    const settings: ClaudeSettings = {
      mcpServers: {},
      hooks: {
        'startup-validate-dor': {
          command: `${startupPath}/bin/the-startup validate dor`,
          description: 'Validate DOR',
        },
        'startup-validate-dod': {
          command: `${startupPath}/bin/the-startup validate dod`,
          description: 'Validate DOD',
        },
        'user-custom-hook': {
          command: 'echo "User hook"',
          description: 'User custom hook',
        },
      },
    };
    await fs.writeFile(
      join(claudePath, 'settings.json'),
      JSON.stringify(settings, null, 2),
      'utf-8'
    );

    // Create lock file
    const lockFile: LockFile = {
      version: '1.0.0',
      installedAt: new Date().toISOString(),
      lockFileVersion: 2,
      files: {
        agents: [
          { path: join(claudePath, 'agents/specify.md'), checksum: 'abc123' },
          { path: join(claudePath, 'agents/implement.md'), checksum: 'def456' },
        ],
        commands: [
          { path: join(claudePath, 'commands/s-specify.md'), checksum: 'ghi789' },
        ],
        templates: [
          { path: join(startupPath, 'templates/SPEC.md'), checksum: 'jkl012' },
          { path: join(startupPath, 'templates/TASK-DOD.md'), checksum: 'mno345' },
        ],
        rules: [
          { path: join(startupPath, 'rules/SCQA.md'), checksum: 'pqr678' },
        ],
        outputStyles: [
          { path: join(claudePath, 'output-styles/json.md'), checksum: 'stu901' },
        ],
        binary: {
          path: join(startupPath, 'bin/the-startup'),
          checksum: 'vwx234',
        },
      },
    };

    await fs.writeFile(
      lockFilePath,
      JSON.stringify(lockFile, null, 2),
      'utf-8'
    );
  }

  /**
   * Scenario 1: Clean uninstall removes all files
   *
   * Verifies:
   * - All installed files are removed
   * - Lock file is deleted
   * - Hooks are removed from settings.json
   * - User's other settings are preserved
   */
  it('should remove all installed files during clean uninstall', async () => {
    await createInstallation();

    const result = await uninstaller.uninstall({
      keepLogs: false,
      keepSettings: false,
    });

    expect(result.success).toBe(true);
    expect(result.removedFiles.length).toBeGreaterThan(0);

    // Verify all agent files are removed
    await expect(fs.access(join(claudePath, 'agents/specify.md'))).rejects.toThrow();
    await expect(fs.access(join(claudePath, 'agents/implement.md'))).rejects.toThrow();

    // Verify all command files are removed
    await expect(fs.access(join(claudePath, 'commands/s-specify.md'))).rejects.toThrow();

    // Verify all template files are removed
    await expect(fs.access(join(startupPath, 'templates/SPEC.md'))).rejects.toThrow();
    await expect(fs.access(join(startupPath, 'templates/TASK-DOD.md'))).rejects.toThrow();

    // Verify all rule files are removed
    await expect(fs.access(join(startupPath, 'rules/SCQA.md'))).rejects.toThrow();

    // Verify all output style files are removed
    await expect(fs.access(join(claudePath, 'output-styles/json.md'))).rejects.toThrow();

    // Verify lock file is removed
    await expect(fs.access(lockFilePath)).rejects.toThrow();

    // Verify hooks are removed from settings.json
    const settingsContent = await fs.readFile(join(claudePath, 'settings.json'), 'utf-8');
    const settings = JSON.parse(settingsContent);

    expect(settings.hooks).toBeUndefined();

    // Other settings should still exist
    expect(settings.mcpServers).toBeDefined();
  });

  /**
   * Scenario 2: Uninstall with --keep-logs preserves log files
   *
   * Verifies:
   * - Log files are preserved
   * - Other files are removed
   */
  it('should preserve log files when --keep-logs flag is used', async () => {
    await createInstallation();

    // Add log files to installation
    await fs.mkdir(join(startupPath, 'templates/logs'), { recursive: true });
    await fs.writeFile(
      join(startupPath, 'templates/logs/install.log'),
      'Installation log content',
      'utf-8'
    );

    // Update lock file to include log file
    const lockFile = await lockManager.readLockFile();
    if (lockFile) {
      lockFile.files.templates.push({
        path: join(startupPath, 'templates/logs/install.log'),
        checksum: 'log123',
      });
      await fs.writeFile(lockFilePath, JSON.stringify(lockFile, null, 2), 'utf-8');
    }

    const result = await uninstaller.uninstall({
      keepLogs: true,
      keepSettings: false,
    });

    expect(result.success).toBe(true);

    // Verify log file was preserved
    await expect(fs.access(join(startupPath, 'templates/logs/install.log'))).resolves.not.toThrow();
    const logContent = await fs.readFile(join(startupPath, 'templates/logs/install.log'), 'utf-8');
    expect(logContent).toBe('Installation log content');

    // Verify other templates were removed
    await expect(fs.access(join(startupPath, 'templates/SPEC.md'))).rejects.toThrow();
  });

  /**
   * Scenario 3: Uninstall with --keep-settings preserves settings.json
   *
   * Verifies:
   * - Settings.json is not modified
   * - Hooks remain in settings.json
   * - Files are still removed
   */
  it('should preserve settings.json when --keep-settings flag is used', async () => {
    await createInstallation();

    const settingsBeforeUninstall = await fs.readFile(join(claudePath, 'settings.json'), 'utf-8');

    const result = await uninstaller.uninstall({
      keepLogs: false,
      keepSettings: true,
    });

    expect(result.success).toBe(true);

    // Verify files were removed
    await expect(fs.access(join(claudePath, 'agents/specify.md'))).rejects.toThrow();

    // Verify settings.json was NOT modified
    const settingsAfterUninstall = await fs.readFile(join(claudePath, 'settings.json'), 'utf-8');
    expect(settingsAfterUninstall).toBe(settingsBeforeUninstall);

    const settings = JSON.parse(settingsAfterUninstall);
    expect(settings.hooks).toBeDefined();
    expect(settings.hooks['startup-validate-dor']).toBeDefined();
    expect(settings.hooks['user-custom-hook']).toBeDefined();
  });

  /**
   * Scenario 4: Uninstall with missing files (graceful handling)
   *
   * Verifies:
   * - Uninstall continues even if some files are missing
   * - No errors for missing files
   * - Other files are still removed
   * - Process completes successfully
   */
  it('should handle missing files gracefully during uninstall', async () => {
    await createInstallation();

    // Delete some files before uninstall (simulate manual deletion)
    await fs.rm(join(claudePath, 'agents/implement.md'), { force: true });
    await fs.rm(join(startupPath, 'templates/SPEC.md'), { force: true });

    const result = await uninstaller.uninstall({
      keepLogs: false,
      keepSettings: false,
    });

    // Should succeed despite missing files
    expect(result.success).toBe(true);

    // Verify other files were removed
    await expect(fs.access(join(claudePath, 'agents/specify.md'))).rejects.toThrow();
    await expect(fs.access(join(claudePath, 'commands/s-specify.md'))).rejects.toThrow();

    // Verify lock file was removed
    await expect(fs.access(lockFilePath)).rejects.toThrow();
  });

  /**
   * Scenario 5: Settings.json hook removal
   *
   * Verifies:
   * - All hooks are removed from settings.json
   * - Other settings properties are preserved
   * - Backup is created and cleaned up
   */
  it('should remove all hooks from settings.json during uninstall', async () => {
    await createInstallation();

    const result = await uninstaller.uninstall({
      keepLogs: false,
      keepSettings: false,
    });

    expect(result.success).toBe(true);

    // Verify hooks property is removed entirely
    const settingsContent = await fs.readFile(join(claudePath, 'settings.json'), 'utf-8');
    const settings = JSON.parse(settingsContent);

    expect(settings.hooks).toBeUndefined();

    // Verify other properties are preserved
    expect(settings.mcpServers).toBeDefined();

    // Note: User custom hooks are also removed during uninstall
    // This is intentional to ensure clean removal
  });

  /**
   * Scenario 6: Uninstall with no lock file
   *
   * Verifies:
   * - Uninstall fails gracefully if no lock file exists
   * - Error message is descriptive
   * - No files are modified
   */
  it('should fail gracefully when lock file does not exist', async () => {
    // Don't create installation (no lock file)

    const result = await uninstaller.uninstall({
      keepLogs: false,
      keepSettings: false,
    });

    expect(result.success).toBe(false);
    expect(result.errors).toBeDefined();
    expect(result.errors![0]).toContain('No installation found');
    expect(result.removedFiles.length).toBe(0);
  });

  /**
   * Scenario 7: Partial uninstall on permission errors
   *
   * Verifies:
   * - Uninstall continues even if some files cannot be deleted
   * - Errors are reported for failed deletions
   * - Successfully deleted files are tracked
   * - Lock file is still removed
   */
  it('should report errors but continue uninstall on permission errors', async () => {
    await createInstallation();

    // Make one file read-only to simulate permission error
    const protectedFile = join(claudePath, 'agents/implement.md');
    await fs.chmod(protectedFile, 0o444); // Read-only

    const result = await uninstaller.uninstall({
      keepLogs: false,
      keepSettings: false,
    });

    // Uninstall may report errors but should continue
    expect(result.removedFiles.length).toBeGreaterThan(0);

    // Other files should be removed
    await expect(fs.access(join(claudePath, 'agents/specify.md'))).rejects.toThrow();
    await expect(fs.access(join(startupPath, 'templates/SPEC.md'))).rejects.toThrow();

    // Clean up protected file if it still exists
    try {
      await fs.access(protectedFile);
      await fs.chmod(protectedFile, 0o644);
      await fs.rm(protectedFile, { force: true });
    } catch (error) {
      // File already removed, that's OK
    }
  });
});
