package conf

import (
	"fmt"
	"unicode/utf8"
)

type scanner struct {
	src []byte

	ch    rune // current character
	off   int  // character offset
	rdOff int  // reading offset (position after current character)
}

func (s *scanner) init(data []byte) {
	s.src = data

	s.ch = ' '
	s.off = 0
	s.rdOff = 0

	s.next()
}

func (s *scanner) scanKeys() map[string][]byte {
	keys := make(map[string][]byte)

	var (
		key   string
		value []byte
	)
scanAgain:
	// scan key
	s.skipWhitespace()
	switch ch := s.ch; {
	case isLetter(ch):
		lit := s.scanIdentifier()
		switch lookupToken(lit) {
		case INCLUDE:
			// TODO: deal with `include` directive
		}
		key = lit
	default: // not a IDENT
		switch ch {
		case -1:
			return keys
		}
		panic("key must be a identifier")
	}

	// scan value
	// basic value type: int, float, string, bool
	valueOff := s.off
	s.skipWhitespace()

	switch ch := s.ch; {
	case isDigit(ch): //  int/float
		_, _ = s.scanNumber(false)
		value = s.src[valueOff:s.off]
	case isLetter(ch): // maybe bool
		_ = s.scanIdentifier()
		value = s.src[valueOff:s.off]
	default:
		switch ch {
		case -1: // eof
			return keys
		case '"': // string
			found := s.findRune('"')
			if found {
				value = s.src[valueOff : s.off+1]
			} else {
				panic("expected \" but not found")
			}
		case '[': //array
			if s.findAnotherPair('[', ']') {
				value = s.src[valueOff : s.off+1]
			} else {
				panic("expected ']' but not found")
			}
		case '{': // dict
			if s.findAnotherPair('{', '}') {
				value = s.src[valueOff : s.off+1]
			} else {
				panic("expected '}' but not found")
			}
		default: // error
			panic(fmt.Sprintf("unexpected rune %b", s.ch))
		}
	}

	keys[key] = value

	goto scanAgain

}

func (s *scanner) findAnotherPair(has, want rune) (found bool) {
	var stacked int
	s.next()
	for s.ch != -1 {
		switch s.ch {
		case has:
			stacked++
		case want:
			if stacked == 0 {
				return true
			}
			stacked--
		}
		s.next()
	}
	return false
}

func (s *scanner) findRune(r rune) (found bool) {
	s.next()
	for s.ch != r && s.ch != -1 {
		s.next()
	}
	if s.ch == -1 {
		return false
	}
	return true
}

func (s *scanner) scan() (pos int, tok Token, lit string) {
	s.skipWhitespace()

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		switch lookupToken(lit) {
		case INCLUDE:
			return s.off, INCLUDE, lit
		}
		return s.off, IDENT, lit
	case isDigit(ch):
		tok, lit = s.scanNumber(false)

	}

	return 0, invalid, ""
}

func (s *scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

const bom = 0xFEFF // byte order mark, only permitted as very first character

func (s *scanner) next() {
	if s.rdOff < len(s.src) {
		s.off = s.rdOff
		r, w := rune(s.src[s.rdOff]), 1
		switch {
		case r == 0:
			panic("illegal character NUL")
		case r >= utf8.RuneSelf:
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOff:])
			if r == utf8.RuneError && w == 1 {
				panic("illegal UTF-8 encoding")
			} else if r == bom && s.off > 0 {
				panic("illegal byte order mark")
			}
		}
		s.rdOff += w
		s.ch = r
	} else {
		s.off = len(s.src)
		s.ch = -1 // eof
	}
}

func (s *scanner) scanIdentifier() string {
	off := s.off
	for isLetter(s.ch) || isDigit(s.ch) {
		s.next()
	}
	return string(s.src[off:s.off])
}

func (s *scanner) scanNumber(seenDecimalPoint bool) (Token, string) {
	// digitVal(s.ch) < 10
	offs := s.off
	tok := INT

	if seenDecimalPoint {
		offs--
		tok = FLOAT
		s.scanMantissa(10)
		goto exponent
	}

	if s.ch == '0' {
		// int or float
		offs := s.off
		s.next()
		if s.ch == 'x' || s.ch == 'X' {
			// hexadecimal int
			s.next()
			s.scanMantissa(16)
			if s.off-offs <= 2 {
				// only scanned "0x" or "0X"
				panic("illegal hexadecimal number")
			}
		} else {
			// octal int or float
			seenDecimalDigit := false
			s.scanMantissa(8)
			if s.ch == '8' || s.ch == '9' {
				// illegal octal int or float
				seenDecimalDigit = true
				s.scanMantissa(10)
			}
			if s.ch == '.' || s.ch == 'e' || s.ch == 'E' || s.ch == 'i' {
				goto fraction
			}
			// octal int
			if seenDecimalDigit {
				panic("illegal octal number")
			}
		}
		goto exit
	}

	// decimal int or float
	s.scanMantissa(10)

fraction:
	if s.ch == '.' {
		tok = FLOAT
		s.next()
		s.scanMantissa(10)
	}

exponent:
	if s.ch == 'e' || s.ch == 'E' {
		tok = FLOAT
		s.next()
		if s.ch == '-' || s.ch == '+' {
			s.next()
		}
		if digitVal(s.ch) < 10 {
			s.scanMantissa(10)
		} else {
			panic("illegal floating-point exponent")
		}
	}

	if s.ch == 'i' {
		tok = IMAG
		s.next()
	}

exit:
	return tok, string(s.src[offs:s.off])
}

func (s *scanner) scanMantissa(base int) {
	for digitVal(s.ch) < base {
		s.next()
	}
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	}
	return 16 // larger than any legal digit val
}

func isLetter(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
