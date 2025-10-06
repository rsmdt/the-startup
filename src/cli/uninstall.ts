import React from 'react';
import { render } from 'ink';
import chalk from 'chalk';
import { UninstallWizard } from '../ui/uninstall/UninstallWizard.js';
import type { UninstallResult } from '../ui/uninstall/UninstallWizard.js';
import { LockManager } from '../core/installer/LockManager.js';
import { SettingsMerger } from '../core/installer/SettingsMerger.js';
import type { UninstallCommandOptions } from '../core/types/config.js';
import { resolve, join } from 'path';
import { createSettingsMergerFS } from './fs-adapter.js';

/**
 * Uninstall CLI Command
 *
 * Implements interactive uninstallation with UninstallWizard TUI.
 *
 * Business Rules (PRD lines 177-189):
 * - Rule 1: Launch UninstallWizard TUI
 * - Rule 2: Read lock file to identify installed files
 * - Rule 3: Confirm before deleting
 * - Rule 4: --keep-logs flag: Preserve .the-startup/logs directory
 * - Rule 5: --keep-settings flag: Don't modify settings.json
 *
 * @param options - CLI command options
 *
 * @example
 * ```bash
 * the-agentic-startup uninstall
 * the-agentic-startup uninstall --keep-logs
 * the-agentic-startup uninstall --keep-settings
 * the-agentic-startup uninstall --keep-logs --keep-settings
 * ```
 */
export async function uninstallCommand(
  options: UninstallCommandOptions
): Promise<void> {
  try {
    // Create uninstaller components
    const startupPath = resolve(process.cwd(), '.the-startup');
    const lockFilePath = join(startupPath, '.the-startup.lock');
    const lockManager = new LockManager(lockFilePath);
    const settingsMergerFS = createSettingsMergerFS();
    const settingsMerger = new SettingsMerger(settingsMergerFS);

    // Show header
    console.log(chalk.blue.bold('\nUninstallation Wizard'));
    console.log(chalk.gray('Remove installed components\n'));

    // Render UninstallWizard
    // Note: settingsMerger.removeHooks returns Promise<ClaudeSettings>, but UninstallWizard
    // expects Promise<void>. This is a type mismatch in the wizard that will be fixed in T006.
    // For now, we cast to the expected type.
    const { waitUntilExit } = render(
      React.createElement(UninstallWizard, {
        options,
        lockManager,
        settingsMerger: settingsMerger as any,
        onComplete: (result: UninstallResult) => {
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
  } catch (error) {
    console.log(chalk.red.bold('\nUninstallation error\n'));
    console.log(
      chalk.red(error instanceof Error ? error.message : 'Unknown error')
    );
    process.exit(1);
  }
}
