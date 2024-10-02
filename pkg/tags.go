package pkg

import "strings"

func parseTags(tagString string) []map[string]string {
	var tags []map[string]string

	// Split the input string by spaces
	parts := strings.Split(tagString, " ")
	for _, part := range parts {
		// Further split each part by the colon
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			// Remove surrounding quotes from the value
			key := removeQuotesAndBackticks(kv[0])
			value := strings.Trim(kv[1], `"`)
			tags = append(tags, map[string]string{key: removeQuotesAndBackticks(value)})
		}
	}

	return tags
}

func removeQuotesAndBackticks(input string) string {
	// Replace double quotes with an empty string
	result := strings.ReplaceAll(input, `"`, "")
	// Replace backticks with an empty string
	result = strings.ReplaceAll(result, "`", "")
	return result
}
