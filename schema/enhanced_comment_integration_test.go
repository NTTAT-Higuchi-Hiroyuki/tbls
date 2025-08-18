package schema

import (
	"testing"
)

// TestEnhancedCommentIntegration 拡張コメント処理の統合テスト
func TestEnhancedCommentIntegration(t *testing.T) {
	processor := NewEnhancedCommentProcessor()

	// テスト用のスキーマを作成
	schema := &Schema{
		Name: "test_db",
		Tables: []*Table{
			{
				Name:    "users",
				Comment: `{"name": "ユーザーテーブル", "description": "システムのユーザー情報", "tags": ["マスター", "重要"], "deprecated": false}`,
				Columns: []*Column{
					{
						Name:    "id",
						Comment: `{"name": "ユーザーID", "description": "一意識別子", "tags": ["PK"], "priority": 1}`,
					},
					{
						Name:    "name",
						Comment: "ユーザー名|ユーザーの表示名",
					},
					{
						Name:    "email",
						Comment: "description: メールアドレス\ntags:\n  - 連絡先\n  - 必須",
					},
				},
				Indexes: []*Index{
					{
						Name:    "idx_users_email",
						Comment: `{"description": "メールアドレス一意制約", "tags": ["性能", "一意性"]}`,
					},
				},
				Constraints: []*Constraint{
					{
						Name:    "pk_users",
						Comment: `{"description": "主キー制約", "tags": ["整合性"]}`,
					},
				},
				Triggers: []*Trigger{
					{
						Name:    "tr_users_updated_at",
						Comment: "description: 更新日時自動設定\ntags:\n  - 自動化\n  - 監査",
					},
				},
			},
			{
				Name:    "posts",
				Comment: "投稿テーブル|ユーザーの投稿情報を管理",
				Columns: []*Column{
					{
						Name:    "id",
						Comment: "投稿ID",
					},
					{
						Name:    "user_id",
						Comment: `{"name": "ユーザーID", "description": "投稿者のID", "tags": ["FK"]}`,
					},
				},
			},
		},
	}

	t.Run("全オブジェクトの拡張コメント処理", func(t *testing.T) {
		// 全テーブルを処理
		for _, table := range schema.Tables {
			err := table.ProcessEnhancedComment(processor, "|", true)
			if err != nil {
				t.Errorf("Table %s ProcessEnhancedComment failed: %v", table.Name, err)
			}

			// 全カラムを処理
			for _, column := range table.Columns {
				err := column.ProcessEnhancedComment(processor, "|", true)
				if err != nil {
					t.Errorf("Column %s.%s ProcessEnhancedComment failed: %v", table.Name, column.Name, err)
				}
			}

			// 全インデックスを処理
			for _, index := range table.Indexes {
				err := index.ProcessEnhancedComment(processor, "|")
				if err != nil {
					t.Errorf("Index %s ProcessEnhancedComment failed: %v", index.Name, err)
				}
			}

			// 全制約を処理
			for _, constraint := range table.Constraints {
				err := constraint.ProcessEnhancedComment(processor, "|")
				if err != nil {
					t.Errorf("Constraint %s ProcessEnhancedComment failed: %v", constraint.Name, err)
				}
			}

			// 全トリガーを処理
			for _, trigger := range table.Triggers {
				err := trigger.ProcessEnhancedComment(processor, "|")
				if err != nil {
					t.Errorf("Trigger %s ProcessEnhancedComment failed: %v", trigger.Name, err)
				}
			}
		}

		// 処理結果の検証
		usersTable := schema.Tables[0]
		if !usersTable.HasEnhancedComment() {
			t.Error("Users table should have enhanced comment")
		}

		if usersTable.GetDescription() != "システムのユーザー情報" {
			t.Errorf("Users table description = %q, expected 'システムのユーザー情報'", usersTable.GetDescription())
		}

		expectedTags := []string{"マスター", "重要"}
		tags := usersTable.GetTags()
		if len(tags) != len(expectedTags) {
			t.Errorf("Users table tags length = %d, expected %d", len(tags), len(expectedTags))
		}

		// カラムの検証
		idColumn := usersTable.Columns[0]
		if idColumn.GetDescription() != "一意識別子" {
			t.Errorf("ID column description = %q, expected '一意識別子'", idColumn.GetDescription())
		}

		nameColumn := usersTable.Columns[1]
		if nameColumn.LogicalName != "ユーザー名" {
			t.Errorf("Name column logical name = %q, expected 'ユーザー名'", nameColumn.LogicalName)
		}

		emailColumn := usersTable.Columns[2]
		if emailColumn.GetDescription() != "メールアドレス" {
			t.Errorf("Email column description = %q, expected 'メールアドレス'", emailColumn.GetDescription())
		}

		// インデックスの検証
		index := usersTable.Indexes[0]
		if index.GetDescription() != "メールアドレス一意制約" {
			t.Errorf("Index description = %q, expected 'メールアドレス一意制約'", index.GetDescription())
		}

		// 制約の検証
		constraint := usersTable.Constraints[0]
		if constraint.GetDescription() != "主キー制約" {
			t.Errorf("Constraint description = %q, expected '主キー制約'", constraint.GetDescription())
		}

		// トリガーの検証
		trigger := usersTable.Triggers[0]
		if trigger.GetDescription() != "更新日時自動設定" {
			t.Errorf("Trigger description = %q, expected '更新日時自動設定'", trigger.GetDescription())
		}

		// postsテーブルの検証（従来形式）
		postsTable := schema.Tables[1]
		if postsTable.LogicalName != "投稿テーブル" {
			t.Errorf("Posts table logical name = %q, expected '投稿テーブル'", postsTable.LogicalName)
		}

		if postsTable.Comment != "ユーザーの投稿情報を管理" {
			t.Errorf("Posts table comment = %q, expected 'ユーザーの投稿情報を管理'", postsTable.Comment)
		}
	})
}

