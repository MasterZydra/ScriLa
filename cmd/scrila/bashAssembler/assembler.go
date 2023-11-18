package bashAssembler

import (
	"ScriLa/cmd/scrila/bashAst"
	"ScriLa/cmd/scrila/config"
	"fmt"
	"os"
)

type Assembler struct {
	outputFilename string
	outputFile     *os.File

	testMode      bool
	testPrintMode bool

	nativeScrilaFuncs map[string]nativeScrilaFunc

	isFuncContext bool
	indentDepth   int
}

func NewAssembler() *Assembler {
	assembler := &Assembler{indentDepth: 0}
	assembler.registerNativeScrilaFuncs()
	return assembler
}

func (self *Assembler) Assemble(astNode bashAst.IStatement) error {
	if config.Filename != "" && !self.testMode {
		self.outputFilename = config.Filename + ".sh"
		f, err := os.Create(self.outputFilename)
		if err != nil {
			fmt.Println("Something went wrong creating the output file:", err)
			os.Exit(1)
		}
		defer f.Close()
		self.outputFile = f
	}

	self.writeFileHeader()

	err := self.assemble(astNode)
	if err != nil {
		return err
	}

	return nil
}

func (self *Assembler) assemble(astNode bashAst.IStatement) error {
	switch astNode.GetKind() {
	case bashAst.AssignmentExprNode:
		return self.evalAssignmentExpr(bashAst.StmtToAssignmentExpr(astNode))
	case bashAst.BashStmtNode:
		return self.evalBashStmt(bashAst.StmtToBashStmt(astNode))
	case bashAst.BreakExprNode:
		return self.evalBreakExpr(astNode)
	case bashAst.CallExprNode:
		return self.evalCallExpr(bashAst.StmtToCallExpr(astNode))
	case bashAst.CommentNode:
		return self.evalComment(bashAst.StmtToComment(astNode))
	case bashAst.ContinueExprNode:
		return self.evalContinueExpr(astNode)
	case bashAst.FuncDeclarationNode:
		return self.evalFuncDeclaration(bashAst.StmtToFuncDeclaration(astNode))
	case bashAst.IfStmtNode:
		return self.evalIfStmt(bashAst.StmtToIfStmt(astNode))
	case bashAst.ProgramNode:
		return self.evalProgram(bashAst.StmtToProgram(astNode))
	case bashAst.ReturnExprNode:
		return self.evalReturnExpr(astNode)
	case bashAst.WhileStmtNode:
		return self.evalWhileStmt(bashAst.StmtToWhileStmt(astNode))
	default:
		return fmt.Errorf("The '%s' AST Node has not been setup for interpretion", astNode.GetKind())
	}
}
