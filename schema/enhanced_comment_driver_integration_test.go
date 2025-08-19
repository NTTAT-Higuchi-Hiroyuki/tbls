//go:build integration

package schema

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/k1LoW/tbls/drivers"
	"github.com/k1LoW/tbls/drivers/postgres"
	"github.com/k1LoW/tbls/drivers/sqlite"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xo/dburl"
)

// TestEnhancedCommentDriverIntegration ドライバー統合テスト
func TestEnhancedCommentDriverIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// SQLite統合テスト（ローカル環境）
	t.Run("SQLite統合テスト", func(t *testing.T) {
		testSQLiteIntegration(t)
	})

	// PostgreSQL統合テスト（Docker環境が必要）
	if os.Getenv("POSTGRES_TEST_ENABLED") == "true" {
		t.Run("PostgreSQL統合テスト", func(t *testing.T) {
			testPostgreSQLIntegration(t)
		})
	}
}

// testSQLiteIntegration SQLite統合テスト
func testSQLiteIntegration(t *testing.T) {
	// テスト用SQLiteデータベース作成
	dbPath := "/tmp/test_enhanced_comment.db"
	defer os.Remove(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	// テストデータ作成
	setupSQLiteTestData(t, db)

	// Schemaオブジェクト作成
	schema := &Schema{Name: "test_database"}

	// SQLiteドライバーで解析
	driver := sqlite.New(db)
	driver.SetLogicalNameConfig("|", true)
	driver.SetTableLogicalNameConfig("|", true)

	err = driver.Analyze(schema)
	if err != nil {
		t.Fatalf("Driver analysis failed: %v", err)
	}

	// 拡張コメント処理プロセッサで追加処理
	processor := NewEnhancedCommentProcessor()
	err = processSchemaWithEnhancedComments(schema, processor)
	if err != nil {
		t.Fatalf("Enhanced comment processing failed: %v", err)
	}

	// 統合結果の検証
	validateSQLiteIntegrationResults(t, schema)
}

// setupSQLiteTestData SQLiteテストデータセットアップ
func setupSQLiteTestData(t *testing.T, db *sql.DB) {
	queries := []string{
		// テーブル作成（様々な形式のコメント）
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			profile_data TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// SQLiteのコメント追加（テーブル）
		// 注意: SQLiteはテーブルコメントを直接サポートしていないため、
		// 別の方法でコメント情報を保存する必要がある

		`CREATE TABLE posts (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT,
			status TEXT DEFAULT 'draft',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		`CREATE TABLE tags (
			id INTEGER PRIMARY KEY,
			name TEXT UNIQUE NOT NULL
		)`,

		`CREATE TABLE post_tags (
			post_id INTEGER,
			tag_id INTEGER,
			PRIMARY KEY (post_id, tag_id),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
		)`,

		// インデックス作成
		`CREATE INDEX idx_posts_user_id ON posts(user_id)`,
		`CREATE UNIQUE INDEX idx_users_username ON users(username)`,
		`CREATE INDEX idx_posts_status ON posts(status)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("Failed to execute query: %q, error: %v", query, err)
		}
	}

	// テストデータ挿入
	insertQueries := []string{
		`INSERT INTO users (username, email, profile_data) VALUES 
			('testuser1', 'user1@example.com', '{"name": "テストユーザー1", "tags": ["admin", "active"]}'),
			('testuser2', 'user2@example.com', 'name: テストユーザー2\ntags:\n  - regular\n  - active')`,
		
		`INSERT INTO posts (user_id, title, content, status) VALUES 
			(1, 'テスト投稿1', 'テスト用の投稿内容です', 'published'),
			(2, 'テスト投稿2', 'もう一つのテスト投稿', 'draft')`,
		
		`INSERT INTO tags (name) VALUES ('テクノロジー'), ('プログラミング'), ('テスト')`,
		
		`INSERT INTO post_tags (post_id, tag_id) VALUES (1, 1), (1, 2), (2, 3)`,
	}

	for _, query := range insertQueries {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("Failed to insert test data: %q, error: %v", query, err)
		}
	}
}

