package pkg

import (
	"strings"
	"unicode"
)

type Field struct {
	SourceName Name
	DestName   Name
	Counter    int
	SourceType string
	DestType   string
	Tags       []map[string]string
	Comment    string
}

type Name string

func (fn *Name) Unexported() string {
	runes := []rune(*fn)
	if len(runes) == 0 {
		return ""
	}
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func (fn *Name) Exported() string {
	runes := []rune(*fn)
	if len(runes) == 0 {
		return ""
	}
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func (fn *Name) SnakeCase() string {
	var sb strings.Builder
	for i, r := range *fn {
		if unicode.IsUpper(r) {
			if i > 0 {
				sb.WriteRune('_')
			}
			sb.WriteRune(unicode.ToLower(r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

type TemplateContext struct {
	TemplateName string
	SourceFile   string
	StructName   Name
	Fields       []Field
	Comments     []string
	Config       StructConfig
}
