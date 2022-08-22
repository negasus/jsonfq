package jsonfq

import "fmt"

// findMapElement finds map element with key and returns start and end indexes of map element
// findMapElement expect { at start and string as key
func findMapElement(data []byte, idx int, value string) (int, int, error) {
	var err error
	start := idx
	if idx >= len(data) {
		return start, idx, ErrUnexpectedEndOfData
	}
	if data[idx] != '{' {
		return start, idx, errUnexpectedSymbol(data[idx], "{")
	}
	idx++
	for {
		if idx >= len(data) {
			return start, idx, ErrUnexpectedEndOfData
		}

		if isSpace(data[idx]) {
			idx++
			continue
		}

		var endKeyIdx int
		endKeyIdx, err = getBlockString(data, idx)
		if err != nil {
			return start, idx, fmt.Errorf("expect string, %w", err)
		}

		// Check key match
		// +1 is skip start quote. getBlockString guarantees that string is quoted
		if string(data[idx+1:endKeyIdx]) == value {
			// key equals, skip ":" and returns block start index
			idx, err = skipColonWithSpaces(data, endKeyIdx+1)
			if err != nil {
				return start, idx, fmt.Errorf("error skip colon, %w", err)
			}
			var endIdx int
			endIdx, err = getBlock(data, idx)
			if err != nil {
				return start, endIdx, fmt.Errorf("error get block, %w", err)
			}

			return idx, endIdx, nil
		}

		idx = endKeyIdx + 1

		// key not equals, skip ":" and spaces, get block and skip it
		// if found end of block return error
		// if found a comma, skip it and continue
		idx, err = skipColonWithSpaces(data, idx)
		if err != nil {
			return start, idx, fmt.Errorf("error skip colon, %w", err)
		}

		// skip block value
		idx, err = getBlock(data, idx)
		if err != nil {
			return start, idx, fmt.Errorf("error get block, %w", err)
		}
		idx++ // next sym after block
		idx, err = skipSpaces(data, idx)
		if err != nil {
			return start, idx, fmt.Errorf("error skip spaces, %w", err)
		}

		if data[idx] == ',' {
			idx++
			continue
		}

		if data[idx] == '}' {
			return start, idx, ErrMapFieldNotFound
		}

		return start, idx, errUnexpectedSymbol(data[idx], ", or }")
	}
}
