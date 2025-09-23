package md

import (
	"testing"

	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/schema"
)

// Test data processing with customization
func TestMd_makeTableTemplateDataWithCustomization(t *testing.T) {
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
							Aliases: map[string]string{
								"Name":    "物理名",
								"Type":    "型",
								"Comment": "説明",
							},
						},
					},
				},
			},
			wantColumns: 7,
			wantHeaders: []string{"物理名", "型", "Default", "Nullable", "Children", "Parents", "説明"},
			description: "Should apply column aliases when configured",
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
							Order: []string{"Comment", "Type", "Name"},
						},
					},
				},
			},
			wantColumns: 7,
			wantHeaders: []string{"Comment", "Type", "Name", "Default", "Nullable", "Children", "Parents"},
			description: "Should reorder columns when custom order specified",
		},
		{
			name: "with logical name and aliases",
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
								"Name":         "物理名",
								"Logical Name": "論理名",
								"Type":         "型",
							},
						},
					},
				},
			},
			wantColumns: 8,
			wantHeaders: []string{"物理名", "論理名", "型", "Default", "Nullable", "Children", "Parents", "Comment"},
			description: "Should support both logical name and aliases",
		},
		{
			name: "with table-specific configuration",
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
								Aliases: map[string]string{
									"Name":         "物理名",
									"Logical Name": "論理名",
									"Type":         "型",
								},
							},
						},
					},
				},
			},
			wantColumns: 8,
			wantHeaders: []string{"物理名", "論理名", "型", "Default", "Nullable", "Children", "Parents", "Comment"},
			description: "Should apply table-specific configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.config)
			data := m.makeTableTemplateData(tt.table)

			// Check if template data is properly structured
			if data == nil {
				t.Fatal("makeTableTemplateData returned nil")
			}

			// Check if required keys exist
			requiredKeys := []string{"Table", "Columns", "Viewpoints", "Constraints", "Indexes", "Triggers", "ReferencedTables"}
			for _, key := range requiredKeys {
				if _, exists := data[key]; !exists {
					t.Errorf("Expected key %s not found in template data", key)
				}
			}

			// Check columns data
			columnsData, ok := data["Columns"].([][]string)
			if !ok {
				t.Fatalf("Columns should be [][]string, got %T", data["Columns"])
			}

			if len(columnsData) < 2 {
				t.Fatalf("Columns data should have at least 2 rows (header + separator), got %d", len(columnsData))
			}

			// Check header structure
			headerRow := columnsData[0]
			if len(headerRow) != tt.wantColumns {
				t.Errorf("Expected %d columns, got %d. Headers: %v", tt.wantColumns, len(headerRow), headerRow)
			}

			// Check specific headers if provided
			if tt.wantHeaders != nil {
				if len(headerRow) != len(tt.wantHeaders) {
					t.Errorf("Header count mismatch. Expected %d, got %d", len(tt.wantHeaders), len(headerRow))
				}
				for i, expectedHeader := range tt.wantHeaders {
					if i < len(headerRow) && headerRow[i] != expectedHeader {
						t.Errorf("Header[%d]: expected %q, got %q", i, expectedHeader, headerRow[i])
					}
				}
			}

			// Verify data rows
			if len(columnsData) > 2 {
				dataRow := columnsData[2] // First data row
				if len(dataRow) != tt.wantColumns {
					t.Errorf("Data row should have %d columns, got %d", tt.wantColumns, len(dataRow))
				}
			}

			t.Logf("Test %q passed: %s", tt.name, tt.description)
		})
	}
}

