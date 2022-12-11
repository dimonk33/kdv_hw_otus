package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidString Сообщение об ошибке.
var ErrInvalidString = errors.New("invalid string")

// Unpack Распаковка строки.
func Unpack(inStr string) (string, error) {
	if _, errInt := strconv.Atoi(inStr); errInt == nil {
		return "", ErrInvalidString
	}

	isEscape := false
	prevSym := rune(0)
	builder := strings.Builder{}

	for _, r := range inStr {
		if isEscape {
			if prevSym != rune(0) {
				builder.WriteRune(prevSym)
			}
			prevSym = r
			isEscape = false
			continue
		}

		if r == '\\' {
			isEscape = true
			continue
		}

		if unicode.IsDigit(r) {
			if prevSym == rune(0) {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(r))
			builder.WriteString(strings.Repeat(string(prevSym), count))
			prevSym = rune(0)
			continue
		}

		if prevSym != rune(0) {
			builder.WriteRune(prevSym)

		}
		prevSym = r
	}

	if prevSym != rune(0) {
		builder.WriteRune(prevSym)
	}

	return builder.String(), nil
}
