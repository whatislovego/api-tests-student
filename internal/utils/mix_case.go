package utils

import "unicode"

func MixCase(text string) string {
	runes := []rune(text)

	for i, char := range runes {
		if i%2 == 0 {
			runes[i] = unicode.ToLower(char)
		} else {
			runes[i] = unicode.ToUpper(char)
		}
	}

	return string(runes)
}
