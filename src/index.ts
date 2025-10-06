/**
 * Main entry point for the-agentic-startup npm package
 *
 * This file serves as the CLI entry point when the package is invoked
 * via `npx the-agentic-startup` or after global installation.
 *
 * @ref SDD lines 473 - Main entry point
 */

// Re-export types for programmatic usage
export * from './core/types/lock.js';
export * from './core/types/settings.js';
export * from './core/types/config.js';

// Entry point will be implemented in later phases
// For now, this ensures the build pipeline works
console.log('the-agentic-startup - TypeScript foundation established');
