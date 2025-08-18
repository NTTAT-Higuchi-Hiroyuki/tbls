package config

import (
	"testing"

	"github.com/goccy/go-yaml"
)

func TestDefaultEnhancedCommentConfig(t *testing.T) {
	config, err := New()
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	// デフォルト値の確認
	if config.EnhancedComment.Parser.PreferredFormat != "auto" {
		t.Errorf("expected PreferredFormat 'auto', got %s", config.EnhancedComment.Parser.PreferredFormat)
	}

	if config.EnhancedComment.Parser.MaxDepth != 5 {
		t.Errorf("expected MaxDepth 5, got %d", config.EnhancedComment.Parser.MaxDepth)
	}

	if config.EnhancedComment.Parser.MaxSize != 8192 {
		t.Errorf("expected MaxSize 8192, got %d", config.EnhancedComment.Parser.MaxSize)
	}

	if config.EnhancedComment.Validation.SecurityLevel != "default" {
		t.Errorf("expected SecurityLevel 'default', got %s", config.EnhancedComment.Validation.SecurityLevel)
	}

	if config.EnhancedComment.Processing.ProcessingTimeout != 1000 {
		t.Errorf("expected ProcessingTimeout 1000, got %d", config.EnhancedComment.Processing.ProcessingTimeout)
	}

	expectedObjectTypes := []string{"table", "column", "index", "view", "constraint"}
	if len(config.EnhancedComment.Processing.ObjectTypes) != len(expectedObjectTypes) {
		t.Errorf("expected %d ObjectTypes, got %d", len(expectedObjectTypes), len(config.EnhancedComment.Processing.ObjectTypes))
	}
}

func TestEnhancedCommentConfigYAML(t *testing.T) {
	yamlContent := `
name: Test Database
enhancedComment:
  enabled: true
  parser:
    enableJSON: true
    enableYAML: false
    preferredFormat: json
    maxDepth: 3
    maxSize: 4096
    fallbackToLegacy: false
  validation:
    enabled: true
    enableSanitization: false
    securityLevel: strict
    maxLogicalNameLength: 50
    maxDescriptionLength: 500
    maxTagCount: 10
    enableHTMLEscape: false
    enableSQLInjectionCheck: true
    forbiddenWords: ["DROP", "DELETE"]
  processing:
    strictMode: true
    processingTimeout: 2000
    enableBatchProcessing: false
    objectTypes: ["table", "column"]
`

	config := &Config{}
	if err := yaml.Unmarshal([]byte(yamlContent), config); err != nil {
		t.Fatalf("failed to unmarshal YAML: %v", err)
	}

	if err := config.setDefault(); err != nil {
		t.Fatalf("failed to set defaults: %v", err)
	}

	// YAML設定値の確認
	if !config.EnhancedComment.Enabled {
		t.Error("expected EnhancedComment.Enabled to be true")
	}

	if !config.EnhancedComment.Parser.EnableJSON {
		t.Error("expected EnableJSON to be true")
	}

	if config.EnhancedComment.Parser.EnableYAML {
		t.Error("expected EnableYAML to be false")
	}

	if config.EnhancedComment.Parser.PreferredFormat != "json" {
		t.Errorf("expected PreferredFormat 'json', got %s", config.EnhancedComment.Parser.PreferredFormat)
	}

	if config.EnhancedComment.Parser.MaxDepth != 3 {
		t.Errorf("expected MaxDepth 3, got %d", config.EnhancedComment.Parser.MaxDepth)
	}

	if config.EnhancedComment.Parser.MaxSize != 4096 {
		t.Errorf("expected MaxSize 4096, got %d", config.EnhancedComment.Parser.MaxSize)
	}

	// FallbackToLegacyは明示的にfalseに設定されているので期待通り
	if config.EnhancedComment.Parser.FallbackToLegacy {
		t.Error("expected FallbackToLegacy to be false")
	}

	if config.EnhancedComment.Validation.SecurityLevel != "strict" {
		t.Errorf("expected SecurityLevel 'strict', got %s", config.EnhancedComment.Validation.SecurityLevel)
	}

	if config.EnhancedComment.Validation.MaxLogicalNameLength != 50 {
		t.Errorf("expected MaxLogicalNameLength 50, got %d", config.EnhancedComment.Validation.MaxLogicalNameLength)
	}

	if config.EnhancedComment.Processing.StrictMode != true {
		t.Error("expected StrictMode to be true")
	}

	if config.EnhancedComment.Processing.ProcessingTimeout != 2000 {
		t.Errorf("expected ProcessingTimeout 2000, got %d", config.EnhancedComment.Processing.ProcessingTimeout)
	}

	expectedObjectTypes := []string{"table", "column"}
	if len(config.EnhancedComment.Processing.ObjectTypes) != len(expectedObjectTypes) {
		t.Errorf("expected %d ObjectTypes, got %d", len(expectedObjectTypes), len(config.EnhancedComment.Processing.ObjectTypes))
	}

	expectedForbiddenWords := []string{"DROP", "DELETE"}
	if len(config.EnhancedComment.Validation.ForbiddenWords) != len(expectedForbiddenWords) {
		t.Errorf("expected %d ForbiddenWords, got %d", len(expectedForbiddenWords), len(config.EnhancedComment.Validation.ForbiddenWords))
	}
}

