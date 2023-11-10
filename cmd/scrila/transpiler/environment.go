package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"ScriLa/cmd/scrila/lexer"
	"fmt"

	"golang.org/x/exp/slices"
)

var reservedIdentifiers = []string{"null", "true", "false", "break", "continue"}

func (self *Transpiler) setupScope(env *Environment) {
	// Create Default Global Environment
	env.declareVar("null", NewNullVal(), true, lexer.Identifier)
	env.declareVar("true", NewBoolVal(true), true, lexer.Bool)
	env.declareVar("false", NewBoolVal(false), true, lexer.Bool)
	env.declareVar("break", NewNullVal(), true, lexer.Break)
	env.declareVar("continue", NewNullVal(), true, lexer.Continue)

	// Variables used for internal use
	env.declareVar("tmpStr", NewStrVal(""), false, lexer.StrType)
	env.declareVar("tmpInt", NewIntVal(0), false, lexer.IntType)
	env.declareVar("tmpBool", NewBoolVal(false), false, lexer.BoolType)

	// Define native builtin methods
	self.declareNativeFunctions(env)
}

type Environment struct {
	parent    *Environment
	functions map[string]ast.IRuntimeVal
	variables map[string]ast.IRuntimeVal
	varTypes  map[string]lexer.TokenType
	constants []string
}

func NewEnvironment(parentEnv *Environment, transpiler *Transpiler) *Environment {
	isGlobal := parentEnv == nil
	env := &Environment{
		parent:    parentEnv,
		functions: make(map[string]ast.IRuntimeVal),
		variables: make(map[string]ast.IRuntimeVal),
		varTypes:  make(map[string]lexer.TokenType),
		constants: make([]string, 0),
	}

	if isGlobal {
		transpiler.setupScope(env)
	}
	return env
}

func (self *Environment) declareFunc(funcName string, value ast.IRuntimeVal) (ast.IRuntimeVal, error) {
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

func (self *Environment) lookupFunc(funcName string) (ast.IRuntimeVal, error) {
	env, err := self.resolveFunc(funcName)
	if err != nil {
		return NewNullVal(), err
	}
	return env.functions[funcName], nil
}

func (self *Environment) declareVar(varName string, value ast.IRuntimeVal, isConstant bool, varType lexer.TokenType) (ast.IRuntimeVal, error) {
	if _, ok := self.variables[varName]; ok {
		return NewNullVal(), fmt.Errorf("Cannot declare variable '%s' as it already is defined", varName)
	}

	self.variables[varName] = value
	self.varTypes[varName] = varType

	if isConstant {
		self.constants = append(self.constants, varName)
	}
	return value, nil
}

func (self *Environment) assignVar(varName string, value ast.IRuntimeVal) (ast.IRuntimeVal, error) {
	env, err := self.resolve(varName)
	if err != nil {
		return NewNullVal(), err
	}

	// Cannot assign to constant
	if slices.Contains(self.constants, varName) {
		return NewNullVal(), fmt.Errorf("Cannot reassign to variable '%s' as it was declared constant", varName)
	}

	env.variables[varName] = value
	return value, nil
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

func (self *Environment) lookupVar(varName string) (ast.IRuntimeVal, error) {
	env, err := self.resolve(varName)
	if err != nil {
		return NewNullVal(), err
	}
	return env.variables[varName], nil
}

func (self *Environment) lookupVarType(varName string) (lexer.TokenType, error) {
	env, err := self.resolve(varName)
	if err != nil {
		return lexer.EndOfFile, err
	}
	return env.varTypes[varName], nil
}
