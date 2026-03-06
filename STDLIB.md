# Standard Library Reference

This document describes the built-in functions currently implemented by CartBASIC.

## File.Read(path)

- Signature: `File.Read(path)`
- Arguments:
- `path`: STRING
- Behavior:
- Reads the full file at `path` into memory.
- Return value:
- STRING containing file contents.
- Error conditions:
- Wrong argument count: `File.Read expects 1 argument(s), got=X`
- Wrong argument type: `File.Read argument 1 must be STRING, got=TYPE`
- Read failure (missing file, directory path, permission issues): `File.Read failed: <os error>`
- Example:

```basic
let contents = File.Read("config.json");
let data = Json.Parse(contents);
```

## File.Write(path, content)

- Signature: `File.Write(path, content)`
- Arguments:
- `path`: STRING
- `content`: STRING
- Behavior:
- Writes `content` bytes to `path` with mode `0644`.
- Does not create parent directories.
- Return value:
- `null`
- Error conditions:
- Wrong argument count: `File.Write expects 2 argument(s), got=X`
- Wrong argument type: `File.Write argument 1/2 must be STRING, got=TYPE`
- Write failure (invalid path, missing parent directory, permission issues): `File.Write failed: <os error>`
- Example:

```basic
let payload = Json.Stringify({"name": "cart", "enabled": true});
File.Write("output.json", payload);
```

## Fetch(url)

- Signature: `Fetch(url)`
- Arguments:
- `url`: STRING
- Behavior:
- Performs an HTTP GET request with a 10-second timeout.
- Reads and returns the entire response body on 2xx status codes.
- Return value:
- STRING containing response body.
- Error conditions:
- Wrong argument count: `Fetch expects 1 argument(s), got=X`
- Wrong argument type: `Fetch argument 1 must be STRING, got=TYPE`
- Request setup/network failure: `Fetch request failed: <net error>`
- Body read failure: `Fetch read failed: <io error>`
- Non-2xx status: `Fetch status failed: <code> <status>[: <body>]`
- Example:

```basic
let body = Fetch("https://example.com/data.json");
let parsed = Json.Parse(body);
```

## Json.Parse(jsonString)

- Signature: `Json.Parse(jsonString)`
- Arguments:
- `jsonString`: STRING
- Behavior:
- Parses one JSON value and converts it to interpreter values.
- JSON objects become hash maps with STRING keys.
- JSON arrays become arrays.
- JSON numbers must be integers.
- Return value:
- Converted value: HASH, ARRAY, INTEGER, STRING, BOOLEAN, or `null`.
- Error conditions:
- Wrong argument count: `Json.Parse expects 1 argument(s), got=X`
- Wrong argument type: `Json.Parse argument 1 must be STRING, got=TYPE`
- Invalid JSON: `Json.Parse failed: <decode error>`
- Trailing extra content: `Json.Parse failed: trailing data after JSON value`
- Non-integer number: `Json.Parse failed: unsupported non-integer number ...`
- Out-of-range integer: `Json.Parse failed: unsupported integer number ...`
- Example:

```basic
let raw = "{\"name\":\"cart\",\"items\":[1,2,3]}";
let obj = Json.Parse(raw);
let name = obj["name"];
```

## Json.Stringify(value)

- Signature: `Json.Stringify(value)`
- Arguments:
- `value`: Any supported runtime value
- Behavior:
- Converts interpreter values to JSON text.
- Supported: `null`, BOOLEAN, INTEGER, STRING, ARRAY, HASH with STRING keys.
- Return value:
- STRING containing compact JSON.
- Error conditions:
- Wrong argument count: `Json.Stringify expects 1 argument(s), got=X`
- Unsupported value type: `Json.Stringify failed: unsupported value type TYPE`
- HASH contains non-STRING key: `Json.Stringify failed: hash key must be STRING, got=TYPE`
- Example:

```basic
let doc = {"project": "cart", "count": 3, "ok": true};
let json = Json.Stringify(doc);
File.Write("doc.json", json);
```
