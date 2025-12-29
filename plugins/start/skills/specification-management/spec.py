#!/usr/bin/env python3
"""
The Agentic Startup - Spec Generation Script
Creates numbered spec directories with auto-incrementing IDs

Location: plugins/start/skills/specification-management/spec.py
Template resolution: skills/[template-name]/template.md (primary)
                    templates/[template-name].md (fallback, deprecated)
"""

import argparse
import re
import sys
from pathlib import Path
from typing import Optional


# Get plugin root from script location
# This script is at: plugins/start/skills/specification-management/spec.py
# Plugin root is: plugins/start/
script_dir = Path(__file__).resolve().parent
plugin_root = script_dir.parent.parent

# Specs are created in the current working directory
SPECS_DIR = Path("docs/specs")
# Skills directory for primary template lookup
SKILLS_DIR = plugin_root / "skills"
# Templates directory for fallback (deprecated)
TEMPLATES_DIR = plugin_root / "templates"


def get_template_path(template_name: str) -> Path:
    """
    Resolve template path with skill-first, legacy-fallback pattern.

    Resolution order:
    1. skills/[template-name]/template.md (new location)
    2. templates/[template-name].md (deprecated, backward compat)
    """
    # Primary: Look in skill directory
    skill_template = SKILLS_DIR / template_name / "template.md"
    if skill_template.exists():
        return skill_template

    # Fallback: Legacy templates directory (deprecated)
    legacy_template = TEMPLATES_DIR / f"{template_name}.md"
    if legacy_template.exists():
        print(f"Warning: Using deprecated template location. "
              f"Template should be at: {skill_template}", file=sys.stderr)
        return legacy_template

    raise FileNotFoundError(f"Template not found: {template_name}")


def get_next_spec_id() -> str:
    """Get next available spec ID by scanning existing directories."""
    max_id = 0

    if SPECS_DIR.exists():
        # Find highest existing spec number
        for dir_path in SPECS_DIR.iterdir():
            if dir_path.is_dir():
                # Extract number from pattern: 001-feature-name
                match = re.match(r'^(\d{3})-', dir_path.name)
                if match:
                    num = int(match.group(1))
                    if num > max_id:
                        max_id = num

    # Return next ID with zero-padding
    return f"{max_id + 1:03d}"


def sanitize_name(name: str) -> str:
    """Convert feature name to URL-friendly directory name."""
    # Convert to lowercase
    name = name.lower()
    # Replace special characters with hyphens
    name = re.sub(r'[^a-z0-9]+', '-', name)
    # Remove leading/trailing hyphens
    name = name.strip('-')
    return name


def read_spec(spec_id: str) -> None:
    """Read spec metadata and output TOML format."""
    if not SPECS_DIR.exists():
        print("Error: Specs directory not found", file=sys.stderr)
        sys.exit(1)

    # Find directory matching spec ID
    spec_dir = None
    for dir_path in SPECS_DIR.iterdir():
        if dir_path.is_dir() and dir_path.name.startswith(f"{spec_id}-"):
            spec_dir = dir_path
            break

    if not spec_dir:
        print(f"Error: Spec {spec_id} not found", file=sys.stderr)
        sys.exit(1)

    # Extract name from directory
    name = spec_dir.name[len(spec_id) + 1:]  # Remove "001-" prefix

    # Generate TOML output
    print(f'id = "{spec_id}"')
    print(f'name = "{name}"')
    print(f'dir = "{spec_dir}"')

    # List spec documents
    print()
    print("[spec]")
    if (spec_dir / "product-requirements.md").exists():
        print(f'prd = "{spec_dir / "product-requirements.md"}"')
    if (spec_dir / "solution-design.md").exists():
        print(f'sdd = "{spec_dir / "solution-design.md"}"')
    if (spec_dir / "implementation-plan.md").exists():
        print(f'plan = "{spec_dir / "implementation-plan.md"}"')

    # List quality gates if they exist
    gate_files = [
        ("definition-of-ready.md", "definition_of_ready"),
        ("definition-of-done.md", "definition_of_done"),
        ("task-definition-of-done.md", "task_definition_of_done"),
    ]

    gate_exists = any((spec_dir / file).exists() for file, _ in gate_files)
    if gate_exists:
        print()
        print("[gates]")
        for file, key in gate_files:
            if (spec_dir / file).exists():
                print(f'{key} = "{spec_dir / file}"')

    # List all files
    print()
    print("files = [")
    files = sorted([f.name for f in spec_dir.iterdir() if f.is_file()])
    for i, file in enumerate(files):
        comma = "," if i < len(files) - 1 else ""
        print(f'  "{file}"{comma}')
    print("]")


def create_spec(feature_name: str, template: Optional[str] = None) -> None:
    """Create a new spec directory with optional template."""
    # Check if feature_name is an existing spec ID (3 digits)
    is_spec_id = re.match(r'^\d{3}$', feature_name)

    if is_spec_id and template:
        # Try to find existing directory with this ID
        if SPECS_DIR.exists():
            for dir_path in SPECS_DIR.iterdir():
                if dir_path.is_dir() and dir_path.name.startswith(f"{feature_name}-"):
                    spec_dir = dir_path
                    spec_id = feature_name
                    print(f"Adding template to existing spec: {spec_dir}")

                    # Copy template to existing directory using new resolution
                    try:
                        template_file = get_template_path(template)
                        dest_file = spec_dir / f"{template}.md"
                        dest_file.write_text(template_file.read_text())
                        print(f"Generated template: {template}.md")
                    except FileNotFoundError as e:
                        print(f"Error: {e}", file=sys.stderr)
                        sys.exit(1)
                    return

        # If we get here, the spec ID was not found
        print(f"Error: Spec {feature_name} not found", file=sys.stderr)
        sys.exit(1)

    # Create new spec directory
    spec_id = get_next_spec_id()
    sanitized_name = sanitize_name(feature_name)
    spec_dir = SPECS_DIR / f"{spec_id}-{sanitized_name}"

    # Create spec directory
    spec_dir.mkdir(parents=True, exist_ok=True)

    print(f"Created spec directory: {spec_dir}")
    print(f"Spec ID: {spec_id}")

    # Copy template if requested
    if template:
        try:
            template_file = get_template_path(template)
            dest_file = spec_dir / f"{template}.md"
            dest_file.write_text(template_file.read_text())
            print(f"Generated template: {template}.md")
        except FileNotFoundError as e:
            print(f"Warning: {e}", file=sys.stderr)

    print("Specification directory created successfully")


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="The Agentic Startup - Spec Generation Script"
    )
    parser.add_argument(
        "feature_name",
        help="Feature name or spec ID (for --read mode)"
    )
    parser.add_argument(
        "--add",
        metavar="TEMPLATE",
        help="Add template file to spec directory"
    )
    parser.add_argument(
        "--read",
        action="store_true",
        help="Read spec metadata and output TOML"
    )

    args = parser.parse_args()

    # Handle read mode
    if args.read:
        read_spec(args.feature_name)
    else:
        # Handle spec creation
        create_spec(args.feature_name, args.add)


if __name__ == "__main__":
    main()
