package main

import (
	"fmt"
//	"reflect"
)

type Expr interface {
    // the String method is used to print the expression recursively.
    // for example: (*(+ 1 2) 3)
    String() string

    // the Eval method is used to evaluate the result of expression recursively.
    // in other implementations, maybe use the visitor pattern of AST
    Eval() (ValueType, error)
}

// Variable expression. for example: a + 123
type VarExpr struct {
    token *Token
}

func (e VarExpr) String() string {
    return fmt.Sprintf("(var %s)", e.token.Lexeme)
}

func (e VarExpr) Eval() (ValueType, error) {
    if v, ok := globals[e.token.Lexeme]; ok {
        return v, nil
    } else {
        return NilValue, fmt.Errorf("Undefined variable '%s'.", e.token.Lexeme)
    }
}

// Assignment expression. for example: a = 123
type AssignmentExpr struct {
    token *Token
    expr Expr
}

func (e AssignmentExpr) String() string {
    return fmt.Sprintf("(= %s %s)", e.token.Lexeme, e.expr.String())
}

func (e AssignmentExpr) Eval() (ValueType, error) {
    if _, ok := globals[e.token.Lexeme]; ok {
        if v, err := e.expr.Eval(); err != nil {
            return NilValue, err
        } else {
            globals[e.token.Lexeme] = v
            return v, nil
        }
    }

    return NilValue, fmt.Errorf("Undefined variable '%s'.", e.token.Lexeme)
}

// Literal expression. for example: true, false, nil, 123, "abc"
type LiteralExpr struct {
    token *Token
}

func (e LiteralExpr) String() string {
    switch e.token.Type {
    case TK_NUMBER, TK_STRING:
        return e.token.LiteralString()
    default:
        return e.token.Lexeme
    }
}

func (e LiteralExpr) Eval() (ValueType, error) {
    return e.token.Literal()
}

// Group expression. for example: ("abc"), (1+2)
type GroupExpr struct {
    expr Expr
}

func (e GroupExpr) String() string {
    return fmt.Sprintf("(group %s)", e.expr)
}

func (e GroupExpr) Eval() (ValueType, error) {
    return e.expr.Eval()
}

// Unary expression. for example: -1, !a==b
type UnaryExpr struct {
    expr Expr
    token *Token
}

func (e UnaryExpr) String() string {
    return fmt.Sprintf("(%s %s)", e.token.Lexeme, e.expr)
}

func (e UnaryExpr) Eval() (ValueType, error) {
    val, err := e.expr.Eval()
    if err != nil {
        return nil, err
    }
    switch(e.token.Type) {
    case TK_MINUS:
        if v, ok := val.(NumberType); ok {
            return NumberType{v: -v.v}, nil
        }
    case TK_BANG:
        return BoolType{v:!IsTruthy(val)}, nil
    }

    return nil, fmt.Errorf("Unknown unary operator: %s", e.token.Lexeme)
}

func IsTruthy(val ValueType) bool {
    return val.IsTrue()
}

// Binary expression. for example: 1+2, 3*4
type BinaryExpr struct {
    left Expr
    right Expr
    optr *Token
}

func (e BinaryExpr) String() string {
    return fmt.Sprintf("(%s %s %s)", e.optr.Lexeme, e.left, e.right)
}

func (e BinaryExpr) Eval() (ValueType, error) {
    lhs, err := e.left.Eval()
    if err != nil {
        return nil, err
    }

    rhs, err := e.right.Eval()
    if err != nil {
        return nil, err
    }

    switch (e.optr.Type) {
    case TK_PLUS:
        return EvalIfMatch(
            lhs,
            rhs,
            EvalPlus[NumberType],
            EvalPlus[StringType],
        )
    case TK_MINUS:
        return EvalIfMatch(
            lhs,
            rhs,
            EvalMinus[NumberType],
        )
    case TK_STAR:
        return EvalIfMatch(
            lhs,
            rhs,
            EvalMul[NumberType],
        )
    case TK_SLASH:
        return EvalIfMatch(
            lhs,
            rhs,
            EvalDiv[NumberType],
        )
    case TK_BANG_EQUAL:
        if lhs.Type() != rhs.Type() {
            return TrueValue, nil
        }

        return EvalIfMatch(
            lhs,
            rhs,
            EvalBangEqual[NilType],
            EvalBangEqual[BoolType],
            EvalBangEqual[StringType],
            EvalBangEqual[NumberType],
        )
    case TK_EQUAL_EQUAL:
        if lhs.Type() != rhs.Type() {
            return FalseValue, nil
        }

        return EvalIfMatch(
            lhs,
            rhs,
            EvalEqualEqual[NilType],
            EvalEqualEqual[BoolType],
            EvalEqualEqual[StringType],
            EvalEqualEqual[NumberType],
        )
    case TK_LESS:
        return EvalIfMatch(
            lhs,
            rhs,
            EvalLess[NumberType],
            EvalLess[StringType],
        )
    case TK_LESS_EQUAL:
        return EvalIfMatch(
            lhs,
            rhs,
            EvalLessEqual[NumberType],
            EvalLessEqual[StringType],
        )
    case TK_GREATER:
        return EvalIfMatch(
            lhs,
            rhs,
            EvalGreater[NumberType],
            EvalGreater[StringType],
        )
    case TK_GREATER_EQUAL:
        return EvalIfMatch(
            lhs,
            rhs,
            EvalGreaterEqual[NumberType],
            EvalGreaterEqual[StringType],
        )
    }

    return nil, fmt.Errorf("Unknown binary operator: %s", e.optr.Lexeme)
}


