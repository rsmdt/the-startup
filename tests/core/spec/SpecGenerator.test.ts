import { describe, it, expect, beforeEach, vi } from 'vitest';
import { SpecGenerator } from '../../../src/core/spec/SpecGenerator';
import type { SpecOptions, SpecResult } from '../../../src/core/types/config';

/**
 * Test suite for SpecGenerator core
 *
 * Tests follow SDD requirements:
 * - Spec Flow (lines 870-895)
 * - Spec interfaces (lines 631-653)
 * - TOML output format (PRD line 206)
 *
 * Coverage requirement: 90%+ (SDD line 1241)
 */

describe('SpecGenerator', () => {
  let generator: SpecGenerator;
  let mockFs: any;

  beforeEach(() => {
    // Mock file system (follows Installer pattern)
    mockFs = {
      mkdir: vi.fn().mockResolvedValue(undefined),
      readdir: vi.fn().mockResolvedValue([]),
      writeFile: vi.fn().mockResolvedValue(undefined),
      readFile: vi.fn().mockImplementation((path: string) => {
        // Mock reading template files
        if (path.includes('product-requirements.md')) {
          return Promise.resolve('# Product Requirements\n\nFull template content...');
        }
        if (path.includes('solution-design.md')) {
          return Promise.resolve('# Solution Design\n\nFull template content...');
        }
        if (path.includes('implementation-plan.md')) {
          return Promise.resolve('# Implementation Plan\n\nFull template content...');
        }
        if (path.includes('business-requirements.md')) {
          return Promise.resolve('# Business Requirements\n\nFull template content...');
        }
        return Promise.resolve('');
      }),
      stat: vi.fn().mockResolvedValue({ isDirectory: () => true }),
    };

    generator = new SpecGenerator(mockFs, 'docs/specs', 'assets/the-startup/templates');
  });

  describe('createSpec', () => {
    it('should create spec directory with auto-incremented ID starting at 001', async () => {
      const options: SpecOptions = {
        name: 'user-authentication',
      };

      const result = await generator.createSpec(options);

      expect(result.success).toBe(true);
      expect(result.specId).toBe('001');
      expect(result.directory).toBe('docs/specs/001-user-authentication');
      expect(mockFs.mkdir).toHaveBeenCalledWith(
        'docs/specs/001-user-authentication',
        { recursive: true }
      );
    });

    it('should increment spec ID based on existing directories', async () => {
      // Mock existing spec directories
      mockFs.readdir.mockResolvedValue([
        '001-first-feature',
        '002-second-feature',
        '003-third-feature',
      ]);

      const options: SpecOptions = {
        name: 'fourth-feature',
      };

      const result = await generator.createSpec(options);

      expect(result.success).toBe(true);
      expect(result.specId).toBe('004');
      expect(result.directory).toBe('docs/specs/004-fourth-feature');
    });

    it('should handle non-sequential spec IDs correctly', async () => {
      // Mock spec directories with gaps (001, 003, 005)
      mockFs.readdir.mockResolvedValue([
        '001-feature-one',
        '003-feature-three',
        '005-feature-five',
      ]);

      const options: SpecOptions = {
        name: 'next-feature',
      };

      const result = await generator.createSpec(options);

      // Should use highest + 1 (005 + 1 = 006)
      expect(result.success).toBe(true);
      expect(result.specId).toBe('006');
      expect(result.directory).toBe('docs/specs/006-next-feature');
    });

    it('should ignore invalid directory names when calculating next ID', async () => {
      // Mock mix of valid and invalid directory names
      mockFs.readdir.mockResolvedValue([
        '001-valid-spec',
        'invalid-no-number',
        '.hidden-dir',
        '002-another-valid',
        'README.md',
      ]);

      const options: SpecOptions = {
        name: 'new-feature',
      };

      const result = await generator.createSpec(options);

      expect(result.success).toBe(true);
      expect(result.specId).toBe('003');
    });

    it('should sanitize feature names for directory creation', async () => {
      const options: SpecOptions = {
        name: 'User Authentication & Authorization!',
      };

      const result = await generator.createSpec(options);

      expect(result.success).toBe(true);
      expect(result.specId).toBe('001');
      // Should convert to lowercase and replace spaces/special chars with hyphens
      expect(result.directory).toBe('docs/specs/001-user-authentication-authorization');
    });

    it('should return error when directory creation fails', async () => {
      mockFs.mkdir.mockRejectedValue(new Error('Permission denied'));

      const options: SpecOptions = {
        name: 'test-feature',
      };

      const result = await generator.createSpec(options);

      expect(result.success).toBe(false);
      expect(result.error).toContain('Permission denied');
    });
  });

  describe('readSpec', () => {
    it('should return spec metadata with TOML format', async () => {
      // Mock spec directory exists
      mockFs.readdir.mockResolvedValue(['001-user-authentication']);

      // Mock files in spec directory
      mockFs.readdir.mockImplementation((path: string) => {
        if (path === 'docs/specs') {
          return Promise.resolve(['001-user-authentication']);
        }
        if (path === 'docs/specs/001-user-authentication') {
          return Promise.resolve(['product-requirements.md', 'solution-design.md', 'implementation-plan.md']);
        }
        return Promise.resolve([]);
      });

      const result = await generator.readSpec('001');

      expect(result.success).toBe(true);
      expect(result.specId).toBe('001');
      expect(result.directory).toBe('docs/specs/001-user-authentication');

      // Verify TOML format (PRD line 206)
      expect(result.toml).toBeDefined();
      expect(result.toml).toContain('id = "001"');
      expect(result.toml).toContain('name = "user-authentication"');
      expect(result.toml).toContain('dir = "docs/specs/001-user-authentication"');
      expect(result.toml).toContain('files = ["product-requirements.md", "solution-design.md", "implementation-plan.md"]');
    });

    it('should return error when spec ID does not exist', async () => {
      mockFs.readdir.mockResolvedValue(['001-first-feature']);

      const result = await generator.readSpec('999');

      expect(result.success).toBe(false);
      expect(result.error).toContain('Spec 999 not found');
    });

    it('should handle empty spec directory (no files)', async () => {
      mockFs.readdir.mockImplementation((path: string) => {
        if (path === 'docs/specs') {
          return Promise.resolve(['001-empty-spec']);
        }
        if (path === 'docs/specs/001-empty-spec') {
          return Promise.resolve([]);
        }
        return Promise.resolve([]);
      });

      const result = await generator.readSpec('001');

      expect(result.success).toBe(true);
      expect(result.toml).toContain('files = []');
    });
  });

  describe('createSpec with template', () => {
    it('should create spec directory and copy product-requirements template', async () => {
      const options: SpecOptions = {
        name: 'api-integration',
        template: 'product-requirements',
      };

      const result = await generator.createSpec(options);

      expect(result.success).toBe(true);
      expect(result.specId).toBe('001');
      expect(result.templateGenerated).toBe('docs/specs/001-api-integration/product-requirements.md');

      // Verify readFile was called to read source template
      expect(mockFs.readFile).toHaveBeenCalledWith(
        'assets/the-startup/templates/product-requirements.md',
        'utf-8'
      );

      // Verify writeFile was called to write to destination
      expect(mockFs.writeFile).toHaveBeenCalledWith(
        'docs/specs/001-api-integration/product-requirements.md',
        expect.stringContaining('# Product Requirements'),
        'utf-8'
      );
    });

    it('should create spec directory and copy solution-design template', async () => {
      const options: SpecOptions = {
        name: 'database-migration',
        template: 'solution-design',
      };

      const result = await generator.createSpec(options);

      expect(result.success).toBe(true);
      expect(result.templateGenerated).toBe('docs/specs/001-database-migration/solution-design.md');
      expect(mockFs.readFile).toHaveBeenCalledWith(
        'assets/the-startup/templates/solution-design.md',
        'utf-8'
      );
      expect(mockFs.writeFile).toHaveBeenCalledWith(
        'docs/specs/001-database-migration/solution-design.md',
        expect.stringContaining('# Solution Design'),
        'utf-8'
      );
    });

    it('should create spec directory and copy implementation-plan template', async () => {
      const options: SpecOptions = {
        name: 'ui-redesign',
        template: 'implementation-plan',
      };

      const result = await generator.createSpec(options);

      expect(result.success).toBe(true);
      expect(result.templateGenerated).toBe('docs/specs/001-ui-redesign/implementation-plan.md');
      expect(mockFs.readFile).toHaveBeenCalledWith(
        'assets/the-startup/templates/implementation-plan.md',
        'utf-8'
      );
      expect(mockFs.writeFile).toHaveBeenCalledWith(
        'docs/specs/001-ui-redesign/implementation-plan.md',
        expect.stringContaining('# Implementation Plan'),
        'utf-8'
      );
    });

    it('should create spec directory and copy business-requirements template', async () => {
      const options: SpecOptions = {
        name: 'business-requirements',
        template: 'business-requirements',
      };

      const result = await generator.createSpec(options);

      expect(result.success).toBe(true);
      expect(result.templateGenerated).toBe('docs/specs/001-business-requirements/business-requirements.md');
      expect(mockFs.readFile).toHaveBeenCalledWith(
        'assets/the-startup/templates/business-requirements.md',
        'utf-8'
      );
      expect(mockFs.writeFile).toHaveBeenCalledWith(
        'docs/specs/001-business-requirements/business-requirements.md',
        expect.stringContaining('# Business Requirements'),
        'utf-8'
      );
    });
  });

  describe('parseSpecId', () => {
    it('should extract number from valid spec directory names', () => {
      expect(generator.parseSpecId('001-feature-name')).toBe(1);
      expect(generator.parseSpecId('042-another-feature')).toBe(42);
      expect(generator.parseSpecId('999-last-spec')).toBe(999);
    });

    it('should return null for invalid directory names', () => {
      expect(generator.parseSpecId('invalid-no-number')).toBeNull();
      expect(generator.parseSpecId('.hidden-dir')).toBeNull();
      expect(generator.parseSpecId('README.md')).toBeNull();
      expect(generator.parseSpecId('')).toBeNull();
    });

    it('should handle edge cases', () => {
      expect(generator.parseSpecId('000-zero')).toBe(0);
      expect(generator.parseSpecId('001')).toBe(1); // No feature name
      expect(generator.parseSpecId('123-')).toBe(123); // Trailing hyphen
    });
  });

  describe('getNextSpecId', () => {
    it('should return "001" when no specs exist', async () => {
      mockFs.readdir.mockResolvedValue([]);

      const nextId = await generator.getNextSpecId();

      expect(nextId).toBe('001');
    });

    it('should return next ID with zero-padding', async () => {
      mockFs.readdir.mockResolvedValue([
        '001-first',
        '002-second',
        '003-third',
      ]);

      const nextId = await generator.getNextSpecId();

      expect(nextId).toBe('004');
    });

    it('should handle large spec numbers correctly', async () => {
      mockFs.readdir.mockResolvedValue([
        '099-ninety-nine',
      ]);

      const nextId = await generator.getNextSpecId();

      expect(nextId).toBe('100');
    });

    it('should return "001" when directory read fails', async () => {
      mockFs.readdir.mockRejectedValue(new Error('Directory not found'));

      const nextId = await generator.getNextSpecId();

      expect(nextId).toBe('001');
    });
  });

  describe('error handling', () => {
    it('should handle readSpec errors gracefully', async () => {
      mockFs.readdir.mockRejectedValue(new Error('Permission denied'));

      const result = await generator.readSpec('001');

      expect(result.success).toBe(false);
      expect(result.error).toContain('Permission denied');
    });
  });
});
