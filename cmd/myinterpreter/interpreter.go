package main


type Interpreter struct {
    expr Expr
}

func NewInterpreter(expr Expr) *Interpreter {
    return &Interpreter {
        expr: expr,
    }
}

func (i *Interpreter) Eval() (ValueType, error) {
    return i.expr.Eval()
}
