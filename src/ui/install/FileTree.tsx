import { FC, useState } from 'react';
import { Box, Text, useInput } from 'ink';
import { theme } from '../shared/theme';

export interface TreeNode {
  name: string;
  type: 'file' | 'directory';
  selected: boolean;
  expanded?: boolean;
  children?: TreeNode[];
}

export interface FileTreeProps {
  /**
   * Tree structure to display
   */
  tree: TreeNode;

  /**
   * Optional title to display above the tree
   */
  title?: string;

  /**
   * Callback when user confirms selection
   */
  onSubmit: (tree: TreeNode) => void;
}

/**
 * Interactive file tree component with keyboard navigation
 * Supports arrow keys and vim bindings (hjkl) for navigation
 * Space to toggle selection, Enter to confirm
 *
 * @example
 * ```tsx
 * <FileTree
 *   tree={fileTreeData}
 *   title="Select components to install"
 *   onSubmit={(tree) => console.log('Selected:', tree)}
 * />
 * ```
 */
export const FileTree: FC<FileTreeProps> = ({ tree, title, onSubmit }) => {
  const [currentTree, setCurrentTree] = useState<TreeNode>(tree);
  const [cursorIndex, setCursorIndex] = useState(0);

  // Flatten tree for navigation
  const flattenTree = (node: TreeNode, depth = 0): Array<{ node: TreeNode; depth: number; path: string }> => {
    const result: Array<{ node: TreeNode; depth: number; path: string }> = [];

    const traverse = (n: TreeNode, d: number, parentPath: string) => {
      const currentPath = parentPath ? `${parentPath}/${n.name}` : n.name;
      result.push({ node: n, depth: d, path: currentPath });

      if (n.expanded && n.children) {
        n.children.forEach((child) => traverse(child, d + 1, currentPath));
      }
    };

    traverse(node, depth, '');
    return result;
  };

  const flatItems = flattenTree(currentTree);

  // Keyboard navigation
  useInput((input, key) => {
    // Navigation: Arrow keys and vim bindings
    if (key.downArrow || input === 'j') {
      setCursorIndex((prev) => Math.min(prev + 1, flatItems.length - 1));
    } else if (key.upArrow || input === 'k') {
      setCursorIndex((prev) => Math.max(prev - 1, 0));
    }

    // Expand/collapse: Arrow keys and vim bindings
    else if ((key.rightArrow || input === 'l') && flatItems[cursorIndex]) {
      const item = flatItems[cursorIndex];
      if (item.node.type === 'directory') {
        toggleExpanded(item.path);
      }
    } else if ((key.leftArrow || input === 'h') && flatItems[cursorIndex]) {
      const item = flatItems[cursorIndex];
      if (item.node.type === 'directory' && item.node.expanded) {
        toggleExpanded(item.path);
      }
    }

    // Toggle selection: Space
    else if (input === ' ' && flatItems[cursorIndex]) {
      toggleSelection(flatItems[cursorIndex].path);
    }

    // Submit: Enter
    else if (key.return) {
      onSubmit(currentTree);
    }
  });

  // Toggle selection for a node by path
  const toggleSelection = (path: string) => {
    const updateNode = (node: TreeNode, targetPath: string, currentPath = ''): TreeNode => {
      const nodePath = currentPath ? `${currentPath}/${node.name}` : node.name;

      if (nodePath === targetPath) {
        const newSelected = !node.selected;
        return {
          ...node,
          selected: newSelected,
          children: node.children?.map((child) => updateNodeSelection(child, newSelected)),
        };
      }

      if (node.children) {
        return {
          ...node,
          children: node.children.map((child) => updateNode(child, targetPath, nodePath)),
        };
      }

      return node;
    };

    setCurrentTree((prev) => updateNode(prev, path));
  };

  // Update selection state for all children
  const updateNodeSelection = (node: TreeNode, selected: boolean): TreeNode => {
    return {
      ...node,
      selected,
      children: node.children?.map((child) => updateNodeSelection(child, selected)),
    };
  };

  // Toggle expanded state for a directory
  const toggleExpanded = (path: string) => {
    const updateNode = (node: TreeNode, targetPath: string, currentPath = ''): TreeNode => {
      const nodePath = currentPath ? `${currentPath}/${node.name}` : node.name;

      if (nodePath === targetPath && node.type === 'directory') {
        return {
          ...node,
          expanded: !node.expanded,
        };
      }

      if (node.children) {
        return {
          ...node,
          children: node.children.map((child) => updateNode(child, targetPath, nodePath)),
        };
      }

      return node;
    };

    setCurrentTree((prev) => updateNode(prev, path));
  };

  // Count selected items
  const countSelected = (node: TreeNode): number => {
    let count = node.selected ? 1 : 0;
    if (node.children) {
      count += node.children.reduce((sum, child) => sum + countSelected(child), 0);
    }
    return count;
  };

  const selectedCount = countSelected(currentTree);

  return (
    <Box flexDirection="column" gap={1}>
      {title && (
        <Text color={theme.colors.text} bold>
          {title}
        </Text>
      )}

      {/* Tree display */}
      <Box flexDirection="column">
        {flatItems.map((item, index) => {
          const isCursor = index === cursorIndex;
          const indent = '  '.repeat(item.depth);
          const icon = item.node.type === 'directory'
            ? (item.node.expanded ? '📂' : '📁')
            : '📄';
          const checkbox = item.node.selected ? '[✓]' : '[ ]';
          const cursor = isCursor ? '▸' : ' ';

          return (
            <Box key={item.path}>
              <Text color={isCursor ? theme.colors.primary : theme.colors.text}>
                {cursor} {indent}{checkbox} {icon} {item.node.name}
              </Text>
            </Box>
          );
        })}
      </Box>

      {/* Status and hints */}
      <Box flexDirection="column" gap={1} marginTop={1}>
        <Text color={theme.colors.info}>
          {selectedCount} selected
        </Text>
        <Text color={theme.colors.textMuted}>
          ↑↓/jk: navigate | ←→/hl: expand/collapse | space: select | enter: confirm
        </Text>
      </Box>
    </Box>
  );
};
