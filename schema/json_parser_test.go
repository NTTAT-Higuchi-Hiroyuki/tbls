package schema

import (
	"strings"
	"testing"
)

func TestJSONParser_Name(t *testing.T) {
	parser := NewJSONParser()
	if parser.Name() != "json" {
		t.Errorf("expected name 'json', got %s", parser.Name())
	}
}

func TestJSONParser_Priority(t *testing.T) {
	parser := NewJSONParser()
	if parser.Priority() != 10 {
		t.Errorf("expected priority 10, got %d", parser.Priority())
	}
}

func TestJSONParser_CanParse(t *testing.T) {
	parser := NewJSONParser()

	tests := []struct {
		name    string
		comment string
		want    bool
	}{
		{"空文字列", "", false},
		{"有効なJSONオブジェクト", `{"name": "test"}`, true},
		{"有効なJSON配列", `[{"name": "test"}]`, true},
		{"無効なJSON", `{name: "test"}`, false},
		{"非JSON文字列", "通常のコメント", false},
		{"JSONっぽいが無効", `{"name": "test"`, false},
		{"サイズ制限超過", strings.Repeat("a", 10000), false},
		{"短すぎる文字列", "a", false},
		{"スペースのみ", "   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.CanParse(tt.comment)
			if got != tt.want {
				t.Errorf("CanParse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONParser_ParseComment(t *testing.T) {
	parser := NewJSONParser()

	tests := []struct {
		name        string
		comment     string
		delimiter   string
		wantLogical string
		wantDesc    string
		wantTags    []string
		wantPrio    int
		wantDeprecated bool
		wantError   bool
	}{
		{
			name:        "空文字列",
			comment:     "",
			delimiter:   "|",
			wantLogical: "",
			wantDesc:    "",
			wantTags:    nil,
			wantPrio:    0,
			wantDeprecated: false,
			wantError:   false,
		},
		{
			name:        "基本的なJSONオブジェクト",
			comment:     `{"name": "ユーザー名", "description": "ユーザーの表示名"}`,
			delimiter:   "|",
			wantLogical: "ユーザー名",
			wantDesc:    "ユーザーの表示名",
			wantTags:    nil,
			wantPrio:    0,
			wantDeprecated: false,
			wantError:   false,
		},
		{
			name:        "全フィールド含むJSONオブジェクト",
			comment:     `{"name": "商品名", "description": "商品の名前", "tags": ["重要", "公開"], "priority": 5, "deprecated": true}`,
			delimiter:   "|",
			wantLogical: "商品名",
			wantDesc:    "商品の名前",
			wantTags:    []string{"重要", "公開"},
			wantPrio:    5,
			wantDeprecated: true,
			wantError:   false,
		},
		{
			name:        "代替フィールド名対応",
			comment:     `{"logical_name": "論理名", "desc": "説明文"}`,
			delimiter:   "|",
			wantLogical: "論理名",
			wantDesc:    "説明文",
			wantTags:    nil,
			wantPrio:    0,
			wantDeprecated: false,
			wantError:   false,
		},
		{
			name:        "タグが文字列の場合",
			comment:     `{"name": "テスト", "tags": "単一タグ"}`,
			delimiter:   "|",
			wantLogical: "テスト",
			wantDesc:    "",
			wantTags:    []string{"単一タグ"},
			wantPrio:    0,
			wantDeprecated: false,
			wantError:   false,
		},
		{
			name:        "JSON配列形式",
			comment:     `[{"name": "配列テスト", "description": "配列の最初の要素"}]`,
			delimiter:   "|",
			wantLogical: "配列テスト",
			wantDesc:    "配列の最初の要素",
			wantTags:    nil,
			wantPrio:    0,
			wantDeprecated: false,
			wantError:   false,
		},
		{
			name:      "無効なJSON",
			comment:   `{name: "test"}`,
			delimiter: "|",
			wantError: true,
		},
		{
			name:      "サイズ制限超過",
			comment:   `{"name": "` + strings.Repeat("a", 10000) + `"}`,
			delimiter: "|",
			wantError: true,
		},
		{
			name:      "深すぎるJSON構造",
			comment:   `{"a": {"b": {"c": {"d": {"e": {"f": "too deep"}}}}}}`,
			delimiter: "|",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseComment(tt.comment, tt.delimiter)

			if (err != nil) != tt.wantError {
				t.Errorf("ParseComment() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError {
				return
			}

			if got.LogicalName != tt.wantLogical {
				t.Errorf("ParseComment() LogicalName = %q, want %q", got.LogicalName, tt.wantLogical)
			}

			if got.Description != tt.wantDesc {
				t.Errorf("ParseComment() Description = %q, want %q", got.Description, tt.wantDesc)
			}

			if len(got.Tags) != len(tt.wantTags) {
				t.Errorf("ParseComment() Tags length = %d, want %d", len(got.Tags), len(tt.wantTags))
			} else {
				for i, tag := range got.Tags {
					if tag != tt.wantTags[i] {
						t.Errorf("ParseComment() Tags[%d] = %q, want %q", i, tag, tt.wantTags[i])
					}
				}
			}

			if got.Priority != tt.wantPrio {
				t.Errorf("ParseComment() Priority = %d, want %d", got.Priority, tt.wantPrio)
			}

			if got.Deprecated != tt.wantDeprecated {
				t.Errorf("ParseComment() Deprecated = %v, want %v", got.Deprecated, tt.wantDeprecated)
			}

			if got.Source != tt.comment {
				t.Errorf("ParseComment() Source = %q, want %q", got.Source, tt.comment)
			}
		})
	}
}

func TestJSONParser_convertToStringSlice(t *testing.T) {
	parser := NewJSONParser()

	tests := []struct {
		name      string
		value     interface{}
		want      []string
		wantError bool
	}{
		{
			name:  "文字列配列",
			value: []interface{}{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "混合型配列",
			value: []interface{}{"test", 123, true},
			want:  []string{"test", "123", "true"},
		},
		{
			name:  "文字列スライス",
			value: []string{"x", "y", "z"},
			want:  []string{"x", "y", "z"},
		},
		{
			name:  "単一文字列",
			value: "single",
			want:  []string{"single"},
		},
		{
			name:      "変換不可能な型",
			value:     123,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.convertToStringSlice(tt.value)

			if (err != nil) != tt.wantError {
				t.Errorf("convertToStringSlice() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError {
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("convertToStringSlice() length = %d, want %d", len(got), len(tt.want))
				return
			}

			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("convertToStringSlice()[%d] = %q, want %q", i, v, tt.want[i])
				}
			}
		})
	}
}

func TestJSONParser_getDepth(t *testing.T) {
	parser := NewJSONParser()

	tests := []struct {
		name  string
		data  interface{}
		want  int
	}{
		{
			name: "プリミティブ値",
			data: "test",
			want: 1,
		},
		{
			name: "フラットなオブジェクト",
			data: map[string]interface{}{"a": "value"},
			want: 2,
		},
		{
			name: "ネストしたオブジェクト",
			data: map[string]interface{}{
				"a": map[string]interface{}{
					"b": "value",
				},
			},
			want: 3,
		},
		{
			name: "フラットな配列",
			data: []interface{}{"a", "b", "c"},
			want: 2,
		},
		{
			name: "ネストした配列",
			data: []interface{}{
				[]interface{}{"nested"},
			},
			want: 3,
		},
		{
			name: "複雑な構造",
			data: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": []interface{}{
						map[string]interface{}{
							"level3": "value",
						},
					},
				},
			},
			want: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.getDepth(tt.data)
			if got != tt.want {
				t.Errorf("getDepth() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestJSONParser_withCustomLimits(t *testing.T) {
	parser := NewJSONParserWithLimits(3, 100)

	if parser.maxDepth != 3 {
		t.Errorf("expected maxDepth 3, got %d", parser.maxDepth)
	}

	if parser.maxSize != 100 {
		t.Errorf("expected maxSize 100, got %d", parser.maxSize)
	}

	// 制限を超える深度のJSONをテスト
	deepJSON := `{"a": {"b": {"c": {"d": "too deep"}}}}`
	_, err := parser.ParseComment(deepJSON, "|")
	if err == nil {
		t.Error("expected error for too deep JSON, got nil")
	}

	// 制限を超えるサイズのJSONをテスト
	largeJSON := `{"name": "` + strings.Repeat("a", 200) + `"}`
	canParse := parser.CanParse(largeJSON)
	if canParse {
		t.Error("expected CanParse to return false for large JSON")
	}
}

func TestJSONParser_metadataExtraction(t *testing.T) {
	parser := NewJSONParser()

	comment := `{
		"name": "テストテーブル",
		"description": "テスト用のテーブル",
		"author": "開発者",
		"version": "1.0",
		"created_date": "2024-01-01"
	}`

	result, err := parser.ParseComment(comment, "|")
	if err != nil {
		t.Fatalf("ParseComment() error = %v", err)
	}

	expectedMetadata := map[string]string{
		"author":       "開発者",
		"version":      "1.0",
		"created_date": "2024-01-01",
	}

	if len(result.Metadata) != len(expectedMetadata) {
		t.Errorf("Metadata length = %d, want %d", len(result.Metadata), len(expectedMetadata))
	}

	for key, expectedValue := range expectedMetadata {
		if actualValue, exists := result.Metadata[key]; !exists {
			t.Errorf("Metadata[%s] not found", key)
		} else if actualValue != expectedValue {
			t.Errorf("Metadata[%s] = %q, want %q", key, actualValue, expectedValue)
		}
	}
}

func TestUtilityFunctions(t *testing.T) {
	// IsValidJSON のテスト
	tests := []struct {
		comment string
		want    bool
	}{
		{`{"valid": "json"}`, true},
		{`[{"array": "json"}]`, true},
		{`invalid json`, false},
		{`{"incomplete":`, false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run("IsValidJSON", func(t *testing.T) {
			got := IsValidJSON(tt.comment)
			if got != tt.want {
				t.Errorf("IsValidJSON(%q) = %v, want %v", tt.comment, got, tt.want)
			}
		})
	}

	// QuickParseJSON のテスト
	t.Run("QuickParseJSON", func(t *testing.T) {
		comment := `{"name": "クイックテスト", "description": "クイック解析テスト"}`
		result, err := QuickParseJSON(comment)
		if err != nil {
			t.Fatalf("QuickParseJSON() error = %v", err)
		}

		if result.LogicalName != "クイックテスト" {
			t.Errorf("QuickParseJSON() LogicalName = %q, want %q", result.LogicalName, "クイックテスト")
		}

		if result.Description != "クイック解析テスト" {
			t.Errorf("QuickParseJSON() Description = %q, want %q", result.Description, "クイック解析テスト")
		}
	})
}