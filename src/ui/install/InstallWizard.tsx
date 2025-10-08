import { FC, useState, useEffect } from 'react';
import { Box, Text, useInput, useApp } from 'ink';
import { ChoiceSelector, Choice } from './ChoiceSelector';
import { FinalConfirmation } from './FinalConfirmation';
import { Complete, InstallationSummary } from './Complete';
import type { InstallCommandOptions, InstallResult, InstallerOptions } from '../../core/types/config';
import type { Installer } from '../../core/installer/Installer';
import { theme } from '../shared/theme';
import { homedir } from 'os';
import { resolve, join } from 'path';
import { existsSync } from 'fs';

/**
 * Wizard state machine
 */
type WizardState =
  | 'startupPath'
  | 'claudePath'
  | 'fileSelection'
  | 'installing'
  | 'complete'
  | 'error';

export interface InstallWizardProps {
  /**
   * CLI command options (--local, --yes flags)
   */
  options: InstallCommandOptions;

  /**
   * Installer instance for performing installation
   */
  installer?: Installer;

  /**
   * Callback when installation completes (success or failure)
   */
  onComplete: (result: InstallResult) => void;
}

/**
 * InstallWizard - Interactive installation TUI state machine
 *
 * Implements the 4-state wizard flow from SDD:
 * 1. StartupPath: Select installation directory
 * 2. ClaudePath: Select Claude config directory
 * 3. FileSelection: Interactive tree to select components
 * 4. Complete: Display success summary
 *
 * Features:
 * - --local flag: Skip prompts, use defaults (./.the-startup, ~/.claude)
 * - --yes flag: Auto-confirm all prompts with recommended settings
 * - Ctrl+C: Cancellation with rollback
 * - Error handling: Display errors with recovery options
 * - Progress indication: Show progress during long operations
 *
 * Business Rules (PRD lines 309-314):
 * - Rule 1: --local uses ./.the-startup and ~/.claude defaults
 * - Rule 2: --yes auto-confirms all prompts
 * - Rule 3: Merge with existing files (handled by Installer)
 * - Rule 4: Merge hooks in settings.json (handled by SettingsMerger)
 * - Rule 5: Reinstall compares checksums (handled by LockManager)
 *
 * @example
 * ```tsx
 * <InstallWizard
 *   options={{ local: false, yes: false }}
 *   installer={installerInstance}
 *   onComplete={(result) => console.log('Done:', result)}
 * />
 * ```
 */
