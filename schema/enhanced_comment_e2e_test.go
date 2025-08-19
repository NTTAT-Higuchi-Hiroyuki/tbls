package schema

import (
	"testing"
)

// TestEnhancedCommentEndToEnd エンドツーエンドテスト
func TestEnhancedCommentEndToEnd(t *testing.T) {
	// 実際のアプリケーションスキーマを模擬
	schema := createSampleApplicationSchema()
	processor := NewEnhancedCommentProcessor()

	// スキーマ全体を処理
	err := processCompleteSchema(schema, processor)
	if err != nil {
		t.Fatalf("Schema processing failed: %v", err)
	}

	// 処理結果の詳細検証
	validateProcessedSchema(t, schema)
}

// createSampleApplicationSchema サンプルアプリケーションスキーマを作成
func createSampleApplicationSchema() *Schema {
	return &Schema{
		Name: "sample_app",
		Tables: []*Table{
			// ユーザー管理テーブル
			{
				Name:    "users",
				Comment: `{"name": "ユーザー", "description": "システムユーザーの基本情報", "tags": ["マスター", "認証"], "priority": 1}`,
				Columns: []*Column{
					{
						Name:    "id",
						Comment: `{"name": "ユーザーID", "description": "システム内でユーザーを一意に識別するID", "tags": ["PK", "自動生成"], "priority": 1}`,
					},
					{
						Name:    "username",
						Comment: "ユーザー名|ログイン時に使用する一意のユーザー名",
					},
					{
						Name:    "email",
						Comment: "description: メールアドレス\ntags:\n  - 連絡先\n  - 一意\n  - 必須",
					},
					{
						Name:    "password_hash",
						Comment: `{"name": "パスワードハッシュ", "description": "ハッシュ化されたパスワード", "tags": ["セキュリティ", "機密"], "deprecated": false}`,
					},
					{
						Name:    "created_at",
						Comment: "作成日時",
					},
					{
						Name:    "updated_at",
						Comment: "description: 更新日時\ndeprecated: false",
					},
				},
				Indexes: []*Index{
					{
						Name:    "idx_users_username",
						Comment: `{"description": "ユーザー名一意インデックス", "tags": ["一意性", "ログイン性能"]}`,
					},
					{
						Name:    "idx_users_email",
						Comment: "メールアドレス一意インデックス",
					},
				},
				Constraints: []*Constraint{
					{
						Name:    "pk_users",
						Comment: `{"description": "ユーザーテーブル主キー", "tags": ["整合性"]}`,
					},
					{
						Name:    "uq_users_username",
						Comment: "ユーザー名一意制約",
					},
				},
				Triggers: []*Trigger{
					{
						Name:    "tr_users_updated_at",
						Comment: "description: 更新日時自動設定トリガー\ntags:\n  - 自動化\n  - 監査",
					},
				},
			},
			// 投稿管理テーブル
			{
				Name:    "posts",
				Comment: "投稿|ユーザーの投稿内容を管理するテーブル",
				Columns: []*Column{
					{
						Name:    "id",
						Comment: `{"name": "投稿ID", "description": "投稿の一意識別子", "tags": ["PK"]}`,
					},
					{
						Name:    "user_id",
						Comment: `{"name": "ユーザーID", "description": "投稿者のユーザーID", "tags": ["FK", "必須"]}`,
					},
					{
						Name:    "title",
						Comment: "タイトル|投稿のタイトル",
					},
					{
						Name:    "content",
						Comment: "内容|投稿の本文内容",
					},
					{
						Name:    "status",
						Comment: "description: 投稿ステータス\ntags:\n  - enum\n  - ワークフロー",
					},
				},
				Indexes: []*Index{
					{
						Name:    "idx_posts_user_id",
						Comment: `{"description": "ユーザー別投稿検索用インデックス", "tags": ["性能", "検索"]}`,
					},
				},
				Constraints: []*Constraint{
					{
						Name:    "fk_posts_user_id",
						Comment: `{"description": "投稿者参照外部キー", "tags": ["参照整合性"]}`,
					},
				},
			},
			// タグ管理テーブル（多対多関係）
			{
				Name:    "tags",
				Comment: `{"name": "タグ", "description": "投稿の分類用タグ", "tags": ["マスター"]}`,
				Columns: []*Column{
					{
						Name:    "id",
						Comment: "タグID",
					},
					{
						Name:    "name",
						Comment: "name: タグ名\ndescription: タグの表示名\ntags:\n  - 一意",
					},
				},
			},
			{
				Name:    "post_tags",
				Comment: "投稿タグ関連|投稿とタグの多対多関係を管理",
				Columns: []*Column{
					{
						Name:    "post_id",
						Comment: `{"name": "投稿ID", "description": "関連する投稿のID", "tags": ["FK", "複合PK"]}`,
					},
					{
						Name:    "tag_id",
						Comment: `{"name": "タグID", "description": "関連するタグのID", "tags": ["FK", "複合PK"]}`,
					},
				},
				Constraints: []*Constraint{
					{
						Name:    "pk_post_tags",
						Comment: "description: 投稿タグ複合主キー\ntags:\n  - 複合キー\n  - 一意性",
					},
					{
						Name:    "fk_post_tags_post_id",
						Comment: `{"description": "投稿参照外部キー", "tags": ["参照整合性"]}`,
					},
					{
						Name:    "fk_post_tags_tag_id",
						Comment: `{"description": "タグ参照外部キー", "tags": ["参照整合性"]}`,
					},
				},
			},
		},
	}
}

