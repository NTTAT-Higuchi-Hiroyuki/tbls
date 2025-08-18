package schema

import (
	"fmt"
	"strings"
	"time"
)

// ErrorSeverity エラーの重要度
type ErrorSeverity int

const (
	ErrorSeverityInfo ErrorSeverity = iota
	ErrorSeverityWarning
	ErrorSeverityError
	ErrorSeverityCritical
)

func (s ErrorSeverity) String() string {
	switch s {
	case ErrorSeverityInfo:
		return "INFO"
	case ErrorSeverityWarning:
		return "WARNING"
	case ErrorSeverityError:
		return "ERROR"
	case ErrorSeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// ErrorCategory エラーのカテゴリ
type ErrorCategory int

const (
	ErrorCategoryValidation ErrorCategory = iota
	ErrorCategoryParsing
	ErrorCategoryProcessing
	ErrorCategoryConfiguration
	ErrorCategoryTimeout
	ErrorCategoryMemory
	ErrorCategorySchema
)

func (c ErrorCategory) String() string {
	switch c {
	case ErrorCategoryValidation:
		return "VALIDATION"
	case ErrorCategoryParsing:
		return "PARSING"
	case ErrorCategoryProcessing:
		return "PROCESSING"
	case ErrorCategoryConfiguration:
		return "CONFIGURATION"
	case ErrorCategoryTimeout:
		return "TIMEOUT"
	case ErrorCategoryMemory:
		return "MEMORY"
	case ErrorCategorySchema:
		return "SCHEMA"
	default:
		return "UNKNOWN"
	}
}

// EnhancedCommentError 拡張されたエラー情報
type EnhancedCommentError struct {
	Message      string                 `json:"message"`
	Severity     ErrorSeverity         `json:"severity"`
	Category     ErrorCategory         `json:"category"`
	ErrorCode    string                `json:"error_code"`
	Context      map[string]interface{} `json:"context"`
	Timestamp    time.Time             `json:"timestamp"`
	ObjectType   ObjectType            `json:"object_type,omitempty"`
	ObjectName   string                `json:"object_name,omitempty"`
	ParserName   string                `json:"parser_name,omitempty"`
	SourceComment string               `json:"source_comment,omitempty"`
	Suggestions  []string              `json:"suggestions,omitempty"`
	InnerError   error                 `json:"-"`
}

// Error Error インターフェースの実装
func (e *EnhancedCommentError) Error() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("[%s:%s]", e.Severity, e.Category))
	parts = append(parts, e.ErrorCode)
	parts = append(parts, e.Message)
	
	if e.ObjectName != "" {
		parts = append(parts, fmt.Sprintf("(object: %s)", e.ObjectName))
	}
	
	if e.ParserName != "" {
		parts = append(parts, fmt.Sprintf("(parser: %s)", e.ParserName))
	}
	
	return strings.Join(parts, " ")
}

// Unwrap エラーの展開
func (e *EnhancedCommentError) Unwrap() error {
	return e.InnerError
}

// ErrorBuilder エラー構築ヘルパー
type ErrorBuilder struct {
	error *EnhancedCommentError
}

// NewErrorBuilder 新しいエラービルダーを作成
func NewErrorBuilder() *ErrorBuilder {
	return &ErrorBuilder{
		error: &EnhancedCommentError{
			Context:   make(map[string]interface{}),
			Timestamp: time.Now(),
		},
	}
}

// WithMessage メッセージを設定
func (b *ErrorBuilder) WithMessage(message string) *ErrorBuilder {
	b.error.Message = message
	return b
}

// WithSeverity 重要度を設定
func (b *ErrorBuilder) WithSeverity(severity ErrorSeverity) *ErrorBuilder {
	b.error.Severity = severity
	return b
}

// WithCategory カテゴリを設定
func (b *ErrorBuilder) WithCategory(category ErrorCategory) *ErrorBuilder {
	b.error.Category = category
	return b
}

// WithErrorCode エラーコードを設定
func (b *ErrorBuilder) WithErrorCode(code string) *ErrorBuilder {
	b.error.ErrorCode = code
	return b
}

// WithObjectInfo オブジェクト情報を設定
func (b *ErrorBuilder) WithObjectInfo(objectType ObjectType, objectName string) *ErrorBuilder {
	b.error.ObjectType = objectType
	b.error.ObjectName = objectName
	return b
}

