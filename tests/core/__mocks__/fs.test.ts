import { describe, it, expect, beforeEach } from 'vitest';
import {
  mockFs,
  resetMockFs,
  getMockFsOperations,
  getOperationsByName,
  getLastOperation,
  setMockFileContent,
  setMockDirectory,
  getMockFileContent,
  mockDirectoryExists,
} from './fs';

describe('File System Mocks', () => {
  beforeEach(() => {
    resetMockFs();
  });

  describe('Operation Tracking', () => {
    it('tracks readFile operations', async () => {
      setMockFileContent('/test/file.txt', 'test content');

      await mockFs.readFile('/test/file.txt', 'utf-8');

      const operations = getMockFsOperations();
      expect(operations).toHaveLength(1);
      expect(operations[0].operation).toBe('readFile');
      expect(operations[0].args).toEqual(['/test/file.txt', 'utf-8']);
    });

    it('tracks writeFile operations', async () => {
      await mockFs.writeFile('/test/file.txt', 'new content', 'utf-8');

      const operations = getMockFsOperations();
      expect(operations).toHaveLength(1);
      expect(operations[0].operation).toBe('writeFile');
    });

    it('filters operations by name', async () => {
      setMockFileContent('/test/file1.txt', 'content1');
      setMockFileContent('/test/file2.txt', 'content2');

      await mockFs.readFile('/test/file1.txt', 'utf-8');
      await mockFs.writeFile('/test/file3.txt', 'content3', 'utf-8');
      await mockFs.readFile('/test/file2.txt', 'utf-8');

      const readOps = getOperationsByName('readFile');
      const writeOps = getOperationsByName('writeFile');

      expect(readOps).toHaveLength(2);
      expect(writeOps).toHaveLength(1);
    });

    it('gets last operation', async () => {
      setMockFileContent('/test/file.txt', 'content');

      await mockFs.readFile('/test/file.txt', 'utf-8');
      await mockFs.writeFile('/test/file.txt', 'new content', 'utf-8');

      const lastOp = getLastOperation();

      expect(lastOp).toBeDefined();
      expect(lastOp?.operation).toBe('writeFile');
    });
  });

  describe('Mock State Management', () => {
    it('stores and retrieves file content', async () => {
      setMockFileContent('/test/file.txt', 'test content');

      const content = await mockFs.readFile('/test/file.txt', 'utf-8');

      expect(content).toBe('test content');
    });

    it('throws error for non-existent file', async () => {
      await expect(mockFs.readFile('/non/existent.txt', 'utf-8')).rejects.toThrow(
        'ENOENT'
      );
    });

    it('creates files when writeFile is called', async () => {
      await mockFs.writeFile('/test/newfile.txt', 'new content', 'utf-8');

      const content = getMockFileContent('/test/newfile.txt');
      expect(content).toBe('new content');
    });

    it('tracks directory creation', async () => {
      await mockFs.mkdir('/test/dir', { recursive: true });

      expect(mockDirectoryExists('/test/dir')).toBe(true);
    });

    it('removes files with rm', async () => {
      setMockFileContent('/test/file.txt', 'content');

      await mockFs.rm('/test/file.txt');

      expect(getMockFileContent('/test/file.txt')).toBeUndefined();
    });

    it('removes directories with rm', async () => {
      setMockDirectory('/test/dir');

      await mockFs.rm('/test/dir');

      expect(mockDirectoryExists('/test/dir')).toBe(false);
    });

    it('throws error when removing non-existent file without force', async () => {
      await expect(mockFs.rm('/non/existent.txt')).rejects.toThrow('ENOENT');
    });

    it('does not throw error when removing non-existent file with force', async () => {
      await expect(
        mockFs.rm('/non/existent.txt', { force: true })
      ).resolves.toBeUndefined();
    });
  });

  describe('File System Operations', () => {
    it('access checks file existence', async () => {
      setMockFileContent('/test/file.txt', 'content');

      await expect(mockFs.access('/test/file.txt')).resolves.toBeUndefined();
    });

    it('access throws for non-existent file', async () => {
      await expect(mockFs.access('/non/existent.txt')).rejects.toThrow('ENOENT');
    });

    it('stat returns file information', async () => {
      setMockFileContent('/test/file.txt', 'test content');

      const stats = await mockFs.stat('/test/file.txt');

      expect(stats.size).toBe('test content'.length);
      expect(stats.isFile()).toBe(true);
      expect(stats.isDirectory()).toBe(false);
      expect(stats.mtime).toBeInstanceOf(Date);
    });

    it('stat returns directory information', async () => {
      setMockDirectory('/test/dir');

      const stats = await mockFs.stat('/test/dir');

      expect(stats.size).toBe(0);
      expect(stats.isFile()).toBe(false);
      expect(stats.isDirectory()).toBe(true);
    });

    it('readdir lists files in directory', async () => {
      setMockDirectory('/test/dir');
      setMockFileContent('/test/dir/file1.txt', 'content1');
      setMockFileContent('/test/dir/file2.txt', 'content2');

      const files = await mockFs.readdir('/test/dir');

      expect(files).toContain('file1.txt');
      expect(files).toContain('file2.txt');
      expect(files).toHaveLength(2);
    });

    it('readdir throws for non-existent directory', async () => {
      await expect(mockFs.readdir('/non/existent')).rejects.toThrow('ENOENT');
    });

    it('copyFile copies file content', async () => {
      setMockFileContent('/test/source.txt', 'source content');

      await mockFs.copyFile('/test/source.txt', '/test/dest.txt');

      const destContent = getMockFileContent('/test/dest.txt');
      expect(destContent).toBe('source content');
    });

    it('copyFile throws for non-existent source', async () => {
      await expect(
        mockFs.copyFile('/non/existent.txt', '/test/dest.txt')
      ).rejects.toThrow('ENOENT');
    });

    it('unlink removes file', async () => {
      setMockFileContent('/test/file.txt', 'content');

      await mockFs.unlink('/test/file.txt');

      expect(getMockFileContent('/test/file.txt')).toBeUndefined();
    });

    it('rmdir removes directory', async () => {
      setMockDirectory('/test/dir');

      await mockFs.rmdir('/test/dir');

      expect(mockDirectoryExists('/test/dir')).toBe(false);
    });
  });

  describe('resetMockFs', () => {
    it('clears all operations', async () => {
      setMockFileContent('/test/file.txt', 'content');
      await mockFs.readFile('/test/file.txt', 'utf-8');

      resetMockFs();

      const operations = getMockFsOperations();
      expect(operations).toHaveLength(0);
    });

    it('clears all mock state', async () => {
      setMockFileContent('/test/file.txt', 'content');
      setMockDirectory('/test/dir');

      resetMockFs();

      expect(getMockFileContent('/test/file.txt')).toBeUndefined();
      expect(mockDirectoryExists('/test/dir')).toBe(false);
    });

    it('clears mock call history', async () => {
      setMockFileContent('/test/file.txt', 'content');
      await mockFs.readFile('/test/file.txt', 'utf-8');

      expect(mockFs.readFile).toHaveBeenCalledTimes(1);

      resetMockFs();

      expect(mockFs.readFile).toHaveBeenCalledTimes(0);
    });
  });
});
