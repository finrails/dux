package parser

import (
	"dux/src/ast"
	"dux/src/lexer"
	"dux/src/token"
)

type Parser struct {
	l *lexer.Lexer

	currentToken token.Token
	peekToken token.Token
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken() // Shift ahead two times, to read and set the tokens.

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
