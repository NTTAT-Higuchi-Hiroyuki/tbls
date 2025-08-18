package schema

import (
	"testing"
)

// MockEnhancedCommentConfigurator テスト用の設定モック
type MockEnhancedCommentConfigurator struct {
	enabled                       bool
	jsonEnabled                   bool
	yamlEnabled                   bool
	preferredFormat              string
	validationEnabled            bool
	sanitizationEnabled          bool
	securityLevel                string
	strictMode                   bool
	processingTimeout            int
	objectTypeEnabledMap         map[string]bool
	logicalNameDelimiter         string
}

func NewMockEnhancedCommentConfigurator() *MockEnhancedCommentConfigurator {
	return &MockEnhancedCommentConfigurator{
		enabled:                      true,
		jsonEnabled:                  true,
		yamlEnabled:                  true,
		preferredFormat:              "auto",
		validationEnabled:            true,
		sanitizationEnabled:          true,
		securityLevel:                "default",
		strictMode:                   false,
		processingTimeout:            1000,
		objectTypeEnabledMap:         make(map[string]bool),
		logicalNameDelimiter:         "|",
	}
}

func (m *MockEnhancedCommentConfigurator) IsEnhancedCommentEnabled() bool {
	return m.enabled
}

func (m *MockEnhancedCommentConfigurator) IsEnhancedCommentJSONEnabled() bool {
	return m.enabled && m.jsonEnabled
}

func (m *MockEnhancedCommentConfigurator) IsEnhancedCommentYAMLEnabled() bool {
	return m.enabled && m.yamlEnabled
}

func (m *MockEnhancedCommentConfigurator) EnhancedCommentPreferredFormat() string {
	return m.preferredFormat
}

func (m *MockEnhancedCommentConfigurator) IsEnhancedCommentValidationEnabled() bool {
	return m.enabled && m.validationEnabled
}

func (m *MockEnhancedCommentConfigurator) IsEnhancedCommentSanitizationEnabled() bool {
	return m.enabled && m.sanitizationEnabled
}

func (m *MockEnhancedCommentConfigurator) EnhancedCommentSecurityLevel() string {
	return m.securityLevel
}

func (m *MockEnhancedCommentConfigurator) IsEnhancedCommentStrictMode() bool {
	return m.enabled && m.strictMode
}

func (m *MockEnhancedCommentConfigurator) EnhancedCommentProcessingTimeout() int {
	return m.processingTimeout
}

func (m *MockEnhancedCommentConfigurator) IsEnhancedCommentObjectTypeEnabled(objectType string) bool {
	if enabled, exists := m.objectTypeEnabledMap[objectType]; exists {
		return m.enabled && enabled
	}
	return m.enabled // デフォルトでは有効
}

func (m *MockEnhancedCommentConfigurator) LogicalNameDelimiter() string {
	return m.logicalNameDelimiter
}

func TestNewEnhancedCommentProcessorFromConfig(t *testing.T) {
	config := NewMockEnhancedCommentConfigurator()

	processor := NewEnhancedCommentProcessorFromConfig(config)
	if processor == nil {
		t.Fatal("processor should not be nil")
	}

	// 基本設定の確認
	processingConfig := processor.GetConfig()
	if processingConfig.EnableValidation != config.IsEnhancedCommentValidationEnabled() {
		t.Errorf("expected EnableValidation %v, got %v", config.IsEnhancedCommentValidationEnabled(), processingConfig.EnableValidation)
	}

	if processingConfig.EnableSanitization != config.IsEnhancedCommentSanitizationEnabled() {
		t.Errorf("expected EnableSanitization %v, got %v", config.IsEnhancedCommentSanitizationEnabled(), processingConfig.EnableSanitization)
	}

	if processingConfig.DefaultDelimiter != config.LogicalNameDelimiter() {
		t.Errorf("expected DefaultDelimiter %s, got %s", config.LogicalNameDelimiter(), processingConfig.DefaultDelimiter)
	}

	if processingConfig.StrictMode != config.IsEnhancedCommentStrictMode() {
		t.Errorf("expected StrictMode %v, got %v", config.IsEnhancedCommentStrictMode(), processingConfig.StrictMode)
	}

	if processingConfig.ProcessingTimeout != config.EnhancedCommentProcessingTimeout() {
		t.Errorf("expected ProcessingTimeout %d, got %d", config.EnhancedCommentProcessingTimeout(), processingConfig.ProcessingTimeout)
	}

	// パーサーが登録されているか確認
	formats := processor.GetSupportedFormats()
	if len(formats) == 0 {
		t.Error("no parsers registered")
	}

	// JSONとYAMLパーサーが登録されているか確認
	hasJSON := processor.HasParser("json")
	hasYAML := processor.HasParser("yaml")
	hasLegacy := processor.HasParser("legacy")

	if !hasJSON && config.IsEnhancedCommentJSONEnabled() {
		t.Error("JSON parser should be registered")
	}

	if !hasYAML && config.IsEnhancedCommentYAMLEnabled() {
		t.Error("YAML parser should be registered")
	}

	if !hasLegacy {
		t.Error("Legacy parser should always be registered")
	}
}

