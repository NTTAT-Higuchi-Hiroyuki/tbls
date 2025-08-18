package schema

import (
	"fmt"
	"time"
)

// CommentProcessor 拡張されたコメント処理を行うインターフェース
type CommentProcessor interface {
	// ProcessComment コメントを処理してCommentDataに変換
	ProcessComment(comment string, delimiter string, objectType ObjectType) (*CommentData, error)
	// ProcessCommentWithValidation バリデーション付きでコメントを処理
	ProcessCommentWithValidation(comment string, delimiter string, objectType ObjectType) (*CommentData, error)
	// RegisterParser パーサーを登録
	RegisterParser(parser CommentParser)
	// SetValidator バリデーターを設定
	SetValidator(validator CommentValidator)
	// GetSupportedFormats サポートされているフォーマット一覧を取得
	GetSupportedFormats() []string
}

// ProcessingConfig コメント処理設定
type ProcessingConfig struct {
	// EnableValidation バリデーションを有効にするか
	EnableValidation bool
	// EnableSanitization サニタイゼーションを有効にするか
	EnableSanitization bool
	// DefaultDelimiter デフォルトの区切り文字
	DefaultDelimiter string
	// FallbackToLegacy 他のパーサーが失敗した場合にLegacyParserにフォールバックするか
	FallbackToLegacy bool
	// StrictMode 厳格モード（エラー時に失敗）
	StrictMode bool
	// ProcessingTimeout 処理タイムアウト（ミリ秒）
	ProcessingTimeout int
}

// DefaultProcessingConfig デフォルトの処理設定
func DefaultProcessingConfig() *ProcessingConfig {
	return &ProcessingConfig{
		EnableValidation:   true,
		EnableSanitization: true,
		DefaultDelimiter:   "|",
		FallbackToLegacy:   true,
		StrictMode:         false,
		ProcessingTimeout:  1000, // 1秒
	}
}

// EnhancedCommentProcessor 拡張されたコメントプロセッサーの実装
type EnhancedCommentProcessor struct {
	registry  ParserRegistry
	validator CommentValidator
	config    *ProcessingConfig
}

// NewEnhancedCommentProcessor 新しいEnhancedCommentProcessorを作成
func NewEnhancedCommentProcessor() *EnhancedCommentProcessor {
	processor := &EnhancedCommentProcessor{
		registry:  NewDefaultParserRegistry(),
		validator: NewDefaultCommentValidator(),
		config:    DefaultProcessingConfig(),
	}

	// デフォルトパーサーを登録
	processor.registerDefaultParsers()

	return processor
}

// NewEnhancedCommentProcessorWithConfig 設定を指定してEnhancedCommentProcessorを作成
func NewEnhancedCommentProcessorWithConfig(config *ProcessingConfig) *EnhancedCommentProcessor {
	processor := &EnhancedCommentProcessor{
		registry:  NewDefaultParserRegistry(),
		validator: NewDefaultCommentValidator(),
		config:    config,
	}

	// デフォルトパーサーを登録
	processor.registerDefaultParsers()

	return processor
}

// NewEnhancedCommentProcessorFromConfig config/config.goの設定からEnhancedCommentProcessorを作成
func NewEnhancedCommentProcessorFromConfig(config EnhancedCommentConfigurator) *EnhancedCommentProcessor {
	if !config.IsEnhancedCommentEnabled() {
		// 拡張コメント処理が無効な場合はデフォルト設定で作成
		return NewEnhancedCommentProcessor()
	}

	// 設定から処理設定を構築
	processingConfig := &ProcessingConfig{
		EnableValidation:   config.IsEnhancedCommentValidationEnabled(),
		EnableSanitization: config.IsEnhancedCommentSanitizationEnabled(),
		DefaultDelimiter:   config.LogicalNameDelimiter(), // 既存のlogical name delimiterを使用
		FallbackToLegacy:   true, // 常にlegacyにフォールバック
		StrictMode:         config.IsEnhancedCommentStrictMode(),
		ProcessingTimeout:  config.EnhancedCommentProcessingTimeout(),
	}

	// プロセッサー作成
	processor := &EnhancedCommentProcessor{
		registry:  NewDefaultParserRegistry(),
		validator: createValidatorFromConfig(config),
		config:    processingConfig,
	}

	// 設定に基づいてパーサーを登録
	processor.registerParsersFromConfig(config)

	return processor
}