// testPostgreSQLIntegration PostgreSQL統合テスト
func testPostgreSQLIntegration(t *testing.T) {
	// PostgreSQL接続（環境変数から取得）
	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:pgpass@localhost:5432/testdb?sslmode=disable"
	}

	db, err := dburl.Open(dsn)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// 接続テスト
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	// テストデータセットアップ
	setupPostgreSQLTestData(t, db)

	// Schemaオブジェクト作成
	schema := &Schema{Name: "postgres_test_db"}

	// PostgreSQLドライバーで解析
	driver := postgres.New(db)
	driver.SetLogicalNameConfig("|", true)
	driver.SetTableLogicalNameConfig("|", true)

	err = driver.Analyze(schema)
	if err != nil {
		t.Fatalf("PostgreSQL driver analysis failed: %v", err)
	}

	// 拡張コメント処理プロセッサで追加処理
	processor := NewEnhancedCommentProcessor()
	err = processSchemaWithEnhancedComments(schema, processor)
	if err != nil {
		t.Fatalf("Enhanced comment processing failed: %v", err)
	}

	// 統合結果の検証
	validatePostgreSQLIntegrationResults(t, schema)

	// テストデータクリーンアップ
	cleanupPostgreSQLTestData(t, db)
}

// setupPostgreSQLTestData PostgreSQLテストデータセットアップ
func setupPostgreSQLTestData(t *testing.T, db *sql.DB) {
	// 既存テーブル削除
	cleanupQueries := []string{
		`DROP TABLE IF EXISTS integration_post_tags CASCADE`,
		`DROP TABLE IF EXISTS integration_tags CASCADE`,
		`DROP TABLE IF EXISTS integration_posts CASCADE`,
		`DROP TABLE IF EXISTS integration_users CASCADE`,
	}

	for _, query := range cleanupQueries {
		db.Exec(query) // エラーは無視（テーブルが存在しない場合）
	}

	// テストテーブル作成（様々な形式のコメント付き）
	createQueries := []string{
		`CREATE TABLE integration_users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			profile_data JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// テーブルコメント追加（JSON形式）
		`COMMENT ON TABLE integration_users IS '{"name": "統合テストユーザー", "description": "ドライバー統合テスト用のユーザーテーブル", "tags": ["integration", "test", "user"], "priority": 1}'`,

		// カラムコメント追加（様々な形式）
		`COMMENT ON COLUMN integration_users.id IS '{"name": "ユーザーID", "description": "システム内でユーザーを一意に識別するID", "tags": ["PK", "auto_increment"]}'`,
		`COMMENT ON COLUMN integration_users.username IS 'ユーザー名|ログイン時に使用する一意のユーザー名'`,
		`COMMENT ON COLUMN integration_users.email IS 'description: メールアドレス\ntags:\n  - contact\n  - unique\n  - required'`,
		`COMMENT ON COLUMN integration_users.profile_data IS 'プロフィールデータ（JSON形式）'`,
		`COMMENT ON COLUMN integration_users.created_at IS 'name: 作成日時\ndescription: ユーザーの作成日時\ndefault: CURRENT_TIMESTAMP'`,

		`CREATE TABLE integration_posts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES integration_users(id),
			title VARCHAR(200) NOT NULL,
			content TEXT,
			status VARCHAR(20) DEFAULT 'draft',
			tags JSONB DEFAULT '[]',
			metadata JSONB DEFAULT '{}',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// 投稿テーブルのコメント（Legacy形式）
		`COMMENT ON TABLE integration_posts IS '統合テスト投稿|ドライバー統合テスト用の投稿テーブル'`,

		// カラムコメント
		`COMMENT ON COLUMN integration_posts.id IS '{"name": "投稿ID", "description": "投稿の一意識別子", "tags": ["PK"]}'`,
		`COMMENT ON COLUMN integration_posts.user_id IS '{"name": "ユーザーID", "description": "投稿者のユーザーID", "tags": ["FK", "required"]}'`,
		`COMMENT ON COLUMN integration_posts.title IS 'タイトル|投稿のタイトル'`,
		`COMMENT ON COLUMN integration_posts.content IS 'description: 投稿内容\ntags:\n  - text\n  - content'`,
		`COMMENT ON COLUMN integration_posts.status IS '{"name": "ステータス", "description": "投稿の公開状態", "tags": ["enum"], "metadata": {"values": ["draft", "published", "archived"]}}'`,
		`COMMENT ON COLUMN integration_posts.tags IS 'タグ（JSON形式）'`,
		`COMMENT ON COLUMN integration_posts.metadata IS 'description: メタデータ\ntags:\n  - json\n  - optional'`,

		`CREATE TABLE integration_tags (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50) UNIQUE NOT NULL,
			color VARCHAR(7) DEFAULT '#000000',
			description TEXT
		)`,

		// タグテーブル（YAML形式）
		`COMMENT ON TABLE integration_tags IS 'name: 統合テストタグ
description: ドライバー統合テスト用のタグテーブル
tags:
  - integration
  - test
  - tag
priority: 2'`,

		`COMMENT ON COLUMN integration_tags.id IS 'タグID'`,
		`COMMENT ON COLUMN integration_tags.name IS '{"name": "タグ名", "description": "タグの表示名", "tags": ["unique", "required"]}'`,
		`COMMENT ON COLUMN integration_tags.color IS 'description: カラーコード\nformat: hex\ndefault: "#000000"'`,

		`CREATE TABLE integration_post_tags (
			post_id INTEGER REFERENCES integration_posts(id) ON DELETE CASCADE,
			tag_id INTEGER REFERENCES integration_tags(id) ON DELETE CASCADE,
			PRIMARY KEY (post_id, tag_id)
		)`,

		`COMMENT ON TABLE integration_post_tags IS '{"name": "投稿タグ関連", "description": "投稿とタグの多対多関係", "tags": ["junction", "relation"]}'`,
	}

	for _, query := range createQueries {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("Failed to execute query: %q, error: %v", query, err)
		}
	}

	// インデックス作成
	indexQueries := []string{
		`CREATE INDEX idx_integration_posts_user_id ON integration_posts(user_id)`,
		`CREATE INDEX idx_integration_posts_status ON integration_posts(status)`,
		`CREATE INDEX idx_integration_posts_created_at ON integration_posts(created_at)`,
	}

	for _, query := range indexQueries {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("Failed to create index: %q, error: %v", query, err)
		}
	}

	// インデックスコメント
	commentQueries := []string{
		`COMMENT ON INDEX idx_integration_posts_user_id IS '{"description": "ユーザー別投稿検索用インデックス", "tags": ["performance", "lookup"]}'`,
		`COMMENT ON INDEX idx_integration_posts_status IS 'ステータス検索用インデックス'`,
		`COMMENT ON INDEX idx_integration_posts_created_at IS 'description: 作成日時検索用インデックス\ntags:\n  - performance\n  - time_series'`,
	}

	for _, query := range commentQueries {
		if _, err := db.Exec(query); err != nil {
			t.Logf("Warning: Failed to add index comment: %q, error: %v", query, err)
		}
	}

	// テストデータ挿入
	insertQueries := []string{
		`INSERT INTO integration_users (username, email, profile_data) VALUES 
			('integtest1', 'integtest1@example.com', '{"name": "統合テストユーザー1", "department": "開発"}'),
			('integtest2', 'integtest2@example.com', '{"name": "統合テストユーザー2", "department": "QA"}')`,

		`INSERT INTO integration_posts (user_id, title, content, status, tags, metadata) VALUES 
			(1, '統合テスト投稿1', '統合テスト用の投稿内容です', 'published', '["tech", "test"]', '{"views": 100}'),
			(2, '統合テスト投稿2', 'もう一つの統合テスト投稿', 'draft', '["qa", "testing"]', '{"priority": "high"}')`,

		`INSERT INTO integration_tags (name, color, description) VALUES 
			('統合テスト', '#FF5722', '統合テスト用のタグ'),
			('ドライバーテスト', '#2196F3', 'ドライバー統合テスト用')`,

		`INSERT INTO integration_post_tags (post_id, tag_id) VALUES (1, 1), (2, 2)`,
	}

	for _, query := range insertQueries {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("Failed to insert test data: %q, error: %v", query, err)
		}
	}
}

