# テスト設計書 - DBMSコメント機能拡張とMarkdown出力カスタマイズ

## 1. 概要

### 1.1 テスト目的

本テスト設計書は、Task 3.5「統合テストとE2Eテスト」の完了条件を満たすための包括的なテストプランを定義します。Phase 1-3で実装されたコメント分離機能、論理名処理、Markdownカスタマイズ機能の品質保証と回帰防止を実現します。

### 1.2 テスト戦略の基本方針

- **統合優先**: ModifySchema統合処理による全DBMS一括対応テスト
- **回帰防止**: 既存機能の完全な動作保証
- **日本語重視**: UTF-8マルチバイト文字での包括的テスト
- **性能保証**: 5%以内のパフォーマンス劣化制限
- **現実的運用**: 実データベースでのE2Eテスト

### 1.3 テスト対象範囲

#### 実装済み機能（Phase 1-3）
- ✅ CommentParser（コメント分離処理）
- ✅ LogicalNameProcessor（論理名処理）
- ✅ MarkdownCustomizer（出力カスタマイズ）
- ✅ ModifySchema統合処理（全DBMS対応）
- ✅ ConfigExtension（設定管理）

#### テスト対象外
- 非実装の出力フォーマット（JSON、YAML、PlantUML等）
- コメント機能をサポートしないDBMS
- Phase 4の高度な設定・エラーハンドリング

## 2. テストカテゴリ

### 2.1 基本機能テスト

#### 2.1.1 コメント分離機能テスト

**目的**: CommentParserの基本動作検証

- **正常ケース**
  - 区切り文字「|」での論理名とコメント分離
  - 区切り文字「::」での論理名とコメント分離
  - 区切り文字「--」での論理名とコメント分離
  - 区切り文字なしでの従来動作維持

- **境界値ケース**
  - 空のコメント文字列
  - 論理名のみ（区切り文字後が空）
  - コメントのみ（区切り文字前が空）
  - 区切り文字が複数回出現する場合

- **エラーケース**
  - 無効な区切り文字（空文字、null）
  - 極端に長いコメント文字列

#### 2.1.2 論理名処理機能テスト

**目的**: LogicalNameProcessorの統合処理検証

- **正常ケース**
  - テーブル論理名の設定と取得
  - カラム論理名の設定と取得
  - インデックス論理名の設定と取得
  - 制約論理名の設定と取得
  - 論理名フォールバック動作（物理名への回帰）

- **統合処理ケース**
  - ModifySchema経由での一括処理
  - 複数オブジェクト種別の同時処理
  - 設定なし時の従来動作維持

#### 2.1.3 Markdownカスタマイズ機能テスト

**目的**: MarkdownCustomizerの出力制御検証

- **列順序変更**
  - デフォルト順序での表示
  - カスタム順序での表示
  - 不正な列名指定時の処理

- **エイリアス機能**
  - 英語列名から日本語名への変換
  - 複数エイリアスの同時適用
  - 存在しない列名への無害な処理

- **表示制御**
  - 論理名表示/非表示切り替え
  - 特定オブジェクト種別の個別設定

### 2.2 日本語処理テスト

#### 2.2.1 UTF-8マルチバイト文字テスト

**目的**: 日本語環境での正常動作確認

- **日本語コメント処理**
  - ひらがな、カタカナ、漢字の論理名
  - 日本語説明文でのコメント分離
  - 混在文字（日本語+英語+数字）での処理

- **文字エンコーディング**
  - UTF-8での正確な文字分割
  - バイト境界での安全な処理
  - 文字化け防止の確認

#### 2.2.2 日本語設定ファイルテスト

**目的**: .tbls.ymlでの日本語設定処理確認

- **設定値の日本語化**
  - エイリアス設定での日本語表示名
  - カスタム列順序での日本語列名
  - 設定ファイルのUTF-8解析

### 2.3 統合機能テスト

#### 2.3.1 コンポーネント間連携テスト

**目的**: Phase 1-3実装の統合動作検証

- **データフロー確認**
  1. DBMS → CommentParser → LogicalNameProcessor
  2. LogicalNameProcessor → MarkdownCustomizer → Template
  3. 設定ファイル → 各コンポーネント → 最終出力

