# 要求定義書 - DBMSコメント機能拡張とMarkdown出力カスタマイズ

## 1. 目的

コメント機能を持つDBMSに対して、論理名とコメントを区切り文字で分離して管理し、Markdown出力において列の順序変更、列名のエイリアス表示、論理名表示を可能にする拡張機能を実装する。

## 2. 機能要求

### 2.1 必須要求

#### 2.1.1 DBMSコメント解析機能
- [ ] コメント内に区切り文字が含まれる場合、前半を論理名、後半をコメントとして分離できること
- [ ] 区切り文字がtbls.ymlで設定可能であること
- [ ] 区切り文字が設定されていない場合は従来通りの動作を維持すること
- [ ] 論理名が取得できない場合は物理カラム名をフォールバックとして使用すること

#### 2.1.2 Markdown出力カスタマイズ機能
- [ ] tbls.ymlでMarkdown出力時の列の順序を指定できること
- [ ] tbls.ymlでMarkdown出力時の列名のエイリアスを設定できること
- [ ] 論理名をMarkdown出力に表示できること
- [ ] 設定された順序に従って列が表示されること
- [ ] エイリアスが設定されている場合は物理名の代わりにエイリアスが表示されること

#### 2.1.3 設定管理機能
- [ ] tbls.yml（.yaml）ファイルでコメント区切り文字を設定できること
- [ ] tbls.yml（.yaml）ファイルで列順序を設定できること
- [ ] tbls.yml（.yaml）ファイルで列名エイリアスを設定できること
- [ ] 設定が無効または不正な場合はデフォルト動作を維持すること

### 2.2 対象のコメントフィールド
- [ ] **対象コメントフィールドの拡張**
  - 対象オブジェクト
    - DATABASE（データベース）
    - TABLE(テーブル)
    - INDEX（インデックス）
    - SCHEMA（スキーマ）
    - TABLESPACE（テーブルスペース）
    - COLLATION（照合順序）
    - CONVERSION（変換）
    - LANGUAGE（プロシージャ言語）
    - TABLE（テーブル本体）
    - COLUMN（カラム）
    - CONSTRAINT（制約 — PRIMARY KEY, UNIQUE, CHECK など）
    - VIEW（ビュー）
    - MATERIALIZED VIEW（マテリアライズドビュー）
    - FOREIGN TABLE（外部テーブル）
    - FOREIGN DATA WRAPPER（外部データラッパ）
    - SERVER（外部サーバ）
    - POLICY（行レベルセキュリティポリシー）
    - PUBLICATION / SUBSCRIPTION（論理レプリケーション関連）

## 3. 非機能要求

### 3.1 互換性要求
- 既存のtblsの機能に影響を与えないこと
- 既存の設定ファイルとの後方互換性を維持すること
- コメント区切り文字が設定されていない場合は従来通りの動作をすること

### 3.2 出力形式制約
- Markdown出力のみに適用されること
- 他の出力形式（JSON、YAML等）には影響しないこと

## 4. 制約事項

### 4.1 技術的制約

- Go言語での実装
- 既存のtblsアーキテクチャとの整合性維持
- schema/、config/、output/パッケージの修正が必要
- 各種DBMSドライバーでのコメント取得機能に依存

### 4.2 機能制約

- Markdown出力形式のみが対象
- コメント機能をサポートするDBMSのみが対象
- 区切り文字は1文字に限定

## 5. 成功基準

### 5.1 完了の定義

- [ ] PostgreSQL、MySQL、SQL Serverなどでコメント区切り機能が正常に動作すること
- [ ] tbls.ymlでコメント区切り文字を設定し、論理名とコメントが正しく分離されること
- [ ] tbls.ymlで列順序を設定し、Markdown出力で指定順序で表示されること
- [ ] tbls.ymlで列名エイリアスを設定し、Markdown出力でエイリアスが表示されること
- [ ] 設定が無い場合でも従来通りの動作を維持すること
- [ ] 全ての既存テストがパスすること
- [ ] 新機能に対する包括的なテストケースが追加されていること

### 5.2 品質基準

- [ ] golangci-lintでエラーが無いこと
- [ ] テストカバレッジが既存レベルを維持していること
- [ ] 実際のデータベースでの動作確認が完了していること
- [ ] ドキュメント（README、設定例）が更新されていること

## 6. 設定例

### 6.1 tbls.yml設定例

