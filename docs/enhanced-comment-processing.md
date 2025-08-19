# 拡張コメント処理機能

## 概要

tblsの拡張コメント処理機能は、データベースのコメントをJSON、YAML、または従来の形式で構造化されたメタデータとして解析し、より豊富なドキュメント生成を可能にします。

## 機能一覧

### コメント形式サポート

1. **JSON形式**: 完全な構造化メタデータ
2. **YAML形式**: 可読性の高い構造化メタデータ
3. **Legacy形式**: 既存の `論理名|説明` 形式との互換性

### 対象オブジェクト

- テーブル
- カラム
- インデックス
- 制約
- トリガー

### 主要機能

- **パーサー優先度の自動選択**: コメント内容に基づく最適なパーサーの選択
- **フォールバック処理**: 解析失敗時の自動回復
- **バリデーション**: コメント内容の検証
- **サニタイゼーション**: セキュリティ対応
- **エラーハンドリング**: 詳細なエラー情報と回復戦略
- **パフォーマンス監視**: 処理統計とベンチマーク

## 使用例

### JSON形式のコメント

```sql
COMMENT ON TABLE users IS '{
  "name": "ユーザー",
  "description": "システムユーザーの基本情報",
  "tags": ["マスター", "認証"],
  "priority": 1,
  "metadata": {
    "owner": "開発チーム",
    "created": "2023-01-01"
  }
}';

COMMENT ON COLUMN users.id IS '{
  "name": "ユーザーID",
  "description": "システム内でユーザーを一意に識別するID",
  "tags": ["PK", "自動生成"],
  "priority": 1
}';
```

### YAML形式のコメント

```sql
COMMENT ON TABLE posts IS 'name: 投稿
description: ユーザーの投稿を管理するテーブル
tags:
  - コンテンツ
  - 管理
priority: 2
deprecated: false';
```

### Legacy形式のコメント（従来互換）

```sql
COMMENT ON TABLE comments IS 'コメント|投稿に対するコメントを管理';
COMMENT ON COLUMN comments.content IS 'コメント内容|ユーザーが投稿したコメントの本文';
```

## 設定例

### 基本設定

```go
// 推奨設定（本番環境）
config := &ProcessingConfig{
    EnableValidation:   true,
    EnableSanitization: true,
    DefaultDelimiter:   "|",
    FallbackToLegacy:   true,
    StrictMode:         false,
    ProcessingTimeout:  5000,
}

// 開発環境設定
config := &ProcessingConfig{
    EnableValidation:   true,
    EnableSanitization: false,
    DefaultDelimiter:   "|",
    FallbackToLegacy:   true,
    StrictMode:         true,  // エラーを厳密にチェック
    ProcessingTimeout:  10000,
}
```

### ドライバー統合

```go
// PostgreSQLドライバーでの使用例
driver := postgres.New(db)
adapter := NewEnhancedCommentDriverAdapter(config)

// 通常のドライバー解析
err := driver.Analyze(schema)
if err != nil {
    return err
}

// 拡張コメント処理を追加
err = adapter.ProcessSchemaComments(schema)
if err != nil {
    return err
}

// 処理統計の取得
stats := adapter.GetProcessingStatistics(schema)
log.Printf("処理完了: %s", stats.GetSummary())
```

## テスト

### 基本テスト

```bash
# 全拡張コメント機能のテスト
go test ./schema -run Enhanced -v

# 特定機能のテスト
go test ./schema -run TestJSONParser -v
go test ./schema -run TestYAMLParser -v
go test ./schema -run TestLegacyParser -v
```

### パフォーマンステスト

```bash
# ベンチマークテスト
go test ./schema -bench=BenchmarkEnhanced -v

# パフォーマンス詳細分析
go test ./schema -run TestDetailedPerformanceAnalysis -v

# メモリリーク検出
go test ./schema -run TestMemoryLeakDetection -v
```

### 統合テスト