// processSchemaWithEnhancedComments スキーマ全体に拡張コメント処理を適用
func processSchemaWithEnhancedComments(schema *Schema, processor CommentProcessor) error {
	for _, table := range schema.Tables {
		if err := table.ProcessEnhancedComment(processor, "|", true); err != nil {
			return fmt.Errorf("table %s processing failed: %w", table.Name, err)
		}

		for _, column := range table.Columns {
			if err := column.ProcessEnhancedComment(processor, "|", true); err != nil {
				return fmt.Errorf("column %s.%s processing failed: %w", table.Name, column.Name, err)
			}
		}

		for _, index := range table.Indexes {
			if err := index.ProcessEnhancedComment(processor, "|"); err != nil {
				return fmt.Errorf("index %s processing failed: %w", index.Name, err)
			}
		}

		for _, constraint := range table.Constraints {
			if err := constraint.ProcessEnhancedComment(processor, "|"); err != nil {
				return fmt.Errorf("constraint %s processing failed: %w", constraint.Name, err)
			}
		}

		for _, trigger := range table.Triggers {
			if err := trigger.ProcessEnhancedComment(processor, "|"); err != nil {
				return fmt.Errorf("trigger %s processing failed: %w", trigger.Name, err)
			}
		}
	}
	return nil
}

