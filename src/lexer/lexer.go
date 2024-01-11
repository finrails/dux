package lexer

import "dux/src/token"

type Lexer struct {
	input string
	position int // Current position in input, current char
	readPosition int // Current reading position in input, after current char
	char byte // Current character under analysis
}

func (lex *Lexer) readChar() {
	if lex.readPosition >= len(lex.input) {
		lex.char = 0
	} else {
		lex.char = lex.input[lex.readPosition]
	}

	lex.position = lex.readPosition

	lex.readPosition += 1
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

func (lex *Lexer) eatGhostCharacters() {
	for lex.char == ' ' || lex.char == '\t' || lex.char == '\n' || lex.char == '\r' {
		lex.readChar()
	}
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_' || char == '?'
}

func (lex *Lexer) readIdentifier() string {
	startPosition := lex.position

	for isLetter(lex.char) {
		lex.readChar()
	}

	return lex.input[startPosition:lex.position]
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func (lex *Lexer) readNumber() string {
	startPosition := lex.position

	for isDigit(lex.char) {
		lex.readChar()
	}

	return lex.input[startPosition:lex.position]
}

func (lex *Lexer) NextToken() token.Token {
	var tok token.Token

	lex.eatGhostCharacters()

	switch lex.char {
	case '=':
		tok = newToken(token.ASSIGN, lex.char)
	case ';':
		tok = newToken(token.SEMICOLON, lex.char)
	case '(':
		tok = newToken(token.LPAREN, lex.char)
	case ')':
		tok = newToken(token.RPAREN, lex.char)
	case ',':
		tok = newToken(token.COMMA, lex.char)
	case '+':
		tok = newToken(token.PLUS, lex.char)
	case '{':
		tok = newToken(token.LBRACE, lex.char)
	case '}':
		tok = newToken(token.RBRACE, lex.char)
	case '!':
		tok = newToken(token.EXCLAMATION, lex.char)
	case '-':
		tok = newToken(token.MINUS, lex.char)
	case '/':
		tok = newToken(token.RBAR, lex.char)
	case '*':
		tok = newToken(token.STAR, lex.char)
	case '<':
		tok = newToken(token.STHAN, lex.char)
	case '>':
		tok = newToken(token.GTHAN, lex.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(lex.char) {
			tok.Literal = lex.readIdentifier()
			tok.Type = token.LookupType(tok.Literal)
			return tok
		} else if isDigit(lex.char) {
			tok.Literal = lex.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, lex.char)
		}
	}

	lex.readChar()

	return tok
}

func New(input string) *Lexer {
	lex := &Lexer{input: input}
	lex.readChar()

	return lex
}
