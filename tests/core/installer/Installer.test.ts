import { describe, it, expect, beforeEach, vi } from 'vitest';
import { Installer } from '../../../src/core/installer/Installer';
import { LockManager } from '../../../src/core/installer/LockManager';
import { SettingsMerger } from '../../../src/core/installer/SettingsMerger';
import type { InstallerOptions, InstallResult } from '../../../src/core/types/config';
import type { FileEntry } from '../../../src/core/types/lock';

/**
 * Test suite for Installer core
 *
 * Tests follow SDD requirements:
 * - Installer interface (lines 569-586)
 * - Install flow (lines 677-718)
 * - Non-interactive flow (lines 720-756)
 * - Error handling (lines 899-921)
 *
 * Coverage requirement: 90%+ (SDD line 1241)
 */

describe('Installer', () => {
  let installer: Installer;
  let mockFs: any;
  let mockLockManager: any;
  let mockSettingsMerger: any;
  let mockAssetProvider: any;

  const defaultOptions: InstallerOptions = {
    startupPath: '/test/.the-startup',
    claudePath: '/test/.claude',
    selectedFiles: {
      agents: true,
      commands: true,
      templates: true,
      rules: true,
      outputStyles: true,
    },
  };

  beforeEach(() => {
    // Mock file system
    mockFs = {
      mkdir: vi.fn().mockResolvedValue(undefined),
      copyFile: vi.fn().mockResolvedValue(undefined),
      readFile: vi.fn().mockImplementation((path: string) => {
        // For JSON files, return valid JSON
        if (path.endsWith('.json')) {
          return Promise.resolve(JSON.stringify({ test: 'data' }));
        }
        // For other files, return mock content
        return Promise.resolve('mock content');
      }),
      writeFile: vi.fn().mockResolvedValue(undefined),
      rm: vi.fn().mockResolvedValue(undefined),
      access: vi.fn().mockResolvedValue(undefined),
      stat: vi.fn().mockResolvedValue({ isDirectory: () => true }),
    };

    // Mock LockManager
    mockLockManager = {
      writeLockFile: vi.fn().mockResolvedValue(undefined),
      generateChecksum: vi.fn().mockResolvedValue('abc123'),
    };

    // Mock SettingsMerger
    mockSettingsMerger = {
      mergeSettings: vi.fn().mockResolvedValue({
        hooks: {
          'user-prompt-submit': {
            command: '/test/.the-startup/bin/the-startup statusline',
          },
        },
      }),
      mergeFullSettings: vi.fn().mockResolvedValue({
        permissions: {
          additionalDirectories: ['/test/.the-startup'],
        },
        statusLine: {
          type: 'command',
          command: '/test/.the-startup/bin/the-startup statusline',
        },
        hooks: {
          'user-prompt-submit': {
            command: '/test/.the-startup/bin/the-startup statusline',
          },
        },
      }),
    };

    // Mock AssetProvider (simplified structure)
    mockAssetProvider = {
      getAssetFiles: vi.fn().mockReturnValue([
        { sourcePath: 'claude/agents/the-chief.md', relativePath: 'agents/the-chief.md', targetCategory: 'claude', isJson: false },
        { sourcePath: 'claude/commands/spec.md', relativePath: 'commands/spec.md', targetCategory: 'claude', isJson: false },
        { sourcePath: 'the-startup/templates/PRD.md', relativePath: 'templates/PRD.md', targetCategory: 'startup', isJson: false },
        { sourcePath: 'the-startup/rules/SCQA.md', relativePath: 'rules/SCQA.md', targetCategory: 'startup', isJson: false },
        { sourcePath: 'claude/output-styles/the-startup.md', relativePath: 'output-styles/the-startup.md', targetCategory: 'claude', isJson: false },
        { sourcePath: 'claude/settings.json', relativePath: 'settings.json', targetCategory: 'claude', isJson: true },
        { sourcePath: 'claude/settings.local.json', relativePath: 'settings.local.json', targetCategory: 'claude', isJson: true },
      ]),
    };

    installer = new Installer(
      mockFs,
      mockLockManager,
      mockSettingsMerger,
      mockAssetProvider,
      '1.0.0'
    );
  });

  describe('install', () => {
    it('should successfully install all selected assets', async () => {
      const result = await installer.install(defaultOptions);

      expect(result.success).toBe(true);
      expect(result.installedFiles.length).toBeGreaterThan(0);
      expect(result.errors).toBeUndefined();
    });

    it('should create installation directories', async () => {
      await installer.install(defaultOptions);

      expect(mockFs.mkdir).toHaveBeenCalledWith('/test/.the-startup', {
        recursive: true,
      });
      expect(mockFs.mkdir).toHaveBeenCalledWith('/test/.claude/agents', {
        recursive: true,
      });
    });

    it('should copy selected asset files to correct destinations', async () => {
      await installer.install(defaultOptions);

      // Check that writeFile was called for each asset (files are now written with placeholders replaced)
      expect(mockFs.writeFile).toHaveBeenCalled();
      const writeFileCalls = mockFs.writeFile.mock.calls;

      // Verify agents were written to .claude/agents/
      const agentCalls = writeFileCalls.filter((call: any) =>
        call[0].includes('.claude/agents/')
      );
      expect(agentCalls.length).toBeGreaterThan(0);

      // Verify templates were written to .the-startup/templates/
      const templateCalls = writeFileCalls.filter((call: any) =>
        call[0].includes('.the-startup/templates/')
      );
      expect(templateCalls.length).toBeGreaterThan(0);
    });

    it('should only copy selected file categories', async () => {
      const options: InstallerOptions = {
        ...defaultOptions,
        selectedFiles: {
          agents: true,
          commands: false,
          templates: false,
          rules: false,
          outputStyles: false,
        },
      };

      mockAssetProvider.getAssetFiles.mockReturnValue([
        { sourcePath: 'claude/agents/the-chief.md', relativePath: 'agents/the-chief.md', targetCategory: 'claude', isJson: false },
      ]);

      await installer.install(options);

      // Should write the selected file
      expect(mockFs.writeFile).toHaveBeenCalled();
    });

    it('should merge settings.json during installation', async () => {
      // Mock readFile to return settings.json content
      mockFs.readFile.mockResolvedValue(JSON.stringify({
        permissions: {
          additionalDirectories: ['{{STARTUP_PATH}}'],
        },
        statusLine: {
          type: 'command',
          command: '{{STARTUP_PATH}}/bin/statusline{{SHELL_SCRIPT_EXTENSION}}',
        },
        hooks: {
          'user-prompt-submit': {
            command: '{{STARTUP_PATH}}/bin/statusline{{SHELL_SCRIPT_EXTENSION}}',
          },
        },
      }));

      await installer.install(defaultOptions);

      // Settings merge should be called for both settings.json and settings.local.json
      expect(mockSettingsMerger.mergeFullSettings).toHaveBeenCalled();
    });

    it('should create lock file with installed file entries', async () => {
      await installer.install(defaultOptions);

      expect(mockLockManager.writeLockFile).toHaveBeenCalledWith(
        expect.arrayContaining([
          expect.objectContaining({
            path: expect.any(String),
            checksum: expect.any(String),
          }),
        ]),
        '1.0.0'
      );
    });

    it('should generate checksums for all installed files', async () => {
      await installer.install(defaultOptions);

      expect(mockLockManager.generateChecksum).toHaveBeenCalled();
    });

    it('should return list of installed file paths', async () => {
      const result = await installer.install(defaultOptions);

      expect(result.installedFiles).toEqual(
        expect.arrayContaining([expect.any(String)])
      );
      expect(result.installedFiles.length).toBeGreaterThan(0);
    });
  });

  describe('error handling', () => {
    it('should handle invalid path error with clear message', async () => {
      mockFs.mkdir.mockRejectedValue({
        code: 'ENOENT',
        message: 'no such file or directory',
      });

      const result = await installer.install(defaultOptions);

      expect(result.success).toBe(false);
      expect(result.errors).toBeDefined();
      expect(result.errors![0]).toContain('Please re-enter a valid path');
    });

    it('should handle permission denied error with helpful suggestion', async () => {
      mockFs.mkdir.mockRejectedValue({
        code: 'EACCES',
        message: 'permission denied',
      });

      const result = await installer.install(defaultOptions);

      expect(result.success).toBe(false);
      expect(result.errors).toBeDefined();
      expect(result.errors![0]).toContain('permission');
      expect(result.errors![0]).toMatch(/chmod|permission/i);
    });

    it('should handle disk full error with space information', async () => {
      mockFs.writeFile.mockRejectedValue({
        code: 'ENOSPC',
        message: 'no space left on device',
      });

      const result = await installer.install(defaultOptions);

      expect(result.success).toBe(false);
      expect(result.errors).toBeDefined();
      expect(result.errors![0]).toMatch(/space|disk/i);
    });

    it('should handle settings merge failure with rollback', async () => {
      mockSettingsMerger.mergeFullSettings.mockRejectedValue(
        new Error('Invalid JSON in settings file')
      );

      const result = await installer.install(defaultOptions);

      expect(result.success).toBe(false);
      expect(result.errors).toBeDefined();
    });

    it('should handle asset copy failure with cleanup', async () => {
      let writeCount = 0;
      mockFs.writeFile.mockImplementation(() => {
        writeCount++;
        if (writeCount === 3) {
          throw new Error('Write failed');
        }
        return Promise.resolve();
      });

      const result = await installer.install(defaultOptions);

      expect(result.success).toBe(false);
      // Verify cleanup was attempted
      expect(mockFs.rm).toHaveBeenCalled();
    });
  });

  describe('rollback mechanism', () => {
    it('should rollback all changes on failure', async () => {
      // Simulate failure during lock file creation
      mockLockManager.writeLockFile.mockRejectedValue(
        new Error('Lock write failed')
      );

      const result = await installer.install(defaultOptions);

      expect(result.success).toBe(false);
      // Verify cleanup was called to remove copied files
      expect(mockFs.rm).toHaveBeenCalled();
    });

    it('should track installed files for rollback', async () => {
      // Install 2 files successfully, then fail on 3rd
      let writeCount = 0;
      const writtenFiles: string[] = [];

      mockFs.writeFile.mockImplementation((dest: string) => {
        writeCount++;
        if (writeCount === 3) {
          throw new Error('Write failed');
        }
        writtenFiles.push(dest);
        return Promise.resolve();
      });

      await installer.install(defaultOptions);

      // Should clean up the 2 files that were written
      expect(mockFs.rm).toHaveBeenCalled();
    });

    it('should not rollback if all operations succeed', async () => {
      const result = await installer.install(defaultOptions);

      // Should succeed without errors
      expect(result.success).toBe(true);
      expect(result.installedFiles.length).toBeGreaterThan(0);

      // rm should not be called to remove installed files on success
      // (SettingsMerger may call rm for backup cleanup, which is OK)
      const rmCalls = mockFs.rm.mock.calls;
      const removedInstalledFiles = rmCalls.filter((call: any) =>
        result.installedFiles.includes(call[0])
      );
      expect(removedInstalledFiles.length).toBe(0);
    });

    it('should continue rollback even if individual file deletion fails', async () => {
      // Write some files successfully, then fail
      let writeCount = 0;
      mockFs.writeFile.mockImplementation(() => {
        writeCount++;
        if (writeCount === 3) {
          throw new Error('Write failed');
        }
        return Promise.resolve();
      });

      // Make rm fail for the first file but succeed for others
      let rmCount = 0;
      mockFs.rm.mockImplementation(() => {
        rmCount++;
        if (rmCount === 1) {
          throw new Error('Permission denied');
        }
        return Promise.resolve();
      });

      const result = await installer.install(defaultOptions);

      expect(result.success).toBe(false);
      // Should have attempted to remove all written files despite first failure
      expect(mockFs.rm).toHaveBeenCalled();
    });
  });

  describe('progress reporting', () => {
    it('should report progress for operations taking > 5 seconds', async () => {
      const progressCallback = vi.fn();
      const installerWithProgress = new Installer(
        mockFs,
        mockLockManager,
        mockSettingsMerger,
        mockAssetProvider,
        '1.0.0',
        progressCallback
      );

      await installerWithProgress.install(defaultOptions);

      // Progress should be reported at various stages
      expect(progressCallback).toHaveBeenCalled();
      expect(progressCallback).toHaveBeenCalledWith(
        expect.objectContaining({
          stage: expect.any(String),
          current: expect.any(Number),
          total: expect.any(Number),
        })
      );
    });

    it('should report completion percentage', async () => {
      const progressCallback = vi.fn();
      const installerWithProgress = new Installer(
        mockFs,
        mockLockManager,
        mockSettingsMerger,
        mockAssetProvider,
        '1.0.0',
        progressCallback
      );

      await installerWithProgress.install(defaultOptions);

      // Find the last progress call
      const lastCall =
        progressCallback.mock.calls[progressCallback.mock.calls.length - 1][0];
      expect(lastCall.current).toBe(lastCall.total);
    });
  });

  describe('path handling', () => {
    it('should normalize paths with tilde expansion', async () => {
      const options: InstallerOptions = {
        startupPath: '~/.the-startup',
        claudePath: '~/.claude',
        selectedFiles: defaultOptions.selectedFiles,
      };

      // Mock home directory
      const homeDir = '/Users/testuser';
      const installerWithHome = new Installer(
        mockFs,
        mockLockManager,
        mockSettingsMerger,
        mockAssetProvider,
        '1.0.0',
        undefined,
        homeDir
      );

      await installerWithHome.install(options);

      // Verify paths were expanded
      expect(mockFs.mkdir).toHaveBeenCalledWith(
        `${homeDir}/.the-startup`,
        expect.any(Object)
      );
    });

    it('should convert absolute paths back to tilde notation for placeholders', async () => {
      const options: InstallerOptions = {
        startupPath: '~/.the-startup',
        claudePath: '~/.claude',
        selectedFiles: defaultOptions.selectedFiles,
      };

      // Mock home directory
      const homeDir = '/Users/testuser';
      const installerWithHome = new Installer(
        mockFs,
        mockLockManager,
        mockSettingsMerger,
        mockAssetProvider,
        '1.0.0',
        undefined,
        homeDir
      );

      await installerWithHome.install(options);

      // Verify settings merger was called with tilde paths in placeholders
      expect(mockSettingsMerger.mergeFullSettings).toHaveBeenCalledWith(
        expect.any(String),
        expect.any(Object),
        expect.objectContaining({
          STARTUP_PATH: '~/.the-startup',
          CLAUDE_PATH: '~/.claude',
        })
      );
    });

    it('should keep absolute paths unchanged when not under home directory', async () => {
      const options: InstallerOptions = {
        startupPath: '/opt/the-startup',
        claudePath: '/etc/claude',
        selectedFiles: defaultOptions.selectedFiles,
      };

      const homeDir = '/Users/testuser';
      const installerWithHome = new Installer(
        mockFs,
        mockLockManager,
        mockSettingsMerger,
        mockAssetProvider,
        '1.0.0',
        undefined,
        homeDir
      );

      await installerWithHome.install(options);

      // Verify settings merger was called with absolute paths (not tilde)
      expect(mockSettingsMerger.mergeFullSettings).toHaveBeenCalledWith(
        expect.any(String),
        expect.any(Object),
        expect.objectContaining({
          STARTUP_PATH: '/opt/the-startup',
          CLAUDE_PATH: '/etc/claude',
        })
      );
    });

    it('should handle case-sensitive paths on case-sensitive systems', async () => {
      const result = await installer.install({
        ...defaultOptions,
        startupPath: '/Test/.the-startup',
      });

      expect(result.success).toBe(true);
      expect(mockFs.mkdir).toHaveBeenCalledWith(
        '/Test/.the-startup',
        expect.any(Object)
      );
    });

    it('should resolve relative paths to absolute paths', async () => {
      const options: InstallerOptions = {
        startupPath: './.the-startup',
        claudePath: '~/.claude',
        selectedFiles: defaultOptions.selectedFiles,
      };

      const cwd = '/Users/testuser/project';
      const installerWithCwd = new Installer(
        mockFs,
        mockLockManager,
        mockSettingsMerger,
        mockAssetProvider,
        '1.0.0',
        undefined,
        '/Users/testuser',
        cwd
      );

      await installerWithCwd.install(options);

      expect(mockFs.mkdir).toHaveBeenCalledWith(
        `${cwd}/.the-startup`,
        expect.any(Object)
      );
    });
  });

  describe('atomic operations', () => {
    it('should complete all operations or none (all-or-nothing)', async () => {
      // Simulate failure in the middle
      mockLockManager.writeLockFile.mockRejectedValue(
        new Error('Lock failed')
      );

      const result = await installer.install(defaultOptions);

      expect(result.success).toBe(false);
      // All copied files should be cleaned up
      expect(mockFs.rm).toHaveBeenCalled();
    });

    it('should not create partial installation on failure', async () => {
      mockFs.readFile.mockRejectedValue(new Error('Read failed'));

      const result = await installer.install(defaultOptions);

      expect(result.success).toBe(false);
      expect(result.installedFiles.length).toBe(0);
    });
  });
});