// WithParserName パーサー名を設定
func (b *ErrorBuilder) WithParserName(parserName string) *ErrorBuilder {
	b.error.ParserName = parserName
	return b
}

// WithSourceComment ソースコメントを設定
func (b *ErrorBuilder) WithSourceComment(comment string) *ErrorBuilder {
	b.error.SourceComment = comment
	return b
}

// WithContext コンテキスト情報を追加
func (b *ErrorBuilder) WithContext(key string, value interface{}) *ErrorBuilder {
	b.error.Context[key] = value
	return b
}

// WithSuggestion 提案を追加
func (b *ErrorBuilder) WithSuggestion(suggestion string) *ErrorBuilder {
	b.error.Suggestions = append(b.error.Suggestions, suggestion)
	return b
}

// WithInnerError 内部エラーを設定
func (b *ErrorBuilder) WithInnerError(err error) *ErrorBuilder {
	b.error.InnerError = err
	return b
}

// Build エラーを構築
func (b *ErrorBuilder) Build() *EnhancedCommentError {
	return b.error
}

// ErrorRecoveryStrategy エラー回復戦略
type ErrorRecoveryStrategy interface {
	CanRecover(err *EnhancedCommentError) bool
	Recover(err *EnhancedCommentError, context interface{}) (interface{}, error)
	GetDescription() string
}

// FallbackParsingStrategy フォールバック解析戦略
type FallbackParsingStrategy struct{}

func (s *FallbackParsingStrategy) CanRecover(err *EnhancedCommentError) bool {
	return err.Category == ErrorCategoryParsing && 
		   (err.ErrorCode == "E_JSON_PARSE_FAILED" || err.ErrorCode == "E_YAML_PARSE_FAILED")
}

func (s *FallbackParsingStrategy) Recover(err *EnhancedCommentError, context interface{}) (interface{}, error) {
	if comment, ok := context.(string); ok {
		// Legacy parserでフォールバック
		legacyParser := NewLegacyParser()
		return legacyParser.ParseComment(comment, "|")
	}
	return nil, fmt.Errorf("cannot recover: invalid context type")
}

func (s *FallbackParsingStrategy) GetDescription() string {
	return "フォールバック解析戦略: JSON/YAML解析失敗時にLegacy形式で再試行"
}

// CommentSanitizationStrategy コメントサニタイゼーション戦略
type CommentSanitizationStrategy struct{}

func (s *CommentSanitizationStrategy) CanRecover(err *EnhancedCommentError) bool {
	return err.Category == ErrorCategoryValidation &&
		   err.ErrorCode == "E_UNSAFE_CONTENT"
}

func (s *CommentSanitizationStrategy) Recover(err *EnhancedCommentError, context interface{}) (interface{}, error) {
	if commentData, ok := context.(*CommentData); ok {
		// 安全でない内容をサニタイズ
		sanitized := &CommentData{
			LogicalName: sanitizeString(commentData.LogicalName),
			Description: sanitizeString(commentData.Description),
			Tags:        sanitizeStringSlice(commentData.Tags),
			Metadata:    sanitizeMetadata(commentData.Metadata),
			Priority:    commentData.Priority,
			Deprecated:  commentData.Deprecated,
			Source:      commentData.Source,
		}
		return sanitized, nil
	}
	return nil, fmt.Errorf("cannot recover: invalid context type")
}

func (s *CommentSanitizationStrategy) GetDescription() string {
	return "コメントサニタイゼーション戦略: 安全でない内容を自動的にサニタイズ"
}

// ErrorRecoveryManager エラー回復管理
type ErrorRecoveryManager struct {
	strategies []ErrorRecoveryStrategy
}

// NewErrorRecoveryManager 新しいエラー回復管理を作成
func NewErrorRecoveryManager() *ErrorRecoveryManager {
	return &ErrorRecoveryManager{
		strategies: []ErrorRecoveryStrategy{
			&FallbackParsingStrategy{},
			&CommentSanitizationStrategy{},
		},
	}
}

// AddStrategy 回復戦略を追加
func (m *ErrorRecoveryManager) AddStrategy(strategy ErrorRecoveryStrategy) {
	m.strategies = append(m.strategies, strategy)
}

