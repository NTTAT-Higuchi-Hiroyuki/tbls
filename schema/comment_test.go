package schema

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultCommentParser_ParseComment(t *testing.T) {
	parser := NewDefaultCommentParser()

	tests := []struct {
		name            string
		rawComment      string
		separator       string
		expectedLogical string
		expectedClean   string
	}{
		// T001: 標準区切り文字での論理名抽出
		{
			name:            "標準区切り文字での論理名抽出",
			rawComment:      "ユーザーID|システム内でユーザーを一意に識別するID",
			separator:       "|",
			expectedLogical: "ユーザーID",
			expectedClean:   "システム内でユーザーを一意に識別するID",
		},
		// T002: カスタム区切り文字での論理名抽出
		{
			name:            "カスタム区切り文字での論理名抽出",
			rawComment:      "ユーザー名::ログイン時に使用する名前",
			separator:       "::",
			expectedLogical: "ユーザー名",
			expectedClean:   "ログイン時に使用する名前",
		},
		// T003: 日本語論理名の処理
		{
			name:            "日本語論理名の処理",
			rawComment:      "顧客番号|お客様を識別する番号",
			separator:       "|",
			expectedLogical: "顧客番号",
			expectedClean:   "お客様を識別する番号",
		},
		// T004: 英数字論理名の処理
		{
			name:            "英数字論理名の処理",
			rawComment:      "UserID|Unique identifier for user",
			separator:       "|",
			expectedLogical: "UserID",
			expectedClean:   "Unique identifier for user",
		},
		// T005: 区切り文字なしのコメント処理
		{
			name:            "区切り文字なしのコメント処理",
			rawComment:      "通常のコメント",
			separator:       "|",
			expectedLogical: "",
			expectedClean:   "通常のコメント",
		},
		// T101: 空文字列入力
		{
			name:            "空文字列入力",
			rawComment:      "",
			separator:       "|",
			expectedLogical: "",
			expectedClean:   "",
		},
		// T103: 区切り文字のみ
		{
			name:            "区切り文字のみ",
			rawComment:      "|",
			separator:       "|",
			expectedLogical: "",
			expectedClean:   "",
		},
		// T104: 複数区切り文字
		{
			name:            "複数区切り文字",
			rawComment:      "論理名|説明|追加情報",
			separator:       "|",
			expectedLogical: "論理名",
			expectedClean:   "説明|追加情報",
		},
		// T201: 最小長論理名
		{
			name:            "最小長論理名",
			rawComment:      "A|説明",
			separator:       "|",
			expectedLogical: "A",
			expectedClean:   "説明",
		},
		// T203: 空の説明部分
		{
			name:            "空の説明部分",
			rawComment:      "論理名|",
			separator:       "|",
			expectedLogical: "論理名",
			expectedClean:   "",
		},
		// T204: 特殊文字を含む論理名
		{
			name:            "特殊文字を含む論理名",
			rawComment:      "論理名_01|説明#1",
			separator:       "|",
			expectedLogical: "論理名_01",
			expectedClean:   "説明#1",
		},
		// 空白の処理
		{
			name:            "前後の空白処理",
			rawComment:      " 論理名 | 説明文 ",
			separator:       "|",
			expectedLogical: "論理名",
			expectedClean:   "説明文",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logicalName, cleanComment := parser.ParseComment(tt.rawComment, tt.separator)
			assert.Equal(t, tt.expectedLogical, logicalName)
			assert.Equal(t, tt.expectedClean, cleanComment)
		})
	}
}

