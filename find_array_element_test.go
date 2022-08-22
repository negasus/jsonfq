package jsonfq

import "testing"

func TestJSON_findArrayElement(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		argIdx    int
		argValue  int
		wantStart int
		wantEnd   int
		wantErr   bool
		errValue  string
	}{
		{
			name:      "empty data",
			data:      nil,
			argIdx:    0,
			argValue:  0,
			wantStart: 0,
			wantEnd:   0,
			wantErr:   true,
			errValue:  ErrUnexpectedEndOfData.Error(),
		},
		{
			name:      "start not [",
			data:      []byte(`x[]`),
			argIdx:    0,
			argValue:  0,
			wantStart: 0,
			wantEnd:   0,
			wantErr:   true,
			errValue:  `unexpected symbol x, expect "["`,
		},
		{
			name:      "[  ",
			data:      []byte(`[  `),
			argIdx:    0,
			argValue:  0,
			wantStart: 3,
			wantEnd:   3,
			wantErr:   true,
			errValue:  ErrUnexpectedEndOfData.Error(),
		},
		{
			name:      "wrong block",
			data:      []byte(`[x]`),
			argIdx:    0,
			argValue:  0,
			wantStart: 1,
			wantEnd:   1,
			wantErr:   true,
			errValue:  `error get block, unexpected symbol x, expect "object, array, number, string, bool or null"`,
		},
		{
			name:      "error skip spaces after element",
			data:      []byte(`[true `),
			argIdx:    0,
			argValue:  1,
			wantStart: 6,
			wantEnd:   6,
			wantErr:   true,
			errValue:  `error skip spaces, unexpected end of data`,
		},
		{
			name:      "bad closed array",
			data:      []byte(`[true,10[`),
			argIdx:    0,
			argValue:  5,
			wantStart: 8,
			wantEnd:   8,
			wantErr:   true,
			errValue:  `unexpected symbol [, expect ", or ]"`,
		},
		{
			name:      "element not found",
			data:      []byte(`[true,10]`),
			argIdx:    0,
			argValue:  5,
			wantStart: 0,
			wantEnd:   8,
			wantErr:   true,
			errValue:  ErrArrayElementNotFound.Error(),
		},
		{
			name:      "element found bool",
			data:      []byte(`[true,10]`),
			argIdx:    0,
			argValue:  0,
			wantStart: 1,
			wantEnd:   4,
			wantErr:   false,
		},
		{
			name:      "element found number",
			data:      []byte(`[true,10]`),
			argIdx:    0,
			argValue:  1,
			wantStart: 6,
			wantEnd:   7,
			wantErr:   false,
		},
		{
			name:      "element found string",
			data:      []byte(`[true,10,"foo ' bar] [ ] "]`),
			argIdx:    0,
			argValue:  2,
			wantStart: 9,
			wantEnd:   25,
			wantErr:   false,
		},
		{
			name:      "element found object",
			data:      []byte(`[true,10,{"foo":[],"bar":23},20`),
			argIdx:    0,
			argValue:  2,
			wantStart: 9,
			wantEnd:   27,
			wantErr:   false,
		},
		{
			name:      "element found null",
			data:      []byte(`[true,10,null,20`),
			argIdx:    0,
			argValue:  2,
			wantStart: 9,
			wantEnd:   12,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, err := findArrayElement(tt.data, tt.argIdx, tt.argValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("findArrayElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errValue {
				t.Errorf("findArrayElement() error = %v, wantErr %v", err, tt.errValue)
				return
			}
			if gotStart != tt.wantStart {
				t.Errorf("findArrayElement() gotStart = %v, wantStart %v", gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("findArrayElement() gotEnd = %v, wantEnd %v", gotEnd, tt.wantEnd)
			}
		})
	}
}
