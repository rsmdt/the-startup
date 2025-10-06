import inquirer from 'inquirer';
import ora from 'ora';
import chalk from 'chalk';
import { SpecGenerator } from '../core/spec/SpecGenerator.js';
import type { SpecCommandOptions } from '../core/types/config.js';
import { createSpecGeneratorFS } from './fs-adapter.js';

/**
 * Spec CLI Command
 *
 * Implements interactive spec directory creation with Inquirer prompts.
 * Creates numbered spec directories (001, 002, 003...) with optional templates.
 *
 * Business Rules (PRD lines 201-211):
 * - Rule 1: Auto-increment spec IDs based on existing directories
 * - Rule 2: Create directory with format: [id]-[name]
 * - Rule 3: Support --add flag to generate template file (PRD, SDD, PLAN, BRD)
 * - Rule 4: Support --read flag to output spec state in TOML format
 * - Rule 5: If name argument provided, create spec directly
 *
 * @param name - Optional spec name (creates spec directly if provided)
 * @param options - CLI command options
 *
 * @example
 * ```bash
 * the-agentic-startup spec
 * the-agentic-startup spec user-authentication
 * the-agentic-startup spec user-authentication --add PRD
 * the-agentic-startup spec --read 001
 * ```
 */
export async function specCommand(
  name?: string,
  options?: SpecCommandOptions
): Promise<void> {
  try {
    const fsAdapter = createSpecGeneratorFS();
    const generator = new SpecGenerator(fsAdapter, 'docs/specs');

    // Handle --read flag
    if (options?.read) {
      await handleReadSpec(generator, name);
      return;
    }

    // Handle --add flag or direct creation
    if (name) {
      await handleCreateSpec(generator, name, options?.add);
    } else {
      await handleInteractiveSpec(generator, options?.add);
    }
  } catch (error) {
    console.log(chalk.red.bold('\nError during spec operation\n'));
    console.log(
      chalk.red(error instanceof Error ? error.message : 'Unknown error')
    );
    process.exit(1);
  }
}

/**
 * Handle --read flag: Output spec state in TOML format
 */
async function handleReadSpec(
  generator: SpecGenerator,
  id?: string
): Promise<void> {
  if (!id) {
    console.log(chalk.red('\nError: Spec ID required for --read flag\n'));
    console.log(chalk.gray('Usage: the-agentic-startup spec --read 001\n'));
    process.exit(1);
  }

  const spinner = ora('Reading spec...').start();
  const result = await generator.readSpec(id);
  spinner.stop();

  if (result.success && result.toml) {
    console.log(chalk.blue.bold('\nSpec metadata:\n'));
    console.log(result.toml);
    console.log(); // Empty line
    process.exit(0);
  } else {
    console.log(chalk.red.bold('\nFailed to read spec\n'));
    console.log(chalk.red(result.error || 'Unknown error'));
    process.exit(1);
  }
}

/**
 * Handle direct spec creation (name provided as argument)
 */
async function handleCreateSpec(
  generator: SpecGenerator,
  name: string,
  template?: 'PRD' | 'SDD' | 'PLAN' | 'BRD'
): Promise<void> {
  console.log(chalk.blue.bold('\nCreating spec directory'));
  console.log(chalk.gray(`Name: ${name}\n`));

  const spinner = ora('Creating spec...').start();

  const result = await generator.createSpec({
    name,
    template,
  });

  spinner.stop();

  if (result.success) {
    console.log(chalk.green.bold('\nSpec created successfully!\n'));
    console.log(chalk.blue('Details:'));
    console.log(chalk.gray(`  ID: ${result.specId}`));
    console.log(chalk.gray(`  Directory: ${result.directory}`));

    if (result.templateGenerated) {
      console.log(chalk.gray(`  Template: ${result.templateGenerated}`));
    }

    console.log(); // Empty line
    process.exit(0);
  } else {
    console.log(chalk.red.bold('\nFailed to create spec\n'));
    console.log(chalk.red(result.error || 'Unknown error'));
    process.exit(1);
  }
}

/**
 * Handle interactive spec creation (no name argument)
 */
async function handleInteractiveSpec(
  generator: SpecGenerator,
  template?: 'PRD' | 'SDD' | 'PLAN' | 'BRD'
): Promise<void> {
  console.log(chalk.blue.bold('\nSpec Directory Creation'));
  console.log(chalk.gray('Create a numbered spec directory\n'));

  // Prompt for spec name
  const { specName } = await inquirer.prompt<{ specName: string }>([
    {
      type: 'input',
      name: 'specName',
      message: 'Spec name (will be sanitized for directory):',
      validate: (input) => {
        if (input.trim().length === 0) {
          return 'Spec name cannot be empty';
        }
        return true;
      },
    },
  ]);

  // Prompt for template if not provided via flag
  let selectedTemplate = template;

  if (!selectedTemplate) {
    const { wantsTemplate } = await inquirer.prompt<{ wantsTemplate: boolean }>(
      [
        {
          type: 'confirm',
          name: 'wantsTemplate',
          message: 'Generate a template file?',
          default: false,
        },
      ]
    );

    if (wantsTemplate) {
      const { templateChoice } = await inquirer.prompt<{
        templateChoice: 'PRD' | 'SDD' | 'PLAN' | 'BRD';
      }>([
        {
          type: 'list',
          name: 'templateChoice',
          message: 'Select template:',
          choices: [
            { name: 'Product Requirements Document (PRD)', value: 'PRD' },
            { name: 'System Design Document (SDD)', value: 'SDD' },
            { name: 'Implementation Plan (PLAN)', value: 'PLAN' },
            { name: 'Business Requirements Document (BRD)', value: 'BRD' },
          ],
        },
      ]);

      selectedTemplate = templateChoice;
    }
  }

  // Create spec
  const spinner = ora('Creating spec...').start();

  const result = await generator.createSpec({
    name: specName,
    template: selectedTemplate,
  });

  spinner.stop();

  if (result.success) {
    console.log(chalk.green.bold('\nSpec created successfully!\n'));
    console.log(chalk.blue('Details:'));
    console.log(chalk.gray(`  ID: ${result.specId}`));
    console.log(chalk.gray(`  Directory: ${result.directory}`));

    if (result.templateGenerated) {
      console.log(chalk.gray(`  Template: ${result.templateGenerated}`));
    }

    console.log(); // Empty line
    process.exit(0);
  } else {
    console.log(chalk.red.bold('\nFailed to create spec\n'));
    console.log(chalk.red(result.error || 'Unknown error'));
    process.exit(1);
  }
}
