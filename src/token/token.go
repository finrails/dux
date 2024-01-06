package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	// Identifiers and literals
	IDENT = "IDENT" // ADD, X, Y...
	INT = "INT" // ...-2, -1, 0, 1, 2...

	// Operators
	ASSIGN = "="
	PLUS = "+"

	// Delimiters
	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET = "LET"
)

type Token struct {
	Type TokenType
	Literal string
}