func TestNewEnhancedCommentProcessorFromConfigDisabled(t *testing.T) {
	config := NewMockEnhancedCommentConfigurator()
	config.enabled = false

	processor := NewEnhancedCommentProcessorFromConfig(config)
	if processor == nil {
		t.Fatal("processor should not be nil")
	}

	// 無効な場合はデフォルト設定で作成されるべき
	processingConfig := processor.GetConfig()
	if processingConfig.EnableValidation != true {
		t.Error("expected default EnableValidation to be true")
	}

	if processingConfig.DefaultDelimiter != "|" {
		t.Errorf("expected default DefaultDelimiter '|', got %s", processingConfig.DefaultDelimiter)
	}
}

func TestRegisterParsersFromConfigAutoFormat(t *testing.T) {
	config := NewMockEnhancedCommentConfigurator()
	config.preferredFormat = "auto"

	processor := NewEnhancedCommentProcessorFromConfig(config)

	// autoの場合はJSON、YAML、Legacyの順で登録される
	formats := processor.GetSupportedFormats()
	expectedFormats := []string{"json", "yaml", "legacy"}

	if len(formats) != len(expectedFormats) {
		t.Errorf("expected %d formats, got %d", len(expectedFormats), len(formats))
	}

	for i, expected := range expectedFormats {
		if i < len(formats) && formats[i] != expected {
			t.Errorf("expected format[%d] %s, got %s", i, expected, formats[i])
		}
	}
}

func TestRegisterParsersFromConfigJSONPreferred(t *testing.T) {
	config := NewMockEnhancedCommentConfigurator()
	config.preferredFormat = "json"
	config.yamlEnabled = true

	processor := NewEnhancedCommentProcessorFromConfig(config)

	// JSONが優先される場合の順序確認
	formats := processor.GetSupportedFormats()
	if len(formats) < 3 {
		t.Fatal("expected at least 3 formats")
	}

	// JSONが最初に来るべき
	if formats[0] != "json" {
		t.Errorf("expected first format to be json, got %s", formats[0])
	}
}

func TestRegisterParsersFromConfigYAMLPreferred(t *testing.T) {
	config := NewMockEnhancedCommentConfigurator()
	config.preferredFormat = "yaml"
	config.jsonEnabled = true

	processor := NewEnhancedCommentProcessorFromConfig(config)

	// YAMLが優先される場合の順序確認
	formats := processor.GetSupportedFormats()
	if len(formats) < 3 {
		t.Fatal("expected at least 3 formats")
	}

	// YAMLが最初に来るべき
	if formats[0] != "yaml" {
		t.Errorf("expected first format to be yaml, got %s", formats[0])
	}
}

func TestRegisterParsersFromConfigLegacyOnly(t *testing.T) {
	config := NewMockEnhancedCommentConfigurator()
	config.preferredFormat = "legacy"

	processor := NewEnhancedCommentProcessorFromConfig(config)

	// legacyのみの場合
	formats := processor.GetSupportedFormats()
	if len(formats) != 1 {
		t.Errorf("expected 1 format for legacy only, got %d", len(formats))
	}

	if formats[0] != "legacy" {
		t.Errorf("expected format to be legacy, got %s", formats[0])
	}
}

func TestRegisterParsersFromConfigJSONDisabled(t *testing.T) {
	config := NewMockEnhancedCommentConfigurator()
	config.jsonEnabled = false
	config.yamlEnabled = true

	processor := NewEnhancedCommentProcessorFromConfig(config)

	if processor.HasParser("json") {
		t.Error("JSON parser should not be registered when disabled")
	}

	if !processor.HasParser("yaml") {
		t.Error("YAML parser should be registered when enabled")
	}

	if !processor.HasParser("legacy") {
		t.Error("Legacy parser should always be registered")
	}
}

func TestCreateValidatorFromConfigDefault(t *testing.T) {
	config := NewMockEnhancedCommentConfigurator()
	config.securityLevel = "default"

	validator := createValidatorFromConfig(config)
	if validator == nil {
		t.Fatal("validator should not be nil")
	}

	// デフォルト設定のバリデーターかどうかの確認
	// バリデーター設定の内容確認
	if defaultValidator, ok := validator.(*DefaultCommentValidator); ok {
		validatorConfig := defaultValidator.GetConfig()
		if validatorConfig.MaxLogicalNameLength != 100 {
			t.Errorf("expected default MaxLogicalNameLength 100, got %d", validatorConfig.MaxLogicalNameLength)
		}
	}
}

