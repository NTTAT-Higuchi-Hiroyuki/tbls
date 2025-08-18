# テスト設計書 - データベースコメントからの論理名・説明読み取り機能強化

## 1. テスト概要

### 1.1 テスト目的

- 新しいコメント解析システム（JSON、YAML対応）の動作検証
- 既存機能との後方互換性確認
- エラーハンドリングとフォールバック処理の検証
- パフォーマンス要件（1ms/テーブル）の達成確認
- セキュリティ対策の有効性検証

### 1.2 テスト範囲

- **対象**: 
  - EnhancedCommentProcessor および配下の全パーサー
  - CommentValidator による検証機能
  - 既存schema.Column/Table.SetLogicalNameFromComment拡張
  - 新規追加のIndex/View/Constraintコメント処理
  - 設定システムの拡張部分

- **対象外**: 
  - データベースドライバーのコメント取得処理（既存機能）
  - テンプレートエンジン側の論理名表示（既存機能）

### 1.3 テスト環境

- **環境**: Go 1.23.8+、既存のtblsテスト環境
- **依存関係**: 
  - テスト用SQLiteデータベース
  - MockデータベースドライバーInterface
  - 各種パーサーのテストデータセット

## 2. テストケース設計

### 2.1 CommentParserインターフェースのテストケース

#### 2.1.1 正常系テスト

| ID   | テストケース名 | 入力データ | 期待結果 | 優先度 |
|------|----------------|------------|----------|---------|
| T001 | JSON形式の基本解析 | `{"name": "ユーザー名", "description": "ユーザーの表示名"}` | LogicalName="ユーザー名", Description="ユーザーの表示名" | High |
| T002 | JSON形式タグ付き解析 | `{"name": "顧客ID", "description": "顧客識別子", "tags": ["primary", "required"]}` | Tags=["primary", "required"] | High |
| T003 | YAML形式基本解析 | `name: テーブル名, description: テーブルの説明` | LogicalName="テーブル名", Description="テーブルの説明" | High |
| T004 | Legacy形式継続対応 | `論理名\|説明文` | LogicalName="論理名", Description="説明文" | High |
| T005 | 空コメント処理 | `""` | LogicalName="", Description="" | Medium |
| T006 | 日本語文字対応 | `{"name": "価格（税込）", "description": "商品の税込価格"}` | 正常に日本語処理 | High |

#### 2.1.2 異常系テスト

| ID   | テストケース名 | 入力データ | 期待結果 | 優先度 |
|------|----------------|------------|----------|---------|
| T101 | 不正JSON形式 | `{"name": "test"` | ParseError、フォールバック処理 | High |
| T102 | 深すぎるJSON構造 | 6階層以上のネストJSON | エラーまたは制限による打ち切り | Medium |
| T103 | 長すぎる論理名 | `{"name": "[256文字の文字列]"}` | ValidationError | High |
| T104 | 長すぎる説明 | `{"description": "[2049文字の文字列]"}` | 説明の切り詰め処理 | Medium |
| T105 | 禁止文字を含む論理名 | `{"name": "<script>alert('xss')</script>"}` | サニタイズ処理 | High |
| T106 | NULL文字を含む入力 | `"論理名\x00|説明"` | 適切な処理または拒否 | High |

#### 2.1.3 境界値テスト

| ID   | テストケース名 | 入力データ | 期待結果 | 優先度 |
|------|----------------|------------|----------|---------|
| T201 | 論理名255文字ちょうど | `{"name": "[255文字の文字列]"}` | 正常処理 | Medium |
| T202 | 説明2048文字ちょうど | `{"description": "[2048文字の文字列]"}` | 正常処理 | Medium |
| T203 | 最小限JSON | `{"name": "a"}` | 正常処理 | Medium |
| T204 | 区切り文字のエッジケース | `"論理名\|\|説明文"` | パイプ文字のエスケープ処理 | Medium |

### 2.2 EnhancedCommentProcessorのテストケース

#### 2.2.1 正常系テスト

| ID   | テストケース名 | 入力データ | 期待結果 | 優先度 |
|------|----------------|------------|----------|---------|
| T301 | パーサー自動選択(JSON) | JSON形式コメント | JSONParserが選択される | High |
| T302 | パーサー自動選択(Legacy) | 従来形式コメント | LegacyParserが選択される | High |
| T303 | フォールバック処理 | JSON解析失敗→Legacy解析成功 | Legacy結果を返却 | High |
| T304 | 設定による無効化 | JSON無効設定 + JSONコメント | LegacyParserにフォールバック | Medium |
| T305 | 重複チェック機能 | 同一論理名の複数使用 | 警告ログ出力 | Medium |

