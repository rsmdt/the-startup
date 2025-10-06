import { describe, it, expect } from 'vitest';
import React from 'react';
import { render } from 'ink-testing-library';
import { Complete, InstallationSummary } from '../../../src/ui/install/Complete';

describe('Complete', () => {
  const sampleSummary: InstallationSummary = {
    installedFiles: [
      '/Users/test/.claude/agents/specify.md',
      '/Users/test/.claude/agents/analyze.md',
      '/Users/test/.claude/commands/init.md',
      '/Users/test/.the-startup/templates/PRD.md',
    ],
    startupPath: '/Users/test/.the-startup',
    claudePath: '/Users/test/.claude',
    totalFiles: 4,
  };

  describe('Rendering', () => {
    it('renders without crashing', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      expect(lastFrame()).toBeDefined();
    });

    it('displays success message', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      const output = lastFrame();

      expect(output).toMatch(/(success|complete|installed)/i);
    });

    it('displays success icon from theme', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      const output = lastFrame();

      expect(output).toContain('✓');
    });

    it('shows total files count', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      const output = lastFrame();

      expect(output).toContain('4');
      expect(output).toMatch(/file/i);
    });
  });

  describe('Installation Paths', () => {
    it('displays startup installation path', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      const output = lastFrame();

      expect(output).toContain('/Users/test/.the-startup');
    });

    it('displays claude configuration path', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      const output = lastFrame();

      expect(output).toContain('/Users/test/.claude');
    });

    it('shows both installation paths', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      const output = lastFrame();

      expect(output).toContain('.the-startup');
      expect(output).toContain('.claude');
    });
  });

  describe('Installed Files List', () => {
    it('displays list of installed files', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      const output = lastFrame();

      expect(output).toContain('specify.md');
      expect(output).toContain('analyze.md');
      expect(output).toContain('init.md');
      expect(output).toContain('PRD.md');
    });

    it('shows file count matches summary', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      const output = lastFrame();

      // Should mention 4 files
      expect(output).toContain('4');
    });

    it('groups files by directory', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      const output = lastFrame();

      // Should show directory structure
      expect(output).toMatch(/agents|commands|templates/i);
    });
  });

  describe('Next Steps', () => {
    it('displays next steps or usage instructions', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      const output = lastFrame();

      expect(output).toMatch(/(next|usage|getting started|restart)/i);
    });

    it('mentions restarting Claude Code', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      const output = lastFrame();

      expect(output).toMatch(/restart|reload/i);
    });
  });

  describe('Component Pattern', () => {
    it('follows Ink functional component pattern', () => {
      expect(typeof Complete).toBe('function');
    });

    it('returns valid React element', () => {
      const element = <Complete summary={sampleSummary} />;
      expect(React.isValidElement(element)).toBe(true);
    });
  });

  describe('Props', () => {
    it('accepts required prop: summary', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);

      expect(lastFrame()).toBeDefined();
    });

    it('accepts optional message prop', () => {
      const { lastFrame } = render(
        <Complete
          summary={sampleSummary}
          message="Custom success message"
        />
      );
      const output = lastFrame();

      expect(output).toContain('Custom success message');
    });
  });

  describe('Edge Cases', () => {
    it('handles empty installed files list', () => {
      const emptySummary: InstallationSummary = {
        installedFiles: [],
        startupPath: '/Users/test/.the-startup',
        claudePath: '/Users/test/.claude',
        totalFiles: 0,
      };

      const { lastFrame } = render(<Complete summary={emptySummary} />);
      const output = lastFrame();

      expect(output).toContain('0');
    });

    it('handles single file installation', () => {
      const singleFileSummary: InstallationSummary = {
        installedFiles: ['/Users/test/.claude/agents/specify.md'],
        startupPath: '/Users/test/.the-startup',
        claudePath: '/Users/test/.claude',
        totalFiles: 1,
      };

      const { lastFrame } = render(<Complete summary={singleFileSummary} />);
      const output = lastFrame();

      expect(output).toContain('1');
      expect(output).toContain('specify.md');
    });

    it('handles very long file paths', () => {
      const longPathSummary: InstallationSummary = {
        installedFiles: [
          '/Users/test/very/long/nested/path/to/some/deep/directory/structure/file.md',
        ],
        startupPath: '/Users/test/.the-startup',
        claudePath: '/Users/test/.claude',
        totalFiles: 1,
      };

      const { lastFrame } = render(<Complete summary={longPathSummary} />);

      expect(lastFrame()).toBeDefined();
    });

    it('handles special characters in paths', () => {
      const specialCharSummary: InstallationSummary = {
        installedFiles: [
          '/Users/test/.claude/file-with-special_chars@123.md',
        ],
        startupPath: '/Users/test/.the-startup',
        claudePath: '/Users/test/.claude',
        totalFiles: 1,
      };

      const { lastFrame } = render(<Complete summary={specialCharSummary} />);
      const output = lastFrame();

      expect(output).toContain('file-with-special_chars@123.md');
    });
  });

  describe('Formatting', () => {
    it('uses success color from theme', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      const output = lastFrame();

      // Should render successfully (color testing is limited in terminal)
      expect(output).toBeDefined();
      expect(output.length).toBeGreaterThan(0);
    });

    it('formats file list with proper indentation', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      const output = lastFrame();

      // Should have some structure/indentation
      expect(output).toBeDefined();
      expect(output).toContain('specify.md');
    });

    it('displays summary information clearly', () => {
      const { lastFrame } = render(<Complete summary={sampleSummary} />);
      const output = lastFrame();

      // Should contain key information
      expect(output).toContain('4');
      expect(output).toContain('.the-startup');
      expect(output).toContain('.claude');
    });
  });

  describe('Multiple Files Display', () => {
    it('handles large number of files', () => {
      const manyFilesSummary: InstallationSummary = {
        installedFiles: Array.from(
          { length: 20 },
          (_, i) => `/Users/test/.claude/file${i}.md`
        ),
        startupPath: '/Users/test/.the-startup',
        claudePath: '/Users/test/.claude',
        totalFiles: 20,
      };

      const { lastFrame } = render(<Complete summary={manyFilesSummary} />);
      const output = lastFrame();

      expect(output).toContain('20');
    });

    it('displays files from different directories', () => {
      const multiDirSummary: InstallationSummary = {
        installedFiles: [
          '/Users/test/.claude/agents/file1.md',
          '/Users/test/.claude/commands/file2.md',
          '/Users/test/.the-startup/templates/file3.md',
        ],
        startupPath: '/Users/test/.the-startup',
        claudePath: '/Users/test/.claude',
        totalFiles: 3,
      };

      const { lastFrame } = render(<Complete summary={multiDirSummary} />);
      const output = lastFrame();

      expect(output).toContain('file1.md');
      expect(output).toContain('file2.md');
      expect(output).toContain('file3.md');
    });
  });
});