// processCompleteSchema スキーマ全体を処理
func processCompleteSchema(schema *Schema, processor CommentProcessor) error {
	for _, table := range schema.Tables {
		if err := table.ProcessEnhancedComment(processor, "|", true); err != nil {
			return err
		}

		for _, column := range table.Columns {
			if err := column.ProcessEnhancedComment(processor, "|", true); err != nil {
				return err
			}
		}

		for _, index := range table.Indexes {
			if err := index.ProcessEnhancedComment(processor, "|"); err != nil {
				return err
			}
		}

		for _, constraint := range table.Constraints {
			if err := constraint.ProcessEnhancedComment(processor, "|"); err != nil {
				return err
			}
		}

		for _, trigger := range table.Triggers {
			if err := trigger.ProcessEnhancedComment(processor, "|"); err != nil {
				return err
			}
		}
	}
	return nil
}

// validateProcessedSchema 処理結果を検証
func validateProcessedSchema(t *testing.T, schema *Schema) {
	// ユーザーテーブルの検証
	usersTable := schema.Tables[0]
	if !usersTable.HasEnhancedComment() {
		t.Error("Users table should have enhanced comment")
	}

	if usersTable.LogicalName != "ユーザー" {
		t.Errorf("Users table logical name = %q, expected 'ユーザー'", usersTable.LogicalName)
	}

	if usersTable.GetDescription() != "システムユーザーの基本情報" {
		t.Errorf("Users table description = %q", usersTable.GetDescription())
	}

	expectedUsersTags := []string{"マスター", "認証"}
	usersTags := usersTable.GetTags()
	if len(usersTags) != len(expectedUsersTags) {
		t.Errorf("Users table tags = %v, expected %v", usersTags, expectedUsersTags)
	}

	// ユーザーテーブルのカラム検証
	validateUsersTableColumns(t, usersTable)

	// ユーザーテーブルのインデックス/制約/トリガー検証
	validateUsersTableRelatedObjects(t, usersTable)

	// 投稿テーブルの検証
	postsTable := schema.Tables[1]
	if postsTable.LogicalName != "投稿" {
		t.Errorf("Posts table logical name = %q, expected '投稿'", postsTable.LogicalName)
	}

	if postsTable.Comment != "ユーザーの投稿内容を管理するテーブル" {
		t.Errorf("Posts table comment = %q", postsTable.Comment)
	}

	// 投稿テーブルのカラム検証
	validatePostsTableColumns(t, postsTable)

	// タグテーブルの検証（JSON形式）
	tagsTable := schema.Tables[2]
	if !tagsTable.HasEnhancedComment() {
		t.Error("Tags table should have enhanced comment")
	}

	if tagsTable.LogicalName != "タグ" {
		t.Errorf("Tags table logical name = %q, expected 'タグ'", tagsTable.LogicalName)
	}

	// 投稿タグ関連テーブルの検証（従来形式）
	postTagsTable := schema.Tables[3]
	if postTagsTable.LogicalName != "投稿タグ関連" {
		t.Errorf("Post tags table logical name = %q, expected '投稿タグ関連'", postTagsTable.LogicalName)
	}
}

