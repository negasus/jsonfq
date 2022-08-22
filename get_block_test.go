package jsonfq

import "testing"

func TestParser_getBlockString(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		start    int
		want     int
		wantErr  bool
		errValue string
	}{
		{
			name:    "string with quote \"",
			data:    `"foo"`,
			start:   0,
			want:    4,
			wantErr: false,
		},
		{
			name:    "string with quote '",
			data:    `'foo'`,
			start:   0,
			want:    4,
			wantErr: false,
		},
		{
			name:    "string with tail",
			data:    `"foo"bar`,
			start:   0,
			want:    4,
			wantErr: false,
		},
		{
			name:     "no quote at start",
			data:     `x"foo"`,
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: `unexpected symbol x, expect "quote"`,
		},
		{
			name:     "no tail quote",
			data:     `"foo`,
			start:    0,
			want:     4,
			wantErr:  true,
			errValue: ErrUnexpectedEndOfData.Error(),
		},
		{
			name:    "escaped quote",
			data:    `"foo\" ' bar"`,
			start:   0,
			want:    12,
			wantErr: false,
		},
		{
			name:    "escaped quote 2",
			data:    `"foo\\\" \n \t ' bar" tail`,
			start:   0,
			want:    20,
			wantErr: false,
		},
		{
			name:    "bad escaped quote",
			data:    `"foo\\" bar"`,
			start:   0,
			want:    6,
			wantErr: false,
		},
		{
			name:     "empty input",
			data:     ``,
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: ErrUnexpectedEndOfData.Error(),
		},
		{
			name:    "empty string",
			data:    `""`,
			start:   0,
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBlockString([]byte(tt.data), tt.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBlockString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getBlockString() got = %v, want %v", got, tt.want)
			}
			if err != nil && err.Error() != tt.errValue {
				t.Errorf("getBlockString() error = %v, want %v", err, tt.errValue)
			}
		})
	}
}

func TestParser_getBlockBool(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		start    int
		want     int
		wantErr  bool
		errValue string
	}{
		{
			name:     "empty data",
			data:     "",
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: ErrUnexpectedEndOfData.Error(),
		},
		{
			name:     "not t/f",
			data:     "a",
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: `unexpected symbol a, expect "true or false"`,
		},
		{
			name:     "tru",
			data:     "tru",
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: `expect true`,
		},
		{
			name:     "truff",
			data:     "truff",
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: `expect true`,
		},
		{
			name:     "fals",
			data:     "fals",
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: `expect false`,
		},
		{
			name:     "falsff",
			data:     "falsff",
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: `expect false`,
		},
		{
			name:    "trueXX",
			data:    "trueXX",
			start:   0,
			want:    3,
			wantErr: false,
		},
		{
			name:    "true",
			data:    "true",
			start:   0,
			want:    3,
			wantErr: false,
		},
		{
			name:    "falseXX",
			data:    "falseXX",
			start:   0,
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBlockBool([]byte(tt.data), tt.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBlockBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getBlockBool() got = %v, want %v", got, tt.want)
			}
			if err != nil && err.Error() != tt.errValue {
				t.Errorf("getBlockBool() error = %v, want %v", err, tt.errValue)
			}
		})
	}
}
func TestParser_getBlockNumber(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		start    int
		want     int
		wantErr  bool
		errValue string
	}{
		{
			name:     "empty block",
			data:     "",
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: ErrUnexpectedEndOfData.Error(),
		},
		{
			name:     "-",
			data:     "-",
			start:    0,
			want:     1,
			wantErr:  true,
			errValue: ErrUnexpectedEndOfData.Error(),
		},
		{
			name:     "dot without base",
			data:     ".100",
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: `unexpected symbol ., expect "digit"`,
		},
		{
			name:     "dot without base, -",
			data:     "-.100",
			start:    0,
			want:     1,
			wantErr:  true,
			errValue: `unexpected symbol ., expect "digit"`,
		},
		{
			name:     "many dot",
			data:     "100..100",
			start:    0,
			want:     4,
			wantErr:  true,
			errValue: `unexpected symbol ., expect "digit"`,
		},
		{
			name:     "many dot, -",
			data:     "-100..100",
			start:    0,
			want:     5,
			wantErr:  true,
			errValue: `unexpected symbol ., expect "digit"`,
		},
		{
			name:     "not digit base",
			data:     "x.100",
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: `unexpected symbol x, expect "digit"`,
		},
		{
			name:     "not digit tail",
			data:     "100.x",
			start:    0,
			want:     4,
			wantErr:  true,
			errValue: `unexpected symbol x, expect "digit"`,
		},
		{
			name:    "0",
			data:    "0",
			start:   0,
			want:    0,
			wantErr: false,
		},
		{
			name:    "-0",
			data:    "-0",
			start:   0,
			want:    1,
			wantErr: false,
		},
		{
			name:    "1234",
			data:    "1234",
			start:   0,
			want:    3,
			wantErr: false,
		},
		{
			name:    "-1234",
			data:    "-1234",
			start:   0,
			want:    4,
			wantErr: false,
		},
		{
			name:    "1234.56",
			data:    "1234.56",
			start:   0,
			want:    6,
			wantErr: false,
		},
		{
			name:    "-1234.56",
			data:    "-1234.56",
			start:   0,
			want:    7,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBlockNumber([]byte(tt.data), tt.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBlockNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getBlockNumber() got = %v, want %v", got, tt.want)
			}
			if err != nil && err.Error() != tt.errValue {
				t.Errorf("getBlockNumber() error = %v, want %v", err, tt.errValue)
			}
		})
	}
}
func TestParser_getBlockNull(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		start    int
		want     int
		wantErr  bool
		errValue string
	}{
		{
			name:     "empty block",
			data:     "",
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: ErrUnexpectedEndOfData.Error(),
		},
		{
			name:     "nulx",
			data:     "nulx",
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: `expect null`,
		},
		{
			name:    "null",
			data:    "null",
			start:   0,
			want:    3,
			wantErr: false,
		},
		{
			name:    "nullXX",
			data:    "nullXX",
			start:   0,
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBlockNull([]byte(tt.data), tt.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBlockNull() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getBlockNull() got = %v, want %v", got, tt.want)
			}
			if err != nil && err.Error() != tt.errValue {
				t.Errorf("getBlockNull() error = %v, want %v", err, tt.errValue)
			}
		})
	}
}

