package schema

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// CommentValidator コメントデータの検証とサニタイゼーションを行うインターフェース
type CommentValidator interface {
	// Validate CommentDataの検証を行う
	Validate(data *CommentData) error
	// Sanitize CommentDataのサニタイゼーションを行う
	Sanitize(data *CommentData) *CommentData
	// ValidateAndSanitize 検証とサニタイゼーションを同時に行う
	ValidateAndSanitize(data *CommentData) (*CommentData, error)
}

// ValidationConfig 検証設定
type ValidationConfig struct {
	// MaxLogicalNameLength 論理名の最大長
	MaxLogicalNameLength int
	// MaxDescriptionLength 説明の最大長
	MaxDescriptionLength int
	// MaxTagCount タグの最大数
	MaxTagCount int
	// MaxTagLength 各タグの最大長
	MaxTagLength int
	// MaxMetadataCount メタデータの最大数
	MaxMetadataCount int
	// MaxMetadataKeyLength メタデータキーの最大長
	MaxMetadataKeyLength int
	// MaxMetadataValueLength メタデータ値の最大長
	MaxMetadataValueLength int
	// AllowedCharPattern 許可される文字パターン（正規表現）
	AllowedCharPattern string
	// ForbiddenWords 禁止語句
	ForbiddenWords []string
	// EnableHTMLEscape HTML エスケープを有効にするか
	EnableHTMLEscape bool
	// EnableSQLInjectionCheck SQLインジェクションチェックを有効にするか
	EnableSQLInjectionCheck bool
}

// DefaultValidationConfig デフォルトの検証設定
func DefaultValidationConfig() *ValidationConfig {
	return &ValidationConfig{
		MaxLogicalNameLength:     100,
		MaxDescriptionLength:     1000,
		MaxTagCount:             20,
		MaxTagLength:            50,
		MaxMetadataCount:        50,
		MaxMetadataKeyLength:    100,
		MaxMetadataValueLength:  500,
		AllowedCharPattern:      `^[\p{L}\p{N}\p{P}\p{S}\p{Zs}]+$`, // Unicode文字、数字、句読点、記号、空白
		ForbiddenWords:          []string{"DROP", "DELETE", "INSERT", "UPDATE", "EXEC", "SCRIPT"},
		EnableHTMLEscape:        true,
		EnableSQLInjectionCheck: true,
	}
}

// DefaultCommentValidator デフォルトのコメントバリデーター実装
type DefaultCommentValidator struct {
	config           *ValidationConfig
	allowedCharRegex *regexp.Regexp
	sqlInjectionRegex *regexp.Regexp
}

// NewDefaultCommentValidator 新しいDefaultCommentValidatorを作成
func NewDefaultCommentValidator() *DefaultCommentValidator {
	return NewDefaultCommentValidatorWithConfig(DefaultValidationConfig())
}

// NewDefaultCommentValidatorWithConfig 設定を指定してDefaultCommentValidatorを作成
func NewDefaultCommentValidatorWithConfig(config *ValidationConfig) *DefaultCommentValidator {
	validator := &DefaultCommentValidator{
		config: config,
	}

	// 許可文字パターンのコンパイル
	if config.AllowedCharPattern != "" {
		validator.allowedCharRegex = regexp.MustCompile(config.AllowedCharPattern)
	}

	// SQLインジェクションパターンのコンパイル
	if config.EnableSQLInjectionCheck {
		// 基本的なSQLインジェクションパターン
		sqlPattern := `(?i)(union\s+select|or\s+1\s*=\s*1|and\s+1\s*=\s*1|'|\-\-|\/\*|\*\/|xp_|sp_|exec|execute|drop\s+table|delete\s+from|insert\s+into|update\s+set)`
		validator.sqlInjectionRegex = regexp.MustCompile(sqlPattern)
	}

	return validator
}

