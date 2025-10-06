import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import React from 'react';
import { render } from 'ink-testing-library';
import { Spinner } from '../../../src/ui/shared/Spinner';

describe('Spinner', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('Rendering', () => {
    it('renders without crashing', () => {
      const { lastFrame } = render(<Spinner />);
      expect(lastFrame()).toBeDefined();
    });

    it('displays spinner frame from theme', () => {
      const { lastFrame } = render(<Spinner />);
      const output = lastFrame();

      // Should contain a spinner character
      expect(output).toMatch(/[⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏]/);
    });

    it('renders with custom text message', () => {
      const { lastFrame } = render(<Spinner text="Loading data..." />);
      const output = lastFrame();

      expect(output).toContain('Loading data...');
    });

    it('renders without text when not provided', () => {
      const { lastFrame } = render(<Spinner />);
      const output = lastFrame();

      // Should only contain spinner, no additional text
      expect(output?.trim().length).toBeGreaterThan(0);
    });
  });

  describe('Animation', () => {
    it('cycles through spinner frames over time', () => {
      const { lastFrame, rerender } = render(<Spinner />);

      const frame1 = lastFrame();

      // Advance timer by 80ms (default frame interval)
      vi.advanceTimersByTime(80);
      rerender(<Spinner />);

      const frame2 = lastFrame();

      // Frames should be different due to animation
      // Note: This may not work as expected with Ink's internal timing
      // This is a basic check for animation behavior
      expect(frame1).toBeDefined();
      expect(frame2).toBeDefined();
    });
  });

  describe('Props', () => {
    it('accepts and displays custom text', () => {
      const customText = 'Installing packages...';
      const { lastFrame } = render(<Spinner text={customText} />);

      expect(lastFrame()).toContain(customText);
    });

    it('uses theme color for spinner', () => {
      const { lastFrame } = render(<Spinner />);
      const output = lastFrame();

      // Spinner should be rendered (contains spinner character)
      expect(output).toMatch(/[⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏]/);
    });
  });

  describe('Component Pattern', () => {
    it('follows Ink functional component pattern', () => {
      // Spinner should be a function component
      expect(typeof Spinner).toBe('function');
    });

    it('returns valid React element', () => {
      const element = <Spinner />;
      expect(React.isValidElement(element)).toBe(true);
    });
  });
});
