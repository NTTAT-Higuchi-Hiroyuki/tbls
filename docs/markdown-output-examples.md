# Markdown出力例 - 論理名とカスタマイズ機能

このドキュメントでは、tbls v1.86.0の新機能である論理名表示とMarkdown出力カスタマイズ機能の出力例を示します。

## ⚠️ 重要: `tables`と`columns`の設定の違い

Markdown出力のカスタマイズには、**2つの異なる設定領域**があります：

### `markdown.tables` - テーブルリストの表示制御

**適用先:** `README.md`内の**テーブル一覧表**

```yaml
markdown:
  tables:
    show_logical_name: true
    order: ["name", "LogicalName", "comment", "type"]
```

**出力例（README.md内）:**
```markdown
## Tables

| Name     | Logical Name | Comment              | Type  |
| -------- | ------------ | -------------------- | ----- |
| users    | ユーザー     | システム利用者の管理 | table |
| products | 商品         | 販売商品の管理       | table |
```

### `markdown.columns` - カラム詳細の表示制御

**適用先:** 各テーブルページ（`users.md`等）の**Columnsセクション**

```yaml
markdown:
  columns:
    show_logical_name: true
    order: ["LogicalName", "name", "type", "nullable", "comment"]
```

**出力例（users.md内）:**
```markdown
## Columns

| Logical Name | Name     | Type    | Nullable | Comment          |
| ------------ | -------- | ------- | -------- | ---------------- |
| ユーザーID   | user_id  | integer | false    | Primary key      |
| ユーザー名   | username | varchar | false    | Login user name  |
```

### 両方を設定する場合

通常、両方の設定を行うことで、一貫性のあるドキュメントを生成できます：

```yaml
markdown:
  # テーブルリスト（README.md）の設定
  tables:
    show_logical_name: true
    order: ["name", "LogicalName", "comment", "type"]

  # カラム詳細（各テーブルページ）の設定
  columns:
    show_logical_name: true
    order: ["LogicalName", "name", "type", "nullable", "comment"]
```

---

## 目次

