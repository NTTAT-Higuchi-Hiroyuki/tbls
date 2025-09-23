# Configuration Examples - Comment Parsing and Markdown Customization

This document provides comprehensive configuration examples for the new comment parsing and Markdown customization features in tbls v1.86.0.

## Table of Contents

- [Basic Examples](#basic-examples)
- [Database-Specific Examples](#database-specific-examples)
- [Use Case Examples](#use-case-examples)
- [Advanced Examples](#advanced-examples)
- [Multilingual Examples](#multilingual-examples)

## Basic Examples

### Minimal Comment Parsing

Enable basic comment parsing with a simple separator:

```yaml
# .tbls.yml
name: "My Database"
dsn: "postgres://user:pass@localhost:5432/mydb"

comment:
  separator: "|"
```

**Database comments should follow this format:**
```
LogicalName|Description
```

**Example comments:**
- `"User ID|Unique identifier for users"`
- `"Order Date|Date when the order was placed"`
- `"Product Name|Display name of the product"`

### Basic Markdown Customization

Show logical names in documentation:

```yaml
# .tbls.yml
name: "My Database"
dsn: "postgres://user:pass@localhost:5432/mydb"

comment:
  separator: "|"

markdown:
  tables:
    show_logical_name: true
  columns:
    show_logical_name: true
```

### Column Reordering

Change the order of columns in Markdown tables:

```yaml
# .tbls.yml
markdown:
  columns:
    order: ["logical_name", "name", "type", "nullable", "comment"]
```

**Default order:**
| Name | Type | Nullable | Default | Comment |

**Custom order:**
| Logical Name | Name | Type | Nullable | Comment |

## Database-Specific Examples

### PostgreSQL Configuration

PostgreSQL with schema-specific settings:

```yaml
# .tbls.yml - PostgreSQL with advanced features
name: "PostgreSQL Enterprise Database"
dsn: "postgres://username:password@localhost:5432/enterprise_db?sslmode=require"

comment:
  separator: "::"  # PostgreSQL style separator

markdown:
  database:
    show_logical_name: true
    order: ["name", "logical_name", "comment"]
    aliases:
      name: "Database Name"
      logical_name: "Business Name"
      comment: "Description"

  schemas:
    show_logical_name: true
    order: ["name", "logical_name", "comment"]
    aliases:
      name: "Schema Name"
      logical_name: "Business Area"
      comment: "Purpose"

  tables:
    show_logical_name: true
    order: ["name", "logical_name", "comment", "type"]
    aliases:
      name: "Physical Name"
      logical_name: "Business Entity"
      comment: "Description"
      type: "Object Type"

  columns:
    show_logical_name: true
    order: ["logical_name", "name", "type", "nullable", "default", "comment"]
    aliases:
      name: "Physical Name"
      logical_name: "Business Field"
      type: "Data Type"
      nullable: "Required"
      default: "Default Value"
      comment: "Description"

  indexes:
    show_logical_name: true
    aliases:
      name: "Index Name"
      logical_name: "Purpose"
      columns: "Indexed Columns"
      comment: "Description"

  constraints:
    show_logical_name: true
    aliases:
      name: "Constraint Name"
      logical_name: "Business Rule"
      type: "Rule Type"
      columns: "Affected Columns"

# PostgreSQL-specific exclusions
exclude:
  - "pg_*"
  - "information_schema.*"
  - "*_pkey"
  - "*_fkey"
```

### MySQL Configuration

MySQL with engine-specific information:

```yaml
# .tbls.yml - MySQL with performance focus
name: "MySQL Application Database"
dsn: "mysql://username:password@localhost:3306/app_db?charset=utf8mb4&parseTime=true"

comment:
  separator: "|"

markdown:
  tables:
    show_logical_name: true
    order: ["name", "logical_name", "engine", "comment", "type"]
    aliases:
      name: "Table Name"
      logical_name: "Business Entity"
      engine: "Storage Engine"
      comment: "Description"
      type: "Table Type"

  columns:
    show_logical_name: true
    order: ["name", "logical_name", "type", "nullable", "auto_increment", "default", "comment"]
    aliases:
      name: "Column Name"
      logical_name: "Business Field"
      type: "MySQL Type"
      nullable: "NULL Allowed"
      auto_increment: "Auto Increment"
      default: "Default Value"
      comment: "Description"

# MySQL-specific exclusions
exclude:
  - "mysql.*"
  - "information_schema.*"
  - "performance_schema.*"
  - "sys.*"

# MySQL-specific ER diagram settings
er:
  format: "svg"
  showColumnTypes: true
  comment: true
```

### SQL Server Configuration

SQL Server with extended properties:

```yaml
# .tbls.yml - SQL Server Enterprise
name: "SQL Server Enterprise Database"
dsn: "mssql://username:password@localhost:1433/enterprise_db?database=enterprise_db"

comment:
  separator: "|"

markdown:
  schemas:
    show_logical_name: true
    order: ["name", "logical_name", "comment"]
    aliases:
      name: "Schema Name"
      logical_name: "Business Domain"
      comment: "Purpose"

  tables:
    show_logical_name: true
    order: ["schema_name", "name", "logical_name", "comment", "type"]
    aliases:
      schema_name: "Schema"
      name: "Table Name"
      logical_name: "Business Entity"
      comment: "Description"
      type: "Object Type"

  columns:
    show_logical_name: true
    order: ["name", "logical_name", "type", "nullable", "identity", "default", "comment"]
    aliases:
      name: "Column Name"
      logical_name: "Business Field"
      type: "SQL Server Type"
      nullable: "Nullable"
      identity: "Identity"
      default: "Default Value"
      comment: "Description"

# SQL Server-specific exclusions
exclude:
  - "sys.*"
  - "INFORMATION_SCHEMA.*"
  - "msdb.*"
  - "tempdb.*"
```

## Use Case Examples

### Enterprise Documentation

Complete enterprise-grade configuration:

```yaml
# .tbls.yml - Enterprise Documentation Standard
name: "Enterprise Data Dictionary"
desc: "Comprehensive database documentation for enterprise applications"
dsn: "${DATABASE_URL}"

comment:
  separator: "|"

format:
  adjust: true
  sort: true
  showOnlyFirstParagraph: false
  hideColumnsWithoutValues: false

markdown:
  database:
    show_logical_name: true
    order: ["name", "logical_name", "comment"]
    aliases:
      name: "🗄️ Database Name"
      logical_name: "📋 Business Name"
      comment: "📝 Description"

  schemas:
    show_logical_name: true
    order: ["name", "logical_name", "comment"]
    aliases:
      name: "📂 Schema Name"
      logical_name: "🏢 Business Domain"
      comment: "📋 Purpose"

  tables:
    show_logical_name: true
    order: ["name", "logical_name", "comment", "type"]
    aliases:
      name: "📋 Table Name"
      logical_name: "💼 Business Entity"
      comment: "📝 Description"
      type: "🏷️ Type"
    specific:
      # Critical business tables
      customers:
        aliases:
          name: "👥 Customer Management"
          logical_name: "顧客管理テーブル"
      orders:
        aliases:
          name: "🛒 Order Processing"
          logical_name: "注文処理テーブル"
      products:
        aliases:
          name: "📦 Product Catalog"
          logical_name: "商品カタログテーブル"

  columns:
    show_logical_name: true
    order: ["logical_name", "name", "type", "nullable", "default", "comment"]
    aliases:
      name: "🔧 Physical Name"
      logical_name: "💼 Business Field"
      type: "📊 Data Type"
      nullable: "❓ Required"
      default: "⚙️ Default"
      comment: "📝 Description"
    specific:
      # Customer table columns
      customers:
        order: ["logical_name", "name", "type", "nullable", "comment"]
        aliases:
          logical_name: "顧客項目名"
          name: "物理カラム名"
          type: "データ型"
          nullable: "必須項目"
          comment: "項目説明"

  views:
    show_logical_name: true
    aliases:
      name: "👁️ View Name"
      logical_name: "📊 Business View"
      comment: "📝 Purpose"

  indexes:
    show_logical_name: true
    aliases:
      name: "🗂️ Index Name"
      logical_name: "⚡ Performance Purpose"
      columns: "📋 Indexed Columns"
      comment: "📝 Description"

  constraints:
    show_logical_name: true
    aliases:
      name: "🔒 Constraint Name"
      logical_name: "📏 Business Rule"
      type: "🏷️ Rule Type"
      columns: "📋 Affected Columns"
      comment: "📝 Description"

  functions:
    show_logical_name: true
    aliases:
      name: "⚙️ Function Name"
      logical_name: "💼 Business Function"
      return_type: "📤 Returns"
      arguments: "📥 Parameters"
      comment: "📝 Purpose"

# Enterprise-grade filtering
include:
  - name: "customer*"
    labels: ["core", "customer-data"]
  - name: "order*"
    labels: ["core", "transaction-data"]
  - name: "product*"
    labels: ["core", "catalog-data"]
  - name: "view_*"
    kind: "view"
    labels: ["reporting"]

exclude:
  - "*_temp"
  - "*_backup"
  - "*_archive"
  - "test_*"
  - "tmp_*"

# Comprehensive linting
lint:
  require:
    tableComment:
      enabled: true
    columnComment:
      enabled: true
      exclude: ["id", "created_at", "updated_at", "deleted_at"]
    indexComment:
      enabled: true
    relationName:
      enabled: true
  unrelatedTable:
    enabled: true
    exclude: ["logs", "*_temp", "*_archive"]

# Professional ER diagrams
er:
  format: "svg"
  comment: true
  fontSize: 12
  showColumnTypes:
    related: true
    primary: true
  showColumnComments: false
```

### API Documentation

Configuration for API-focused documentation:

```yaml
# .tbls.yml - API Documentation Focus
name: "API Database Schema"
desc: "Database schema for REST API endpoints"
dsn: "postgres://api_user:${API_DB_PASSWORD}@localhost:5432/api_db"

comment:
  separator: ":"

markdown:
  tables:
    show_logical_name: true
    order: ["name", "logical_name", "comment", "type"]
    aliases:
      name: "🌐 API Resource"
      logical_name: "📋 Resource Name"
      comment: "📝 API Purpose"
      type: "🏷️ Type"
    specific:
      users:
        aliases:
          name: "👥 /api/users"
          logical_name: "User Accounts"
      posts:
        aliases:
          name: "📝 /api/posts"
          logical_name: "Blog Posts"
      comments:
        aliases:
          name: "💬 /api/comments"
          logical_name: "User Comments"

  columns:
    show_logical_name: true
    order: ["name", "logical_name", "type", "nullable", "comment"]
    aliases:
      name: "🔧 JSON Field"
      logical_name: "📋 API Field"
      type: "📊 JSON Type"
      nullable: "❓ Optional"
      comment: "📝 Field Description"
    specific:
      users:
        order: ["logical_name", "name", "type", "nullable", "comment"]
      posts:
        order: ["logical_name", "name", "type", "nullable", "comment"]
      comments:
        order: ["logical_name", "name", "type", "nullable", "comment"]

# API-focused filtering
include:
  - "users"
  - "posts"
  - "comments"
  - "categories"
  - "tags"
  - "user_*"
  - "post_*"

exclude:
  - "*_logs"
  - "*_audit"
  - "*_sessions"
  - "admin_*"
```

### Analytics Database

Configuration for analytics and reporting:

```yaml
# .tbls.yml - Analytics Database
name: "Data Warehouse Schema"
desc: "Analytics and reporting database documentation"
dsn: "bigquery://my-project/analytics_dataset?credentialsFile=/path/to/key.json"

comment:
  separator: " | "

markdown:
  tables:
    show_logical_name: true
    order: ["name", "logical_name", "comment", "type"]
    aliases:
      name: "📊 Table Name"
      logical_name: "📈 Business Metric"
      comment: "📝 Analytics Purpose"
      type: "🏷️ Table Type"
    specific:
      fact_sales:
        aliases:
          name: "💰 Sales Facts"
          logical_name: "売上実績ファクト"
      dim_customers:
        aliases:
          name: "👥 Customer Dimension"
          logical_name: "顧客ディメンション"
      dim_products:
        aliases:
          name: "📦 Product Dimension"
          logical_name: "商品ディメンション"
      dim_time:
        aliases:
          name: "📅 Time Dimension"
          logical_name: "時間ディメンション"

  columns:
    show_logical_name: true
    order: ["logical_name", "name", "type", "nullable", "comment"]
    aliases:
      name: "🔧 Column Name"
      logical_name: "📊 Metric Name"
      type: "📈 Data Type"
      nullable: "❓ Optional"
      comment: "📝 Calculation Logic"

# Analytics-focused filtering
include:
  - "fact_*"
  - "dim_*"
  - "staging_*"
  - "mart_*"

exclude:
  - "*_temp"
  - "*_backup"
  - "test_*"

# Analytics-specific ER diagram
er:
  format: "mermaid"
  showColumnTypes: true
  comment: true
```

## Advanced Examples

### Multi-Schema Enterprise

Configuration for complex multi-schema environments:

```yaml
# .tbls.yml - Multi-Schema Enterprise
name: "Enterprise Multi-Schema Database"
dsn: "postgres://user:pass@localhost:5432/enterprise"

comment:
  separator: "::"

markdown:
  schemas:
    show_logical_name: true
    order: ["name", "logical_name", "comment"]
    aliases:
      name: "Schema Name"
      logical_name: "Business Domain"
      comment: "Domain Purpose"

  tables:
    show_logical_name: true
    order: ["schema_name", "name", "logical_name", "comment", "type"]
    aliases:
      schema_name: "📂 Domain"
      name: "📋 Table"
      logical_name: "💼 Entity"
      comment: "📝 Description"
      type: "🏷️ Type"
    specific:
      # Sales domain
      sales.customers:
        aliases:
          name: "👥 Customer Master"
          logical_name: "顧客マスタ"
      sales.orders:
        aliases:
          name: "🛒 Order Transactions"
          logical_name: "注文取引"
      # HR domain
      hr.employees:
        aliases:
          name: "👤 Employee Records"
          logical_name: "従業員記録"
      hr.departments:
        aliases:
          name: "🏢 Department Structure"
          logical_name: "部門構造"
      # Finance domain
      finance.accounts:
        aliases:
          name: "💰 Chart of Accounts"
          logical_name: "勘定科目"
      finance.transactions:
        aliases:
          name: "📊 Financial Transactions"
          logical_name: "財務取引"

  columns:
    show_logical_name: true
    order: ["logical_name", "name", "type", "nullable", "default", "comment"]
    aliases:
      name: "🔧 Column"
      logical_name: "💼 Field"
      type: "📊 Type"
      nullable: "❓ Required"
      default: "⚙️ Default"
      comment: "📝 Description"
    specific:
      # Schema-specific column settings
      sales.customers:
        order: ["logical_name", "name", "type", "nullable", "comment"]
        aliases:
          logical_name: "顧客項目"
          name: "物理項目"
      hr.employees:
        order: ["logical_name", "name", "type", "nullable", "comment"]
        aliases:
          logical_name: "人事項目"
          name: "物理項目"

# Schema-based filtering
include:
  - name: "*"
    schema: "sales"
    labels: ["sales", "customer-facing"]
  - name: "*"
    schema: "hr"
    labels: ["hr", "internal"]
  - name: "*"
    schema: "finance"
    labels: ["finance", "compliance"]

exclude:
  - name: "*"
    schema: "temp"
  - name: "*"
    schema: "staging"
  - name: "*_archive"
```

### Internationalization Example

Full internationalization setup:

```yaml
# .tbls.yml - International Documentation
name: "国際化データベース設計書"
desc: "多言語対応システムのデータベース設計書"
dsn: "postgres://user:pass@localhost:5432/international_db"

comment:
  separator: "|"

markdown:
  database:
    show_logical_name: true
    order: ["name", "logical_name", "comment"]
    aliases:
      name: "データベース名"
      logical_name: "論理名"
      comment: "説明"

  tables:
    show_logical_name: true
    order: ["name", "logical_name", "comment", "type"]
    aliases:
      name: "テーブル名"
      logical_name: "論理名"
      comment: "説明"
      type: "種別"
    specific:
      users:
        aliases:
          name: "👥 ユーザーテーブル"
          logical_name: "利用者マスタ"
      products:
        aliases:
          name: "📦 商品テーブル"
          logical_name: "商品マスタ"
      orders:
        aliases:
          name: "🛒 注文テーブル"
          logical_name: "注文データ"
      categories:
        aliases:
          name: "📂 カテゴリテーブル"
          logical_name: "分類マスタ"

  columns:
    show_logical_name: true
    order: ["logical_name", "name", "type", "nullable", "default", "comment"]
    aliases:
      name: "カラム名"
      logical_name: "論理名"
      type: "データ型"
      nullable: "NULL許可"
      default: "デフォルト値"
      comment: "説明"
    specific:
      users:
        order: ["logical_name", "name", "type", "nullable", "comment"]
        aliases:
          logical_name: "項目名"
          name: "物理名"
          type: "型"
          nullable: "必須"
          comment: "説明"
      products:
        aliases:
          logical_name: "商品項目名"
          name: "物理カラム名"

  views:
    show_logical_name: true
    aliases:
      name: "ビュー名"
      logical_name: "論理名"
      comment: "説明"
      definition: "定義"

  indexes:
    show_logical_name: true
    aliases:
      name: "インデックス名"
      logical_name: "論理名"
      columns: "対象カラム"
      comment: "説明"

  constraints:
    show_logical_name: true
    aliases:
      name: "制約名"
      logical_name: "論理名"
      type: "制約タイプ"
      columns: "対象カラム"
      comment: "説明"

  functions:
    show_logical_name: true
    aliases:
      name: "関数名"
      logical_name: "論理名"
      return_type: "戻り値型"
      arguments: "引数"
      comment: "説明"
```

## Multilingual Examples

### Japanese-English Mixed

Perfect for international teams:

```yaml
# .tbls.yml - Japanese-English Documentation
name: "国際プロジェクト Database Schema"
desc: "International project database with Japanese business terms"
dsn: "postgres://user:pass@localhost:5432/international_app"

comment:
  separator: "|"

markdown:
  tables:
    show_logical_name: true
    order: ["name", "logical_name", "comment", "type"]
    aliases:
      name: "Table Name"
      logical_name: "日本語名"
      comment: "Description"
      type: "Type"

  columns:
    show_logical_name: true
    order: ["name", "logical_name", "type", "nullable", "comment"]
    aliases:
      name: "Column Name"
      logical_name: "項目名"
      type: "Data Type"
      nullable: "Required"
      comment: "Description"
```

**Example database comments:**
```sql
COMMENT ON TABLE users IS 'ユーザー|User account information';
COMMENT ON COLUMN users.id IS 'ユーザーID|Unique user identifier';
COMMENT ON COLUMN users.name IS 'ユーザー名|Display name of the user';
COMMENT ON COLUMN users.email IS 'メールアドレス|Email address for login';
```

### Spanish Documentation

```yaml
# .tbls.yml - Documentación en Español
name: "Esquema de Base de Datos"
desc: "Documentación completa de la base de datos en español"
dsn: "mysql://usuario:contraseña@localhost:3306/aplicacion"

comment:
  separator: "|"

markdown:
  tables:
    show_logical_name: true
    aliases:
      name: "Nombre de Tabla"
      logical_name: "Nombre Lógico"
      comment: "Descripción"
      type: "Tipo"

  columns:
    show_logical_name: true
    aliases:
      name: "Nombre de Columna"
      logical_name: "Nombre Lógico"
      type: "Tipo de Dato"
      nullable: "Permite NULL"
      default: "Valor Predeterminado"
      comment: "Descripción"
```

### French Documentation

```yaml
# .tbls.yml - Documentation Française
name: "Schéma de Base de Données"
desc: "Documentation complète de la base de données en français"
dsn: "postgres://utilisateur:motdepasse@localhost:5432/application"

comment:
  separator: "|"

markdown:
  tables:
    show_logical_name: true
    aliases:
      name: "Nom de Table"
      logical_name: "Nom Logique"
      comment: "Description"
      type: "Type"

  columns:
    show_logical_name: true
    aliases:
      name: "Nom de Colonne"
      logical_name: "Nom Logique"
      type: "Type de Données"
      nullable: "Autorise NULL"
      default: "Valeur par Défaut"
      comment: "Description"
```

## Best Practices

### Configuration Organization

1. **Separate by Environment**
   ```yaml
   # .tbls.production.yml
   # .tbls.staging.yml
   # .tbls.development.yml
   ```

2. **Use Environment Variables**
   ```yaml
   dsn: "${DATABASE_URL}"
   name: "${PROJECT_NAME} Database"
   ```

3. **Hierarchical Configuration**
   ```yaml
   # Global settings
   markdown:
     columns:
       show_logical_name: true
       # Specific overrides
       specific:
         important_table:
           order: ["logical_name", "name"]
   ```

### Performance Optimization

```yaml
# Optimize for large databases
include:
  - "core_*"     # Only include core tables
exclude:
  - "*_temp"     # Exclude temporary tables
  - "*_log"      # Exclude log tables
  - "*_archive"  # Exclude archive tables

er:
  skip: false
  showOnlyRelatedTables: true  # Reduce ER diagram complexity
```

### Team Collaboration

```yaml
# Team-friendly configuration
comment:
  separator: "|"  # Standardize across team

markdown:
  # Use consistent aliases across projects
  columns:
    aliases:
      name: "Column Name"
      logical_name: "Business Name"
      type: "Data Type"
      nullable: "Required"
      comment: "Description"

# Consistent filtering
exclude:
  - "*_temp"
  - "*_backup"
  - "test_*"
```