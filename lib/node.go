package boolParser

type Node interface {
}

type BinaryNode interface {
	Node

	LeftExpression() Node
	RightExpression() Node

	SetLeftExpression(Node)
	SetRightExpression(Node)
}

type Identifier struct {
	name string
}

type Const struct {
	value string
}

type NegationNode struct {
	expr Node
}

type BinaryNodeStruct struct {
	BinaryNode
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

type DifferenceNode struct {
	BinaryNodeStruct
}

type SymDifferenceNode struct {
	BinaryNodeStruct
}

type IntersectionNode struct {
	BinaryNodeStruct
}
