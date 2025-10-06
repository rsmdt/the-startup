import { describe, it, expect } from 'vitest';
import React from 'react';
import { render } from 'ink-testing-library';
import { FileTree, TreeNode } from '../../../src/ui/install/FileTree';

describe('FileTree', () => {
  const sampleTree: TreeNode = {
    name: 'root',
    type: 'directory',
    selected: true,
    expanded: true,
    children: [
      {
        name: 'agents',
        type: 'directory',
        selected: true,
        expanded: true,
        children: [
          {
            name: 'specify.md',
            type: 'file',
            selected: true,
          },
          {
            name: 'analyze.md',
            type: 'file',
            selected: true,
          },
        ],
      },
      {
        name: 'commands',
        type: 'directory',
        selected: true,
        expanded: false,
        children: [
          {
            name: 'init.md',
            type: 'file',
            selected: true,
          },
        ],
      },
      {
        name: 'templates',
        type: 'directory',
        selected: true,
        expanded: true,
        children: [
          {
            name: 'PRD.md',
            type: 'file',
            selected: false,
          },
        ],
      },
    ],
  };

  describe('Rendering', () => {
    it('renders without crashing', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      expect(lastFrame()).toBeDefined();
    });

    it('displays tree structure with directories and files', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      expect(output).toContain('agents');
      expect(output).toContain('commands');
      expect(output).toContain('templates');
    });

    it('shows expanded directories with children', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      // Expanded directory should show children
      expect(output).toContain('specify.md');
      expect(output).toContain('analyze.md');
    });

    it('hides collapsed directory children', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      // Collapsed directory should not show children
      // 'init.md' is in collapsed 'commands' directory
      expect(output).not.toContain('init.md');
    });

    it('displays selection indicators for selected items', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      // Should show selection indicators (e.g., checkboxes)
      // Using regex to match checkbox patterns
      expect(output).toMatch(/\[[\sxX✓]\]/);
    });

    it('displays folder icons for directories', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      // Should contain folder icons
      expect(output).toMatch(/[📁📂]/);
    });

    it('displays file icons for files', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      // Should contain file icons
      expect(output).toMatch(/[📄]/);
    });
  });

  describe('Selection State', () => {
    it('shows selected items with checked indicator', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      // Selected items should have checked indicator
      expect(output).toContain('✓');
    });

    it('shows unselected items with unchecked indicator', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      // Should show both checked and unchecked states
      // PRD.md is unselected
      expect(output).toContain('[ ]');
    });

    it('displays current selection count', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      // Should display count of selected items
      expect(output).toMatch(/\d+\s+(selected|item)/i);
    });
  });

  describe('Keyboard Navigation', () => {
    it('accepts keyboard navigation instructions', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      // Should show keyboard navigation hints
      expect(output).toMatch(/(arrow|↑|↓|j|k)/i);
    });

    it('shows space for selection toggle hint', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      // Should mention space key for selection
      expect(output).toMatch(/(space|select)/i);
    });

    it('shows enter for submission hint', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      // Should mention enter key for submission
      expect(output).toMatch(/(enter|confirm|submit)/i);
    });
  });

  describe('Component Pattern', () => {
    it('follows Ink functional component pattern', () => {
      expect(typeof FileTree).toBe('function');
    });

    it('returns valid React element', () => {
      const element = <FileTree tree={sampleTree} onSubmit={() => {}} />;
      expect(React.isValidElement(element)).toBe(true);
    });
  });

  describe('Props', () => {
    it('accepts required props: tree, onSubmit', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );

      expect(lastFrame()).toBeDefined();
    });

    it('accepts optional title prop', () => {
      const { lastFrame } = render(
        <FileTree
          tree={sampleTree}
          title="Select components to install"
          onSubmit={() => {}}
        />
      );
      const output = lastFrame();

      expect(output).toContain('Select components to install');
    });
  });

  describe('Tree Structure', () => {
    it('handles empty tree', () => {
      const emptyTree: TreeNode = {
        name: 'root',
        type: 'directory',
        selected: false,
        expanded: true,
        children: [],
      };

      const { lastFrame } = render(
        <FileTree tree={emptyTree} onSubmit={() => {}} />
      );

      expect(lastFrame()).toBeDefined();
    });

    it('handles single file', () => {
      const singleFile: TreeNode = {
        name: 'README.md',
        type: 'file',
        selected: true,
      };

      const { lastFrame } = render(
        <FileTree tree={singleFile} onSubmit={() => {}} />
      );
      const output = lastFrame();

      expect(output).toContain('README.md');
    });

    it('handles deeply nested tree', () => {
      const deepTree: TreeNode = {
        name: 'root',
        type: 'directory',
        selected: true,
        expanded: true,
        children: [
          {
            name: 'level1',
            type: 'directory',
            selected: true,
            expanded: true,
            children: [
              {
                name: 'level2',
                type: 'directory',
                selected: true,
                expanded: true,
                children: [
                  {
                    name: 'deep-file.md',
                    type: 'file',
                    selected: true,
                  },
                ],
              },
            ],
          },
        ],
      };

      const { lastFrame } = render(
        <FileTree tree={deepTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      expect(output).toContain('level1');
      expect(output).toContain('level2');
      expect(output).toContain('deep-file.md');
    });
  });

  describe('Visual Indicators', () => {
    it('shows indentation for nested items', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      // Nested items should have visual indentation
      // This is hard to test directly, but we can verify the structure exists
      expect(output).toBeDefined();
      expect(output.length).toBeGreaterThan(0);
    });

    it('uses theme colors for different item types', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      // Component should render (colors not easily testable in terminal)
      expect(output).toBeDefined();
    });

    it('highlights current cursor position', () => {
      const { lastFrame } = render(
        <FileTree tree={sampleTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      // Should have some form of cursor indicator
      expect(output).toMatch(/[>›▸]/);
    });
  });

  describe('Edge Cases', () => {
    it('handles tree with all items unselected', () => {
      const unselectedTree: TreeNode = {
        name: 'root',
        type: 'directory',
        selected: false,
        expanded: true,
        children: [
          {
            name: 'file.md',
            type: 'file',
            selected: false,
          },
        ],
      };

      const { lastFrame } = render(
        <FileTree tree={unselectedTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      expect(output).toContain('0 selected');
    });

    it('handles very long file names', () => {
      const longNameTree: TreeNode = {
        name: 'a-very-long-file-name-that-might-wrap-or-truncate.md',
        type: 'file',
        selected: true,
      };

      const { lastFrame } = render(
        <FileTree tree={longNameTree} onSubmit={() => {}} />
      );

      expect(lastFrame()).toBeDefined();
    });

    it('handles special characters in names', () => {
      const specialTree: TreeNode = {
        name: 'file-with-special_chars@123.md',
        type: 'file',
        selected: true,
      };

      const { lastFrame } = render(
        <FileTree tree={specialTree} onSubmit={() => {}} />
      );
      const output = lastFrame();

      expect(output).toContain('file-with-special_chars@123.md');
    });
  });
});
