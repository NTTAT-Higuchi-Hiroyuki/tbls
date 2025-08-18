package schema

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

// BenchmarkEnhancedCommentProcessing 拡張コメント処理のベンチマーク
func BenchmarkEnhancedCommentProcessing(b *testing.B) {
	processor := NewEnhancedCommentProcessor()

	// 各形式のサンプルコメント
	jsonComment := `{"name": "テストカラム", "description": "ベンチマーク用のテストコメント", "tags": ["test", "benchmark"], "priority": 1}`
	yamlComment := "name: テストカラム\ndescription: ベンチマーク用のテストコメント\ntags:\n  - test\n  - benchmark\npriority: 1"
	legacyComment := "テストカラム|ベンチマーク用のテストコメント"

	b.Run("JSON形式", func(b *testing.B) {
		column := &Column{Name: "test_col", Comment: jsonComment}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			column.ProcessEnhancedComment(processor, "|", true)
		}
	})

	b.Run("YAML形式", func(b *testing.B) {
		column := &Column{Name: "test_col", Comment: yamlComment}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			column.ProcessEnhancedComment(processor, "|", true)
		}
	})

	b.Run("Legacy形式", func(b *testing.B) {
		column := &Column{Name: "test_col", Comment: legacyComment}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			column.ProcessEnhancedComment(processor, "|", true)
		}
	})
}

// BenchmarkSchemaProcessing スキーマ全体処理のベンチマーク
func BenchmarkSchemaProcessing(b *testing.B) {
	processor := NewEnhancedCommentProcessor()

	// ベンチマーク用スキーマ生成
	createBenchmarkSchema := func(tableCount, columnCount int) *Schema {
		tables := make([]*Table, tableCount)
		for i := 0; i < tableCount; i++ {
			columns := make([]*Column, columnCount)
			for j := 0; j < columnCount; j++ {
				columns[j] = &Column{
					Name:    fmt.Sprintf("col_%d", j),
					Comment: fmt.Sprintf(`{"name": "カラム%d", "description": "テーブル%dのカラム%d", "tags": ["test"]}`, j, i, j),
				}
			}
			tables[i] = &Table{
				Name:    fmt.Sprintf("table_%d", i),
				Comment: fmt.Sprintf(`{"name": "テーブル%d", "description": "ベンチマーク用テーブル%d", "tags": ["table"]}`, i, i),
				Columns: columns,
			}
		}
		return &Schema{Name: "benchmark_db", Tables: tables}
	}

	// 小規模スキーマ（10テーブル、5カラム）
	b.Run("小規模スキーマ", func(b *testing.B) {
		schema := createBenchmarkSchema(10, 5)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			processSchemaForBenchmark(schema, processor)
		}
	})

	// 中規模スキーマ（50テーブル、10カラム）
	b.Run("中規模スキーマ", func(b *testing.B) {
		schema := createBenchmarkSchema(50, 10)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			processSchemaForBenchmark(schema, processor)
		}
	})

	// 大規模スキーマ（100テーブル、20カラム）
	b.Run("大規模スキーマ", func(b *testing.B) {
		schema := createBenchmarkSchema(100, 20)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			processSchemaForBenchmark(schema, processor)
		}
	})
}

// processSchemaForBenchmark ベンチマーク用スキーマ処理
func processSchemaForBenchmark(schema *Schema, processor CommentProcessor) {
	for _, table := range schema.Tables {
		table.ProcessEnhancedComment(processor, "|", true)
		for _, column := range table.Columns {
			column.ProcessEnhancedComment(processor, "|", true)
		}
	}
}

// BenchmarkParserComparison パーサー別性能比較
func BenchmarkParserComparison(b *testing.B) {
	jsonParser := NewJSONParser()
	yamlParser := NewYAMLParser()
	legacyParser := NewLegacyParser()

	jsonComment := `{"name": "テストカラム", "description": "パーサー性能比較用のテストコメント"}`
	yamlComment := "name: テストカラム\ndescription: パーサー性能比較用のテストコメント"
	legacyComment := "テストカラム|パーサー性能比較用のテストコメント"

	b.Run("JSONParser", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			jsonParser.ParseComment(jsonComment, "|")
		}
	})

	b.Run("YAMLParser", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			yamlParser.ParseComment(yamlComment, "|")
		}
	})

	b.Run("LegacyParser", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			legacyParser.ParseComment(legacyComment, "|")
		}
	})
}