```yaml
# コメント区切り文字の設定
comment:
  separator: "|"  # 論理名|コメントの形式

# Markdown出力カスタマイズ
markdown:
  # データベースレベル設定
  database:
    show_logical_name: true
    aliases:
      database: "データベース"
      logical_name: "論理名"
      comment: "説明"

  # スキーマレベル設定
  schemas:
    show_logical_name: true
    aliases:
      schema: "スキーマ"
      logical_name: "論理名"
      comment: "説明"

  # テーブルレベル設定
  tables:
    # テーブル設定
    show_logical_name: true
    order: ["name", "logical_name", "comment", "type"]
    aliases:
      name: "テーブル名"
      logical_name: "論理名"
      comment: "説明"
      type: "種別"

  # カラムレベル設定
  columns:
    # カラム設定
    show_logical_name: true
    order: ["name", "logical_name", "type", "nullable", "default", "comment"]
    aliases:
      name: "カラム名"
      logical_name: "論理名"
      type: "データ型"
      nullable: "NULL許可"
      default: "デフォルト値"
      comment: "説明"


  # ビューレベル設定
  views:
    show_logical_name: true
    order: ["name", "logical_name", "comment", "definition"]
    aliases:
      name: "ビュー名"
      logical_name: "論理名"
      comment: "説明"
      definition: "定義"

  # インデックスレベル設定
  indexes:
    show_logical_name: true
    order: ["name", "logical_name", "table", "columns", "unique", "comment"]
    aliases:
      name: "インデックス名"
      logical_name: "論理名"
      table: "対象テーブル"
      columns: "対象カラム"
      unique: "ユニーク"
      comment: "説明"

  # 制約レベル設定
  constraints:
    show_logical_name: true
    order: ["name", "logical_name", "type", "table", "columns", "comment"]
    aliases:
      name: "制約名"
      logical_name: "論理名"
      type: "制約種別"
      table: "対象テーブル"
      columns: "対象カラム"
      comment: "説明"

  # 関数レベル設定
  functions:
    show_logical_name: true
    order: ["name", "logical_name", "schema", "return_type", "comment"]
    aliases:
      name: "関数名"
      logical_name: "論理名"
      schema: "スキーマ"
      return_type: "戻り値型"
      comment: "説明"


```

### 6.2 データベースコメント例

```sql
-- PostgreSQL例

-- データベース
COMMENT ON DATABASE myapp IS 'アプリケーションDB|メインアプリケーションのデータベース';

-- スキーマ
COMMENT ON SCHEMA public IS '公開スキーマ|デフォルトの公開スキーマ';
COMMENT ON SCHEMA auth IS '認証スキーマ|ユーザー認証関連のスキーマ';

-- テーブル
COMMENT ON TABLE users IS 'ユーザーテーブル|システムユーザーの基本情報を管理';
COMMENT ON TABLE posts IS '投稿テーブル|ユーザーによる投稿データを管理';

-- カラム
COMMENT ON COLUMN users.user_id IS 'ユーザーID|システム内でユーザーを一意に識別するID';
COMMENT ON COLUMN users.user_name IS 'ユーザー名|ログイン時に使用する名前';
COMMENT ON COLUMN users.email IS 'メールアドレス|ユーザーの連絡先メールアドレス';

-- ビュー
COMMENT ON VIEW active_users IS 'アクティブユーザービュー|現在アクティブなユーザーのみを表示';

-- インデックス
COMMENT ON INDEX idx_users_email IS 'メールインデックス|メールアドレスでの高速検索用';
COMMENT ON INDEX idx_posts_created_at IS '投稿日時インデックス|投稿日時での並び替え用';

-- 制約
-- PostgreSQLは制約に直接コメントできないため、テーブルコメントまたは
-- カラムコメントで制約の説明を含める

-- 関数
COMMENT ON FUNCTION calculate_age(birth_date date) IS '年齢計算関数|生年月日から現在の年齢を計算';
COMMENT ON FUNCTION get_user_posts(user_id integer) IS 'ユーザー投稿取得|指定ユーザーの全投稿を取得';

-- トリガー（PostgreSQLではトリガー関数にコメント）
COMMENT ON FUNCTION update_modified_time() IS '更新時間トリガー|レコード更新時に自動で更新時間を設定';

-- シーケンス
COMMENT ON SEQUENCE users_id_seq IS 'ユーザーIDシーケンス|ユーザーテーブルの主キー生成用';
```

```sql
-- MySQL例

-- データベース
-- MySQL 8.0以降
CREATE DATABASE myapp COMMENT = 'アプリケーションDB|メインアプリケーションのデータベース';

-- テーブル
CREATE TABLE users (
  user_id INT PRIMARY KEY AUTO_INCREMENT COMMENT 'ユーザーID|システム内でユーザーを一意に識別するID',
  user_name VARCHAR(100) NOT NULL COMMENT 'ユーザー名|ログイン時に使用する名前',
  email VARCHAR(255) NOT NULL COMMENT 'メールアドレス|ユーザーの連絡先メールアドレス'
) COMMENT = 'ユーザーテーブル|システムユーザーの基本情報を管理';

-- ビュー
CREATE VIEW active_users AS
SELECT * FROM users WHERE status = 'active'
COMMENT = 'アクティブユーザービュー|現在アクティブなユーザーのみを表示';
```

```sql
-- SQL Server例

-- データベース
-- SQL Serverでは拡張プロパティを使用
EXEC sp_addextendedproperty
  @name = N'MS_Description',
  @value = N'アプリケーションDB|メインアプリケーションのデータベース',
  @level0type = N'DATABASE',
  @level0name = N'MyApp';

-- テーブル
EXEC sp_addextendedproperty
  @name = N'MS_Description',
  @value = N'ユーザーテーブル|システムユーザーの基本情報を管理',
  @level0type = N'SCHEMA', @level0name = N'dbo',
  @level1type = N'TABLE', @level1name = N'users';

-- カラム
EXEC sp_addextendedproperty
  @name = N'MS_Description',
  @value = N'ユーザーID|システム内でユーザーを一意に識別するID',
  @level0type = N'SCHEMA', @level0name = N'dbo',
  @level1type = N'TABLE', @level1name = N'users',
  @level2type = N'COLUMN', @level2name = N'user_id';
```