```bash
# データベース統合テスト（Docker環境が必要）
export POSTGRES_TEST_ENABLED=true
go test ./schema -tags integration -run DriverIntegration -v

# SQLite統合テスト
go test ./schema -tags integration -run SQLite -v
```

## エラーハンドリング

### エラータイプ

1. **PARSING**: JSON/YAMLの構文エラー
2. **VALIDATION**: コメント内容の検証エラー
3. **PROCESSING**: 処理実行時のエラー
4. **CONFIGURATION**: 設定エラー
5. **TIMEOUT**: 処理タイムアウト
6. **MEMORY**: メモリ使用量エラー
7. **SCHEMA**: スキーマ構造エラー

### 自動回復戦略

1. **FallbackParsingStrategy**: JSON/YAML解析失敗時のLegacy形式へのフォールバック
2. **CommentSanitizationStrategy**: 危険なコンテンツの自動サニタイゼーション

### エラー監視

```go
reporter := NewDefaultErrorReporter()

// エラー処理中にレポートを収集
err := processor.ProcessComment(comment)
if err != nil {
    reporter.ReportError(err)
}

// 回復処理
if manager.CanRecover(err) {
    recovered, recoveryErr := manager.TryRecover(err, context)
    if recoveryErr == nil {
        reporter.ReportRecovery(err, recovered)
    }
}

// エラーサマリーの取得
summary := reporter.GetErrorSummary()
log.Printf("エラー統計: %+v", summary)
```

## データベース固有の制限事項

### PostgreSQL ✅ フルサポート
- テーブル、カラム、インデックス、制約のコメント
- JSON/YAML/Legacy形式すべて対応
- 特別な制限なし

### MySQL ⚠️ 制限あり
- コメント長制限: 1024文字
- 文字エンコーディング問題の可能性
- **推奨**: サニタイゼーション有効化

### SQLite ⚠️ 制限あり
- テーブルコメント: 未サポート
- カラムコメント: 制限的サポート
- **推奨**: カラムコメントのみ使用

### その他のデータベース
各データベースの制限に応じた互換性チェック機能を提供。詳細は `ValidateDriverCompatibility` 関数を参照。

## パフォーマンス考慮事項

### ベンチマーク結果（参考値）

- **JSONParser**: ~36.5µs/operation
- **LegacyParser**: ~1.9ns/operation
- **YAMLParser**: ~74.5µs/operation

### 最適化のヒント

1. **大規模スキーマ**: バッチ処理とタイムアウト設定の調整
2. **パフォーマンス重視**: Legacy形式の優先使用
3. **メモリ使用量**: 並行処理数の制限

## トラブルシューティング

### よくある問題

1. **JSON解析エラー**: 構文チェックとエスケープ処理の確認
2. **YAML解析エラー**: インデントとエンコーディングの確認
3. **タイムアウト**: 処理時間制限の調整
4. **メモリ不足**: バッチサイズの縮小

### デバッグ

```go
// デバッグモードでの実行
config.StrictMode = true
config.ProcessingTimeout = 60000  // 長めのタイムアウト

// エラー詳細の出力
if err != nil {
    enhancedErr, ok := err.(*EnhancedCommentError)
    if ok {
        log.Printf("Error details: %+v", enhancedErr)
        log.Printf("Context: %+v", enhancedErr.Context)
        log.Printf("Suggestions: %v", enhancedErr.Suggestions)
    }
}
```

## 今後の展望

1. **カスタムパーサー**: プラグイン式パーサーの追加
2. **出力フォーマット連携**: ドキュメント生成への統合
3. **設定UI**: Web UIでの設定管理
4. **AI連携**: 自動コメント生成・改善提案

## 関連ファイル

- `schema/enhanced_comment_*.go`: 核となる実装
- `schema/enhanced_comment_*_test.go`: テストコード
- `CLAUDE.md`: 開発者向け詳細ドキュメント
- `docs/enhanced-comment-processing.md`: このファイル