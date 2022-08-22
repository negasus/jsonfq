package ast

import (
	"bytes"
)

type StmtType int

const (
	StmtTypeMapKey StmtType = iota
	StmtTypeArrayIndex
	StmtTypeFilter
)

type Stmt interface {
	GetStmtType() StmtType
	GetPosition() int
	Marshal() ([]byte, error)
	Unmarshal(*bytes.Buffer) error
}

type StmtMapKey struct {
	Position int
	Value    string
}

func (*StmtMapKey) GetStmtType() StmtType { return StmtTypeMapKey }
func (st *StmtMapKey) GetPosition() int   { return st.Position }

type StmtArrayIndex struct {
	Position int
	Value    int
}

func (*StmtArrayIndex) GetStmtType() StmtType { return StmtTypeArrayIndex }
func (st *StmtArrayIndex) GetPosition() int   { return st.Position }

type StmtFilter struct {
	Position int
	Expr     Expr
}

func (*StmtFilter) GetStmtType() StmtType { return StmtTypeFilter }
func (st *StmtFilter) GetPosition() int   { return st.Position }
