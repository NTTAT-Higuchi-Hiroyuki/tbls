package schema

import (
	"fmt"
)

// EnhancedCommentDriverAdapter 拡張コメント処理対応ドライバーアダプター
type EnhancedCommentDriverAdapter struct {
	processor CommentProcessor
	config    *ProcessingConfig
}

// NewEnhancedCommentDriverAdapter 新しいドライバーアダプターを作成
func NewEnhancedCommentDriverAdapter(config *ProcessingConfig) *EnhancedCommentDriverAdapter {
	if config == nil {
		config = &ProcessingConfig{
			EnableValidation:   true,
			EnableSanitization: false,
			DefaultDelimiter:   "|",
			FallbackToLegacy:   true,
			StrictMode:         false,
			ProcessingTimeout:  1000, // 1秒
		}
	}

	return &EnhancedCommentDriverAdapter{
		processor: NewEnhancedCommentProcessorWithConfig(config),
		config:    config,
	}
}

// ProcessSchemaComments スキーマ全体のコメントを拡張処理
func (adapter *EnhancedCommentDriverAdapter) ProcessSchemaComments(schema *Schema) error {
	if schema == nil {
		return fmt.Errorf("schema cannot be nil")
	}

	// スキーマレベルの拡張コメント処理
	// 注意: 現在のSchemaオブジェクトにはCommentフィールドが存在しないため、
	// 将来的にスキーマレベルのコメント処理が必要な場合のプレースホルダー

	// テーブルレベルの処理
	for _, table := range schema.Tables {
		if err := adapter.processTableComments(table); err != nil {
			return fmt.Errorf("table %s processing failed: %w", table.Name, err)
		}
	}

	return nil
}

// processTableComments テーブルとその関連オブジェクトのコメント処理
func (adapter *EnhancedCommentDriverAdapter) processTableComments(table *Table) error {
	delimiter := adapter.config.DefaultDelimiter
	fallback := adapter.config.FallbackToLegacy

	// テーブルコメント処理
	if err := table.ProcessEnhancedComment(adapter.processor, delimiter, fallback); err != nil {
		if adapter.config.StrictMode {
			return fmt.Errorf("table comment processing failed: %w", err)
		}
		// 非厳密モードではログ出力のみ（実装は呼び出し側に依存）
	}

	// カラムコメント処理
	for _, column := range table.Columns {
		if err := column.ProcessEnhancedComment(adapter.processor, delimiter, fallback); err != nil {
			if adapter.config.StrictMode {
				return fmt.Errorf("column %s comment processing failed: %w", column.Name, err)
			}
		}
	}

	// インデックスコメント処理
	for _, index := range table.Indexes {
		if err := index.ProcessEnhancedComment(adapter.processor, delimiter); err != nil {
			if adapter.config.StrictMode {
				return fmt.Errorf("index %s comment processing failed: %w", index.Name, err)
			}
		}
	}

	// 制約コメント処理
	for _, constraint := range table.Constraints {
		if err := constraint.ProcessEnhancedComment(adapter.processor, delimiter); err != nil {
			if adapter.config.StrictMode {
				return fmt.Errorf("constraint %s comment processing failed: %w", constraint.Name, err)
			}
		}
	}

	// トリガーコメント処理
	for _, trigger := range table.Triggers {
		if err := trigger.ProcessEnhancedComment(adapter.processor, delimiter); err != nil {
			if adapter.config.StrictMode {
				return fmt.Errorf("trigger %s comment processing failed: %w", trigger.Name, err)
			}
		}
	}

	return nil
}

// GetProcessingStatistics 処理統計情報を取得
func (adapter *EnhancedCommentDriverAdapter) GetProcessingStatistics(schema *Schema) *ProcessingStatistics {
	stats := &ProcessingStatistics{
		TotalTables:     len(schema.Tables),
		ProcessedTables: 0,
		TotalColumns:    0,
		ProcessedColumns: 0,
		TotalIndexes:    0,
		ProcessedIndexes: 0,
		TotalConstraints: 0,
		ProcessedConstraints: 0,
		TotalTriggers:   0,
		ProcessedTriggers: 0,
		ProcessingErrors: make([]string, 0),
	}

	for _, table := range schema.Tables {
		// テーブル統計
		if table.HasEnhancedComment() {
			stats.ProcessedTables++
		}

		// カラム統計
		stats.TotalColumns += len(table.Columns)
		for _, column := range table.Columns {
			if column.HasEnhancedComment() {
				stats.ProcessedColumns++
			}
		}

		// インデックス統計
		stats.TotalIndexes += len(table.Indexes)
		for _, index := range table.Indexes {
			if index.HasEnhancedComment() {
				stats.ProcessedIndexes++
			}
		}

		// 制約統計
		stats.TotalConstraints += len(table.Constraints)
		for _, constraint := range table.Constraints {
			if constraint.HasEnhancedComment() {
				stats.ProcessedConstraints++
			}
		}

		// トリガー統計
		stats.TotalTriggers += len(table.Triggers)
		for _, trigger := range table.Triggers {
			if trigger.HasEnhancedComment() {
				stats.ProcessedTriggers++
			}
		}
	}

	return stats
}

// ProcessingStatistics 処理統計情報
type ProcessingStatistics struct {
	TotalTables         int      `json:"total_tables"`
	ProcessedTables     int      `json:"processed_tables"`
	TotalColumns        int      `json:"total_columns"`
	ProcessedColumns    int      `json:"processed_columns"`
	TotalIndexes        int      `json:"total_indexes"`
	ProcessedIndexes    int      `json:"processed_indexes"`
	TotalConstraints    int      `json:"total_constraints"`
	ProcessedConstraints int     `json:"processed_constraints"`
	TotalTriggers       int      `json:"total_triggers"`
	ProcessedTriggers   int      `json:"processed_triggers"`
	ProcessingErrors    []string `json:"processing_errors"`
}

