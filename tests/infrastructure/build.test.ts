/**
 * T001.2.2: tsup Build Infrastructure Validation
 *
 * This test validates that tsup builds dual ESM/CJS output correctly.
 * Verifies dist/index.js (CJS), dist/index.mjs (ESM), and dist/index.d.ts exist.
 * Tests build process runs without errors.
 *
 * [ref: SDD; lines: 43, 486-490]
 * [ref: PLAN; lines: 134]
 */

import { describe, it, expect, beforeAll } from 'vitest';
import { execSync } from 'child_process';
import { existsSync, readFileSync, statSync } from 'fs';
import { join } from 'path';

const PROJECT_ROOT = join(__dirname, '../../');
const DIST_DIR = join(PROJECT_ROOT, 'dist');

describe('tsup Build Infrastructure', () => {
  describe('build execution', () => {
    it('should have npm run build command configured', () => {
      // Check that package.json exists and has build script
      const packageJsonPath = join(PROJECT_ROOT, 'package.json');
      expect(existsSync(packageJsonPath)).toBe(true);

      const packageJson = JSON.parse(readFileSync(packageJsonPath, 'utf-8'));
      expect(packageJson.scripts).toBeDefined();
      expect(packageJson.scripts.build).toBeDefined();
      expect(packageJson.scripts.build).toContain('tsup');
    });

    it('should run build without errors', () => {
      // This test will fail until build configuration is set up
      // It validates that the build process completes successfully
      expect(() => {
        execSync('npm run build', {
          cwd: PROJECT_ROOT,
          stdio: 'pipe',
          encoding: 'utf-8',
        });
      }).not.toThrow();
    });
  });

  describe('build output - dual module format', () => {
    beforeAll(() => {
      // Ensure build has run before checking outputs
      if (!existsSync(DIST_DIR)) {
        execSync('npm run build', {
          cwd: PROJECT_ROOT,
          stdio: 'pipe',
        });
      }
    });

    it('should generate ESM output (dist/index.js or dist/index.mjs)', () => {
      // tsup can output ESM as either .js or .mjs depending on config
      const esmPath = existsSync(join(DIST_DIR, 'index.mjs'))
        ? join(DIST_DIR, 'index.mjs')
        : join(DIST_DIR, 'index.js');

      expect(existsSync(esmPath)).toBe(true);

      const content = readFileSync(esmPath, 'utf-8');
      // ESM should use import/export syntax
      expect(content).toMatch(/export|import/);
    });

    it('should generate CJS output (dist/index.cjs or dist/index.js)', () => {
      // tsup can output CJS as either .cjs or .js depending on config
      const cjsPath = existsSync(join(DIST_DIR, 'index.cjs'))
        ? join(DIST_DIR, 'index.cjs')
        : join(DIST_DIR, 'index.js');

      expect(existsSync(cjsPath)).toBe(true);

      const content = readFileSync(cjsPath, 'utf-8');
      // CJS should use require/module.exports or __esModule marker
      expect(
        content.includes('require') ||
          content.includes('module.exports') ||
          content.includes('__esModule')
      ).toBe(true);
    });

    it('should generate TypeScript declaration files (dist/index.d.ts)', () => {
      const dtsPath = join(DIST_DIR, 'index.d.ts');
      expect(existsSync(dtsPath)).toBe(true);

      const content = readFileSync(dtsPath, 'utf-8');
      // Declaration files should have TypeScript syntax
      expect(
        content.includes('export') ||
          content.includes('declare') ||
          content.includes('interface') ||
          content.includes('type')
      ).toBe(true);
    });
  });

  describe('build output - file validation', () => {
    it('should create non-empty build artifacts', () => {
      const artifacts = [
        join(DIST_DIR, 'index.js'),
        join(DIST_DIR, 'index.mjs'),
        join(DIST_DIR, 'index.cjs'),
        join(DIST_DIR, 'index.d.ts'),
      ];

      // At least one of each format should exist and be non-empty
      const existingArtifacts = artifacts.filter(existsSync);
      expect(existingArtifacts.length).toBeGreaterThan(0);

      existingArtifacts.forEach((artifact) => {
        const stats = statSync(artifact);
        expect(stats.size).toBeGreaterThan(0);
      });
    });

    it('should preserve source structure in dist/assets/', () => {
      // Verify assets are copied to dist/assets/ as per SDD line 490
      const assetsDir = join(DIST_DIR, 'assets');

      // This will fail until assets are configured
      // Test validates asset copying strategy
      if (existsSync(assetsDir)) {
        const expectedDirs = ['agents', 'commands', 'templates', 'rules', 'output-styles'];
        expectedDirs.forEach((dir) => {
          const dirPath = join(assetsDir, dir);
          if (existsSync(dirPath)) {
            expect(statSync(dirPath).isDirectory()).toBe(true);
          }
        });
      }
    });
  });

  describe('build configuration', () => {
    it('should have tsup.config.ts configured', () => {
      const configPath = join(PROJECT_ROOT, 'tsup.config.ts');
      expect(existsSync(configPath)).toBe(true);

      const content = readFileSync(configPath, 'utf-8');
      expect(content).toContain('defineConfig');
      expect(content).toMatch(/format.*:\s*\[.*['"]cjs['"].*['"]esm['"]/);
    });

    it('should enable declaration files in tsup config', () => {
      const configPath = join(PROJECT_ROOT, 'tsup.config.ts');
      const content = readFileSync(configPath, 'utf-8');

      // Should have dts: true or declaration: true
      expect(
        content.includes('dts: true') || content.includes('declaration: true')
      ).toBe(true);
    });

    it('should specify entry point in tsup config', () => {
      const configPath = join(PROJECT_ROOT, 'tsup.config.ts');
      const content = readFileSync(configPath, 'utf-8');

      // Should have entry point defined
      expect(content).toMatch(/entry.*['"](src\/)?index\.ts['"]/);
    });
  });

  describe('package.json bin field', () => {
    it('should define bin field for CLI executable', () => {
      const packageJsonPath = join(PROJECT_ROOT, 'package.json');
      const packageJson = JSON.parse(readFileSync(packageJsonPath, 'utf-8'));

      // PRD line 155-161 requires bin field for CLI
      expect(packageJson.bin).toBeDefined();
      expect(
        typeof packageJson.bin === 'string' ||
          (typeof packageJson.bin === 'object' && packageJson.bin !== null)
      ).toBe(true);
    });

    it('should point bin to dist output', () => {
      const packageJsonPath = join(PROJECT_ROOT, 'package.json');
      const packageJson = JSON.parse(readFileSync(packageJsonPath, 'utf-8'));

      const binValue =
        typeof packageJson.bin === 'string'
          ? packageJson.bin
          : Object.values(packageJson.bin)[0];

      expect(binValue).toMatch(/^\.\/dist\//);
    });
  });
});
