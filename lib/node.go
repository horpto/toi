package boolParser

import "errors"

type Namespace map[string]bool

type Node interface {
	String() string
	Calculate(Namespace) (bool, error)
}

type BinaryNode interface {
	Node

	LeftExpression() Node
	RightExpression() Node

	SetLeftExpression(Node)
	SetRightExpression(Node)
}

type Identifier struct {
	Name string
}

func (id Identifier) Calculate(n Namespace) (bool, error) {
	val, ok := n[id.Name]
	if !ok {
		return false, errors.New("Var '" + id.Name + "' not found")
	}
	return val, nil
}

func (id Identifier) String() string {
	return id.Name
}

type Const struct {
	Value string
}

func (c Const) Calculate(n Namespace) (bool, error) {
	return c.Value == "1", nil
}

func (c Const) String() string {
	return c.Value
}

type NegationNode struct {
	expr Node
}

func (nn NegationNode) Calculate(ns Namespace) (bool, error) {
	val, err := nn.expr.Calculate(ns)
	if err != nil {
		return false, err
	}
	return !val, nil
}

func (nn NegationNode) String() string {
	return "!" + nn.expr.String()
}

type BinaryNodeStruct struct {
	LExpr Node
	RExpr Node
}

func (bns *BinaryNodeStruct) LeftExpression() Node {
	return bns.LExpr
}

func (bns *BinaryNodeStruct) RightExpression() Node {
	return bns.RExpr
}

func (bns *BinaryNodeStruct) SetLeftExpression(n Node) {
	bns.LExpr = n
}

func (bns *BinaryNodeStruct) SetRightExpression(n Node) {
	bns.RExpr = n
}

type UnionNode struct {
	BinaryNodeStruct
}

func (un UnionNode) Calculate(ns Namespace) (bool, error) {
	lexpr, err := un.LExpr.Calculate(ns)
	if err != nil {
		return lexpr, err
	}
	rexpr, err := un.RExpr.Calculate(ns)
	if err != nil {
		return rexpr, err
	}
	return (lexpr || rexpr), nil
}

func (un UnionNode) String() string {
	return "(" + un.LExpr.String() + " + " + un.RExpr.String() + ")"
}

type IntersectionNode struct {
	BinaryNodeStruct
}

func (in IntersectionNode) Calculate(ns Namespace) (bool, error) {
	lexpr, err := in.LExpr.Calculate(ns)
	if err != nil {
		return lexpr, err
	}
	rexpr, err := in.RExpr.Calculate(ns)
	if err != nil {
		return rexpr, err
	}
	return (lexpr && rexpr), nil
}

func (in IntersectionNode) String() string {
	return "(" + in.LExpr.String() + " * " + in.RExpr.String() + ")"
}
