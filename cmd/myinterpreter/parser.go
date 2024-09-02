package main

import (
    "fmt"
)

type Parser struct {
    Tokens []Token
    Current int
}

func NewParser(tokens []Token) *Parser {
    return &Parser {
        Tokens: tokens,
        Current: 0,
    }
}

func (p *Parser) MatchAny(token_types ...string) bool {
    for _, typ := range token_types {
        if p.Check(typ) {
            p.Advance()
            return true 
        }
    }

    return false 
}


func (p *Parser) Advance() *Token {
    if p.IsEnd() {
        return &p.Tokens[len(p.Tokens)-1]
    }
    
    token := p.Peek()
    p.Current++
    return token
}

func (p *Parser) Check(token_type string) bool {
    if p.IsEnd() {
        return token_type == TK_EOF
    }
    return p.Peek().Type == token_type
}

func (p *Parser) Peek() *Token {
    // TODO: assert p.Current < len(p.Tokens)
    return &p.Tokens[p.Current]
}

func (p *Parser) Previous() *Token {
    // TODO: assert p.Current - 1 >= 0
    return &p.Tokens[p.Current-1]
}

func (p *Parser) IsEnd() bool {
    return p.Current >= len(p.Tokens) || p.Peek().Type == TK_EOF
}

// expect and consume the token type. if mismatch, return an error
func (p *Parser) Expect(tk string, msg string) (*Token, error) {
    if p.Check(tk) {
        return p.Advance(), nil
    }

    return nil, fmt.Errorf("%s", msg)
}

func (p *Parser) Parse() ([]Stmt, error) {
    stmts := []Stmt{}
    for !p.IsEnd() {
        if stmt, err := p.ParseStatement(); err != nil {
            return stmts, err
        } else {
            stmts = append(stmts, stmt)
        }
    }

    return stmts, nil
}

func (p *Parser) ParseStatement() (Stmt, error) {
    if p.MatchAny(KW_PRINT, KW_VAR, TK_LEFT_BRACE) {
        switch p.Previous().Type {
        case KW_PRINT:
            return p.ParsePrintStatement()
        case KW_VAR:
            return p.ParseVarStatement()
        case TK_LEFT_BRACE:
            return p.ParseBlock()
        }
    }

    return p.ParseExpressionStatement()
}

func (p *Parser) ParseBlock() (Stmt, error) {
    stmts := []Stmt{}
    for !p.IsEnd() && !p.Check(TK_RIGHT_BRACE) {
        if stmt, err := p.ParseStatement(); err != nil {
            return nil, err
        } else {
            stmts = append(stmts, stmt)
        }
    }

    if _, err := p.Expect(TK_RIGHT_BRACE, "Expect '}' after block."); err != nil {
        return nil, err
    }

    return BlockStmt{stmts: stmts}, nil
}

func (p *Parser) ParseVarStatement() (Stmt, error) {
    var tk *Token
    var err error

    // we can handle the following use cases:
    //  1) var a;
    //  2) var a=1;
    //  3) var a=b;
    //  4) var a=b=4;
    if tk, err = p.Expect(TK_IDENTIFIER, "Expect an Identifier"); err != nil {
        return nil, err
    }

    var expr Expr

    if p.MatchAny(TK_EQUAL) {
        if expr, err = p.ParseExpression(); err != nil {
            return nil, err
        }
    }

    if _, err = p.Expect(TK_SEMICOLON, "Expect ';' after expression"); err != nil {
        return nil, err
    }

    return VarStmt{v: expr, tk: tk}, nil
}

func (p *Parser) ParseExpressionStatement() (Stmt, error) {
    expr, err := p.ParseExpression()
    if err != nil {
        return nil, err
    }

    if _, err = p.Expect(TK_SEMICOLON, "Expect ';' after expression."); err != nil {
        return nil, err
    }

    return ExprStmt{v: expr}, nil
}

func (p *Parser) ParsePrintStatement() (Stmt, error) {
    expr, err := p.ParseExpression()
    if err != nil {
        return nil, err
    }

    if _, err = p.Expect(TK_SEMICOLON, "Expect ';' after value."); err != nil {
        return nil, err
    }

    return PrintStmt{v: expr}, nil
}