// Validate CommentDataの検証を行う
func (v *DefaultCommentValidator) Validate(data *CommentData) error {
	if data == nil {
		return NewCommentValidationError("data", "", "CommentData is nil", nil)
	}

	// 論理名の検証
	if err := v.validateLogicalName(data.LogicalName); err != nil {
		return err
	}

	// 説明の検証
	if err := v.validateDescription(data.Description); err != nil {
		return err
	}

	// タグの検証
	if err := v.validateTags(data.Tags); err != nil {
		return err
	}

	// メタデータの検証
	if err := v.validateMetadata(data.Metadata); err != nil {
		return err
	}

	return nil
}

// Sanitize CommentDataのサニタイゼーションを行う
func (v *DefaultCommentValidator) Sanitize(data *CommentData) *CommentData {
	if data == nil {
		return nil
	}

	sanitized := data.Clone()

	// 論理名のサニタイゼーション
	sanitized.LogicalName = v.sanitizeString(sanitized.LogicalName)

	// 説明のサニタイゼーション
	sanitized.Description = v.sanitizeString(sanitized.Description)

	// タグのサニタイゼーション
	if len(sanitized.Tags) > 0 {
		sanitizedTags := make([]string, 0, len(sanitized.Tags))
		for _, tag := range sanitized.Tags {
			sanitizedTag := v.sanitizeString(tag)
			if sanitizedTag != "" {
				sanitizedTags = append(sanitizedTags, sanitizedTag)
			}
		}
		sanitized.Tags = sanitizedTags
	}

	// メタデータのサニタイゼーション
	if len(sanitized.Metadata) > 0 {
		sanitizedMetadata := make(map[string]string)
		for key, value := range sanitized.Metadata {
			sanitizedKey := v.sanitizeString(key)
			sanitizedValue := v.sanitizeString(value)
			if sanitizedKey != "" && sanitizedValue != "" {
				sanitizedMetadata[sanitizedKey] = sanitizedValue
			}
		}
		sanitized.Metadata = sanitizedMetadata
	}

	return sanitized
}

// ValidateAndSanitize 検証とサニタイゼーションを同時に行う
func (v *DefaultCommentValidator) ValidateAndSanitize(data *CommentData) (*CommentData, error) {
	// まずサニタイゼーションを実行
	sanitized := v.Sanitize(data)

	// サニタイゼーション後のデータを検証
	if err := v.Validate(sanitized); err != nil {
		return nil, err
	}

	return sanitized, nil
}

// validateLogicalName 論理名の検証
func (v *DefaultCommentValidator) validateLogicalName(logicalName string) error {
	if logicalName == "" {
		return nil // 空の論理名は許可
	}

	// 長さチェック
	if utf8.RuneCountInString(logicalName) > v.config.MaxLogicalNameLength {
		return NewCommentValidationError("logical_name", logicalName, 
			fmt.Sprintf("logical name too long (max: %d)", v.config.MaxLogicalNameLength), 
			ErrLogicalNameTooLong)
	}

	// 文字パターンチェック
	if err := v.validateCharPattern(logicalName, "logical_name"); err != nil {
		return err
	}

	// 禁止語句チェック
	if err := v.validateForbiddenWords(logicalName, "logical_name"); err != nil {
		return err
	}

	// SQLインジェクションチェック
	if err := v.validateSQLInjection(logicalName, "logical_name"); err != nil {
		return err
	}

	return nil
}

// validateDescription 説明の検証
func (v *DefaultCommentValidator) validateDescription(description string) error {
	if description == "" {
		return nil // 空の説明は許可
	}

	// 長さチェック
	if utf8.RuneCountInString(description) > v.config.MaxDescriptionLength {
		return NewCommentValidationError("description", description,
			fmt.Sprintf("description too long (max: %d)", v.config.MaxDescriptionLength),
			ErrDescriptionTooLong)
	}

	// 文字パターンチェック
	if err := v.validateCharPattern(description, "description"); err != nil {
		return err
	}

	// 禁止語句チェック
	if err := v.validateForbiddenWords(description, "description"); err != nil {
		return err
	}

	// SQLインジェクションチェック
	if err := v.validateSQLInjection(description, "description"); err != nil {
		return err
	}

	return nil
}