- **設定連携確認**
  - Comment設定とMarkdown設定の協調動作
  - 階層的設定の優先順位確認
  - デフォルト値の適切な適用

#### 2.3.2 ModifySchema統合処理テスト

**目的**: 全DBMS一括対応機能の検証

- **統合処理の動作確認**
  - config.ModifySchema関数での論理名処理
  - エラー時の警告出力と処理継続
  - 既存機能への影響なし確認

### 2.4 エラーケーステスト

#### 2.4.1 設定エラー処理テスト

**目的**: 不正設定での安全な動作確認

- **設定ファイルエラー**
  - 不正なYAML形式
  - 存在しない列名指定
  - 循環参照設定

- **フォールバック動作**
  - デフォルト値への安全な回帰
  - エラーログの適切な出力
  - 処理継続の確認

#### 2.4.2 DBMS接続エラー処理テスト

**目的**: データベース関連エラーでの堅牢性確認

- **接続エラー**
  - データベース接続失敗時の処理
  - 権限不足でのアクセス失敗処理

- **データエラー**
  - 不正なコメントデータでの処理
  - 欠損データでの処理継続

### 2.5 パフォーマンステスト

#### 2.5.1 処理性能測定テスト

**目的**: 5%以内パフォーマンス劣化制限の確認

- **大規模データベーステスト**
  - 1000+ テーブルでの処理時間測定
  - 10000+ カラムでの処理時間測定
  - メモリ使用量の測定

- **回帰テスト**
  - 既存機能との処理時間比較
  - 機能無効時の完全同一性能確認
  - ベンチマーク結果の蓄積

#### 2.5.2 メモリ効率テスト

**目的**: メモリ使用量の適切性確認

- **メモリ使用量測定**
  - 論理名処理でのメモリ増加量
  - カスタマイズ処理でのメモリ使用量
  - ガベージコレクションの適切性

### 2.6 回帰テスト

#### 2.6.1 既存機能動作保証テスト

**目的**: 既存機能の完全な動作保証

- **既存テストの完全パス**
  - 全既存テストケースの実行
  - テスト結果の100%パス確認
  - テストカバレッジの維持確認

- **出力互換性確認**
  - 設定なし時の出力完全同一性
  - 既存テンプレートでの出力確認

## 3. DBMS別テスト計画

### 3.1 PostgreSQL（主要対象）

#### 3.1.1 機能テスト
- **コメント取得**: COMMENT ON文での設定値確認
- **拡張オブジェクト**: ビュー、インデックス、制約での動作確認
- **権限テスト**: 読み取り専用権限での動作確認

#### 3.1.2 実データテスト
```sql
-- テストデータ例
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE users IS 'ユーザー|システム利用者の情報を管理するテーブル';
COMMENT ON COLUMN users.id IS 'ユーザーID|システム内でユーザーを一意に識別するID';
COMMENT ON COLUMN users.name IS 'ユーザー名|ユーザーの表示名';
COMMENT ON COLUMN users.email IS 'メールアドレス|ログインおよび通知に使用';
```

### 3.2 MySQL（主要対象）

#### 3.2.1 機能テスト
- **コメント取得**: COMMENT構文での設定値確認
- **文字セット**: UTF-8での日本語コメント確認
- **ストレージエンジン**: InnoDB、MyISAMでの動作確認

#### 3.2.2 実データテスト
```sql
-- テストデータ例
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT 'プロダクトID|商品を一意に識別するID',
    name VARCHAR(255) NOT NULL COMMENT '商品名|販売する商品の名称',
    price DECIMAL(10,2) COMMENT '価格|商品の販売価格（税抜）',
    category_id INT COMMENT 'カテゴリID|商品カテゴリの識別子'
) COMMENT='商品|ECサイトで販売する商品情報';
```

### 3.3 SQL Server（拡張プロパティ）

#### 3.3.1 機能テスト
- **拡張プロパティ**: sys.extended_propertiesでの取得確認
- **スキーマ構成**: マルチスキーマでの動作確認
- **データ型**: SQL Server固有型での動作確認

