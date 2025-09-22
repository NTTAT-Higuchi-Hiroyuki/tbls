# テスト設計書 - DBMSコメント機能拡張とMarkdown出力カスタマイズ

## 1. テスト概要

### 1.1 テスト目的

- DBMSコメントからの論理名分離機能の品質保証
- Markdown出力カスタマイズ機能の正常動作確認
- 既存機能への影響がないことの確認
- パフォーマンス劣化が5%以内であることの確認
- 全対応DBMSでの互換性確認

### 1.2 テスト範囲

**対象:**
- CommentParser: コメント区切り・論理名抽出
- MarkdownCustomizer: Markdown出力カスタマイズ
- ConfigExtension: 新設定項目の管理
- SchemaExtension: スキーマ構造体拡張
- TemplateDataProcessor: テンプレートデータ前処理
- 全対応DBMS（PostgreSQL, MySQL, SQL Server, SQLite, BigQuery, Snowflake, ClickHouse, Spanner, DynamoDB, MongoDB）

**対象外:**
- 既存機能の詳細テスト（回帰テストのみ）
- 他の出力形式（JSON, YAML等）への影響

### 1.3 テスト環境

- **言語**: Go 1.19+
- **テストフレームワーク**: Go標準testing + testify
- **DB環境**: Docker Compose による各DBMS環境
- **CI環境**: GitHub Actions（または同等）
- **カバレッジ目標**: 新規コード90%以上、全体で既存レベル維持

## 2. テストケース設計

### 2.1 CommentParserのテストケース

#### 2.1.1 正常系テスト

| ID   | テストケース名 | 入力データ | 期待結果 | 優先度 |
|------|----------------|------------|----------|---------|
| T001 | 標準区切り文字での論理名抽出 | "ユーザーID\|システム内でユーザーを一意に識別するID" | LogicalName: "ユーザーID", CleanComment: "システム内でユーザーを一意に識別するID" | High |
| T002 | カスタム区切り文字での論理名抽出 | "ユーザー名::ログイン時に使用する名前", separator: "::" | LogicalName: "ユーザー名", CleanComment: "ログイン時に使用する名前" | High |
| T003 | 日本語論理名の処理 | "顧客番号\|お客様を識別する番号" | 正しい日本語文字列の抽出 | High |
| T004 | 英数字論理名の処理 | "UserID\|Unique identifier for user" | LogicalName: "UserID", CleanComment: "Unique identifier for user" | Medium |
| T005 | 区切り文字なしのコメント処理 | "通常のコメント" | LogicalName: "", CleanComment: "通常のコメント" | High |

#### 2.1.2 異常系テスト

| ID   | テストケース名 | 入力データ | 期待結果 | 優先度 |
|------|----------------|------------|----------|---------|
| T101 | 空文字列入力 | "" | LogicalName: "", CleanComment: "" | High |
| T102 | nil入力 | nil | エラーハンドリング | High |
| T103 | 区切り文字のみ | "\|" | LogicalName: "", CleanComment: "" | Medium |
| T104 | 複数区切り文字 | "論理名\|説明\|追加情報" | 最初の区切り文字で分割 | Medium |
| T105 | 非常に長いコメント | 10000文字のコメント | 適切な処理とメモリ効率 | Low |

#### 2.1.3 境界値テスト

| ID   | テストケース名 | 入力データ | 期待結果 | 優先度 |
|------|----------------|------------|----------|---------|
| T201 | 最小長論理名 | "A\|説明" | LogicalName: "A" | Medium |
| T202 | 最大長論理名 | 255文字の論理名 | 正常処理 | Medium |
| T203 | 空の説明部分 | "論理名\|" | LogicalName: "論理名", CleanComment: "" | Medium |
| T204 | 特殊文字を含む論理名 | "論理名_01\|説明" | 正常処理 | Medium |

### 2.2 MarkdownCustomizerのテストケース

#### 2.2.1 正常系テスト

| ID   | テストケース名 | 入力データ | 期待結果 | 優先度 |
|------|----------------|------------|----------|---------|
| T301 | 基本的な列順序カスタマイズ | order: ["logical_name", "name", "type"] | 指定順序での列表示 | High |
| T302 | エイリアス適用 | aliases: {"name": "カラム名", "type": "データ型"} | 日本語エイリアスでの表示 | High |
| T303 | 論理名表示制御 | show_logical_name: true | 論理名列の表示 | High |
| T304 | 個別テーブル設定 | specific設定でのテーブル個別カスタマイズ | テーブル固有の表示 | Medium |
| T305 | 全体設定と個別設定の併用 | 全体設定 + 個別設定 | 個別設定の優先適用 | Medium |

#### 2.2.2 異常系テスト

| ID   | テストケース名 | 入力データ | 期待結果 | 優先度 |
|------|----------------|------------|----------|---------|
| T401 | 存在しない列名指定 | order: ["nonexistent_column"] | デフォルト順序でのフォールバック | High |
| T402 | 空の順序設定 | order: [] | デフォルト順序での表示 | Medium |
| T403 | nil設定 | config: nil | デフォルト動作 | High |
| T404 | 循環参照設定 | 不正な参照設定 | エラー検出とフォールバック | Medium |

