package schema

import (
	"testing"
)

func TestLegacyParser_Name(t *testing.T) {
	parser := NewLegacyParser()
	if parser.Name() != "legacy" {
		t.Errorf("expected name 'legacy', got %s", parser.Name())
	}
}

func TestLegacyParser_Priority(t *testing.T) {
	parser := NewLegacyParser()
	if parser.Priority() != 1000 {
		t.Errorf("expected priority 1000, got %d", parser.Priority())
	}
}

func TestLegacyParser_CanParse(t *testing.T) {
	parser := NewLegacyParser()
	
	tests := []struct {
		name    string
		comment string
		want    bool
	}{
		{"空文字列", "", true},
		{"通常のコメント", "論理名|説明", true},
		{"JSON形式", `{"name": "test"}`, true},
		{"YAML形式", "name: test", true},
		{"特殊文字", "テスト<>&", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.CanParse(tt.comment)
			if got != tt.want {
				t.Errorf("CanParse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLegacyParser_ParseComment(t *testing.T) {
	parser := NewLegacyParser()
	
	tests := []struct {
		name      string
		comment   string
		delimiter string
		wantLogicalName string
		wantDescription string
		wantError bool
	}{
		{
			name:            "空文字列",
			comment:         "",
			delimiter:       "|",
			wantLogicalName: "",
			wantDescription: "",
			wantError:       false,
		},
		{
			name:            "基本的な分割",
			comment:         "論理名|説明文",
			delimiter:       "|",
			wantLogicalName: "論理名",
			wantDescription: "説明文",
			wantError:       false,
		},
		{
			name:            "論理名のみ",
			comment:         "論理名のみ",
			delimiter:       "|",
			wantLogicalName: "論理名のみ",
			wantDescription: "",
			wantError:       false,
		},
		{
			name:            "空白の正規化",
			comment:         "  論理名  |  説明文  ",
			delimiter:       "|",
			wantLogicalName: "論理名",
			wantDescription: "説明文",
			wantError:       false,
		},
		{
			name:            "連続空白の処理",
			comment:         "論理名   複数空白|説明   複数   空白",
			delimiter:       "|",
			wantLogicalName: "論理名 複数空白",
			wantDescription: "説明 複数 空白",
			wantError:       false,
		},
		{
			name:            "カスタム区切り文字",
			comment:         "論理名::説明文",
			delimiter:       "::",
			wantLogicalName: "論理名",
			wantDescription: "説明文",
			wantError:       false,
		},
		{
			name:            "エスケープ処理",
			comment:         "論理名\\|エスケープ|説明文",
			delimiter:       "|",
			wantLogicalName: "論理名|エスケープ",
			wantDescription: "説明文",
			wantError:       false,
		},
		{
			name:            "日本語文字",
			comment:         "ユーザー名|ユーザーの表示名を格納する",
			delimiter:       "|",
			wantLogicalName: "ユーザー名",
			wantDescription: "ユーザーの表示名を格納する",
			wantError:       false,
		},
		{
			name:            "複数の区切り文字",
			comment:         "論理名|説明|追加情報",
			delimiter:       "|",
			wantLogicalName: "論理名",
			wantDescription: "説明",
			wantError:       false,
		},
		{
			name:            "区切り文字なし",
			comment:         "論理名のみ説明なし",
			delimiter:       "|",
			wantLogicalName: "論理名のみ説明なし",
			wantDescription: "",
			wantError:       false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseComment(tt.comment, tt.delimiter)
			
			if (err != nil) != tt.wantError {
				t.Errorf("ParseComment() error = %v, wantError %v", err, tt.wantError)
				return
			}
			
			if err != nil {
				return
			}
			
			if got.LogicalName != tt.wantLogicalName {
				t.Errorf("ParseComment() LogicalName = %q, want %q", got.LogicalName, tt.wantLogicalName)
			}
			
			if got.Description != tt.wantDescription {
				t.Errorf("ParseComment() Description = %q, want %q", got.Description, tt.wantDescription)
			}
			
			if got.Source != tt.comment {
				t.Errorf("ParseComment() Source = %q, want %q", got.Source, tt.comment)
			}
		})
	}
}

func TestLegacyParser_normalizeWhitespace(t *testing.T) {
	parser := NewLegacyParser()
	
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"通常の文字列", "通常の文字列", "通常の文字列"},
		{"前後の空白", "  前後の空白  ", "前後の空白"},
		{"連続する空白", "連続  する   空白", "連続 する 空白"},
		{"タブ文字", "タブ\t文字", "タブ 文字"},
		{"改行文字", "改行\n文字", "改行 文字"},
		{"混合空白", "混合\t \n空白", "混合 空白"},
		{"全角空白", "全角　空白", "全角 空白"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.normalizeWhitespace(tt.input)
			if got != tt.want {
				t.Errorf("normalizeWhitespace() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLegacyParser_splitWithEscape(t *testing.T) {
	parser := NewLegacyParser()
	
	tests := []struct {
		name      string
		input     string
		delimiter string
		want      []string
	}{
		{
			name:      "基本的な分割",
			input:     "a|b|c",
			delimiter: "|",
			want:      []string{"a", "b", "c"},
		},
		{
			name:      "エスケープあり",
			input:     "a\\|b|c",
			delimiter: "|",
			want:      []string{"a|b", "c"},
		},
		{
			name:      "複数エスケープ",
			input:     "a\\|b\\|c|d",
			delimiter: "|",
			want:      []string{"a|b|c", "d"},
		},
		{
			name:      "複数文字区切り",
			input:     "a::b::c",
			delimiter: "::",
			want:      []string{"a", "b", "c"},
		},
		{
			name:      "エスケープと複数文字区切り",
			input:     "a\\::b::c",
			delimiter: "::",
			want:      []string{"a::b", "c"},
		},
		{
			name:      "区切り文字なし",
			input:     "abc",
			delimiter: "|",
			want:      []string{"abc"},
		},
		{
			name:      "空文字列",
			input:     "",
			delimiter: "|",
			want:      []string{""},
		},
		{
			name:      "連続する区切り文字",
			input:     "a||b",
			delimiter: "|",
			want:      []string{"a", "", "b"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parser.splitWithEscape(tt.input, tt.delimiter)
			if len(got) != len(tt.want) {
				t.Errorf("splitWithEscape() length = %d, want %d", len(got), len(tt.want))
				return
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("splitWithEscape()[%d] = %q, want %q", i, v, tt.want[i])
				}
			}
		})
	}
}

func TestLegacyCommentSplitter(t *testing.T) {
	tests := []struct {
		name      string
		comment   string
		delimiter string
		wantLogicalName string
		wantComment     string
	}{
		{
			name:            "基本的な分割",
			comment:         "論理名|説明文",
			delimiter:       "|",
			wantLogicalName: "論理名",
			wantComment:     "説明文",
		},
		{
			name:            "空文字列",
			comment:         "",
			delimiter:       "|",
			wantLogicalName: "",
			wantComment:     "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LegacyCommentSplitter(tt.comment, tt.delimiter)
			if got.LogicalName != tt.wantLogicalName {
				t.Errorf("LegacyCommentSplitter() LogicalName = %q, want %q", got.LogicalName, tt.wantLogicalName)
			}
			if got.Comment != tt.wantComment {
				t.Errorf("LegacyCommentSplitter() Comment = %q, want %q", got.Comment, tt.wantComment)
			}
		})
	}
}

func TestEnhancedSplitComment(t *testing.T) {
	tests := []struct {
		name      string
		comment   string
		delimiter string
		wantLogicalName string
		wantDescription string
	}{
		{
			name:            "基本的な分割",
			comment:         "論理名|説明文",
			delimiter:       "|",
			wantLogicalName: "論理名",
			wantDescription: "説明文",
		},
		{
			name:            "空文字列",
			comment:         "",
			delimiter:       "|",
			wantLogicalName: "",
			wantDescription: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EnhancedSplitComment(tt.comment, tt.delimiter)
			if got.LogicalName != tt.wantLogicalName {
				t.Errorf("EnhancedSplitComment() LogicalName = %q, want %q", got.LogicalName, tt.wantLogicalName)
			}
			if got.Description != tt.wantDescription {
				t.Errorf("EnhancedSplitComment() Description = %q, want %q", got.Description, tt.wantDescription)
			}
		})
	}
}

func TestExtractLogicalNameEnhanced(t *testing.T) {
	tests := []struct {
		name           string
		comment        string
		delimiter      string
		physicalName   string
		fallbackToName bool
		want           string
	}{
		{
			name:           "論理名あり",
			comment:        "論理名|説明",
			delimiter:      "|",
			physicalName:   "physical_name",
			fallbackToName: false,
			want:           "論理名",
		},
		{
			name:           "論理名なし、フォールバックあり",
			comment:        "",
			delimiter:      "|",
			physicalName:   "physical_name",
			fallbackToName: true,
			want:           "physical_name",
		},
		{
			name:           "論理名なし、フォールバックなし",
			comment:        "",
			delimiter:      "|",
			physicalName:   "physical_name",
			fallbackToName: false,
			want:           "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractLogicalNameEnhanced(tt.comment, tt.delimiter, tt.physicalName, tt.fallbackToName)
			if got != tt.want {
				t.Errorf("ExtractLogicalNameEnhanced() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractCleanCommentEnhanced(t *testing.T) {
	tests := []struct {
		name      string
		comment   string
		delimiter string
		want      string
	}{
		{
			name:      "説明あり",
			comment:   "論理名|説明文",
			delimiter: "|",
			want:      "説明文",
		},
		{
			name:      "説明なし",
			comment:   "論理名のみ",
			delimiter: "|",
			want:      "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractCleanCommentEnhanced(tt.comment, tt.delimiter)
			if got != tt.want {
				t.Errorf("ExtractCleanCommentEnhanced() = %q, want %q", got, tt.want)
			}
		})
	}
}