func TestMd_BackwardCompatibility(t *testing.T) {
	// Test that existing functionality works without any configuration
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
	}

	m := New(config)
	result := m.makeTableTemplateData(table)

	// Should work exactly as before
	if _, ok := result["Table"]; !ok {
		t.Error("Result should contain 'Table' key")
	}
	if _, ok := result["Columns"]; !ok {
		t.Error("Result should contain 'Columns' key")
	}

	columnsData, ok := result["Columns"].([][]string)
	if !ok {
		t.Fatalf("Columns data should be [][]string, got %T", result["Columns"])
	}

	if len(columnsData) < 2 {
		t.Fatalf("Columns data should have at least 2 rows, got %d", len(columnsData))
	}

	// Verify standard column count (without customization)
	headerRow := columnsData[0]
	expectedColumns := 7 // Name, Type, Default, Nullable, Children, Parents, Comment
	if len(headerRow) != expectedColumns {
		t.Errorf("Expected %d columns, got %d", expectedColumns, len(headerRow))
	}

	t.Log("Backward compatibility test passed")
}

func TestMd_CustomizerIntegration(t *testing.T) {
	// Test integration with MarkdownCustomizer
	config := &config.Config{
		Format: config.Format{
			Number: false,
			Adjust: false,
		},
		Markdown: &config.MarkdownConfig{
			Tables: &config.TableCustomConfig{
				ObjectCustomConfig: &config.ObjectCustomConfig{
					ShowLogicalName: true,
					Aliases: map[string]string{
						"Name": "カラム名",
						"Type": "データ型",
					},
				},
			},
		},
	}

	m := New(config)

	// Verify that customizer is properly initialized
	if m.customizer == nil {
		t.Fatal("MarkdownCustomizer not initialized")
	}

	table := &schema.Table{
		Name: "test_table",
		Columns: []*schema.Column{
			{
				Name:        "id",
				Type:        "int",
				LogicalName: "識別子",
				Comment:     "Primary key",
			},
		},
	}

	result := m.makeTableTemplateData(table)
	columnsData := result["Columns"].([][]string)
	headerRow := columnsData[0]

	// Check if aliases are applied
	found := false
	for _, header := range headerRow {
		if header == "カラム名" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Alias 'カラム名' not found in headers")
	}

	t.Log("Customizer integration test passed")
}

func TestMd_ProcessTemplateDataWithCustomization(t *testing.T) {
	table := &schema.Table{
		Name: "users",
		Columns: []*schema.Column{
			{
				Name:        "id",
				Type:        "int",
				LogicalName: "ユーザーID",
				Comment:     "Primary key",
			},
			{
				Name:        "name",
				Type:        "varchar(100)",
				LogicalName: "ユーザー名",
				Comment:     "User name",
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
					Order:           []string{"Name", "Logical Name", "Type", "Comment"},
					Aliases: map[string]string{
						"Name":         "カラム名",
						"Logical Name": "論理名",
						"Type":         "データ型",
						"Comment":      "説明",
					},
				},
			},
		},
	}

	m := New(config)
	result := m.makeTableTemplateData(table)

	// Verify structure
	columnsData, ok := result["Columns"].([][]string)
	if !ok {
		t.Fatalf("Expected [][]string, got %T", result["Columns"])
	}

	headerRow := columnsData[0]
	expectedHeaders := []string{"カラム名", "論理名", "データ型", "説明", "Default", "Nullable", "Children", "Parents"}

	if len(headerRow) != len(expectedHeaders) {
		t.Fatalf("Expected %d headers, got %d", len(expectedHeaders), len(headerRow))
	}

	for i, expected := range expectedHeaders {
		if i < len(headerRow) && headerRow[i] != expected {
			t.Errorf("Header[%d]: expected %q, got %q", i, expected, headerRow[i])
		}
	}

	// Verify logical name data
	if len(columnsData) > 2 {
		firstRow := columnsData[2]
		logicalNameIndex := 1 // Second column should be logical name
		if firstRow[logicalNameIndex] != "ユーザーID" {
			t.Errorf("Expected logical name 'ユーザーID', got '%s'", firstRow[logicalNameIndex])
		}
	}

	t.Log("Template data processing test passed")
}
