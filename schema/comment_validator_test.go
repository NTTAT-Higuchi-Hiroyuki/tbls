package schema

import (
	"strings"
	"testing"
)

func TestDefaultValidationConfig(t *testing.T) {
	config := DefaultValidationConfig()
	
	if config.MaxLogicalNameLength != 100 {
		t.Errorf("expected MaxLogicalNameLength 100, got %d", config.MaxLogicalNameLength)
	}
	
	if config.MaxDescriptionLength != 1000 {
		t.Errorf("expected MaxDescriptionLength 1000, got %d", config.MaxDescriptionLength)
	}
	
	if !config.EnableHTMLEscape {
		t.Error("expected EnableHTMLEscape to be true")
	}
	
	if !config.EnableSQLInjectionCheck {
		t.Error("expected EnableSQLInjectionCheck to be true")
	}
}

func TestNewDefaultCommentValidator(t *testing.T) {
	validator := NewDefaultCommentValidator()
	
	if validator.config == nil {
		t.Error("expected config to be set")
	}
	
	if validator.allowedCharRegex == nil {
		t.Error("expected allowedCharRegex to be compiled")
	}
	
	if validator.sqlInjectionRegex == nil {
		t.Error("expected sqlInjectionRegex to be compiled")
	}
}

func TestDefaultCommentValidator_Validate(t *testing.T) {
	validator := NewDefaultCommentValidator()
	
	tests := []struct {
		name      string
		data      *CommentData
		wantError bool
	}{
		{
			name:      "nilデータ",
			data:      nil,
			wantError: true,
		},
		{
			name: "有効なデータ",
			data: &CommentData{
				LogicalName: "ユーザー名",
				Description: "ユーザーの表示名を格納します",
				Tags:        []string{"重要", "公開"},
				Metadata:    map[string]string{"author": "開発者"},
			},
			wantError: false,
		},
		{
			name: "空のデータ",
			data: &CommentData{},
			wantError: false,
		},
		{
			name: "論理名が長すぎる",
			data: &CommentData{
				LogicalName: strings.Repeat("あ", 101),
			},
			wantError: true,
		},
		{
			name: "説明が長すぎる",
			data: &CommentData{
				Description: strings.Repeat("あ", 1001),
			},
			wantError: true,
		},
		{
			name: "タグが多すぎる",
			data: &CommentData{
				Tags: make([]string, 21),
			},
			wantError: true,
		},
		{
			name: "SQLインジェクション検出",
			data: &CommentData{
				LogicalName: "test'; DROP TABLE users; --",
			},
			wantError: true,
		},
		{
			name: "禁止語句検出",
			data: &CommentData{
				Description: "この説明にはDROPコマンドが含まれています",
			},
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestDefaultCommentValidator_Sanitize(t *testing.T) {
	validator := NewDefaultCommentValidator()
	
	tests := []struct {
		name     string
		input    *CommentData
		expected *CommentData
	}{
		{
			name:     "nilデータ",
			input:    nil,
			expected: nil,
		},
		{
			name: "前後の空白除去",
			input: &CommentData{
				LogicalName: "  ユーザー名  ",
				Description: "  説明文  ",
			},
			expected: &CommentData{
				LogicalName: "ユーザー名",
				Description: "説明文",
			},
		},
		{
			name: "HTMLエスケープ",
			input: &CommentData{
				LogicalName: "<script>alert('test')</script>",
				Description: "A & B",
			},
			expected: &CommentData{
				LogicalName: "&lt;script&gt;alert(&#39;test&#39;)&lt;/script&gt;",
				Description: "A &amp; B",
			},
		},
		{
			name: "連続空白の正規化",
			input: &CommentData{
				LogicalName: "テスト   名前",
				Description: "複数    空白    あり",
			},
			expected: &CommentData{
				LogicalName: "テスト 名前",
				Description: "複数 空白 あり",
			},
		},
		{
			name: "空タグの除去",
			input: &CommentData{
				Tags: []string{"有効", "", "タグ", "  "},
			},
			expected: &CommentData{
				Tags: []string{"有効", "タグ"},
			},
		},
		{
			name: "空メタデータの除去",
			input: &CommentData{
				Metadata: map[string]string{
					"valid_key":   "valid_value",
					"":           "empty_key",
					"empty_value": "",
					"  spaces  ":  "  spaces  ",
				},
			},
			expected: &CommentData{
				Metadata: map[string]string{
					"valid_key": "valid_value",
					"spaces":    "spaces",
				},
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.Sanitize(tt.input)
			
			if tt.expected == nil {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}
			
			if got == nil {
				t.Errorf("expected non-nil result")
				return
			}
			
			if got.LogicalName != tt.expected.LogicalName {
				t.Errorf("LogicalName = %q, want %q", got.LogicalName, tt.expected.LogicalName)
			}
			
			if got.Description != tt.expected.Description {
				t.Errorf("Description = %q, want %q", got.Description, tt.expected.Description)
			}
			
			if len(got.Tags) != len(tt.expected.Tags) {
				t.Errorf("Tags length = %d, want %d", len(got.Tags), len(tt.expected.Tags))
			} else {
				for i, tag := range got.Tags {
					if tag != tt.expected.Tags[i] {
						t.Errorf("Tags[%d] = %q, want %q", i, tag, tt.expected.Tags[i])
					}
				}
			}
			
			if len(got.Metadata) != len(tt.expected.Metadata) {
				t.Errorf("Metadata length = %d, want %d", len(got.Metadata), len(tt.expected.Metadata))
			} else {
				for key, value := range tt.expected.Metadata {
					if gotValue, exists := got.Metadata[key]; !exists {
						t.Errorf("Metadata[%q] not found", key)
					} else if gotValue != value {
						t.Errorf("Metadata[%q] = %q, want %q", key, gotValue, value)
					}
				}
			}
		})
	}
}

func TestDefaultCommentValidator_ValidateAndSanitize(t *testing.T) {
	validator := NewDefaultCommentValidator()
	
	tests := []struct {
		name      string
		input     *CommentData
		wantError bool
	}{
		{
			name: "有効なデータのサニタイゼーションと検証",
			input: &CommentData{
				LogicalName: "  ユーザー名  ",
				Description: "  説明文  ",
				Tags:        []string{"タグ1", "", "タグ2"},
			},
			wantError: false,
		},
		{
			name: "サニタイゼーション後に無効になるデータ",
			input: &CommentData{
				LogicalName: strings.Repeat("a", 101), // 長すぎる
			},
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validator.ValidateAndSanitize(tt.input)
			
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateAndSanitize() error = %v, wantError %v", err, tt.wantError)
				return
			}
			
			if !tt.wantError && got == nil {
				t.Error("expected non-nil result for valid input")
			}
		})
	}
}

func TestDefaultCommentValidator_validateLogicalName(t *testing.T) {
	validator := NewDefaultCommentValidator()
	
	tests := []struct {
		name        string
		logicalName string
		wantError   bool
	}{
		{"空文字列", "", false},
		{"通常の論理名", "ユーザー名", false},
		{"英数字", "User Name 123", false},
		{"長すぎる論理名", strings.Repeat("a", 101), true},
		{"SQLインジェクション", "test'; DROP TABLE users; --", true},
		{"禁止語句", "DROP something", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateLogicalName(tt.logicalName)
			if (err != nil) != tt.wantError {
				t.Errorf("validateLogicalName() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestDefaultCommentValidator_validateTags(t *testing.T) {
	validator := NewDefaultCommentValidator()
	
	tests := []struct {
		name      string
		tags      []string
		wantError bool
	}{
		{"空スライス", []string{}, false},
		{"nilスライス", nil, false},
		{"有効なタグ", []string{"タグ1", "タグ2"}, false},
		{"空タグ", []string{"有効", ""}, true},
		{"多すぎるタグ", make([]string, 21), true},
		{"長すぎるタグ", []string{strings.Repeat("a", 51)}, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 多すぎるタグのテストケースで有効な値を設定
			if tt.name == "多すぎるタグ" {
				for i := range tt.tags {
					tt.tags[i] = "タグ"
				}
			}
			
			err := validator.validateTags(tt.tags)
			if (err != nil) != tt.wantError {
				t.Errorf("validateTags() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestDefaultCommentValidator_validateMetadata(t *testing.T) {
	validator := NewDefaultCommentValidator()
	
	tests := []struct {
		name      string
		metadata  map[string]string
		wantError bool
	}{
		{"空マップ", map[string]string{}, false},
		{"nilマップ", nil, false},
		{"有効なメタデータ", map[string]string{"key": "value"}, false},
		{"空キー", map[string]string{"": "value"}, true},
		{"長すぎるキー", map[string]string{strings.Repeat("a", 101): "value"}, true},
		{"長すぎる値", map[string]string{"key": strings.Repeat("a", 501)}, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateMetadata(tt.metadata)
			if (err != nil) != tt.wantError {
				t.Errorf("validateMetadata() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestDefaultCommentValidator_sanitizeString(t *testing.T) {
	validator := NewDefaultCommentValidator()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"空文字列", "", ""},
		{"前後の空白", "  テスト  ", "テスト"},
		{"連続空白", "テスト   文字列", "テスト 文字列"},
		{"HTMLエスケープ", "<script>", "&lt;script&gt;"},
		{"制御文字除去", "テスト\x00文字列", "テスト文字列"},
		{"改行とタブ", "行1\n行2\t列", "行1 行2 列"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.sanitizeString(tt.input)
			if got != tt.expected {
				t.Errorf("sanitizeString() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestCustomValidationConfig(t *testing.T) {
	config := &ValidationConfig{
		MaxLogicalNameLength: 10,
		MaxDescriptionLength: 20,
		ForbiddenWords:      []string{"禁止"},
		EnableHTMLEscape:    false,
	}
	
	validator := NewDefaultCommentValidatorWithConfig(config)
	
	// 設定が反映されているかテスト
	data := &CommentData{
		LogicalName: strings.Repeat("a", 11), // 制限を超える
	}
	
	err := validator.Validate(data)
	if err == nil {
		t.Error("expected validation error for long logical name")
	}
	
	// HTMLエスケープが無効になっているかテスト
	data2 := &CommentData{
		LogicalName: "<test>",
	}
	
	sanitized := validator.Sanitize(data2)
	if sanitized.LogicalName != "<test>" {
		t.Errorf("expected HTML not to be escaped, got %q", sanitized.LogicalName)
	}
}

func TestStrictValidationConfig(t *testing.T) {
	config := StrictValidationConfig()
	validator := NewDefaultCommentValidatorWithConfig(config)
	
	// より厳格な制限をテスト
	data := &CommentData{
		LogicalName: strings.Repeat("a", 51), // strict設定では50文字制限
	}
	
	err := validator.Validate(data)
	if err == nil {
		t.Error("expected validation error for strict config")
	}
}

func TestPermissiveValidationConfig(t *testing.T) {
	config := PermissiveValidationConfig()
	validator := NewDefaultCommentValidatorWithConfig(config)
	
	// 緩い制限をテスト
	data := &CommentData{
		LogicalName: "test'; DROP TABLE users; --", // SQL注入
	}
	
	err := validator.Validate(data)
	if err != nil {
		t.Errorf("expected no validation error for permissive config, got %v", err)
	}
}

func TestValidatorConfigManagement(t *testing.T) {
	validator := NewDefaultCommentValidator()
	
	// 初期設定の確認
	initialConfig := validator.GetConfig()
	if initialConfig.MaxLogicalNameLength != 100 {
		t.Errorf("expected initial MaxLogicalNameLength 100, got %d", initialConfig.MaxLogicalNameLength)
	}
	
	// 新しい設定を適用
	newConfig := &ValidationConfig{
		MaxLogicalNameLength: 50,
		EnableHTMLEscape:    false,
	}
	
	validator.SetConfig(newConfig)
	
	// 設定が更新されているか確認
	updatedConfig := validator.GetConfig()
	if updatedConfig.MaxLogicalNameLength != 50 {
		t.Errorf("expected updated MaxLogicalNameLength 50, got %d", updatedConfig.MaxLogicalNameLength)
	}
	
	if updatedConfig.EnableHTMLEscape {
		t.Error("expected EnableHTMLEscape to be false")
	}
}

func TestSQLInjectionDetection(t *testing.T) {
	validator := NewDefaultCommentValidator()
	
	sqlInjectionPatterns := []string{
		"'; DROP TABLE users; --",
		"' OR 1=1 --",
		"UNION SELECT * FROM passwords",
		"'; EXEC xp_cmdshell('dir'); --",
		"/* comment */ DROP TABLE test",
	}
	
	for _, pattern := range sqlInjectionPatterns {
		t.Run("SQL injection: "+pattern, func(t *testing.T) {
			data := &CommentData{
				LogicalName: pattern,
			}
			
			err := validator.Validate(data)
			if err == nil {
				t.Errorf("expected SQL injection to be detected in: %s", pattern)
			}
		})
	}
}

func TestControlCharacterRemoval(t *testing.T) {
	validator := NewDefaultCommentValidator()
	
	// 制御文字を含む文字列
	input := "テスト\x00\x01\x02文字列"
	expected := "テスト文字列"
	
	result := validator.removeControlChars(input)
	if result != expected {
		t.Errorf("removeControlChars() = %q, want %q", result, expected)
	}
	
	// 許可される制御文字（タブ、改行など）
	input2 := "行1\n行2\t列"
	result2 := validator.removeControlChars(input2)
	if result2 != input2 {
		t.Errorf("removeControlChars() should preserve tabs and newlines")
	}
}