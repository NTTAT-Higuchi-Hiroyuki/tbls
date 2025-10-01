# 詳細設計書 - DBMSコメント機能拡張とMarkdown出力カスタマイズ

## 1. アーキテクチャ概要

### 1.1 システム構成図

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     DBMS        │    │   Config        │    │   Schema        │
│  (PostgreSQL,   │───▶│  (.tbls.yml)    │───▶│   (構造体)      │
│   MySQL, etc.)  │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Comment        │    │   Markdown      │    │   Output        │
│  Parser         │◀───│   Customizer    │◀───│   Generator     │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 1.2 技術スタック

- **言語**: Go 1.19+（既存コードベース）
- **フレームワーク**: 標準ライブラリ + Cobra CLI
- **テンプレート**: Go template（既存Markdownテンプレート拡張）
- **設定**: YAML（既存.tbls.yml拡張）
- **テスト**: Go標準testing + testify（既存テスト環境維持）

## 2. コンポーネント設計

### 2.1 コンポーネント一覧

| コンポーネント名 | 責務 | 依存関係 |
|------------------|------|----------|
| CommentParser | コメント区切り・論理名抽出 | config.Config |
| MarkdownCustomizer | Markdown出力カスタマイズ | schema.Schema, config.Config |
| ConfigExtension | 新設定項目の管理 | 既存config.Config |
| SchemaExtension | スキーマ構造体拡張 | 既存schema.Schema |
| TemplateDataProcessor | テンプレートデータ前処理 | output/md パッケージ |

### 2.2 各コンポーネントの詳細

#### CommentParser

- **目的**: DBMSコメントから論理名とコメントを分離
- **公開インターフェース**:
  ```go
  type CommentParser interface {
      ParseComment(rawComment, separator string) (logicalName, cleanComment string)
      ExtractLogicalName(rawComment, separator string) string
      ExtractCleanComment(rawComment, separator string) string
      HasLogicalName(rawComment, separator string) bool
  }

  type DefaultCommentParser struct{}

  func (p *DefaultCommentParser) ParseComment(rawComment, separator string) (string, string)
  ```
- **内部実装方針**:
  - 区切り文字による文字列分割
  - UTF-8安全な文字列処理
  - 空文字・不正値のハンドリング

#### MarkdownCustomizer

- **目的**: Markdown出力の列順序・エイリアス・表示制御
- **公開インターフェース**:
  ```go
  type MarkdownCustomizer interface {
      CustomizeTableOutput(table *schema.Table, config *MarkdownConfig) *CustomizedTableData
      CustomizeColumnOutput(columns []*schema.Column, config *ColumnConfig) []*CustomizedColumnData
      ApplyAliases(fieldName string, config *AliasConfig) string
      GetDisplayOrder(objectType string, config *OrderConfig) []string
  }

  type CustomizedTableData struct {
      Table    *schema.Table
      Columns  []*CustomizedColumnData
      Order    []string
      Aliases  map[string]string
  }

  type CustomizedColumnData struct {
      Column      *schema.Column
      DisplayName string
      LogicalName string
      Show        bool
  }
  ```

#### ConfigExtension

- **目的**: 新設定項目の統合管理
- **公開インターフェース**:
  ```go
  // config/config.go に追加
  type Config struct {
      // 既存フィールド...
      Comment  *CommentConfig  `yaml:"comment,omitempty"`
      Markdown *MarkdownConfig `yaml:"markdown,omitempty"`
  }

  type CommentConfig struct {
      Separator string `yaml:"separator,omitempty"`
  }

  type MarkdownConfig struct {
      Database    *ObjectCustomConfig           `yaml:"database,omitempty"`
      Schemas     *ObjectCustomConfig           `yaml:"schemas,omitempty"`
      Tables      *TableCustomConfig            `yaml:"tables,omitempty"`
      Columns     *ColumnCustomConfig           `yaml:"columns,omitempty"`
      Views       *ObjectCustomConfig           `yaml:"views,omitempty"`
      Indexes     *ObjectCustomConfig           `yaml:"indexes,omitempty"`
      Constraints *ObjectCustomConfig           `yaml:"constraints,omitempty"`
      Functions   *ObjectCustomConfig           `yaml:"functions,omitempty"`
      Others      map[string]*ObjectCustomConfig `yaml:"others,omitempty"`
  }

  type ObjectCustomConfig struct {
      ShowLogicalName bool              `yaml:"show_logical_name,omitempty"`
      Order          []string           `yaml:"order,omitempty"`
      Aliases        map[string]string  `yaml:"aliases,omitempty"`
  }

  type TableCustomConfig struct {
      *ObjectCustomConfig
      Specific map[string]*ObjectCustomConfig `yaml:"specific,omitempty"`
  }

  type ColumnCustomConfig struct {
      *ObjectCustomConfig
      Specific map[string]*ObjectCustomConfig `yaml:"specific,omitempty"`
  }
  ```

