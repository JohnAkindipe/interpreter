# Monkey Language Interpreter

A tree-walking interpreter written in Go. This project implements a complete interpreter pipeline — from raw source text to a working REPL while employing a test-driven development approach throughout.

---

## Table of Contents

- [Interpretation Technique](#interpretation-technique)
- [Project Structure](#project-structure)
- [Components](#components)
  - [Token](#1-token)
  - [Lexer](#2-lexer)
  - [AST](#3-abstract-syntax-tree-ast)
  - [Parser](#4-parser)
  - [REPL](#5-repl)
- [Parsing: Pratt Parsing](#parsing-pratt-parsing)
- [Test-Driven Development](#test-driven-development)
- [Getting Started](#getting-started)

---

## Interpretation Technique

This interpreter uses the **tree-walking** strategy, also known as an **AST interpreter**. It does not compile source code to bytecode or machine code. Instead, it:

1. **Lexes** the source text into a flat stream of tokens.
2. **Parses** the token stream into an Abstract Syntax Tree (AST).
3. **Walks** the AST nodes directly to evaluate the program (evaluation phase — in progress).

This approach is a simple and approachable way to build an interpreter and is ideal for getting a glimpse of how programming languages work under the hood.

```
Source Code (string)
        │
        ▼
    [ Lexer ]  ──────────────▶  Token Stream
        │
        ▼
    [ Parser ] ──────────────▶  Abstract Syntax Tree (AST)
        │
        ▼
  [ Evaluator ]
        │
        ▼
    Result / Value
```

---

## Project Structure

```
interpreter/
├── main.go          # Entry point — launches the REPL
├── go.mod
├── token/
│   └── token.go     # Token type definitions and keyword lookup
├── lexer/
│   ├── lexer.go     # Lexical analysis (source → tokens)
│   └── lexer_test.go
├── ast/
│   ├── ast.go       # AST node definitions and interfaces
│   └── ast_test.go
├── parser/
│   ├── parser.go    # Pratt parser (tokens → AST)
│   └── parser_test.go
└── repl/
    └── repl.go      # Read-Eval-Print Loop
```

---

## Components

### 1. Token

**Package:** `token`  
**File:** [token/token.go](token/token.go)

The `token` package defines token types, as well as the `Token` struct that pairs a type with its literal string value.

**Token categories:**

| Category | Examples |
|---|---|
| Identifiers & Literals | `IDENT`, `INT` |
| Operators | `+`, `-`, `*`, `/`, `!`, `==`, `!=`, `<`, `>` |
| Delimiters | `;`, `,`, `(`, `)`, `{`, `}` |
| Keywords | `let`, `fn`, `if`, `else`, `true`, `false`, `return` |
| Special | `EOF`, `ILLEGAL` |

The `LookupIdent` function differentiates user-defined identifiers from reserved keywords, ensuring `fn` is classified as `FUNCTION` and `foobar` as `IDENT`.

```go
type Token struct {
    Type    TokenType
    Literal string
}
```

---

### 2. Lexer

**Package:** `lexer`  
**File:** [lexer/lexer.go](lexer/lexer.go)

The lexer (also called a *tokenizer* or *scanner*) performs **lexical analysis**: it reads the raw source string character by character as input and produces tokens as output.

**Key design details:**

- Uses a **two-pointer approach**: `position` (current character) and `readPosition` (next character) — enabling single-character lookahead via `peekChar()` without consuming the character.
- Handles **multi-character operators** like `==` and `!=` by peeking at the next character before deciding the token type.
- **Skips whitespace** (`' '`, `'\t'`, `'\n'`, `'\r'`) between tokens.
- Reads **identifiers** and **integers** by consuming characters while the predicate (`isLetter`, `isDigit`) holds, then classifying via `LookupIdent`.
- Returns an `EOF` token when the input is exhausted and an `ILLEGAL` token for unrecognised characters.

```go
// Two-pointer state inside the Lexer
type Lexer struct {
    input        string
    position     int  // current character index
    readPosition int  // next character index (lookahead)
    ch           byte // character currently under examination
}
```

---

### 3. Abstract Syntax Tree (AST)

**Package:** `ast`  
**File:** [ast/ast.go](ast/ast.go)

The AST represents the structure of the source program as a tree of Go structs. Every node in the tree implements the `Node` interface:

```go
type Node interface {
    TokenLiteral() string  // the literal value of the token
    String() string        // human-readable representation
}
```

Nodes are further divided into two categories:

- **`Statement`** — constructs that perform an action but do not produce a value (e.g. `let x = 5;`, `return 5;`).
- **`Expression`** — constructs that produce a value (e.g. `5`, `x + y`, `!true`, `fn(x) { x; }`).

---

### 4. Parser

**Package:** `parser`  
**File:** [parser/parser.go](parser/parser.go)

The parser consumes the token stream produced by the lexer and constructs the AST. It implements a **Pratt parser** (see below).

The `Parser` struct maintains two tokens at all times — the current token and a lookahead peek token — enabling it to make decisions based on what comes next without backtracking:

```go
type Parser struct {
    l          *lexer.Lexer
    curToken   token.Token
    peekToken  token.Token
    errors     []string
    prefixParseFns map[token.TokenType]prefixParseFn
    infixParseFns  map[token.TokenType]infixParseFn
}
```

**Error handling:** The parser collects all errors into a slice rather than panicking on the first issue. After parsing, callers can retrieve all errors via `p.Errors()`, and test helpers use `checkParserErrors` to surface them immediately.

**Statement dispatch:** `parseStatement` acts as a router — it inspects the current token to decide whether to delegate to `parseLetStatement`, `parseReturnStatement`, or the general `parseExpressionStatement`.

---

### 5. REPL

**Package:** `repl`  
**File:** [repl/repl.go](repl/repl.go)

The REPL (Read-Eval-Print Loop) is the interactive interface to the interpreter. It reads a line of input, runs it through the lexer, and prints every token produced — one per line.

---

## Parsing: Pratt Parsing

The expression parser uses **Pratt parsing** (also called *top-down operator precedence* parsing).

### Core Idea

Instead of encoding grammar rules rigidly into recursive functions, Pratt parsing associates each token type with one or two **parse functions**:

- **`prefixParseFn`** — called when the token appears at the *start* of an expression (e.g. an integer literal, an identifier, or the `-` in `-5`).
- **`infixParseFn`** — called when the token appears *between* two expressions (e.g. `+` in `5 + 10`). It receives the left-hand expression as an argument.

These functions are stored in maps keyed by `token.TokenType`:

```go
prefixParseFns map[token.TokenType]prefixParseFn
infixParseFns  map[token.TokenType]infixParseFn
```

## Test-Driven Development

The project uses a test-driven development approach. Every major component has a corresponding `_test.go` file written **before or alongside** the implementation.

---