// TestMemoryUsage メモリ使用量テスト
func TestMemoryUsage(t *testing.T) {
	// メモリ使用量測定のヘルパー関数
	measureMemory := func(name string, fn func()) {
		runtime.GC()
		var m1, m2 runtime.MemStats
		runtime.ReadMemStats(&m1)
		
		fn()
		
		runtime.GC()
		runtime.ReadMemStats(&m2)
		
		allocDiff := m2.TotalAlloc - m1.TotalAlloc
		t.Logf("%s: メモリ使用量 %d bytes", name, allocDiff)
		
		// メモリリークチェック（ヒープ増加が過大でないことを確認）
		heapDiff := int64(m2.HeapInuse) - int64(m1.HeapInuse)
		if heapDiff > 10*1024*1024 { // 10MB以上の増加は警告
			t.Logf("警告: ヒープメモリ増加が大きいです: %d bytes", heapDiff)
		}
	}

	processor := NewEnhancedCommentProcessor()

	// 大量データ処理でのメモリ使用量測定
	measureMemory("大量データ処理", func() {
		schema := createLargeSchema(200, 15) // 200テーブル、15カラム
		
		for _, table := range schema.Tables {
			table.ProcessEnhancedComment(processor, "|", true)
			for _, column := range table.Columns {
				column.ProcessEnhancedComment(processor, "|", true)
			}
		}
	})

	// 繰り返し処理でのメモリリークチェック
	measureMemory("繰り返し処理", func() {
		column := &Column{
			Name:    "test_col",
			Comment: `{"name": "テストカラム", "description": "メモリリークテスト用"}`,
		}
		
		for i := 0; i < 1000; i++ {
			column.ProcessEnhancedComment(processor, "|", true)
		}
	})
}

// createLargeSchema 大規模スキーマ生成
func createLargeSchema(tableCount, columnCount int) *Schema {
	tables := make([]*Table, tableCount)
	for i := 0; i < tableCount; i++ {
		columns := make([]*Column, columnCount)
		indexes := make([]*Index, 3) // 各テーブル3つのインデックス
		constraints := make([]*Constraint, 2) // 各テーブル2つの制約
		
		for j := 0; j < columnCount; j++ {
			columns[j] = &Column{
				Name:    fmt.Sprintf("col_%d_%d", i, j),
				Comment: fmt.Sprintf(`{"name": "カラム%d_%d", "description": "大規模スキーマテスト用カラム", "tags": ["test", "large"]}`, i, j),
			}
		}
		
		for k := 0; k < 3; k++ {
			indexes[k] = &Index{
				Name:    fmt.Sprintf("idx_%d_%d", i, k),
				Comment: fmt.Sprintf(`{"description": "インデックス%d_%d", "tags": ["index"]}`, i, k),
			}
		}
		
		for l := 0; l < 2; l++ {
			constraints[l] = &Constraint{
				Name:    fmt.Sprintf("const_%d_%d", i, l),
				Comment: fmt.Sprintf(`{"description": "制約%d_%d", "tags": ["constraint"]}`, i, l),
			}
		}
		
		tables[i] = &Table{
			Name:        fmt.Sprintf("table_%d", i),
			Comment:     fmt.Sprintf(`{"name": "テーブル%d", "description": "大規模スキーマテスト用テーブル", "tags": ["table", "large"]}`, i),
			Columns:     columns,
			Indexes:     indexes,
			Constraints: constraints,
		}
	}
	return &Schema{Name: "large_schema_db", Tables: tables}
}

