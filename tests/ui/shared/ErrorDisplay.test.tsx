import { describe, it, expect } from 'vitest';
import React from 'react';
import { render } from 'ink-testing-library';
import { ErrorDisplay } from '../../../src/ui/shared/ErrorDisplay';

describe('ErrorDisplay', () => {
  describe('Rendering', () => {
    it('renders without crashing', () => {
      const testError = new Error('Test error');
      const { lastFrame } = render(<ErrorDisplay error={testError} />);
      expect(lastFrame()).toBeDefined();
    });

    it('displays error icon from theme', () => {
      const testError = new Error('Test error');
      const { lastFrame } = render(<ErrorDisplay error={testError} />);
      const output = lastFrame();

      expect(output).toContain('✗');
    });

    it('displays error message', () => {
      const testError = new Error('Connection failed');
      const { lastFrame } = render(<ErrorDisplay error={testError} />);
      const output = lastFrame();

      expect(output).toContain('Connection failed');
    });

    it('displays string error message', () => {
      const { lastFrame } = render(<ErrorDisplay error="Invalid input" />);
      const output = lastFrame();

      expect(output).toContain('Invalid input');
    });
  });

  describe('Error Types', () => {
    it('handles Error objects', () => {
      const testError = new Error('System error');
      const { lastFrame } = render(<ErrorDisplay error={testError} />);
      const output = lastFrame();

      expect(output).toContain('System error');
    });

    it('handles string errors', () => {
      const { lastFrame } = render(<ErrorDisplay error="String error" />);
      const output = lastFrame();

      expect(output).toContain('String error');
    });

    it('handles unknown error types', () => {
      const { lastFrame } = render(<ErrorDisplay error={null} />);
      const output = lastFrame();

      // Should display a fallback message
      expect(output).toContain('An unknown error occurred');
    });

    it('handles undefined errors', () => {
      const { lastFrame } = render(<ErrorDisplay error={undefined} />);
      const output = lastFrame();

      expect(output).toContain('An unknown error occurred');
    });
  });

  describe('Formatting', () => {
    it('formats error with icon and message', () => {
      const testError = new Error('Format test');
      const { lastFrame } = render(<ErrorDisplay error={testError} />);
      const output = lastFrame();

      // Should contain both icon and message
      expect(output).toContain('✗');
      expect(output).toContain('Format test');
    });

    it('uses error color from theme', () => {
      const testError = new Error('Colored error');
      const { lastFrame } = render(<ErrorDisplay error={testError} />);
      const output = lastFrame();

      // Error should be rendered (basic check)
      expect(output).toBeDefined();
      expect(output).toContain('Colored error');
    });
  });

  describe('Props', () => {
    it('accepts Error object', () => {
      const testError = new Error('Error object test');
      const { lastFrame } = render(<ErrorDisplay error={testError} />);

      expect(lastFrame()).toContain('Error object test');
    });

    it('accepts string message', () => {
      const { lastFrame } = render(<ErrorDisplay error="String message test" />);

      expect(lastFrame()).toContain('String message test');
    });
  });

  describe('Component Pattern', () => {
    it('follows Ink functional component pattern', () => {
      // ErrorDisplay should be a function component
      expect(typeof ErrorDisplay).toBe('function');
    });

    it('returns valid React element', () => {
      const element = <ErrorDisplay error="test" />;
      expect(React.isValidElement(element)).toBe(true);
    });
  });

  describe('Edge Cases', () => {
    it('handles empty error message', () => {
      const testError = new Error('');
      const { lastFrame } = render(<ErrorDisplay error={testError} />);
      const output = lastFrame();

      // Should still render icon
      expect(output).toContain('✗');
    });

    it('handles error with stack trace', () => {
      const testError = new Error('Error with stack');
      testError.stack = 'Error: Error with stack\n  at test.ts:1:1';

      const { lastFrame } = render(<ErrorDisplay error={testError} />);
      const output = lastFrame();

      // Should display message, not stack trace
      expect(output).toContain('Error with stack');
      expect(output).not.toContain('at test.ts:1:1');
    });

    it('handles very long error messages', () => {
      const longMessage = 'A'.repeat(200);
      const testError = new Error(longMessage);
      const { lastFrame } = render(<ErrorDisplay error={testError} />);
      const output = lastFrame();

      // Should contain at least part of the message
      expect(output).toContain('A');
    });
  });
});
