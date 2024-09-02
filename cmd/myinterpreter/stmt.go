package main

import (
    "fmt"
)


type Stmt interface {
    Run() error
}


type PrintStmt struct {
    v Expr
}


func (s PrintStmt) Run() error {
    v, err := s.v.Eval()
    if err != nil {
        return err
    }

    fmt.Println(v)
    return nil
}


type ExprStmt struct {
    v Expr
}


func (s ExprStmt) Run() error {
    _, err := s.v.Eval()
    return err
}


var globals = make(map[string]ValueType)


type VarStmt struct {
    v Expr
    tk *Token
}


func (s VarStmt) Run() error {
    var err error
    var v ValueType

    if s.v == nil {
        v = NilValue
    } else {
        if v, err = s.v.Eval(); err != nil {
            return err
        }
    }

    // store the variable
    globals[s.tk.Lexeme] = v

    return nil
}

type BlockStmt struct {
    stmts []Stmt

    env map[string]ValueType
}

// TODO: optimize performance
// maybe we can use a local mapping in block statement
func (s *BlockStmt) EnterBlock() {
    // deep copy
    s.env = make(map[string]ValueType)
    for k, v := range globals {
        s.env[k] = v
    }
}

func (s *BlockStmt) ExitBlock() {
    globals = s.env
}

func (s BlockStmt) Run() error {
    s.EnterBlock()

    for _, stmt := range(s.stmts) {
        if err := stmt.Run(); err != nil {
            return err
        }
    }

    s.ExitBlock()

    return nil
}
