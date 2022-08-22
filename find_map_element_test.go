package jsonfq

import "testing"

func TestJSON_findMapElement(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		argIdx    int
		argValue  string
		wantStart int
		wantEnd   int
		wantErr   bool
		errValue  string
	}{
		{
			name:      "empty data",
			data:      nil,
			argIdx:    0,
			argValue:  "",
			wantStart: 0,
			wantEnd:   0,
			wantErr:   true,
			errValue:  ErrUnexpectedEndOfData.Error(),
		},
		{
			name:      "start not {",
			data:      []byte(`x`),
			argIdx:    0,
			argValue:  "",
			wantStart: 0,
			wantEnd:   0,
			wantErr:   true,
			errValue:  `unexpected symbol x, expect "{"`,
		},
		{
			name:      "empty object",
			data:      []byte(`{`),
			argIdx:    0,
			argValue:  "",
			wantStart: 0,
			wantEnd:   1,
			wantErr:   true,
			errValue:  ErrUnexpectedEndOfData.Error(),
		},
		{
			name:      "error get key",
			data:      []byte(`{ 10`),
			argIdx:    0,
			argValue:  "",
			wantStart: 0,
			wantEnd:   2,
			wantErr:   true,
			errValue:  `expect string, unexpected symbol 1, expect "quote"`,
		},
		{
			name:      "error skip colon with spaces after key",
			data:      []byte(`{ "foo" `),
			argIdx:    0,
			argValue:  "baz",
			wantStart: 0,
			wantEnd:   8,
			wantErr:   true,
			errValue:  `error skip colon, unexpected end of data`,
		},
		{
			name:      "error get block after key",
			data:      []byte(`{ "foo" : x`),
			argIdx:    0,
			argValue:  "baz",
			wantStart: 0,
			wantEnd:   10,
			wantErr:   true,
			errValue:  `error get block, unexpected symbol x, expect "object, array, number, string, bool or null"`,
		},
		{
			name:      "error skip spaces after block",
			data:      []byte(`{ "foo" : 10  `),
			argIdx:    0,
			argValue:  "baz",
			wantStart: 0,
			wantEnd:   14,
			wantErr:   true,
			errValue:  `error skip spaces, unexpected end of data`,
		},
		{
			name:      "bad close block",
			data:      []byte(`{ "foo" : 10{`),
			argIdx:    0,
			argValue:  "baz",
			wantStart: 0,
			wantEnd:   12,
			wantErr:   true,
			errValue:  `unexpected symbol {, expect ", or }"`,
		},
		{
			name:      "not found",
			data:      []byte(`{"foo":10}`),
			argIdx:    0,
			argValue:  "baz",
			wantStart: 0,
			wantEnd:   9,
			wantErr:   true,
			errValue:  ErrMapFieldNotFound.Error(),
		},
		{
			name:      "found number",
			data:      []byte(`{"foo":10,"baz":20}`),
			argIdx:    0,
			argValue:  "baz",
			wantStart: 16,
			wantEnd:   17,
			wantErr:   false,
		},
		{
			name:      "found bool",
			data:      []byte(`{"foo":10,"baz":false,"bar":20}`),
			argIdx:    0,
			argValue:  "baz",
			wantStart: 16,
			wantEnd:   20,
			wantErr:   false,
		},
		{
			name:      "found string",
			data:      []byte(`{"foo":10,"baz":"foo { } \" } ]"}`),
			argIdx:    0,
			argValue:  "baz",
			wantStart: 16,
			wantEnd:   31,
			wantErr:   false,
		},
		{
			name:      "found object",
			data:      []byte(`{"foo":10,"baz":{}}`),
			argIdx:    0,
			argValue:  "baz",
			wantStart: 16,
			wantEnd:   17,
			wantErr:   false,
		},
		{
			name:      "found array (not closed object)",
			data:      []byte(`{"foo":10,"baz":[10]`),
			argIdx:    0,
			argValue:  "baz",
			wantStart: 16,
			wantEnd:   19,
			wantErr:   false,
		},
		{
			name:      "found, error skip colon",
			data:      []byte(`{"foo":10,"baz" `),
			argIdx:    0,
			argValue:  "baz",
			wantStart: 0,
			wantEnd:   16,
			wantErr:   true,
			errValue:  `error skip colon, unexpected end of data`,
		},
		{
			name:      "found, error get block",
			data:      []byte(`{"foo":10,"baz" : x`),
			argIdx:    0,
			argValue:  "baz",
			wantStart: 0,
			wantEnd:   18,
			wantErr:   true,
			errValue:  `error get block, unexpected symbol x, expect "object, array, number, string, bool or null"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, err := findMapElement(tt.data, tt.argIdx, tt.argValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("findMapElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errValue {
				t.Errorf("findMapElement() error = %v, wantErr %v", err, tt.errValue)
				return
			}
			if gotStart != tt.wantStart {
				t.Errorf("findMapElement() gotStart = %v, wantStart %v", gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("findMapElement() gotEnd = %v, wantEnd %v", gotEnd, tt.wantEnd)
			}
		})
	}
}
