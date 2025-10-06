import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { LockManager } from '../../../src/core/installer/LockManager';
import { LockFile, FileEntry } from '../../../src/core/types/lock';
import {
  createTempDir,
  cleanupTempDir,
  createMockLockFile,
} from '../../shared/testUtils';
import { promises as fs } from 'fs';
import * as path from 'path';
import * as crypto from 'crypto';

describe('LockManager', () => {
  let tempDir: string;
  let lockFilePath: string;
  let lockManager: LockManager;

  beforeEach(async () => {
    tempDir = await createTempDir();
    lockFilePath = path.join(tempDir, 'lock.json');
    lockManager = new LockManager(lockFilePath);
  });

  afterEach(async () => {
    await cleanupTempDir(tempDir);
  });

  describe('readLockFile', () => {
    it('returns null when lock file does not exist', async () => {
      const result = await lockManager.readLockFile();
      expect(result).toBeNull();
    });

    it('reads v2 lock file with checksums', async () => {
      const mockLockV2 = createMockLockFile('1.0.0', 2);
      await fs.writeFile(lockFilePath, JSON.stringify(mockLockV2, null, 2));

      const result = await lockManager.readLockFile();

      expect(result).not.toBeNull();
      expect(result?.lockFileVersion).toBe(2);
      expect(result?.version).toBe('1.0.0');
      expect(result?.files.agents).toHaveLength(2);
      expect(result?.files.agents[0]).toHaveProperty('path');
      expect(result?.files.agents[0]).toHaveProperty('checksum');
    });

    it('reads v1 legacy lock file and auto-migrates to v2 format', async () => {
      const mockLockV1 = createMockLockFile('0.9.0', 1);
      await fs.writeFile(lockFilePath, JSON.stringify(mockLockV1, null, 2));

      const result = await lockManager.readLockFile();

      expect(result).not.toBeNull();
      expect(result?.lockFileVersion).toBe(1);
      expect(result?.version).toBe('0.9.0');

      // Verify migration to FileEntry[] format
      expect(Array.isArray(result?.files.agents)).toBe(true);
      expect(result?.files.agents[0]).toHaveProperty('path');
      expect(result?.files.agents[0]).toHaveProperty('checksum');
      expect(result?.files.agents[0].checksum).toBeUndefined();

      // Verify binary is converted to FileEntry
      expect(result?.files.binary).toHaveProperty('path');
      expect(result?.files.binary).toHaveProperty('checksum');
      expect(result?.files.binary.checksum).toBeUndefined();
    });

    it('handles corrupted JSON gracefully', async () => {
      await fs.writeFile(lockFilePath, 'invalid json{{{');

      await expect(lockManager.readLockFile()).rejects.toThrow();
    });

    it('detects v1 format by checking if files.agents[0] is string', async () => {
      const legacyLock = {
        version: '0.8.0',
        installedAt: new Date().toISOString(),
        files: {
          agents: ['.claude/agents/specify.md'],
          commands: ['.claude/commands/s-specify.md'],
          templates: ['.the-startup/templates/SPEC.md'],
          rules: [],
          outputStyles: [],
          binary: '.the-startup/bin/the-startup',
        },
      };
      await fs.writeFile(lockFilePath, JSON.stringify(legacyLock, null, 2));

      const result = await lockManager.readLockFile();

      expect(result).not.toBeNull();
      expect(result?.lockFileVersion).toBe(1);
      expect(result?.files.agents[0].path).toBe('.claude/agents/specify.md');
      expect(result?.files.agents[0].checksum).toBeUndefined();
    });

    it('treats lock file with empty agents array as v2 format', async () => {
      const emptyArrayLock = {
        version: '1.0.0',
        installedAt: new Date().toISOString(),
        files: {
          agents: [],
          commands: [],
          templates: [],
          rules: [],
          outputStyles: [],
          binary: { path: '.the-startup/bin/the-startup', checksum: 'abc123' },
        },
      };
      await fs.writeFile(lockFilePath, JSON.stringify(emptyArrayLock, null, 2));

      const result = await lockManager.readLockFile();

      expect(result).not.toBeNull();
      expect(result?.files.agents).toEqual([]);
      expect(result?.files.binary.path).toBe('.the-startup/bin/the-startup');
    });
  });

  describe('writeLockFile', () => {
    it('writes v2 lock file with SHA-256 checksums', async () => {
      const testFilePath = path.join(tempDir, 'test-file.txt');
      const testContent = 'test file content';
      await fs.writeFile(testFilePath, testContent);

      const expectedChecksum = crypto
        .createHash('sha256')
        .update(testContent)
        .digest('hex');

      const fileEntries: FileEntry[] = [
        { path: testFilePath, checksum: expectedChecksum },
      ];

      await lockManager.writeLockFile(fileEntries, '1.0.0');

      const writtenContent = await fs.readFile(lockFilePath, 'utf-8');
      const writtenLock: LockFile = JSON.parse(writtenContent);

      expect(writtenLock.lockFileVersion).toBe(2);
      expect(writtenLock.version).toBe('1.0.0');
      expect(writtenLock.installedAt).toBeDefined();

      // Verify ISO timestamp format
      expect(() => new Date(writtenLock.installedAt)).not.toThrow();
    });

    it('categorizes files correctly by path patterns', async () => {
      const fileEntries: FileEntry[] = [
        { path: '.claude/agents/specify.md', checksum: 'abc123' },
        { path: '.claude/commands/s-specify.md', checksum: 'def456' },
        { path: '.the-startup/templates/SPEC.md', checksum: 'ghi789' },
        { path: '.the-startup/bin/the-startup', checksum: 'jkl012' },
      ];

      await lockManager.writeLockFile(fileEntries, '1.0.0');

      const result = await lockManager.readLockFile();

      expect(result?.files.agents).toHaveLength(1);
      expect(result?.files.commands).toHaveLength(1);
      expect(result?.files.templates).toHaveLength(1);
      expect(result?.files.binary.path).toBe('.the-startup/bin/the-startup');
    });

    it('overwrites existing lock file', async () => {
      const firstEntries: FileEntry[] = [
        { path: '.claude/agents/first.md', checksum: 'abc' },
      ];
      await lockManager.writeLockFile(firstEntries, '1.0.0');

      const secondEntries: FileEntry[] = [
        { path: '.claude/agents/second.md', checksum: 'def' },
      ];
      await lockManager.writeLockFile(secondEntries, '1.1.0');

      const result = await lockManager.readLockFile();

      expect(result?.version).toBe('1.1.0');
      expect(result?.files.agents).toHaveLength(1);
      expect(result?.files.agents[0].path).toBe('.claude/agents/second.md');
    });
  });

  describe('generateChecksum', () => {
    it('generates SHA-256 checksum for file content', async () => {
      const testFilePath = path.join(tempDir, 'checksum-test.txt');
      const testContent = 'checksum test content';
      await fs.writeFile(testFilePath, testContent);

      const expectedChecksum = crypto
        .createHash('sha256')
        .update(testContent)
        .digest('hex');

      const result = await lockManager.generateChecksum(testFilePath);

      expect(result).toBe(expectedChecksum);
      expect(result).toHaveLength(64); // SHA-256 is 64 hex characters
    });

    it('generates different checksums for different content', async () => {
      const file1 = path.join(tempDir, 'file1.txt');
      const file2 = path.join(tempDir, 'file2.txt');

      await fs.writeFile(file1, 'content 1');
      await fs.writeFile(file2, 'content 2');

      const checksum1 = await lockManager.generateChecksum(file1);
      const checksum2 = await lockManager.generateChecksum(file2);

      expect(checksum1).not.toBe(checksum2);
    });

    it('generates same checksum for identical content', async () => {
      const file1 = path.join(tempDir, 'identical1.txt');
      const file2 = path.join(tempDir, 'identical2.txt');
      const sameContent = 'identical content';

      await fs.writeFile(file1, sameContent);
      await fs.writeFile(file2, sameContent);

      const checksum1 = await lockManager.generateChecksum(file1);
      const checksum2 = await lockManager.generateChecksum(file2);

      expect(checksum1).toBe(checksum2);
    });

    it('throws error for non-existent file', async () => {
      const nonExistentPath = path.join(tempDir, 'does-not-exist.txt');

      await expect(lockManager.generateChecksum(nonExistentPath))
        .rejects.toThrow();
    });
  });

  describe('compareChecksums', () => {
    it('returns true when checksums match', async () => {
      const testFilePath = path.join(tempDir, 'compare-test.txt');
      await fs.writeFile(testFilePath, 'test content');

      const checksum = await lockManager.generateChecksum(testFilePath);
      const fileEntry: FileEntry = { path: testFilePath, checksum };

      const result = await lockManager.compareChecksums(fileEntry);

      expect(result).toBe(true);
    });

    it('returns false when checksums do not match (file modified)', async () => {
      const testFilePath = path.join(tempDir, 'modified-test.txt');
      await fs.writeFile(testFilePath, 'original content');

      const originalChecksum = await lockManager.generateChecksum(testFilePath);

      // Modify the file
      await fs.writeFile(testFilePath, 'modified content');

      const fileEntry: FileEntry = { path: testFilePath, checksum: originalChecksum };
      const result = await lockManager.compareChecksums(fileEntry);

      expect(result).toBe(false);
    });

    it('returns false when FileEntry has no checksum (legacy)', async () => {
      const testFilePath = path.join(tempDir, 'legacy-test.txt');
      await fs.writeFile(testFilePath, 'test content');

      const fileEntry: FileEntry = { path: testFilePath }; // No checksum
      const result = await lockManager.compareChecksums(fileEntry);

      expect(result).toBe(false);
    });

    it('throws error when file does not exist', async () => {
      const nonExistentPath = path.join(tempDir, 'missing.txt');
      const fileEntry: FileEntry = { path: nonExistentPath, checksum: 'abc123' };

      await expect(lockManager.compareChecksums(fileEntry))
        .rejects.toThrow();
    });
  });

  describe('getFilesNeedingReinstall', () => {
    it('returns all files when no lock file exists', async () => {
      const targetFiles: FileEntry[] = [
        { path: '.claude/agents/specify.md', checksum: 'abc' },
        { path: '.claude/commands/s-specify.md', checksum: 'def' },
      ];

      const result = await lockManager.getFilesNeedingReinstall(targetFiles);

      expect(result).toHaveLength(2);
      expect(result).toEqual(targetFiles);
    });

    it('skips files with matching checksums', async () => {
      const file1 = path.join(tempDir, '.claude', 'agents', 'unchanged.md');
      const file2 = path.join(tempDir, '.claude', 'agents', 'modified.md');

      await fs.mkdir(path.dirname(file1), { recursive: true });
      await fs.writeFile(file1, 'unchanged content');
      await fs.writeFile(file2, 'original content');

      const checksum1 = await lockManager.generateChecksum(file1);
      const checksum2Original = await lockManager.generateChecksum(file2);

      const existingEntries: FileEntry[] = [
        { path: file1, checksum: checksum1 },
        { path: file2, checksum: checksum2Original },
      ];
      await lockManager.writeLockFile(existingEntries, '1.0.0');

      // Modify file2
      await fs.writeFile(file2, 'modified content');
      const checksum2New = await lockManager.generateChecksum(file2);

      const targetFiles: FileEntry[] = [
        { path: file1, checksum: checksum1 },
        { path: file2, checksum: checksum2New },
      ];

      const result = await lockManager.getFilesNeedingReinstall(targetFiles);

      expect(result).toHaveLength(1);
      expect(result[0].path).toBe(file2);
    });

    it('reinstalls files from legacy lock file (no checksums)', async () => {
      const mockLockV1 = createMockLockFile('0.9.0', 1);
      await fs.writeFile(lockFilePath, JSON.stringify(mockLockV1, null, 2));

      const targetFiles: FileEntry[] = [
        { path: '.claude/agents/specify.md', checksum: 'new-checksum' },
      ];

      const result = await lockManager.getFilesNeedingReinstall(targetFiles);

      // Legacy entries have no checksums, so they should be marked for reinstall
      expect(result).toHaveLength(1);
    });

    it('marks new files not in lock for install', async () => {
      const existingFile = path.join(tempDir, '.claude', 'commands', 'existing.md');
      await fs.mkdir(path.dirname(existingFile), { recursive: true });
      await fs.writeFile(existingFile, 'existing content');
      const existingChecksum = await lockManager.generateChecksum(existingFile);

      await lockManager.writeLockFile(
        [{ path: existingFile, checksum: existingChecksum }],
        '1.0.0'
      );

      const newFile = path.join(tempDir, '.claude', 'commands', 'new-file.md');
      await fs.writeFile(newFile, 'new content');
      const newChecksum = await lockManager.generateChecksum(newFile);

      const targetFiles: FileEntry[] = [
        { path: existingFile, checksum: existingChecksum },
        { path: newFile, checksum: newChecksum },
      ];

      const result = await lockManager.getFilesNeedingReinstall(targetFiles);

      expect(result).toHaveLength(1);
      expect(result[0].path).toBe(newFile);
    });
  });

  describe('backward compatibility integration', () => {
    it('supports full v1 to v2 migration workflow', async () => {
      // Step 1: Create v1 lock file (simulating old installation)
      const mockLockV1 = createMockLockFile('0.9.0', 1);
      await fs.writeFile(lockFilePath, JSON.stringify(mockLockV1, null, 2));

      // Step 2: Read v1 lock file (auto-migrates)
      const migratedLock = await lockManager.readLockFile();
      expect(migratedLock?.lockFileVersion).toBe(1);
      expect(migratedLock?.files.agents[0].checksum).toBeUndefined();

      // Step 3: Create actual files for new installation
      const agentFile = path.join(tempDir, '.claude', 'agents', 'specify.md');
      await fs.mkdir(path.dirname(agentFile), { recursive: true });
      await fs.writeFile(agentFile, 'new agent content');

      const newChecksum = await lockManager.generateChecksum(agentFile);

      // Step 4: Write new v2 lock file with checksums
      await lockManager.writeLockFile(
        [{ path: agentFile, checksum: newChecksum }],
        '1.0.0'
      );

      // Step 5: Read v2 lock file
      const v2Lock = await lockManager.readLockFile();
      expect(v2Lock?.lockFileVersion).toBe(2);
      expect(v2Lock?.version).toBe('1.0.0');
      expect(v2Lock?.files.agents[0].checksum).toBeDefined();
    });

    it('enables idempotent reinstall with checksum comparison', async () => {
      // Initial installation
      const file1 = path.join(tempDir, '.the-startup', 'templates', 'file1.md');
      const file2 = path.join(tempDir, '.the-startup', 'templates', 'file2.md');
      const file3 = path.join(tempDir, '.the-startup', 'templates', 'file3.md');

      await fs.mkdir(path.dirname(file1), { recursive: true });
      await fs.writeFile(file1, 'content 1');
      await fs.writeFile(file2, 'content 2');
      await fs.writeFile(file3, 'content 3');

      const checksum1 = await lockManager.generateChecksum(file1);
      const checksum2 = await lockManager.generateChecksum(file2);
      const checksum3 = await lockManager.generateChecksum(file3);

      await lockManager.writeLockFile(
        [
          { path: file1, checksum: checksum1 },
          { path: file2, checksum: checksum2 },
          { path: file3, checksum: checksum3 },
        ],
        '1.0.0'
      );

      // User modifies file2
      await fs.writeFile(file2, 'modified content 2');
      const newChecksum2 = await lockManager.generateChecksum(file2);

      // Reinstall attempt
      const targetFiles: FileEntry[] = [
        { path: file1, checksum: checksum1 },
        { path: file2, checksum: newChecksum2 },
        { path: file3, checksum: checksum3 },
      ];

      const needReinstall = await lockManager.getFilesNeedingReinstall(targetFiles);

      // Only file2 should need reinstall (modified by user)
      expect(needReinstall).toHaveLength(1);
      expect(needReinstall[0].path).toBe(file2);
    });
  });
});
