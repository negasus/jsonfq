%{
package parser

import (
	"strconv"

	"github.com/negasus/jsonfq/ast"
)
%}

%type<nil> program
%type<stmt> stmt
%type<stmt> path
%type<expr> expr
%type<stmtlist> getblockv
%type<stmtlist> getblockk
%type<exprlist> exprlist
%type<exprlist> args

%union {
	token 		*ast.Token

	exprlist []ast.Expr
	stmtlist []ast.Stmt

	stmt 		ast.Stmt
	expr 		ast.Expr
	path 		[]string
}

%token<token> TIdent TString TValue TKey TGte TLte LNe TInt TFloat

%left '>' '<' TGte TLte '=' TNe
%left '+' '-'
%left '*' '/'

%%

program: {}
	|
	program stmt {
		l := yylex.(*scanner)
		l.stmts = append(l.stmts, $2)
	}

stmt:
	path {
		$$ = $1
	}
	|
	'{' expr '}'  {
		$$ = &ast.StmtFilter{Expr: $2, Position: $2.GetPosition()}
	}
	|
	'[' TInt ']'  {
           	ii, ee := strconv.Atoi($2.Value)
           	if ee != nil {
           		yylex.Error(__yyfmt__.Sprintf("%s is not a valid integer", $2.Value))
           	} else {
			$$ = &ast.StmtArrayIndex{Value: ii, Position: $2.Position}
		}
	}

exprlist:
        expr {
           	$$ = []ast.Expr{$1}
        }
        |
        exprlist ',' expr {
           	$$ = append($1, $3)
        }

getblockv:
	TValue {
		$$ = []ast.Stmt{}
	}
	|
	getblockv '.' path {
		$$ = append($1, $3)
	}

getblockk:
	TKey {
		$$ = []ast.Stmt{}
	}
	|
	getblockk '.' path {
		$$ = append($1, $3)
	}


expr:
	TIdent args  {
		$$ = &ast.ExprFn{Name: $1.Value, Args: $2, Position: $1.Position}
	}
	|
	getblockv {
		$$ = &ast.ExprGetBlock{T: ast.ExprGetBlockTypeValue, Stmts: $1}
	}
	|
	getblockk {
		$$ = &ast.ExprGetBlock{T: ast.ExprGetBlockTypeKey, Stmts: $1}
	}
	|
        TString {
		$$ = &ast.ExprValue{T:ast.ExprValueTypeString, S: $1.Value, Position: $1.Position}
        }
        |
        TInt {
           	ii, ee := strconv.Atoi($1.Value)
           	if ee != nil {
           		yylex.Error(__yyfmt__.Sprintf("%s is not a valid integer", $1.Value))
           	} else {
           		$$ = &ast.ExprValue{T:ast.ExprValueTypeInt, I: ii, Position: $1.Position}
           	}
        }
        |
        TFloat {
           	ii, ee := strconv.ParseFloat($1.Value, 64)
           	if ee != nil {
           		yylex.Error(__yyfmt__.Sprintf("%s is not a valid float", $1.Value))
           	} else {
           		$$ = &ast.ExprValue{T:ast.ExprValueTypeFloat, F: ii, Position: $1.Position}
           	}
        }
        |
        expr '>' expr {
           	$$ = &ast.ExprBinaryOp{Left: $1, Op: ">", Right: $3, Position: $1.GetPosition()}
        }
        |
        expr TGte expr {
           	$$ = &ast.ExprBinaryOp{Left: $1, Op: ">=", Right: $3, Position: $1.GetPosition()}
        }
        |
        expr '<' expr {
           	$$ = &ast.ExprBinaryOp{Left: $1, Op: "<", Right: $3, Position: $1.GetPosition()}
        }
        |
        expr TLte expr {
           	$$ = &ast.ExprBinaryOp{Left: $1, Op: "<=", Right: $3, Position: $1.GetPosition()}
        }
        |
        expr '=' expr {
           	$$ = &ast.ExprBinaryOp{Left: $1, Op: "=", Right: $3, Position: $1.GetPosition()}
        }
        |
        expr TNe expr {
           	$$ = &ast.ExprBinaryOp{Left: $1, Op: "!=", Right: $3, Position: $1.GetPosition()}
        }
        |
        expr '+' expr {
           	$$ = &ast.ExprBinaryOp{Left: $1, Op: "+", Right: $3, Position: $1.GetPosition()}
        }
        |
        expr '*' expr {
           	$$ = &ast.ExprBinaryOp{Left: $1, Op: "*", Right: $3, Position: $1.GetPosition()}
        }

path:
	TIdent {
		$$ = &ast.StmtMapKey{Value: $1.Value, Position: $1.Position}
	}
	|
	TString {
		$$ = &ast.StmtMapKey{Value: $1.Value, Position: $1.Position}
	}
	|
	'.' TIdent {
		$$ = &ast.StmtMapKey{Value: $2.Value, Position: $2.Position}
	}
	|
	'.' TString {
		$$ = &ast.StmtMapKey{Value: $2.Value, Position: $2.Position}
	}

args:
        '(' ')' {
           	$$ = []ast.Expr{}
        }
        |
        '(' exprlist ')' {
           	$$ = $2
        }


%%