import { describe, it, expect } from 'vitest';
import {
  theme,
  Spinner,
  ErrorDisplay,
} from '../../../src/ui/shared/index';

describe('Shared UI Index Exports', () => {
  it('exports theme', () => {
    expect(theme).toBeDefined();
    expect(theme.colors).toBeDefined();
    expect(theme.icons).toBeDefined();
    expect(theme.spacing).toBeDefined();
  });

  it('exports Spinner component', () => {
    expect(Spinner).toBeDefined();
    expect(typeof Spinner).toBe('function');
  });

  it('exports ErrorDisplay component', () => {
    expect(ErrorDisplay).toBeDefined();
    expect(typeof ErrorDisplay).toBe('function');
  });
});
