package schema

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestEnhancedCommentError 拡張エラーの基本機能テスト
func TestEnhancedCommentError(t *testing.T) {
	t.Run("基本エラー構築", func(t *testing.T) {
		err := NewErrorBuilder().
			WithMessage("テストエラー").
			WithSeverity(ErrorSeverityError).
			WithCategory(ErrorCategoryParsing).
			WithErrorCode("E_TEST_ERROR").
			WithObjectInfo(ObjectTypeColumn, "test_column").
			WithParserName("json").
			WithSourceComment(`{"invalid": json}`).
			WithSuggestion("JSON形式を確認してください").
			Build()

		if err.Message != "テストエラー" {
			t.Errorf("Message = %q, expected 'テストエラー'", err.Message)
		}

		if err.Severity != ErrorSeverityError {
			t.Errorf("Severity = %v, expected %v", err.Severity, ErrorSeverityError)
		}

		if err.Category != ErrorCategoryParsing {
			t.Errorf("Category = %v, expected %v", err.Category, ErrorCategoryParsing)
		}

		if err.ErrorCode != "E_TEST_ERROR" {
			t.Errorf("ErrorCode = %q, expected 'E_TEST_ERROR'", err.ErrorCode)
		}

		if err.ObjectType != ObjectTypeColumn {
			t.Errorf("ObjectType = %v, expected %v", err.ObjectType, ObjectTypeColumn)
		}

		if err.ObjectName != "test_column" {
			t.Errorf("ObjectName = %q, expected 'test_column'", err.ObjectName)
		}

		if len(err.Suggestions) != 1 {
			t.Errorf("Suggestions length = %d, expected 1", len(err.Suggestions))
		}

		expectedErrorString := "[ERROR:PARSING] E_TEST_ERROR テストエラー (object: test_column) (parser: json)"
		if err.Error() != expectedErrorString {
			t.Errorf("Error() = %q, expected %q", err.Error(), expectedErrorString)
		}
	})

	t.Run("コンテキスト付きエラー", func(t *testing.T) {
		err := NewErrorBuilder().
			WithMessage("コンテキストテスト").
			WithSeverity(ErrorSeverityWarning).
			WithCategory(ErrorCategoryValidation).
			WithErrorCode("E_CONTEXT_TEST").
			WithContext("line_number", 42).
			WithContext("column_position", 15).
			WithContext("attempted_value", "invalid_data").
			Build()

		if len(err.Context) != 3 {
			t.Errorf("Context length = %d, expected 3", len(err.Context))
		}

		if err.Context["line_number"] != 42 {
			t.Errorf("Context line_number = %v, expected 42", err.Context["line_number"])
		}
	})
}

// TestErrorSeverityAndCategory 重要度とカテゴリのテスト
func TestErrorSeverityAndCategory(t *testing.T) {
	severityTests := []struct {
		severity ErrorSeverity
		expected string
	}{
		{ErrorSeverityInfo, "INFO"},
		{ErrorSeverityWarning, "WARNING"},
		{ErrorSeverityError, "ERROR"},
		{ErrorSeverityCritical, "CRITICAL"},
	}

	for _, test := range severityTests {
		if test.severity.String() != test.expected {
			t.Errorf("Severity.String() = %q, expected %q", test.severity.String(), test.expected)
		}
	}

	categoryTests := []struct {
		category ErrorCategory
		expected string
	}{
		{ErrorCategoryValidation, "VALIDATION"},
		{ErrorCategoryParsing, "PARSING"},
		{ErrorCategoryProcessing, "PROCESSING"},
		{ErrorCategoryConfiguration, "CONFIGURATION"},
		{ErrorCategoryTimeout, "TIMEOUT"},
		{ErrorCategoryMemory, "MEMORY"},
		{ErrorCategorySchema, "SCHEMA"},
	}

	for _, test := range categoryTests {
		if test.category.String() != test.expected {
			t.Errorf("Category.String() = %q, expected %q", test.category.String(), test.expected)
		}
	}
}

