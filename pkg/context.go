package pkg

type Field struct {
	SourceName string
	DestName   string
	SourceType string
	DestType   string
	Tags       []map[string]string
	Comment    string
}

type TemplateContext struct {
	TemplateName string
	SourceFile   string
	StructName   string
	Fields       []Field
	Comments     []string
	Config       StructConfig
}
