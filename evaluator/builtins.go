package evaluator

import (
	"encoding/json"
	"io"
	"math"
	"modern-basic/object"
	"net/http"
	"os"
	"strings"
	"time"
)

var builtins = map[string]*object.Builtin{
	"Fetch":          {Fn: builtinFetch},
	"File.Read":      {Fn: builtinFileRead},
	"File.Write":     {Fn: builtinFileWrite},
	"Json.Parse":     {Fn: builtinJSONParse},
	"Json.Stringify": {Fn: builtinJSONStringify},
}

func getBuiltin(name string) (*object.Builtin, bool) {
	builtin, ok := builtins[name]
	return builtin, ok
}

func builtinFileRead(args ...object.Object) object.Object {
	if err := expectArgCount("File.Read", args, 1); err != nil {
		return err
	}

	path, errObj := expectStringArg("File.Read", args[0], 1)
	if errObj != nil {
		return errObj
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return newError("File.Read failed: %s", err.Error())
	}

	return &object.String{Value: string(content)}
}

func builtinFileWrite(args ...object.Object) object.Object {
	if err := expectArgCount("File.Write", args, 2); err != nil {
		return err
	}

	path, errObj := expectStringArg("File.Write", args[0], 1)
	if errObj != nil {
		return errObj
	}

	content, errObj := expectStringArg("File.Write", args[1], 2)
	if errObj != nil {
		return errObj
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return newError("File.Write failed: %s", err.Error())
	}

	return nullObj
}

func builtinFetch(args ...object.Object) object.Object {
	if err := expectArgCount("Fetch", args, 1); err != nil {
		return err
	}

	url, errObj := expectStringArg("Fetch", args[0], 1)
	if errObj != nil {
		return errObj
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return newError("Fetch request failed: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return newError("Fetch read failed: %s", err.Error())
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		message := strings.TrimSpace(string(body))
		if message == "" {
			return newError("Fetch status failed: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		}
		return newError("Fetch status failed: %d %s: %s", resp.StatusCode, http.StatusText(resp.StatusCode), message)
	}

	return &object.String{Value: string(body)}
}

func builtinJSONParse(args ...object.Object) object.Object {
	if err := expectArgCount("Json.Parse", args, 1); err != nil {
		return err
	}

	input, errObj := expectStringArg("Json.Parse", args[0], 1)
	if errObj != nil {
		return errObj
	}

	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.UseNumber()

	var payload interface{}
	if err := decoder.Decode(&payload); err != nil {
		return newError("Json.Parse failed: %s", err.Error())
	}

	var trailing interface{}
	if err := decoder.Decode(&trailing); err != io.EOF {
		return newError("Json.Parse failed: trailing data after JSON value")
	}

	value, convErr := jsonValueToObject(payload)
	if convErr != nil {
		return convErr
	}

	return value
}

func builtinJSONStringify(args ...object.Object) object.Object {
	if err := expectArgCount("Json.Stringify", args, 1); err != nil {
		return err
	}

	payload, convErr := objectToJSONValue(args[0])
	if convErr != nil {
		return convErr
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return newError("Json.Stringify failed: %s", err.Error())
	}

	return &object.String{Value: string(bytes)}
}

func expectArgCount(name string, args []object.Object, expected int) *object.Error {
	if len(args) != expected {
		return newError("%s expects %d argument(s), got=%d", name, expected, len(args))
	}
	return nil
}

func expectStringArg(name string, arg object.Object, position int) (string, *object.Error) {
	str, ok := arg.(*object.String)
	if !ok {
		return "", newError("%s argument %d must be STRING, got=%s", name, position, arg.Type())
	}
	return str.Value, nil
}

func jsonValueToObject(value interface{}) (object.Object, *object.Error) {
	switch v := value.(type) {
	case nil:
		return nullObj, nil
	case bool:
		return nativeBoolToBooleanObject(v), nil
	case string:
		return &object.String{Value: v}, nil
	case json.Number:
		if strings.ContainsAny(v.String(), ".eE") {
			return nil, newError("Json.Parse failed: unsupported non-integer number %q", v.String())
		}
		i, err := v.Int64()
		if err != nil {
			return nil, newError("Json.Parse failed: unsupported integer number %q", v.String())
		}
		return &object.Integer{Value: i}, nil
	case float64:
		if math.Trunc(v) != v {
			return nil, newError("Json.Parse failed: unsupported non-integer number %v", v)
		}
		if v > math.MaxInt64 || v < math.MinInt64 {
			return nil, newError("Json.Parse failed: integer out of range %v", v)
		}
		return &object.Integer{Value: int64(v)}, nil
	case []interface{}:
		elements := make([]object.Object, 0, len(v))
		for _, item := range v {
			converted, errObj := jsonValueToObject(item)
			if errObj != nil {
				return nil, errObj
			}
			elements = append(elements, converted)
		}
		return &object.Array{Elements: elements}, nil
	case map[string]interface{}:
		pairs := make(map[object.HashKey]object.HashPair, len(v))
		for key, item := range v {
			converted, errObj := jsonValueToObject(item)
			if errObj != nil {
				return nil, errObj
			}

			keyObj := &object.String{Value: key}
			pairs[keyObj.HashKey()] = object.HashPair{Key: keyObj, Value: converted}
		}
		return &object.Hash{Pairs: pairs}, nil
	default:
		return nil, newError("Json.Parse failed: unsupported JSON type %T", value)
	}
}

func objectToJSONValue(value object.Object) (interface{}, *object.Error) {
	switch v := value.(type) {
	case *object.Null:
		return nil, nil
	case *object.Boolean:
		return v.Value, nil
	case *object.Integer:
		return v.Value, nil
	case *object.String:
		return v.Value, nil
	case *object.Array:
		converted := make([]interface{}, 0, len(v.Elements))
		for _, element := range v.Elements {
			item, errObj := objectToJSONValue(element)
			if errObj != nil {
				return nil, errObj
			}
			converted = append(converted, item)
		}
		return converted, nil
	case *object.Hash:
		converted := make(map[string]interface{}, len(v.Pairs))
		for _, pair := range v.Pairs {
			strKey, ok := pair.Key.(*object.String)
			if !ok {
				return nil, newError("Json.Stringify failed: hash key must be STRING, got=%s", pair.Key.Type())
			}

			item, errObj := objectToJSONValue(pair.Value)
			if errObj != nil {
				return nil, errObj
			}
			converted[strKey.Value] = item
		}
		return converted, nil
	default:
		return nil, newError("Json.Stringify failed: unsupported value type %s", value.Type())
	}
}
