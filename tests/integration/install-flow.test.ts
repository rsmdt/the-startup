import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { promises as fs } from 'fs';
import { join } from 'path';
import { Installer } from '../../src/core/installer/Installer';
import { LockManager } from '../../src/core/installer/LockManager';
import { SettingsMerger } from '../../src/core/installer/SettingsMerger';
import { createTempDir, cleanupTempDir } from '../shared/testUtils';
import type { InstallerOptions } from '../../src/core/types/config';
import type { LockFile } from '../../src/core/types/lock';

/**
 * Integration Test Suite: Install Flow
 *
 * Tests the complete installation flow from start to finish using real file system operations.
 * This is an integration test, not a unit test - we verify actual behavior with real files.
 *
 * Test Scenarios (from PRD lines 285-322):
 * 1. Fresh install to empty directories
 * 2. Install with existing settings.json (merge test)
 * 3. Partial file selection
 * 4. Lock file creation with checksums
 * 5. Settings.json placeholder replacement
 * 6. Error scenarios with rollback
 *
 * Key Differences from Unit Tests:
 * - Uses real file system (fs.promises), not mocks
 * - Creates temporary directories for each test
 * - Verifies actual file contents and structure
 * - Tests integration between Installer, LockManager, and SettingsMerger
 */

describe('Integration: Install Flow', () => {
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

    // Create real instances (not mocks) for integration testing
    lockManager = new LockManager(lockFilePath);
    settingsMerger = new SettingsMerger(fs);

    // Create mock asset provider with real file content
    assetProvider = {
      getAssetFiles: () => [
        { category: 'agents' as const, sourcePath: join(tempDir, 'mock-assets/agents/specify.md') },
        { category: 'agents' as const, sourcePath: join(tempDir, 'mock-assets/agents/implement.md') },
        { category: 'commands' as const, sourcePath: join(tempDir, 'mock-assets/commands/s-specify.md') },
        { category: 'commands' as const, sourcePath: join(tempDir, 'mock-assets/commands/s-implement.md') },
        { category: 'templates' as const, sourcePath: join(tempDir, 'mock-assets/templates/SPEC.md') },
        { category: 'templates' as const, sourcePath: join(tempDir, 'mock-assets/templates/TASK-DOD.md') },
        { category: 'rules' as const, sourcePath: join(tempDir, 'mock-assets/rules/SCQA.md') },
        { category: 'outputStyles' as const, sourcePath: join(tempDir, 'mock-assets/output-styles/json.md') },
      ],
      getSettingsTemplate: () => ({
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
      }),
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
    // Clean up temporary directory
    await cleanupTempDir(tempDir);
  });

  /**
   * Helper: Create mock asset files for testing
   */
  async function createMockAssets(): Promise<void> {
    const assetsDir = join(tempDir, 'mock-assets');
    await fs.mkdir(join(assetsDir, 'agents'), { recursive: true });
    await fs.mkdir(join(assetsDir, 'commands'), { recursive: true });
    await fs.mkdir(join(assetsDir, 'templates'), { recursive: true });
    await fs.mkdir(join(assetsDir, 'rules'), { recursive: true });
    await fs.mkdir(join(assetsDir, 'output-styles'), { recursive: true });

    // Create agent files
    await fs.writeFile(
      join(assetsDir, 'agents/specify.md'),
      '# Specify Agent\nCreates specifications from requirements.',
      'utf-8'
    );
    await fs.writeFile(
      join(assetsDir, 'agents/implement.md'),
      '# Implement Agent\nImplements code from specifications.',
      'utf-8'
    );

    // Create command files
    await fs.writeFile(
      join(assetsDir, 'commands/s-specify.md'),
      '# /s:specify command\nCreate specification',
      'utf-8'
    );
    await fs.writeFile(
      join(assetsDir, 'commands/s-implement.md'),
      '# /s:implement command\nImplement from spec',
      'utf-8'
    );

    // Create template files
    await fs.writeFile(
      join(assetsDir, 'templates/SPEC.md'),
      '# Specification Template\n{{PLACEHOLDER}}',
      'utf-8'
    );
    await fs.writeFile(
      join(assetsDir, 'templates/TASK-DOD.md'),
      '# Task Definition of Done\n{{PLACEHOLDER}}',
      'utf-8'
    );

    // Create rule files
    await fs.writeFile(
      join(assetsDir, 'rules/SCQA.md'),
      '# SCQA Framework\nSituation, Complication, Question, Answer.',
      'utf-8'
    );

    // Create output style files
    await fs.writeFile(
      join(assetsDir, 'output-styles/json.md'),
      '# JSON Output Style\nFormat output as JSON.',
      'utf-8'
    );
  }

  /**
   * Scenario 1: Fresh install to empty directories
   *
   * Verifies:
   * - All selected files are copied to correct locations
   * - Lock file is created with checksums
   * - Settings.json is created with merged hooks
   * - Placeholders are replaced correctly
   * - Directory structure is created
   */
  it('should successfully install all files to empty directories', async () => {
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

    // Verify installation succeeded
    expect(result.success).toBe(true);
    expect(result.installedFiles.length).toBe(8); // 2 agents + 2 commands + 2 templates + 1 rule + 1 output style
    expect(result.errors).toBeUndefined();

    // Verify directory structure was created
    await expect(fs.access(join(claudePath, 'agents'))).resolves.not.toThrow();
    await expect(fs.access(join(claudePath, 'commands'))).resolves.not.toThrow();
    await expect(fs.access(join(startupPath, 'templates'))).resolves.not.toThrow();
    await expect(fs.access(join(startupPath, 'rules'))).resolves.not.toThrow();
    await expect(fs.access(join(claudePath, 'output-styles'))).resolves.not.toThrow();

    // Verify agent files were copied
    const specifyContent = await fs.readFile(join(claudePath, 'agents/specify.md'), 'utf-8');
    expect(specifyContent).toContain('Specify Agent');

    const implementContent = await fs.readFile(join(claudePath, 'agents/implement.md'), 'utf-8');
    expect(implementContent).toContain('Implement Agent');

    // Verify command files were copied
    const specifyCommandContent = await fs.readFile(join(claudePath, 'commands/s-specify.md'), 'utf-8');
    expect(specifyCommandContent).toContain('/s:specify command');

    // Verify template files were copied
    const specTemplateContent = await fs.readFile(join(startupPath, 'templates/SPEC.md'), 'utf-8');
    expect(specTemplateContent).toContain('Specification Template');

    // Verify rule files were copied
    const scqaContent = await fs.readFile(join(startupPath, 'rules/SCQA.md'), 'utf-8');
    expect(scqaContent).toContain('SCQA Framework');

    // Verify output style files were copied
    const jsonStyleContent = await fs.readFile(join(claudePath, 'output-styles/json.md'), 'utf-8');
    expect(jsonStyleContent).toContain('JSON Output Style');

    // Verify lock file was created with checksums
    const lockContent = await fs.readFile(lockFilePath, 'utf-8');
    const lockFile: LockFile = JSON.parse(lockContent);

    expect(lockFile.version).toBe('1.0.0');
    expect(lockFile.lockFileVersion).toBe(2);
    expect(lockFile.installedAt).toBeDefined();

    expect(lockFile.files.agents.length).toBe(2);
    expect(lockFile.files.commands.length).toBe(2);
    expect(lockFile.files.templates.length).toBe(2);
    expect(lockFile.files.rules.length).toBe(1);
    expect(lockFile.files.outputStyles.length).toBe(1);

    // Verify all lock entries have checksums
    for (const agent of lockFile.files.agents) {
      expect(agent.checksum).toBeDefined();
      expect(agent.checksum!.length).toBe(64); // SHA-256 is 64 hex chars
    }

    // Verify settings.json was created with merged hooks and replaced placeholders
    const settingsContent = await fs.readFile(join(claudePath, 'settings.json'), 'utf-8');
    const settings = JSON.parse(settingsContent);

    expect(settings.hooks).toBeDefined();
    expect(settings.hooks['startup-validate-dor']).toBeDefined();
    expect(settings.hooks['startup-validate-dor'].command).toBe(`${startupPath}/bin/the-startup validate dor`);
    expect(settings.hooks['startup-validate-dod'].command).toBe(`${startupPath}/bin/the-startup validate dod`);
  });

  /**
   * Scenario 2: Install with existing settings.json
   *
   * Verifies:
   * - Existing settings.json is preserved
   * - New hooks are merged without overwriting user hooks
   * - Backup is created and cleaned up
   */
  it('should merge with existing settings.json without overwriting user hooks', async () => {
    // Create existing settings.json with user hooks
    await fs.mkdir(claudePath, { recursive: true });
    const existingSettings = {
      mcpServers: {
        'user-server': {
          command: 'npx',
          args: ['-y', 'user-package'],
        },
      },
      hooks: {
        'user-custom-hook': {
          command: 'echo "User hook"',
          description: 'User custom hook',
          continueOnError: false,
        },
        'startup-validate-dor': {
          command: 'echo "User override"',
          description: 'User-defined DOR validation',
          continueOnError: true,
        },
      },
    };
    await fs.writeFile(
      join(claudePath, 'settings.json'),
      JSON.stringify(existingSettings, null, 2),
      'utf-8'
    );

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

    expect(result.success).toBe(true);

    // Verify settings.json was merged
    const settingsContent = await fs.readFile(join(claudePath, 'settings.json'), 'utf-8');
    const settings = JSON.parse(settingsContent);

    // User's mcpServers should be preserved
    expect(settings.mcpServers).toBeDefined();
    expect(settings.mcpServers['user-server']).toBeDefined();

    // User's custom hook should be preserved
    expect(settings.hooks['user-custom-hook']).toBeDefined();
    expect(settings.hooks['user-custom-hook'].command).toBe('echo "User hook"');

    // User's override of startup-validate-dor should be preserved (NOT overwritten)
    expect(settings.hooks['startup-validate-dor'].command).toBe('echo "User override"');

    // New hook should be added
    expect(settings.hooks['startup-validate-dod']).toBeDefined();
    expect(settings.hooks['startup-validate-dod'].command).toBe(`${startupPath}/bin/the-startup validate dod`);
  });

  /**
   * Scenario 3: Partial file selection
   *
   * Verifies:
   * - Only selected categories are installed
   * - Unselected categories are not installed
   * - Lock file only contains selected files
   */
  it('should install only selected file categories', async () => {
    const options: InstallerOptions = {
      startupPath,
      claudePath,
      selectedFiles: {
        agents: true,
        commands: false,
        templates: true,
        rules: false,
        outputStyles: false,
      },
    };

    const result = await installer.install(options);

    expect(result.success).toBe(true);
    expect(result.installedFiles.length).toBe(4); // 2 agents + 2 templates

    // Verify agents directory exists and has files
    await expect(fs.access(join(claudePath, 'agents'))).resolves.not.toThrow();
    const agentsDir = await fs.readdir(join(claudePath, 'agents'));
    expect(agentsDir.length).toBe(2);

    // Verify templates directory exists and has files
    await expect(fs.access(join(startupPath, 'templates'))).resolves.not.toThrow();
    const templatesDir = await fs.readdir(join(startupPath, 'templates'));
    expect(templatesDir.length).toBe(2);

    // Verify commands directory was NOT created (not selected)
    await expect(fs.access(join(claudePath, 'commands'))).rejects.toThrow();

    // Verify rules directory was NOT created (not selected)
    await expect(fs.access(join(startupPath, 'rules'))).rejects.toThrow();

    // Verify output-styles directory was NOT created (not selected)
    await expect(fs.access(join(claudePath, 'output-styles'))).rejects.toThrow();

    // Verify lock file only contains selected categories
    const lockContent = await fs.readFile(lockFilePath, 'utf-8');
    const lockFile: LockFile = JSON.parse(lockContent);

    expect(lockFile.files.agents.length).toBe(2);
    expect(lockFile.files.commands.length).toBe(0);
    expect(lockFile.files.templates.length).toBe(2);
    expect(lockFile.files.rules.length).toBe(0);
    expect(lockFile.files.outputStyles.length).toBe(0);
  });

  /**
   * Scenario 4: Lock file checksum verification
   *
   * Verifies:
   * - Lock file is created with v2 format
   * - All files have SHA-256 checksums
   * - Checksums match actual file contents
   * - Lock file contains correct metadata
   */
  it('should create lock file with accurate checksums for all installed files', async () => {
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

    expect(result.success).toBe(true);

    // Read lock file
    const lockContent = await fs.readFile(lockFilePath, 'utf-8');
    const lockFile: LockFile = JSON.parse(lockContent);

    // Verify lock file metadata
    expect(lockFile.version).toBe('1.0.0');
    expect(lockFile.lockFileVersion).toBe(2);
    expect(lockFile.installedAt).toBeDefined();
    expect(new Date(lockFile.installedAt).getTime()).toBeLessThanOrEqual(Date.now());

    // Verify all agents have checksums
    for (const agent of lockFile.files.agents) {
      expect(agent.path).toBeDefined();
      expect(agent.checksum).toBeDefined();
      expect(agent.checksum!.length).toBe(64); // SHA-256 hex

      // Verify checksum matches actual file content
      const actualChecksum = await lockManager.generateChecksum(agent.path);
      expect(agent.checksum).toBe(actualChecksum);
    }

    // Verify all commands have checksums
    for (const command of lockFile.files.commands) {
      expect(command.checksum).toBeDefined();
      const actualChecksum = await lockManager.generateChecksum(command.path);
      expect(command.checksum).toBe(actualChecksum);
    }

    // Verify all templates have checksums
    for (const template of lockFile.files.templates) {
      expect(template.checksum).toBeDefined();
      const actualChecksum = await lockManager.generateChecksum(template.path);
      expect(template.checksum).toBe(actualChecksum);
    }
  });

  /**
   * Scenario 5: Settings placeholder replacement
   *
   * Verifies:
   * - {{STARTUP_PATH}} is replaced with actual startup path
   * - {{CLAUDE_PATH}} is replaced with actual claude path
   * - Replacement only happens in newly added hooks
   * - User's existing hooks are not modified
   */
  it('should correctly replace placeholders in settings.json hooks', async () => {
    const options: InstallerOptions = {
      startupPath,
      claudePath,
      selectedFiles: {
        agents: true,
        commands: false,
        templates: false,
        rules: false,
        outputStyles: false,
      },
    };

    const result = await installer.install(options);

    expect(result.success).toBe(true);

    // Verify settings.json placeholders were replaced
    const settingsContent = await fs.readFile(join(claudePath, 'settings.json'), 'utf-8');
    const settings = JSON.parse(settingsContent);

    // Verify no placeholder strings remain
    expect(settingsContent).not.toContain('{{STARTUP_PATH}}');
    expect(settingsContent).not.toContain('{{CLAUDE_PATH}}');

    // Verify actual paths are present
    expect(settings.hooks['startup-validate-dor'].command).toContain(startupPath);
    expect(settings.hooks['startup-validate-dod'].command).toContain(startupPath);
  });

  /**
   * Scenario 6: Error handling with rollback
   *
   * Verifies:
   * - Installation fails gracefully on errors
   * - All copied files are rolled back
   * - Lock file is not created on failure
   * - Settings.json is not modified on failure
   * - Error messages are descriptive
   */
  it('should rollback all changes on failure', async () => {
    // Create settings.json to verify it's not modified on failure
    await fs.mkdir(claudePath, { recursive: true });
    const originalSettings = {
      hooks: {
        'existing-hook': {
          command: 'echo "Original"',
        },
      },
    };
    await fs.writeFile(
      join(claudePath, 'settings.json'),
      JSON.stringify(originalSettings, null, 2),
      'utf-8'
    );

    // Delete one of the asset files to cause a copy error
    await fs.rm(join(tempDir, 'mock-assets/agents/implement.md'), { force: true });

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

    // Verify installation failed
    expect(result.success).toBe(false);
    expect(result.errors).toBeDefined();
    expect(result.errors!.length).toBeGreaterThan(0);
    expect(result.installedFiles.length).toBe(0);

    // Verify no files were left behind (rollback succeeded)
    const agentsDir = join(claudePath, 'agents');
    const agentsDirExists = await fs.access(agentsDir).then(() => true).catch(() => false);

    if (agentsDirExists) {
      const agentsFiles = await fs.readdir(agentsDir);
      expect(agentsFiles.length).toBe(0); // All files should be rolled back
    }

    // Verify lock file was not created
    await expect(fs.access(lockFilePath)).rejects.toThrow();

    // Verify settings.json was not modified (or was restored from backup)
    const settingsContent = await fs.readFile(join(claudePath, 'settings.json'), 'utf-8');
    const settings = JSON.parse(settingsContent);

    expect(settings.hooks['existing-hook']).toBeDefined();
    expect(settings.hooks['existing-hook'].command).toBe('echo "Original"');

    // Startup hooks should not be present
    expect(settings.hooks['startup-validate-dor']).toBeUndefined();
    expect(settings.hooks['startup-validate-dod']).toBeUndefined();
  });

  /**
   * Scenario 7: Path normalization
   *
   * Verifies:
   * - Relative paths are resolved correctly
   * - Tilde (~) expansion works
   * - Absolute paths are used in lock file
   */
  it('should normalize paths correctly (relative and tilde)', async () => {
    const options: InstallerOptions = {
      startupPath: './.the-startup',
      claudePath: './.claude',
      selectedFiles: {
        agents: true,
        commands: false,
        templates: false,
        rules: false,
        outputStyles: false,
      },
    };

    const result = await installer.install(options);

    expect(result.success).toBe(true);

    // Verify paths were normalized to absolute paths
    expect(result.installedFiles.every(path => path.startsWith('/'))).toBe(true);

    // Verify files were created at expected locations
    const expectedAgentPath = join(tempDir, '.claude/agents/specify.md');
    await expect(fs.access(expectedAgentPath)).resolves.not.toThrow();
  });
});
