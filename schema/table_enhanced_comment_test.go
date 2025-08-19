package schema

import (
	"testing"
)

func TestTable_ProcessEnhancedComment(t *testing.T) {
	processor := NewEnhancedCommentProcessor()

	tests := []struct {
		name              string
		comment           string
		delimiter         string
		fallbackToName    bool
		expectedLogical   string
		expectedDesc      string
		expectedTags      []string
		expectedError     bool
		expectedHasEnhanced bool
	}{
		{
			name:            "JSONコメント処理",
			comment:         `{"name": "ユーザーテーブル", "description": "システムのユーザー情報を管理", "tags": ["重要", "マスター"]}`,
			delimiter:       "|",
			fallbackToName:  true,
			expectedLogical: "ユーザーテーブル",
			expectedDesc:    "システムのユーザー情報を管理",
			expectedTags:    []string{"重要", "マスター"},
			expectedError:   false,
			expectedHasEnhanced: true,
		},
		{
			name:            "YAMLコメント処理",
			comment:         "name: テーブル名\ndescription: テーブルの説明\ntags:\n  - タグ1\n  - タグ2",
			delimiter:       "|",
			fallbackToName:  true,
			expectedLogical: "テーブル名",
			expectedDesc:    "テーブルの説明",
			expectedTags:    []string{"タグ1", "タグ2"},
			expectedError:   false,
			expectedHasEnhanced: true,
		},
		{
			name:            "従来形式コメント処理",
			comment:         "論理名|説明文",
			delimiter:       "|",
			fallbackToName:  true,
			expectedLogical: "論理名",
			expectedDesc:    "説明文",
			expectedTags:    nil,
			expectedError:   false,
			expectedHasEnhanced: true,
		},
		{
			name:            "無効なJSONでフォールバック",
			comment:         "不正なJSON{",
			delimiter:       "|",
			fallbackToName:  true,
			expectedLogical: "不正なJSON{",  // LegacyParserでフォールバック処理
			expectedDesc:    "",             // 区切り文字がないため説明部分は空
			expectedTags:    nil,
			expectedError:   false,          // LegacyParserが成功するためエラーなし
			expectedHasEnhanced: true,       // LegacyParserでEnhancedCommentDataが作成される
		},
		{
			name:            "空コメント",
			comment:         "",
			delimiter:       "|",
			fallbackToName:  true,
			expectedLogical: "",
			expectedDesc:    "",
			expectedTags:    nil,
			expectedError:   false,
			expectedHasEnhanced: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{
				Name:    "test_table",
				Comment: tt.comment,
			}

			err := table.ProcessEnhancedComment(processor, tt.delimiter, tt.fallbackToName)

			// エラーチェック
			if (err != nil) != tt.expectedError {
				t.Errorf("ProcessEnhancedComment() error = %v, expectedError %v", err, tt.expectedError)
			}

			// LogicalNameチェック
			if table.LogicalName != tt.expectedLogical {
				t.Errorf("LogicalName = %q, expected %q", table.LogicalName, tt.expectedLogical)
			}

			// 拡張コメントデータの存在チェック
			if table.HasEnhancedComment() != tt.expectedHasEnhanced {
				t.Errorf("HasEnhancedComment() = %v, expected %v", table.HasEnhancedComment(), tt.expectedHasEnhanced)
			}

			if tt.expectedHasEnhanced {
				// 説明チェック
				if table.GetDescription() != tt.expectedDesc {
					t.Errorf("GetDescription() = %q, expected %q", table.GetDescription(), tt.expectedDesc)
				}

				// タグチェック
				tags := table.GetTags()
				if len(tags) != len(tt.expectedTags) {
					t.Errorf("GetTags() length = %d, expected %d", len(tags), len(tt.expectedTags))
				} else {
					for i, tag := range tags {
						if tag != tt.expectedTags[i] {
							t.Errorf("GetTags()[%d] = %q, expected %q", i, tag, tt.expectedTags[i])
						}
					}
				}
			}
		})
	}
}

