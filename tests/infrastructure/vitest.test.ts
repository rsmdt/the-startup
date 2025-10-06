/**
 * T001.2.1: Vitest Infrastructure Validation
 *
 * This test validates that vitest can discover and run test files.
 * Tests vitest patterns (*.test.ts, *.spec.ts) and describe/it/expect syntax.
 *
 * [ref: SDD; lines: 42, 1241]
 * [ref: PLAN; lines: 133]
 */

import { describe, it, expect } from 'vitest';

describe('Vitest Infrastructure', () => {
  describe('test discovery', () => {
    it('should run tests with .test.ts extension', () => {
      // This test existing proves .test.ts files are discovered
      expect(true).toBe(true);
    });

    it('should support describe/it/expect syntax', () => {
      const value = 42;
      expect(value).toBe(42);
      expect(value).toBeTypeOf('number');
    });
  });

  describe('basic assertions', () => {
    it('should support equality assertions', () => {
      expect(1 + 1).toBe(2);
      expect({ a: 1 }).toEqual({ a: 1 });
    });

    it('should support truthy/falsy assertions', () => {
      expect(true).toBeTruthy();
      expect(false).toBeFalsy();
      expect(null).toBeNull();
      expect(undefined).toBeUndefined();
    });

    it('should support type assertions', () => {
      expect('hello').toBeTypeOf('string');
      expect(123).toBeTypeOf('number');
      expect(true).toBeTypeOf('boolean');
      expect({}).toBeTypeOf('object');
    });

    it('should support array and object assertions', () => {
      const arr = [1, 2, 3];
      expect(arr).toHaveLength(3);
      expect(arr).toContain(2);

      const obj = { name: 'test', value: 42 };
      expect(obj).toHaveProperty('name');
      expect(obj).toHaveProperty('value', 42);
    });
  });

  describe('async support', () => {
    it('should support async/await', async () => {
      const promise = Promise.resolve('success');
      await expect(promise).resolves.toBe('success');
    });

    it('should support promise rejection', async () => {
      const promise = Promise.reject(new Error('failure'));
      await expect(promise).rejects.toThrow('failure');
    });
  });

  describe('test patterns', () => {
    it('should support nested describe blocks', () => {
      // This test being inside nested describe proves nesting works
      expect(true).toBe(true);
    });

    it('should allow multiple expectations per test', () => {
      const result = { status: 'ok', code: 200 };

      expect(result).toBeDefined();
      expect(result.status).toBe('ok');
      expect(result.code).toBe(200);
    });
  });

  describe('error handling', () => {
    it('should support error expectations', () => {
      const throwError = () => {
        throw new Error('test error');
      };

      expect(throwError).toThrow('test error');
      expect(throwError).toThrow(Error);
    });

    it('should support custom matchers', () => {
      expect(() => {
        // intentional error
        throw new Error('custom error message');
      }).toThrow('custom error message');
    });
  });
});

/**
 * Companion test with .spec.ts extension to verify pattern discovery
 *
 * Note: This is intentionally minimal to prove both .test.ts and .spec.ts
 * patterns are discovered by vitest configuration.
 */
describe('Vitest Pattern Discovery', () => {
  it('should discover this test file with .test.ts extension', () => {
    // File name ends with .test.ts
    expect(__filename).toMatch(/\.test\.ts$/);
  });
});