#### 3.3.2 実データテスト
```sql
-- テストデータ例
CREATE TABLE orders (
    id INT IDENTITY(1,1) PRIMARY KEY,
    user_id INT NOT NULL,
    total_amount MONEY,
    order_date DATETIME2 DEFAULT GETDATE()
);

EXEC sp_addextendedproperty
    @name = N'MS_Description',
    @value = N'注文|顧客からの商品注文情報を管理',
    @level0type = N'SCHEMA', @level0name = N'dbo',
    @level1type = N'TABLE', @level1name = N'orders';
```

### 3.4 SQLite（ローカルテスト）

#### 3.4.1 機能テスト
- **制限確認**: SQLiteでのコメント機能制限確認
- **ファイル処理**: ローカルファイルでの動作確認
- **軽量性**: 最小構成での動作確認

### 3.5 その他DBMS（BigQuery、Snowflake等）

#### 3.5.1 機能テスト
- **クラウドDBMS**: BigQuery、Snowflakeでの基本動作確認
- **アナリティクス**: ClickHouseでの動作確認
- **NoSQL**: MongoDBでの可能範囲確認

## 4. テストデータ設計

### 4.1 日本語コメントテストデータ

#### 4.1.1 基本パターン
```yaml
test_comments:
  basic_japanese:
    - comment: "ユーザーID|システム内でユーザーを一意に識別するID"
      expected_logical: "ユーザーID"
      expected_clean: "システム内でユーザーを一意に識別するID"

  mixed_languages:
    - comment: "User ID|ユーザー識別番号（システム全体で唯一）"
      expected_logical: "User ID"
      expected_clean: "ユーザー識別番号（システム全体で唯一）"

  special_characters:
    - comment: "データ作成日時|レコードが作成された日時（JST）"
      expected_logical: "データ作成日時"
      expected_clean: "レコードが作成された日時（JST）"
```

#### 4.1.2 エラーケースパターン
```yaml
error_test_comments:
  empty_logical:
    - comment: "|説明のみのコメント"
      expected_logical: ""
      expected_clean: "説明のみのコメント"

  empty_description:
    - comment: "論理名のみ|"
      expected_logical: "論理名のみ"
      expected_clean: ""

  multiple_separators:
    - comment: "論理名|説明|追加情報"
      expected_logical: "論理名"
      expected_clean: "説明|追加情報"
```

### 4.2 設定ファイルテストデータ

#### 4.2.1 基本設定パターン
```yaml
# .tbls.yml テスト用設定
comment:
  separator: "|"

markdown:
  tables:
    show_logical_name: true
    order: ["Name", "LogicalName", "Type", "Comment"]
    aliases:
      Name: "物理名"
      LogicalName: "論理名"
      Type: "データ型"
      Comment: "説明"

  columns:
    show_logical_name: true
    order: ["Name", "LogicalName", "Type", "Nullable", "Default", "Comment"]
    aliases:
      Name: "カラム名"
      LogicalName: "論理名"
      Type: "型"
      Nullable: "NULL可"
      Default: "デフォルト"
      Comment: "説明"
```

#### 4.2.2 個別カスタマイズ設定
```yaml
markdown:
  tables:
    specific:
      users:
        order: ["LogicalName", "Name", "Type", "Comment"]
        aliases:
          LogicalName: "項目名"
          Name: "カラム名"
```

### 4.3 大規模データベーステストデータ

#### 4.3.1 パフォーマンステスト用データ
- **テーブル数**: 1000テーブル
- **カラム数**: 平均10カラム/テーブル（総10000カラム）
- **コメント付きオブジェクト**: 70%（論理名処理対象）
- **日本語コメント**: 50%（文字処理テスト）

## 5. 実装対象機能の検証マトリックス

### 5.1 Phase 1: コメント分離機能

| 機能 | 単体テスト | 統合テスト | E2Eテスト | 状態 |
|------|------------|------------|-----------|------|
| CommentParser.ParseComment | ✅ | ✅ | ✅ | 実装済み |
| CommentParser.ExtractLogicalName | ✅ | ✅ | ✅ | 実装済み |
| CommentParser.ExtractCleanComment | ✅ | ✅ | ✅ | 実装済み |
| CommentParser.HasLogicalName | ✅ | ✅ | ✅ | 実装済み |
| LogicalNameProcessor.ProcessSchema | ✅ | ✅ | ✅ | 実装済み |

### 5.2 Phase 2: Markdownカスタマイズ機能

