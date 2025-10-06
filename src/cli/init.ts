import inquirer from 'inquirer';
import ora from 'ora';
import chalk from 'chalk';
import { Initializer } from '../core/init/Initializer.js';
import type { InitCommandOptions, InitOptions } from '../core/types/config.js';
import { createInitializerFS } from './fs-adapter.js';

/**
 * Asset Provider for embedded templates
 */
class TemplateAssetProvider {
  /**
   * Get template content by name
   * In production, this would read from embedded assets
   * For now, returns minimal template boilerplate
   */
  getTemplateContent(template: string): string {
    switch (template) {
      case 'DOR':
        return `# Definition of Ready (DOR)

## Purpose
Ensure all requirements are clear before work begins.

## Checklist
- [ ] User story is clearly defined
- [ ] Acceptance criteria are documented
- [ ] Dependencies are identified
- [ ] Technical approach is outlined
- [ ] Estimates are provided
- [ ] Design artifacts are available (if applicable)
`;
      case 'DOD':
        return `# Definition of Done (DOD)

## Purpose
Ensure all work meets quality standards before completion.

## Checklist
- [ ] Code is written and peer-reviewed
- [ ] Tests are written and passing
- [ ] Documentation is updated
- [ ] Code is merged to main branch
- [ ] Feature is deployed to staging
- [ ] Acceptance criteria are met
- [ ] Product owner approves
`;
      case 'TASK-DOD':
        return `# Task Definition of Done (TASK-DOD)

## Purpose
Ensure individual tasks meet quality standards.

## Checklist
- [ ] Implementation is complete
- [ ] Unit tests are written and passing
- [ ] Code follows project style guide
- [ ] Comments explain complex logic
- [ ] No linting errors
- [ ] Changes are committed with descriptive message
`;
      default:
        return '# Template\n\n';
    }
  }

  /**
   * Get template filename by name
   */
  getTemplateFileName(template: string): string {
    return `${template.toLowerCase()}.md`;
  }
}

/**
 * Init CLI Command
 *
 * Implements interactive template initialization with Inquirer prompts.
 * Copies DOR, DOD, and TASK-DOD templates to docs/ directory.
 *
 * Business Rules (PRD lines 190-200):
 * - Rule 1: Prompt user for templates to initialize
 * - Rule 2: Allow custom values for template placeholders
 * - Rule 3: Support --dry-run for preview without creating files
 * - Rule 4: Support --force to overwrite existing files
 * - Rule 5: Default target directory is docs/
 *
 * @param options - CLI command options
 *
 * @example
 * ```bash
 * the-agentic-startup init
 * the-agentic-startup init --dry-run
 * the-agentic-startup init --force
 * ```
 */
export async function initCommand(options: InitCommandOptions): Promise<void> {
  try {
    console.log(chalk.blue.bold('\nTemplate Initialization'));
    console.log(chalk.gray('Initialize DOR, DOD, and TASK-DOD templates\n'));

    // Prompt for templates to initialize
    const { selectedTemplates } = await inquirer.prompt<{
      selectedTemplates: ('DOR' | 'DOD' | 'TASK-DOD')[];
    }>([
      {
        type: 'checkbox',
        name: 'selectedTemplates',
        message: 'Select templates to initialize:',
        choices: [
          { name: 'Definition of Ready (DOR)', value: 'DOR', checked: true },
          { name: 'Definition of Done (DOD)', value: 'DOD', checked: true },
          { name: 'Task Definition of Done (TASK-DOD)', value: 'TASK-DOD', checked: true },
        ],
        validate: (answer) => {
          if (answer.length === 0) {
            return 'You must select at least one template.';
          }
          return true;
        },
      },
    ]);

    // Prompt for target directory
    const { targetDirectory } = await inquirer.prompt<{ targetDirectory: string }>([
      {
        type: 'input',
        name: 'targetDirectory',
        message: 'Target directory:',
        default: 'docs',
      },
    ]);

    // Prompt for custom values (optional)
    const { wantsCustomValues } = await inquirer.prompt<{ wantsCustomValues: boolean }>([
      {
        type: 'confirm',
        name: 'wantsCustomValues',
        message: 'Do you want to customize template placeholders?',
        default: false,
      },
    ]);

    let customValues: Record<string, string> | undefined;

    if (wantsCustomValues) {
      const { projectName, teamName } = await inquirer.prompt<{
        projectName: string;
        teamName: string;
      }>([
        {
          type: 'input',
          name: 'projectName',
          message: 'Project name:',
          default: 'My Project',
        },
        {
          type: 'input',
          name: 'teamName',
          message: 'Team name:',
          default: 'Development Team',
        },
      ]);

      customValues = { projectName, teamName };
    }

    // Build init options
    const initOptions: InitOptions = {
      templates: selectedTemplates,
      customValues,
      targetDirectory,
      dryRun: options.dryRun,
      force: options.force,
    };

    // Show dry-run notice
    if (options.dryRun) {
      console.log(chalk.yellow('\nDry-run mode: No files will be created\n'));
    }

    // Show force notice
    if (options.force) {
      console.log(chalk.yellow('\nForce mode: Existing files will be overwritten\n'));
    }

    // Execute initialization
    const spinner = ora('Initializing templates...').start();

    const assetProvider = new TemplateAssetProvider();
    const fsAdapter = createInitializerFS();
    const initializer = new Initializer(fsAdapter, assetProvider);
    const result = await initializer.initialize(initOptions);

    spinner.stop();

    // Display results
    if (result.success) {
      console.log(chalk.green.bold('\nInitialization complete!\n'));

      if (result.filesCreated && result.filesCreated.length > 0) {
        console.log(chalk.blue('Created files:'));
        result.filesCreated.forEach((file) => {
          console.log(chalk.gray(`  - ${file}`));
        });
      }

      if (result.filesPreview && result.filesPreview.length > 0) {
        console.log(chalk.blue('\nFiles that would be created:'));
        result.filesPreview.forEach((file) => {
          console.log(chalk.gray(`  - ${file}`));
        });
      }

      if (result.skipped && result.skipped.length > 0) {
        console.log(chalk.yellow('\nSkipped (already exists):'));
        result.skipped.forEach((file) => {
          console.log(chalk.gray(`  - ${file}`));
        });
        console.log(chalk.gray('\nUse --force to overwrite existing files'));
      }

      console.log(); // Empty line
      process.exit(0);
    } else {
      console.log(chalk.red.bold('\nInitialization failed\n'));
      console.log(chalk.red(result.error || 'Unknown error'));
      process.exit(1);
    }
  } catch (error) {
    console.log(chalk.red.bold('\nError during initialization\n'));
    console.log(
      chalk.red(error instanceof Error ? error.message : 'Unknown error')
    );
    process.exit(1);
  }
}
