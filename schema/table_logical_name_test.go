package schema

import (
	"testing"
)

func TestTable_SetLogicalNameFromComment(t *testing.T) {
	tests := []struct {
		name             string
		comment          string
		delimiter        string
		fallbackToName   bool
		expectedLogical  string
		expectedComment  string
	}{
		{
			name:             "コメントから論理名とコメントを分離",
			comment:          "ユーザーマスタ|システムのユーザー情報を管理",
			delimiter:        "|",
			fallbackToName:   false,
			expectedLogical:  "ユーザーマスタ",
			expectedComment:  "システムのユーザー情報を管理",
		},
		{
			name:             "区切り文字なし、フォールバック無効",
			comment:          "ユーザーテーブル",
			delimiter:        "|",
			fallbackToName:   false,
			expectedLogical:  "ユーザーテーブル",
			expectedComment:  "",
		},
		{
			name:             "区切り文字なし、フォールバック有効",
			comment:          "ユーザーテーブル",
			delimiter:        "|",
			fallbackToName:   true,
			expectedLogical:  "ユーザーテーブル",
			expectedComment:  "",
		},
		{
			name:             "空のコメント、フォールバック有効",
			comment:          "",
			delimiter:        "|",
			fallbackToName:   true,
			expectedLogical:  "users",
			expectedComment:  "",
		},
		{
			name:             "空のコメント、フォールバック無効",
			comment:          "",
			delimiter:        "|",
			fallbackToName:   false,
			expectedLogical:  "",
			expectedComment:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{
				Name:    "users",
				Comment: tt.comment,
			}
			
			table.SetLogicalNameFromComment(tt.delimiter, tt.fallbackToName)
			
			if table.LogicalName != tt.expectedLogical {
				t.Errorf("LogicalName = %v, want %v", table.LogicalName, tt.expectedLogical)
			}
			if table.Comment != tt.expectedComment {
				t.Errorf("Comment = %v, want %v", table.Comment, tt.expectedComment)
			}
		})
	}
}

func TestTable_GetDisplayName(t *testing.T) {
	tests := []struct {
		name          string
		table         Table
		displayFormat string
		expected      string
	}{
		{
			name: "physical_logical形式",
			table: Table{
				Name:        "public.users",
				LogicalName: "ユーザーマスタ",
			},
			displayFormat: "physical_logical",
			expected:      "public.users（ユーザーマスタ）",
		},
		{
			name: "logical_physical形式",
			table: Table{
				Name:        "public.users",
				LogicalName: "ユーザーマスタ",
			},
			displayFormat: "logical_physical",
			expected:      "ユーザーマスタ（public.users）",
		},
		{
			name: "論理名なし",
			table: Table{
				Name:        "public.users",
				LogicalName: "",
			},
			displayFormat: "physical_logical",
			expected:      "public.users",
		},
		{
			name: "不明な形式",
			table: Table{
				Name:        "public.users",
				LogicalName: "ユーザーマスタ",
			},
			displayFormat: "unknown",
			expected:      "public.users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.table.GetDisplayName(tt.displayFormat)
			if result != tt.expected {
				t.Errorf("GetDisplayName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTable_HasLogicalName(t *testing.T) {
	tests := []struct {
		name        string
		logicalName string
		expected    bool
	}{
		{
			name:        "論理名が設定されている場合",
			logicalName: "ユーザーマスタ",
			expected:    true,
		},
		{
			name:        "論理名が空の場合",
			logicalName: "",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := Table{
				LogicalName: tt.logicalName,
			}
			
			result := table.HasLogicalName()
			if result != tt.expected {
				t.Errorf("HasLogicalName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTable_GetLogicalNameOrFallback(t *testing.T) {
	tests := []struct {
		name           string
		table          Table
		fallbackToName bool
		expected       string
	}{
		{
			name: "論理名が設定されている場合",
			table: Table{
				Name:        "users",
				LogicalName: "ユーザーマスタ",
			},
			fallbackToName: false,
			expected:       "ユーザーマスタ",
		},
		{
			name: "論理名が空、フォールバック有効",
			table: Table{
				Name:        "users",
				LogicalName: "",
			},
			fallbackToName: true,
			expected:       "users",
		},
		{
			name: "論理名が空、フォールバック無効",
			table: Table{
				Name:        "users",
				LogicalName: "",
			},
			fallbackToName: false,
			expected:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.table.GetLogicalNameOrFallback(tt.fallbackToName)
			if result != tt.expected {
				t.Errorf("GetLogicalNameOrFallback() = %v, want %v", result, tt.expected)
			}
		})
	}
}