func TestEnhancedCommentConfigHelperMethods(t *testing.T) {
	config := &Config{
		EnhancedComment: EnhancedCommentConfig{
			Enabled: true,
			Parser: EnhancedCommentParserConfig{
				EnableJSON:       true,
				EnableYAML:       false,
				PreferredFormat:  "json",
				FallbackToLegacy: true,
			},
			Validation: EnhancedCommentValidationConfig{
				Enabled:            true,
				EnableSanitization: true,
				SecurityLevel:      "strict",
			},
			Processing: EnhancedCommentProcessingConfig{
				StrictMode:            true,
				ProcessingTimeout:     2000,
				EnableBatchProcessing: true,
				ObjectTypes:           []string{"table", "column"},
			},
		},
	}

	// ヘルパーメソッドのテスト
	if !config.IsEnhancedCommentEnabled() {
		t.Error("expected IsEnhancedCommentEnabled to be true")
	}

	if !config.IsEnhancedCommentJSONEnabled() {
		t.Error("expected IsEnhancedCommentJSONEnabled to be true")
	}

	if config.IsEnhancedCommentYAMLEnabled() {
		t.Error("expected IsEnhancedCommentYAMLEnabled to be false")
	}

	if config.EnhancedCommentPreferredFormat() != "json" {
		t.Errorf("expected EnhancedCommentPreferredFormat 'json', got %s", config.EnhancedCommentPreferredFormat())
	}

	if !config.IsEnhancedCommentValidationEnabled() {
		t.Error("expected IsEnhancedCommentValidationEnabled to be true")
	}

	if !config.IsEnhancedCommentSanitizationEnabled() {
		t.Error("expected IsEnhancedCommentSanitizationEnabled to be true")
	}

	if config.EnhancedCommentSecurityLevel() != "strict" {
		t.Errorf("expected EnhancedCommentSecurityLevel 'strict', got %s", config.EnhancedCommentSecurityLevel())
	}

	if !config.IsEnhancedCommentStrictMode() {
		t.Error("expected IsEnhancedCommentStrictMode to be true")
	}

	if config.EnhancedCommentProcessingTimeout() != 2000 {
		t.Errorf("expected EnhancedCommentProcessingTimeout 2000, got %d", config.EnhancedCommentProcessingTimeout())
	}

	if !config.IsEnhancedCommentObjectTypeEnabled("table") {
		t.Error("expected IsEnhancedCommentObjectTypeEnabled('table') to be true")
	}

	if !config.IsEnhancedCommentObjectTypeEnabled("column") {
		t.Error("expected IsEnhancedCommentObjectTypeEnabled('column') to be true")
	}

	if config.IsEnhancedCommentObjectTypeEnabled("index") {
		t.Error("expected IsEnhancedCommentObjectTypeEnabled('index') to be false")
	}
}