func EvalIfMatch(lhs, rhs ValueType, functors ...func(ValueType, ValueType)(ValueType, error)) (ValueType, error) {
    for _, functor := range(functors) {
        res, err := functor(lhs, rhs)
        if err == nil {
            return res, nil
        }
    }

    return NilValue, fmt.Errorf("Unmatch functors")
}


type GenericType interface {
    NilType | NumberType | StringType | BoolType
}


type  ComparableType interface {
    NumberType | StringType
}


func EvalPlus[T NumberType | StringType](lhs, rhs ValueType) (ValueType, error) {
    return EvalGeneric[T](lhs, rhs, func(v1, v2 ValueType) ValueType {
        switch v1.Type() {
        case VT_Number:
            return NumberType{v: v1.(NumberType).v + v2.(NumberType).v}
        case VT_String:
            return StringType{v: v1.(StringType).v + v2.(StringType).v}
       }

        return NilValue
    })
}


func EvalMinus[T NumberType](lhs, rhs ValueType) (ValueType, error) {
    return EvalGeneric[T](lhs, rhs, func(v1, v2 ValueType) ValueType {
        return NumberType{v: v1.(NumberType).v - v2.(NumberType).v}
    })
}


func EvalMul[T NumberType](lhs, rhs ValueType) (ValueType, error) {
    return EvalGeneric[T](lhs, rhs, func(v1, v2 ValueType) ValueType {
        return NumberType{v: v1.(NumberType).v * v2.(NumberType).v}
    })
}


func EvalDiv[T NumberType](lhs, rhs ValueType) (ValueType, error) {
    return EvalGeneric[T](lhs, rhs, func(v1, v2 ValueType) ValueType {
        return NumberType{v: v1.(NumberType).v / v2.(NumberType).v}
    })
}


func EvalLess[T ComparableType](lhs, rhs ValueType) (ValueType, error) {
    return EvalGeneric[T](lhs, rhs, func(v1, v2 ValueType) ValueType {
        switch v1.Type() {
        case VT_Number:
            return BoolType{v: v1.(NumberType).v < v2.(NumberType).v}
        case VT_String:
            return BoolType{v: v1.(StringType).v < v2.(StringType).v}
       }

        return FalseValue
    })
}


func EvalLessEqual[T ComparableType](lhs, rhs ValueType) (ValueType, error) {
    return EvalGeneric[T](lhs, rhs, func(v1, v2 ValueType) ValueType {
        switch v1.Type() {
        case VT_Number:
            return BoolType{v: v1.(NumberType).v <= v2.(NumberType).v}
        case VT_String:
            return BoolType{v: v1.(StringType).v <= v2.(StringType).v}
       }

        return FalseValue
    })
}


func EvalGreater[T ComparableType](lhs, rhs ValueType) (ValueType, error) {
    return EvalGeneric[T](lhs, rhs, func(v1, v2 ValueType) ValueType {
        switch v1.Type() {
        case VT_Number:
            return BoolType{v: v1.(NumberType).v > v2.(NumberType).v}
        case VT_String:
            return BoolType{v: v1.(StringType).v > v2.(StringType).v}
       }

        return FalseValue
    })
}

func EvalGreaterEqual[T ComparableType](lhs, rhs ValueType) (ValueType, error) {
    return EvalGeneric[T](lhs, rhs, func(v1, v2 ValueType) ValueType {
        switch v1.Type() {
        case VT_Number:
            return BoolType{v: v1.(NumberType).v >= v2.(NumberType).v}
        case VT_String:
            return BoolType{v: v1.(StringType).v >= v2.(StringType).v}
       }

        return FalseValue
    })
}


func EvalEqualEqual[T GenericType](lhs, rhs ValueType) (ValueType, error) {
    return EvalGeneric[T](lhs, rhs, func(v1, v2 ValueType) ValueType {
        switch v1.Type() {
        case VT_Nil:
            return FalseValue
        case VT_Bool:
            return BoolType{v: v1.(BoolType).v == v2.(BoolType).v}
        case VT_Number:
            return BoolType{v: v1.(NumberType).v == v2.(NumberType).v}
        case VT_String:
            return BoolType{v: v1.(StringType).v == v2.(StringType).v}
       }

        return FalseValue
    })
}


func EvalBangEqual[T GenericType](lhs, rhs ValueType) (ValueType, error) {
    // Hack: type of T in generic function cannot use T.field or T.method
    // we have to provide interfaces of T to acquire type name and value
    return EvalGeneric[T](lhs, rhs, func(v1, v2 ValueType) ValueType {
        switch v1.Type() {
        case VT_Nil:
            return TrueValue
        case VT_Bool:
            return BoolType{v: v1.(BoolType).v != v2.(BoolType).v}
        case VT_Number:
            return BoolType{v: v1.(NumberType).v != v2.(NumberType).v}
        case VT_String:
            return BoolType{v: v1.(StringType).v != v2.(StringType).v}
       }

        return FalseValue
    })
}

func EvalGeneric[T GenericType](lhs, rhs ValueType, op func(ValueType, ValueType) ValueType) (ValueType, error){
    var result ValueType

    if _, ok := lhs.(T); !ok {
        return result, fmt.Errorf("can't convert '%v'", lhs)
    }

    if _, ok := rhs.(T); !ok {
        return result, fmt.Errorf("can't convert '%v'", rhs)
    }

    result = op(lhs, rhs)

    return result, nil
}
