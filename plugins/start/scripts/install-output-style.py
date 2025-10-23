#!/usr/bin/env python3

import json
import sys
from pathlib import Path
import shutil

def main():
    """Install output style for The Agentic Startup framework."""

    # Get plugin root from script location
    # This script is at: plugins/start/scripts/install-output-style.py
    # Plugin root is: plugins/start/
    script_dir = Path(__file__).resolve().parent
    plugin_root = script_dir.parent

    # Define all paths at the top
    source_output_style = plugin_root / 'output-styles' / 'the-startup.md'
    target_dir = Path.home() / '.claude'
    target_output_styles_dir = target_dir / 'output-styles'
    target_output_style = target_output_styles_dir / 'the-startup.md'
    target_settings = target_dir / 'settings.json'

    # Handle --check flag
    if len(sys.argv) > 1 and sys.argv[1] == '--check':
        if target_output_style.exists():
            print("INSTALLED")
            print(f"Location: {target_output_style}")
        else:
            print("NOT_INSTALLED")
        sys.exit(0)

    # Create directories
    target_dir.mkdir(parents=True, exist_ok=True)
    target_output_styles_dir.mkdir(exist_ok=True)

    # Copy output style file
    shutil.copy2(source_output_style, target_output_style)

    # Update settings.json
    if target_settings.exists():
        with open(target_settings, 'r') as f:
            settings = json.load(f)
    else:
        settings = {}

    settings['outputStyle'] = 'The Startup'

    with open(target_settings, 'w') as f:
        json.dump(settings, f, indent=2, ensure_ascii=False)

    # Print success
    print("✓ Output style installed successfully!\n")
    print(f"Location: {target_output_style}")
    print(f"Configuration: {target_settings}")
    print("Output Style: the-startup\n")
    print("Changes:")
    print("• Installed output-styles/the-startup.md")
    print("• Updated settings.json with outputStyle field\n")
    print("⚠️  Use /output-style command to activate immediately.")

if __name__ == '__main__':
    main()