- [⚠️ 重要: tablesとcolumnsの設定の違い](#️-重要-tablesとcolumnsの設定の違い)
- [従来の出力 vs 新機能の出力](#従来の出力-vs-新機能の出力)
- [日本語環境での出力例](#日本語環境での出力例)
- [エンタープライズ環境での出力例](#エンタープライズ環境での出力例)
- [列順序カスタマイズ例](#列順序カスタマイズ例)

## 従来の出力 vs 新機能の出力

### 従来の出力（論理名なし）

```markdown
# public.users

## Columns

| Name     | Type    | Default | Nullable | Comment                               |
| -------- | ------- | ------- | -------- | ------------------------------------- |
| user_id  | integer |         | false    | システム内でユーザーを一意に識別するID |
| username | varchar |         | false    | ログイン時に使用する名前              |
| email    | varchar |         | false    | ユーザーの連絡先メールアドレス        |
```

### 新機能の出力（論理名あり）

**設定ファイル (.tbls.yml):**
```yaml
comment:
  separator: "|"

markdown:
  tables:
    show_logical_name: true
  columns:
    show_logical_name: true
    order: ["LogicalName", "name", "type", "nullable", "comment"]
    aliases:
      name: "物理名"
      LogicalName: "論理名"
      type: "データ型"
      nullable: "必須"
      comment: "説明"
```

**データベースコメント:**
```sql
COMMENT ON TABLE users IS 'ユーザーテーブル|システムユーザーの基本情報を管理';
COMMENT ON COLUMN users.user_id IS 'ユーザーID|システム内でユーザーを一意に識別するID';
COMMENT ON COLUMN users.username IS 'ユーザー名|ログイン時に使用する名前';
COMMENT ON COLUMN users.email IS 'メールアドレス|ユーザーの連絡先メールアドレス';
```

**出力結果:**
```markdown
# public.users

## Description

**論理名:** ユーザーテーブル

システムユーザーの基本情報を管理

## Columns

| 論理名         | 物理名   | データ型 | 必須  | 説明                           |
| -------------- | -------- | -------- | ----- | ------------------------------ |
| ユーザーID     | user_id  | integer  | false | システム内でユーザーを一意に識別するID |
| ユーザー名     | username | varchar  | false | ログイン時に使用する名前       |
| メールアドレス | email    | varchar  | false | ユーザーの連絡先メールアドレス |
```

## 日本語環境での出力例

### 設定ファイル

```yaml
# .tbls.yml - 日本語環境
name: "サンプルアプリケーション データベース"
dsn: "postgres://user:pass@localhost:5432/sample_app"

comment:
  separator: "|"

markdown:
  database:
    show_logical_name: true
    aliases:
      name: "データベース名"
      LogicalName: "論理名"
      comment: "説明"

  tables:
    show_logical_name: true
    order: ["name", "LogicalName", "comment", "type"]
    aliases:
      name: "テーブル名"
      LogicalName: "論理名"
      comment: "説明"
      type: "種別"

  columns:
    show_logical_name: true
    order: ["LogicalName", "name", "type", "nullable", "default", "comment"]
    aliases:
      name: "カラム名"
      LogicalName: "論理名"
      type: "データ型"
      nullable: "NULL許可"
      default: "デフォルト値"
      comment: "説明"
```

### 出力結果

```markdown
# サンプルアプリケーション データベース

## Description

**論理名:** ECサイトDB

Eコマースアプリケーションのメインデータベース

## Tables

| テーブル名 | 論理名       | 説明                     | 種別  |
| ---------- | ------------ | ------------------------ | ----- |
| users      | ユーザー     | システム利用者の管理     | table |
| products   | 商品         | 販売商品の管理           | table |
| orders     | 注文         | 注文情報の管理           | table |
| categories | カテゴリ     | 商品分類の管理           | table |

---

# public.users

## Description

**論理名:** ユーザー

システム利用者の管理

## Columns

| 論理名           | カラム名    | データ型                | NULL許可 | デフォルト値                 | 説明                           |
| ---------------- | ----------- | ----------------------- | -------- | ---------------------------- | ------------------------------ |
| ユーザーID       | user_id     | integer                 | false    |                              | システム内でユーザーを一意に識別するID |
| ユーザー名       | username    | character varying(100)  | false    |                              | ログイン時に使用する名前       |
| メールアドレス   | email       | character varying(255)  | false    |                              | ユーザーの連絡先メールアドレス |
| パスワードハッシュ | password_hash | character varying(255) | false    |                              | 暗号化されたパスワード         |
| 作成日時         | created_at  | timestamp               | false    | CURRENT_TIMESTAMP            | レコード作成日時               |
| 更新日時         | updated_at  | timestamp               | true     | CURRENT_TIMESTAMP            | レコード最終更新日時           |
| 削除フラグ       | deleted_at  | timestamp               | true     |                              | 論理削除用のタイムスタンプ     |

## Constraints

| 制約名              | 種別        | 定義                                    |
| ------------------- | ----------- | --------------------------------------- |
| users_pkey          | PRIMARY KEY | PRIMARY KEY (user_id)                   |
| users_username_key  | UNIQUE      | UNIQUE (username)                       |
| users_email_key     | UNIQUE      | UNIQUE (email)                          |

## Indexes

| インデックス名           | 定義                                                    |
| ----------------------- | ------------------------------------------------------- |
| users_pkey              | CREATE UNIQUE INDEX users_pkey ON users USING btree (user_id) |
| users_username_key      | CREATE UNIQUE INDEX users_username_key ON users USING btree (username) |
| users_email_key         | CREATE UNIQUE INDEX users_email_key ON users USING btree (email) |
| idx_users_created_at    | CREATE INDEX idx_users_created_at ON users USING btree (created_at) |

---

# public.products

## Description

**論理名:** 商品

販売商品の管理

## Columns

| 論理名       | カラム名      | データ型               | NULL許可 | デフォルト値      | 説明                     |
| ------------ | ------------- | ---------------------- | -------- | ----------------- | ------------------------ |
| 商品ID       | product_id    | integer                | false    |                   | 商品を一意に識別するID   |
| 商品名       | product_name  | character varying(200) | false    |                   | 商品の表示名             |
| 商品説明     | description   | text                   | true     |                   | 商品の詳細説明           |
| 価格         | price         | decimal(10,2)          | false    |                   | 商品の販売価格（税込み） |
| 在庫数       | stock_count   | integer                | false    | 0                 | 現在の在庫数             |
| カテゴリID   | category_id   | integer                | false    |                   | 商品カテゴリの外部キー   |
| 作成日時     | created_at    | timestamp              | false    | CURRENT_TIMESTAMP | レコード作成日時         |
| 更新日時     | updated_at    | timestamp              | true     | CURRENT_TIMESTAMP | レコード最終更新日時     |
```

## エンタープライズ環境での出力例

### 設定ファイル

```yaml
# .tbls.yml - エンタープライズ環境
name: "企業基幹システム データベース"

comment:
  separator: "::"

markdown:
  tables:
    show_logical_name: true
    order: ["schema_name", "name", "LogicalName", "comment", "type"]
    aliases:
      schema_name: "📂 スキーマ"
      name: "📋 テーブル名"
      LogicalName: "💼 業務名"
      comment: "📝 説明"
      type: "🏷️ 種別"

  columns:
    show_logical_name: true
    order: ["LogicalName", "name", "type", "nullable", "default", "comment"]
    aliases:
      name: "🔧 物理名"
      LogicalName: "💼 業務項目名"
      type: "📊 データ型"
      nullable: "❓ 必須"
      default: "⚙️ デフォルト"
      comment: "📝 詳細説明"
```

### 出力結果

```markdown
# 企業基幹システム データベース

## Tables

| 📂 スキーマ | 📋 テーブル名 | 💼 業務名       | 📝 説明                   | 🏷️ 種別 |
| ----------- | ------------- | --------------- | ------------------------- | -------- |
| sales       | customers     | 顧客マスタ      | 取引先顧客の基本情報      | table    |
| sales       | orders        | 受注データ      | 顧客からの注文情報        | table    |
| hr          | employees     | 従業員マスタ    | 社員の基本情報と所属      | table    |
| finance     | accounts      | 勘定科目マスタ  | 会計システムの勘定科目    | table    |

---

# sales.customers

## Description

**💼 業務名:** 顧客マスタ

取引先顧客の基本情報

## Columns

| 💼 業務項目名     | 🔧 物理名        | 📊 データ型            | ❓ 必須 | ⚙️ デフォルト     | 📝 詳細説明                      |
| ----------------- | --------------- | ---------------------- | ------- | ----------------- | -------------------------------- |
| 顧客コード       | customer_code   | character varying(10)  | false   |                   | 顧客を一意に識別する業務コード   |
| 顧客名           | customer_name   | character varying(100) | false   |                   | 正式な顧客名（法人名・個人名）   |
| 顧客名カナ       | customer_kana   | character varying(100) | true    |                   | 顧客名のカタカナ表記             |
| 郵便番号         | postal_code     | character varying(8)   | true    |                   | 住所の郵便番号（ハイフンあり）   |
| 住所1            | address1        | character varying(100) | true    |                   | 都道府県・市区町村               |
| 住所2            | address2        | character varying(100) | true    |                   | 町名・番地・建物名               |
| 電話番号         | phone_number    | character varying(15)  | true    |                   | 代表電話番号                     |
| FAX番号          | fax_number      | character varying(15)  | true    |                   | FAX番号                          |
| メールアドレス   | email_address   | character varying(100) | true    |                   | 連絡用メールアドレス             |
| 担当者名         | contact_person  | character varying(50)  | true    |                   | 主要連絡担当者名                 |
| 顧客区分         | customer_type   | integer                | false   | 1                 | 1:法人, 2:個人, 3:その他         |
| 与信限度額       | credit_limit    | decimal(12,0)          | true    |                   | 取引可能な上限金額（円）         |
| 取引開始日       | contract_date   | date                   | false   |                   | 初回取引開始日                   |
| 取引停止フラグ   | is_suspended    | boolean                | false   | false             | true:取引停止, false:取引可能    |
| 備考             | remarks         | text                   | true    |                   | 顧客に関する特記事項             |
| 登録日時         | created_at      | timestamp              | false   | CURRENT_TIMESTAMP | レコード作成日時                 |
| 更新日時         | updated_at      | timestamp              | true    | CURRENT_TIMESTAMP | レコード最終更新日時             |
| 更新者ID         | updated_by      | integer                | true    |                   | 最終更新者の従業員ID             |
```

## 列順序カスタマイズ例

`order`パラメータを使用すると、出力する列の順序と表示する列を制御できます。

**重要:** `order`配列に指定した列のみが出力されます。指定されていない列は非表示になります。これにより、必要な情報だけを選択的に表示できます。

### 開発者向け設定（技術情報優先）

```yaml
markdown:
  columns:
    order: ["name", "type", "nullable", "default", "LogicalName", "comment"]
    aliases:
      name: "Column"
      type: "Type"
      nullable: "Null"
      default: "Default"
      LogicalName: "Business Name"
      comment: "Description"
```

**出力:**
```markdown
| Column       | Type          | Null  | Default           | Business Name    | Description                  |
| ------------ | ------------- | ----- | ----------------- | ---------------- | ---------------------------- |
| user_id      | integer       | false |                   | ユーザーID       | システム内でユーザーを一意に識別するID |
| username     | varchar(100)  | false |                   | ユーザー名       | ログイン時に使用する名前     |
| email        | varchar(255)  | false |                   | メールアドレス   | ユーザーの連絡先メールアドレス |
| created_at   | timestamp     | false | CURRENT_TIMESTAMP | 作成日時         | レコード作成日時             |
```

### ビジネス向け設定（業務情報優先）

```yaml
markdown:
  columns:
    order: ["LogicalName", "comment", "name", "type", "nullable"]
    aliases:
      LogicalName: "項目名"
      comment: "説明"
      name: "システム項目名"
      type: "データ形式"
      nullable: "必須項目"
```

**出力:**
```markdown
| 項目名         | 説明                           | システム項目名 | データ形式   | 必須項目 |
| -------------- | ------------------------------ | -------------- | ------------ | -------- |
| ユーザーID     | システム内でユーザーを一意に識別するID | user_id        | integer      | 必須     |
| ユーザー名     | ログイン時に使用する名前       | username       | varchar(100) | 必須     |
| メールアドレス | ユーザーの連絡先メールアドレス | email          | varchar(255) | 必須     |
| 作成日時       | レコード作成日時               | created_at     | timestamp    | 必須     |
```

## API仕様書向け設定

```yaml
markdown:
  columns:
    order: ["LogicalName", "name", "type", "nullable", "comment"]
    aliases:
      LogicalName: "🌐 API Field"
      name: "🔧 JSON Key"
      type: "📊 JSON Type"
      nullable: "❓ Required"
      comment: "📝 Description"
```

**出力:**
```markdown
| 🌐 API Field   | 🔧 JSON Key | 📊 JSON Type | ❓ Required | 📝 Description                  |
| -------------- | ----------- | ------------ | ----------- | -------------------------------- |
| ユーザーID     | userId      | number       | Required    | システム内でユーザーを一意に識別するID |
| ユーザー名     | username    | string       | Required    | ログイン時に使用する名前         |
| メールアドレス | email       | string       | Required    | ユーザーの連絡先メールアドレス   |
| 作成日時       | createdAt   | string       | Required    | レコード作成日時（ISO 8601形式） |
```

## Constraints と Indexes のカスタマイズ

### ⚠️ 重要な注意点

**`show_logical_name: true` だけではLogicalNameは表示されません。**

`order` パラメータで明示的に `LogicalName` を指定する必要があります。

### Constraints のカスタマイズ例

```yaml
markdown:
  constraints:
    show_logical_name: true
    order: ["Name", "LogicalName", "Type", "Definition", "Comment"]
    aliases:
      Name: "制約名"
      LogicalName: "論理名"
      Type: "種類"
      Definition: "定義"
      Comment: "説明"
```

**出力:**
```markdown
| 制約名 | 論理名 | 種類 | 定義 | 説明 |
| --- | --- | --- | --- | --- |
| users_pkey | ユーザーID主キー | PRIMARY KEY | PRIMARY KEY (id) | Primary key for user identification |
```

### Indexes のカスタマイズ例

```yaml
markdown:
  indexes:
    show_logical_name: true
    order: ["Name", "LogicalName", "Definition", "Comment"]
    aliases:
      Name: "インデックス名"
      LogicalName: "論理名"
      Definition: "定義"
      Comment: "説明"
```

**出力:**
```markdown
| インデックス名 | 論理名 | 定義 | 説明 |
| --- | --- | --- | --- |
| idx_users_email | メールアドレス索引 | CREATE INDEX idx_users_email ON users(email) | Email lookup optimization |
```

### LogicalNameが空の場合の動作

LogicalNameが設定されていない場合、自動的に `Name` の値がフォールバックとして使用されます：

```markdown
| 制約名 | 論理名 | 種類 | 定義 | 説明 |
| --- | --- | --- | --- | --- |
| users_pkey | users_pkey | PRIMARY KEY | PRIMARY KEY (id) | |
```

### デフォルトの出力（orderを指定しない場合）

`order` を指定しない場合、以下の列がデフォルトで表示されます：

**Constraints:**
- `Name`
- `Type`
- `Definition`
- `Comment` (制約にコメントがある場合のみ)

**Indexes:**
- `Name`
- `Definition`
- `Comment` (インデックスにコメントがある場合のみ)

**LogicalName は含まれません。**

## まとめ

tbls v1.86.0の新機能により、以下のカスタマイズが可能になりました：

- **論理名表示**: 業務観点での項目名を明確に表示
- **列順序変更**: 用途に応じて重要な情報を優先表示
- **多言語エイリアス**: 日本語やその他の言語での項目名表示
- **用途別設定**: 開発者向け、ビジネス向け、API仕様書向けなど
- **Constraints/Indexesカスタマイズ**: 制約とインデックスの論理名表示（明示的な`order`指定が必要）

これらの機能により、データベース設計書がより読みやすく、実用的なドキュメントとして活用できます。