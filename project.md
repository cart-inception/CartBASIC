# Design Document: Modern BASIC Dialect

## 1. Project Vision

The goal of this project is to create a modern, capable, and highly readable programming language inspired by the simplicity of early 2000s BASIC (specifically JustBASIC). While preserving the accessible, beginner-friendly spirit of classic BASIC, this new dialect will introduce modern programming paradigms, first-class data structures, and seamless interaction with today's web-centric ecosystem. The interpreter will be built entirely in Go, leveraging Go's fast compilation, strong typing, and native concurrency.

## 2. Core Goals

- **Zero-Friction Onboarding:** Retain the highly readable, English-like syntax of classic BASIC.
- **Modern Capabilities:** Ensure the language can easily parse JSON, make HTTP requests, and handle modern string manipulation out of the box.
- **Safe Execution:** Eliminate legacy paradigms that lead to spaghetti code or unpredictable state (e.g., global-only variables, GOTO).
- **Leverage Go:** Expose Go's powerful concurrency and standard library through simplified BASIC keywords.

## 3. Syntax & Structural Enhancements

To modernize the language, we will intentionally break backward compatibility with legacy BASIC in a few key areas to enforce better software design.

**Deprecation of GOTO/GOSUB:** Line numbers, line labels, and jump statements will not be supported. Control flow will rely exclusively on structured blocks:
- If / ElseIf / Else / End If
- While / End While
- For / Next

**True Functions and Local Scope:** Global subroutines will be replaced by standard functions that accept parameters and return values. Variables declared within a function will be locally scoped.

```basic
Function CalculateTax(amount, rate)
    Let tax = amount * rate
    Return tax
End Function
```

**Modern Variable Declaration:** We will abandon implicit global variables and type sigils (like `name$` or `age%`). Variables will be declared explicitly with type inference.

```basic
Let name = "Alice"  ' Infers String
Let age = 30        ' Infers Integer/Number
```

## 4. Data Structures & Types

Traditional fixed-size arrays (`Dim arr(10)`) are insufficient for modern data processing. The language will feature dynamic, first-class data structures.

**Dynamic Lists (Arrays):** Arrays that automatically resize, similar to Python lists or Go slices.

```basic
Let names = ["Alice", "Bob", "Charlie"]
List.Append(names, "David")
```

**Dictionaries (Hash Maps):** Key-value pairs for structured data representation.

```basic
Let user = {"name": "Alice", "age": 30}
```

**Native JSON Support:** Because modern programming relies heavily on APIs, converting between Dictionaries/Lists and JSON strings will be built directly into the language.

```basic
Let jsonString = Json.Stringify(user)
Let parsedData = Json.Parse(jsonString)
```

## 5. The "Batteries-Included" Standard Library

The standard library will abstract away complexity, making common tasks single-line operations.

**Network and Web Built-ins:**

```basic
Let response = Fetch("https://api.example.com/data")  ' Executes an HTTP GET request synchronously
```

**File I/O:** Simplified read/write operations replacing the old `Open As #1` syntax.

```basic
Let content = File.Read("config.json")
File.Write("output.txt", "Hello World")
```

## 6. Error Handling Strategy

The classic `On Error Goto` paradigm is unpredictable and difficult to trace. We will adopt a more modern approach.

**Proposed Solution:** A Try / Catch block system, which is familiar and cleanly separates the "happy path" from error management.

```basic
Try
    Let data = File.Read("missing_file.txt")
Catch error
    Print "Could not read file: " + error
End Try
```

## 7. Concurrency (The Go Advantage)

We will expose Go's goroutines to the user through a simplified syntax, making asynchronous tasks trivial even for beginners.

**The Spawn Keyword:** Allows a user to run a function in the background without manually managing threads or standard library async packages.

```basic
Function DownloadLargeFile(url)
    ' Downloading logic here
End Function

' Runs immediately in the background
Spawn DownloadLargeFile("https://example.com/huge.zip")
Print "Download started in background..."
```

## 8. Implementation Architecture (Go)

The Go-based interpreter will follow a standard pipeline:

- **Lexer (Scanner):** Reads the raw `.bas` source text and breaks it down into a stream of Tokens (e.g., TOKEN_LET, TOKEN_IDENTIFIER, TOKEN_STRING).
- **Parser:** Consumes the tokens and builds an Abstract Syntax Tree (AST) based on our grammar rules.
- **Evaluator / Tree-Walker:** Walks the AST and executes the corresponding Go code. Basic types (Strings, Numbers, Booleans, Lists, Dictionaries) will be represented using Go `interface{}` or a custom Object interface to keep the evaluator cleanly typed.
