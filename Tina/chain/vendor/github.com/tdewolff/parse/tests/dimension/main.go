// +build gofuzz
package fuzz

import "github.com/tdewolff/parse/v2"

// Fuzz is a fuzz test.
func Fuzz(data []byte) int {
	_, _ = parse.Dimension(data)
	return 1
}
