package lexer

import "fmt"

var singleCharTokens = map[string]TokenType{
	"-": BinaryOperator,
	":": Colon,
	",": Comma,
	";": Semicolon,
	"(": OpenParen,
	")": CloseParen,
	"{": OpenBrace,
	"}": CloseBrace,
	"*": BinaryOperator,
	"/": BinaryOperator,
	"+": BinaryOperator,
	"=": Equals,
}

var keywords = map[string]TokenType{
	"bool":  BoolType,
	"const": Const,
	"int":   IntType,
	"obj":   ObjType,
	"str":   StrType,
}

type TokenType string

const (
	Semicolon  TokenType = "Semicolon"
	Comma      TokenType = "Comma"
	Colon      TokenType = "Colon"
	OpenBrace  TokenType = "OpenBrace"  // {
	CloseBrace TokenType = "CloseBrace" // }
	EndOfFile  TokenType = "EOF"
	// --- Operations ---
	BinaryOperator TokenType = "BinaryOperator"
	// --- Priority ---
	OpenParen  TokenType = "OpenParen"  // (
	CloseParen TokenType = "CloseParen" // )
	// --- Variables ---
	Identifier TokenType = "Identifier"
	Equals     TokenType = "Equals"
	// Variable types
	Bool     TokenType = "BoolValue"
	BoolType TokenType = "BoolType"
	Const    TokenType = "Const"
	Int      TokenType = "IntValue"
	IntType  TokenType = "IntType"
	ObjType  TokenType = "ObjType"
	Str      TokenType = "StrValue"
	StrType  TokenType = "StrType"
)

type Token struct {
	TokenType TokenType
	Value     string
}

func (self *Token) String() string {
	return fmt.Sprintf("&{Token %s %s}", self.TokenType, self.Value)
}