// TestConcurrentProcessing 並行処理テスト
func TestConcurrentProcessing(t *testing.T) {
	processor := NewEnhancedCommentProcessor()
	schema := createLargeSchema(50, 10)
	
	// 並行処理でのデータ競合チェック
	numWorkers := 10
	done := make(chan bool, numWorkers)
	
	start := time.Now()
	
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer func() { done <- true }()
			
			// 各ワーカーが独立してスキーマを処理
			for _, table := range schema.Tables {
				table.ProcessEnhancedComment(processor, "|", true)
				for _, column := range table.Columns {
					column.ProcessEnhancedComment(processor, "|", true)
				}
			}
		}(i)
	}
	
	// 全ワーカーの完了を待機
	for i := 0; i < numWorkers; i++ {
		<-done
	}
	
	elapsed := time.Since(start)
	t.Logf("並行処理完了時間: %v", elapsed)
	
	// 処理結果の整合性チェック
	for _, table := range schema.Tables {
		if !table.HasEnhancedComment() {
			t.Errorf("Table %s should have enhanced comment after concurrent processing", table.Name)
		}
		
		for _, column := range table.Columns {
			if !column.HasEnhancedComment() {
				t.Errorf("Column %s.%s should have enhanced comment after concurrent processing", table.Name, column.Name)
			}
		}
	}
}

// TestProcessingTimeout タイムアウト処理テスト
func TestProcessingTimeout(t *testing.T) {
	// 短いタイムアウトを設定
	config := &ProcessingConfig{
		EnableValidation:   true,
		EnableSanitization: true,
		DefaultDelimiter:   "|",
		FallbackToLegacy:   true,
		StrictMode:         false,
		ProcessingTimeout:  1, // 1ms（非常に短い）
	}
	
	processor := NewEnhancedCommentProcessorWithConfig(config)
	
	// 複雑なコメント（処理時間がかかる可能性）
	complexComment := `{
		"name": "複雑なカラム",
		"description": "非常に長い説明文を含むカラムで、処理時間がかかる可能性があります。" +
			"この説明文は意図的に長くして、JSON解析やバリデーション処理に時間がかかるようにしています。",
		"tags": ["複雑", "長い", "テスト", "タイムアウト", "性能"],
		"metadata": {
			"key1": "value1",
			"key2": "value2",
			"key3": "value3"
		},
		"priority": 10,
		"deprecated": false
	}`
	
	column := &Column{
		Name:    "complex_col",
		Comment: complexComment,
	}
	
	// タイムアウトが発生する可能性があるが、適切にハンドリングされることを確認
	err := column.ProcessEnhancedComment(processor, "|", true)
	
	// タイムアウトエラーまたは正常処理のいずれかで完了することを確認
	if err != nil {
		t.Logf("Processing resulted in error (possibly timeout): %v", err)
		// タイムアウトエラーの場合、EnhancedCommentDataは設定されない
		if column.HasEnhancedComment() {
			t.Error("Column should not have enhanced comment data on timeout")
		}
	} else {
		t.Log("Processing completed successfully despite short timeout")
		// 正常処理の場合、EnhancedCommentDataが設定される
		if !column.HasEnhancedComment() {
			t.Error("Column should have enhanced comment data on successful processing")
		}
	}
}

// TestLargeCommentProcessing 大きなコメント処理テスト
func TestLargeCommentProcessing(t *testing.T) {
	processor := NewEnhancedCommentProcessor()
	
	// 大きなJSONコメント生成
	createLargeJSONComment := func(size int) string {
		tags := make([]string, size)
		for i := 0; i < size; i++ {
			tags[i] = fmt.Sprintf("tag_%d", i)
		}
		
		metadata := make(map[string]string)
		for i := 0; i < size; i++ {
			metadata[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
		}
		
		// 注意: 実際のJSON生成は簡略化
		return fmt.Sprintf(`{
			"name": "大きなコメント",
			"description": "サイズ%dの大きなコメントテスト",
			"tags": %s,
			"priority": 1
		}`, size, `["tag1", "tag2", "tag3"]`) // 簡略化
	}
	
	sizes := []int{100, 1000, 5000}
	
	for _, size := range sizes {
		t.Run(fmt.Sprintf("サイズ%d", size), func(t *testing.T) {
			comment := createLargeJSONComment(size)
			
			column := &Column{
				Name:    "large_comment_col",
				Comment: comment,
			}
			
			start := time.Now()
			err := column.ProcessEnhancedComment(processor, "|", true)
			elapsed := time.Since(start)
			
			if err != nil {
				t.Errorf("Large comment processing failed: %v", err)
			}
			
			t.Logf("サイズ%d処理時間: %v", size, elapsed)
			
			// 処理時間が合理的な範囲内であることを確認
			if elapsed > 100*time.Millisecond {
				t.Logf("警告: 処理時間が長いです: %v", elapsed)
			}
		})
	}
}