// EnhancedCommentConfigurator 設定インターフェース（config/config.goの依存を避けるため）
type EnhancedCommentConfigurator interface {
	IsEnhancedCommentEnabled() bool
	IsEnhancedCommentJSONEnabled() bool
	IsEnhancedCommentYAMLEnabled() bool
	EnhancedCommentPreferredFormat() string
	IsEnhancedCommentValidationEnabled() bool
	IsEnhancedCommentSanitizationEnabled() bool
	EnhancedCommentSecurityLevel() string
	IsEnhancedCommentStrictMode() bool
	EnhancedCommentProcessingTimeout() int
	IsEnhancedCommentObjectTypeEnabled(objectType string) bool
	LogicalNameDelimiter() string
}

// registerParsersFromConfig 設定に基づいてパーサーを登録
func (p *EnhancedCommentProcessor) registerParsersFromConfig(config EnhancedCommentConfigurator) {
	// 優先フォーマットに基づいてパーサーの登録順序を決定
	preferredFormat := config.EnhancedCommentPreferredFormat()

	switch preferredFormat {
	case "json":
		if config.IsEnhancedCommentJSONEnabled() {
			// JSON優先の場合、JSONの優先度を最高に設定
			jsonParser := NewJSONParser()
			jsonParser.SetPriority(5) // 最高優先度
			p.registry.RegisterParser(jsonParser)
		}
		if config.IsEnhancedCommentYAMLEnabled() {
			yamlParser := NewYAMLParser()
			yamlParser.SetPriority(15) // 通常優先度
			p.registry.RegisterParser(yamlParser)
		}
	case "yaml":
		if config.IsEnhancedCommentYAMLEnabled() {
			// YAML優先の場合、YAMLの優先度を最高に設定
			yamlParser := NewYAMLParser()
			yamlParser.SetPriority(5) // 最高優先度
			p.registry.RegisterParser(yamlParser)
		}
		if config.IsEnhancedCommentJSONEnabled() {
			jsonParser := NewJSONParser()
			jsonParser.SetPriority(10) // 通常優先度
			p.registry.RegisterParser(jsonParser)
		}
	case "legacy":
		// legacyパーサーのみ
		legacyParser := NewLegacyParser()
		legacyParser.SetPriority(5) // 最高優先度
		p.registry.RegisterParser(legacyParser)
		return
	default: // "auto" または未知の値
		// JSON -> YAML -> Legacyの順で登録（優先度順）
		if config.IsEnhancedCommentJSONEnabled() {
			jsonParser := NewJSONParser()
			jsonParser.SetPriority(10) // 高優先度
			p.registry.RegisterParser(jsonParser)
		}
		if config.IsEnhancedCommentYAMLEnabled() {
			yamlParser := NewYAMLParser()
			yamlParser.SetPriority(15) // 中優先度
			p.registry.RegisterParser(yamlParser)
		}
	}

	// 常にlegacyパーサーをフォールバックとして登録
	legacyParser := NewLegacyParser()
	legacyParser.SetPriority(20) // 最低優先度
	p.registry.RegisterParser(legacyParser)
}

// createValidatorFromConfig 設定からバリデーターを作成
func createValidatorFromConfig(config EnhancedCommentConfigurator) CommentValidator {
	if !config.IsEnhancedCommentValidationEnabled() {
		// バリデーション無効の場合はデフォルトバリデーターを使用
		return NewDefaultCommentValidator()
	}

	// セキュリティレベルに基づいてバリデーション設定を決定
	var validationConfig *ValidationConfig
	switch config.EnhancedCommentSecurityLevel() {
	case "strict":
		validationConfig = StrictValidationConfig()
	case "permissive":
		validationConfig = PermissiveValidationConfig()
	default: // "default"
		validationConfig = DefaultValidationConfig()
	}

	// サニタイゼーション設定を反映
	validationConfig.EnableHTMLEscape = config.IsEnhancedCommentSanitizationEnabled()
	validationConfig.EnableSQLInjectionCheck = config.IsEnhancedCommentValidationEnabled()

	return NewDefaultCommentValidatorWithConfig(validationConfig)
}