// validateSQLiteIntegrationResults SQLite統合結果検証
func validateSQLiteIntegrationResults(t *testing.T, schema *Schema) {
	t.Log("=== SQLite統合テスト結果検証 ===")

	// テーブル数確認
	if len(schema.Tables) < 4 {
		t.Errorf("Expected at least 4 tables, got %d", len(schema.Tables))
	}

	// usersテーブル検証
	usersTable, err := schema.FindTableByName("users")
	if err != nil {
		t.Fatal("Users table not found")
	}

	// SQLiteではテーブルコメントがサポートされていないため、
	// カラムコメントのみ検証
	validateSQLiteUsersTable(t, usersTable)

	// postsテーブル検証
	postsTable, err := schema.FindTableByName("posts")
	if err != nil {
		t.Fatal("Posts table not found")
	}

	validateSQLitePostsTable(t, postsTable)

	t.Log("SQLite統合テスト検証完了")
}

// validateSQLiteUsersTable SQLite usersテーブル検証
func validateSQLiteUsersTable(t *testing.T, table *Table) {
	// ID カラムの検証
	idColumn, err := table.FindColumnByName("id")
	if err != nil {
		t.Error("ID column not found in users table")
		return
	}

	// SQLiteの場合、コメントがデータベースレベルで保存されていないため、
	// 拡張コメント処理は行われない（Commentフィールドが空）
	if idColumn.Comment != "" {
		t.Logf("ID column has comment: %s", idColumn.Comment)
	}

	// username カラムの検証
	usernameColumn, err := table.FindColumnByName("username")
	if err != nil {
		t.Error("Username column not found in users table")
		return
	}

	if usernameColumn.Comment != "" {
		t.Logf("Username column has comment: %s", usernameColumn.Comment)
	}

	t.Logf("Users table validation completed - columns: %d", len(table.Columns))
}

