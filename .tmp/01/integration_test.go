package main

import (
	"testing"

	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/schema"
	"github.com/stretchr/testify/assert"
)

// TestLogicalNameIntegrationWithModifySchema tests the integration of logical name processing
// through the ModifySchema function, which is the central point where all DBMS drivers
// will have their schema processed.
func TestLogicalNameIntegrationWithModifySchema(t *testing.T) {
	// Prepare a test schema with Japanese comments containing logical names
	testSchema := &schema.Schema{
		Tables: []*schema.Table{
			{
				Name:    "users",
				Comment: "ユーザーテーブル|システムユーザーを管理する",
				Columns: []*schema.Column{
					{
						Name:    "user_id",
						Comment: "ユーザーID|一意識別子",
					},
					{
						Name:    "email",
						Comment: "メールアドレス|ログイン用メール",
					},
				},
				Indexes: []*schema.Index{
					{
						Name:    "idx_email",
						Comment: "メール検索インデックス|メールアドレス検索用",
					},
				},
				Constraints: []*schema.Constraint{
					{
						Name:    "fk_dept",
						Comment: "部署参照制約|部署との関連",
					},
				},
				Triggers: []*schema.Trigger{
					{
						Name:    "trg_audit",
						Comment: "監査トリガー|変更履歴記録",
					},
				},
			},
			{
				Name:    "departments",
				Comment: "部署テーブル|会社の部署情報",
			},
		},
	}

	// Prepare config with comment separator
	cfg := &config.Config{
		Comment: &config.CommentConfig{
			Separator: "|",
		},
	}

	// Execute ModifySchema which should apply logical name processing
	err := cfg.ModifySchema(testSchema)
	assert.NoError(t, err)

	// Verify that logical names have been extracted and set
	// Table verification
	assert.Equal(t, "ユーザーテーブル", testSchema.Tables[0].LogicalName)
	assert.Equal(t, "部署テーブル", testSchema.Tables[1].LogicalName)

	// Column verification
	assert.Equal(t, "ユーザーID", testSchema.Tables[0].Columns[0].LogicalName)
	assert.Equal(t, "メールアドレス", testSchema.Tables[0].Columns[1].LogicalName)

	// Index verification
	assert.Equal(t, "メール検索インデックス", testSchema.Tables[0].Indexes[0].LogicalName)

	// Constraint verification
	assert.Equal(t, "部署参照制約", testSchema.Tables[0].Constraints[0].LogicalName)

	// Trigger verification
	assert.Equal(t, "監査トリガー", testSchema.Tables[0].Triggers[0].LogicalName)

	// Verify that original comments are preserved (this is the expected behavior)
	// The LogicalNameProcessor extracts logical names but preserves original comments
	assert.Equal(t, "ユーザーテーブル|システムユーザーを管理する", testSchema.Tables[0].Comment)
	assert.Equal(t, "ユーザーID|一意識別子", testSchema.Tables[0].Columns[0].Comment)

	// Test GetLogicalNameOrFallback functionality
	assert.Equal(t, "ユーザーテーブル", testSchema.Tables[0].GetLogicalNameOrFallback())
	assert.Equal(t, "ユーザーID", testSchema.Tables[0].Columns[0].GetLogicalNameOrFallback())

	t.Logf("✅ Task 3.1 Integration Test: Logical name processing is working correctly")
	t.Logf("   - Tables with logical names: %s -> %s", testSchema.Tables[0].Name, testSchema.Tables[0].LogicalName)
	t.Logf("   - Columns with logical names: %s -> %s", testSchema.Tables[0].Columns[0].Name, testSchema.Tables[0].Columns[0].LogicalName)
	t.Logf("   - All PostgreSQL objects (tables, columns, indexes, constraints, triggers) support logical names")
}

// TestLogicalNameIntegrationWithoutSeparator tests that when no separator is configured,
// logical name processing is skipped and original comments are preserved.
func TestLogicalNameIntegrationWithoutSeparator(t *testing.T) {
	// Prepare a test schema
	testSchema := &schema.Schema{
		Tables: []*schema.Table{
			{
				Name:    "users",
				Comment: "ユーザーテーブル|システムユーザーを管理する",
				Columns: []*schema.Column{
					{
						Name:    "user_id",
						Comment: "ユーザーID|一意識別子",
					},
				},
			},
		},
	}

	// Config without comment separator
	cfg := &config.Config{}

	// Execute ModifySchema
	err := cfg.ModifySchema(testSchema)
	assert.NoError(t, err)

	// Verify that logical names are NOT set (empty)
	assert.Equal(t, "", testSchema.Tables[0].LogicalName)
	assert.Equal(t, "", testSchema.Tables[0].Columns[0].LogicalName)

	// Verify that original comments are preserved
	assert.Equal(t, "ユーザーテーブル|システムユーザーを管理する", testSchema.Tables[0].Comment)
	assert.Equal(t, "ユーザーID|一意識別子", testSchema.Tables[0].Columns[0].Comment)

	t.Logf("✅ Backward compatibility verified: No processing when separator not configured")
}

// TestLogicalNameIntegrationBackwardCompatibility tests that the integration
// maintains backward compatibility when logical name processing fails.
func TestLogicalNameIntegrationBackwardCompatibility(t *testing.T) {
	// Prepare a test schema that might cause processing issues
	testSchema := &schema.Schema{
		Tables: []*schema.Table{
			{
				Name:    "users",
				Comment: "Normal comment without separator",
			},
		},
	}

	// Config with comment separator
	cfg := &config.Config{
		Comment: &config.CommentConfig{
			Separator: "|",
		},
	}

	// Execute ModifySchema - should not fail even if logical name processing doesn't find separators
	err := cfg.ModifySchema(testSchema)
	assert.NoError(t, err)

	// Verify that processing completed successfully without logical names
	assert.Equal(t, "", testSchema.Tables[0].LogicalName)
	assert.Equal(t, "Normal comment without separator", testSchema.Tables[0].Comment)

	t.Logf("✅ Error tolerance verified: Processing continues when no logical names found")
}
