package config

import (
	"fmt"
	"strings"
	"unicode/utf8"
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

// Validate validates the MarkdownConfig and returns validation errors with fallback support.
func (mc *MarkdownConfig) Validate() error {
	if mc == nil {
		return nil
	}

	var warnings []string
	hasErrors := false

	// Validate database config
	if mc.Database != nil {
		if issues := validateObjectConfigWithFallback(mc.Database, "database"); len(issues) > 0 {
			warnings = append(warnings, issues...)
		}
	}

	// Validate schemas config
	if mc.Schemas != nil {
		if issues := validateObjectConfigWithFallback(mc.Schemas, "schemas"); len(issues) > 0 {
			warnings = append(warnings, issues...)
		}
	}

	// Validate tables config
	if mc.Tables != nil {
		if mc.Tables.ObjectCustomConfig != nil {
			if issues := validateObjectConfigWithFallback(mc.Tables.ObjectCustomConfig, "tables"); len(issues) > 0 {
				warnings = append(warnings, issues...)
			}
		}
		for tableName, config := range mc.Tables.Specific {
			if config != nil {
				if issues := validateObjectConfigWithFallback(config, fmt.Sprintf("tables.specific.%s", tableName)); len(issues) > 0 {
					warnings = append(warnings, issues...)
				}
			}
		}
	}

	// Validate columns config
	if mc.Columns != nil {
		if mc.Columns.ObjectCustomConfig != nil {
			if issues := validateObjectConfigWithFallback(mc.Columns.ObjectCustomConfig, "columns"); len(issues) > 0 {
				warnings = append(warnings, issues...)
			}
		}
		for tableName, config := range mc.Columns.Specific {
			if config != nil {
				if issues := validateObjectConfigWithFallback(config, fmt.Sprintf("columns.specific.%s", tableName)); len(issues) > 0 {
					warnings = append(warnings, issues...)
				}
			}
		}
	}

	// Validate views config
	if mc.Views != nil {
		if issues := validateObjectConfigWithFallback(mc.Views, "views"); len(issues) > 0 {
			warnings = append(warnings, issues...)
		}
	}

	// Validate indexes config
	if mc.Indexes != nil {
		if issues := validateObjectConfigWithFallback(mc.Indexes, "indexes"); len(issues) > 0 {
			warnings = append(warnings, issues...)
		}
	}

	// Validate constraints config
	if mc.Constraints != nil {
		if issues := validateObjectConfigWithFallback(mc.Constraints, "constraints"); len(issues) > 0 {
			warnings = append(warnings, issues...)
		}
	}

	// Validate functions config
	if mc.Functions != nil {
		if issues := validateObjectConfigWithFallback(mc.Functions, "functions"); len(issues) > 0 {
			warnings = append(warnings, issues...)
		}
	}

	// Validate others config
	for objectType, config := range mc.Others {
		if config != nil {
			if issues := validateObjectConfigWithFallback(config, fmt.Sprintf("others.%s", objectType)); len(issues) > 0 {
				warnings = append(warnings, issues...)
			}
		}
	}

	// Print warnings
	for _, warning := range warnings {
		fmt.Printf("警告: %s\n", warning)
	}

	if hasErrors {
		return fmt.Errorf("markdown設定に重大なエラーがあります。デフォルト値で続行します")
	}

	return nil
}

// validateObjectConfigWithFallback validates a single ObjectCustomConfig with fallback support.
func validateObjectConfigWithFallback(config *ObjectCustomConfig, context string) []string {
	if config == nil {
		return nil
	}

	var warnings []string

	// Validate order fields - check for duplicates and fix them
	seen := make(map[string]bool)
	cleanOrder := make([]string, 0, len(config.Order))

	for _, field := range config.Order {
		if field == "" {
			warnings = append(warnings, fmt.Sprintf("%s: 空のフィールド名が順序設定に含まれています（スキップします）", context))
			continue
		}

		// Check UTF-8 validity
		if !utf8.ValidString(field) {
			warnings = append(warnings, fmt.Sprintf("%s: 無効なUTF-8文字を含むフィールド '%s' をスキップします", context, field))
			continue
		}

		if seen[field] {
			warnings = append(warnings, fmt.Sprintf("%s: 重複するフィールド '%s' が順序設定にあります（最初のもののみ使用）", context, field))
			continue
		}

		seen[field] = true
		cleanOrder = append(cleanOrder, field)
	}

	// Update order with cleaned version
	config.Order = cleanOrder

	// Validate aliases - check for empty keys or values and clean up
	cleanAliases := make(map[string]string)
	for key, value := range config.Aliases {
		if key == "" {
			warnings = append(warnings, fmt.Sprintf("%s: 空のキーがエイリアス設定にあります（スキップします）", context))
			continue
		}
		if value == "" {
			warnings = append(warnings, fmt.Sprintf("%s: エイリアスキー '%s' の値が空です（スキップします）", context, key))
			continue
		}

		// Check UTF-8 validity
		if !utf8.ValidString(key) || !utf8.ValidString(value) {
			warnings = append(warnings, fmt.Sprintf("%s: 無効なUTF-8文字を含むエイリアス '%s' -> '%s' をスキップします", context, key, value))
			continue
		}

		cleanAliases[key] = value
	}

	// Update aliases with cleaned version
	config.Aliases = cleanAliases

	return warnings
}

// validateObjectConfig validates a single ObjectCustomConfig (legacy method for compatibility).
func validateObjectConfig(config *ObjectCustomConfig, context string) error {
	warnings := validateObjectConfigWithFallback(config, context)
	if len(warnings) > 0 {
		return fmt.Errorf(strings.Join(warnings, "; "))
	}
	return nil
}

// IsValid checks if the MarkdownConfig is valid without modifying it.
func (mc *MarkdownConfig) IsValid() bool {
	if mc == nil {
		return true
	}

	// Check each configuration
	configs := []*ObjectCustomConfig{
		mc.Database, mc.Schemas, mc.Views, mc.Indexes, mc.Constraints, mc.Functions,
	}

	if mc.Tables != nil {
		configs = append(configs, mc.Tables.ObjectCustomConfig)
		for _, config := range mc.Tables.Specific {
			configs = append(configs, config)
		}
	}

	if mc.Columns != nil {
		configs = append(configs, mc.Columns.ObjectCustomConfig)
		for _, config := range mc.Columns.Specific {
			configs = append(configs, config)
		}
	}

	for _, config := range mc.Others {
		configs = append(configs, config)
	}

	for _, config := range configs {
		if !isObjectConfigValid(config) {
			return false
		}
	}

	return true
}

// isObjectConfigValid checks if a single ObjectCustomConfig is valid.
func isObjectConfigValid(config *ObjectCustomConfig) bool {
	if config == nil {
		return true
	}

	// Check for duplicate fields in order
	seen := make(map[string]bool)
	for _, field := range config.Order {
		if field == "" || !utf8.ValidString(field) {
			return false
		}
		if seen[field] {
			return false
		}
		seen[field] = true
	}

	// Check aliases
	for key, value := range config.Aliases {
		if key == "" || value == "" || !utf8.ValidString(key) || !utf8.ValidString(value) {
			return false
		}
	}

	return true
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