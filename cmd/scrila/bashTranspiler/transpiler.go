package bashTranspiler

import (
	"ScriLa/cmd/scrila/bashAst"
	"ScriLa/cmd/scrila/config"
	"ScriLa/cmd/scrila/scrilaAst"
	"fmt"
	"runtime"
	"strings"

	"golang.org/x/exp/slices"
)

type Context string

const (
	NoContext        Context = "NoContext"
	FunctionContext  Context = "FunctionContext"
	WhileLoopContext Context = "WhileLoopContext"
	IfStmtContext    Context = "IfStmtContext"
)

type Transpiler struct {
	usedNativeFunctions  []string
	userScriptTranspilat string

	contexts        []Context
	bashContexts    []bashAst.IAppendBody
	currentFunc     IFunctionVal
	currentBashFunc bashAst.IFuncDeclaration

	bashStmtStack map[int]bashAst.IStatement

	bashProgram bashAst.IProgram
}

func NewTranspiler() *Transpiler {
	return &Transpiler{
		usedNativeFunctions: []string{},
		contexts:            []Context{NoContext},
		bashContexts:        []bashAst.IAppendBody{},
		bashStmtStack:       make(map[int]bashAst.IStatement),
	}
}

func (self *Transpiler) getPos(astNode scrilaAst.IStatement) string {
	return fmt.Sprintf("%s:%d:%d", config.Filename, astNode.GetLn(), astNode.GetCol())
}

func (self *Transpiler) Transpile(astNode scrilaAst.IStatement, env *Environment) (bashAst.IProgram, error) {
	self.bashProgram = bashAst.NewProgram()

	_, err := self.transpile(astNode, env)
	if err != nil {
		return self.bashProgram, err
	}

	return self.bashProgram, nil
}

func (self *Transpiler) transpile(astNode scrilaAst.IStatement, env *Environment) (scrilaAst.IRuntimeVal, error) {
	switch astNode.GetKind() {
	// Handle Expressions
	case scrilaAst.ArrayLiteralNode:
		return self.evalArray(scrilaAst.ExprToArray(astNode), env)
	case scrilaAst.IntLiteralNode:
		return NewIntVal(scrilaAst.ExprToIntLit(astNode).GetValue()), nil
	case scrilaAst.StrLiteralNode:
		return NewStrVal(scrilaAst.ExprToStrLit(astNode).GetValue()), nil
	case scrilaAst.BoolLiteralNode:
		return NewBoolVal(scrilaAst.ExprToBoolLit(astNode).GetValue()), nil
	case scrilaAst.IdentifierNode:
		return self.evalIdentifier(scrilaAst.ExprToIdent(astNode), env)
	case scrilaAst.ObjectLiteralNode:
		return self.evalObjectExpr(scrilaAst.ExprToObjLit(astNode), env)
	case scrilaAst.CallExprNode:
		return self.evalCallExpr(scrilaAst.ExprToCallExpr(astNode), env)
	case scrilaAst.AssignmentExprNode:
		return self.evalAssignment(scrilaAst.ExprToAssignmentExpr(astNode), env)
	case scrilaAst.BinaryExprNode:
		return self.evalBinaryExpr(scrilaAst.ExprToBinExpr(astNode), env)
	case scrilaAst.MemberExprNode:
		return self.evalMemberExpr(scrilaAst.ExprToMemberExpr(astNode), env)
	case scrilaAst.ReturnExprNode:
		return self.evalReturnExpr(scrilaAst.ExprToReturnExpr(astNode), env)
	case scrilaAst.BreakExprNode, scrilaAst.ContinueExprNode:
		return self.evalWhileExitKeywords(astNode, env)

	// Handle Statements
	case scrilaAst.CommentNode:
		self.appendUserBody(bashAst.NewComment(scrilaAst.ExprToComment(astNode).GetComment()))
		return NewNullVal(), nil
	case scrilaAst.ProgramNode:
		return self.evalProgram(scrilaAst.ExprToProgram(astNode), env)
	case scrilaAst.VarDeclarationNode:
		return self.evalVarDeclaration(scrilaAst.ExprToVarDecl(astNode), env)
	case scrilaAst.IfStatementNode:
		return self.evalIfStatement(scrilaAst.ExprToIfStmt(astNode), env)
	case scrilaAst.WhileStatementNode:
		return self.evalWhileStatement(scrilaAst.ExprToWhileStmt(astNode), env)
	case scrilaAst.FunctionDeclarationNode:
		return self.evalFunctionDeclaration(scrilaAst.ExprToFuncDecl(astNode), env)

	default:
		return NewNullVal(), fmt.Errorf("%s: This AST Node has not been setup for interpretion: %s", self.getPos(astNode), astNode.GetKind())
	}
}

func (self *Transpiler) pushContext(context Context) {
	self.contexts = append(self.contexts, context)
}

func (self *Transpiler) popContext() {
	self.contexts = self.contexts[:len(self.contexts)-1]
}

func (self *Transpiler) pushBashContext(context bashAst.IAppendBody) {
	self.bashContexts = append(self.bashContexts, context)
}

func (self *Transpiler) popBashContext() {
	self.bashContexts = self.bashContexts[:len(self.bashContexts)-1]
}

func (self *Transpiler) currentContext() Context {
	return self.contexts[len(self.contexts)-1]
}

func (self *Transpiler) currentBashContext() bashAst.IAppendBody {
	return self.bashContexts[len(self.bashContexts)-1]
}

func (self *Transpiler) contextContains(context Context) bool {
	return slices.Contains(self.contexts, context)
}

func (self *Transpiler) indent(offset int) string {
	return strings.Repeat("\t", len(self.contexts)-1-offset)
}

func (self *Transpiler) printFuncName(msg string) {
	if config.ShowCallStackScrila {
		pc, _, _, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc).Name()
		funcName = strings.Replace(funcName, "ScriLa/cmd/scrila/transpiler.(*Transpiler).", "", -1)

		if msg == "" {
			fmt.Printf("%s()\n", funcName)
		} else {
			fmt.Printf("%s(): %s\n", funcName, msg)
		}
	}
}

func (self *Transpiler) appendUserBody(stmt bashAst.IStatement) {
	if len(self.bashContexts) != 0 {
		self.currentBashContext().AppendBody(stmt)
	} else if self.currentFunc != nil {
		self.currentBashFunc.AppendBody(stmt)
	} else {
		self.bashProgram.AppendUserBody(stmt)
	}
}