// validateTags タグの検証
func (v *DefaultCommentValidator) validateTags(tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	// タグ数チェック
	if len(tags) > v.config.MaxTagCount {
		return NewCommentValidationError("tags", fmt.Sprintf("%v", tags),
			fmt.Sprintf("too many tags (max: %d)", v.config.MaxTagCount),
			nil)
	}

	// 各タグの検証
	for i, tag := range tags {
		if tag == "" {
			return NewCommentValidationError("tags", tag,
				fmt.Sprintf("empty tag at index %d", i),
				nil)
		}

		// 長さチェック
		if utf8.RuneCountInString(tag) > v.config.MaxTagLength {
			return NewCommentValidationError("tags", tag,
				fmt.Sprintf("tag too long at index %d (max: %d)", i, v.config.MaxTagLength),
				nil)
		}

		// 文字パターンチェック
		if err := v.validateCharPattern(tag, fmt.Sprintf("tags[%d]", i)); err != nil {
			return err
		}

		// 禁止語句チェック
		if err := v.validateForbiddenWords(tag, fmt.Sprintf("tags[%d]", i)); err != nil {
			return err
		}

		// SQLインジェクションチェック
		if err := v.validateSQLInjection(tag, fmt.Sprintf("tags[%d]", i)); err != nil {
			return err
		}
	}

	return nil
}

// validateMetadata メタデータの検証
func (v *DefaultCommentValidator) validateMetadata(metadata map[string]string) error {
	if len(metadata) == 0 {
		return nil
	}

	// メタデータ数チェック
	if len(metadata) > v.config.MaxMetadataCount {
		return NewCommentValidationError("metadata", fmt.Sprintf("%v", metadata),
			fmt.Sprintf("too many metadata entries (max: %d)", v.config.MaxMetadataCount),
			nil)
	}

	// 各メタデータエントリの検証
	for key, value := range metadata {
		// キーの検証
		if key == "" {
			return NewCommentValidationError("metadata", key,
				"empty metadata key",
				nil)
		}

		if utf8.RuneCountInString(key) > v.config.MaxMetadataKeyLength {
			return NewCommentValidationError("metadata", key,
				fmt.Sprintf("metadata key too long (max: %d)", v.config.MaxMetadataKeyLength),
				nil)
		}

		// 値の検証
		if utf8.RuneCountInString(value) > v.config.MaxMetadataValueLength {
			return NewCommentValidationError("metadata", value,
				fmt.Sprintf("metadata value too long for key '%s' (max: %d)", key, v.config.MaxMetadataValueLength),
				nil)
		}

		// キーの文字パターンチェック
		if err := v.validateCharPattern(key, fmt.Sprintf("metadata[%s].key", key)); err != nil {
			return err
		}

		// 値の文字パターンチェック
		if err := v.validateCharPattern(value, fmt.Sprintf("metadata[%s].value", key)); err != nil {
			return err
		}

		// キーの禁止語句チェック
		if err := v.validateForbiddenWords(key, fmt.Sprintf("metadata[%s].key", key)); err != nil {
			return err
		}

		// 値の禁止語句チェック
		if err := v.validateForbiddenWords(value, fmt.Sprintf("metadata[%s].value", key)); err != nil {
			return err
		}

		// キーのSQLインジェクションチェック
		if err := v.validateSQLInjection(key, fmt.Sprintf("metadata[%s].key", key)); err != nil {
			return err
		}

		// 値のSQLインジェクションチェック
		if err := v.validateSQLInjection(value, fmt.Sprintf("metadata[%s].value", key)); err != nil {
			return err
		}
	}

	return nil
}

// validateCharPattern 文字パターンの検証
func (v *DefaultCommentValidator) validateCharPattern(text, field string) error {
	if v.allowedCharRegex == nil || text == "" {
		return nil
	}

	if !v.allowedCharRegex.MatchString(text) {
		return NewCommentValidationError(field, text,
			"contains invalid characters",
			ErrInvalidCharacters)
	}

	return nil
}

