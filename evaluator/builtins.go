package evaluator

import (
	"dux/object"
	"fmt"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return newError("wrong number of arguments. got=%d, want=%d", len(args), 1) }
			
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return newError("wrong number of arguments. got=%d, want=%d", len(args), 1) }

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 { return arg.Elements[0] }
			case *object.String:
				if len(arg.Value) > 0 { return &object.String{Value: string(arg.Value[0])} }
			default:
				return newError("invalid argument %s to 'first', must be ARRAY or STRING.", arg.Type())
			}

			return NIL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return newError("wrong number of arguments. got=%d, want=%d", len(args), 1) }

			switch arg := args[0].(type) {
			case *object.Array:
				arrayLen := len(arg.Elements)
				if arrayLen > 0 { return arg.Elements[arrayLen - 1] }
			case *object.String:
				length := len(arg.Value)
				if length > 0 { return &object.String{Value: string(arg.Value[length - 1])} }
			default:
				return newError("invalid argument %s to 'last', must be ARRAY or STRING", arg.Type())
			}

			return NIL
		},
	},
	"tail": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return newError("wrong number of arguments. got=%d, want=%d", len(args), 1) }

			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					//if length == 1 { return &object.Array{Elements: []object.Object{}} }

					newElements := make([]object.Object, length-1, length-1)
					copy(newElements, arg.Elements[1:length])
					
					return &object.Array{Elements: newElements}
				}
			default:
				return newError("invalid argument %s to 'tail', must be ARRAY", arg.Type())
			}

			return NIL
		},
	},
	"head": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return newError("wrong number of arguments. got=%d, want=%d", len(args), 1) }

			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					newElements := make([]object.Object, length-1, length-1)
					copy(newElements, arg.Elements[0:length-1])

					return &object.Array{Elements: newElements}
				}
			default:
				return newError("invalid argument %s to 'tail', must be ARRAY", arg.Type())
			}
			return NIL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 { return newError("wrong number of arguments. got=%d, want=%d", len(args), 2) }

			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)

				newElements := make([]object.Object, length+1, length+1)
				copy(newElements, arg.Elements)

				newElements[length] = args[1]

				return &object.Array{Elements: newElements}
			default:
				return newError("invalid first argument %s to 'head', must be ARRAY", arg.Type())
			}
		},
	},
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NIL
		},
	},
}
