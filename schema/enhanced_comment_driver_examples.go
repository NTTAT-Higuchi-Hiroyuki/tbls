package schema

import (
	"database/sql"
	"fmt"
	"log"
)

// ExampleDriverUsage ドライバーでの拡張コメント処理の使用例
//
// この関数は、実際のドライバーで拡張コメント処理を使用する方法を示します。
// 各データベースドライバーは、このパターンに従って拡張コメント処理を統合できます。
func ExampleDriverUsage() {
	// PostgreSQLの例
	examplePostgreSQLDriver()

	// SQLiteの例
	exampleSQLiteDriver()

	// 汎用的な使用例
	exampleGenericDriverIntegration()
}

// examplePostgreSQLDriver PostgreSQLドライバーでの使用例
func examplePostgreSQLDriver() {
	// 接続設定（実際の使用では環境変数から取得）
	dsn := "postgres://user:password@localhost/dbname?sslmode=disable"
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("PostgreSQL connection failed: %v", err)
		return
	}
	defer db.Close()

	// スキーマオブジェクト作成
	schema := &Schema{Name: "example_database"}

	// 従来のドライバー解析を模擬（実際には postgres.New(db).Analyze(schema)）
	simulateDriverAnalysis(schema, "postgres")

	// 拡張コメント処理の統合
	config := &ProcessingConfig{
		EnableValidation:   true,
		EnableSanitization: false,
		DefaultDelimiter:   "|",
		FallbackToLegacy:   true,
		StrictMode:         false,
	}

	adapter := NewEnhancedCommentDriverAdapter(config)
	
	// 拡張コメント処理実行
	err = adapter.ProcessSchemaComments(schema)
	if err != nil {
		log.Printf("Enhanced comment processing failed: %v", err)
		return
	}

	// 処理結果の統計表示
	stats := adapter.GetProcessingStatistics(schema)
	log.Printf("PostgreSQL processing completed: %s", stats.GetSummary())

	// 処理されたスキーマの利用例
	displayProcessedSchema(schema, "PostgreSQL")
}

// exampleSQLiteDriver SQLiteドライバーでの使用例
func exampleSQLiteDriver() {
	// インメモリSQLiteの例
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Printf("SQLite connection failed: %v", err)
		return
	}
	defer db.Close()

	// テストテーブル作成
	createSQLiteTestTables(db)

	// スキーマオブジェクト作成
	schema := &Schema{Name: "sqlite_example"}

	// 従来のドライバー解析を模擬
	simulateDriverAnalysis(schema, "sqlite")

	// SQLite用の設定（テーブルコメントは無効）
	config := &ProcessingConfig{
		EnableValidation:   true,
		EnableSanitization: true,  // SQLiteでは安全性を重視
		DefaultDelimiter:   "|",
		FallbackToLegacy:   true,
		StrictMode:         false, // SQLiteは制限が多いため非厳密モード
	}

	adapter := NewEnhancedCommentDriverAdapter(config)
	
	// 拡張コメント処理実行
	err = adapter.ProcessSchemaComments(schema)
	if err != nil {
		log.Printf("Enhanced comment processing failed: %v", err)
		return
	}

	// SQLite特有の互換性チェック
	helper := NewDriverIntegrationHelper()
	issues := helper.ValidateDriverCompatibility("sqlite", schema)
	if len(issues) > 0 {
		log.Printf("SQLite compatibility issues found: %v", issues)
	}

	// 処理結果の統計表示
	stats := adapter.GetProcessingStatistics(schema)
	log.Printf("SQLite processing completed: %s", stats.GetSummary())

	displayProcessedSchema(schema, "SQLite")
}

