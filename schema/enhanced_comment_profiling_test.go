package schema

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

// PerformanceMetrics パフォーマンス測定結果
type PerformanceMetrics struct {
	ProcessingTime    time.Duration
	AllocatedMemory   uint64
	AllocationsCount  uint64
	GCPauses          uint64
	ObjectsProcessed  int
	ThroughputPerSec  float64
}

// measurePerformance パフォーマンス測定ヘルパー
func measurePerformance(name string, objectCount int, fn func()) *PerformanceMetrics {
	runtime.GC()
	
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	
	start := time.Now()
	fn()
	elapsed := time.Since(start)
	
	runtime.ReadMemStats(&m2)
	
	metrics := &PerformanceMetrics{
		ProcessingTime:    elapsed,
		AllocatedMemory:   m2.TotalAlloc - m1.TotalAlloc,
		AllocationsCount:  m2.Mallocs - m1.Mallocs,
		GCPauses:          m2.PauseTotalNs - m1.PauseTotalNs,
		ObjectsProcessed:  objectCount,
		ThroughputPerSec:  float64(objectCount) / elapsed.Seconds(),
	}
	
	return metrics
}

// TestDetailedPerformanceAnalysis 詳細パフォーマンス分析
func TestDetailedPerformanceAnalysis(t *testing.T) {
	processor := NewEnhancedCommentProcessor()
	
	testCases := []struct {
		name         string
		tableCount   int
		columnCount  int
		commentType  string
	}{
		{"小規模JSON", 10, 5, "json"},
		{"中規模JSON", 50, 10, "json"},
		{"大規模JSON", 100, 20, "json"},
		{"小規模YAML", 10, 5, "yaml"},
		{"中規模YAML", 50, 10, "yaml"},
		{"大規模YAML", 100, 20, "yaml"},
		{"小規模Legacy", 10, 5, "legacy"},
		{"中規模Legacy", 50, 10, "legacy"},
		{"大規模Legacy", 100, 20, "legacy"},
	}
	
	results := make(map[string]*PerformanceMetrics)
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			schema := createTypedSchema(tc.tableCount, tc.columnCount, tc.commentType)
			objectCount := tc.tableCount + (tc.tableCount * tc.columnCount)
			
			metrics := measurePerformance(tc.name, objectCount, func() {
				for _, table := range schema.Tables {
					table.ProcessEnhancedComment(processor, "|", true)
					for _, column := range table.Columns {
						column.ProcessEnhancedComment(processor, "|", true)
					}
				}
			})
			
			results[tc.name] = metrics
			
			t.Logf("=== %s パフォーマンス結果 ===", tc.name)
			t.Logf("処理時間: %v", metrics.ProcessingTime)
			t.Logf("メモリ使用量: %d bytes (%.2f MB)", metrics.AllocatedMemory, float64(metrics.AllocatedMemory)/1024/1024)
			t.Logf("メモリ割り当て回数: %d", metrics.AllocationsCount)
			t.Logf("GC停止時間: %d ns", metrics.GCPauses)
			t.Logf("処理オブジェクト数: %d", metrics.ObjectsProcessed)
			t.Logf("スループット: %.2f objects/sec", metrics.ThroughputPerSec)
			t.Logf("オブジェクト当たり処理時間: %.2f µs", float64(metrics.ProcessingTime.Nanoseconds())/float64(metrics.ObjectsProcessed)/1000)
			
			// パフォーマンス閾値チェック
			checkPerformanceThresholds(t, tc.name, metrics)
		})
	}
	
	// 形式別比較分析
	analyzeFormatComparison(t, results)
}

