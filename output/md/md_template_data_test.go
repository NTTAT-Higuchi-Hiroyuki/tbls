package md

import (
	"testing"

	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/schema"
)

func TestMdExtended_makeTableTemplateDataWithCustomization(t *testing.T) {
	tests := []struct {
		name        string
		table       *schema.Table
		config      *config.Config
		wantColumns int
		wantHeaders []string
		description string
	}{
		{
			name: "default configuration without customization",
			table: &schema.Table{
				Name: "users",
				Columns: []*schema.Column{
					{
						Name:        "id",
						Type:        "int",
						LogicalName: "ユーザーID",
						Comment:     "システム内でユーザーを一意に識別するID",
					},
					{
						Name:        "name",
						Type:        "varchar(100)",
						LogicalName: "ユーザー名",
						Comment:     "ユーザーの表示名",
					},
				},
			},
			config: &config.Config{
				Format: config.Format{
					Number: false,
					Adjust: false,
				},
			},
			wantColumns: 7, // Name, Type, Default, Nullable, Children, Parents, Comment
			description: "Should work without customization config",
		},
		{
			name: "with logical name display enabled",
			table: &schema.Table{
				Name: "users",
				Columns: []*schema.Column{
					{
						Name:        "id",
						Type:        "int",
						LogicalName: "ユーザーID",
						Comment:     "システム内でユーザーを一意に識別するID",
					},
					{
						Name:        "name",
						Type:        "varchar(100)",
						LogicalName: "ユーザー名",
						Comment:     "ユーザーの表示名",
					},
				},
			},
			config: &config.Config{
				Format: config.Format{
					Number: false,
					Adjust: false,
				},
				Markdown: &config.MarkdownConfig{
					Tables: &config.TableCustomConfig{
						ObjectCustomConfig: &config.ObjectCustomConfig{
							ShowLogicalName: true,
						},
					},
				},
			},
			wantColumns: 8, // Name, Logical Name, Type, Default, Nullable, Children, Parents, Comment
			wantHeaders: []string{"Name", "Logical Name", "Type", "Default", "Nullable", "Children", "Parents", "Comment"},
			description: "Should add logical name column when enabled",
		},
		{
			name: "with custom column order",
			table: &schema.Table{
				Name: "users",
				Columns: []*schema.Column{
					{
						Name:        "id",
						Type:        "int",
						LogicalName: "ユーザーID",
						Comment:     "システム内でユーザーを一意に識別するID",
					},
				},
			},
			config: &config.Config{
				Format: config.Format{
					Number: false,
					Adjust: false,
				},
				Markdown: &config.MarkdownConfig{
					Tables: &config.TableCustomConfig{
						ObjectCustomConfig: &config.ObjectCustomConfig{
							ShowLogicalName: true,
							Order:          []string{"Type", "Name", "Logical Name"},
						},
					},
				},
			},
			wantColumns: 8, // Type, Name, Logical Name, Default, Nullable, Children, Parents, Comment (custom order)
			wantHeaders: []string{"Type", "Name", "Logical Name", "Default", "Nullable", "Children", "Parents", "Comment"},
			description: "Should reorder columns according to configuration",
		},
		{
			name: "with column aliases",
			table: &schema.Table{
				Name: "users",
				Columns: []*schema.Column{
					{
						Name:        "id",
						Type:        "int",
						LogicalName: "ユーザーID",
						Comment:     "システム内でユーザーを一意に識別するID",
					},
				},
			},
			config: &config.Config{
				Format: config.Format{
					Number: false,
					Adjust: false,
				},
				Markdown: &config.MarkdownConfig{
					Tables: &config.TableCustomConfig{
						ObjectCustomConfig: &config.ObjectCustomConfig{
							ShowLogicalName: true,
							Aliases: map[string]string{
								"Name":         "カラム名",
								"Logical Name": "論理名",
								"Type":         "データ型",
								"Default":      "デフォルト値",
								"Nullable":     "NULL許可",
								"Comment":      "コメント",
							},
						},
					},
				},
			},
			wantColumns: 8, // All columns with Japanese aliases
			wantHeaders: []string{"カラム名", "論理名", "データ型", "デフォルト値", "NULL許可", "Children", "Parents", "コメント"},
			description: "Should apply aliases to column headers",
		},
		{
			name: "table-specific configuration",
			table: &schema.Table{
				Name: "users",
				Columns: []*schema.Column{
					{
						Name:        "id",
						Type:        "int",
						LogicalName: "ユーザーID",
						Comment:     "システム内でユーザーを一意に識別するID",
					},
				},
			},
			config: &config.Config{
				Format: config.Format{
					Number: false,
					Adjust: false,
				},
				Markdown: &config.MarkdownConfig{
					Tables: &config.TableCustomConfig{
						ObjectCustomConfig: &config.ObjectCustomConfig{
							ShowLogicalName: false,
						},
						Specific: map[string]*config.ObjectCustomConfig{
							"users": {
								ShowLogicalName: true,
								Order:          []string{"Logical Name", "Name", "Type"},
								Aliases: map[string]string{
									"Logical Name": "論理名",
									"Name":         "物理名",
									"Type":         "型",
								},
							},
						},
					},
				},
			},
			wantColumns: 8, // All columns with specific config
			wantHeaders: []string{"論理名", "物理名", "型", "Default", "Nullable", "Children", "Parents", "Comment"},
			description: "Should apply table-specific configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewExtended(tt.config)
			data := m.makeTableTemplateDataWithCustomization(tt.table)

			// Check if template data is properly structured
			if data == nil {
				t.Fatal("makeTableTemplateDataWithCustomization returned nil")
			}

			// Check if required keys exist
			requiredKeys := []string{"Table", "Columns", "Viewpoints", "Constraints", "Indexes", "Triggers", "ReferencedTables"}
			for _, key := range requiredKeys {
				if _, exists := data[key]; !exists {
					t.Errorf("Expected key %s not found in template data", key)
				}
			}

			// Check columns data structure
			columnsData, ok := data["Columns"].([][]string)
			if !ok {
				t.Fatal("Columns data is not [][]string")
			}

			if len(columnsData) < 2 {
				t.Fatal("Columns data should have at least header and separator rows")
			}

			// Check column count
			headerRow := columnsData[0]
			if len(headerRow) != tt.wantColumns {
				t.Errorf("Expected %d columns, got %d. Headers: %v", tt.wantColumns, len(headerRow), headerRow)
			}

			// Check specific headers if provided
			if len(tt.wantHeaders) > 0 {
				if len(headerRow) != len(tt.wantHeaders) {
					t.Errorf("Header count mismatch. Expected %d, got %d", len(tt.wantHeaders), len(headerRow))
				} else {
					for i, expectedHeader := range tt.wantHeaders {
						if headerRow[i] != expectedHeader {
							t.Errorf("Header mismatch at index %d. Expected %q, got %q", i, expectedHeader, headerRow[i])
						}
					}
				}
			}

			// Check that we have data rows for each column
			expectedDataRows := len(tt.table.Columns)
			actualDataRows := len(columnsData) - 2 // Exclude header and separator rows
			if actualDataRows != expectedDataRows {
				t.Errorf("Expected %d data rows, got %d", expectedDataRows, actualDataRows)
			}

			// Check logical name data if enabled
			if tt.config.Markdown != nil && tt.config.Markdown.Tables != nil {
				var customConfig *config.ObjectCustomConfig
				if tt.config.Markdown.Tables.Specific != nil {
					if specific, exists := tt.config.Markdown.Tables.Specific[tt.table.Name]; exists {
						customConfig = specific
					}
				}
				if customConfig == nil {
					customConfig = tt.config.Markdown.Tables.ObjectCustomConfig
				}

				if customConfig != nil && customConfig.ShowLogicalName {
					// Find logical name column index
					logicalNameIndex := -1
					for i, header := range headerRow {
						if header == "Logical Name" || header == "論理名" {
							logicalNameIndex = i
							break
						}
					}

					if logicalNameIndex == -1 {
						t.Error("Logical name column not found in header when ShowLogicalName is true")
					} else {
						// Check that logical name data is present in data rows
						for i := 2; i < len(columnsData); i++ {
							dataRow := columnsData[i]
							if len(dataRow) > logicalNameIndex {
								logicalNameValue := dataRow[logicalNameIndex]
								columnIndex := i - 2
								if columnIndex < len(tt.table.Columns) {
									expectedLogical := tt.table.Columns[columnIndex].GetLogicalNameOrFallback()
									if logicalNameValue != expectedLogical {
										t.Errorf("Row %d: expected logical name %q, got %q", i, expectedLogical, logicalNameValue)
									}
								}
							}
						}
					}
				}
			}

			t.Logf("Test '%s' passed: %s", tt.name, tt.description)
		})
	}
}

