package jsonfq

import (
	"fmt"
	"strconv"

	"github.com/negasus/jsonfq/ast"
)

func executeFn(data []byte, el filterElement, e *ast.ExprFn) (ast.Expr, error) {
	switch e.Name {
	case "int":
		return fnInt(data, el, e.Args)
	default:
		return nil, fmt.Errorf("unknown function: %s", e.Name)
	}
}

func fnInt(data []byte, el filterElement, args []ast.Expr) (ast.Expr, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("fnInt expects 1 argument, got %d", len(args))
	}

	v, err := executeExp(data, el, args[0])
	if err != nil {
		return nil, fmt.Errorf("error execute argument for fnInt, %w", err)
	}

	s, ok := v.(*ast.ExprValue)
	if !ok {
		return nil, fmt.Errorf("fnInt expects argument to be ident, got %T", v)
	}
	if s.T != ast.ExprValueTypeBytes {
		return nil, fmt.Errorf("fnInt expects argument to be ident of type identTypeBytes, got %v", s.T)
	}
	// try to unescape string, if it fails, continue with string as is
	var ss string
	ss, err = strconv.Unquote(string(s.D))
	if err != nil {
		ss = string(s.D)
	}
	var i int
	i, err = strconv.Atoi(ss)
	if err != nil {
		return nil, fmt.Errorf("error convert string to int for fnInt, %w", err)
	}

	return &ast.ExprValue{T: ast.ExprValueTypeInt, I: i}, nil
}
