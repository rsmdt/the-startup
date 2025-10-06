import { FC, useState, useEffect } from 'react';
import { Box, Text, useInput, useApp } from 'ink';
import { Spinner } from '../shared/Spinner';
import { ErrorDisplay } from '../shared/ErrorDisplay';
import { theme } from '../shared/theme';
import type { UninstallCommandOptions } from '../../core/types/config';
import type { LockManager } from '../../core/installer/LockManager';
import type { SettingsMerger } from '../../core/installer/SettingsMerger';
import type { LockFile, FileEntry } from '../../core/types/lock';
import { promises as fs } from 'fs';

/**
 * Wizard state machine
 * confirmation -> uninstalling -> complete -> error
 */
type WizardState = 'confirmation' | 'uninstalling' | 'complete' | 'error';

/**
 * Uninstall result returned to onComplete callback
 */
export interface UninstallResult {
  success: boolean;
  filesRemoved: number;
  errors?: string[];
}

/**
 * File System Interface (for testing)
 */
interface FileSystem {
  rm(path: string, options?: { force?: boolean; recursive?: boolean }): Promise<void>;
}

export interface UninstallWizardProps {
  /**
   * CLI command options (--keep-logs, --keep-settings flags)
   */
  options: UninstallCommandOptions;

  /**
   * Lock manager for reading installed files
   */
  lockManager: LockManager;

  /**
   * Settings merger for removing hooks
   */
  settingsMerger: SettingsMerger & {
    removeHooks?: (settingsPath: string) => Promise<void>;
  };

  /**
   * Callback when uninstall completes (success or failure)
   */
  onComplete: (result: UninstallResult) => void;

  /**
   * Optional file system abstraction for testing
   */
  fileSystem?: FileSystem;
}

/**
 * UninstallWizard - Interactive uninstall TUI state machine
 *
 * Implements the 3-state wizard flow:
 * 1. Confirmation: Show file list and ask for confirmation
 * 2. Uninstalling: Remove files and clean up
 * 3. Complete: Display success summary
 *
 * Features:
 * - --keep-logs flag: Preserve .the-startup/logs directory
 * - --keep-settings flag: Don't modify settings.json
 * - Ctrl+C: Cancellation support
 * - Error handling: Continue on missing files (ENOENT)
 * - Progress indication: Show progress during uninstall
 *
 * Business Rules (PRD lines 177-189):
 * - Rule 1: Read lock file to identify installed files
 * - Rule 2: Confirm before deleting
 * - Rule 3: Handle missing files gracefully (skip ENOENT)
 * - Rule 4: Remove hooks from settings.json (unless --keep-settings)
 * - Rule 5: Delete lock file after successful uninstall
 *
 * @example
 * ```tsx
 * <UninstallWizard
 *   options={{ keepLogs: false, keepSettings: false }}
 *   lockManager={lockManagerInstance}
 *   settingsMerger={settingsMergerInstance}
 *   onComplete={(result) => console.log('Done:', result)}
 * />
 * ```
 */