func TestMdExtended_BackwardCompatibility(t *testing.T) {
	// Test that existing functionality still works without customization config
	table := &schema.Table{
		Name: "test_table",
		Columns: []*schema.Column{
			{
				Name:    "id",
				Type:    "int",
				Comment: "Primary key",
			},
			{
				Name:    "name",
				Type:    "varchar(100)",
				Comment: "Name field",
			},
		},
	}

	config := &config.Config{
		Format: config.Format{
			Number: false,
			Adjust: false,
		},
		// No Markdown config - should work with defaults
	}

	m := NewExtended(config)
	data := m.makeTableTemplateDataWithCustomization(table)

	if data == nil {
		t.Fatal("makeTableTemplateDataWithCustomization returned nil")
	}

	columnsData, ok := data["Columns"].([][]string)
	if !ok {
		t.Fatal("Columns data is not [][]string")
	}

	// Should have standard columns: Name, Type, Default, Nullable, ...
	headerRow := columnsData[0]
	expectedMinColumns := 4
	if len(headerRow) < expectedMinColumns {
		t.Errorf("Expected at least %d columns, got %d", expectedMinColumns, len(headerRow))
	}

	// Should have data for each column
	expectedDataRows := len(table.Columns)
	actualDataRows := len(columnsData) - 2
	if actualDataRows != expectedDataRows {
		t.Errorf("Expected %d data rows, got %d", expectedDataRows, actualDataRows)
	}
}