// createTypedSchema 指定された形式のコメントを持つスキーマ生成
func createTypedSchema(tableCount, columnCount int, commentType string) *Schema {
	tables := make([]*Table, tableCount)
	
	for i := 0; i < tableCount; i++ {
		columns := make([]*Column, columnCount)
		for j := 0; j < columnCount; j++ {
			var comment string
			switch commentType {
			case "json":
				comment = fmt.Sprintf(`{"name": "カラム%d_%d", "description": "JSONテスト用カラム", "tags": ["test", "json"], "priority": %d}`, i, j, j%5+1)
			case "yaml":
				comment = fmt.Sprintf("name: カラム%d_%d\ndescription: YAMLテスト用カラム\ntags:\n  - test\n  - yaml\npriority: %d", i, j, j%5+1)
			case "legacy":
				comment = fmt.Sprintf("カラム%d_%d|Legacyテスト用カラム", i, j)
			}
			
			columns[j] = &Column{
				Name:    fmt.Sprintf("col_%d_%d", i, j),
				Comment: comment,
			}
		}
		
		var tableComment string
		switch commentType {
		case "json":
			tableComment = fmt.Sprintf(`{"name": "テーブル%d", "description": "JSONテスト用テーブル", "tags": ["table", "json"]}`, i)
		case "yaml":
			tableComment = fmt.Sprintf("name: テーブル%d\ndescription: YAMLテスト用テーブル\ntags:\n  - table\n  - yaml", i)
		case "legacy":
			tableComment = fmt.Sprintf("テーブル%d|Legacyテスト用テーブル", i)
		}
		
		tables[i] = &Table{
			Name:    fmt.Sprintf("table_%d", i),
			Comment: tableComment,
			Columns: columns,
		}
	}
	
	return &Schema{Name: fmt.Sprintf("%s_test_db", commentType), Tables: tables}
}

// checkPerformanceThresholds パフォーマンス閾値チェック
func checkPerformanceThresholds(t *testing.T, testName string, metrics *PerformanceMetrics) {
	// 基本的な性能閾値（実環境に応じて調整）
	thresholds := map[string]struct {
		maxProcessingTimePerObject time.Duration
		maxMemoryPerObject         uint64
		minThroughput              float64
	}{
		// JSON形式の閾値
		"小規模JSON": {maxProcessingTimePerObject: 100 * time.Microsecond, maxMemoryPerObject: 10000, minThroughput: 1000},
		"中規模JSON": {maxProcessingTimePerObject: 150 * time.Microsecond, maxMemoryPerObject: 15000, minThroughput: 500},
		"大規模JSON": {maxProcessingTimePerObject: 200 * time.Microsecond, maxMemoryPerObject: 20000, minThroughput: 300},
		
		// YAML形式の閾値（JSONより処理が重い）
		"小規模YAML": {maxProcessingTimePerObject: 200 * time.Microsecond, maxMemoryPerObject: 15000, minThroughput: 500},
		"中規模YAML": {maxProcessingTimePerObject: 300 * time.Microsecond, maxMemoryPerObject: 20000, minThroughput: 300},
		"大規模YAML": {maxProcessingTimePerObject: 400 * time.Microsecond, maxMemoryPerObject: 25000, minThroughput: 200},
		
		// Legacy形式の閾値（最も軽い）
		"小規模Legacy": {maxProcessingTimePerObject: 10 * time.Microsecond, maxMemoryPerObject: 1000, minThroughput: 10000},
		"中規模Legacy": {maxProcessingTimePerObject: 15 * time.Microsecond, maxMemoryPerObject: 1500, minThroughput: 5000},
		"大規模Legacy": {maxProcessingTimePerObject: 20 * time.Microsecond, maxMemoryPerObject: 2000, minThroughput: 3000},
	}
	
	threshold, exists := thresholds[testName]
	if !exists {
		t.Logf("警告: %s の性能閾値が未定義", testName)
		return
	}
	
	// オブジェクト当たりの処理時間チェック
	avgProcessingTime := metrics.ProcessingTime / time.Duration(metrics.ObjectsProcessed)
	if avgProcessingTime > threshold.maxProcessingTimePerObject {
		t.Logf("警告: 処理時間が閾値を超過 - 実測: %v, 閾値: %v", avgProcessingTime, threshold.maxProcessingTimePerObject)
	}
	
	// オブジェクト当たりのメモリ使用量チェック
	avgMemoryUsage := metrics.AllocatedMemory / uint64(metrics.ObjectsProcessed)
	if avgMemoryUsage > threshold.maxMemoryPerObject {
		t.Logf("警告: メモリ使用量が閾値を超過 - 実測: %d bytes, 閾値: %d bytes", avgMemoryUsage, threshold.maxMemoryPerObject)
	}
	
	// スループットチェック
	if metrics.ThroughputPerSec < threshold.minThroughput {
		t.Logf("警告: スループットが閾値を下回る - 実測: %.2f ops/sec, 閾値: %.2f ops/sec", metrics.ThroughputPerSec, threshold.minThroughput)
	}
}

