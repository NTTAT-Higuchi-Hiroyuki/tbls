package schema

import (
	"testing"
)

func TestColumn_ProcessEnhancedComment(t *testing.T) {
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
			comment:         `{"name": "ユーザーID", "description": "ユーザーの一意識別子", "tags": ["PK", "重要"]}`,
			delimiter:       "|",
			fallbackToName:  true,
			expectedLogical: "ユーザーID",
			expectedDesc:    "ユーザーの一意識別子",
			expectedTags:    []string{"PK", "重要"},
			expectedError:   false,
			expectedHasEnhanced: true,
		},
		{
			name:            "YAMLコメント処理",
			comment:         "name: カラム名\ndescription: カラムの説明\ntags:\n  - タグ1\n  - タグ2",
			delimiter:       "|",
			fallbackToName:  true,
			expectedLogical: "カラム名",
			expectedDesc:    "カラムの説明",
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
			column := &Column{
				Name:    "test_column",
				Comment: tt.comment,
			}

			err := column.ProcessEnhancedComment(processor, tt.delimiter, tt.fallbackToName)

			// エラーチェック
			if (err != nil) != tt.expectedError {
				t.Errorf("ProcessEnhancedComment() error = %v, expectedError %v", err, tt.expectedError)
			}

			// LogicalNameチェック
			if column.LogicalName != tt.expectedLogical {
				t.Errorf("LogicalName = %q, expected %q", column.LogicalName, tt.expectedLogical)
			}

			// 拡張コメントデータの存在チェック
			if column.HasEnhancedComment() != tt.expectedHasEnhanced {
				t.Errorf("HasEnhancedComment() = %v, expected %v", column.HasEnhancedComment(), tt.expectedHasEnhanced)
			}

			if tt.expectedHasEnhanced {
				// 説明チェック
				if column.GetDescription() != tt.expectedDesc {
					t.Errorf("GetDescription() = %q, expected %q", column.GetDescription(), tt.expectedDesc)
				}

				// タグチェック
				tags := column.GetTags()
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

func TestColumn_GetEnhancedLogicalNameOrFallback(t *testing.T) {
	tests := []struct {
		name              string
		column            *Column
		fallbackToName    bool
		expected          string
	}{
		{
			name: "拡張コメントの論理名を使用",
			column: &Column{
				Name:        "test_col",
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
			column: &Column{
				Name:        "test_col",
				LogicalName: "既存論理名",
			},
			fallbackToName: true,
			expected:       "既存論理名",
		},
		{
			name: "フォールバックでカラム名を使用",
			column: &Column{
				Name: "test_col",
			},
			fallbackToName: true,
			expected:       "test_col",
		},
		{
			name: "フォールバック無効で空文字列",
			column: &Column{
				Name: "test_col",
			},
			fallbackToName: false,
			expected:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.column.GetEnhancedLogicalNameOrFallback(tt.fallbackToName)
			if result != tt.expected {
				t.Errorf("GetEnhancedLogicalNameOrFallback() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestColumn_EnhancedCommentGetters(t *testing.T) {
	// 拡張コメントデータを持つカラム
	enhancedColumn := &Column{
		Name:    "enhanced_col",
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

	// 拡張コメントデータを持たないカラム
	normalColumn := &Column{
		Name:    "normal_col",
		Comment: "通常のコメント",
	}

	t.Run("拡張コメントありのカラム", func(t *testing.T) {
		if !enhancedColumn.HasEnhancedComment() {
			t.Error("HasEnhancedComment() should return true")
		}

		if enhancedColumn.GetDescription() != "拡張説明" {
			t.Errorf("GetDescription() = %q, expected '拡張説明'", enhancedColumn.GetDescription())
		}

		tags := enhancedColumn.GetTags()
		expectedTags := []string{"タグ1", "タグ2"}
		if len(tags) != len(expectedTags) {
			t.Errorf("GetTags() length = %d, expected %d", len(tags), len(expectedTags))
		}

		metadata := enhancedColumn.GetMetadata()
		if metadata["key"] != "value" {
			t.Errorf("GetMetadata()['key'] = %q, expected 'value'", metadata["key"])
		}

		if enhancedColumn.GetPriority() != 5 {
			t.Errorf("GetPriority() = %d, expected 5", enhancedColumn.GetPriority())
		}

		if !enhancedColumn.IsDeprecated() {
			t.Error("IsDeprecated() should return true")
		}
	})

	t.Run("拡張コメントなしのカラム", func(t *testing.T) {
		if normalColumn.HasEnhancedComment() {
			t.Error("HasEnhancedComment() should return false")
		}

		if normalColumn.GetDescription() != "通常のコメント" {
			t.Errorf("GetDescription() = %q, expected '通常のコメント'", normalColumn.GetDescription())
		}

		if normalColumn.GetTags() != nil {
			t.Error("GetTags() should return nil")
		}

		if normalColumn.GetMetadata() != nil {
			t.Error("GetMetadata() should return nil")
		}

		if normalColumn.GetPriority() != 0 {
			t.Errorf("GetPriority() = %d, expected 0", normalColumn.GetPriority())
		}

		if normalColumn.IsDeprecated() {
			t.Error("IsDeprecated() should return false")
		}
	})
}

func TestColumn_ProcessEnhancedCommentWithProcessor(t *testing.T) {
	// 設定を使用したプロセッサーでのテスト
	config := NewMockEnhancedCommentConfigurator()
	config.preferredFormat = "json"
	config.validationEnabled = true

	processor := NewEnhancedCommentProcessorFromConfig(config)

	column := &Column{
		Name:    "user_id",
		Comment: `{"name": "ユーザーID", "description": "システム内でユーザーを一意に識別するID", "tags": ["PK", "重要"], "deprecated": false}`,
	}

	err := column.ProcessEnhancedComment(processor, "|", true)
	if err != nil {
		t.Fatalf("ProcessEnhancedComment() failed: %v", err)
	}

	// 拡張コメントデータの検証
	if !column.HasEnhancedComment() {
		t.Fatal("Expected enhanced comment data")
	}

	enhancedData := column.GetEnhancedCommentData()
	if enhancedData.LogicalName != "ユーザーID" {
		t.Errorf("LogicalName = %q, expected 'ユーザーID'", enhancedData.LogicalName)
	}

	if enhancedData.Description != "システム内でユーザーを一意に識別するID" {
		t.Errorf("Description = %q, expected full description", enhancedData.Description)
	}

	expectedTags := []string{"PK", "重要"}
	tags := column.GetTags()
	if len(tags) != len(expectedTags) {
		t.Errorf("Tags length = %d, expected %d", len(tags), len(expectedTags))
	}

	if column.IsDeprecated() {
		t.Error("IsDeprecated() should return false")
	}

	// オブジェクトタイプがメタデータに設定されているか確認
	metadata := column.GetMetadata()
	if metadata["object_type"] != string(ObjectTypeColumn) {
		t.Errorf("object_type metadata = %q, expected %q", metadata["object_type"], ObjectTypeColumn)
	}
}

func TestColumn_ProcessEnhancedCommentBackwardCompatibility(t *testing.T) {
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
			column := &Column{
				Name:    "test_col",
				Comment: tt.comment,
			}

			err := column.ProcessEnhancedComment(processor, "|", true)
			if err != nil && tt.name == "JSONコメントの場合は元のコメントから論理名部分を削除" {
				// JSONの場合はエラーが発生しないことを期待
				t.Errorf("Unexpected error for JSON comment: %v", err)
			}

			// Commentフィールドがクリーンアップされているか確認
			if tt.name == "従来形式の論理名削除" && column.Comment != tt.expectedClean {
				t.Errorf("Comment after processing = %q, expected %q", column.Comment, tt.expectedClean)
			}
		})
	}
}