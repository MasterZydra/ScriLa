package bashAst

import (
	"fmt"
	"strings"
)

func StmtToArray(stmt IStatement) IArray {
	var i interface{} = stmt
	return i.(IArray)
}

func StmtToAssignmentExpr(stmt IStatement) IAssignmentExpr {
	var i interface{} = stmt
	return i.(IAssignmentExpr)
}

func StmtToBashStmt(stmt IStatement) IBashStmt {
	var i interface{} = stmt
	return i.(IBashStmt)
}

func StmtToBinaryOpExpr(stmt IStatement) IBinaryOpExpr {
	var i interface{} = stmt
	return i.(IBinaryOpExpr)
}

func StmtToBoolLiteral(stmt IStatement) IBoolLiteral {
	var i interface{} = stmt
	return i.(IBoolLiteral)
}

func StmtToCallExpr(stmt IStatement) ICallExpr {
	var i interface{} = stmt
	return i.(ICallExpr)
}

func StmtToComment(stmt IStatement) IComment {
	var i interface{} = stmt
	return i.(IComment)
}

func StmtToFuncDeclaration(stmt IStatement) IFuncDeclaration {
	var i interface{} = stmt
	return i.(IFuncDeclaration)
}

func StmtToIfStmt(stmt IStatement) IIfStmt {
	var i interface{} = stmt
	return i.(IIfStmt)
}

func StmtToIntLiteral(stmt IStatement) IIntLiteral {
	var i interface{} = stmt
	return i.(IIntLiteral)
}

func StmtToProgram(stmt IStatement) IProgram {
	var i interface{} = stmt
	return i.(IProgram)
}

func StmtToStrLiteral(stmt IStatement) IStrLiteral {
	var i interface{} = stmt
	return i.(IStrLiteral)
}

func StmtToVarLiteral(stmt IStatement) IVarLiteral {
	var i interface{} = stmt
	return i.(IVarLiteral)
}

func StmtToWhileStmt(stmt IStatement) IWhileStmt {
	var i interface{} = stmt
	return i.(IWhileStmt)
}

var tmpVarNameToBashNodeTypeMapping = map[string]NodeType{
	"tmpBool": BoolLiteralNode,
	"tmpInt":  IntLiteralNode,
	"tmpStr":  StrLiteralNode,
}

func TmpVarNameToBashNodeType(tmpVarName string) (NodeType, error) {
	value, ok := tmpVarNameToBashNodeTypeMapping[tmpVarName]
	if !ok {
		return ProgramNode, fmt.Errorf("TmpVarNamesToBashNodeType(): VarName '%s' is not in mapping", tmpVarName)
	}
	return value, nil
}

var indentDepth int = 0

func indent() string {
	return strings.Repeat("  ", indentDepth+1)
}

func SprintAst(program IProgram) string {
	astString := ""
	if len(program.GetNativeBody()) > 0 {
		astString += "Native body:\n"
		for _, stmt := range program.GetNativeBody() {
			astString += fmt.Sprintf("%s%s\n", indent(), stmt)
		}
	}
	if len(program.GetUserBody()) > 0 {
		astString += "User body:\n"
		for _, stmt := range program.GetUserBody() {
			astString += fmt.Sprintf("%s%s\n", indent(), stmt)
		}
	}
	return astString
}
