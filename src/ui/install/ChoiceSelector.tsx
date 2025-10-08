import { FC, useState } from 'react';
import { Box, Text, useInput } from 'ink';
import TextInput from 'ink-text-input';
import { theme } from '../shared/theme';
import { AppBanner } from '../shared/Banner';
import { homedir } from 'os';
import { resolve } from 'path';

export interface Choice {
  /**
   * Display label for the choice
   */
  label: string;

  /**
   * Actual value (path, 'CUSTOM', or 'CANCEL')
   */
  value: string;

  /**
   * Optional description/help text
   */
  description?: string;
}

export interface ChoiceSelectorProps {
  /**
   * Title displayed below banner
   */
  title: string;

  /**
   * Subtitle/info message
   */
  subtitle?: string;

  /**
   * Predefined choices to show
   */
  choices: Choice[];

  /**
   * Callback when user submits a selection
   */
  onSubmit: (value: string) => void;

  /**
   * Optional callback when user presses ESC to go back
   */
  onBack?: () => void;
}

/**
 * ChoiceSelector - Arrow key navigation menu with custom input mode
 *
 * Matches the Go installer UX with:
 * - Predefined choices with arrow key navigation
 * - "Custom location" option triggers text input mode
 * - Tab autocomplete in custom mode (future enhancement)
 * - ESC to go back from custom mode to choices
 *
 * @example
 * ```tsx
 * <ChoiceSelector
 *   title="Select .the-startup installation location"
 *   subtitle="This is where The Startup's templates will be installed"
 *   choices={[
 *     { label: '~/.config/the-startup (recommended)', value: '~/.config/the-startup' },
 *     { label: '.the-startup (local)', value: './.the-startup' },
 *     { label: 'Custom location', value: 'CUSTOM' },
 *     { label: 'Cancel', value: 'CANCEL' },
 *   ]}
 *   onSubmit={(path) => console.log('Selected:', path)}
 * />
 * ```
 */
export const ChoiceSelector: FC<ChoiceSelectorProps> = ({
  title,
  subtitle,
  choices,
  onSubmit,
  onBack,
}) => {
  const [cursor, setCursor] = useState(0);
  const [inputMode, setInputMode] = useState(false);
  const [customPath, setCustomPath] = useState('');

  // Expand tilde to home directory
  const expandPath = (path: string): string => {
    if (path.startsWith('~')) {
      return resolve(homedir(), path.slice(1));
    }
    return path;
  };

  // Keyboard navigation in choice mode
  useInput((input, key) => {
    if (inputMode) {
      // In custom input mode, ESC goes back to choices
      if (key.escape) {
        setInputMode(false);
        setCustomPath('');
      }
      return;
    }

    // ESC or Backspace to go back (if onBack callback provided)
    if (key.escape || key.backspace || key.delete) {
      if (onBack) {
        onBack();
      }
      return;
    }

    // Arrow key navigation
    if (key.upArrow || input === 'k') {
      setCursor((prev) => Math.max(0, prev - 1));
    } else if (key.downArrow || input === 'j') {
      setCursor((prev) => Math.min(choices.length - 1, prev + 1));
    }

    // Enter to select
    else if (key.return) {
      const selected = choices[cursor];

      if (selected.value === 'CUSTOM') {
        // Enter custom input mode
        setInputMode(true);
      } else if (selected.value === 'CANCEL') {
        onSubmit('CANCEL');
      } else {
        // Regular choice - expand path and submit
        const expandedPath = expandPath(selected.value);
        onSubmit(expandedPath);
      }
    }
  });

  // Handle custom path submission
  const handleCustomSubmit = () => {
    if (customPath.trim()) {
      const expandedPath = expandPath(customPath);
      onSubmit(expandedPath);
    }
  };

  return (
    <Box flexDirection="column" padding={1}>
      {/* Banner */}
      <AppBanner />

      <Box marginTop={1} marginBottom={1} flexDirection="column">
        {/* Title */}
        <Text bold color={theme.colors.info}>
          {title}
        </Text>

        {/* Subtitle */}
        {subtitle && (
          <Box marginTop={1}>
            <Text color={theme.colors.textMuted}>{subtitle}</Text>
          </Box>
        )}
      </Box>

      {/* Choice mode: Show predefined options */}
      {!inputMode && (
        <Box flexDirection="column" gap={0}>
          {choices.map((choice, index) => {
            const isCursor = index === cursor;
            const cursor_icon = isCursor ? '▸' : ' ';
            const color = isCursor ? theme.colors.primary : theme.colors.text;

            return (
              <Box key={index}>
                <Text color={color}>
                  {cursor_icon} {choice.label}
                </Text>
              </Box>
            );
          })}

          {/* Help text */}
          <Box marginTop={2}>
            <Text color={theme.colors.textMuted}>
              ↑↓ navigate • Enter: select{onBack ? ' • ESC: back' : ''} • Ctrl-C: quit
            </Text>
          </Box>
        </Box>
      )}

      {/* Input mode: Show custom path input */}
      {inputMode && (
        <Box flexDirection="column" gap={1}>
          <Text color={theme.colors.text}>Enter custom path:</Text>

          <Box>
            <Text color={theme.colors.textMuted}>{'> '}</Text>
            <TextInput
              value={customPath}
              placeholder="Enter path (Tab for autocomplete)"
              onChange={setCustomPath}
              onSubmit={handleCustomSubmit}
            />
          </Box>

          {/* Help text */}
          <Box marginTop={1}>
            <Text color={theme.colors.textMuted}>
              Tab: autocomplete • Enter: confirm • Escape: cancel
            </Text>
          </Box>
        </Box>
      )}
    </Box>
  );
};