// GetProcessingRate 処理率を計算
func (stats *ProcessingStatistics) GetProcessingRate() float64 {
	totalObjects := stats.TotalTables + stats.TotalColumns + stats.TotalIndexes + stats.TotalConstraints + stats.TotalTriggers
	processedObjects := stats.ProcessedTables + stats.ProcessedColumns + stats.ProcessedIndexes + stats.ProcessedConstraints + stats.ProcessedTriggers

	if totalObjects == 0 {
		return 0.0
	}

	return float64(processedObjects) / float64(totalObjects) * 100.0
}

// GetSummary 処理サマリーを取得
func (stats *ProcessingStatistics) GetSummary() string {
	return fmt.Sprintf(
		"拡張コメント処理統計: テーブル %d/%d (%.1f%%), カラム %d/%d (%.1f%%), インデックス %d/%d (%.1f%%), 制約 %d/%d (%.1f%%), トリガー %d/%d (%.1f%%), 全体処理率: %.1f%%",
		stats.ProcessedTables, stats.TotalTables, stats.getTableRate(),
		stats.ProcessedColumns, stats.TotalColumns, stats.getColumnRate(),
		stats.ProcessedIndexes, stats.TotalIndexes, stats.getIndexRate(),
		stats.ProcessedConstraints, stats.TotalConstraints, stats.getConstraintRate(),
		stats.ProcessedTriggers, stats.TotalTriggers, stats.getTriggerRate(),
		stats.GetProcessingRate(),
	)
}

func (stats *ProcessingStatistics) getTableRate() float64 {
	if stats.TotalTables == 0 {
		return 0.0
	}
	return float64(stats.ProcessedTables) / float64(stats.TotalTables) * 100.0
}

func (stats *ProcessingStatistics) getColumnRate() float64 {
	if stats.TotalColumns == 0 {
		return 0.0
	}
	return float64(stats.ProcessedColumns) / float64(stats.TotalColumns) * 100.0
}

func (stats *ProcessingStatistics) getIndexRate() float64 {
	if stats.TotalIndexes == 0 {
		return 0.0
	}
	return float64(stats.ProcessedIndexes) / float64(stats.TotalIndexes) * 100.0
}

func (stats *ProcessingStatistics) getConstraintRate() float64 {
	if stats.TotalConstraints == 0 {
		return 0.0
	}
	return float64(stats.ProcessedConstraints) / float64(stats.TotalConstraints) * 100.0
}

func (stats *ProcessingStatistics) getTriggerRate() float64 {
	if stats.TotalTriggers == 0 {
		return 0.0
	}
	return float64(stats.ProcessedTriggers) / float64(stats.TotalTriggers) * 100.0
}

// DriverIntegrationHelper ドライバー統合ヘルパー関数群
type DriverIntegrationHelper struct{}

// NewDriverIntegrationHelper 新しいドライバー統合ヘルパーを作成
func NewDriverIntegrationHelper() *DriverIntegrationHelper {
	return &DriverIntegrationHelper{}
}

// WrapAnalyzeWithEnhancedComments ドライバーのAnalyzeメソッドに拡張コメント処理を追加
func (helper *DriverIntegrationHelper) WrapAnalyzeWithEnhancedComments(
	originalAnalyze func(*Schema) error,
	config *ProcessingConfig,
) func(*Schema) error {
	return func(schema *Schema) error {
		// 元のドライバー解析を実行
		if err := originalAnalyze(schema); err != nil {
			return err
		}

		// 拡張コメント処理を追加
		adapter := NewEnhancedCommentDriverAdapter(config)
		return adapter.ProcessSchemaComments(schema)
	}
}

// CreateMigrationScript 拡張コメント対応移行スクリプト作成
func (helper *DriverIntegrationHelper) CreateMigrationScript(schema *Schema, targetDriver string) (string, error) {
	// 将来的な機能: 既存のコメントを拡張コメント形式に移行するスクリプト生成
	// 現在はプレースホルダー実装
	return fmt.Sprintf("-- Migration script for %s\n-- TODO: Implement migration logic\n", targetDriver), nil
}

// ValidateDriverCompatibility ドライバー互換性検証
func (helper *DriverIntegrationHelper) ValidateDriverCompatibility(driverName string, schema *Schema) []string {
	issues := make([]string, 0)

	// ドライバー固有の互換性チェック
	switch driverName {
	case "sqlite":
		// SQLiteはテーブルコメントをサポートしていない
		for _, table := range schema.Tables {
			if table.Comment != "" {
				issues = append(issues, fmt.Sprintf("SQLite does not support table comments: table %s", table.Name))
			}
		}

	case "mysql":
		// MySQLの文字エンコーディングチェック
		for _, table := range schema.Tables {
			for _, column := range table.Columns {
				if len(column.Comment) > 1024 {
					issues = append(issues, fmt.Sprintf("MySQL comment length limit exceeded: table %s, column %s", table.Name, column.Name))
				}
			}
		}

	case "postgres":
		// PostgreSQLの場合は特別な制限は少ない
		break

	default:
		issues = append(issues, fmt.Sprintf("Unknown driver: %s", driverName))
	}

	return issues
}

// ConvertCommentFormat コメント形式変換
func (helper *DriverIntegrationHelper) ConvertCommentFormat(comment string, fromFormat, toFormat string) (string, error) {
	// プレースホルダー実装
	// 実際の実装では、Legacy -> JSON -> YAML の相互変換を行う
	return comment, fmt.Errorf("comment format conversion not yet implemented")
}