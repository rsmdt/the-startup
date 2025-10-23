#!/usr/bin/env python3

import json
import os
import sys
from pathlib import Path
import shutil
import stat

def main():
    """Install statusline for The Agentic Startup framework."""

    # Get plugin root from script location
    # This script is at: plugins/start/scripts/install-statusline.py
    # Plugin root is: plugins/start/
    script_dir = Path(__file__).resolve().parent
    plugin_root = script_dir.parent

    # Define all paths at the top
    source_statusline = plugin_root / 'hooks' / 'statusline.sh'
    target_dir = Path.home() / '.claude'
    target_statusline = target_dir / 'statusline.sh'
    target_settings = target_dir / 'settings.json'

    # Handle --check flag
    if len(sys.argv) > 1 and sys.argv[1] == '--check':
        # Check if statusline.sh exists
        file_exists = target_statusline.exists()

        # Check if settings.json is configured
        statusline_configured = False

        if target_settings.exists():
            try:
                with open(target_settings, 'r') as f:
                    settings = json.load(f)

                    # Check if statusLine is configured
                    if 'statusLine' in settings:
                        statusline_configured = True
            except (json.JSONDecodeError, KeyError):
                pass

        # Both must be true for full installation
        if file_exists and statusline_configured:
            print("INSTALLED")
            print(f"StatusLine: {target_statusline}")
            print(f"Settings: {target_settings} (statusLine configured)")
        else:
            print("NOT_INSTALLED")
        sys.exit(0)

    # Create directory
    target_dir.mkdir(parents=True, exist_ok=True)

    # Copy statusline.sh to .claude/ directory
    shutil.copy2(source_statusline, target_statusline)

    # Set execute permissions (Unix/macOS only)
    if os.name != 'nt':
        target_statusline.chmod(target_statusline.stat().st_mode | stat.S_IEXEC)

    # Update settings.json with statusLine configuration
    if target_settings.exists():
        with open(target_settings, 'r') as f:
            settings = json.load(f)
    else:
        settings = {}

    # Configure statusLine
    settings['statusLine'] = {
        "type": "command",
        "command": str(target_statusline)
    }

    with open(target_settings, 'w') as f:
        json.dump(settings, f, indent=2, ensure_ascii=False)

    # Print success
    print("✓ Statusline installed successfully!\n")
    print(f"StatusLine: {target_statusline}")
    print(f"Configuration: {target_settings}\n")
    print("Configured:")
    print("• statusLine: Shows git branch in Claude Code statusline\n")
    print("Changes:")
    print("• Installed statusline.sh to ~/.claude/")
    print("• Set execute permissions on statusline.sh")
    print("• Updated settings.json with statusLine field\n")
    print("⚠️  Changes take effect on next Claude Code session.")
    print("    Exit and restart Claude Code to see the statusline.")

if __name__ == '__main__':
    main()
