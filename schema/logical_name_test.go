package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogicalNameProcessor(t *testing.T) {
	processor := NewLogicalNameProcessor("|")

	assert.NotNil(t, processor)
	assert.Equal(t, "|", processor.separator)
	assert.NotNil(t, processor.parser)
}

func TestLogicalNameProcessor_ProcessSchema(t *testing.T) {
	processor := NewLogicalNameProcessor("|")

	t.Run("正常なスキーマ処理", func(t *testing.T) {
		schema := &Schema{
			Tables: []*Table{
				{
					Name:    "users",
					Comment: "ユーザーテーブル|システムユーザーを管理する",
					Columns: []*Column{
						{
							Name:    "user_id",
							Comment: "ユーザーID|一意識別子",
						},
					},
				},
			},
		}

		err := processor.ProcessSchema(schema)
		assert.NoError(t, err)

		// テーブルの論理名が設定されていることを確認
		assert.Equal(t, "ユーザーテーブル", schema.Tables[0].LogicalName)

		// カラムの論理名が設定されていることを確認
		assert.Equal(t, "ユーザーID", schema.Tables[0].Columns[0].LogicalName)
	})

	t.Run("nilスキーマのエラー処理", func(t *testing.T) {
		err := processor.ProcessSchema(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "schema is nil")
	})
}

func TestLogicalNameProcessor_ProcessTable(t *testing.T) {
	processor := NewLogicalNameProcessor("|")

	t.Run("正常なテーブル処理", func(t *testing.T) {
		table := &Table{
			Name:    "users",
			Comment: "ユーザーテーブル|システムユーザーを管理する",
			Columns: []*Column{
				{
					Name:    "user_id",
					Comment: "ユーザーID|一意識別子",
				},
			},
			Indexes: []*Index{
				{
					Name:    "idx_email",
					Comment: "メール検索インデックス|メールアドレス検索用",
				},
			},
			Constraints: []*Constraint{
				{
					Name:    "fk_dept",
					Comment: "部署参照制約|部署との関連",
				},
			},
			Triggers: []*Trigger{
				{
					Name:    "trg_audit",
					Comment: "監査トリガー|変更履歴記録",
				},
			},
		}

		err := processor.ProcessTable(table)
		assert.NoError(t, err)

		// 各オブジェクトの論理名が設定されていることを確認
		assert.Equal(t, "ユーザーテーブル", table.LogicalName)
		assert.Equal(t, "ユーザーID", table.Columns[0].LogicalName)
		assert.Equal(t, "メール検索インデックス", table.Indexes[0].LogicalName)
		assert.Equal(t, "部署参照制約", table.Constraints[0].LogicalName)
		assert.Equal(t, "監査トリガー", table.Triggers[0].LogicalName)
	})

	t.Run("nilテーブルのエラー処理", func(t *testing.T) {
		err := processor.ProcessTable(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "table is nil")
	})
}

func TestNewLogicalNameDisplayController(t *testing.T) {
	controller := NewLogicalNameDisplayController(true, true)

	assert.NotNil(t, controller)
	assert.True(t, controller.showLogicalName)
	assert.True(t, controller.fallbackToPhysical)
	assert.Equal(t, "%s", controller.format)
}

func TestLogicalNameDisplayController_SetFormat(t *testing.T) {
	controller := NewLogicalNameDisplayController(true, true)

	controller.SetFormat("%l (%p)")
	assert.Equal(t, "%l (%p)", controller.format)

	// 空のフォーマットは設定されない
	controller.SetFormat("")
	assert.Equal(t, "%l (%p)", controller.format)
}

