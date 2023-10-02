package runtime

import (
	"ScriLa/cmd/scrila/lexer"
	"fmt"

	"golang.org/x/exp/slices"
)

func setupScope(env *Environment) {
	// Create Default Global Environment
	env.declareVar("null", NewNullVal(), true, lexer.Identifier)
	env.declareVar("true", NewBoolVal(true), true, lexer.Bool)
	env.declareVar("false", NewBoolVal(false), true, lexer.Bool)

	// Define native builtin methods
	env.declareFunc("input", NewNativeFunc(nativeInput))
	env.declareFunc("print", NewNativeFunc(nativePrint))
	env.declareFunc("printLn", NewNativeFunc(nativePrintLn))
}

type Environment struct {
	parent    *Environment
	functions map[string]IRuntimeVal
	variables map[string]IRuntimeVal
	varTypes  map[string]lexer.TokenType
	constants []string
}

func NewEnvironment(parentEnv *Environment) *Environment {
	isGlobal := parentEnv == nil
	env := &Environment{
		parent:    parentEnv,
		functions: make(map[string]IRuntimeVal),
		variables: make(map[string]IRuntimeVal),
		varTypes:  make(map[string]lexer.TokenType),
		constants: make([]string, 0),
	}

	if isGlobal {
		setupScope(env)
	}
	return env
}

func (self *Environment) declareFunc(funcName string, value IRuntimeVal) (IRuntimeVal, error) {
	if self.isFuncDeclared(funcName) {
		return NewNullVal(), fmt.Errorf("Cannot declare function '%s'. As it already is defined.", funcName)
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
		return nil, fmt.Errorf("Cannot resolve function '%s' as it does not exist.", funcName)
	}

	return self.parent.resolveFunc(funcName)
}

func (self *Environment) lookupFunc(funcName string) (IRuntimeVal, error) {
	env, err := self.resolveFunc(funcName)
	if err != nil {
		return NewNullVal(), err
	}
	return env.functions[funcName], nil
}

func (self *Environment) declareVar(varName string, value IRuntimeVal, isConstant bool, varType lexer.TokenType) (IRuntimeVal, error) {
	if _, ok := self.variables[varName]; ok {
		return NewNullVal(), fmt.Errorf("Cannot declare variable '%s'. As it already is defined.", varName)
	}

	self.variables[varName] = value
	self.varTypes[varName] = varType

	if isConstant {
		self.constants = append(self.constants, varName)
	}
	return value, nil
}

func (self *Environment) assignVar(varName string, value IRuntimeVal) (IRuntimeVal, error) {
	env, err := self.resolve(varName)
	if err != nil {
		return NewNullVal(), err
	}

	// Cannot assign to constant
	if slices.Contains(self.constants, varName) {
		return NewNullVal(), fmt.Errorf("Cannot reasign to variable '%s' as it was declared constant.", varName)
	}

	env.variables[varName] = value
	return value, nil
}

func (self *Environment) resolve(varName string) (*Environment, error) {
	if _, ok := self.variables[varName]; ok {
		return self, nil
	}

	if self.parent == nil {
		return nil, fmt.Errorf("Cannot resolve variable '%s' as it does not exist.", varName)
	}

	return self.parent.resolve(varName)
}

func (self *Environment) lookupVar(varName string) (IRuntimeVal, error) {
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
