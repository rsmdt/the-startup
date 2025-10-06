import { describe, it, expect, vi } from 'vitest';
import React from 'react';
import { render } from 'ink-testing-library';
import { UninstallWizard } from '../../../src/ui/uninstall/UninstallWizard';
import type { LockFile } from '../../../src/core/types/lock';

describe('UninstallWizard', () => {
  // Mock lock manager
  const createMockLockManager = (lockFile: LockFile | null = null) => ({
    readLockFile: vi.fn().mockResolvedValue(lockFile),
    writeLockFile: vi.fn().mockResolvedValue(undefined),
    generateChecksum: vi.fn().mockResolvedValue('mock-checksum'),
    compareChecksums: vi.fn().mockResolvedValue(true),
    getFilesNeedingReinstall: vi.fn().mockResolvedValue([]),
  });

  // Mock settings merger
  const createMockSettingsMerger = () => ({
    mergeSettings: vi.fn().mockResolvedValue({}),
    removeHooks: vi.fn().mockResolvedValue({}),
  });

  // Mock file system
  const createMockFileSystem = () => ({
    rm: vi.fn().mockResolvedValue(undefined),
  });

  // Sample lock file
  const sampleLockFile: LockFile = {
    version: '1.0.0',
    installedAt: '2025-01-01T00:00:00.000Z',
    lockFileVersion: 2,
    files: {
      agents: [
        { path: '/test/.claude/agents/test-agent.md', checksum: 'abc123' },
      ],
      commands: [
        { path: '/test/.claude/commands/test-command.md', checksum: 'def456' },
      ],
      templates: [
        { path: '/test/.the-startup/templates/test-template.md', checksum: 'ghi789' },
      ],
      rules: [
        { path: '/test/.the-startup/rules/test-rule.md', checksum: 'jkl012' },
      ],
      outputStyles: [
        { path: '/test/.the-startup/outputStyles/test-style.md', checksum: 'mno345' },
      ],
      binary: {
        path: '/test/.the-startup/bin/the-startup',
        checksum: 'pqr678',
      },
    },
  };

  describe('State Machine', () => {
    it('renders confirmation state with file list', async () => {
      const mockLockManager = createMockLockManager(sampleLockFile);
      const mockSettingsMerger = createMockSettingsMerger();
      const mockFileSystem = createMockFileSystem();

      const { lastFrame } = render(
        <UninstallWizard
          options={{}}
          lockManager={mockLockManager as any}
          settingsMerger={mockSettingsMerger as any}
          onComplete={() => {}}
          fileSystem={mockFileSystem as any}
        />
      );

      // Wait for lock file to load
      await vi.waitFor(() => {
        const output = lastFrame();
        expect(output).toContain('files will be removed');
        expect(output).toContain('6');
      });
    });

    it('displays error when lock file does not exist', async () => {
      const mockLockManager = createMockLockManager(null);
      const mockSettingsMerger = createMockSettingsMerger();
      const mockFileSystem = createMockFileSystem();

      const { lastFrame } = render(
        <UninstallWizard
          options={{}}
          lockManager={mockLockManager as any}
          settingsMerger={mockSettingsMerger as any}
          onComplete={() => {}}
          fileSystem={mockFileSystem as any}
        />
      );

      await vi.waitFor(() => {
        const output = lastFrame();
        expect(output).toContain('No installation found');
      });
    });

    it('transitions to uninstalling state after yes confirmation', async () => {
      const mockLockManager = createMockLockManager(sampleLockFile);
      const mockSettingsMerger = createMockSettingsMerger();
      const mockFileSystem = createMockFileSystem();

      const { lastFrame, stdin } = render(
        <UninstallWizard
          options={{}}
          lockManager={mockLockManager as any}
          settingsMerger={mockSettingsMerger as any}
          onComplete={() => {}}
          fileSystem={mockFileSystem as any}
        />
      );

      // Wait for confirmation
      await vi.waitFor(() => {
        expect(lastFrame()).toContain('files will be removed');
      });

      // Confirm uninstall
      stdin.write('y');

      // Should transition to uninstalling
      await vi.waitFor(() => {
        const output = lastFrame();
        expect(output).toMatch(/uninstalling|removing/i);
      });
    });

    it('calls onComplete after successful uninstall', async () => {
      const mockLockManager = createMockLockManager(sampleLockFile);
      const mockSettingsMerger = createMockSettingsMerger();
      const mockFileSystem = createMockFileSystem();
      const onComplete = vi.fn();

      const { stdin, lastFrame } = render(
        <UninstallWizard
          options={{}}
          lockManager={mockLockManager as any}
          settingsMerger={mockSettingsMerger as any}
          onComplete={onComplete}
          fileSystem={mockFileSystem as any}
        />
      );

      // Wait for confirmation
      await vi.waitFor(() => {
        expect(lastFrame()).toContain('files will be removed');
      });

      // Confirm
      stdin.write('y');

      // Should complete
      await vi.waitFor(() => {
        expect(onComplete).toHaveBeenCalledWith(
          expect.objectContaining({ success: true })
        );
      }, { timeout: 3000 });
    });

    it('cancels uninstall on n input', async () => {
      const mockLockManager = createMockLockManager(sampleLockFile);
      const mockSettingsMerger = createMockSettingsMerger();
      const mockFileSystem = createMockFileSystem();
      const onComplete = vi.fn();

      const { stdin, lastFrame } = render(
        <UninstallWizard
          options={{}}
          lockManager={mockLockManager as any}
          settingsMerger={mockSettingsMerger as any}
          onComplete={onComplete}
          fileSystem={mockFileSystem as any}
        />
      );

      // Wait for confirmation
      await vi.waitFor(() => {
        expect(lastFrame()).toContain('files will be removed');
      });

      // Cancel
      stdin.write('n');

      // Should call onComplete with cancelled
      await vi.waitFor(() => {
        expect(onComplete).toHaveBeenCalledWith(
          expect.objectContaining({ success: false })
        );
      });
    });
  });

  describe('File Deletion', () => {
    it('deletes all files from lock file', async () => {
      const mockLockManager = createMockLockManager(sampleLockFile);
      const mockSettingsMerger = createMockSettingsMerger();
      const mockFileSystem = createMockFileSystem();

      const { stdin, lastFrame } = render(
        <UninstallWizard
          options={{}}
          lockManager={mockLockManager as any}
          settingsMerger={mockSettingsMerger as any}
          onComplete={() => {}}
          fileSystem={mockFileSystem as any}
        />
      );

      await vi.waitFor(() => {
        expect(lastFrame()).toContain('files will be removed');
      });

      stdin.write('y');

      // Should delete files
      await vi.waitFor(() => {
        expect(mockFileSystem.rm).toHaveBeenCalled();
      }, { timeout: 3000 });
    });

    it('handles missing files gracefully (ENOENT)', async () => {
      const mockLockManager = createMockLockManager(sampleLockFile);
      const mockSettingsMerger = createMockSettingsMerger();
      const mockFileSystem = {
        rm: vi.fn().mockRejectedValue({ code: 'ENOENT' }),
      };
      const onComplete = vi.fn();

      const { stdin, lastFrame } = render(
        <UninstallWizard
          options={{}}
          lockManager={mockLockManager as any}
          settingsMerger={mockSettingsMerger as any}
          onComplete={onComplete}
          fileSystem={mockFileSystem as any}
        />
      );

      await vi.waitFor(() => {
        expect(lastFrame()).toContain('files will be removed');
      });

      stdin.write('y');

      // Should complete successfully despite missing files
      await vi.waitFor(() => {
        expect(onComplete).toHaveBeenCalledWith(
          expect.objectContaining({ success: true })
        );
      }, { timeout: 3000 });
    });
  });

  describe('--keep-logs Flag', () => {
    it('preserves logs when --keep-logs is provided', async () => {
      const mockLockManager = createMockLockManager({
        ...sampleLockFile,
        files: {
          ...sampleLockFile.files,
          templates: [
            ...sampleLockFile.files.templates,
            { path: '/test/.the-startup/logs/test.log', checksum: 'log123' },
          ],
        },
      });
      const mockSettingsMerger = createMockSettingsMerger();
      const mockFileSystem = createMockFileSystem();

      const { stdin, lastFrame } = render(
        <UninstallWizard
          options={{ keepLogs: true }}
          lockManager={mockLockManager as any}
          settingsMerger={mockSettingsMerger as any}
          onComplete={() => {}}
          fileSystem={mockFileSystem as any}
        />
      );

      await vi.waitFor(() => {
        expect(lastFrame()).toContain('Logs will be preserved');
      });

      stdin.write('y');

      await vi.waitFor(() => {
        // Should not delete log files
        const calls = mockFileSystem.rm.mock.calls;
        const logDeleted = calls.some((call: any[]) =>
          call[0].includes('/logs/')
        );
        expect(logDeleted).toBe(false);
      }, { timeout: 3000 });
    });
  });

  describe('--keep-settings Flag', () => {
    it('does not modify settings when --keep-settings is provided', async () => {
      const mockLockManager = createMockLockManager(sampleLockFile);
      const mockSettingsMerger = createMockSettingsMerger();
      const mockFileSystem = createMockFileSystem();

      const { stdin, lastFrame } = render(
        <UninstallWizard
          options={{ keepSettings: true }}
          lockManager={mockLockManager as any}
          settingsMerger={mockSettingsMerger as any}
          onComplete={() => {}}
          fileSystem={mockFileSystem as any}
        />
      );

      await vi.waitFor(() => {
        expect(lastFrame()).toContain('Settings.json will not be modified');
      });

      stdin.write('y');

      await vi.waitFor(() => {
        expect(mockSettingsMerger.removeHooks).not.toHaveBeenCalled();
      }, { timeout: 3000 });
    });

    it('removes hooks when --keep-settings is not provided', async () => {
      const mockLockManager = createMockLockManager(sampleLockFile);
      const mockSettingsMerger = createMockSettingsMerger();
      const mockFileSystem = createMockFileSystem();

      const { stdin, lastFrame } = render(
        <UninstallWizard
          options={{}}
          lockManager={mockLockManager as any}
          settingsMerger={mockSettingsMerger as any}
          onComplete={() => {}}
          fileSystem={mockFileSystem as any}
        />
      );

      await vi.waitFor(() => {
        expect(lastFrame()).toContain('files will be removed');
      });

      stdin.write('y');

      await vi.waitFor(() => {
        expect(mockSettingsMerger.removeHooks).toHaveBeenCalled();
      }, { timeout: 3000 });
    });
  });

  describe('Component Pattern', () => {
    it('follows Ink functional component pattern', () => {
      expect(typeof UninstallWizard).toBe('function');
    });

    it('returns valid React element', () => {
      const mockLockManager = createMockLockManager(sampleLockFile);
      const mockSettingsMerger = createMockSettingsMerger();
      const mockFileSystem = createMockFileSystem();

      const element = (
        <UninstallWizard
          options={{}}
          lockManager={mockLockManager as any}
          settingsMerger={mockSettingsMerger as any}
          onComplete={() => {}}
          fileSystem={mockFileSystem as any}
        />
      );
      expect(React.isValidElement(element)).toBe(true);
    });
  });
});
