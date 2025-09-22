package schema

import (
	"fmt"
	"strings"
)

// LogicalNameProcessor は論理名処理のためのユーティリティ構造体
type LogicalNameProcessor struct {
	separator string
	parser    CommentParser
}

// NewLogicalNameProcessor は新しいLogicalNameProcessorを作成する
func NewLogicalNameProcessor(separator string) *LogicalNameProcessor {
	return &LogicalNameProcessor{
		separator: separator,
		parser:    NewDefaultCommentParser(),
	}
}

// ProcessSchema はスキーマ全体に論理名処理を適用する
func (p *LogicalNameProcessor) ProcessSchema(schema *Schema) error {
	if schema == nil {
		return fmt.Errorf("schema is nil")
	}

	// テーブルの論理名処理
	for _, table := range schema.Tables {
		if err := p.ProcessTable(table); err != nil {
			return fmt.Errorf("failed to process table %s: %w", table.Name, err)
		}
	}

	return nil
}

// ProcessTable はテーブルとその関連オブジェクトに論理名処理を適用する
func (p *LogicalNameProcessor) ProcessTable(table *Table) error {
	if table == nil {
		return fmt.Errorf("table is nil")
	}

	// テーブル自体の論理名処理
	table.SetLogicalNameFromComment(p.separator)

	// カラムの論理名処理
	for _, column := range table.Columns {
		column.SetLogicalNameFromComment(p.separator)
	}

	// インデックスの論理名処理
	for _, index := range table.Indexes {
		index.SetLogicalNameFromComment(p.separator)
	}

	// 制約の論理名処理
	for _, constraint := range table.Constraints {
		constraint.SetLogicalNameFromComment(p.separator)
	}

	// トリガーの論理名処理
	for _, trigger := range table.Triggers {
		trigger.SetLogicalNameFromComment(p.separator)
	}

	return nil
}

// LogicalNameDisplayController は論理名表示制御のためのユーティリティ構造体
type LogicalNameDisplayController struct {
	showLogicalName    bool
	fallbackToPhysical bool
	format             string
}

// NewLogicalNameDisplayController は新しいLogicalNameDisplayControllerを作成する
func NewLogicalNameDisplayController(showLogicalName, fallbackToPhysical bool) *LogicalNameDisplayController {
	return &LogicalNameDisplayController{
		showLogicalName:    showLogicalName,
		fallbackToPhysical: fallbackToPhysical,
		format:             "%s", // デフォルトフォーマット
	}
}

// SetFormat は表示フォーマットを設定する
// %s: 論理名または物理名
// %l: 論理名のみ
// %p: 物理名のみ
func (c *LogicalNameDisplayController) SetFormat(format string) {
	if format != "" {
		c.format = format
	}
}

// GetDisplayName はオブジェクトの表示名を取得する
func (c *LogicalNameDisplayController) GetDisplayName(logicalName, physicalName string) string {
	if !c.showLogicalName {
		return physicalName
	}

	targetName := logicalName
	if targetName == "" && c.fallbackToPhysical {
		targetName = physicalName
	}

	if targetName == "" {
		return physicalName
	}

	// フォーマットに応じた表示
	switch c.format {
	case "%l":
		if logicalName != "" {
			return logicalName
		}
		if c.fallbackToPhysical {
			return physicalName
		}
		return ""
	case "%p":
		return physicalName
	case "%l (%p)":
		if logicalName != "" {
			return fmt.Sprintf("%s (%s)", logicalName, physicalName)
		}
		return physicalName
	case "%p (%l)":
		if logicalName != "" {
			return fmt.Sprintf("%s (%s)", physicalName, logicalName)
		}
		return physicalName
	default:
		return targetName
	}
}

// GetTableDisplayName はテーブルの表示名を取得する
func (c *LogicalNameDisplayController) GetTableDisplayName(table *Table) string {
	if table == nil {
		return ""
	}
	return c.GetDisplayName(table.LogicalName, table.Name)
}

// GetColumnDisplayName はカラムの表示名を取得する
func (c *LogicalNameDisplayController) GetColumnDisplayName(column *Column) string {
	if column == nil {
		return ""
	}
	return c.GetDisplayName(column.LogicalName, column.Name)
}

