package config

import (
	"fmt"
	"strings"
)

// MarkdownConfig is the configuration for Markdown output customization.
type MarkdownConfig struct {
	Database    *ObjectCustomConfig            `yaml:"database,omitempty"`
	Schemas     *ObjectCustomConfig            `yaml:"schemas,omitempty"`
	Tables      *TableCustomConfig             `yaml:"tables,omitempty"`
	Columns     *ColumnCustomConfig            `yaml:"columns,omitempty"`
	Views       *ObjectCustomConfig            `yaml:"views,omitempty"`
	Indexes     *ObjectCustomConfig            `yaml:"indexes,omitempty"`
	Constraints *ObjectCustomConfig            `yaml:"constraints,omitempty"`
	Functions   *ObjectCustomConfig            `yaml:"functions,omitempty"`
	Others      map[string]*ObjectCustomConfig `yaml:"others,omitempty"`
}

// ObjectCustomConfig is the base configuration for object customization.
type ObjectCustomConfig struct {
	ShowLogicalName bool              `yaml:"show_logical_name,omitempty"`
	Order           []string          `yaml:"order,omitempty"`
	Aliases         map[string]string `yaml:"aliases,omitempty"`
}

// TableCustomConfig is the configuration for table customization with specific table settings.
type TableCustomConfig struct {
	*ObjectCustomConfig
	Specific map[string]*ObjectCustomConfig `yaml:"specific,omitempty"`
}

// ColumnCustomConfig is the configuration for column customization with specific table column settings.
type ColumnCustomConfig struct {
	*ObjectCustomConfig
	Specific map[string]*ObjectCustomConfig `yaml:"specific,omitempty"`
}

// GetObjectConfig returns the appropriate object configuration.
func (mc *MarkdownConfig) GetObjectConfig(objectType string) *ObjectCustomConfig {
	if mc == nil {
		return nil
	}

	switch strings.ToLower(objectType) {
	case "database":
		return mc.Database
	case "schema", "schemas":
		return mc.Schemas
	case "table", "tables":
		if mc.Tables != nil {
			return mc.Tables.ObjectCustomConfig
		}
		return nil
	case "column", "columns":
		if mc.Columns != nil {
			return mc.Columns.ObjectCustomConfig
		}
		return nil
	case "view", "views":
		return mc.Views
	case "index", "indexes":
		return mc.Indexes
	case "constraint", "constraints":
		return mc.Constraints
	case "function", "functions":
		return mc.Functions
	default:
		if mc.Others != nil {
			return mc.Others[objectType]
		}
		return nil
	}
}

// GetTableSpecificConfig returns the specific configuration for a table.
func (mc *MarkdownConfig) GetTableSpecificConfig(tableName string) *ObjectCustomConfig {
	if mc == nil || mc.Tables == nil || mc.Tables.Specific == nil {
		return nil
	}
	return mc.Tables.Specific[tableName]
}

// GetColumnSpecificConfig returns the specific configuration for a table's columns.
func (mc *MarkdownConfig) GetColumnSpecificConfig(tableName string) *ObjectCustomConfig {
	if mc == nil || mc.Columns == nil || mc.Columns.Specific == nil {
		return nil
	}
	return mc.Columns.Specific[tableName]
}

// GetEffectiveConfig returns the effective configuration by merging global, object-type, and specific configurations.
// Priority: specific > object-type > global
func (mc *MarkdownConfig) GetEffectiveConfig(objectType, specificName string) *ObjectCustomConfig {
	var globalConfig, specificConfig *ObjectCustomConfig

	// Get object-type configuration
	globalConfig = mc.GetObjectConfig(objectType)

	// Get specific configuration
	switch strings.ToLower(objectType) {
	case "table", "tables":
		specificConfig = mc.GetTableSpecificConfig(specificName)
	case "column", "columns":
		specificConfig = mc.GetColumnSpecificConfig(specificName)
	}

	// Merge configurations
	return mergeObjectConfigs(globalConfig, specificConfig)
}

// mergeObjectConfigs merges two ObjectCustomConfig instances.
// Values from specific config take precedence over global config.
func mergeObjectConfigs(global, specific *ObjectCustomConfig) *ObjectCustomConfig {
	if global == nil && specific == nil {
		return nil
	}
	if global == nil {
		return specific
	}
	if specific == nil {
		return global
	}

	merged := &ObjectCustomConfig{
		ShowLogicalName: global.ShowLogicalName,
		Order:           make([]string, len(global.Order)),
		Aliases:         make(map[string]string),
	}

	// Copy global order and aliases
	copy(merged.Order, global.Order)
	for k, v := range global.Aliases {
		merged.Aliases[k] = v
	}

	// Override with specific values
	if specific.ShowLogicalName != global.ShowLogicalName {
		merged.ShowLogicalName = specific.ShowLogicalName
	}
	if len(specific.Order) > 0 {
		merged.Order = make([]string, len(specific.Order))
		copy(merged.Order, specific.Order)
	}
	for k, v := range specific.Aliases {
		merged.Aliases[k] = v
	}

	return merged
}

