package main

import (
	"fmt"
	"os"
)

func Tokenize(fileContents []byte) {
    scanner := NewScanner(string(fileContents))
    tokens := scanner.ScanTokens()

    for _, token := range(tokens) {
        fmt.Println(token.ToString())
    }

    if scanner.HasError() {
        os.Exit(65)
    }
}

func Parse(fileContents []byte) {
    scanner := NewScanner(string(fileContents))
    tokens := scanner.ScanTokens()

    parser := NewParser(tokens)
    expr, err := parser.ParseExpression()

    if expr != nil {
        fmt.Println(expr.String())
    }

    if err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err.Error())
        os.Exit(65)
    }
}

func Evaluate(fileContents []byte) {
    scanner := NewScanner(string(fileContents))
    tokens := scanner.ScanTokens()

    parser := NewParser(tokens)
    expr, _ := parser.ParseExpression()

    interpreter := NewInterpreter(expr)
    value, err := interpreter.Eval()

    if err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err.Error())
        os.Exit(70)
    }

    fmt.Println(value)
}

func Run(fileContents []byte) {
    scanner := NewScanner(string(fileContents))
    tokens := scanner.ScanTokens()

    parser := NewParser(tokens)
    stmts, err := parser.Parse()

    if err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err.Error())
        os.Exit(65)
    }

    for _, stmt := range(stmts) {
        if err := stmt.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "%s\n", err.Error())
            os.Exit(70)
        }
    }
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

    switch command {
    case "tokenize":
        Tokenize(fileContents)
        return
    case "parse":
        Parse(fileContents)
        return
    case "evaluate":
        Evaluate(fileContents)
        return
    case "run":
        Run(fileContents)
        return
    }

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
	os.Exit(1)
}