// TestFallbackParsingStrategy フォールバック解析戦略のテスト
func TestFallbackParsingStrategy(t *testing.T) {
	strategy := &FallbackParsingStrategy{}

	t.Run("回復可能エラーの判定", func(t *testing.T) {
		// 回復可能なエラー
		recoverableError := NewErrorBuilder().
			WithCategory(ErrorCategoryParsing).
			WithErrorCode("E_JSON_PARSE_FAILED").
			Build()

		if !strategy.CanRecover(recoverableError) {
			t.Error("Should be able to recover from JSON parse failed error")
		}

		// 回復不可能なエラー
		nonRecoverableError := NewErrorBuilder().
			WithCategory(ErrorCategoryValidation).
			WithErrorCode("E_VALIDATION_FAILED").
			Build()

		if strategy.CanRecover(nonRecoverableError) {
			t.Error("Should not be able to recover from validation error")
		}
	})

	t.Run("エラー回復の実行", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCategory(ErrorCategoryParsing).
			WithErrorCode("E_JSON_PARSE_FAILED").
			Build()

		comment := "テスト論理名|テスト説明"
		result, recoveryErr := strategy.Recover(err, comment)

		if recoveryErr != nil {
			t.Fatalf("Recovery failed: %v", recoveryErr)
		}

		commentData, ok := result.(*CommentData)
		if !ok {
			t.Fatalf("Recovery result is not CommentData: %T", result)
		}

		if commentData.LogicalName != "テスト論理名" {
			t.Errorf("LogicalName = %q, expected 'テスト論理名'", commentData.LogicalName)
		}

		if commentData.Description != "テスト説明" {
			t.Errorf("Description = %q, expected 'テスト説明'", commentData.Description)
		}
	})

	t.Run("無効なコンテキストでの回復", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCategory(ErrorCategoryParsing).
			WithErrorCode("E_JSON_PARSE_FAILED").
			Build()

		invalidContext := 12345 // 文字列でない
		_, recoveryErr := strategy.Recover(err, invalidContext)

		if recoveryErr == nil {
			t.Error("Should fail to recover with invalid context type")
		}
	})
}

// TestCommentSanitizationStrategy コメントサニタイゼーション戦略のテスト
func TestCommentSanitizationStrategy(t *testing.T) {
	strategy := &CommentSanitizationStrategy{}

	t.Run("回復可能エラーの判定", func(t *testing.T) {
		recoverableError := NewErrorBuilder().
			WithCategory(ErrorCategoryValidation).
			WithErrorCode("E_UNSAFE_CONTENT").
			Build()

		if !strategy.CanRecover(recoverableError) {
			t.Error("Should be able to recover from unsafe content error")
		}
	})

	t.Run("サニタイゼーション実行", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCategory(ErrorCategoryValidation).
			WithErrorCode("E_UNSAFE_CONTENT").
			Build()

		unsafeData := &CommentData{
			LogicalName: "テスト<script>alert('xss')</script>カラム",
			Description: "説明javascript:void(0)",
			Tags:        []string{"タグ1", "タグ<script>2"},
			Metadata: map[string]string{
				"key<script>": "value</script>",
			},
		}

		result, recoveryErr := strategy.Recover(err, unsafeData)

		if recoveryErr != nil {
			t.Fatalf("Sanitization failed: %v", recoveryErr)
		}

		sanitized, ok := result.(*CommentData)
		if !ok {
			t.Fatalf("Sanitization result is not CommentData: %T", result)
		}

		if sanitized.LogicalName != "テストカラム" {
			t.Errorf("LogicalName = %q, expected 'テストカラム'", sanitized.LogicalName)
		}

		if sanitized.Description != "説明" {
			t.Errorf("Description = %q, expected '説明'", sanitized.Description)
		}

		if len(sanitized.Tags) != 2 || sanitized.Tags[1] != "タグ2" {
			t.Errorf("Tags = %v, expected sanitized tags", sanitized.Tags)
		}

		for key := range sanitized.Metadata {
			if key == "key<script>" {
				t.Error("Metadata key should be sanitized")
			}
		}
	})
}

// TestErrorRecoveryManager エラー回復管理のテスト
func TestErrorRecoveryManager(t *testing.T) {
	manager := NewErrorRecoveryManager()

	t.Run("利用可能戦略の取得", func(t *testing.T) {
		strategies := manager.GetAvailableStrategies()
		if len(strategies) < 2 {
			t.Errorf("Available strategies = %d, expected at least 2", len(strategies))
		}

		t.Logf("Available strategies: %v", strategies)
	})

	t.Run("回復可能エラーの処理", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCategory(ErrorCategoryParsing).
			WithErrorCode("E_JSON_PARSE_FAILED").
			Build()

		comment := "フォールバック論理名|フォールバック説明"
		result, recoveryErr := manager.TryRecover(err, comment)

		if recoveryErr != nil {
			t.Fatalf("Recovery failed: %v", recoveryErr)
		}

		if result == nil {
			t.Fatal("Recovery result should not be nil")
		}

		t.Logf("Recovery successful: %T", result)
	})

	t.Run("回復不可能エラーの処理", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCategory(ErrorCategoryTimeout).
			WithErrorCode("E_PROCESSING_TIMEOUT").
			Build()

		_, recoveryErr := manager.TryRecover(err, "context")

		if recoveryErr == nil {
			t.Error("Should fail to recover from timeout error")
		}

		t.Logf("Expected recovery failure: %v", recoveryErr)
	})
}

