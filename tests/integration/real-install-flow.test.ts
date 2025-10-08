import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { existsSync, readdirSync, readFileSync } from 'fs';
import { mkdir, rm, readdir, readFile } from 'fs/promises';
import { join } from 'path';
import { Installer } from '../../src/core/installer/Installer';
import { LockManager } from '../../src/core/installer/LockManager';
import { SettingsMerger } from '../../src/core/installer/SettingsMerger';
import { FileSystemAssetProvider } from '../../src/cli/asset-provider';
import type { InstallerOptions } from '../../src/core/types/config';
import type { LockFile } from '../../src/core/types/lock';
import { createInstallerFS, createSettingsMergerFS } from '../../src/cli/fs-adapter';

/**
 * REAL Integration Test Suite: Install Flow with Actual Files
 *
 * This test suite:
 * - Uses REAL assets from assets/ directory (not mocks!)
 * - Creates REAL files in tests/fixtures/
 * - Allows inspection of installed files after tests
 * - Only mocks at the edge (filesystem adapter), not internal logic
 *
 * Test Philosophy:
 * - "Only mock at the edge" - we mock filesystem operations but use real business logic
 * - All tests create real files that can be inspected
 * - Fixtures are preserved in tests/fixtures/ for manual verification
 */

describe('Real Integration: Install Flow with Actual Files', () => {
  const fixturesDir = join(process.cwd(), 'tests', 'fixtures');
  let testId: string;
  let testDir: string;
  let startupPath: string;
  let claudePath: string;
  let lockFilePath: string;

  beforeEach(async () => {
    // Create unique test directory for this test run
    testId = `test-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
    testDir = join(fixturesDir, testId);

    startupPath = join(testDir, '.the-startup');
    claudePath = join(testDir, '.claude');
    lockFilePath = join(startupPath, '.the-startup.lock');

    // Create test directories
    await mkdir(testDir, { recursive: true });
    await mkdir(startupPath, { recursive: true });
    await mkdir(claudePath, { recursive: true });
  });

  afterEach(async () => {
    // Keep fixtures for inspection - don't delete!
    // To clean up old fixtures manually: rm -rf tests/fixtures/test-*
    console.log(`\n📁 Test artifacts preserved in: ${testDir}`);
  });

  it('should perform complete installation with real assets from assets/ directory', async () => {
    // Create real instances using actual asset provider
    const assetProvider = new FileSystemAssetProvider();
    const lockManager = new LockManager(lockFilePath);
    const installerFS = createInstallerFS();
    const settingsMergerFS = createSettingsMergerFS();
    const settingsMerger = new SettingsMerger(settingsMergerFS);

    const installer = new Installer(
      installerFS,
      lockManager,
      settingsMerger,
      assetProvider,
      '1.0.0'
    );

    // Install all components
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
    expect(result.errors || []).toHaveLength(0);
    expect(result.installedFiles.length).toBeGreaterThan(0);

    // Verify agents directory exists and has files
    const agentsDir = join(claudePath, 'agents');
    expect(existsSync(agentsDir)).toBe(true);

    const agentFiles = await readdir(agentsDir);
    expect(agentFiles.length).toBeGreaterThan(0);

    // Verify specific agents exist (from actual assets)
    expect(existsSync(join(agentsDir, 'the-chief.md'))).toBe(true);
    expect(existsSync(join(agentsDir, 'the-meta-agent.md'))).toBe(true);

    // Verify activity files exist (flattened structure as expected by Claude Code)
    expect(existsSync(join(agentsDir, 'requirements-analysis.md'))).toBe(true);
    expect(existsSync(join(agentsDir, 'system-architecture.md'))).toBe(true);
    expect(existsSync(join(agentsDir, 'api-development.md'))).toBe(true);

    // Verify commands directory (flattened structure)
    const commandsDir = join(claudePath, 'commands');
    expect(existsSync(commandsDir)).toBe(true);

    const commandFiles = await readdir(commandsDir);
    expect(commandFiles).toContain('specify.md');
    expect(commandFiles).toContain('implement.md');
    expect(commandFiles).toContain('analyze.md');
    expect(commandFiles).toContain('refactor.md');
    expect(commandFiles).toContain('init.md');

    // Verify templates directory
    const templatesDir = join(startupPath, 'templates');
    expect(existsSync(templatesDir)).toBe(true);

    // Verify rules directory
    const rulesDir = join(startupPath, 'rules');
    expect(existsSync(rulesDir)).toBe(true);

    // Verify output styles
    const outputStylesDir = join(claudePath, 'output-styles');
    expect(existsSync(outputStylesDir)).toBe(true);
    expect(existsSync(join(outputStylesDir, 'the-startup.md'))).toBe(true);

    // Verify lock file was created with correct format
    expect(existsSync(lockFilePath)).toBe(true);

    const lockFileContent = await readFile(lockFilePath, 'utf-8');
    const lockFile: LockFile = JSON.parse(lockFileContent);

    expect(lockFile.version).toBe('1.0.0');
    expect(lockFile.lockFileVersion).toBe(2);
    expect(lockFile.installedAt).toBeDefined();
    expect(lockFile.files).toBeDefined();
    expect(lockFile.files.agents).toBeDefined();
    expect(lockFile.files.commands).toBeDefined();

    // Verify checksums exist for all files
    const allFiles = [
      ...lockFile.files.agents,
      ...lockFile.files.commands,
      ...lockFile.files.templates,
      ...lockFile.files.rules,
      ...lockFile.files.outputStyles,
    ];

    allFiles.forEach((file) => {
      expect(file.checksum).toBeDefined();
      expect(file.checksum).toMatch(/^[a-f0-9]{64}$/); // SHA-256 hash (64 hex chars)
      expect(file.path).toBeDefined();
    });

    // Verify settings.json was created
    const settingsPath = join(claudePath, 'settings.json');
    expect(existsSync(settingsPath)).toBe(true);

    const settingsContent = await readFile(settingsPath, 'utf-8');
    const settings = JSON.parse(settingsContent);

    expect(settings.hooks).toBeDefined();
    expect(settings.hooks['user-prompt-submit']).toBeDefined();

    console.log(`\n✅ Installed ${result.installedFiles.length} files`);
    console.log(`📁 Agents: ${agentFiles.length} files`);
    console.log(`📁 Commands: ${commandFiles.length} files`);
    console.log(`📄 Lock file: ${lockFilePath}`);
    console.log(`⚙️  Settings: ${settingsPath}`);
  });

  it('should correctly handle file content with real agent definitions', async () => {
    // Install only agents
    const assetProvider = new FileSystemAssetProvider();
    const lockManager = new LockManager(lockFilePath);
    const installerFS = createInstallerFS();
    const settingsMergerFS = createSettingsMergerFS();
    const settingsMerger = new SettingsMerger(settingsMergerFS);

    const installer = new Installer(
      installerFS,
      lockManager,
      settingsMerger,
      assetProvider,
      '1.0.0'
    );

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

    // Read and verify the-chief.md content
    const chiefPath = join(claudePath, 'agents', 'the-chief.md');
    expect(existsSync(chiefPath)).toBe(true);

    const chiefContent = readFileSync(chiefPath, 'utf-8');

    // Verify it's the actual agent definition (not mock data)
    expect(chiefContent).toContain('the-chief');
    expect(chiefContent).toContain('Core Responsibilities');
    expect(chiefContent.length).toBeGreaterThan(1000); // Real agent files are substantial

    // Read and verify activity file (api-development.md - flattened structure)
    const apiDevPath = join(claudePath, 'agents', 'api-development.md');
    expect(existsSync(apiDevPath)).toBe(true);

    const apiDevContent = readFileSync(apiDevPath, 'utf-8');
    expect(apiDevContent).toContain('API');
    expect(apiDevContent.length).toBeGreaterThan(500);

    console.log(`\n✅ Verified real agent content`);
    console.log(`📄 the-chief.md: ${chiefContent.length} bytes`);
    console.log(`📄 api-development.md: ${apiDevContent.length} bytes`);
  });

  it('should create inspectable directory structure matching expected layout', async () => {
    const assetProvider = new FileSystemAssetProvider();
    const lockManager = new LockManager(lockFilePath);
    const installerFS = createInstallerFS();
    const settingsMergerFS = createSettingsMergerFS();
    const settingsMerger = new SettingsMerger(settingsMergerFS);

    const installer = new Installer(
      installerFS,
      lockManager,
      settingsMerger,
      assetProvider,
      '1.0.0'
    );

    await installer.install({
      startupPath,
      claudePath,
      selectedFiles: {
        agents: true,
        commands: true,
        templates: true,
        rules: true,
        outputStyles: true,
      },
    });

    // Expected structure:
    // tests/fixtures/test-{id}/
    //   .the-startup/
    //     templates/
    //     rules/
    //     .the-startup.lock
    //   .claude/
    //     agents/
    //       the-chief.md
    //       the-meta-agent.md
    //       requirements-analysis.md
    //       api-development.md
    //       ... (all agents flattened)
    //     commands/
    //       s/
    //         specify.md
    //         implement.md
    //         ...
    //     output-styles/
    //       the-startup.md
    //     settings.json

    // Verify startup structure
    expect(existsSync(startupPath)).toBe(true);
    expect(existsSync(join(startupPath, 'templates'))).toBe(true);
    expect(existsSync(join(startupPath, 'rules'))).toBe(true);
    expect(existsSync(lockFilePath)).toBe(true);

    // Verify claude structure
    expect(existsSync(claudePath)).toBe(true);
    expect(existsSync(join(claudePath, 'agents'))).toBe(true);
    expect(existsSync(join(claudePath, 'commands'))).toBe(true);
    expect(existsSync(join(claudePath, 'output-styles'))).toBe(true);
    expect(existsSync(join(claudePath, 'settings.json'))).toBe(true);

    // Print directory tree for manual inspection
    console.log(`\n📁 Directory structure created in: ${testDir}`);
    console.log(`\nTo inspect:`);
    console.log(`  cd ${testDir}`);
    console.log(`  tree .`);
    console.log(`  cat .claude/agents/the-chief.md`);
    console.log(`  cat .claude/settings.json`);
    console.log(`  cat .the-startup/.the-startup.lock | jq`);
  });

  it('should handle partial installation (agents only) with real files', async () => {
    const assetProvider = new FileSystemAssetProvider();
    const lockManager = new LockManager(lockFilePath);
    const installerFS = createInstallerFS();
    const settingsMergerFS = createSettingsMergerFS();
    const settingsMerger = new SettingsMerger(settingsMergerFS);

    const installer = new Installer(
      installerFS,
      lockManager,
      settingsMerger,
      assetProvider,
      '1.0.0'
    );

    // Install only agents and commands
    const result = await installer.install({
      startupPath,
      claudePath,
      selectedFiles: {
        agents: true,
        commands: true,
        templates: false, // Don't install templates
        rules: false,     // Don't install rules
        outputStyles: false,
      },
    });

    expect(result.success).toBe(true);

    // Verify agents were installed
    expect(existsSync(join(claudePath, 'agents'))).toBe(true);
    expect(existsSync(join(claudePath, 'commands'))).toBe(true);

    // Verify templates were NOT installed
    expect(existsSync(join(startupPath, 'templates'))).toBe(false);
    expect(existsSync(join(startupPath, 'rules'))).toBe(false);
    expect(existsSync(join(claudePath, 'output-styles'))).toBe(false);

    // Verify lock file only includes installed categories
    const lockContent = await readFile(lockFilePath, 'utf-8');
    const lock: LockFile = JSON.parse(lockContent);

    expect(lock.files.agents.length).toBeGreaterThan(0);
    expect(lock.files.commands.length).toBeGreaterThan(0);
    expect(lock.files.templates.length).toBe(0);
    expect(lock.files.rules.length).toBe(0);
    expect(lock.files.outputStyles.length).toBe(0);

    console.log(`\n✅ Partial installation verified`);
    console.log(`📁 Installed: agents (${lock.files.agents.length}), commands (${lock.files.commands.length})`);
    console.log(`❌ Skipped: templates, rules, output-styles`);
  });
});
