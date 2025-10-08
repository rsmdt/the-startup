import { FC, useState } from 'react';
import { Box, Text, useInput } from 'ink';
import { theme } from '../shared/theme';
import { AppBanner } from '../shared/Banner';
import { homedir } from 'os';

export interface FinalConfirmationProps {
  /**
   * Installation directory path
   */
  startupPath: string;

  /**
   * Claude configuration directory path
   */
  claudePath: string;

  /**
   * List of files to be installed
   */
  files: string[];

  /**
   * Whether settings.json already exists
   */
  settingsExists: boolean;

  /**
   * Mode: 'install' or 'uninstall'
   */
  mode: 'install' | 'uninstall';

  /**
   * Callback when user confirms
   */
  onConfirm: () => void;

  /**
   * Callback when user cancels
   */
  onCancel: () => void;
}

/**
 * FinalConfirmation - Final screen showing files and yes/no confirmation
 *
 * Matches the Go installer's file selection screen:
 * - Shows banner
 * - Displays both paths (startup + claude)
 * - Shows static tree of files to be installed/removed
 * - Shows settings.json with update/create status
 * - Simple two-choice confirmation (arrow key navigation)
 *
 * @example
 * ```tsx
 * <FinalConfirmation
 *   startupPath="~/.the-startup"
 *   claudePath="~/.claude"
 *   files={['agents/the-chief.md', 'commands/prd.md']}
 *   settingsExists={true}
 *   mode="install"
 *   onConfirm={() => console.log('Installing...')}
 *   onCancel={() => console.log('Cancelled')}
 * />
 * ```
 */
