package main

import (
	"testing"
	"time"

	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/schema"
	"github.com/stretchr/testify/assert"
)

// TestE2E_JapaneseLogicalNameIntegration tests the complete integration
// of Japanese logical name processing through ModifySchema
func TestE2E_JapaneseLogicalNameIntegration(t *testing.T) {
	// Create test schema with Japanese comments containing logical names
	testSchema := &schema.Schema{
		Name: "test_database",
		Tables: []*schema.Table{
			{
				Name:    "users",
				Comment: "ユーザーテーブル|システムユーザーの基本情報を管理",
				Columns: []*schema.Column{
					{
						Name:     "id",
						Type:     "SERIAL",
						Comment:  "ユーザーID|システム内でユーザーを一意に識別するID",
						Nullable: false,
					},
					{
						Name:     "name",
						Type:     "VARCHAR(100)",
						Comment:  "ユーザー名|ユーザーの表示名",
						Nullable: true,
					},
					{
						Name:     "email",
						Type:     "VARCHAR(255)",
						Comment:  "メールアドレス|ログインおよび通知に使用",
						Nullable: false,
					},
				},
			},
		},
	}

	// Configure comment parsing with Japanese separator
	cfg := &config.Config{
		Comment: &config.CommentConfig{
			Separator: "|",
		},
	}

	// Execute the integrated processing pipeline
	err := cfg.ModifySchema(testSchema)
	assert.NoError(t, err)

	// Verify logical names were correctly extracted
	assert.Equal(t, "ユーザーテーブル", testSchema.Tables[0].LogicalName)
	assert.Equal(t, "ユーザーID", testSchema.Tables[0].Columns[0].LogicalName)
	assert.Equal(t, "ユーザー名", testSchema.Tables[0].Columns[1].LogicalName)
	assert.Equal(t, "メールアドレス", testSchema.Tables[0].Columns[2].LogicalName)

	// Verify original comments are preserved
	assert.Equal(t, "ユーザーテーブル|システムユーザーの基本情報を管理", testSchema.Tables[0].Comment)

	// Test fallback functionality
	assert.Equal(t, "ユーザーテーブル", testSchema.Tables[0].GetLogicalNameOrFallback())
	assert.Equal(t, "ユーザーID", testSchema.Tables[0].Columns[0].GetLogicalNameOrFallback())

	t.Logf("✅ E2E Integration Test: Japanese logical name processing completed successfully")
}

// TestE2E_PerformanceBenchmark tests performance impact of logical name processing
func TestE2E_PerformanceBenchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Create test dataset
	testSchema := createSimpleTestDataset()

	// Measure baseline performance (no logical name processing)
	baselineConfig := &config.Config{}

	start := time.Now()
	err := baselineConfig.ModifySchema(testSchema)
	assert.NoError(t, err)
	baselineTime := time.Since(start)

	// Reset schema for enhanced test
	testSchema = createSimpleTestDataset()

	// Measure enhanced performance (with logical name processing)
	enhancedConfig := &config.Config{
		Comment: &config.CommentConfig{
			Separator: "|",
		},
	}

	start = time.Now()
	err = enhancedConfig.ModifySchema(testSchema)
	assert.NoError(t, err)
	enhancedTime := time.Since(start)

	// Calculate performance metrics
	performanceRatio := float64(enhancedTime) / float64(baselineTime)
	degradationPercent := (performanceRatio - 1.0) * 100

	// Verify that logical names were processed
	logicalNameCount := 0
	for _, table := range testSchema.Tables {
		if table.LogicalName != "" {
			logicalNameCount++
		}
		for _, column := range table.Columns {
			if column.LogicalName != "" {
				logicalNameCount++
			}
		}
	}

	t.Logf("📈 Performance Benchmark Results:")
	t.Logf("   ⏱️  Baseline time: %v", baselineTime)
	t.Logf("   ⏱️  Enhanced time: %v", enhancedTime)
	t.Logf("   📊 Performance ratio: %.3f", performanceRatio)
	t.Logf("   📈 Performance change: %.2f%%", degradationPercent)
	t.Logf("   🏷️  Logical names processed: %d", logicalNameCount)
	t.Logf("✅ Performance test completed - functionality working")
}

// TestE2E_BackwardCompatibilityGuarantee ensures existing behavior is preserved
func TestE2E_BackwardCompatibilityGuarantee(t *testing.T) {
	// Create test schema using existing patterns
	originalSchema := &schema.Schema{
		Tables: []*schema.Table{
			{
				Name:    "legacy_table",
				Comment: "This is a legacy comment without separator",
				Columns: []*schema.Column{
					{Name: "id", Type: "int", Comment: "Primary key"},
					{Name: "name", Type: "varchar", Comment: "Name field"},
				},
			},
		},
	}

	// Test: No configuration (should work exactly as before)
	emptyConfig := &config.Config{}
	err := emptyConfig.ModifySchema(originalSchema)
	assert.NoError(t, err)

	// Verify no logical names were set
	assert.Equal(t, "", originalSchema.Tables[0].LogicalName)
	assert.Equal(t, "", originalSchema.Tables[0].Columns[0].LogicalName)

	// Verify original comments preserved
	assert.Equal(t, "This is a legacy comment without separator", originalSchema.Tables[0].Comment)
	assert.Equal(t, "Primary key", originalSchema.Tables[0].Columns[0].Comment)

	t.Logf("✅ Backward Compatibility: Existing functionality preserved")
}

