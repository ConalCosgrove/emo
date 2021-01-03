package parser

import (
	"emo/ast"
	"emo/lexer"
	"emo/token"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST // Zero value used when no statements have been parsed
	EQUALS // ==
	LESSGREATER // > or <
	SUM // +
	PRODUCT // *
	PREFIX // -X or !X
	CALL // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ: EQUALS,
	token.NOTEQ: EQUALS,
	token.LT: LESSGREATER,
	token.GT: LESSGREATER,
	token.PLUS: SUM,
	token.SUB: SUM,
	token.FWDSLASH: PRODUCT,
	token.ASTRX: PRODUCT,
}

type Parser struct {
	l *lexer.Lexer

	curToken token.Token
	nextToken token.Token
	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) readNextToken() {
	p.curToken = p.nextToken
	p.nextToken = p.l.NextToken()
}

func (p *Parser) readNextPrecedence() int {
	if precedence, ok := precedences[p.nextToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if precedence, ok := precedences[p.curToken.Type]; ok {
		return precedence;
	}
	return LOWEST;
}

	
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
		case token.LET:
			return p.parseLetStatement()
		case token.RETURN: 
			return p.parseReturnStatement()
		default:
			return p.parseExpressionStatement()
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.readNextToken()
	}
	return program
}

func (p *Parser) Errors() []string {
	return p.errors;
}

func (p *Parser) addError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.nextToken.Literal)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	fmt.Printf("parseExpression with %s. Next token is %s\n", p.curToken, p.nextToken)
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.readNextPrecedence() {
		infix := p.infixParseFns[p.nextToken.Type]
		if infix == nil {
			return leftExp
		}

		p.readNextToken()

		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	integer, err := strconv.ParseInt(p.curToken.Literal,10, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = integer
	return lit
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	prefEx := &ast.PrefixExpression{
		Token: p.curToken, 
		Operator: p.curToken.Literal,
	}
	p.readNextToken()
	prefEx.Right = p.parseExpression(PREFIX)
	return prefEx
}

func (p * Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infExp := &ast.InfixExpression{
		Token: p.curToken,
		Operator: p.curToken.Literal,
		Left: left,
	}

	precedence := p.curPrecedence()
	p.readNextToken()
	infExp.Right = p.parseExpression(precedence)
	return infExp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.readNextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.readNextToken()
	}
	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	fmt.Print("parseLetStatement\n")
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.readNextToken()
	stmt.Value = p.parseExpression(LOWEST)
	for !p.curTokenIs(token.SEMICOLON) {
		p.readNextToken()
	}
		return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	fmt.Print("parseReturnStatement\n")
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.readNextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)
	for !p.curTokenIs(token.SEMICOLON) {
		p.readNextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.readNextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.readNextToken()
	}
	return block
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.readNextToken()
		return identifiers
	}

	p.readNextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.readNextToken()
		p.readNextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.readNextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.readNextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}
	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	functionLiteral := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	functionLiteral.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	functionLiteral.Body = p.parseBlockStatement()

	return functionLiteral
	
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.nextToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.readNextToken()
		return true
	} else {
		p.addError(t)
		return false
	}
}

func New(l *lexer.Lexer) * Parser {
	p := &Parser{
		l: l,
		errors: []string{},
	}

	//prefix parse functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.SUB, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	// infix parse functions
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.SUB, p.parseInfixExpression)
	p.registerInfix(token.FWDSLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTRX, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOTEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	p.readNextToken()
	p.readNextToken()
	return p;
}