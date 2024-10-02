package pkg

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	GoPath       string                  `yaml:"go_path"`
	TemplatePath string                  `yaml:"template_path"`
	Name         string                  `yaml:"name"`
	OutputPath   string                  `yaml:"output_path"`
	RemoveFields RemoveFieldsTransformer `yaml:"remove_fields"`
}

func loadYAMLConfigs(filePath string) ([]Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	documents := bytes.Split(data, []byte("\n---\n"))

	var configs []Config
	for _, doc := range documents {
		var config Config
		if err := yaml.Unmarshal(doc, &config); err != nil {
			return nil, fmt.Errorf("error unmarshaling document: %w", err)
		}
		configs = append(configs, config)
	}

	return configs, nil
}
