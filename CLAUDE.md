# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

このプロジェクトの担当者は、日本語でのコミュニケーションを希望しています。ドキュメント・コメントは全て日本語で記載してください。

## プロジェクト概要

tblsは、CI-Friendlyなデータベースドキュメント生成ツールです。単一バイナリで動作し、多数のデータベースをサポートしています。

## 開発コマンド

### ビルドとテスト
```bash
# バイナリのビルド
make build

# 全テストの実行（各種データベースタグ付き）
go test ./... -tags 'bq clickhouse dynamo mariadb mongodb mssql mysql postgres redshift snowflake spanner sqlite' -coverprofile=coverage.out -covermode=count

# 特定のパッケージのテスト実行
go test ./datasource -v
go test ./drivers/mysql -v

# 特定のテストケースの実行
go test -run TestAnalyzeSchema ./datasource -v

# Lintの実行
make lint
# または
golangci-lint run ./...

# CI環境での完全なテスト実行
make ci
```

### 開発用データベースのセットアップ
```bash
# テスト用データベースの起動（Docker Compose使用）
docker compose up -d

# テスト用データベースの初期化
make db

# SQLiteのみの初期化
make db_sqlite
```

### ドキュメント生成とテスト
```bash
# サンプルドキュメントの生成
make doc

# ドキュメントとデータベースの差分確認
make testdoc
```

## アーキテクチャとコード構造

### コマンド構造（cmd/）
- Cobraライブラリを使用したCLI実装
- 各サブコマンドは独立したファイルで実装
- 外部サブコマンドのサポート（PATH上の`tbls-*`実行可能ファイル）

### データベースドライバー（drivers/）
- 各データベース固有の実装を含む
- 共通インターフェース: `Driver`を実装
- 新しいドライバー追加時は`drivers.go`への登録が必要

### 出力フォーマット（output/）
- 各出力形式は独立したパッケージ
- テンプレートエンジンを使用（Go template）
- カスタムテンプレートのサポート

### スキーマ表現（schema/）
- `Schema`構造体がデータベース全体を表現
- `Table`、`Column`、`Relation`等の基本構造体
- JSONとYAMLでのシリアライズ/デシリアライズ対応
- **拡張コメント処理機能**: JSON/YAML/Legacy形式のコメント解析とメタデータ抽出

### 設定管理（config/）
- `.tbls.yml`ファイルでの設定
- 環境変数の展開サポート（`${}`構文）
- DSN、フィルタ、Lint規則等の設定

## 重要な実装パターン

### エラーハンドリング
- 基本的に早期リターンでエラーを伝播
- CLIレベルでのエラー表示とexit処理

### テストデータ
- `testdata/`ディレクトリにテスト用SQL、設定ファイル、期待値を配置
- Goldenファイルパターンの使用（`.golden`拡張子）

### データベース接続
- DSN形式: `driver://user:pass@host:port/dbname?option=value`
- 環境変数での設定も可能（`TBLS_DSN`）

### 外部コマンド/ドライバー
- PATH上の`tbls-*`実行可能ファイルを外部サブコマンドとして認識
- 標準入出力を通じた通信

## デバッグとトラブルシューティング

### 詳細ログの表示
```bash
# --debugオプションでデバッグ情報を表示
tbls doc --debug postgres://...
```

### 特定のテーブルのみ処理
```bash
# --table/-tオプションでテーブルを指定
tbls doc postgres://... -t users -t posts
```

### 設定ファイルの検証
```bash
# --config/-cで設定ファイルを明示的に指定
tbls doc --config custom.yml postgres://...
```

## 拡張コメント処理機能

### 概要
データベースのテーブル、カラム、インデックス、制約、トリガーのコメントを構造化されたメタデータとして解析・処理する機能です。

### サポートされるコメント形式

#### 1. JSON形式
```json
{
  "name": "ユーザー",
  "description": "システムユーザーの基本情報",
  "tags": ["マスター", "認証"],
  "priority": 1,
  "deprecated": false,
  "metadata": {
    "owner": "開発チーム",
    "version": "1.0"
  }
}
```

#### 2. YAML形式
```yaml
name: ユーザー
description: システムユーザーの基本情報
tags:
  - マスター
  - 認証
priority: 1
deprecated: false
metadata:
  owner: 開発チーム
  version: "1.0"
```

#### 3. Legacy形式（従来互換）
```
論理名|説明文
```

### 拡張コメント処理のテスト実行

```bash
# 拡張コメント機能のテスト実行
go test ./schema -run Enhanced -v

# パフォーマンステスト
go test ./schema -run BenchmarkEnhanced -bench=. -v

# 統合テスト（データベース接続が必要）
go test ./schema -tags integration -run DriverIntegration -v

# エラーハンドリングテスト
go test ./schema -run ErrorHandling -v
```

### 主要コンポーネント

#### パーサー
- `JSONParser`: JSON形式コメントの解析
- `YAMLParser`: YAML形式コメントの解析  
- `LegacyParser`: 従来形式（`論理名|説明`）の解析

#### プロセッサー
- `EnhancedCommentProcessor`: 複数パーサーの統合管理
- `CommentValidator`: コメント内容の検証
- パーサー優先度の自動選択

#### エラーハンドリング
- `EnhancedCommentError`: 詳細なエラー情報とコンテキスト
- `ErrorRecoveryManager`: 自動回復戦略
- `FallbackParsingStrategy`: 解析失敗時のフォールバック処理

#### ドライバー統合
- `EnhancedCommentDriverAdapter`: ドライバーとの統合アダプター
- 既存ドライバーとの透過的な連携
- 統計情報とパフォーマンス監視

### 設定

拡張コメント処理は以下のように設定できます：

```go
config := &ProcessingConfig{
    EnableValidation:   true,    // バリデーション有効化
    EnableSanitization: false,   // サニタイゼーション有効化
    DefaultDelimiter:   "|",     // デフォルト区切り文字
    FallbackToLegacy:   true,    // Legacy形式へのフォールバック
    StrictMode:         false,   // 厳密モード
    ProcessingTimeout:  5000,    // タイムアウト（ミリ秒）
}
```

### データベース固有の考慮事項

#### PostgreSQL
- テーブル・カラム・インデックス・制約のコメントをフルサポート
- JSON/YAML/Legacy形式すべてに対応

#### MySQL
- 文字エンコーディングの問題対策でサニタイゼーション推奨
- コメント長制限（1024文字）に注意

#### SQLite
- テーブルコメントは未サポート（データベース制限）
- カラムコメントは制限的サポート

#### その他のデータベース
- 各データベースの制限に応じた互換性チェック機能を提供

## 注意事項

- 新機能追加時は必ずテストを追加
- データベース固有の処理はドライバー層に実装
- 出力フォーマット固有の処理は出力層に実装
- 環境変数の展開は`${VAR_NAME}`形式を使用
- Lint機能追加時は`config/lint.go`を更新
- **拡張コメント機能使用時は適切な設定でパフォーマンス影響を最小化**

## リポジトリ運用

- ISSUE、PRは、NTTAT-Higuchi-Hiroyuki/tblsリポジトリに対して行うこと