// validateForbiddenWords 禁止語句の検証
func (v *DefaultCommentValidator) validateForbiddenWords(text, field string) error {
	if len(v.config.ForbiddenWords) == 0 || text == "" {
		return nil
	}

	upperText := strings.ToUpper(text)
	for _, word := range v.config.ForbiddenWords {
		if strings.Contains(upperText, strings.ToUpper(word)) {
			return NewCommentValidationError(field, text,
				fmt.Sprintf("contains forbidden word: %s", word),
				ErrInvalidCharacters)
		}
	}

	return nil
}

// validateSQLInjection SQLインジェクションの検証
func (v *DefaultCommentValidator) validateSQLInjection(text, field string) error {
	if !v.config.EnableSQLInjectionCheck || v.sqlInjectionRegex == nil || text == "" {
		return nil
	}

	if v.sqlInjectionRegex.MatchString(text) {
		return NewCommentValidationError(field, text,
			"potentially contains SQL injection patterns",
			ErrInvalidCharacters)
	}

	return nil
}

// sanitizeString 文字列のサニタイゼーション
func (v *DefaultCommentValidator) sanitizeString(text string) string {
	if text == "" {
		return text
	}

	// 前後の空白を除去
	text = strings.TrimSpace(text)

	// 制御文字を除去
	text = v.removeControlChars(text)

	// HTMLエスケープ
	if v.config.EnableHTMLEscape {
		text = html.EscapeString(text)
	}

	// 連続する空白を正規化
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	return text
}

// removeControlChars 制御文字を除去
func (v *DefaultCommentValidator) removeControlChars(text string) string {
	var result strings.Builder
	for _, r := range text {
		// 印刷可能文字またはタブ・改行・復帰のみ許可
		if unicode.IsPrint(r) || r == '\t' || r == '\n' || r == '\r' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// GetConfig 設定を取得
func (v *DefaultCommentValidator) GetConfig() *ValidationConfig {
	return v.config
}

// SetConfig 設定を更新
func (v *DefaultCommentValidator) SetConfig(config *ValidationConfig) {
	v.config = config
	
	// 正規表現を再コンパイル
	if config.AllowedCharPattern != "" {
		v.allowedCharRegex = regexp.MustCompile(config.AllowedCharPattern)
	} else {
		v.allowedCharRegex = nil
	}

	if config.EnableSQLInjectionCheck {
		sqlPattern := `(?i)(union\s+select|or\s+1\s*=\s*1|and\s+1\s*=\s*1|'|\-\-|\/\*|\*\/|xp_|sp_|exec|execute|drop\s+table|delete\s+from|insert\s+into|update\s+set)`
		v.sqlInjectionRegex = regexp.MustCompile(sqlPattern)
	} else {
		v.sqlInjectionRegex = nil
	}
}

// StrictValidationConfig より厳格な検証設定
func StrictValidationConfig() *ValidationConfig {
	return &ValidationConfig{
		MaxLogicalNameLength:     50,
		MaxDescriptionLength:     500,
		MaxTagCount:             10,
		MaxTagLength:            30,
		MaxMetadataCount:        20,
		MaxMetadataKeyLength:    50,
		MaxMetadataValueLength:  200,
		AllowedCharPattern:      `^[\p{L}\p{N}\p{Zs}\-_.,()[\]{}:;'"!?]+$`, // より制限的な文字セット
		ForbiddenWords:          []string{"DROP", "DELETE", "INSERT", "UPDATE", "EXEC", "SCRIPT", "UNION", "SELECT", "FROM", "WHERE"},
		EnableHTMLEscape:        true,
		EnableSQLInjectionCheck: true,
	}
}

// PermissiveValidationConfig より緩い検証設定
func PermissiveValidationConfig() *ValidationConfig {
	return &ValidationConfig{
		MaxLogicalNameLength:     200,
		MaxDescriptionLength:     2000,
		MaxTagCount:             50,
		MaxTagLength:            100,
		MaxMetadataCount:        100,
		MaxMetadataKeyLength:    200,
		MaxMetadataValueLength:  1000,
		AllowedCharPattern:      "", // パターンチェックなし
		ForbiddenWords:          []string{}, // 禁止語句なし
		EnableHTMLEscape:        false,
		EnableSQLInjectionCheck: false,
	}
}