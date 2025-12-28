package utils

import (
	"fmt"
	"pplx2api/config"
)

func searchShowDetails(index int, title, url, snippet string) string {
	return fmt.Sprintf("<details>\n<summary>[%d] %s</summary>\n\n%s\n\n[Link](%s)\n\n</details>", index, title, snippet, url)
}

func searchShowCompatible(index int, title, url, snippet string) string {
	return fmt.Sprintf("[%d] [%s](%s):\n%s\n", index, title, url, snippet)
}

func SearchShow(index int, title, url, snippet string) string {
	index++
	if len([]rune(snippet)) > 150 {
		runeSnippet := []rune(snippet)
		snippet = fmt.Sprintf("%s ……", string(runeSnippet[:150]))
	}
	if config.ConfigInstance.SearchResultCompatible {
		return searchShowCompatible(index, title, url, snippet)
	}
	return searchShowDetails(index, title, url, snippet)
}
