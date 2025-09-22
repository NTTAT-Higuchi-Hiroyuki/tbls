package schema

import (
	"strings"
	"unicode/utf8"

	"github.com/k1LoW/errors"
)

// CommentParser はコメントから論理名とクリーンなコメントを分離するためのインターフェース
type CommentParser interface {
	// ParseComment は rawComment を separator で分割し、論理名とクリーンなコメントを返す
	ParseComment(rawComment, separator string) (logicalName, cleanComment string)

	// ExtractLogicalName は rawComment から論理名部分のみを抽出する
	ExtractLogicalName(rawComment, separator string) string

	// ExtractCleanComment は rawComment からクリーンなコメント部分のみを抽出する
	ExtractCleanComment(rawComment, separator string) string

	// HasLogicalName は rawComment に separator で分離された論理名が含まれているかを判定する
	HasLogicalName(rawComment, separator string) bool
}

// DefaultCommentParser はデフォルトのコメント解析実装
type DefaultCommentParser struct{}

// NewDefaultCommentParser は新しい DefaultCommentParser を作成する
func NewDefaultCommentParser() *DefaultCommentParser {
	return &DefaultCommentParser{}
}

// ParseComment は rawComment を separator で分割し、論理名とクリーンなコメントを返す
func (p *DefaultCommentParser) ParseComment(rawComment, separator string) (logicalName, cleanComment string) {
	// 入力検証
	if rawComment == "" || separator == "" {
		return "", rawComment
	}

	// UTF-8安全な文字列処理で分割
	parts := p.splitComment(rawComment, separator)

	if len(parts) < 2 {
		// 区切り文字が見つからない場合は論理名なし
		return "", rawComment
	}

	// 最初の区切り文字で分割（複数の区切り文字がある場合）
	logicalName = strings.TrimSpace(parts[0])

	// 残りの部分を結合（区切り文字が複数ある場合の対応）
	if len(parts) > 2 {
		cleanComment = strings.TrimSpace(strings.Join(parts[1:], separator))
	} else {
		cleanComment = strings.TrimSpace(parts[1])
	}

	return logicalName, cleanComment
}

// ExtractLogicalName は rawComment から論理名部分のみを抽出する
func (p *DefaultCommentParser) ExtractLogicalName(rawComment, separator string) string {
	logicalName, _ := p.ParseComment(rawComment, separator)
	return logicalName
}

// ExtractCleanComment は rawComment からクリーンなコメント部分のみを抽出する
func (p *DefaultCommentParser) ExtractCleanComment(rawComment, separator string) string {
	_, cleanComment := p.ParseComment(rawComment, separator)
	return cleanComment
}

// HasLogicalName は rawComment に separator で分離された論理名が含まれているかを判定する
func (p *DefaultCommentParser) HasLogicalName(rawComment, separator string) bool {
	if rawComment == "" || separator == "" {
		return false
	}

	logicalName, _ := p.ParseComment(rawComment, separator)
	return logicalName != ""
}

// splitComment はUTF-8安全な文字列分割を行う
func (p *DefaultCommentParser) splitComment(comment, separator string) []string {
	// UTF-8文字列として正しく処理
	if !utf8.ValidString(comment) || !utf8.ValidString(separator) {
		// 無効なUTF-8の場合はそのまま返す
		return []string{comment}
	}

	return strings.Split(comment, separator)
}

// ParseCommentWithParser は指定されたパーサーを使用してコメントを解析する
// パーサーがnilの場合はデフォルトパーサーを使用する
func ParseCommentWithParser(rawComment, separator string, parser CommentParser) (logicalName, cleanComment string) {
	if parser == nil {
		parser = NewDefaultCommentParser()
	}
	return parser.ParseComment(rawComment, separator)
}

// ParseComment はデフォルトパーサーを使用してコメントを解析する便利関数
func ParseComment(rawComment, separator string) (logicalName, cleanComment string) {
	parser := NewDefaultCommentParser()
	return parser.ParseComment(rawComment, separator)
}

// ExtractLogicalName はデフォルトパーサーを使用して論理名を抽出する便利関数
func ExtractLogicalName(rawComment, separator string) string {
	parser := NewDefaultCommentParser()
	return parser.ExtractLogicalName(rawComment, separator)
}

// ExtractCleanComment はデフォルトパーサーを使用してクリーンなコメントを抽出する便利関数
func ExtractCleanComment(rawComment, separator string) string {
	parser := NewDefaultCommentParser()
	return parser.ExtractCleanComment(rawComment, separator)
}

// ValidateCommentSeparator は区切り文字の妥当性をチェックする
func ValidateCommentSeparator(separator string) error {
	if separator == "" {
		return errors.New("comment separator cannot be empty")
	}

	if !utf8.ValidString(separator) {
		return errors.New("comment separator must be valid UTF-8")
	}

	// 制御文字やタブ、改行文字をチェック
	for _, r := range separator {
		if r < 32 && r != 9 { // 制御文字（タブ以外）はエラー
			return errors.New("comment separator cannot contain control characters")
		}
	}

	return nil
}
