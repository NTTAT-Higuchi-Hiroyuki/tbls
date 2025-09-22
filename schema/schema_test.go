package schema

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNormalizeTableName(t *testing.T) {
	tests := []struct {
		s    *Schema
		name string
		want string
	}{
		{&Schema{}, "testtable", "testtable"},
		{&Schema{Driver: &Driver{Name: "postgres", Meta: &DriverMeta{CurrentSchema: "public"}}}, "testtable", "public.testtable"},
		{&Schema{Driver: &Driver{Name: "mysql", Meta: &DriverMeta{CurrentSchema: "public"}}}, "testtable", "testtable"},
		{&Schema{Driver: &Driver{Name: "postgres", Meta: &DriverMeta{CurrentSchema: ""}}}, "testtable", "testtable"},
		{&Schema{Driver: &Driver{Name: "postgres", Meta: &DriverMeta{CurrentSchema: "public"}}}, "other.testtable", "other.testtable"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := tt.s.NormalizeTableName(tt.name)
			if got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestSchema_FindTableByName(t *testing.T) {
	schema := Schema{
		Name: "testschema",
		Tables: []*Table{
			&Table{
				Name:    "a",
				Comment: "table a",
			},
			&Table{
				Name:    "b",
				Comment: "table b",
			},
		},
	}
	table, _ := schema.FindTableByName("b")
	want := "table b"
	got := table.Comment
	if got != want {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

func TestTable_FindColumnByName(t *testing.T) {
	table := Table{
		Name: "testtable",
		Columns: []*Column{
			&Column{
				Name:    "a",
				Comment: "column a",
			},
			&Column{
				Name:    "b",
				Comment: "column b",
			},
		},
	}
	column, _ := table.FindColumnByName("b")
	want := "column b"
	got := column.Comment
	if got != want {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

func TestTable_FindConstrainsByColumnName(t *testing.T) {
	table := Table{
		Name: "testtable",
		Columns: []*Column{
			&Column{
				Name:    "a",
				Comment: "column a",
			},
			&Column{
				Name:    "b",
				Comment: "column b",
			},
		},
	}
	table.Constraints = []*Constraint{
		&Constraint{
			Name:              "PRIMARY",
			Type:              "PRIMARY KEY",
			Def:               "PRIMARY KEY(a)",
			ReferencedTable:   nil,
			Table:             &table.Name,
			Columns:           []string{"a"},
			ReferencedColumns: []string{},
		},
		&Constraint{
			Name:              "UNIQUE",
			Type:              "UNIQUE",
			Def:               "UNIQUE KEY a (b)",
			ReferencedTable:   nil,
			Table:             &table.Name,
			Columns:           []string{"b"},
			ReferencedColumns: []string{},
		},
	}

	got := table.FindConstrainsByColumnName("a")
	if want := 1; len(got) != want {
		t.Errorf("got %v\nwant %v", len(got), want)
	}
	if want := "PRIMARY"; got[0].Name != want {
		t.Errorf("got %v\nwant %v", got[0].Name, want)
	}
}

func TestTable_hasColumnWithValues(t *testing.T) {
	tests := []struct {
		testName  string
		name      string
		addColumn *Column
		want      bool
	}{
		{"Has no ExtraDef value.", ColumnExtraDef, &Column{Name: "b"}, false},
		{"Has ExtraDef value", ColumnExtraDef, &Column{Name: "b", ExtraDef: "ExtraDef"}, true},
		{"Occurrences is invalid", ColumnOccurrences, &Column{Name: "b", Occurrences: sql.NullInt32{Valid: false}}, false},
		{"Occurrences is valid", ColumnOccurrences, &Column{Name: "b", Occurrences: sql.NullInt32{Valid: true}}, true},
		{"Percents is invalid", ColumnPercents, &Column{Name: "b", Percents: sql.NullFloat64{Valid: false}}, false},
		{"Percents is valid", ColumnPercents, &Column{Name: "b", Percents: sql.NullFloat64{Valid: true}}, true},
		{"Has no ChildRelations", ColumnChildren, &Column{Name: "b"}, false},
		{"Has ChildRelations", ColumnChildren, &Column{Name: "b", ChildRelations: []*Relation{{}}}, true},
		{"Has no ParentRelations", ColumnParents, &Column{Name: "b"}, false},
		{"Has ParentRelations", ColumnParents, &Column{Name: "b", ParentRelations: []*Relation{{}}}, true},
		{"Has no Comment", ColumnComment, &Column{Name: "b"}, false},
		{"Has Comment", ColumnComment, &Column{Name: "b", Comment: "comment"}, true},
		{"Has no Labels", ColumnLabels, &Column{Name: "b"}, false},
		{"Has Labels", ColumnLabels, &Column{Name: "b", Labels: Labels{{Name: "TestLabel"}}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			table := Table{
				Name: "testTable",
				Columns: []*Column{
					&Column{
						Name: "a",
					},
				},
			}
			table.Columns = append(table.Columns, tt.addColumn)

			got := table.hasColumnWithValues(tt.name)
			if got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestTable_ShowColumn(t *testing.T) {
	tests := []struct {
		testName    string
		table       Table
		name        string
		hideColumns []string
		want        bool
	}{
		{"No hideColumns", Table{Name: "testTable"}, ColumnComment, []string{}, true},
		{"hideColumns without value", Table{Name: "testTable"}, ColumnComment, []string{ColumnComment}, false},
		{"hideColumns with value", Table{Name: "testTable", Columns: []*Column{{Name: "testColumn", Comment: "comment"}}}, ColumnComment, []string{ColumnComment}, true},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := tt.table.ShowColumn(tt.name, tt.hideColumns)
			if got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestSchema_Sort(t *testing.T) {
	schema := Schema{
		Name: "testschema",
		Tables: []*Table{
			&Table{
				Name:    "b",
				Comment: "table b",
			},
			&Table{
				Name:    "a",
				Comment: "table a",
				Columns: []*Column{
					&Column{
						Name:    "b",
						Comment: "column b",
					},
					&Column{
						Name:    "a",
						Comment: "column a",
					},
				},
			},
		},
		Functions: []*Function{
			&Function{
				Name:      "b",
				Arguments: "arg b",
			},
			&Function{
				Name:      "b",
				Arguments: "arg a",
			},
		},
	}
	if err := schema.Sort(); err != nil {
		t.Error(err)
	}
	want := "a"
	got := schema.Tables[0].Name
	if got != want {
		t.Errorf("got %v\nwant %v", got, want)
	}
	want2 := "a"
	got2 := schema.Tables[0].Columns[0].Name
	if got2 != want2 {
		t.Errorf("got %v\nwant %v", got2, want2)
	}
	want3 := "arg a"
	got3 := schema.Functions[0].Arguments
	if got3 != want3 {
		t.Errorf("got %v\nwant %v", got3, want3)
	}
}

func TestRepair(t *testing.T) {
	got := &Schema{}
	f, err := os.Open(filepath.Join(testdataDir(), "test_repair.golden"))
	if err != nil {
		t.Error(err)
	}
	dec := json.NewDecoder(f)
	if err := dec.Decode(got); err != nil {
		t.Error(err)
	}
	if err := got.Repair(); err != nil {
		t.Error(err)
	}

	want := newTestSchema(t)

	if diff := cmp.Diff(got, want, nil); diff != "" {
		t.Errorf("%s", diff)
	}

	b, err := json.Marshal(want)
	if err != nil {
		t.Error(err)
	}
	want2 := &Schema{}
	if err := json.Unmarshal(b, want2); err != nil {
		t.Error(err)
	}
	if err := want2.Repair(); err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(want, want2, nil); diff != "" {
		t.Errorf("%s", diff)
	}
}

func TestClone(t *testing.T) {
	want := newTestSchema(t)
	got, err := want.Clone()
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(got, want, nil); diff != "" {
		t.Errorf("%s", diff)
	}
}

// LogicalName機能のテスト

func TestTable_SetLogicalNameFromComment(t *testing.T) {
	tests := []struct {
		name      string
		table     *Table
		separator string
		want      string
	}{
		{
			name:      "正常な論理名とコメントの分離",
			table:     &Table{Name: "users", Comment: "ユーザーテーブル|システムの利用者を管理する"},
			separator: "|",
			want:      "ユーザーテーブル",
		},
		{
			name:      "区切り文字が含まれていない場合",
			table:     &Table{Name: "users", Comment: "ユーザーテーブル"},
			separator: "|",
			want:      "",
		},
		{
			name:      "コメントが空の場合",
			table:     &Table{Name: "users", Comment: ""},
			separator: "|",
			want:      "",
		},
		{
			name:      "区切り文字が空の場合",
			table:     &Table{Name: "users", Comment: "ユーザーテーブル|説明"},
			separator: "",
			want:      "",
		},
		{
			name:      "複数の区切り文字がある場合",
			table:     &Table{Name: "users", Comment: "ユーザーテーブル|説明|詳細説明"},
			separator: "|",
			want:      "ユーザーテーブル",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.table.SetLogicalNameFromComment(tt.separator)
			if tt.table.LogicalName != tt.want {
				t.Errorf("Table.SetLogicalNameFromComment() LogicalName = %v, want %v", tt.table.LogicalName, tt.want)
			}
		})
	}
}

func TestTable_GetLogicalNameOrFallback(t *testing.T) {
	tests := []struct {
		name    string
		table   *Table
		want    string
	}{
		{
			name:    "論理名が設定されている場合",
			table:   &Table{Name: "users", LogicalName: "ユーザーテーブル"},
			want:    "ユーザーテーブル",
		},
		{
			name:    "論理名が空の場合",
			table:   &Table{Name: "users", LogicalName: ""},
			want:    "users",
		},
		{
			name:    "論理名が未設定の場合",
			table:   &Table{Name: "users"},
			want:    "users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.table.GetLogicalNameOrFallback()
			if got != tt.want {
				t.Errorf("Table.GetLogicalNameOrFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColumn_SetLogicalNameFromComment(t *testing.T) {
	tests := []struct {
		name      string
		column    *Column
		separator string
		want      string
	}{
		{
			name:      "正常な論理名とコメントの分離",
			column:    &Column{Name: "user_id", Comment: "ユーザーID|ユーザーの一意識別子"},
			separator: "|",
			want:      "ユーザーID",
		},
		{
			name:      "区切り文字が含まれていない場合",
			column:    &Column{Name: "user_id", Comment: "ユーザーID"},
			separator: "|",
			want:      "",
		},
		{
			name:      "コメントが空の場合",
			column:    &Column{Name: "user_id", Comment: ""},
			separator: "|",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.column.SetLogicalNameFromComment(tt.separator)
			if tt.column.LogicalName != tt.want {
				t.Errorf("Column.SetLogicalNameFromComment() LogicalName = %v, want %v", tt.column.LogicalName, tt.want)
			}
		})
	}
}

func TestColumn_GetLogicalNameOrFallback(t *testing.T) {
	tests := []struct {
		name   string
		column *Column
		want   string
	}{
		{
			name:   "論理名が設定されている場合",
			column: &Column{Name: "user_id", LogicalName: "ユーザーID"},
			want:   "ユーザーID",
		},
		{
			name:   "論理名が空の場合",
			column: &Column{Name: "user_id", LogicalName: ""},
			want:   "user_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.column.GetLogicalNameOrFallback()
			if got != tt.want {
				t.Errorf("Column.GetLogicalNameOrFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_SetLogicalNameFromComment(t *testing.T) {
	tests := []struct {
		name      string
		index     *Index
		separator string
		want      string
	}{
		{
			name:      "正常な論理名とコメントの分離",
			index:     &Index{Name: "idx_user_email", Comment: "メール検索インデックス|ユーザーのメールアドレス検索用"},
			separator: "|",
			want:      "メール検索インデックス",
		},
		{
			name:      "区切り文字が含まれていない場合",
			index:     &Index{Name: "idx_user_email", Comment: "メール検索インデックス"},
			separator: "|",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.index.SetLogicalNameFromComment(tt.separator)
			if tt.index.LogicalName != tt.want {
				t.Errorf("Index.SetLogicalNameFromComment() LogicalName = %v, want %v", tt.index.LogicalName, tt.want)
			}
		})
	}
}

func TestIndex_GetLogicalNameOrFallback(t *testing.T) {
	tests := []struct {
		name  string
		index *Index
		want  string
	}{
		{
			name:  "論理名が設定されている場合",
			index: &Index{Name: "idx_user_email", LogicalName: "メール検索インデックス"},
			want:  "メール検索インデックス",
		},
		{
			name:  "論理名が空の場合",
			index: &Index{Name: "idx_user_email", LogicalName: ""},
			want:  "idx_user_email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.index.GetLogicalNameOrFallback()
			if got != tt.want {
				t.Errorf("Index.GetLogicalNameOrFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstraint_SetLogicalNameFromComment(t *testing.T) {
	tests := []struct {
		name       string
		constraint *Constraint
		separator  string
		want       string
	}{
		{
			name:       "正常な論理名とコメントの分離",
			constraint: &Constraint{Name: "fk_user_dept", Comment: "部署参照制約|ユーザーと部署の関連"},
			separator:  "|",
			want:       "部署参照制約",
		},
		{
			name:       "区切り文字が含まれていない場合",
			constraint: &Constraint{Name: "fk_user_dept", Comment: "部署参照制約"},
			separator:  "|",
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.constraint.SetLogicalNameFromComment(tt.separator)
			if tt.constraint.LogicalName != tt.want {
				t.Errorf("Constraint.SetLogicalNameFromComment() LogicalName = %v, want %v", tt.constraint.LogicalName, tt.want)
			}
		})
	}
}

func TestConstraint_GetLogicalNameOrFallback(t *testing.T) {
	tests := []struct {
		name       string
		constraint *Constraint
		want       string
	}{
		{
			name:       "論理名が設定されている場合",
			constraint: &Constraint{Name: "fk_user_dept", LogicalName: "部署参照制約"},
			want:       "部署参照制約",
		},
		{
			name:       "論理名が空の場合",
			constraint: &Constraint{Name: "fk_user_dept", LogicalName: ""},
			want:       "fk_user_dept",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.constraint.GetLogicalNameOrFallback()
			if got != tt.want {
				t.Errorf("Constraint.GetLogicalNameOrFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrigger_SetLogicalNameFromComment(t *testing.T) {
	tests := []struct {
		name      string
		trigger   *Trigger
		separator string
		want      string
	}{
		{
			name:      "正常な論理名とコメントの分離",
			trigger:   &Trigger{Name: "trg_user_audit", Comment: "ユーザー監査トリガー|ユーザー変更履歴の記録"},
			separator: "|",
			want:      "ユーザー監査トリガー",
		},
		{
			name:      "区切り文字が含まれていない場合",
			trigger:   &Trigger{Name: "trg_user_audit", Comment: "ユーザー監査トリガー"},
			separator: "|",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.trigger.SetLogicalNameFromComment(tt.separator)
			if tt.trigger.LogicalName != tt.want {
				t.Errorf("Trigger.SetLogicalNameFromComment() LogicalName = %v, want %v", tt.trigger.LogicalName, tt.want)
			}
		})
	}
}

func TestTrigger_GetLogicalNameOrFallback(t *testing.T) {
	tests := []struct {
		name    string
		trigger *Trigger
		want    string
	}{
		{
			name:    "論理名が設定されている場合",
			trigger: &Trigger{Name: "trg_user_audit", LogicalName: "ユーザー監査トリガー"},
			want:    "ユーザー監査トリガー",
		},
		{
			name:    "論理名が空の場合",
			trigger: &Trigger{Name: "trg_user_audit", LogicalName: ""},
			want:    "trg_user_audit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.trigger.GetLogicalNameOrFallback()
			if got != tt.want {
				t.Errorf("Trigger.GetLogicalNameOrFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

// JSON/YAMLシリアライズのテスト
func TestLogicalNameSerialization(t *testing.T) {
	tests := []struct {
		name string
		obj  interface{}
	}{
		{
			name: "Tableのシリアライズ",
			obj: &Table{
				Name:        "users",
				LogicalName: "ユーザーテーブル",
				Comment:     "システムユーザー",
			},
		},
		{
			name: "Columnのシリアライズ",
			obj: &Column{
				Name:        "user_id",
				LogicalName: "ユーザーID",
				Comment:     "一意識別子",
				Type:        "bigint",
			},
		},
		{
			name: "Indexのシリアライズ",
			obj: &Index{
				Name:        "idx_user_email",
				LogicalName: "メール検索インデックス",
				Comment:     "メール検索用",
			},
		},
		{
			name: "Constraintのシリアライズ",
			obj: &Constraint{
				Name:        "fk_user_dept",
				LogicalName: "部署参照制約",
				Comment:     "部署との関連",
			},
		},
		{
			name: "Triggerのシリアライズ",
			obj: &Trigger{
				Name:        "trg_user_audit",
				LogicalName: "ユーザー監査トリガー",
				Comment:     "変更履歴記録",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// JSONシリアライズテスト
			jsonData, err := json.Marshal(tt.obj)
			if err != nil {
				t.Errorf("JSON Marshal error: %v", err)
				return
			}

			// JSONデシリアライズテスト
			switch orig := tt.obj.(type) {
			case *Table:
				var unmarshaled Table
				if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
					t.Errorf("JSON Unmarshal error: %v", err)
					return
				}
				if unmarshaled.LogicalName != orig.LogicalName {
					t.Errorf("JSON serialization failed: got %v, want %v", unmarshaled.LogicalName, orig.LogicalName)
				}
			case *Column:
				var unmarshaled Column
				if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
					t.Errorf("JSON Unmarshal error: %v", err)
					return
				}
				if unmarshaled.LogicalName != orig.LogicalName {
					t.Errorf("JSON serialization failed: got %v, want %v", unmarshaled.LogicalName, orig.LogicalName)
				}
			case *Index:
				var unmarshaled Index
				if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
					t.Errorf("JSON Unmarshal error: %v", err)
					return
				}
				if unmarshaled.LogicalName != orig.LogicalName {
					t.Errorf("JSON serialization failed: got %v, want %v", unmarshaled.LogicalName, orig.LogicalName)
				}
			case *Constraint:
				var unmarshaled Constraint
				if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
					t.Errorf("JSON Unmarshal error: %v", err)
					return
				}
				if unmarshaled.LogicalName != orig.LogicalName {
					t.Errorf("JSON serialization failed: got %v, want %v", unmarshaled.LogicalName, orig.LogicalName)
				}
			case *Trigger:
				var unmarshaled Trigger
				if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
					t.Errorf("JSON Unmarshal error: %v", err)
					return
				}
				if unmarshaled.LogicalName != orig.LogicalName {
					t.Errorf("JSON serialization failed: got %v, want %v", unmarshaled.LogicalName, orig.LogicalName)
				}
			}
		})
	}
}

func testdataDir() string {
	wd, _ := os.Getwd()
	dir, _ := filepath.Abs(filepath.Join(filepath.Dir(wd), "testdata"))
	return dir
}

func newTestSchema(t *testing.T) *Schema {
	t.Helper()
	ca := &Column{
		Name:     "a",
		Type:     "bigint(20)",
		Comment:  "column a",
		Nullable: false,
	}
	cb := &Column{
		Name:     "b",
		Type:     "text",
		Comment:  "column b",
		Nullable: true,
	}

	ta := &Table{
		Name:    "a",
		Type:    "BASE TABLE",
		Comment: "table a",
		Columns: []*Column{
			ca,
			&Column{
				Name:     "a2",
				Type:     "datetime",
				Comment:  "column a2",
				Nullable: false,
				Default: sql.NullString{
					String: "CURRENT_TIMESTAMP",
					Valid:  true,
				},
			},
		},
	}

	tb := &Table{
		Name:    "b",
		Type:    "BASE TABLE",
		Comment: "table b",
		Columns: []*Column{
			cb,
			&Column{
				Name:     "b2",
				Comment:  "column b2",
				Type:     "text",
				Nullable: true,
			},
		},
	}
	r := &Relation{
		Table:         ta,
		Columns:       []*Column{ca},
		ParentTable:   tb,
		ParentColumns: []*Column{cb},
	}
	ca.ParentRelations = []*Relation{r}
	cb.ChildRelations = []*Relation{r}

	s := &Schema{
		Name: "testschema",
		Tables: []*Table{
			ta,
			tb,
		},
		Relations: []*Relation{
			r,
		},
		Driver: &Driver{
			Name:            "testdriver",
			DatabaseVersion: "1.0.0",
			Meta:            &DriverMeta{},
		},
	}
	return s
}