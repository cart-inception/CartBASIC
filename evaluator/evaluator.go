package evaluator

import (
	"fmt"
	"modern-basic/ast"
	"modern-basic/object"
	"sync"
)

var (
	trueObj  = &object.Boolean{Value: true}
	falseObj = &object.Boolean{Value: false}
	nullObj  = &object.Null{}
	spawnWG  sync.WaitGroup
)

// Eval evaluates an AST node and returns the resulting runtime object.
func Eval(node ast.Node, env *object.Environment) object.Object {
	if env == nil {
		env = object.NewEnvironment()
	}

	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return nullObj

	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return nullObj

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.DotExpression:
		return evalDotExpression(node)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.TryExpression:
		return evalTryExpression(node, env)

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.SpawnStatement:
		return evalSpawnStatement(node, env)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	}

	return nullObj
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = nullObj

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE || rt == object.ERROR {
				return result
			}
		}
	}

	return result
}

func evalIfExpression(expression *ast.IfExpression, env *object.Environment) object.Object {
	for _, branch := range expression.Branches {
		condition := Eval(branch.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			return Eval(branch.Consequence, env)
		}
	}

	if expression.Alternative != nil {
		return Eval(expression.Alternative, env)
	}

	return nullObj
}

func evalTryExpression(expression *ast.TryExpression, env *object.Environment) object.Object {
	result := Eval(expression.TryBlock, env)
	errObj, ok := result.(*object.Error)
	if !ok {
		return result
	}

	catchEnv := object.NewEnclosedEnvironment(env)
	catchEnv.Set(expression.CatchIdent.Value, &object.CaughtError{Err: errObj})

	return Eval(expression.CatchBlock, catchEnv)
}

func evalSpawnStatement(statement *ast.SpawnStatement, env *object.Environment) object.Object {
	function := Eval(statement.Call.Function, env)
	if isError(function) {
		return function
	}

	args := evalExpressions(statement.Call.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	if validationError := validateSpawnCall(function, args); validationError != nil {
		return validationError
	}

	spawnWG.Add(1)
	go func(fn object.Object, callArgs []object.Object) {
		defer spawnWG.Done()
		_ = applyFunction(fn, callArgs)
	}(function, args)

	return nullObj
}

func evalWhileExpression(expression *ast.WhileExpression, env *object.Environment) object.Object {
	result := object.Object(nullObj)

	for {
		condition := Eval(expression.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(expression.Body, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE || rt == object.ERROR {
				return result
			}
		}
	}

	return result
}

func evalForExpression(expression *ast.ForExpression, env *object.Environment) object.Object {
	result := object.Object(nullObj)

	if expression.Init != nil {
		initResult := Eval(expression.Init, env)
		if isError(initResult) {
			return initResult
		}
	}

	for {
		if expression.Condition != nil {
			condition := Eval(expression.Condition, env)
			if isError(condition) {
				return condition
			}

			if !isTruthy(condition) {
				break
			}
		}

		result = Eval(expression.Body, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE || rt == object.ERROR {
				return result
			}
		}

		if expression.Post != nil {
			postResult := Eval(expression.Post, env)
			if isError(postResult) {
				return postResult
			}
			if postResult != nil && postResult.Type() == object.RETURN_VALUE {
				return postResult
			}
		}
	}

	return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	result := []object.Object{}

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	if builtin, ok := fn.(*object.Builtin); ok {
		return builtin.Fn(args...)
	}

	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	if len(args) != len(function.Parameters) {
		return newError("wrong number of arguments: got=%d, want=%d", len(args), len(function.Parameters))
	}

	extendedEnv := object.NewEnclosedEnvironment(function.Env)
	for i, param := range function.Parameters {
		extendedEnv.Set(param.Value, args[i])
	}

	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func validateSpawnCall(fn object.Object, args []object.Object) *object.Error {
	if _, ok := fn.(*object.Builtin); ok {
		return nil
	}

	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	if len(args) != len(function.Parameters) {
		return newError("wrong number of arguments: got=%d, want=%d", len(args), len(function.Parameters))
	}

	return nil
}

// WaitForSpawnTasks blocks until all spawned calls finish.
func WaitForSpawnTasks() {
	spawnWG.Wait()
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object = nullObj

	for _, statement := range statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return nativeBoolToBooleanObject(!isTruthy(right))
	case "-":
		if right.Type() != object.INTEGER {
			return newError("unknown operator: -%s", right.Type())
		}
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left.(*object.Integer), right.(*object.Integer))
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(operator, left.(*object.String), right.(*object.String))
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left *object.Integer, right *object.Integer) object.Object {
	switch operator {
	case "+":
		return &object.Integer{Value: left.Value + right.Value}
	case "-":
		return &object.Integer{Value: left.Value - right.Value}
	case "*":
		return &object.Integer{Value: left.Value * right.Value}
	case "/":
		if right.Value == 0 {
			return newError("division by zero")
		}
		return &object.Integer{Value: left.Value / right.Value}
	case "<":
		return nativeBoolToBooleanObject(left.Value < right.Value)
	case ">":
		return nativeBoolToBooleanObject(left.Value > right.Value)
	case "==":
		return nativeBoolToBooleanObject(left.Value == right.Value)
	case "!=":
		return nativeBoolToBooleanObject(left.Value != right.Value)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left *object.String, right *object.String) object.Object {
	switch operator {
	case "+":
		return &object.String{Value: left.Value + right.Value}
	case "==":
		return nativeBoolToBooleanObject(left.Value == right.Value)
	case "!=":
		return nativeBoolToBooleanObject(left.Value != right.Value)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if obj, ok := env.Get(node.Value); ok {
		return obj
	}

	if builtin, ok := getBuiltin(node.Value); ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
}

func evalDotExpression(node *ast.DotExpression) object.Object {
	name, ok := dotExpressionName(node)
	if !ok {
		return newError("unsupported dotted expression")
	}

	if builtin, ok := getBuiltin(name); ok {
		return builtin
	}

	return newError("identifier not found: %s", name)
}

func dotExpressionName(exp ast.Expression) (string, bool) {
	switch node := exp.(type) {
	case *ast.Identifier:
		return node.Value, true
	case *ast.DotExpression:
		left, ok := dotExpressionName(node.Left)
		if !ok || node.Right == nil {
			return "", false
		}
		return left + "." + node.Right.Value, true
	default:
		return "", false
	}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array object.Object, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return nullObj
	}

	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash object.Object, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return nullObj
	}

	return pair.Value
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return trueObj
	}

	return falseObj
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case nullObj:
		return false
	case trueObj:
		return true
	case falseObj:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR
}
