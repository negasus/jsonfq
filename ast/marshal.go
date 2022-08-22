package ast

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

// todo: WIP

func (st *StmtMapKey) Marshal() ([]byte, error) {
	var res []byte
	buf := make([]byte, 10)
	n := binary.PutUvarint(buf, uint64(st.Position))
	res = append(res, buf[:n]...)
	n = binary.PutUvarint(buf, uint64(len(st.Value)))
	res = append(res, buf[:n]...)
	return append(res, st.Value...), nil
}

func (st *StmtMapKey) Unmarshal(data *bytes.Buffer) error {
	v, e := binary.ReadUvarint(data)
	if e != nil {
		return e
	}
	st.Position = int(v)
	v, e = binary.ReadUvarint(data)
	if e != nil {
		return e
	}
	b := make([]byte, int(v))
	_, e = data.Read(b)
	if e != nil {
		return e
	}
	st.Value = string(b)
	return nil
}

func (st *StmtArrayIndex) Marshal() ([]byte, error) {
	var res []byte
	buf := make([]byte, 10)
	n := binary.PutUvarint(buf, uint64(st.Position))
	res = append(res, buf[:n]...)
	n = binary.PutUvarint(buf, uint64(st.Value))
	return append(res, buf[:n]...), nil
}

func (st *StmtArrayIndex) Unmarshal(data *bytes.Buffer) error {
	v, e := binary.ReadUvarint(data)
	if e != nil {
		return e
	}
	st.Position = int(v)
	v, e = binary.ReadUvarint(data)
	if e != nil {
		return e
	}
	st.Value = int(v)
	return nil
}

func (st *StmtFilter) Marshal() ([]byte, error) {
	var res []byte
	buf := make([]byte, 10)
	n := binary.PutUvarint(buf, uint64(st.Position))
	res = append(res, buf[:n]...)
	d, err := st.Expr.Marshal()
	if err != nil {
		return nil, err
	}
	return append(res, d...), nil
}

func (st *StmtFilter) Unmarshal(data *bytes.Buffer) error {
	v, e := binary.ReadUvarint(data)
	if e != nil {
		return e
	}
	st.Position = int(v)
	v, e = binary.ReadUvarint(data)
	if e != nil {
		return e
	}
	var ex Expr
	switch ExprType(v) {
	case ExprTypeValue:
		ex = &ExprValue{}
	case ExprTypeGetBlock:
		ex = &ExprGetBlock{}
	case ExprTypeFn:
		ex = &ExprFn{}
	case ExprTypeBinaryOp:
		ex = &ExprBinaryOp{}
	case ExprTypeUnaryOp:
		ex = &ExprUnaryOp{}
	}
	e = ex.Unmarshal(data)
	if e != nil {
		return e
	}
	st.Expr = ex
	return nil
}

func (e *ExprValue) Marshal() ([]byte, error) {
	var res []byte
	buf := make([]byte, 10)
	// Position
	n := binary.PutUvarint(buf, uint64(e.Position))
	res = append(res, buf[:n]...)
	// T
	res = append(res, byte(e.T))
	// S
	n = binary.PutUvarint(buf, uint64(len(e.S)))
	res = append(res, buf[:n]...)
	res = append(res, e.S...)
	// B
	if e.B {
		res = append(res, 0x01)
	} else {
		res = append(res, 0x00)
	}
	// I
	n = binary.PutVarint(buf, int64(e.I))
	res = append(res, buf[:n]...)
	// F
	ff := strconv.FormatFloat(e.F, 'f', 16, 64)
	n = binary.PutUvarint(buf, uint64(len(ff)))
	res = append(res, buf[:n]...)
	res = append(res, ff...)
	// D
	n = binary.PutUvarint(buf, uint64(len(e.D)))
	res = append(res, buf[:n]...)
	res = append(res, e.D...)

	return res, nil
}
func (e *ExprValue) Unmarshal(data *bytes.Buffer) error {
	return nil
}

func (e *ExprGetBlock) Marshal() ([]byte, error) {
	var res []byte
	return res, nil
}
func (e *ExprGetBlock) Unmarshal(data *bytes.Buffer) error {
	return nil
}

func (e *ExprFn) Marshal() ([]byte, error) {
	var res []byte
	return res, nil
}
func (e *ExprFn) Unmarshal(data *bytes.Buffer) error {
	return nil
}

func (e *ExprBinaryOp) Marshal() ([]byte, error) {
	var res []byte
	return res, nil
}
func (e *ExprBinaryOp) Unmarshal(data *bytes.Buffer) error {
	return nil
}

func (e *ExprUnaryOp) Marshal() ([]byte, error) {
	var res []byte
	return res, nil
}
func (e *ExprUnaryOp) Unmarshal(data *bytes.Buffer) error {
	return nil
}
