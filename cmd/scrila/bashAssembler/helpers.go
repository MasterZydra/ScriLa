package bashAssembler

import (
	"ScriLa/cmd/scrila/bashAst"
	"fmt"
)

func (self *Assembler) assembleBody(stmts []bashAst.IStatement) error {
	isBodyEmpty := true
	for _, stmt := range stmts {
		// The body of if/while/function with just Bash comment still counts as empty
		// and throws an error if it is executed.
		if stmt.GetKind() != bashAst.CommentNode {
			isBodyEmpty = false
		}

		if err := self.assemble(stmt); err != nil {
			return err
		}
	}
	// Empty functions/if blocks/while blocks are not allowed in Bash
	if isBodyEmpty {
		self.writeLnWithTabsToFile(":")
		return nil
	}
	return nil
}

func stmtToBashConditionStr(stmt bashAst.IStatement) (string, error) {
	bash, err := stmtToRhsBashStr(stmt)
	if err != nil {
		return "", err
	}
	switch stmt.GetKind() {
	case bashAst.BinaryCompExprNode, bashAst.BinaryOpExprNode:
		return bash, nil
	case bashAst.BoolLiteralNode, bashAst.VarLiteralNode:
		return strToBashBoolComparison(bash), nil
	default:
		return "", fmt.Errorf("stmtToBashConditionStr(): Kind '%s' is not implemented", stmt.GetKind())
	}
}

// Returns the bash equivalent for the given stmt on the right hand side
func stmtToRhsBashStr(stmt bashAst.IStatement) (string, error) {
	bash, err := stmtToBashStr(stmt)
	if err != nil {
		return "", err
	}

	switch stmt.GetKind() {
	case bashAst.BinaryOpExprNode:
		switch bashAst.StmtToBinaryOpExpr(stmt).GetDataType() {
		case bashAst.StrLiteralNode:
			return strToBashStr(bash), nil
		}
	case bashAst.BoolLiteralNode, bashAst.StrLiteralNode:
		return strToBashStr(bash), nil
	case bashAst.VarLiteralNode:
		switch varType := bashAst.StmtToVarLiteral(stmt).GetDataType(); varType {
		case bashAst.BoolLiteralNode, bashAst.StrLiteralNode:
			return strToBashStr(bash), nil
		}
	}
	return bash, err
}

// Returns the bash equivalent for the given stmt e.g. without wrapping in double quotes
func stmtToBashStr(stmt bashAst.IStatement) (string, error) {
	switch stmt.GetKind() {
	case bashAst.BoolLiteralNode:
		// e.g.: "true"
		return boolToBashStr(bashAst.StmtToBoolLiteral(stmt).GetValue()), nil
	case bashAst.BinaryCompExprNode:
		// e.g.: [[ 1 -gt 2 ]]
		return binCompToBashStr(bashAst.StmtToBinaryOpExpr(stmt))
	case bashAst.BinaryOpExprNode:
		// e.g.: $((1 + 2))
		return binOpToBashStr(bashAst.StmtToBinaryOpExpr(stmt))
	case bashAst.IntLiteralNode:
		// e.g.: 42
		return fmt.Sprintf("%d", bashAst.StmtToIntLiteral(stmt).GetValue()), nil
	case bashAst.StrLiteralNode:
		// e.g.: "hello world"
		return bashAst.StmtToStrLiteral(stmt).GetValue(), nil
	case bashAst.VarLiteralNode:
		switch varType := bashAst.StmtToVarLiteral(stmt).GetDataType(); varType {
		case bashAst.BoolLiteralNode, bashAst.StrLiteralNode:
			// e.g.: "${var}"
			return strToBashVar(bashAst.StmtToVarLiteral(stmt).GetValue()), nil
		case bashAst.IntLiteralNode:
			// e.g.: ${var}
			return strToBashVar(bashAst.StmtToVarLiteral(stmt).GetValue()), nil
		default:
			return "", fmt.Errorf("stmtToBashStr(): Var type '%s' is not implemented", varType)
		}
	default:
		return "", fmt.Errorf("stmtToBashStr(): Kind '%s' is not implemented", stmt.GetKind())
	}
}

