package parser

import (
	"fmt"
	"strconv"

	"github.com/rafa-ribeiro/brasalang/internal/ast"
	"github.com/rafa-ribeiro/brasalang/internal/lexer"
	"github.com/rafa-ribeiro/brasalang/internal/token"
)

type Parser struct {
	tokens []token.Token
	curr   int
	errs   []error
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func NewFromSource(src string) *Parser {
	l := lexer.New(src)
	return New(l.Tokens())
}

func (p *Parser) Errors() []error {
	return p.errs
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	for !p.check(token.EOF) {
		stmt := p.parseExpressionStatement()
		if stmt == nil {
			p.synchronize()
			continue
		}
		program.Statements = append(program.Statements, stmt)
	}
	return program
}

func (p *Parser) parseExpressionStatement() ast.Stmt {
	expr := p.parsePrimary()
	if expr == nil {
		return nil
	}
	semi, ok := p.expect(token.SEMICOLON, "expected ';' after expression") // i may want to remove it
	if !ok {
		return nil
	}
	return &ast.ExprStmt{Expression: expr, Semicolon: semi}
}

func (p *Parser) parsePrimary() ast.Expr {
	tok := p.peek()
	switch tok.Type {
	case token.INT:
		p.advance()
		v, err := strconv.ParseInt(tok.Lexeme, 10, 64)
		if err != nil {
			p.errs = append(p.errs, fmt.Errorf("invalid integer %q at %d:%d", tok.Lexeme, tok.Position.Line, tok.Position.Column))
			return nil
		}
		return &ast.IntLiteral{Token: tok, Value: v}
	case token.TRUE:
		p.advance()
		return &ast.BoolLiteral{Token: tok, Value: true}
	case token.FALSE:
		p.advance()
		return &ast.BoolLiteral{Token: tok, Value: false}
	case token.IDENT:
		p.advance()
		return &ast.Identifier{Token: tok, Name: tok.Lexeme}
	default:
		p.errs = append(p.errs, fmt.Errorf("unexpected token %s (%q) at %d:%d", tok.Type, tok.Lexeme, tok.Position.Line, tok.Position.Column))
		return nil
	}
}

func (p *Parser) synchronize() {
	for !p.check(token.EOF) {
		if p.previous().Type == token.SEMICOLON {
			return
		}
		p.advance()
	}
}

func (p *Parser) expect(tt token.Type, msg string) (token.Token, bool) {
	if p.check(tt) {
		return p.advance(), true
	}
	at := p.peek()
	p.errs = append(p.errs, fmt.Errorf("%s at %d:%d", msg, at.Position.Line, at.Position.Column))
	return token.Token{}, false
}

func (p *Parser) check(tt token.Type) bool {
	return p.peek().Type == tt
}

func (p *Parser) peek() token.Token {
	if p.curr >= len(p.tokens) {
		return token.Token{Type: token.EOF}
	}
	return p.tokens[p.curr]
}

func (p *Parser) previous() token.Token {
	if p.curr == 0 {
		return token.Token{Type: token.EOF}
	}
	return p.tokens[p.curr-1]
}

func (p *Parser) advance() token.Token {
	tok := p.peek()
	if p.curr < len(p.tokens) {
		p.curr++
	}
	return tok
}
