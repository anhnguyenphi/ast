package main

import (
	"errors"
	"fmt"
)

const (
	AdsQueryNot = "not"
	AdsQueryAnd = "and"
	AdsQueryOr = "or"
	AdsQueryLT = "<"
	AdsQueryLTE = "<="
	AdsQueryGT = ">"
	AdsQueryGTE = ">="
	AdsQueryEQ = "="
	AdsQueryNEQ = "<>"
)

type Node interface {
	GetParent() *Node
	SetParent(parent *Node)
	GetValue() interface{}
	SetValue(value interface{})
	GetChildren() []*Node
	SetChildren(nodes... *Node)
	Accept(AdsQueryVisitor) // Using visitor pattern for some additional logic
}

type Metadata struct {
	Segment  bool
	Property bool
}

type ValueNode struct {
	Metadata *Metadata
	Parent *Node
	Value interface{}
}

func (node *ValueNode) GetParent() *Node {
	return node.Parent
}

func (node *ValueNode) SetParent(parent *Node) {
	node.Parent = parent
}

func (node *ValueNode) GetValue() interface{} {
	return node.Value
}

func (node *ValueNode) SetValue(value interface{}) {
	node.Value = value
}

func (node *ValueNode) GetChildren() []*Node {
	return nil
}

func (node *ValueNode) SetChildren(nodes... *Node) {
}

func (node *ValueNode) Accept(visitor AdsQueryVisitor)  {
	visitor.VisitValueNode(node)
}

type BinaryNode struct {
	ValueNode
	Left *Node
	Right *Node
}

func (node *BinaryNode) GetChildren() []*Node {
	if node.Left == nil && node.Right == nil {
		return nil
	}

	var result []*Node
	result = append(result, node.Left)
	result = append(result, node.Right)
	return result
}

func (node *BinaryNode) SetChildren(nodes... *Node) {
	if nodes == nil {
		return
	}

	for i, n := range nodes {
		if i == 0 {
			node.Left = n
			continue
		}
		if i == 1 {
			node.Right = n
			break
		}
	}
}

func (node *BinaryNode) Accept(visitor AdsQueryVisitor)  {
	visitor.VisitBinaryOpNode(node)

	if node.Left != nil {
		(*node.Left).Accept(visitor)
	}

	if node.Right != nil {
		(*node.Right).Accept(visitor)
	}
}

type UnaryNode struct {
	ValueNode
	Child *Node
}

func (node *UnaryNode) GetChildren() []*Node {
	if node.Child == nil {
		return nil
	}

	var result []*Node
	result = append(result, node.Child)
	return result
}

func (node *UnaryNode) SetChildren(nodes... *Node) {
	if nodes == nil {
		return
	}

	for i, n := range nodes {
		if i == 0 {
			node.Child = n
			break
		}
	}
}

func (node *UnaryNode) Accept(visitor AdsQueryVisitor)  {
	visitor.VisitSingleOpNode(node)

	child := node.Child
	if child != nil {
		(*child).Accept(visitor)
	}
}

func traversal(node Node, level int)  {
	if node == nil {
		return
	}

	for i := 0; i < level; i++ {
		fmt.Print("-")
	}

	// TODO: do something
	fmt.Println(node.GetValue())

	children := node.GetChildren()
	if children != nil {
		for _, child := range children {
			traversal(*child, level + 1)
		}
	}
}

// TODO: visitor for validation, etc
type AdsQueryVisitor interface {
	Visit(Node)
	VisitValueNode(Node)
	VisitBinaryOpNode(Node)
	VisitSingleOpNode(Node)
}


type ValidationVisitor struct {
	errors []error
}

func (visitor *ValidationVisitor) Visit(node Node)  {
	//TODO
}

func (visitor *ValidationVisitor) VisitValueNode(node Node)  {
	visitor.errors = append(visitor.errors, errors.New("new error"))
}


func (visitor *ValidationVisitor) VisitBinaryOpNode(node Node)  {
	//TODO
}

func (visitor *ValidationVisitor) VisitSingleOpNode(node Node)  {
	//TODO
}

type A struct {

}

type B struct {
	A
}

func main() {
	var bin Node = &BinaryNode{}
	bin.SetValue("AND")

	var left Node = &ValueNode{}
	left.SetValue("female")
	left.SetParent(&bin)

	var sin Node = &UnaryNode{}
	sin.SetValue("NOT")

	var child Node = &ValueNode{}
	child.SetParent(&sin)
	child.SetValue("vietnamese")

	sin.SetChildren(&child)
	sin.SetParent(&bin)

	bin.SetChildren(&left, &sin)

	traversal(bin,0)

	fmt.Println("------------")

	visitor := &ValidationVisitor{}
	bin.Accept(visitor)
	fmt.Println(visitor.errors)
}
