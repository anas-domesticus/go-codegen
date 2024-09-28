package pkg

import "strings"

type StructConfig struct {
	Values map[string]string
	Flags  []string
}

func parseCommentAndLoadConfig(comment string) StructConfig {
	config := StructConfig{
		Values: make(map[string]string),
	}

	// Remove '@codegen' and split by whitespace
	comment = strings.TrimSpace(comment)
	comment = strings.ReplaceAll(comment, "//", "")
	comment = strings.ReplaceAll(comment, "@codegen", "")
	comment = strings.TrimSpace(comment)
	parts := strings.Fields(comment)

	for _, part := range parts {
		if strings.Contains(part, "=") {
			// It's a key-value pair
			keyValue := strings.SplitN(part, "=", 2)
			if len(keyValue) == 2 {
				key := strings.TrimSpace(keyValue[0])
				value := strings.TrimSpace(keyValue[1])
				config.Values[key] = value
			}
		} else {
			// It's a single flag
			config.Flags = append(config.Flags, part)
		}
	}

	return config
}
