package schema

import (
	"errors"
	"fmt"
)

// CommentParser はコメント解析を行うインターフェース
type CommentParser interface {
	// ParseComment コメントを解析してCommentDataに変換
	ParseComment(comment string, delimiter string) (*CommentData, error)
	// CanParse このパーサーが指定されたコメントを解析可能かを判定
	CanParse(comment string) bool
	// Priority パーサーの優先度（小さいほど高優先度）
	Priority() int
	// Name パーサーの名前
	Name() string
}

// ParserRegistry はコメントパーサーを管理するインターフェース
type ParserRegistry interface {
	// RegisterParser パーサーを登録
	RegisterParser(parser CommentParser)
	// GetParsers 登録済みのパーサー一覧を取得（優先度順）
	GetParsers() []CommentParser
	// ParseWithFallback フォールバック付きでコメントを解析
	ParseWithFallback(comment string, delimiter string) (*CommentData, error)
	// Clear 全てのパーサーをクリア
	Clear()
}

// CommentParseError コメント解析エラー
type CommentParseError struct {
	ParserName string
	Comment    string
	Reason     string
	Cause      error
}

func (e *CommentParseError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("comment parse error in %s: %s (comment: %q, cause: %v)", 
			e.ParserName, e.Reason, e.Comment, e.Cause)
	}
	return fmt.Sprintf("comment parse error in %s: %s (comment: %q)", 
		e.ParserName, e.Reason, e.Comment)
}

func (e *CommentParseError) Unwrap() error {
	return e.Cause
}

// NewCommentParseError CommentParseErrorを作成
func NewCommentParseError(parserName, comment, reason string, cause error) *CommentParseError {
	return &CommentParseError{
		ParserName: parserName,
		Comment:    comment,
		Reason:     reason,
		Cause:      cause,
	}
}

// CommentValidationError コメント検証エラー
type CommentValidationError struct {
	Field   string
	Value   string
	Reason  string
	Cause   error
}

func (e *CommentValidationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("comment validation error for %s: %s (value: %q, cause: %v)", 
			e.Field, e.Reason, e.Value, e.Cause)
	}
	return fmt.Sprintf("comment validation error for %s: %s (value: %q)", 
		e.Field, e.Reason, e.Value)
}

func (e *CommentValidationError) Unwrap() error {
	return e.Cause
}

// NewCommentValidationError CommentValidationErrorを作成
func NewCommentValidationError(field, value, reason string, cause error) *CommentValidationError {
	return &CommentValidationError{
		Field:  field,
		Value:  value,
		Reason: reason,
		Cause:  cause,
	}
}

// 事前定義されたエラー
var (
	// ErrInvalidCommentFormat 不正なコメント形式
	ErrInvalidCommentFormat = errors.New("invalid comment format")
	// ErrUnsupportedFormat サポートされていない形式
	ErrUnsupportedFormat = errors.New("unsupported comment format")
	// ErrCommentTooLong コメントが長すぎる
	ErrCommentTooLong = errors.New("comment too long")
	// ErrLogicalNameTooLong 論理名が長すぎる
	ErrLogicalNameTooLong = errors.New("logical name too long")
	// ErrDescriptionTooLong 説明が長すぎる
	ErrDescriptionTooLong = errors.New("description too long")
	// ErrInvalidCharacters 無効な文字が含まれている
	ErrInvalidCharacters = errors.New("invalid characters")
	// ErrJSONTooDeep JSON構造が深すぎる
	ErrJSONTooDeep = errors.New("JSON structure too deep")
	// ErrYAMLParseError YAML解析エラー
	ErrYAMLParseError = errors.New("YAML parse error")
)

// DefaultParserRegistry デフォルトのパーサーレジストリ実装
type DefaultParserRegistry struct {
	parsers []CommentParser
}

// NewDefaultParserRegistry 新しいDefaultParserRegistryを作成
func NewDefaultParserRegistry() *DefaultParserRegistry {
	return &DefaultParserRegistry{
		parsers: make([]CommentParser, 0),
	}
}

// RegisterParser パーサーを登録（優先度順でソート）
func (r *DefaultParserRegistry) RegisterParser(parser CommentParser) {
	if parser == nil {
		return
	}
	
	// 優先度順に挿入
	inserted := false
	for i, p := range r.parsers {
		if parser.Priority() < p.Priority() {
			// i番目に挿入
			r.parsers = append(r.parsers[:i], append([]CommentParser{parser}, r.parsers[i:]...)...)
			inserted = true
			break
		}
	}
	
	if !inserted {
		r.parsers = append(r.parsers, parser)
	}
}

// GetParsers 登録済みのパーサー一覧を取得
func (r *DefaultParserRegistry) GetParsers() []CommentParser {
	// スライスのコピーを返す（外部から変更されないように）
	result := make([]CommentParser, len(r.parsers))
	copy(result, r.parsers)
	return result
}

// ParseWithFallback フォールバック付きでコメントを解析
func (r *DefaultParserRegistry) ParseWithFallback(comment string, delimiter string) (*CommentData, error) {
	if comment == "" {
		return &CommentData{Source: comment}, nil
	}
	
	var lastError error
	
	// 優先度順にパーサーを試行
	for _, parser := range r.parsers {
		if parser.CanParse(comment) {
			result, err := parser.ParseComment(comment, delimiter)
			if err == nil && result != nil {
				result.Source = comment
				return result, nil
			}
			// エラーの場合は記録して次のパーサーを試行
			lastError = err
		}
	}
	
	// 全てのパーサーで解析に失敗した場合
	if lastError != nil {
		return nil, fmt.Errorf("all parsers failed to parse comment: %w", lastError)
	}
	
	// 解析可能なパーサーが見つからなかった場合
	return nil, NewCommentParseError("registry", comment, "no suitable parser found", ErrUnsupportedFormat)
}

// Clear 全てのパーサーをクリア
func (r *DefaultParserRegistry) Clear() {
	r.parsers = r.parsers[:0]
}

// Count 登録されているパーサーの数を取得
func (r *DefaultParserRegistry) Count() int {
	return len(r.parsers)
}

// HasParser 指定された名前のパーサーが登録されているかを確認
func (r *DefaultParserRegistry) HasParser(name string) bool {
	for _, parser := range r.parsers {
		if parser.Name() == name {
			return true
		}
	}
	return false
}

// GetParser 指定された名前のパーサーを取得
func (r *DefaultParserRegistry) GetParser(name string) CommentParser {
	for _, parser := range r.parsers {
		if parser.Name() == name {
			return parser
		}
	}
	return nil
}