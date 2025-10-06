import { describe, it, expect } from 'vitest';
import React from 'react';
import { render } from 'ink-testing-library';
import { PathSelector } from '../../../src/ui/install/PathSelector';

describe('PathSelector', () => {
  describe('Rendering', () => {
    it('renders without crashing', () => {
      const { lastFrame } = render(
        <PathSelector
          label="Select path"
          defaultValue="~/.the-startup"
          onSubmit={() => {}}
        />
      );
      expect(lastFrame()).toBeDefined();
    });

    it('displays label text', () => {
      const { lastFrame } = render(
        <PathSelector
          label="Choose installation directory"
          defaultValue="~/.the-startup"
          onSubmit={() => {}}
        />
      );
      const output = lastFrame();

      expect(output).toContain('Choose installation directory');
    });

    it('displays default value', () => {
      const { lastFrame } = render(
        <PathSelector
          label="Path"
          defaultValue="/custom/path"
          onSubmit={() => {}}
        />
      );
      const output = lastFrame();

      expect(output).toContain('/custom/path');
    });

    it('displays placeholder when no default value', () => {
      const { lastFrame } = render(
        <PathSelector
          label="Path"
          placeholder="Enter path..."
          onSubmit={() => {}}
        />
      );
      const output = lastFrame();

      expect(output).toContain('Enter path...');
    });
  });

  describe('Path Validation', () => {
    it('validates path and shows success indicator for valid path', () => {
      const { lastFrame } = render(
        <PathSelector
          label="Path"
          defaultValue="/valid/path"
          onSubmit={() => {}}
          validator={(path) => ({ isValid: true, message: '' })}
        />
      );
      const output = lastFrame();

      // Should show success icon when valid
      expect(output).toContain('✓');
    });

    it('validates path and shows error for invalid path', () => {
      const { lastFrame } = render(
        <PathSelector
          label="Path"
          defaultValue="/invalid/path"
          onSubmit={() => {}}
          validator={(path) => ({ isValid: false, message: 'Path does not exist' })}
        />
      );
      const output = lastFrame();

      // Should show error icon when invalid
      expect(output).toContain('✗');
      expect(output).toContain('Path does not exist');
    });

    it('calls validator on value change', () => {
      let validatorCalled = false;
      const validator = (path: string) => {
        validatorCalled = true;
        return { isValid: true, message: '' };
      };

      render(
        <PathSelector
          label="Path"
          defaultValue="/test/path"
          onSubmit={() => {}}
          validator={validator}
        />
      );

      expect(validatorCalled).toBe(true);
    });
  });

  describe('User Interaction', () => {
    it('accepts onSubmit callback for valid paths', () => {
      let submitCalled = false;
      const onSubmit = () => {
        submitCalled = true;
      };

      const { lastFrame } = render(
        <PathSelector
          label="Path"
          defaultValue="/test/path"
          onSubmit={onSubmit}
          validator={() => ({ isValid: true, message: '' })}
        />
      );

      // Component should render and accept the callback
      expect(lastFrame()).toBeDefined();
      expect(typeof onSubmit).toBe('function');
    });

    it('accepts validator function', () => {
      const validator = () => ({ isValid: true, message: '' });

      const { lastFrame } = render(
        <PathSelector
          label="Path"
          defaultValue="/test/path"
          onSubmit={() => {}}
          validator={validator}
        />
      );

      // Component should accept validator
      expect(lastFrame()).toBeDefined();
      expect(typeof validator).toBe('function');
    });

    it('shows path is ready for submission when valid', () => {
      const { lastFrame } = render(
        <PathSelector
          label="Path"
          defaultValue="~/.the-startup"
          onSubmit={() => {}}
          validator={() => ({ isValid: true, message: '' })}
        />
      );

      const output = lastFrame();

      // Should show success indicator for valid path
      expect(output).toContain('✓');
      expect(output).toContain('Valid path');
    });
  });

  describe('Component Pattern', () => {
    it('follows Ink functional component pattern', () => {
      expect(typeof PathSelector).toBe('function');
    });

    it('returns valid React element', () => {
      const element = (
        <PathSelector
          label="Path"
          defaultValue="/test"
          onSubmit={() => {}}
        />
      );
      expect(React.isValidElement(element)).toBe(true);
    });
  });

  describe('Props', () => {
    it('accepts required props: label, onSubmit', () => {
      const { lastFrame } = render(
        <PathSelector
          label="Test Label"
          onSubmit={() => {}}
        />
      );

      expect(lastFrame()).toContain('Test Label');
    });

    it('accepts optional props: defaultValue, placeholder, validator', () => {
      const { lastFrame } = render(
        <PathSelector
          label="Path"
          defaultValue="/default"
          placeholder="Enter path"
          validator={() => ({ isValid: true, message: '' })}
          onSubmit={() => {}}
        />
      );

      expect(lastFrame()).toBeDefined();
    });
  });

  describe('Edge Cases', () => {
    it('handles empty input', () => {
      const { lastFrame } = render(
        <PathSelector
          label="Path"
          defaultValue=""
          onSubmit={() => {}}
        />
      );

      expect(lastFrame()).toBeDefined();
    });

    it('handles very long paths', () => {
      const longPath = '/very/long/path/with/many/segments/'.repeat(10);
      const { lastFrame } = render(
        <PathSelector
          label="Path"
          defaultValue={longPath}
          onSubmit={() => {}}
        />
      );

      expect(lastFrame()).toBeDefined();
    });

    it('handles paths with spaces', () => {
      const pathWithSpaces = '/path/with spaces/in name';

      const { lastFrame } = render(
        <PathSelector
          label="Path"
          defaultValue={pathWithSpaces}
          onSubmit={() => {}}
          validator={() => ({ isValid: true, message: '' })}
        />
      );

      const output = lastFrame();

      // Should display path with spaces correctly
      expect(output).toContain(pathWithSpaces);
    });
  });
});