// validateSQLitePostsTable SQLite postsテーブル検証
func validateSQLitePostsTable(t *testing.T, table *Table) {
	// 外部キー制約の検証
	if len(table.Constraints) > 0 {
		t.Logf("Posts table has %d constraints", len(table.Constraints))
	}

	// インデックスの検証
	if len(table.Indexes) > 0 {
		t.Logf("Posts table has %d indexes", len(table.Indexes))
	}

	t.Logf("Posts table validation completed - columns: %d", len(table.Columns))
}

// validatePostgreSQLIntegrationResults PostgreSQL統合結果検証
func validatePostgreSQLIntegrationResults(t *testing.T, schema *Schema) {
	t.Log("=== PostgreSQL統合テスト結果検証 ===")

	// integration_usersテーブル検証
	usersTable, err := schema.FindTableByName("integration_users")
	if err != nil {
		t.Fatal("integration_users table not found")
	}

	validatePostgreSQLUsersTable(t, usersTable)

	// integration_postsテーブル検証
	postsTable, err := schema.FindTableByName("integration_posts")
	if err != nil {
		t.Fatal("integration_posts table not found")
	}

	validatePostgreSQLPostsTable(t, postsTable)

	// integration_tagsテーブル検証
	tagsTable, err := schema.FindTableByName("integration_tags")
	if err != nil {
		t.Fatal("integration_tags table not found")
	}

	validatePostgreSQLTagsTable(t, tagsTable)

	t.Log("PostgreSQL統合テスト検証完了")
}

// validatePostgreSQLUsersTable PostgreSQL usersテーブル検証
func validatePostgreSQLUsersTable(t *testing.T, table *Table) {
	// テーブルコメント（JSON形式）の検証
	if !table.HasEnhancedComment() {
		t.Error("Users table should have enhanced comment")
		return
	}

	if table.LogicalName != "統合テストユーザー" {
		t.Errorf("Users table logical name = %q, expected '統合テストユーザー'", table.LogicalName)
	}

	if table.GetDescription() != "ドライバー統合テスト用のユーザーテーブル" {
		t.Errorf("Users table description = %q", table.GetDescription())
	}

	expectedTags := []string{"integration", "test", "user"}
	tags := table.GetTags()
	if len(tags) != len(expectedTags) {
		t.Errorf("Users table tags = %v, expected %v", tags, expectedTags)
	}

	// ID カラム（JSON形式）の検証
	idColumn, err := table.FindColumnByName("id")
	if err != nil {
		t.Error("ID column not found")
		return
	}

	if !idColumn.HasEnhancedComment() {
		t.Error("ID column should have enhanced comment")
		return
	}

	if idColumn.LogicalName != "ユーザーID" {
		t.Errorf("ID column logical name = %q, expected 'ユーザーID'", idColumn.LogicalName)
	}

	// username カラム（Legacy形式）の検証
	usernameColumn, err := table.FindColumnByName("username")
	if err != nil {
		t.Error("Username column not found")
		return
	}

	if usernameColumn.LogicalName != "ユーザー名" {
		t.Errorf("Username column logical name = %q, expected 'ユーザー名'", usernameColumn.LogicalName)
	}

	// email カラム（YAML形式）の検証
	emailColumn, err := table.FindColumnByName("email")
	if err != nil {
		t.Error("Email column not found")
		return
	}

	if !emailColumn.HasEnhancedComment() {
		t.Error("Email column should have enhanced comment")
		return
	}

	if emailColumn.GetDescription() != "メールアドレス" {
		t.Errorf("Email column description = %q", emailColumn.GetDescription())
	}

	expectedEmailTags := []string{"contact", "unique", "required"}
	emailTags := emailColumn.GetTags()
	if len(emailTags) != len(expectedEmailTags) {
		t.Errorf("Email column tags = %v, expected %v", emailTags, expectedEmailTags)
	}

	t.Logf("PostgreSQL Users table validation completed successfully")
}