func TestLogicalNameDisplayController_GetDisplayName(t *testing.T) {
	tests := []struct {
		name               string
		showLogicalName    bool
		fallbackToPhysical bool
		format             string
		logicalName        string
		physicalName       string
		expected           string
	}{
		{
			name:               "論理名表示有効・論理名あり",
			showLogicalName:    true,
			fallbackToPhysical: true,
			format:             "%s",
			logicalName:        "ユーザーID",
			physicalName:       "user_id",
			expected:           "ユーザーID",
		},
		{
			name:               "論理名表示無効",
			showLogicalName:    false,
			fallbackToPhysical: true,
			format:             "%s",
			logicalName:        "ユーザーID",
			physicalName:       "user_id",
			expected:           "user_id",
		},
		{
			name:               "論理名なし・フォールバック有効",
			showLogicalName:    true,
			fallbackToPhysical: true,
			format:             "%s",
			logicalName:        "",
			physicalName:       "user_id",
			expected:           "user_id",
		},
		{
			name:               "論理名なし・フォールバック無効",
			showLogicalName:    true,
			fallbackToPhysical: false,
			format:             "%s",
			logicalName:        "",
			physicalName:       "user_id",
			expected:           "user_id",
		},
		{
			name:               "フォーマット：論理名のみ",
			showLogicalName:    true,
			fallbackToPhysical: true,
			format:             "%l",
			logicalName:        "ユーザーID",
			physicalName:       "user_id",
			expected:           "ユーザーID",
		},
		{
			name:               "フォーマット：物理名のみ",
			showLogicalName:    true,
			fallbackToPhysical: true,
			format:             "%p",
			logicalName:        "ユーザーID",
			physicalName:       "user_id",
			expected:           "user_id",
		},
		{
			name:               "フォーマット：論理名(物理名)",
			showLogicalName:    true,
			fallbackToPhysical: true,
			format:             "%l (%p)",
			logicalName:        "ユーザーID",
			physicalName:       "user_id",
			expected:           "ユーザーID (user_id)",
		},
		{
			name:               "フォーマット：物理名(論理名)",
			showLogicalName:    true,
			fallbackToPhysical: true,
			format:             "%p (%l)",
			logicalName:        "ユーザーID",
			physicalName:       "user_id",
			expected:           "user_id (ユーザーID)",
		},
		{
			name:               "論理名なし・フォーマット：論理名(物理名)",
			showLogicalName:    true,
			fallbackToPhysical: true,
			format:             "%l (%p)",
			logicalName:        "",
			physicalName:       "user_id",
			expected:           "user_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewLogicalNameDisplayController(tt.showLogicalName, tt.fallbackToPhysical)
			controller.SetFormat(tt.format)

			result := controller.GetDisplayName(tt.logicalName, tt.physicalName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLogicalNameDisplayController_ObjectMethods(t *testing.T) {
	controller := NewLogicalNameDisplayController(true, true)

	t.Run("GetTableDisplayName", func(t *testing.T) {
		table := &Table{
			Name:        "users",
			LogicalName: "ユーザーテーブル",
		}

		result := controller.GetTableDisplayName(table)
		assert.Equal(t, "ユーザーテーブル", result)

		// nilテーブル
		result = controller.GetTableDisplayName(nil)
		assert.Equal(t, "", result)
	})

	t.Run("GetColumnDisplayName", func(t *testing.T) {
		column := &Column{
			Name:        "user_id",
			LogicalName: "ユーザーID",
		}

		result := controller.GetColumnDisplayName(column)
		assert.Equal(t, "ユーザーID", result)

		// nilカラム
		result = controller.GetColumnDisplayName(nil)
		assert.Equal(t, "", result)
	})

	t.Run("GetIndexDisplayName", func(t *testing.T) {
		index := &Index{
			Name:        "idx_email",
			LogicalName: "メール検索インデックス",
		}

		result := controller.GetIndexDisplayName(index)
		assert.Equal(t, "メール検索インデックス", result)

		// nilインデックス
		result = controller.GetIndexDisplayName(nil)
		assert.Equal(t, "", result)
	})

	t.Run("GetConstraintDisplayName", func(t *testing.T) {
		constraint := &Constraint{
			Name:        "fk_dept",
			LogicalName: "部署参照制約",
		}

		result := controller.GetConstraintDisplayName(constraint)
		assert.Equal(t, "部署参照制約", result)

		// nil制約
		result = controller.GetConstraintDisplayName(nil)
		assert.Equal(t, "", result)
	})

	t.Run("GetTriggerDisplayName", func(t *testing.T) {
		trigger := &Trigger{
			Name:        "trg_audit",
			LogicalName: "監査トリガー",
		}

		result := controller.GetTriggerDisplayName(trigger)
		assert.Equal(t, "監査トリガー", result)

		// nilトリガー
		result = controller.GetTriggerDisplayName(nil)
		assert.Equal(t, "", result)
	})
}

func TestNewLogicalNameValidator(t *testing.T) {
	validator := NewLogicalNameValidator()

	assert.NotNil(t, validator)
	assert.Equal(t, 255, validator.maxLength)
	assert.Equal(t, "", validator.allowedChars)
	assert.Empty(t, validator.forbiddenWords)
	assert.True(t, validator.caseSensitive)
}

func TestLogicalNameValidator_SetMethods(t *testing.T) {
	validator := NewLogicalNameValidator()

	// SetMaxLength
	validator.SetMaxLength(100)
	assert.Equal(t, 100, validator.maxLength)

	// 無効な値は設定されない
	validator.SetMaxLength(-1)
	assert.Equal(t, 100, validator.maxLength)

	// SetAllowedChars
	validator.SetAllowedChars("あいうえおアイウエオ")
	assert.Equal(t, "あいうえおアイウエオ", validator.allowedChars)

	// AddForbiddenWord
	validator.AddForbiddenWord("禁止")
	validator.AddForbiddenWord("ダメ")
	assert.Contains(t, validator.forbiddenWords, "禁止")
	assert.Contains(t, validator.forbiddenWords, "ダメ")

	// 空文字は追加されない
	validator.AddForbiddenWord("")
	assert.Len(t, validator.forbiddenWords, 2)

	// SetCaseSensitive
	validator.SetCaseSensitive(false)
	assert.False(t, validator.caseSensitive)
}

func TestLogicalNameValidator_Validate(t *testing.T) {
	tests := []struct {
		name           string
		logicalName    string
		maxLength      int
		allowedChars   string
		forbiddenWords []string
		caseSensitive  bool
		expectError    bool
		errorContains  string
	}{
		{
			name:        "空の論理名は有効",
			logicalName: "",
			expectError: false,
		},
		{
			name:        "通常の論理名",
			logicalName: "ユーザーID",
			expectError: false,
		},
		{
			name:          "最大長超過",
			logicalName:   "これはとても長い論理名です",
			maxLength:     10,
			expectError:   true,
			errorContains: "too long",
		},
		{
			name:          "禁止文字",
			logicalName:   "ユーザーID@",
			allowedChars:  "あいうえおアイウエオユーザーID",
			expectError:   true,
			errorContains: "forbidden character",
		},
		{
			name:           "禁止単語（大文字小文字区別あり）",
			logicalName:    "禁止ユーザーID",
			forbiddenWords: []string{"禁止"},
			caseSensitive:  true,
			expectError:    true,
			errorContains:  "forbidden word",
		},
		{
			name:           "禁止単語（大文字小文字区別なし）",
			logicalName:    "ADMIN_USER",
			forbiddenWords: []string{"admin"},
			caseSensitive:  false,
			expectError:    true,
			errorContains:  "forbidden word",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewLogicalNameValidator()

			if tt.maxLength > 0 {
				validator.SetMaxLength(tt.maxLength)
			}
			if tt.allowedChars != "" {
				validator.SetAllowedChars(tt.allowedChars)
			}
			for _, word := range tt.forbiddenWords {
				validator.AddForbiddenWord(word)
			}
			validator.SetCaseSensitive(tt.caseSensitive)

			err := validator.Validate(tt.logicalName)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewLogicalNameHelper(t *testing.T) {
	helper := NewLogicalNameHelper()
	assert.NotNil(t, helper)
}

func TestLogicalNameHelper_HasLogicalName(t *testing.T) {
	helper := NewLogicalNameHelper()

	tests := []struct {
		name        string
		logicalName string
		expected    bool
	}{
		{
			name:        "論理名あり",
			logicalName: "ユーザーID",
			expected:    true,
		},
		{
			name:        "空文字",
			logicalName: "",
			expected:    false,
		},
		{
			name:        "空白のみ",
			logicalName: "   ",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.HasLogicalName(tt.logicalName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLogicalNameHelper_IsLogicalNameAvailable(t *testing.T) {
	helper := NewLogicalNameHelper()

	tests := []struct {
		name      string
		comment   string
		separator string
		expected  bool
	}{
		{
			name:      "論理名取得可能",
			comment:   "ユーザーID|一意識別子",
			separator: "|",
			expected:  true,
		},
		{
			name:      "論理名取得不可能",
			comment:   "通常のコメント",
			separator: "|",
			expected:  false,
		},
		{
			name:      "空コメント",
			comment:   "",
			separator: "|",
			expected:  false,
		},
		{
			name:      "空区切り文字",
			comment:   "ユーザーID|一意識別子",
			separator: "",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.IsLogicalNameAvailable(tt.comment, tt.separator)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLogicalNameHelper_CleanLogicalName(t *testing.T) {
	helper := NewLogicalNameHelper()

	tests := []struct {
		name        string
		logicalName string
		expected    string
	}{
		{
			name:        "前後の空白削除",
			logicalName: "  ユーザーID  ",
			expected:    "ユーザーID",
		},
		{
			name:        "連続する空白の正規化",
			logicalName: "ユーザー  ID",
			expected:    "ユーザー ID",
		},
		{
			name:        "複雑な空白パターン",
			logicalName: "  ユーザー   ID  ",
			expected:    "ユーザー ID",
		},
		{
			name:        "空白なし",
			logicalName: "ユーザーID",
			expected:    "ユーザーID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.CleanLogicalName(tt.logicalName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLogicalNameHelper_NormalizeLogicalName(t *testing.T) {
	helper := NewLogicalNameHelper()

	tests := []struct {
		name        string
		logicalName string
		expected    string
	}{
		{
			name:        "通常の正規化",
			logicalName: "  ユーザー  ID  ",
			expected:    "ユーザー ID",
		},
		{
			name:        "制御文字の削除",
			logicalName: "ユーザー\tID\n",
			expected:    "ユーザー ID",
		},
		{
			name:        "DEL文字の削除",
			logicalName: "ユーザー\x7fID",
			expected:    "ユーザーID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.NormalizeLogicalName(tt.logicalName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLogicalNameHelper_CompareLogicalNames(t *testing.T) {
	helper := NewLogicalNameHelper()

	tests := []struct {
		name          string
		name1         string
		name2         string
		caseSensitive bool
		expected      bool
	}{
		{
			name:          "大文字小文字区別あり・同じ",
			name1:         "UserID",
			name2:         "UserID",
			caseSensitive: true,
			expected:      true,
		},
		{
			name:          "大文字小文字区別あり・異なる",
			name1:         "UserID",
			name2:         "userid",
			caseSensitive: true,
			expected:      false,
		},
		{
			name:          "大文字小文字区別なし・同じ",
			name1:         "UserID",
			name2:         "userid",
			caseSensitive: false,
			expected:      true,
		},
		{
			name:          "日本語比較",
			name1:         "ユーザーID",
			name2:         "ユーザーID",
			caseSensitive: true,
			expected:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.CompareLogicalNames(tt.name1, tt.name2, tt.caseSensitive)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLogicalNameHelper_GetPreferredName(t *testing.T) {
	helper := NewLogicalNameHelper()

	tests := []struct {
		name          string
		logicalName   string
		physicalName  string
		preferLogical bool
		expected      string
	}{
		{
			name:          "論理名優先・論理名あり",
			logicalName:   "ユーザーID",
			physicalName:  "user_id",
			preferLogical: true,
			expected:      "ユーザーID",
		},
		{
			name:          "論理名優先・論理名なし",
			logicalName:   "",
			physicalName:  "user_id",
			preferLogical: true,
			expected:      "user_id",
		},
		{
			name:          "物理名優先",
			logicalName:   "ユーザーID",
			physicalName:  "user_id",
			preferLogical: false,
			expected:      "user_id",
		},
		{
			name:          "両方空の場合",
			logicalName:   "",
			physicalName:  "",
			preferLogical: true,
			expected:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.GetPreferredName(tt.logicalName, tt.physicalName, tt.preferLogical)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateLogicalNameStatistics(t *testing.T) {
	t.Run("正常な統計計算", func(t *testing.T) {
		schema := &Schema{
			Tables: []*Table{
				{
					Name:        "users",
					LogicalName: "ユーザーテーブル",
					Columns: []*Column{
						{
							Name:        "user_id",
							LogicalName: "ユーザーID",
						},
						{
							Name:        "email",
							LogicalName: "", // 論理名なし
						},
					},
					Indexes: []*Index{
						{
							Name:        "idx_email",
							LogicalName: "メール検索インデックス",
						},
					},
					Constraints: []*Constraint{
						{
							Name:        "pk_user",
							LogicalName: "", // 論理名なし
						},
					},
					Triggers: []*Trigger{
						{
							Name:        "trg_audit",
							LogicalName: "監査トリガー",
						},
					},
				},
				{
					Name:        "departments",
					LogicalName: "", // 論理名なし
				},
			},
		}

		stats := CalculateLogicalNameStatistics(schema)

		assert.Equal(t, 2, stats.TotalTables)
		assert.Equal(t, 1, stats.TablesWithLogical)
		assert.Equal(t, 2, stats.TotalColumns)
		assert.Equal(t, 1, stats.ColumnsWithLogical)
		assert.Equal(t, 1, stats.TotalIndexes)
		assert.Equal(t, 1, stats.IndexesWithLogical)
		assert.Equal(t, 1, stats.TotalConstraints)
		assert.Equal(t, 0, stats.ConstraintsWithLogical)
		assert.Equal(t, 1, stats.TotalTriggers)
		assert.Equal(t, 1, stats.TriggersWithLogical)
	})

	t.Run("nilスキーマ", func(t *testing.T) {
		stats := CalculateLogicalNameStatistics(nil)

		assert.Equal(t, 0, stats.TotalTables)
		assert.Equal(t, 0, stats.TablesWithLogical)
	})

	t.Run("空のスキーマ", func(t *testing.T) {
		schema := &Schema{Tables: []*Table{}}
		stats := CalculateLogicalNameStatistics(schema)

		assert.Equal(t, 0, stats.TotalTables)
		assert.Equal(t, 0, stats.TablesWithLogical)
	})
}

func TestLogicalNameStatistics_GetCoveragePercentage(t *testing.T) {
	stats := &LogicalNameStatistics{
		TotalTables:            4,
		TablesWithLogical:      2,
		TotalColumns:           10,
		ColumnsWithLogical:     8,
		TotalIndexes:           2,
		IndexesWithLogical:     1,
		TotalConstraints:       0, // 0の場合はカバレッジに含まれない
		ConstraintsWithLogical: 0,
		TotalTriggers:          2,
		TriggersWithLogical:    2,
	}

	coverage := stats.GetCoveragePercentage()

	assert.Equal(t, 50.0, coverage["tables"])
	assert.Equal(t, 80.0, coverage["columns"])
	assert.Equal(t, 50.0, coverage["indexes"])
	assert.Equal(t, 100.0, coverage["triggers"])
	assert.NotContains(t, coverage, "constraints") // 総数が0のため含まれない
}

func TestLogicalNameStatistics_GetTotalCoveragePercentage(t *testing.T) {
	t.Run("正常なカバレッジ計算", func(t *testing.T) {
		stats := &LogicalNameStatistics{
			TotalTables:            2,
			TablesWithLogical:      1,
			TotalColumns:           4,
			ColumnsWithLogical:     3,
			TotalIndexes:           2,
			IndexesWithLogical:     1,
			TotalConstraints:       2,
			ConstraintsWithLogical: 0,
			TotalTriggers:          2,
			TriggersWithLogical:    1,
		}

		// 総数: 2+4+2+2+2 = 12
		// 論理名あり: 1+3+1+0+1 = 6
		// カバレッジ: 6/12 * 100 = 50%
		coverage := stats.GetTotalCoveragePercentage()
		assert.Equal(t, 50.0, coverage)
	})

	t.Run("総数が0の場合", func(t *testing.T) {
		stats := &LogicalNameStatistics{}
		coverage := stats.GetTotalCoveragePercentage()
		assert.Equal(t, 0.0, coverage)
	})
}

func TestLogicalNameIntegration(t *testing.T) {
	t.Run("統合シナリオテスト", func(t *testing.T) {
		// 1. スキーマの準備
		schema := &Schema{
			Tables: []*Table{
				{
					Name:    "users",
					Comment: "ユーザーテーブル|システムユーザーを管理する",
					Columns: []*Column{
						{
							Name:    "user_id",
							Comment: "ユーザーID|一意識別子",
						},
						{
							Name:    "name",
							Comment: "氏名|ユーザーの氏名",
						},
					},
				},
			},
		}

		// 2. 論理名処理の実行
		processor := NewLogicalNameProcessor("|")
		err := processor.ProcessSchema(schema)
		require.NoError(t, err)

		// 3. 表示制御のテスト
		controller := NewLogicalNameDisplayController(true, true)
		controller.SetFormat("%l (%p)")

		tableDisplayName := controller.GetTableDisplayName(schema.Tables[0])
		assert.Equal(t, "ユーザーテーブル (users)", tableDisplayName)

		columnDisplayName := controller.GetColumnDisplayName(schema.Tables[0].Columns[0])
		assert.Equal(t, "ユーザーID (user_id)", columnDisplayName)

		// 4. 統計情報の計算
		stats := CalculateLogicalNameStatistics(schema)
		assert.Equal(t, 1, stats.TotalTables)
		assert.Equal(t, 1, stats.TablesWithLogical)
		assert.Equal(t, 2, stats.TotalColumns)
		assert.Equal(t, 2, stats.ColumnsWithLogical)

		totalCoverage := stats.GetTotalCoveragePercentage()
		assert.Equal(t, 100.0, totalCoverage)

		// 5. ヘルパー関数のテスト
		helper := NewLogicalNameHelper()
		assert.True(t, helper.HasLogicalName(schema.Tables[0].LogicalName))
		assert.True(t, helper.IsLogicalNameAvailable(schema.Tables[0].Comment, "|"))

		preferredName := helper.GetPreferredName(
			schema.Tables[0].LogicalName,
			schema.Tables[0].Name,
			true,
		)
		assert.Equal(t, "ユーザーテーブル", preferredName)
	})
}
