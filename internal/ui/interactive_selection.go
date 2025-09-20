package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// InteractiveSelector provides Huh-based selection for agents and commands
type InteractiveSelector struct {
	discovery *AssetDiscovery
	items     []SelectableItem
	style     *SelectionStyle
}

// SelectionStyle defines styling for the selection interface
type SelectionStyle struct {
	Title       lipgloss.Style
	Category    lipgloss.Style
	Item        lipgloss.Style
	Description lipgloss.Style
	Selected    lipgloss.Style
	Help        lipgloss.Style
}

// NewInteractiveSelector creates a new interactive selector
func NewInteractiveSelector(discovery *AssetDiscovery) *InteractiveSelector {
	return &InteractiveSelector{
		discovery: discovery,
		style:     NewSelectionStyle(),
	}
}

// NewSelectionStyle creates styled selection interface
func NewSelectionStyle() *SelectionStyle {
	return &SelectionStyle{
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true).
			Margin(1, 0),
		Category: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Margin(1, 0, 0, 2),
		Item: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Margin(0, 0, 0, 4),
		Description: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Margin(0, 0, 0, 6),
		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")).
			Bold(true),
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Margin(1, 0),
	}
}

// RunInteractiveSelection presents multi-step selection using Huh forms
func (is *InteractiveSelector) RunInteractiveSelection() ([]SelectableItem, error) {
	// Discover all available items
	items, err := is.discovery.DiscoverAllItems()
	if err != nil {
		return nil, fmt.Errorf("failed to discover items: %w", err)
	}

	is.items = items

	// Run selection flow
	selectedItems, err := is.runSelectionFlow()
	if err != nil {
		return nil, fmt.Errorf("selection failed: %w", err)
	}

	return selectedItems, nil
}

// runSelectionFlow orchestrates the multi-step selection process
func (is *InteractiveSelector) runSelectionFlow() ([]SelectableItem, error) {
	// Step 1: Choose selection approach
	approach, err := is.selectApproach()
	if err != nil {
		return nil, err
	}

	switch approach {
	case "custom":
		return is.customSelection()
	case "role-based":
		return is.roleBasedSelection()
	case "recommended":
		return is.recommendedSelection()
	case "minimal":
		return is.minimalSelection()
	default:
		return is.recommendedSelection()
	}
}

// selectApproach presents the initial approach selection
func (is *InteractiveSelector) selectApproach() (string, error) {
	var approach string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🚀 How would you like to select components?").
				Description("Choose your installation approach for The (Agentic) Startup").
				Options(
					huh.NewOption("🎯 Recommended Setup - Essential agents and commands for most developers", "recommended"),
					huh.NewOption("👤 Role-based Selection - Select components based on your development role", "role-based"),
					huh.NewOption("✨ Custom Selection - Manually choose specific agents and commands", "custom"),
					huh.NewOption("📦 Minimal Install - Only essential components (settings + core commands)", "minimal"),
				).
				Value(&approach),
		),
	).WithTheme(huh.ThemeCharm())

	err := form.Run()
	return approach, err
}

// recommendedSelection provides a curated recommended set
func (is *InteractiveSelector) recommendedSelection() ([]SelectableItem, error) {
	// Get recommended items (most commonly used agents + all commands)
	recommended := is.getRecommendedItems()

	// Show what will be installed
	var confirm bool
	items := is.formatItemsForDisplay(recommended)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("🎯 Recommended Installation").
				Description(fmt.Sprintf("This will install %d carefully selected components for The (Agentic) Startup:\n\n%s",
					len(recommended), items)),
			huh.NewConfirm().
				Title("🚀 Install recommended components?").
				Value(&confirm),
		),
	).WithTheme(huh.ThemeCharm())

	err := form.Run()
	if err != nil {
		return nil, err
	}

	if !confirm {
		// User declined, go back to approach selection
		return is.runSelectionFlow()
	}

	return recommended, nil
}

