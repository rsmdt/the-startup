import { describe, it, expect } from 'vitest';
import { theme } from '../../../src/ui/shared/theme';

describe('Theme', () => {
  describe('Color Scheme', () => {
    it('defines success color mapped from Go CharmTheme', () => {
      expect(theme.colors.success).toBe('#04B575');
    });

    it('defines error color mapped from Go CharmTheme', () => {
      expect(theme.colors.error).toBe('#FF4444');
    });

    it('defines warning color mapped from Go CharmTheme', () => {
      expect(theme.colors.warning).toBe('#FFA500');
    });

    it('defines info color mapped from Go CharmTheme', () => {
      expect(theme.colors.info).toBe('#3C7EFF');
    });

    it('defines text color mapped from Go CharmTheme', () => {
      expect(theme.colors.text).toBe('#FAFAFA');
    });

    it('defines text muted color mapped from Go CharmTheme', () => {
      expect(theme.colors.textMuted).toBe('#606060');
    });

    it('defines text bright color mapped from Go CharmTheme', () => {
      expect(theme.colors.textBright).toBe('#42FF76');
    });

    it('defines primary color mapped from Go CharmTheme', () => {
      expect(theme.colors.primary).toBe('#FF06B7');
    });
  });

  describe('Icons', () => {
    it('defines success icon', () => {
      expect(theme.icons.success).toBe('✓');
    });

    it('defines error icon', () => {
      expect(theme.icons.error).toBe('✗');
    });

    it('defines warning icon', () => {
      expect(theme.icons.warning).toBe('⚠');
    });

    it('defines info icon', () => {
      expect(theme.icons.info).toBe('ℹ');
    });

    it('defines spinner frames for loading animation', () => {
      expect(Array.isArray(theme.icons.spinner)).toBe(true);
      expect(theme.icons.spinner.length).toBeGreaterThan(0);
    });
  });

  describe('Styles', () => {
    it('defines consistent spacing units', () => {
      expect(theme.spacing).toBeDefined();
      expect(theme.spacing.small).toBeDefined();
      expect(theme.spacing.medium).toBeDefined();
      expect(theme.spacing.large).toBeDefined();
    });
  });
});