// exampleGenericDriverIntegration 汎用的なドライバー統合例
func exampleGenericDriverIntegration() {
	// 任意のドライバーでの使用を想定した汎用例

	// カスタム解析関数（実際のドライバーのAnalyzeメソッド）
	customAnalyze := func(schema *Schema) error {
		// 模擬的なテーブル/カラム追加
		table := &Table{
			Name:    "custom_table",
			Comment: `{"name": "カスタムテーブル", "description": "汎用ドライバー統合例", "tags": ["custom", "example"]}`,
			Columns: []*Column{
				{
					Name:    "id",
					Comment: `{"name": "ID", "description": "主キー", "tags": ["pk", "auto_increment"]}`,
				},
				{
					Name:    "name",
					Comment: "名前|エンティティの名前",
				},
				{
					Name:    "metadata",
					Comment: "description: メタデータ\nformat: json\ntags:\n  - optional\n  - flexible",
				},
			},
		}
		schema.Tables = append(schema.Tables, table)
		return nil
	}

	// ドライバー統合ヘルパーを使用
	helper := NewDriverIntegrationHelper()
	
	// 拡張コメント処理付きの解析関数を作成
	enhancedAnalyze := helper.WrapAnalyzeWithEnhancedComments(customAnalyze, nil)

	// スキーマ解析実行
	schema := &Schema{Name: "generic_example"}
	err := enhancedAnalyze(schema)
	if err != nil {
		log.Printf("Generic driver integration failed: %v", err)
		return
	}

	log.Printf("Generic driver integration completed successfully")
	displayProcessedSchema(schema, "Generic")
}

// simulateDriverAnalysis ドライバー解析の模擬
func simulateDriverAnalysis(schema *Schema, driverType string) {
	switch driverType {
	case "postgres":
		simulatePostgreSQLAnalysis(schema)
	case "sqlite":
		simulateSQLiteAnalysis(schema)
	default:
		simulateGenericAnalysis(schema)
	}
}

// simulatePostgreSQLAnalysis PostgreSQL解析の模擬
func simulatePostgreSQLAnalysis(schema *Schema) {
	// PostgreSQLの特徴を模擬したテーブル作成
	tables := []*Table{
		{
			Name:    "users",
			Comment: `{"name": "ユーザー", "description": "システムユーザー管理テーブル", "tags": ["master", "auth"], "priority": 1}`,
			Columns: []*Column{
				{
					Name:    "id",
					Comment: `{"name": "ユーザーID", "description": "一意識別子", "tags": ["pk", "serial"]}`,
				},
				{
					Name:    "username",
					Comment: "ユーザー名|ログイン用の一意ユーザー名",
				},
				{
					Name:    "email",
					Comment: "description: メールアドレス\nvalidation: email\ntags:\n  - contact\n  - unique",
				},
				{
					Name:    "profile",
					Comment: `{"name": "プロフィール", "description": "JSON形式のプロフィールデータ", "tags": ["json", "optional"]}`,
				},
			},
			Indexes: []*Index{
				{
					Name:    "idx_users_username",
					Comment: `{"description": "ユーザー名一意インデックス", "tags": ["unique", "performance"]}`,
				},
				{
					Name:    "idx_users_email",
					Comment: "メールアドレス検索用インデックス",
				},
			},
			Constraints: []*Constraint{
				{
					Name:    "pk_users",
					Comment: `{"description": "ユーザーテーブル主キー", "tags": ["integrity"]}`,
				},
			},
		},
		{
			Name:    "posts",
			Comment: "投稿|ユーザーの投稿を管理するテーブル",
			Columns: []*Column{
				{
					Name:    "id",
					Comment: `{"name": "投稿ID", "description": "投稿の一意識別子", "tags": ["pk"]}`,
				},
				{
					Name:    "user_id",
					Comment: `{"name": "ユーザーID", "description": "投稿者のID", "tags": ["fk", "required"]}`,
				},
				{
					Name:    "title",
					Comment: "タイトル|投稿のタイトル",
				},
				{
					Name:    "content",
					Comment: "description: 投稿内容\nformat: text\ntags:\n  - content\n  - searchable",
				},
			},
		},
	}

	schema.Tables = append(schema.Tables, tables...)
}

// simulateSQLiteAnalysis SQLite解析の模擬
func simulateSQLiteAnalysis(schema *Schema) {
	// SQLiteの特徴（テーブルコメント制限）を考慮
	tables := []*Table{
		{
			Name:    "simple_table",
			Comment: "", // SQLiteはテーブルコメントをサポートしない
			Columns: []*Column{
				{
					Name:    "id",
					Comment: "", // SQLiteはカラムコメントも制限的
				},
				{
					Name:    "data",
					Comment: "", // アプリケーション側で管理されたコメントがあると仮定
				},
			},
		},
	}

	schema.Tables = append(schema.Tables, tables...)
}

