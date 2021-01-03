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

To Be Done: 
- Evaluation of AST


## To use the REPL

1. git clone this repo
2. run ```go run main.go```
3. type whatever you want