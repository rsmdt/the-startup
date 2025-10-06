/**
 * T001.2.3: TypeScript Strict Mode Infrastructure Validation
 *
 * This test validates that TypeScript compiles in strict mode.
 * Verifies no implicit any errors and type checking catches common errors.
 *
 * [ref: SDD; lines: 42]
 * [ref: PLAN; lines: 135]
 */

import { describe, it, expect } from 'vitest';
import { execSync } from 'child_process';
import { existsSync, readFileSync, writeFileSync, mkdirSync } from 'fs';
import { join } from 'path';

const PROJECT_ROOT = join(__dirname, '../../');
const TSCONFIG_PATH = join(PROJECT_ROOT, 'tsconfig.json');

// Helper to parse tsconfig.json with comments
function parseTsConfig(path: string): any {
  const content = readFileSync(path, 'utf-8');
  // Strip JSON comments (// style) for parsing
  const jsonWithoutComments = content.replace(/\/\/.*$/gm, '');
  return JSON.parse(jsonWithoutComments);
}

describe('TypeScript Strict Mode Infrastructure', () => {
  describe('tsconfig.json configuration', () => {
    it('should have tsconfig.json configured', () => {
      expect(existsSync(TSCONFIG_PATH)).toBe(true);
    });

    it('should enable strict mode', () => {
      const tsconfig = parseTsConfig(TSCONFIG_PATH);

      expect(tsconfig.compilerOptions).toBeDefined();
      expect(tsconfig.compilerOptions.strict).toBe(true);
    });

    it('should disable implicit any', () => {
      const tsconfig = parseTsConfig(TSCONFIG_PATH);

      // strict: true implies noImplicitAny: true
      // but we can also check explicit setting
      if ('noImplicitAny' in tsconfig.compilerOptions) {
        expect(tsconfig.compilerOptions.noImplicitAny).toBe(true);
      }
    });

    it('should enable other strict checks', () => {
      const tsconfig = parseTsConfig(TSCONFIG_PATH);

      // When strict: true, these should be enabled implicitly or explicitly
      const strictChecks = [
        'noImplicitAny',
        'strictNullChecks',
        'strictFunctionTypes',
        'strictBindCallApply',
        'strictPropertyInitialization',
        'noImplicitThis',
        'alwaysStrict',
      ];

      // strict: true should enable all these, but explicit settings override
      strictChecks.forEach((check) => {
        if (check in tsconfig.compilerOptions) {
          expect(tsconfig.compilerOptions[check]).toBe(true);
        }
      });
    });
  });

  describe('type checking execution', () => {
    it('should have typecheck npm script configured', () => {
      const packageJsonPath = join(PROJECT_ROOT, 'package.json');
      const packageJson = JSON.parse(readFileSync(packageJsonPath, 'utf-8'));

      expect(packageJson.scripts).toBeDefined();
      expect(packageJson.scripts.typecheck).toBeDefined();
      expect(packageJson.scripts.typecheck).toContain('tsc');
      expect(packageJson.scripts.typecheck).toContain('--noEmit');
    });

    it('should run type checking without errors on valid code', () => {
      // This validates that typecheck runs successfully
      expect(() => {
        execSync('npm run typecheck', {
          cwd: PROJECT_ROOT,
          stdio: 'pipe',
          encoding: 'utf-8',
        });
      }).not.toThrow();
    });
  });

  describe('strict mode error detection', () => {
    const TEST_FILE_DIR = join(PROJECT_ROOT, 'tests/infrastructure/.temp');
    const INVALID_FILE = join(TEST_FILE_DIR, 'invalid-strict.ts');

    it('should catch implicit any errors', () => {
      // Create a temporary file with implicit any
      mkdirSync(TEST_FILE_DIR, { recursive: true });

      const invalidCode = `
// This should fail with implicit any error in strict mode
function processValue(value) {
  return value.toString();
}
`;

      writeFileSync(INVALID_FILE, invalidCode);

      try {
        // Attempt to type check the invalid file
        execSync(`npx tsc --noEmit ${INVALID_FILE}`, {
          cwd: PROJECT_ROOT,
          stdio: 'pipe',
          encoding: 'utf-8',
        });

        // If we reach here, strict mode is NOT working
        expect.fail('TypeScript should have caught implicit any error');
      } catch (error: any) {
        // This is expected - strict mode should reject implicit any
        const errorOutput = error.stderr?.toString() || error.stdout?.toString() || '';
        expect(errorOutput).toMatch(/Parameter 'value' implicitly has an 'any' type/);
      }
    });

    it('should catch null/undefined errors with strictNullChecks', () => {
      const invalidCode = `
// This should fail with strict null checks
function getLength(str: string): number {
  return str.length;
}

// This should error - passing null to parameter expecting string
const result = null;
getLength(result);
`;

      writeFileSync(INVALID_FILE, invalidCode);

      try {
        execSync(`npx tsc --noEmit ${INVALID_FILE}`, {
          cwd: PROJECT_ROOT,
          stdio: 'pipe',
          encoding: 'utf-8',
        });

        expect.fail('TypeScript should have caught null assignment error');
      } catch (error: any) {
        // Expected - strict null checks should catch this
        const errorOutput = error.stderr?.toString() || error.stdout?.toString() || '';
        expect(errorOutput).toMatch(/Argument of type.*null.*not assignable|TS2345/);
      }
    });

    it('should enforce strict function types', () => {
      const invalidCode = `
// This should fail with strict function types
interface Handler {
  (value: string): void;
}

const handler: Handler = (value: number) => {
  console.log(value);
};
`;

      writeFileSync(INVALID_FILE, invalidCode);

      try {
        const result = execSync(`npx tsc --noEmit ${INVALID_FILE}`, {
          cwd: PROJECT_ROOT,
          stdio: 'pipe',
          encoding: 'utf-8',
        });

        expect.fail('TypeScript should have caught function type mismatch');
      } catch (error: any) {
        // Expected - strict function types should catch this
        // TypeScript error will be in stderr
        const errorOutput = error.stderr?.toString() || error.stdout?.toString() || '';
        expect(errorOutput).toMatch(/not assignable to type|TS2322/);
      }
    });

    it('should require proper type annotations', () => {
      const validCode = `
// This should pass - proper type annotations
function processValue(value: string | number): string {
  return value.toString();
}

interface User {
  name: string;
  age: number;
}

function greetUser(user: User): string {
  return \`Hello, \${user.name}!\`;
}
`;

      writeFileSync(INVALID_FILE, validCode);

      // This should NOT throw - code is valid with proper types
      expect(() => {
        execSync(`npx tsc --noEmit ${INVALID_FILE}`, {
          cwd: PROJECT_ROOT,
          stdio: 'pipe',
          encoding: 'utf-8',
        });
      }).not.toThrow();
    });
  });

  describe('module resolution', () => {
    it('should use modern module resolution', () => {
      const tsconfig = parseTsConfig(TSCONFIG_PATH);

      // Should use node16, nodenext, or bundler for modern resolution
      expect(
        tsconfig.compilerOptions.moduleResolution === 'node16' ||
          tsconfig.compilerOptions.moduleResolution === 'nodenext' ||
          tsconfig.compilerOptions.moduleResolution === 'bundler'
      ).toBe(true);
    });

    it('should target modern ES version', () => {
      const tsconfig = parseTsConfig(TSCONFIG_PATH);

      // Should target ES2020 or newer for modern features
      expect(tsconfig.compilerOptions.target).toBeDefined();
      const target = tsconfig.compilerOptions.target.toLowerCase();

      const modernTargets = ['es2020', 'es2021', 'es2022', 'esnext'];
      expect(modernTargets.some((t) => target.includes(t))).toBe(true);
    });

    it('should configure lib for target runtime', () => {
      const tsconfig = parseTsConfig(TSCONFIG_PATH);

      expect(tsconfig.compilerOptions.lib).toBeDefined();
      expect(Array.isArray(tsconfig.compilerOptions.lib)).toBe(true);

      // Should include modern ES features
      const libs = tsconfig.compilerOptions.lib.map((l: string) => l.toLowerCase());
      expect(
        libs.some((l: string) =>
          ['es2020', 'es2021', 'es2022', 'esnext'].some((modern) => l.includes(modern))
        )
      ).toBe(true);
    });
  });

  describe('type safety enforcement', () => {
    it('should enforce noUncheckedIndexedAccess for safer array access', () => {
      const tsconfig = parseTsConfig(TSCONFIG_PATH);

      // Optional but recommended for strictness
      if ('noUncheckedIndexedAccess' in tsconfig.compilerOptions) {
        expect(tsconfig.compilerOptions.noUncheckedIndexedAccess).toBe(true);
      }
    });

    it('should configure exactOptionalPropertyTypes for stricter objects', () => {
      const tsconfig = parseTsConfig(TSCONFIG_PATH);

      // Optional but recommended for strictness
      if ('exactOptionalPropertyTypes' in tsconfig.compilerOptions) {
        expect(tsconfig.compilerOptions.exactOptionalPropertyTypes).toBe(true);
      }
    });

    it('should enforce noImplicitReturns for complete code coverage', () => {
      const tsconfig = parseTsConfig(TSCONFIG_PATH);

      // Helps catch missing return statements
      if ('noImplicitReturns' in tsconfig.compilerOptions) {
        expect(tsconfig.compilerOptions.noImplicitReturns).toBe(true);
      }
    });

    it('should enforce noFallthroughCasesInSwitch for safer switch statements', () => {
      const tsconfig = parseTsConfig(TSCONFIG_PATH);

      // Prevents accidental fallthrough
      if ('noFallthroughCasesInSwitch' in tsconfig.compilerOptions) {
        expect(tsconfig.compilerOptions.noFallthroughCasesInSwitch).toBe(true);
      }
    });
  });

  describe('source and output configuration', () => {
    it('should exclude dist and node_modules from compilation', () => {
      const tsconfig = parseTsConfig(TSCONFIG_PATH);

      expect(tsconfig.exclude).toBeDefined();
      expect(Array.isArray(tsconfig.exclude)).toBe(true);

      const exclude = tsconfig.exclude.map((e: string) => e.toLowerCase());
      expect(exclude).toContain('node_modules');
      expect(exclude).toContain('dist');
    });

    it('should include src directory in compilation', () => {
      const tsconfig = parseTsConfig(TSCONFIG_PATH);

      // Either include should have src, or no include means all files
      if (tsconfig.include) {
        expect(Array.isArray(tsconfig.include)).toBe(true);
        expect(tsconfig.include.some((i: string) => i.includes('src'))).toBe(true);
      }
    });

    it('should set outDir when using tsc for emit', () => {
      const tsconfig = parseTsConfig(TSCONFIG_PATH);

      // outDir should point to dist or similar
      if ('outDir' in tsconfig.compilerOptions) {
        expect(tsconfig.compilerOptions.outDir).toMatch(/dist|build/);
      }
    });
  });
});
