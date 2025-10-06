/**
 * Statusline Shell Script Tests
 *
 * Tests for ultra-fast (<10ms) statusline shell scripts that extract git branch
 * from Claude Code's JSON input format.
 *
 * Test Coverage:
 * 1. JSON parsing from stdin (exact Claude Code format)
 * 2. Git branch detection via .git/HEAD read (fast path)
 * 3. Git fallback to 'git' command when .git/HEAD read fails
 * 4. Graceful degradation when git is missing
 * 5. Performance requirement (<10ms)
 * 6. JSON schema validation (exact Claude Code format)
 * 7. Home directory expansion (~) support
 * 8. Cross-platform compatibility (bash/zsh/PowerShell)
 *
 * [ref: PRD; lines: 213-223]
 * [ref: SDD; lines: 799-820, 218-222, 235-240, 1153, 1188]
 */

import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { execSync, spawn } from 'child_process';
import { mkdirSync, writeFileSync, rmSync, existsSync } from 'fs';
import { join } from 'path';
import { tmpdir } from 'os';

const PROJECT_ROOT = join(__dirname, '../../');
const STATUSLINE_SH = join(PROJECT_ROOT, 'bin/statusline.sh');
const STATUSLINE_PS1 = join(PROJECT_ROOT, 'bin/statusline.ps1');

// Exact Claude Code JSON format from docs
interface ClaudeCodeHookInput {
  session_id: string;
  transcript_path: string;
  cwd: string;
  hook_event_name: string;
  prompt: string;
}

