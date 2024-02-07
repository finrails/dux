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
	MINUS = "-"
	RBAR = "/"
	STAR = "*"
	EXCLAMATION = "!"
	STHAN = "<"
	GTHAN = ">"

	// Double Operators
	EQUAL = "=="
	NEQUAL = "!="

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
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"
	TRUE = "TRUE"
	FALSE = "FALSE"

	// Records
	STRING = "STRING"
)

type Token struct {
	Type TokenType
	Literal string
}

var keywords = map[string]TokenType {
	"fn": FUNCTION,
	"let": LET,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
	"true": TRUE,
	"false": FALSE,
}

func LookupType(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	
	return IDENT
}
