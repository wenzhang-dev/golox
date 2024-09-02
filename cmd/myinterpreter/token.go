package main

import (
    "fmt"
    "strconv"
)


const (
    TK_LEFT_PAREN = "LEFT_PAREN"      // (
    TK_RIGHT_PAREN = "RIGHT_PAREN"    // )
    TK_LEFT_BRACE = "LEFT_BRACE"      // {
    TK_RIGHT_BRACE = "RIGHT_BRACE"    // }
    TK_STAR = "STAR"                  // *
    TK_DOT = "DOT"                    // .
    TK_COMMA = "COMMA"                // ,
    TK_PLUS = "PLUS"                  // +
    TK_MINUS = "MINUS"                // -
    TK_SEMICOLON = "SEMICOLON"        // ;
    TK_EQUAL = "EQUAL"                // =
    TK_EQUAL_EQUAL = "EQUAL_EQUAL"    // ==
    TK_BANG = "BANG"                  // !
    TK_BANG_EQUAL = "BANG_EQUAL"      // !=
    TK_LESS = "LESS"                  // <
    TK_LESS_EQUAL = "LESS_EQUAL"      // <=
    TK_GREATER = "GREATER"            // >
    TK_GREATER_EQUAL = "GREATER_EQUAL"// >=
    TK_SLASH = "SLASH"                // /
    TK_STRING = "STRING"              // "abc"
    TK_NUMBER = "NUMBER"              // 3.14
    TK_IDENTIFIER = "IDENTIFIER"      // abc
    TK_EOF = "EOF"                    // EOF
    TK_INVALID = "INVALID TOKEN"      // INVALID TOKEN

    KW_AND = "AND"                    // and
    KW_CLASS = "CLASS"                // class
    KW_ELSE = "ELSE"                  // else
    KW_FALSE = "FALSE"                // false
    KW_TRUE = "TRUE"                  // true
    KW_FOR = "FOR"                    // for
    KW_FUN = "FUN"                    // fun
    KW_IF = "IF"                      // if
    KW_NIL = "NIL"                    // nil
    KW_OR = "OR"                      // or
    KW_PRINT = "PRINT"                // print
    KW_RETURN = "RETURN"              // return
    KW_SUPER = "SUPER"                // super
    KW_THIS = "THIS"                  // this
    KW_VAR = "VAR"                    // var
    KW_WHILE = "WHILE"                // where
)

type Token struct {
    Line int
    Lexeme string
    Type string
}

func NewToken(typ string, lexeme string, line int) Token {
    return Token {
        Line: line,
        Lexeme: lexeme,
        Type: typ,
    }
}

func (t Token) Literal() (ValueType, error) {
    switch (t.Type) {
    case TK_STRING:
        return StringType{v: t.Lexeme[1:len(t.Lexeme)-1]}, nil
    case TK_NUMBER:
        if v, err := strconv.ParseFloat(t.Lexeme, 64); err != nil {
            return NilValue, nil
        } else {
            return NumberType{v: v}, nil
        }
    case KW_FALSE:
        return FalseValue, nil
    case KW_TRUE:
        return TrueValue, nil
    case KW_NIL:
        return NilValue, nil
    }

    return nil, fmt.Errorf("Expect a literal type")
}

func (t Token) LiteralString() string {
    if t.Type == TK_STRING {
        return t.Lexeme[1:len(t.Lexeme)-1]
    }

    if t.Type == TK_NUMBER {
        num, _ := strconv.ParseFloat(t.Lexeme, 64)
        if num == float64(int(num)) {
            return fmt.Sprintf("%.1f", num)
        } else {
            return fmt.Sprintf("%v", num)
        }
    }

    return "null"
}

func (t Token) ToString() string {
    // token_type lexeme literal
    return fmt.Sprintf("%s %s %s", t.Type, t.Lexeme, t.LiteralString())
}
