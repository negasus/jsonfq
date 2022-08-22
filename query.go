package jsonfq

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/negasus/jsonfq/ast"
)

type Query struct {
	stmts []ast.Stmt
}

func (q *Query) Marshal() ([]byte, error) {
	if len(q.stmts) == 0 {
		return nil, nil
	}

	var res []byte

	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, uint64(len(q.stmts)))
	res = append(res, buf[:n]...)

	for _, s := range q.stmts {
		res = append(res, byte(s.GetStmtType()))
		d, err := s.Marshal()
		if err != nil {
			return nil, err
		}
		res = append(res, d...)
	}

	return res, nil
}

func (q *Query) Unmarshal(data []byte) error {
	buf := bytes.NewBuffer(data)

	count, errCount := binary.ReadUvarint(buf)
	if errCount != nil {
		return errCount
	}

	for i := 0; i < int(count); i++ {
		t, errT := buf.ReadByte()
		if errT != nil {
			return errT
		}
		var s ast.Stmt
		switch ast.StmtType(t) {
		case ast.StmtTypeMapKey:
			s = &ast.StmtMapKey{}
		case ast.StmtTypeArrayIndex:
			s = &ast.StmtArrayIndex{}
		case ast.StmtTypeFilter:
			s = &ast.StmtFilter{}
		default:
			return fmt.Errorf("unexpected statement type %v", t)
		}

		e := s.Unmarshal(buf)
		if e != nil {
			return e
		}
		q.stmts = append(q.stmts, s)
	}

	return nil
}