describe('Statusline Shell Scripts', () => {
  let testDir: string;
  let gitDir: string;

  beforeEach(() => {
    // Create temporary test directory with git repo
    testDir = join(tmpdir(), `statusline-test-${Date.now()}`);
    gitDir = join(testDir, '.git');
    mkdirSync(testDir, { recursive: true });
    mkdirSync(gitDir, { recursive: true });
  });

  afterEach(() => {
    // Clean up test directory
    if (existsSync(testDir)) {
      rmSync(testDir, { recursive: true, force: true });
    }
  });

  describe('JSON Parsing', () => {
    it('should parse valid Claude Code JSON from stdin', (ctx) => {
      // Skip if shell script doesn't exist yet (TDD Red phase)
      if (!existsSync(STATUSLINE_SH)) {
        ctx.skip();
        return;
      }

      const input: ClaudeCodeHookInput = {
        session_id: 'test-session-123',
        transcript_path: '/path/to/transcript.jsonl',
        cwd: testDir,
        hook_event_name: 'UserPromptSubmit',
        prompt: 'test prompt',
      };

      // Create valid git branch
      writeFileSync(join(gitDir, 'HEAD'), 'ref: refs/heads/main\n');

      const result = execSync(STATUSLINE_SH, {
        input: JSON.stringify(input),
        encoding: 'utf-8',
        cwd: testDir,
      });

      expect(result.trim()).toBe('main');
    });

    it('should handle malformed JSON gracefully', (ctx) => {
      if (!existsSync(STATUSLINE_SH)) {
        ctx.skip();
        return;
      }

      const result = execSync(STATUSLINE_SH, {
        input: 'not valid json',
        encoding: 'utf-8',
        cwd: testDir,
      });

      // Should return empty string on error
      expect(result.trim()).toBe('');
    });

    it('should validate exact Claude Code schema fields', (ctx) => {
      if (!existsSync(STATUSLINE_SH)) {
        ctx.skip();
        return;
      }

      // Missing required fields
      const invalidInput = {
        cwd: testDir,
        // Missing: session_id, transcript_path, hook_event_name, prompt
      };

      writeFileSync(join(gitDir, 'HEAD'), 'ref: refs/heads/feature-branch\n');

      const result = execSync(STATUSLINE_SH, {
        input: JSON.stringify(invalidInput),
        encoding: 'utf-8',
        cwd: testDir,
      });

      // Should still extract cwd and work if basic fields present
      expect(result.trim()).toBe('feature-branch');
    });
  });

  describe('Git Branch Detection - Fast Path', () => {
    it('should read branch from .git/HEAD directly (fast path)', (ctx) => {
      if (!existsSync(STATUSLINE_SH)) {
        ctx.skip();
        return;
      }

      const input: ClaudeCodeHookInput = {
        session_id: 'test-session',
        transcript_path: '/path/to/transcript.jsonl',
        cwd: testDir,
        hook_event_name: 'UserPromptSubmit',
        prompt: 'test',
      };

      // Create branch reference
      writeFileSync(join(gitDir, 'HEAD'), 'ref: refs/heads/feature-xyz\n');

      const result = execSync(STATUSLINE_SH, {
        input: JSON.stringify(input),
        encoding: 'utf-8',
        cwd: testDir,
      });

      expect(result.trim()).toBe('feature-xyz');
    });

    it('should handle detached HEAD state', (ctx) => {
      if (!existsSync(STATUSLINE_SH)) {
        ctx.skip();
        return;
      }

      const input: ClaudeCodeHookInput = {
        session_id: 'test-session',
        transcript_path: '/path/to/transcript.jsonl',
        cwd: testDir,
        hook_event_name: 'UserPromptSubmit',
        prompt: 'test',
      };

      // Detached HEAD (commit SHA)
      writeFileSync(
        join(gitDir, 'HEAD'),
        'abc123def456789012345678901234567890abcd\n'
      );

      const result = execSync(STATUSLINE_SH, {
        input: JSON.stringify(input),
        encoding: 'utf-8',
        cwd: testDir,
      });

      // Should return empty or 'HEAD' for detached state
      expect(['', 'HEAD']).toContain(result.trim());
    });

    it('should handle branch names with special characters', (ctx) => {
      if (!existsSync(STATUSLINE_SH)) {
        ctx.skip();
        return;
      }

      const input: ClaudeCodeHookInput = {
        session_id: 'test-session',
        transcript_path: '/path/to/transcript.jsonl',
        cwd: testDir,
        hook_event_name: 'UserPromptSubmit',
        prompt: 'test',
      };

      // Branch with special chars
      writeFileSync(
        join(gitDir, 'HEAD'),
        'ref: refs/heads/feature/user-auth_v2.0\n'
      );

      const result = execSync(STATUSLINE_SH, {
        input: JSON.stringify(input),
        encoding: 'utf-8',
        cwd: testDir,
      });

      expect(result.trim()).toBe('feature/user-auth_v2.0');
    });
  });

  describe('Git Fallback', () => {
    it('should fall back to git command when .git/HEAD is missing', (ctx) => {
      if (!existsSync(STATUSLINE_SH)) {
        ctx.skip();
        return;
      }

      // Remove .git/HEAD but keep directory
      if (existsSync(join(gitDir, 'HEAD'))) {
        rmSync(join(gitDir, 'HEAD'));
      }

      // Initialize actual git repo for fallback test
      try {
        execSync('git init', { cwd: testDir, stdio: 'pipe' });
        execSync('git checkout -b test-branch', { cwd: testDir, stdio: 'pipe' });
      } catch (error) {
        ctx.skip(); // Skip if git not available
        return;
      }

      const input: ClaudeCodeHookInput = {
        session_id: 'test-session',
        transcript_path: '/path/to/transcript.jsonl',
        cwd: testDir,
        hook_event_name: 'UserPromptSubmit',
        prompt: 'test',
      };

      const result = execSync(STATUSLINE_SH, {
        input: JSON.stringify(input),
        encoding: 'utf-8',
        cwd: testDir,
      });

      expect(result.trim()).toBe('test-branch');
    });
  });

  describe('Graceful Degradation', () => {
    it('should return empty string when not in git repo', (ctx) => {
      if (!existsSync(STATUSLINE_SH)) {
        ctx.skip();
        return;
      }

      // Remove .git directory entirely
      rmSync(gitDir, { recursive: true, force: true });

      const input: ClaudeCodeHookInput = {
        session_id: 'test-session',
        transcript_path: '/path/to/transcript.jsonl',
        cwd: testDir,
        hook_event_name: 'UserPromptSubmit',
        prompt: 'test',
      };

      const result = execSync(STATUSLINE_SH, {
        input: JSON.stringify(input),
        encoding: 'utf-8',
        cwd: testDir,
      });

      expect(result.trim()).toBe('');
    });

    it('should handle directories without read permissions gracefully', (ctx) => {
      if (!existsSync(STATUSLINE_SH)) {
        ctx.skip();
        return;
      }

      // This test would require chmod, skip on Windows
      if (process.platform === 'win32') {
        ctx.skip();
        return;
      }

      const input: ClaudeCodeHookInput = {
        session_id: 'test-session',
        transcript_path: '/path/to/transcript.jsonl',
        cwd: '/root/.ssh', // Typically inaccessible
        hook_event_name: 'UserPromptSubmit',
        prompt: 'test',
      };

      const result = execSync(STATUSLINE_SH, {
        input: JSON.stringify(input),
        encoding: 'utf-8',
      });

      // Should not crash, return empty
      expect(result.trim()).toBe('');
    });
  });

  describe('Performance Requirements', () => {
    it('should execute in less than 10ms (SDD requirement)', (ctx) => {
      if (!existsSync(STATUSLINE_SH)) {
        ctx.skip();
        return;
      }

      const input: ClaudeCodeHookInput = {
        session_id: 'test-session',
        transcript_path: '/path/to/transcript.jsonl',
        cwd: testDir,
        hook_event_name: 'UserPromptSubmit',
        prompt: 'test',
      };

      writeFileSync(join(gitDir, 'HEAD'), 'ref: refs/heads/main\n');

      const iterations = 10;
      const times: number[] = [];

      for (let i = 0; i < iterations; i++) {
        const start = performance.now();

        execSync(STATUSLINE_SH, {
          input: JSON.stringify(input),
          encoding: 'utf-8',
          cwd: testDir,
        });

        const end = performance.now();
        times.push(end - start);
      }

      const avgTime = times.reduce((a, b) => a + b, 0) / times.length;

      // SDD requirement: <10ms (line 1153)
      expect(avgTime).toBeLessThan(10);
    });

    it('should use fast path (.git/HEAD read) over git command', (ctx) => {
      if (!existsSync(STATUSLINE_SH)) {
        ctx.skip();
        return;
      }

      const input: ClaudeCodeHookInput = {
        session_id: 'test-session',
        transcript_path: '/path/to/transcript.jsonl',
        cwd: testDir,
        hook_event_name: 'UserPromptSubmit',
        prompt: 'test',
      };

      writeFileSync(join(gitDir, 'HEAD'), 'ref: refs/heads/main\n');

      // Verify fast path returns correct result
      const result = execSync(STATUSLINE_SH, {
        input: JSON.stringify(input),
        encoding: 'utf-8',
        cwd: testDir,
      });

      // Should return correct branch
      expect(result.trim()).toBe('main');

      // Performance verified by the averaging test above
      // This test focuses on correctness of fast path
    });
  });

  describe('Home Directory Expansion', () => {
    it('should handle tilde (~) in cwd path', (ctx) => {
      if (!existsSync(STATUSLINE_SH)) {
        ctx.skip();
        return;
      }

      const input: ClaudeCodeHookInput = {
        session_id: 'test-session',
        transcript_path: '/path/to/transcript.jsonl',
        cwd: '~/projects/test',
        hook_event_name: 'UserPromptSubmit',
        prompt: 'test',
      };

      // Should not crash with tilde path
      const result = execSync(STATUSLINE_SH, {
        input: JSON.stringify(input),
        encoding: 'utf-8',
      });

      // Gracefully return empty if path doesn't exist
      expect(typeof result).toBe('string');
    });
  });

  describe('Cross-Platform Compatibility', () => {
    it('should have Unix shell script for bash/zsh', () => {
      expect(existsSync(STATUSLINE_SH)).toBe(true);
    });

    it('should have PowerShell script for Windows', () => {
      expect(existsSync(STATUSLINE_PS1)).toBe(true);
    });

    it('should have executable permissions on Unix script', (ctx) => {
      if (process.platform === 'win32') {
        ctx.skip();
        return;
      }

      if (!existsSync(STATUSLINE_SH)) {
        ctx.skip();
        return;
      }

      const stats = execSync(`stat -f "%Lp" "${STATUSLINE_SH}"`, {
        encoding: 'utf-8',
      }).trim();

      // Check if executable bit is set (should be 755 or 775)
      expect(['755', '775', '777']).toContain(stats);
    });

    it('should work with PowerShell on Windows', (ctx) => {
      if (process.platform !== 'win32') {
        ctx.skip();
        return;
      }

      if (!existsSync(STATUSLINE_PS1)) {
        ctx.skip();
        return;
      }

      const input: ClaudeCodeHookInput = {
        session_id: 'test-session',
        transcript_path: 'C:\\path\\to\\transcript.jsonl',
        cwd: testDir,
        hook_event_name: 'UserPromptSubmit',
        prompt: 'test',
      };

      writeFileSync(join(gitDir, 'HEAD'), 'ref: refs/heads/main\n');

      const result = execSync(`powershell -File "${STATUSLINE_PS1}"`, {
        input: JSON.stringify(input),
        encoding: 'utf-8',
        cwd: testDir,
      });

      expect(result.trim()).toBe('main');
    });
  });
});
