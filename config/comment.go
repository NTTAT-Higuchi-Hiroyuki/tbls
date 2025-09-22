package config

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

// Validate validates CommentConfig values.
func (c *CommentConfig) Validate() error {
	// Basic validation - separator should not be empty after setting defaults
	if c.Separator == "" {
		c.Separator = DefaultCommentSeparator
	}
	return nil
}