// analyzeFormatComparison 形式別比較分析
func analyzeFormatComparison(t *testing.T, results map[string]*PerformanceMetrics) {
	t.Log("\n=== 形式別パフォーマンス比較分析 ===")
	
	scales := []string{"小規模", "中規模", "大規模"}
	formats := []string{"JSON", "YAML", "Legacy"}
	
	for _, scale := range scales {
		t.Logf("\n--- %s比較 ---", scale)
		
		baselineKey := fmt.Sprintf("%sLegacy", scale)
		baseline := results[baselineKey]
		
		if baseline == nil {
			continue
		}
		
		for _, format := range formats {
			key := fmt.Sprintf("%s%s", scale, format)
			metrics := results[key]
			if metrics == nil {
				continue
			}
			
			processingRatio := float64(metrics.ProcessingTime) / float64(baseline.ProcessingTime)
			memoryRatio := float64(metrics.AllocatedMemory) / float64(baseline.AllocatedMemory)
			throughputRatio := metrics.ThroughputPerSec / baseline.ThroughputPerSec
			
			t.Logf("%s: 処理時間比 %.2fx, メモリ比 %.2fx, スループット比 %.2fx",
				format, processingRatio, memoryRatio, throughputRatio)
		}
	}
}

// TestMemoryLeakDetection メモリリーク検出テスト
func TestMemoryLeakDetection(t *testing.T) {
	processor := NewEnhancedCommentProcessor()
	
	// 繰り返し処理でのメモリリーク検出
	iterations := []int{100, 500, 1000}
	memoryUsages := make([]uint64, len(iterations))
	
	for i, iteration := range iterations {
		t.Run(fmt.Sprintf("反復回数_%d", iteration), func(t *testing.T) {
			runtime.GC()
			var m1, m2 runtime.MemStats
			runtime.ReadMemStats(&m1)
			
			column := &Column{
				Name:    "leak_test_col",
				Comment: `{"name": "リークテスト", "description": "メモリリーク検出用テスト", "tags": ["test", "leak"]}`,
			}
			
			for j := 0; j < iteration; j++ {
				column.ProcessEnhancedComment(processor, "|", true)
			}
			
			runtime.GC()
			runtime.ReadMemStats(&m2)
			
			memoryIncrease := m2.TotalAlloc - m1.TotalAlloc
			memoryUsages[i] = memoryIncrease
			
			t.Logf("反復回数 %d: メモリ増加量 %d bytes", iteration, memoryIncrease)
		})
	}
	
	// メモリ使用量の増加傾向分析
	if len(memoryUsages) >= 2 {
		for i := 1; i < len(memoryUsages); i++ {
			ratio := float64(memoryUsages[i]) / float64(memoryUsages[i-1])
			expectedRatio := float64(iterations[i]) / float64(iterations[i-1])
			
			t.Logf("メモリ増加比 %.2f (期待値: %.2f)", ratio, expectedRatio)
			
			// メモリ増加が線形でない場合（リーク疑い）
			if ratio > expectedRatio*1.5 {
				t.Logf("警告: メモリリークの可能性があります - 実測比: %.2f, 期待比: %.2f", ratio, expectedRatio)
			}
		}
	}
}

// TestCPUProfileAnalysis CPU使用率分析テスト
func TestCPUProfileAnalysis(t *testing.T) {
	processor := NewEnhancedCommentProcessor()
	
	testCases := []struct {
		name        string
		commentType string
		complexity  string
	}{
		{"シンプルJSON", "json", "simple"},
		{"複雑JSON", "json", "complex"},
		{"シンプルYAML", "yaml", "simple"},
		{"複雑YAML", "yaml", "complex"},
		{"シンプルLegacy", "legacy", "simple"},
		{"複雑Legacy", "legacy", "complex"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			comment := generateComment(tc.commentType, tc.complexity)
			column := &Column{Name: "cpu_test_col", Comment: comment}
			
			// CPU時間測定
			iterations := 1000
			start := time.Now()
			
			for i := 0; i < iterations; i++ {
				column.ProcessEnhancedComment(processor, "|", true)
			}
			
			elapsed := time.Since(start)
			avgCPUTime := elapsed / time.Duration(iterations)
			
			t.Logf("%s: 平均CPU時間 %v (1000回反復)", tc.name, avgCPUTime)
			
			// CPU使用効率の基本チェック
			if avgCPUTime > 1*time.Millisecond {
				t.Logf("警告: CPU使用時間が長いです: %v", avgCPUTime)
			}
		})
	}
}

