package cmd

import (
	"fmt"

	"github.com/rsmdt/the-startup/internal/log"
	"github.com/spf13/cobra"
)

// NewLogCommand creates the log command for hook processing
func NewLogCommand() *cobra.Command {
	var flags log.LogFlags

	cmd := &cobra.Command{
		Use:   "log",
		Short: "Process hook data from Claude Code or read agent context",
		Long: `Process JSON hook data from Claude Code via stdin or read agent context.

Write modes:
  --assistant: Process PreToolUse hooks (agent_start events)
  --user:      Process PostToolUse hooks (agent_complete events)

Read mode:
  --read:      Read agent context with specified agent-id
  
Examples:
  # Process hook data
  echo '{"tool_name":"Task","tool_input":{"subagent_type":"the-architect"}}' | the-startup log --assistant
  
  # Read agent context  
  the-startup log --read --agent-id arch-001 --lines 20 --format json

This command reads JSON from stdin and writes JSONL logs to .the-startup/ directory.
On any error, exits silently with code 0 unless DEBUG_HOOKS is set.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate flags using the new validation function
			if err := log.ValidateLogCommand(flags); err != nil {
				return err
			}

			// Handle read mode
			if flags.Read {
				return handleReadMode(cmd, flags)
			}

			// Handle write modes (existing functionality)
			return handleWriteMode(cmd, flags)
		},
		// Silent usage and errors to match hook behavior
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add write mode flags
	cmd.Flags().BoolVarP(&flags.Assistant, "assistant", "a", false, "Process PreToolUse hook (agent_start event)")
	cmd.Flags().BoolVarP(&flags.User, "user", "u", false, "Process PostToolUse hook (agent_complete event)")

	// Add read mode flags
	cmd.Flags().BoolVarP(&flags.Read, "read", "r", false, "Read agent context (requires --agent-id)")
	cmd.Flags().StringVarP(&flags.AgentID, "agent-id", "i", "", "Agent ID to read context for (required with --read)")
	cmd.Flags().IntVar(&flags.Lines, "lines", 50, "Number of recent lines to return (1-1000, default 50)")
	cmd.Flags().StringVar(&flags.Session, "session", "", "Specific session ID (optional, defaults to latest)")
	cmd.Flags().StringVar(&flags.Format, "format", "json", "Output format: json or text (default json)")
	cmd.Flags().BoolVar(&flags.IncludeMetadata, "include-metadata", false, "Include file metadata in output (default false)")

	return cmd
}

// handleReadMode processes the --read flag and outputs agent context
func handleReadMode(cmd *cobra.Command, flags log.LogFlags) error {
	query := log.ContextQuery{
		AgentID:         flags.AgentID,
		SessionID:       flags.Session,
		MaxLines:        flags.Lines,
		Format:          flags.Format,
		IncludeMetadata: flags.IncludeMetadata,
	}

	result, err := log.ReadContext(query)
	if err != nil {
		// Silent error handling for hook compatibility
		log.DebugError(err)
		return nil
	}

	// Output the result
	fmt.Fprint(cmd.OutOrStdout(), result)
	return nil
}

// handleWriteMode processes the --assistant and --user flags (existing functionality)
func handleWriteMode(cmd *cobra.Command, flags log.LogFlags) error {
	// Determine hook type
	isPostHook := flags.User

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

	// Write to agent-specific context file (NEW)
	if err := log.WriteAgentContext(hookData.SessionID, hookData.AgentID, hookData); err != nil {
		log.DebugError(fmt.Errorf("failed to write agent context: %w", err))
		// Continue even if agent context write fails
	}

	// Note: Orchestrator routing removed - all agents get their own files based on AgentID

	// Write to session log (for backward compatibility)
	if err := log.WriteSessionLog(hookData.SessionID, hookData); err != nil {
		log.DebugError(fmt.Errorf("failed to write session log: %w", err))
		// Silent exit on errors
		return nil
	}

	log.DebugLog("Successfully processed %s role for agent %s", hookData.Role, hookData.AgentID)
	return nil
}
