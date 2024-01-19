package ast

import "dux/src/token"

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

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct {
	Token				token.Token // The 'return' Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}