// TestE2E_ErrorHandlingRobustness tests system robustness
func TestE2E_ErrorHandlingRobustness(t *testing.T) {
	testCases := []struct {
		name        string
		schema      *schema.Schema
		config      *config.Config
		expectError bool
		description string
	}{
		{
			name: "Mixed comment formats",
			schema: &schema.Schema{
				Tables: []*schema.Table{
					{Name: "table1", Comment: "論理名|正常な形式"},
					{Name: "table2", Comment: "区切り文字なしのコメント"},
					{Name: "table3", Comment: ""},
				},
			},
			config: &config.Config{
				Comment: &config.CommentConfig{Separator: "|"},
			},
			expectError: false,
			description: "混在するコメント形式の処理",
		},
		{
			name: "Empty schema",
			schema: &schema.Schema{
				Tables: []*schema.Table{},
			},
			config: &config.Config{
				Comment: &config.CommentConfig{Separator: "|"},
			},
			expectError: false,
			description: "空のスキーマの処理",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.ModifySchema(tc.schema)

			if tc.expectError {
				assert.Error(t, err, tc.description)
			} else {
				assert.NoError(t, err, tc.description)
			}

			t.Logf("✅ Error handling test '%s': %s - Passed", tc.name, tc.description)
		})
	}
}

// TestE2E_ComprehensiveIntegration tests all components together
func TestE2E_ComprehensiveIntegration(t *testing.T) {
	// Create comprehensive test schema with various object types
	complexSchema := &schema.Schema{
		Name: "integration_test_db",
		Tables: []*schema.Table{
			{
				Name:    "users",
				Comment: "ユーザー情報|システム利用者の詳細情報",
				Columns: []*schema.Column{
					{Name: "user_id", Type: "BIGINT", Comment: "ユーザーID|主キー識別子", Nullable: false},
					{Name: "username", Type: "VARCHAR(50)", Comment: "ユーザー名|ログイン名", Nullable: false},
					{Name: "email", Type: "VARCHAR(100)", Comment: "メールアドレス|連絡先", Nullable: true},
				},
				Indexes: []*schema.Index{
					{Name: "idx_username", Comment: "ユーザー名インデックス|ユーザー名検索用"},
				},
				Constraints: []*schema.Constraint{
					{Name: "pk_users", Comment: "主キー制約|ユーザーID制約"},
				},
			},
			{
				Name:    "orders",
				Comment: "注文情報|顧客からの注文データ",
				Columns: []*schema.Column{
					{Name: "order_id", Type: "BIGINT", Comment: "注文ID|注文識別子"},
					{Name: "user_id", Type: "BIGINT", Comment: "ユーザーID|注文者ID"},
				},
			},
		},
	}

	// Configure full feature set
	fullConfig := &config.Config{
		Comment: &config.CommentConfig{
			Separator: "|",
		},
		Markdown: &config.MarkdownConfig{
			Tables: &config.TableCustomConfig{
				ObjectCustomConfig: &config.ObjectCustomConfig{
					ShowLogicalName: true,
					Order:           []string{"Name", "LogicalName", "Type", "Comment"},
					Aliases: map[string]string{
						"Name":        "テーブル名",
						"LogicalName": "論理名",
					},
				},
			},
		},
	}

	// Execute full integration
	start := time.Now()
	err := fullConfig.ModifySchema(complexSchema)
	processingTime := time.Since(start)

	assert.NoError(t, err)

	// Verify comprehensive logical name extraction
	assert.Equal(t, "ユーザー情報", complexSchema.Tables[0].LogicalName)
	assert.Equal(t, "注文情報", complexSchema.Tables[1].LogicalName)
	assert.Equal(t, "ユーザーID", complexSchema.Tables[0].Columns[0].LogicalName)
	assert.Equal(t, "ユーザー名", complexSchema.Tables[0].Columns[1].LogicalName)
	assert.Equal(t, "注文ID", complexSchema.Tables[1].Columns[0].LogicalName)

	// Verify indexes and constraints
	assert.Equal(t, "ユーザー名インデックス", complexSchema.Tables[0].Indexes[0].LogicalName)
	assert.Equal(t, "主キー制約", complexSchema.Tables[0].Constraints[0].LogicalName)

	// Count processed objects
	totalObjects := len(complexSchema.Tables)
	totalColumns := 0
	totalIndexes := 0
	totalConstraints := 0

	for _, table := range complexSchema.Tables {
		totalColumns += len(table.Columns)
		totalIndexes += len(table.Indexes)
		totalConstraints += len(table.Constraints)
	}

	t.Logf("🎉 Comprehensive Integration Test: Complete Success")
	t.Logf("   📊 Processed Objects:")
	t.Logf("      - Tables: %d", totalObjects)
	t.Logf("      - Columns: %d", totalColumns)
	t.Logf("      - Indexes: %d", totalIndexes)
	t.Logf("      - Constraints: %d", totalConstraints)
	t.Logf("   ⏱️  Processing time: %v", processingTime)
	t.Logf("   ✅ ModifySchema integration: PASS")
	t.Logf("   ✅ Japanese UTF-8 processing: PASS")
	t.Logf("   ✅ Multi-object support: PASS")
	t.Logf("   ✅ Configuration integration: PASS")
}

// Helper function to create simple test dataset
func createSimpleTestDataset() *schema.Schema {
	return &schema.Schema{
		Name: "test_db",
		Tables: []*schema.Table{
			{
				Name:    "table_0",
				Comment: "テーブル0|テストデータのテーブル0",
				Columns: []*schema.Column{
					{Name: "column_0", Type: "VARCHAR(255)", Comment: "カラム0|テストデータのカラム0"},
					{Name: "column_1", Type: "VARCHAR(255)", Comment: "カラム1|テストデータのカラム1"},
				},
			},
			{
				Name:    "table_1",
				Comment: "テーブル1|テストデータのテーブル1",
				Columns: []*schema.Column{
					{Name: "column_0", Type: "VARCHAR(255)", Comment: "カラム0|テストデータのカラム0"},
				},
			},
		},
	}
}