/**
 * Public API for programmatic usage of the-agentic-startup
 *
 * This module exports the core components for use as a library,
 * allowing other packages to integrate the-agentic-startup programmatically
 * without using the CLI.
 *
 * Usage:
 * ```typescript
 * import { Installer, Initializer, SpecGenerator } from 'the-agentic-startup/lib';
 *
 * const installer = new Installer(fs, lockManager, settingsMerger, assetProvider, '1.0.0');
 * const result = await installer.install(options);
 * ```
 */

// Core installer components
export { Installer } from '../core/installer/Installer.js';
export { LockManager } from '../core/installer/LockManager.js';
export { SettingsMerger } from '../core/installer/SettingsMerger.js';

// Core initialization components
export { Initializer } from '../core/init/Initializer.js';

// Core spec generation components
export { SpecGenerator } from '../core/spec/SpecGenerator.js';

// Type definitions
export type {
  // Config types
  InstallCommandOptions,
  UninstallCommandOptions,
  InitCommandOptions,
  SpecCommandOptions,
  // Core options and results
  InstallerOptions,
  InstallResult,
  InitOptions,
  InitResult,
  SpecOptions,
  SpecResult,
  SpecNumbering,
} from '../core/types/config.js';

export type {
  // Lock file types
  LockFile,
  FileEntry,
  LegacyLockFile,
} from '../core/types/lock.js';

export type {
  // Settings types
  ClaudeSettings,
  PlaceholderMap,
} from '../core/types/settings.js';
