package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(packString string) (string, error) {
	var unpackString strings.Builder
	var prevSymbol symbol

	for symbolIndex, currentSymbol := range packString {
		if symbolIndex == 0 {
			if unicode.IsDigit(currentSymbol) {
				return "", ErrInvalidString
			}

			prevSymbol = symbol{currentSymbol, false}
		} else {
			unpackedSymbol, err := unpackSymbol(currentSymbol, prevSymbol)
			if err != nil {
				return "", ErrInvalidString
			}

			prevSymbol = symbol{
				char:      currentSymbol,
				isEscaped: prevSymbol.isEscaping() && !prevSymbol.isEscaped,
			}
			unpackString.WriteString(unpackedSymbol)
		}
	}

	if prevSymbol.char != 0 && (!prevSymbol.isDigit() || prevSymbol.isEscaped) {
		unpackString.WriteRune(prevSymbol.char)
	}
	if prevSymbol.isEscaping() && !prevSymbol.isEscaped {
		return "", ErrInvalidString
	}

	return unpackString.String(), nil
}

func unpackSymbol(currentSymbol rune, prevSymbol symbol) (string, error) {
	if !unicode.IsDigit(currentSymbol) {
		if !prevSymbol.isEscaped && (prevSymbol.isDigit() || prevSymbol.isEscaping()) {
			return "", nil
		}

		return string(prevSymbol.char), nil
	}
	if !prevSymbol.isEscaping() && !prevSymbol.isEscaped && prevSymbol.isDigit() {
		return "", ErrInvalidString
	}
	if prevSymbol.isEscaping() && !prevSymbol.isEscaped {
		return "", nil
	}

	number, err := strconv.Atoi(string(currentSymbol))
	if err != nil {
		return "", err
	}

	return strings.Repeat(string(prevSymbol.char), number), nil
}
