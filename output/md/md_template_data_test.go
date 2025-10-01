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
							Order:           []string{"name", "logical_name", "type", "default", "nullable", "children", "parents", "comment"},
						},
					},
				},
			},
			wantColumns: 8, // Name, Logical Name, Type, Default, Nullable, Children, Parents, Comment
			wantHeaders: []string{"Name", "Logical Name", "Type", "Default", "Nullable", "Children", "Parents", "Comment"},
			description: "Should add logical name column when order explicitly includes logical_name",
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
					Columns: &config.ColumnCustomConfig{
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
							Order:           []string{"name", "logical_name", "type", "default", "nullable", "children", "parents", "comment"},
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
			description: "Should support both logical name and aliases when order includes logical_name",
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
								Order:           []string{"name", "logical_name", "type", "default", "nullable", "children", "parents", "comment"},
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
			Columns: &config.ColumnCustomConfig{
				ObjectCustomConfig: &config.ObjectCustomConfig{
					ShowLogicalName: true,
					Order:           []string{"Name", "LogicalName", "Type", "Comment"},
					Aliases: map[string]string{
						"Name":        "カラム名",
						"LogicalName": "論理名",
						"Type":        "データ型",
						"Comment":     "説明",
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
	// When Order is explicitly specified, only those fields are shown (line 428: Order: []string{"Name", "Logical Name", "Type", "Comment"})
	expectedHeaders := []string{"カラム名", "論理名", "データ型", "説明"}

	if len(headerRow) != len(expectedHeaders) {
		t.Fatalf("Expected %d headers, got %d (headers: %v)", len(expectedHeaders), len(headerRow), headerRow)
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

// TestMd_ConstraintsWithExplicitLogicalName tests constraints with explicit LogicalName in order
func TestMd_ConstraintsWithExplicitLogicalName(t *testing.T) {
	// Create test table with constraints
	table := &schema.Table{
		Name: "test_table",
		Constraints: []*schema.Constraint{
			{
				Name:        "users_pkey",
				Type:        "PRIMARY KEY",
				Def:         "PRIMARY KEY (id)",
				Comment:     "Primary key for user identification",
				LogicalName: "ユーザーID主キー",
			},
		},
	}

	cfg := &config.Config{
		Format: config.Format{
			Number: false,
			Adjust: false,
		},
		Markdown: &config.MarkdownConfig{
			Constraints: &config.ObjectCustomConfig{
				ShowLogicalName: true,
				Order:           []string{"Name", "LogicalName", "Type", "Definition", "Comment"},
				Aliases: map[string]string{
					"Name":        "制約名",
					"LogicalName": "論理名",
					"Type":        "種類",
					"Definition":  "定義",
					"Comment":     "説明",
				},
			},
		},
	}
	md := New(cfg)

	// Test constraints with explicit LogicalName
	constraintsData := md.constraintsData(table, cfg.Markdown.Constraints)

	t.Logf("Constraints data rows: %d", len(constraintsData))
	for i, row := range constraintsData {
		t.Logf("Row %d: %v", i, row)
	}

	if len(constraintsData) < 3 {
		t.Fatal("constraintsData should have at least header, separator, and data row")
	}

	// Check header
	header := constraintsData[0]
	expectedHeader := []string{"制約名", "論理名", "種類", "定義", "説明"}
	if len(header) != len(expectedHeader) {
		t.Errorf("Header length: expected %d, got %d", len(expectedHeader), len(header))
		return
	}
	for i, exp := range expectedHeader {
		if header[i] != exp {
			t.Errorf("Header[%d]: expected %q, got %q", i, exp, header[i])
		}
	}

	// Check data row
	dataRow := constraintsData[2]
	if len(dataRow) != 5 {
		t.Errorf("Data row should have 5 columns, got %d", len(dataRow))
		return
	}

	// Check that all fields are populated
	if dataRow[0] == "" {
		t.Error("Name should not be empty")
	}
	if dataRow[1] == "" {
		t.Error("LogicalName should not be empty")
	}
	if dataRow[2] == "" {
		t.Error("Type should not be empty")
	}
	if dataRow[3] == "" {
		t.Error("Definition should not be empty")
	}

	t.Logf("Name: %q", dataRow[0])
	t.Logf("LogicalName: %q", dataRow[1])
	t.Logf("Type: %q", dataRow[2])
	t.Logf("Definition: %q", dataRow[3])
	t.Logf("Comment: %q", dataRow[4])
}

// TestMd_ConstraintsWithEmptyLogicalName tests constraints with empty LogicalName (real PostgreSQL scenario)
func TestMd_ConstraintsWithEmptyLogicalName(t *testing.T) {
	// Simulate real PostgreSQL scenario: LogicalName is empty, Comment is empty
	// GetLogicalNameOrFallback() should return Name
	table := &schema.Table{
		Name: "test_table",
		Constraints: []*schema.Constraint{
			{
				Name:        "users_pkey",
				Type:        "PRIMARY KEY",
				Def:         "PRIMARY KEY (id)",
				Comment:     "",  // Empty comment (real scenario)
				LogicalName: "",  // Empty (real scenario)
			},
		},
	}

	cfg := &config.Config{
		Format: config.Format{
			Number: false,
			Adjust: false,
		},
		Markdown: &config.MarkdownConfig{
			Constraints: &config.ObjectCustomConfig{
				ShowLogicalName: true,
				Order:           []string{"Name", "LogicalName", "Type", "Definition", "Comment"},
				Aliases: map[string]string{
					"Name":        "制約名",
					"LogicalName": "論理名",
					"Type":        "種類",
					"Definition":  "定義",
					"Comment":     "説明",
				},
			},
		},
	}

	// Verify GetLogicalNameOrFallback returns Name when LogicalName is empty
	constraint := table.Constraints[0]
	logicalName := constraint.GetLogicalNameOrFallback()
	if logicalName != "users_pkey" {
		t.Errorf("GetLogicalNameOrFallback() should return Name when LogicalName is empty, got %q", logicalName)
	}

	md := New(cfg)

	// Test constraints with empty LogicalName
	constraintsData := md.constraintsData(table, cfg.Markdown.Constraints)

	t.Logf("Constraints data rows: %d", len(constraintsData))
	for i, row := range constraintsData {
		t.Logf("Row %d: %v", i, row)
	}

	if len(constraintsData) < 3 {
		t.Fatal("constraintsData should have at least header, separator, and data row")
	}

	// Check header
	header := constraintsData[0]
	expectedHeader := []string{"制約名", "論理名", "種類", "定義", "説明"}
	if len(header) != len(expectedHeader) {
		t.Errorf("Header length: expected %d, got %d", len(expectedHeader), len(header))
		return
	}
	for i, exp := range expectedHeader {
		if header[i] != exp {
			t.Errorf("Header[%d]: expected %q, got %q", i, exp, header[i])
		}
	}

	// Check data row
	dataRow := constraintsData[2]
	if len(dataRow) != 5 {
		t.Errorf("Data row should have 5 columns, got %d: %v", len(dataRow), dataRow)
		return
	}

	// Check that Name is populated
	if dataRow[0] == "" {
		t.Error("Name should not be empty")
	}

	// CRITICAL: Check that LogicalName is populated with Name (fallback)
	if dataRow[1] == "" {
		t.Error("LogicalName should not be empty (should fallback to Name)")
	}
	if dataRow[1] != "users_pkey" {
		t.Errorf("LogicalName should be 'users_pkey' (fallback), got %q", dataRow[1])
	}

	// Check other fields
	if dataRow[2] == "" {
		t.Error("Type should not be empty")
	}
	if dataRow[3] == "" {
		t.Error("Definition should not be empty")
	}

	t.Logf("Name: %q", dataRow[0])
	t.Logf("LogicalName: %q (should be 'users_pkey')", dataRow[1])
	t.Logf("Type: %q", dataRow[2])
	t.Logf("Definition: %q", dataRow[3])
	t.Logf("Comment: %q", dataRow[4])
}
