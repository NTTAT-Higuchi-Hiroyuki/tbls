package config

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// DefaultCommentSeparator is the default separator for logical name in comment.
const DefaultCommentSeparator = ":"

// CommentConfig is the configuration for comment processing.
type CommentConfig struct {
	// Separator for logical name extraction from comment
	Separator string `yaml:"separator,omitempty"`
}

// setDefaultComment sets default values for CommentConfig.
func (c *CommentConfig) setDefaultComment() {
	if c.Separator == "" {
		c.Separator = DefaultCommentSeparator
	}
}

// Validate validates CommentConfig values with enhanced validation.
func (c *CommentConfig) Validate() error {
	if c == nil {
		return nil
	}

	// Set default if empty
	if c.Separator == "" {
		c.Separator = DefaultCommentSeparator
		fmt.Printf("警告: コメント区切り文字が設定されていません。デフォルト値 '%s' を使用します。\n", DefaultCommentSeparator)
		return nil
	}

	// Validate separator length (should not be too long)
	if len(c.Separator) > 10 {
		fmt.Printf("警告: コメント区切り文字が長すぎます（%d文字）。デフォルト値 '%s' を使用します。\n", len(c.Separator), DefaultCommentSeparator)
		c.Separator = DefaultCommentSeparator
		return nil
	}

	// Validate UTF-8 validity
	if !utf8.ValidString(c.Separator) {
		fmt.Printf("警告: コメント区切り文字に無効なUTF-8文字が含まれています。デフォルト値 '%s' を使用します。\n", DefaultCommentSeparator)
		c.Separator = DefaultCommentSeparator
		return nil
	}

	// Check for problematic characters that might cause parsing issues
	problematicChars := []string{"\n", "\r", "\t", "\"", "'", "`"}
	for _, char := range problematicChars {
		if strings.Contains(c.Separator, char) {
			fmt.Printf("警告: コメント区切り文字に問題のある文字が含まれています。デフォルト値 '%s' を使用します。\n", DefaultCommentSeparator)
			c.Separator = DefaultCommentSeparator
			return nil
		}
	}

	// Success validation message for debug
	if c.Separator != DefaultCommentSeparator {
		fmt.Printf("情報: カスタムコメント区切り文字 '%s' を使用します。\n", c.Separator)
	}

	return nil
}

// IsValid checks if the CommentConfig is valid without modifying it.
func (c *CommentConfig) IsValid() bool {
	if c == nil {
		return true // nil config is considered valid with defaults
	}

	if c.Separator == "" {
		return true // empty separator will use default
	}

	// Check length
	if len(c.Separator) > 10 {
		return false
	}

	// Check UTF-8 validity
	if !utf8.ValidString(c.Separator) {
		return false
	}

	// Check for problematic characters
	problematicChars := []string{"\n", "\r", "\t", "\"", "'", "`"}
	for _, char := range problematicChars {
		if strings.Contains(c.Separator, char) {
			return false
		}
	}

	return true
}
