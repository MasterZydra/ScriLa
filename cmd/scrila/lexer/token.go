package lexer

import "fmt"

var comparisons = []string{"|", "&", "="}

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
}

var keywords = map[string]TokenType{
	"bool":   BoolType,
	"const":  Const,
	"func":   Function,
	"int":    IntType,
	"obj":    ObjType,
	"return": Return,
	"str":    StrType,
	"void":   VoidType,
}

type TokenType string

const (
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
	OpenBrace    TokenType = "OpenBrace"  // {
	CloseBrace   TokenType = "CloseBrace" // }
	OpenBracket  TokenType = "OpenBrace"  // [
	CloseBracket TokenType = "CloseBrace" // ]
	OpenParen    TokenType = "OpenParen"  // (
	CloseParen   TokenType = "CloseParen" // )
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