// roleBasedSelection allows selection by development role
func (is *InteractiveSelector) roleBasedSelection() ([]SelectableItem, error) {
	var selectedRoles []string

	// Define role categories
	roleOptions := []huh.Option[string]{
		huh.NewOption("Frontend Developer - UI/UX, React, Vue, styling", "frontend"),
		huh.NewOption("Backend Developer - APIs, databases, microservices", "backend"),
		huh.NewOption("Full Stack Developer - Frontend + Backend + DevOps", "fullstack"),
		huh.NewOption("DevOps Engineer - Infrastructure, CI/CD, monitoring", "devops"),
		huh.NewOption("Mobile Developer - iOS, Android, React Native", "mobile"),
		huh.NewOption("Data/ML Engineer - Machine learning, data pipelines", "ml"),
		huh.NewOption("Product Manager - Requirements, analysis, documentation", "product"),
		huh.NewOption("QA Engineer - Testing, quality assurance", "qa"),
		huh.NewOption("Security Engineer - Security assessment, compliance", "security"),
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("👤 Select your development roles").
				Description("Choose all roles that apply to you - The (Agentic) Startup will customize your toolkit").
				Options(roleOptions...).
				Value(&selectedRoles),
		),
	).WithTheme(huh.ThemeCharm())

	err := form.Run()
	if err != nil {
		return nil, err
	}

	if len(selectedRoles) == 0 {
		return is.runSelectionFlow() // Go back if no roles selected
	}

	// Get items for selected roles
	selectedItems := is.getItemsForRoles(selectedRoles)

	// Confirm selection
	var confirm bool
	items := is.formatItemsForDisplay(selectedItems)

	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("👤 Role-based Selection").
				Description(fmt.Sprintf("Based on your roles, The (Agentic) Startup will install %d components:\n\n%s",
					len(selectedItems), items)),
			huh.NewConfirm().
				Title("🚀 Install these components?").
				Value(&confirm),
		),
	).WithTheme(huh.ThemeCharm())

	err = confirmForm.Run()
	if err != nil {
		return nil, err
	}

	if !confirm {
		return is.runSelectionFlow()
	}

	return selectedItems, nil
}

// customSelection provides granular item-by-item selection
func (is *InteractiveSelector) customSelection() ([]SelectableItem, error) {
	// Group items by category for easier selection
	categories := is.groupItemsByCategory()

	var selectedItems []SelectableItem

	// Select by category
	for categoryName, items := range categories {
		if len(items) == 0 {
			continue
		}

		// For required categories, just add all items
		if categoryName == "Configuration" {
			selectedItems = append(selectedItems, items...)
			continue
		}

		// Create options for this category
		var options []huh.Option[string]
		for _, item := range items {
			desc := item.Description
			if len(desc) > 80 {
				desc = desc[:77] + "..."
			}
			displayText := fmt.Sprintf("%s - %s", item.Name, desc)
			options = append(options, huh.NewOption(displayText, item.ID))
		}

		var selectedIDs []string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title(fmt.Sprintf("Select from %s", categoryName)).
					Description(fmt.Sprintf("Choose which items to install (%d available)", len(items))).
					Options(options...).
					Value(&selectedIDs),
			),
		).WithTheme(huh.ThemeCharm())

		err := form.Run()
		if err != nil {
			return nil, err
		}

		// Add selected items
		for _, item := range items {
			for _, selectedID := range selectedIDs {
				if item.ID == selectedID {
					item.Selected = true
					selectedItems = append(selectedItems, item)
					break
				}
			}
		}
	}

	if len(selectedItems) == 0 {
		return is.runSelectionFlow() // Go back if nothing selected
	}

	// Final confirmation
	var confirm bool
	itemsList := is.formatItemsForDisplay(selectedItems)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("✨ Custom Selection").
				Description(fmt.Sprintf("You've selected %d components:\n\n%s",
					len(selectedItems), itemsList)),
			huh.NewConfirm().
				Title("Install selected components?").
				Value(&confirm),
		),
	).WithTheme(huh.ThemeCharm())

	err := form.Run()
	if err != nil {
		return nil, err
	}

	if !confirm {
		return is.runSelectionFlow()
	}

	return selectedItems, nil
}

// minimalSelection provides just the essentials
func (is *InteractiveSelector) minimalSelection() ([]SelectableItem, error) {
	minimal := is.getMinimalItems()

	var confirm bool
	items := is.formatItemsForDisplay(minimal)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("📦 Minimal Installation").
				Description(fmt.Sprintf("This will install only essential components (%d items):\n\n%s",
					len(minimal), items)),
			huh.NewConfirm().
				Title("Install minimal setup?").
				Value(&confirm),
		),
	).WithTheme(huh.ThemeCharm())

	err := form.Run()
	if err != nil {
		return nil, err
	}

	if !confirm {
		return is.runSelectionFlow()
	}

	return minimal, nil
}

// Helper methods for different selection strategies