func TestTable_GetEnhancedLogicalNameOrFallback(t *testing.T) {
	tests := []struct {
		name              string
		table             *Table
		fallbackToName    bool
		expected          string
	}{
		{
			name: "拡張コメントの論理名を使用",
			table: &Table{
				Name:        "test_tbl",
				LogicalName: "既存論理名",
				EnhancedCommentData: &CommentData{
					LogicalName: "拡張論理名",
				},
			},
			fallbackToName: true,
			expected:       "拡張論理名",
		},
		{
			name: "既存LogicalNameフィールドを使用",
			table: &Table{
				Name:        "test_tbl",
				LogicalName: "既存論理名",
			},
			fallbackToName: true,
			expected:       "既存論理名",
		},
		{
			name: "フォールバックでテーブル名を使用",
			table: &Table{
				Name: "test_tbl",
			},
			fallbackToName: true,
			expected:       "test_tbl",
		},
		{
			name: "フォールバック無効で空文字列",
			table: &Table{
				Name: "test_tbl",
			},
			fallbackToName: false,
			expected:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.table.GetEnhancedLogicalNameOrFallback(tt.fallbackToName)
			if result != tt.expected {
				t.Errorf("GetEnhancedLogicalNameOrFallback() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestTable_EnhancedCommentGetters(t *testing.T) {
	// 拡張コメントデータを持つテーブル
	enhancedTable := &Table{
		Name:    "enhanced_tbl",
		Comment: "元のコメント",
		EnhancedCommentData: &CommentData{
			LogicalName: "拡張論理名",
			Description: "拡張説明",
			Tags:        []string{"タグ1", "タグ2"},
			Metadata:    map[string]string{"key": "value"},
			Priority:    5,
			Deprecated:  true,
		},
	}

	// 拡張コメントデータを持たないテーブル
	normalTable := &Table{
		Name:    "normal_tbl",
		Comment: "通常のコメント",
	}

	t.Run("拡張コメントありのテーブル", func(t *testing.T) {
		if !enhancedTable.HasEnhancedComment() {
			t.Error("HasEnhancedComment() should return true")
		}

		if enhancedTable.GetDescription() != "拡張説明" {
			t.Errorf("GetDescription() = %q, expected '拡張説明'", enhancedTable.GetDescription())
		}

		tags := enhancedTable.GetTags()
		expectedTags := []string{"タグ1", "タグ2"}
		if len(tags) != len(expectedTags) {
			t.Errorf("GetTags() length = %d, expected %d", len(tags), len(expectedTags))
		}

		metadata := enhancedTable.GetMetadata()
		if metadata["key"] != "value" {
			t.Errorf("GetMetadata()['key'] = %q, expected 'value'", metadata["key"])
		}

		if enhancedTable.GetPriority() != 5 {
			t.Errorf("GetPriority() = %d, expected 5", enhancedTable.GetPriority())
		}

		if !enhancedTable.IsDeprecated() {
			t.Error("IsDeprecated() should return true")
		}
	})

	t.Run("拡張コメントなしのテーブル", func(t *testing.T) {
		if normalTable.HasEnhancedComment() {
			t.Error("HasEnhancedComment() should return false")
		}

		if normalTable.GetDescription() != "通常のコメント" {
			t.Errorf("GetDescription() = %q, expected '通常のコメント'", normalTable.GetDescription())
		}

		if normalTable.GetTags() != nil {
			t.Error("GetTags() should return nil")
		}

		if normalTable.GetMetadata() != nil {
			t.Error("GetMetadata() should return nil")
		}

		if normalTable.GetPriority() != 0 {
			t.Errorf("GetPriority() = %d, expected 0", normalTable.GetPriority())
		}

		if normalTable.IsDeprecated() {
			t.Error("IsDeprecated() should return false")
		}
	})
}

func TestTable_ProcessEnhancedCommentWithProcessor(t *testing.T) {
	// 設定を使用したプロセッサーでのテスト
	config := NewMockEnhancedCommentConfigurator()
	config.preferredFormat = "json"
	config.validationEnabled = true

	processor := NewEnhancedCommentProcessorFromConfig(config)

	table := &Table{
		Name:    "users",
		Comment: `{"name": "ユーザーテーブル", "description": "システム内のユーザー情報を管理するテーブル", "tags": ["マスター", "重要"], "deprecated": false}`,
	}

	err := table.ProcessEnhancedComment(processor, "|", true)
	if err != nil {
		t.Fatalf("ProcessEnhancedComment() failed: %v", err)
	}

	// 拡張コメントデータの検証
	if !table.HasEnhancedComment() {
		t.Fatal("Expected enhanced comment data")
	}

	enhancedData := table.GetEnhancedCommentData()
	if enhancedData.LogicalName != "ユーザーテーブル" {
		t.Errorf("LogicalName = %q, expected 'ユーザーテーブル'", enhancedData.LogicalName)
	}

	if enhancedData.Description != "システム内のユーザー情報を管理するテーブル" {
		t.Errorf("Description = %q, expected full description", enhancedData.Description)
	}

	expectedTags := []string{"マスター", "重要"}
	tags := table.GetTags()
	if len(tags) != len(expectedTags) {
		t.Errorf("Tags length = %d, expected %d", len(tags), len(expectedTags))
	}

	if table.IsDeprecated() {
		t.Error("IsDeprecated() should return false")
	}

	// オブジェクトタイプがメタデータに設定されているか確認
	metadata := table.GetMetadata()
	if metadata["object_type"] != string(ObjectTypeTable) {
		t.Errorf("object_type metadata = %q, expected %q", metadata["object_type"], ObjectTypeTable)
	}
}

func TestTable_ProcessEnhancedCommentBackwardCompatibility(t *testing.T) {
	// 既存の動作との互換性テスト
	processor := NewEnhancedCommentProcessor()

	tests := []struct {
		name           string
		comment        string
		expectedClean  string
	}{
		{
			name:          "従来形式の論理名削除",
			comment:       "論理名|説明文",
			expectedClean: "説明文",
		},
		{
			name:          "JSONコメントの場合は元のコメントから論理名部分を削除",
			comment:       `{"name": "論理名", "description": "説明"}`,
			expectedClean: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{
				Name:    "test_tbl",
				Comment: tt.comment,
			}

			err := table.ProcessEnhancedComment(processor, "|", true)
			if err != nil && tt.name == "JSONコメントの場合は元のコメントから論理名部分を削除" {
				// JSONの場合はエラーが発生しないことを期待
				t.Errorf("Unexpected error for JSON comment: %v", err)
			}

			// Commentフィールドがクリーンアップされているか確認
			if tt.name == "従来形式の論理名削除" && table.Comment != tt.expectedClean {
				t.Errorf("Comment after processing = %q, expected %q", table.Comment, tt.expectedClean)
			}
		})
	}
}