## The emo Programming Language

This repo contains the source code for the interpreter, and will contain the compiler of emo code. 

Currently this project supports: 

- Lexing of Input tokens
- Interpreter Can Parse:
  - Identifiers
  - Let statement
  - Return statement
  - Integer literal
  - Prefix expressions (-, !)
  - Infix expressions (+, -, *, /)
  - Boolean literals 
  - if / else statements
  - function literals
  - call expressions

- Evaluator Can Evaluate: 
  - Integer literals 
  - Boolean literals 
  - Null
  - Prefix expressions (- and !)
  - Infix expressions (+, -, *, /) currently only between two integers 


To Be Done: 
- Evaluator cannot yet evaluate: 
  - Infix expressions containing non integers eg. function calls
  - Conditionals 
  - Return statements 
  - Functions and Function Calls 
- Error handling during evaluation
- Support for more dataTypes:
  - Strings
  - Built-in functions
  - Arrays 
  - Hashes 
  - etc.


## To use the REPL

1. git clone this repo
2. run ```go run main.go```
3. type whatever you want