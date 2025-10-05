package cmd

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// SetupAnswers holds the user's configuration choices
type SetupAnswers struct {
	UsesTDD        bool
	CoverageTarget int
	BuildCommand   string
	TestCommand    string
	LintCommand    string
	FormatCommand  string
}

// NewInitCommand creates the init command
func NewInitCommand(startupAssets *embed.FS) *cobra.Command {
	var skipPrompts bool
	var force bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "init [template]",
		Short: "Initialize quality gate templates (definition-of-ready, definition-of-done, task-definition-of-done)",
		Long: `Initialize Definition of Ready and Definition of Done templates
for validating document creation and task completion in The Startup workflow.

Examples:
  # Initialize all templates
  the-startup init

  # Check what would be created (dry-run)
  the-startup init --dry-run

  # Initialize specific template
  the-startup init definition-of-ready
  the-startup init definition-of-done
  the-startup init task-definition-of-done

  # Skip prompts and use defaults
  the-startup init --skip-prompts`,
		RunE: func(cmd *cobra.Command, args []string) error {
			template := ""
			if len(args) > 0 {
				template = args[0]
			}
			return runInit(startupAssets, template, skipPrompts, force, dryRun)
		},
	}

	cmd.Flags().BoolVarP(&skipPrompts, "skip-prompts", "s", false, "Skip guided setup questions, use defaults")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing files without prompting")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Check what files exist without creating/overwriting")

	return cmd
}

func runInit(assets *embed.FS, template string, skipPrompts, force, dryRun bool) error {
	// Determine which templates to process
	templates := []string{"definition-of-ready", "definition-of-done", "task-definition-of-done"}
	if template != "" {
		// Validate template name (case-insensitive)
		valid := false
		normalizedTemplate := strings.ToLower(template)
		for _, t := range templates {
			if normalizedTemplate == t {
				valid = true
				template = t // Use the normalized form
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid template name: %s (valid: definition-of-ready, definition-of-done, task-definition-of-done)", template)
		}
		templates = []string{strings.ToLower(template)}
	}

	// Step 1: Check existing files
	existing := make(map[string]bool)
	for _, tmpl := range templates {
		filePath := fmt.Sprintf("./docs/%s.md", tmpl)
		if _, err := os.Stat(filePath); err == nil {
			existing[tmpl] = true
		}
	}

	// If dry-run, just report status and exit
	if dryRun {
		fmt.Println("📋 Validation Template Status:")
		fmt.Println()
		for _, tmpl := range templates {
			filePath := fmt.Sprintf("docs/%s.md", tmpl)
			if existing[tmpl] {
				fmt.Printf("  ✓ %s exists\n", filePath)
			} else {
				fmt.Printf("  ○ %s not found\n", filePath)
			}
		}
		fmt.Println()
		if len(existing) > 0 {
			fmt.Println("💡 Use --force to overwrite existing files")
			fmt.Println("💡 Or initialize specific templates: the-startup init definition-of-ready")
		}
		return nil
	}

	// Step 2: Create docs directory
	if err := os.MkdirAll("./docs", 0755); err != nil {
		return fmt.Errorf("failed to create docs directory: %w", err)
	}

	// Step 3: Check for overwrites (unless force flag)
	if !force && len(existing) > 0 {
		fmt.Printf("Files already exist:\n")
		for tmpl := range existing {
			fmt.Printf("  • docs/%s.md\n", tmpl)
		}
		fmt.Println()
		if !promptYesNo("Overwrite existing files?", false) {
			return fmt.Errorf("cancelled by user")
		}
	}

	// Step 4: Copy templates
	if template == "" {
		fmt.Println("Copying validation templates...")
	} else {
		fmt.Printf("Copying %s template...\n", template)
	}

	for _, tmpl := range templates {
		templateFile := fmt.Sprintf("%s.md", tmpl)
		destPath := fmt.Sprintf("./docs/%s.md", tmpl)

		if err := copyTemplate(assets, templateFile, destPath); err != nil {
			return err
		}
		fmt.Printf("✓ Created docs/%s.md\n", tmpl)
	}
	fmt.Println()

	// Step 5: Run guided prompts (only if initializing all templates and not skipped)
	var answers SetupAnswers
	if template == "" && !skipPrompts {
		answers = runGuidedPrompts()
		displayAdvice(answers)
	} else if skipPrompts {
		fmt.Println("Skipped guided setup. Using default values.")
		fmt.Println()
	}

	// Step 6: Display next steps
	if template == "" {
		displayNextSteps(answers, skipPrompts)
	} else {
		fmt.Printf("💡 Next: Review docs/%s.md and customize [NEEDS CLARIFICATION: ...] markers\n", template)
		fmt.Println()
	}

	return nil
}

func copyTemplate(assets *embed.FS, templateName, destPath string) error {
	// Read template from embedded assets
	templatePath := filepath.Join("assets/the-startup/templates", templateName)
	data, err := assets.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templateName, err)
	}

	// Write to destination
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", destPath, err)
	}

	return nil
}

