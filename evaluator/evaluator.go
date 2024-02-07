package evaluator

import (
	"dux/ast"
	"dux/object"
	"fmt"
	"strings"
)


var (
	NIL    = &object.Nil{}
	TRUE   = &object.Boolean{Value: true}
	FALSE  = &object.Boolean{Value: false}
	ZERO   = &object.Integer{Value: 0}
)
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		if node.Value == 0 { return ZERO }
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)

		if isError(right) { return right }
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) { return left }

		right := Eval(node.Right, env)
		if isError(right) { return right }

		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)

		if isError(val) { return val }
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) { return val }
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) { return function }

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	/**
		The for loop will only return a evaluated result object.Object when it
		find out a *object.ReturnValue object.Object

		If it find out a *object.ReturnValue, it will stop the for loop execution
		and it's going to return the unwrapped *object.ReturnValue.Value value.

		If it doesn't find one, the loop will repeat running through all the 
		statements in stmts []ast.Statement until it encounter the first 
		successfully evaluated *Object.ReturnValue ocurrency

		If it doesn't any, it's going to return the last evaluated Object
	**/
	for _, statement := range program.Statements {
		result = Eval(statement, env)

		// Unwrap result if it's a return value and return its literal value
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
		return evalExclamationOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalExclamationOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NIL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	/**
		Need to be higher in the switch branches, because can't assure correct 
		pointer comparasion for Integers.
	**/
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case (left.Type() == object.STRING_OBJ || right.Type() == object.STRING_OBJ) && (left.Type() == object.INTEGER_OBJ || right.Type() == object.INTEGER_OBJ):
		return evalStringInfixExpression(operator, left, right)
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

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {

	switch operator {
	case "+":
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value

		return &object.String{Value: leftVal + rightVal}
	case "*":
		var out strings.Builder
		var index int64
		var edge int64
		var str string
		
		integer, ok := left.(*object.Integer)
		if ok {
			edge = integer.Value
			str = right.(*object.String).Value
		} else {
			edge = right.(*object.Integer).Value
			str = left.(*object.String).Value
		}

		for index < edge {
			out.WriteString(str)
			index++
		}

		return &object.String{Value: out.String()}
	}
	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if right == ZERO { return newError("division by zero: it is impossible to divide by zero") }
		return &object.Integer{Value: leftVal / rightVal}
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) { return condition }

	if truthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NIL
	}
}

func evalBlockStatement(bs *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range bs.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	
	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	if !ok { return newError("identifier not found: " + node.Value)}

	return val
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

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
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramId, param := range fn.Parameters {
		env.Set(param.Value, args[paramId])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func truthy(obj object.Object) bool {
	switch obj {
	case NIL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	case ZERO:
		return false
	default:
		return true
	}
}

func nativeBoolToBooleanObject(bol bool) *object.Boolean {
	if bol { return TRUE } else { return FALSE }
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil { return obj.Type() == object.ERROR_OBJ }

	return false
}
