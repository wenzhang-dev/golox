package main

import (
	"fmt"
	"os"
)

type Scanner struct {
    Source string
    Start int
    Current int
    Line int

    Tokens []Token

    HasErr bool 

    ReservedKws map[string]string
}

func (s *Scanner) HasError() bool {
    return s.HasErr
} 

func NewScanner(source string) *Scanner {
    return &Scanner {
        Source: source,
        Start: 0,
        Current: 0,
        Line: 1,
        Tokens: make([]Token, 0),
        HasErr: false,
        ReservedKws: map[string]string{
            "and": KW_AND,
            "class": KW_CLASS,
            "else": KW_ELSE,
            "false": KW_FALSE,
            "true": KW_TRUE,
            "for": KW_FOR,
            "fun": KW_FUN,
            "if": KW_IF,
            "nil": KW_NIL,
            "or": KW_OR,
            "print": KW_PRINT,
            "return": KW_RETURN,
            "super": KW_SUPER,
            "this": KW_THIS,
            "var": KW_VAR,
            "while": KW_WHILE,
        },
    }
}

func (s *Scanner) IsEnd() bool {
    return s.Current >= len(s.Source)
}

func (s *Scanner) ScanTokens() []Token {
    for !s.IsEnd() {
        s.Start = s.Current
        s.ScanToken()
    }

    s.Start = s.Current
    s.AddToken(TK_EOF)
    
    return s.Tokens
}

func (s *Scanner) Advance() rune {
    if s.IsEnd() {
        return 0
    }

    var ch = rune(s.Source[s.Current])
    s.Current = s.Current + 1
    return ch
}

func (s *Scanner) AddToken(token_type string) Token {
    token := Token {
        Line: s.Line,
        Lexeme: s.String(),
        Type: token_type,
    }
    s.Tokens = append(s.Tokens, token)
    return token
}

func (s *Scanner) ScanToken() {
    var c = s.Advance()
    switch {
    case c == '(':
        s.AddToken(TK_LEFT_PAREN)
    case c == ')':
        s.AddToken(TK_RIGHT_PAREN)
    case c == '{':
        s.AddToken(TK_LEFT_BRACE)
    case c == '}':
        s.AddToken(TK_RIGHT_BRACE)
    case c == '.':
        s.AddToken(TK_DOT)
    case c == '+':
        s.AddToken(TK_PLUS)
    case c == '-':
        s.AddToken(TK_MINUS)
    case c == '*':
        s.AddToken(TK_STAR)
    case c == ',':
        s.AddToken(TK_COMMA)
    case c == ';':
        s.AddToken(TK_SEMICOLON)
    case c == ' ':
    case c == '\r':
    case c == '\t':
        // Ingore white space
    case c == '\n':
        s.Line++
    case c == '=':
        if s.Match("=") {
            s.AddToken(TK_EQUAL_EQUAL)
        } else {
            s.AddToken(TK_EQUAL)
        }
    case c == '!':
        if s.Match("=") {
            s.AddToken(TK_BANG_EQUAL)
        } else {
            s.AddToken(TK_BANG)
        }
    case c == '>':
        if s.Match("=") {
            s.AddToken(TK_GREATER_EQUAL)
        } else {
            s.AddToken(TK_GREATER)
        }
    case c == '<':
        if s.Match("=") {
            s.AddToken(TK_LESS_EQUAL)
        } else {
            s.AddToken(TK_LESS)
        }
    case c == '/':
        if s.Match("/") {
            s.SkipOneLineComment()
        } else {
            s.AddToken(TK_SLASH)
        }
    case c == '"':
        if s.FindChar('"') {
            s.AddToken(TK_STRING)
        }
    case c >= '0' && c <= '9':
        s.ResolveNum()
    case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_':
        s.ResolveIdAndKeyword()
    default:
        s.Report("Unexpected character: %c", c)
    }
}

func (s *Scanner) ResolveIdAndKeyword() {
    isIdChar := func (ch rune) bool {
        return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_'
    }

    for isIdChar(s.Peek()) {
        s.Current++
    }

    if s.ReservedKws[s.String()] != "" {
        s.AddToken(s.ReservedKws[s.String()])
    } else {
        s.AddToken(TK_IDENTIFIER)
    }
}


func (s *Scanner) String() string {
    return s.Source[s.Start : s.Current]
}


func (s *Scanner) ResolveNum() {
    isDigit := func (ch rune) bool {
        return ch >= '0' && ch <= '9'
    }

    for isDigit(s.Peek()) {
        s.Current++
    }

    if s.Peek() == '.' && isDigit(s.PeekNext()) {
        s.Current++
    }

    for isDigit(s.Peek()) {
        s.Current++
    }

    s.AddToken(TK_NUMBER)
}

func (s *Scanner) Peek() rune {
    if s.IsEnd() {
        return 0
    }
    return rune(s.Source[s.Current])
}

func (s *Scanner) PeekNext() rune {
    if s.Current + 1 >= len(s.Source) {
        return 0
    }

    return rune(s.Source[s.Current + 1])
}

func (s *Scanner) FindChar(ch byte) bool {
    isFind := false 
    for !s.IsEnd() {
        if s.Source[s.Current] == ch {
            isFind = true 
            s.Current++
            break
        }
        // TODO: maybe handle newline
        s.Current++
    }

    if !isFind {
        s.Report("Unterminated string.")
    }
    
    return isFind
}

func (s *Scanner) SkipOneLineComment() {
    for !s.IsEnd() {
        if s.Source[s.Current] == '\n' {
            s.Line++
            s.Current++ // Skip the newline
            break
        }
        s.Current++
    }
}

func (s *Scanner) Match(expect string) bool {
    if s.Current + len(expect) > len(s.Source) {
        return false
    }

    if expect == s.Source[s.Current : s.Current + len(expect)] {
        s.Current += len(expect)
        return true
    }

    return false
}

func (s *Scanner) Report(format string, args ...any) {
    s.HasErr = true
    fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", s.Line, fmt.Sprintf(format, args...))
}
