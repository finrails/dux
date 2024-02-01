package ast

import (
	"bytes"
	"dux/src/token"
)

/*
	Interfaces are really simple, the user can use a interface as expected parameter,
	but the user needs to implement all the interface methods into a structure if he
	wants to use it as a argument.

	So if a user sees a function that expects a Statement parameter, he only can call
	the function passing a structure that implements Statement interface methods.

	That's why the user needs to implement all Statement interface methods if he wants
	to use a structure as Statement argument.

	A interface can also nest other interfaces, for instance, the Statement interface
	nests Node interface, so the user will needs to implement Node interface as well.
*/

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type Identifier struct {
	Token token.Token
	Value string
}

type Boolean struct {
	Token token.Token
	Value bool
}

type InfixExpression struct {
	Token    token.Token // Operator Token (i.e. +, -, >, <...)
	Left     Expression
	Operator string
	Right    Expression
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct {
	Token				token.Token // The 'return' Token
	ReturnValue Expression
}

type ExpressionStatement struct {
	Token token.Token
	Expression Expression
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

func (ifex *IfExpression) expressionNode() {}
func (ifex *IfExpression) TokenLiteral() string { return ifex.Token.Literal }
func (ifex *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ifex.Condition.String())
	out.WriteString(" ")
	out.WriteString(ifex.Consequence.String())

	if ifex.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ifex.Alternative.String())
	}

	return out.String()
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string { return b.Token.Literal }

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string { return il.Token.Literal }

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var output bytes.Buffer

	output.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		output.WriteString(rs.ReturnValue.String())
	}

	output.WriteString(";")

	return output.String()
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var output bytes.Buffer

	output.WriteString(ls.TokenLiteral() + " " + ls.Name.String() + " = " )

	if ls.Value != nil {
		output.WriteString(ls.Value.String())
	}

	output.WriteString(";")

	return output.String()
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string { return i.Value }

func (p *Program) String() string {
	var output bytes.Buffer

	for _, stmt := range p.Statements {
		output.WriteString(stmt.String())
	}

	// Refactor it later to use StringBuilder, to stringfy efficiently
	return output.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}