// simulateGenericAnalysis 汎用解析の模擬
func simulateGenericAnalysis(schema *Schema) {
	table := &Table{
		Name:    "generic_table",
		Comment: "汎用テーブル|データベース非依存のテーブル",
		Columns: []*Column{
			{
				Name:    "common_id",
				Comment: `{"name": "共通ID", "description": "汎用識別子"}`,
			},
		},
	}

	schema.Tables = append(schema.Tables, table)
}

// createSQLiteTestTables SQLiteテストテーブル作成
func createSQLiteTestTables(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS example_users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS example_posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT,
			FOREIGN KEY (user_id) REFERENCES example_users(id)
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Failed to create SQLite table: %v", err)
		}
	}
}

// displayProcessedSchema 処理済みスキーマの表示
func displayProcessedSchema(schema *Schema, driverName string) {
	log.Printf("=== %s スキーマ処理結果 ===", driverName)
	log.Printf("データベース名: %s", schema.Name)
	log.Printf("テーブル数: %d", len(schema.Tables))

	for _, table := range schema.Tables {
		log.Printf("  テーブル: %s", table.Name)
		
		if table.HasEnhancedComment() {
			log.Printf("    論理名: %s", table.LogicalName)
			log.Printf("    説明: %s", table.GetDescription())
			tags := table.GetTags()
			if len(tags) > 0 {
				log.Printf("    タグ: %v", tags)
			}
			priority := table.GetPriority()
			if priority > 0 {
				log.Printf("    優先度: %d", priority)
			}
		} else if table.Comment != "" {
			log.Printf("    コメント: %s", table.Comment)
		}

		// カラム情報の表示
		for _, column := range table.Columns {
			log.Printf("    カラム: %s", column.Name)
			
			if column.HasEnhancedComment() {
				if column.LogicalName != "" {
					log.Printf("      論理名: %s", column.LogicalName)
				}
				if column.GetDescription() != "" {
					log.Printf("      説明: %s", column.GetDescription())
				}
				tags := column.GetTags()
				if len(tags) > 0 {
					log.Printf("      タグ: %v", tags)
				}
			} else if column.Comment != "" {
				log.Printf("      コメント: %s", column.Comment)
			}
		}

		// インデックス情報の表示
		if len(table.Indexes) > 0 {
			log.Printf("    インデックス数: %d", len(table.Indexes))
			for _, index := range table.Indexes {
				if index.HasEnhancedComment() && index.GetDescription() != "" {
					log.Printf("      %s: %s", index.Name, index.GetDescription())
				}
			}
		}

		// 制約情報の表示
		if len(table.Constraints) > 0 {
			log.Printf("    制約数: %d", len(table.Constraints))
			for _, constraint := range table.Constraints {
				if constraint.HasEnhancedComment() && constraint.GetDescription() != "" {
					log.Printf("      %s: %s", constraint.Name, constraint.GetDescription())
				}
			}
		}
	}

	log.Printf("=== %s 処理完了 ===\n", driverName)
}

// DriverEnhancedCommentConfiguration ドライバー拡張コメント設定
type DriverEnhancedCommentConfiguration struct {
	DriverName           string            `json:"driver_name"`
	EnableEnhancedComments bool            `json:"enable_enhanced_comments"`
	ProcessingConfig     *ProcessingConfig `json:"processing_config"`
	CompatibilityMode    bool              `json:"compatibility_mode"`
	CustomDelimiter      string            `json:"custom_delimiter"`
	ErrorHandling        string            `json:"error_handling"` // "strict", "permissive", "silent"
}