// validateUsersTableColumns ユーザーテーブルのカラム検証
func validateUsersTableColumns(t *testing.T, usersTable *Table) {
	// IDカラム（JSON形式）
	idColumn := usersTable.Columns[0]
	if !idColumn.HasEnhancedComment() {
		t.Error("ID column should have enhanced comment")
	}

	if idColumn.LogicalName != "ユーザーID" {
		t.Errorf("ID column logical name = %q, expected 'ユーザーID'", idColumn.LogicalName)
	}

	if idColumn.GetDescription() != "システム内でユーザーを一意に識別するID" {
		t.Errorf("ID column description = %q", idColumn.GetDescription())
	}

	expectedIdTags := []string{"PK", "自動生成"}
	idTags := idColumn.GetTags()
	if len(idTags) != len(expectedIdTags) {
		t.Errorf("ID column tags = %v, expected %v", idTags, expectedIdTags)
	}

	// ユーザー名カラム（従来形式）
	usernameColumn := usersTable.Columns[1]
	if usernameColumn.LogicalName != "ユーザー名" {
		t.Errorf("Username column logical name = %q, expected 'ユーザー名'", usernameColumn.LogicalName)
	}

	if usernameColumn.Comment != "ログイン時に使用する一意のユーザー名" {
		t.Errorf("Username column comment = %q", usernameColumn.Comment)
	}

	// メールカラム（YAML形式）
	emailColumn := usersTable.Columns[2]
	if !emailColumn.HasEnhancedComment() {
		t.Error("Email column should have enhanced comment")
	}

	if emailColumn.GetDescription() != "メールアドレス" {
		t.Errorf("Email column description = %q", emailColumn.GetDescription())
	}

	expectedEmailTags := []string{"連絡先", "一意", "必須"}
	emailTags := emailColumn.GetTags()
	if len(emailTags) != len(expectedEmailTags) {
		t.Errorf("Email column tags = %v, expected %v", emailTags, expectedEmailTags)
	}

	// パスワードハッシュカラム（JSON形式、非推奨フラグあり）
	passwordColumn := usersTable.Columns[3]
	if passwordColumn.IsDeprecated() {
		t.Error("Password column should not be deprecated")
	}

	expectedPasswordTags := []string{"セキュリティ", "機密"}
	passwordTags := passwordColumn.GetTags()
	if len(passwordTags) != len(expectedPasswordTags) {
		t.Errorf("Password column tags = %v, expected %v", passwordTags, expectedPasswordTags)
	}

	// 作成日時カラム（シンプル形式）
	createdAtColumn := usersTable.Columns[4]
	if !createdAtColumn.HasEnhancedComment() {
		t.Error("Created at column should have enhanced comment data from legacy parser")
	}

	// 更新日時カラム（YAML形式）
	updatedAtColumn := usersTable.Columns[5]
	if updatedAtColumn.GetDescription() != "更新日時" {
		t.Errorf("Updated at column description = %q", updatedAtColumn.GetDescription())
	}

	if updatedAtColumn.IsDeprecated() {
		t.Error("Updated at column should not be deprecated")
	}
}

// validateUsersTableRelatedObjects ユーザーテーブルの関連オブジェクト検証
func validateUsersTableRelatedObjects(t *testing.T, usersTable *Table) {
	// インデックス検証
	if len(usersTable.Indexes) != 2 {
		t.Errorf("Users table should have 2 indexes, got %d", len(usersTable.Indexes))
	}

	usernameIndex := usersTable.Indexes[0]
	if !usernameIndex.HasEnhancedComment() {
		t.Error("Username index should have enhanced comment")
	}

	if usernameIndex.GetDescription() != "ユーザー名一意インデックス" {
		t.Errorf("Username index description = %q", usernameIndex.GetDescription())
	}

	expectedIndexTags := []string{"一意性", "ログイン性能"}
	indexTags := usernameIndex.GetTags()
	if len(indexTags) != len(expectedIndexTags) {
		t.Errorf("Username index tags = %v, expected %v", indexTags, expectedIndexTags)
	}

	emailIndex := usersTable.Indexes[1]
	if !emailIndex.HasEnhancedComment() {
		t.Error("Email index should have enhanced comment data from legacy parser")
	}

	// 制約検証
	if len(usersTable.Constraints) != 2 {
		t.Errorf("Users table should have 2 constraints, got %d", len(usersTable.Constraints))
	}

	pkConstraint := usersTable.Constraints[0]
	if !pkConstraint.HasEnhancedComment() {
		t.Error("PK constraint should have enhanced comment")
	}

	if pkConstraint.GetDescription() != "ユーザーテーブル主キー" {
		t.Errorf("PK constraint description = %q", pkConstraint.GetDescription())
	}

	// トリガー検証
	if len(usersTable.Triggers) != 1 {
		t.Errorf("Users table should have 1 trigger, got %d", len(usersTable.Triggers))
	}

	trigger := usersTable.Triggers[0]
	if !trigger.HasEnhancedComment() {
		t.Error("Trigger should have enhanced comment")
	}

	if trigger.GetDescription() != "更新日時自動設定トリガー" {
		t.Errorf("Trigger description = %q", trigger.GetDescription())
	}

	expectedTriggerTags := []string{"自動化", "監査"}
	triggerTags := trigger.GetTags()
	if len(triggerTags) != len(expectedTriggerTags) {
		t.Errorf("Trigger tags = %v, expected %v", triggerTags, expectedTriggerTags)
	}
}