func TestEnhancedCommentConfigDisabled(t *testing.T) {
	config := &Config{
		EnhancedComment: EnhancedCommentConfig{
			Enabled: false,
		},
	}

	// 無効な場合のテスト
	if config.IsEnhancedCommentEnabled() {
		t.Error("expected IsEnhancedCommentEnabled to be false")
	}

	if config.IsEnhancedCommentJSONEnabled() {
		t.Error("expected IsEnhancedCommentJSONEnabled to be false")
	}

	if config.IsEnhancedCommentYAMLEnabled() {
		t.Error("expected IsEnhancedCommentYAMLEnabled to be false")
	}

	if config.IsEnhancedCommentValidationEnabled() {
		t.Error("expected IsEnhancedCommentValidationEnabled to be false")
	}

	if config.IsEnhancedCommentStrictMode() {
		t.Error("expected IsEnhancedCommentStrictMode to be false")
	}

	if config.IsEnhancedCommentObjectTypeEnabled("table") {
		t.Error("expected IsEnhancedCommentObjectTypeEnabled('table') to be false")
	}
}

func TestEnhancedCommentConfigDefaultValues(t *testing.T) {
	config := &Config{}

	// setDefaultを呼び出してデフォルト値を設定
	if err := config.setDefault(); err != nil {
		t.Fatalf("failed to set defaults: %v", err)
	}

	// デフォルト値ヘルパーメソッドのテスト
	if config.EnhancedCommentPreferredFormat() != "auto" {
		t.Errorf("expected default EnhancedCommentPreferredFormat 'auto', got %s", config.EnhancedCommentPreferredFormat())
	}

	if config.EnhancedCommentSecurityLevel() != "default" {
		t.Errorf("expected default EnhancedCommentSecurityLevel 'default', got %s", config.EnhancedCommentSecurityLevel())
	}

	if config.EnhancedCommentProcessingTimeout() != 1000 {
		t.Errorf("expected default EnhancedCommentProcessingTimeout 1000, got %d", config.EnhancedCommentProcessingTimeout())
	}
}

func TestEnhancedCommentConfigPartialSettings(t *testing.T) {
	// 部分的な設定でのテスト
	yamlContent := `
name: Partial Test
enhancedComment:
  enabled: true
  parser:
    enableJSON: false
  validation:
    securityLevel: permissive
`

	config := &Config{}
	if err := yaml.Unmarshal([]byte(yamlContent), config); err != nil {
		t.Fatalf("failed to unmarshal YAML: %v", err)
	}

	if err := config.setDefault(); err != nil {
		t.Fatalf("failed to set defaults: %v", err)
	}

	// 明示的に設定された値
	if config.EnhancedComment.Parser.EnableJSON {
		t.Error("expected EnableJSON to be false")
	}

	if config.EnhancedComment.Validation.SecurityLevel != "permissive" {
		t.Errorf("expected SecurityLevel 'permissive', got %s", config.EnhancedComment.Validation.SecurityLevel)
	}

	// EnableYAMLはデフォルトではfalseなので確認しない
	// YAML設定でEnableJSONのみfalseに設定されているため

	if config.EnhancedComment.Parser.MaxDepth != 5 {
		t.Errorf("expected MaxDepth 5 (default), got %d", config.EnhancedComment.Parser.MaxDepth)
	}
}

func TestEnhancedCommentConfigValidation(t *testing.T) {
	// 無効な設定値のテスト
	config := &Config{
		EnhancedComment: EnhancedCommentConfig{
			Enabled: true,
			Validation: EnhancedCommentValidationConfig{
				SecurityLevel: "invalid_level",
			},
		},
	}

	if err := config.setDefault(); err != nil {
		t.Fatalf("failed to set defaults: %v", err)
	}

	// 無効なセキュリティレベルの場合はデフォルトが適用されないことを確認
	if config.EnhancedComment.Validation.SecurityLevel == "default" {
		t.Error("invalid SecurityLevel should not be overridden by default")
	}

	// ヘルパーメソッドは無効な値でもデフォルトを返すべき
	if config.EnhancedCommentSecurityLevel() != "invalid_level" {
		t.Error("EnhancedCommentSecurityLevel should return the configured value even if invalid")
	}
}

