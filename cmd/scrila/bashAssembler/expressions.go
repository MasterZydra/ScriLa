package bashAssembler

import (
	"ScriLa/cmd/scrila/bashAst"
	"fmt"
)

func (self *Assembler) evalAssignmentExpr(assignment bashAst.IAssignmentExpr) error {
	// e.g.: solution=42
	bash, err := stmtToRhsBashStr(assignment.GetValue())
	if err != nil {
		return err
	}

	format := "%s=%s"
	if assignment.IsDeclaration() && self.isFuncContext {
		format = "local " + format
	}
	self.writeLnWithTabsToFile(fmt.Sprintf(format, assignment.GetVarname().GetValue(), bash))

	return nil
}

func (self *Assembler) evalCallExpr(callExpr bashAst.ICallExpr) error {
	// Check for native ScriLa functions that do not exist as Batch function
	if fn, ok := self.nativeScrilaFuncs[callExpr.GetFuncName()]; ok {
		return fn(callExpr.GetArgs())
	}

	// Call the function
	// e.g.: concat "hello" "world"
	self.writeWithTabsToFile(callExpr.GetFuncName())
	for _, arg := range callExpr.GetArgs() {
		bash, err := stmtToRhsBashStr(arg)
		if err != nil {
			return err
		}
		self.writeToFile(fmt.Sprintf(" %s", bash))
	}
	self.writeLnToFile("")
	return nil
}

func (self *Assembler) evalReturnExpr(returnExpr bashAst.IStatement) error {
	self.writeLnWithTabsToFile("return")
	return nil
}
