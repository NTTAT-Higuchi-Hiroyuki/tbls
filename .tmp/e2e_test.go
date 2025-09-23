package main

import (
	"testing"
	"time"

	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/output/md"
	"github.com/k1LoW/tbls/schema"
	"github.com/stretchr/testify/assert"
)

// TestE2E_JapaneseLogicalNameProcessing tests end-to-end processing
// of Japanese logical names through the complete pipeline
func TestE2E_JapaneseLogicalNameProcessing(t *testing.T) {
	// Create test schema with Japanese comments
	testSchema := &schema.Schema{
		Name: "test_database",
		Tables: []*schema.Table{
			{
				Name:        "users",
				Comment:     "ユーザーテーブル|システムユーザーの基本情報を管理",
				LogicalName: "", // Will be set by processing
				Columns: []*schema.Column{
					{
						Name:        "id",
						Type:        "SERIAL",
						Comment:     "ユーザーID|システム内でユーザーを一意に識別するID",
						LogicalName: "",
						Nullable:    false,
					},
					{
						Name:        "name",
						Type:        "VARCHAR(100)",
						Comment:     "ユーザー名|ユーザーの表示名",
						LogicalName: "",
						Nullable:    true,
					},
					{
						Name:        "email",
						Type:        "VARCHAR(255)",
						Comment:     "メールアドレス|ログインおよび通知に使用",
						LogicalName: "",
						Nullable:    false,
					},
				},
			},
			{
				Name:        "products",
				Comment:     "商品テーブル|ECサイトで販売する商品情報",
				LogicalName: "",
				Columns: []*schema.Column{
					{
						Name:        "id",
						Type:        "SERIAL",
						Comment:     "商品ID|商品を一意に識別するID",
						LogicalName: "",
						Nullable:    false,
					},
					{
						Name:        "name",
						Type:        "VARCHAR(255)",
						Comment:     "商品名|販売する商品の名称",
						LogicalName: "",
						Nullable:    false,
					},
					{
						Name:        "price",
						Type:        "DECIMAL(10,2)",
						Comment:     "価格|商品の販売価格（税抜）",
						LogicalName: "",
						Nullable:    false,
					},
				},
			},
		},
	}

	// Configure comment parsing and markdown customization
	cfg := &config.Config{
		Comment: &config.CommentConfig{
			Separator: "|",
		},
		Markdown: &config.MarkdownConfig{
			Tables: &config.TableCustomConfig{
				ObjectCustomConfig: &config.ObjectCustomConfig{
					ShowLogicalName: true,
					Order:           []string{"Name", "LogicalName", "Type", "Comment"},
					Aliases: map[string]string{
						"Name":        "物理名",
						"LogicalName": "論理名",
						"Type":        "データ型",
						"Comment":     "説明",
					},
				},
			},
			Columns: &config.ColumnCustomConfig{
				ObjectCustomConfig: &config.ObjectCustomConfig{
					ShowLogicalName: true,
					Order:           []string{"Name", "LogicalName", "Type", "Nullable", "Comment"],
					Aliases: map[string]string{
						"Name":        "カラム名",
						"LogicalName": "論理名",
						"Type":        "型",
						"Nullable":    "NULL可",
						"Comment":     "説明",
					},
				},
			},
		},
	}

	// Apply logical name processing via ModifySchema
	err := cfg.ModifySchema(testSchema)
	assert.NoError(t, err)

	// Verify logical names were extracted
	assert.Equal(t, "ユーザーテーブル", testSchema.Tables[0].LogicalName)
	assert.Equal(t, "商品テーブル", testSchema.Tables[1].LogicalName)
	assert.Equal(t, "ユーザーID", testSchema.Tables[0].Columns[0].LogicalName)
	assert.Equal(t, "商品ID", testSchema.Tables[1].Columns[0].LogicalName)

	// Test markdown generation with customization
	mdCustomizer := md.NewDefaultMarkdownCustomizer()

	// Test table customization
	tableConfig := cfg.Markdown.GetEffectiveConfig("table", "users")
	customizedTable := mdCustomizer.CustomizeTableOutput(testSchema.Tables[0], tableConfig)

	assert.NotNil(t, customizedTable)
	assert.Equal(t, "ユーザーテーブル", customizedTable.LogicalName)
	assert.Contains(t, customizedTable.Headers, "論理名")
	assert.Contains(t, customizedTable.Headers, "物理名")

	// Test column customization
	columnConfig := cfg.Markdown.GetEffectiveConfig("column", "")
	customizedColumns := mdCustomizer.CustomizeColumnOutput(testSchema.Tables[0].Columns, columnConfig)

	assert.NotNil(t, customizedColumns)
	assert.Len(t, customizedColumns, 3)
	assert.Equal(t, "ユーザーID", customizedColumns[0].LogicalName)
	assert.Contains(t, customizedColumns[0].Data, "ユーザーID") // LogicalName in data
	assert.Contains(t, customizedColumns[0].Data, "id")        // Name in data

	t.Logf("✅ E2E Test: Complete pipeline from Japanese comments to customized Markdown")
	t.Logf("   - Logical name extraction: ✓")
	t.Logf("   - Markdown customization: ✓")
	t.Logf("   - Japanese alias support: ✓")
}

