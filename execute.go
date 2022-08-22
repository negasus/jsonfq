package jsonfq

import (
	"errors"
	"fmt"

	"github.com/negasus/jsonfq/ast"
)

type filterElement struct {
	keyStart   int
	keyEnd     int
	valueStart int
	valueEnd   int
	arrayIndex int
}

func selectIndexes(data []byte, idx int, stmts []ast.Stmt) (int, int, error) {
	var err error
	var startIdx, endIdx int

	for _, st := range stmts {
		if idx >= len(data) {
			return 0, idx, ErrUnexpectedEndOfData
		}

		switch st.(type) {
		case *ast.StmtMapKey:
			startIdx, endIdx, err = findMapElement(data, idx, st.(*ast.StmtMapKey).Value)
			if err != nil {
				if errors.Is(err, ErrMapFieldNotFound) {
					return 0, idx, ErrMapFieldNotFound
				}
				return 0, idx, fmt.Errorf("position: %d, error find map element %q, %w", endIdx, st.(*ast.StmtMapKey).Value, err)
			}
		case *ast.StmtArrayIndex:
			startIdx, endIdx, err = findArrayElement(data, idx, st.(*ast.StmtArrayIndex).Value)
			if err != nil {
				if errors.Is(err, ErrArrayElementNotFound) {
					return 0, idx, ErrArrayElementNotFound
				}
				return 0, idx, fmt.Errorf("position: %d, error find array element %d, %w", endIdx, st.(*ast.StmtArrayIndex).Value, err)
			}

		default:
			return 0, idx, fmt.Errorf("unexpected stmt type %T", st)
		}

		idx = startIdx
	}

	return startIdx, endIdx + 1, nil
}

func executeStmt(data []byte, idx int, n ast.Stmt) ([]byte, error) {
	var err error
	var start, end int

	switch v := n.(type) {
	case *ast.StmtFilter:
		data, err = filter(data, v.Expr)
		if err != nil {
			return nil, fmt.Errorf("error filter data, %w", err)
		}
		return data, nil
	case *ast.StmtArrayIndex:
		start, end, err = findArrayElement(data, idx, v.Value)
		if err != nil {
			if errors.Is(err, ErrArrayElementNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("error find array element, %w", err)
		}
		return data[start : end+1], nil
	case *ast.StmtMapKey:
		start, end, err = findMapElement(data, idx, v.Value)
		if err != nil {
			if errors.Is(err, ErrMapFieldNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("error find map element, %w", err)
		}
		return data[start : end+1], nil
	default:
		return nil, fmt.Errorf("unknown statement type: %T", n)
	}
}

func executeExp(data []byte, el filterElement, e ast.Expr) (ast.Expr, error) {
	switch v := e.(type) {
	case *ast.ExprGetBlock:
		return executeGetBlock(data, el, v)
	case *ast.ExprValue:
		return v, nil
	case *ast.ExprFn:
		return executeFn(data, el, v)
	case *ast.ExprBinaryOp:
		return executeBinaryOp(data, el, v)
	default:
		return nil, fmt.Errorf("unknown expression type: %T", v)
	}
}

func executeGetBlock(data []byte, el filterElement, n *ast.ExprGetBlock) (ast.Expr, error) {
	var d []byte
	switch n.T {
	case ast.ExprGetBlockTypeValue:
		d = data[el.valueStart : el.valueEnd+1]
	case ast.ExprGetBlockTypeKey:
		d = data[el.keyStart : el.keyEnd+1]
	default:
		d = data
	}

	start, end, err := selectIndexes(d, 0, n.Stmts)
	if err != nil {
		return nil, fmt.Errorf("error select indexes, %w", err)
	}

	return &ast.ExprValue{T: ast.ExprValueTypeBytes, D: d[start:end]}, nil
}

func executeBinaryOp(data []byte, el filterElement, n *ast.ExprBinaryOp) (ast.Expr, error) {
	var l, r ast.Expr
	var err error
	l, err = executeExp(data, el, n.Left)
	if err != nil {
		return nil, fmt.Errorf("error execute left expression, %w", err)
	}
	r, err = executeExp(data, el, n.Right)
	if err != nil {
		return nil, fmt.Errorf("error execute right expression, %w", err)
	}

	if l.GetExprType() != ast.ExprTypeValue {
		return nil, fmt.Errorf("left expression must be ExprValue type, got %T", l)
	}

	if r.GetExprType() != ast.ExprTypeValue {
		return nil, fmt.Errorf("right expression must be ExprValue type, got %T", l)
	}

	ll := l.(*ast.ExprValue)
	rr := r.(*ast.ExprValue)

	if ll.T != rr.T {
		return nil, fmt.Errorf("left and right expression must be same type, got %q and %q", ll.T, rr.T)
	}

	var ors bool

	switch n.Op {
	case "+":
		err = ll.ADD(rr)
		return ll, err
	case "-":
		err = ll.SUB(rr)
		return ll, err
	case "*":
		err = ll.MUL(rr)
		return ll, err
	case "/":
		err = ll.DIV(rr)
		return ll, err
	case "=":
		ors = ll.EQ(rr)
	case "!=":
		ors = ll.NE(rr)
	case ">":
		ors, err = ll.GT(rr)
	case ">=":
		ors, err = ll.GTE(rr)
	case "<":
		ors, err = ll.LT(rr)
	case "<=":
		ors, err = ll.LTE(rr)
	default:
		return nil, fmt.Errorf("unexpected binary operator, %q", n.Op)
	}
	if err != nil {
		return nil, fmt.Errorf("error execute binary operator for op %q, %w", n.Op, err)
	}

	return &ast.ExprValue{T: ast.ExprValueTypeBool, B: ors}, nil

}
