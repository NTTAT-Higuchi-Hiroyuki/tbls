package config

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestCommentConfigDefault(t *testing.T) {
	c := &CommentConfig{}
	c.setDefaultComment()

	if want := DefaultCommentSeparator; c.Separator != want {
		t.Errorf("got %v\nwant %v", c.Separator, want)
	}
}

func TestCommentConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *CommentConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &CommentConfig{
				Separator: ":",
			},
			wantErr: false,
		},
		{
			name: "empty separator gets default",
			config: &CommentConfig{
				Separator: "",
			},
			wantErr: false,
		},
		{
			name: "custom separator",
			config: &CommentConfig{
				Separator: "|",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CommentConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			// After validation, separator should never be empty
			if tt.config.Separator == "" {
				t.Errorf("CommentConfig.Separator should not be empty after validation")
			}
		})
	}
}

func TestConfigWithCommentField(t *testing.T) {
	// Test that CommentConfig is properly initialized in Config
	config, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if config.Comment == nil {
		t.Error("Config.Comment should be initialized")
	}

	if want := DefaultCommentSeparator; config.Comment.Separator != want {
		t.Errorf("got %v\nwant %v", config.Comment.Separator, want)
	}
}

func TestConfigCommentYAMLUnmarshal(t *testing.T) {
	yamlData := `
name: test
comment:
  separator: "|"
`
	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatal(err)
	}

	if config.Comment == nil {
		t.Error("Config.Comment should be loaded from YAML")
	}

	if want := "|"; config.Comment.Separator != want {
		t.Errorf("got %v\nwant %v", config.Comment.Separator, want)
	}
}

func TestConfigCommentYAMLMarshal(t *testing.T) {
	config := &Config{
		Name: "test",
		Comment: &CommentConfig{
			Separator: "|",
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the comment section is in the YAML output
	yamlStr := string(data)
	if !strings.Contains(yamlStr, "comment:") {
		t.Error("YAML output should contain comment section")
	}
	// Check for pipe character in different formats
	if !strings.Contains(yamlStr, "separator: \"|\"") && !strings.Contains(yamlStr, "separator: '|'") && !strings.Contains(yamlStr, "separator: |") {
		t.Errorf("YAML output should contain separator field with pipe character, got: %s", yamlStr)
	}
}

func TestConfigWithEmptyCommentConfig(t *testing.T) {
	// Test loading empty config file
	configFilepath := filepath.Join(testdataDir(), "empty.yml")
	config, err := New()
	if err != nil {
		t.Fatal(err)
	}
	err = config.Load(configFilepath)
	if err != nil {
		t.Fatal(err)
	}

	// Comment should be initialized with defaults
	if config.Comment == nil {
		t.Error("Config.Comment should be initialized even with empty config file")
	}

	if want := DefaultCommentSeparator; config.Comment.Separator != want {
		t.Errorf("got %v\nwant %v", config.Comment.Separator, want)
	}
}