// validatePostgreSQLPostsTable PostgreSQL postsテーブル検証
func validatePostgreSQLPostsTable(t *testing.T, table *Table) {
	// テーブルコメント（Legacy形式）の検証
	if table.LogicalName != "統合テスト投稿" {
		t.Errorf("Posts table logical name = %q, expected '統合テスト投稿'", table.LogicalName)
	}

	if table.Comment != "ドライバー統合テスト用の投稿テーブル" {
		t.Errorf("Posts table comment = %q", table.Comment)
	}

	// status カラム（メタデータ付きJSON）の検証
	statusColumn, err := table.FindColumnByName("status")
	if err != nil {
		t.Error("Status column not found")
		return
	}

	if !statusColumn.HasEnhancedComment() {
		t.Error("Status column should have enhanced comment")
		return
	}

	metadata := statusColumn.GetMetadata()
	if values, exists := metadata["values"]; !exists {
		t.Error("Status column should have values metadata")
	} else {
		t.Logf("Status column values metadata: %v", values)
	}

	t.Logf("PostgreSQL Posts table validation completed successfully")
}

// validatePostgreSQLTagsTable PostgreSQL tagsテーブル検証
func validatePostgreSQLTagsTable(t *testing.T, table *Table) {
	// テーブルコメント（YAML形式）の検証
	if !table.HasEnhancedComment() {
		t.Error("Tags table should have enhanced comment")
		return
	}

	if table.LogicalName != "統合テストタグ" {
		t.Errorf("Tags table logical name = %q, expected '統合テストタグ'", table.LogicalName)
	}

	if table.GetDescription() != "ドライバー統合テスト用のタグテーブル" {
		t.Errorf("Tags table description = %q", table.GetDescription())
	}

	if table.GetPriority() != 2 {
		t.Errorf("Tags table priority = %d, expected 2", table.GetPriority())
	}

	t.Logf("PostgreSQL Tags table validation completed successfully")
}

// cleanupPostgreSQLTestData PostgreSQLテストデータクリーンアップ
func cleanupPostgreSQLTestData(t *testing.T, db *sql.DB) {
	cleanupQueries := []string{
		`DROP TABLE IF EXISTS integration_post_tags CASCADE`,
		`DROP TABLE IF EXISTS integration_tags CASCADE`,
		`DROP TABLE IF EXISTS integration_posts CASCADE`,
		`DROP TABLE IF EXISTS integration_users CASCADE`,
	}

	for _, query := range cleanupQueries {
		if _, err := db.Exec(query); err != nil {
			t.Logf("Warning: Failed to cleanup: %q, error: %v", query, err)
		}
	}

	t.Log("PostgreSQL test data cleanup completed")
}

// TestDriverCompatibility ドライバー互換性テスト
func TestDriverCompatibility(t *testing.T) {
	// ConfigurableDriverインターフェースとの互換性確認
	t.Run("ConfigurableDriver互換性", func(t *testing.T) {
		// SQLiteドライバーのテスト
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("Failed to open in-memory SQLite: %v", err)
		}
		defer db.Close()

		driver := sqlite.New(db)

		// ConfigurableDriverインターフェースの確認
		if configurableDriver, ok := driver.(drivers.ConfigurableDriver); ok {
			// 拡張コメント処理用の設定
			configurableDriver.SetLogicalNameConfig("|", true)
			configurableDriver.SetTableLogicalNameConfig("|", true)

			t.Log("SQLite driver implements ConfigurableDriver interface")
		} else {
			t.Error("SQLite driver should implement ConfigurableDriver interface")
		}
	})

	// 論理名設定との互換性確認
	t.Run("論理名設定互換性", func(t *testing.T) {
		// テスト用のダミードライバー設定
		testLogicalNameCompatibility(t)
	})
}