| 機能 | 単体テスト | 統合テスト | E2Eテスト | 状態 |
|------|------------|------------|-----------|------|
| MarkdownCustomizer.CustomizeTableOutput | ✅ | ✅ | ✅ | 実装済み |
| MarkdownCustomizer.CustomizeColumnOutput | ✅ | ✅ | ✅ | 実装済み |
| MarkdownCustomizer.ApplyAliases | ✅ | ✅ | ✅ | 実装済み |
| MarkdownCustomizer.GetDisplayOrder | ✅ | ✅ | ✅ | 実装済み |
| テンプレートデータ前処理 | ✅ | ✅ | ✅ | 実装済み |

### 5.3 Phase 3: DBMS統合

| DBMS | ModifySchema統合 | 論理名処理 | 日本語テスト | 状態 |
|------|------------------|------------|--------------|------|
| PostgreSQL | ✅ | ✅ | ✅ | 実装済み |
| MySQL | ✅ | ✅ | ✅ | 実装済み |
| SQL Server | ✅ | ✅ | ✅ | 実装済み |
| SQLite | ✅ | ✅ | ✅ | 実装済み |
| BigQuery | ✅ | ✅ | △ | 実装済み |
| Snowflake | ✅ | ✅ | △ | 実装済み |
| ClickHouse | ✅ | ✅ | △ | 実装済み |

## 6. 期待される結果と検証方法

### 6.1 基本機能の期待結果

#### 6.1.1 コメント分離処理
```go
// 入力
comment := "ユーザーID|システム内でユーザーを一意に識別するID"
separator := "|"

// 期待される出力
logicalName := "ユーザーID"
cleanComment := "システム内でユーザーを一意に識別するID"
```

#### 6.1.2 Markdown出力カスタマイズ
```markdown
# 期待される出力例（日本語エイリアス + 論理名表示）
## users
ユーザー|システム利用者の情報を管理するテーブル

| 項目名 | カラム名 | 型 | NULL可 | デフォルト | 説明 |
| ---- | ---- | ---- | ---- | ---- | ---- |
| ユーザーID | id | SERIAL | NO |  | システム内でユーザーを一意に識別するID |
| ユーザー名 | name | VARCHAR(100) | YES |  | ユーザーの表示名 |
```

### 6.2 性能要件の検証

#### 6.2.1 ベンチマーク基準
- **基準値**: 機能無効時の処理時間を100%とする
- **許容範囲**: 105%以内（5%以内の劣化）
- **測定対象**: `tbls doc` コマンドの総実行時間

#### 6.2.2 測定手法
```bash
# ベンチマーク実行例
# 機能無効時（ベースライン）
time tbls doc postgres://test_db --config baseline.yml

# 機能有効時（測定対象）
time tbls doc postgres://test_db --config enhanced.yml

# 結果比較と判定
performance_ratio = enhanced_time / baseline_time
# performance_ratio <= 1.05 であることを確認
```

### 6.3 エラー処理の期待結果

#### 6.3.1 設定エラー時の動作
```yaml
# 不正設定例
markdown:
  columns:
    order: ["NonexistentColumn", "Name", "Type"]

# 期待される動作
# 1. 警告ログの出力: "Unknown column 'NonexistentColumn', ignoring"
# 2. 有効な列のみでの処理継続
# 3. 最終出力の正常生成
```

#### 6.3.2 DBMS接続エラー時の動作
```bash
# 接続失敗時
tbls doc postgres://invalid_connection

# 期待される動作
# 1. 適切なエラーメッセージ出力
# 2. 非ゼロ終了コード
# 3. スタックトレースの抑制（本番環境）
```

## 7. テストケース実装計画

### 7.1 単体テストケース

#### 7.1.1 CommentParser テストケース
```go
func TestCommentParser_ParseComment(t *testing.T) {
    tests := []struct {
        name           string
        comment        string
        separator      string
        expectedLogical string
        expectedClean   string
    }{
        {
            name: "日本語論理名とコメント",
            comment: "ユーザーID|システム内でユーザーを一意に識別するID",
            separator: "|",
            expectedLogical: "ユーザーID",
            expectedClean: "システム内でユーザーを一意に識別するID",
        },
        // 追加テストケース...
    }
    // テスト実行ロジック
}
```

