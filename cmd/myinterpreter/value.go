package main

import (
    "fmt"
)

var VT_Nil = "nil"
var VT_Bool = "bool"
var VT_String = "string"
var VT_Number = "number"

type ValueType interface {
    String() string
    Literal() any
    Type() string
    IsTrue() bool
}

type NilType struct {
    v ValueType // v is nil by default
}

func (t NilType) String() string {
    return "nil"
}

func (t NilType) Literal() any {
    return nil
}

func (t NilType) Type() string {
    return "nil"
}

func (t NilType) IsTrue() bool {
    return false
}

type BoolType struct {
    v bool
}

func (t BoolType) String() string {
    if t.v {
        return "true"
    } else {
        return "false"
    }
}

func (t BoolType) Literal() any {
    return t.v
}

func (t BoolType) Type() string {
    return "bool"
}

func (t BoolType) IsTrue() bool {
    return t.v
}

type StringType struct {
    v string
}

func (t StringType) String() string {
    return t.v
}

func (t StringType) Literal() any {
    return t.v
}

func (t StringType) Type() string {
    return "string"
}

func (t StringType) IsTrue() bool {
    return len(t.v) > 0
}

type NumberType struct {
    v float64
}

func (t NumberType) String() string {
    return fmt.Sprintf("%v", t.v)
}

func (t NumberType) Literal() any {
    return t.v
}

func (t NumberType) Type() string {
    return "number"
}

func (t NumberType) IsTrue() bool {
    return t.v != 0
}

var NilValue = NilType{}
var TrueValue = BoolType{v: true}
var FalseValue = BoolType{v: false}
var EmptyStringValue = StringType{v: ""}
var ZeroNumberValue = NumberType{v: 0}
