import { defineConfig } from 'tsup';

export default defineConfig({
  // Entry point
  entry: ['src/index.ts'],

  // Output formats (dual ESM/CJS)
  format: ['esm', 'cjs'],

  // Output directory
  outDir: 'dist',

  // Generate TypeScript declarations
  dts: true,

  // Clean output directory before build
  clean: true,

  // Minify for production
  minify: true,

  // Source maps for debugging
  sourcemap: true,

  // Split output by entry
  splitting: false,

  // Shims for Node.js globals
  shims: true,

  // Target environment
  target: 'node18',

  // Bundle dependencies (exclude peer dependencies)
  noExternal: [
    'chalk',
    'commander',
    'fs-extra',
    'ink',
    'inquirer',
    'ora'
  ],

  // Platform
  platform: 'node',

  // Treeshake
  treeshake: true,

  // Bundle node_modules
  bundle: true,

  // Preserve binary permissions
  onSuccess: 'chmod +x dist/index.js'
});
