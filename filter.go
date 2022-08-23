package jsonfq

import (
	"fmt"

	"github.com/negasus/jsonfq/ast"
)

type filterContext struct {
	keyStart   int
	keyEnd     int
	valueStart int
	valueEnd   int
	arrayIndex int
}

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

	var resultElements []filterContext

	for {
		if idx >= len(data) {
			return nil, ErrUnexpectedEndOfData
		}

		el := filterContext{-1, -1, -1, -1, -1}

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

	var ii int

	data[0] = startSym
	ii++

	first := true
	for _, el := range resultElements {
		if !first {
			data[ii] = ','
			ii++
		}
		if startSym == '{' {
			copy(data[ii:], data[el.keyStart:el.keyEnd+1])
			ii += el.keyEnd - el.keyStart + 1
			data[ii] = ':'
			ii++
		}
		copy(data[ii:], data[el.valueStart:el.valueEnd+1])
		ii += el.valueEnd - el.valueStart + 1
		first = false
	}
	data[ii] = endSym
	return data[:ii+1], nil
}

func matchElement(data []byte, el filterContext, e ast.Expr) (bool, error) {
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
