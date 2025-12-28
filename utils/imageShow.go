package utils

import "fmt"

func ImageShow(index int, modelName, url string) string {
	index++
	return fmt.Sprintf("![%s](%s)", modelName, url)
}