// validatePostsTableColumns 投稿テーブルのカラム検証
func validatePostsTableColumns(t *testing.T, postsTable *Table) {
	// 投稿IDカラム（JSON形式）
	idColumn := postsTable.Columns[0]
	if idColumn.LogicalName != "投稿ID" {
		t.Errorf("Post ID column logical name = %q, expected '投稿ID'", idColumn.LogicalName)
	}

	// ユーザーIDカラム（JSON形式、外部キー）
	userIdColumn := postsTable.Columns[1]
	if userIdColumn.LogicalName != "ユーザーID" {
		t.Errorf("Post user ID column logical name = %q, expected 'ユーザーID'", userIdColumn.LogicalName)
	}

	expectedUserIdTags := []string{"FK", "必須"}
	userIdTags := userIdColumn.GetTags()
	if len(userIdTags) != len(expectedUserIdTags) {
		t.Errorf("Post user ID column tags = %v, expected %v", userIdTags, expectedUserIdTags)
	}

	// タイトルカラム（従来形式）
	titleColumn := postsTable.Columns[2]
	if titleColumn.LogicalName != "タイトル" {
		t.Errorf("Title column logical name = %q, expected 'タイトル'", titleColumn.LogicalName)
	}

	if titleColumn.Comment != "投稿のタイトル" {
		t.Errorf("Title column comment = %q, expected '投稿のタイトル'", titleColumn.Comment)
	}

	// 内容カラム（従来形式）
	contentColumn := postsTable.Columns[3]
	if contentColumn.LogicalName != "内容" {
		t.Errorf("Content column logical name = %q, expected '内容'", contentColumn.LogicalName)
	}

	// ステータスカラム（YAML形式）
	statusColumn := postsTable.Columns[4]
	if statusColumn.GetDescription() != "投稿ステータス" {
		t.Errorf("Status column description = %q", statusColumn.GetDescription())
	}

	expectedStatusTags := []string{"enum", "ワークフロー"}
	statusTags := statusColumn.GetTags()
	if len(statusTags) != len(expectedStatusTags) {
		t.Errorf("Status column tags = %v, expected %v", statusTags, expectedStatusTags)
	}
}

// TestEnhancedCommentMetadataConsistency メタデータ一貫性テスト
func TestEnhancedCommentMetadataConsistency(t *testing.T) {
	schema := createSampleApplicationSchema()
	processor := NewEnhancedCommentProcessor()

	err := processCompleteSchema(schema, processor)
	if err != nil {
		t.Fatalf("Schema processing failed: %v", err)
	}

	// 全オブジェクトのメタデータでobject_typeが正しく設定されているか確認
	for _, table := range schema.Tables {
		if table.HasEnhancedComment() {
			metadata := table.GetMetadata()
			if metadata["object_type"] != string(ObjectTypeTable) {
				t.Errorf("Table %s object_type = %q, expected %q", table.Name, metadata["object_type"], ObjectTypeTable)
			}
		}

		for _, column := range table.Columns {
			if column.HasEnhancedComment() {
				metadata := column.GetMetadata()
				if metadata["object_type"] != string(ObjectTypeColumn) {
					t.Errorf("Column %s.%s object_type = %q, expected %q", table.Name, column.Name, metadata["object_type"], ObjectTypeColumn)
				}
			}
		}

		for _, index := range table.Indexes {
			if index.HasEnhancedComment() {
				metadata := index.GetMetadata()
				if metadata["object_type"] != string(ObjectTypeIndex) {
					t.Errorf("Index %s object_type = %q, expected %q", index.Name, metadata["object_type"], ObjectTypeIndex)
				}
			}
		}

		for _, constraint := range table.Constraints {
			if constraint.HasEnhancedComment() {
				metadata := constraint.GetMetadata()
				if metadata["object_type"] != string(ObjectTypeConstraint) {
					t.Errorf("Constraint %s object_type = %q, expected %q", constraint.Name, metadata["object_type"], ObjectTypeConstraint)
				}
			}
		}

		for _, trigger := range table.Triggers {
			if trigger.HasEnhancedComment() {
				metadata := trigger.GetMetadata()
				if metadata["object_type"] != string(ObjectTypeTrigger) {
					t.Errorf("Trigger %s object_type = %q, expected %q", trigger.Name, metadata["object_type"], ObjectTypeTrigger)
				}
			}
		}
	}
}