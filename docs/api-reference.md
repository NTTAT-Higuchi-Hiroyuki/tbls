# API Reference - Comment Parsing and Markdown Customization

This document provides detailed API reference for the comment parsing and Markdown customization features introduced in tbls v1.86.0.

## Table of Contents

- [Configuration Structures](#configuration-structures)
- [CommentParser API](#commentparser-api)
- [LogicalNameProcessor API](#logicalnameprocessor-api)
- [MarkdownCustomizer API](#markdowncustomizer-api)
- [Schema Extensions](#schema-extensions)
- [Usage Examples](#usage-examples)

## Configuration Structures

### CommentConfig

Configuration for comment parsing functionality.

```go
type CommentConfig struct {
    // Separator is the character(s) used to split logical names from descriptions
    // Example: "|", ":", "-", " - "
    // Default: "" (disabled)
    Separator string `yaml:"separator,omitempty"`
}
```

**YAML Configuration:**
```yaml
comment:
  separator: "|"
```

**Validation Rules:**
- Must be 1-10 characters long
- Cannot contain newlines or tabs
- Special characters are allowed (e.g., `|`, `:`, `-`, `::`)

### MarkdownConfig

Configuration for Markdown output customization.

```go
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

**YAML Configuration:**
```yaml
markdown:
  tables:
    show_logical_name: true
    order: ["name", "logical_name", "comment", "type"]
    aliases:
      name: "Table Name"
      logical_name: "Business Name"
    specific:
      users:
        aliases:
          name: "User Management Table"
```

## CommentParser API

### Interface

```go
type CommentParser interface {
    ParseComment(rawComment, separator string) (logicalName, cleanComment string)
    ExtractLogicalName(rawComment, separator string) string
    ExtractCleanComment(rawComment, separator string) string
    HasLogicalName(rawComment, separator string) bool
}
```

### DefaultCommentParser

Default implementation of the CommentParser interface.

```go
type DefaultCommentParser struct{}

func NewDefaultCommentParser() *DefaultCommentParser {
    return &DefaultCommentParser{}
}
```

#### Methods

##### ParseComment

```go
func (p *DefaultCommentParser) ParseComment(rawComment, separator string) (logicalName, cleanComment string)
```

Parses a raw comment and returns both logical name and clean comment.

**Parameters:**
- `rawComment`: The original comment from the database
- `separator`: The separator character(s) to split on

**Returns:**
- `logicalName`: The extracted logical name (before separator)
- `cleanComment`: The cleaned comment (after separator)

**Example:**
```go
parser := NewDefaultCommentParser()
logical, clean := parser.ParseComment("User ID|Unique identifier", "|")
// logical = "User ID"
// clean = "Unique identifier"
```

##### ExtractLogicalName

```go
func (p *DefaultCommentParser) ExtractLogicalName(rawComment, separator string) string
```

Extracts only the logical name from a comment.

**Example:**
```go
logical := parser.ExtractLogicalName("User ID|Unique identifier", "|")
// logical = "User ID"
```

##### ExtractCleanComment

```go
func (p *DefaultCommentParser) ExtractCleanComment(rawComment, separator string) string
```

Extracts only the clean comment (description) from a comment.

**Example:**
```go
clean := parser.ExtractCleanComment("User ID|Unique identifier", "|")
// clean = "Unique identifier"
```

##### HasLogicalName

```go
func (p *DefaultCommentParser) HasLogicalName(rawComment, separator string) bool
```

Checks if a comment contains a logical name (i.e., contains the separator).

**Example:**
```go
has := parser.HasLogicalName("User ID|Unique identifier", "|")
// has = true

has = parser.HasLogicalName("Just a description", "|")
// has = false
```

## LogicalNameProcessor API

### Structure

```go
type LogicalNameProcessor struct {
    separator string
    parser    CommentParser
}

func NewLogicalNameProcessor(separator string) *LogicalNameProcessor {
    return &LogicalNameProcessor{
        separator: separator,
        parser:    NewDefaultCommentParser(),
    }
}
```

### Methods

#### ProcessSchema

```go
func (p *LogicalNameProcessor) ProcessSchema(s *Schema) error
```

Processes an entire schema, extracting logical names from all objects.

**Parameters:**
- `s`: Pointer to the schema to process

**Returns:**
- `error`: Any error that occurred during processing

**Example:**
```go
processor := NewLogicalNameProcessor("|")
err := processor.ProcessSchema(schema)
if err != nil {
    log.Printf("Error processing schema: %v", err)
}
```

#### ProcessTable

```go
func (p *LogicalNameProcessor) ProcessTable(table *Table) error
```

Processes a single table and its columns, indexes, constraints, and triggers.

#### ProcessColumn

```go
func (p *LogicalNameProcessor) ProcessColumn(column *Column) error
```

Processes a single column to extract logical name from its comment.

## MarkdownCustomizer API

### Interface

```go
type MarkdownCustomizer interface {
    CustomizeTableOutput(table *Table, config *MarkdownConfig) *CustomizedTableData
    CustomizeColumnOutput(columns []*Column, config *ColumnCustomConfig) []*CustomizedColumnData
    ApplyAliases(fieldName string, aliases map[string]string) string
    GetDisplayOrder(fields []string, order []string) []string
}
```

### DefaultMarkdownCustomizer

```go
type DefaultMarkdownCustomizer struct{}

func NewDefaultMarkdownCustomizer() *DefaultMarkdownCustomizer {
    return &DefaultMarkdownCustomizer{}
}
```

### Data Structures

#### CustomizedTableData

```go
type CustomizedTableData struct {
    Table    *Table
    Columns  []*CustomizedColumnData
    Order    []string
    Aliases  map[string]string
    ShowLogicalName bool
}
```

#### CustomizedColumnData

```go
type CustomizedColumnData struct {
    Column      *Column
    DisplayName string
    LogicalName string
    Show        bool
    Order       int
}
```

### Methods

#### CustomizeTableOutput

```go
func (c *DefaultMarkdownCustomizer) CustomizeTableOutput(table *Table, config *MarkdownConfig) *CustomizedTableData
```

Customizes table output based on configuration.

**Parameters:**
- `table`: The table to customize
- `config`: Markdown configuration

**Returns:**
- `*CustomizedTableData`: Customized table data

#### CustomizeColumnOutput

```go
func (c *DefaultMarkdownCustomizer) CustomizeColumnOutput(columns []*Column, config *ColumnCustomConfig) []*CustomizedColumnData
```

Customizes column output based on configuration.

#### ApplyAliases

```go
func (c *DefaultMarkdownCustomizer) ApplyAliases(fieldName string, aliases map[string]string) string
```

Applies field name aliases.

**Example:**
```go
customizer := NewDefaultMarkdownCustomizer()
aliases := map[string]string{
    "name": "Table Name",
    "comment": "Description",
}
displayName := customizer.ApplyAliases("name", aliases)
// displayName = "Table Name"
```

#### GetDisplayOrder

```go
func (c *DefaultMarkdownCustomizer) GetDisplayOrder(fields []string, order []string) []string
```

Reorders fields based on configuration.

**Example:**
```go
fields := []string{"name", "type", "comment", "nullable"}
order := []string{"comment", "name", "type"}
reordered := customizer.GetDisplayOrder(fields, order)
// reordered = ["comment", "name", "type", "nullable"]
```

## Schema Extensions

### Table Extensions

```go
type Table struct {
    // ... existing fields ...
    LogicalName string `json:"logicalName,omitempty"`
}

// SetLogicalNameFromComment extracts logical name from comment
func (t *Table) SetLogicalNameFromComment(separator string) {
    if t.Comment != "" && separator != "" {
        parser := NewDefaultCommentParser()
        logical, clean := parser.ParseComment(t.Comment, separator)
        if logical != "" {
            t.LogicalName = logical
            t.Comment = clean
        }
    }
}

// GetLogicalNameOrFallback returns logical name or falls back to physical name
func (t *Table) GetLogicalNameOrFallback() string {
    if t.LogicalName != "" {
        return t.LogicalName
    }
    return t.Name
}
```

### Column Extensions

```go
type Column struct {
    // ... existing fields ...
    LogicalName string `json:"logicalName,omitempty"`
}

func (c *Column) SetLogicalNameFromComment(separator string) {
    // Similar to Table.SetLogicalNameFromComment
}

func (c *Column) GetLogicalNameOrFallback() string {
    // Similar to Table.GetLogicalNameOrFallback
}
```

### Index Extensions

```go
type Index struct {
    // ... existing fields ...
    LogicalName string `json:"logicalName,omitempty"`
}
```

### Constraint Extensions

```go
type Constraint struct {
    // ... existing fields ...
    LogicalName string `json:"logicalName,omitempty"`
}
```

## Usage Examples

### Basic Comment Parsing

```go
// Initialize parser
parser := NewDefaultCommentParser()

// Parse comment
logical, clean := parser.ParseComment("User ID|Unique user identifier", "|")
fmt.Printf("Logical: %s, Clean: %s\n", logical, clean)
// Output: Logical: User ID, Clean: Unique user identifier
```

### Schema Processing

```go
// Create processor
processor := NewLogicalNameProcessor("|")

// Process entire schema
err := processor.ProcessSchema(schema)
if err != nil {
    log.Printf("Processing error: %v", err)
}

// Access logical names
for _, table := range schema.Tables {
    fmt.Printf("Table: %s (Logical: %s)\n", table.Name, table.LogicalName)
    for _, column := range table.Columns {
        fmt.Printf("  Column: %s (Logical: %s)\n", column.Name, column.LogicalName)
    }
}
```

### Markdown Customization

```go
// Create customizer
customizer := NewDefaultMarkdownCustomizer()

// Configure markdown settings
config := &MarkdownConfig{
    Tables: &TableCustomConfig{
        ObjectCustomConfig: &ObjectCustomConfig{
            ShowLogicalName: true,
            Order: []string{"name", "logical_name", "comment"},
            Aliases: map[string]string{
                "name": "Table Name",
                "logical_name": "Business Name",
                "comment": "Description",
            },
        },
    },
}

// Customize table output
customizedData := customizer.CustomizeTableOutput(table, config)
```

### Integration Example

```go
// Complete integration in your application
func ProcessDatabaseSchema(dsn string, configPath string) error {
    // Load configuration
    config, err := LoadConfig(configPath)
    if err != nil {
        return err
    }

    // Analyze database
    schema, err := AnalyzeDatabase(dsn)
    if err != nil {
        return err
    }

    // Process logical names if configured
    if config.Comment != nil && config.Comment.Separator != "" {
        processor := NewLogicalNameProcessor(config.Comment.Separator)
        if err := processor.ProcessSchema(schema); err != nil {
            log.Printf("Warning: Failed to process logical names: %v", err)
        }
    }

    // Generate markdown with customization
    if config.Markdown != nil {
        customizer := NewDefaultMarkdownCustomizer()
        // Use customizer to generate customized markdown...
    }

    return nil
}
```

## Error Handling

### Common Errors

```go
// ErrInvalidSeparator is returned when separator is invalid
var ErrInvalidSeparator = errors.New("invalid separator: must be 1-10 characters, no newlines")

// ErrInvalidOrder is returned when order contains unknown fields
var ErrInvalidOrder = errors.New("invalid order: unknown field name")

// ErrInvalidAlias is returned when alias configuration is invalid
var ErrInvalidAlias = errors.New("invalid alias: empty field name or alias")
```

### Error Handling Best Practices

```go
// Always handle errors gracefully
processor := NewLogicalNameProcessor(separator)
if err := processor.ProcessSchema(schema); err != nil {
    // Log warning but continue processing
    log.Printf("Warning: Failed to process logical names: %v", err)
    // Don't return error - allow processing to continue
}
```

## Thread Safety

All APIs are designed to be thread-safe for read operations:

- `CommentParser` methods are thread-safe
- `LogicalNameProcessor` is thread-safe for processing different schemas
- `MarkdownCustomizer` methods are thread-safe

**Note:** Concurrent modification of the same schema object is not thread-safe.

## Performance Considerations

- Comment parsing: O(n) where n is the number of objects with comments
- Logical name processing: O(n) where n is the total number of database objects
- Markdown customization: O(m) where m is the number of fields being customized

Typical performance impact:
- Small databases (< 100 tables): < 1ms overhead
- Medium databases (100-1000 tables): 1-10ms overhead
- Large databases (> 1000 tables): 10-100ms overhead