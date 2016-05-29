package boolParser

import (
	"errors"
)

func parseIdentifier(ts TokenStream) Node {
	token := ts.topToken()
	if token._type != tokIdent {
		return nil
	}

	ts.popToken()
	return Identifier{name: token.value}
}

func parseConst(ts TokenStream) Node {
	token := ts.topToken()
	if token._type != tokConst {
		return nil
	}

	ts.popToken()
	return Const{value: token.value}
}

func parseNegation(ts TokenStream) Node {
	token := ts.topToken()
	if token._type != tokNegation {
		return nil
	}
	ts.popToken()
	return NegationNode{expr: parseExpression(ts)}
}

func parseExpression(ts TokenStream) Node {
	var node Node
	if node = parseConst(ts); node != nil {
		return node
	}
	if node = parseIdentifier(ts); node != nil {
		return node
	}
	if node = parseNegation(ts); node != nil {
		return node
	}
	oType := ts.topToken()._type
	if oType != tokOpeningParenthesis && oType != tokOpeningBraces {
		// maybe panic ???
		return nil
	}

	ts.popToken()
	if node = parseStatement(ts); node == nil {
		return nil
		// panic() // TODO
	}
	cType := ts.topToken()._type
	if oType == tokOpeningBraces && cType != tokClosingBraces {
		return nil
		// panic() // TODO
	}
	if oType == tokOpeningParenthesis && cType != tokClosingParenthesis {
		return nil
		// panic() // TODO
	}

	ts.popToken()
	return node
}

func parseIntersectionExpression(ts TokenStream) Node {
	lExpr := parseExpression(ts)
	if lExpr == nil {
		return nil
		// panic() // TODO
	}

	t_type := ts.topToken()._type
	var node BinaryNode
	switch t_type {
	case tokIntersection:
		node = &IntersectionNode{}
		node.SetLeftExpression(lExpr)
	case tokEOF:
		return lExpr
	default:
		return lExpr
	}

	ts.popToken()
	rExpr := parseIntersectionExpression(ts)
	if rExpr == nil {
		return nil
		// panic() // TODO
	}
	node.SetRightExpression(rExpr)
	return node
}

func parseStatement(ts TokenStream) Node {
	lExpr := parseIntersectionExpression(ts)
	if lExpr == nil {
		return nil
		// panic() // TODO
	}

	t_type := ts.topToken()._type
	var bn BinaryNode
	switch t_type {
	case tokUnion:
		bn = &UnionNode{}
	case tokDifference:
		bn = &DifferenceNode{}
	case tokSymmDifference:
		bn = &SymDifferenceNode{}
	default:
		return lExpr
	}
	bn.SetLeftExpression(lExpr)
	ts.popToken()

	rExpr := parseStatement(ts)
	if rExpr == nil {
		return nil
		// panic() // TODO
	}
	bn.SetRightExpression(rExpr)
	ts.popToken()
	return bn
}

func Parser(ts TokenStream) (node Node, err error) {
	node = parseStatement(ts)

	if node == nil {
		err = errors.New("fail to parse")
	}
	return
}
