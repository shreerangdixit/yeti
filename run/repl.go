package run

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/shreerangdixit/yeti/ast"
	"github.com/shreerangdixit/yeti/build"
	"github.com/shreerangdixit/yeti/eval"
	"github.com/shreerangdixit/yeti/lex"
)

const Logo = `
 __     ________ _______ _____
 \ \   / /  ____|__   __|_   _|
  \ \_/ /| |__     | |    | |
   \   / |  __|    | |    | |
    | |  | |____   | |   _| |_
    |_|  |______|  |_|  |_____|
`

func REPL() {
	r := newRepl()
	r.Start()
}

type repl struct {
	in     io.Reader
	out    io.Writer
	errout io.Writer
}

func newRepl() *repl {
	return &repl{
		in:     os.Stdin,
		out:    os.Stdout,
		errout: os.Stderr,
	}
}

func (r *repl) Start() {
	fmt.Fprintf(r.out, "%s\n", Logo)
	fmt.Fprintf(r.out, "%s", build.Info)

	scanner := bufio.NewScanner(r.in)
	e := eval.NewEvaluator()
	for {
		fmt.Fprintf(r.out, "yeti >>> ")

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		cmd := scanner.Text()

		if len(cmd) == 0 {
			continue
		}

		root, err := ast.New(lex.New(cmd)).RootNode()
		if err != nil {
			r.printErr(cmd, err)
			continue
		}

		// If the input is a single expression, eval and print the result
		// Otherwise run statements
		exp, ok := isSingleExpression(root)
		if !ok {
			_, err = e.Evaluate(root)
			if err != nil {
				r.printErr(cmd, err)
				continue
			}
		} else {
			val, err := e.Evaluate(exp)
			if err != nil {
				r.printErr(cmd, err)
				continue
			} else if val != eval.NIL {
				fmt.Fprintf(r.out, "%s\n", val)
			}
		}
	}
}

func (r *repl) printErr(cmd string, err error) {
	if formatter, ok := eval.NewErrorFormatter(err, eval.NewInMemoryModule("<repl>", "<repl>", cmd)); ok {
		fmt.Fprintf(r.out, "%s", formatter.Format())
		return
	}
	fmt.Fprintf(r.out, "%s\n", err)
}

func isSingleExpression(node ast.Node) (ast.Node, bool) {
	programNode, ok := node.(ast.ProgramNode)
	if !ok {
		return nil, false
	}

	if len(programNode.Declarations) == 0 || len(programNode.Declarations) > 1 {
		return nil, false
	}

	expStat, ok := programNode.Declarations[0].(ast.ExpStmtNode)
	if !ok {
		return nil, false
	}

	return expStat.Exp, true
}