export const UninstallWizard: FC<UninstallWizardProps> = ({
  options,
  lockManager,
  settingsMerger,
  onComplete,
  fileSystem,
}) => {
  // State machine
  const [state, setState] = useState<WizardState>('confirmation');
  const [lockFile, setLockFile] = useState<LockFile | null>(null);
  const [filesRemoved, setFilesRemoved] = useState<number>(0);
  const [error, setError] = useState<string | null>(null);
  const [warnings, setWarnings] = useState<string[]>([]);

  const { exit } = useApp();

  // Use provided fileSystem or default to fs.promises
  const fsImpl = fileSystem || fs;

  // Load lock file on mount
  useEffect(() => {
    loadLockFile();
  }, []);

  /**
   * Load lock file to get installed files
   */
  const loadLockFile = async () => {
    try {
      const lock = await lockManager.readLockFile();

      if (!lock) {
        setError('No installation found. Nothing to uninstall.');
        setState('error');
        return;
      }

      setLockFile(lock);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to read lock file';
      setError(errorMessage);
      setState('error');
    }
  };

  /**
   * Handle user confirmation input
   */
  useInput((input, key) => {
    if (state !== 'confirmation') {
      return;
    }

    // Handle Ctrl+C
    if (key.ctrl && input === 'c') {
      handleCancellation();
      return;
    }

    // Handle yes/no
    if (input.toLowerCase() === 'y') {
      setState('uninstalling');
    } else if (input.toLowerCase() === 'n') {
      handleCancellation();
    }
  });

  /**
   * Perform uninstall when in uninstalling state
   */
  useEffect(() => {
    if (state === 'uninstalling' && lockFile) {
      performUninstall();
    }
  }, [state, lockFile]);

  /**
   * Handle cancellation
   */
  const handleCancellation = () => {
    const result: UninstallResult = {
      success: false,
      filesRemoved: 0,
      errors: ['Uninstall cancelled by user'],
    };

    onComplete(result);
    exit();
  };

  /**
   * Perform the actual uninstall
   */
  const performUninstall = async () => {
    if (!lockFile) {
      setError('No lock file loaded');
      setState('error');
      return;
    }

    try {
      let removedCount = 0;
      const newWarnings: string[] = [];

      // Collect all file paths from lock file
      const allFiles = collectAllFiles(lockFile);

      // Delete files
      for (const fileEntry of allFiles) {
        const shouldSkip = shouldSkipFile(fileEntry.path);

        if (shouldSkip) {
          continue;
        }

        try {
          await fsImpl.rm(fileEntry.path, { force: true });
          removedCount++;
        } catch (err) {
          const error = err as NodeJS.ErrnoException;

          // Skip missing files (already deleted)
          if (error.code === 'ENOENT') {
            continue;
          }

          // Log warning for other errors but continue
          newWarnings.push(`Failed to delete ${fileEntry.path}: ${error.message}`);
        }
      }

      // Remove hooks from settings.json (unless --keep-settings)
      if (!options.keepSettings && settingsMerger.removeHooks) {
        try {
          const settingsPath = deriveSettingsPath(lockFile);
          await settingsMerger.removeHooks(settingsPath);
        } catch (err) {
          const error = err instanceof Error ? err.message : 'Unknown error';
          newWarnings.push(`Failed to remove hooks: ${error}`);
        }
      }

      // Delete lock file
      try {
        const lockFilePath = deriveLockFilePath(lockFile);
        await fsImpl.rm(lockFilePath, { force: true });
      } catch (err) {
        const error = err as NodeJS.ErrnoException;
        if (error.code !== 'ENOENT') {
          newWarnings.push(`Failed to delete lock file: ${error.message}`);
        }
      }

      setFilesRemoved(removedCount);
      setWarnings(newWarnings);
      setState('complete');

      onComplete({
        success: true,
        filesRemoved: removedCount,
        errors: newWarnings.length > 0 ? newWarnings : undefined,
      });
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Uninstall failed';
      setError(errorMessage);
      setState('error');

      onComplete({
        success: false,
        filesRemoved: 0,
        errors: [errorMessage],
      });
    }
  };

  /**
   * Collect all file paths from lock file
   */
  const collectAllFiles = (lock: LockFile): FileEntry[] => {
    const files: FileEntry[] = [];

    files.push(...lock.files.agents);
    files.push(...lock.files.commands);
    files.push(...lock.files.templates);
    files.push(...lock.files.rules);
    files.push(...lock.files.outputStyles);

    if (lock.files.binary.path) {
      files.push(lock.files.binary);
    }

    return files;
  };

  /**
   * Determine if a file should be skipped based on flags
   */
  const shouldSkipFile = (filePath: string): boolean => {
    // Skip logs if --keep-logs flag is set
    if (options.keepLogs && filePath.includes('/logs/')) {
      return true;
    }

    return false;
  };

  /**
   * Derive settings.json path from lock file
   */
  const deriveSettingsPath = (lock: LockFile): string => {
    // Find any .claude file to derive the path
    const claudeFile =
      lock.files.agents[0]?.path ||
      lock.files.commands[0]?.path;

    if (claudeFile) {
      // Extract .claude directory path
      const claudeDir = claudeFile.substring(0, claudeFile.indexOf('.claude') + 7);
      return `${claudeDir}/settings.json`;
    }

    // Fallback to common path
    return `${process.env.HOME}/.claude/settings.json`;
  };

  /**
   * Derive lock file path from lock file data
   */
  const deriveLockFilePath = (lock: LockFile): string => {
    // Find any .the-startup file to derive the path
    const startupFile =
      lock.files.templates[0]?.path ||
      lock.files.rules[0]?.path ||
      lock.files.binary?.path;

    if (startupFile) {
      // Extract .the-startup directory path
      const startupDir = startupFile.substring(
        0,
        startupFile.indexOf('.the-startup') + 12
      );
      return `${startupDir}/.the-startup.lock`;
    }

    // Fallback to current directory
    return './.the-startup/.the-startup.lock';
  };

  /**
   * Count total files
   */
  const countFiles = (lock: LockFile | null): number => {
    if (!lock) {
      return 0;
    }

    return (
      lock.files.agents.length +
      lock.files.commands.length +
      lock.files.templates.length +
      lock.files.rules.length +
      lock.files.outputStyles.length +
      (lock.files.binary.path ? 1 : 0)
    );
  };

  /**
   * Render confirmation state
   */
  const renderConfirmation = () => {
    if (!lockFile) {
      return <Spinner text="Loading installation info..." />;
    }

    const totalFiles = countFiles(lockFile);

    return (
      <Box flexDirection="column" gap={1}>
        <Text color={theme.colors.warning}>
          {theme.icons.warning} Uninstall confirmation
        </Text>
        <Text color={theme.colors.text}>
          The following {totalFiles} files will be removed:
        </Text>
        <Box flexDirection="column" marginLeft={2}>
          <Text color={theme.colors.textMuted}>
            {lockFile.files.agents.length} agent files
          </Text>
          <Text color={theme.colors.textMuted}>
            {lockFile.files.commands.length} command files
          </Text>
          <Text color={theme.colors.textMuted}>
            {lockFile.files.templates.length} template files
          </Text>
          <Text color={theme.colors.textMuted}>
            {lockFile.files.rules.length} rule files
          </Text>
          <Text color={theme.colors.textMuted}>
            {lockFile.files.outputStyles.length} output style files
          </Text>
          {lockFile.files.binary.path && (
            <Text color={theme.colors.textMuted}>1 binary file</Text>
          )}
        </Box>
        {options.keepLogs && (
          <Text color={theme.colors.info}>Logs will be preserved (--keep-logs)</Text>
        )}
        {options.keepSettings && (
          <Text color={theme.colors.info}>
            Settings.json will not be modified (--keep-settings)
          </Text>
        )}
        <Text color={theme.colors.text}>
          Continue with uninstall? [y/n]
        </Text>
      </Box>
    );
  };

  /**
   * Render uninstalling state
   */
  const renderUninstalling = () => {
    return (
      <Box flexDirection="column" gap={1}>
        <Spinner text="Uninstalling..." />
        <Text color={theme.colors.textMuted}>Removing installed files...</Text>
      </Box>
    );
  };

  /**
   * Render complete state
   */
  const renderComplete = () => {
    return (
      <Box flexDirection="column" gap={1}>
        <Text color={theme.colors.success}>
          {theme.icons.success} Uninstall completed successfully
        </Text>
        <Text color={theme.colors.text}>Removed {filesRemoved} files</Text>
        {warnings.length > 0 && (
          <Box flexDirection="column" marginTop={1}>
            <Text color={theme.colors.warning}>Warnings:</Text>
            {warnings.map((warning, i) => (
              <Text key={i} color={theme.colors.textMuted}>
                - {warning}
              </Text>
            ))}
          </Box>
        )}
      </Box>
    );
  };

  /**
   * Render error state
   */
  const renderError = () => {
    return (
      <Box flexDirection="column" gap={1}>
        <ErrorDisplay error={error || 'Unknown error'} />
        <Text color={theme.colors.textMuted}>Press any key to exit...</Text>
      </Box>
    );
  };

  /**
   * Render current state
   */
  const renderState = () => {
    switch (state) {
      case 'confirmation':
        return renderConfirmation();
      case 'uninstalling':
        return renderUninstalling();
      case 'complete':
        return renderComplete();
      case 'error':
        return renderError();
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