func TestMdExtended_CustomizerIntegration(t *testing.T) {
	// Test integration with MarkdownCustomizer
	table := &schema.Table{
		Name: "integration_test",
		Columns: []*schema.Column{
			{
				Name:        "user_id",
				Type:        "bigint",
				LogicalName: "ユーザーID",
				Comment:     "ユーザーを識別するID",
			},
		},
	}

	config := &config.Config{
		Format: config.Format{
			Number: false,
			Adjust: false,
		},
		Markdown: &config.MarkdownConfig{
			Tables: &config.TableCustomConfig{
				ObjectCustomConfig: &config.ObjectCustomConfig{
					ShowLogicalName: true,
					Order:          []string{"Logical Name", "Name", "Type"},
					Aliases: map[string]string{
						"Name":         "物理名",
						"Logical Name": "論理名",
						"Type":         "データ型",
					},
				},
			},
		},
	}

	m := NewExtended(config)

	// Verify that customizer is properly initialized
	if m.customizer == nil {
		t.Fatal("MarkdownCustomizer not initialized")
	}

	data := m.makeTableTemplateDataWithCustomization(table)
	columnsData := data["Columns"].([][]string)
	headerRow := columnsData[0]

	// Check that aliases are applied in headers
	expectedHeaders := []string{"論理名", "物理名", "データ型"}
	foundHeaders := 0
	for _, expected := range expectedHeaders {
		for _, actual := range headerRow {
			if actual == expected {
				foundHeaders++
				break
			}
		}
	}

	if foundHeaders != len(expectedHeaders) {
		t.Errorf("Expected all aliases to be applied. Found %d out of %d. Headers: %v", foundHeaders, len(expectedHeaders), headerRow)
	}
}

func TestMdExtended_ProcessTemplateDataWithCustomization(t *testing.T) {
	// Test the public API method
	table := &schema.Table{
		Name: "public_api_test",
		Columns: []*schema.Column{
			{
				Name:        "test_column",
				Type:        "text",
				LogicalName: "テストカラム",
			},
		},
	}

	config := &config.Config{
		Format: config.Format{
			Number: false,
			Adjust: false,
		},
		Markdown: &config.MarkdownConfig{
			Tables: &config.TableCustomConfig{
				ObjectCustomConfig: &config.ObjectCustomConfig{
					ShowLogicalName: true,
				},
			},
		},
	}

	m := NewExtended(config)
	data := m.ProcessTemplateDataWithCustomization(table)

	if data == nil {
		t.Fatal("ProcessTemplateDataWithCustomization returned nil")
	}

	columnsData := data["Columns"].([][]string)
	if len(columnsData) < 3 { // header, separator, data
		t.Fatal("Expected at least 3 rows in columns data")
	}

	// Check that logical name is in the header
	headerRow := columnsData[0]
	hasLogicalName := false
	for _, header := range headerRow {
		if header == "Logical Name" {
			hasLogicalName = true
			break
		}
	}

	if !hasLogicalName {
		t.Error("Expected 'Logical Name' column in header")
	}
}