#!/usr/bin/env bash
# The Agentic Startup - Spec Generation Script
# Creates numbered spec directories with auto-incrementing IDs

set -euo pipefail

SPECS_DIR="docs/specs"
TEMPLATES_DIR="templates"

# Function to get next spec ID
get_next_spec_id() {
  local max_id=0

  if [ -d "$SPECS_DIR" ]; then
    # Find highest existing spec number
    for dir in "$SPECS_DIR"/*/ ; do
      if [ -d "$dir" ]; then
        basename=$(basename "$dir")
        # Extract number from pattern: 001-feature-name
        if [[ $basename =~ ^([0-9]{3})- ]]; then
          num="${BASH_REMATCH[1]}"
          # Remove leading zeros for comparison
          num=$((10#$num))
          if [ $num -gt $max_id ]; then
            max_id=$num
          fi
        fi
      fi
    done
  fi

  # Return next ID with zero-padding
  printf "%03d" $((max_id + 1))
}

# Function to sanitize feature name
sanitize_name() {
  local name="$1"
  # Convert to lowercase, replace special chars with hyphens
  echo "$name" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]\+/-/g' | sed 's/^-\+\|-\+$//g'
}

# Function to read spec metadata
read_spec() {
  local spec_id="$1"

  if [ ! -d "$SPECS_DIR" ]; then
    echo "Error: Specs directory not found" >&2
    exit 1
  fi

  # Find directory matching spec ID
  local spec_dir=""
  for dir in "$SPECS_DIR"/*/ ; do
    if [ -d "$dir" ]; then
      basename=$(basename "$dir")
      if [[ $basename =~ ^$spec_id- ]]; then
        spec_dir="$dir"
        break
      fi
    fi
  done

  if [ -z "$spec_dir" ]; then
    echo "Error: Spec $spec_id not found" >&2
    exit 1
  fi

  # Extract name from directory
  local full_name=$(basename "$spec_dir")
  local name="${full_name#$spec_id-}"

  # Generate TOML output
  echo "id = \"$spec_id\""
  echo "name = \"$name\""
  echo "dir = \"$spec_dir\""

  # List spec documents
  echo ""
  echo "[spec]"
  [ -f "$spec_dir/product-requirements.md" ] && echo "prd = \"$spec_dir/product-requirements.md\""
  [ -f "$spec_dir/solution-design.md" ] && echo "sdd = \"$spec_dir/solution-design.md\""
  [ -f "$spec_dir/implementation-plan.md" ] && echo "plan = \"$spec_dir/implementation-plan.md\""

  # List quality gates if they exist
  if [ -f "$spec_dir/definition-of-ready.md" ] || [ -f "$spec_dir/definition-of-done.md" ] || [ -f "$spec_dir/task-definition-of-done.md" ]; then
    echo ""
    echo "[gates]"
    [ -f "$spec_dir/definition-of-ready.md" ] && echo "definition_of_ready = \"$spec_dir/definition-of-ready.md\""
    [ -f "$spec_dir/definition-of-done.md" ] && echo "definition_of_done = \"$spec_dir/definition-of-done.md\""
    [ -f "$spec_dir/task-definition-of-done.md" ] && echo "task_definition_of_done = \"$spec_dir/task-definition-of-done.md\""
  fi

  # List all files
  echo ""
  echo "files = ["
  first=true
  for file in "$spec_dir"/* ; do
    if [ -f "$file" ]; then
      if [ "$first" = true ]; then
        first=false
      else
        echo ","
      fi
      echo -n "  \"$(basename "$file")\""
    fi
  done
  echo ""
  echo "]"
}

# Main script logic
main() {
  # Parse command line arguments
  local feature_name=""
  local template=""
  local read_mode=false

  while [ $# -gt 0 ]; do
    case "$1" in
      --add)
        template="$2"
        shift 2
        ;;
      --read)
        read_mode=true
        shift
        ;;
      *)
        feature_name="$1"
        shift
        ;;
    esac
  done

  # Handle --read mode
  if [ "$read_mode" = true ]; then
    if [ -z "$feature_name" ]; then
      echo "Error: Spec ID required for --read" >&2
      exit 1
    fi
    read_spec "$feature_name"
    exit 0
  fi

  # Handle spec creation
  if [ -z "$feature_name" ]; then
    echo "Error: Feature name required" >&2
    echo "Usage: $0 <feature-name> [--add <template>]" >&2
    exit 1
  fi

  # Get next spec ID
  spec_id=$(get_next_spec_id)
  sanitized_name=$(sanitize_name "$feature_name")
  spec_dir="$SPECS_DIR/$spec_id-$sanitized_name"

  # Create spec directory
  mkdir -p "$spec_dir"

  echo "Created spec directory: $spec_dir"
  echo "Spec ID: $spec_id"

  # Copy template if requested
  if [ -n "$template" ]; then
    template_file="$TEMPLATES_DIR/${template}.md"
    if [ ! -f "$template_file" ]; then
      echo "Warning: Template $template_file not found" >&2
    else
      cp "$template_file" "$spec_dir/${template}.md"
      echo "Generated template: ${template}.md"
    fi
  fi

  echo "Specification directory created successfully"
}

main "$@"
