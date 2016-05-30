package boolParser

import (
	"io"
)

type TokenType int
type TokenOffset int

const (
	tokUndefined TokenType = iota // 0 //actually not used
	tokIdent     TokenType = iota

	// operations
	// binary
	tokUnion        // +
	tokIntersection // *

	// unary
	tokNegation // '!'

	// other symbols
	tokOpeningParenthesis // '('
	tokClosingParenthesis // ')'

	tokOpeningBraces // "["
	tokClosingBraces // "]"
	tokConst         // '0' or '1'

	// special
	tokEOF
	tokError
)

func (t TokenType) String() string {
	switch t {
	case tokUndefined:
		return "Undefined"
	case tokIdent:
		return "Identifier"
	case tokUnion:
		return "Union"
	case tokIntersection:
		return "Intersection"
	case tokNegation:
		return "Negation"
	case tokOpeningBraces:
		return "OpeningBraces"
	case tokOpeningParenthesis:
		return "OpeningParenthesis"
	case tokClosingBraces:
		return "ClosingBraces"
	case tokClosingParenthesis:
		return "ClosingParenthesis"
	case tokConst:
		return "Const"
	case tokEOF:
		return "EOF"
	case tokError:
		return "Error"
	default:
		return "Unknown"
	}
}

type Token struct {
	_type  TokenType
	offset TokenOffset
	value  string
}

// some reusable constant tokens
var singleCharTokens map[byte]TokenType = map[byte]TokenType{
	'+': tokUnion,
	'*': tokIntersection,
	'!': tokNegation,
	'(': tokOpeningParenthesis,
	')': tokClosingParenthesis,
	'[': tokOpeningBraces,
	']': tokClosingBraces,
}

func isBinDigit(c byte) bool {
	return '0' == c || c == '1'
}

func isAlpha(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func isAlphaDig(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func readIdentifier(p []byte, n int, i int, r io.Reader) (string, int, int, error) {
	startIndex := i
	str := ""
	i++

	for {
		for ; i < n; i++ {
			if !isAlphaDig(p[i]) {
				str += string(p[startIndex:i])
				return str, n, i - 1, nil
			}
		}
		i = 0
		startIndex = 0

		var err error
		n, err = r.Read(p)
		if err != nil {
			return str, n, i, err
		}
	}
}

func Lexer(r io.Reader, out chan<- Token) {
	var err error
	var p []byte = make([]byte, 256)
	offset := TokenOffset(0)

	for {
		var n int
		n, err = r.Read(p)
		if err != nil {
			break
		}

		for i := 0; i < n; i++ {
			char := p[i]

			switch {
			case char <= ' ':
				offset++
			case singleCharTokens[char] != tokUndefined:
				out <- Token{_type: singleCharTokens[char], offset: offset, value: string(char)}
				offset++
			case isBinDigit(char):
				out <- Token{_type: tokConst, offset: offset, value: string(char)}
				offset++
			case isAlpha(char):
				var ident string
				ident, n, i, err = readIdentifier(p, n, i, r)
				out <- Token{_type: tokIdent, offset: offset, value: ident}
				offset += TokenOffset(len(ident))

				if err != nil {
					break
				}
			default:
				break
			}
		}
	}
	if err == io.EOF {
		out <- Token{_type: tokEOF, offset: offset}
		return
	}
	out <- Token{_type: tokError, offset: offset, value: err.Error()}
}
