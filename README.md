<p align="center">
<br>
<img src="https://github.com/k1LoW/tbls/raw/main/img/logo.png" width="200" alt="tbls">
<br><br>
</p>

[![Build Status](https://github.com/k1LoW/tbls/workflows/build/badge.svg)](https://github.com/k1LoW/tbls/actions) [![GitHub release](https://img.shields.io/github/release/k1LoW/tbls.svg)](https://github.com/k1LoW/tbls/releases) [![Go Report Card](https://goreportcard.com/badge/github.com/k1LoW/tbls)](https://goreportcard.com/report/github.com/k1LoW/tbls) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/tbls/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/tbls/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/tbls/time.svg)

`tbls` (pronounced /ˈteɪbl̩z/) is a CI-Friendly tool to document a database, written in Go.

Key features of `tbls` are:

- **Document a database automatically in [GFM](https://github.github.com/gfm/) format. Output database schema [in many formats](#output-formats).**
- **Single binary = CI-Friendly.**
- **[Support many databases](#support-datasource).**
- **Work as linter for database**
- **Advanced comment parsing and logical name extraction**
- **Customizable Markdown output with flexible column ordering and aliases**

### Table of Contents

  - [Quick Start](#quick-start)
  - [Install](#install)
  - [Getting Started](#getting-started)
    - [Document a database](#document-a-database)
    - [Diff database and (document or database)](#diff-database-and-document-or-database)
    - [Lint a database](#lint-a-database)
    - [Measure document coverage](#measure-document-coverage)
    - [Continuous Integration](#continuous-integration)
  - [Configuration](#configuration)
    - [Name](#name)
    - [Description](#description)
    - [Labels](#labels)
    - [DSN](#dsn)
      - [Support Datasource](#support-datasource)
    - [Document path](#document-path)
    - [Document format](#document-format)
    - [ER diagram](#er-diagram)
    - [Filter tables](#filter-tables)
    - [Lint](#lint)
    - [Comments](#comments)
    - [Comment Parsing and Logical Names](#comment-parsing-and-logical-names)
    - [Markdown Output Customization](#markdown-output-customization)
    - [Relations](#relations)
    - [Viewpoints](#viewpoints)
    - [Dictionary](#dictionary)
    - [Personalized Templates](#personalized-templates)
    - [Required Version](#required-version)
  - [Expand environment variables](#expand-environment-variables)
  - [Output formats](#output-formats)
  - [Command arguments](#command-arguments)
  - [Environment variables](#environment-variables)

<br>

## Quick Start

1. Install tbls to macOS via [Homebrew](https://brew.sh/) or [MacPorts](https://www.macports.org/)

   **Homebrew**

   ```console
   $ brew install k1LoW/tap/tbls
   ```

   **MacPorts**

   ```console
   $ sudo port install tbls
   ```

   <br>

2. Setup database

   ```console
   $ make setup
   ```

   <br>

3. Run tbls doc to analyze a database and generate document in GitHub Flavored Markdown format.

   ```console
   $ tbls doc postgres://dbuser:dbpass@hostname:port/dbname
   ```

   or

   ```console
   $ tbls doc my://dbuser:dbpass@hostname:port/dbname
   ```

   <br>

4. Commit generated document

   ```console
   $ git add .
   $ git commit -m 'Update database document'
   ```

   <br>

5. View the document on GitHub

   ![sample document](https://user-images.githubusercontent.com/157344/62417503-0f5dd900-b6a3-11e9-9b88-25de2a1f6b97.png)

   **[ Show sample document on GitHub](https://github.com/k1LoW/tbls/tree/main/sample/mysql)**

## Install

**deb:**

Use [dpkg](https://en.wikipedia.org/wiki/Dpkg):

``` console
$ export TBLS_VERSION=X.X.X
$ curl -o tbls.deb -L https://github.com/k1LoW/tbls/releases/download/v$TBLS_VERSION/tbls_$TBLS_VERSION-1_amd64.deb
$ dpkg -i tbls.deb
```

**RPM:**

``` console
$ export TBLS_VERSION=X.X.X
$ yum install https://github.com/k1LoW/tbls/releases/download/v$TBLS_VERSION/tbls_$TBLS_VERSION-1_amd64.rpm
```

**apk:**

Use [apk](https://en.wikipedia.org/wiki/Alpine_Linux#APK_software_package_management):

``` console
$ export TBLS_VERSION=X.X.X
$ curl -o tbls.apk -L https://github.com/k1LoW/tbls/releases/download/v$TBLS_VERSION/tbls_$TBLS_VERSION-1_amd64.apk
$ apk add tbls.apk
```

**homebrew tap:**

```console
$ brew install k1LoW/tap/tbls
```

**macports:**

```console
$ sudo port install tbls
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/tbls/releases)

**go install:**

```console
$ go install github.com/k1LoW/tbls@latest
```

**docker:**

```console
$ docker pull ghcr.io/k1low/tbls:latest
```

**install by aqua:**

```console
$ aqua g -i k1LoW/tbls
```

## Getting Started

### Document a database

Add a `.tbls.yml` ( or `.tbls.yaml` ) file to your repository.

```yaml
# .tbls.yml

# Database document title
name: myschema

# Database to document
dsn: my://dbuser:dbpass@localhost:3306/myschema

# Path to generate document
# Default: `dbdoc`
docPath: doc/schema
```

Run `tbls doc` to analyze a database and generate document in GitHub Flavored Markdown format.

```console
$ tbls doc
```

Commit the generated document and publish it on GitHub.

```console
$ git add .
$ git commit -m 'Add database document'
$ git push origin main
```

### Diff database and (document or database)

#### `tbls diff` shows the difference between database schema and generated document.

```console
$ tbls diff
```

Currently, `tbls diff` shows the difference between:

1. `Add table`
2. `Delete table`
3. `Add column`
4. `Delete column`
5. `Change column ( type / not null )`

#### `tbls diff` also compares between databases.

```console
$ tbls diff my://root:mypass@localhost:3306/myschema my://root:mypass@localhost:3306/myschema2
```

### Lint a database

`tbls lint` work as linter for database.

```yaml
# .tbls.yml
lint:
  require:
    tableComment:
      enabled: true
      exclude:
        - logs
        - comment_*
    columnComment:
      enabled: true
      exclude:
        - id
        - created_at
        - updated_at
    indexComment:
      enabled: true
    relationName:
      enabled: true
  columnType:
    int:
      enabled: true
    varchar:
      enabled: true
  duplicateRelations:
    enabled: true
  requireForeignKeyIndex:
    enabled: true
  unrelatedTable:
    enabled: true
    exclude:
      - logs
      - comment_*
  labelStyleBigQuery:
    enabled: true
```

```console
$ tbls lint
```

### Measure document coverage

`tbls coverage` measures and show document coverage ( description, comments ).

```console
$ tbls coverage
```

### Continuous Integration

`tbls` is a CI-Friendly tool.

For example, you can add following step to GitHub Actions.

```yaml
name: Document

on:
  pull_request:

jobs:
  tbls:
    runs-on: ubuntu-latest
    steps:
      -
        uses: actions/checkout@v2
      -
        uses: k1LoW/setup-tbls@v1
      -
        run: tbls diff
        env:
          TBLS_DSN: my://root:mypass@localhost:3306/myschema
```

## Configuration

### Name

Database document title ( or `tbls doc -t title` )

```yaml
# .tbls.yml
name: myschema
```

### Description

Database document description

```yaml
# .tbls.yml
desc: This is database document for myschema
```

### Labels

Database labels are used to label databases.

The labels are used to filter databases in the index page, sort databases in the index page, and add metadata to each page.

```yaml
# .tbls.yml
labels:
  - user
  - analytics
```

### DSN

#### Support Datasource

**PostgreSQL:**

```yaml
# .tbls.yml
dsn: postgres://dbuser:dbpass@hostname:port/dbname
```

```yaml
# .tbls.yml
dsn:
  url: postgres://dbuser:dbpass@hostname:port/dbname
  # or `pg://`
  # or `postgresql://`
```

**MySQL:**

```yaml
# .tbls.yml
dsn: mysql://dbuser:dbpass@hostname:port/dbname
```

```yaml
# .tbls.yml
dsn:
  url: mysql://dbuser:dbpass@hostname:port/dbname
  # or `my://`
```

**SQLite:**

```yaml
# .tbls.yml
dsn: sqlite:///path/to/dbname.db
```

```yaml
# .tbls.yml
dsn:
  url: sqlite:///path/to/dbname.db
  # or `sq://`
```

**SQL Server:**

```yaml
# .tbls.yml
dsn: mssql://dbuser:dbpass@hostname:port/dbname
```

```yaml
# .tbls.yml
dsn:
  url: mssql://dbuser:dbpass@hostname:port/dbname
  # or `sqlserver://`
  # or `ms://`
```

**BigQuery:**

```yaml
# .tbls.yml
dsn: bigquery://project-id/dataset-id?credentialsFile=/path/to/key.json
```

```yaml
# .tbls.yml
dsn:
  url: bigquery://project-id/dataset-id?credentialsFile=/path/to/key.json
  # or `bq://`
```

**Snowflake:**

```yaml
# .tbls.yml
dsn: snowflake://user:pass@account/database/schema?warehouse=warehouse&role=role
```

**Cloud Spanner:**

```yaml
# .tbls.yml
dsn: spanner://project-id/instance-id/database-id?credentialsFile=/path/to/key.json
```

**Amazon DynamoDB:**

```yaml
# .tbls.yml
dsn: dynamodb://us-west-2?accessKeyId=XXXXXxxxxxXXXXXXXXX&secretAccessKey=XXXXXxxxxxXXXXXXXXX
```

**Amazon Redshift:**

```yaml
# .tbls.yml
dsn: redshift://user:pass@hostname:port/dbname
```

**MongoDB:**

```yaml
# .tbls.yml
dsn: mongodb://user:pass@hostname:port/dbname
```

**ClickHouse:**

```yaml
# .tbls.yml
dsn: clickhouse://dbuser:dbpass@hostname:9000/dbname
```

See also: https://pkg.go.dev/github.com/ClickHouse/clickhouse-go

**JSON:**

The JSON file output by the `tbls out -t json` command can be read as a datasource (JSON Schema is [here](spec/tbls.schema.json_schema.json)).

```yaml
---
# .tbls.yml
dsn: json://path/to/testdb.json
```

**HTTP:**

```yaml
---
# .tbls.yml
dsn: https://hostname/path/to/testdb.json
```

```yaml
---
# .tbls.yml
dsn:
  url: https://hostname/path/to/testdb.json
  headers:
    Authorization: token GITHUB_OAUTH_TOKEN
```

**GitHub:**

```yaml
---
# .tbls.yml
dsn: github://k1LoW/tbls/sample/mysql/schema.json
```

### External database driver

tbls can integrate with external database drivers. If an executable with the pattern `tbls-driver-*` is on the PATH, tbls will recognize the corresponding scheme.

For example, if you have an executable named `tbls-driver-foodb`, tbls will recognize the `foodb://` scheme.

`tbls-driver-foodb` receives the DSN at runtime via the environment variable `TBLS_DSN`. By outputting [schema.json](spec/tbls.schema.json_schema.json) via STDOUT, tbls will work with it.

### Document path

`tbls doc` generates document in the directory specified by `docPath:`.

```yaml
# .tbls.yml
# Default is `dbdoc`
docPath: doc/schema
```

### Document format

`format:` is used to change the document format.

```yaml
# .tbls.yml
format:
  # Adjust the column width of Markdown format table
  # Default is false
  adjust: true
  # Sort the order of table list and columns
  # Default is false
  sort: false
  # Display sequential numbers in table rows
  # Default is false
  number: false
  # The comments for each table in the Tables section of the index page will display the text up to the first double newline (first paragraph).
  # Default is false
  showOnlyFirstParagraph: true
  # Hide table columns without values
  # Default is false
  hideColumnsWithoutValues: true
  # It can be boolean or array
  # hideColumnsWithoutValues: ["Parents", "Children"]
```

### ER diagram

`tbls doc` generate ER diagram images at the same time.

```yaml
# .tbls.yml
er:
  # Skip generation of ER diagram
  # Default is false
  skip: false
  # ER diagram image format (`png`, `jpg`, `svg`, `mermaid`)
  # Default is `svg`
  format: svg
  # ER diagram font
  # Default is `Arial`
  font: "Trebuchet MS"
  # ER diagram font size
  # Default is `12`
  fontSize: 14
  # Show column types on ER diagram
  # Default is true
  showColumnTypes: true
  # Show only related tables on ER diagram
  showOnlyRelatedTables: false
  showColumnNullability: false
  showColumnDefaultValues: false
  showColumnComments: false
  hideColumnsWithoutValues: false
  # Set the ER diagram theme
  # Options are: 'default', 'forest', 'dark', 'neutral', 'base'
  # Default is 'default'
  mermaidTheme: default
  # Mermaid configuration JSON file path
  mermaidConfig: path/to/mermaid_config.json
  showColumnTypes:
    # Show related table column types on ER diagram
    related: true
    # Show primary key table column types on ER diagram
    primary: true
```

> Notice: `tbls` generate ER diagram images using Graphviz. Please install Graphviz or Docker.

The ER diagram can be rendered in Mermaid format by setting `format: mermaid` in the configuration.

```yaml
# .tbls.yml
er:
  format: mermaid
```

**Sample ER diagram**

<img src="https://user-images.githubusercontent.com/157344/72404768-9a6da480-378f-11ea-9ca3-84bf64077154.png" width="40%">

<details>
<summary>Show svg</summary>

<img src="https://raw.githubusercontent.com/k1LoW/tbls/main/sample/postgres/public.svg">

</details>

<details>
<summary>Show Mermaid ER diagram</summary>

~~~
erDiagram
  tables {
    table text
    options json
    table_type text
    column_value_types json
    column_value_sizes json
    schema_name text
    table_name text
  }
  columns {
    table_name text
    column_name text
    column_value_types text
    column_value_nullable bool
    column_value_default text
    column_value_primary bool
    table text
    column text
    options json
    schema_name text
  }
  relations {
    table_name text
    column_name text
    referenced_table_name text
    referenced_column_name text
    table text
    column text
    referenced_table text
    referenced_column text
    options json
    schema_name text
  }
  tables ||--o{ columns : "table"
  tables ||--o{ relations : "table"
  columns }o--|| relations : ""
~~~

</details>

### Filter tables

`include:` and `exclude:` are used to filter tables and functions.

> **Notice:** By default, views are included and functions are excluded.

```yaml
# .tbls.yml
# include tables/views/functions
include:
  - users
  - posts
  - comments
  - view_*
  # or
  - name: log_*
    labels:
      - log
  # or
  - name: products
    schema: public
  # include only tables
  - name: "function_*"
    kind: "function"
```

```yaml
# .tbls.yml
# exclude tables/views/functions
exclude:
  - logs
  - users_old
  - temp_*
```

`filterOption:` is used to change the behavior of table filtering.

```yaml
# .tbls.yml
filterOption:
  # Enable to use regular expression for include/exclude
  useRegexpFilter: true
```

If you want to filter by schema name (PostgreSQL, SQL Server):

```yaml
# .tbls.yml
# include `public` schema only
include:
  - name: "*"
    schema: public
```

### Lint

`lint:` work as linter for database.

```yaml
# .tbls.yml
lint:
  require:
    tableComment:
      enabled: true
      exclude:
        - logs
        - comment_*
    columnComment:
      enabled: true
      exclude:
        - id
        - created_at
        - updated_at
  columnType:
    int:
      enabled: true
    varchar:
      enabled: true
  duplicateRelations:
    enabled: true
  requireForeignKeyIndex:
    enabled: true
  unrelatedTable:
    enabled: true
    exclude:
      - logs
      - comment_*
  labelStyleBigQuery:
    enabled: true
```

Rules provided:

- `require.tableComment`: Require table comments
- `require.columnComment`: Require column comments
- `require.indexComment`: Require index comments
- `require.relationName`: Require relation names
- `columnType.int`: Check for integer column types
- `columnType.varchar`: Check for varchar column types
- `duplicateRelations`: Detect duplicate relations
- `requireForeignKeyIndex`: Require foreign key indexes
- `unrelatedTable`: Detect unrelated tables
- `labelStyleBigQuery`: Check for BigQuery label style

### Comments

`comments:` is used to add table/column comment to database document without `ALTER TABLE`.

For example, you can add comment about VIEW TABLE or SQLite tables/columns.

> **Notice:** Comments defined in `.tbls.yml` will override existing comments in the schema.

```yaml
# .tbls.yml
comments:
  -
    table: users
    # table comment
    tableComment: Users table
    # column comments
    columnComments:
      email: Email address as login id. ex. user@example.com
    # labels for tables
    labels:
      - privary data
      - backup:true
  -
    table: post_comments
    tableComment: post and comments View table
    columnComments:
      id: comments.id
      title: posts.title
      post_user: posts.users.username
      comment_user: comments.users.username
      created: comments.created
      updated: comments.updated
  -
    table: posts
    # index comments
    indexComments:
      posts_user_id_idx: user.id index
    # constraints comments
    constraintComments:
      posts_id_pk: PRIMARY KEY
    # triggers comments
    triggerComments:
      update_posts_updated: Update updated when posts update
```

### Comment Parsing and Logical Names

`tbls` can automatically parse database comments to extract logical names and clean descriptions using a configurable separator character. This feature helps create more readable documentation by separating business logic names from technical descriptions.

#### Configuration

```yaml
# .tbls.yml
comment:
  # Set the separator character to split logical names from descriptions
  # Default: "" (disabled)
  separator: "|"
```

#### How it works

When a separator is configured, `tbls` will parse comments in the following format:

```
LogicalName|Description
```

**Example:**

Database comment: `"User ID|System-wide unique identifier for users"`

Result:
- **Logical Name**: `User ID`
- **Clean Comment**: `System-wide unique identifier for users`

#### Benefits

- **Bilingual Support**: Perfect for international projects where logical names are in local language and descriptions in English
- **Better Documentation**: Cleaner separation between what a field represents (logical name) and how it works (description)
- **Consistent Formatting**: Automatic parsing ensures consistent documentation format

#### Usage Examples

**PostgreSQL:**
```sql
COMMENT ON COLUMN users.id IS 'ユーザーID|System-wide unique identifier for users';
COMMENT ON TABLE users IS 'ユーザーテーブル|Stores user account information';
```

**MySQL:**
```sql
CREATE TABLE users (
  id INT PRIMARY KEY COMMENT 'ユーザーID|System-wide unique identifier for users',
  name VARCHAR(100) COMMENT 'ユーザー名|Display name for the user'
);
```

### Markdown Output Customization

`tbls` provides powerful customization options for Markdown output, allowing you to control column ordering, field aliases, and logical name display for different database objects.

#### Configuration Structure

```yaml
# .tbls.yml
markdown:
  # Database-level customization
  database:
    show_logical_name: true
    order: ["name", "logical_name", "comment"]
    aliases:
      name: "Database Name"
      logical_name: "Business Name"
      comment: "Description"

  # Table customization
  tables:
    show_logical_name: true
    order: ["name", "logical_name", "comment", "type"]
    aliases:
      name: "Table Name"
      logical_name: "Business Name"
      comment: "Description"
      type: "Type"
    # Table-specific overrides
    specific:
      users:
        show_logical_name: true
        aliases:
          name: "User Table"

  # Column customization
  columns:
    show_logical_name: true
    order: ["name", "logical_name", "type", "nullable", "default", "comment"]
    aliases:
      name: "Column Name"
      logical_name: "Business Name"
      type: "Data Type"
      nullable: "Nullable"
      default: "Default Value"
      comment: "Description"
    # Table-specific column settings
    specific:
      users:
        order: ["logical_name", "name", "type", "nullable", "comment"]

  # Other object types
  views:
    show_logical_name: true
    # ... similar configuration

  indexes:
    show_logical_name: true
    # ... similar configuration

  constraints:
    show_logical_name: true
    # ... similar configuration

  functions:
    show_logical_name: true
    # ... similar configuration
```

#### Features

**1. Column Ordering**
- Customize the order of columns in Markdown tables
- Different orders for different object types
- Table-specific column ordering overrides

**2. Field Aliases**
- Rename column headers to more user-friendly names
- Support for multiple languages (e.g., English headers → Japanese headers)
- Consistent terminology across documentation

**3. Logical Name Display**
- Show/hide logical names extracted from comments
- Logical names appear as separate columns when enabled
- Fallback to physical names when logical names are not available

**4. Object-Specific Customization**
- Global settings for each object type (tables, columns, views, etc.)
- Specific overrides for individual tables or objects
- Hierarchical configuration (global → object type → specific)

#### Example Output

With the above configuration, a table documentation might look like:

| Table Name | Business Name | Description | Type |
|------------|---------------|-------------|------|
| users      | ユーザーテーブル | Stores user information | BASE TABLE |

| Column Name | Business Name | Data Type | Nullable | Default Value | Description |
|-------------|---------------|-----------|----------|---------------|-------------|
| id          | ユーザーID    | int       | NO       | NULL          | Unique user identifier |
| name        | ユーザー名    | varchar(100) | YES   | NULL          | User display name |

### Relations

`relations:` is used to add or override table relation to database document without `FOREIGN KEY`.

You can create ER diagrams with relations without having foreign key constraints.

```yaml
relations:
  -
    table: logs
    columns:
      - user_id
    parentTable: users
    parentColumns:
      - id
    # Relation definition
    def: logs->users
  -
    table: logs
    columns:
      - post_id
    parentTable: posts
    parentColumns:
      - id
    def: logs->posts
```

### Viewpoints

`viewpoints:` is used to generate additional documents (add/change ER diagram, table list).

```yaml
viewpoints:
  -
    name: users
    desc: Users and posts
    tables:
      - users
      - posts
      - comments
    # Relations viewpoint
    relations:
      - table: posts
        columns:
          - user_id
        parentTable: users
        parentColumns:
          - id
        def: posts->users
  -
    name: payment
    desc: Payment
    tables:
      - name: orders
        comment: Order table
      - payments
    labels:
      - payment
      - public
```

### Dictionary

`dict:` is used to generate additional document with merged table and column info.

```yaml
# .tbls.yml
dict:
  -
    name: analytics
    desc: Analytics tables list
    tables:
      - analytics_*
    labels:
      - analytics
```

**Generated additional document**

![sample dictionary](https://github.com/user-attachments/assets/dcf64e8c-c264-4ad9-8b2b-7eb6c5b2f4be)

### Personalized Templates

`tbls` uses Go templates to generate Markdown documents. You can create personalized templates.

See [here](output/md/templates) for the default templates.

```yaml
# .tbls.yml
templates:
  md:
    index: path/to/index.md.tmpl
    table: path/to/table.md.tmpl
```

### Required Version

The `requiredVersion` can be used to set the minimum required version of tbls.

```yaml
# .tbls.yml
requiredVersion: ">= 1.42.0"
```

## Expand environment variables

`tbls` expand environment variables using `${VAR}` in `.tbls.yml`

```yaml
# .tbls.yml
dsn: postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}
name: ${POSTGRES_DATABASE}
```

## Output formats

`tbls` can output in various formats by using `tbls out` command.

```console
$ tbls out [format]
```

### Markdown

```console
$ tbls out md
```

**Output directory structure**

```console
docs/ # `docPath:`
├── README.md # Database
├── users.md # Table
└── posts.md
```

### DOT

```console
$ tbls out dot
```

### PlantUML

```console
$ tbls out plantuml
```

The PlantUML file can be converted to PNG format.

```console
$ cat dbdoc/schema.puml | docker run --rm -i think/plantuml -tpng > schema.png
```

**Schema PNG**

<img src="https://user-images.githubusercontent.com/157344/62372115-a7d8c080-b571-11e9-9431-a8a2b2bd7c64.png" alt="schema.png" width="100%">

### JSON

```console
$ tbls out json
```

### YAML

```console
$ tbls out yaml
```

### XML

```console
$ tbls out xml
```

### XLSX

```console
$ tbls out xlsx
```

### Mermaid

```console
$ tbls out mermaid
```

## Command arguments

### Common arguments

- `--config` (`-c`): config file path.
- `--dsn` (`-d`): data source name.
- `--schema` (`-s`): schema name.
- `--when`: when to run `tbls`.
- `--debug`: debug mode.

### `tbls doc`

Generate database document.

#### Arguments

- `--force`: force-run (skip asking)
- `--rm-dist`: remove files in docPath before generating documents
- `--adjust` (`-a`): adjust column width
- `--sort`: sort
- `--number`: add number to table rows
- `--ER` (`-j`): also generate ER diagram
- `--add-relation-from-er`: add relations from ER diagram
- `--without-er`: generate without ER diagram

### `tbls diff`

Show diff between database and generated document.

#### Arguments

- `--with-color`: show colored diff
- `--unique`: show unique tables only

If you run `tbls diff path/to/docPath.md`, it only outputs the diff for the table of the specified file.

### `tbls lint`

Lint the database.

#### Arguments

- `--format`: lint rule format
- `--fix`: fix lint (beta)

### `tbls coverage`

Show document coverage.

#### Arguments

- `--format`: output format for coverage (default is table)

### `tbls out`

Output in various formats.

#### Arguments

- `--template` (`-t`): template type
- `--output` (`-o`): output file path

### `tbls md2dot`

Convert Markdown to DOT.

#### Arguments

- `--output` (`-o`): output file path

### `tbls completion`

Generate shell completion.

```console
$ source <(tbls completion bash)
```

```console
$ source <(tbls completion zsh)
```

```powershell
PS > tbls completion powershell | Out-String | Invoke-Expression
```

```console
$ tbls completion fish | source
```

### `tbls version`

Print the version number.

### `tbls help`

Show help.

## Environment variables

tbls supports environment variables.

### `TBLS_DSN`

You can use the environment variable `TBLS_DSN` instead of the command argument `--dsn`.

```console
$ env TBLS_DSN=my://dbuser:dbpass@localhost:3306/myschema tbls doc
```

### `TBLS_DOC_PATH`

You can use the environment variable `TBLS_DOC_PATH` instead of the `.tbls.yml` setting `docPath`.

```console
$ env TBLS_DOC_PATH=./custom-doc-path tbls doc
```

### `TBLS_CONFIG_PATH`

You can use the environment variable `TBLS_CONFIG_PATH` instead of the command argument `--config`.

```console
$ env TBLS_CONFIG_PATH=./custom.tbls.yml tbls doc
```

### `TBLS_WHEN`

You can use the environment variable `TBLS_WHEN` instead of the command argument `--when`.

```console
$ env TBLS_WHEN="command arg" tbls doc
```

### `TBLS_DEBUG`

You can use the environment variable `TBLS_DEBUG` instead of the command argument `--debug`.

```console
$ env TBLS_DEBUG=true tbls doc
```

## External subcommands

`tbls` supports external subcommands.

If an executable with the pattern `tbls-*` is on the PATH, `tbls` will treat it as a subcommand.

```console
$ tbls subcommand
```

For example, you can use [tbls-ask](https://github.com/k1LoW/tbls-ask) by installing it on your PATH.

```console
$ tbls ask "What is the average age of users?"
```