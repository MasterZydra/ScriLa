package transpiler

import (
	"ScriLa/cmd/scrila/parser"
	"fmt"
	"strings"
	"testing"
)

func setTestPrintMode() {
	testPrintMode = true
}

func transpileTest(code string) error {
	testMode = true
	fileName = "test.scri"
	parser := parser.NewParser()
	env := NewEnvironment(nil)

	program, err := parser.ProduceAST(code, fileName)
	if err != nil {
		return err
	}
	return Transpile(program, env, "")
}

func TestErrorLexerUnrecognizedChar(t *testing.T) {
	err := transpileTest(`~`)
	expected := fmt.Errorf("test.scri:1:1: Unrecognized character '~' found")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExmaplePrint() {
	setTestPrintMode()
	transpileTest(`
		print("Hello ");
		printLn("World");
		printLn("!");
	`)

	// Output:
	// #!/bin/bash
	// echo -n "Hello "
	// echo "World"
	// echo "!"
}

func ExamplePrintBaseTypes() {
	setTestPrintMode()
	transpileTest(`printLn(42, "str", true, false, null);`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// echo "42 str true false null"
}

func ExamplePrintVariables() {
	setTestPrintMode()
	transpileTest(`
		int i = 42;
		str s = "hello world";
		bool b = false;
		printLn(i, s, b);
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// i=42
	// s="hello world"
	// b="false"
	// echo "${i} ${s} ${b}"
}

func ExampleIntDeclaration() {
	setTestPrintMode()
	transpileTest(`int i = 42;`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// i=42
}

func ExampleIntAssignment() {
	setTestPrintMode()
	transpileTest(`
		int i = 42;
		i = 101;
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// i=42
	// i=101
}

func TestErrorAssignWrongLeftSide(t *testing.T) {
	err := transpileTest(`12 = 34;`)
	expected := fmt.Errorf("test.scri:1:1: Left side of an assignment must be a variable. Got 'IntLiteral'")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorUnsupportedStringOperation(t *testing.T) {
	err := transpileTest(`str s = "str" - "str";`)
	expected := fmt.Errorf("test.scri:1:15: Binary string expression with unsupported operator '-'")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorBinaryExprWithUnsupportedCombination(t *testing.T) {
	err := transpileTest(`int i = "str" - 123;`)
	expected := fmt.Errorf("test.scri:1:15: No support for binary expressions of type 'str' and 'int'")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorIntDeclarationWithMissingSemicolon(t *testing.T) {
	err := transpileTest(`int i = 42`)
	expected := fmt.Errorf("test.scri:1:11: Expression must end with a semicolon")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorDoubleVariableDeclration(t *testing.T) {
	err := transpileTest(`
		int i = 42;
		int i = 42;
	`)
	expected := fmt.Errorf("test.scri:3:7: Cannot declare variable 'i' as it already is defined")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorAssignContVar(t *testing.T) {
	err := transpileTest(`
		const int i = 42;
		i = 43;
	`)
	expected := fmt.Errorf("test.scri:3:3: Cannot reassign to variable 'i' as it was declared constant")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorIntAssignmentWithMissingDeclaration(t *testing.T) {
	err := transpileTest(`i = 42;`)
	expected := fmt.Errorf("test.scri:1:1: Cannot resolve variable 'i' as it does not exist")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleIntAssignmentBinaryExpr() {
	setTestPrintMode()
	transpileTest(`int i = 42 * 2;`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// i=$((42 * 2))
}

func ExampleIntAssignmentBinaryExprWithVar() {
	setTestPrintMode()
	transpileTest(`
		int i = 42;
		int j = i * 2;
		j = (i + 2) * i;
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// i=42
	// j=$((${i} * 2))
	// j=$(($((${i} + 2)) * ${i}))
}

func ExampleStrAssignmentBinaryExprWithVar() {
	setTestPrintMode()
	transpileTest(`
		str a = "Hello";
		str b = "World";
		str c = a + " " + b;
		str d = a + " World";
		d = a + " World";
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// a="Hello"
	// b="World"
	// c="${a} ${b}"
	// d="${a} World"
	// d="${a} World"
}

func ExampleVarDeclarationAndAssignmentWithVariable() {
	setTestPrintMode()
	transpileTest(`
		int i = 123;
		int j = i;
		j = i;

		str s = "str";
		str t = s;
		t = s;
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// i=123
	// j=${i}
	// j=${i}
	// s="str"
	// t="${s}"
	// t="${s}"
}

func TestErrorAssignDifferentVarTypes(t *testing.T) {
	err := transpileTest(`
		int i = 123;
		str s = "str";
		s = i;
	`)
	expected := fmt.Errorf("test.scri:4:7: Cannot assign a value of type 'IntType' to a var of type 'StrType'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorUnsupportedVarType(t *testing.T) {
	err := transpileTest(`const func i = 13;`)
	expected := fmt.Errorf("test.scri:1:7: Variable type 'func' not given or supported")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorDeclareDifferentVarTypes(t *testing.T) {
	err := transpileTest(`
		int i = 123;
		str s = i;
	`)
	expected := fmt.Errorf("test.scri:3:11: Cannot assign a value of type 'IntType' to a var of type 'StrType'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorDeclareDifferentType(t *testing.T) {
	err := transpileTest(`int i = "123";`)
	expected := fmt.Errorf("test.scri:1:11: Cannot assign a value of type 'StrType' to a var of type 'IntType'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorAssignDifferentType(t *testing.T) {
	err := transpileTest(`
		int i = 123;
		i = "456";
	`)
	expected := fmt.Errorf("test.scri:3:9: Cannot assign a value of type 'StrType' to a var of type 'IntType'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleComment() {
	setTestPrintMode()
	transpileTest(`
		# Comment 1
		int i = 42;
		#  Comment 2
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// # Comment 1
	// i=42
	// # Comment 2
}

func TestErrorReturnOutsideOfFunction(t *testing.T) {
	err := transpileTest(`return true;`)
	expected := fmt.Errorf("test.scri:1:1: Return is only allowed inside a function")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorInvalidFuncCallName(t *testing.T) {
	err := transpileTest(`12();`)
	expected := fmt.Errorf("test.scri:1:1: Function name must be an identifier. Got: 'IntLiteral'")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorFuncParamsWithUnexpectedToken(t *testing.T) {
	err := transpileTest(`func fn(int a const) void {}`)
	expected := fmt.Errorf("test.scri:1:15: Unexpected token 'const' in parameter list")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorMissingFuncParamType(t *testing.T) {
	err := transpileTest(`func fn(a) void {}`)
	expected := fmt.Errorf("test.scri:1:9: Expected param type but got Identifier 'a'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorNonexistentFunc(t *testing.T) {
	err := transpileTest(`fn();`)
	expected := fmt.Errorf("test.scri:1:1: Cannot resolve function 'fn' as it does not exist")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorAlreadyDefinedFunction(t *testing.T) {
	err := transpileTest(`
		func print() int {
			return 0;
		}
	`)
	expected := fmt.Errorf("test.scri:2:3: Cannot declare function 'print' as it already is defined")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorSleepFuncCallWithWrongParamVarType(t *testing.T) {
	err := transpileTest(`
		str s = "123";
		sleep(s);
	`)
	expected := fmt.Errorf("test.scri:3:3: sleep() - Parameter seconds must be an int or a variable of type int. Got 'StrType'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorSleepFuncCallWithWrongParamType(t *testing.T) {
	err := transpileTest(`sleep("123");`)
	expected := fmt.Errorf("test.scri:1:1: sleep() - Parameter seconds must be an int or a variable of type int. Got 'StrLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorInputFuncCallWithWrongParamVarType(t *testing.T) {
	err := transpileTest(`
		int i = 42;
		input(i);
	`)
	expected := fmt.Errorf("test.scri:3:3: input() - Parameter prompt must be a string or a variable of type string. Got 'IntType'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorInputFuncCallWithWrongParamType(t *testing.T) {
	err := transpileTest(`input(42);`)
	expected := fmt.Errorf("test.scri:1:1: input() - Parameter prompt must be a string or a variable of type string. Got 'IntLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleFuncDeclarationWithCall() {
	setTestPrintMode()
	transpileTest(`
		# Function without params
		func funcWithoutParams() void {
			str str1 = "Test";
			printLn(str1);
		}

		# Function with params
		func funcWithParams(int a, str s) void {
			int b = a;
			str t = s;
			printLn(a, b, s, t);
		}

		funcWithoutParams();
		funcWithParams(123, "abc");
		int i = 123;
		str s = "abc";
		funcWithParams(i, s);

		# Function with return value
		func add(int a, int b) int {
			return a + b;
		}
		int sum = add(i, 321);
		printLn(add(123, 321));
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// # Function without params
	// funcWithoutParams () {
	// 	local str1="Test"
	// 	echo "${str1}"
	// }
	//
	// # Function with params
	// funcWithParams () {
	// 	local a=$1
	// 	local s=$2
	// 	local b=${a}
	// 	local t="${s}"
	// 	echo "${a} ${b} ${s} ${t}"
	// }
	//
	// funcWithoutParams
	// funcWithParams 123 "abc"
	// i=123
	// s="abc"
	// funcWithParams ${i} "${s}"
	// # Function with return value
	// add () {
	// 	local a=$1
	// 	local b=$2
	// 	tmpInt=$((${a} + ${b}))
	// 	return
	// }
	//
	// add ${i} 321
	// sum=${tmpInt}
	// add 123 321
	// echo "${tmpInt}"
}

func TestErrorFuncCallWithWrongTypes(t *testing.T) {
	err := transpileTest(`
		func fn(int a) void {
			printLn(a);
		}
		fn("hello");
	`)
	expected := fmt.Errorf("test.scri:5:3: fn(): Parameter 'a' type does not match. Expected: IntType, Got: str")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorFuncCallWithWrongAmountOfArgs(t *testing.T) {
	err := transpileTest(`
		func fn(int a) void {
			printLn(a);
		}
		fn(1, 2);
	`)
	expected := fmt.Errorf("test.scri:5:3: fn(): The amount of passed parameters does not match with the function declaration. Expected: 1, Got: 2")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorFuncDeclWithMissingType(t *testing.T) {
	err := transpileTest(`
		func fn(int a) {
			printLn(a);
		}
	`)
	expected := fmt.Errorf("test.scri:2:18: Return type is missing")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorFuncDeclWithUnsupportedType(t *testing.T) {
	err := transpileTest(`
		func fn(int a) const {
			printLn(a);
		}
	`)
	expected := fmt.Errorf("test.scri:2:18: Unsupported return type 'const'")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorFuncVoidReturnValUsed(t *testing.T) {
	err := transpileTest(`
		func fn() void {
			printLn(123);
		}
		int i = fn();
	`)
	expected := fmt.Errorf("test.scri:5:11: Func 'fn' does not have a return value")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorInvalidPropertyName(t *testing.T) {
	err := transpileTest(`
		obj o = { a: 1, };
		o.1 = 32;
	`)
	expected := fmt.Errorf("test.scri:3:5: Cannot use dot operator without right hand side being an identifier")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleObject() {
	setTestPrintMode()
	transpileTest(`
		obj o = { p1: 123, p2: "str", p3: false, };
		o.p1 = 321;
		printLn(o.p2);
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// declare -A o
	// o["p1"]=123
	// o["p2"]="str"
	// o["p3"]=false
	// o["p1"]=321
	// echo "${o["p2"]}"
}

func TestErrorObjectWithMissingComma(t *testing.T) {
	err := transpileTest(`
		obj o = { p1: 123 };
	`)
	expected := fmt.Errorf("test.scri:2:21: Expected comma following Property")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorObjectWithMissingColon(t *testing.T) {
	err := transpileTest(`
		obj o = { p1 };
	`)
	expected := fmt.Errorf("test.scri:2:16: Missing colon following identifier in ObjectExpr")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorObjectWithMissingValue(t *testing.T) {
	err := transpileTest(`
		obj o = { p1: , };
	`)
	expected := fmt.Errorf("test.scri:2:17: Unexpected token 'Comma' (',') found during parsing")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorMemberExprWithObjectOfWrongType(t *testing.T) {
	err := transpileTest(`
		int i = 42;
		i.a = 1;
	`)
	expected := fmt.Errorf("test.scri:3:3: Variable 'i' is not of type 'object'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorMemberExprWithWrongObjectNameType(t *testing.T) {
	err := transpileTest(`
		int i = 42;
		1.a = 1;
	`)
	expected := fmt.Errorf("test.scri:3:3: Object name is not the right type. Got 'IntLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorInputWithoutPrompt(t *testing.T) {
	err := transpileTest(`
		input();
	`)
	expected := fmt.Errorf("test.scri:2:3: Expected syntax: input(str prompt)")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleInput() {
	setTestPrintMode()
	transpileTest(`
		str s = input("Enter username:");
		input(s);
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// read -p "Enter username: " tmpStr
	// s="${tmpStr}"
	// read -p "${s} " tmpStr
}

func TestErrorSleepWithoutSeconds(t *testing.T) {
	err := transpileTest(`
		sleep();
	`)
	expected := fmt.Errorf("test.scri:2:3: Expected syntax: sleep(int seconds)")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleSleep() {
	setTestPrintMode()
	transpileTest(`
		sleep(10);
		int i = 10;
		sleep(i);
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// sleep 10
	// i=10
	// sleep ${i}
}
