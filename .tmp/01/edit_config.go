#!/bin/bash
# Edit config.go to add logical name parsing in ModifySchema function

# First, let's read the content around line 566 where we need to add the logical name processing
sed -i '' '566i\
\
	// Apply logical name parsing if comment separator is configured\
	if c.Comment != nil && c.Comment.Separator != "" {\
		processor := schema.NewLogicalNameProcessor(c.Comment.Separator)\
		if err := processor.ProcessSchema(s); err != nil {\
			// Log warning but continue processing to avoid breaking existing functionality\
			// This ensures backward compatibility\
			fmt.Printf("Warning: Failed to process logical names: %v\\n", err)\
		}\
	}\
' /Users/higu/Project/tbls/config/config.go