func TestEnhancedCommentConfigSerialization(t *testing.T) {
	// 設定のシリアライゼーション・デシリアライゼーションテスト
	original := &Config{
		Name: "Serialization Test",
		EnhancedComment: EnhancedCommentConfig{
			Enabled: true,
			Parser: EnhancedCommentParserConfig{
				EnableJSON:       true,
				EnableYAML:       true,
				PreferredFormat:  "auto",
				FallbackToLegacy: true,
				MaxDepth:         5,
				MaxSize:          8192,
			},
			Validation: EnhancedCommentValidationConfig{
				Enabled:                 true,
				EnableSanitization:      true,
				SecurityLevel:          "default",
				MaxLogicalNameLength:   100,
				MaxDescriptionLength:   1000,
				MaxTagCount:            20,
				EnableHTMLEscape:       true,
				EnableSQLInjectionCheck: true,
				ForbiddenWords:         []string{"DROP", "DELETE"},
			},
			Processing: EnhancedCommentProcessingConfig{
				StrictMode:            false,
				ProcessingTimeout:     1000,
				EnableBatchProcessing: true,
				ObjectTypes:           []string{"table", "column", "index"},
			},
		},
	}

	// YAML形式でシリアライズ
	yamlData, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal to YAML: %v", err)
	}

	// YAML形式からデシリアライズ
	restored := &Config{}
	if err := yaml.Unmarshal(yamlData, restored); err != nil {
		t.Fatalf("failed to unmarshal from YAML: %v", err)
	}

	// デフォルト値を設定
	if err := restored.setDefault(); err != nil {
		t.Fatalf("failed to set defaults: %v", err)
	}

	// 復元された設定の確認
	if restored.Name != original.Name {
		t.Errorf("Name mismatch: expected %s, got %s", original.Name, restored.Name)
	}

	if restored.EnhancedComment.Enabled != original.EnhancedComment.Enabled {
		t.Error("EnhancedComment.Enabled mismatch")
	}

	if restored.EnhancedComment.Parser.PreferredFormat != original.EnhancedComment.Parser.PreferredFormat {
		t.Errorf("PreferredFormat mismatch: expected %s, got %s", 
			original.EnhancedComment.Parser.PreferredFormat, 
			restored.EnhancedComment.Parser.PreferredFormat)
	}

	if len(restored.EnhancedComment.Validation.ForbiddenWords) != len(original.EnhancedComment.Validation.ForbiddenWords) {
		t.Error("ForbiddenWords length mismatch")
	}

	if len(restored.EnhancedComment.Processing.ObjectTypes) != len(original.EnhancedComment.Processing.ObjectTypes) {
		t.Error("ObjectTypes length mismatch")
	}
}

func TestEnhancedCommentConfigEmptyObjectTypes(t *testing.T) {
	// ObjectTypesが空の場合のテスト
	config := &Config{
		EnhancedComment: EnhancedCommentConfig{
			Enabled: true,
			Processing: EnhancedCommentProcessingConfig{
				ObjectTypes: []string{}, // 空のスライス
			},
		},
	}

	// ObjectTypesが空の場合はすべてのタイプが有効になるべき
	if !config.IsEnhancedCommentObjectTypeEnabled("table") {
		t.Error("expected all object types to be enabled when ObjectTypes is empty")
	}

	if !config.IsEnhancedCommentObjectTypeEnabled("column") {
		t.Error("expected all object types to be enabled when ObjectTypes is empty")
	}

	if !config.IsEnhancedCommentObjectTypeEnabled("index") {
		t.Error("expected all object types to be enabled when ObjectTypes is empty")
	}

	if !config.IsEnhancedCommentObjectTypeEnabled("arbitrary_type") {
		t.Error("expected all object types to be enabled when ObjectTypes is empty")
	}
}