// registerDefaultParsers デフォルトパーサーを登録
func (p *EnhancedCommentProcessor) registerDefaultParsers() {
	// JSONParserを登録（高優先度）
	p.registry.RegisterParser(NewJSONParser())

	// LegacyParserを登録（低優先度、フォールバック用）
	p.registry.RegisterParser(NewLegacyParser())
}

// ProcessComment コメントを処理してCommentDataに変換
func (p *EnhancedCommentProcessor) ProcessComment(comment string, delimiter string, objectType ObjectType) (*CommentData, error) {
	if comment == "" {
		return &CommentData{Source: comment}, nil
	}

	// デフォルト区切り文字の適用
	if delimiter == "" {
		delimiter = p.config.DefaultDelimiter
	}

	// タイムアウト処理（基本的な実装）
	done := make(chan *ProcessingResult, 1)
	go func() {
		result, err := p.processCommentInternal(comment, delimiter, objectType)
		done <- &ProcessingResult{Data: result, Error: err}
	}()

	select {
	case result := <-done:
		return result.Data, result.Error
	case <-time.After(time.Duration(p.config.ProcessingTimeout) * time.Millisecond):
		return nil, fmt.Errorf("comment processing timeout after %d ms", p.config.ProcessingTimeout)
	}
}

// ProcessingResult 処理結果
type ProcessingResult struct {
	Data  *CommentData
	Error error
}

// processCommentInternal 内部的なコメント処理
func (p *EnhancedCommentProcessor) processCommentInternal(comment string, delimiter string, objectType ObjectType) (*CommentData, error) {
	// パーサーレジストリを使用してコメントを解析
	result, err := p.registry.ParseWithFallback(comment, delimiter)
	if err != nil {
		if p.config.StrictMode {
			return nil, fmt.Errorf("failed to parse comment: %w", err)
		}
		// 非厳格モードでは空のCommentDataを返す
		return &CommentData{Source: comment}, nil
	}

	// オブジェクトタイプをメタデータに追加
	if result.Metadata == nil {
		result.Metadata = make(map[string]string)
	}
	result.Metadata["object_type"] = string(objectType)

	return result, nil
}

// ProcessCommentWithValidation バリデーション付きでコメントを処理
func (p *EnhancedCommentProcessor) ProcessCommentWithValidation(comment string, delimiter string, objectType ObjectType) (*CommentData, error) {
	// 基本処理
	result, err := p.ProcessComment(comment, delimiter, objectType)
	if err != nil {
		return nil, err
	}

	// バリデーションとサニタイゼーション
	if p.config.EnableValidation && p.validator != nil {
		if p.config.EnableSanitization {
			// サニタイゼーション + バリデーション
			validated, err := p.validator.ValidateAndSanitize(result)
			if err != nil {
				if p.config.StrictMode {
					return nil, fmt.Errorf("validation failed: %w", err)
				}
				// 非厳格モードでは警告ログとして記録し、サニタイゼーションのみ適用
				result = p.validator.Sanitize(result)
			} else {
				result = validated
			}
		} else {
			// バリデーションのみ
			err := p.validator.Validate(result)
			if err != nil && p.config.StrictMode {
				return nil, fmt.Errorf("validation failed: %w", err)
			}
		}
	} else if p.config.EnableSanitization && p.validator != nil {
		// サニタイゼーションのみ
		result = p.validator.Sanitize(result)
	}

	return result, nil
}

// RegisterParser パーサーを登録
func (p *EnhancedCommentProcessor) RegisterParser(parser CommentParser) {
	p.registry.RegisterParser(parser)
}

// SetValidator バリデーターを設定
func (p *EnhancedCommentProcessor) SetValidator(validator CommentValidator) {
	p.validator = validator
}

// GetSupportedFormats サポートされているフォーマット一覧を取得
func (p *EnhancedCommentProcessor) GetSupportedFormats() []string {
	parsers := p.registry.GetParsers()
	formats := make([]string, len(parsers))
	for i, parser := range parsers {
		formats[i] = parser.Name()
	}
	return formats
}

// GetConfig 設定を取得
func (p *EnhancedCommentProcessor) GetConfig() *ProcessingConfig {
	return p.config
}

// SetConfig 設定を更新
func (p *EnhancedCommentProcessor) SetConfig(config *ProcessingConfig) {
	p.config = config
}

// GetRegistry パーサーレジストリを取得
func (p *EnhancedCommentProcessor) GetRegistry() ParserRegistry {
	return p.registry
}

