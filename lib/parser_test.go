package boolParser

import (
	"strings"
	"testing"
)

// --- TEST UTILITIES

type TestTokenStream struct {
	stack []*Token
}

func (ts TestTokenStream) topToken() Token {
	return *ts.stack[len(ts.stack)-1]
}

func (ts *TestTokenStream) popToken() {
	ts.stack = ts.stack[:len(ts.stack)-1]
}

func (ts *TestTokenStream) putToken(tok Token) {
	ts.stack = append(ts.stack, &tok)
}

func concat(arrs ...[]*Token) []*Token {
	outArr := []*Token{}
	for _, arr := range arrs {
		outArr = append(outArr, arr...)
	}
	return outArr
}

func mkNegationStack(_type TokenType, value string) []*Token {
	return []*Token{
		&Token{_type: _type, offset: 1, value: value},
		&Token{_type: tokNegation, offset: 0, value: "!"},
	}
}

func mkBracketsStack(stack []*Token) []*Token {
	return concat(
		[]*Token{&Token{_type: tokClosingBraces, value: ")"}},
		stack,
		[]*Token{&Token{_type: tokOpeningBraces, value: "("}},
	)
}

func mkParenthesisStack(stack []*Token) []*Token {
	return concat(
		[]*Token{&Token{_type: tokClosingBraces, value: ")"}},
		stack,
		[]*Token{&Token{_type: tokOpeningBraces, value: "("}},
	)
}

// --- CHECKS

func checkIdentifier(t *testing.T, node Node, value string) {
	if node == nil {
		t.Error("expected Node, actual: nil")
	}
	ident, ok := node.(Identifier)
	if !ok {
		t.Errorf("expected value of Identifier type, actual: %T", node)
	}
	if ident.Name != value {
		t.Error("expected value:", value, ", actual:", ident.Name)
	}
}

func checkConst(t *testing.T, node Node, value string) {
	if node == nil {
		t.Error("expected Node, actual: nil")
		return
	}
	_const, ok := node.(Const)
	if !ok {
		t.Errorf("expected value of Const type, actual: %T", node)
		return
	}
	if _const.Value != value {
		t.Error("expected value:", value, ", actual:", _const.Value)
	}
}

func checkIdentOrConst(t *testing.T, node Node, _type TokenType, value string) {
	if _type == tokIdent {
		checkIdentifier(t, node, value)
	} else {
		checkConst(t, node, value)
	}
}

func checkTokenStream(t *testing.T, ts *TestTokenStream) {
	if len(ts.stack) > 0 {
		t.Error("token stream not empty")
	}
}

func checkExpression(t *testing.T, ts *TestTokenStream, _type TokenType, value string) {
	node := parseExpression(ts)
	checkIdentOrConst(t, node, _type, value)
	checkTokenStream(t, ts)
}

// --- TESTS

func TestParseIdentifier(t *testing.T) {
	//t.Parallel()
	value := "x"
	ts := &TestTokenStream{
		stack: []*Token{
			&Token{_type: tokIdent, offset: 0, value: value},
		},
	}
	node := parseIdentifier(ts)
	checkIdentifier(t, node, value)
	checkTokenStream(t, ts)
}

func TestParseConst(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"0", "1"} {
		ts := &TestTokenStream{
			stack: []*Token{
				&Token{_type: tokConst, offset: 0, value: value},
			},
		}
		node := parseConst(ts)
		checkConst(t, node, value)
		checkTokenStream(t, ts)
	}
}

func TestParseExpressionCanParseIdentifier(t *testing.T) {
	t.Parallel()
	for _, value := range "yx" {
		value := string(value)
		ts := &TestTokenStream{
			stack: []*Token{
				&Token{_type: tokIdent, offset: 0, value: value},
			},
		}
		checkExpression(t, ts, tokIdent, value)
	}
}

func TestParseExpressionCanParseConst(t *testing.T) {
	t.Parallel()
	for _, value := range "01" {
		value := string(value)
		ts := &TestTokenStream{
			stack: []*Token{
				&Token{_type: tokConst, offset: 0, value: value},
			},
		}
		checkExpression(t, ts, tokConst, value)
	}
}

var identsAndConsts = map[string]TokenType{
	"x": tokIdent,
	"z": tokIdent,
	"0": tokConst,
	"1": tokConst,
}

