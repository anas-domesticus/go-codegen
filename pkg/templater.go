package pkg

import (
	"fmt"
	"os"
	"text/template"
)

type Templater struct {
	configs      []Config
	Transformers map[string]TransformerFn
}

type TransformerFn func(ctx *TemplateContext) error

func NewTemplaterFromPath(path string) (*Templater, error) {
	cfgs, err := loadYAMLConfigs(path)
	if err != nil {
		return nil, err
	}
	return NewTemplater(cfgs), nil
}

func NewTemplater(configs []Config) *Templater {
	return &Templater{
		configs: configs,
	}
}

func (t *Templater) AddTransformer(name string, fn TransformerFn) {
	if t.Transformers == nil {
		t.Transformers = make(map[string]TransformerFn)
	}
	t.Transformers[name] = fn
}

func (t *Templater) RunTransformers() {
	for i := range t.configs {
		if t.configs[i].RemoveFields.Enable {

		}
	}
}

func (t *Templater) GenerateFiles() error {
	for _, cfg := range t.configs {
		contexts, err := parseGoFileOrModule(cfg.GoPath, &cfg)
		if err != nil {
			fmt.Printf("Error parsing file: %s\n", err)
			return err
		}

		// Field removal transformer
		if cfg.RemoveFields.Enable {
			t.AddTransformer("remove-fields", cfg.RemoveFields.Transform)
		}

		// Transformer for struct tags to update types
		newTypeTrans := FieldTypeAndNameTransformer{}
		t.AddTransformer("type-name-type", newTypeTrans.Transform)

		tmpl, err := template.ParseFiles(cfg.TemplatePath)
		if err != nil {
			panic(err)
		}
		for _, context := range contexts {
			// Run the transformers if they exist
			for _, fn := range t.Transformers {
				err = fn(&context)
				if err != nil {
					return err
				}
			}

			// Create the output file
			outputFile, err := os.Create(fmt.Sprintf(cfg.OutputPath, context.StructName))
			if err != nil {
				panic(err)
			}
			defer outputFile.Close()

			// Execute the template and write output to the file
			err = tmpl.Execute(outputFile, context)
			if err != nil {
				panic(err)
			}
		}
	}
	return nil
}
