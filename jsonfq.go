package jsonfq

import (
	"fmt"

	"github.com/negasus/jsonfq/parser"
)

var (
	ErrUnexpectedEndOfData  = fmt.Errorf("unexpected end of data")
	ErrMapFieldNotFound     = fmt.Errorf("map field not found")
	ErrArrayElementNotFound = fmt.Errorf("array element not found")
)

func ParseAndSelect(data []byte, s string) ([]byte, error) {
	q, err := Parse(s)
	if err != nil {
		return nil, fmt.Errorf("error parse query, %w", err)
	}
	return Select(data, q)
}

func Parse(q string) (*Query, error) {
	stmts, err := parser.Parse([]byte(q))
	if err != nil {
		return nil, err
	}
	return &Query{stmts: stmts}, nil
}

func Select(data []byte, q *Query) ([]byte, error) {
	var err error
	for _, st := range q.stmts {
		data, err = executeStmt(data, 0, st)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func Compact(data []byte) []byte {
	var idx int
	for i, b := range data {
		if isSpace(b) {
			continue
		}
		data[idx] = data[i]
		idx++
	}
	data = data[:idx:idx]
	return data
}

func skipSpaces(data []byte, idx int) (int, error) {
	for {
		if idx >= len(data) {
			return idx, ErrUnexpectedEndOfData
		}
		if isSpace(data[idx]) {
			idx++
			continue
		}
		return idx, nil
	}
}

func skipColonWithSpaces(data []byte, idx int) (int, error) {
	var found bool
	for {
		if idx >= len(data) {
			return idx, ErrUnexpectedEndOfData
		}
		if data[idx] == ':' {
			found = true
			idx++
			continue
		}
		if isSpace(data[idx]) {
			idx++
			continue
		}
		if !found {
			return idx, errUnexpectedSymbol(data[idx], "colon or space")
		}
		return idx, nil
	}
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func errUnexpectedSymbol(sym byte, expect string) error {
	return fmt.Errorf("unexpected symbol %c, expect %q", sym, expect)
}