### 2.3 ConfigExtensionのテストケース

#### 2.3.1 設定ファイル読み込みテスト

| ID   | テストケース名 | 入力データ | 期待結果 | 優先度 |
|------|----------------|------------|----------|---------|
| T501 | 基本設定の読み込み | 有効なYAML設定 | 正常な設定値の読み込み | High |
| T502 | 不正YAML | 構文エラーのあるYAML | エラーハンドリングとデフォルト値 | High |
| T503 | 部分設定 | comment設定のみ | 他はデフォルト値 | Medium |
| T504 | 空設定ファイル | 空のYAMLファイル | 全てデフォルト値 | Medium |
| T505 | 後方互換性 | 既存の.tbls.yml | 既存機能の維持 | High |

### 2.4 SchemaExtensionのテストケース

#### 2.4.1 論理名フィールドテスト

| ID   | テストケース名 | 入力データ | 期待結果 | 優先度 |
|------|----------------|------------|----------|---------|
| T601 | LogicalNameフィールド追加 | Table/Column構造体 | LogicalNameフィールドの存在 | High |
| T602 | SetLogicalNameFromCommentメソッド | コメント付きTable/Column | 論理名の正常設定 | High |
| T603 | GetLogicalNameOrFallbackメソッド | 論理名なしのTable/Column | 物理名でのフォールバック | High |
| T604 | JSON/YAMLシリアライズ | 論理名付きスキーマ | 正常なシリアライズ | Medium |

### 2.5 統合テストシナリオ

#### シナリオ1: PostgreSQL基本機能シナリオ

1. **前提条件**: PostgreSQLデータベースにコメント付きテーブル・カラムが存在
2. **テスト手順**:
   - Step 1: tbls.ymlでコメント区切り文字を"|"に設定
   - Step 2: Markdown出力でカラム順序を["logical_name", "name", "type", "comment"]に設定
   - Step 3: tbls docコマンド実行
   - Step 4: 生成されたMarkdownファイルの確認
3. **期待結果**: 論理名列が最初に表示され、指定順序でカラムが表示される

#### シナリオ2: マルチDBMS互換性シナリオ

1. **前提条件**: PostgreSQL, MySQL, SQL Serverに同等のテストデータが存在
2. **テスト手順**:
   - Step 1: 各DBMSで同一の設定ファイルを使用
   - Step 2: 各DBMSでtbls docコマンド実行
   - Step 3: 生成されたMarkdownファイルの比較
3. **期待結果**: 各DBMSで一貫した論理名表示とカスタマイズが適用される

#### シナリオ3: 大規模データベースパフォーマンスシナリオ

1. **前提条件**: 1000テーブル、10000カラムを持つデータベース
2. **テスト手順**:
   - Step 1: 新機能無効でのベンチマーク実行
   - Step 2: 新機能有効でのベンチマーク実行
   - Step 3: 実行時間とメモリ使用量の比較
3. **期待結果**: 性能劣化が5%以内に収まる

## 3. テストデータ設計

### 3.1 基本テストデータ

```yaml
# テスト用設定ファイル
comment:
  separator: "|"

markdown:
  columns:
    show_logical_name: true
    order: ["logical_name", "name", "type", "nullable", "comment"]
    aliases:
      logical_name: "論理名"
      name: "カラム名"
      type: "データ型"
      nullable: "NULL許可"
      comment: "説明"
```

```sql
-- PostgreSQL テストデータ
COMMENT ON TABLE users IS 'ユーザーテーブル|システムユーザーの基本情報を管理';
COMMENT ON COLUMN users.user_id IS 'ユーザーID|システム内でユーザーを一意に識別するID';
COMMENT ON COLUMN users.user_name IS 'ユーザー名|ログイン時に使用する名前';
COMMENT ON COLUMN users.email IS 'メールアドレス|ユーザーの連絡先メールアドレス';
```

### 3.2 境界値テストデータ

```go
// 境界値テストケース
type BoundaryTestCase struct {
    Name           string
    Comment        string
    Separator      string
    ExpectedLogical string
    ExpectedComment string
}

var boundaryTestCases = []BoundaryTestCase{
    {"最小論理名", "A|説明", "|", "A", "説明"},
    {"最大論理名", strings.Repeat("A", 255) + "|説明", "|", strings.Repeat("A", 255), "説明"},
    {"空説明", "論理名|", "|", "論理名", ""},
    {"特殊文字", "論理名_01|説明#1", "|", "論理名_01", "説明#1"},
}
```

### 3.3 エラーケーステストデータ

```go
// エラーケーステストデータ
type ErrorTestCase struct {
    Name        string
    Input       interface{}
    ExpectedErr string
}

var errorTestCases = []ErrorTestCase{
    {"nil入力", nil, "invalid input"},
    {"空文字列", "", ""},
    {"不正設定", map[string]interface{}{"invalid": "config"}, "validation error"},
}
```