func (is *InteractiveSelector) getRecommendedItems() []SelectableItem {
	var recommended []SelectableItem

	for _, item := range is.items {
		// Always include settings and core commands
		if item.Required || item.Type == ItemTypeCommand {
			recommended = append(recommended, item)
			continue
		}

		// Include most commonly used agents
		if item.Type == ItemTypeAgent {
			commonAgents := []string{
				"the-chief", "the-software-engineer", "the-analyst",
				"the-designer", "the-platform-engineer", "the-qa-engineer",
			}

			for _, common := range commonAgents {
				if strings.Contains(item.ID, common) {
					recommended = append(recommended, item)
					break
				}
			}
		}

		// Include output styles
		if item.Type == ItemTypeOutputStyle {
			recommended = append(recommended, item)
		}
	}

	return recommended
}

func (is *InteractiveSelector) getItemsForRoles(roles []string) []SelectableItem {
	var selected []SelectableItem
	roleAgentMap := map[string][]string{
		"frontend":  {"the-designer", "the-software-engineer"},
		"backend":   {"the-software-engineer", "the-platform-engineer", "the-architect"},
		"fullstack": {"the-software-engineer", "the-platform-engineer", "the-designer", "the-analyst"},
		"devops":    {"the-platform-engineer", "the-architect", "the-security-engineer"},
		"mobile":    {"the-mobile-engineer", "the-software-engineer", "the-designer"},
		"ml":        {"the-ml-engineer", "the-platform-engineer", "the-analyst"},
		"product":   {"the-analyst", "the-designer"},
		"qa":        {"the-qa-engineer", "the-software-engineer"},
		"security":  {"the-security-engineer", "the-platform-engineer"},
	}

	// Collect all relevant agents for selected roles
	agentSet := make(map[string]bool)
	for _, role := range roles {
		if agents, exists := roleAgentMap[role]; exists {
			for _, agent := range agents {
				agentSet[agent] = true
			}
		}
	}

	// Add the chief for any role (orchestration)
	agentSet["the-chief"] = true

	for _, item := range is.items {
		// Always include required items and commands
		if item.Required || item.Type == ItemTypeCommand || item.Type == ItemTypeOutputStyle {
			selected = append(selected, item)
			continue
		}

		// Include agents that match selected roles
		if item.Type == ItemTypeAgent {
			for agent := range agentSet {
				if strings.Contains(item.ID, agent) {
					selected = append(selected, item)
					break
				}
			}
		}
	}

	return selected
}

func (is *InteractiveSelector) getMinimalItems() []SelectableItem {
	var minimal []SelectableItem

	for _, item := range is.items {
		// Only required items and core commands
		if item.Required {
			minimal = append(minimal, item)
		} else if item.Type == ItemTypeCommand && strings.Contains(item.ID, "specify") {
			// Just the specify command as essential
			minimal = append(minimal, item)
		} else if item.ID == "agents/the-chief.md" {
			// Always include the chief for orchestration
			minimal = append(minimal, item)
		}
	}

	return minimal
}

func (is *InteractiveSelector) groupItemsByCategory() map[string][]SelectableItem {
	categories := make(map[string][]SelectableItem)

	for _, item := range is.items {
		category := item.Category
		categories[category] = append(categories[category], item)
	}

	// Sort items within each category
	for category := range categories {
		sort.Slice(categories[category], func(i, j int) bool {
			return categories[category][i].Name < categories[category][j].Name
		})
	}

	return categories
}

func (is *InteractiveSelector) formatItemsForDisplay(items []SelectableItem) string {
	categories := make(map[string][]SelectableItem)

	// Group by category
	for _, item := range items {
		categories[item.Category] = append(categories[item.Category], item)
	}

	var result []string

	// Sort categories for consistent display
	var categoryNames []string
	for category := range categories {
		categoryNames = append(categoryNames, category)
	}
	sort.Strings(categoryNames)

	for _, category := range categoryNames {
		items := categories[category]
		if len(items) == 0 {
			continue
		}

		// Add category header
		result = append(result, fmt.Sprintf("• %s (%d items)", category, len(items)))

		// Add up to 3 items, then "... and X more"
		for i, item := range items {
			if i < 3 {
				result = append(result, fmt.Sprintf("  - %s", item.Name))
			} else {
				remaining := len(items) - 3
				result = append(result, fmt.Sprintf("  - ... and %d more", remaining))
				break
			}
		}
	}

	return strings.Join(result, "\n")
}