// TestDefaultErrorReporter デフォルトエラーレポーターのテスト
func TestDefaultErrorReporter(t *testing.T) {
	reporter := NewDefaultErrorReporter()

	t.Run("エラー報告", func(t *testing.T) {
		// 複数のエラーを報告
		errors := []*EnhancedCommentError{
			NewErrorBuilder().
				WithMessage("エラー1").
				WithSeverity(ErrorSeverityError).
				WithCategory(ErrorCategoryParsing).
				Build(),
			NewErrorBuilder().
				WithMessage("警告1").
				WithSeverity(ErrorSeverityWarning).
				WithCategory(ErrorCategoryValidation).
				Build(),
			NewErrorBuilder().
				WithMessage("情報1").
				WithSeverity(ErrorSeverityInfo).
				WithCategory(ErrorCategoryProcessing).
				Build(),
		}

		for _, err := range errors {
			reporter.ReportError(err)
		}

		summary := reporter.GetErrorSummary()
		if summary.TotalErrors != 3 {
			t.Errorf("TotalErrors = %d, expected 3", summary.TotalErrors)
		}

		if summary.ErrorsBySeverity[ErrorSeverityError] != 1 {
			t.Errorf("Error count = %d, expected 1", summary.ErrorsBySeverity[ErrorSeverityError])
		}

		if summary.ErrorsByCategory[ErrorCategoryParsing] != 1 {
			t.Errorf("Parsing error count = %d, expected 1", summary.ErrorsByCategory[ErrorCategoryParsing])
		}
	})

	t.Run("回復報告", func(t *testing.T) {
		originalErr := NewErrorBuilder().
			WithMessage("回復対象エラー").
			WithSeverity(ErrorSeverityError).
			WithCategory(ErrorCategoryParsing).
			WithObjectInfo(ObjectTypeColumn, "test_col").
			Build()

		recoveredData := &CommentData{LogicalName: "回復済み"}

		reporter.ReportRecovery(originalErr, recoveredData)

		summary := reporter.GetErrorSummary()
		if summary.RecoveredErrors != 1 {
			t.Errorf("RecoveredErrors = %d, expected 1", summary.RecoveredErrors)
		}

		// 回復成功の情報エラーが追加されているかチェック
		found := false
		for _, err := range summary.RecentErrors {
			if err.ErrorCode == "E_RECOVERY_SUCCESS" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Recovery success error should be reported")
		}
	})

	t.Run("最新エラー取得", func(t *testing.T) {
		// 15個のエラーを追加（10個を超える）
		for i := 0; i < 15; i++ {
			err := NewErrorBuilder().
				WithMessage(fmt.Sprintf("エラー%d", i)).
				WithSeverity(ErrorSeverityInfo).
				WithCategory(ErrorCategoryProcessing).
				Build()
			reporter.ReportError(err)
		}

		summary := reporter.GetErrorSummary()
		
		// 最新10件まで取得されることを確認
		if len(summary.RecentErrors) > 10 {
			t.Errorf("RecentErrors length = %d, expected <= 10", len(summary.RecentErrors))
		}
	})
}

// TestPredefinedErrorCreators 事前定義エラー作成関数のテスト
func TestPredefinedErrorCreators(t *testing.T) {
	t.Run("パーシングエラー作成", func(t *testing.T) {
		innerErr := fmt.Errorf("invalid JSON syntax")
		err := NewParsingError("json", `{"invalid": json}`, innerErr)

		if err.Category != ErrorCategoryParsing {
			t.Errorf("Category = %v, expected %v", err.Category, ErrorCategoryParsing)
		}

		if err.ErrorCode != ErrorCodeJSONParseFailed {
			t.Errorf("ErrorCode = %q, expected %q", err.ErrorCode, ErrorCodeJSONParseFailed)
		}

		if err.ParserName != "json" {
			t.Errorf("ParserName = %q, expected 'json'", err.ParserName)
		}

		if err.InnerError != innerErr {
			t.Error("InnerError should be preserved")
		}

		if len(err.Suggestions) == 0 {
			t.Error("Suggestions should be provided")
		}
	})

	t.Run("バリデーションエラー作成", func(t *testing.T) {
		err := NewValidationError("無効な値です", ObjectTypeTable, "test_table")

		if err.Category != ErrorCategoryValidation {
			t.Errorf("Category = %v, expected %v", err.Category, ErrorCategoryValidation)
		}

		if err.Severity != ErrorSeverityWarning {
			t.Errorf("Severity = %v, expected %v", err.Severity, ErrorSeverityWarning)
		}

		if err.ObjectType != ObjectTypeTable {
			t.Errorf("ObjectType = %v, expected %v", err.ObjectType, ObjectTypeTable)
		}

		if err.ObjectName != "test_table" {
			t.Errorf("ObjectName = %q, expected 'test_table'", err.ObjectName)
		}
	})

	t.Run("タイムアウトエラー作成", func(t *testing.T) {
		timeout := 5 * time.Second
		err := NewProcessingTimeoutError(timeout, ObjectTypeColumn, "test_col")

		if err.Category != ErrorCategoryTimeout {
			t.Errorf("Category = %v, expected %v", err.Category, ErrorCategoryTimeout)
		}

		if err.Severity != ErrorSeverityError {
			t.Errorf("Severity = %v, expected %v", err.Severity, ErrorSeverityError)
		}

		if err.Context["timeout_duration"] != timeout.String() {
			t.Errorf("Timeout context = %v, expected %v", err.Context["timeout_duration"], timeout.String())
		}
	})

	t.Run("設定エラー作成", func(t *testing.T) {
		err := NewConfigurationError("無効な設定値", "enhanced_comment.enabled")

		if err.Category != ErrorCategoryConfiguration {
			t.Errorf("Category = %v, expected %v", err.Category, ErrorCategoryConfiguration)
		}

		if err.Severity != ErrorSeverityCritical {
			t.Errorf("Severity = %v, expected %v", err.Severity, ErrorSeverityCritical)
		}

		if err.Context["config_key"] != "enhanced_comment.enabled" {
			t.Errorf("Config key = %v, expected 'enhanced_comment.enabled'", err.Context["config_key"])
		}
	})
}

