#!/usr/bin/env python3
"""
The Agentic Startup - Spec Generation Script
Creates numbered spec directories with auto-incrementing IDs

Location: plugins/start/skills/specify-meta/spec.py
Template resolution: skills/[template-name]/template.md (primary)
                    templates/[template-name].md (fallback, deprecated)
"""

import argparse
import re
import sys
from pathlib import Path
from typing import Optional


# Get plugin root from script location
# This script is at: plugins/start/skills/specify-meta/spec.py
# Plugin root is: plugins/start/
script_dir = Path(__file__).resolve().parent
plugin_root = script_dir.parent.parent

# Specs are created in the current working directory
SPECS_DIR = Path(".start/specs")
LEGACY_SPECS_DIR = Path("docs/specs")
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


def resolve_specs_dir() -> Path:
    """
    Resolve the active specs directory.

    Resolution order:
    1. .start/specs/ (primary)
    2. docs/specs/ (legacy fallback)
    """
    if SPECS_DIR.exists():
        return SPECS_DIR
    if LEGACY_SPECS_DIR.exists():
        return LEGACY_SPECS_DIR
    return SPECS_DIR  # default to new location for creation


def get_next_spec_id() -> str:
    """Get next available spec ID by scanning existing directories."""
    max_id = 0

    # Check both directories for existing specs
    for specs_dir in [SPECS_DIR, LEGACY_SPECS_DIR]:
        if specs_dir.exists():
            for dir_path in specs_dir.iterdir():
                if dir_path.is_dir():
                    match = re.match(r'^(\d{3})-', dir_path.name)
                    if match:
                        num = int(match.group(1))
                        if num > max_id:
                            max_id = num

    # Return next ID with zero-padding
    return f"{max_id + 1:03d}"


def sanitize_name(name: str) -> str:
    """Convert feature name to URL-friendly directory name."""
    name = name.lower()
    name = re.sub(r'[^a-z0-9]+', '-', name)
    name = name.strip('-')
    return name


def find_spec_dir(spec_id: str) -> Optional[Path]:
    """
    Find a spec directory by ID, checking both locations.

    Resolution order:
    1. .start/specs/[NNN]-*/ (primary)
    2. docs/specs/[NNN]-*/ (legacy fallback)
    """
    for specs_dir in [SPECS_DIR, LEGACY_SPECS_DIR]:
        if specs_dir.exists():
            for dir_path in specs_dir.iterdir():
                if dir_path.is_dir() and dir_path.name.startswith(f"{spec_id}-"):
                    return dir_path
    return None


def resolve_doc_path(spec_dir: Path, new_name: str, legacy_name: str) -> Optional[Path]:
    """Resolve document path, checking new name first, then legacy."""
    new_path = spec_dir / new_name
    if new_path.exists():
        return new_path
    legacy_path = spec_dir / legacy_name
    if legacy_path.exists():
        return legacy_path
    return None


def read_spec(spec_id: str) -> None:
    """Read spec metadata and output TOML format."""
    spec_dir = find_spec_dir(spec_id)

    if not spec_dir:
        print(f"Error: Spec {spec_id} not found", file=sys.stderr)
        sys.exit(1)

    # Extract name from directory
    name = spec_dir.name[len(spec_id) + 1:]  # Remove "001-" prefix

    # Generate TOML output
    print(f'id = "{spec_id}"')
    print(f'name = "{name}"')
    print(f'dir = "{spec_dir}"')

    # List spec documents (check new names first, then legacy)
    print()
    print("[spec]")

    prd = resolve_doc_path(spec_dir, "requirements.md", "product-requirements.md")
    if prd:
        print(f'prd = "{prd}"')

    sdd = resolve_doc_path(spec_dir, "solution.md", "solution-design.md")
    if sdd:
        print(f'sdd = "{sdd}"')

    # Factory artifacts: manifest, units, scenarios
    manifest = spec_dir / "manifest.md"
    if manifest.exists():
        print(f'manifest = "{manifest}"')

    units_dir = spec_dir / "units"
    if units_dir.is_dir():
        unit_files = sorted(units_dir.glob("*.md"))
        if unit_files:
            units_str = ", ".join(f'"{f}"' for f in unit_files)
            print(f"units = [{units_str}]")

    scenarios_dir = spec_dir / "scenarios"
    if scenarios_dir.is_dir():
        scenario_dirs = sorted([d for d in scenarios_dir.iterdir() if d.is_dir()])
        if scenario_dirs:
            for scenario_dir in scenario_dirs:
                scenario_files = sorted(scenario_dir.glob("*.md"))
                if scenario_files:
                    files_str = ", ".join(f'"{f}"' for f in scenario_files)
                    print(f'scenarios_{scenario_dir.name} = [{files_str}]')

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

    # List all files (including units/ and scenarios/ directory contents)
    print()
    print("files = [")
    files = sorted([f.name for f in spec_dir.iterdir() if f.is_file()])
    if (spec_dir / "units").is_dir():
        unit_files = sorted([f"units/{f.name}" for f in (spec_dir / "units").iterdir() if f.is_file()])
        files.extend(unit_files)
    if (spec_dir / "scenarios").is_dir():
        for scenario_dir in sorted((spec_dir / "scenarios").iterdir()):
            if scenario_dir.is_dir():
                scenario_files = sorted([f"scenarios/{scenario_dir.name}/{f.name}" for f in scenario_dir.iterdir() if f.is_file()])
                files.extend(scenario_files)
    for i, file in enumerate(files):
        comma = "," if i < len(files) - 1 else ""
        print(f'  "{file}"{comma}')
    print("]")


def create_factory_directories(spec_dir: Path) -> None:
    """Create units/ and scenarios/ directories for Dark Factory specs."""
    (spec_dir / "units").mkdir(parents=True, exist_ok=True)
    (spec_dir / "scenarios").mkdir(parents=True, exist_ok=True)
    print("Created factory directories: units/, scenarios/")


def create_spec(feature_name: str, template: Optional[str] = None) -> None:
    """Create a new spec directory with optional template."""
    # Check if feature_name is an existing spec ID (3 digits)
    is_spec_id = re.match(r'^\d{3}$', feature_name)

    if is_spec_id and template:
        # Try to find existing directory with this ID
        spec_dir = find_spec_dir(feature_name)

        if spec_dir:
            print(f"Adding template to existing spec: {spec_dir}")

            # Special handling for factory directories
            if template == "specify-factory":
                create_factory_directories(spec_dir)
                return

            # Map template names to short filenames
            filename_map = {
                "specify-requirements": "requirements.md",
                "product-requirements": "requirements.md",
                "specify-solution": "solution.md",
                "solution-design": "solution.md",
            }
            dest_name = filename_map.get(template, f"{template}.md")

            try:
                template_file = get_template_path(template)
                dest_file = spec_dir / dest_name
                dest_file.write_text(template_file.read_text())
                print(f"Generated template: {dest_name}")
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

    # Create spec directory with factory subdirectories
    spec_dir.mkdir(parents=True, exist_ok=True)
    create_factory_directories(spec_dir)

    print(f"Created spec directory: {spec_dir}")
    print(f"Spec ID: {spec_id}")

    # Copy template if requested
    if template:
        filename_map = {
            "specify-requirements": "requirements.md",
            "product-requirements": "requirements.md",
            "specify-solution": "solution.md",
            "solution-design": "solution.md",
        }
        dest_name = filename_map.get(template, f"{template}.md")

        try:
            template_file = get_template_path(template)
            dest_file = spec_dir / dest_name
            dest_file.write_text(template_file.read_text())
            print(f"Generated template: {dest_name}")
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
