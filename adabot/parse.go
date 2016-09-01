package adabot

import (
	"fmt"
	"strings"
	"text/scanner"
)

// ---- lexer ----

type lexer struct {
	scan  scanner.Scanner
	token rune // current lookahead token
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

type lexPanic string

// describe returns a string describing the current token, for use in errors.
func (lex *lexer) describe() string {
	switch lex.token {
	case scanner.EOF:
		return "end of file"
	case scanner.Ident:
		return fmt.Sprintf("identifier %s", lex.text())
	}
	return fmt.Sprintf("%q", rune(lex.token)) // any other rune
}

// ---- parser ----

type parser struct {
}

// Parse parses the input string as an arithmetic expression.
//
//   expr = num                         a literal number, e.g., 3.14159
//        | id                          a variable name, e.g., x
//        | id '(' expr ',' ... ')'     a function call
//        | '-' expr                    a unary operator (+-)
//        | expr '+' expr               a binary operator (+-*/)
//
func (p *parser) Parse(input string) (_ Expr, err error) {
	defer func() {
		switch x := recover().(type) {
		case nil:
			// no panic
		case lexPanic:
			err = fmt.Errorf("%s", x)
		default:
			// unexpected panic: resume state of panic.
			panic(x)
		}
	}()
	lex := new(lexer)
	lex.scan.Init(strings.NewReader(input))
	lex.scan.Mode = scanner.ScanIdents
	lex.next() // initial lookahead
	e := p.parseExpr(lex)
	if lex.token != scanner.EOF {
		return nil, fmt.Errorf("unexpected %s", lex.describe())
	}
	return e, nil
}

func (p *parser) parseExpr(lex *lexer) Expr {
	return p.parsePrimary(lex)
}

// primary = id
//         | id '(' expr ',' ... ',' expr ')'
//         | num
//         | '(' expr ')'
func (p *parser) parsePrimary(lex *lexer) Expr {
	switch lex.token {
	case scanner.Ident:
		id := lex.text()
		lex.next() // consume Identifier
		if lex.token != '(' {
			return Var(id)
		}
	}
	msg := fmt.Sprintf("unexpected %s", lex.describe())
	panic(lexPanic(msg))
}
