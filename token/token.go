package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343456
	// Operators
	ASSIGN   = "="
	BANG     = "!"
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="
	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"
)

var infixOperators = []TokenType{
	PLUS, MINUS, ASTERISK, SLASH, LT, GT, EQ, NOT_EQ,
}

var keywords = map[string]TokenType{
	"fn":     TokenType(FUNCTION),
	"let":    TokenType(LET),
	"if":     TokenType(IF),
	"else":   TokenType(ELSE),
	"true":   TokenType(TRUE),
	"false":  TokenType(FALSE),
	"return": TokenType(RETURN),
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TokenType(IDENT)
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func NewToken(typ TokenType, lit byte) Token {
	return Token{
		Type:    typ,
		Literal: string(lit),
	}
}

func InfixOperators() []TokenType {
	infixOps := make([]TokenType, len(infixOperators))
	copy(infixOps, infixOperators)
	return infixOps
}