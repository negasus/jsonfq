package ast

import (
	"bytes"
)

type ExprType int

const (
	ExprTypeValue ExprType = iota
	ExprTypeGetBlock
	ExprTypeFn
	ExprTypeBinaryOp
	ExprTypeUnaryOp
)

type Expr interface {
	GetExprType() ExprType
	GetPosition() int
	Marshal() ([]byte, error)
	Unmarshal(*bytes.Buffer) error
}

type ExprValueType int

func (e ExprValueType) String() string {
	return exprValueTypeS[e]
}

const (
	ExprValueTypeString ExprValueType = iota
	ExprValueTypeBool
	ExprValueTypeInt
	ExprValueTypeFloat
	ExprValueTypeBytes
)

var (
	exprValueTypeS = [...]string{
		ExprValueTypeString: "string",
		ExprValueTypeBool:   "bool",
		ExprValueTypeInt:    "int",
		ExprValueTypeFloat:  "float",
		ExprValueTypeBytes:  "bytes",
	}
)

type ExprValue struct {
	Position int
	T        ExprValueType
	S        string
	B        bool
	I        int
	F        float64
	D        []byte
}

func (*ExprValue) GetExprType() ExprType { return ExprTypeValue }
func (e *ExprValue) GetPosition() int    { return e.Position }

type ExprGetBlockType int

const (
	ExprGetBlockTypeKey ExprGetBlockType = iota
	ExprGetBlockTypeValue
)

type ExprGetBlock struct {
	Position int
	T        ExprGetBlockType
	Stmts    []Stmt
}

func (*ExprGetBlock) GetExprType() ExprType { return ExprTypeGetBlock }
func (e *ExprGetBlock) GetPosition() int    { return e.Position }

type ExprFn struct {
	Position int
	Name     string
	Args     []Expr
}

func (*ExprFn) GetExprType() ExprType { return ExprTypeFn }
func (e *ExprFn) GetPosition() int    { return e.Position }

type ExprBinaryOp struct {
	Position int
	Op       string
	Left     Expr
	Right    Expr
}

func (*ExprBinaryOp) GetExprType() ExprType { return ExprTypeBinaryOp }
func (e *ExprBinaryOp) GetPosition() int    { return e.Position }

type ExprUnaryOp struct {
	Position int
	Op       string
	Right    Expr
}

func (*ExprUnaryOp) GetExprType() ExprType { return ExprTypeUnaryOp }
func (e *ExprUnaryOp) GetPosition() int    { return e.Position }
