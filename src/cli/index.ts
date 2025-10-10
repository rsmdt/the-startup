import { Command } from 'commander';
import type {
  InstallCommandOptions,
  UninstallCommandOptions,
  InitCommandOptions,
  SpecCommandOptions,
} from '../core/types/config.js';
import { initCommand } from './init.js';
import { specCommand } from './spec.js';
import { installCommand } from './install.js';
import { uninstallCommand } from './uninstall.js';

/**
 * CLI Entry Point
 *
 * Creates and configures Commander.js program with all available commands.
 *
 * Commands:
 * - init: Initialize DOR, DOD, TASK-DOD templates
 * - spec: Create numbered spec directories with optional templates
 * - install: Install components (interactive or non-interactive)
 * - uninstall: Remove installed components
 *
 * @returns Configured Commander.js program
 */
export function createCLI(): Command {
  const program = new Command();

  program
    .name('the-agentic-startup')
    .description('Enterprise-grade AI development agents as frictionless as any npm package')
    .version('1.0.0');

  // Init command (PRD lines 190-200)
  program
    .command('init')
    .description('Initialize DOR, DOD, and TASK-DOD templates')
    .option('--dry-run', 'Preview changes without creating files')
    .option('--force', 'Overwrite existing files without prompting')
    .action(async (options: InitCommandOptions) => {
      await initCommand(options);
    });

  // Spec command (PRD lines 201-211)
  program
    .command('spec [name]')
    .description('Create numbered spec directory with auto-incrementing ID')
    .option('--add <template>', 'Generate template file (product-requirements, solution-design, implementation-plan, business-requirements)')
    .option('--read', 'Output spec state in TOML format (requires name/ID argument)')
    .action(async (name: string | undefined, options: SpecCommandOptions) => {
      await specCommand(name, options);
    });

  // Install command (PRD lines 162-175)
  program
    .command('install')
    .description('Install components interactively or with defaults')
    .option('--local', 'Use local paths (./.the-startup, ~/.claude)')
    .option('--yes', 'Auto-confirm all prompts with recommended settings')
    .action(async (options: InstallCommandOptions) => {
      await installCommand(options);
    });

  // Uninstall command (PRD lines 177-189)
  program
    .command('uninstall')
    .description('Remove installed components')
    .option('--keep-logs', 'Preserve .the-startup/logs directory')
    .option('--keep-settings', "Don't modify ~/.claude/settings.json")
    .action(async (options: UninstallCommandOptions) => {
      await uninstallCommand(options);
    });

  return program;
}

/**
 * Parse and execute CLI commands
 *
 * This function is called from the main entry point.
 * It creates the CLI and parses the provided arguments.
 *
 * @param argv - Command line arguments (default: process.argv)
 */
export async function runCLI(argv: string[] = process.argv): Promise<void> {
  const program = createCLI();
  await program.parseAsync(argv);
}