#### 2.2.2 異常系テスト

| ID   | テストケース名 | 入力データ | 期待結果 | 優先度 |
|------|----------------|------------|----------|---------|
| T401 | 全パーサー解析失敗 | 解析不能なコメント | 空のCommentDataまたはデフォルト値 | High |
| T402 | 設定ファイル不正 | 不正な設定値 | デフォルト設定で継続 | Medium |
| T403 | メモリ不足状況 | 巨大なコメントデータ | 適切なエラーハンドリング | Low |

### 2.3 CommentValidatorのテストケース

#### 2.3.1 正常系テスト

| ID   | テストケース名 | 入力データ | 期待結果 | 優先度 |
|------|----------------|------------|----------|---------|
| T501 | 論理名文字種検証 | 英数字、ひらがな、カタカナ、漢字 | 正常通過 | High |
| T502 | 説明サニタイズ | HTMLタグを含む説明 | HTMLタグ除去 | High |
| T503 | 論理名長さ検証 | 規定内文字数 | 正常通過 | Medium |

#### 2.3.2 異常系テスト

| ID   | テストケース名 | 入力データ | 期待結果 | 優先度 |
|------|----------------|------------|----------|---------|
| T601 | SQLインジェクション検出 | `'; DROP TABLE users; --` | サニタイズまたは拒否 | High |
| T602 | XSS検出 | `<script>alert('xss')</script>` | サニタイズ処理 | High |
| T603 | パストラバーサル検出 | `../../../etc/passwd` | 拒否または無害化 | High |

### 2.4 統合テストシナリオ

#### シナリオ1: 既存機能との互換性

1. **前提条件**: 従来の区切り文字ベース設定
2. **テスト手順**:
   - Step 1: 既存の`論理名|説明`形式でテーブル作成
   - Step 2: EnhancedCommentProcessorで処理
   - Step 3: 従来のSetLogicalNameFromCommentと同じ結果を確認
3. **期待結果**: 既存機能と完全に同じ動作

#### シナリオ2: 新機能の段階的有効化

1. **前提条件**: JSON機能が無効な設定
2. **テスト手順**:
   - Step 1: JSON形式コメントを入力
   - Step 2: レガシーパーサーにフォールバック確認
   - Step 3: JSON機能を有効化
   - Step 4: 同じコメントがJSON解析される確認
3. **期待結果**: 設定による段階的有効化の動作

#### シナリオ3: 大規模データベースでの性能

1. **前提条件**: 1000テーブル、10000カラムのテストDB
2. **テスト手順**:
   - Step 1: 各テーブル・カラムに多様なコメント設定
   - Step 2: 全体スキーマ解析の実行
   - Step 3: 処理時間とメモリ使用量の計測
3. **期待結果**: 1テーブル平均1ms以内、メモリ1.1倍以内

## 3. テストデータ設計

### 3.1 マスタデータ

```json
{
  "testData": {
    "valid": {
      "json_basic": "{\"name\": \"ユーザー\", \"description\": \"ユーザー情報\"}",
      "json_with_tags": "{\"name\": \"商品\", \"description\": \"商品マスタ\", \"tags\": [\"master\", \"public\"]}",
      "yaml_basic": "name: 注文, description: 注文情報",
      "legacy_basic": "顧客|顧客マスタテーブル",
      "japanese_chars": "{\"name\": \"価格（税抜）\", \"description\": \"商品の税抜価格を格納\"}"
    },
    "invalid": {
      "malformed_json": "{\"name\": \"test\"",
      "too_long_name": "{\"name\": \"" + "a".repeat(256) + "\"}",
      "xss_attempt": "{\"name\": \"<script>alert('xss')</script>\"}",
      "sql_injection": "'; DROP TABLE users; --",
      "null_chars": "test\x00name|description"
    },
    "boundary": {
      "max_name_length": "{\"name\": \"" + "a".repeat(255) + "\"}",
      "max_desc_length": "{\"description\": \"" + "a".repeat(2048) + "\"}",
      "empty_comment": "",
      "only_spaces": "   ",
      "unicode_boundary": "{\"name\": \"🚀テスト\"}"
    }
  }
}
```

### 3.2 モックデータ

```go
// データベースドライバーのモック
type MockDriver struct {
    tables []TableWithComments
    columns []ColumnWithComments
}

type TableWithComments struct {
    Name    string
    Comment string
    Type    string // TABLE, VIEW, etc.
}

type ColumnWithComments struct {
    TableName string
    Name      string
    Comment   string
    Type      string
}
```

## 4. パフォーマンステスト

### 4.1 負荷テスト

