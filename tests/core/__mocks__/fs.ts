import { vi } from 'vitest';

/**
 * File System Mocks
 *
 * Provides comprehensive mocking for Node.js fs module operations.
 * Tracks all file system operations for test assertions and verification.
 *
 * This mock supports both legacy callback-based APIs and promise-based APIs.
 * All operations are tracked in the `operations` array for verification.
 *
 * @example
 * import { mockFs, resetMockFs, getMockFsOperations } from '../__mocks__/fs';
 *
 * beforeEach(() => {
 *   resetMockFs();
 *   mockFs.readFile.mockResolvedValue('file content');
 * });
 *
 * test('reads file', async () => {
 *   const content = await readFile('/path/to/file', 'utf-8');
 *   const ops = getMockFsOperations();
 *   expect(ops).toContainEqual({ operation: 'readFile', args: ['/path/to/file', 'utf-8'] });
 * });
 */

// Operation tracking
interface FsOperation {
  operation: string;
  args: unknown[];
  timestamp: number;
}

const operations: FsOperation[] = [];

/**
 * Records a file system operation for later verification
 *
 * @param operation - Name of the fs operation (e.g., 'readFile', 'writeFile')
 * @param args - Arguments passed to the operation
 */
function recordOperation(operation: string, args: unknown[]): void {
  operations.push({
    operation,
    args,
    timestamp: Date.now(),
  });
}

/**
 * Gets all recorded file system operations
 *
 * @returns Array of recorded operations
 */
export function getMockFsOperations(): FsOperation[] {
  return [...operations];
}

/**
 * Clears all recorded operations
 */
export function clearMockFsOperations(): void {
  operations.length = 0;
}

/**
 * Gets operations filtered by operation name
 *
 * @param operationName - Name of operation to filter by
 * @returns Array of operations matching the name
 *
 * @example
 * const readOps = getOperationsByName('readFile');
 */
export function getOperationsByName(operationName: string): FsOperation[] {
  return operations.filter((op) => op.operation === operationName);
}

/**
 * Gets the most recent operation
 *
 * @returns Most recent operation or undefined if no operations
 */
export function getLastOperation(): FsOperation | undefined {
  return operations[operations.length - 1];
}

// Mock file system state (in-memory storage)
interface MockFileSystemState {
  files: Map<string, string>;
  directories: Set<string>;
}

const mockState: MockFileSystemState = {
  files: new Map(),
  directories: new Set(),
};

/**
 * Sets the content of a file in the mock file system
 *
 * @param filePath - Path to the file
 * @param content - File content
 */
export function setMockFileContent(filePath: string, content: string): void {
  mockState.files.set(filePath, content);
}

/**
 * Creates a directory in the mock file system
 *
 * @param dirPath - Path to the directory
 */
export function setMockDirectory(dirPath: string): void {
  mockState.directories.add(dirPath);
}

/**
 * Gets the content of a file from the mock file system
 *
 * @param filePath - Path to the file
 * @returns File content or undefined if not found
 */
export function getMockFileContent(filePath: string): string | undefined {
  return mockState.files.get(filePath);
}

/**
 * Checks if a directory exists in the mock file system
 *
 * @param dirPath - Path to check
 * @returns True if directory exists
 */
export function mockDirectoryExists(dirPath: string): boolean {
  return mockState.directories.has(dirPath);
}

/**
 * Resets the entire mock file system state
 */
export function resetMockFsState(): void {
  mockState.files.clear();
  mockState.directories.clear();
}

