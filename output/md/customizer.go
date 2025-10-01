package md

import (
	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/schema"
)

// MarkdownCustomizer is the interface for customizing Markdown output.
type MarkdownCustomizer interface {
	// CustomizeTableOutput customizes table output data with configuration.
	CustomizeTableOutput(table *schema.Table, config *config.ObjectCustomConfig) *CustomizedTableData

	// CustomizeColumnOutput customizes column output data with configuration.
	CustomizeColumnOutput(columns []*schema.Column, config *config.ObjectCustomConfig) []*CustomizedColumnData

	// ApplyAliases applies field name aliases to a field name.
	ApplyAliases(fieldName string, aliases map[string]string) string

	// GetDisplayOrder returns the display order based on default and configuration order.
	GetDisplayOrder(defaultOrder []string, configOrder []string) []string
}

// CustomizedTableData represents the customized table data for template rendering.
type CustomizedTableData struct {
	Table   *schema.Table
	Columns []*CustomizedColumnData
	Order   []string
	Aliases map[string]string
}

// CustomizedColumnData represents the customized column data for template rendering.
type CustomizedColumnData struct {
	Column      *schema.Column
	DisplayName string
	LogicalName string
	Show        bool
}

// DefaultMarkdownCustomizer is the default implementation of MarkdownCustomizer.
type DefaultMarkdownCustomizer struct{}

// NewDefaultMarkdownCustomizer creates a new DefaultMarkdownCustomizer.
func NewDefaultMarkdownCustomizer() *DefaultMarkdownCustomizer {
	return &DefaultMarkdownCustomizer{}
}

// CustomizeTableOutput customizes table output data with configuration.
func (c *DefaultMarkdownCustomizer) CustomizeTableOutput(table *schema.Table, config *config.ObjectCustomConfig) *CustomizedTableData {
	if table == nil {
		return nil
	}

	// Create customized column data
	customizedColumns := c.CustomizeColumnOutput(table.Columns, config)

	// Get display order for table fields
	defaultOrder := []string{"name", "LogicalName", "comment", "type"}
	var displayOrder []string
	if config != nil && len(config.Order) > 0 {
		displayOrder = c.GetDisplayOrder(defaultOrder, config.Order)
	} else {
		displayOrder = defaultOrder
	}

	// Get aliases
	var aliases map[string]string
	if config != nil {
		aliases = config.Aliases
	}
	if aliases == nil {
		aliases = make(map[string]string)
	}

	return &CustomizedTableData{
		Table:   table,
		Columns: customizedColumns,
		Order:   displayOrder,
		Aliases: aliases,
	}
}

// CustomizeColumnOutput customizes column output data with configuration.
func (c *DefaultMarkdownCustomizer) CustomizeColumnOutput(columns []*schema.Column, config *config.ObjectCustomConfig) []*CustomizedColumnData {
	if columns == nil {
		return nil
	}

	result := make([]*CustomizedColumnData, len(columns))

	for i, column := range columns {
		if column == nil {
			continue
		}

		// Determine display name
		displayName := column.Name
		logicalName := column.GetLogicalNameOrFallback()

		// Apply aliases if available
		if config != nil && config.Aliases != nil {
			if alias, exists := config.Aliases[column.Name]; exists && alias != "" {
				displayName = alias
			}
		}

		// All columns are shown by default
		show := true

		result[i] = &CustomizedColumnData{
			Column:      column,
			DisplayName: displayName,
			LogicalName: logicalName,
			Show:        show,
		}
	}

	return result
}

// ApplyAliases applies field name aliases to a field name.
func (c *DefaultMarkdownCustomizer) ApplyAliases(fieldName string, aliases map[string]string) string {
	if aliases == nil {
		return fieldName
	}

	if alias, exists := aliases[fieldName]; exists && alias != "" {
		return alias
	}

	return fieldName
}

// GetDisplayOrder returns the display order based on default and configuration order.
// If configOrder is provided and not empty, it returns only the specified fields in that order.
// If configOrder is empty, it returns the defaultOrder.
func (c *DefaultMarkdownCustomizer) GetDisplayOrder(defaultOrder []string, configOrder []string) []string {
	if len(configOrder) == 0 {
		return defaultOrder
	}
	// When configOrder is specified, return only those fields
	return configOrder
}
