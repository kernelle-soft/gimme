package log

import (
	"fmt"
	"os"
)

func Slice[T any](items []T) {
	result := ""
	for _, item := range items {
		result += fmt.Sprintf("%v\n", item)
	}
	fmt.Fprintln(os.Stdout, result)
}