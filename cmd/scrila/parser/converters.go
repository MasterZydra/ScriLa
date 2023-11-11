package parser

import (
	"ScriLa/cmd/scrila/lexer"
	"ScriLa/cmd/scrila/scrilaAst"
	"fmt"
)

var lexerTokenTypeToScrilaNodeTypeMapping = map[lexer.TokenType]scrilaAst.NodeType{
	lexer.BoolType: scrilaAst.BoolLiteralNode,
	lexer.IntType:  scrilaAst.IntLiteralNode,
	lexer.StrType:  scrilaAst.StrLiteralNode,
	lexer.VoidType: scrilaAst.VoidNode,
}

func lexerTokenTypeToScrilaNodeType(tokenType lexer.TokenType) (scrilaAst.NodeType, error) {
	value, ok := lexerTokenTypeToScrilaNodeTypeMapping[tokenType]
	if !ok {
		return scrilaAst.ProgramNode, fmt.Errorf("lexerTokenTypeToScrilaNodeTypeMapping(): Type '%s' is not in mapping", tokenType)
	}
	return value, nil
}
