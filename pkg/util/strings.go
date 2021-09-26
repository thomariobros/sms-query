package util

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// StandardizeStringLower strip accents, remove useless spaces, replace " by '
// and to lower case
func StandardizeStringLower(str string) string {
	return StandardizeString(str, true)
}

// StandardizeString strip accents, remove useless spaces, replace " by '
// and to lower case if lower is true
func StandardizeString(str string, lower bool) string {
	// strip accents
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, str)
	// remove useless spaces
	result = strings.Join(strings.Fields(result), " ")
	// replace all " by '
	result = strings.Replace(result, "\"", "'", -1)
	// replace all ’ by '
	result = strings.Replace(result, "’", "'", -1)

	// to lower
	if lower {
		return strings.ToLower(result)
	}
	return result
}
