package bashAssembler

import (
	"ScriLa/cmd/scrila/bashAst"
	"fmt"
)

func (self *Assembler) evalBashStmt(bashStmt bashAst.IBashStmt) error {
	self.writeLnWithTabsToFile(bashStmt.GetValue())
	return nil
}

func (self *Assembler) evalComment(comment bashAst.IComment) error {
	self.writeLnWithTabsToFile(fmt.Sprintf("# %s", comment.GetValue()))
	return nil
}

func (self *Assembler) evalFuncDeclaration(funcDecl bashAst.IFuncDeclaration) error {
	// Create a documentation header with the functions signature written in ScriLa syntax
	sig, err := self.getFuncSignature(funcDecl)
	if err != nil {
		return err
	}
	self.writeLnToFile(fmt.Sprintf("# %s", sig))

	self.writeLnToFile(fmt.Sprintf("%s () {", funcDecl.GetName()))
	self.isFuncContext = true
	self.incTabs()

	// Setup parameters
	for i, param := range funcDecl.GetParams() {
		self.writeLnWithTabsToFile(fmt.Sprintf("local %s=$%d", param.GetName(), i+1))
	}

	// Assemble body line by line
	if err = self.assembleBody(funcDecl.GetBody()); err != nil {
		return err
	}

	self.isFuncContext = false
	self.decTabs()
	self.writeLnToFile("}\n")
	return nil
}

func (self *Assembler) evalIfStmt(ifStmt bashAst.IIfStmt) error {
	bash, err := stmtToBashConditionStr(ifStmt.GetCondition())
	if err != nil {
		return err
	}
	self.writeLnWithTabsToFile(fmt.Sprintf("if %s", bash))
	self.writeLnWithTabsToFile("then")
	self.incTabs()

	// Assemble body line by line
	if err = self.assembleBody(ifStmt.GetBody()); err != nil {
		return err
	}
	self.decTabs()

	if ifStmt.GetElse() != nil {
		if err = self.evalIfStmtElse(ifStmt.GetElse()); err != nil {
			return err
		}
	}

	self.writeLnWithTabsToFile("fi")
	return nil
}

func (self *Assembler) evalIfStmtElse(ifStmt bashAst.IIfStmt) error {
	if ifStmt.GetCondition() != nil {
		bash, err := stmtToBashConditionStr(ifStmt.GetCondition())
		if err != nil {
			return err
		}
		self.writeLnWithTabsToFile(fmt.Sprintf("elif %s", bash))
		self.writeLnWithTabsToFile("then")
	} else {
		self.writeLnWithTabsToFile("else")
	}
	self.incTabs()

	// Assemble body line by line
	if err := self.assembleBody(ifStmt.GetBody()); err != nil {
		return err
	}
	self.decTabs()

	if ifStmt.GetElse() != nil {
		if err := self.evalIfStmtElse(ifStmt.GetElse()); err != nil {
			return err
		}
	}

	return nil
}

func (self *Assembler) evalProgram(program bashAst.IProgram) error {
	var err error

	if len(program.GetNativeBody()) > 0 {
		self.writeLnToFile("# Native function implementations\n")
		for _, stmt := range program.GetNativeBody() {
			if err = self.assemble(stmt); err != nil {
				return err
			}
		}
	}

	self.writeLnToFile("# User script\n")
	for _, stmt := range program.GetUserBody() {
		if err = self.assemble(stmt); err != nil {
			return err
		}
	}

	return nil
}

func (self *Assembler) evalWhileStmt(whileStmt bashAst.IWhileStmt) error {
	bash, err := stmtToBashConditionStr(whileStmt.GetCondition())
	if err != nil {
		return err
	}
	self.writeLnWithTabsToFile(fmt.Sprintf("while %s", bash))
	self.writeLnWithTabsToFile("do")
	self.incTabs()

	// Assemble body line by line
	if err = self.assembleBody(whileStmt.GetBody()); err != nil {
		return err
	}
	self.decTabs()

	self.writeLnWithTabsToFile("done")
	return nil
}