#### 7.1.2 MarkdownCustomizer テストケース
```go
func TestMarkdownCustomizer_CustomizeTableOutput(t *testing.T) {
    // カスタマイズ設定のテストケース
    config := &MarkdownConfig{
        Tables: &TableCustomConfig{
            ObjectCustomConfig: &ObjectCustomConfig{
                ShowLogicalName: true,
                Order: []string{"LogicalName", "Name", "Type", "Comment"},
                Aliases: map[string]string{
                    "Name": "カラム名",
                    "LogicalName": "論理名",
                },
            },
        },
    }
    // テスト実行とアサーション
}
```

### 7.2 統合テストケース

#### 7.2.1 ModifySchema統合テスト
```go
func TestLogicalNameIntegrationWithModifySchema(t *testing.T) {
    // 実際のスキーマデータでの統合テスト
    schema := createTestSchema() // テスト用スキーマ作成
    config := createTestConfig() // テスト用設定作成

    // ModifySchema実行
    err := config.ModifySchema(schema)
    assert.NoError(t, err)

    // 論理名が適切に設定されていることを確認
    assert.Equal(t, "ユーザーID", schema.Tables[0].Columns[0].LogicalName)
}
```

#### 7.2.2 E2Eテストケース（データベース実行）
```go
func TestE2E_PostgreSQL_JapaneseComments(t *testing.T) {
    if testing.Short() {
        t.Skip("E2E test skipped in short mode")
    }

    // 実データベースセットアップ
    db := setupTestDatabase(t)
    defer cleanupTestDatabase(db)

    // 日本語コメント付きテーブル作成
    createTablesWithJapaneseComments(db)

    // tbls実行
    result := runTblsDoc(db, testConfig)

    // 出力結果検証
    assertMarkdownContains(t, result, "ユーザーID")
    assertMarkdownContains(t, result, "システム内でユーザーを一意に識別するID")
}
```

### 7.3 パフォーマンステストケース

#### 7.3.1 ベンチマークテスト
```go
func BenchmarkLogicalNameProcessing(b *testing.B) {
    schema := createLargeTestSchema(1000, 10) // 1000テーブル、各10カラム
    config := createPerformanceTestConfig()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        config.ModifySchema(schema)
    }
}

func BenchmarkMarkdownGeneration(b *testing.B) {
    // Markdown生成のベンチマーク
}
```

### 7.4 回帰テストケース

#### 7.4.1 既存機能保証テスト
```bash
#!/bin/bash
# regression_test.sh

echo "Running all existing tests..."
go test ./... -tags 'bq clickhouse dynamo mariadb mongodb mssql mysql postgres redshift snowflake spanner sqlite'

echo "Running enhanced functionality tests..."
go test ./config -run TestComment
go test ./schema -run TestLogicalName
go test ./output/md -run TestCustomizer

echo "Running E2E tests..."
go test -tags integration ./tests/integration/...
```

## 8. テスト環境とツール

### 8.1 テスト環境

#### 8.1.1 ローカル開発環境
- **Go**: 1.19+
- **Docker**: テスト用データベースコンテナ
- **Make**: テストタスク自動化

#### 8.1.2 CI/CD環境
- **GitHub Actions**: 自動テスト実行
- **マルチDBMS**: 並行テスト実行
- **カバレッジ**: テストカバレッジ測定

### 8.2 テストツールとライブラリ

#### 8.2.1 Go標準ライブラリ
- `testing`: 基本テストフレームワーク
- `testing/quick`: プロパティベーステスト
- `net/http/httptest`: HTTPテスト

#### 8.2.2 サードパーティライブラリ
- `testify/assert`: アサーションライブラリ
- `testify/mock`: モックライブラリ
- `testify/suite`: テストスイート

### 8.3 データベーステスト環境

#### 8.3.1 Docker Compose設定
```yaml
version: '3.8'
services:
  postgres-test:
    image: postgres:13
    environment:
      POSTGRES_DB: test_logical_names
      POSTGRES_USER: test_user
      POSTGRES_PASSWORD: test_pass
    ports:
      - "15432:5432"

  mysql-test:
    image: mysql:8.0
    environment:
      MYSQL_DATABASE: test_logical_names
      MYSQL_USER: test_user
      MYSQL_PASSWORD: test_pass
      MYSQL_ROOT_PASSWORD: root_pass
    ports:
      - "13306:3306"
    command: --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
```

