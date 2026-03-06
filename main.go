package main

import (
	"bufio"
	"fmt"
	"io"
	"modern-basic/evaluator"
	"modern-basic/lexer"
	"modern-basic/object"
	"modern-basic/parser"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, os.Args[0]))
}

type cliMode int

const (
	modeREPL cliMode = iota
	modeRunFile
)

func run(args []string, stdin io.Reader, stdout, stderr io.Writer, binaryName string) int {
	mode, target, err := parseCommand(args)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		printUsage(stderr, binaryName)
		return 1
	}

	switch mode {
	case modeREPL:
		repl(stdin, stdout, stderr)
		return 0
	case modeRunFile:
		return runFile(target, stderr)
	default:
		fmt.Fprintln(stderr, "internal error: unsupported mode")
		return 1
	}
}

func parseCommand(args []string) (cliMode, string, error) {
	if len(args) == 0 {
		return modeREPL, "", nil
	}

	if args[0] != "run" {
		return modeREPL, "", fmt.Errorf("unknown command %q", args[0])
	}

	if len(args) < 2 || strings.TrimSpace(args[1]) == "" {
		return modeREPL, "", fmt.Errorf("missing script path for run command")
	}

	if len(args) > 2 {
		return modeREPL, "", fmt.Errorf("run expects exactly one script path")
	}

	return modeRunFile, args[1], nil
}

func printUsage(out io.Writer, binaryName string) {
	name := filepath.Base(binaryName)
	fmt.Fprintf(out, "Usage: %s run <file.bas>\n", name)
}

func runFile(path string, stderr io.Writer) int {
	contents, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(stderr, "file error: cannot read %q: %v\n", path, err)
		return 1
	}

	l := lexer.New(string(contents))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Fprintf(stderr, "parse failure in %q:\n", path)
		printParserErrors(stderr, p.Errors())
		return 1
	}

	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)
	evaluator.WaitForSpawnTasks()

	if errObj, ok := result.(*object.Error); ok {
		fmt.Fprintf(stderr, "runtime failure in %q: %s\n", path, errObj.Message)
		return 1
	}

	return 0
}

func repl(input io.Reader, output, errOut io.Writer) {
	scanner := bufio.NewScanner(input)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(output, "mb> ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(errOut, "input error: %v\n", err)
			}
			fmt.Fprintln(output)
			fmt.Fprintln(output, "Goodbye")
			return
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.EqualFold(line, "exit") || strings.EqualFold(line, "quit") {
			fmt.Fprintln(output, "Goodbye")
			return
		}

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(output, p.Errors())
			continue
		}

		if program != nil {
			evaluated := evaluator.Eval(program, env)
			if evaluated != nil {
				fmt.Fprintln(output, evaluated.Inspect())
			}
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	fmt.Fprintln(out, "parser errors:")
	for _, msg := range errors {
		fmt.Fprintf(out, "  %s\n", msg)
	}
}
