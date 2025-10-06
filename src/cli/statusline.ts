import { spawn } from 'child_process';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

/**
 * Statusline CLI Command
 *
 * Implements simple stdio passthrough to shell scripts:
 * - Unix: bin/statusline.sh
 * - Windows: bin/statusline.ps1
 *
 * No flags, just stdio passthrough for cross-platform compatibility.
 *
 * Business Rules (PRD lines 213-223):
 * - Rule 1: Detect platform (win32 vs unix)
 * - Rule 2: Pass stdin to appropriate shell script
 * - Rule 3: Return stdout to user
 * - Rule 4: Exit with shell script exit code
 *
 * @example
 * ```bash
 * echo "custom data" | the-agentic-startup statusline
 * ```
 */

/**
 * Get the project root directory (where package.json and bin/ are located)
 */
function getProjectRoot(): string {
  // In ESM, __filename and __dirname are not available
  // We need to derive them from import.meta.url
  const currentFile = fileURLToPath(import.meta.url);
  const currentDir = dirname(currentFile);

  // Navigate from src/cli/ to project root
  return join(currentDir, '..', '..');
}

/**
 * Execute statusline command with stdio passthrough
 *
 * @returns Promise that resolves when script completes
 */
export async function statuslineCommand(): Promise<void> {
  return new Promise((resolve, reject) => {
    // Detect platform
    const isWindows = process.platform === 'win32';
    const scriptName = isWindows ? 'statusline.ps1' : 'statusline.sh';
    const shell = isWindows ? 'powershell.exe' : undefined;

    // Resolve script path from project root
    const projectRoot = getProjectRoot();
    const scriptPath = join(projectRoot, 'bin', scriptName);

    // Spawn shell script with stdio passthrough
    const args = isWindows ? ['-File', scriptPath] : [scriptPath];
    const child = spawn(
      isWindows ? shell! : scriptPath,
      isWindows ? args : [],
      {
        stdio: ['inherit', 'inherit', 'inherit'], // Pass stdin, stdout, stderr
        shell: !isWindows, // Unix needs shell for script execution
      }
    );

    // Handle exit
    child.on('close', (code) => {
      if (code === 0) {
        resolve();
      } else {
        reject(new Error(`Statusline script exited with code ${code}`));
      }
    });

    // Handle errors
    child.on('error', (error) => {
      reject(new Error(`Failed to execute statusline script: ${error.message}`));
    });
  });
}
