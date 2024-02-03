package evaluator

import (
	"dux/ast"
	"dux/object"
	"fmt"
)


var (
	NIL    = &object.Nil{}
	TRUE   = &object.Boolean{Value: true}
	FALSE  = &object.Boolean{Value: false}
	ZERO   = &object.Integer{Value: 0}
)
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		if node.Value == 0 { return ZERO }
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)

		if isError(right) { return right }
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) { return left }

		right := Eval(node.Right)
		if isError(right) { return right }

		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)

		if isError(val) { return val }
		return &object.ReturnValue{Value: val}
	}

	return nil
}

func evalProgram(program *ast.Program) object.Object {
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
		result = Eval(statement)

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

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isError(condition) { return condition }

	if truthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NIL
	}
}

func evalBlockStatement(bs *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range bs.Statements {
		result = Eval(statement)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	
	return result
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