// testLogicalNameCompatibility 論理名設定互換性テスト
func testLogicalNameCompatibility(t *testing.T) {
	// 拡張コメント処理プロセッサの設定
	config := &ProcessingConfig{
		EnableValidation:   true,
		EnableSanitization: false,
		DefaultDelimiter:   "|",
		FallbackToLegacy:   true,
		StrictMode:         false,
	}

	processor := NewEnhancedCommentProcessorWithConfig(config)

	// テストコラム作成
	column := &Column{
		Name:    "test_column",
		Comment: `{"name": "テストカラム", "description": "互換性テスト用", "tags": ["test"]}`,
	}

	// 拡張コメント処理実行
	err := column.ProcessEnhancedComment(processor, "|", true)
	if err != nil {
		t.Fatalf("Enhanced comment processing failed: %v", err)
	}

	// 結果検証
	if !column.HasEnhancedComment() {
		t.Error("Column should have enhanced comment after processing")
	}

	if column.LogicalName != "テストカラム" {
		t.Errorf("Column logical name = %q, expected 'テストカラム'", column.LogicalName)
	}

	t.Log("Logical name compatibility test passed")
}

// TestDriverIntegrationPerformance ドライバー統合パフォーマンステスト
func TestDriverIntegrationPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// SQLiteパフォーマンステスト
	t.Run("SQLite統合パフォーマンス", func(t *testing.T) {
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("Failed to open in-memory SQLite: %v", err)
		}
		defer db.Close()

		// 大量テーブル作成
		createLargeSchemaInSQLite(t, db, 50, 10)

		// パフォーマンス測定
		start := time.Now()

		schema := &Schema{Name: "performance_test"}
		driver := sqlite.New(db)
		driver.SetLogicalNameConfig("|", true)

		// ドライバー解析
		err = driver.Analyze(schema)
		if err != nil {
			t.Fatalf("Driver analysis failed: %v", err)
		}

		driverElapsed := time.Since(start)

		// 拡張コメント処理
		start = time.Now()
		processor := NewEnhancedCommentProcessor()
		err = processSchemaWithEnhancedComments(schema, processor)
		if err != nil {
			t.Fatalf("Enhanced comment processing failed: %v", err)
		}

		commentElapsed := time.Since(start)

		t.Logf("Driver analysis time: %v", driverElapsed)
		t.Logf("Enhanced comment processing time: %v", commentElapsed)
		t.Logf("Total tables processed: %d", len(schema.Tables))

		// パフォーマンス基準チェック
		totalObjects := countSchemaObjects(schema)
		avgProcessingTime := commentElapsed / time.Duration(totalObjects)

		if avgProcessingTime > 10*time.Millisecond {
			t.Logf("Warning: Average processing time per object is high: %v", avgProcessingTime)
		}

		t.Logf("Average processing time per object: %v", avgProcessingTime)
	})
}

// createLargeSchemaInSQLite SQLiteで大規模スキーマ作成
func createLargeSchemaInSQLite(t *testing.T, db *sql.DB, tableCount, columnCount int) {
	for i := 0; i < tableCount; i++ {
		// テーブル作成
		var columns []string
		for j := 0; j < columnCount; j++ {
			columnType := "TEXT"
			if j == 0 {
				columnType = "INTEGER PRIMARY KEY"
			}
			columns = append(columns, fmt.Sprintf("col_%d %s", j, columnType))
		}

		query := fmt.Sprintf("CREATE TABLE perf_table_%d (%s)", i, strings.Join(columns, ", "))
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("Failed to create table perf_table_%d: %v", i, err)
		}
	}
}

// countSchemaObjects スキーマ内のオブジェクト数をカウント
func countSchemaObjects(schema *Schema) int {
	count := len(schema.Tables)
	for _, table := range schema.Tables {
		count += len(table.Columns)
		count += len(table.Indexes)
		count += len(table.Constraints)
		count += len(table.Triggers)
	}
	return count
}