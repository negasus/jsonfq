package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"

	"github.com/negasus/jsonfq/ast"
)

var eof rune = -1

type scanner struct {
	buf *bytes.Buffer
	ch  rune
	idx int

	stmts []ast.Stmt

	stmtlist []ast.Stmt
	exprList []ast.Expr
	path     []string

	errors []string
}

func newScanner(data []byte) *scanner {
	s := &scanner{
		buf: bytes.NewBuffer(data),
		idx: 0,
	}
	return s
}

func (s *scanner) prev() {
	s.idx--
	err := s.buf.UnreadRune()
	if err != nil {
		s.Error(fmt.Sprintf("error unread rune, %s", err))
	}
}

func (s *scanner) next() {
	ch, _, err := s.buf.ReadRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			s.ch = eof
			return
		}
		s.Error(fmt.Sprintf("error read rune, %s", err))
		return
	}
	s.ch = ch
	s.idx++
}

func (s *scanner) Error(v string) {
	s.errors = append(s.errors, fmt.Sprintf("error %q at position %d", v, s.idx))
}

func (s *scanner) Lex(lval *yySymType) int {

	for {
		s.next()
		if !isSpace(s.ch) {
			break
		}
	}

	t := &ast.Token{
		Position: s.idx - 1,
	}

	switch ch := s.ch; {
	case ch == eof:
		return 0
	case isDigit(ch):
		t.Type, t.Value = s.scanNumber(false)
	case isLetter(ch):
		t.Type = TIdent
		t.Value = s.scanString(0)
	default:
		switch ch {
		case '"':
			t.Type = TString
			t.Value = s.scanString('"')
		case '$':
			s.next()
			v := s.scanString(0)
			switch v {
			case "value":
				t.Type = TValue
			case "key":
				t.Type = TKey
			default:
				return 0
			}
		case ',':
			t.Type = ','
		case '.':
			t.Type = '.'
		case '{':
			t.Type = '{'
		case '}':
			t.Type = '}'
		case '[':
			t.Type = '['
		case ']':
			t.Type = ']'
		case '(':
			t.Type = '('
		case ')':
			t.Type = ')'
		case '=':
			t.Type = '='
		case '+':
			t.Type = '+'
		case '*':
			t.Type = '*'
		case '/':
			t.Type = '/'
		case '-':
			s.next()
			if isDigit(s.ch) {
				s.prev()
				t.Type, t.Value = s.scanNumber(true)
			} else {
				s.prev()
				t.Type = '-'
			}
		case '>':
			t.Type = s.switch2('>', TGte)
		case '<':
			t.Type = s.switch2('<', TLte)
		case '!':
			t.Type = s.switch2('!', TNe)
		}
	}

	lval.token = t

	return lval.token.Type
}

func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}
func lower(ch rune) rune   { return ('a' - 'A') | ch } // returns lower-case ch if ch is ASCII letter
func isDigit(ch rune) bool { return '0' <= ch && ch <= '9' }
func isSpace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

func (s *scanner) scanString(sym rune) (value string) {
	var escaped bool
	if sym > 0 {
		s.next()
	}
	for {
		if sym > 0 && s.ch == '\\' {
			escaped = !escaped
			s.next()
			continue
		}
		if sym > 0 && s.ch == sym && !escaped {
			return
		}
		escaped = false

		if sym == 0 && !isLetter(s.ch) && !isDigit(s.ch) {
			if value == "" {
				s.Error("expect letter")
				return
			}
			if s.ch != eof {
				s.prev()
			}
			return
		}
		value += string(s.ch)
		s.next()
	}
}

func (s *scanner) scanNumber(negative bool) (int, string) {
	var value string
	if s.ch == '-' {
		negative = true
		s.next()
	}
	if negative {
		value += "-"
	}

	var hasDot bool
	for {
		if s.ch == '.' {
			if hasDot {
				s.Error("unexpected dot")
				return 0, ""
			}
			hasDot = true
			value += "."
			s.next()
			continue
		}
		if isDigit(s.ch) {
			value += string(s.ch)
			s.next()
			continue
		}
		if s.ch != eof {
			s.prev()
		}
		if value == "-" || value == "" {
			s.Error("expect digit")
			return 0, ""
		}
		if hasDot {
			return TFloat, value
		}
		return TInt, value
	}
}

func (s *scanner) switch2(t int, t2 int) int {
	s.next()
	if s.ch == '=' {
		return t2
	}
	s.prev()
	return t
}
