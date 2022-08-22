package jsonfq

import (
	"bytes"
	"fmt"
)

// getBlock returns end index of JSON block from start index
// JSON Block it is one of: string, number, bool, object, array or null
func getBlock(data []byte, idx int) (int, error) {
	if idx >= len(data) {
		return idx, ErrUnexpectedEndOfData
	}
	switch data[idx] {
	case '"', '\'':
		return getBlockString(data, idx)
	case '{':
		return getBlockSym(data, idx, '{', '}')
	case '[':
		return getBlockSym(data, idx, '[', ']')
	case 't', 'f':
		return getBlockBool(data, idx)
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return getBlockNumber(data, idx)
	case 'n':
		return getBlockNull(data, idx)
	}

	return idx, errUnexpectedSymbol(data[idx], "object, array, number, string, bool or null")
}

// getBlockString expect string start with quote
// and returns last index of quote or error if end quote not found
func getBlockString(data []byte, idx int) (int, error) {
	if idx >= len(data) {
		return idx, ErrUnexpectedEndOfData
	}

	var quoteSym byte

	if data[idx] != '"' && data[idx] != '\'' {
		return idx, errUnexpectedSymbol(data[idx], "quote")
	}

	if data[idx] == '"' {
		quoteSym = '"'
	} else {
		quoteSym = '\''
	}

	idx++
	var escape bool
	for {
		if idx >= len(data) {
			return idx, ErrUnexpectedEndOfData
		}
		if data[idx] == quoteSym && !escape {
			return idx, nil
		}
		if data[idx] == '\\' {
			escape = !escape
			idx++
			continue
		}
		escape = false
		idx++
	}
}

func getBlockBool(data []byte, idx int) (int, error) {
	if idx >= len(data) {
		return idx, ErrUnexpectedEndOfData
	}
	if data[idx] == 't' {
		if len(data) >= idx+4 && bytes.Equal(data[idx:idx+4], []byte("true")) {
			return idx + 3, nil
		}
		return idx, fmt.Errorf("expect true")
	}
	if data[idx] == 'f' {
		if len(data) >= idx+5 && bytes.Equal(data[idx:idx+5], []byte("false")) {
			return idx + 4, nil
		}
		return idx, fmt.Errorf("expect false")
	}
	return idx, errUnexpectedSymbol(data[idx], "true or false")
}

// getBlockNumber find number and returns last index of number or error if end of number not found
func getBlockNumber(data []byte, idx int) (int, error) {
	if idx >= len(data) {
		return idx, ErrUnexpectedEndOfData
	}

	if data[idx] == '-' {
		idx++
	}
	var baseCount int
	var tailCount int
	var hasDot bool
	for {
		if idx >= len(data) {
			if baseCount > 0 {
				return idx - 1, nil
			}
			return idx, ErrUnexpectedEndOfData
		}
		if data[idx] == '.' {
			// number must be started with digit, not dot
			if baseCount == 0 {
				return idx, errUnexpectedSymbol(data[idx], "digit")
			}
			if hasDot {
				return idx, errUnexpectedSymbol(data[idx], "digit")
			}
			hasDot = true
			idx++
			continue
		}
		if data[idx] < '0' || data[idx] > '9' {
			if baseCount == 0 || (hasDot && tailCount == 0) {
				return idx, errUnexpectedSymbol(data[idx], "digit")
			}
			return idx - 1, nil
		}
		if hasDot {
			tailCount++
		} else {
			baseCount++
		}
		idx++
	}
}

func getBlockNull(data []byte, idx int) (int, error) {
	if len(data) < idx+4 {
		return idx, ErrUnexpectedEndOfData
	}
	if string(data[idx:idx+4]) != "null" {
		return idx, fmt.Errorf("expect null")
	}
	return idx + 3, nil
}

func getBlockSym(data []byte, idx int, startSym, endSym byte) (int, error) {
	if idx >= len(data) {
		return idx, ErrUnexpectedEndOfData
	}

	if data[idx] != startSym {
		return idx, errUnexpectedSymbol(data[idx], string(startSym))
	}

	var err error
	var depth int
	for {
		if idx >= len(data) {
			return idx, ErrUnexpectedEndOfData
		}
		if data[idx] == '\'' || data[idx] == '"' {
			idx, err = getBlockString(data, idx)
			if err != nil {
				return idx, fmt.Errorf("error get string block: %v", err)
			}
			idx++
			continue
		}

		if isSpace(data[idx]) {
			idx++
			continue
		}
		if data[idx] == startSym {
			depth++
			idx++
			continue
		}

		if data[idx] == endSym {
			depth--
			if depth == 0 {
				return idx, nil
			}
			idx++
			continue
		}
		idx++
	}
}