func (p *Parser) ParseExpression() (Expr, error) {
    return p.ParseAssignment()
}

func (p *Parser) ParseAssignment() (Expr, error) {
    expr, err := p.ParseEquality()
    if err != nil {
        return nil, err
    }

    // recursive descent parse
    // we can handle the following use cases:
    //   1) a=b=123
    //   2) a=b
    //   3) a=123
    if p.MatchAny(TK_EQUAL) {
        val, err := p.ParseAssignment()
        if err != nil {
            return nil, err
        }

        if v, ok := expr.(VarExpr); ok {
            expr = AssignmentExpr{token: v.token, expr: val}
        } else {
            return nil, fmt.Errorf("Invalid assignment expression.")
        }
    }

    return expr, nil
}

func (p *Parser) ParseEquality() (Expr, error) {
    expr, err := p.ParseComparsion()
    if err != nil {
        return nil, err
    }

    // recursive descent parse
    for p.MatchAny(TK_EQUAL_EQUAL, TK_BANG_EQUAL) {
        optr := p.Previous()
        right, err := p.ParseComparsion()
        if err != nil {
            return nil, err
        }
        expr = BinaryExpr{optr: optr, left:expr, right:right}
    }

    return expr, nil
}

func (p *Parser) ParseComparsion() (Expr, error) {
    expr, err := p.ParseTerm()
    if err != nil {
        return nil, err
    }

    // recursive descent parse
    for p.MatchAny(TK_GREATER, TK_GREATER_EQUAL, TK_LESS, TK_LESS_EQUAL) {
        optr := p.Previous()
        right, err := p.ParseTerm()
        if err != nil {
            return nil, err
        }
        expr = BinaryExpr{optr: optr, left:expr, right: right}
    }

    return expr, nil
}

func (p *Parser) ParseTerm() (Expr, error) {
    expr, err := p.ParseFactor()
    if err != nil {
        return nil, err
    }

    // recursive descent parse
    for p.MatchAny(TK_MINUS, TK_PLUS) {
        optr := p.Previous()
        right, err := p.ParseFactor()
        if err != nil {
            return nil, err
        }
        expr = BinaryExpr{optr:optr, left: expr, right: right}
    }

    return expr, nil
}

func (p *Parser) ParseFactor() (Expr, error) {
    expr, err := p.ParseUnary()
    if err != nil {
        return nil, err
    }

    // recursive descent parse
    for p.MatchAny(TK_STAR, TK_SLASH) {
        optr := p.Previous()
        right, err := p.ParseUnary()
        if err != nil {
            return nil, err
        }
        expr = BinaryExpr{optr: optr, left: expr, right: right}
    }

    return expr, nil
}

func (p *Parser) ParseUnary() (Expr, error) {
    if p.MatchAny(TK_BANG, TK_MINUS) {
        optr := p.Previous()
        expr, err := p.ParseUnary()
        if err != nil {
            return nil, err
        }
        return UnaryExpr{token: optr, expr: expr}, nil
    }

    return p.ParsePrimary()
}

func (p *Parser) ParsePrimary() (Expr, error) {
    if p.MatchAny(TK_NUMBER, TK_STRING, KW_TRUE, KW_FALSE, KW_NIL) {
        return LiteralExpr{token: p.Previous()}, nil
    } else if p.MatchAny(TK_LEFT_PAREN) {
        expr, err := p.ParseExpression()
        if err != nil {
            return nil, err
        }
        if p.MatchAny(TK_RIGHT_PAREN) {
            return GroupExpr{expr: expr}, nil
        }
        return nil, fmt.Errorf("Error: unmatched parenthesis")
    } else if p.MatchAny(TK_IDENTIFIER) {
        return VarExpr{token: p.Previous()}, nil
    }

    // p.Advance()
    return nil, fmt.Errorf("[line %d] Error at '%s': Expect expression", p.Peek().Line, p.Peek().Lexeme)
}
