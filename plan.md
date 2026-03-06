# Implementation Plan: Modern BASIC Interpreter in Go

This document outlines the step-by-step phases for building the modern BASIC interpreter. It is structured to provide incremental victories, allowing you to test and run basic parts of the language before moving on to complex features.

## Phase 1: Project Setup and Foundation

**Goal:** Initialize the Go workspace and define the core architecture.

**Status:** Completed (March 1, 2026).

- **Task 1.1: Initialize Go Module.** Run `go mod init modern-basic` to set up the project.
- **Task 1.2: Establish Directory Structure.** Create packages for the main components: `lexer`, `parser`, `ast` (Abstract Syntax Tree), `evaluator`, and `object` (the type system).
- **Task 1.3: Build a Basic REPL.** Create a simple Read-Eval-Print Loop in the main package that reads user input from the command line and echoes it back. This will eventually feed into the lexer.

## Phase 2: Lexical Analysis (The Lexer)

**Goal:** Convert raw source code text into a stream of meaningful tokens.

**Status:** Completed (March 1, 2026).

- **Task 2.1: Define Token Types.** Create a `token` package. Define constants for keywords (LET, FUNCTION, IF), operators (+, -, =, ==), and types (IDENTIFIER, INT, STRING).
- **Task 2.2: Build the Scanner.** Implement a Lexer struct that reads a string character by character.
- **Task 2.3: Implement NextToken().** Write the logic to advance through the characters and identify the current token, handling whitespace and illegal characters.
- **Task 2.4: Write Lexer Tests.** Create Go tests providing sample BASIC code and asserting that the correct sequence of tokens is generated.

## Phase 3: Parsing and the Abstract Syntax Tree (AST)

**Goal:** Convert the flat stream of tokens into a hierarchical tree representing the structure of the program.

**Status:** Completed (March 1, 2026).

- **Task 3.1: Define AST Nodes.** In the `ast` package, define Go interfaces and structs for Statements (e.g., LetStatement, ReturnStatement) and Expressions (e.g., Identifier, IntegerLiteral, InfixExpression).
- **Task 3.2: Build the Parser Shell.** Create a Parser struct that consumes the Lexer's output.
- **Task 3.3: Parse Simple Statements.** Implement parsing for basic statements like `Let x = 5`.
- **Task 3.4: Parse Expressions.** Implement a recursive descent parser (like a Pratt parser) to handle operator precedence (e.g., ensuring `2 + 3 * 4` parses correctly).
- **Task 3.5: Write Parser Tests.** Feed tokens into the parser and verify the generated AST structure is correct.

## Phase 4: Evaluation and the Object System

**Goal:** Bring the AST to life by writing the engine that executes the tree.

**Status:** Completed (March 4, 2026).

- **Task 4.1: Define the Object Interface.** In the `object` package, define how the interpreter represents data. Create structs for IntegerObject, StringObject, and BooleanObject.
- **Task 4.2: Build the Eval Function.** Write the core recursive function that takes an AST node, determines its type, and evaluates it to an Object.
- **Task 4.3: Evaluate Expressions.** Implement the logic to perform actual math and string concatenation when the evaluator encounters infix expressions.
- **Task 4.4: Implement the Environment.** Create an Environment struct (essentially a wrapper around a Go map) to store variable names and their associated Objects, enabling state.

## Phase 5: Control Flow and Functions

**Goal:** Add logic, branching, and reusable code blocks.

**Status:** Completed (March 4, 2026).

- **Task 5.1: If/ElseIf/Else.** Update the parser to recognize conditional blocks and the evaluator to selectively execute branches based on boolean evaluation.
- **Task 5.2: Loops.** Implement While and For loops. Ensure the evaluator can repeatedly evaluate a block of statements until a condition is met.
- **Task 5.3: Function Parsing.** Add AST nodes for FunctionDefinition and FunctionCall.
- **Task 5.4: Function Evaluation & Local Scope.** Update the evaluator to handle function calls. Crucially, when a function is called, create a new, enclosed Environment to ensure local variables do not leak into the global scope.

## Phase 6: Modern Data Structures

**Goal:** Implement dynamic lists and dictionaries.

**Status:** Completed (March 4, 2026).

- **Task 6.1: Array/List Implementation.** Add an ArrayObject (backed by a Go slice) to the object system. Update the parser to handle `[1, 2, 3]` syntax and index operators `myList[0]`.
- **Task 6.2: Dictionary Implementation.** Add a HashObject (backed by a Go map). Implement parsing for `{"key": "value"}` and retrieval logic.

## Phase 7: The Standard Library

**Goal:** Provide the built-in tools that make the language useful.

**Status:** Completed (March 4, 2026).

- **Task 7.1: Setup Built-in Function Registry.** Create a way to register native Go functions that can be called from within the BASIC code.
- **Task 7.2: Implement I/O.** Write the `File.Read` and `File.Write` built-ins using Go's `os` and `io` packages.
- **Task 7.3: Implement Networking.** Write the `Fetch` built-in wrapping Go's `net/http` package.
- **Task 7.4: Implement JSON Support.** Write `Json.Parse` and `Json.Stringify` built-ins wrapping Go's `encoding/json` package to interface with Arrays and Dictionaries.

## Phase 8: Advanced Features (Concurrency and Error Handling)

**Goal:** Add the distinguishing modern features detailed in the design document.

**Status:** Completed (March 5, 2026).

- **Task 8.1: Try/Catch.** Implement the Try block in the parser. In the evaluator, modify how errors bubble up. If an error object is returned inside a Try block, bind it to the variable in the Catch block instead of halting execution.
- **Task 8.2: The Spawn Keyword.** Implement the Spawn statement. In the evaluator, when this node is encountered, execute the function call inside a Go goroutine (`go evalFunction(...)`), allowing the main interpreter to continue immediately.

## Phase 9: Finalization

**Goal:** Package the language for actual use.

**Status:** Completed (March 5, 2026).

- **Task 9.1: File Execution.** Modify the main package to read and execute `.bas` files from the command line (e.g., `cart-basic run script.bas`).
- **Task 9.2: Standard Library Documentation.** Document the built-in functions.
- **Task 9.3: Refactoring and Edge Cases.** Run comprehensive tests, check for memory leaks (especially in the Environment scopes), and handle edge cases in syntax errors.