func binCompToBashStr(binOp bashAst.IBinaryOpExpr) (string, error) {
	lhs, err := stmtToRhsBashStr(binOp.GetLeft())
	if err != nil {
		return "", err
	}
	rhs, err := stmtToRhsBashStr(binOp.GetRight())
	if err != nil {
		return "", err
	}

	// https://devmanual.gentoo.org/tools-reference/bash/index.html
	switch binOp.GetDataType() {
	case bashAst.BoolLiteralNode, bashAst.StrLiteralNode:
		return fmt.Sprintf("[[ %s %s %s ]]", lhs, binOp.GetOperator(), rhs), nil
	case bashAst.IntLiteralNode:
		opMapping := map[string]string{">": "-gt", "<": "-lt", ">=": "-ge", "<=": "-le", "==": "-eq", "!=": "-ne"}
		return fmt.Sprintf("[[ %s %s %s ]]", lhs, opMapping[binOp.GetOperator()], rhs), nil
	default:
		return "", fmt.Errorf("binCompToBashStr(): Kind '%s' is not implemented", binOp.GetDataType())
	}
}

func binOpToBashStr(binOp bashAst.IBinaryOpExpr) (string, error) {
	lhs, err := stmtToBashStr(binOp.GetLeft())
	if err != nil {
		return "", err
	}
	rhs, err := stmtToBashStr(binOp.GetRight())
	if err != nil {
		return "", err
	}

	switch binOp.GetDataType() {
	case bashAst.BoolLiteralNode:
		if binOp.GetLeft().GetKind() == bashAst.BoolLiteralNode {
			lhs = strToBashBoolComparison(strToBashStr(lhs))
		}
		if binOp.GetRight().GetKind() == bashAst.BoolLiteralNode {
			rhs = strToBashBoolComparison(strToBashStr(rhs))
		}
		return fmt.Sprintf("%s %s %s", lhs, binOp.GetOperator(), rhs), nil
	case bashAst.IntLiteralNode:
		return fmt.Sprintf("$((%s %s %s))", lhs, binOp.GetOperator(), rhs), nil
	case bashAst.StrLiteralNode:
		switch binOp.GetOperator() {
		case "+":
			return fmt.Sprintf("%s%s", lhs, rhs), nil
		default:
			return "", fmt.Errorf("binOpToBashStr(): String operations '%s' is not implemented", binOp.GetOperator())
		}
	default:
		return "", fmt.Errorf("binOpToBashStr(): Kind '%s' is not implemented", binOp.GetDataType())
	}
}

// Return a bash comparision to represent a bool (true|false)
func strToBashBoolComparison(value string) string {
	return fmt.Sprintf("[[ %s == \"true\" ]]", value)
}

// Returns the given string wrapped in double quotes
func strToBashStr(value string) string {
	return fmt.Sprintf("\"%s\"", value)
}

// Returns the string "true" or "false" wrapped in double quotes
func boolToBashStr(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

// Return the variable name as Bash variable
func strToBashVar(value string) string {
	return fmt.Sprintf("${%s}", value)
}

var nodeTypeToVarTypeKeywordMapping = map[bashAst.NodeType]string{
	bashAst.BoolLiteralNode: "bool",
	bashAst.IntLiteralNode:  "int",
	bashAst.StrLiteralNode:  "str",
	bashAst.VoidNode:        "void",
}

// Returns the ScriLa variable type keyword for the given Bash NodeType
func nodeTypeToVarTypeKeyword(varType bashAst.NodeType) (string, error) {
	value, ok := nodeTypeToVarTypeKeywordMapping[varType]
	if !ok {
		return "", fmt.Errorf("nodeTypeToVarTypeKeyword(): Type '%s' is not in mapping", varType)
	}
	return value, nil
}
