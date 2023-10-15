package transpiler

import (
	"ScriLa/cmd/scrila/parser"
	"fmt"
	"strings"
	"testing"
)

var testTranspiler *Transpiler

func initTest() {
	testTranspiler = NewTranspiler()
	testTranspiler.filename = "test.scri"
	testTranspiler.testMode = true
}

func initTestForPrintMode() {
	initTest()
	testTranspiler.testPrintMode = true
}

func transpileTest(code string) error {
	parser := parser.NewParser()
	env := NewEnvironment(nil, testTranspiler)
	program, err := parser.ProduceAST(code, testTranspiler.filename)
	if err != nil {
		return err
	}
	return testTranspiler.Transpile(program, env, "")
}

func TestErrorLexerUnrecognizedChar(t *testing.T) {
	initTest()
	err := transpileTest(`~`)
	expected := fmt.Errorf("test.scri:1:1: Unrecognized character '~' found")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExamplePrint() {
	initTestForPrintMode()
	transpileTest(`
		# Print with(out) linebreaks
		print("Hello ");
		printLn("World");
		printLn("!");
		# Print base types
		printLn(42, "str", true, false, null);
		# Print variables
		int i = 42;
		str s = "hello world";
		bool b = false;
		printLn(i, s, b);
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	//
	// # User script
	// # Print with(out) linebreaks
	// echo -n "Hello "
	// echo "World"
	// echo "!"
	// # Print base types
	// echo "42 str true false null"
	// # Print variables
	// i=42
	// s="hello world"
	// b="false"
	// echo "${i} ${s} ${b}"

}

func ExampleIntVar() {
	initTestForPrintMode()
	transpileTest(`
		# Declare and assign new value
		int i = 42;
		i = 101;
		# Declare with binary expr
		int j = 42 * 2;
		# Declare with binary expr with var
		int k = i * 2;
		k = (i + 2) * i;
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	//
	// # User script
	// # Declare and assign new value
	// i=42
	// i=101
	// # Declare with binary expr
	// j=$((42 * 2))
	// # Declare with binary expr with var
	// k=$((${i} * 2))
	// k=$(($((${i} + 2)) * ${i}))
}

func TestErrorAssignWrongLeftSide(t *testing.T) {
	initTest()
	err := transpileTest(`12 = 34;`)
	expected := fmt.Errorf("test.scri:1:1: Left side of an assignment must be a variable. Got 'IntLiteral'")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorUnsupportedStringOperation(t *testing.T) {
	initTest()
	err := transpileTest(`str s = "str" - "str";`)
	expected := fmt.Errorf("test.scri:1:15: Binary string expression with unsupported operator '-'")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorBinaryExprWithUnsupportedCombination(t *testing.T) {
	initTest()
	err := transpileTest(`int i = "str" - 123;`)
	expected := fmt.Errorf("test.scri:1:15: No support for binary expressions of type 'str' and 'int'")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorIntDeclarationWithMissingSemicolon(t *testing.T) {
	initTest()
	err := transpileTest(`int i = 42`)
	expected := fmt.Errorf("test.scri:1:11: Expression must end with a semicolon")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorDoubleVariableDeclration(t *testing.T) {
	initTest()
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
	initTest()
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
	initTest()
	err := transpileTest(`i = 42;`)
	expected := fmt.Errorf("test.scri:1:1: Cannot resolve variable 'i' as it does not exist")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleStrAssignmentBinaryExprWithVar() {
	initTestForPrintMode()
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
	//
	// # User script
	// a="Hello"
	// b="World"
	// c="${a} ${b}"
	// d="${a} World"
	// d="${a} World"
}

func ExampleVarDeclarationAndAssignmentWithVariable() {
	initTestForPrintMode()
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
	//
	// # User script
	// i=123
	// j=${i}
	// j=${i}
	// s="str"
	// t="${s}"
	// t="${s}"
}

func TestErrorAssignDifferentVarTypes(t *testing.T) {
	initTest()
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
	initTest()
	err := transpileTest(`const func i = 13;`)
	expected := fmt.Errorf("test.scri:1:7: Variable type 'func' not given or supported")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorDeclareDifferentVarTypes(t *testing.T) {
	initTest()
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
	initTest()
	err := transpileTest(`int i = "123";`)
	expected := fmt.Errorf("test.scri:1:11: Cannot assign a value of type 'StrType' to a var of type 'IntType'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorAssignDifferentType(t *testing.T) {
	initTest()
	err := transpileTest(`
		int i = 123;
		i = "456";
	`)
	expected := fmt.Errorf("test.scri:3:9: Cannot assign a value of type 'StrType' to a var of type 'IntType'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorReturnOutsideOfFunction(t *testing.T) {
	initTest()
	err := transpileTest(`return true;`)
	expected := fmt.Errorf("test.scri:1:1: Return is only allowed inside a function")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorInvalidFuncCallName(t *testing.T) {
	initTest()
	err := transpileTest(`12();`)
	expected := fmt.Errorf("test.scri:1:1: Function name must be an identifier. Got: 'IntLiteral'")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorFuncParamsWithUnexpectedToken(t *testing.T) {
	initTest()
	err := transpileTest(`func fn(int a const) void {}`)
	expected := fmt.Errorf("test.scri:1:15: Unexpected token 'const' in parameter list")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorMissingFuncParamType(t *testing.T) {
	initTest()
	err := transpileTest(`func fn(a) void {}`)
	expected := fmt.Errorf("test.scri:1:9: Expected param type but got Identifier 'a'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorNonexistentFunc(t *testing.T) {
	initTest()
	err := transpileTest(`fn();`)
	expected := fmt.Errorf("test.scri:1:1: Cannot resolve function 'fn' as it does not exist")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorAlreadyDefinedFunction(t *testing.T) {
	initTest()
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
	initTest()
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
	initTest()
	err := transpileTest(`sleep("123");`)
	expected := fmt.Errorf("test.scri:1:1: sleep() - Parameter seconds must be an int or a variable of type int. Got 'StrLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorInputFuncCallWithWrongParamVarType(t *testing.T) {
	initTest()
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
	initTest()
	err := transpileTest(`input(42);`)
	expected := fmt.Errorf("test.scri:1:1: input() - Parameter prompt must be a string or a variable of type string. Got 'IntLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleFuncDeclarationWithCall() {
	initTestForPrintMode()
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
	//
	// # User script
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
	initTest()
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
	initTest()
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
	initTest()
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
	initTest()
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
	initTest()
	err := transpileTest(`
		func fn() void {
			printLn(123);
		}
		int i = fn();
	`)
	expected := fmt.Errorf("test.scri:5:11: Cannot assign a value of type 'VoidType' to a var of type 'IntType'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorNativeFuncReturnTypeUnequalVarType(t *testing.T) {
	initTest()
	err := transpileTest(`int i = input("prompt");`)
	expected := fmt.Errorf("test.scri:1:9: Cannot assign a value of type 'StrType' to a var of type 'IntType'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorInvalidPropertyName(t *testing.T) {
	initTest()
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
	initTestForPrintMode()
	transpileTest(`
		obj o = { p1: 123, p2: "str", p3: false, };
		o.p1 = 321;
		printLn(o.p2);
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	//
	// # User script
	// declare -A o
	// o["p1"]=123
	// o["p2"]="str"
	// o["p3"]=false
	// o["p1"]=321
	// echo "${o["p2"]}"
}

func TestErrorObjectWithMissingComma(t *testing.T) {
	initTest()
	err := transpileTest(`
		obj o = { p1: 123 };
	`)
	expected := fmt.Errorf("test.scri:2:21: Expected comma following Property")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorObjectWithMissingColon(t *testing.T) {
	initTest()
	err := transpileTest(`
		obj o = { p1 };
	`)
	expected := fmt.Errorf("test.scri:2:16: Missing colon following identifier in ObjectExpr")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorObjectWithMissingValue(t *testing.T) {
	initTest()
	err := transpileTest(`
		obj o = { p1: , };
	`)
	expected := fmt.Errorf("test.scri:2:17: Unexpected token 'Comma' (',') found during parsing")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorMemberExprWithObjectOfWrongType(t *testing.T) {
	initTest()
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
	initTest()
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
	initTest()
	err := transpileTest(`input();`)
	expected := fmt.Errorf("test.scri:1:1: Expected syntax: input(str prompt)")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleInput() {
	initTestForPrintMode()
	transpileTest(`
		str s = input("Enter username:");
		input(s);
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	//
	// # User script
	// read -p "Enter username: " tmpStr
	// s="${tmpStr}"
	// read -p "${s} " tmpStr
}

func TestErrorSleepWithoutSeconds(t *testing.T) {
	initTest()
	err := transpileTest(`sleep();`)
	expected := fmt.Errorf("test.scri:1:1: Expected syntax: sleep(int seconds)")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleSleep() {
	initTestForPrintMode()
	transpileTest(`
		sleep(10);
		int i = 10;
		sleep(i);
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	//
	// # User script
	// sleep 10
	// i=10
	// sleep ${i}
}

func TestErrorIsIntWithoutValue(t *testing.T) {
	initTest()
	err := transpileTest(`isInt();`)
	expected := fmt.Errorf("test.scri:1:1: Expected syntax: isInt(mixed value)")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleIsInt() {
	initTestForPrintMode()
	transpileTest(`
		bool b = isInt(10);
		b = isInt("10");
		b = isInt("str");
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	//
	// # Native function implementations
	// isInt () {
	//	case $1 in
	//		''|*[!0-9]*) tmpBool="false" ;;
	// 		*) tmpBool="true" ;;
	// 	esac
	// }
	//
	// # User script
	// isInt 10
	// b="${tmpBool}"
	// isInt "10"
	// b="${tmpBool}"
	// isInt "str"
	// b="${tmpBool}"
}

func TestErrorStrToIntWithoutValue(t *testing.T) {
	initTest()
	err := transpileTest(`strToInt();`)
	expected := fmt.Errorf("test.scri:1:1: Expected syntax: strToInt(str value)")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorStrToIntWithWrongArgType(t *testing.T) {
	initTest()
	err := transpileTest(`strToInt(123);`)
	expected := fmt.Errorf("test.scri:1:1: strToInt() - Parameter value must be a string or a variable of type string. Got 'IntLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleStrToInt() {
	initTestForPrintMode()
	transpileTest(`int i = strToInt("123");`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	//
	// # User script
	// tmpInt="123"
	// i=${tmpInt}
}

func TestErrorIfWithoutOpenParen(t *testing.T) {
	initTest()
	err := transpileTest(`
		if true {}
	`)
	expected := fmt.Errorf("test.scri:2:6: Expected condition wrapped in parentheses")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorIfWithoutOpenBrace(t *testing.T) {
	initTest()
	err := transpileTest(`
		if (true)
		printLn("str");
	`)
	expected := fmt.Errorf("test.scri:3:3: Expected block following condition")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorIfWithWrongBinaryExprType(t *testing.T) {
	initTest()
	err := transpileTest(`
		if (1 + 1) {
			printLn("str");
		}
	`)
	expected := fmt.Errorf("test.scri:2:9: Condition is no boolean expression. Got int")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorIfWithWrongVarType(t *testing.T) {
	initTest()
	err := transpileTest(`
		int i = 42;
		if (i) {
			printLn("str");
		}
	`)
	expected := fmt.Errorf("test.scri:3:7: Condition is not of type bool. Got IntType")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleIf() {
	initTestForPrintMode()
	transpileTest(`
		if (true) {
			printLn("true");
		}
		if (true && false) {
			printLn("true");
		}
		if (true || false) {
			printLn("true");
		}
		bool b = true;
		if (b) {
			printLn("true");
		}
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	//
	// # User script
	// if [[ "true" == "true" ]]; then
	// 	echo "true"
	// fi
	// if [[ "true" == "true" ]] && [[ "false" == "true" ]]; then
	// 	echo "true"
	// fi
	// if [[ "true" == "true" ]] || [[ "false" == "true" ]]; then
	// 	echo "true"
	// fi
	// b="true"
	// if [[ "${b}" == "true" ]]; then
	// 	echo "true"
	// fi
}

func ExampleIfComparisons() {
	initTestForPrintMode()
	transpileTest(`
		int i = 123;
		if (i > 122) {
			printLn(true);
		}
		if (i < 124) {
			printLn(true);
		}
		if (i >= 122) {
			printLn(true);
		}
		if (122 <= i) {
			printLn(true);
		}
		if (123 == i) {
			printLn(true);
		}
		if (i != 321) {
			printLn(true);
		}
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	//
	// # User script
	// i=123
	// if [[ ${i} -gt 122 ]]; then
	// 	echo "true"
	// fi
	// if [[ ${i} -lt 124 ]]; then
	// 	echo "true"
	// fi
	// if [[ ${i} -ge 122 ]]; then
	// 	echo "true"
	// fi
	// if [[ 122 -le ${i} ]]; then
	// 	echo "true"
	// fi
	// if [[ 123 -eq ${i} ]]; then
	// 	echo "true"
	// fi
	// if [[ ${i} -ne 321 ]]; then
	// 	echo "true"
	// fi
}

func TestErrorCompareDiffVarTypes(t *testing.T) {
	initTest()
	err := transpileTest(`bool b = 42 > "123";`)
	expected := fmt.Errorf("test.scri:1:13: Cannot compare type 'int' and 'str'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}
