import { FC } from 'react';
import { Box, Text } from 'ink';
import { theme } from '../shared/theme';

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
    </Box>
  );
};
