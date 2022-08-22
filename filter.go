package jsonfq

import (
	"fmt"

	"github.com/negasus/jsonfq/ast"
)

func filter(data []byte, e ast.Expr) ([]byte, error) {
	if len(data) == 0 {
		return nil, ErrUnexpectedEndOfData
	}
	if data[0] != '{' && data[0] != '[' {
		return nil, fmt.Errorf("filter can be applied only to array or object, unexpected symbol: %c", data[0])
	}

	startSym, endSym := byte('{'), byte('}')
	if data[0] == '[' {
		startSym, endSym = '[', ']'
	}
	idx := 1
	var err error
	var arrayIndex int
	var ok bool

	var resultElements []filterElement

	for {
		if idx >= len(data) {
			return nil, ErrUnexpectedEndOfData
		}

		el := filterElement{-1, -1, -1, -1, -1}

		idx, err = skipSpaces(data, idx)
		if err != nil {
			return nil, fmt.Errorf("error skip at block start spaces, %w", err)
		}

		if startSym == '{' {
			el.keyStart = idx
			el.keyEnd, err = getBlockString(data, idx)
			if err != nil {
				return nil, fmt.Errorf("error get object key, %w", err)
			}
			idx = el.keyEnd + 1
			idx, err = skipColonWithSpaces(data, idx)
			if err != nil {
				return nil, fmt.Errorf("error skip colon after object key, %w", err)
			}
		} else {
			el.arrayIndex = arrayIndex
			arrayIndex++
		}
		el.valueStart = idx
		el.valueEnd, err = getBlock(data, idx)
		if err != nil {
			return nil, fmt.Errorf("error get value, %w", err)
		}
		idx = el.valueEnd + 1

		ok, err = matchElement(data, el, e)
		if err != nil {
			return nil, fmt.Errorf("error match element, %w", err)
		}
		if ok {
			resultElements = append(resultElements, el)
		}

		idx, err = skipSpaces(data, idx)
		if err != nil {
			return nil, fmt.Errorf("error skip spaces after block, %w", err)
		}

		if data[idx] == ',' {
			idx++
			continue
		}
		if data[idx] == endSym {
			break
		}

		return nil, errUnexpectedSymbol(data[idx], ", or "+string(endSym))
	}

	res := make([]byte, 0, len(data))
	res = append(res, startSym)

	first := true
	for _, el := range resultElements {
		if !first {
			res = append(res, ',')
		}
		if startSym == '{' {
			res = append(res, data[el.keyStart:el.keyEnd+1]...)
			res = append(res, ':')
		}
		res = append(res, data[el.valueStart:el.valueEnd+1]...)

		first = false
	}
	res = append(res, endSym)

	return res, nil
}

func matchElement(data []byte, el filterElement, e ast.Expr) (bool, error) {
	res, err := executeExp(data, el, e)
	if err != nil {
		return false, fmt.Errorf("error execute expression, %w", err)
	}
	v, ok := res.(*ast.ExprValue)
	if !ok {
		return false, fmt.Errorf("expression result expect ExprValue, got %T", res)
	}
	if v.T != ast.ExprValueTypeBool {
		return false, fmt.Errorf("expression result expect ExprValue bool, got %v", v.T)
	}
	return v.B, nil
}