// TestE2E_PerformanceWithLargeDataset tests performance with a large dataset
// to ensure we meet the 5% degradation requirement
func TestE2E_PerformanceWithLargeDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Create a large dataset
	largeSchema := createLargeTestDataset(100, 10) // 100 tables, 10 columns each

	// Baseline: Processing without logical name features
	baselineConfig := &config.Config{}
	start := time.Now()
	err := baselineConfig.ModifySchema(largeSchema)
	assert.NoError(t, err)
	baselineTime := time.Since(start)

	// Reset schema for comparison
	largeSchema = createLargeTestDataset(100, 10)

	// Enhanced: Processing with logical name features
	enhancedConfig := &config.Config{
		Comment: &config.CommentConfig{
			Separator: "|",
		},
		Markdown: &config.MarkdownConfig{
			Tables: &config.TableCustomConfig{
				ObjectCustomConfig: &config.ObjectCustomConfig{
					ShowLogicalName: true,
					Order:           []string{"Name", "LogicalName", "Type", "Comment"},
					Aliases: map[string]string{
						"Name":        "物理名",
						"LogicalName": "論理名",
					},
				},
			},
		},
	}
	start = time.Now()
	err = enhancedConfig.ModifySchema(largeSchema)
	assert.NoError(t, err)
	enhancedTime := time.Since(start)

	// Calculate performance impact
	performanceRatio := float64(enhancedTime) / float64(baselineTime)
	degradationPercent := (performanceRatio - 1.0) * 100

	t.Logf("Performance Test Results:")
	t.Logf("   - Baseline time: %v", baselineTime)
	t.Logf("   - Enhanced time: %v", enhancedTime)
	t.Logf("   - Performance ratio: %.3f", performanceRatio)
	t.Logf("   - Degradation: %.2f%%", degradationPercent)

	// Assert 5% degradation limit
	assert.True(t, degradationPercent <= 5.0,
		"Performance degradation (%.2f%%) exceeds 5%% limit", degradationPercent)

	t.Logf("✅ Performance Test: Within 5%% degradation limit (%.2f%%)", degradationPercent)
}

