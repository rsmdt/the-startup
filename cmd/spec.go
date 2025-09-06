package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// NewSpecCommand creates a new spec command
func NewSpecCommand(startupAssets *embed.FS) *cobra.Command {
	var (
		readMode string
		addMode  string
	)

	cmd := &cobra.Command{
		Use:   "spec [feature description or ID]",
		Short: "Manage specification directories",
		Long: `Creates new specification directories or manages existing ones.
This command is used by the /s:specify command to setup and manage specification documents.`,
		Example: `  # Create new specification
  the-startup spec "user authentication system"
  
  # Read existing specification
  the-startup spec --read 010
  the-startup spec --read 010-user-auth
  
  # Add template to specification
  the-startup spec 010 --add PRD
  the-startup spec 010-user-auth --add SDD`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Handle --read mode
			if readMode != "" {
				return handleReadMode(readMode)
			}

			// Handle --add mode
			if addMode != "" {
				if len(args) == 0 {
					return fmt.Errorf("spec ID required when using --add")
				}
				return handleAddMode(startupAssets, args[0], addMode)
			}

			// Default: create new spec
			if len(args) == 0 {
				return fmt.Errorf("feature description required")
			}
			
			description := strings.Join(args, " ")
			return handleCreateMode(startupAssets, description)
		},
	}

	cmd.Flags().StringVar(&readMode, "read", "", "Read existing specification by ID")
	cmd.Flags().StringVar(&addMode, "add", "", "Add template to specification (PRD, SDD, PLAN, BRD)")

	return cmd
}

// handleReadMode reads and displays an existing specification
func handleReadMode(specID string) error {
	specsDir := filepath.Join(".", "docs", "specs")
	
	// Find the spec directory
	specDir, id, name, err := findSpecDirectory(specsDir, specID)
	if err != nil {
		fmt.Printf("error = \"spec not found\"\n")
		fmt.Printf("id = %q\n", specID)
		return nil
	}

	// Scan for files
	files, err := scanSpecFiles(specDir)
	if err != nil {
		return fmt.Errorf("failed to scan spec directory: %w", err)
	}

	// Output in TOML format
	fmt.Printf("id = %q\n", id)
	fmt.Printf("name = %q\n", name)
	fmt.Printf("dir = %q\n", specDir)
	fmt.Println()
	fmt.Println("[files]")
	
	// Sort and output files
	keys := make([]string, 0, len(files))
	for k := range files {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	for _, key := range keys {
		fmt.Printf("%s = %q\n", key, files[key])
	}

	return nil
}

// handleAddMode adds a template to an existing specification
func handleAddMode(assets *embed.FS, specID, templateType string) error {
	// Validate template type
	templateType = strings.ToUpper(templateType)
	validTemplates := map[string]bool{
		"PRD": true, "SDD": true, "PLAN": true, "BRD": true,
	}
	if !validTemplates[templateType] {
		return fmt.Errorf("invalid template type: %s (valid: PRD, SDD, PLAN, BRD)", templateType)
	}

	specsDir := filepath.Join(".", "docs", "specs")
	
	// Try to find existing spec directory
	specDir, id, name, err := findSpecDirectory(specsDir, specID)
	if err != nil {
		// Directory doesn't exist, create it
		id = extractSpecID(specID)
		name = extractSpecName(specID)
		
		if id == "" {
			// If not a valid spec ID format, use as-is
			id = specID
			specDir = filepath.Join(specsDir, specID)
		} else if name != "" {
			specDir = filepath.Join(specsDir, fmt.Sprintf("%s-%s", id, name))
		} else {
			specDir = filepath.Join(specsDir, id)
		}
		
		// Extract name from directory if created
		if name == "" && strings.Contains(filepath.Base(specDir), "-") {
			parts := strings.SplitN(filepath.Base(specDir), "-", 2)
			if len(parts) > 1 {
				name = parts[1]
			}
		}
	}

	// Create directory if needed
	if err := os.MkdirAll(specDir, 0755); err != nil {
		return fmt.Errorf("failed to create spec directory: %w", err)
	}

	// Copy template
	destPath := filepath.Join(specDir, fmt.Sprintf("%s.md", templateType))
	
	// Check if file already exists
	if _, err := os.Stat(destPath); err == nil {
		return fmt.Errorf("file already exists: %s", destPath)
	}

	if err := copyTemplateFile(assets, templateType, destPath); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}

	// Output in TOML format with [files.new]
	fmt.Printf("id = %q\n", id)
	fmt.Printf("name = %q\n", name)
	fmt.Printf("dir = %q\n", specDir)
	fmt.Println()
	fmt.Println("[files.new]")
	fmt.Printf("%s = %q\n", strings.ToLower(templateType), destPath)

	return nil
}

