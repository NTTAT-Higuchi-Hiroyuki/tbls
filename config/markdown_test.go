package config

import (
	"testing"

	"github.com/goccy/go-yaml"
)

func TestMarkdownConfig_GetObjectConfig(t *testing.T) {
	config := &MarkdownConfig{
		Database: &ObjectCustomConfig{ShowLogicalName: true},
		Tables: &TableCustomConfig{
			ObjectCustomConfig: &ObjectCustomConfig{ShowLogicalName: false},
		},
	}

	tests := []struct {
		name       string
		objectType string
		want       *ObjectCustomConfig
	}{
		{
			name:       "database config",
			objectType: "database",
			want:       config.Database,
		},
		{
			name:       "table config",
			objectType: "tables",
			want:       config.Tables.ObjectCustomConfig,
		},
		{
			name:       "unknown object type",
			objectType: "unknown",
			want:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.GetObjectConfig(tt.objectType)
			if got != tt.want {
				t.Errorf("GetObjectConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkdownConfig_GetTableSpecificConfig(t *testing.T) {
	config := &MarkdownConfig{
		Tables: &TableCustomConfig{
			ObjectCustomConfig: &ObjectCustomConfig{ShowLogicalName: true},
			Specific: map[string]*ObjectCustomConfig{
				"users": {ShowLogicalName: false},
			},
		},
	}

	tests := []struct {
		name      string
		tableName string
		want      *ObjectCustomConfig
	}{
		{
			name:      "existing specific config",
			tableName: "users",
			want:      config.Tables.Specific["users"],
		},
		{
			name:      "non-existing specific config",
			tableName: "posts",
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.GetTableSpecificConfig(tt.tableName)
			if got != tt.want {
				t.Errorf("GetTableSpecificConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkdownConfig_GetEffectiveConfig(t *testing.T) {
	config := &MarkdownConfig{
		Tables: &TableCustomConfig{
			ObjectCustomConfig: &ObjectCustomConfig{
				ShowLogicalName: true,
				Order:           []string{"name", "comment"},
				Aliases:         map[string]string{"name": "テーブル名"},
			},
			Specific: map[string]*ObjectCustomConfig{
				"users": {
					ShowLogicalName: false,
					Order:           []string{"logical_name", "name"},
				},
			},
		},
	}

	tests := []struct {
		name         string
		objectType   string
		specificName string
		wantShow     bool
		wantOrder    []string
	}{
		{
			name:         "with specific config - overrides global",
			objectType:   "tables",
			specificName: "users",
			wantShow:     false,
			wantOrder:    []string{"logical_name", "name"},
		},
		{
			name:         "without specific config - uses global",
			objectType:   "tables",
			specificName: "posts",
			wantShow:     true,
			wantOrder:    []string{"name", "comment"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.GetEffectiveConfig(tt.objectType, tt.specificName)
			if got.ShowLogicalName != tt.wantShow {
				t.Errorf("GetEffectiveConfig().ShowLogicalName = %v, want %v", got.ShowLogicalName, tt.wantShow)
			}
			if len(got.Order) != len(tt.wantOrder) {
				t.Errorf("GetEffectiveConfig().Order length = %v, want %v", len(got.Order), len(tt.wantOrder))
				return
			}
			for i, o := range got.Order {
				if o != tt.wantOrder[i] {
					t.Errorf("GetEffectiveConfig().Order[%d] = %v, want %v", i, o, tt.wantOrder[i])
				}
			}
		})
	}
}

func TestMarkdownConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *MarkdownConfig
		wantErr bool
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: false,
		},
		{
			name: "valid config",
			config: &MarkdownConfig{
				Tables: &TableCustomConfig{
					ObjectCustomConfig: &ObjectCustomConfig{
						ShowLogicalName: true,
						Order:           []string{"name", "logical_name", "comment"},
						Aliases:         map[string]string{"name": "テーブル名"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid config with warnings - duplicate order fields (auto-fixed)",
			config: &MarkdownConfig{
				Tables: &TableCustomConfig{
					ObjectCustomConfig: &ObjectCustomConfig{
						Order: []string{"name", "name", "comment"},
					},
				},
			},
			wantErr: false, // 新しいバリデーションは警告を表示して自動修正するため、エラーは返さない
		},
		{
			name: "invalid config with warnings - empty alias key (auto-fixed)",
			config: &MarkdownConfig{
				Tables: &TableCustomConfig{
					ObjectCustomConfig: &ObjectCustomConfig{
						Aliases: map[string]string{"": "empty key"},
					},
				},
			},
			wantErr: false, // 新しいバリデーションは警告を表示して自動修正するため、エラーは返さない
		},
		{
			name: "invalid config with warnings - empty alias value (auto-fixed)",
			config: &MarkdownConfig{
				Tables: &TableCustomConfig{
					ObjectCustomConfig: &ObjectCustomConfig{
						Aliases: map[string]string{"name": ""},
					},
				},
			},
			wantErr: false, // 新しいバリデーションは警告を表示して自動修正するため、エラーは返さない
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarkdownConfig_SetDefaults(t *testing.T) {
	tests := []struct {
		name   string
		config *MarkdownConfig
		check  func(*MarkdownConfig) bool
	}{
		{
			name:   "nil config",
			config: nil,
			check: func(mc *MarkdownConfig) bool {
				return mc == nil
			},
		},
		{
			name:   "empty config",
			config: &MarkdownConfig{},
			check: func(mc *MarkdownConfig) bool {
				return mc.Database != nil &&
					mc.Tables != nil &&
					mc.Columns != nil &&
					len(mc.Tables.Order) > 0 &&
					len(mc.Columns.Order) > 0
			},
		},
		{
			name: "existing config - should not override",
			config: &MarkdownConfig{
				Tables: &TableCustomConfig{
					ObjectCustomConfig: &ObjectCustomConfig{
						Order: []string{"custom", "order"},
					},
				},
			},
			check: func(mc *MarkdownConfig) bool {
				return len(mc.Tables.Order) == 2 &&
					mc.Tables.Order[0] == "custom" &&
					mc.Tables.Order[1] == "order"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.SetDefaults()
			if !tt.check(tt.config) {
				t.Errorf("SetDefaults() did not set expected defaults")
			}
		})
	}
}

func TestMarkdownConfig_YAML_Serialization(t *testing.T) {
	config := &MarkdownConfig{
		Tables: &TableCustomConfig{
			ObjectCustomConfig: &ObjectCustomConfig{
				ShowLogicalName: true,
				Order:           []string{"name", "logical_name", "comment"},
				Aliases:         map[string]string{"name": "テーブル名"},
			},
			Specific: map[string]*ObjectCustomConfig{
				"users": {
					ShowLogicalName: false,
					Order:           []string{"logical_name", "name"},
					Aliases:         map[string]string{"logical_name": "論理名"},
				},
			},
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Unmarshal from YAML
	var unmarshaled MarkdownConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Check that the unmarshaled config matches the original
	if unmarshaled.Tables.ShowLogicalName != config.Tables.ShowLogicalName {
		t.Errorf("ShowLogicalName mismatch: got %v, want %v", unmarshaled.Tables.ShowLogicalName, config.Tables.ShowLogicalName)
	}
	if len(unmarshaled.Tables.Order) != len(config.Tables.Order) {
		t.Errorf("Order length mismatch: got %v, want %v", len(unmarshaled.Tables.Order), len(config.Tables.Order))
	}
	if unmarshaled.Tables.Aliases["name"] != config.Tables.Aliases["name"] {
		t.Errorf("Aliases mismatch: got %v, want %v", unmarshaled.Tables.Aliases["name"], config.Tables.Aliases["name"])
	}
}

func TestMarkdownConfig_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		config *MarkdownConfig
		want   bool
	}{
		{
			name:   "nil config",
			config: nil,
			want:   true,
		},
		{
			name: "valid config",
			config: &MarkdownConfig{
				Tables: &TableCustomConfig{
					ObjectCustomConfig: &ObjectCustomConfig{
						ShowLogicalName: true,
						Order:           []string{"name", "logical_name", "comment"},
						Aliases:         map[string]string{"name": "テーブル名"},
					},
				},
			},
			want: true,
		},
		{
			name: "invalid config - duplicate order fields",
			config: &MarkdownConfig{
				Tables: &TableCustomConfig{
					ObjectCustomConfig: &ObjectCustomConfig{
						Order: []string{"name", "name", "comment"},
					},
				},
			},
			want: false,
		},
		{
			name: "invalid config - empty alias key",
			config: &MarkdownConfig{
				Tables: &TableCustomConfig{
					ObjectCustomConfig: &ObjectCustomConfig{
						Aliases: map[string]string{"": "empty key"},
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