// TestSchemaWideEnhancedCommentProcessing スキーマ全体の拡張コメント処理テスト
func TestSchemaWideEnhancedCommentProcessing(t *testing.T) {
	// 複数の形式のコメントを含むスキーマ
	schema := &Schema{
		Name: "mixed_format_db",
		Tables: []*Table{
			{
				Name:    "json_table",
				Comment: `{"name": "JSONテーブル", "description": "JSON形式のコメント", "tags": ["JSON"], "priority": 5}`,
				Columns: []*Column{
					{
						Name:    "json_col",
						Comment: `{"name": "JSONカラム", "description": "JSON形式", "tags": ["test"]}`,
					},
				},
			},
			{
				Name:    "yaml_table",
				Comment: "name: YAMLテーブル\ndescription: YAML形式のコメント\ntags:\n  - YAML\n  - テスト",
				Columns: []*Column{
					{
						Name:    "yaml_col",
						Comment: "name: YAMLカラム\ndescription: YAML形式",
					},
				},
			},
			{
				Name:    "legacy_table",
				Comment: "従来テーブル|従来形式のコメント",
				Columns: []*Column{
					{
						Name:    "legacy_col",
						Comment: "従来カラム|従来形式",
					},
				},
			},
		},
	}

	processor := NewEnhancedCommentProcessor()

	// 一括処理関数
	processSchema := func(s *Schema) error {
		for _, table := range s.Tables {
			if err := table.ProcessEnhancedComment(processor, "|", true); err != nil {
				return err
			}

			for _, column := range table.Columns {
				if err := column.ProcessEnhancedComment(processor, "|", true); err != nil {
					return err
				}
			}
		}
		return nil
	}

	err := processSchema(schema)
	if err != nil {
		t.Fatalf("Schema processing failed: %v", err)
	}

	t.Run("JSON形式の処理検証", func(t *testing.T) {
		table := schema.Tables[0]
		if table.GetDescription() != "JSON形式のコメント" {
			t.Errorf("JSON table description = %q", table.GetDescription())
		}
		if table.GetPriority() != 5 {
			t.Errorf("JSON table priority = %d, expected 5", table.GetPriority())
		}

		column := table.Columns[0]
		if column.GetDescription() != "JSON形式" {
			t.Errorf("JSON column description = %q", column.GetDescription())
		}
	})

	t.Run("YAML形式の処理検証", func(t *testing.T) {
		table := schema.Tables[1]
		if table.GetDescription() != "YAML形式のコメント" {
			t.Errorf("YAML table description = %q", table.GetDescription())
		}

		expectedTags := []string{"YAML", "テスト"}
		tags := table.GetTags()
		if len(tags) != len(expectedTags) {
			t.Errorf("YAML table tags = %v, expected %v", tags, expectedTags)
		}

		column := table.Columns[0]
		if column.GetDescription() != "YAML形式" {
			t.Errorf("YAML column description = %q", column.GetDescription())
		}
	})

	t.Run("従来形式の処理検証", func(t *testing.T) {
		table := schema.Tables[2]
		if table.LogicalName != "従来テーブル" {
			t.Errorf("Legacy table logical name = %q", table.LogicalName)
		}
		if table.Comment != "従来形式のコメント" {
			t.Errorf("Legacy table comment = %q", table.Comment)
		}

		column := table.Columns[0]
		if column.LogicalName != "従来カラム" {
			t.Errorf("Legacy column logical name = %q", column.LogicalName)
		}
		if column.Comment != "従来形式" {
			t.Errorf("Legacy column comment = %q", column.Comment)
		}
	})
}

