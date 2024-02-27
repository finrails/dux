Dux is a programming language interpreter, it interpreters dx programming language. Dux is not a production-ready programming language interpreter, due its lack on standard library functions and mainly because of its parsing design choice.

The Dux interpreter was implemented using the AST(a.k.a. Abstract Syntax Tree) walk parsing strategy, it means that Dux will parse each dx programming language character into tokens, then it will build an AST, parsing each token into a node(i.e. a record), feeding the AST with it. Finally, after the AST is built, the Dux interpreter will evaluates each AST node, returning the result at the end of evaluation(or if it finds a return statement).

Each of these parts(lexer, parser and evaluator) that compose Dux programming language interpreter were written in different packages, and you can skim and check each of those singly in this repository. There're a good amount of tests, so if you had asked to me, i would recommend you to read the tests.

Here below, you will find a brief definition of what each Dux programming language packages roles, you will find the definition of 'dx' programming language as well.

0. token: it's responsible for defining each dx token.

1. lexer: it's responsible for parsing each dx source code character into tokens, for that forward Dux interpreter parse these tokens into AST nodes. It's also useful for keep tracking of current parsing character and number line.

2. ast: it's responsible for defining all AST nodes data structures.

3. parser: it's job is done by parsing each of the dx tokens into AST records and attaching it to the AST record. 

4. object: it's responsible for defining all object representation of the AST records, an evaluated version of an ast record, easier to interact with, holding the final value of the expression.

5. evaluator: it's job is to evaluates an ast node, and transform it into an object, returning it to the main program.

6. repl: an read, eval, print and loop for the dx programming language, it uses dux as it core.

### Features:

- [x] computation (i.e. number operations like: +, -, * and /)
- [x] logic operators (>, <, >=, <=, ==)
- [x] integers, booleans, strings, arrays and hashes
- [x] let statements
- [x] if-else statements
- [x] function statements
- [x] explicit and implicit return statements
- [x] first-class functions
- [x] function calls
- [x] closures
- [x] recursion
- [x] allow variable names to have '?'
- [x] length, first, last, tail, head, push and puts builtin functions
- [ ] floats
- [ ] loop statements
- [ ] else if statement
- [ ] switch statement
- [ ] variable assign

### dx programming language definition:

* variable definition: let variable_name = expression
* function definition: let function_name = fn(parameterx, parametery, ...) { expression block }
* function call: function_name(argumentx, argumenty, ...)
* if-else definition: if (expression) { expression block } else { expression block }

### dx code example

![dx code example](https://i.imgur.com/FKJ8led.png)

### Interpreting dx source code with Dux

You'll have two ways to run dux code, either you can use builtin REPL inputting 'dux' in the shell or 'dux file.dx'.
