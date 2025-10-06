import { FC } from 'react';
import { Text } from 'ink';
import { theme } from './theme';

export interface ErrorDisplayProps {
  /**
   * Error to display - can be Error object, string, or unknown type
   */
  error: Error | string | unknown;
}

/**
 * Error display component following Ink component pattern
 * Displays errors with consistent formatting and theme colors
 *
 * @example
 * ```tsx
 * <ErrorDisplay error={new Error("Connection failed")} />
 * <ErrorDisplay error="Invalid input" />
 * ```
 */
export const ErrorDisplay: FC<ErrorDisplayProps> = ({ error }) => {
  const errorMessage = getErrorMessage(error);

  return (
    <Text color={theme.colors.error}>
      {theme.icons.error} {errorMessage}
    </Text>
  );
};

/**
 * Extracts error message from various error types
 * Handles Error objects, strings, and unknown types
 */
function getErrorMessage(error: Error | string | unknown): string {
  if (error instanceof Error) {
    return error.message || 'An unknown error occurred';
  }

  if (typeof error === 'string') {
    return error || 'An unknown error occurred';
  }

  return 'An unknown error occurred';
}
