/**
 * FileSystem adapter for CLI commands
 *
 * Adapts Node.js fs.promises API to the FileSystem interfaces expected
 * by Installer, Initializer, and SpecGenerator.
 *
 * This adapter ensures type compatibility between Node's fs API
 * and our simplified FileSystem interfaces.
 */

import { promises as fs } from 'fs';

/**
 * Create FileSystem adapter for Installer
 */
export function createInstallerFS() {
  return {
    mkdir: async (path: string, options: { recursive: boolean }): Promise<void> => {
      await fs.mkdir(path, options);
    },
    copyFile: async (src: string, dest: string): Promise<void> => {
      await fs.copyFile(src, dest);
    },
    readFile: async (path: string, encoding: string): Promise<string> => {
      return await fs.readFile(path, encoding as BufferEncoding);
    },
    writeFile: async (path: string, content: string, encoding: string): Promise<void> => {
      await fs.writeFile(path, content, encoding as BufferEncoding);
    },
    rm: async (path: string, options?: { force?: boolean; recursive?: boolean }): Promise<void> => {
      await fs.rm(path, options);
    },
    access: async (path: string): Promise<void> => {
      await fs.access(path);
    },
    stat: async (path: string) => {
      return await fs.stat(path);
    },
  };
}

/**
 * Create FileSystem adapter for Initializer
 */
export function createInitializerFS() {
  return {
    mkdir: async (path: string, options: { recursive: boolean }): Promise<void> => {
      await fs.mkdir(path, options);
    },
    readFile: async (path: string, encoding: string): Promise<string> => {
      return await fs.readFile(path, encoding as BufferEncoding);
    },
    writeFile: async (path: string, content: string, encoding: string): Promise<void> => {
      await fs.writeFile(path, content, encoding as BufferEncoding);
    },
    access: async (path: string): Promise<void> => {
      await fs.access(path);
    },
    stat: async (path: string) => {
      return await fs.stat(path);
    },
  };
}

/**
 * Create FileSystem adapter for SpecGenerator
 */
export function createSpecGeneratorFS() {
  return {
    mkdir: async (path: string, options: { recursive: boolean }): Promise<void> => {
      await fs.mkdir(path, options);
    },
    readdir: async (path: string): Promise<string[]> => {
      return await fs.readdir(path);
    },
    writeFile: async (path: string, content: string, encoding: string): Promise<void> => {
      await fs.writeFile(path, content, encoding as BufferEncoding);
    },
    readFile: async (path: string, encoding: string): Promise<string> => {
      return await fs.readFile(path, encoding as BufferEncoding);
    },
    stat: async (path: string) => {
      return await fs.stat(path);
    },
  };
}

/**
 * Create FileSystem adapter for SettingsMerger
 */
export function createSettingsMergerFS() {
  return {
    readFile: async (path: string, encoding: string): Promise<string> => {
      return await fs.readFile(path, encoding as BufferEncoding);
    },
    writeFile: async (path: string, content: string, encoding: string): Promise<void> => {
      await fs.writeFile(path, content, encoding as BufferEncoding);
    },
    copyFile: async (src: string, dest: string): Promise<void> => {
      await fs.copyFile(src, dest);
    },
    rm: async (path: string, options?: { force?: boolean }): Promise<void> => {
      await fs.rm(path, options);
    },
    access: async (path: string): Promise<void> => {
      await fs.access(path);
    },
  };
}
