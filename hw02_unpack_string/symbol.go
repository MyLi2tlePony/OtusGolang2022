package hw02unpackstring

import "unicode"

type symbol struct {
	char      rune
	isEscaped bool
}

func (s *symbol) isEscaping() bool {
	return s.char == '\\'
}

func (s *symbol) isDigit() bool {
	return unicode.IsDigit(s.char)
}