func TestParseNegation(t *testing.T) {
	t.Parallel()
	for value, _type := range identsAndConsts {
		ts := &TestTokenStream{stack: mkNegationStack(_type, value)}
		node, ok := parseNegation(ts).(NegationNode)
		if !ok {
			t.Errorf("Expected Negation type, actual %T", node)
		}
		checkIdentOrConst(t, node.expr, _type, value)
		checkTokenStream(t, ts)
	}
}

func TestParseExpressionCanParseNegation(t *testing.T) {
	t.Parallel()
	for value, _type := range identsAndConsts {
		ts := &TestTokenStream{stack: mkNegationStack(_type, value)}
		node, ok := parseExpression(ts).(NegationNode)
		if !ok {
			t.Errorf("Expected Negation type, actual %T", node)
		}
		checkIdentOrConst(t, node.expr, _type, value)
		checkTokenStream(t, ts)
	}
}

func TestParseExpressionCanParseBrackets(t *testing.T) {
	t.Parallel()

	for value, _type := range identsAndConsts {
		var (
			stack []*Token
			ts    *TestTokenStream
		)
		stack = mkBracketsStack([]*Token{&Token{_type: _type, value: value}})
		ts = &TestTokenStream{stack: stack}
		checkExpression(t, ts, _type, value)

		stack = mkParenthesisStack([]*Token{&Token{_type: _type, value: value}})
		ts = &TestTokenStream{stack: stack}
		checkExpression(t, ts, _type, value)

		// double brackets
		stack = mkBracketsStack([]*Token{&Token{_type: _type, value: value}})
		stack = mkParenthesisStack(stack)
		ts = &TestTokenStream{stack: stack}
		checkExpression(t, ts, _type, value)

		stack = mkParenthesisStack([]*Token{&Token{_type: _type, value: value}})
		stack = mkBracketsStack(stack)
		ts = &TestTokenStream{stack: stack}
		checkExpression(t, ts, _type, value)

		// negation in brackets
		{
			stack = mkNegationStack(_type, value)
			stack = mkBracketsStack(stack)
			ts = &TestTokenStream{stack: stack}
			node, ok := parseExpression(ts).(NegationNode)
			if ok {
				checkIdentOrConst(t, node.expr, _type, value)
				checkTokenStream(t, ts)
			} else {
				t.Errorf("Expected type Negation get %t", node)
			}
		}

		// negation in brackets
		{
			stack = mkNegationStack(_type, value)
			stack = mkParenthesisStack(stack)
			ts = &TestTokenStream{stack: stack}
			node, ok := parseExpression(ts).(NegationNode)
			if ok {
				checkIdentOrConst(t, node.expr, _type, value)
				checkTokenStream(t, ts)
			} else {
				t.Errorf("Expected type Negation get %t", node)
			}
		}
	}
}

// TODO: write tests for parseIntersectionStatement

func TestParseStatementPriority(t *testing.T) {
	t.Parallel()
	//a + b * c
	//(a + (b * c))
	stack := []*Token{
		&Token{_type: tokEOF},
		&Token{_type: tokIdent, value: "c"},
		&Token{_type: tokIntersection, value: "*"},
		&Token{_type: tokIdent, value: "b"},
		&Token{_type: tokUnion, value: "+"},
		&Token{_type: tokIdent, value: "a"},
	}
	ts := &TestTokenStream{stack: stack}
	node := parseStatement(ts)
	if node == nil {
		t.Error("expected expression, actual nil")
	}
	unode, ok := node.(*UnionNode)
	if !ok {
		t.Errorf("expected Union, actual %T", node)
	}
	checkIdentifier(t, unode.LExpr, "a")
	rnode, ok := unode.RExpr.(*IntersectionNode)
	if !ok {
		t.Errorf("expected Intersection, actual %T", unode.RExpr)
	}
	checkIdentifier(t, rnode.LExpr, "b")
	checkIdentifier(t, rnode.RExpr, "c")
}

func TestParser(t *testing.T) {
	t.Parallel()
	ch := make(chan Token)
	r := strings.NewReader("a + c * d - !(q \\ 1) ")
	go Lexer(r, ch)

	ts := NewArrayTokenStream(ch)
	ast, err := Parser(&ts)
	if err != nil {
		t.Error(err.Error())
	}
	if ast == nil {
		t.Error("expected AST, actual nil are got")
	}
	node, ok := ast.(*UnionNode)
	if !ok {
		t.Errorf("expected Union, actual got: %T", ast)
	}
	checkIdentifier(t, node.LExpr, "a")
}