- **測定項目**: 
  - 1テーブルあたりの解析時間
  - 1000テーブル一括処理時間
  - メモリ使用量（処理前後の差分）
- **合格基準**:
  - 平均解析時間: 1ms/テーブル以内
  - メモリ増加: 既存処理の1.1倍以内
  - CPU使用率: 既存処理と同等

### 4.2 ストレステスト

- **テスト条件**:
  - 最大10,000テーブル、100,000カラム
  - 各コメントにJSONデータ（平均200文字）
- **期待される挙動**:
  - メモリリークの発生なし
  - 処理時間の線形増加
  - エラー率5%以下（設計上のタイムアウト等を含む）

## 5. セキュリティテスト

### 5.1 入力検証テスト

| テストケース | 攻撃パターン | 期待される防御 |
|-------------|-------------|---------------|
| SQLインジェクション | `'; DROP TABLE users; --` | コメント内容のサニタイズ |
| XSS攻撃 | `<script>alert('xss')</script>` | HTML特殊文字エスケープ |
| パストラバーサル | `../../../etc/passwd` | パス文字の検出・拒否 |
| JSONインジェクション | `{"name": "test\",\"malicious\":\"evil"}` | JSON構造の検証 |
| バッファオーバーフロー | 極端に長い文字列 | 長さ制限による保護 |

### 5.2 制限事項テスト

| 制限項目 | テスト値 | 期待動作 |
|---------|---------|----------|
| JSON深度 | 6階層のネスト | 制限による拒否 |
| 解析タイムアウト | 1秒超の処理 | タイムアウトエラー |
| 文字列長制限 | 論理名256文字 | ValidationError |
| 同時解析数 | 大量並行処理 | リソース制御 |

## 6. テスト実行計画

### 6.1 実行順序

1. **単体テスト** (フェーズ1)
   - CommentParser各実装
   - CommentValidator
   - EnhancedCommentProcessor

2. **統合テスト** (フェーズ2)
   - パーサーチェーンの動作
   - 既存機能との統合
   - エラーハンドリング

3. **パフォーマンステスト** (フェーズ3)
   - 負荷テスト
   - メモリ使用量測定
   - ストレステスト

4. **セキュリティテスト** (フェーズ4)
   - 入力検証
   - セキュリティ制限
   - 脆弱性スキャン

### 6.2 合格基準

- **機能テスト**: 全テストケースの95%以上合格
- **カバレッジ**: コード行カバレッジ90%以上、分岐カバレッジ85%以上
- **パフォーマンス**: 設計書記載の性能基準達成
- **セキュリティ**: 重要度Highのセキュリティテスト100%合格
- **互換性**: 既存機能のリグレッションゼロ

## 7. リスクと対策

| リスク | 影響度 | 発生確率 | 対策 |
|--------|--------|----------|------|
| JSONライブラリの予期しない動作 | High | Low | 複数のJSONテストケースとモックライブラリ |
| 既存機能の互換性破綻 | High | Medium | 包括的な回帰テストスイート |
| パフォーマンス劣化 | Medium | Medium | ベンチマークテストとプロファイリング |
| セキュリティホール | High | Low | セキュリティ専門家によるレビュー |
| 設定システムの複雑化 | Medium | High | 設定パターンの網羅テスト |

## 8. テスト自動化戦略

### 8.1 自動化対象

- **完全自動化**:
  - 全ての単体テスト
  - 統合テストの大部分
  - パフォーマンス基準チェック
  - セキュリティテストの基本項目

- **半自動化**:
  - 大規模データでのストレステスト
  - メモリリーク検出
  - 設定ファイルパターンテスト

### 8.2 CI/CD統合

```yaml
# GitHub Actions設定例
name: Enhanced Comment Processing Tests

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - name: Unit Tests
        run: go test ./schema/... -v -race -coverprofile=coverage.out
      
  integration-tests:
    runs-on: ubuntu-latest
    needs: unit-tests
    steps:
      - name: Integration Tests
        run: go test ./... -tags=integration -v
        
  performance-tests:
    runs-on: ubuntu-latest
    needs: unit-tests
    steps:
      - name: Performance Benchmark
        run: go test ./schema/... -bench=BenchmarkCommentProcessing -benchmem
        
  security-tests:
    runs-on: ubuntu-latest
    needs: unit-tests
    steps:
      - name: Security Tests
        run: go test ./schema/... -tags=security -v
```

### 8.3 継続的品質管理

- **コードカバレッジ**: codecov.ioとの統合
- **性能監視**: ベンチマーク結果の履歴追跡
- **セキュリティスキャン**: gosecによる自動脆弱性検出
- **依存関係監視**: Dependabotによる依存ライブラリ更新