func runGuidedPrompts() SetupAnswers {
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Templates created successfully!                             │")
	fmt.Println("│                                                              │")
	fmt.Println("│ Let's customize them for your workflow.                     │")
	fmt.Println("│ (You can skip and edit docs/DOR.md and docs/DOD.md later)  │")
	fmt.Println("└─────────────────────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Println("❓ Quick setup questions (press Enter to use defaults):")
	fmt.Println()

	var answers SetupAnswers

	// Question 1: TDD
	answers.UsesTDD = promptYesNo("Do you use Test-Driven Development (TDD)?", false)

	// Question 2: Coverage
	answers.CoverageTarget = promptInt("What's your test coverage target?", 80, 0, 100)

	// Question 3: Build command
	answers.BuildCommand = promptString("What command builds your project?", "go build ./...")

	// Question 4: Test command
	answers.TestCommand = promptString("What command runs your tests?", "go test ./...")

	// Question 5: Lint command
	answers.LintCommand = promptString("What command runs linting?", "golangci-lint run")

	// Question 6: Format command
	answers.FormatCommand = promptString("What command checks formatting?", "gofmt -l .")

	fmt.Println()

	return answers
}

func promptYesNo(question string, defaultValue bool) bool {
	defaultStr := "y/N"
	if defaultValue {
		defaultStr = "Y/n"
	}

	fmt.Printf("%s [%s]: ", question, defaultStr)

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	if response == "" {
		return defaultValue
	}

	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}

func promptInt(question string, defaultValue, min, max int) int {
	fmt.Printf("%s [%d]: ", question, defaultValue)

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	if response == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(response)
	if err != nil || value < min || value > max {
		fmt.Printf("Invalid input, using default: %d\n", defaultValue)
		return defaultValue
	}

	return value
}

func promptString(question, defaultValue string) string {
	fmt.Printf("%s [%s]: ", question, defaultValue)

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	if response == "" {
		return defaultValue
	}

	return response
}

func displayAdvice(answers SetupAnswers) {
	hasCustomizations := false

	fmt.Println("📝 Customization Advice:")
	fmt.Println()

	// TDD advice
	if !answers.UsesTDD {
		hasCustomizations = true
		fmt.Println("  ⚠️  TDD is disabled")
		fmt.Println("     Edit docs/definition-of-done.md and docs/task-definition-of-done.md:")
		fmt.Println("     • Remove or comment out sections marked with:")
		fmt.Println("       <!-- OPTIONAL: TDD -->")
		fmt.Println("       (everything between <!-- OPTIONAL: TDD --> and <!-- END OPTIONAL: TDD -->)")
		fmt.Println()
	} else {
		fmt.Println("  ✓ TDD is enabled")
		fmt.Println("    The templates include TDD RED→GREEN cycle enforcement")
		fmt.Println()
	}

	// Coverage advice
	if answers.CoverageTarget != 80 {
		hasCustomizations = true
		fmt.Printf("  📊 Coverage target: %d%%\n", answers.CoverageTarget)
		fmt.Println("     Edit docs/definition-of-done.md and docs/task-definition-of-done.md:")
		fmt.Println("     • Find: [NEEDS CLARIFICATION: coverage target]")
		fmt.Printf("     • Replace with: %d\n", answers.CoverageTarget)
		fmt.Println()
	}

	// Build command advice
	if answers.BuildCommand != "go build ./..." {
		hasCustomizations = true
		fmt.Printf("  🔨 Build command: %s\n", answers.BuildCommand)
		fmt.Println("     Edit docs/task-definition-of-done.md:")
		fmt.Println("     • Find: [NEEDS CLARIFICATION: build command]")
		fmt.Printf("     • Replace with: %s\n", answers.BuildCommand)
		fmt.Println()
	}

	// Test command advice
	if answers.TestCommand != "go test ./..." {
		hasCustomizations = true
		fmt.Printf("  🧪 Test command: %s\n", answers.TestCommand)
		fmt.Println("     Edit docs/definition-of-ready.md and docs/task-definition-of-done.md:")
		fmt.Println("     • Find: [NEEDS CLARIFICATION: test command]")
		fmt.Printf("     • Replace with: %s\n", answers.TestCommand)
		fmt.Println()
	}

	// Lint command advice
	if answers.LintCommand != "golangci-lint run" {
		hasCustomizations = true
		fmt.Printf("  🔍 Lint command: %s\n", answers.LintCommand)
		fmt.Println("     Edit docs/task-definition-of-done.md:")
		fmt.Println("     • Find: [NEEDS CLARIFICATION: lint command]")
		fmt.Printf("     • Replace with: %s\n", answers.LintCommand)
		fmt.Println()
	}

	// Format command advice
	if answers.FormatCommand != "gofmt -l ." {
		hasCustomizations = true
		fmt.Printf("  💅 Format command: %s\n", answers.FormatCommand)
		fmt.Println("     Edit docs/task-definition-of-done.md:")
		fmt.Println("     • Find: [NEEDS CLARIFICATION: format command]")
		fmt.Printf("     • Replace with: %s\n", answers.FormatCommand)
		fmt.Println()
	}

	// If using all defaults
	if !hasCustomizations {
		fmt.Println("  ✓ Using all defaults")
		fmt.Println("    Templates are ready to use as-is!")
		fmt.Println("    Review and customize later if needed.")
		fmt.Println()
	}
}

func displayNextSteps(answers SetupAnswers, skipPrompts bool) {
	fmt.Println("✅ Initialization complete!")
	fmt.Println()
	fmt.Println("📂 Files created:")
	fmt.Println("   • docs/definition-of-ready.md - Definition of Ready")
	fmt.Println("   • docs/definition-of-done.md - Definition of Done")
	fmt.Println("   • docs/task-definition-of-done.md - Task Definition of Done")
	fmt.Println()
	fmt.Println("📖 What's next:")
	fmt.Println()
	fmt.Println("   1. Review the templates:")
	fmt.Println("      cat docs/definition-of-ready.md")
	fmt.Println("      cat docs/definition-of-done.md")
	fmt.Println("      cat docs/task-definition-of-done.md")
	fmt.Println()

	if !skipPrompts {
		fmt.Println("   2. Apply the customizations shown above")
		fmt.Println()
		fmt.Println("   3. Start using in your workflow:")
	} else {
		fmt.Println("   2. Customize the templates:")
		fmt.Println("      • Replace [NEEDS CLARIFICATION: ...] markers")
		fmt.Println("      • Remove <!-- OPTIONAL: ... --> sections if not needed")
		fmt.Println()
		fmt.Println("   3. Start using in your workflow:")
	}

	fmt.Println("      /s:specify \"your feature description\"")
	fmt.Println("      /s:implement <spec-id>")
	fmt.Println()
	fmt.Println("💡 Tip: The templates include inline comments explaining each section.")
	fmt.Println()
}