// NewDriverConfiguration 新しいドライバー設定を作成
func NewDriverConfiguration(driverName string) *DriverEnhancedCommentConfiguration {
	config := &DriverEnhancedCommentConfiguration{
		DriverName:           driverName,
		EnableEnhancedComments: true,
		CompatibilityMode:    true,
		CustomDelimiter:      "|",
		ErrorHandling:        "permissive",
	}

	// ドライバー固有のデフォルト設定
	switch driverName {
	case "sqlite":
		config.ProcessingConfig = &ProcessingConfig{
			EnableValidation:   true,
			EnableSanitization: true,
			DefaultDelimiter:   "|",
			FallbackToLegacy:   true,
			StrictMode:         false,
		}
		config.CompatibilityMode = true // SQLiteは制限が多いため互換モード有効

	case "postgres":
		config.ProcessingConfig = &ProcessingConfig{
			EnableValidation:   true,
			EnableSanitization: false,
			DefaultDelimiter:   "|",
			FallbackToLegacy:   true,
			StrictMode:         false,
		}

	case "mysql":
		config.ProcessingConfig = &ProcessingConfig{
			EnableValidation:   true,
			EnableSanitization: true, // MySQLの文字エンコーディング問題対策
			DefaultDelimiter:   "|",
			FallbackToLegacy:   true,
			StrictMode:         false,
		}

	default:
		config.ProcessingConfig = &ProcessingConfig{
			EnableValidation:   true,
			EnableSanitization: false,
			DefaultDelimiter:   "|",
			FallbackToLegacy:   true,
			StrictMode:         false,
		}
	}

	return config
}

// ApplyConfiguration 設定をドライバーに適用
func (config *DriverEnhancedCommentConfiguration) ApplyConfiguration() *EnhancedCommentDriverAdapter {
	if !config.EnableEnhancedComments {
		return nil
	}

	// カスタムデリミターの適用
	if config.CustomDelimiter != "" {
		config.ProcessingConfig.DefaultDelimiter = config.CustomDelimiter
	}

	// エラーハンドリング設定の適用
	switch config.ErrorHandling {
	case "strict":
		config.ProcessingConfig.StrictMode = true
	case "permissive":
		config.ProcessingConfig.StrictMode = false
	case "silent":
		config.ProcessingConfig.StrictMode = false
		// 将来的にサイレントモードの実装を追加
	}

	return NewEnhancedCommentDriverAdapter(config.ProcessingConfig)
}

// ValidateConfiguration 設定の妥当性検証
func (config *DriverEnhancedCommentConfiguration) ValidateConfiguration() []string {
	issues := make([]string, 0)

	if config.DriverName == "" {
		issues = append(issues, "Driver name is required")
	}

	if config.ProcessingConfig == nil {
		issues = append(issues, "Processing config is required")
		return issues // これ以上の検証は不可能
	}

	if config.CustomDelimiter == "" {
		issues = append(issues, "Custom delimiter cannot be empty")
	}

	if config.ErrorHandling != "strict" && config.ErrorHandling != "permissive" && config.ErrorHandling != "silent" {
		issues = append(issues, fmt.Sprintf("Invalid error handling mode: %s", config.ErrorHandling))
	}

	// ドライバー固有の検証
	switch config.DriverName {
	case "sqlite":
		if config.ProcessingConfig.StrictMode {
			issues = append(issues, "SQLite driver should not use strict mode due to limited comment support")
		}

	case "mysql":
		if !config.ProcessingConfig.EnableSanitization {
			issues = append(issues, "MySQL driver should enable sanitization for character encoding safety")
		}
	}

	return issues
}

// GetRecommendedConfiguration 推奨設定の取得
func GetRecommendedConfiguration(driverName string, useCase string) *DriverEnhancedCommentConfiguration {
	config := NewDriverConfiguration(driverName)

	// 用途別の調整
	switch useCase {
	case "production":
		config.ProcessingConfig.EnableValidation = true
		config.ProcessingConfig.EnableSanitization = true
		config.ErrorHandling = "permissive"
		config.CompatibilityMode = true

	case "development":
		config.ProcessingConfig.EnableValidation = true
		config.ProcessingConfig.EnableSanitization = false
		config.ErrorHandling = "strict"
		config.CompatibilityMode = false

	case "testing":
		config.ProcessingConfig.EnableValidation = true
		config.ProcessingConfig.EnableSanitization = false
		config.ErrorHandling = "strict"
		config.CompatibilityMode = false

	default:
		// デフォルト設定をそのまま使用
	}

	return config
}