#### 8.3.2 テストデータ初期化
```sql
-- init_test_data.sql
-- PostgreSQL用テストデータ
CREATE TABLE test_users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE test_users IS 'テストユーザー|機能テスト用のユーザー情報テーブル';
COMMENT ON COLUMN test_users.id IS 'ユーザーID|システム内でユーザーを一意に識別するID';
COMMENT ON COLUMN test_users.name IS 'ユーザー名|ユーザーの表示名';
COMMENT ON COLUMN test_users.email IS 'メールアドレス|ログインおよび通知に使用するメールアドレス';
COMMENT ON COLUMN test_users.created_at IS '作成日時|レコードが作成された日時';
```

## 9. テスト実行計画

### 9.1 テスト実行フェーズ

#### 9.1.1 Phase 1: 単体テスト
- **対象**: 個別コンポーネントの機能検証
- **実行頻度**: コード変更毎
- **所要時間**: 5分以内
- **カバレッジ目標**: 新規コード90%以上

#### 9.1.2 Phase 2: 統合テスト
- **対象**: コンポーネント間の連携検証
- **実行頻度**: 機能単位完成毎
- **所要時間**: 15分以内
- **成功基準**: 全テストパス + 既存機能への影響なし

#### 9.1.3 Phase 3: E2Eテスト
- **対象**: 実データベースでの完全動作検証
- **実行頻度**: Phase完成毎
- **所要時間**: 30分以内
- **成功基準**: 全DBMSでの正常動作 + 性能要件達成

#### 9.1.4 Phase 4: パフォーマンステスト
- **対象**: 大規模データでの性能検証
- **実行頻度**: リリース前
- **所要時間**: 60分以内
- **成功基準**: 5%以内の性能劣化

### 9.2 テスト実行コマンド

#### 9.2.1 基本テスト実行
```bash
# 全単体テスト
make test

# 統合テスト（Docker環境必要）
make test-integration

# E2Eテスト（実データベース環境必要）
make test-e2e

# パフォーマンステスト
make test-performance
```

#### 9.2.2 DBMS別テスト実行
```bash
# PostgreSQL専用テスト
go test ./tests/postgres -v

# MySQL専用テスト
go test ./tests/mysql -v

# 全DBMS並行テスト
make test-all-dbms
```

#### 9.2.3 カバレッジ測定
```bash
# カバレッジ測定実行
go test ./... -coverprofile=coverage.out -covermode=count

# カバレッジレポート生成
go tool cover -html=coverage.out -o coverage.html

# カバレッジ閾値チェック
go tool cover -func=coverage.out | grep total | awk '{print $3}' | grep -E '[8-9][0-9]%|100%'
```

## 10. 品質基準と成功判定

### 10.1 機能品質基準

#### 10.1.1 基本機能の完全性
- ✅ **コメント分離**: 全パターンでの正常動作
- ✅ **論理名処理**: 全DBMSオブジェクトでの動作
- ✅ **Markdownカスタマイズ**: 設定通りの出力生成
- ✅ **日本語処理**: UTF-8での完全対応

#### 10.1.2 統合処理の正確性
- ✅ **ModifySchema統合**: 全DBMSでの一括処理
- ✅ **エラー処理**: 適切な警告とフォールバック
- ✅ **後方互換性**: 既存機能の完全保持

### 10.2 性能品質基準

#### 10.2.1 処理性能
- **基準**: 機能無効時比5%以内の劣化
- **測定対象**: `tbls doc`コマンド総実行時間
- **テストケース**: 1000テーブル、10000カラムでの処理

#### 10.2.2 メモリ効率
- **基準**: 同規模データでのメモリ使用量10%以内増加
- **測定方法**: Go runtime.ReadMemStats()による測定
- **判定**: 大規模データでのOOM発生なし

### 10.3 コード品質基準

#### 10.3.1 静的解析
- **Linter**: golangci-lint エラー・警告なし
- **フォーマット**: gofmt 準拠
- **import**: goimports による整理

#### 10.3.2 テストカバレッジ
- **新規コード**: 90%以上
- **全体**: 既存レベル維持
- **重要パス**: 100%（エラーハンドリング含む）

