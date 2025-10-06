import { FC, useState, useEffect } from 'react';
import { Box, Text } from 'ink';
import TextInput from 'ink-text-input';
import { theme } from '../shared/theme';
import { homedir } from 'os';
import { resolve } from 'path';

export interface PathValidationResult {
  isValid: boolean;
  message: string;
}

export interface PathSelectorProps {
  /**
   * Label to display above the input
   */
  label: string;

  /**
   * Default path value
   */
  defaultValue?: string;

  /**
   * Placeholder text when input is empty
   */
  placeholder?: string;

  /**
   * Validation function to check if path is valid
   */
  validator?: (path: string) => PathValidationResult;

  /**
   * Callback when user submits a valid path
   */
  onSubmit: (path: string) => void;
}

/**
 * Path selector component with validation
 * Allows user to input a file system path with validation feedback
 *
 * @example
 * ```tsx
 * <PathSelector
 *   label="Installation directory"
 *   defaultValue="~/.the-startup"
 *   validator={(path) => ({
 *     isValid: existsSync(path),
 *     message: existsSync(path) ? '' : 'Path does not exist'
 *   })}
 *   onSubmit={(path) => console.log('Selected:', path)}
 * />
 * ```
 */
export const PathSelector: FC<PathSelectorProps> = ({
  label,
  defaultValue = '',
  placeholder = '',
  validator,
  onSubmit,
}) => {
  // Expand tilde to home directory
  const expandPath = (path: string): string => {
    if (path.startsWith('~')) {
      return resolve(homedir(), path.slice(1));
    }
    return path;
  };

  // Initial validation with default value
  const getInitialValidation = (): PathValidationResult => {
    if (validator && defaultValue) {
      return validator(defaultValue);
    }
    return { isValid: true, message: '' };
  };

  const [value, setValue] = useState(defaultValue);
  const [validation, setValidation] = useState<PathValidationResult>(
    getInitialValidation()
  );

  // Validate path whenever value changes
  useEffect(() => {
    if (validator) {
      const result = validator(value);
      setValidation(result);
    }
  }, [value, validator]);

  const handleSubmit = () => {
    // Only submit if path is valid
    if (validation.isValid && value) {
      const expandedPath = expandPath(value);
      onSubmit(expandedPath);
    }
  };

  return (
    <Box flexDirection="column" gap={1}>
      <Text color={theme.colors.text}>{label}</Text>

      <Box>
        <Text color={theme.colors.textMuted}>{'> '}</Text>
        <TextInput
          value={value}
          placeholder={placeholder}
          onChange={setValue}
          onSubmit={handleSubmit}
        />
      </Box>

      {/* Validation feedback */}
      {value && validation.isValid && (
        <Box>
          <Text color={theme.colors.success}>
            {theme.icons.success} Valid path
          </Text>
        </Box>
      )}

      {value && !validation.isValid && validation.message && (
        <Box>
          <Text color={theme.colors.error}>
            {theme.icons.error} {validation.message}
          </Text>
        </Box>
      )}
    </Box>
  );
};