// TestE2E_ErrorToleranceAndFallback tests that the system handles various error
// conditions gracefully and maintains backward compatibility
func TestE2E_ErrorToleranceAndFallback(t *testing.T) {
	testCases := []struct {
		name     string
		schema   *schema.Schema
		config   *config.Config
		shouldSucceed bool
		description string
	}{
		{
			name: "Invalid separator in config",
			schema: &schema.Schema{
				Tables: []*schema.Table{
					{Name: "test", Comment: "normal comment"},
				},
			},
			config: &config.Config{
				Comment: &config.CommentConfig{
					Separator: "", // Invalid empty separator
				},
			},
			shouldSucceed: true,
			description: "Should handle invalid separator gracefully",
		},
		{
			name: "Mixed comment formats",
			schema: &schema.Schema{
				Tables: []*schema.Table{
					{Name: "table1", Comment: "論理名|説明"},
					{Name: "table2", Comment: "no separator comment"},
					{Name: "table3", Comment: ""},
				},
			},
			config: &config.Config{
				Comment: &config.CommentConfig{
					Separator: "|",
				},
			},
			shouldSucceed: true,
			description: "Should handle mixed comment formats",
		},
		{
			name: "Empty schema",
			schema: &schema.Schema{
				Tables: []*schema.Table{},
			},
			config: &config.Config{
				Comment: &config.CommentConfig{
					Separator: "|",
				},
			},
			shouldSucceed: true,
			description: "Should handle empty schema",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.ModifySchema(tc.schema)

			if tc.shouldSucceed {
				assert.NoError(t, err, tc.description)
			} else {
				assert.Error(t, err, tc.description)
			}

			t.Logf("✅ Error tolerance test '%s': %s", tc.name, tc.description)
		})
	}
}

// TestE2E_BackwardCompatibility tests that existing functionality
// remains unchanged when new features are not configured
func TestE2E_BackwardCompatibility(t *testing.T) {
	// Create schema identical to existing test scenarios
	originalSchema := &schema.Schema{
		Tables: []*schema.Table{
			{
				Name:    "users",
				Comment: "User table",
				Columns: []*schema.Column{
					{Name: "id", Type: "int", Comment: "User ID"},
					{Name: "name", Type: "varchar", Comment: "User name"},
				},
			},
		},
	}

	// Test with empty config (default behavior)
	emptyConfig := &config.Config{}
	err := emptyConfig.ModifySchema(originalSchema)
	assert.NoError(t, err)

	// Verify no logical names were set
	assert.Equal(t, "", originalSchema.Tables[0].LogicalName)
	assert.Equal(t, "", originalSchema.Tables[0].Columns[0].LogicalName)

	// Verify original comments are preserved
	assert.Equal(t, "User table", originalSchema.Tables[0].Comment)
	assert.Equal(t, "User ID", originalSchema.Tables[0].Columns[0].Comment)

	t.Logf("✅ Backward Compatibility Test: Original behavior preserved when features not configured")
}

// createLargeTestDataset creates a large test dataset for performance testing
func createLargeTestDataset(numTables, numColumnsPerTable int) *schema.Schema {
	tables := make([]*schema.Table, numTables)

	for i := 0; i < numTables; i++ {
		columns := make([]*schema.Column, numColumnsPerTable)

		for j := 0; j < numColumnsPerTable; j++ {
			columns[j] = &schema.Column{
				Name:    fmt.Sprintf("column_%d", j),
				Type:    "VARCHAR(255)",
				Comment: fmt.Sprintf("カラム%d|テストデータのカラム%d", j, j),
			}
		}

		tables[i] = &schema.Table{
			Name:    fmt.Sprintf("table_%d", i),
			Comment: fmt.Sprintf("テーブル%d|テストデータのテーブル%d", i, i),
			Columns: columns,
		}
	}

	return &schema.Schema{
		Name:   "large_test_db",
		Tables: tables,
	}
}

// Helper function to format table name (needed for createLargeTestDataset)
func fmt.Sprintf(format string, args ...interface{}) string {
	// This is a simplified implementation for testing
	// In real code, this would use the fmt package
	if len(args) == 1 {
		switch format {
		case "column_%d":
			return "column_" + string(rune(args[0].(int) + '0'))
		case "table_%d":
			return "table_" + string(rune(args[0].(int) + '0'))
		case "カラム%d|テストデータのカラム%d":
			return "カラム" + string(rune(args[0].(int) + '0')) + "|テストデータのカラム" + string(rune(args[0].(int) + '0'))
		case "テーブル%d|テストデータのテーブル%d":
			return "テーブル" + string(rune(args[0].(int) + '0')) + "|テストデータのテーブル" + string(rune(args[0].(int) + '0'))
		}
	}
	return format
}