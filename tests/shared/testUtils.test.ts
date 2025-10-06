import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import {
  createTempDir,
  cleanupTempDir,
  createMockLockFile,
  createMockSettings,
  createMockAssetStructure,
  mockFileSystem,
} from './testUtils';
import { promises as fs } from 'fs';

describe('Test Utilities', () => {
  describe('createTempDir and cleanupTempDir', () => {
    let tempDir: string;

    afterEach(async () => {
      if (tempDir) {
        await cleanupTempDir(tempDir);
      }
    });

    it('creates a temporary directory', async () => {
      tempDir = await createTempDir();

      expect(tempDir).toBeDefined();
      expect(tempDir).toContain('the-startup-test-');

      // Verify directory exists
      const stats = await fs.stat(tempDir);
      expect(stats.isDirectory()).toBe(true);
    });

    it('cleans up temporary directory', async () => {
      tempDir = await createTempDir();
      await cleanupTempDir(tempDir);

      // Verify directory no longer exists
      await expect(fs.access(tempDir)).rejects.toThrow();
    });
  });

  describe('createMockLockFile', () => {
    it('creates v1 lock file with string arrays', () => {
      const lockFile = createMockLockFile('0.9.0', 1);

      expect(lockFile.version).toBe('0.9.0');
      expect(lockFile.installedAt).toBeDefined();
      expect(Array.isArray(lockFile.files.agents)).toBe(true);
      expect(Array.isArray(lockFile.files.commands)).toBe(true);
      expect(typeof lockFile.files.binary).toBe('string');
    });

    it('creates v2 lock file with FileEntry objects', () => {
      const lockFile = createMockLockFile('1.0.0', 2);

      expect(lockFile.version).toBe('1.0.0');
      expect(lockFile.installedAt).toBeDefined();
      expect('lockFileVersion' in lockFile).toBe(true);

      if ('lockFileVersion' in lockFile) {
        expect(lockFile.lockFileVersion).toBe(2);
        expect(Array.isArray(lockFile.files.agents)).toBe(true);

        const firstAgent = lockFile.files.agents[0];
        expect(firstAgent).toHaveProperty('path');
        expect(firstAgent).toHaveProperty('checksum');
      }
    });
  });

  describe('createMockSettings', () => {
    it('creates empty settings without hooks', () => {
      const settings = createMockSettings(false);

      expect(settings).toBeDefined();
      expect(settings.mcpServers).toBeDefined();
      expect(settings.hooks).toBeUndefined();
    });

    it('creates settings with user hooks', () => {
      const settings = createMockSettings(true);

      expect(settings).toBeDefined();
      expect(settings.hooks).toBeDefined();

      if (settings.hooks) {
        expect(Object.keys(settings.hooks).length).toBeGreaterThan(0);

        const firstHook = Object.values(settings.hooks)[0];
        expect(firstHook).toHaveProperty('command');
      }
    });
  });

  describe('createMockAssetStructure', () => {
    it('creates asset map with all required files', () => {
      const assets = createMockAssetStructure();

      expect(assets.size).toBeGreaterThan(0);
      expect(assets.has('agents/specify.md')).toBe(true);
      expect(assets.has('commands/s-specify.md')).toBe(true);
      expect(assets.has('templates/SPEC.md')).toBe(true);
      expect(assets.has('settings.json')).toBe(true);
    });

    it('all asset files have content', () => {
      const assets = createMockAssetStructure();

      for (const [path, content] of assets.entries()) {
        expect(content).toBeDefined();
        expect(content.length).toBeGreaterThan(0);
      }
    });
  });

  describe('mockFileSystem', () => {
    let mockFs: ReturnType<typeof mockFileSystem>;

    beforeEach(() => {
      mockFs = mockFileSystem();
    });

    it('creates mock file system with all required methods', () => {
      expect(mockFs.readFile).toBeDefined();
      expect(mockFs.writeFile).toBeDefined();
      expect(mockFs.mkdir).toBeDefined();
      expect(mockFs.rm).toBeDefined();
      expect(mockFs.access).toBeDefined();
      expect(mockFs.stat).toBeDefined();
      expect(mockFs.readdir).toBeDefined();
      expect(mockFs.copyFile).toBeDefined();
    });

    it('allows mocking file operations', async () => {
      const testContent = 'test file content';
      mockFs.readFile.mockResolvedValue(testContent);

      const result = await mockFs.readFile('/test/file.txt', 'utf-8');

      expect(result).toBe(testContent);
      expect(mockFs.readFile).toHaveBeenCalledWith('/test/file.txt', 'utf-8');
    });
  });
});
