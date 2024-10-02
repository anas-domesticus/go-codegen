package pkg

type FieldTypeAndNameTransformer struct{}

func (t *FieldTypeAndNameTransformer) Transform(ctx *TemplateContext) error {
	for i := range ctx.Fields {
		val, ok := t.hasTag("codegen-new-type", ctx.Fields[i].Tags)
		if ok {
			ctx.Fields[i].DestType = val
		}
		val, ok = t.hasTag("codegen-new-name", ctx.Fields[i].Tags)
		if ok {
			ctx.Fields[i].DestName = val
		}
	}
	return nil
}

func (t *FieldTypeAndNameTransformer) hasTag(tagName string, tags []map[string]string) (string, bool) {
	for _, tag := range tags {
		v, ok := tag[tagName]
		if ok {
			return v, true
		}
	}
	return "", false
}