func TestDefaultCommentParser_ExtractLogicalName(t *testing.T) {
	parser := NewDefaultCommentParser()

	tests := []struct {
		name       string
		rawComment string
		separator  string
		expected   string
	}{
		{
			name:       "論理名の抽出",
			rawComment: "ユーザーID|システム内でユーザーを一意に識別するID",
			separator:  "|",
			expected:   "ユーザーID",
		},
		{
			name:       "論理名なし",
			rawComment: "通常のコメント",
			separator:  "|",
			expected:   "",
		},
		{
			name:       "空文字列",
			rawComment: "",
			separator:  "|",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.ExtractLogicalName(tt.rawComment, tt.separator)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultCommentParser_ExtractCleanComment(t *testing.T) {
	parser := NewDefaultCommentParser()

	tests := []struct {
		name       string
		rawComment string
		separator  string
		expected   string
	}{
		{
			name:       "クリーンコメントの抽出",
			rawComment: "ユーザーID|システム内でユーザーを一意に識別するID",
			separator:  "|",
			expected:   "システム内でユーザーを一意に識別するID",
		},
		{
			name:       "論理名なしの場合",
			rawComment: "通常のコメント",
			separator:  "|",
			expected:   "通常のコメント",
		},
		{
			name:       "空文字列",
			rawComment: "",
			separator:  "|",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.ExtractCleanComment(tt.rawComment, tt.separator)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultCommentParser_HasLogicalName(t *testing.T) {
	parser := NewDefaultCommentParser()

	tests := []struct {
		name       string
		rawComment string
		separator  string
		expected   bool
	}{
		{
			name:       "論理名あり",
			rawComment: "ユーザーID|システム内でユーザーを一意に識別するID",
			separator:  "|",
			expected:   true,
		},
		{
			name:       "論理名なし",
			rawComment: "通常のコメント",
			separator:  "|",
			expected:   false,
		},
		{
			name:       "空文字列",
			rawComment: "",
			separator:  "|",
			expected:   false,
		},
		{
			name:       "区切り文字のみ",
			rawComment: "|",
			separator:  "|",
			expected:   false,
		},
		{
			name:       "空の区切り文字",
			rawComment: "論理名|説明",
			separator:  "",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.HasLogicalName(tt.rawComment, tt.separator)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultCommentParser_UTF8Safety(t *testing.T) {
	parser := NewDefaultCommentParser()

	tests := []struct {
		name            string
		rawComment      string
		separator       string
		expectedLogical string
		expectedClean   string
	}{
		{
			name:            "日本語文字",
			rawComment:      "テーブル名前|これは日本語のコメントです",
			separator:       "|",
			expectedLogical: "テーブル名前",
			expectedClean:   "これは日本語のコメントです",
		},
		{
			name:            "絵文字を含む",
			rawComment:      "ユーザー😀|ユーザー情報📝",
			separator:       "|",
			expectedLogical: "ユーザー😀",
			expectedClean:   "ユーザー情報📝",
		},
		{
			name:            "中国語文字",
			rawComment:      "用户表|用户基本信息",
			separator:       "|",
			expectedLogical: "用户表",
			expectedClean:   "用户基本信息",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logicalName, cleanComment := parser.ParseComment(tt.rawComment, tt.separator)
			assert.Equal(t, tt.expectedLogical, logicalName)
			assert.Equal(t, tt.expectedClean, cleanComment)
		})
	}
}

func TestDefaultCommentParser_EdgeCases(t *testing.T) {
	parser := NewDefaultCommentParser()

	t.Run("非常に長いコメント", func(t *testing.T) {
		longLogical := strings.Repeat("A", 1000)
		longComment := strings.Repeat("説明", 1000)
		rawComment := longLogical + "|" + longComment

		logicalName, cleanComment := parser.ParseComment(rawComment, "|")
		assert.Equal(t, longLogical, logicalName)
		assert.Equal(t, longComment, cleanComment)
	})

	t.Run("最大長論理名", func(t *testing.T) {
		longLogical := strings.Repeat("A", 255)
		rawComment := longLogical + "|説明"

		logicalName, cleanComment := parser.ParseComment(rawComment, "|")
		assert.Equal(t, longLogical, logicalName)
		assert.Equal(t, "説明", cleanComment)
	})

	t.Run("複数文字の区切り文字", func(t *testing.T) {
		rawComment := "論理名::説明文"
		logicalName, cleanComment := parser.ParseComment(rawComment, "::")
		assert.Equal(t, "論理名", logicalName)
		assert.Equal(t, "説明文", cleanComment)
	})
}

// 便利関数のテスト
func TestParseCommentHelpers(t *testing.T) {
	t.Run("ParseComment便利関数", func(t *testing.T) {
		logicalName, cleanComment := ParseComment("論理名|説明", "|")
		assert.Equal(t, "論理名", logicalName)
		assert.Equal(t, "説明", cleanComment)
	})

	t.Run("ExtractLogicalName便利関数", func(t *testing.T) {
		result := ExtractLogicalName("論理名|説明", "|")
		assert.Equal(t, "論理名", result)
	})

	t.Run("ExtractCleanComment便利関数", func(t *testing.T) {
		result := ExtractCleanComment("論理名|説明", "|")
		assert.Equal(t, "説明", result)
	})
}

func TestParseCommentWithParser(t *testing.T) {
	t.Run("通常のパーサー", func(t *testing.T) {
		parser := NewDefaultCommentParser()
		logicalName, cleanComment := ParseCommentWithParser("論理名|説明", "|", parser)
		assert.Equal(t, "論理名", logicalName)
		assert.Equal(t, "説明", cleanComment)
	})

	t.Run("nilパーサー", func(t *testing.T) {
		logicalName, cleanComment := ParseCommentWithParser("論理名|説明", "|", nil)
		assert.Equal(t, "論理名", logicalName)
		assert.Equal(t, "説明", cleanComment)
	})
}

func TestValidateCommentSeparator(t *testing.T) {
	tests := []struct {
		name        string
		separator   string
		expectError bool
	}{
		{
			name:        "正常な区切り文字",
			separator:   "|",
			expectError: false,
		},
		{
			name:        "複数文字の区切り文字",
			separator:   "::",
			expectError: false,
		},
		{
			name:        "空文字列",
			separator:   "",
			expectError: true,
		},
		{
			name:        "制御文字",
			separator:   "\x01",
			expectError: true,
		},
		{
			name:        "タブ文字（許可）",
			separator:   "\t",
			expectError: false,
		},
		{
			name:        "改行文字",
			separator:   "\n",
			expectError: true,
		},
		{
			name:        "日本語区切り文字",
			separator:   "・",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommentSeparator(tt.separator)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ベンチマークテスト
func BenchmarkDefaultCommentParser_ParseComment(b *testing.B) {
	parser := NewDefaultCommentParser()
	comment := "ユーザーID|システム内でユーザーを一意に識別するID"
	separator := "|"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.ParseComment(comment, separator)
	}
}

func BenchmarkParseComment_LongComment(b *testing.B) {
	parser := NewDefaultCommentParser()
	longComment := strings.Repeat("論理名", 100) + "|" + strings.Repeat("説明文", 100)
	separator := "|"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.ParseComment(longComment, separator)
	}
}

// インターフェースの実装確認
func TestCommentParserInterface(t *testing.T) {
	var parser CommentParser = NewDefaultCommentParser()

	logicalName, cleanComment := parser.ParseComment("論理名|説明", "|")
	assert.Equal(t, "論理名", logicalName)
	assert.Equal(t, "説明", cleanComment)

	assert.Equal(t, "論理名", parser.ExtractLogicalName("論理名|説明", "|"))
	assert.Equal(t, "説明", parser.ExtractCleanComment("論理名|説明", "|"))
	assert.True(t, parser.HasLogicalName("論理名|説明", "|"))
	assert.False(t, parser.HasLogicalName("説明のみ", "|"))
}
