package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const escapeLetter = `\`

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var unpackedString strings.Builder
	var previousLetter rune
	var isExcludedCase bool

	for i, letter := range str {
		if !isValidLetter(letter) {
			return "", ErrInvalidString
		}

		if i == 0 && !isValidEscapeLetter(letter) {
			return "", ErrInvalidString
		}

		if isValidEscapeLetter(letter) {
			if isEscapedLetter(previousLetter) {
				if !isEscapedLetter(letter) {
					return "", ErrInvalidString
				}

				isExcludedCase = true

				continue
			}

			if previousLetter != 0 {
				unpackedString.WriteRune(previousLetter)
			}

			isExcludedCase = false
			previousLetter = letter
		}

		if unicode.IsDigit(letter) {
			if previousLetter == 0 {
				return "", ErrInvalidString
			}

			if unicode.IsDigit(previousLetter) && !isExcludedCase {
				return "", ErrInvalidString
			}

			if isEscapedLetter(previousLetter) {
				if !isExcludedCase {
					isExcludedCase = true

					previousLetter = letter

					continue
				}
			}

			multipliedLetter, isError := multipleLetter(letter, previousLetter)

			if isError {
				return "", ErrInvalidString
			}

			unpackedString.WriteString(multipliedLetter)
			previousLetter = 0
		}
	}

	if previousLetter != 0 {
		unpackedString.WriteRune(previousLetter)
	}

	return unpackedString.String(), nil
}

func isValidLetter(letter rune) bool {
	return isValidEscapeLetter(letter) || unicode.IsDigit(letter)
}

func isValidEscapeLetter(letter rune) bool {
	return unicode.IsLetter(letter) || isEscapedLetter(letter)
}

func isEscapedLetter(letter rune) bool {
	return escapeLetter == string(letter)
}

func multipleLetter(letter, previousLetter rune) (string, bool) {
	digit, errAtoi := strconv.Atoi(string(letter))

	if errAtoi != nil {
		return "", true
	}

	return strings.Repeat(string(previousLetter), digit), false
}
