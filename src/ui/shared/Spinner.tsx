import { FC, useState, useEffect } from 'react';
import { Text } from 'ink';
import { theme } from './theme';

export interface SpinnerProps {
  /**
   * Optional text to display alongside the spinner
   */
  text?: string;
}

/**
 * Loading spinner component following Ink component pattern
 * Displays an animated spinner with optional text message
 *
 * @example
 * ```tsx
 * <Spinner text="Loading data..." />
 * ```
 */
export const Spinner: FC<SpinnerProps> = ({ text }) => {
  const [frameIndex, setFrameIndex] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => {
      setFrameIndex((prev) => (prev + 1) % theme.icons.spinner.length);
    }, 80);

    return () => clearInterval(interval);
  }, []);

  const spinnerFrame = theme.icons.spinner[frameIndex];

  return (
    <Text color={theme.colors.info}>
      {spinnerFrame}
      {text && ` ${text}`}
    </Text>
  );
};
