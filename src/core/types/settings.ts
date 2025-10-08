/**
 * Claude Code settings structure for hook configuration.
 *
 * This file defines the shape of ~/.claude/settings.json,
 * specifically the hooks section that the-agentic-startup modifies.
 *
 * Deep merge requirements:
 * - Preserve ALL existing user settings (never overwrite)
 * - Only add new hooks if they don't already exist
 * - Keep all non-hook settings untouched
 * - Atomic operation with backup and rollback
 *
 * Placeholder replacement:
 * - {{STARTUP_PATH}}: Replaced with actual installation directory path
 * - {{CLAUDE_PATH}}: Replaced with Claude config directory path
 * - Replacement happens at install time before writing settings.json
 */

/**
 * Claude Code settings structure.
 *
 * This represents the full settings.json file structure,
 * though only the hooks section is modified by the-agentic-startup.
 */
export interface ClaudeSettings {
  /**
   * Additional directories Claude Code can access.
   *
   * Allows Claude Code to read/write files outside default directories.
   */
  permissions?: {
    additionalDirectories?: string[];
  };

  /**
   * Status line configuration for terminal display.
   *
   * Configures what information is shown in the Claude Code status line.
   */
  statusLine?: {
    /** Type of status line (e.g., "command") */
    type?: string;

    /**
     * Command to execute to get status line content.
     *
     * May contain placeholders:
     * - {{STARTUP_PATH}}: Installation directory
     * - {{CLAUDE_PATH}}: Claude config directory
     */
    command?: string;
  };

  /**
   * Hook configurations for Claude Code events.
   *
   * Each hook is triggered by a specific Claude Code event
   * (e.g., user-prompt-submit for statusline).
   */
  hooks?: {
    [hookName: string]: {
      /**
       * Command to execute when hook is triggered.
       *
       * May contain placeholders:
       * - {{STARTUP_PATH}}: Installation directory (e.g., /Users/name/.the-startup)
       * - {{CLAUDE_PATH}}: Claude config directory (e.g., /Users/name/.claude)
       *
       * Placeholders are replaced at install time with absolute paths.
       */
      command: string;

      /** Human-readable description of what the hook does */
      description?: string;

      /** If true, Claude Code continues even if hook fails */
      continueOnError?: boolean;
    };
  };

  /**
   * Other Claude settings (preserved during merge).
   *
   * This allows for any other settings that Claude Code supports,
   * ensuring we don't lose user customizations.
   */
  [key: string]: unknown;
}

/**
 * Placeholder replacement mapping.
 *
 * Maps placeholder names to their actual runtime values.
 * Used during settings.json merge to replace placeholders in hook commands.
 */
export interface PlaceholderMap {
  /** Actual installation directory path (e.g., /Users/name/.the-startup) */
  STARTUP_PATH: string;

  /** Claude config directory path (e.g., /Users/name/.claude) */
  CLAUDE_PATH: string;
}
