import { describe, it, expect, beforeEach } from 'vitest';
import {
  mockFs,
  resetMockFs,
  getMockFileContent,
  setMockFileContent,
  getOperationsByName,
} from '../../core/__mocks__/fs';
import { SettingsMerger } from '../../../src/core/installer/SettingsMerger';
import type { ClaudeSettings, PlaceholderMap } from '../../../src/core/types/settings';

/**
 * SettingsMerger Test Suite
 *
 * Tests the deep merge algorithm for Claude settings.json with:
 * - User data preservation (never overwrite existing hooks)
 * - Placeholder replacement ({{STARTUP_PATH}}, {{CLAUDE_PATH}})
 * - Backup and rollback on failure
 * - Atomic operations (all-or-nothing)
 *
 * TDD Approach: Red-Green-Refactor
 * - Write failing tests first (RED)
 * - Implement minimal code to pass (GREEN)
 * - Refactor while keeping tests green (REFACTOR)
 */

describe('SettingsMerger', () => {
  const settingsPath = '/Users/test/.claude/settings.json';
  const placeholders: PlaceholderMap = {
    STARTUP_PATH: '/Users/test/.the-startup',
    CLAUDE_PATH: '/Users/test/.claude',
  };

  let merger: SettingsMerger;

  beforeEach(() => {
    resetMockFs();
    merger = new SettingsMerger(mockFs as any);
  });

  describe('Deep Merge with User Preservation', () => {
    it('should preserve existing user hooks and add new hooks', async () => {
      // Arrange: User has custom hook configuration
      const existingSettings: ClaudeSettings = {
        hooks: {
          'user-prompt-submit': {
            command: '/custom/statusline.sh',
            description: 'Custom user statusline',
            continueOnError: true,
          },
        },
        someOtherSetting: 'user value',
      };

      setMockFileContent(settingsPath, JSON.stringify(existingSettings, null, 2));

      const newHooks: ClaudeSettings['hooks'] = {
        'user-prompt-submit': {
          command: '{{STARTUP_PATH}}/bin/statusline.sh',
          description: 'Default statusline',
        },
        'new-hook': {
          command: '{{CLAUDE_PATH}}/scripts/new-script.sh',
          description: 'New hook',
        },
      };

      // Act
      const result = await merger.mergeSettings(settingsPath, newHooks, placeholders);

      // Assert: User's custom hook preserved, new hook added
      expect(result.hooks?.['user-prompt-submit']).toEqual({
        command: '/custom/statusline.sh',
        description: 'Custom user statusline',
        continueOnError: true,
      });

      expect(result.hooks?.['new-hook']).toEqual({
        command: '/Users/test/.claude/scripts/new-script.sh',
        description: 'New hook',
      });

      expect(result.someOtherSetting).toBe('user value');
    });

    it('should create hooks object if it does not exist', async () => {
      // Arrange: User has no hooks
      const existingSettings: ClaudeSettings = {
        someOtherSetting: 'user value',
      };

      setMockFileContent(settingsPath, JSON.stringify(existingSettings, null, 2));

      const newHooks: ClaudeSettings['hooks'] = {
        'user-prompt-submit': {
          command: '{{STARTUP_PATH}}/bin/statusline.sh',
        },
      };

      // Act
      const result = await merger.mergeSettings(settingsPath, newHooks, placeholders);

      // Assert: Hooks object created with new hook
      expect(result.hooks).toBeDefined();
      expect(result.hooks?.['user-prompt-submit']).toEqual({
        command: '/Users/test/.the-startup/bin/statusline.sh',
      });
      expect(result.someOtherSetting).toBe('user value');
    });

    it('should preserve all non-hook settings unchanged', async () => {
      // Arrange: User has various settings
      const existingSettings: ClaudeSettings = {
        hooks: {},
        theme: 'dark',
        editor: { fontSize: 14, fontFamily: 'Monaco' },
        customArray: [1, 2, 3],
      };

      setMockFileContent(settingsPath, JSON.stringify(existingSettings, null, 2));

      const newHooks: ClaudeSettings['hooks'] = {
        'test-hook': { command: 'test' },
      };

      // Act
      const result = await merger.mergeSettings(settingsPath, newHooks, placeholders);

      // Assert: All non-hook settings preserved exactly
      expect(result.theme).toBe('dark');
      expect(result.editor).toEqual({ fontSize: 14, fontFamily: 'Monaco' });
      expect(result.customArray).toEqual([1, 2, 3]);
    });
  });

  describe('Placeholder Replacement', () => {
    it('should replace {{STARTUP_PATH}} placeholder in hook commands', async () => {
      // Arrange
      setMockFileContent(settingsPath, JSON.stringify({ hooks: {} }, null, 2));

      const newHooks: ClaudeSettings['hooks'] = {
        'test-hook': {
          command: '{{STARTUP_PATH}}/bin/script.sh',
          description: 'Test hook',
        },
      };

      // Act
      const result = await merger.mergeSettings(settingsPath, newHooks, placeholders);

      // Assert
      expect(result.hooks?.['test-hook']?.command).toBe(
        '/Users/test/.the-startup/bin/script.sh'
      );
    });

    it('should replace {{CLAUDE_PATH}} placeholder in hook commands', async () => {
      // Arrange
      setMockFileContent(settingsPath, JSON.stringify({ hooks: {} }, null, 2));

      const newHooks: ClaudeSettings['hooks'] = {
        'test-hook': {
          command: '{{CLAUDE_PATH}}/scripts/script.sh',
        },
      };

      // Act
      const result = await merger.mergeSettings(settingsPath, newHooks, placeholders);

      // Assert
      expect(result.hooks?.['test-hook']?.command).toBe(
        '/Users/test/.claude/scripts/script.sh'
      );
    });

    it('should replace multiple placeholders in same command', async () => {
      // Arrange
      setMockFileContent(settingsPath, JSON.stringify({ hooks: {} }, null, 2));

      const newHooks: ClaudeSettings['hooks'] = {
        'test-hook': {
          command:
            '{{STARTUP_PATH}}/bin/script.sh --config {{CLAUDE_PATH}}/config.json',
        },
      };

      // Act
      const result = await merger.mergeSettings(settingsPath, newHooks, placeholders);

      // Assert
      expect(result.hooks?.['test-hook']?.command).toBe(
        '/Users/test/.the-startup/bin/script.sh --config /Users/test/.claude/config.json'
      );
    });

    it('should not replace placeholders in user-preserved hooks', async () => {
      // Arrange: User already has hook with placeholder-like text
      const existingSettings: ClaudeSettings = {
        hooks: {
          'test-hook': {
            command: 'echo {{STARTUP_PATH}} is my custom text',
          },
        },
      };

      setMockFileContent(settingsPath, JSON.stringify(existingSettings, null, 2));

      const newHooks: ClaudeSettings['hooks'] = {
        'test-hook': {
          command: '{{STARTUP_PATH}}/bin/script.sh',
        },
      };

      // Act
      const result = await merger.mergeSettings(settingsPath, newHooks, placeholders);

      // Assert: User's hook preserved exactly (no replacement)
      expect(result.hooks?.['test-hook']?.command).toBe(
        'echo {{STARTUP_PATH}} is my custom text'
      );
    });
  });

  describe('Backup and Rollback', () => {
    it('should create backup before modifying settings', async () => {
      // Arrange
      const existingSettings: ClaudeSettings = { hooks: {} };
      setMockFileContent(settingsPath, JSON.stringify(existingSettings, null, 2));

      const newHooks: ClaudeSettings['hooks'] = {
        'test-hook': { command: 'test' },
      };

      // Act
      await merger.mergeSettings(settingsPath, newHooks, placeholders);

      // Assert: Backup file created
      const copyOps = getOperationsByName('copyFile');
      expect(copyOps.length).toBeGreaterThan(0);

      const backupOp = copyOps.find((op) =>
        (op.args[1] as string).includes('settings.json.backup')
      );
      expect(backupOp).toBeDefined();
    });

    it('should rollback on write failure', async () => {
      // Arrange
      const existingSettings: ClaudeSettings = { hooks: {} };
      setMockFileContent(settingsPath, JSON.stringify(existingSettings, null, 2));

      // Mock writeFile to fail
      mockFs.writeFile.mockRejectedValueOnce(new Error('Write failed'));

      const newHooks: ClaudeSettings['hooks'] = {
        'test-hook': { command: 'test' },
      };

      // Act & Assert
      await expect(
        merger.mergeSettings(settingsPath, newHooks, placeholders)
      ).rejects.toThrow('Write failed');

      // Verify rollback: original content restored from backup
      const copyOps = getOperationsByName('copyFile');
      expect(copyOps.length).toBeGreaterThan(0);
    });

    it('should cleanup backup after successful merge', async () => {
      // Arrange
      const existingSettings: ClaudeSettings = { hooks: {} };
      setMockFileContent(settingsPath, JSON.stringify(existingSettings, null, 2));

      const newHooks: ClaudeSettings['hooks'] = {
        'test-hook': { command: 'test' },
      };

      // Act
      await merger.mergeSettings(settingsPath, newHooks, placeholders);

      // Assert: Backup removed after success
      const rmOps = getOperationsByName('rm');
      const backupRemoved = rmOps.some((op) =>
        (op.args[0] as string).includes('settings.json.backup')
      );
      expect(backupRemoved).toBe(true);
    });
  });

  describe('Settings Creation from Scratch', () => {
    it('should create settings.json if it does not exist', async () => {
      // Arrange: No settings file exists
      const newHooks: ClaudeSettings['hooks'] = {
        'test-hook': {
          command: '{{STARTUP_PATH}}/bin/script.sh',
        },
      };

      // Act
      const result = await merger.mergeSettings(settingsPath, newHooks, placeholders);

      // Assert: New settings created with hooks
      expect(result.hooks).toBeDefined();
      expect(result.hooks?.['test-hook']?.command).toBe(
        '/Users/test/.the-startup/bin/script.sh'
      );

      // Verify file was written
      const writeOps = getOperationsByName('writeFile');
      expect(writeOps.length).toBeGreaterThan(0);
    });

    it('should handle settings with existing empty hooks object', async () => {
      // Arrange: Settings file exists with empty hooks
      const existingSettings: ClaudeSettings = {
        hooks: {},
        otherSetting: 'preserved',
      };

      setMockFileContent(settingsPath, JSON.stringify(existingSettings, null, 2));

      const newHooks: ClaudeSettings['hooks'] = {
        'new-hook': {
          command: '{{STARTUP_PATH}}/bin/script.sh',
        },
      };

      // Act
      const result = await merger.mergeSettings(settingsPath, newHooks, placeholders);

      // Assert: Hook added to existing empty hooks
      expect(result.hooks).toBeDefined();
      expect(result.hooks?.['new-hook']?.command).toBe(
        '/Users/test/.the-startup/bin/script.sh'
      );
      expect(result.otherSetting).toBe('preserved');
    });
  });

  describe('Error Handling', () => {
    it('should throw clear error for invalid JSON in existing settings', async () => {
      // Arrange: Invalid JSON in settings file
      setMockFileContent(settingsPath, '{ invalid json }');

      const newHooks: ClaudeSettings['hooks'] = {
        'test-hook': { command: 'test' },
      };

      // Act & Assert
      await expect(
        merger.mergeSettings(settingsPath, newHooks, placeholders)
      ).rejects.toThrow(/JSON/);
    });

    it('should provide clear error message with file path on JSON parse failure', async () => {
      // Arrange
      setMockFileContent(settingsPath, '{ invalid }');

      const newHooks: ClaudeSettings['hooks'] = {
        'test-hook': { command: 'test' },
      };

      // Act & Assert
      try {
        await merger.mergeSettings(settingsPath, newHooks, placeholders);
        expect.fail('Should have thrown error');
      } catch (error) {
        expect(error).toBeInstanceOf(Error);
        expect((error as Error).message).toContain(settingsPath);
      }
    });

    it('should validate settings structure after merge', async () => {
      // Arrange
      setMockFileContent(settingsPath, JSON.stringify({ hooks: {} }, null, 2));

      const newHooks: ClaudeSettings['hooks'] = {
        'test-hook': { command: 'test' },
      };

      // Act
      const result = await merger.mergeSettings(settingsPath, newHooks, placeholders);

      // Assert: Result has valid structure
      expect(result).toBeTypeOf('object');
      expect(result.hooks).toBeTypeOf('object');
    });
  });

  describe('Atomic Operations', () => {
    it('should be all-or-nothing (no partial merge on failure)', async () => {
      // Arrange
      const existingSettings: ClaudeSettings = {
        hooks: {
          'existing-hook': { command: 'existing' },
        },
      };
      setMockFileContent(settingsPath, JSON.stringify(existingSettings, null, 2));

      // Mock writeFile to fail after backup
      mockFs.writeFile.mockRejectedValueOnce(new Error('Disk full'));

      const newHooks: ClaudeSettings['hooks'] = {
        'new-hook': { command: 'new' },
      };

      // Act & Assert
      await expect(
        merger.mergeSettings(settingsPath, newHooks, placeholders)
      ).rejects.toThrow();

      // Verify original file restored (rollback occurred)
      const content = getMockFileContent(settingsPath);
      expect(content).toBe(JSON.stringify(existingSettings, null, 2));
    });
  });
});
