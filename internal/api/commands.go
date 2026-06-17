package api

import (
	"net/url"
	"strings"
)

func isCloseCommand(transcript string) bool {
	normalized := strings.ToLower(transcript)
	phrases := []string{"close it", "close this", "okay close it", "hide jarvis"}
	for _, phrase := range phrases {
		if strings.Contains(normalized, phrase) {
			return true
		}
	}

	return false
}

func toImageURLs(queries []string) []string {
	urls := make([]string, 0, len(queries))
	for _, query := range queries {
		trimmed := strings.TrimSpace(query)
		if trimmed == "" {
			continue
		}

		urls = append(urls, "https://loremflickr.com/800/600/"+url.QueryEscape(trimmed))
	}

	return urls
}