// handleCreateMode creates a new specification
func handleCreateMode(assets *embed.FS, description string) error {
	// Create docs/specs directory if it doesn't exist
	specsDir := filepath.Join(".", "docs", "specs")
	if err := os.MkdirAll(specsDir, 0755); err != nil {
		return fmt.Errorf("failed to create specs directory: %w", err)
	}

	// Find the highest spec ID
	highestID, err := findHighestSpecID(specsDir)
	if err != nil {
		return fmt.Errorf("failed to find highest spec ID: %w", err)
	}

	// Generate next ID with 3-digit padding
	nextID := fmt.Sprintf("%03d", highestID+1)

	// Create feature name from description
	featureName := sanitizeFeatureName(description)

	// Create spec directory
	specDirName := fmt.Sprintf("%s-%s", nextID, featureName)
	specDir := filepath.Join(specsDir, specDirName)
	if err := os.MkdirAll(specDir, 0755); err != nil {
		return fmt.Errorf("failed to create spec directory: %w", err)
	}

	// Copy PRD template (without any modifications)
	prdPath := filepath.Join(specDir, "PRD.md")
	if err := copyTemplateFile(assets, "PRD", prdPath); err != nil {
		return fmt.Errorf("failed to copy PRD template: %w", err)
	}

	// Output in TOML format
	fmt.Printf("id = %q\n", nextID)
	fmt.Printf("name = %q\n", featureName)
	fmt.Printf("dir = %q\n", specDir)
	fmt.Println()
	fmt.Println("[files]")
	fmt.Printf("prd = %q\n", prdPath)

	return nil
}

// findSpecDirectory finds a spec directory by ID or full name
func findSpecDirectory(specsDir, specID string) (dir, id, name string, err error) {
	// Check if specID is numeric (e.g., "010")
	if matched, _ := regexp.MatchString(`^\d{3}$`, specID); matched {
		// Look for directory starting with this ID
		entries, readErr := os.ReadDir(specsDir)
		if readErr != nil {
			return "", "", "", fmt.Errorf("failed to read specs directory: %w", readErr)
		}

		prefix := specID + "-"
		for _, entry := range entries {
			if entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) {
				dir = filepath.Join(specsDir, entry.Name())
				id = specID
				name = strings.TrimPrefix(entry.Name(), prefix)
				return dir, id, name, nil
			}
		}
		
		// Also check for exact match (e.g., "010" directory)
		exactPath := filepath.Join(specsDir, specID)
		if info, err := os.Stat(exactPath); err == nil && info.IsDir() {
			return exactPath, specID, "", nil
		}
	} else {
		// Check for exact directory name match
		fullPath := filepath.Join(specsDir, specID)
		if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
			// Extract ID and name from directory name
			id = extractSpecID(specID)
			name = extractSpecName(specID)
			return fullPath, id, name, nil
		}
	}

	return "", "", "", fmt.Errorf("specification not found: %s", specID)
}

// extractSpecID extracts the numeric ID from a spec directory name
func extractSpecID(dirName string) string {
	re := regexp.MustCompile(`^(\d{3})`)
	matches := re.FindStringSubmatch(dirName)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractSpecName extracts the feature name from a spec directory name
func extractSpecName(dirName string) string {
	re := regexp.MustCompile(`^\d{3}-(.+)$`)
	matches := re.FindStringSubmatch(dirName)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// scanSpecFiles scans a spec directory for documentation files
func scanSpecFiles(specDir string) (map[string]string, error) {
	files := make(map[string]string)
	
	// Check for each known template type
	templates := []string{"PRD", "SDD", "PLAN", "BRD"}
	for _, tmpl := range templates {
		path := filepath.Join(specDir, fmt.Sprintf("%s.md", tmpl))
		if _, err := os.Stat(path); err == nil {
			files[strings.ToLower(tmpl)] = path
		}
	}

	// Also check for other .md files
	entries, err := os.ReadDir(specDir)
	if err != nil {
		return files, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			baseName := strings.TrimSuffix(entry.Name(), ".md")
			// Skip if already added as standard template
			lowerName := strings.ToLower(baseName)
			if _, exists := files[lowerName]; !exists {
				files[lowerName] = filepath.Join(specDir, entry.Name())
			}
		}
	}

	return files, nil
}

// findHighestSpecID scans the specs directory for the highest numbered specification
func findHighestSpecID(specsDir string) (int, error) {
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil // No specs yet
		}
		return 0, err
	}

	highest := 0
	re := regexp.MustCompile(`^(\d{3})`)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		matches := re.FindStringSubmatch(entry.Name())
		if len(matches) > 1 {
			id, err := strconv.Atoi(matches[1])
			if err == nil && id > highest {
				highest = id
			}
		}
	}

	return highest, nil
}

// sanitizeFeatureName converts a description into a valid directory name
func sanitizeFeatureName(description string) string {
	// Convert to lowercase
	name := strings.ToLower(description)

	// Replace non-alphanumeric with hyphens
	re := regexp.MustCompile(`[^a-z0-9]+`)
	name = re.ReplaceAllString(name, "-")

	// Remove leading/trailing hyphens
	name = strings.Trim(name, "-")

	// Limit to 3-4 meaningful words
	words := strings.Split(name, "-")
	if len(words) > 4 {
		words = words[:4]
	}

	// Remove common words
	filtered := []string{}
	skipWords := map[string]bool{
		"a": true, "an": true, "the": true, "and": true, "or": true,
		"for": true, "with": true, "to": true, "of": true, "in": true,
	}

	for _, word := range words {
		if !skipWords[word] && word != "" {
			filtered = append(filtered, word)
		}
	}

	if len(filtered) == 0 {
		return "feature"
	}

	return strings.Join(filtered, "-")
}

// copyTemplateFile copies a template file without any modifications
func copyTemplateFile(assets *embed.FS, templateName, destPath string) error {
	// Read template bytes
	templatePath := fmt.Sprintf("assets/the-startup/templates/%s.md", templateName)
	content, err := assets.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read %s template: %w", templateName, err)
	}

	// Create directory if needed
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write exact bytes (no modification)
	if err := os.WriteFile(destPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}