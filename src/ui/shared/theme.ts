/**
 * Theme configuration for Ink-based TUI
 * Maps from Go internal/ui/theme.go CharmTheme
 */

export interface ThemeColors {
  primary: string;
  success: string;
  error: string;
  warning: string;
  info: string;
  text: string;
  textMuted: string;
  textBright: string;
}

export interface ThemeIcons {
  success: string;
  error: string;
  warning: string;
  info: string;
  spinner: string[];
}

export interface ThemeSpacing {
  small: number;
  medium: number;
  large: number;
}

export interface Theme {
  colors: ThemeColors;
  icons: ThemeIcons;
  spacing: ThemeSpacing;
}

/**
 * Default theme matching Go CharmTheme
 * Colors mapped from internal/ui/theme.go lines 24-35
 */
export const theme: Theme = {
  colors: {
    primary: '#FF06B7', // Pink/Magenta
    success: '#04B575', // Green
    error: '#FF4444', // Red
    warning: '#FFA500', // Orange
    info: '#3C7EFF', // Blue
    text: '#FAFAFA', // Light gray
    textMuted: '#606060', // Dark gray
    textBright: '#42FF76', // Bright green
  },

  icons: {
    success: '✓',
    error: '✗',
    warning: '⚠',
    info: 'ℹ',
    // Spinner frames for loading animation
    spinner: ['⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'],
  },

  spacing: {
    small: 1,
    medium: 2,
    large: 3,
  },
};
