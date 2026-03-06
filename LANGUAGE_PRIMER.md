# CartBASIC Language Primer

This is a brief guide for writing and running programs in the current interpreter.

## Running Code

- Interactive mode (REPL):
- Run the binary with no arguments, then type expressions/statements.
- Prompt: `mb>`
- Exit commands: `exit` or `quit`
- Script mode:
- Run `cart-basic run script.bas`
- Script mode executes the file and reports parse/runtime errors to stderr.

## Core Syntax

### Variables

```basic
let name = "Cart";
let count = 3;
```

### Reassignment

```basic
let total = 1;
total = total + 2;
```

### Expressions

```basic
1 + 2 * 3;
"hello" + " " + "world";
```

### Conditionals

```basic
if (count > 5) {
  "large";
} elseif (count == 5) {
  "equal";
} else {
  "small";
}
```

### Loops

```basic
let i = 0;
while (i < 3) {
  i = i + 1;
}

let sum = 0;
for (let j = 0; j < 5; j = j + 1) {
  sum = sum + j;
}
```

### Functions

```basic
function add(a, b) {
  return a + b;
}

let result = add(2, 4);
```

Anonymous function form:

```basic
let twice = fn(x) { x * 2; };
let v = twice(10);
```

### Arrays and Hashes

```basic
let nums = [1, 2, 3];
let second = nums[1];

let user = {"name": "Cart", "active": true};
let uname = user["name"];
```

## Error Handling

Use `try/catch` to recover from runtime errors:

```basic
try {
  1 / 0;
} catch err {
  err;
}
```

## Concurrency

Use `spawn` to run a function call asynchronously:

```basic
let work = fn(x) { x + 1; };
spawn work(41);
```

## Built-ins You Can Use Today

- `File.Read(path)`
- `File.Write(path, content)`
- `Fetch(url)`
- `Json.Parse(jsonString)`
- `Json.Stringify(value)`

See `STDLIB.md` for full signatures and behavior.

## Notes

- Numbers are currently integer-based.
- In REPL mode, evaluated results are printed.
- In script mode (`run`), execution is non-interactive and focused on side effects and error reporting.
