package pkg

import (
	"reflect"
	"testing"
)

func TestParseCommentAndLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		comment  string
		expected StructConfig
	}{
		{
			name:    "EmptyComment",
			comment: "",
			expected: StructConfig{
				Values: make(map[string]string),
			},
		},
		{
			name:    "SimpleFlag",
			comment: "// @codegen flag",
			expected: StructConfig{
				Values: make(map[string]string),
				Flags:  []string{"flag"},
			},
		},
		{
			name:    "MultipleFlags",
			comment: "// @codegen flag1 flag2",
			expected: StructConfig{
				Values: make(map[string]string),
				Flags:  []string{"flag1", "flag2"},
			},
		},
		{
			name:    "SingleKeyValue",
			comment: "// @codegen key=value",
			expected: StructConfig{
				Values: map[string]string{"key": "value"},
			},
		},
		{
			name:    "MultipleKeyValues",
			comment: "// @codegen key1=value1 key2=value2",
			expected: StructConfig{
				Values: map[string]string{"key1": "value1", "key2": "value2"},
			},
		},
		{
			name:    "MixedFlagsAndKeyValues",
			comment: "// @codegen flag key=value",
			expected: StructConfig{
				Values: map[string]string{"key": "value"},
				Flags:  []string{"flag"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCommentAndLoadConfig(tt.comment)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("parseCommentAndLoadConfig() = %v, want %v", got, tt.expected)
			}
		})
	}
}
