import React from 'react';
import { render } from 'ink';
import chalk from 'chalk';
import { InstallWizard } from '../ui/install/InstallWizard.js';
import { Installer } from '../core/installer/Installer.js';
import { LockManager } from '../core/installer/LockManager.js';
import { SettingsMerger } from '../core/installer/SettingsMerger.js';
import type { InstallCommandOptions, InstallResult } from '../core/types/config.js';
import { homedir } from 'os';
import { resolve, join } from 'path';
import { createInstallerFS, createSettingsMergerFS } from './fs-adapter.js';
import { createAssetProvider } from './asset-provider.js';

/**
 * Install CLI Command
 *
 * Implements interactive and non-interactive installation flows.
 *
 * Business Rules (PRD lines 162-175, SDD lines 677-718):
 * - Rule 1: Interactive mode (default): Launch InstallWizard TUI
 * - Rule 2: Non-interactive mode (--local or --yes): Direct Installer call
 * - Rule 3: --local flag: Use defaults (./.the-startup, ~/.claude)
 * - Rule 4: --yes flag: Auto-confirm all prompts
 * - Rule 5: --force flag: Overwrite existing installation
 *
 * @param options - CLI command options
 *
 * @example
 * ```bash
 * the-agentic-startup install              # Interactive TUI
 * the-agentic-startup install --local      # Non-interactive, local paths
 * the-agentic-startup install --yes        # Auto-confirm prompts
 * the-agentic-startup install --local --yes # Fully automated
 * ```
 */
export async function installCommand(options: InstallCommandOptions): Promise<void> {
  try {
    // Check for non-interactive mode (--local or --yes without TUI)
    const isNonInteractive = options.local && options.yes;

    if (isNonInteractive) {
      await handleNonInteractiveInstall(options);
    } else {
      await handleInteractiveInstall(options);
    }
  } catch (error) {
    console.log(chalk.red.bold('\nInstallation error\n'));
    console.log(
      chalk.red(error instanceof Error ? error.message : 'Unknown error')
    );
    process.exit(1);
  }
}

/**
 * Handle non-interactive installation (direct Installer call)
 */
async function handleNonInteractiveInstall(
  _options: InstallCommandOptions
): Promise<void> {
  console.log(chalk.blue.bold('\nNon-interactive installation'));
  console.log(chalk.gray('Using default paths and settings\n'));

  // Use default paths
  const startupPath = resolve(process.cwd(), '.the-startup');
  const claudePath = resolve(homedir(), '.claude');

  console.log(chalk.gray(`Installation directory: ${startupPath}`));
  console.log(chalk.gray(`Claude config directory: ${claudePath}\n`));

  // Create installer components
  const lockFilePath = join(startupPath, '.the-startup.lock');
  const lockManager = new LockManager(lockFilePath);
  const settingsMergerFS = createSettingsMergerFS();
  const settingsMerger = new SettingsMerger(settingsMergerFS);
  const assetProvider = createAssetProvider();
  const installerFS = createInstallerFS();
  const installer = new Installer(
    installerFS,
    lockManager,
    settingsMerger,
    assetProvider,
    '1.0.0'
  );

  // Execute installation
  console.log(chalk.blue('Installing...\n'));

  const result = await installer.install({
    startupPath,
    claudePath,
    selectedFiles: {
      agents: true,
      commands: true,
      templates: true,
      rules: true,
      outputStyles: true,
    },
  });

  // Display results
  if (result.success) {
    console.log(chalk.green.bold('\nInstallation complete!\n'));
    console.log(chalk.blue(`Installed ${result.installedFiles.length} files\n`));
    process.exit(0);
  } else {
    console.log(chalk.red.bold('\nInstallation failed\n'));
    if (result.errors && result.errors.length > 0) {
      result.errors.forEach((error) => {
        console.log(chalk.red(`  ${error}`));
      });
    }
    console.log(); // Empty line
    process.exit(1);
  }
}

/**
 * Handle interactive installation (InstallWizard TUI)
 */
async function handleInteractiveInstall(
  options: InstallCommandOptions
): Promise<void> {
  // Create installer components
  const startupPath = resolve(process.cwd(), '.the-startup');
  const lockFilePath = join(startupPath, '.the-startup.lock');
  const lockManager = new LockManager(lockFilePath);
  const settingsMergerFS = createSettingsMergerFS();
  const settingsMerger = new SettingsMerger(settingsMergerFS);
  const assetProvider = createAssetProvider();
  const installerFS = createInstallerFS();
  const installer = new Installer(
    installerFS,
    lockManager,
    settingsMerger,
    assetProvider,
    '1.0.0'
  );

  // Render InstallWizard
  const { waitUntilExit } = render(
    React.createElement(InstallWizard, {
      options,
      installer,
      onComplete: (result: InstallResult) => {
        // onComplete callback is handled by the wizard
        // Exit code is set based on result
        if (!result.success) {
          process.exitCode = 1;
        }
      },
    })
  );

  // Wait for wizard to complete
  await waitUntilExit();
}
