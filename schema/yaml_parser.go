package schema

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// YAMLParser YAML形式のコメントを解析するパーサー
type YAMLParser struct {
	name        string
	priority    int
	maxDepth    int
	maxSize     int
}

// NewYAMLParser 新しいYAMLParserを作成
func NewYAMLParser() *YAMLParser {
	return &YAMLParser{
		name:     "yaml",
		priority: 15, // JSONParserより少し低い優先度
		maxDepth: 5,  // YAML構造の最大深度
		maxSize:  8192, // YAML文字列の最大サイズ（バイト）
	}
}

// NewYAMLParserWithLimits 制限を指定してYAMLParserを作成
func NewYAMLParserWithLimits(maxDepth, maxSize int) *YAMLParser {
	parser := NewYAMLParser()
	parser.maxDepth = maxDepth
	parser.maxSize = maxSize
	return parser
}

// Name パーサーの名前を返す
func (p *YAMLParser) Name() string {
	return p.name
}

// Priority パーサーの優先度を返す
func (p *YAMLParser) Priority() int {
	return p.priority
}

// SetPriority パーサーの優先度を設定
func (p *YAMLParser) SetPriority(priority int) {
	p.priority = priority
}

// CanParse このパーサーでコメントを解析可能かを判定
func (p *YAMLParser) CanParse(comment string) bool {
	if comment == "" {
		return false
	}
	
	// サイズ制限チェック
	if len(comment) > p.maxSize {
		return false
	}
	
	comment = strings.TrimSpace(comment)
	if len(comment) == 0 {
		return false
	}
	
	// YAML形式の基本的な特徴をチェック
	if p.looksLikeYAML(comment) {
		// 実際にYAMLとして解析可能かチェック
		var data interface{}
		err := yaml.Unmarshal([]byte(comment), &data)
		return err == nil
	}
	
	return false
}

// looksLikeYAML 文字列がYAML形式に見えるかを判定
func (p *YAMLParser) looksLikeYAML(comment string) bool {
	// YAML特有のパターンをチェック
	yamlPatterns := []string{
		`^\s*\w+\s*:\s*.+`,           // key: value
		`^\s*-\s+\w+\s*:\s*.+`,       // - key: value (list item)
		`^\s*\w+\s*:\s*$`,            // key: (empty value)
		`^\s*-\s+.*`,                 // - item (list)
		`^\s*\|\s*$`,                 // | (literal block)
		`^\s*>\s*$`,                  // > (folded block)
	}
	
	lines := strings.Split(comment, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue // 空行やコメント行はスキップ
		}
		
		for _, pattern := range yamlPatterns {
			matched, _ := regexp.MatchString(pattern, line)
			if matched {
				return true
			}
		}
	}
	
	return false
}

// ParseComment コメントを解析してCommentDataに変換
func (p *YAMLParser) ParseComment(comment string, delimiter string) (*CommentData, error) {
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
		return nil, NewCommentParseError(p.name, comment, "not a valid YAML format", ErrInvalidCommentFormat)
	}
	
	// YAML解析
	var rawData interface{}
	if err := yaml.Unmarshal([]byte(comment), &rawData); err != nil {
		return nil, NewCommentParseError(p.name, comment, "YAML unmarshal failed", err)
	}
	
	// 深度チェック
	if p.getDepth(rawData) > p.maxDepth {
		return nil, NewCommentParseError(p.name, comment, "YAML structure too deep", ErrYAMLParseError)
	}
	
	// CommentDataに変換
	result, err := p.convertToCommentData(rawData)
	if err != nil {
		return nil, NewCommentParseError(p.name, comment, "failed to convert to CommentData", err)
	}
	
	result.Source = comment
	return result, nil
}

// convertToCommentData 解析されたYAMLデータをCommentDataに変換
func (p *YAMLParser) convertToCommentData(data interface{}) (*CommentData, error) {
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
	case string:
		// 単純な文字列の場合は説明として扱う
		return &CommentData{Description: v}, nil
	default:
		return &CommentData{}, nil
	}
}

// convertMapToCommentData マップデータをCommentDataに変換
func (p *YAMLParser) convertMapToCommentData(data map[string]interface{}) (*CommentData, error) {
	result := &CommentData{}
	
	// 論理名の抽出（複数のキー名に対応）
	logicalNameKeys := []string{"name", "logical_name", "logicalName", "title", "label", "display_name"}
	for _, key := range logicalNameKeys {
		if value, exists := data[key]; exists {
			if str, ok := value.(string); ok && str != "" {
				result.LogicalName = str
				break
			}
		}
	}
	
	// 説明の抽出（複数のキー名に対応）
	descriptionKeys := []string{"description", "desc", "comment", "note", "summary", "details"}
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
		switch v := value.(type) {
		case int:
			result.Priority = v
		case float64:
			result.Priority = int(v)
		case string:
			// 文字列の場合は数値に変換を試行
			if priority, err := parseIntFromString(v); err == nil {
				result.Priority = priority
			}
		}
	}
	
	// 非推奨フラグの抽出
	if value, exists := data["deprecated"]; exists {
		switch v := value.(type) {
		case bool:
			result.Deprecated = v
		case string:
			result.Deprecated = strings.ToLower(v) == "true" || v == "1" || strings.ToLower(v) == "yes"
		}
	}
	
	// メタデータの抽出（その他のフィールド）
	result.Metadata = make(map[string]string)
	excludeKeys := map[string]bool{
		"name": true, "logical_name": true, "logicalName": true, "title": true, "label": true, "display_name": true,
		"description": true, "desc": true, "comment": true, "note": true, "summary": true, "details": true,
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
func (p *YAMLParser) convertToStringSlice(value interface{}) ([]string, error) {
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
		// 単一の文字列をスライスとして扱う、またはカンマ区切りで分割
		if strings.Contains(v, ",") {
			parts := strings.Split(v, ",")
			result := make([]string, len(parts))
			for i, part := range parts {
				result[i] = strings.TrimSpace(part)
			}
			return result, nil
		}
		return []string{v}, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to string slice", value)
	}
}

// getDepth YAML構造の深度を取得
func (p *YAMLParser) getDepth(data interface{}) int {
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

// parseIntFromString 文字列から整数を解析
func parseIntFromString(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}
	
	// 浮動小数点が含まれていないかチェック
	if strings.Contains(s, ".") {
		return 0, fmt.Errorf("contains decimal point")
	}
	
	var result int
	n, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, fmt.Errorf("failed to parse complete string")
	}
	
	return result, nil
}

// IsValidYAML YAML文字列が有効かどうかを判定するユーティリティ関数
func IsValidYAML(comment string) bool {
	parser := NewYAMLParser()
	return parser.CanParse(comment)
}

// QuickParseYAML YAML文字列を素早く解析するユーティリティ関数
func QuickParseYAML(comment string) (*CommentData, error) {
	parser := NewYAMLParser()
	return parser.ParseComment(comment, "")
}

// SupportedYAMLFormats サポートされているYAML形式の例を返す
func SupportedYAMLFormats() []string {
	return []string{
		"name: ユーザー名\ndescription: ユーザーの表示名",
		"- name: テーブル名\n  description: テーブルの説明",
		"logical_name: 論理名\ntags:\n  - 重要\n  - 公開\npriority: 1",
		"title: タイトル\ndetails: |\n  複数行の\n  詳細説明",
	}
}