// GetIndexDisplayName はインデックスの表示名を取得する
func (c *LogicalNameDisplayController) GetIndexDisplayName(index *Index) string {
	if index == nil {
		return ""
	}
	return c.GetDisplayName(index.LogicalName, index.Name)
}

// GetConstraintDisplayName は制約の表示名を取得する
func (c *LogicalNameDisplayController) GetConstraintDisplayName(constraint *Constraint) string {
	if constraint == nil {
		return ""
	}
	return c.GetDisplayName(constraint.LogicalName, constraint.Name)
}

// GetTriggerDisplayName はトリガーの表示名を取得する
func (c *LogicalNameDisplayController) GetTriggerDisplayName(trigger *Trigger) string {
	if trigger == nil {
		return ""
	}
	return c.GetDisplayName(trigger.LogicalName, trigger.Name)
}

// LogicalNameValidator は論理名の妥当性検証を行う
type LogicalNameValidator struct {
	maxLength      int
	allowedChars   string
	forbiddenWords []string
	caseSensitive  bool
}

// NewLogicalNameValidator は新しいLogicalNameValidatorを作成する
func NewLogicalNameValidator() *LogicalNameValidator {
	return &LogicalNameValidator{
		maxLength:      255,
		allowedChars:   "", // 空の場合は制限なし
		forbiddenWords: []string{},
		caseSensitive:  true,
	}
}

// SetMaxLength は最大長を設定する
func (v *LogicalNameValidator) SetMaxLength(length int) {
	if length > 0 {
		v.maxLength = length
	}
}

// SetAllowedChars は許可文字を設定する
func (v *LogicalNameValidator) SetAllowedChars(chars string) {
	v.allowedChars = chars
}

// AddForbiddenWord は禁止単語を追加する
func (v *LogicalNameValidator) AddForbiddenWord(word string) {
	if word != "" {
		v.forbiddenWords = append(v.forbiddenWords, word)
	}
}

// SetCaseSensitive は大文字小文字の区別を設定する
func (v *LogicalNameValidator) SetCaseSensitive(sensitive bool) {
	v.caseSensitive = sensitive
}

// Validate は論理名の妥当性を検証する
func (v *LogicalNameValidator) Validate(logicalName string) error {
	if logicalName == "" {
		return nil // 空の論理名は有効
	}

	// 長さチェック
	if len(logicalName) > v.maxLength {
		return fmt.Errorf("logical name too long: %d characters (max: %d)", len(logicalName), v.maxLength)
	}

	// 許可文字チェック
	if v.allowedChars != "" {
		for _, char := range logicalName {
			if !strings.ContainsRune(v.allowedChars, char) {
				return fmt.Errorf("logical name contains forbidden character: %c", char)
			}
		}
	}

	// 禁止単語チェック
	targetName := logicalName
	if !v.caseSensitive {
		targetName = strings.ToLower(targetName)
	}

	for _, word := range v.forbiddenWords {
		checkWord := word
		if !v.caseSensitive {
			checkWord = strings.ToLower(checkWord)
		}
		if strings.Contains(targetName, checkWord) {
			return fmt.Errorf("logical name contains forbidden word: %s", word)
		}
	}

	return nil
}

// LogicalNameHelper は論理名処理のヘルパー関数を提供する
type LogicalNameHelper struct{}

// NewLogicalNameHelper は新しいLogicalNameHelperを作成する
func NewLogicalNameHelper() *LogicalNameHelper {
	return &LogicalNameHelper{}
}

// HasLogicalName はオブジェクトが論理名を持つかどうかを判定する
func (h *LogicalNameHelper) HasLogicalName(logicalName string) bool {
	return strings.TrimSpace(logicalName) != ""
}

// IsLogicalNameAvailable はコメントから論理名が取得可能かどうかを判定する
func (h *LogicalNameHelper) IsLogicalNameAvailable(comment, separator string) bool {
	if comment == "" || separator == "" {
		return false
	}
	parser := NewDefaultCommentParser()
	return parser.HasLogicalName(comment, separator)
}

// CleanLogicalName は論理名をクリーンアップする
func (h *LogicalNameHelper) CleanLogicalName(logicalName string) string {
	// 前後の空白を削除
	cleaned := strings.TrimSpace(logicalName)

	// 連続する空白を単一の空白に変換
	cleaned = strings.Join(strings.Fields(cleaned), " ")

	return cleaned
}

