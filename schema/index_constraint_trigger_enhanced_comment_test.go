package schema

import (
	"testing"
)

func TestIndex_ProcessEnhancedComment(t *testing.T) {
	processor := NewEnhancedCommentProcessor()

	tests := []struct {
		name              string
		comment           string
		delimiter         string
		expectedDesc      string
		expectedTags      []string
		expectedError     bool
		expectedHasEnhanced bool
	}{
		{
			name:            "JSONコメント処理",
			comment:         `{"description": "ユーザーIDインデックス", "tags": ["性能", "重要"]}`,
			delimiter:       "|",
			expectedDesc:    "ユーザーIDインデックス",
			expectedTags:    []string{"性能", "重要"},
			expectedError:   false,
			expectedHasEnhanced: true,
		},
		{
			name:            "YAMLコメント処理",
			comment:         "description: インデックスの説明\ntags:\n  - タグ1\n  - タグ2",
			delimiter:       "|",
			expectedDesc:    "インデックスの説明",
			expectedTags:    []string{"タグ1", "タグ2"},
			expectedError:   false,
			expectedHasEnhanced: true,
		},
		{
			name:            "従来形式コメント処理",
			comment:         "説明文",
			delimiter:       "|",
			expectedDesc:    "",         // 区切り文字がないため、LegacyParserでは論理名として処理され、説明は空
			expectedTags:    nil,
			expectedError:   false,
			expectedHasEnhanced: true,
		},
		{
			name:            "空コメント",
			comment:         "",
			delimiter:       "|",
			expectedDesc:    "",
			expectedTags:    nil,
			expectedError:   false,
			expectedHasEnhanced: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := &Index{
				Name:    "test_index",
				Comment: tt.comment,
			}

			err := index.ProcessEnhancedComment(processor, tt.delimiter)

			// エラーチェック
			if (err != nil) != tt.expectedError {
				t.Errorf("ProcessEnhancedComment() error = %v, expectedError %v", err, tt.expectedError)
			}

			// 拡張コメントデータの存在チェック
			if index.HasEnhancedComment() != tt.expectedHasEnhanced {
				t.Errorf("HasEnhancedComment() = %v, expected %v", index.HasEnhancedComment(), tt.expectedHasEnhanced)
			}

			if tt.expectedHasEnhanced {
				// 説明チェック
				if index.GetDescription() != tt.expectedDesc {
					t.Errorf("GetDescription() = %q, expected %q", index.GetDescription(), tt.expectedDesc)
				}

				// タグチェック
				tags := index.GetTags()
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

func TestIndex_EnhancedCommentGetters(t *testing.T) {
	// 拡張コメントデータを持つインデックス
	enhancedIndex := &Index{
		Name:    "enhanced_idx",
		Comment: "元のコメント",
		EnhancedCommentData: &CommentData{
			Description: "拡張説明",
			Tags:        []string{"タグ1", "タグ2"},
			Metadata:    map[string]string{"key": "value"},
			Priority:    5,
			Deprecated:  true,
		},
	}

	// 拡張コメントデータを持たないインデックス
	normalIndex := &Index{
		Name:    "normal_idx",
		Comment: "通常のコメント",
	}

	t.Run("拡張コメントありのインデックス", func(t *testing.T) {
		if !enhancedIndex.HasEnhancedComment() {
			t.Error("HasEnhancedComment() should return true")
		}

		if enhancedIndex.GetDescription() != "拡張説明" {
			t.Errorf("GetDescription() = %q, expected '拡張説明'", enhancedIndex.GetDescription())
		}

		tags := enhancedIndex.GetTags()
		expectedTags := []string{"タグ1", "タグ2"}
		if len(tags) != len(expectedTags) {
			t.Errorf("GetTags() length = %d, expected %d", len(tags), len(expectedTags))
		}

		metadata := enhancedIndex.GetMetadata()
		if metadata["key"] != "value" {
			t.Errorf("GetMetadata()['key'] = %q, expected 'value'", metadata["key"])
		}

		if enhancedIndex.GetPriority() != 5 {
			t.Errorf("GetPriority() = %d, expected 5", enhancedIndex.GetPriority())
		}

		if !enhancedIndex.IsDeprecated() {
			t.Error("IsDeprecated() should return true")
		}
	})

	t.Run("拡張コメントなしのインデックス", func(t *testing.T) {
		if normalIndex.HasEnhancedComment() {
			t.Error("HasEnhancedComment() should return false")
		}

		if normalIndex.GetDescription() != "通常のコメント" {
			t.Errorf("GetDescription() = %q, expected '通常のコメント'", normalIndex.GetDescription())
		}

		if normalIndex.GetTags() != nil {
			t.Error("GetTags() should return nil")
		}

		if normalIndex.GetMetadata() != nil {
			t.Error("GetMetadata() should return nil")
		}

		if normalIndex.GetPriority() != 0 {
			t.Errorf("GetPriority() = %d, expected 0", normalIndex.GetPriority())
		}

		if normalIndex.IsDeprecated() {
			t.Error("IsDeprecated() should return false")
		}
	})
}

func TestConstraint_ProcessEnhancedComment(t *testing.T) {
	processor := NewEnhancedCommentProcessor()

	tests := []struct {
		name              string
		comment           string
		delimiter         string
		expectedDesc      string
		expectedTags      []string
		expectedError     bool
		expectedHasEnhanced bool
	}{
		{
			name:            "JSONコメント処理",
			comment:         `{"description": "外部キー制約", "tags": ["整合性", "重要"]}`,
			delimiter:       "|",
			expectedDesc:    "外部キー制約",
			expectedTags:    []string{"整合性", "重要"},
			expectedError:   false,
			expectedHasEnhanced: true,
		},
		{
			name:            "YAMLコメント処理",
			comment:         "description: 制約の説明\ntags:\n  - 制約タグ1\n  - 制約タグ2",
			delimiter:       "|",
			expectedDesc:    "制約の説明",
			expectedTags:    []string{"制約タグ1", "制約タグ2"},
			expectedError:   false,
			expectedHasEnhanced: true,
		},
		{
			name:            "従来形式コメント処理",
			comment:         "制約説明文",
			delimiter:       "|",
			expectedDesc:    "",         // 区切り文字がないため、LegacyParserでは論理名として処理され、説明は空
			expectedTags:    nil,
			expectedError:   false,
			expectedHasEnhanced: true,
		},
		{
			name:            "空コメント",
			comment:         "",
			delimiter:       "|",
			expectedDesc:    "",
			expectedTags:    nil,
			expectedError:   false,
			expectedHasEnhanced: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := &Constraint{
				Name:    "test_constraint",
				Comment: tt.comment,
			}

			err := constraint.ProcessEnhancedComment(processor, tt.delimiter)

			// エラーチェック
			if (err != nil) != tt.expectedError {
				t.Errorf("ProcessEnhancedComment() error = %v, expectedError %v", err, tt.expectedError)
			}

			// 拡張コメントデータの存在チェック
			if constraint.HasEnhancedComment() != tt.expectedHasEnhanced {
				t.Errorf("HasEnhancedComment() = %v, expected %v", constraint.HasEnhancedComment(), tt.expectedHasEnhanced)
			}

			if tt.expectedHasEnhanced {
				// 説明チェック
				if constraint.GetDescription() != tt.expectedDesc {
					t.Errorf("GetDescription() = %q, expected %q", constraint.GetDescription(), tt.expectedDesc)
				}

				// タグチェック
				tags := constraint.GetTags()
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

func TestTrigger_ProcessEnhancedComment(t *testing.T) {
	processor := NewEnhancedCommentProcessor()

	tests := []struct {
		name              string
		comment           string
		delimiter         string
		expectedDesc      string
		expectedTags      []string
		expectedError     bool
		expectedHasEnhanced bool
	}{
		{
			name:            "JSONコメント処理",
			comment:         `{"description": "更新時刻トリガー", "tags": ["自動化", "監査"]}`,
			delimiter:       "|",
			expectedDesc:    "更新時刻トリガー",
			expectedTags:    []string{"自動化", "監査"},
			expectedError:   false,
			expectedHasEnhanced: true,
		},
		{
			name:            "YAMLコメント処理",
			comment:         "description: トリガーの説明\ntags:\n  - トリガータグ1\n  - トリガータグ2",
			delimiter:       "|",
			expectedDesc:    "トリガーの説明",
			expectedTags:    []string{"トリガータグ1", "トリガータグ2"},
			expectedError:   false,
			expectedHasEnhanced: true,
		},
		{
			name:            "従来形式コメント処理",
			comment:         "トリガー説明文",
			delimiter:       "|",
			expectedDesc:    "",         // 区切り文字がないため、LegacyParserでは論理名として処理され、説明は空
			expectedTags:    nil,
			expectedError:   false,
			expectedHasEnhanced: true,
		},
		{
			name:            "空コメント",
			comment:         "",
			delimiter:       "|",
			expectedDesc:    "",
			expectedTags:    nil,
			expectedError:   false,
			expectedHasEnhanced: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trigger := &Trigger{
				Name:    "test_trigger",
				Comment: tt.comment,
			}

			err := trigger.ProcessEnhancedComment(processor, tt.delimiter)

			// エラーチェック
			if (err != nil) != tt.expectedError {
				t.Errorf("ProcessEnhancedComment() error = %v, expectedError %v", err, tt.expectedError)
			}

			// 拡張コメントデータの存在チェック
			if trigger.HasEnhancedComment() != tt.expectedHasEnhanced {
				t.Errorf("HasEnhancedComment() = %v, expected %v", trigger.HasEnhancedComment(), tt.expectedHasEnhanced)
			}

			if tt.expectedHasEnhanced {
				// 説明チェック
				if trigger.GetDescription() != tt.expectedDesc {
					t.Errorf("GetDescription() = %q, expected %q", trigger.GetDescription(), tt.expectedDesc)
				}

				// タグチェック
				tags := trigger.GetTags()
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

func TestConstraint_EnhancedCommentGetters(t *testing.T) {
	// 拡張コメントデータを持つ制約
	enhancedConstraint := &Constraint{
		Name:    "enhanced_constraint",
		Comment: "元のコメント",
		EnhancedCommentData: &CommentData{
			Description: "拡張説明",
			Tags:        []string{"タグ1", "タグ2"},
			Metadata:    map[string]string{"key": "value"},
			Priority:    3,
			Deprecated:  false,
		},
	}

	// 拡張コメントデータを持たない制約
	normalConstraint := &Constraint{
		Name:    "normal_constraint",
		Comment: "通常のコメント",
	}

	t.Run("拡張コメントありの制約", func(t *testing.T) {
		if !enhancedConstraint.HasEnhancedComment() {
			t.Error("HasEnhancedComment() should return true")
		}

		if enhancedConstraint.GetDescription() != "拡張説明" {
			t.Errorf("GetDescription() = %q, expected '拡張説明'", enhancedConstraint.GetDescription())
		}

		if enhancedConstraint.GetPriority() != 3 {
			t.Errorf("GetPriority() = %d, expected 3", enhancedConstraint.GetPriority())
		}

		if enhancedConstraint.IsDeprecated() {
			t.Error("IsDeprecated() should return false")
		}
	})

	t.Run("拡張コメントなしの制約", func(t *testing.T) {
		if normalConstraint.HasEnhancedComment() {
			t.Error("HasEnhancedComment() should return false")
		}

		if normalConstraint.GetDescription() != "通常のコメント" {
			t.Errorf("GetDescription() = %q, expected '通常のコメント'", normalConstraint.GetDescription())
		}
	})
}

func TestTrigger_EnhancedCommentGetters(t *testing.T) {
	// 拡張コメントデータを持つトリガー
	enhancedTrigger := &Trigger{
		Name:    "enhanced_trigger",
		Comment: "元のコメント",
		EnhancedCommentData: &CommentData{
			Description: "拡張説明",
			Tags:        []string{"タグ1", "タグ2"},
			Metadata:    map[string]string{"key": "value"},
			Priority:    1,
			Deprecated:  true,
		},
	}

	// 拡張コメントデータを持たないトリガー
	normalTrigger := &Trigger{
		Name:    "normal_trigger",
		Comment: "通常のコメント",
	}

	t.Run("拡張コメントありのトリガー", func(t *testing.T) {
		if !enhancedTrigger.HasEnhancedComment() {
			t.Error("HasEnhancedComment() should return true")
		}

		if enhancedTrigger.GetDescription() != "拡張説明" {
			t.Errorf("GetDescription() = %q, expected '拡張説明'", enhancedTrigger.GetDescription())
		}

		if enhancedTrigger.GetPriority() != 1 {
			t.Errorf("GetPriority() = %d, expected 1", enhancedTrigger.GetPriority())
		}

		if !enhancedTrigger.IsDeprecated() {
			t.Error("IsDeprecated() should return true")
		}
	})

	t.Run("拡張コメントなしのトリガー", func(t *testing.T) {
		if normalTrigger.HasEnhancedComment() {
			t.Error("HasEnhancedComment() should return false")
		}

		if normalTrigger.GetDescription() != "通常のコメント" {
			t.Errorf("GetDescription() = %q, expected '通常のコメント'", normalTrigger.GetDescription())
		}
	})
}