// TestErrorHandlingIntegration エラーハンドリング統合テスト
func TestErrorHandlingIntegration(t *testing.T) {
	reporter := NewDefaultErrorReporter()
	manager := NewErrorRecoveryManager()

	t.Run("完全なエラーハンドリングフロー", func(t *testing.T) {
		// 1. エラー発生
		originalErr := NewParsingError("json", `{"invalid": json}`, fmt.Errorf("syntax error"))
		reporter.ReportError(originalErr)

		// 2. 回復試行
		comment := "フォールバック論理名|フォールバック説明"
		recoveredData, recoveryErr := manager.TryRecover(originalErr, comment)

		if recoveryErr != nil {
			t.Fatalf("Recovery should succeed: %v", recoveryErr)
		}

		// 3. 回復成功報告
		reporter.ReportRecovery(originalErr, recoveredData)

		// 4. 最終結果確認
		summary := reporter.GetErrorSummary()
		if summary.TotalErrors < 2 { // 元のエラー + 回復成功情報
			t.Errorf("TotalErrors = %d, expected at least 2", summary.TotalErrors)
		}

		if summary.RecoveredErrors != 1 {
			t.Errorf("RecoveredErrors = %d, expected 1", summary.RecoveredErrors)
		}

		// 回復されたデータの検証
		commentData, ok := recoveredData.(*CommentData)
		if !ok {
			t.Fatalf("Recovered data is not CommentData: %T", recoveredData)
		}

		if commentData.LogicalName != "フォールバック論理名" {
			t.Errorf("Recovered LogicalName = %q, expected 'フォールバック論理名'", commentData.LogicalName)
		}

		t.Logf("Error handling flow completed successfully")
		t.Logf("Final summary: %+v", summary)
	})
}

// TestSanitizationFunctions サニタイゼーション関数のテスト
func TestSanitizationFunctions(t *testing.T) {
	t.Run("文字列サニタイゼーション", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"普通の文字列", "普通の文字列"},
			{"<script>alert('xss')</script>", "alert('xss')"},
			{"javascript:void(0)", ""},
			{"前<script>中</script>後", "前中後"},
		}

		for _, test := range tests {
			result := sanitizeString(test.input)
			if result != test.expected {
				t.Errorf("sanitizeString(%q) = %q, expected %q", test.input, result, test.expected)
			}
		}
	})

	t.Run("文字列スライスサニタイゼーション", func(t *testing.T) {
		input := []string{"正常", "<script>悪意</script>", "javascript:void(0)"}
		expected := []string{"正常", "悪意", ""}

		result := sanitizeStringSlice(input)
		if len(result) != len(expected) {
			t.Fatalf("Result length = %d, expected %d", len(result), len(expected))
		}

		for i, exp := range expected {
			if result[i] != exp {
				t.Errorf("Result[%d] = %q, expected %q", i, result[i], exp)
			}
		}
	})

	t.Run("メタデータサニタイゼーション", func(t *testing.T) {
		input := map[string]string{
			"正常キー":             "正常値",
			"key<script>":       "value</script>",
			"javascript:bad":    "javascript:void(0)",
		}

		result := sanitizeMetadata(input)
		
		if len(result) != len(input) {
			t.Fatalf("Result length = %d, expected %d", len(result), len(input))
		}

		// 危険なスクリプトが除去されていることを確認
		for key, value := range result {
			if strings.Contains(key, "<script>") || strings.Contains(value, "</script>") {
				t.Errorf("Unsafe content found in result: %s = %s", key, value)
			}
		}
	})
}