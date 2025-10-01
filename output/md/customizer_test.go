package md

import (
	"reflect"
	"testing"

	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/schema"
)

func TestNewDefaultMarkdownCustomizer(t *testing.T) {
	customizer := NewDefaultMarkdownCustomizer()
	if customizer == nil {
		t.Error("NewDefaultMarkdownCustomizer should return a non-nil instance")
	}
}

func TestDefaultMarkdownCustomizer_CustomizeTableOutput(t *testing.T) {
	customizer := NewDefaultMarkdownCustomizer()

	tests := []struct {
		name     string
		table    *schema.Table
		config   *config.ObjectCustomConfig
		expected *CustomizedTableData
	}{
		{
			name:     "nil table",
			table:    nil,
			config:   nil,
			expected: nil,
		},
		{
			name: "table with no config",
			table: &schema.Table{
				Name:        "users",
				Comment:     "ユーザーテーブル",
				LogicalName: "ユーザー",
				Columns: []*schema.Column{
					{
						Name:        "id",
						Type:        "integer",
						LogicalName: "ID",
						Comment:     "ユーザーID",
					},
					{
						Name:        "name",
						Type:        "varchar(100)",
						LogicalName: "名前",
						Comment:     "ユーザー名",
					},
				},
			},
			config: nil,
			expected: &CustomizedTableData{
				Table: &schema.Table{
					Name:        "users",
					Comment:     "ユーザーテーブル",
					LogicalName: "ユーザー",
					Columns: []*schema.Column{
						{
							Name:        "id",
							Type:        "integer",
							LogicalName: "ID",
							Comment:     "ユーザーID",
						},
						{
							Name:        "name",
							Type:        "varchar(100)",
							LogicalName: "名前",
							Comment:     "ユーザー名",
						},
					},
				},
				Columns: []*CustomizedColumnData{
					{
						Column: &schema.Column{
							Name:        "id",
							Type:        "integer",
							LogicalName: "ID",
							Comment:     "ユーザーID",
						},
						DisplayName: "id",
						LogicalName: "ID",
						Show:        true,
					},
					{
						Column: &schema.Column{
							Name:        "name",
							Type:        "varchar(100)",
							LogicalName: "名前",
							Comment:     "ユーザー名",
						},
						DisplayName: "name",
						LogicalName: "名前",
						Show:        true,
					},
				},
				Order:   []string{"name", "LogicalName", "comment", "type"},
				Aliases: map[string]string{},
			},
		},
		{
			name: "table with custom config",
			table: &schema.Table{
				Name:        "products",
				Comment:     "商品テーブル",
				LogicalName: "商品",
				Columns: []*schema.Column{
					{
						Name:        "product_id",
						Type:        "integer",
						LogicalName: "商品ID",
						Comment:     "商品の一意識別子",
					},
				},
			},
			config: &config.ObjectCustomConfig{
				ShowLogicalName: true,
				Order:           []string{"LogicalName", "name", "type", "comment"},
				Aliases: map[string]string{
					"product_id": "商品ID",
					"name":       "名前",
				},
			},
			expected: &CustomizedTableData{
				Table: &schema.Table{
					Name:        "products",
					Comment:     "商品テーブル",
					LogicalName: "商品",
					Columns: []*schema.Column{
						{
							Name:        "product_id",
							Type:        "integer",
							LogicalName: "商品ID",
							Comment:     "商品の一意識別子",
						},
					},
				},
				Columns: []*CustomizedColumnData{
					{
						Column: &schema.Column{
							Name:        "product_id",
							Type:        "integer",
							LogicalName: "商品ID",
							Comment:     "商品の一意識別子",
						},
						DisplayName: "商品ID",
						LogicalName: "商品ID",
						Show:        true,
					},
				},
				Order: []string{"LogicalName", "name", "type", "comment"},
				Aliases: map[string]string{
					"product_id": "商品ID",
					"name":       "名前",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := customizer.CustomizeTableOutput(tt.table, tt.config)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CustomizeTableOutput() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestDefaultMarkdownCustomizer_CustomizeColumnOutput(t *testing.T) {
	customizer := NewDefaultMarkdownCustomizer()

	tests := []struct {
		name     string
		columns  []*schema.Column
		config   *config.ObjectCustomConfig
		expected []*CustomizedColumnData
	}{
		{
			name:     "nil columns",
			columns:  nil,
			config:   nil,
			expected: nil,
		},
		{
			name:     "empty columns",
			columns:  []*schema.Column{},
			config:   nil,
			expected: []*CustomizedColumnData{},
		},
		{
			name: "columns with no config",
			columns: []*schema.Column{
				{
					Name:        "id",
					Type:        "integer",
					LogicalName: "ID",
					Comment:     "識別子",
				},
				{
					Name:        "email",
					Type:        "varchar(255)",
					LogicalName: "メールアドレス",
					Comment:     "ユーザーのメールアドレス",
				},
			},
			config: nil,
			expected: []*CustomizedColumnData{
				{
					Column: &schema.Column{
						Name:        "id",
						Type:        "integer",
						LogicalName: "ID",
						Comment:     "識別子",
					},
					DisplayName: "id",
					LogicalName: "ID",
					Show:        true,
				},
				{
					Column: &schema.Column{
						Name:        "email",
						Type:        "varchar(255)",
						LogicalName: "メールアドレス",
						Comment:     "ユーザーのメールアドレス",
					},
					DisplayName: "email",
					LogicalName: "メールアドレス",
					Show:        true,
				},
			},
		},
		{
			name: "columns with aliases",
			columns: []*schema.Column{
				{
					Name:        "user_id",
					Type:        "bigint",
					LogicalName: "ユーザーID",
					Comment:     "ユーザー識別子",
				},
			},
			config: &config.ObjectCustomConfig{
				Aliases: map[string]string{
					"user_id": "利用者ID",
				},
			},
			expected: []*CustomizedColumnData{
				{
					Column: &schema.Column{
						Name:        "user_id",
						Type:        "bigint",
						LogicalName: "ユーザーID",
						Comment:     "ユーザー識別子",
					},
					DisplayName: "利用者ID",
					LogicalName: "ユーザーID",
					Show:        true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := customizer.CustomizeColumnOutput(tt.columns, tt.config)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CustomizeColumnOutput() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestDefaultMarkdownCustomizer_ApplyAliases(t *testing.T) {
	customizer := NewDefaultMarkdownCustomizer()

	tests := []struct {
		name      string
		fieldName string
		aliases   map[string]string
		expected  string
	}{
		{
			name:      "nil aliases",
			fieldName: "name",
			aliases:   nil,
			expected:  "name",
		},
		{
			name:      "empty aliases",
			fieldName: "name",
			aliases:   map[string]string{},
			expected:  "name",
		},
		{
			name:      "field not in aliases",
			fieldName: "name",
			aliases: map[string]string{
				"id": "ID",
			},
			expected: "name",
		},
		{
			name:      "field in aliases",
			fieldName: "name",
			aliases: map[string]string{
				"name": "名前",
				"id":   "ID",
			},
			expected: "名前",
		},
		{
			name:      "empty alias value",
			fieldName: "name",
			aliases: map[string]string{
				"name": "",
			},
			expected: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := customizer.ApplyAliases(tt.fieldName, tt.aliases)
			if result != tt.expected {
				t.Errorf("ApplyAliases() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDefaultMarkdownCustomizer_GetDisplayOrder(t *testing.T) {
	customizer := NewDefaultMarkdownCustomizer()

	tests := []struct {
		name         string
		defaultOrder []string
		configOrder  []string
		expected     []string
	}{
		{
			name:         "empty config order",
			defaultOrder: []string{"name", "type", "comment"},
			configOrder:  []string{},
			expected:     []string{"name", "type", "comment"},
		},
		{
			name:         "nil config order",
			defaultOrder: []string{"name", "type", "comment"},
			configOrder:  nil,
			expected:     []string{"name", "type", "comment"},
		},
		{
			name:         "config order subset of default",
			defaultOrder: []string{"name", "type", "comment", "nullable"},
			configOrder:  []string{"type", "name"},
			expected:     []string{"type", "name"},
		},
		{
			name:         "config order with extra fields",
			defaultOrder: []string{"name", "type"},
			configOrder:  []string{"custom_field", "name", "extra"},
			expected:     []string{"custom_field", "name", "extra"},
		},
		{
			name:         "completely different orders",
			defaultOrder: []string{"a", "b", "c"},
			configOrder:  []string{"x", "y", "z"},
			expected:     []string{"x", "y", "z"},
		},
		{
			name:         "overlapping orders",
			defaultOrder: []string{"name", "type", "comment", "nullable"},
			configOrder:  []string{"comment", "name", "custom"},
			expected:     []string{"comment", "name", "custom"},
		},
	{
		name:         "order specifies only subset - should not include defaults",
		defaultOrder: []string{"Name", "Type", "Default", "Nullable", "Comment"},
		configOrder:  []string{"Name", "Type"},
		expected:     []string{"Name", "Type"},
	},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := customizer.GetDisplayOrder(tt.defaultOrder, tt.configOrder)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetDisplayOrder() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test that DefaultMarkdownCustomizer implements MarkdownCustomizer interface
func TestDefaultMarkdownCustomizer_ImplementsInterface(t *testing.T) {
	var _ MarkdownCustomizer = (*DefaultMarkdownCustomizer)(nil)
}

// Test handling of nil columns in column customization
func TestDefaultMarkdownCustomizer_CustomizeColumnOutput_NilColumns(t *testing.T) {
	customizer := NewDefaultMarkdownCustomizer()

	columns := []*schema.Column{
		{
			Name:        "id",
			Type:        "integer",
			LogicalName: "ID",
		},
		nil, // nil column
		{
			Name:        "name",
			Type:        "varchar(100)",
			LogicalName: "名前",
		},
	}

	result := customizer.CustomizeColumnOutput(columns, nil)

	// Check that the nil column is handled correctly
	if len(result) != 3 {
		t.Errorf("Expected 3 result items, got %d", len(result))
	}

	// First column should be properly processed
	if result[0] == nil || result[0].Column.Name != "id" {
		t.Error("First column not processed correctly")
	}

	// Second column (nil input) should result in nil or be skipped
	if result[1] != nil && result[1].Column != nil {
		t.Error("Nil column should not be processed")
	}

	// Third column should be properly processed
	if result[2] == nil || result[2].Column.Name != "name" {
		t.Error("Third column not processed correctly")
	}
}
