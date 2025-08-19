package schema

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSONParser JSON形式のコメントを解析するパーサー
type JSONParser struct {
	name        string
	priority    int
	maxDepth    int
	maxSize     int
}

// NewJSONParser 新しいJSONParserを作成
func NewJSONParser() *JSONParser {
	return &JSONParser{
		name:     "json",
		priority: 10, // 高優先度（LegacyParserより先に試行）
		maxDepth: 5,  // JSON構造の最大深度
		maxSize:  8192, // JSON文字列の最大サイズ（バイト）
	}
}

// NewJSONParserWithLimits 制限を指定してJSONParserを作成
func NewJSONParserWithLimits(maxDepth, maxSize int) *JSONParser {
	parser := NewJSONParser()
	parser.maxDepth = maxDepth
	parser.maxSize = maxSize
	return parser
}

// Name パーサーの名前を返す
func (p *JSONParser) Name() string {
	return p.name
}

// Priority パーサーの優先度を返す
func (p *JSONParser) Priority() int {
	return p.priority
}

// SetPriority パーサーの優先度を設定
func (p *JSONParser) SetPriority(priority int) {
	p.priority = priority
}

// CanParse このパーサーでコメントを解析可能かを判定
func (p *JSONParser) CanParse(comment string) bool {
	if comment == "" {
		return false
	}
	
	// サイズ制限チェック
	if len(comment) > p.maxSize {
		return false
	}
	
	// 基本的なJSON形式チェック（最初と最後の文字）
	comment = strings.TrimSpace(comment)
	if len(comment) < 2 {
		return false
	}
	
	// JSON形式の基本的なチェック
	hasValidBraces := (comment[0] == '{' && comment[len(comment)-1] == '}') ||
		              (comment[0] == '[' && comment[len(comment)-1] == ']')
	
	if !hasValidBraces {
		return false
	}
	
	// 実際にJSONとして解析可能かチェック
	var rawData interface{}
	err := json.Unmarshal([]byte(comment), &rawData)
	return err == nil
}

// ParseComment コメントを解析してCommentDataに変換
func (p *JSONParser) ParseComment(comment string, delimiter string) (*CommentData, error) {
	if comment == "" {
		return &CommentData{Source: comment}, nil
	}
	
	// 前処理
	comment = strings.TrimSpace(comment)
	
	// サイズ制限チェック
	if len(comment) > p.maxSize {
		return nil, NewCommentParseError(p.name, comment, "comment too large", ErrCommentTooLong)
	}
	
	// 基本的な形式チェック
	if !p.CanParse(comment) {
		return nil, NewCommentParseError(p.name, comment, "not a valid JSON format", ErrInvalidCommentFormat)
	}
	
	// JSON解析
	var rawData interface{}
	if err := json.Unmarshal([]byte(comment), &rawData); err != nil {
		return nil, NewCommentParseError(p.name, comment, "JSON unmarshal failed", err)
	}
	
	// 深度チェック
	if p.getDepth(rawData) > p.maxDepth {
		return nil, NewCommentParseError(p.name, comment, "JSON structure too deep", ErrJSONTooDeep)
	}
	
	// CommentDataに変換
	result, err := p.convertToCommentData(rawData)
	if err != nil {
		return nil, NewCommentParseError(p.name, comment, "failed to convert to CommentData", err)
	}
	
	result.Source = comment
	return result, nil
}

// convertToCommentData 解析されたJSONデータをCommentDataに変換
func (p *JSONParser) convertToCommentData(data interface{}) (*CommentData, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		return p.convertMapToCommentData(v)
	case []interface{}:
		// 配列の場合は最初の要素をオブジェクトとして扱う
		if len(v) > 0 {
			if mapData, ok := v[0].(map[string]interface{}); ok {
				return p.convertMapToCommentData(mapData)
			}
		}
		return &CommentData{}, nil
	default:
		return &CommentData{}, nil
	}
}

// convertMapToCommentData マップデータをCommentDataに変換
func (p *JSONParser) convertMapToCommentData(data map[string]interface{}) (*CommentData, error) {
	result := &CommentData{}
	
	// 論理名の抽出（複数のキー名に対応）
	logicalNameKeys := []string{"name", "logical_name", "logicalName", "title", "label"}
	for _, key := range logicalNameKeys {
		if value, exists := data[key]; exists {
			if str, ok := value.(string); ok && str != "" {
				result.LogicalName = str
				break
			}
		}
	}
	
	// 説明の抽出（複数のキー名に対応）
	descriptionKeys := []string{"description", "desc", "comment", "note", "summary"}
	for _, key := range descriptionKeys {
		if value, exists := data[key]; exists {
			if str, ok := value.(string); ok && str != "" {
				result.Description = str
				break
			}
		}
	}
	
	// タグの抽出
	if value, exists := data["tags"]; exists {
		tags, err := p.convertToStringSlice(value)
		if err == nil {
			result.Tags = tags
		}
	}
	
	// 優先度の抽出
	if value, exists := data["priority"]; exists {
		if num, ok := value.(float64); ok {
			result.Priority = int(num)
		}
	}
	
	// 非推奨フラグの抽出
	if value, exists := data["deprecated"]; exists {
		if deprecated, ok := value.(bool); ok {
			result.Deprecated = deprecated
		}
	}
	
	// メタデータの抽出（その他のフィールド）
	result.Metadata = make(map[string]string)
	excludeKeys := map[string]bool{
		"name": true, "logical_name": true, "logicalName": true, "title": true, "label": true,
		"description": true, "desc": true, "comment": true, "note": true, "summary": true,
		"tags": true, "priority": true, "deprecated": true,
	}
	
	for key, value := range data {
		if !excludeKeys[key] {
			if str, ok := value.(string); ok {
				result.Metadata[key] = str
			} else {
				// 文字列以外の値は文字列に変換
				result.Metadata[key] = fmt.Sprintf("%v", value)
			}
		}
	}
	
	return result, nil
}

// convertToStringSlice インターフェースを文字列スライスに変換
func (p *JSONParser) convertToStringSlice(value interface{}) ([]string, error) {
	switch v := value.(type) {
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			} else {
				result = append(result, fmt.Sprintf("%v", item))
			}
		}
		return result, nil
	case []string:
		return v, nil
	case string:
		// 単一の文字列をスライスとして扱う
		return []string{v}, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to string slice", value)
	}
}

// getDepth JSON構造の深度を取得
func (p *JSONParser) getDepth(data interface{}) int {
	switch v := data.(type) {
	case map[string]interface{}:
		maxDepth := 0
		for _, value := range v {
			depth := p.getDepth(value)
			if depth > maxDepth {
				maxDepth = depth
			}
		}
		return maxDepth + 1
	case []interface{}:
		maxDepth := 0
		for _, value := range v {
			depth := p.getDepth(value)
			if depth > maxDepth {
				maxDepth = depth
			}
		}
		return maxDepth + 1
	default:
		return 1
	}
}

// IsValidJSON JSON文字列が有効かどうかを判定するユーティリティ関数
func IsValidJSON(comment string) bool {
	parser := NewJSONParser()
	return parser.CanParse(comment)
}

// QuickParseJSON JSON文字列を素早く解析するユーティリティ関数
func QuickParseJSON(comment string) (*CommentData, error) {
	parser := NewJSONParser()
	return parser.ParseComment(comment, "")
}