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
	p.skipNewlines()

	for !p.check(token.EOF) {
		stmt := p.parseStatement()
		if stmt == nil {
			p.synchronize()
			p.skipNewlines()
			continue
		}
		program.Statements = append(program.Statements, stmt)

		if !p.consumeStatementTerminator() {
			p.synchronize()
		}
		p.skipNewlines()
	}

	return program
}

func (p *Parser) parseStatement() ast.Stmt {
	switch {
	case p.check(token.LBRACE):
		return p.parseBlockStatement()
	case p.isVarDeclStart():
		return p.parseVarDeclStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) isVarDeclStart() bool {
	return p.check(token.IDENT) && p.peekN(1).Type == token.IDENT && p.peekN(2).Type == token.EQUAL
}

func (p *Parser) parseVarDeclStatement() ast.Stmt {
	nameTok, ok := p.expect(token.IDENT, "expected variable name")
	if !ok {
		return nil
	}

	typeTok, ok := p.expect(token.IDENT, "expected type name after variable name")
	if !ok {
		return nil
	}

	if _, ok := p.expect(token.EQUAL, "expected '=' after variable type"); !ok {
		return nil
	}

	initializer := p.parseExpression()
	if initializer == nil {
		return nil
	}

	return &ast.VarDeclStmt{Name: nameTok, TypeName: typeTok, Initializer: initializer}
}

func (p *Parser) parseBlockStatement() ast.Stmt {
	lbrace, _ := p.expect(token.LBRACE, "expected '{'")
	block := &ast.BlockStmt{LBrace: lbrace}

	p.skipNewlines()
	for !p.check(token.RBRACE) && !p.check(token.EOF) {
		stmt := p.parseStatement()
		if stmt == nil {
			p.synchronize()
			p.skipNewlines()
			continue
		}
		block.Statements = append(block.Statements, stmt)

		if !p.consumeStatementTerminator() {
			p.synchronize()
		}
		p.skipNewlines()
	}

	if _, ok := p.expect(token.RBRACE, "expected '}' to close block"); !ok {
		return nil
	}

	return block
}

func (p *Parser) parseExpressionStatement() ast.Stmt {
	expr := p.parseExpression()
	if expr == nil {
		return nil
	}
	return &ast.ExprStmt{Expression: expr}
}

func (p *Parser) consumeStatementTerminator() bool {
	if p.check(token.NEWLINE) {
		p.skipNewlines()
		return true
	}
	if p.check(token.EOF) || p.check(token.RBRACE) {
		return true
	}

	at := p.peek()
	p.errs = append(p.errs, fmt.Errorf("expected newline after statement at %d:%d", at.Position.Line, at.Position.Column))
	return false
}

func (p *Parser) skipNewlines() {
	for p.check(token.NEWLINE) {
		p.advance()
	}
}

func (p *Parser) parseExpression() ast.Expr {
	return p.parseOr()
}

func (p *Parser) parseOr() ast.Expr {
	left := p.parseAnd()
	for p.check(token.OR_OR) {
		op := p.advance()
		right := p.parseAnd()
		if left == nil || right == nil {
			return nil
		}
		left = &ast.BinaryExpr{Left: left, Operator: op, Right: right}
	}
	return left
}

func (p *Parser) parseAnd() ast.Expr {
	left := p.parseEquality()
	for p.check(token.AND_AND) {
		op := p.advance()
		right := p.parseEquality()
		if left == nil || right == nil {
			return nil
		}
		left = &ast.BinaryExpr{Left: left, Operator: op, Right: right}
	}
	return left
}

func (p *Parser) parseEquality() ast.Expr {
	left := p.parseComparison()
	for p.check(token.EQUAL_EQUAL) || p.check(token.NOT_EQUAL) {
		op := p.advance()
		right := p.parseComparison()
		if left == nil || right == nil {
			return nil
		}
		left = &ast.BinaryExpr{Left: left, Operator: op, Right: right}
	}
	return left
}

func (p *Parser) parseComparison() ast.Expr {
	left := p.parseTerm()
	for p.check(token.GREATER) || p.check(token.GREATER_EQ) || p.check(token.LESS) || p.check(token.LESS_EQ) {
		op := p.advance()
		right := p.parseTerm()
		if left == nil || right == nil {
			return nil
		}
		left = &ast.BinaryExpr{Left: left, Operator: op, Right: right}
	}
	return left
}

func (p *Parser) parseTerm() ast.Expr {
	left := p.parseFactor()
	for p.check(token.PLUS) || p.check(token.MINUS) {
		op := p.advance()
		right := p.parseFactor()
		if left == nil || right == nil {
			return nil
		}
		left = &ast.BinaryExpr{Left: left, Operator: op, Right: right}
	}
	return left
}

func (p *Parser) parseFactor() ast.Expr {
	left := p.parseUnary()
	for p.check(token.STAR) || p.check(token.SLASH) {
		op := p.advance()
		right := p.parseUnary()
		if left == nil || right == nil {
			return nil
		}
		left = &ast.BinaryExpr{Left: left, Operator: op, Right: right}
	}
	return left
}

func (p *Parser) parseUnary() ast.Expr {
	if p.check(token.NOT) || p.check(token.MINUS) {
		op := p.advance()
		right := p.parseUnary()
		if right == nil {
			return nil
		}
		return &ast.UnaryExpr{Operator: op, Right: right}
	}

	return p.parsePrimary()
}

func (p *Parser) parsePrimary() ast.Expr {
	tok := p.peek()
	switch tok.Type {
	case token.LPAREN:
		p.advance()
		expr := p.parseExpression()
		if expr == nil {
			return nil
		}
		if _, ok := p.expect(token.RPAREN, "expected ')' after expression"); !ok {
			return nil
		}
		return expr
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
		if p.previous().Type == token.NEWLINE || p.check(token.RBRACE) {
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
	return p.peekN(0)
}

func (p *Parser) peekN(offset int) token.Token {
	idx := p.curr + offset
	if idx >= len(p.tokens) {
		return token.Token{Type: token.EOF}
	}

	return p.tokens[idx]
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
