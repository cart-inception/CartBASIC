package token

import "strings"

// TokenType identifies the kind of token.
type TokenType string

// Token represents a lexical token and its source literal.
type Token struct {
	Type    TokenType
	Literal string
}

const (
	// Special tokens.
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals.
	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	STRING TokenType = "STRING"

	// Operators.
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	BANG     TokenType = "!"
	LT       TokenType = "<"
	GT       TokenType = ">"
	EQ       TokenType = "=="
	NOT_EQ   TokenType = "!="

	// Delimiters.
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"
	DOT       TokenType = "."
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	LBRACKET  TokenType = "["
	RBRACKET  TokenType = "]"

	// Keywords.
	FUNCTION TokenType = "FUNCTION"
	LET      TokenType = "LET"
	IF       TokenType = "IF"
	ELSEIF   TokenType = "ELSEIF"
	ELSE     TokenType = "ELSE"
	WHILE    TokenType = "WHILE"
	FOR      TokenType = "FOR"
	RETURN   TokenType = "RETURN"
	TRY      TokenType = "TRY"
	CATCH    TokenType = "CATCH"
	SPAWN    TokenType = "SPAWN"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
)

var keywords = map[string]TokenType{
	"fn":       FUNCTION,
	"function": FUNCTION,
	"let":      LET,
	"if":       IF,
	"elseif":   ELSEIF,
	"else":     ELSE,
	"while":    WHILE,
	"for":      FOR,
	"return":   RETURN,
	"try":      TRY,
	"catch":    CATCH,
	"spawn":    SPAWN,
	"true":     TRUE,
	"false":    FALSE,
}

// LookupIdent returns a keyword token type for language keywords, or IDENT otherwise.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}

	return IDENT
}
