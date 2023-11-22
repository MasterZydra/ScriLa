package lexer

import "fmt"

// TODO Move the following two arrays in separate package so that they are delcared once.
// Now BooleanOps and ComparisonOps exist here and in ScrilaAst
var BooleanOps = []string{"||", "&&"}

var ComparisonOps = []string{"<", ">", "<=", ">=", "!=", "=="}

var singleCharTokens = map[string]TokenType{
	"-": BinaryOperator,
	":": Colon,
	",": Comma,
	".": Dot,
	";": Semicolon,
	"(": OpenParen,
	")": CloseParen,
	"{": OpenBrace,
	"}": CloseBrace,
	"[": OpenBracket,
	"]": CloseBracket,
	"*": BinaryOperator,
	"/": BinaryOperator,
	"+": BinaryOperator,
	"=": Equals,
	"<": BinaryOperator,
	">": BinaryOperator,
}

var keywords = map[string]TokenType{
	"bool":     BoolType,
	"break":    Break,
	"const":    Const,
	"continue": Continue,
	"else":     Else,
	"func":     Function,
	"if":       If,
	"int":      IntType,
	"obj":      ObjType,
	"return":   Return,
	"str":      StrType,
	"void":     VoidType,
	"while":    While,
	"true":     Bool,
	"false":    Bool,
}

type TokenType string

const (
	If             TokenType = "If"
	Else           TokenType = "Else"
	While          TokenType = "While"
	Break          TokenType = "Break"
	Continue       TokenType = "Continue"
	Function       TokenType = "Function"
	Comment        TokenType = "Comment"
	BinaryOperator TokenType = "BinaryOperator"
	Return         TokenType = "Return"
	// Characters
	Semicolon    TokenType = "Semicolon"
	Comma        TokenType = "Comma"
	Colon        TokenType = "Colon"
	Dot          TokenType = "Dot"
	Equals       TokenType = "Equals"
	OpenBrace    TokenType = "OpenBrace"    // {
	CloseBrace   TokenType = "CloseBrace"   // }
	OpenBracket  TokenType = "OpenBracket"  // [
	CloseBracket TokenType = "CloseBracket" // ]
	OpenParen    TokenType = "OpenParen"    // (
	CloseParen   TokenType = "CloseParen"   // )
	EndOfFile    TokenType = "EOF"
	// Variables
	Identifier TokenType = "Identifier"
	Bool       TokenType = "BoolValue"
	BoolType   TokenType = "BoolType"
	Const      TokenType = "Const"
	Int        TokenType = "IntValue"
	IntType    TokenType = "IntType"
	ObjType    TokenType = "ObjType"
	Str        TokenType = "StrValue"
	StrType    TokenType = "StrType"
	VoidType   TokenType = "VoidType"
)

type Token struct {
	TokenType TokenType
	Value     string
	Ln        int
	Col       int
}

func (self *Token) String() string {
	return fmt.Sprintf("&{Token %s %s %d %d}", self.TokenType, self.Value, self.Ln, self.Col)
}
