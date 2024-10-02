package pkg

type RemoveFieldsTransformer struct {
	Enable bool     `yaml:"enable"`
	Fields []string `yaml:"fields"`
}

func (t *RemoveFieldsTransformer) Transform(c *TemplateContext) error {
	var newFields []Field
	for i := range c.Fields {
		if !t.contains(t.Fields, c.Fields[i].Name) {
			newFields = append(newFields, c.Fields[i])
		}
	}
	c.Fields = newFields
	return nil
}

func (t *RemoveFieldsTransformer) contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
