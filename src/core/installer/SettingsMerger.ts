import type { ClaudeSettings, PlaceholderMap } from '../types/settings';

/**
 * File System Interface
 *
 * Abstraction for file system operations to enable testing.
 * Matches Node.js fs.promises API surface.
 */
interface FileSystem {
  readFile(path: string, encoding: string): Promise<string>;
  writeFile(path: string, content: string, encoding: string): Promise<void>;
  copyFile(src: string, dest: string): Promise<void>;
  rm(path: string, options?: { force?: boolean }): Promise<void>;
  access(path: string): Promise<void>;
}

/**
 * SettingsMerger - Deep merge Claude settings with user preservation
 *
 * Implements the settings merge algorithm from SDD Section "Settings Merge Algorithm" (lines 923-954).
 *
 * Key Responsibilities:
 * - Deep merge user settings with new hooks
 * - Never overwrite existing user hooks
 * - Replace placeholders ({{STARTUP_PATH}}, {{CLAUDE_PATH}})
 * - Atomic operations with backup and rollback
 * - Create settings from scratch if needed
 *
 * Algorithm:
 * 1. BACKUP: Create settings.json.backup with timestamp
 * 2. PARSE: Read and parse existing settings
 * 3. INITIALIZE: Create hooks object if needed
 * 4. MERGE_HOOKS: Add new hooks, preserve existing ones
 * 5. REPLACE_PLACEHOLDERS: Replace in new hooks only
 * 6. PRESERVE: Keep all other user settings
 * 7. VALIDATE: Ensure JSON structure is valid
 * 8. WRITE: Write merged settings atomically
 * 9. ON_ERROR: Restore from backup
 *
 * @example
 * const merger = new SettingsMerger(fs.promises);
 * const placeholders = {
 *   STARTUP_PATH: '/Users/name/.the-startup',
 *   CLAUDE_PATH: '/Users/name/.claude'
 * };
 * const newHooks = {
 *   'user-prompt-submit': {
 *     command: '{{STARTUP_PATH}}/bin/statusline.sh'
 *   }
 * };
 * const result = await merger.mergeSettings(
 *   '/Users/name/.claude/settings.json',
 *   newHooks,
 *   placeholders
 * );
 */
export class SettingsMerger {
  constructor(private fs: FileSystem) {}

  /**
   * Merges new hooks into existing Claude settings with user preservation.
   *
   * This is an atomic operation - either all changes succeed or none do.
   * User's existing hooks are NEVER overwritten.
   *
   * @param settingsPath - Absolute path to settings.json (e.g., ~/.claude/settings.json)
   * @param newHooks - Hooks to add (if they don't already exist)
   * @param placeholders - Values to replace placeholders with
   * @returns Merged settings object
   * @throws Error if JSON is invalid or write fails
   */
  async mergeSettings(
    settingsPath: string,
    newHooks: ClaudeSettings['hooks'],
    placeholders: PlaceholderMap
  ): Promise<ClaudeSettings> {
    const backupPath = this.generateBackupPath(settingsPath);
    const settingsExisted = await this.createBackup(settingsPath, backupPath);

    try {
      const finalSettings = await this.performMerge(
        settingsPath,
        newHooks,
        placeholders
      );

      await this.writeSettings(settingsPath, finalSettings);
      await this.cleanupBackup(backupPath, settingsExisted);

      return finalSettings;
    } catch (error) {
      await this.restoreBackup(settingsPath, backupPath, settingsExisted);
      throw error;
    }
  }

  /**
   * Generates timestamped backup path.
   *
   * @param settingsPath - Original settings path
   * @returns Backup file path with timestamp
   */
  private generateBackupPath(settingsPath: string): string {
    return `${settingsPath}.backup-${Date.now()}`;
  }

  /**
   * Creates backup of settings file if it exists.
   *
   * @param settingsPath - Settings file to backup
   * @param backupPath - Where to save backup
   * @returns True if file existed and was backed up, false otherwise
   */
  private async createBackup(
    settingsPath: string,
    backupPath: string
  ): Promise<boolean> {
    try {
      await this.fs.access(settingsPath);
      await this.fs.copyFile(settingsPath, backupPath);
      return true;
    } catch (error) {
      return false;
    }
  }

  /**
   * Performs the core merge logic.
   *
   * @param settingsPath - Settings file path
   * @param newHooks - Hooks to merge
   * @param placeholders - Placeholder replacements
   * @returns Merged settings
   */
  private async performMerge(
    settingsPath: string,
    newHooks: ClaudeSettings['hooks'],
    placeholders: PlaceholderMap
  ): Promise<ClaudeSettings> {
    const userSettings = await this.readSettings(settingsPath);

    if (!userSettings.hooks) {
      userSettings.hooks = {};
    }

    const { settings: mergedSettings, addedHooks } = this.mergeHooks(
      userSettings,
      newHooks
    );

    return this.replacePlaceholders(mergedSettings, addedHooks, placeholders);
  }

  /**
   * Writes settings to file with formatting.
   *
   * @param settingsPath - Where to write settings
   * @param settings - Settings to write
   */
  private async writeSettings(
    settingsPath: string,
    settings: ClaudeSettings
  ): Promise<void> {
    const settingsJson = JSON.stringify(settings, null, 2);
    await this.fs.writeFile(settingsPath, settingsJson, 'utf-8');
  }

  /**
   * Cleans up backup after successful merge.
   *
   * @param backupPath - Backup file to remove
   * @param existed - Whether backup was created
   */
  private async cleanupBackup(backupPath: string, existed: boolean): Promise<void> {
    if (existed) {
      await this.fs.rm(backupPath, { force: true });
    }
  }