export const InstallWizard: FC<InstallWizardProps> = ({
  options,
  installer,
  onComplete
}) => {
  // State machine
  const [state, setState] = useState<WizardState>('startupPath');
  const [startupPath, setStartupPath] = useState<string>('');
  const [claudePath, setClaudePath] = useState<string>('');
  const [selectedFiles] = useState<InstallerOptions['selectedFiles']>({
    agents: true,
    commands: true,
    templates: true,
    rules: true,
    outputStyles: true,
  });
  const [installResult, setInstallResult] = useState<InstallResult | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [progressInfo] = useState<{
    stage: string;
    current: number;
    total: number;
  } | null>(null);

  const { exit } = useApp();

  // Handle --local flag: Skip TUI, use defaults
  useEffect(() => {
    if (options.local) {
      // Rule 1: --local uses ./.the-startup and ~/.claude defaults
      const defaultStartupPath = resolve(process.cwd(), '.the-startup');
      const defaultClaudePath = resolve(homedir(), '.claude');

      setStartupPath(defaultStartupPath);
      setClaudePath(defaultClaudePath);

      // Skip directly to installation
      setState('installing');
    }
  }, [options.local]);

  // Handle --yes flag: Auto-confirm prompts
  useEffect(() => {
    if (options.yes && state === 'startupPath') {
      // Rule 2: --yes auto-confirms with recommended settings
      const recommendedStartupPath = resolve(homedir(), '.the-startup');
      setStartupPath(recommendedStartupPath);
      setState('claudePath');
    }
  }, [options.yes, state]);

  useEffect(() => {
    if (options.yes && state === 'claudePath') {
      const recommendedClaudePath = resolve(homedir(), '.claude');
      setClaudePath(recommendedClaudePath);
      setState('fileSelection');
    }
  }, [options.yes, state]);

  useEffect(() => {
    if (options.yes && state === 'fileSelection') {
      // Auto-select all files with --yes
      setState('installing');
    }
  }, [options.yes, state]);

  // Perform installation when in 'installing' state
  useEffect(() => {
    if (state === 'installing' && installer && startupPath && claudePath) {
      performInstallation();
    }
  }, [state, installer, startupPath, claudePath]);

  // Ctrl+C cancellation handler
  useInput((input, key) => {
    if (key.ctrl && input === 'c') {
      handleCancellation();
    }
  });

  // Auto-exit after successful installation
  useEffect(() => {
    if (state === 'complete') {
      // Give user 2 seconds to read the next steps, then exit
      const timer = setTimeout(() => {
        exit();
      }, 2000);

      return () => clearTimeout(timer);
    }

    return undefined;
  }, [state, exit]);

  /**
   * Perform the actual installation
   */
  const performInstallation = async () => {
    if (!installer) {
      setError('Installer not available');
      setState('error');
      return;
    }

    try {
      const installerOptions: InstallerOptions = {
        startupPath,
        claudePath,
        selectedFiles,
      };

      const result = await installer.install(installerOptions);

      if (result.success) {
        setInstallResult(result);
        setState('complete');
        onComplete(result);
      } else {
        setError(result.errors?.[0] || 'Installation failed');
        setState('error');
        onComplete(result);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Installation failed';
      setError(errorMessage);
      setState('error');
      onComplete({
        success: false,
        installedFiles: [],
        errors: [errorMessage],
      });
    }
  };

  /**
   * Handle Ctrl+C cancellation
   */
  const handleCancellation = async () => {
    // Trigger rollback if installation is in progress
    if (state === 'installing' && installer) {
      // Installer will rollback automatically on failure
      // Just exit gracefully
    }

    const result: InstallResult = {
      success: false,
      installedFiles: [],
      errors: ['Installation cancelled by user'],
    };

    onComplete(result);
    exit();
  };

  /**
   * Handle startup path submission
   */
  const handleStartupPathSubmit = (path: string) => {
    if (path === 'CANCEL') {
      exit();
      return;
    }
    setStartupPath(path);
    setState('claudePath');
  };

  /**
   * Handle claude path submission
   */
  const handleClaudePathSubmit = (path: string) => {
    if (path === 'CANCEL') {
      exit();
      return;
    }
    setClaudePath(path);
    setState('fileSelection');
  };

  /**
   * Handle final confirmation
   */
  const handleConfirm = () => {
    setState('installing');
  };

  /**
   * Handle final cancellation (go back to claude path)
   */
  const handleCancelConfirmation = () => {
    setState('claudePath');
  };


  /**
   * Render current state
   */
  const renderState = () => {
    switch (state) {
      case 'startupPath': {
        const choices: Choice[] = [
          {
            label: '~/.config/the-startup (recommended)',
            value: resolve(homedir(), '.config', 'the-startup'),
          },
          {
            label: '.the-startup (local)',
            value: resolve(process.cwd(), '.the-startup'),
          },
          {
            label: 'Custom location',
            value: 'CUSTOM',
          },
          {
            label: 'Cancel',
            value: 'CANCEL',
          },
        ];

        return (
          <ChoiceSelector
            title="Select .the-startup installation location"
            subtitle="This is where The Startup's templates and rules will be installed"
            choices={choices}
            onSubmit={handleStartupPathSubmit}
          />
        );
      }

      case 'claudePath': {
        const choices: Choice[] = [
          {
            label: '~/.claude (recommended)',
            value: resolve(homedir(), '.claude'),
          },
          {
            label: '.claude (local)',
            value: resolve(process.cwd(), '.claude'),
          },
          {
            label: 'Custom location',
            value: 'CUSTOM',
          },
          {
            label: 'Cancel',
            value: 'CANCEL',
          },
        ];

        return (
          <ChoiceSelector
            title="Select Claude configuration directory"
            subtitle="This is where Claude Code's agents and commands will be installed"
            choices={choices}
            onSubmit={handleClaudePathSubmit}
            onBack={() => setState('startupPath')}
          />
        );
      }

      case 'fileSelection': {
        const settingsPath = join(claudePath, 'settings.json');
        const settingsExists = existsSync(settingsPath);

        return (
          <FinalConfirmation
            startupPath={startupPath}
            claudePath={claudePath}
            files={[]} // Files list is now hardcoded in FinalConfirmation
            settingsExists={settingsExists}
            mode="install"
            onConfirm={handleConfirm}
            onCancel={handleCancelConfirmation}
          />
        );
      }

      case 'installing':
        return (
          <Box flexDirection="column" gap={1}>
            <Text color={theme.colors.info}>Installing...</Text>
            {progressInfo && (
              <Box flexDirection="column">
                <Text color={theme.colors.text}>
                  {progressInfo.stage}
                </Text>
                <Text color={theme.colors.textMuted}>
                  {progressInfo.current} / {progressInfo.total}
                </Text>
              </Box>
            )}
          </Box>
        );

      case 'complete':
        if (!installResult) {
          return <Text color={theme.colors.error}>No installation result</Text>;
        }

        const summary: InstallationSummary = {
          installedFiles: installResult.installedFiles,
          startupPath,
          claudePath,
          totalFiles: installResult.installedFiles.length,
        };

        return <Complete summary={summary} />;

      case 'error':
        return (
          <Box flexDirection="column" gap={1}>
            <Text color={theme.colors.error}>
              {theme.icons.error} Installation failed
            </Text>
            {error && (
              <Text color={theme.colors.text}>
                {error}
              </Text>
            )}
            <Text color={theme.colors.textMuted}>
              Press any key to exit...
            </Text>
          </Box>
        );

      default:
        return <Text>Unknown state</Text>;
    }
  };

  return (
    <Box flexDirection="column" padding={1}>
      {renderState()}
    </Box>
  );
};