// Validate validates the MarkdownConfig and returns validation errors.
func (mc *MarkdownConfig) Validate() error {
	if mc == nil {
		return nil
	}

	var errors []string

	// Validate database config
	if err := validateObjectConfig(mc.Database, "database"); err != nil {
		errors = append(errors, err.Error())
	}

	// Validate schemas config
	if err := validateObjectConfig(mc.Schemas, "schemas"); err != nil {
		errors = append(errors, err.Error())
	}

	// Validate tables config
	if mc.Tables != nil {
		if err := validateObjectConfig(mc.Tables.ObjectCustomConfig, "tables"); err != nil {
			errors = append(errors, err.Error())
		}
		for tableName, config := range mc.Tables.Specific {
			if err := validateObjectConfig(config, fmt.Sprintf("tables.specific.%s", tableName)); err != nil {
				errors = append(errors, err.Error())
			}
		}
	}

	// Validate columns config
	if mc.Columns != nil {
		if err := validateObjectConfig(mc.Columns.ObjectCustomConfig, "columns"); err != nil {
			errors = append(errors, err.Error())
		}
		for tableName, config := range mc.Columns.Specific {
			if err := validateObjectConfig(config, fmt.Sprintf("columns.specific.%s", tableName)); err != nil {
				errors = append(errors, err.Error())
			}
		}
	}

	// Validate views config
	if err := validateObjectConfig(mc.Views, "views"); err != nil {
		errors = append(errors, err.Error())
	}

	// Validate indexes config
	if err := validateObjectConfig(mc.Indexes, "indexes"); err != nil {
		errors = append(errors, err.Error())
	}

	// Validate constraints config
	if err := validateObjectConfig(mc.Constraints, "constraints"); err != nil {
		errors = append(errors, err.Error())
	}

	// Validate functions config
	if err := validateObjectConfig(mc.Functions, "functions"); err != nil {
		errors = append(errors, err.Error())
	}

	// Validate others config
	for objectType, config := range mc.Others {
		if err := validateObjectConfig(config, fmt.Sprintf("others.%s", objectType)); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("markdown config validation errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// validateObjectConfig validates a single ObjectCustomConfig.
func validateObjectConfig(config *ObjectCustomConfig, context string) error {
	if config == nil {
		return nil
	}

	var errors []string

	// Validate order fields - check for duplicates
	seen := make(map[string]bool)
	for _, field := range config.Order {
		if field == "" {
			errors = append(errors, fmt.Sprintf("%s: empty field name in order", context))
			continue
		}
		if seen[field] {
			errors = append(errors, fmt.Sprintf("%s: duplicate field '%s' in order", context, field))
		}
		seen[field] = true
	}

	// Validate aliases - check for empty keys or values
	for key, value := range config.Aliases {
		if key == "" {
			errors = append(errors, fmt.Sprintf("%s: empty key in aliases", context))
		}
		if value == "" {
			errors = append(errors, fmt.Sprintf("%s: empty value for alias key '%s'", context, key))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}

	return nil
}

// SetDefaults sets default values for MarkdownConfig.
func (mc *MarkdownConfig) SetDefaults() {
	if mc == nil {
		return
	}

	// Set default values for each object type if not already set
	if mc.Database == nil {
		mc.Database = &ObjectCustomConfig{}
	}
	setObjectDefaults(mc.Database, "database")

	if mc.Schemas == nil {
		mc.Schemas = &ObjectCustomConfig{}
	}
	setObjectDefaults(mc.Schemas, "schemas")

	if mc.Tables == nil {
		mc.Tables = &TableCustomConfig{ObjectCustomConfig: &ObjectCustomConfig{}}
	}
	if mc.Tables.ObjectCustomConfig == nil {
		mc.Tables.ObjectCustomConfig = &ObjectCustomConfig{}
	}
	setObjectDefaults(mc.Tables.ObjectCustomConfig, "tables")

	if mc.Columns == nil {
		mc.Columns = &ColumnCustomConfig{ObjectCustomConfig: &ObjectCustomConfig{}}
	}
	if mc.Columns.ObjectCustomConfig == nil {
		mc.Columns.ObjectCustomConfig = &ObjectCustomConfig{}
	}
	setObjectDefaults(mc.Columns.ObjectCustomConfig, "columns")

	if mc.Views == nil {
		mc.Views = &ObjectCustomConfig{}
	}
	setObjectDefaults(mc.Views, "views")

	if mc.Indexes == nil {
		mc.Indexes = &ObjectCustomConfig{}
	}
	setObjectDefaults(mc.Indexes, "indexes")

	if mc.Constraints == nil {
		mc.Constraints = &ObjectCustomConfig{}
	}
	setObjectDefaults(mc.Constraints, "constraints")

	if mc.Functions == nil {
		mc.Functions = &ObjectCustomConfig{}
	}
	setObjectDefaults(mc.Functions, "functions")
}

// setObjectDefaults sets default values for a single ObjectCustomConfig.
func setObjectDefaults(config *ObjectCustomConfig, objectType string) {
	if config == nil {
		return
	}

	// Set default show_logical_name to false
	// Note: ShowLogicalName is false by default in Go

	// Set default order based on object type
	if len(config.Order) == 0 {
		switch objectType {
		case "database":
			config.Order = []string{"name", "logical_name", "comment"}
		case "schemas":
			config.Order = []string{"name", "logical_name", "comment"}
		case "tables":
			config.Order = []string{"name", "logical_name", "comment", "type"}
		case "columns":
			config.Order = []string{"name", "logical_name", "type", "nullable", "default", "comment"}
		case "views":
			config.Order = []string{"name", "logical_name", "comment", "definition"}
		case "indexes":
			config.Order = []string{"name", "logical_name", "columns", "comment"}
		case "constraints":
			config.Order = []string{"name", "logical_name", "type", "columns", "comment"}
		case "functions":
			config.Order = []string{"name", "logical_name", "return_type", "arguments", "comment"}
		default:
			config.Order = []string{"name", "logical_name", "comment"}
		}
	}

	// Initialize aliases map if nil
	if config.Aliases == nil {
		config.Aliases = make(map[string]string)
	}
}