// GetValidator バリデーターを取得
func (p *EnhancedCommentProcessor) GetValidator() CommentValidator {
	return p.validator
}

// ProcessMultipleComments 複数のコメントを一括処理
func (p *EnhancedCommentProcessor) ProcessMultipleComments(comments []CommentInput) ([]*CommentData, []error) {
	results := make([]*CommentData, len(comments))
	errors := make([]error, len(comments))

	for i, input := range comments {
		result, err := p.ProcessCommentWithValidation(input.Comment, input.Delimiter, input.ObjectType)
		results[i] = result
		errors[i] = err
	}

	return results, errors
}

// CommentInput コメント入力
type CommentInput struct {
	Comment    string
	Delimiter  string
	ObjectType ObjectType
}

// ProcessingStats 処理統計
type ProcessingStats struct {
	TotalProcessed   int
	SuccessCount     int
	ErrorCount       int
	ValidationErrors int
	ParsingErrors    int
	ProcessingTime   time.Duration
}

// ProcessCommentsWithStats 統計付きでコメントを処理
func (p *EnhancedCommentProcessor) ProcessCommentsWithStats(comments []CommentInput) ([]*CommentData, *ProcessingStats) {
	startTime := time.Now()
	stats := &ProcessingStats{
		TotalProcessed: len(comments),
	}

	results := make([]*CommentData, len(comments))

	for i, input := range comments {
		result, err := p.ProcessCommentWithValidation(input.Comment, input.Delimiter, input.ObjectType)
		results[i] = result

		if err != nil {
			stats.ErrorCount++
			// エラーの種類を判定
			if IsValidationError(err) {
				stats.ValidationErrors++
			} else {
				stats.ParsingErrors++
			}
		} else {
			stats.SuccessCount++
		}
	}

	stats.ProcessingTime = time.Since(startTime)
	return results, stats
}

// IsValidationError エラーがバリデーションエラーかどうかを判定
func IsValidationError(err error) bool {
	_, ok := err.(*CommentValidationError)
	return ok
}

// ProcessCommentForObject オブジェクト固有のコメント処理
func (p *EnhancedCommentProcessor) ProcessCommentForObject(comment string, delimiter string, objectType ObjectType, objectName string) (*CommentData, error) {
	result, err := p.ProcessCommentWithValidation(comment, delimiter, objectType)
	if err != nil {
		return nil, err
	}

	// オブジェクト名をメタデータに追加
	if result.Metadata == nil {
		result.Metadata = make(map[string]string)
	}
	result.Metadata["object_name"] = objectName

	return result, nil
}

// ProcessTableComment テーブルコメント専用処理
func (p *EnhancedCommentProcessor) ProcessTableComment(comment string, delimiter string, tableName string) (*CommentData, error) {
	return p.ProcessCommentForObject(comment, delimiter, ObjectTypeTable, tableName)
}

// ProcessColumnComment カラムコメント専用処理
func (p *EnhancedCommentProcessor) ProcessColumnComment(comment string, delimiter string, columnName string) (*CommentData, error) {
	return p.ProcessCommentForObject(comment, delimiter, ObjectTypeColumn, columnName)
}

// ProcessIndexComment インデックスコメント専用処理
func (p *EnhancedCommentProcessor) ProcessIndexComment(comment string, delimiter string, indexName string) (*CommentData, error) {
	return p.ProcessCommentForObject(comment, delimiter, ObjectTypeIndex, indexName)
}

// Clear 全てのパーサーをクリア
func (p *EnhancedCommentProcessor) Clear() {
	p.registry.Clear()
}

// Reset デフォルト状態にリセット
func (p *EnhancedCommentProcessor) Reset() {
	p.Clear()
	p.registerDefaultParsers()
	p.config = DefaultProcessingConfig()
	p.validator = NewDefaultCommentValidator()
}

// HasParser 指定された名前のパーサーが登録されているかを確認
func (p *EnhancedCommentProcessor) HasParser(name string) bool {
	if registry, ok := p.registry.(*DefaultParserRegistry); ok {
		return registry.HasParser(name)
	}
	return false
}

// GetParser 指定された名前のパーサーを取得
func (p *EnhancedCommentProcessor) GetParser(name string) CommentParser {
	if registry, ok := p.registry.(*DefaultParserRegistry); ok {
		return registry.GetParser(name)
	}
	return nil
}