// Promise-based fs.promises API mocks
export const mockFs = {
  /**
   * Mocks fs.promises.readFile
   * Automatically tracks operation when called
   */
  readFile: vi.fn(async (filePath: string, encoding?: string) => {
    recordOperation('readFile', [filePath, encoding]);

    const content = mockState.files.get(filePath);
    if (content === undefined) {
      const error = new Error(`ENOENT: no such file or directory, open '${filePath}'`) as NodeJS.ErrnoException;
      error.code = 'ENOENT';
      throw error;
    }

    return content;
  }),

  /**
   * Mocks fs.promises.writeFile
   * Automatically tracks operation and updates mock state
   */
  writeFile: vi.fn(async (filePath: string, content: string, encoding?: string) => {
    recordOperation('writeFile', [filePath, content, encoding]);
    mockState.files.set(filePath, content);
  }),

  /**
   * Mocks fs.promises.mkdir
   * Automatically tracks operation and updates mock state
   */
  mkdir: vi.fn(async (dirPath: string, options?: { recursive?: boolean }) => {
    recordOperation('mkdir', [dirPath, options]);
    mockState.directories.add(dirPath);
  }),

  /**
   * Mocks fs.promises.rm
   * Automatically tracks operation and updates mock state
   */
  rm: vi.fn(async (targetPath: string, options?: { recursive?: boolean; force?: boolean }) => {
    recordOperation('rm', [targetPath, options]);

    // Remove file
    if (mockState.files.has(targetPath)) {
      mockState.files.delete(targetPath);
      return;
    }

    // Remove directory
    if (mockState.directories.has(targetPath)) {
      mockState.directories.delete(targetPath);
      return;
    }

    // If not found and force is not set, throw error
    if (!options?.force) {
      const error = new Error(`ENOENT: no such file or directory, rm '${targetPath}'`) as NodeJS.ErrnoException;
      error.code = 'ENOENT';
      throw error;
    }
  }),

  /**
   * Mocks fs.promises.access
   * Checks if file or directory exists in mock state
   */
  access: vi.fn(async (targetPath: string, mode?: number) => {
    recordOperation('access', [targetPath, mode]);

    const exists = mockState.files.has(targetPath) || mockState.directories.has(targetPath);
    if (!exists) {
      const error = new Error(`ENOENT: no such file or directory, access '${targetPath}'`) as NodeJS.ErrnoException;
      error.code = 'ENOENT';
      throw error;
    }
  }),

  /**
   * Mocks fs.promises.stat
   * Returns stat information for files and directories
   */
  stat: vi.fn(async (targetPath: string) => {
    recordOperation('stat', [targetPath]);

    const isFile = mockState.files.has(targetPath);
    const isDir = mockState.directories.has(targetPath);

    if (!isFile && !isDir) {
      const error = new Error(`ENOENT: no such file or directory, stat '${targetPath}'`) as NodeJS.ErrnoException;
      error.code = 'ENOENT';
      throw error;
    }

    const content = mockState.files.get(targetPath) || '';
    const size = isFile ? content.length : 0;
    const mtime = new Date();

    return {
      size,
      mtime,
      isFile: () => isFile,
      isDirectory: () => isDir,
    };
  }),

  /**
   * Mocks fs.promises.readdir
   * Returns list of files in a directory
   */
  readdir: vi.fn(async (dirPath: string, options?: { withFileTypes?: boolean }) => {
    recordOperation('readdir', [dirPath, options]);

    if (!mockState.directories.has(dirPath)) {
      const error = new Error(`ENOENT: no such file or directory, scandir '${dirPath}'`) as NodeJS.ErrnoException;
      error.code = 'ENOENT';
      throw error;
    }

    // Return files that start with dirPath
    const files: string[] = [];
    const filePathsArray = Array.from(mockState.files.keys());
    for (const filePath of filePathsArray) {
      if (filePath.startsWith(dirPath + '/')) {
        const relativePath = filePath.substring(dirPath.length + 1);
        // Only include immediate children (not nested)
        if (!relativePath.includes('/')) {
          files.push(relativePath);
        }
      }
    }

    return files;
  }),

  /**
   * Mocks fs.promises.copyFile
   * Copies file content in mock state
   */
  copyFile: vi.fn(async (src: string, dest: string, flags?: number) => {
    recordOperation('copyFile', [src, dest, flags]);

    const content = mockState.files.get(src);
    if (content === undefined) {
      const error = new Error(`ENOENT: no such file or directory, copyfile '${src}'`) as NodeJS.ErrnoException;
      error.code = 'ENOENT';
      throw error;
    }

    mockState.files.set(dest, content);
  }),

  /**
   * Mocks fs.promises.unlink
   * Removes a file from mock state
   */
  unlink: vi.fn(async (filePath: string) => {
    recordOperation('unlink', [filePath]);

    if (!mockState.files.has(filePath)) {
      const error = new Error(`ENOENT: no such file or directory, unlink '${filePath}'`) as NodeJS.ErrnoException;
      error.code = 'ENOENT';
      throw error;
    }

    mockState.files.delete(filePath);
  }),

  /**
   * Mocks fs.promises.rmdir
   * Removes a directory from mock state
   */
  rmdir: vi.fn(async (dirPath: string) => {
    recordOperation('rmdir', [dirPath]);

    if (!mockState.directories.has(dirPath)) {
      const error = new Error(`ENOENT: no such file or directory, rmdir '${dirPath}'`) as NodeJS.ErrnoException;
      error.code = 'ENOENT';
      throw error;
    }

    mockState.directories.delete(dirPath);
  }),
};

/**
 * Resets all mock functions and clears operations
 * Use this in beforeEach/afterEach hooks
 */
export function resetMockFs(): void {
  mockFs.readFile.mockClear();
  mockFs.writeFile.mockClear();
  mockFs.mkdir.mockClear();
  mockFs.rm.mockClear();
  mockFs.access.mockClear();
  mockFs.stat.mockClear();
  mockFs.readdir.mockClear();
  mockFs.copyFile.mockClear();
  mockFs.unlink.mockClear();
  mockFs.rmdir.mockClear();
  clearMockFsOperations();
  resetMockFsState();
}

/**
 * Stats object returned by fs.stat
 */
interface MockStats {
  size: number;
  mtime: Date;
  isFile: () => boolean;
  isDirectory: () => boolean;
}

// Export the mock as default for use with vi.mock()
export default mockFs;
