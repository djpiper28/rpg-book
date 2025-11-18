package normalisation

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func Normalise(text string) string {
	normText := strings.ToLower(norm.NFKD.String(text))
	output := ""
	for _, c := range normText {
		if unicode.IsMark(c) {
			continue
		}

		if unicode.IsPrint(c) {
			output += string(c)
		}
	}

	return output
}
