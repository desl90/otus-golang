package unpacker

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const escapeLetter = '\\'

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var unpackedString strings.Builder
	var previousLetter rune
	var isExcluded bool

	for i, letter := range str {
		switch {
		case !isValidLetter(letter):
			return "", ErrInvalidString

		case i == 0 && !isValidEscapeLetter(letter):
			return "", ErrInvalidString

		case !unicode.IsDigit(letter) && !isEscapedLetter(letter) && isExcluded:
			return "", ErrInvalidString

		case unicode.IsLetter(letter) && isExcluded:
			return "", ErrInvalidString

		case isEscapedLetter(letter) && !isExcluded:
			isExcluded = true

		case previousLetter == 0 && (!unicode.IsDigit(letter) || isExcluded):
			previousLetter = letter
			isExcluded = false

		case unicode.IsDigit(letter) && !isExcluded:
			if previousLetter == 0 {
				return "", ErrInvalidString
			}

			multipliedLetter, isError := multipleLetter(letter, previousLetter)

			if isError {
				return "", ErrInvalidString
			}

			unpackedString.WriteString(multipliedLetter)
			previousLetter = 0
			
		default:
			unpackedString.WriteRune(previousLetter)
			previousLetter = letter
			isExcluded = false
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
	return escapeLetter == letter
}

func multipleLetter(letter, previousLetter rune) (string, bool) {
	digit, errAtoi := strconv.Atoi(string(letter))

	if errAtoi != nil {
		return "", true
	}

	return strings.Repeat(string(previousLetter), digit), false
}
