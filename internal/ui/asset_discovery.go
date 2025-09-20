package ui

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// SelectableItem represents an agent, command, or other component that can be selected for installation
type SelectableItem struct {
	ID          string // unique identifier (file path)
	Name        string // display name extracted from metadata
	Description string // description from metadata
	Type        ItemType
	Category    string // parent category (role for agents, group for commands)
	FilePath    string // path in embedded assets
	Selected    bool   // current selection state
	Required    bool   // cannot be deselected (e.g., settings)
}

// ItemType represents the type of selectable item
type ItemType string

const (
	ItemTypeAgent       ItemType = "agent"
	ItemTypeCommand     ItemType = "command"
	ItemTypeOutputStyle ItemType = "output-style"
	ItemTypeSettings    ItemType = "settings"
)

// AssetDiscovery handles dynamic discovery of agents, commands, and other assets
type AssetDiscovery struct {
	claudeAssets  *embed.FS
	startupAssets *embed.FS
}

// NewAssetDiscovery creates a new asset discovery instance
func NewAssetDiscovery(claudeAssets, startupAssets *embed.FS) *AssetDiscovery {
	return &AssetDiscovery{
		claudeAssets:  claudeAssets,
		startupAssets: startupAssets,
	}
}

// DiscoverAllItems discovers all selectable items from embedded assets
func (ad *AssetDiscovery) DiscoverAllItems() ([]SelectableItem, error) {
	var items []SelectableItem

	// Discover agents
	agents, err := ad.DiscoverAgents()
	if err != nil {
		return nil, fmt.Errorf("failed to discover agents: %w", err)
	}
	items = append(items, agents...)

	// Discover commands
	commands, err := ad.DiscoverCommands()
	if err != nil {
		return nil, fmt.Errorf("failed to discover commands: %w", err)
	}
	items = append(items, commands...)

	// Discover output styles
	outputStyles, err := ad.DiscoverOutputStyles()
	if err != nil {
		return nil, fmt.Errorf("failed to discover output styles: %w", err)
	}
	items = append(items, outputStyles...)

	// Add required settings items
	settingsItems := ad.DiscoverSettings()
	items = append(items, settingsItems...)

	return items, nil
}

// DiscoverAgents discovers all agent files from assets/claude/agents
func (ad *AssetDiscovery) DiscoverAgents() ([]SelectableItem, error) {
	var agents []SelectableItem

	if ad.claudeAssets == nil {
		return agents, nil
	}

	err := fs.WalkDir(ad.claudeAssets, "assets/claude/agents", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors, continue walking
		}

		// Only process .md files
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Extract relative path from assets/claude/
		relPath := strings.TrimPrefix(path, "assets/claude/")

		// Parse metadata from the file
		metadata, err := ad.parseAgentMetadata(path)
		if err != nil {
			// If parsing fails, create basic item from filename
			metadata = AgentMetadata{
				Name:        ad.extractNameFromPath(path),
				Description: "Agent for specialized development tasks",
			}
		}

		// Determine category from directory structure
		category := ad.extractAgentCategory(path)

		agents = append(agents, SelectableItem{
			ID:          relPath,
			Name:        metadata.Name,
			Description: metadata.Description,
			Type:        ItemTypeAgent,
			Category:    category,
			FilePath:    relPath,
			Selected:    true, // Default to selected
			Required:    false,
		})

		return nil
	})

	// Sort agents by category and name for consistent display
	sort.Slice(agents, func(i, j int) bool {
		if agents[i].Category != agents[j].Category {
			return agents[i].Category < agents[j].Category
		}
		return agents[i].Name < agents[j].Name
	})

	return agents, err
}

// DiscoverCommands discovers all command files from assets/claude/commands
func (ad *AssetDiscovery) DiscoverCommands() ([]SelectableItem, error) {
	var commands []SelectableItem

	if ad.claudeAssets == nil {
		return commands, nil
	}

	err := fs.WalkDir(ad.claudeAssets, "assets/claude/commands", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors, continue walking
		}

		// Only process .md files
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Extract relative path from assets/claude/
		relPath := strings.TrimPrefix(path, "assets/claude/")

		// Parse metadata from the command file
		metadata, err := ad.parseCommandMetadata(path)
		if err != nil {
			// If parsing fails, create basic item from filename
			metadata = CommandMetadata{
				Description: "Command for development workflow",
			}
		}

		// Convert file path to command name (e.g., s/specify.md -> /s:specify)
		commandName := ad.extractCommandName(path)

		commands = append(commands, SelectableItem{
			ID:          relPath,
			Name:        commandName,
			Description: metadata.Description,
			Type:        ItemTypeCommand,
			Category:    "Development Commands",
			FilePath:    relPath,
			Selected:    true, // Default to selected
			Required:    false,
		})

		return nil
	})

	// Sort commands by name for consistent display
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	return commands, err
}

