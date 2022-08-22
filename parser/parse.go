package parser

import (
	"fmt"

	"github.com/negasus/jsonfq/ast"
)

func Parse(q []byte) ([]ast.Stmt, error) {
	sc := newScanner(q)
	p := yyNewParser()
	ret := p.Parse(sc)

	if ret != 0 {
		return nil, fmt.Errorf("parse error")
	}

	if len(sc.errors) > 0 {
		return nil, fmt.Errorf("parse error: %v", sc.errors)
	}

	return sc.stmts, nil
}
