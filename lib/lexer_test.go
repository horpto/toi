package boolParser

import "testing"
import "strings"

func TestsIsbinDigit(t *testing.T) {
	t.Parallel()
	if !isBinDigit('0') {
		t.Error("Expected true, but actual is false for '0'")
	}
	if !isBinDigit('1') {
		t.Error("Expected true, but actual is false for '1'")
	}
	if isBinDigit('3') {
		t.Error("Expected false, but actual is true for '3'")
	}
}

func TestIsAlpha(t *testing.T) {
	t.Parallel()
	for i := byte('a'); i <= 'z'; i++ {
		if !isAlpha(i) {
			t.Error("Expected true, but actual is false for '", string(i), "'")
		}
	}
	for i := byte('A'); i <= 'Z'; i++ {
		if !isAlpha(i) {
			t.Error("Expected true, but actual is false for '", string(i), "'")
		}
	}
}

func TestIsAlphaDig(t *testing.T) {
	t.Parallel()
	for i := byte('a'); i <= 'z'; i++ {
		if !isAlphaDig(i) {
			t.Error("Expected true, but actual is false for '", string(i), "'")
		}
	}
	for i := byte('A'); i <= 'Z'; i++ {
		if !isAlphaDig(i) {
			t.Error("Expected true, but actual is false for '", string(i), "'")
		}
	}

	for i := byte('0'); i <= '9'; i++ {
		if !isAlphaDig(i) {
			t.Error("Expected true, but actual is false for '", string(i), "'")
		}
	}
}

var testsLexer = map[string]([]Token){
	"": {
		Token{_type: tokEOF},
	},
	"1": {
		Token{_type: tokConst, value: "1"},
		Token{_type: tokEOF},
	},
	"ax1231 ": {
		Token{_type: tokIdent, value: "ax1231"},
		Token{_type: tokEOF},
	},
	"   (  !\t\r\n  [a1QzRT* q + x*y] )    1": {
		Token{_type: tokOpeningParenthesis, value: "("},
		Token{_type: tokNegation, value: "!"},
		Token{_type: tokOpeningBraces, value: "["},
		Token{_type: tokIdent, value: "a1QzRT"},
		Token{_type: tokIntersection, value: "*"},
		Token{_type: tokIdent, value: "q"},
		Token{_type: tokUnion, value: "+"},
		Token{_type: tokIdent, value: "x"},
		Token{_type: tokIntersection, value: "*"},
		Token{_type: tokIdent, value: "y"},
		Token{_type: tokClosingBraces, value: "]"},
		Token{_type: tokClosingParenthesis, value: ")"},
		Token{_type: tokConst, value: "1"},
		Token{_type: tokEOF},
	},
	"a + c * d + !(q * 1) ": {
		Token{_type: tokIdent, value: "a"},
		Token{_type: tokUnion, value: "+"},
		Token{_type: tokIdent, value: "c"},
		Token{_type: tokIntersection, value: "*"},
		Token{_type: tokIdent, value: "d"},
		Token{_type: tokUnion, value: "+"},
		Token{_type: tokNegation, value: "!"},
		Token{_type: tokOpeningParenthesis, value: "("},
		Token{_type: tokIdent, value: "q"},
		Token{_type: tokIntersection, value: "*"},
		Token{_type: tokConst, value: "1"},
		Token{_type: tokClosingParenthesis, value: ")"},
		Token{_type: tokEOF},
	},
}

func TestLexer(t *testing.T) {
	t.Parallel()

	for k, v := range testsLexer {
		ch := make(chan Token)
		r := strings.NewReader(k)

		go Lexer(r, ch)

		for _, tok := range v {
			rtok := <-ch
			//t.Log("get token:", rtok._type, rtok.value)
			if rtok._type != tok._type {
				t.Error("Expected token type", tok._type, "actual:", rtok._type, "test:", k)
				break
			}
			if rtok.value != tok.value {
				t.Error("Expected value", tok.value, "actual:", rtok.value, "test:", k)
				break
			}
			// don't check offset
		}
	}
}