// DiscoverOutputStyles discovers output style files
func (ad *AssetDiscovery) DiscoverOutputStyles() ([]SelectableItem, error) {
	var styles []SelectableItem

	if ad.claudeAssets == nil {
		return styles, nil
	}

	err := fs.WalkDir(ad.claudeAssets, "assets/claude/output-styles", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors, continue walking
		}

		if d.IsDir() {
			return nil
		}

		// Extract relative path from assets/claude/
		relPath := strings.TrimPrefix(path, "assets/claude/")

		// Create display name from filename
		name := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
		name = strings.ReplaceAll(name, "-", " ")
		name = strings.ReplaceAll(name, "_", " ")

		styles = append(styles, SelectableItem{
			ID:          relPath,
			Name:        name,
			Description: "Output styling configuration",
			Type:        ItemTypeOutputStyle,
			Category:    "Output Styles",
			FilePath:    relPath,
			Selected:    true, // Default to selected
			Required:    false,
		})

		return nil
	})

	return styles, err
}

// DiscoverSettings creates settings items (these are always required)
func (ad *AssetDiscovery) DiscoverSettings() []SelectableItem {
	return []SelectableItem{
		{
			ID:          "settings.json",
			Name:        "Claude Settings",
			Description: "Main Claude configuration file",
			Type:        ItemTypeSettings,
			Category:    "Configuration",
			FilePath:    "settings.json",
			Selected:    true,
			Required:    true, // Always required
		},
		{
			ID:          "settings.local.json",
			Name:        "Local Settings",
			Description: "Local Claude configuration overrides",
			Type:        ItemTypeSettings,
			Category:    "Configuration",
			FilePath:    "settings.local.json",
			Selected:    true,
			Required:    true, // Always required
		},
	}
}

// AgentMetadata represents parsed agent metadata
type AgentMetadata struct {
	Name        string
	Description string
}

// CommandMetadata represents parsed command metadata
type CommandMetadata struct {
	Description string
}

// parseAgentMetadata extracts metadata from agent file frontmatter
func (ad *AssetDiscovery) parseAgentMetadata(path string) (AgentMetadata, error) {
	var metadata AgentMetadata

	data, err := ad.claudeAssets.ReadFile(path)
	if err != nil {
		return metadata, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	// Look for YAML frontmatter
	inFrontmatter := false
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			} else {
				break // End of frontmatter
			}
		}

		if !inFrontmatter {
			continue
		}

		// Parse YAML-like fields
		if strings.HasPrefix(line, "name:") {
			metadata.Name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
		} else if strings.HasPrefix(line, "description:") {
			// Handle multiline descriptions
			desc := strings.TrimSpace(strings.TrimPrefix(line, "description:"))
			if strings.HasPrefix(desc, "\"") && strings.HasSuffix(desc, "\"") {
				desc = strings.Trim(desc, "\"")
			}
			// Take only the first sentence for brevity
			if idx := strings.Index(desc, "."); idx > 0 && idx < 100 {
				desc = desc[:idx+1]
			}
			metadata.Description = desc
		}
	}

	// If no name found, extract from filename
	if metadata.Name == "" {
		metadata.Name = ad.extractNameFromPath(path)
	}

	return metadata, nil
}

// parseCommandMetadata extracts metadata from command file frontmatter
func (ad *AssetDiscovery) parseCommandMetadata(path string) (CommandMetadata, error) {
	var metadata CommandMetadata

	data, err := ad.claudeAssets.ReadFile(path)
	if err != nil {
		return metadata, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	inFrontmatter := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			} else {
				break // End of frontmatter
			}
		}

		if !inFrontmatter {
			continue
		}

		if strings.HasPrefix(line, "description:") {
			desc := strings.TrimSpace(strings.TrimPrefix(line, "description:"))
			if strings.HasPrefix(desc, "\"") && strings.HasSuffix(desc, "\"") {
				desc = strings.Trim(desc, "\"")
			}
			metadata.Description = desc
			break
		}
	}

	return metadata, scanner.Err()
}

// extractNameFromPath extracts a display name from file path
func (ad *AssetDiscovery) extractNameFromPath(path string) string {
	// Get filename without extension
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	// Convert hyphens to spaces and title case
	name = strings.ReplaceAll(name, "-", " ")
	parts := strings.Fields(name)
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}

	return strings.Join(parts, " ")
}

// extractAgentCategory determines the category from agent directory structure
func (ad *AssetDiscovery) extractAgentCategory(path string) string {
	// Extract directory structure from path
	// assets/claude/agents/the-analyst/requirements-documentation.md -> "the-analyst"
	parts := strings.Split(path, "/")

	if len(parts) >= 4 && parts[2] == "agents" {
		category := parts[3]
		if strings.HasPrefix(category, "the-") {
			// Convert the-analyst to "The Analyst"
			category = strings.TrimPrefix(category, "the-")
			category = strings.ReplaceAll(category, "-", " ")
			words := strings.Fields(category)
			for i, word := range words {
				words[i] = strings.Title(word)
			}
			return "The " + strings.Join(words, " ")
		}
		return category
	}

	return "General Agents"
}

// extractCommandName converts file path to command name
func (ad *AssetDiscovery) extractCommandName(path string) string {
	// Convert s/specify.md to /s:specify
	relPath := strings.TrimPrefix(path, "assets/claude/commands/")
	relPath = strings.TrimSuffix(relPath, ".md")

	// Replace directory separators with colons
	parts := strings.Split(relPath, "/")
	return "/" + strings.Join(parts, ":")
}

// GetSelectedFilePaths converts selected items back to file paths for installer
func (ad *AssetDiscovery) GetSelectedFilePaths(items []SelectableItem) []string {
	var paths []string
	for _, item := range items {
		if item.Selected {
			paths = append(paths, item.FilePath)
		}
	}
	return paths
}