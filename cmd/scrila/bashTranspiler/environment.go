package bashTranspiler

import (
	"ScriLa/cmd/scrila/scrilaAst"
	"fmt"

	"golang.org/x/exp/slices"
)

func (self *Transpiler) setupScope(env *Environment) {
	// Create Default Global Environment
	env.declareVar("true", true, scrilaAst.BoolLiteralNode)
	env.declareVar("false", true, scrilaAst.BoolLiteralNode)
	env.declareVar("break", true, scrilaAst.BreakExprNode)
	env.declareVar("continue", true, scrilaAst.ContinueExprNode)

	// Variables used for internal use
	env.declareVar("tmpStrs", false, scrilaAst.StrArrayNode)
	env.declareVar("tmpInts", false, scrilaAst.IntArrayNode)
	env.declareVar("tmpBools", false, scrilaAst.BoolArrayNode)

	// Define native builtin methods
	self.declareNativeFunctions(env)
}

type Environment struct {
	parent    *Environment
	functions map[string]scrilaAst.IRuntimeVal
	variables map[string]scrilaAst.NodeType
	constants []string
}

func NewEnvironment(parentEnv *Environment, transpiler *Transpiler) *Environment {
	isGlobal := parentEnv == nil
	env := &Environment{
		parent:    parentEnv,
		functions: make(map[string]scrilaAst.IRuntimeVal),
		variables: make(map[string]scrilaAst.NodeType),
		constants: make([]string, 0),
	}

	if isGlobal {
		transpiler.setupScope(env)
	}
	return env
}

func (self *Environment) declareFunc(funcName string, value scrilaAst.IRuntimeVal) (scrilaAst.IRuntimeVal, error) {
	if self.isFuncDeclared(funcName) {
		return NewNullVal(), fmt.Errorf("Cannot declare function '%s' as it already is defined", funcName)
	}

	self.functions[funcName] = value

	return value, nil
}

func (self *Environment) isFuncDeclared(funcName string) bool {
	if _, ok := self.functions[funcName]; ok {
		return true
	}

	if self.parent == nil {
		return false
	}

	return self.parent.isFuncDeclared(funcName)
}

func (self *Environment) resolveFunc(funcName string) (*Environment, error) {
	if _, ok := self.functions[funcName]; ok {
		return self, nil
	}

	if self.parent == nil {
		return nil, fmt.Errorf("Cannot resolve function '%s' as it does not exist", funcName)
	}

	return self.parent.resolveFunc(funcName)
}

func (self *Environment) lookupFunc(funcName string) (scrilaAst.IRuntimeVal, error) {
	env, err := self.resolveFunc(funcName)
	if err != nil {
		return NewNullVal(), err
	}
	return env.functions[funcName], nil
}

func (self *Environment) declareVar(varName string, isConstant bool, varType scrilaAst.NodeType) (scrilaAst.IRuntimeVal, error) {
	if _, ok := self.variables[varName]; ok {
		return NewNullVal(), fmt.Errorf("Cannot declare variable '%s' as it already is defined", varName)
	}

	self.variables[varName] = varType

	if isConstant {
		self.constants = append(self.constants, varName)
	}

	return scrilaNodeTypeToRuntimeVal(varType)
}

func (self *Environment) assignVar(varName string) (scrilaAst.IRuntimeVal, error) {
	_, err := self.resolve(varName)
	if err != nil {
		return NewNullVal(), err
	}

	// Cannot assign to constant
	if slices.Contains(self.constants, varName) {
		return NewNullVal(), fmt.Errorf("Cannot reassign to variable '%s' as it was declared constant", varName)
	}

	varType, _ := self.lookupVarType(varName)
	return scrilaNodeTypeToRuntimeVal(varType)
}

func (self *Environment) resolve(varName string) (*Environment, error) {
	if _, ok := self.variables[varName]; ok {
		return self, nil
	}

	if self.parent == nil {
		return nil, fmt.Errorf("Cannot resolve variable '%s' as it does not exist", varName)
	}

	return self.parent.resolve(varName)
}

func (self *Environment) lookupVar(varName string) (scrilaAst.IRuntimeVal, error) {
	_, err := self.resolve(varName)
	if err != nil {
		return NewNullVal(), err
	}
	varType, _ := self.lookupVarType(varName)
	return scrilaNodeTypeToRuntimeVal(varType)
}

func (self *Environment) lookupVarType(varName string) (scrilaAst.NodeType, error) {
	env, err := self.resolve(varName)
	if err != nil {
		return "", err
	}
	return env.variables[varName], nil
}
