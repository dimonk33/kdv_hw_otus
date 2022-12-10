package main

import (
	"errors"
	"fmt"
	"strconv"
)

// ErrInvalidString Сообщение об ошибке.
var ErrInvalidString = errors.New("invalid string")

func main() {
	str, err := Unpack(string("aaa10b"))
	fmt.Println(str, err)
}

// Unpack Распаковка строки.
func Unpack(inStr string) (string, error) {
	_, errInt := strconv.Atoi(inStr)
	if errInt == nil {
		return "", ErrInvalidString
	}

	isEscape := false
	outStr := ""
	prevSym := rune(0)

	for _, r := range inStr {
		if isEscape {
			if prevSym != rune(0) {
				outStr += string(prevSym)
			}
			prevSym = r
			isEscape = false
			continue
		}

		if r == '\\' {
			isEscape = true
			continue
		}

		if r >= '0' && r <= '9' {
			if prevSym == rune(0) {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(r))
			for i := 0; i < count; i++ {
				outStr += string(prevSym)
			}
			prevSym = rune(0)
			continue
		}

		if prevSym != rune(0) {
			outStr += string(prevSym)
		}
		prevSym = r
	}

	if prevSym != rune(0) {
		outStr += string(prevSym)
	}

	return outStr, nil
}
