package evaluator

import (
	"encoding/json"
	"modern-basic/object"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuiltinDispatchAndValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		errorSubstr string
	}{
		{name: "fetch via identifier missing args", input: `Fetch()`, errorSubstr: "Fetch expects 1 argument(s), got=0"},
		{name: "file read missing arg", input: `File.Read()`, errorSubstr: "File.Read expects 1 argument(s), got=0"},
		{name: "file write wrong type", input: `File.Write("x.txt", 1)`, errorSubstr: "File.Write argument 2 must be STRING, got=INTEGER"},
		{name: "json parse wrong type", input: `Json.Parse(1)`, errorSubstr: "Json.Parse argument 1 must be STRING, got=INTEGER"},
		{name: "json stringify unsupported function", input: `Json.Stringify(fn(x) { x; })`, errorSubstr: "Json.Stringify failed: unsupported value type FUNCTION"},
		{name: "json stringify success", input: `Json.Stringify([1, 2, 3])`, expected: "[1,2,3]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEval(t, tt.input)
			if tt.errorSubstr != "" {
				assertErrorContains(t, result, tt.errorSubstr)
				return
			}

			strObj, ok := result.(*object.String)
			if !ok {
				t.Fatalf("expected *object.String, got=%T", result)
			}
			if strObj.Value != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, strObj.Value)
			}
		})
	}
}

func TestFileReadWriteBuiltins(t *testing.T) {
	t.Run("write then read", func(t *testing.T) {
		tmpDir := t.TempDir()
		path := filepath.ToSlash(filepath.Join(tmpDir, "sample.txt"))

		writeInput := `File.Write("` + escapeForStringLiteral(path) + `", "hello file")`
		writeResult := testEval(t, writeInput)
		testNullObject(t, writeResult)

		readInput := `File.Read("` + escapeForStringLiteral(path) + `")`
		readResult := testEval(t, readInput)

		strObj, ok := readResult.(*object.String)
		if !ok {
			t.Fatalf("expected *object.String, got=%T", readResult)
		}
		if strObj.Value != "hello file" {
			t.Fatalf("expected file content hello file, got %q", strObj.Value)
		}
	})

	t.Run("read missing file", func(t *testing.T) {
		missing := filepath.ToSlash(filepath.Join(t.TempDir(), "missing.txt"))
		result := testEval(t, `File.Read("`+escapeForStringLiteral(missing)+`")`)
		assertErrorContains(t, result, "File.Read failed:")
	})

	t.Run("write to missing directory fails", func(t *testing.T) {
		path := filepath.ToSlash(filepath.Join(t.TempDir(), "nested", "out.txt"))
		result := testEval(t, `File.Write("`+escapeForStringLiteral(path)+`", "x")`)
		assertErrorContains(t, result, "File.Write failed:")
	})
}

func TestFetchBuiltin(t *testing.T) {
	successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok-response"))
	}))
	defer successServer.Close()

	statusServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request body"))
	}))
	defer statusServer.Close()

	tests := []struct {
		name        string
		input       string
		expected    string
		errorSubstr string
	}{
		{name: "success", input: `Fetch("` + successServer.URL + `")`, expected: "ok-response"},
		{name: "status failure", input: `Fetch("` + statusServer.URL + `")`, errorSubstr: "Fetch status failed: 400 Bad Request: bad request body"},
		{name: "invalid URL", input: `Fetch("::not-a-url::")`, errorSubstr: "Fetch request failed:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEval(t, tt.input)
			if tt.errorSubstr != "" {
				assertErrorContains(t, result, tt.errorSubstr)
				return
			}

			strObj, ok := result.(*object.String)
			if !ok {
				t.Fatalf("expected *object.String, got=%T", result)
			}
			if strObj.Value != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, strObj.Value)
			}
		})
	}
}

func TestJSONBuiltins(t *testing.T) {
	t.Run("parse object and array", func(t *testing.T) {
		builtin, ok := getBuiltin("Json.Parse")
		if !ok {
			t.Fatalf("Json.Parse builtin not found")
		}

		result := builtin.Fn(&object.String{Value: `{"name":"cart","items":[1,true,null]}`})
		hashObj, ok := result.(*object.Hash)
		if !ok {
			t.Fatalf("expected *object.Hash, got=%T", result)
		}
		if len(hashObj.Pairs) != 2 {
			t.Fatalf("expected 2 pairs, got %d", len(hashObj.Pairs))
		}
	})

	t.Run("parse invalid json", func(t *testing.T) {
		result := testEval(t, `Json.Parse("{bad json}")`)
		assertErrorContains(t, result, "Json.Parse failed:")
	})

	t.Run("parse non-integer number rejected", func(t *testing.T) {
		result := testEval(t, `Json.Parse("[1.5]")`)
		assertErrorContains(t, result, "unsupported non-integer number")
	})

	t.Run("stringify and parse round trip", func(t *testing.T) {
		parseBuiltin, ok := getBuiltin("Json.Parse")
		if !ok {
			t.Fatalf("Json.Parse builtin not found")
		}
		stringifyBuiltin, ok := getBuiltin("Json.Stringify")
		if !ok {
			t.Fatalf("Json.Stringify builtin not found")
		}

		parsed := parseBuiltin.Fn(&object.String{Value: `{"a":1,"b":[2,"x",false,null]}`})
		if errObj, isErr := parsed.(*object.Error); isErr {
			t.Fatalf("unexpected parse error: %s", errObj.Message)
		}

		result := stringifyBuiltin.Fn(parsed)
		jsonStr, ok := result.(*object.String)
		if !ok {
			t.Fatalf("expected *object.String, got=%T", result)
		}

		var actual map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr.Value), &actual); err != nil {
			t.Fatalf("expected valid JSON string, got error: %v", err)
		}

		if actual["a"] != float64(1) {
			t.Fatalf("expected a=1, got=%v", actual["a"])
		}

		items, ok := actual["b"].([]interface{})
		if !ok {
			t.Fatalf("expected b array, got=%T", actual["b"])
		}
		if len(items) != 4 {
			t.Fatalf("expected b length 4, got=%d", len(items))
		}
	})

	t.Run("stringify hash with non-string key fails", func(t *testing.T) {
		result := testEval(t, `Json.Stringify({1: "one"})`)
		assertErrorContains(t, result, "hash key must be STRING")
	})
}

func assertErrorContains(t *testing.T, obj object.Object, expectedSubstr string) {
	t.Helper()

	errObj, ok := obj.(*object.Error)
	if !ok {
		t.Fatalf("expected *object.Error, got=%T (%+v)", obj, obj)
	}

	if !strings.Contains(errObj.Message, expectedSubstr) {
		t.Fatalf("expected error containing %q, got %q", expectedSubstr, errObj.Message)
	}
}

func escapeForStringLiteral(path string) string {
	escaped := strings.ReplaceAll(path, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
	return escaped
}

func TestFileWriteCreatesExpectedBytes(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.ToSlash(filepath.Join(tmpDir, "bytes.txt"))

	result := testEval(t, `File.Write("`+escapeForStringLiteral(path)+`", "abc")`)
	testNullObject(t, result)

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}

	if string(content) != "abc" {
		t.Fatalf("expected file content abc, got %q", string(content))
	}
}