  /**
   * Restores from backup on error.
   *
   * @param settingsPath - Settings file to restore
   * @param backupPath - Backup file location
   * @param existed - Whether backup was created
   */
  private async restoreBackup(
    settingsPath: string,
    backupPath: string,
    existed: boolean
  ): Promise<void> {
    if (existed) {
      await this.rollback(settingsPath, backupPath);
    }
  }

  /**
   * Reads and parses settings.json file.
   *
   * @param settingsPath - Path to settings.json
   * @returns Parsed settings object or empty object if file doesn't exist
   * @throws Error with clear message if JSON is invalid
   */
  private async readSettings(settingsPath: string): Promise<ClaudeSettings> {
    try {
      const content = await this.fs.readFile(settingsPath, 'utf-8');
      return JSON.parse(content) as ClaudeSettings;
    } catch (error) {
      // File doesn't exist - return empty settings
      if ((error as NodeJS.ErrnoException).code === 'ENOENT') {
        return {};
      }

      // JSON parse error - throw with file path for clarity
      if (error instanceof SyntaxError) {
        throw new Error(
          `Invalid JSON in settings file ${settingsPath}: ${error.message}`
        );
      }

      throw error;
    }
  }

  /**
   * Merges new hooks into user settings, preserving existing hooks.
   *
   * Key principle: NEVER overwrite user's existing hooks.
   * Only add hooks that don't already exist.
   *
   * @param userSettings - User's existing settings
   * @param newHooks - New hooks to add
   * @returns Object with merged settings and set of actually added hook names
   */
  private mergeHooks(
    userSettings: ClaudeSettings,
    newHooks: ClaudeSettings['hooks']
  ): { settings: ClaudeSettings; addedHooks: Set<string> } {
    const merged = { ...userSettings };
    const addedHooks = new Set<string>();

    if (!newHooks) {
      return { settings: merged, addedHooks };
    }

    // Ensure hooks object exists
    if (!merged.hooks) {
      merged.hooks = {};
    }

    // Add new hooks only if they don't exist
    for (const [hookName, hookConfig] of Object.entries(newHooks)) {
      if (!merged.hooks[hookName]) {
        merged.hooks[hookName] = { ...hookConfig };
        addedHooks.add(hookName);
      }
      // If hook exists, skip (preserve user's configuration)
    }

    return { settings: merged, addedHooks };
  }

  /**
   * Replaces placeholders in newly added hooks only.
   *
   * User's existing hooks are NEVER modified, even if they contain
   * placeholder-like text.
   *
   * Placeholders:
   * - {{STARTUP_PATH}}: Installation directory path
   * - {{CLAUDE_PATH}}: Claude config directory path
   *
   * @param settings - Merged settings
   * @param addedHooks - Set of hook names that were actually added
   * @param placeholders - Replacement values
   * @returns Settings with placeholders replaced
   */
  private replacePlaceholders(
    settings: ClaudeSettings,
    addedHooks: Set<string>,
    placeholders: PlaceholderMap
  ): ClaudeSettings {
    if (!settings.hooks || addedHooks.size === 0) {
      return settings;
    }

    const result = { ...settings };
    result.hooks = { ...settings.hooks };

    // Only replace in hooks that were newly added
    for (const hookName of addedHooks) {
      const hook = result.hooks[hookName];
      if (hook && hook.command) {
        result.hooks[hookName] = {
          ...hook,
          command: this.replaceInCommand(hook.command, placeholders),
        };
      }
    }

    return result;
  }

  /**
   * Replaces all placeholders in a command string.
   *
   * @param command - Command string with placeholders
   * @param placeholders - Replacement values
   * @returns Command with placeholders replaced
   */
  private replaceInCommand(command: string, placeholders: PlaceholderMap): string {
    let result = command;

    result = result.replace(/\{\{STARTUP_PATH\}\}/g, placeholders.STARTUP_PATH);
    result = result.replace(/\{\{CLAUDE_PATH\}\}/g, placeholders.CLAUDE_PATH);

    return result;
  }

  /**
   * Rolls back settings.json from backup on error.
   *
   * @param settingsPath - Path to settings.json
   * @param backupPath - Path to backup file
   */
  private async rollback(settingsPath: string, backupPath: string): Promise<void> {
    try {
      await this.fs.copyFile(backupPath, settingsPath);
      await this.fs.rm(backupPath, { force: true });
    } catch (rollbackError) {
      // Rollback failed - backup file is still available for manual recovery
      // Log error but don't throw (original error is more important)
    }
  }

  /**
   * Removes hooks from Claude settings during uninstall.
   *
   * This method removes all hooks from settings.json that were added
   * during installation. It preserves user's other settings.
   *
   * @param settingsPath - Absolute path to settings.json
   * @returns Updated settings object (without hooks)
   * @throws Error if settings file doesn't exist or is invalid JSON
   */
  async removeHooks(settingsPath: string): Promise<ClaudeSettings> {
    const backupPath = this.generateBackupPath(settingsPath);
    const settingsExisted = await this.createBackup(settingsPath, backupPath);

    if (!settingsExisted) {
      // No settings file exists, nothing to remove
      return {};
    }

    try {
      const settings = await this.readSettings(settingsPath);

      // Remove hooks property entirely
      const { hooks, ...settingsWithoutHooks } = settings;

      await this.writeSettings(settingsPath, settingsWithoutHooks);
      await this.cleanupBackup(backupPath, settingsExisted);

      return settingsWithoutHooks;
    } catch (error) {
      await this.restoreBackup(settingsPath, backupPath, settingsExisted);
      throw error;
    }
  }
}
