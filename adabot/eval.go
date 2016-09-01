// Package adabot provides an expression evaluator and driver-specific robot code
package adabot

import (
	"fmt"
	"log"
)

//!+env

type controlFunc struct {
	Fn    func(int) error
	Param int
}
type Env map[Var]controlFunc

/*	// Map of function name to func
	funcs := map[string]fn{
		"sin": math.Sin,
	}
	return funcs[c.fn](c.args[0].eval(env))
*/
// An Expr is an expression to command the robot
type Expr interface {
	eval(env Env)
}

// Eval contains the parser to maintain the implicit interface satisfaction
// by the expression types.
type Eval struct {
	parser parser
	bot    *Robot
	env    Env
}

// A Var identifies a command variable
type Var string

func (v Var) eval(env Env) {
}

// NewEval constructs an unexported parser object to store the Env as state.
func NewEval() *Eval {
	bot := NewRobot()
	env := Env{
		// Robot control function map: WASD
		"w":  controlFunc{Fn: bot.Forward, Param: 1},
		"ww": controlFunc{Fn: bot.Forward, Param: 3},
		"a":  controlFunc{Fn: bot.Left, Param: 1},
		"aa": controlFunc{Fn: bot.Left, Param: 3},
		"s":  controlFunc{Fn: bot.Backward, Param: 1},
		"ss": controlFunc{Fn: bot.Backward, Param: 3},
		"d":  controlFunc{Fn: bot.Right, Param: 1},
		"dd": controlFunc{Fn: bot.Right, Param: 3},
	}
	p := parser{}
	e := Eval{env: env, parser: p}
	return &e
}

// Run parses the given string expression then retrieves and executes the
// accepted Robot function if any.
func (e *Eval) Run(input string) {
	// parse
	expr, err := e.parser.Parse(input)
	if err != nil {
		panic(fmt.Sprintf("unsupported expression: %s. [Error: %s]",
			input, err.Error()))
	}
	// type assertion, retrieve the accepted Robot functions from the Env
	controlFunc, ok := e.env[expr.(Var)]
	if ok {
		err = controlFunc.Fn(controlFunc.Param)
		if err != nil {
			log.Printf("%s\n", err.Error())
		}
	}
}