export const FinalConfirmation: FC<FinalConfirmationProps> = ({
  startupPath,
  claudePath,
  files,
  settingsExists,
  mode,
  onConfirm,
  onCancel,
}) => {
  const [cursor, setCursor] = useState(0);

  // Confirmation choices
  const choices =
    mode === 'install'
      ? ['Yes, give me awesome', 'Huh? I did not sign up for this']
      : ['Yes, remove everything', 'No, keep everything as-is'];

  // Format paths for display (replace home directory with ~)
  const formatPath = (path: string): string => {
    const home = homedir();
    if (path.startsWith(home)) {
      return path.replace(home, '~');
    }
    return path;
  };

  // Keyboard navigation
  useInput((input, key) => {
    if (key.upArrow || input === 'k') {
      setCursor((prev) => Math.max(0, prev - 1));
    } else if (key.downArrow || input === 'j') {
      setCursor((prev) => Math.min(choices.length - 1, prev + 1));
    } else if (key.return) {
      if (cursor === 0) {
        onConfirm();
      } else {
        onCancel();
      }
    } else if (key.escape || key.backspace || key.delete) {
      onCancel();
    }
  });

  // Define the actual file structure
  // Note: In a real implementation, we'd check checksums from the lock file
  // For now, we show all files without "will update" status since we can't
  // easily check checksums in the UI layer without the installer's lock manager
  const getFileStructure = () => {
    return {
      agents: [
        { name: 'the-analyst', activities: 3 },
        { name: 'the-architect', activities: 5 },
        { name: 'the-chief', isFile: true },
        { name: 'the-designer', activities: 4 },
        { name: 'the-meta-agent', isFile: true },
        { name: 'the-ml-engineer', activities: 4 },
        { name: 'the-mobile-engineer', activities: 3 },
        { name: 'the-platform-engineer', activities: 7 },
        { name: 'the-qa-engineer', activities: 3 },
        { name: 'the-security-engineer', activities: 3 },
        { name: 'the-software-engineer', activities: 5 },
      ],
      commands: [
        { name: '/s:analyze' },
        { name: '/s:implement' },
        { name: '/s:init' },
        { name: '/s:refactor' },
        { name: '/s:specify' },
      ],
      outputStyles: [
        { name: 'the-startup', isFile: true },
      ],
      settings: {
        exists: settingsExists,
      },
      settingsLocal: {
        exists: false, // We'd need to check this
      },
    };
  };

  const structure = getFileStructure();

  return (
    <Box flexDirection="column" padding={1}>
      {/* Banner */}
      <AppBanner />

      <Box marginTop={1} marginBottom={1} flexDirection="column">
        {/* Paths header */}
        <Text bold color={mode === 'install' ? theme.colors.info : theme.colors.warning}>
          {mode === 'install' ? 'Installation Paths:' : 'Uninstallation Paths:'}
        </Text>
        <Box marginTop={1} flexDirection="column">
          <Text color={theme.colors.text}>  Startup: {formatPath(startupPath)}</Text>
          <Text color={theme.colors.text}>  Claude:  {formatPath(claudePath)}</Text>
        </Box>
      </Box>

      {/* Files section */}
      <Box marginTop={1} marginBottom={1} flexDirection="column">
        <Text bold color={theme.colors.info}>
          {mode === 'install' ? 'Files to be installed to .claude' : 'Files to be removed'}
        </Text>
        <Text color={mode === 'install' ? theme.colors.info : theme.colors.warning}>
          {mode === 'install'
            ? 'The following files will be installed to your Claude directory:'
            : `${files.length} files will be removed:`}
        </Text>
      </Box>

      {/* File tree display - matches Go version structure */}
      <Box marginTop={1} flexDirection="column">
        <Text color={theme.colors.primary}>{'# ~/.claude'}</Text>

        {/* Agents */}
        <Text color={theme.colors.text}>{'├─ agents/'}</Text>
        {structure.agents.map((agent, idx) => {
          const isLast = idx === structure.agents.length - 1;
          const prefix = isLast ? '└─' : '├─';

          if (agent.activities) {
            // Agent directory with activities
            const activityLabel = agent.activities === 1
              ? 'specialized activity'
              : 'specialized activities';
            return (
              <Text key={agent.name} color={theme.colors.info}>
                {'│  '}{prefix} {agent.name} ({agent.activities} {activityLabel})
              </Text>
            );
          } else {
            // Single file agent (no status shown unless we have checksum comparison)
            return (
              <Text key={agent.name} color={theme.colors.info}>
                {'│  '}{prefix} {agent.name}
              </Text>
            );
          }
        })}

        {/* Commands */}
        <Text color={theme.colors.text}>{'├─ commands/'}</Text>
        {structure.commands.map((cmd, idx) => {
          const isLast = idx === structure.commands.length - 1;
          const prefix = isLast ? '└─' : '├─';

          return (
            <Text key={cmd.name} color={theme.colors.info}>
              {'│  '}{prefix} {cmd.name}
            </Text>
          );
        })}

        {/* Output Styles */}
        <Text color={theme.colors.text}>{'├─ output-styles/'}</Text>
        {structure.outputStyles.map((style, idx) => {
          const isLast = idx === structure.outputStyles.length - 1;
          const prefix = isLast ? '└─' : '├─';
          return (
            <Text key={style.name} color={theme.colors.info}>
              {'│  '}{prefix} {style.name}
            </Text>
          );
        })}

        {/* Settings files */}
        {mode === 'install' && (
          <>
            <Text color={structure.settings.exists ? theme.colors.warning : theme.colors.text}>
              {'├─ settings.json'}{structure.settings.exists ? ' (will update)' : ''}
            </Text>
            <Text color={structure.settingsLocal.exists ? theme.colors.warning : theme.colors.text}>
              {'└─ settings.local.json'}{structure.settingsLocal.exists ? ' (will update)' : ''}
            </Text>
          </>
        )}
      </Box>

      {/* Ready to install/uninstall */}
      <Box marginTop={2} marginBottom={1} flexDirection="column">
        <Text bold color={theme.colors.info}>
          {mode === 'install' ? 'Ready to install?' : 'Ready to uninstall?'}
        </Text>
        <Text color={mode === 'install' ? theme.colors.info : theme.colors.warning}>
          {mode === 'install'
            ? 'This will install The (Agentic) Startup to the selected directories.'
            : 'This will remove The (Agentic) Startup from the selected directories.'}
        </Text>
      </Box>

      {/* Confirmation choices */}
      <Box flexDirection="column" gap={0}>
        {choices.map((choice, index) => {
          const isCursor = index === cursor;
          const cursor_icon = isCursor ? '▸' : ' ';
          const color = isCursor ? theme.colors.primary : theme.colors.text;

          return (
            <Box key={index}>
              <Text color={color}>
                {cursor_icon} {choice}
              </Text>
            </Box>
          );
        })}
      </Box>

      {/* Help text */}
      <Box marginTop={2}>
        <Text color={theme.colors.textMuted}>
          ↑↓ navigate • Enter: select • ESC: back • Ctrl-C: quit
        </Text>
      </Box>
    </Box>
  );
};
