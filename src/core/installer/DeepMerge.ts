/**
 * DeepMerge - Generic deep merge utility for settings
 *
 * Provides a fully generic merge algorithm that:
 * - Recursively merges nested objects
 * - Deduplicates arrays
 * - Handles primitives with sensible defaults
 * - Works with ANY structure without hardcoding keys
 *
 * This allows adding new settings keys without code changes.
 */

/**
 * Deep merge two objects with array deduplication.
 *
 * Merge strategy:
 * - Objects: Recursively merge properties
 * - Arrays: Concatenate and deduplicate (by value equality)
 * - Primitives: Target value takes precedence (new overwrites existing)
 *
 * Special behaviors:
 * - Arrays are deduplicated to prevent duplicate entries
 * - Nested objects are merged deeply (not replaced)
 * - null/undefined are treated as "no value" and don't overwrite
 *
 * @param target - Existing object (user's current settings)
 * @param source - New object to merge in (template settings)
 * @returns Merged object (target modified in place and returned)
 *
 * @example
 * const existing = {
 *   permissions: { additionalDirectories: ['/old'] },
 *   hooks: { 'existing-hook': { command: 'echo old' } }
 * };
 * const newSettings = {
 *   permissions: { additionalDirectories: ['/new', '/old'] },
 *   statusLine: { command: 'status.sh' },
 *   hooks: { 'new-hook': { command: 'echo new' } }
 * };
 * const merged = deepMerge(existing, newSettings);
 * // Result:
 * // {
 * //   permissions: { additionalDirectories: ['/old', '/new'] }, // deduplicated
 * //   statusLine: { command: 'status.sh' },                      // added
 * //   hooks: {
 * //     'existing-hook': { command: 'echo old' },                // preserved
 * //     'new-hook': { command: 'echo new' }                      // added
 * //   }
 * // }
 */
export function deepMerge<T extends Record<string, any>>(
  target: T,
  source: Partial<T>
): T {
  // Handle null/undefined
  if (!source) {
    return target;
  }

  // Iterate through all keys in source
  for (const key in source) {
    if (!Object.prototype.hasOwnProperty.call(source, key)) {
      continue;
    }

    const sourceValue = source[key];
    const targetValue = target[key];

    // Case 1: Source value is null/undefined - skip
    if (sourceValue === null || sourceValue === undefined) {
      continue;
    }

    // Case 2: Target doesn't have this key - just assign
    if (!(key in target)) {
      target[key] = sourceValue;
      continue;
    }

    // Case 3: Both are arrays - concatenate and deduplicate
    if (Array.isArray(targetValue) && Array.isArray(sourceValue)) {
      target[key] = deduplicateArray([...targetValue, ...sourceValue]) as any;
      continue;
    }

    // Case 4: Both are objects (but not arrays) - recursively merge
    if (isPlainObject(targetValue) && isPlainObject(sourceValue)) {
      target[key] = deepMerge(targetValue, sourceValue);
      continue;
    }

    // Case 5: Primitive or type mismatch - source overwrites target
    // This handles: string, number, boolean, or when types don't match
    target[key] = sourceValue;
  }

  return target;
}

/**
 * Deduplicates an array by value equality.
 *
 * Uses JSON serialization for deep equality of objects.
 * For primitives, uses Set for O(n) deduplication.
 *
 * @param arr - Array to deduplicate
 * @returns Array with duplicates removed (order preserved)
 *
 * @example
 * deduplicateArray([1, 2, 2, 3]) // => [1, 2, 3]
 * deduplicateArray(['/path1', '/path2', '/path1']) // => ['/path1', '/path2']
 * deduplicateArray([{a: 1}, {a: 1}, {a: 2}]) // => [{a: 1}, {a: 2}]
 */
export function deduplicateArray<T>(arr: T[]): T[] {
  // Fast path for primitives
  if (arr.length === 0 || isPrimitive(arr[0])) {
    return [...new Set(arr)];
  }

  // Slow path for objects/arrays - use JSON serialization
  const seen = new Set<string>();
  const result: T[] = [];

  for (const item of arr) {
    const key = JSON.stringify(item);
    if (!seen.has(key)) {
      seen.add(key);
      result.push(item);
    }
  }

  return result;
}

/**
 * Checks if a value is a plain object (not array, not null, not Date, etc.)
 *
 * @param value - Value to check
 * @returns True if value is a plain object
 */
function isPlainObject(value: unknown): value is Record<string, any> {
  return (
    typeof value === 'object' &&
    value !== null &&
    !Array.isArray(value) &&
    Object.getPrototypeOf(value) === Object.prototype
  );
}

/**
 * Checks if a value is a primitive (string, number, boolean, null, undefined)
 *
 * @param value - Value to check
 * @returns True if value is a primitive
 */
function isPrimitive(value: unknown): boolean {
  return (
    value === null ||
    value === undefined ||
    typeof value === 'string' ||
    typeof value === 'number' ||
    typeof value === 'boolean'
  );
}
