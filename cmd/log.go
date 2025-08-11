package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the-startup/the-startup/internal/log"
)

// NewLogCommand creates the log command for hook processing
func NewLogCommand() *cobra.Command {
	var assistant, user bool

	cmd := &cobra.Command{
		Use:   "log",
		Short: "Process hook data from Claude Code",
		Long: `Process JSON hook data from Claude Code via stdin.
Use --assistant for PreToolUse hooks (agent_start events).
Use --user for PostToolUse hooks (agent_complete events).

This command reads JSON from stdin and writes JSONL logs to .the-startup/ directory.
On any error, exits silently with code 0 unless DEBUG_HOOKS is set.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate flags - exactly one must be specified
			if (assistant && user) || (!assistant && !user) {
				return fmt.Errorf("exactly one of --assistant or --user must be specified")
			}

			// Determine hook type
			isPostHook := user

			// Process the tool call from stdin
			hookData, err := log.ProcessToolCall(cmd.InOrStdin(), isPostHook)
			if err != nil {
				log.DebugError(err)
				// Silent exit on errors (matches Python behavior)
				return nil
			}

			// Skip processing if hookData is nil (filtered out)
			if hookData == nil {
				log.DebugLog("Tool call filtered out, skipping")
				return nil
			}

			// Write to session log
			if err := log.WriteSessionLog(hookData.SessionID, hookData); err != nil {
				log.DebugError(fmt.Errorf("failed to write session log: %w", err))
				// Continue to try global log even if session log fails
			}

			// Write to global log
			if err := log.WriteGlobalLog(hookData); err != nil {
				log.DebugError(fmt.Errorf("failed to write global log: %w", err))
				// Silent exit on errors
				return nil
			}

			log.DebugLog("Successfully processed %s event for agent %s", hookData.Event, hookData.AgentType)
			return nil
		},
		// Silent usage and errors to match hook behavior
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add flags
	cmd.Flags().BoolVarP(&assistant, "assistant", "a", false, "Process PreToolUse hook (agent_start event)")
	cmd.Flags().BoolVarP(&user, "user", "u", false, "Process PostToolUse hook (agent_complete event)")

	return cmd
}