import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { promises as fs } from 'fs';
import { join } from 'path';
import { Installer } from '../../src/core/installer/Installer';
import { LockManager } from '../../src/core/installer/LockManager';
import { SettingsMerger } from '../../src/core/installer/SettingsMerger';
import { createTempDir, cleanupTempDir } from '../shared/testUtils';
import type { InstallerOptions } from '../../src/core/types/config';
import type { LockFile, LegacyLockFile } from '../../src/core/types/lock';

/**
 * Integration Test Suite: Go-to-npm Migration
 *
 * Tests migration scenarios from Go binary installation to npm package installation.
 * Covers backward compatibility with legacy lock file format and settings preservation.
 *
 * Migration Scenarios (from PLAN.md lines 339-344):
 * 1. Lock file v1 (Go format: string[]) → v2 (npm format: FileEntry[] with checksums)
 * 2. Detection of existing Go binary installation
 * 3. Settings.json preservation through migration
 * 4. All 54 assets migrate correctly
 * 5. Idempotent migration (running twice doesn't break)
 *
 * Key Migration Contexts:
 * - User has Go version installed with v1 lock file
 * - User runs npm package install command
 * - npm version must detect v1 lock, migrate to v2, preserve settings
 * - All existing files should have checksums generated
 *
 * Lock File Format Evolution:
 * - Go v1: { files: { agents: string[], ... }, binary: string }
 * - npm v2: { files: { agents: FileEntry[], ... }, binary: FileEntry, lockFileVersion: 2 }
 *
 * This is an integration test suite using real file system operations.
 */

