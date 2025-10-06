import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    // Test environment
    environment: 'node',

    // Include test files
    include: ['tests/**/*.test.ts', 'tests/**/*.test.tsx'],

    // Coverage configuration
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html', 'json', 'lcov'],
      reportsDirectory: './coverage',

      // Coverage thresholds (90% requirement from SDD line 1241)
      thresholds: {
        lines: 90,
        functions: 90,
        branches: 90,
        statements: 90
      },

      // Include only source files
      include: ['src/**/*.ts', 'src/**/*.tsx'],

      // Exclude test files and generated files
      exclude: [
        'node_modules',
        'dist',
        'tests',
        '**/*.test.ts',
        '**/*.test.tsx',
        '**/types/**',
        'src/index.ts'
      ]
    },

    // Global test timeout
    testTimeout: 10000,

    // Hooks timeout
    hookTimeout: 10000,

    // Globals (for describe, it, expect)
    globals: true,

    // Clear mocks between tests
    clearMocks: true,

    // Reset mocks between tests
    resetMocks: true,

    // Restore mocks between tests
    restoreMocks: true
  }
});