func TestCreateValidatorFromConfigStrict(t *testing.T) {
	config := NewMockEnhancedCommentConfigurator()
	config.securityLevel = "strict"

	validator := createValidatorFromConfig(config)
	if validator == nil {
		t.Fatal("validator should not be nil")
	}

	// 厳格設定のバリデーターかどうかの確認
	if defaultValidator, ok := validator.(*DefaultCommentValidator); ok {
		validatorConfig := defaultValidator.GetConfig()
		if validatorConfig.MaxLogicalNameLength != 50 {
			t.Errorf("expected strict MaxLogicalNameLength 50, got %d", validatorConfig.MaxLogicalNameLength)
		}
	}
}

func TestCreateValidatorFromConfigPermissive(t *testing.T) {
	config := NewMockEnhancedCommentConfigurator()
	config.securityLevel = "permissive"

	validator := createValidatorFromConfig(config)
	if validator == nil {
		t.Fatal("validator should not be nil")
	}

	// 緩い設定のバリデーターかどうかの確認
	if defaultValidator, ok := validator.(*DefaultCommentValidator); ok {
		validatorConfig := defaultValidator.GetConfig()
		if validatorConfig.MaxLogicalNameLength != 200 {
			t.Errorf("expected permissive MaxLogicalNameLength 200, got %d", validatorConfig.MaxLogicalNameLength)
		}
	}
}

func TestCreateValidatorFromConfigValidationDisabled(t *testing.T) {
	config := NewMockEnhancedCommentConfigurator()
	config.validationEnabled = false

	validator := createValidatorFromConfig(config)
	if validator == nil {
		t.Fatal("validator should not be nil even when validation is disabled")
	}

	// バリデーション無効でもバリデーターは作成される（デフォルト設定で）
	if defaultValidator, ok := validator.(*DefaultCommentValidator); ok {
		validatorConfig := defaultValidator.GetConfig()
		if validatorConfig.MaxLogicalNameLength != 100 {
			t.Errorf("expected default MaxLogicalNameLength 100 when validation disabled, got %d", validatorConfig.MaxLogicalNameLength)
		}
	}
}

func TestEnhancedCommentProcessorConfigIntegration(t *testing.T) {
	// 統合テスト：設定からプロセッサーを作成し、実際にコメントを処理
	config := NewMockEnhancedCommentConfigurator()
	config.preferredFormat = "json"
	config.strictMode = false
	config.validationEnabled = true
	config.sanitizationEnabled = true

	processor := NewEnhancedCommentProcessorFromConfig(config)

	// JSONコメントの処理テスト
	jsonComment := `{"name": "ユーザー名", "description": "ユーザーの表示名"}`
	result, err := processor.ProcessCommentWithValidation(jsonComment, "|", ObjectTypeColumn)

	if err != nil {
		t.Fatalf("failed to process JSON comment: %v", err)
	}

	if result.LogicalName != "ユーザー名" {
		t.Errorf("expected LogicalName 'ユーザー名', got %s", result.LogicalName)
	}

	if result.Description != "ユーザーの表示名" {
		t.Errorf("expected Description 'ユーザーの表示名', got %s", result.Description)
	}

	// オブジェクトタイプがメタデータに設定されているか確認
	if result.Metadata["object_type"] != string(ObjectTypeColumn) {
		t.Errorf("expected object_type metadata to be %s, got %s", ObjectTypeColumn, result.Metadata["object_type"])
	}
}

func TestEnhancedCommentProcessorConfigStrictMode(t *testing.T) {
	// 厳格モードでのエラーハンドリングテスト
	config := NewMockEnhancedCommentConfigurator()
	config.strictMode = true
	config.validationEnabled = true
	config.securityLevel = "strict" // 厳格なバリデーション設定

	processor := NewEnhancedCommentProcessorFromConfig(config)

	// バリデーションエラーを引き起こすコメント（strictでは50文字制限）
	longName := ""
	for i := 0; i < 51; i++ {
		longName += "あ" // 51文字の日本語文字
	}
	invalidComment := `{"name": "` + longName + `"}`
	
	_, err := processor.ProcessCommentWithValidation(invalidComment, "|", ObjectTypeColumn)

	if err == nil {
		t.Error("expected error in strict mode for comment with too long name")
	}
}

func TestEnhancedCommentProcessorConfigNonStrictMode(t *testing.T) {
	// 非厳格モードでのエラーハンドリングテスト
	config := NewMockEnhancedCommentConfigurator()
	config.strictMode = false
	config.validationEnabled = true

	processor := NewEnhancedCommentProcessorFromConfig(config)

	// 無効なコメント
	invalidComment := `{invalid json`
	result, err := processor.ProcessCommentWithValidation(invalidComment, "|", ObjectTypeColumn)

	if err != nil {
		t.Errorf("expected no error in non-strict mode, got: %v", err)
	}

	// 非厳格モードでは空のCommentDataが返される
	if result == nil {
		t.Error("expected result not to be nil in non-strict mode")
	}

	if result.Source != invalidComment {
		t.Errorf("expected Source to be preserved: %s", invalidComment)
	}
}