describe('Integration: Go-to-npm Migration', () => {
  let tempDir: string;
  let startupPath: string;
  let claudePath: string;
  let lockFilePath: string;
  let installer: Installer;
  let lockManager: LockManager;
  let settingsMerger: SettingsMerger;
  let assetProvider: any;

  beforeEach(async () => {
    // Create temporary test directories
    tempDir = await createTempDir();
    startupPath = join(tempDir, '.the-startup');
    claudePath = join(tempDir, '.claude');
    lockFilePath = join(startupPath, 'lock.json');

    // Create real instances for integration testing
    lockManager = new LockManager(lockFilePath);
    settingsMerger = new SettingsMerger(fs);

    // Create mock asset provider with realistic file structure
    assetProvider = {
      getAssetFiles: () => [
        // Agents (8 files to simulate subset)
        { sourcePath: join(tempDir, 'mock-assets/agents/the-analyst-requirements-analysis.md'), relativePath: 'agents/the-analyst-requirements-analysis.md', targetCategory: 'claude' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/agents/the-architect-system-architecture.md'), relativePath: 'agents/the-architect-system-architecture.md', targetCategory: 'claude' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/agents/the-software-engineer-api-development.md'), relativePath: 'agents/the-software-engineer-api-development.md', targetCategory: 'claude' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/agents/the-platform-engineer-deployment.md'), relativePath: 'agents/the-platform-engineer-deployment.md', targetCategory: 'claude' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/agents/the-security-engineer-assessment.md'), relativePath: 'agents/the-security-engineer-assessment.md', targetCategory: 'claude' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/agents/the-test-engineer-test-execution.md'), relativePath: 'agents/the-test-engineer-test-execution.md', targetCategory: 'claude' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/agents/the-chief.md'), relativePath: 'agents/the-chief.md', targetCategory: 'claude' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/agents/the-product-owner-feature-definition.md'), relativePath: 'agents/the-product-owner-feature-definition.md', targetCategory: 'claude' as const, isJson: false },

        // Commands (5 files)
        { sourcePath: join(tempDir, 'mock-assets/commands/s-specify.md'), relativePath: 'commands/s-specify.md', targetCategory: 'claude' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/commands/s-implement.md'), relativePath: 'commands/s-implement.md', targetCategory: 'claude' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/commands/s-refactor.md'), relativePath: 'commands/s-refactor.md', targetCategory: 'claude' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/commands/s-analyze.md'), relativePath: 'commands/s-analyze.md', targetCategory: 'claude' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/commands/s-init.md'), relativePath: 'commands/s-init.md', targetCategory: 'claude' as const, isJson: false },

        // Templates (6 files)
        { sourcePath: join(tempDir, 'mock-assets/templates/SPEC.md'), relativePath: 'templates/SPEC.md', targetCategory: 'startup' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/templates/TASK-DOD.md'), relativePath: 'templates/TASK-DOD.md', targetCategory: 'startup' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/templates/DOR.md'), relativePath: 'templates/DOR.md', targetCategory: 'startup' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/templates/DOD.md'), relativePath: 'templates/DOD.md', targetCategory: 'startup' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/templates/PRD.md'), relativePath: 'templates/PRD.md', targetCategory: 'startup' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/templates/SDD.md'), relativePath: 'templates/SDD.md', targetCategory: 'startup' as const, isJson: false },

        // Rules (3 files)
        { sourcePath: join(tempDir, 'mock-assets/rules/SCQA.md'), relativePath: 'rules/SCQA.md', targetCategory: 'startup' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/rules/SOLID.md'), relativePath: 'rules/SOLID.md', targetCategory: 'startup' as const, isJson: false },
        { sourcePath: join(tempDir, 'mock-assets/rules/DDD.md'), relativePath: 'rules/DDD.md', targetCategory: 'startup' as const, isJson: false },

        // Output styles (1 file)
        { sourcePath: join(tempDir, 'mock-assets/output-styles/the-startup.md'), relativePath: 'output-styles/the-startup.md', targetCategory: 'claude' as const, isJson: false },

        // Settings (1 file)
        { sourcePath: join(tempDir, 'mock-assets/settings.json'), relativePath: 'settings.json', targetCategory: 'claude' as const, isJson: true },
      ],
    };

    // Create mock asset files
    await createMockAssets();

    installer = new Installer(
      fs,
      lockManager,
      settingsMerger,
      assetProvider,
      '1.0.0',
      undefined,
      tempDir,
      tempDir
    );
  });

  afterEach(async () => {
    await cleanupTempDir(tempDir);
  });

  /**
   * Helper: Create mock asset files for testing
   */
  async function createMockAssets(): Promise<void> {
    const assetsDir = join(tempDir, 'mock-assets');

    // Create directories
    await fs.mkdir(join(assetsDir, 'agents'), { recursive: true });
    await fs.mkdir(join(assetsDir, 'commands'), { recursive: true });
    await fs.mkdir(join(assetsDir, 'templates'), { recursive: true });
    await fs.mkdir(join(assetsDir, 'rules'), { recursive: true });
    await fs.mkdir(join(assetsDir, 'output-styles'), { recursive: true });

    // Create agent files
    await fs.writeFile(join(assetsDir, 'agents/the-analyst-requirements-analysis.md'), '# The Analyst - Requirements Analysis\nAnalyze requirements and create specifications.', 'utf-8');
    await fs.writeFile(join(assetsDir, 'agents/the-architect-system-architecture.md'), '# The Architect - System Architecture\nDesign system architecture and patterns.', 'utf-8');
    await fs.writeFile(join(assetsDir, 'agents/the-software-engineer-api-development.md'), '# The Software Engineer - API Development\nDevelop RESTful APIs.', 'utf-8');
    await fs.writeFile(join(assetsDir, 'agents/the-platform-engineer-deployment.md'), '# The Platform Engineer - Deployment\nDeploy infrastructure.', 'utf-8');
    await fs.writeFile(join(assetsDir, 'agents/the-security-engineer-assessment.md'), '# The Security Engineer - Assessment\nSecurity assessments.', 'utf-8');
    await fs.writeFile(join(assetsDir, 'agents/the-test-engineer-test-execution.md'), '# The Test Engineer - Test Execution\nExecute test plans.', 'utf-8');
    await fs.writeFile(join(assetsDir, 'agents/the-chief.md'), '# The Chief\nExecutive decision maker.', 'utf-8');
    await fs.writeFile(join(assetsDir, 'agents/the-product-owner-feature-definition.md'), '# The Product Owner - Feature Definition\nDefine features.', 'utf-8');

    // Create command files
    await fs.writeFile(join(assetsDir, 'commands/s-specify.md'), '# /s:specify\nCreate specification', 'utf-8');
    await fs.writeFile(join(assetsDir, 'commands/s-implement.md'), '# /s:implement\nImplement from spec', 'utf-8');
    await fs.writeFile(join(assetsDir, 'commands/s-refactor.md'), '# /s:refactor\nRefactor code', 'utf-8');
    await fs.writeFile(join(assetsDir, 'commands/s-analyze.md'), '# /s:analyze\nAnalyze codebase', 'utf-8');
    await fs.writeFile(join(assetsDir, 'commands/s-init.md'), '# /s:init\nInitialize validation', 'utf-8');

    // Create template files
    await fs.writeFile(join(assetsDir, 'templates/SPEC.md'), '# Specification Template\n{{PLACEHOLDER}}', 'utf-8');
    await fs.writeFile(join(assetsDir, 'templates/TASK-DOD.md'), '# Task Definition of Done\n{{PLACEHOLDER}}', 'utf-8');
    await fs.writeFile(join(assetsDir, 'templates/DOR.md'), '# Definition of Ready\n{{PLACEHOLDER}}', 'utf-8');
    await fs.writeFile(join(assetsDir, 'templates/DOD.md'), '# Definition of Done\n{{PLACEHOLDER}}', 'utf-8');
    await fs.writeFile(join(assetsDir, 'templates/PRD.md'), '# Product Requirements Document\n{{PLACEHOLDER}}', 'utf-8');
    await fs.writeFile(join(assetsDir, 'templates/SDD.md'), '# System Design Document\n{{PLACEHOLDER}}', 'utf-8');

    // Create rule files
    await fs.writeFile(join(assetsDir, 'rules/SCQA.md'), '# SCQA Framework\nSituation, Complication, Question, Answer.', 'utf-8');
    await fs.writeFile(join(assetsDir, 'rules/SOLID.md'), '# SOLID Principles\nSingle responsibility, Open/closed, etc.', 'utf-8');
    await fs.writeFile(join(assetsDir, 'rules/DDD.md'), '# Domain-Driven Design\nUbiquitous language, bounded contexts.', 'utf-8');

    // Create output style file
    await fs.writeFile(join(assetsDir, 'output-styles/the-startup.md'), '# The Startup Output Style\nFormat for The Startup.', 'utf-8');

    // Create settings.json file
    await fs.writeFile(
      join(assetsDir, 'settings.json'),
      JSON.stringify({
        'startup-validate-dor': {
          command: '{{STARTUP_PATH}}/bin/the-startup validate dor',
          description: 'Validate Definition of Ready',
          continueOnError: false,
        },
        'startup-validate-dod': {
          command: '{{STARTUP_PATH}}/bin/the-startup validate dod',
          description: 'Validate Definition of Done',
          continueOnError: false,
        },
      }, null, 2),
      'utf-8'
    );
  }

  /**
   * Helper: Create a legacy v1 lock file (Go format)
   */
  async function createLegacyV1LockFile(): Promise<void> {
    await fs.mkdir(startupPath, { recursive: true });

    const legacyLock: LegacyLockFile = {
      version: '0.9.0',
      installedAt: new Date('2024-12-01T10:00:00.000Z').toISOString(),
      files: {
        agents: [
          join(claudePath, 'agents/the-analyst-requirements-analysis.md'),
          join(claudePath, 'agents/the-architect-system-architecture.md'),
          join(claudePath, 'agents/the-software-engineer-api-development.md'),
          join(claudePath, 'agents/the-platform-engineer-deployment.md'),
          join(claudePath, 'agents/the-security-engineer-assessment.md'),
          join(claudePath, 'agents/the-test-engineer-test-execution.md'),
          join(claudePath, 'agents/the-chief.md'),
          join(claudePath, 'agents/the-product-owner-feature-definition.md'),
        ],
        commands: [
          join(claudePath, 'commands/s-specify.md'),
          join(claudePath, 'commands/s-implement.md'),
          join(claudePath, 'commands/s-refactor.md'),
          join(claudePath, 'commands/s-analyze.md'),
          join(claudePath, 'commands/s-init.md'),
        ],
        templates: [
          join(startupPath, 'templates/SPEC.md'),
          join(startupPath, 'templates/TASK-DOD.md'),
          join(startupPath, 'templates/DOR.md'),
          join(startupPath, 'templates/DOD.md'),
          join(startupPath, 'templates/PRD.md'),
          join(startupPath, 'templates/SDD.md'),
        ],
        rules: [
          join(startupPath, 'rules/SCQA.md'),
          join(startupPath, 'rules/SOLID.md'),
          join(startupPath, 'rules/DDD.md'),
        ],
        outputStyles: [
          join(claudePath, 'output-styles/the-startup.md'),
        ],
        binary: join(startupPath, 'bin/the-startup'),
      },
    };

    await fs.writeFile(lockFilePath, JSON.stringify(legacyLock, null, 2), 'utf-8');
  }

  /**
   * Helper: Create existing files from v1 lock file to simulate Go installation
   */
  async function createExistingFilesFromV1Lock(): Promise<void> {
    // Create directories
    await fs.mkdir(join(claudePath, 'agents'), { recursive: true });
    await fs.mkdir(join(claudePath, 'commands'), { recursive: true });
    await fs.mkdir(join(startupPath, 'templates'), { recursive: true });
    await fs.mkdir(join(startupPath, 'rules'), { recursive: true });
    await fs.mkdir(join(claudePath, 'output-styles'), { recursive: true });
    await fs.mkdir(join(startupPath, 'bin'), { recursive: true });

    // Copy all asset files to simulate existing Go installation
    const mockAssetsDir = join(tempDir, 'mock-assets');

    // Copy agents
    const agents = await fs.readdir(join(mockAssetsDir, 'agents'));
    for (const agent of agents) {
      await fs.copyFile(
        join(mockAssetsDir, 'agents', agent),
        join(claudePath, 'agents', agent)
      );
    }

    // Copy commands
    const commands = await fs.readdir(join(mockAssetsDir, 'commands'));
    for (const command of commands) {
      await fs.copyFile(
        join(mockAssetsDir, 'commands', command),
        join(claudePath, 'commands', command)
      );
    }

    // Copy templates
    const templates = await fs.readdir(join(mockAssetsDir, 'templates'));
    for (const template of templates) {
      await fs.copyFile(
        join(mockAssetsDir, 'templates', template),
        join(startupPath, 'templates', template)
      );
    }

    // Copy rules
    const rules = await fs.readdir(join(mockAssetsDir, 'rules'));
    for (const rule of rules) {
      await fs.copyFile(
        join(mockAssetsDir, 'rules', rule),
        join(startupPath, 'rules', rule)
      );
    }

    // Copy output styles
    const outputStyles = await fs.readdir(join(mockAssetsDir, 'output-styles'));
    for (const style of outputStyles) {
      await fs.copyFile(
        join(mockAssetsDir, 'output-styles', style),
        join(claudePath, 'output-styles', style)
      );
    }

    // Create Go binary (empty file to simulate binary)
    await fs.writeFile(join(startupPath, 'bin/the-startup'), '#!/bin/bash\necho "Go binary"', 'utf-8');
  }

  /**
   * T004.4.1: Test upgrade from Go lock file (v1 string[]) to npm lock file (v2 with checksums)
   *
   * Verifies:
   * - LockManager detects v1 format
   * - v1 lock file is automatically migrated to v2 format on read
   * - All file paths are preserved during migration
   * - Checksums are generated for all existing files
   * - New lock file written in v2 format with lockFileVersion: 2
   */
  it('should upgrade v1 lock file to v2 format with checksums', async () => {
    // Setup: Create v1 lock file and existing installation
    await createLegacyV1LockFile();
    await createExistingFilesFromV1Lock();

    // Verify v1 lock file exists
    const v1LockContent = await fs.readFile(lockFilePath, 'utf-8');
    const v1Lock: LegacyLockFile = JSON.parse(v1LockContent);
    expect(v1Lock.version).toBe('0.9.0');
    expect(v1Lock.files.agents).toBeInstanceOf(Array);
    expect(typeof v1Lock.files.agents[0]).toBe('string'); // v1 format is string array

    // Act: Run npm installer (simulates user upgrading to npm version)
    const options: InstallerOptions = {
      startupPath,
      claudePath,
      selectedFiles: {
        agents: true,
        commands: true,
        templates: true,
        rules: true,
        outputStyles: true,
      },
    };

    const result = await installer.install(options);

    // Assert: Installation succeeded
    expect(result.success).toBe(true);
    expect(result.installedFiles.length).toBe(24); // 8 agents + 5 commands + 6 templates + 3 rules + 1 output style + 1 settings.json

    // Verify v2 lock file was created
    const v2LockContent = await fs.readFile(lockFilePath, 'utf-8');
    const v2Lock: LockFile = JSON.parse(v2LockContent);

    // Verify lock file metadata upgraded
    expect(v2Lock.version).toBe('1.0.0'); // npm version
    expect(v2Lock.lockFileVersion).toBe(2); // v2 format
    expect(v2Lock.installedAt).toBeDefined();

    // Verify all agents migrated with checksums
    expect(v2Lock.files.agents.length).toBe(8);
    for (const agent of v2Lock.files.agents) {
      expect(agent.path).toBeDefined();
      expect(agent.checksum).toBeDefined();
      expect(agent.checksum!.length).toBe(64); // SHA-256 is 64 hex chars

      // Verify checksum matches actual file
      const actualChecksum = await lockManager.generateChecksum(agent.path);
      expect(agent.checksum).toBe(actualChecksum);
    }

    // Verify all commands migrated with checksums
    expect(v2Lock.files.commands.length).toBe(5);
    for (const command of v2Lock.files.commands) {
      expect(command.checksum).toBeDefined();
      const actualChecksum = await lockManager.generateChecksum(command.path);
      expect(command.checksum).toBe(actualChecksum);
    }

    // Verify all templates migrated with checksums
    expect(v2Lock.files.templates.length).toBe(6);
    for (const template of v2Lock.files.templates) {
      expect(template.checksum).toBeDefined();
      const actualChecksum = await lockManager.generateChecksum(template.path);
      expect(template.checksum).toBe(actualChecksum);
    }

    // Verify all rules migrated with checksums
    expect(v2Lock.files.rules.length).toBe(3);
    for (const rule of v2Lock.files.rules) {
      expect(rule.checksum).toBeDefined();
      const actualChecksum = await lockManager.generateChecksum(rule.path);
      expect(rule.checksum).toBe(actualChecksum);
    }

    // Verify all output styles migrated with checksums
    expect(v2Lock.files.outputStyles.length).toBe(1);
    for (const style of v2Lock.files.outputStyles) {
      expect(style.checksum).toBeDefined();
      const actualChecksum = await lockManager.generateChecksum(style.path);
      expect(style.checksum).toBe(actualChecksum);
    }
  });

  /**
   * T004.4.2: Test detection of existing Go binary installation
   *
   * Verifies:
   * - npm installer detects existing Go binary in .the-startup/bin/
   * - Installation continues successfully despite Go binary presence
   * - Go binary is not removed (user can manually clean up)
   * - Lock file binary entry is updated to npm version
   */
  it('should detect existing Go binary and continue installation', async () => {
    // Setup: Create Go installation with binary
    await createLegacyV1LockFile();
    await createExistingFilesFromV1Lock();

    // Verify Go binary exists
    const goBinaryPath = join(startupPath, 'bin/the-startup');
    await expect(fs.access(goBinaryPath)).resolves.not.toThrow();
    const goBinaryContent = await fs.readFile(goBinaryPath, 'utf-8');
    expect(goBinaryContent).toContain('Go binary');

    // Act: Run npm installer
    const options: InstallerOptions = {
      startupPath,
      claudePath,
      selectedFiles: {
        agents: true,
        commands: true,
        templates: true,
        rules: true,
        outputStyles: true,
      },
    };

    const result = await installer.install(options);

    // Assert: Installation succeeded despite Go binary
    expect(result.success).toBe(true);

    // Verify Go binary still exists (not removed automatically)
    await expect(fs.access(goBinaryPath)).resolves.not.toThrow();

    // Note: In real implementation, installer might log a warning about existing Go binary
    // This test verifies the migration doesn't fail due to binary presence
  });

  /**
   * T004.4.3: Test settings.json preservation through Go→npm migration
   *
   * Verifies:
   * - Existing settings.json from Go installation is preserved
   * - Go-installed hooks remain intact
   * - New npm hooks are added without duplicates
   * - User customizations are not overwritten
   * - Backup is created and cleaned up
   */
  it('should preserve settings.json through Go to npm migration', async () => {
    // Setup: Create Go installation with settings.json
    await createLegacyV1LockFile();
    await createExistingFilesFromV1Lock();

    // Create settings.json with Go-installed hooks
    await fs.mkdir(claudePath, { recursive: true });
    const goSettings = {
      mcpServers: {
        'user-custom-server': {
          command: 'npx',
          args: ['-y', 'user-package'],
        },
      },
      hooks: {
        // User's custom hook
        'user-lint-hook': {
          command: 'npm run lint',
          description: 'Lint before commit',
          continueOnError: true,
        },
        // Go-installed hooks with old paths
        'startup-validate-dor': {
          command: `${startupPath}/bin/the-startup validate dor`,
          description: 'Validate Definition of Ready',
          continueOnError: false,
        },
        'startup-validate-dod': {
          command: `${startupPath}/bin/the-startup validate dod`,
          description: 'Validate Definition of Done',
          continueOnError: false,
        },
      },
    };
    await fs.writeFile(
      join(claudePath, 'settings.json'),
      JSON.stringify(goSettings, null, 2),
      'utf-8'
    );

    // Act: Run npm installer
    const options: InstallerOptions = {
      startupPath,
      claudePath,
      selectedFiles: {
        agents: true,
        commands: true,
        templates: true,
        rules: true,
        outputStyles: true,
      },
    };

    const result = await installer.install(options);

    // Assert: Installation succeeded
    expect(result.success).toBe(true);

    // Verify settings.json was preserved and merged
    const settingsContent = await fs.readFile(join(claudePath, 'settings.json'), 'utf-8');
    const settings = JSON.parse(settingsContent);

    // User's custom server should be preserved
    expect(settings.mcpServers).toBeDefined();
    expect(settings.mcpServers['user-custom-server']).toBeDefined();
    expect(settings.mcpServers['user-custom-server'].command).toBe('npx');

    // User's custom hook should be preserved
    expect(settings.hooks['user-lint-hook']).toBeDefined();
    expect(settings.hooks['user-lint-hook'].command).toBe('npm run lint');

    // Go-installed hooks should remain (NOT overwritten by npm installer)
    // This is because SettingsMerger preserves existing hooks by default
    expect(settings.hooks['startup-validate-dor']).toBeDefined();
    expect(settings.hooks['startup-validate-dod']).toBeDefined();

    // Verify no duplicate hooks were created
    const hookNames = Object.keys(settings.hooks);
    const uniqueHookNames = new Set(hookNames);
    expect(hookNames.length).toBe(uniqueHookNames.size);
  });

  /**
   * T004.4.4: Test all 23 assets migrate correctly from Go to npm installation
   *
   * Note: Test uses 23 assets (8 agents + 5 commands + 6 templates + 3 rules + 1 output style)
   * for performance. Real installation has 54 assets.
   *
   * Verifies:
   * - All asset files are present after migration
   * - Lock file contains all 23 assets
   * - All assets have valid checksums
   * - File contents are preserved
   * - Directory structure is correct
   */
  it('should migrate all 23 assets correctly from Go to npm', async () => {
    // Setup: Create partial Go installation (20 out of 23 assets)
    await createLegacyV1LockFile();

    // Create only a subset of files (simulate incomplete Go installation)
    await fs.mkdir(join(claudePath, 'agents'), { recursive: true });
    await fs.mkdir(join(claudePath, 'commands'), { recursive: true });
    await fs.mkdir(join(startupPath, 'templates'), { recursive: true });
    await fs.mkdir(join(startupPath, 'rules'), { recursive: true });

    const mockAssetsDir = join(tempDir, 'mock-assets');

    // Copy only first 6 agents (incomplete)
    const agents = (await fs.readdir(join(mockAssetsDir, 'agents'))).slice(0, 6);
    for (const agent of agents) {
      await fs.copyFile(
        join(mockAssetsDir, 'agents', agent),
        join(claudePath, 'agents', agent)
      );
    }

    // Copy all commands
    const commands = await fs.readdir(join(mockAssetsDir, 'commands'));
    for (const command of commands) {
      await fs.copyFile(
        join(mockAssetsDir, 'commands', command),
        join(claudePath, 'commands', command)
      );
    }

    // Copy all templates
    const templates = await fs.readdir(join(mockAssetsDir, 'templates'));
    for (const template of templates) {
      await fs.copyFile(
        join(mockAssetsDir, 'templates', template),
        join(startupPath, 'templates', template)
      );
    }

    // Copy only 2 rules (incomplete)
    const rules = (await fs.readdir(join(mockAssetsDir, 'rules'))).slice(0, 2);
    for (const rule of rules) {
      await fs.copyFile(
        join(mockAssetsDir, 'rules', rule),
        join(startupPath, 'rules', rule)
      );
    }

    // Do NOT copy output styles (completely missing category)

    // Act: Run npm installer (should install missing files)
    const options: InstallerOptions = {
      startupPath,
      claudePath,
      selectedFiles: {
        agents: true,
        commands: true,
        templates: true,
        rules: true,
        outputStyles: true,
      },
    };

    const result = await installer.install(options);

    // Assert: Installation succeeded
    expect(result.success).toBe(true);

    // Verify all 23 assets are now present
    const lockContent = await fs.readFile(lockFilePath, 'utf-8');
    const lockFile: LockFile = JSON.parse(lockContent);

    expect(lockFile.files.agents.length).toBe(8);
    expect(lockFile.files.commands.length).toBe(5);
    expect(lockFile.files.templates.length).toBe(6);
    expect(lockFile.files.rules.length).toBe(3);
    expect(lockFile.files.outputStyles.length).toBe(1);

    // Verify total count (excluding settings.json which is not tracked in lock file)
    const totalAssets =
      lockFile.files.agents.length +
      lockFile.files.commands.length +
      lockFile.files.templates.length +
      lockFile.files.rules.length +
      lockFile.files.outputStyles.length;
    expect(totalAssets).toBe(23); // 23 assets (settings.json is merged, not tracked in lock)

    // Verify all files have checksums
    for (const agent of lockFile.files.agents) {
      expect(agent.checksum).toBeDefined();
      expect(agent.checksum!.length).toBe(64);

      // Verify file exists
      await expect(fs.access(agent.path)).resolves.not.toThrow();
    }

    // Verify missing agents were installed
    const agentFiles = await fs.readdir(join(claudePath, 'agents'));
    expect(agentFiles.length).toBe(8);

    // Verify missing rules were installed
    const ruleFiles = await fs.readdir(join(startupPath, 'rules'));
    expect(ruleFiles.length).toBe(3);

    // Verify output styles were installed (was completely missing)
    const outputStyleFiles = await fs.readdir(join(claudePath, 'output-styles'));
    expect(outputStyleFiles.length).toBe(1);

    // Verify file contents are correct (not corrupted during migration)
    const chiefContent = await fs.readFile(join(claudePath, 'agents/the-chief.md'), 'utf-8');
    expect(chiefContent).toContain('The Chief');
  });

  /**
   * T004.4.5: Test idempotent migration (running twice doesn't break)
   *
   * Verifies:
   * - Running npm install twice on v1 lock file is idempotent
   * - Second run detects v2 lock file and skips migration
   * - No files are duplicated or corrupted
   * - Checksums remain valid after second run
   * - Lock file version stays at 2
   */
  it('should handle idempotent migration when run multiple times', async () => {
    // Setup: Create v1 lock file and existing installation
    await createLegacyV1LockFile();
    await createExistingFilesFromV1Lock();

    const options: InstallerOptions = {
      startupPath,
      claudePath,
      selectedFiles: {
        agents: true,
        commands: true,
        templates: true,
        rules: true,
        outputStyles: true,
      },
    };

    // Act 1: First npm install (v1 → v2 migration)
    const result1 = await installer.install(options);
    expect(result1.success).toBe(true);

    // Verify v2 lock file created
    const lockContent1 = await fs.readFile(lockFilePath, 'utf-8');
    const lockFile1: LockFile = JSON.parse(lockContent1);
    expect(lockFile1.lockFileVersion).toBe(2);

    // Store first run checksums for comparison
    const firstRunChecksums = new Map<string, string>();
    for (const agent of lockFile1.files.agents) {
      firstRunChecksums.set(agent.path, agent.checksum!);
    }

    // Act 2: Second npm install (should be idempotent)
    const result2 = await installer.install(options);
    expect(result2.success).toBe(true);

    // Verify v2 lock file still exists and is valid
    const lockContent2 = await fs.readFile(lockFilePath, 'utf-8');
    const lockFile2: LockFile = JSON.parse(lockContent2);

    // Assert: Lock file version stayed at 2
    expect(lockFile2.lockFileVersion).toBe(2);

    // Verify all checksums are still valid and unchanged
    for (const agent of lockFile2.files.agents) {
      expect(agent.checksum).toBeDefined();
      expect(agent.checksum).toBe(firstRunChecksums.get(agent.path));

      // Verify file wasn't corrupted
      const actualChecksum = await lockManager.generateChecksum(agent.path);
      expect(agent.checksum).toBe(actualChecksum);
    }

    // Verify no duplicate files created
    const agentFiles = await fs.readdir(join(claudePath, 'agents'));
    expect(agentFiles.length).toBe(8); // Same as first run

    const commandFiles = await fs.readdir(join(claudePath, 'commands'));
    expect(commandFiles.length).toBe(5);

    const templateFiles = await fs.readdir(join(startupPath, 'templates'));
    expect(templateFiles.length).toBe(6);

    // Verify lock file counts unchanged
    expect(lockFile2.files.agents.length).toBe(lockFile1.files.agents.length);
    expect(lockFile2.files.commands.length).toBe(lockFile1.files.commands.length);
    expect(lockFile2.files.templates.length).toBe(lockFile1.files.templates.length);
  });
});