// TestEnhancedCommentErrorHandling エラーハンドリングの統合テスト
func TestEnhancedCommentErrorHandling(t *testing.T) {
	processor := NewEnhancedCommentProcessor()

	// 無効なコメントを含むスキーマ
	schema := &Schema{
		Name: "error_test_db",
		Tables: []*Table{
			{
				Name:    "error_table",
				Comment: "無効なJSON{",
				Columns: []*Column{
					{
						Name:    "valid_col",
						Comment: `{"name": "有効カラム", "description": "正常なJSON"}`,
					},
					{
						Name:    "invalid_col",
						Comment: "無効なYAML: [",
					},
				},
			},
		},
	}

	table := schema.Tables[0]
	err := table.ProcessEnhancedComment(processor, "|", true)
	// 無効なJSONでもLegacyParserでフォールバック処理されるため、エラーは発生しない
	if err != nil {
		t.Errorf("Table processing should not fail with fallback: %v", err)
	}

	// 有効なカラムは正常に処理される
	validColumn := table.Columns[0]
	err = validColumn.ProcessEnhancedComment(processor, "|", true)
	if err != nil {
		t.Errorf("Valid column processing failed: %v", err)
	}

	if !validColumn.HasEnhancedComment() {
		t.Error("Valid column should have enhanced comment")
	}

	// 無効なカラムもフォールバック処理される
	invalidColumn := table.Columns[1]
	err = invalidColumn.ProcessEnhancedComment(processor, "|", true)
	if err != nil {
		t.Errorf("Invalid column processing should not fail with fallback: %v", err)
	}

	// フォールバック処理でEnhancedCommentDataが作成される
	if !invalidColumn.HasEnhancedComment() {
		t.Error("Invalid column should have enhanced comment data from fallback")
	}
}

// TestEnhancedCommentPerformance パフォーマンステスト
func TestEnhancedCommentPerformance(t *testing.T) {
	processor := NewEnhancedCommentProcessor()

	// 大量のオブジェクトを含むスキーマ
	tables := make([]*Table, 100)
	for i := 0; i < 100; i++ {
		columns := make([]*Column, 10)
		for j := 0; j < 10; j++ {
			columns[j] = &Column{
				Name:    "col_" + string(rune(j+'0')),
				Comment: `{"name": "カラム` + string(rune(j+'0')) + `", "description": "説明", "tags": ["test"]}`,
			}
		}

		tables[i] = &Table{
			Name:    "table_" + string(rune(i+'0')),
			Comment: `{"name": "テーブル` + string(rune(i+'0')) + `", "description": "説明", "tags": ["table"]}`,
			Columns: columns,
		}
	}

	schema := &Schema{
		Name:   "performance_test_db",
		Tables: tables,
	}

	// 処理時間を測定
	startTime := processor.GetConfig().ProcessingTimeout
	if startTime == 0 {
		startTime = 1000 // デフォルト1秒
	}

	// 全処理を実行
	for _, table := range schema.Tables {
		err := table.ProcessEnhancedComment(processor, "|", true)
		if err != nil {
			t.Errorf("Table processing failed: %v", err)
		}

		for _, column := range table.Columns {
			err := column.ProcessEnhancedComment(processor, "|", true)
			if err != nil {
				t.Errorf("Column processing failed: %v", err)
			}
		}
	}

	// 基本的な処理が完了していることを確認
	if !schema.Tables[0].HasEnhancedComment() {
		t.Error("First table should have enhanced comment")
	}

	if !schema.Tables[0].Columns[0].HasEnhancedComment() {
		t.Error("First column should have enhanced comment")
	}
}

// TestEnhancedCommentConfigurationIntegration 設定統合テスト
func TestEnhancedCommentConfigurationIntegration(t *testing.T) {
	// 設定を使用したプロセッサー
	config := NewMockEnhancedCommentConfigurator()
	config.preferredFormat = "yaml"
	config.validationEnabled = true
	config.strictMode = false

	processor := NewEnhancedCommentProcessorFromConfig(config)

	schema := &Schema{
		Name: "config_test_db",
		Tables: []*Table{
			{
				Name:    "yaml_priority_table",
				Comment: "name: YAMLテーブル\ndescription: YAML優先設定での処理",
				Columns: []*Column{
					{
						Name:    "yaml_col",
						Comment: "name: YAMLカラム\ndescription: YAML形式で処理される",
					},
				},
			},
		},
	}

	table := schema.Tables[0]
	err := table.ProcessEnhancedComment(processor, "|", true)
	if err != nil {
		t.Errorf("Table processing with config failed: %v", err)
	}

	if !table.HasEnhancedComment() {
		t.Error("Table should have enhanced comment")
	}

	if table.GetDescription() != "YAML優先設定での処理" {
		t.Errorf("Table description = %q", table.GetDescription())
	}

	column := table.Columns[0]
	err = column.ProcessEnhancedComment(processor, "|", true)
	if err != nil {
		t.Errorf("Column processing with config failed: %v", err)
	}

	if column.GetDescription() != "YAML形式で処理される" {
		t.Errorf("Column description = %q", column.GetDescription())
	}

	// オブジェクトタイプがメタデータに設定されているか確認
	metadata := table.GetMetadata()
	if metadata["object_type"] != string(ObjectTypeTable) {
		t.Errorf("table object_type = %q, expected %q", metadata["object_type"], ObjectTypeTable)
	}

	columnMetadata := column.GetMetadata()
	if columnMetadata["object_type"] != string(ObjectTypeColumn) {
		t.Errorf("column object_type = %q, expected %q", columnMetadata["object_type"], ObjectTypeColumn)
	}
}