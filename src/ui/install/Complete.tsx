import { FC } from 'react';
import { Box, Text } from 'ink';
import { theme } from '../shared/theme';
import { basename, dirname } from 'path';

export interface InstallationSummary {
  /**
   * List of absolute paths to installed files
   */
  installedFiles: string[];

  /**
   * Path where the-startup was installed
   */
  startupPath: string;

  /**
   * Path to Claude configuration directory
   */
  claudePath: string;

  /**
   * Total number of files installed
   */
  totalFiles: number;
}

export interface CompleteProps {
  /**
   * Installation summary information
   */
  summary: InstallationSummary;

  /**
   * Optional custom success message
   */
  message?: string;
}

/**
 * Installation completion screen component
 * Displays success message and summary of installed files
 *
 * @example
 * ```tsx
 * <Complete
 *   summary={{
 *     installedFiles: ['/path/to/file1.md', '/path/to/file2.md'],
 *     startupPath: '~/.the-startup',
 *     claudePath: '~/.claude',
 *     totalFiles: 2
 *   }}
 * />
 * ```
 */
export const Complete: FC<CompleteProps> = ({ summary, message }) => {
  const defaultMessage = 'Installation completed successfully!';

  // Group files by directory for better display
  const groupFilesByDirectory = (files: string[]): Map<string, string[]> => {
    const grouped = new Map<string, string[]>();

    files.forEach((file) => {
      const dir = dirname(file);
      const fileName = basename(file);

      if (!grouped.has(dir)) {
        grouped.set(dir, []);
      }

      grouped.get(dir)?.push(fileName);
    });

    return grouped;
  };

  const groupedFiles = groupFilesByDirectory(summary.installedFiles);

  // Get relative directory name for display
  const getRelativeDirName = (fullPath: string): string => {
    // Extract meaningful directory name relative to startup or claude path
    if (fullPath.includes('.claude')) {
      const parts = fullPath.split('.claude/')[1];
      return parts ? `.claude/${parts}` : '.claude';
    }
    if (fullPath.includes('.the-startup')) {
      const parts = fullPath.split('.the-startup/')[1];
      return parts ? `.the-startup/${parts}` : '.the-startup';
    }
    return dirname(fullPath);
  };

  return (
    <Box flexDirection="column" gap={1} paddingY={1}>
      {/* Success header */}
      <Box gap={1}>
        <Text color={theme.colors.success} bold>
          {theme.icons.success} {message || defaultMessage}
        </Text>
      </Box>

      {/* Summary */}
      <Box flexDirection="column" gap={0} marginTop={1}>
        <Text color={theme.colors.text}>
          Installed {summary.totalFiles} {summary.totalFiles === 1 ? 'file' : 'files'}
        </Text>
        <Text color={theme.colors.textMuted}>
          Startup: {summary.startupPath}
        </Text>
        <Text color={theme.colors.textMuted}>
          Claude: {summary.claudePath}
        </Text>
      </Box>

      {/* Installed files by directory */}
      {summary.installedFiles.length > 0 && (
        <Box flexDirection="column" gap={0} marginTop={1}>
          <Text color={theme.colors.text} bold>
            Installed Files:
          </Text>
          {Array.from(groupedFiles.entries()).map(([dir, files]) => (
            <Box key={dir} flexDirection="column" gap={0} marginLeft={2}>
              <Text color={theme.colors.info}>
                {getRelativeDirName(dir)}
              </Text>
              {files.map((file) => (
                <Box key={file} marginLeft={2}>
                  <Text color={theme.colors.textMuted}>• {file}</Text>
                </Box>
              ))}
            </Box>
          ))}
        </Box>
      )}

      {/* Next steps */}
      <Box flexDirection="column" gap={0} marginTop={1}>
        <Text color={theme.colors.warning} bold>
          Next Steps:
        </Text>
        <Text color={theme.colors.text}>
          1. Restart Claude Code to load the new agents and commands
        </Text>
        <Text color={theme.colors.text}>
          2. Use slash commands like /s:specify or /s:analyze
        </Text>
        <Text color={theme.colors.text}>
          3. Check the templates in {summary.startupPath}/templates
        </Text>
      </Box>

      {/* Footer */}
      <Box marginTop={1}>
        <Text color={theme.colors.textMuted}>
          Press any key to exit...
        </Text>
      </Box>
    </Box>
  );
};