// TryRecover エラー回復を試行
func (m *ErrorRecoveryManager) TryRecover(err *EnhancedCommentError, context interface{}) (interface{}, error) {
	for _, strategy := range m.strategies {
		if strategy.CanRecover(err) {
			result, recoveryErr := strategy.Recover(err, context)
			if recoveryErr == nil {
				return result, nil
			}
		}
	}
	return nil, fmt.Errorf("no recovery strategy available for error: %s", err.ErrorCode)
}

// GetAvailableStrategies 利用可能な戦略一覧を取得
func (m *ErrorRecoveryManager) GetAvailableStrategies() []string {
	descriptions := make([]string, len(m.strategies))
	for i, strategy := range m.strategies {
		descriptions[i] = strategy.GetDescription()
	}
	return descriptions
}

// ErrorReporter エラー報告インターフェース
type ErrorReporter interface {
	ReportError(err *EnhancedCommentError)
	ReportRecovery(originalErr *EnhancedCommentError, recoveredData interface{})
	GetErrorSummary() *ErrorSummary
}

// ErrorSummary エラー概要
type ErrorSummary struct {
	TotalErrors    int                        `json:"total_errors"`
	ErrorsBySeverity map[ErrorSeverity]int    `json:"errors_by_severity"`
	ErrorsByCategory map[ErrorCategory]int    `json:"errors_by_category"`
	RecoveredErrors  int                      `json:"recovered_errors"`
	RecentErrors     []*EnhancedCommentError  `json:"recent_errors"`
}

// DefaultErrorReporter デフォルトエラーレポーター
type DefaultErrorReporter struct {
	errors    []*EnhancedCommentError
	recovered int
}

// NewDefaultErrorReporter 新しいデフォルトエラーレポーターを作成
func NewDefaultErrorReporter() *DefaultErrorReporter {
	return &DefaultErrorReporter{
		errors: make([]*EnhancedCommentError, 0),
	}
}

// ReportError エラーを報告
func (r *DefaultErrorReporter) ReportError(err *EnhancedCommentError) {
	r.errors = append(r.errors, err)
}

// ReportRecovery 回復を報告
func (r *DefaultErrorReporter) ReportRecovery(originalErr *EnhancedCommentError, recoveredData interface{}) {
	r.recovered++
	
	// 回復成功を示すエラーエントリを追加
	recoveryErr := NewErrorBuilder().
		WithMessage(fmt.Sprintf("Recovered from error: %s", originalErr.Message)).
		WithSeverity(ErrorSeverityInfo).
		WithCategory(ErrorCategoryProcessing).
		WithErrorCode("E_RECOVERY_SUCCESS").
		WithObjectInfo(originalErr.ObjectType, originalErr.ObjectName).
		WithContext("original_error", originalErr.ErrorCode).
		WithContext("recovered_data_type", fmt.Sprintf("%T", recoveredData)).
		Build()
	
	r.errors = append(r.errors, recoveryErr)
}

// GetErrorSummary エラー概要を取得
func (r *DefaultErrorReporter) GetErrorSummary() *ErrorSummary {
	summary := &ErrorSummary{
		TotalErrors:      len(r.errors),
		ErrorsBySeverity: make(map[ErrorSeverity]int),
		ErrorsByCategory: make(map[ErrorCategory]int),
		RecoveredErrors:  r.recovered,
		RecentErrors:     make([]*EnhancedCommentError, 0),
	}
	
	// 重要度別カウント
	for _, err := range r.errors {
		summary.ErrorsBySeverity[err.Severity]++
		summary.ErrorsByCategory[err.Category]++
	}
	
	// 最新10件のエラー
	start := len(r.errors) - 10
	if start < 0 {
		start = 0
	}
	summary.RecentErrors = append(summary.RecentErrors, r.errors[start:]...)
	
	return summary
}

// サニタイゼーションヘルパー関数
func sanitizeString(s string) string {
	// 基本的なサニタイゼーション（実装は要件に応じて調整）
	s = strings.ReplaceAll(s, "<script>", "")
	s = strings.ReplaceAll(s, "</script>", "")
	s = strings.ReplaceAll(s, "javascript:", "")
	return s
}

