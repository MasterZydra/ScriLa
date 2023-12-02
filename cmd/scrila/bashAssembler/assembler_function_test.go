package bashAssembler

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
	//
	// # Native function implementations
	//
	// # exec(str command) void
	// exec () {
	// 	local command=$1
	// 	${command}
	// }
	//
	// # User script
	//
	// exec "echo hi"
	// cmd="echo hi"
	// exec "${cmd}"
}

// -------- Native function "Exit" --------

func TestErrorExitWithoutValue(t *testing.T) {
	initTest()
	err := transpileTest(`exit();`)
	expected := fmt.Errorf("test.scri:1:1: Expected syntax: exit(int code)")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorExitWithWrongArgType(t *testing.T) {
	initTest()
	err := transpileTest(`exit("123");`)
	expected := fmt.Errorf("test.scri:1:1: exit() - Parameter value must be a int or a variable of type int. Got 'StrLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleExit() {
	initTestForPrintMode()
	transpileTest(`
		exit(0);
		int code = 1;
		exit(code);
	`)

	// Output:
	// #!/bin/bash
	//
	// # User script
	//
	// exit 0
	// code=1
	// exit ${code}
}

// -------- Native function "Input" --------

func TestErrorInputFuncCallWithWrongParamVarType(t *testing.T) {
	initTest()
	err := transpileTest(`
		int i = 42;
		input(i);
	`)
	expected := fmt.Errorf("test.scri:3:3: input() - Parameter prompt must be a string or a variable of type string. Got 'IntLiteral'")
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
	expected := fmt.Errorf("test.scri:1:9: Cannot assign a value of type 'StrLiteral' to a var of type 'IntLiteral'")
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
	//
	// # Native function implementations
	//
	// # input(str prompt) str
	// input () {
	// 	local prompt=$1
	// 	read -p "${prompt} " tmpStrs[${tmpInts[0]}]
	// }
	//
	// # User script
	//
	// tmpInts[0]=1
	// input "Enter username:"
	// s="${tmpStrs[1]}"
	// input "${s}"
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
		printLn(42, "str", true, false);
		# Print variables
		int i = 42;
		str s = "hello world";
		bool b = false;
		printLn(i, s, b);
		printLn();
		# Print with function call and binary op
		printLn(strToInt("123")+2);
	`)

	// Output:
	// #!/bin/bash
	//
	// # Native function implementations
	//
	// # strToInt(str value) int
	// strToInt () {
	// 	local value=$1
	// 	tmpInts[${tmpInts[0]}]=${value}
	// }
	//
	// # User script
	//
	// # Print with(out) linebreaks
	// echo -n "Hello "
	// echo "World"
	// echo "!"
	// # Print base types
	// echo "42 str true false"
	// # Print variables
	// i=42
	// s="hello world"
	// b="false"
	// echo "${i} ${s} ${b}"
	// echo ""
	// # Print with function call and binary op
	// tmpInts[0]=1
	// strToInt "123"
	// echo "$((${tmpInts[1]} + 2))"
}

// -------- Native function "Sleep" --------

func TestErrorSleepFuncCallWithWrongParamVarType(t *testing.T) {
	initTest()
	err := transpileTest(`
		str s = "123";
		sleep(s);
	`)
	expected := fmt.Errorf("test.scri:3:3: sleep() - Parameter seconds must be an int or a variable of type int. Got 'StrLiteral'")
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
	//
	// # User script
	//
	// sleep 10
	// i=10
	// sleep ${i}
}

// -------- Native function "StrIsBool" --------

func TestErrorStrIsBoolWithoutValue(t *testing.T) {
	initTest()
	err := transpileTest(`strIsBool();`)
	expected := fmt.Errorf("test.scri:1:1: Expected syntax: strIsBool(str value)")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorStrIsBoolWithWrongArgType(t *testing.T) {
	initTest()
	err := transpileTest(`strIsBool(123);`)
	expected := fmt.Errorf("test.scri:1:1: strIsBool() - Parameter value must be a string or a variable of type string. Got 'IntLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleStrIsBool() {
	initTestForPrintMode()
	transpileTest(`
		bool b = strIsBool("10");
		b = strIsBool("true");
	`)

	// Output:
	// #!/bin/bash
	//
	// # Native function implementations
	//
	// # strIsBool(str value) bool
	// strIsBool () {
	// 	local value=$1
	// 	if [[ "${value}" == "true" ]] || [[ "${value}" == "false" ]]
	// 	then
	// 		tmpBools[${tmpInts[0]}]="true"
	// 	else
	// 		tmpBools[${tmpInts[0]}]="false"
	// 	fi
	// }
	//
	// # User script
	//
	// tmpInts[0]=1
	// strIsBool "10"
	// b="${tmpBools[1]}"
	// strIsBool "true"
	// b="${tmpBools[1]}"
}

// -------- Native function "StrIsInt" --------

func TestErrorStrIsIntWithoutValue(t *testing.T) {
	initTest()
	err := transpileTest(`strIsInt();`)
	expected := fmt.Errorf("test.scri:1:1: Expected syntax: strIsInt(str value)")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorStrIsIntWithWrongArgType(t *testing.T) {
	initTest()
	err := transpileTest(`strIsInt(123);`)
	expected := fmt.Errorf("test.scri:1:1: strIsInt() - Parameter value must be a string or a variable of type string. Got 'IntLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleStrIsInt() {
	initTestForPrintMode()
	transpileTest(`
		bool b = strIsInt("10");
		b = strIsInt("str");
	`)

	// Output:
	// #!/bin/bash
	//
	// # Native function implementations
	//
	// # strIsInt(str value) bool
	// strIsInt () {
	// 	local value=$1
	// 	case ${value} in
	// 		''|*[!0-9]*) tmpBools[${tmpInts[0]}]="false" ;;
	// 		*) tmpBools[${tmpInts[0]}]="true" ;;
	// 	esac
	// }
	//
	// # User script
	//
	// tmpInts[0]=1
	// strIsInt "10"
	// b="${tmpBools[1]}"
	// strIsInt "str"
	// b="${tmpBools[1]}"
}

// -------- Native function "StrToBool" --------

func TestErrorStrToBoolWithoutValue(t *testing.T) {
	initTest()
	err := transpileTest(`strToBool();`)
	expected := fmt.Errorf("test.scri:1:1: Expected syntax: strToBool(str value)")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorStrToBoolWithWrongArgType(t *testing.T) {
	initTest()
	err := transpileTest(`strToBool(123);`)
	expected := fmt.Errorf("test.scri:1:1: strToBool() - Parameter value must be a string or a variable of type string. Got 'IntLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleStrToBool() {
	initTestForPrintMode()
	transpileTest(`
		bool b1 = strToBool("true");
		bool b2 = strToBool("false");
	`)

	// Output:
	// #!/bin/bash
	//
	// # Native function implementations
	//
	// # strToBool(str value) bool
	// strToBool () {
	// 	local value=$1
	// 	if [[ "${value}" == "true" ]]
	// 	then
	// 		tmpBools[${tmpInts[0]}]="true"
	// 	else
	// 		tmpBools[${tmpInts[0]}]="false"
	// 	fi
	// }
	//
	// # User script
	//
	// tmpInts[0]=1
	// strToBool "true"
	// b1="${tmpBools[1]}"
	// strToBool "false"
	// b2="${tmpBools[1]}"
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
	//
	// # Native function implementations
	//
	// # strToInt(str value) int
	// strToInt () {
	// 	local value=$1
	// 	tmpInts[${tmpInts[0]}]=${value}
	// }
	//
	// # User script
	//
	// tmpInts[0]=1
	// strToInt "123"
	// i=${tmpInts[1]}
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

		# Functions without body
		func emptyBody() void {
		}
		func emptyBodyJustComment() void {
			# Do smth
		}
	`)

	// Output:
	// #!/bin/bash
	//
	// # User script
	//
	// # Function without params
	// # funcWithoutParams() void
	// funcWithoutParams () {
	// 	local str1="Test"
	// 	echo "${str1}"
	// }
	//
	// # Function with params
	// # funcWithParams(int a, str s) void
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
	// # add(int a, int b) int
	// add () {
	// 	local a=$1
	// 	local b=$2
	// 	tmpInts[${tmpInts[0]}]=$((${a} + ${b}))
	// 	return
	// }
	//
	// tmpInts[0]=1
	// add ${i} 321
	// sum=${tmpInts[1]}
	// add 123 321
	// echo "${tmpInts[1]}"
	// # Functions without body
	// # emptyBody() void
	// emptyBody () {
	// 	:
	// }
	//
	// # emptyBodyJustComment() void
	// emptyBodyJustComment () {
	// 	# Do smth
	// 	:
	// }
}

func TestErrorFuncCallWithWrongTypes(t *testing.T) {
	initTest()
	err := transpileTest(`
		func fn(int a) void {
			printLn(a);
		}
		fn("hello");
	`)
	expected := fmt.Errorf("test.scri:5:3: fn(): Parameter 'a' type does not match. Expected: IntLiteral, Got: str")
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
	expected := fmt.Errorf("test.scri:5:11: Cannot assign a value of type 'Void' to a var of type 'IntLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorVoidFnReturnValue(t *testing.T) {
	initTest()
	err := transpileTest(`func retVoid() void { return 1; }`)
	expected := fmt.Errorf("test.scri:1:23: retVoid(): Cannot return value if function type is 'void'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorFuncReturnWrongType(t *testing.T) {
	initTest()
	err := transpileTest(`func retInt() int { return "123"; }`)
	expected := fmt.Errorf("test.scri:1:21: retInt(): Return type does not match with function type. Expected: IntLiteral, Got: str")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorFuncWithoutReturnValue(t *testing.T) {
	initTest()
	err := transpileTest(`func retInt() int { return; }`)
	expected := fmt.Errorf("test.scri:1:21: retInt(): Cannot return without a value for a function with return value")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorFuncInFunc(t *testing.T) {
	initTest()
	err := transpileTest(`func a() void { func b() void {} }`)
	expected := fmt.Errorf("test.scri:1:17: Cannot declare a function inside a function")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}
