package schema

import (
	"strings"
	"testing"
)

func TestYAMLParser_Name(t *testing.T) {
	parser := NewYAMLParser()
	if parser.Name() != "yaml" {
		t.Errorf("expected name 'yaml', got %s", parser.Name())
	}
}

func TestYAMLParser_Priority(t *testing.T) {
	parser := NewYAMLParser()
	if parser.Priority() != 15 {
		t.Errorf("expected priority 15, got %d", parser.Priority())
	}
}

func TestYAMLParser_CanParse(t *testing.T) {
	parser := NewYAMLParser()

	tests := []struct {
		name    string
		comment string
		want    bool
	}{
		{"空文字列", "", false},
		{"有効なYAMLオブジェクト", "name: テスト\ndescription: 説明", true},
		{"有効なYAMLリスト", "- name: アイテム1\n- name: アイテム2", true},
		{"単純なキー・バリュー", "title: タイトル", true},
		{"無効なYAML", "name: [unclosed", false},
		{"JSON形式", `{"name": "test"}`, false},
		{"通常のテキスト", "これは普通のコメントです", false},
		{"サイズ制限超過", strings.Repeat("name: test\n", 1000), false},
		{"空白のみ", "   \n\t  ", false},
		{"コメント行のみ", "# これはコメント", false},
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

func TestYAMLParser_looksLikeYAML(t *testing.T) {
	parser := NewYAMLParser()

	tests := []struct {
		name    string
		comment string
		want    bool
	}{
		{"キー・バリュー形式", "name: value", true},
		{"リストアイテム", "- item1", true},
		{"ネストしたキー・バリュー", "- name: value", true},
		{"空値", "name:", true},
		{"リテラルブロック", "description: |", true},
		{"フォールドブロック", "text: >", true},
		{"インデント付きキー・バリュー", "  nested: value", true},
		{"コメント行", "# comment", false},
		{"空行", "", false},
		{"普通のテキスト", "これは普通のテキスト", false},
		{"JSON形式", `{"key": "value"}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.looksLikeYAML(tt.comment)
			if got != tt.want {
				t.Errorf("looksLikeYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYAMLParser_ParseComment(t *testing.T) {
	parser := NewYAMLParser()

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
			name:        "基本的なYAMLオブジェクト",
			comment:     "name: ユーザー名\ndescription: ユーザーの表示名",
			delimiter:   "|",
			wantLogical: "ユーザー名",
			wantDesc:    "ユーザーの表示名",
			wantTags:    nil,
			wantPrio:    0,
			wantDeprecated: false,
			wantError:   false,
		},
		{
			name:        "全フィールド含むYAMLオブジェクト",
			comment:     "name: 商品名\ndescription: 商品の名前\ntags:\n  - 重要\n  - 公開\npriority: 5\ndeprecated: true",
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
			comment:     "logical_name: 論理名\ndesc: 説明文",
			delimiter:   "|",
			wantLogical: "論理名",
			wantDesc:    "説明文",
			wantTags:    nil,
			wantPrio:    0,
			wantDeprecated: false,
			wantError:   false,
		},
		{
			name:        "タグがカンマ区切り文字列の場合",
			comment:     "name: テスト\ntags: \"タグ1, タグ2, タグ3\"",
			delimiter:   "|",
			wantLogical: "テスト",
			wantDesc:    "",
			wantTags:    []string{"タグ1", "タグ2", "タグ3"},
			wantPrio:    0,
			wantDeprecated: false,
			wantError:   false,
		},
		{
			name:        "タグが単一文字列の場合",
			comment:     "name: テスト\ntags: 単一タグ",
			delimiter:   "|",
			wantLogical: "テスト",
			wantDesc:    "",
			wantTags:    []string{"単一タグ"},
			wantPrio:    0,
			wantDeprecated: false,
			wantError:   false,
		},
		{
			name:        "YAML配列形式",
			comment:     "- name: 配列テスト\n  description: 配列の最初の要素",
			delimiter:   "|",
			wantLogical: "配列テスト",
			wantDesc:    "配列の最初の要素",
			wantTags:    nil,
			wantPrio:    0,
			wantDeprecated: false,
			wantError:   false,
		},
		{
			name:        "単純な文字列（YAMLパーサーでは処理不可）",
			comment:     "これは説明文です",
			delimiter:   "|",
			wantError:   true,
		},
		{
			name:        "優先度が文字列の場合",
			comment:     "name: テスト\npriority: \"10\"",
			delimiter:   "|",
			wantLogical: "テスト",
			wantDesc:    "",
			wantTags:    nil,
			wantPrio:    10,
			wantDeprecated: false,
			wantError:   false,
		},
		{
			name:        "非推奨フラグが文字列の場合",
			comment:     "name: テスト\ndeprecated: \"true\"",
			delimiter:   "|",
			wantLogical: "テスト",
			wantDesc:    "",
			wantTags:    nil,
			wantPrio:    0,
			wantDeprecated: true,
			wantError:   false,
		},
		{
			name:      "無効なYAML",
			comment:   "name: [unclosed",
			delimiter: "|",
			wantError: true,
		},
		{
			name:      "サイズ制限超過",
			comment:   strings.Repeat("name: test\n", 1000),
			delimiter: "|",
			wantError: true,
		},
		{
			name:      "深すぎるYAML構造",
			comment:   "a:\n  b:\n    c:\n      d:\n        e:\n          f: too deep",
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

func TestYAMLParser_convertToStringSlice(t *testing.T) {
	parser := NewYAMLParser()

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
			name:  "カンマ区切り文字列",
			value: "a, b, c",
			want:  []string{"a", "b", "c"},
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

func TestYAMLParser_getDepth(t *testing.T) {
	parser := NewYAMLParser()

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

func TestYAMLParser_withCustomLimits(t *testing.T) {
	parser := NewYAMLParserWithLimits(3, 100)

	if parser.maxDepth != 3 {
		t.Errorf("expected maxDepth 3, got %d", parser.maxDepth)
	}

	if parser.maxSize != 100 {
		t.Errorf("expected maxSize 100, got %d", parser.maxSize)
	}

	// 制限を超える深度のYAMLをテスト
	deepYAML := "a:\n  b:\n    c:\n      d: too deep"
	_, err := parser.ParseComment(deepYAML, "|")
	if err == nil {
		t.Error("expected error for too deep YAML, got nil")
	}

	// 制限を超えるサイズのYAMLをテスト
	largeYAML := strings.Repeat("name: test\n", 20)
	canParse := parser.CanParse(largeYAML)
	if canParse {
		t.Error("expected CanParse to return false for large YAML")
	}
}

func TestYAMLParser_metadataExtraction(t *testing.T) {
	parser := NewYAMLParser()

	comment := `name: テストテーブル
description: テスト用のテーブル
author: 開発者
version: "1.0"
created_date: "2024-01-01"`

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

func TestYAMLParser_literalAndFoldedBlocks(t *testing.T) {
	parser := NewYAMLParser()

	tests := []struct {
		name    string
		comment string
		wantDesc string
	}{
		{
			name: "リテラルブロック",
			comment: `name: テスト
description: |
  これは複数行の
  説明文です
  改行が保持されます`,
			wantDesc: "これは複数行の\n説明文です\n改行が保持されます",
		},
		{
			name: "フォールドブロック",
			comment: `name: テスト
description: >
  これは複数行の
  説明文ですが
  改行は空白になります`,
			wantDesc: "これは複数行の 説明文ですが 改行は空白になります",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.ParseComment(tt.comment, "|")
			if err != nil {
				t.Fatalf("ParseComment() error = %v", err)
			}

			if result.Description != tt.wantDesc {
				t.Errorf("Description = %q, want %q", result.Description, tt.wantDesc)
			}
		})
	}
}

func TestParseIntFromString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      int
		wantError bool
	}{
		{"正の整数", "123", 123, false},
		{"負の整数", "-456", -456, false},
		{"前後の空白", "  789  ", 789, false},
		{"ゼロ", "0", 0, false},
		{"空文字列", "", 0, true},
		{"無効な文字列", "abc", 0, true},
		{"浮動小数点", "12.34", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIntFromString(tt.input)

			if (err != nil) != tt.wantError {
				t.Errorf("parseIntFromString() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError && got != tt.want {
				t.Errorf("parseIntFromString() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestYAMLUtilityFunctions(t *testing.T) {
	// IsValidYAML のテスト
	tests := []struct {
		comment string
		want    bool
	}{
		{"name: valid", true},
		{"- item", true},
		{"invalid: [unclosed", false},
		{"plain text", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run("IsValidYAML", func(t *testing.T) {
			got := IsValidYAML(tt.comment)
			if got != tt.want {
				t.Errorf("IsValidYAML(%q) = %v, want %v", tt.comment, got, tt.want)
			}
		})
	}

	// QuickParseYAML のテスト
	t.Run("QuickParseYAML", func(t *testing.T) {
		comment := "name: クイックテスト\ndescription: クイック解析テスト"
		result, err := QuickParseYAML(comment)
		if err != nil {
			t.Fatalf("QuickParseYAML() error = %v", err)
		}

		if result.LogicalName != "クイックテスト" {
			t.Errorf("QuickParseYAML() LogicalName = %q, want %q", result.LogicalName, "クイックテスト")
		}

		if result.Description != "クイック解析テスト" {
			t.Errorf("QuickParseYAML() Description = %q, want %q", result.Description, "クイック解析テスト")
		}
	})

	// SupportedYAMLFormats のテスト
	t.Run("SupportedYAMLFormats", func(t *testing.T) {
		formats := SupportedYAMLFormats()
		if len(formats) == 0 {
			t.Error("SupportedYAMLFormats() should return non-empty slice")
		}

		// 各フォーマットが有効なYAMLかチェック
		for i, format := range formats {
			if !IsValidYAML(format) {
				t.Errorf("SupportedYAMLFormats()[%d] is not valid YAML: %q", i, format)
			}
		}
	})
}