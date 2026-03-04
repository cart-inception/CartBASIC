package main

import (
	"bufio"
	"fmt"
	"modern-basic/evaluator"
	"modern-basic/lexer"
	"modern-basic/object"
	"modern-basic/parser"
	"os"
	"strings"
)

func main() {
	repl()
}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnvironment()

	for {
		fmt.Print("mb> ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "input error: %v\n", err)
			}
			fmt.Println()
			fmt.Println("Goodbye")
			return
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.EqualFold(line, "exit") || strings.EqualFold(line, "quit") {
			fmt.Println("Goodbye")
			return
		}

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(p.Errors())
			continue
		}

		if program != nil {
			evaluated := evaluator.Eval(program, env)
			if evaluated != nil {
				fmt.Println(evaluated.Inspect())
			}
		}
	}
}

func printParserErrors(errors []string) {
	fmt.Println("parser errors:")
	for _, msg := range errors {
		fmt.Printf("  %s\n", msg)
	}
}
