package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ErrorType represents different categories of errors
type ErrorType int

const (
	ErrorTypePermission ErrorType = iota
	ErrorTypeNotFound
	ErrorTypeAlreadyExists
	ErrorTypeNetwork
	ErrorTypeConfig
	ErrorTypeGeneral
)

// ErrorContext provides detailed error information
type ErrorContext struct {
	Type        ErrorType
	Title       string
	Description string
	Context     map[string]string
	Suggestions []string
}

// RenderError creates a formatted error message with context
func RenderError(err error, context *ErrorContext) string {
	styles := GetStyles()
	var b strings.Builder

	// Error header
	b.WriteString("\n")
	b.WriteString(styles.Error.Bold(true).Render(fmt.Sprintf("%s %s", IconError, context.Title)))
	b.WriteString("\n\n")

	// Error description
	if context.Description != "" {
		b.WriteString(styles.Normal.Render(context.Description))
		b.WriteString("\n\n")
	}

	// Original error
	b.WriteString(styles.Help.Render("Error details:"))
	b.WriteString("\n")
	b.WriteString(styles.Error.Render(fmt.Sprintf("  %s", err.Error())))
	b.WriteString("\n")

	// Context information
	if len(context.Context) > 0 {
		b.WriteString("\n")
		b.WriteString(styles.Help.Render("Context:"))
		b.WriteString("\n")
		for key, value := range context.Context {
			b.WriteString(fmt.Sprintf("  %s: %s\n",
				styles.Info.Render(key),
				styles.Normal.Render(value)))
		}
	}

	// Suggestions
	if len(context.Suggestions) > 0 {
		b.WriteString("\n")
		b.WriteString(styles.Success.Render("Possible solutions:"))
		b.WriteString("\n")
		for i, suggestion := range context.Suggestions {
			b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, suggestion))
		}
	}

	b.WriteString("\n")
	return b.String()
}

// GetErrorContext analyzes an error and provides context
func GetErrorContext(err error, operation string) *ErrorContext {
	errStr := err.Error()

	// Permission denied
	if strings.Contains(errStr, "permission denied") || strings.Contains(errStr, "access denied") {
		return &ErrorContext{
			Type:        ErrorTypePermission,
			Title:       "Permission Denied",
			Description: fmt.Sprintf("Unable to %s due to insufficient permissions.", operation),
			Context: map[string]string{
				"Operation": operation,
				"User":      getCurrentUser(),
			},
			Suggestions: []string{
				"Run the command with appropriate permissions (e.g., sudo)",
				"Check if the target directory is writable",
				"Verify you own the files/directories involved",
				"Try installing to a different location (use --path flag)",
			},
		}
	}

	// File/Directory not found
	if strings.Contains(errStr, "no such file") || strings.Contains(errStr, "not found") {
		return &ErrorContext{
			Type:        ErrorTypeNotFound,
			Title:       "File or Directory Not Found",
			Description: fmt.Sprintf("Required file or directory was not found while %s.", operation),
			Context: map[string]string{
				"Operation": operation,
			},
			Suggestions: []string{
				"Verify the path exists and is accessible",
				"Check for typos in the file/directory name",
				"Ensure you're in the correct directory",
				"Run 'the-startup install' to set up missing components",
			},
		}
	}

	// Already exists
	if strings.Contains(errStr, "already exists") || strings.Contains(errStr, "file exists") {
		return &ErrorContext{
			Type:        ErrorTypeAlreadyExists,
			Title:       "File Already Exists",
			Description: fmt.Sprintf("Cannot %s because the target already exists.", operation),
			Context: map[string]string{
				"Operation": operation,
			},
			Suggestions: []string{
				"Use 'the-startup update' to update existing installation",
				"Remove the existing file/directory first",
				"Choose a different installation location with --path",
				"Use --force flag to overwrite (if available)",
			},
		}
	}

	// Network errors
	if strings.Contains(errStr, "connection") || strings.Contains(errStr, "network") {
		return &ErrorContext{
			Type:        ErrorTypeNetwork,
			Title:       "Network Error",
			Description: fmt.Sprintf("Network issue encountered while %s.", operation),
			Context: map[string]string{
				"Operation": operation,
			},
			Suggestions: []string{
				"Check your internet connection",
				"Verify proxy settings if behind a firewall",
				"Try again in a few moments",
				"Check if the remote server is accessible",
			},
		}
	}

	// Configuration errors
	if strings.Contains(errStr, "config") || strings.Contains(errStr, "settings") {
		return &ErrorContext{
			Type:        ErrorTypeConfig,
			Title:       "Configuration Error",
			Description: fmt.Sprintf("Configuration issue detected while %s.", operation),
			Context: map[string]string{
				"Operation": operation,
			},
			Suggestions: []string{
				"Check your configuration file for errors",
				"Run 'the-startup validate' to check configuration",
				"Reset to default configuration",
				"Review the documentation for configuration options",
			},
		}
	}

	// Generic error
	return &ErrorContext{
		Type:        ErrorTypeGeneral,
		Title:       "Operation Failed",
		Description: fmt.Sprintf("An error occurred while %s.", operation),
		Context: map[string]string{
			"Operation": operation,
		},
		Suggestions: []string{
			"Check the error message for more details",
			"Ensure all requirements are met",
			"Try running with --verbose flag for more information",
			"Report this issue if it persists",
		},
	}
}

// getCurrentUser gets the current username
func getCurrentUser() string {
	if user := os.Getenv("USER"); user != "" {
		return user
	}
	if user := os.Getenv("USERNAME"); user != "" {
		return user
	}
	return "unknown"
}

// ErrorBox creates a bordered error message
func ErrorBox(title, message string) string {
	styles := GetStyles()

	errorStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Error.GetForeground()).
		Padding(1, 2).
		Width(60)

	content := fmt.Sprintf("%s %s\n\n%s", IconError, title, message)
	return errorStyle.Render(content)
}

// WarningBox creates a bordered warning message
func WarningBox(title, message string) string {
	styles := GetStyles()

	warningStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Warning.GetForeground()).
		Padding(1, 2).
		Width(60)

	content := fmt.Sprintf("%s %s\n\n%s", IconWarning, title, message)
	return warningStyle.Render(content)
}

// InfoBox creates a bordered info message
func InfoBox(title, message string) string {
	styles := GetStyles()

	infoStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Info.GetForeground()).
		Padding(1, 2).
		Width(60)

	content := fmt.Sprintf("%s %s\n\n%s", IconInfo, title, message)
	return infoStyle.Render(content)
}

// SuccessBox creates a bordered success message
func SuccessBox(title, message string) string {
	styles := GetStyles()

	successStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Success.GetForeground()).
		Padding(1, 2).
		Width(60)

	content := fmt.Sprintf("%s %s\n\n%s", IconSuccess, title, message)
	return successStyle.Render(content)
}
