package md

import (
	"strings"

	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/schema"
)

// MdExtended extends the basic Md struct with customization capabilities
type MdExtended struct {
	*Md
	customizer MarkdownCustomizer
}

// NewExtended creates a new MdExtended instance with customization support
func NewExtended(c *config.Config) *MdExtended {
	return &MdExtended{
		Md:         New(c),
		customizer: NewDefaultMarkdownCustomizer(),
	}
}

// makeTableTemplateDataWithCustomization extends the basic makeTableTemplateData
// with customization support for logical names, column ordering, and aliases
func (m *MdExtended) makeTableTemplateDataWithCustomization(t *schema.Table) map[string]interface{} {
	// Start with the basic template data
	basicData := m.Md.makeTableTemplateData(t)

	// Get markdown customization configuration
	var markdownConfig *config.MarkdownConfig
	if m.config.Markdown != nil {
		markdownConfig = m.config.Markdown
	}

	// Apply table-specific customization if available
	var tableCustomConfig *config.ObjectCustomConfig
	if markdownConfig != nil && markdownConfig.Tables != nil {
		// Use general table configuration
		tableCustomConfig = markdownConfig.Tables.ObjectCustomConfig

		// Override with specific table configuration if exists
		if markdownConfig.Tables.Specific != nil {
			if specificConfig, exists := markdownConfig.Tables.Specific[t.Name]; exists {
				tableCustomConfig = specificConfig
			}
		}
	}

	// If no customization is configured, return basic data
	if tableCustomConfig == nil {
		return basicData
	}

	// Apply customization to columns data
	customizedData := m.customizeColumnsData(t, basicData, tableCustomConfig)

	return customizedData
}

// customizeColumnsData applies customization settings to the columns data
func (m *MdExtended) customizeColumnsData(t *schema.Table, basicData map[string]interface{}, config *config.ObjectCustomConfig) map[string]interface{} {
	// Extract basic columns data
	columnsData, ok := basicData["Columns"].([][]string)
	if !ok || len(columnsData) < 2 {
		return basicData
	}

	// Parse current structure
	headerRow := columnsData[0]
	dataRows := columnsData[2:]

	// Build customized columns structure
	customizedData := m.buildCustomizedColumns(t, headerRow, dataRows, config)

	// Replace columns data in the template data
	result := make(map[string]interface{})
	for k, v := range basicData {
		result[k] = v
	}
	result["Columns"] = customizedData

	return result
}

// buildCustomizedColumns builds the customized columns table data
func (m *MdExtended) buildCustomizedColumns(t *schema.Table, originalHeaders []string, originalDataRows [][]string, config *config.ObjectCustomConfig) [][]string {
	// Determine the final column order and structure
	finalStructure := m.determineFinalColumnStructure(originalHeaders, config)

	// Build new header row with aliases applied
	newHeaderRow := make([]string, len(finalStructure))
	newSeparatorRow := make([]string, len(finalStructure))

	for i, colInfo := range finalStructure {
		headerText := colInfo.DisplayName
		if config.Aliases != nil {
			if alias, exists := config.Aliases[colInfo.OriginalName]; exists && alias != "" {
				headerText = alias
			}
		}
		newHeaderRow[i] = headerText
		newSeparatorRow[i] = strings.Repeat("-", len(headerText))
	}

	// Build new data rows
	newDataRows := make([][]string, len(originalDataRows))
	for rowIdx, originalRow := range originalDataRows {
		newRow := make([]string, len(finalStructure))
		for colIdx, colInfo := range finalStructure {
			if colInfo.SourceIndex >= 0 && colInfo.SourceIndex < len(originalRow) {
				newRow[colIdx] = originalRow[colInfo.SourceIndex]
			} else if colInfo.IsLogicalName && rowIdx < len(t.Columns) {
				// Add logical name data
				newRow[colIdx] = t.Columns[rowIdx].GetLogicalNameOrFallback()
			} else {
				newRow[colIdx] = ""
			}
		}
		newDataRows[rowIdx] = newRow
	}

	// Combine all rows
	result := [][]string{newHeaderRow, newSeparatorRow}
	result = append(result, newDataRows...)

	return result
}

// ColumnInfo represents information about a column in the final structure
type ColumnInfo struct {
	OriginalName string // Original column name (e.g., "Name", "Type")
	DisplayName  string // Display name for the column
	SourceIndex  int    // Index in the original data (-1 for new columns)
	IsLogicalName bool  // Whether this is a logical name column
}

// determineFinalColumnStructure determines the final column structure based on configuration
func (m *MdExtended) determineFinalColumnStructure(originalHeaders []string, config *config.ObjectCustomConfig) []ColumnInfo {
	// Create mapping of original headers to their indices
	headerIndexMap := make(map[string]int)
	for i, header := range originalHeaders {
		headerIndexMap[header] = i
	}

	// Start with the original headers as default columns, preserving existing order
	defaultColumns := make([]string, len(originalHeaders))
	copy(defaultColumns, originalHeaders)

	// Add logical name column if enabled
	if config.ShowLogicalName {
		// Insert after Name column
		var newDefault []string
		for _, header := range defaultColumns {
			newDefault = append(newDefault, header)
			if header == "Name" {
				newDefault = append(newDefault, "Logical Name")
			}
		}
		// If no "Name" column found, just append at the beginning
		if len(newDefault) == len(defaultColumns) {
			newDefault = append([]string{"Logical Name"}, defaultColumns...)
		}
		defaultColumns = newDefault
	}

	// Apply custom order if specified
	var finalOrder []string
	if len(config.Order) > 0 {
		finalOrder = m.customizer.GetDisplayOrder(defaultColumns, config.Order)
	} else {
		finalOrder = defaultColumns
	}

	// Build final structure
	result := make([]ColumnInfo, 0, len(finalOrder))
	for _, columnName := range finalOrder {
		colInfo := ColumnInfo{
			OriginalName: columnName,
			DisplayName:  columnName,
			SourceIndex:  -1,
			IsLogicalName: false,
		}

		if columnName == "Logical Name" {
			colInfo.IsLogicalName = true
		} else if idx, exists := headerIndexMap[columnName]; exists {
			colInfo.SourceIndex = idx
		}

		result = append(result, colInfo)
	}

	return result
}

// ProcessTemplateDataWithCustomization is a public method to process template data with customization
func (m *MdExtended) ProcessTemplateDataWithCustomization(t *schema.Table) map[string]interface{} {
	return m.makeTableTemplateDataWithCustomization(t)
}