package config

import (
	"testing"

	"github.com/goccy/go-yaml"
)

func TestMarkdownConfig_GetObjectConfig(t *testing.T) {
	mc := &MarkdownConfig{
		Database: &ObjectCustomConfig{
			ShowLogicalName: true,
			Order:           []string{"name", "logical_name"},
			Aliases:         map[string]string{"name": "データベース名"},
		},
		Tables: &TableCustomConfig{
			ObjectCustomConfig: &ObjectCustomConfig{
				ShowLogicalName: true,
				Order:           []string{"name", "logical_name", "comment"},
				Aliases:         map[string]string{"name": "テーブル名"},
			},
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
			want: &ObjectCustomConfig{
				ShowLogicalName: true,
				Order:           []string{"name", "logical_name"},
				Aliases:         map[string]string{"name": "データベース名"},
			},
		},
		{
			name:       "table config",
			objectType: "table",
			want: &ObjectCustomConfig{
				ShowLogicalName: true,
				Order:           []string{"name", "logical_name", "comment"},
				Aliases:         map[string]string{"name": "テーブル名"},
			},
		},
		{
			name:       "unknown object type",
			objectType: "unknown",
			want:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mc.GetObjectConfig(tt.objectType)
			if !equalObjectCustomConfig(got, tt.want) {
				t.Errorf("GetObjectConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkdownConfig_GetTableSpecificConfig(t *testing.T) {
	mc := &MarkdownConfig{
		Tables: &TableCustomConfig{
			ObjectCustomConfig: &ObjectCustomConfig{
				ShowLogicalName: true,
				Order:           []string{"name", "logical_name", "comment"},
				Aliases:         map[string]string{"name": "テーブル名"},
			},
			Specific: map[string]*ObjectCustomConfig{
				"users": {
					ShowLogicalName: false,
					Order:           []string{"name", "type", "comment"},
					Aliases:         map[string]string{"name": "ユーザー名"},
				},
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
			want: &ObjectCustomConfig{
				ShowLogicalName: false,
				Order:           []string{"name", "type", "comment"},
				Aliases:         map[string]string{"name": "ユーザー名"},
			},
		},
		{
			name:      "non-existing specific config",
			tableName: "posts",
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mc.GetTableSpecificConfig(tt.tableName)
			if !equalObjectCustomConfig(got, tt.want) {
				t.Errorf("GetTableSpecificConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkdownConfig_GetEffectiveConfig(t *testing.T) {
	mc := &MarkdownConfig{
		Tables: &TableCustomConfig{
			ObjectCustomConfig: &ObjectCustomConfig{
				ShowLogicalName: true,
				Order:           []string{"name", "logical_name", "comment"},
				Aliases:         map[string]string{"name": "テーブル名", "comment": "コメント"},
			},
			Specific: map[string]*ObjectCustomConfig{
				"users": {
					ShowLogicalName: false,
					Order:           []string{"name", "type", "comment"},
					Aliases:         map[string]string{"name": "ユーザー名"},
				},
			},
		},
	}

	tests := []struct {
		name         string
		objectType   string
		specificName string
		want         *ObjectCustomConfig
	}{
		{
			name:         "with specific config - overrides global",
			objectType:   "table",
			specificName: "users",
			want: &ObjectCustomConfig{
				ShowLogicalName: false,
				Order:           []string{"name", "type", "comment"},
				Aliases:         map[string]string{"name": "ユーザー名", "comment": "コメント"},
			},
		},
		{
			name:         "without specific config - uses global",
			objectType:   "table",
			specificName: "posts",
			want: &ObjectCustomConfig{
				ShowLogicalName: true,
				Order:           []string{"name", "logical_name", "comment"},
				Aliases:         map[string]string{"name": "テーブル名", "comment": "コメント"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mc.GetEffectiveConfig(tt.objectType, tt.specificName)
			if !equalObjectCustomConfig(got, tt.want) {
				t.Errorf("GetEffectiveConfig() = %v, want %v", got, tt.want)
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
			name: "invalid config - duplicate order fields",
			config: &MarkdownConfig{
				Tables: &TableCustomConfig{
					ObjectCustomConfig: &ObjectCustomConfig{
						Order: []string{"name", "name", "comment"},
					},
				},
			},
			wantErr: true,
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
			wantErr: true,
		},
		{
			name: "invalid config - empty alias value",
			config: &MarkdownConfig{
				Tables: &TableCustomConfig{
					ObjectCustomConfig: &ObjectCustomConfig{
						Aliases: map[string]string{"name": ""},
					},
				},
			},
			wantErr: true,
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
					len(mc.Database.Order) > 0 &&
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
				t.Errorf("SetDefaults() failed validation check")
			}
		})
	}
}

func TestMarkdownConfig_YAML_Serialization(t *testing.T) {
	config := &MarkdownConfig{
		Database: &ObjectCustomConfig{
			ShowLogicalName: true,
			Order:           []string{"name", "logical_name", "comment"},
			Aliases:         map[string]string{"name": "データベース名"},
		},
		Tables: &TableCustomConfig{
			ObjectCustomConfig: &ObjectCustomConfig{
				ShowLogicalName: true,
				Order:           []string{"name", "logical_name", "comment", "type"},
				Aliases:         map[string]string{"name": "テーブル名"},
			},
			Specific: map[string]*ObjectCustomConfig{
				"users": {
					Order:   []string{"name", "type", "comment"},
					Aliases: map[string]string{"name": "ユーザー名"},
				},
			},
		},
		Columns: &ColumnCustomConfig{
			ObjectCustomConfig: &ObjectCustomConfig{
				ShowLogicalName: true,
				Order:           []string{"name", "logical_name", "type", "nullable"},
				Aliases:         map[string]string{"name": "カラム名", "type": "データ型"},
			},
		},
	}

	// Serialize to YAML
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal YAML: %v", err)
	}

	// Deserialize from YAML
	var deserialized MarkdownConfig
	err = yaml.Unmarshal(yamlData, &deserialized)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Verify the deserialized config
	if !deserialized.Database.ShowLogicalName {
		t.Error("Database.ShowLogicalName should be true")
	}

	if len(deserialized.Database.Order) != 3 {
		t.Errorf("Database.Order length should be 3, got %d", len(deserialized.Database.Order))
	}

	if deserialized.Database.Aliases["name"] != "データベース名" {
		t.Errorf("Database alias for 'name' should be 'データベース名', got %s", deserialized.Database.Aliases["name"])
	}

	if deserialized.Tables.Specific["users"].Aliases["name"] != "ユーザー名" {
		t.Errorf("Table specific alias should be 'ユーザー名', got %s", deserialized.Tables.Specific["users"].Aliases["name"])
	}
}

func TestMergeObjectConfigs(t *testing.T) {
	global := &ObjectCustomConfig{
		ShowLogicalName: true,
		Order:           []string{"name", "logical_name", "comment"},
		Aliases:         map[string]string{"name": "グローバル名", "comment": "コメント"},
	}

	specific := &ObjectCustomConfig{
		ShowLogicalName: false,
		Order:           []string{"name", "type"},
		Aliases:         map[string]string{"name": "特定名", "type": "タイプ"},
	}

	tests := []struct {
		name     string
		global   *ObjectCustomConfig
		specific *ObjectCustomConfig
		want     *ObjectCustomConfig
	}{
		{
			name:     "both nil",
			global:   nil,
			specific: nil,
			want:     nil,
		},
		{
			name:     "global nil",
			global:   nil,
			specific: specific,
			want:     specific,
		},
		{
			name:     "specific nil",
			global:   global,
			specific: nil,
			want:     global,
		},
		{
			name:     "merge both",
			global:   global,
			specific: specific,
			want: &ObjectCustomConfig{
				ShowLogicalName: false,                    // from specific
				Order:           []string{"name", "type"}, // from specific
				Aliases: map[string]string{
					"name":    "特定名",  // from specific
					"comment": "コメント", // from global
					"type":    "タイプ",  // from specific
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeObjectConfigs(tt.global, tt.specific)
			if !equalObjectCustomConfig(got, tt.want) {
				t.Errorf("mergeObjectConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare ObjectCustomConfig instances
func equalObjectCustomConfig(a, b *ObjectCustomConfig) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if a.ShowLogicalName != b.ShowLogicalName {
		return false
	}

	if len(a.Order) != len(b.Order) {
		return false
	}
	for i, v := range a.Order {
		if v != b.Order[i] {
			return false
		}
	}

	if len(a.Aliases) != len(b.Aliases) {
		return false
	}
	for k, v := range a.Aliases {
		if b.Aliases[k] != v {
			return false
		}
	}

	return true
}