func TestParser_getBlockSym(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		startSym byte
		endSym   byte
		start    int
		want     int
		wantErr  bool
		errValue string
	}{
		{
			name:     "empty block",
			data:     "",
			startSym: 0,
			endSym:   0,
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: ErrUnexpectedEndOfData.Error(),
		},
		{
			name:     "start with bad sym",
			data:     `{"foo":"bar"}`,
			startSym: '[',
			endSym:   ']',
			start:    0,
			want:     0,
			wantErr:  true,
			errValue: `unexpected symbol {, expect "["`,
		},
		{
			name:     "error get string block",
			data:     `{" `,
			startSym: '{',
			endSym:   '}',
			start:    0,
			want:     3,
			wantErr:  true,
			errValue: `error get string block: unexpected end of data`,
		},
		{
			name:     "not closed block",
			data:     `{" ", 10 { [] }`,
			startSym: '{',
			endSym:   '}',
			start:    0,
			want:     15,
			wantErr:  true,
			errValue: ErrUnexpectedEndOfData.Error(),
		},
		{
			name:     "block {}",
			data:     `{" ", 10 { [] } },{}`,
			startSym: '{',
			endSym:   '}',
			start:    0,
			want:     16,
			wantErr:  false,
		},
		{
			name:     "block []",
			data:     `[" ", 10 { [] } ],{}`,
			startSym: '[',
			endSym:   ']',
			start:    0,
			want:     16,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBlockSym([]byte(tt.data), tt.start, tt.startSym, tt.endSym)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBlockSym() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getBlockSym() got = %v, want %v", got, tt.want)
			}
			if err != nil && err.Error() != tt.errValue {
				t.Errorf("getBlockSym() error = %v, want %v", err, tt.errValue)
			}
		})
	}
}

func TestJSON_getBlock(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		idx      int
		want     int
		wantErr  bool
		errValue string
	}{
		{
			name:     "empty block",
			data:     nil,
			idx:      0,
			want:     0,
			wantErr:  true,
			errValue: ErrUnexpectedEndOfData.Error(),
		},
		{
			name:     "bad symbol",
			data:     []byte(`x`),
			idx:      0,
			want:     0,
			wantErr:  true,
			errValue: `unexpected symbol x, expect "object, array, number, string, bool or null"`,
		},
		{
			name:    "string",
			data:    []byte(`"foo"`),
			idx:     0,
			want:    4,
			wantErr: false,
		},
		{
			name:    "object",
			data:    []byte(`{"foo"}`),
			idx:     0,
			want:    6,
			wantErr: false,
		},
		{
			name:    "array",
			data:    []byte(`["foo"]`),
			idx:     0,
			want:    6,
			wantErr: false,
		},
		{
			name:    "true",
			data:    []byte(`true`),
			idx:     0,
			want:    3,
			wantErr: false,
		},
		{
			name:    "false",
			data:    []byte(`false`),
			idx:     0,
			want:    4,
			wantErr: false,
		},
		{
			name:    "null",
			data:    []byte(`null`),
			idx:     0,
			want:    3,
			wantErr: false,
		},
		{
			name:    "number",
			data:    []byte(`-100.23`),
			idx:     0,
			want:    6,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBlock(tt.data, tt.idx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errValue {
				t.Errorf("getBlock() error = %v, want %v", err, tt.errValue)
			}
			if got != tt.want {
				t.Errorf("getBlock() got = %v, want %v", got, tt.want)
			}
		})
	}
}
