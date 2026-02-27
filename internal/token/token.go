package token

type Type string

const (
	// Special
	EOF     Type = "EOF"
	ILLEGAL Type = "ILLEGAL"
	NEWLINE Type = "NEWLINE"

	// Literals
	IDENT Type = "IDENT"
	INT   Type = "INT"

	// Keywords
	TRUE  Type = "TRUE"
	FALSE Type = "FALSE"
	IF    Type = "IF"
	ELSE  Type = "ELSE"

	// Delimiters
	LPAREN Type = "LPAREN"
	RPAREN Type = "RPAREN"
	LBRACE Type = "LBRACE"
	RBRACE Type = "RBRACE"

	// Operators
	PLUS        Type = "PLUS"
	MINUS       Type = "MINUS"
	STAR        Type = "STAR"
	SLASH       Type = "SLASH"
	NOT         Type = "NOT"
	EQUAL       Type = "EQUAL"
	EQUAL_EQUAL Type = "EQUAL_EQUAL"
	NOT_EQUAL   Type = "NOT_EQUAL"
	GREATER     Type = "GREATER"
	GREATER_EQ  Type = "GREATER_EQ"
	LESS        Type = "LESS"
	LESS_EQ     Type = "LESS_EQ"
	AND_AND     Type = "AND_AND"
	OR_OR       Type = "OR_OR"
)

type Position struct {
	Line   int
	Column int
}

type Token struct {
	Type     Type
	Lexeme   string
	Position Position
}

var keywords = map[string]Type{
	"if":    IF,
	"else":  ELSE,
	"true":  TRUE,
	"false": FALSE,
}

func LookupIdent(ident string) Type {
	if typ, ok := keywords[ident]; ok {
		return typ
	}
	return IDENT
}
