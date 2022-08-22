package ast

import (
	"bytes"
	"fmt"
)

// Arithmetic operations with two ExprValue

func (e *ExprValue) ADD(v *ExprValue) error {
	if e.T == ExprValueTypeBool {
		return fmt.Errorf("ExprValue.ADD not support Bool type")
	}
	e.S += v.S
	e.I += v.I
	e.F += v.F
	e.D = append(e.D, v.D...)
	return nil
}

func (e *ExprValue) SUB(v *ExprValue) error {
	if e.T != ExprValueTypeInt && e.T != ExprValueTypeFloat {
		return fmt.Errorf("ExprValue.SUB operation supports only int and float")
	}
	e.I -= v.I
	e.F -= v.F
	return nil
}

func (e *ExprValue) MUL(v *ExprValue) error {
	if e.T != ExprValueTypeInt && e.T != ExprValueTypeFloat {
		return fmt.Errorf("ExprValue.MUL operation supports only int and float")
	}
	e.I *= v.I
	e.F *= v.F
	return nil
}

func (e *ExprValue) DIV(v *ExprValue) error {
	if e.T != ExprValueTypeInt && e.T != ExprValueTypeFloat {
		return fmt.Errorf("ExprValue.DIV operation supports only int and float")
	}
	e.I /= v.I
	e.F /= v.F
	return nil
}

func (e *ExprValue) EQ(v *ExprValue) bool {
	return e.S == v.S && e.B == v.B && e.I == v.I && e.F == v.F && bytes.Equal(e.D, v.D)
}

func (e *ExprValue) NE(v *ExprValue) bool {
	return e.S != v.S || e.B != v.B || e.I != v.I || e.F != v.F || !bytes.Equal(e.D, v.D)
}

func (e *ExprValue) GT(v *ExprValue) (bool, error) {
	switch e.T {
	case ExprValueTypeInt:
		return e.I > v.I, nil
	case ExprValueTypeFloat:
		return e.F > v.F, nil
	default:
		return false, fmt.Errorf("unsupported ExprValue type")
	}
}

func (e *ExprValue) GTE(v *ExprValue) (bool, error) {
	switch e.T {
	case ExprValueTypeInt:
		return e.I >= v.I, nil
	case ExprValueTypeFloat:
		return e.F >= v.F, nil
	default:
		return false, fmt.Errorf("unsupported ExprValue type")
	}
}

func (e *ExprValue) LT(v *ExprValue) (bool, error) {
	switch e.T {
	case ExprValueTypeInt:
		return e.I < v.I, nil
	case ExprValueTypeFloat:
		return e.F < v.F, nil
	default:
		return false, fmt.Errorf("unsupported ExprValue type")
	}
}

func (e *ExprValue) LTE(v *ExprValue) (bool, error) {
	switch e.T {
	case ExprValueTypeInt:
		return e.I <= v.I, nil
	case ExprValueTypeFloat:
		return e.F <= v.F, nil
	default:
		return false, fmt.Errorf("unsupported ExprValue type")
	}
}
