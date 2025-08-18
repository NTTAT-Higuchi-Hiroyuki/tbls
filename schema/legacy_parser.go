package schema

import (
	"strings"
	"unicode"
)

// LegacyParser 従来の区切り文字ベースのコメント解析を行うパーサー
type LegacyParser struct {
	name     string
	priority int
}

// NewLegacyParser 新しいLegacyParserを作成
func NewLegacyParser() *LegacyParser {
	return &LegacyParser{
		name:     "legacy",
		priority: 1000, // 低優先度（フォールバック用）
	}
}

// Name パーサーの名前を返す
func (p *LegacyParser) Name() string {
	return p.name
}

// Priority パーサーの優先度を返す
func (p *LegacyParser) Priority() int {
	return p.priority
}

// SetPriority パーサーの優先度を設定
func (p *LegacyParser) SetPriority(priority int) {
	p.priority = priority
}

// CanParse このパーサーでコメントを解析可能かを判定
// LegacyParserは全てのコメントを解析可能（フォールバック）
func (p *LegacyParser) CanParse(comment string) bool {
	return true
}

// ParseComment コメントを解析してCommentDataに変換
func (p *LegacyParser) ParseComment(comment string, delimiter string) (*CommentData, error) {
	if comment == "" {
		return &CommentData{Source: comment}, nil
	}
	
	// デフォルト区切り文字
	if delimiter == "" {
		delimiter = "|"
	}
	
	// コメントの正規化
	normalizedComment := p.normalizeComment(comment)
	
	// 区切り文字での分割
	parts := p.splitWithEscape(normalizedComment, delimiter)
	
	result := &CommentData{
		Source: comment,
	}
	
	if len(parts) > 0 {
		result.LogicalName = strings.TrimSpace(parts[0])
	}
	
	if len(parts) > 1 {
		result.Description = strings.TrimSpace(parts[1])
	}
	
	return result, nil
}

// normalizeComment コメント文字列を正規化
func (p *LegacyParser) normalizeComment(comment string) string {
	// 前後の空白を除去
	comment = strings.TrimSpace(comment)
	
	// 連続する空白文字を単一の空白に変換
	comment = p.normalizeWhitespace(comment)
	
	return comment
}

// normalizeWhitespace 空白文字を正規化
func (p *LegacyParser) normalizeWhitespace(s string) string {
	var result strings.Builder
	var prevIsSpace bool
	
	for _, r := range s {
		isSpace := unicode.IsSpace(r)
		if isSpace {
			if !prevIsSpace {
				result.WriteRune(' ')
			}
		} else {
			result.WriteRune(r)
		}
		prevIsSpace = isSpace
	}
	
	// 前後の空白を除去
	return strings.TrimSpace(result.String())
}

// splitWithEscape エスケープ処理を考慮した文字列分割
func (p *LegacyParser) splitWithEscape(s string, delimiter string) []string {
	if delimiter == "" || s == "" {
		return []string{s}
	}
	
	var parts []string
	var current strings.Builder
	var escaped bool
	
	for i := 0; i < len(s); {
		if escaped {
			// エスケープされた文字をそのまま追加
			current.WriteByte(s[i])
			escaped = false
			i++
			continue
		}
		
		// エスケープ文字をチェック
		if s[i] == '\\' && i+1 < len(s) {
			// 次の文字が区切り文字の場合のみエスケープとして扱う
			if strings.HasPrefix(s[i+1:], delimiter) {
				escaped = true
				i++
				continue
			}
		}
		
		// 区切り文字をチェック
		if strings.HasPrefix(s[i:], delimiter) {
			// 現在の部分を追加
			parts = append(parts, current.String())
			current.Reset()
			i += len(delimiter)
			continue
		}
		
		// 通常の文字を追加
		current.WriteByte(s[i])
		i++
	}
	
	// 最後の部分を追加
	parts = append(parts, current.String())
	
	return parts
}

// LegacyCommentSplitter 既存のSplitComment機能との互換性を保つための関数
// 新しい実装では非推奨だが、既存コードとの互換性のために残す
func LegacyCommentSplitter(comment string, delimiter string) SplitCommentParts {
	parser := NewLegacyParser()
	data, err := parser.ParseComment(comment, delimiter)
	if err != nil || data == nil {
		return SplitCommentParts{}
	}
	
	return SplitCommentParts{
		LogicalName: data.LogicalName,
		Comment:     data.Description,
	}
}

// EnhancedSplitComment 拡張されたコメント分割機能
// 従来のSplitCommentより高機能で、CommentDataを返す
func EnhancedSplitComment(comment string, delimiter string) *CommentData {
	parser := NewLegacyParser()
	data, err := parser.ParseComment(comment, delimiter)
	if err != nil {
		return &CommentData{
			Source: comment,
		}
	}
	return data
}

// ExtractLogicalNameEnhanced 拡張された論理名抽出機能
func ExtractLogicalNameEnhanced(comment string, delimiter string, physicalName string, fallbackToName bool) string {
	data := EnhancedSplitComment(comment, delimiter)
	if data != nil && data.LogicalName != "" {
		return data.LogicalName
	}
	
	if fallbackToName {
		return physicalName
	}
	
	return ""
}

// ExtractCleanCommentEnhanced 拡張されたクリーンなコメント抽出機能
func ExtractCleanCommentEnhanced(comment string, delimiter string) string {
	data := EnhancedSplitComment(comment, delimiter)
	if data != nil && data.Description != "" {
		return data.Description
	}
	return ""
}