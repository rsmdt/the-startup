import { describe, it, expect, beforeEach, vi } from 'vitest';
import { Initializer } from '../../../src/core/init/Initializer';
import type { InitOptions, InitResult } from '../../../src/core/types/config';

/**
 * Test suite for Initializer core
 *
 * Tests follow SDD requirements:
 * - Init interfaces (lines 614-629)
 * - Init flow (lines 822-868)
 * - PRD acceptance criteria (lines 192-199)
 *
 * Coverage requirement: 90%+ (SDD line 1241)
 */

describe('Initializer', () => {
  let initializer: Initializer;
  let mockFs: any;
  let mockAssetProvider: any;

  const defaultOptions: InitOptions = {
    templates: ['DOR', 'DOD', 'TASK-DOD'],
    targetDirectory: '/test/docs',
  };

  beforeEach(() => {
    // Mock file system
    mockFs = {
      mkdir: vi.fn().mockResolvedValue(undefined),
      copyFile: vi.fn().mockResolvedValue(undefined),
      readFile: vi.fn().mockResolvedValue('mock template content'),
      writeFile: vi.fn().mockResolvedValue(undefined),
      access: vi.fn().mockRejectedValue(new Error('File not found')),
      stat: vi.fn().mockResolvedValue({ isDirectory: () => true }),
    };

    // Mock asset provider
    mockAssetProvider = {
      getTemplateContent: vi.fn((template: string) => {
        const templates: Record<string, string> = {
          'DOR': '# Definition of Ready\n\nContent with {{PROJECT_NAME}}',
          'DOD': '# Definition of Done\n\nContent with {{PROJECT_NAME}}',
          'TASK-DOD': '# Task Definition of Done\n\nContent with {{BUILD_COMMAND}}',
        };
        return templates[template] || '';
      }),
      getTemplateFileName: vi.fn((template: string) => {
        const fileNames: Record<string, string> = {
          'DOR': 'definition-of-ready.md',
          'DOD': 'definition-of-done.md',
          'TASK-DOD': 'task-definition-of-done.md',
        };
        return fileNames[template] || '';
      }),
    };

    initializer = new Initializer(mockFs, mockAssetProvider);
  });

  describe('initialize', () => {
    it('should copy templates to target directory', async () => {
      const result = await initializer.initialize(defaultOptions);

      expect(result.success).toBe(true);
      expect(result.filesCreated).toHaveLength(3);
      expect(result.filesCreated).toContain('/test/docs/definition-of-ready.md');
      expect(result.filesCreated).toContain('/test/docs/definition-of-done.md');
      expect(result.filesCreated).toContain('/test/docs/task-definition-of-done.md');

      expect(mockFs.mkdir).toHaveBeenCalledWith('/test/docs', { recursive: true });
      expect(mockFs.writeFile).toHaveBeenCalledTimes(3);
    });

    it('should replace custom values in template content', async () => {
      const options: InitOptions = {
        templates: ['DOR'],
        targetDirectory: '/test/docs',
        customValues: {
          PROJECT_NAME: 'My Awesome Project',
        },
      };

      await initializer.initialize(options);

      const writeCall = mockFs.writeFile.mock.calls[0];
      expect(writeCall[1]).toContain('My Awesome Project');
      expect(writeCall[1]).not.toContain('{{PROJECT_NAME}}');
    });

    it('should preview files in dry-run mode without creating them', async () => {
      const options: InitOptions = {
        ...defaultOptions,
        dryRun: true,
      };

      const result = await initializer.initialize(options);

      expect(result.success).toBe(true);
      expect(result.filesPreview).toHaveLength(3);
      expect(result.filesPreview).toContain('/test/docs/definition-of-ready.md');
      expect(result.filesCreated).toBeUndefined();

      expect(mockFs.writeFile).not.toHaveBeenCalled();
      expect(mockFs.mkdir).toHaveBeenCalledWith('/test/docs', { recursive: true });
    });

    it('should skip existing files when force is false', async () => {
      // Simulate existing file
      mockFs.access = vi.fn()
        .mockResolvedValueOnce(undefined) // DOR exists
        .mockRejectedValueOnce(new Error('Not found')) // DOD doesn't exist
        .mockRejectedValueOnce(new Error('Not found')); // TASK-DOD doesn't exist

      const result = await initializer.initialize(defaultOptions);

      expect(result.success).toBe(true);
      expect(result.filesCreated).toHaveLength(2);
      expect(result.skipped).toHaveLength(1);
      expect(result.skipped).toContain('/test/docs/definition-of-ready.md');

      expect(mockFs.writeFile).toHaveBeenCalledTimes(2);
    });

    it('should overwrite existing files when force is true', async () => {
      // Simulate existing files
      mockFs.access = vi.fn().mockResolvedValue(undefined);

      const options: InitOptions = {
        ...defaultOptions,
        force: true,
      };

      const result = await initializer.initialize(options);

      expect(result.success).toBe(true);
      expect(result.filesCreated).toHaveLength(3);
      expect(result.skipped).toBeUndefined();

      expect(mockFs.writeFile).toHaveBeenCalledTimes(3);
    });

    it('should create target directory if it does not exist', async () => {
      await initializer.initialize(defaultOptions);

      expect(mockFs.mkdir).toHaveBeenCalledWith('/test/docs', { recursive: true });
    });

    it('should handle errors gracefully', async () => {
      mockFs.mkdir = vi.fn().mockRejectedValue(new Error('Permission denied'));

      const result = await initializer.initialize(defaultOptions);

      expect(result.success).toBe(false);
      expect(result.error).toContain('Permission denied');
      expect(result.filesCreated).toBeUndefined();
    });

    it('should handle empty template list', async () => {
      const options: InitOptions = {
        templates: [],
        targetDirectory: '/test/docs',
      };

      const result = await initializer.initialize(options);

      expect(result.success).toBe(true);
      expect(result.filesCreated).toHaveLength(0);
      expect(mockFs.writeFile).not.toHaveBeenCalled();
    });

    it('should use default target directory when not specified', async () => {
      const options: InitOptions = {
        templates: ['DOR'],
      };

      // Need to provide current working directory context
      const initializerWithCwd = new Initializer(mockFs, mockAssetProvider, '/test/project');

      await initializerWithCwd.initialize(options);

      expect(mockFs.mkdir).toHaveBeenCalledWith('/test/project/docs', { recursive: true });
    });

    it('should replace multiple placeholders in same template', async () => {
      mockAssetProvider.getTemplateContent = vi.fn(() =>
        'Project: {{PROJECT_NAME}}, Build: {{BUILD_COMMAND}}, Test: {{TEST_COMMAND}}'
      );

      const options: InitOptions = {
        templates: ['DOR'],
        targetDirectory: '/test/docs',
        customValues: {
          PROJECT_NAME: 'MyProject',
          BUILD_COMMAND: 'npm run build',
          TEST_COMMAND: 'npm test',
        },
      };

      await initializer.initialize(options);

      const writeCall = mockFs.writeFile.mock.calls[0];
      expect(writeCall[1]).toBe('Project: MyProject, Build: npm run build, Test: npm test');
    });
  });
});