func sanitizeStringSlice(slice []string) []string {
	if slice == nil {
		return nil
	}
	result := make([]string, len(slice))
	for i, s := range slice {
		result[i] = sanitizeString(s)
	}
	return result
}

func sanitizeMetadata(metadata map[string]string) map[string]string {
	if metadata == nil {
		return nil
	}
	result := make(map[string]string)
	for k, v := range metadata {
		result[sanitizeString(k)] = sanitizeString(v)
	}
	return result
}

// 事前定義されたエラーコード
const (
	// パーシングエラー
	ErrorCodeJSONParseFailed   = "E_JSON_PARSE_FAILED"
	ErrorCodeYAMLParseFailed   = "E_YAML_PARSE_FAILED"
	ErrorCodeLegacyParseFailed = "E_LEGACY_PARSE_FAILED"
	
	// バリデーションエラー
	ErrorCodeValidationFailed = "E_VALIDATION_FAILED"
	ErrorCodeUnsafeContent    = "E_UNSAFE_CONTENT"
	ErrorCodeInvalidFormat    = "E_INVALID_FORMAT"
	
	// 処理エラー
	ErrorCodeProcessingTimeout = "E_PROCESSING_TIMEOUT"
	ErrorCodeMemoryLimit       = "E_MEMORY_LIMIT"
	ErrorCodeInternalError     = "E_INTERNAL_ERROR"
	
	// 設定エラー
	ErrorCodeInvalidConfig     = "E_INVALID_CONFIG"
	ErrorCodeMissingConfig     = "E_MISSING_CONFIG"
	
	// スキーマエラー
	ErrorCodeInvalidObjectType = "E_INVALID_OBJECT_TYPE"
	ErrorCodeMissingObject     = "E_MISSING_OBJECT"
)

// 便利な関数群
func NewParsingError(parserName, comment string, innerErr error) *EnhancedCommentError {
	var errorCode string
	switch parserName {
	case "json":
		errorCode = ErrorCodeJSONParseFailed
	case "yaml":
		errorCode = ErrorCodeYAMLParseFailed
	case "legacy":
		errorCode = ErrorCodeLegacyParseFailed
	default:
		errorCode = "E_UNKNOWN_PARSER_FAILED"
	}
	
	return NewErrorBuilder().
		WithMessage(fmt.Sprintf("%s parser failed to parse comment", parserName)).
		WithSeverity(ErrorSeverityError).
		WithCategory(ErrorCategoryParsing).
		WithErrorCode(errorCode).
		WithParserName(parserName).
		WithSourceComment(comment).
		WithInnerError(innerErr).
		WithSuggestion("コメント形式を確認してください").
		WithSuggestion("フォールバック解析を有効にしてください").
		Build()
}

func NewValidationError(message string, objectType ObjectType, objectName string) *EnhancedCommentError {
	return NewErrorBuilder().
		WithMessage(message).
		WithSeverity(ErrorSeverityWarning).
		WithCategory(ErrorCategoryValidation).
		WithErrorCode(ErrorCodeValidationFailed).
		WithObjectInfo(objectType, objectName).
		WithSuggestion("コメント内容を見直してください").
		Build()
}

func NewProcessingTimeoutError(timeout time.Duration, objectType ObjectType, objectName string) *EnhancedCommentError {
	return NewErrorBuilder().
		WithMessage(fmt.Sprintf("Processing timeout after %v", timeout)).
		WithSeverity(ErrorSeverityError).
		WithCategory(ErrorCategoryTimeout).
		WithErrorCode(ErrorCodeProcessingTimeout).
		WithObjectInfo(objectType, objectName).
		WithContext("timeout_duration", timeout.String()).
		WithSuggestion("処理タイムアウト値を増やしてください").
		WithSuggestion("コメントの複雑さを減らしてください").
		Build()
}

func NewConfigurationError(message string, configKey string) *EnhancedCommentError {
	return NewErrorBuilder().
		WithMessage(message).
		WithSeverity(ErrorSeverityCritical).
		WithCategory(ErrorCategoryConfiguration).
		WithErrorCode(ErrorCodeInvalidConfig).
		WithContext("config_key", configKey).
		WithSuggestion("設定ファイルを確認してください").
		Build()
}