#### SchemaExtension

- **目的**: スキーマ構造体への論理名フィールド追加
- **公開インターフェース**:
  ```go
  // schema/schema.go に追加
  type Table struct {
      // 既存フィールド...
      LogicalName string `json:"logicalName,omitempty"`
  }

  type Column struct {
      // 既存フィールド...
      LogicalName string `json:"logicalName,omitempty"`
  }

  type Index struct {
      // 既存フィールド...
      LogicalName string `json:"logicalName,omitempty"`
  }

  // 各構造体にメソッド追加
  func (t *Table) SetLogicalNameFromComment(separator string)
  func (t *Table) GetLogicalNameOrFallback() string
  func (c *Column) SetLogicalNameFromComment(separator string)
  func (c *Column) GetLogicalNameOrFallback() string
  ```

#### TemplateDataProcessor

- **目的**: テンプレートファイルを変更せずにデータ前処理でカスタマイズ実現
- **公開インターフェース**:
  ```go
  type TemplateDataProcessor interface {
      ProcessTableData(table *schema.Table, config *MarkdownConfig) interface{}
      ProcessColumnHeaders(headers []string, config *MarkdownConfig) []string
      ApplyCustomOrder(data []interface{}, order []string) []interface{}
      ApplyAliases(data interface{}, aliases map[string]string) interface{}
  }
  ```

## 3. データフロー

### 3.1 データフロー図

```
[DBMS]
   │
   ▼ (1) Raw Comment
[CommentParser]
   │
   ▼ (2) Parsed Data (LogicalName + CleanComment)
[Schema Objects]
   │
   ▼ (3) Extended Schema
[MarkdownCustomizer]
   │
   ▼ (4) Customized Data
[TemplateDataProcessor]
   │
   ▼ (5) Processed Template Data
[Existing Templates] ← テンプレートファイル自体は変更なし
   │
   ▼ (6) Final Markdown
[File Output]
```

### 3.2 データ変換

1. **Raw Comment → Parsed Data**
   - 入力: `"ユーザーID|システム内でユーザーを一意に識別するID"`
   - 処理: 区切り文字「|」で分割
   - 出力: LogicalName=`"ユーザーID"`, CleanComment=`"システム内でユーザーを一意に識別するID"`

2. **Schema Objects → Extended Schema**
   - 入力: schema.Column{Name: "user_id", Comment: "..."}
   - 処理: SetLogicalNameFromComment()実行
   - 出力: schema.Column{Name: "user_id", LogicalName: "ユーザーID", Comment: "..."}

3. **Extended Schema → Customized Data**
   - 入力: []*schema.Column
   - 処理: 設定に基づく順序変更・エイリアス適用
   - 出力: CustomizedTableData

## 4. APIインターフェース

### 4.1 内部API

```go
// パッケージ間の主要インターフェース

// config パッケージ
func (c *Config) GetCommentSeparator() string
func (c *Config) GetMarkdownConfig() *MarkdownConfig
func (c *Config) IsLogicalNameEnabled() bool

// schema パッケージ
func (s *Schema) ApplyCommentParsing(separator string)
func (t *Table) SetLogicalNameFromComment(separator string)
func (c *Column) SetLogicalNameFromComment(separator string)

// output/md パッケージ
func (m *Md) makeTableTemplateDataWithCustomization(table *schema.Table, config *MarkdownConfig) interface{}
func (m *Md) processTemplateData(data interface{}, processor TemplateDataProcessor) interface{}
```