## 4. パフォーマンステスト

### 4.1 負荷テスト

- **対象データ規模**:
  - 小規模: 10テーブル、100カラム
  - 中規模: 100テーブル、1000カラム
  - 大規模: 1000テーブル、10000カラム
- **測定指標**:
  - 実行時間（従来比105%以内）
  - メモリ使用量（従来比110%以内）
  - CPU使用率
- **測定方法**: Go benchmarkとpprof使用

### 4.2 ストレステスト

- **条件**:
  - 10000テーブル以上
  - 複雑なコメント構造
  - 全DBMSオブジェクト種別
- **期待動作**:
  - メモリリーク無し
  - 適切なエラーハンドリング
  - 処理の完了

## 5. セキュリティテスト

### 5.1 入力検証テスト

| ID   | テストケース | 入力データ | 期待結果 |
|------|--------------|------------|----------|
| S001 | SQLインジェクション | コメント内にSQL文 | 適切なエスケープ |
| S002 | XSS | HTMLタグを含むコメント | エスケープ処理 |
| S003 | 長大入力 | 異常に長いコメント | 適切な制限 |
| S004 | 制御文字 | 制御文字を含む入力 | フィルタリング |

### 5.2 情報漏洩テスト

- **機密コメント**: パスワード、APIキー等を含むコメントでの情報漏洩チェック
- **エラーメッセージ**: エラー時の機密情報露出チェック
- **ログ出力**: デバッグ情報での機密情報チェック

## 6. テスト実行計画

### 6.1 実行順序

1. **単体テスト** (各Phase実装時)
   - CommentParser単体テスト
   - MarkdownCustomizer単体テスト
   - ConfigExtension単体テスト
   - SchemaExtension単体テスト

2. **統合テスト** (Phase 3)
   - PostgreSQL統合テスト
   - MySQL統合テスト
   - SQL Server統合テスト
   - その他DBMS統合テスト

3. **E2Eテスト** (Phase 3-4)
   - 全体機能シナリオ
   - マルチDBMS互換性テスト

4. **パフォーマンステスト** (Phase 3-4)
   - 負荷テスト
   - ストレステスト

5. **セキュリティテスト** (Phase 4)
   - 入力検証テスト
   - 情報漏洩テスト

### 6.2 合格基準

- **単体テスト**: 全テストケース合格、カバレッジ90%以上
- **統合テスト**: 全DBMS対応テスト合格
- **パフォーマンス**: 既存比105%以内の実行時間
- **セキュリティ**: 全脆弱性テストクリア
- **回帰テスト**: 既存テストスイート100%パス

### 6.3 テスト環境構築

```yaml
# docker-compose.test.yml
version: '3.8'
services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: testdb
      POSTGRES_USER: testuser
      POSTGRES_PASSWORD: testpass

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_DATABASE: testdb
      MYSQL_USER: testuser
      MYSQL_PASSWORD: testpass
      MYSQL_ROOT_PASSWORD: rootpass

  sqlserver:
    image: mcr.microsoft.com/mssql/server:2019-latest
    environment:
      SA_PASSWORD: TestPass123
      ACCEPT_EULA: Y
```

## 7. リスクと対策

| リスク | 影響度 | 発生確率 | 対策 |
|--------|--------|----------|------|
| DBMS固有の互換性問題 | High | Medium | 各DBMS専用テスト強化、早期検出 |
| パフォーマンス劣化 | High | Medium | 継続的ベンチマーク、最適化実装 |
| 既存機能への影響 | High | Low | 包括的回帰テスト、段階的デプロイ |
| マルチバイト文字化け | Medium | Medium | UTF-8統一処理、文字エンコーディングテスト |
| 設定複雑化によるエラー | Medium | High | 設定検証機能、詳細なドキュメント |

## 8. テスト自動化戦略

### 8.1 自動化対象

- **全単体テスト**: CI実行時に自動実行
- **回帰テスト**: PR作成時に自動実行
- **パフォーマンステスト**: 夜間バッチで定期実行
- **マルチDBMSテスト**: リリース前に自動実行

### 8.2 CI/CD統合

```yaml
# .github/workflows/test.yml
name: Comprehensive Test Suite

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.19
      - run: go test ./... -v -race -coverprofile=coverage.out
      - run: go tool cover -html=coverage.out -o coverage.html

  integration-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:14
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v3
      - name: Run integration tests
        run: go test ./... -tags integration

  performance-tests:
    runs-on: ubuntu-latest
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v3
      - name: Run benchmarks
        run: go test -bench=. -benchmem ./...
```

### 8.3 品質ゲート

- **PR マージ条件**:
  - 全単体テスト合格
  - カバレッジ閾値達成
  - golangci-lint合格
  - 少なくとも1つのDBMSでの統合テスト合格

- **リリース条件**:
  - 全テストスイート合格
  - パフォーマンス基準達成
  - セキュリティテスト合格
  - マニュアルレビュー完了