package log

import "fmt"

// ValidateLogCommand validates the command-line flags for the log command
// Enforces exactly one mode (assistant, user, or read) and validates read mode requirements
func ValidateLogCommand(flags LogFlags) error {
	// Count how many modes are selected
	modeCount := 0
	if flags.Assistant {
		modeCount++
	}
	if flags.User {
		modeCount++
	}
	if flags.Read {
		modeCount++
	}

	// Exactly one mode is required
	if modeCount != 1 {
		return fmt.Errorf("exactly one mode required: --assistant, --user, or --read")
	}

	// Read mode specific validation
	if flags.Read {
		// Agent ID is required for read mode
		if flags.AgentID == "" {
			return fmt.Errorf("--agent-id required when using --read")
		}

		// Validate agent ID format and reserved words
		if !isValidAgentID(flags.AgentID) || isReservedWord(flags.AgentID) {
			return fmt.Errorf("invalid agent-id format")
		}

		// Validate lines range
		if flags.Lines < 1 || flags.Lines > 1000 {
			return fmt.Errorf("--lines must be between 1 and 1000")
		}

		// Validate format
		if flags.Format != "" && flags.Format != "json" && flags.Format != "text" {
			return fmt.Errorf("invalid format: must be 'json' or 'text'")
		}
	}

	return nil
}