// generateComment 指定された複雑さのコメント生成
func generateComment(commentType, complexity string) string {
	switch commentType {
	case "json":
		if complexity == "complex" {
			return `{
				"name": "複雑なJSONコメント",
				"description": "非常に詳細で長い説明文を含む複雑なJSONコメントです。この説明文は複数の文から構成され、様々な情報を含んでいます。処理時間やメモリ使用量の測定に使用されます。",
				"tags": ["複雑", "JSON", "テスト", "パフォーマンス", "分析", "詳細", "測定"],
				"metadata": {
					"author": "test_user",
					"created_at": "2023-01-01T00:00:00Z",
					"version": "1.0.0",
					"category": "test",
					"subcategory": "performance",
					"importance": "high",
					"review_status": "approved"
				},
				"priority": 10,
				"deprecated": false,
				"validation_rules": ["required", "unique", "indexed"]
			}`
		}
		return `{"name": "シンプルJSON", "description": "シンプルなテスト", "tags": ["test"]}`
		
	case "yaml":
		if complexity == "complex" {
			return `name: 複雑なYAMLコメント
description: |
  非常に詳細で長い説明文を含む複雑なYAMLコメントです。
  この説明文は複数行にわたって記述され、様々な情報を含んでいます。
  処理時間やメモリ使用量の測定に使用されます。
tags:
  - 複雑
  - YAML
  - テスト
  - パフォーマンス
  - 分析
  - 詳細
  - 測定
metadata:
  author: test_user
  created_at: "2023-01-01T00:00:00Z"
  version: "1.0.0"
  category: test
  subcategory: performance
  importance: high
  review_status: approved
priority: 10
deprecated: false
validation_rules:
  - required
  - unique
  - indexed`
		}
		return "name: シンプルYAML\ndescription: シンプルなテスト\ntags:\n  - test"
		
	case "legacy":
		if complexity == "complex" {
			return "複雑なLegacyコメント|非常に詳細で長い説明文を含む複雑なLegacyコメントです。この説明文は処理時間やメモリ使用量の測定に使用され、様々な特殊文字や記号を含んでいます。"
		}
		return "シンプルLegacy|シンプルなテスト"
		
	default:
		return "デフォルトコメント"
	}
}

// TestGCImpactAnalysis ガベージコレクション影響分析
func TestGCImpactAnalysis(t *testing.T) {
	processor := NewEnhancedCommentProcessor()
	
	// GC前後でのパフォーマンス測定
	schema := createLargeSchema(100, 15)
	
	// GC強制実行後の測定
	runtime.GC()
	runtime.GC() // 2回実行して安定化
	
	var beforeGC, afterGC runtime.MemStats
	runtime.ReadMemStats(&beforeGC)
	
	start := time.Now()
	for _, table := range schema.Tables {
		table.ProcessEnhancedComment(processor, "|", true)
		for _, column := range table.Columns {
			column.ProcessEnhancedComment(processor, "|", true)
		}
	}
	processingTime := time.Since(start)
	
	runtime.ReadMemStats(&afterGC)
	
	t.Logf("=== GC影響分析 ===")
	t.Logf("処理時間: %v", processingTime)
	t.Logf("メモリ割り当て: %d bytes", afterGC.TotalAlloc-beforeGC.TotalAlloc)
	t.Logf("GC実行回数: %d", afterGC.NumGC-beforeGC.NumGC)
	t.Logf("GC停止時間合計: %d ns", afterGC.PauseTotalNs-beforeGC.PauseTotalNs)
	
	if afterGC.NumGC > beforeGC.NumGC {
		avgGCPause := (afterGC.PauseTotalNs - beforeGC.PauseTotalNs) / uint64(afterGC.NumGC-beforeGC.NumGC)
		t.Logf("平均GC停止時間: %d ns", avgGCPause)
		
		// GC停止時間が処理時間に与える影響分析
		gcImpactRatio := float64(afterGC.PauseTotalNs-beforeGC.PauseTotalNs) / float64(processingTime.Nanoseconds()) * 100
		t.Logf("GCによる影響: %.2f%% of total processing time", gcImpactRatio)
		
		if gcImpactRatio > 10 {
			t.Logf("警告: GCの影響が大きいです (%.2f%%)", gcImpactRatio)
		}
	}
}