package lexer

import (
	"unicode"

	"github.com/rafa-ribeiro/brasalang/internal/token"
)

type Lexer struct {
	src  []rune
	pos  int
	line int
	col  int
}

func New(input string) *Lexer {
	return &Lexer{src: []rune(input), line: 1, col: 1}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	start := token.Position{Line: l.line, Column: l.col}

	if l.isAtEnd() {
		return token.Token{Type: token.EOF, Lexeme: "", Position: start}
	}

	ch := l.advance()

	switch ch {
	case '(':
		return token.Token{Type: token.LPAREN, Lexeme: "(", Position: start}
	case ')':
		return token.Token{Type: token.RPAREN, Lexeme: ")", Position: start}
	case '{':
		return token.Token{Type: token.LBRACE, Lexeme: "{", Position: start}
	case '}':
		return token.Token{Type: token.RBRACE, Lexeme: "}", Position: start}
	case ';':
		return token.Token{Type: token.SEMICOLON, Lexeme: ";", Position: start}
	case '+':
		return token.Token{Type: token.PLUS, Lexeme: "+", Position: start}
	case '-':
		return token.Token{Type: token.MINUS, Lexeme: "-", Position: start}
	case '*':
		return token.Token{Type: token.STAR, Lexeme: "*", Position: start}
	case '/':
		return token.Token{Type: token.SLASH, Lexeme: "/", Position: start}
	case '!':
		if l.match('=') {
			return token.Token{Type: token.NOT_EQUAL, Lexeme: "!=", Position: start}
		}
		return token.Token{Type: token.NOT, Lexeme: "!", Position: start}
	case '=':
		if l.match('=') {
			return token.Token{Type: token.EQUAL_EQUAL, Lexeme: "==", Position: start}
		}
		return token.Token{Type: token.EQUAL, Lexeme: "=", Position: start}
	case '>':
		if l.match('=') {
			return token.Token{Type: token.GREATER_EQ, Lexeme: ">=", Position: start}
		}
		return token.Token{Type: token.GREATER, Lexeme: ">", Position: start}
	case '<':
		if l.match('=') {
			return token.Token{Type: token.LESS_EQ, Lexeme: "<=", Position: start}
		}
		return token.Token{Type: token.LESS, Lexeme: "<", Position: start}
	case '&':
		if l.match('&') {
			return token.Token{Type: token.AND_AND, Lexeme: "&&", Position: start}
		}
		return token.Token{Type: token.ILLEGAL, Lexeme: "&", Position: start}
	case '|':
		if l.match('|') {
			return token.Token{Type: token.OR_OR, Lexeme: "||", Position: start}
		}
		return token.Token{Type: token.ILLEGAL, Lexeme: "|", Position: start}
	}

	if unicode.IsDigit(ch) {
		lex := []rune{ch}
		for !l.isAtEnd() && unicode.IsDigit(l.peek()) {
			lex = append(lex, l.advance())
		}
		return token.Token{Type: token.INT, Lexeme: string(lex), Position: start}
	}

	if unicode.IsLetter(ch) || ch == '_' {
		lex := []rune{ch}
		for !l.isAtEnd() && (unicode.IsLetter(l.peek()) || unicode.IsDigit(l.peek()) || l.peek() == '_') {
			lex = append(lex, l.advance())
		}
		ident := string(lex)
		return token.Token{Type: token.LookupIdent(ident), Lexeme: ident, Position: start}
	}

	return token.Token{Type: token.ILLEGAL, Lexeme: string(ch), Position: start}
}

func (l *Lexer) Tokens() []token.Token {
	var out []token.Token
	for {
		tok := l.NextToken()
		out = append(out, tok)
		if tok.Type == token.EOF {
			break
		}
	}
	return out
}

func (l *Lexer) skipWhitespace() {
	for !l.isAtEnd() {
		ch := l.peek()
		// checks for whitespace
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.advance()
			continue
		}
		// checks for newline
		if ch == '\n' {
			l.advance()
			continue
		}
		return
	}
}

func (l *Lexer) isAtEnd() bool {
	return l.pos >= len(l.src)
}

func (l *Lexer) peek() rune {
	if l.isAtEnd() {
		return '\x00'
	}
	return l.src[l.pos]
}

func (l *Lexer) advance() rune {
	ch := l.src[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return ch
}

func (l *Lexer) match(expected rune) bool {
	if l.isAtEnd() || l.src[l.pos] != expected {
		return false
	}
	l.advance()
	return true
}
