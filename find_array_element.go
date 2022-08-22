package jsonfq

import "fmt"

// findArrayElement finds element in array and returns start and last indexes of array element
// findArrayElement expect [ at start and int as key
func findArrayElement(data []byte, idx int, value int) (int, int, error) {
	var err error

	start := idx

	if idx >= len(data) {
		return idx, idx, ErrUnexpectedEndOfData
	}
	if data[idx] != '[' {
		return idx, idx, errUnexpectedSymbol(data[idx], "[")
	}
	idx++
	var blockNum int
	var blockEnd int
	for {
		if idx >= len(data) {
			return idx, idx, ErrUnexpectedEndOfData
		}

		if isSpace(data[idx]) {
			idx++
			continue
		}

		blockEnd, err = getBlock(data, idx)
		if err != nil {
			return idx, idx, fmt.Errorf("error get block, %w", err)
		}
		if blockNum == value {
			return idx, blockEnd, nil
		}
		idx = blockEnd + 1

		idx, err = skipSpaces(data, idx)
		if err != nil {
			return idx, idx, fmt.Errorf("error skip spaces, %w", err)
		}

		if data[idx] == ',' {
			blockNum++
			idx++
			continue
		}

		if data[idx] == ']' {
			return start, idx, ErrArrayElementNotFound
		}

		return idx, idx, errUnexpectedSymbol(data[idx], ", or ]")
	}
}