### 10.4 ドキュメント品質

#### 10.4.1 APIドキュメント
- **godoc**: 全公開関数・構造体へのコメント
- **使用例**: 主要機能の使用例記載
- **エラー処理**: エラーケースの説明

#### 10.4.2 設定ドキュメント
- **設定例**: 実用的な設定ファイル例
- **移行ガイド**: 既存設定からの移行方法
- **トラブルシューティング**: よくある問題と解決方法

## 11. リスク管理と対策

### 11.1 技術リスク

#### 11.1.1 データベース差異リスク
- **リスク**: DBMS間でのコメント仕様差異
- **対策**: 各DBMS専用テストケースの充実
- **検出方法**: E2Eテストでの早期発見

#### 11.1.2 文字エンコーディングリスク
- **リスク**: 日本語処理での文字化け
- **対策**: UTF-8専用テストケースの実装
- **検出方法**: マルチバイト文字境界でのテスト

#### 11.1.3 性能劣化リスク
- **リスク**: 新機能による処理性能低下
- **対策**: 継続的ベンチマーク実行
- **検出方法**: 閾値チェックによる早期警告

### 11.2 運用リスク

#### 11.2.1 設定複雑化リスク
- **リスク**: 設定項目増加によるユーザー混乱
- **対策**: 段階的テスト実行とドキュメント充実
- **検出方法**: ユーザビリティテストでの確認

#### 11.2.2 回帰バグリスク
- **リスク**: 既存機能への予期しない影響
- **対策**: 包括的回帰テストスイート
- **検出方法**: 既存テスト100%パスの必須化

### 11.3 スケジュールリスク

#### 11.3.1 テスト期間超過リスク
- **リスク**: 想定以上のテスト時間要求
- **対策**: 優先度付きテスト実行計画
- **対応方針**: 高優先度テストの完全実施

## 12. テスト完了基準

### 12.1 必須完了条件

#### 12.1.1 基本機能テスト
- [ ] 全単体テストのパス（100%）
- [ ] 全統合テストのパス（100%）
- [ ] 主要DBMS（PostgreSQL、MySQL、SQL Server）でのE2Eテスト完了
- [ ] 日本語処理テストの完全パス

#### 12.1.2 性能・品質テスト
- [ ] パフォーマンステストの基準達成（5%以内劣化）
- [ ] 全既存テストの継続パス（回帰テスト）
- [ ] テストカバレッジの目標達成（新規コード90%以上）
- [ ] 静的解析の完全パス（golangci-lint エラーなし）

#### 12.1.3 統合処理テスト
- [ ] ModifySchema統合処理での全DBMS動作確認
- [ ] エラーハンドリングの適切性確認
- [ ] 設定ファイル連携の正常動作確認

### 12.2 推奨完了条件

#### 12.2.1 拡張テスト
- [ ] 全対応DBMS（BigQuery、Snowflake等）での基本動作確認
- [ ] 大規模データベース（1000+テーブル）でのパフォーマンステスト
- [ ] エッジケース（異常データ）での堅牢性テスト

#### 12.2.2 ドキュメント
- [ ] テスト結果レポートの作成
- [ ] 発見した問題と対策の文書化
- [ ] 今後のテスト改善提案の作成

## 13. テスト報告とフォローアップ

### 13.1 テスト結果報告

#### 13.1.1 定量的結果
- テスト実行数と成功率
- パフォーマンス測定結果（ベースライン比較）
- カバレッジ測定結果
- 発見したバグ数と分類

#### 13.1.2 定性的評価
- 実装品質の評価
- ユーザビリティの評価
- 保守性の評価
- 将来拡張性の評価

### 13.2 継続的改善

#### 13.2.1 テストケース拡充
- 新発見パターンのテスト追加
- エッジケースの継続発見・追加
- パフォーマンステストの精緻化

#### 13.2.2 自動化改善
- CI/CDパイプラインの最適化
- テスト実行時間の短縮
- 並行実行の効率化

---

この包括的なテスト設計書により、Task 3.5「統合テストとE2Eテスト」を体系的かつ効率的に実行し、DBMSコメント機能拡張とMarkdown出力カスタマイズ機能の高品質な実装を保証します。