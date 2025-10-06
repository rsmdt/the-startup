#!/usr/bin/env node

/**
 * Main CLI entry point for the-agentic-startup npm package
 *
 * This file serves as the CLI entry point when the package is invoked
 * via `npx the-agentic-startup` or after global installation.
 *
 * Implements the complete CLI with all commands:
 * - statusline: Cross-platform stdio passthrough
 * - init: Initialize templates (DOR, DOD, TASK-DOD)
 * - spec: Create numbered spec directories
 * - install: Interactive/non-interactive installation
 * - uninstall: Remove installed components
 *
 * @ref SDD lines 473 - Main entry point
 */

import { runCLI } from './cli/index.js';

// Run CLI and handle any uncaught errors
runCLI().catch((error) => {
  console.error('Fatal error:', error instanceof Error ? error.message : error);
  process.exit(1);
});