// NormalizeLogicalName は論理名を正規化する
func (h *LogicalNameHelper) NormalizeLogicalName(logicalName string) string {
	// クリーンアップ
	normalized := h.CleanLogicalName(logicalName)

	// 制御文字を削除
	var result strings.Builder
	for _, r := range normalized {
		if r >= 32 && r != 127 { // 印刷可能文字のみ保持
			result.WriteRune(r)
		}
	}

	return result.String()
}

// CompareLogicalNames は論理名を比較する
func (h *LogicalNameHelper) CompareLogicalNames(name1, name2 string, caseSensitive bool) bool {
	if !caseSensitive {
		return strings.EqualFold(name1, name2)
	}
	return name1 == name2
}

// GetPreferredName は論理名と物理名から優先される名前を取得する
func (h *LogicalNameHelper) GetPreferredName(logicalName, physicalName string, preferLogical bool) string {
	if preferLogical && h.HasLogicalName(logicalName) {
		return logicalName
	}
	if physicalName != "" {
		return physicalName
	}
	return logicalName // 物理名も空の場合は論理名を返す
}

// LogicalNameStatistics は論理名の統計情報を保持する
type LogicalNameStatistics struct {
	TotalTables            int
	TablesWithLogical      int
	TotalColumns           int
	ColumnsWithLogical     int
	TotalIndexes           int
	IndexesWithLogical     int
	TotalConstraints       int
	ConstraintsWithLogical int
	TotalTriggers          int
	TriggersWithLogical    int
}

// CalculateLogicalNameStatistics はスキーマの論理名統計を計算する
func CalculateLogicalNameStatistics(schema *Schema) *LogicalNameStatistics {
	stats := &LogicalNameStatistics{}
	helper := NewLogicalNameHelper()

	if schema == nil {
		return stats
	}

	for _, table := range schema.Tables {
		stats.TotalTables++
		if helper.HasLogicalName(table.LogicalName) {
			stats.TablesWithLogical++
		}

		for _, column := range table.Columns {
			stats.TotalColumns++
			if helper.HasLogicalName(column.LogicalName) {
				stats.ColumnsWithLogical++
			}
		}

		for _, index := range table.Indexes {
			stats.TotalIndexes++
			if helper.HasLogicalName(index.LogicalName) {
				stats.IndexesWithLogical++
			}
		}

		for _, constraint := range table.Constraints {
			stats.TotalConstraints++
			if helper.HasLogicalName(constraint.LogicalName) {
				stats.ConstraintsWithLogical++
			}
		}

		for _, trigger := range table.Triggers {
			stats.TotalTriggers++
			if helper.HasLogicalName(trigger.LogicalName) {
				stats.TriggersWithLogical++
			}
		}
	}

	return stats
}

// GetCoveragePercentage は論理名のカバレッジ率を取得する
func (s *LogicalNameStatistics) GetCoveragePercentage() map[string]float64 {
	coverage := make(map[string]float64)

	if s.TotalTables > 0 {
		coverage["tables"] = float64(s.TablesWithLogical) / float64(s.TotalTables) * 100
	}
	if s.TotalColumns > 0 {
		coverage["columns"] = float64(s.ColumnsWithLogical) / float64(s.TotalColumns) * 100
	}
	if s.TotalIndexes > 0 {
		coverage["indexes"] = float64(s.IndexesWithLogical) / float64(s.TotalIndexes) * 100
	}
	if s.TotalConstraints > 0 {
		coverage["constraints"] = float64(s.ConstraintsWithLogical) / float64(s.TotalConstraints) * 100
	}
	if s.TotalTriggers > 0 {
		coverage["triggers"] = float64(s.TriggersWithLogical) / float64(s.TotalTriggers) * 100
	}

	return coverage
}

// GetTotalCoveragePercentage は全体のカバレッジ率を取得する
func (s *LogicalNameStatistics) GetTotalCoveragePercentage() float64 {
	total := s.TotalTables + s.TotalColumns + s.TotalIndexes + s.TotalConstraints + s.TotalTriggers
	withLogical := s.TablesWithLogical + s.ColumnsWithLogical + s.IndexesWithLogical + s.ConstraintsWithLogical + s.TriggersWithLogical

	if total == 0 {
		return 0
	}

	return float64(withLogical) / float64(total) * 100
}
