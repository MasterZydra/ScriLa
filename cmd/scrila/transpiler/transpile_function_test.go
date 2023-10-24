package transpiler

import (
	"fmt"
	"strings"
	"testing"
)

// -------- Native function "Exec" --------

func TestErrorExecWithoutValue(t *testing.T) {
	initTest()
	err := transpileTest(`exec();`)
	expected := fmt.Errorf("test.scri:1:1: Expected syntax: exec(str command)")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorExecWithWrongArgType(t *testing.T) {
	initTest()
	err := transpileTest(`exec(123);`)
	expected := fmt.Errorf("test.scri:1:1: exec() - Parameter value must be a string or a variable of type string. Got 'IntLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleExec() {
	initTestForPrintMode()
	transpileTest(`
		exec("echo hi");
		str cmd = "echo hi";
		exec(cmd);
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	//
	// # User script
	// echo hi
	// cmd="echo hi"
	// ${cmd}
}

// -------- Native function "Input" --------

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

func TestErrorNativeFuncReturnTypeUnequalVarType(t *testing.T) {
	initTest()
	err := transpileTest(`int i = input("prompt");`)
	expected := fmt.Errorf("test.scri:1:9: Cannot assign a value of type 'StrType' to a var of type 'IntType'")
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

// -------- Native function "Print" --------

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
		printLn();
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
	// echo ""
}

// -------- Native function "Sleep" --------

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

// -------- Native function "StrIsInt" --------

func TestErrorStrIsIntWithoutValue(t *testing.T) {
	initTest()
	err := transpileTest(`strIsInt();`)
	expected := fmt.Errorf("test.scri:1:1: Expected syntax: strIsInt(mixed value)")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleStrIsInt() {
	initTestForPrintMode()
	transpileTest(`
		bool b = strIsInt(10);
		b = strIsInt("10");
		b = strIsInt("str");
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	//
	// # Native function implementations
	// strIsInt () {
	//	case $1 in
	//		''|*[!0-9]*) tmpBool="false" ;;
	// 		*) tmpBool="true" ;;
	// 	esac
	// }
	//
	// # User script
	// strIsInt 10
	// b="${tmpBool}"
	// strIsInt "10"
	// b="${tmpBool}"
	// strIsInt "str"
	// b="${tmpBool}"
}

// -------- Native function "StrToInt" --------

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

// -------- User defined functions --------

func TestErrorInvalidFuncCallName(t *testing.T) {
	initTest()
	err := transpileTest(`12();`)
	expected := fmt.Errorf("test.scri:1:1: Function name must be an identifier. Got: 'IntLiteral'")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
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