### 4.2 外部API

```go
// CLI コマンド拡張
// cmd/doc.go での新オプション
var (
    commentSeparator   string
    enableLogicalName  bool
    customColumnOrder  []string
)
```

## 5. エラーハンドリング

### 5.1 エラー分類

- **設定エラー**: 不正なYAML設定値
  - 対処: デフォルト値使用 + 警告ログ
- **コメント解析エラー**: 想定外のコメント形式
  - 対処: 従来通りの処理継続 + デバッグログ
- **テンプレートエラー**: Markdown生成時のエラー
  - 対処: エラー内容を含むMarkdown出力

### 5.2 エラー通知

```go
type CommentParsingError struct {
    Object    string
    Comment   string
    Separator string
    Err       error
}

func (e *CommentParsingError) Error() string {
    return fmt.Sprintf("comment parsing failed for %s: %v", e.Object, e.Err)
}
```

## 6. セキュリティ設計

### 6.1 データ保護

- 入力値検証: YAML設定値の型・範囲チェック
- ログ出力制御: デバッグ時のみ詳細情報出力

## 7. テスト設計

**詳細なテスト設計については、`/spec:test-design_jp`コマンドを実行してテスト設計書を作成してください。**

テスト設計書では以下の内容を定義します：

- **単体テスト**: CommentParser, MarkdownCustomizerの各メソッド
- **統合テスト**: 各DBMS（PostgreSQL, MySQL, SQL Server）での動作確認
- **回帰テスト**: 既存機能への影響なし確認


## 9. デプロイメント

### 9.1 デプロイ構成

- **段階リリース**:
  - Phase 1: コメント分離機能
  - Phase 2: Markdown出力カスタマイズ
  - Phase 3: 全DBMSオブジェクト対応
  - Phase 4: 高度な設定・エラーハンドリング

### 9.2 設定管理

- **後方互換性**: 既存.tbls.yml設定の完全サポート
- **設定移行**: 自動的なデフォルト値適用
- **設定検証**: 起動時の設定値チェック

## 10. 実装上の注意事項

### 10.1 アーキテクチャ原則

- **既存コード最小影響**: 既存のConfig, Schema構造体への拡張のみ
- **段階的実装**: Phase分けによる段階的デリバリ
- **テスト駆動**: 新機能は全てテストファースト

### 10.2 コーディング規約

- **Go標準**: gofmt, golangci-lint準拠
- **命名規則**: 既存コードベースとの一貫性維持
- **エラーハンドリング**: Go慣用的なエラー処理パターン

### 10.3 パッケージ構成

```
config/
├── config.go          (既存 + CommentConfig, MarkdownConfig追加)
├── comment.go         (新規: コメント解析設定)
└── markdown.go        (新規: Markdown出力設定)

schema/
├── schema.go          (既存 + LogicalName フィールド追加)
├── comment.go         (新規: コメント解析ロジック)
└── logical_name.go    (新規: 論理名処理ロジック)

output/md/
├── md.go              (既存 + カスタマイズ機能追加)
├── customizer.go      (新規: 出力カスタマイズ)
└── templates/
    ├── table.md.tmpl  (既存 + 論理名・カスタマイズ対応)
    └── index.md.tmpl  (既存テンプレート)
```

### 10.4 データベースドライバー対応

各データベースドライバーでの実装パターン：

```go
// drivers/postgres/postgres.go 例
func (p *Postgres) analyzeTableComment(table *schema.Table) error {
    // 既存のコメント取得処理
    if table.Comment != "" && p.config.GetCommentSeparator() != "" {
        table.SetLogicalNameFromComment(p.config.GetCommentSeparator())
    }
    return nil
}
```

### 10.5 マイグレーション戦略

- **設定ファイル**: 自動的なデフォルト値適用（明示的移行不要）
- **出力フォーマット**: 既存テンプレートとの互換性維持
- **CLI引数**: 新オプション追加（既存オプション維持）
