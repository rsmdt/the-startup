import { defineConfig } from 'tsup';

export default defineConfig({
  // Entry point
  entry: ['src/index.ts'],

  // Output formats (ESM only - CLI uses top-level await)
  format: ['esm'],

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

  // Platform
  platform: 'node',

  // Treeshake
  treeshake: true,

  // Bundle our code but keep dependencies external
  bundle: true,

  // Mark all node_modules as external
  noExternal: [],

  // Preserve binary permissions
  onSuccess: 'chmod +x dist/index.js'
});
