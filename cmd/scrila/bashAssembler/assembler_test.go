package bashAssembler

import (
	"ScriLa/cmd/scrila/bashTranspiler"
	"ScriLa/cmd/scrila/config"
	"ScriLa/cmd/scrila/parser"
	"fmt"
	"strings"
	"testing"
)

var testAssembler *Assembler

func initTest() {
	testAssembler = NewAssembler()
	testAssembler.testMode = true
	config.Filename = "test.scri"
}

func initTestForPrintMode() {
	initTest()
	testAssembler.testPrintMode = true
}

func transpileTest(code string) error {
	transpiler := bashTranspiler.NewTranspiler()
	env := bashTranspiler.NewEnvironment(nil, transpiler)
	scrilaProgram, err := parser.NewParser().ProduceAST(code)
	if err != nil {
		if testAssembler.testPrintMode {
			fmt.Println(err)
		}
		return err
	}
	bashProgram, err := transpiler.Transpile(scrilaProgram, env)
	if err != nil {
		if testAssembler.testPrintMode {
			fmt.Println(err)
		}
		return err
	}
	err = testAssembler.Assemble(bashProgram)
	if testAssembler.testPrintMode && err != nil {
		fmt.Println(err)
	}
	return err
}

func TestErrorLexerUnrecognizedChar(t *testing.T) {
	initTest()
	err := transpileTest(`~`)
	expected := fmt.Errorf("test.scri:1:1: Unrecognized character '~' found")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
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
	//
	// # User script
	//
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
		str b = "\"World\"";
		str c = a + " " + b;
		str d = a + " World";
		d = a + " World";
	`)

	// Output:
	// #!/bin/bash
	//
	// # User script
	//
	// a="Hello"
	// b="\"World\""
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
	//
	// # User script
	//
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
	expected := fmt.Errorf("test.scri:4:7: Cannot assign a value of type 'IntLiteral' to a var of type 'StrLiteral'")
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
	expected := fmt.Errorf("test.scri:3:11: Cannot assign a value of type 'IntLiteral' to a var of type 'StrLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorDeclareDifferentType(t *testing.T) {
	initTest()
	err := transpileTest(`int i = "123";`)
	expected := fmt.Errorf("test.scri:1:11: Cannot assign a value of type 'StrLiteral' to a var of type 'IntLiteral'")
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
	expected := fmt.Errorf("test.scri:3:9: Cannot assign a value of type 'StrLiteral' to a var of type 'IntLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

// func TestErrorInvalidPropertyName(t *testing.T) {
// 	initTest()
// 	err := transpileTest(`
// 		obj o = { a: 1, };
// 		o.1 = 32;
// 	`)
// 	expected := fmt.Errorf("test.scri:3:5: Cannot use dot operator without right hand side being an identifier")
// 	if !strings.HasPrefix(err.Error(), expected.Error()) {
// 		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
// 	}
// }

// func ExampleObject() {
// 	initTestForPrintMode()
// 	transpileTest(`
// 		obj o = { p1: 123, p2: "str", p3: false, };
// 		o.p1 = 321;
// 		printLn(o.p2);
// 	`)

// 	// Output:
// 	// #!/bin/bash
// 	//
// 	// # User script
// 	// declare -A o
// 	// o["p1"]=123
// 	// o["p2"]="str"
// 	// o["p3"]=false
// 	// o["p1"]=321
// 	// echo "${o["p2"]}"
// }

// func TestErrorObjectWithMissingComma(t *testing.T) {
// 	initTest()
// 	err := transpileTest(`
// 		obj o = { p1: 123 };
// 	`)
// 	expected := fmt.Errorf("test.scri:2:21: Expected comma following Property")
// 	if !strings.HasPrefix(err.Error(), expected.Error()) {
// 		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
// 	}
// }

// func TestErrorObjectWithMissingColon(t *testing.T) {
// 	initTest()
// 	err := transpileTest(`
// 		obj o = { p1 };
// 	`)
// 	expected := fmt.Errorf("test.scri:2:16: Missing colon following identifier in ObjectExpr")
// 	if !strings.HasPrefix(err.Error(), expected.Error()) {
// 		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
// 	}
// }

// func TestErrorObjectWithMissingValue(t *testing.T) {
// 	initTest()
// 	err := transpileTest(`
// 		obj o = { p1: , };
// 	`)
// 	expected := fmt.Errorf("test.scri:2:17: Unexpected token 'Comma' (',') found during parsing")
// 	if err.Error() != expected.Error() {
// 		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
// 	}
// }

// func TestErrorMemberExprWithObjectOfWrongType(t *testing.T) {
// 	initTest()
// 	err := transpileTest(`
// 		int i = 42;
// 		i.a = 1;
// 	`)
// 	expected := fmt.Errorf("test.scri:3:3: Variable 'i' is not of type 'object'")
// 	if err.Error() != expected.Error() {
// 		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
// 	}
// }

// func TestErrorMemberExprWithWrongObjectNameType(t *testing.T) {
// 	initTest()
// 	err := transpileTest(`
// 		int i = 42;
// 		1.a = 1;
// 	`)
// 	expected := fmt.Errorf("test.scri:3:3: Object name is not the right type. Got 'IntLiteral'")
// 	if err.Error() != expected.Error() {
// 		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
// 	}
// }

func TestErrorCompareDiffVarTypes(t *testing.T) {
	initTest()
	err := transpileTest(`bool b = 42 > "123";`)
	expected := fmt.Errorf("test.scri:1:13: Cannot compare type 'int' and 'str'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleBoolAssignComparision() {
	initTestForPrintMode()
	transpileTest(`
		bool b = 42 > 13;
		bool b1 = true && false;

		func b(bool b) bool {
			return 42 > 13;
		}

		b = b(42 > 13);
	`)

	// Output:
	// #!/bin/bash
	//
	// # User script
	//
	// if [[ 42 -gt 13 ]]
	// then
	// 	tmpBools[0]="true"
	// else
	// 	tmpBools[0]="false"
	// fi
	// b="${tmpBools[0]}"
	// if [[ "true" == "true" ]] && [[ "false" == "true" ]]
	// then
	// 	tmpBools[0]="true"
	// else
	// 	tmpBools[0]="false"
	// fi
	// b1="${tmpBools[0]}"
	// # b(bool b) bool
	// b () {
	// 	local b=$1
	// 	if [[ 42 -gt 13 ]]
	// 	then
	// 		tmpBools[${tmpIndex}]="true"
	// 	else
	// 		tmpBools[${tmpIndex}]="false"
	// 	fi
	// 	return
	// }
	//
	// if [[ 42 -gt 13 ]]
	// then
	// 	tmpBools[0]="true"
	// else
	// 	tmpBools[0]="false"
	// fi
	// tmpIndex=1
	// tmpIndex=0
	// b "${tmpBools[0]}"
	// b="${tmpBools[0]}"
}

// Array

func TestErrorArrayAssignWrongArrayVarType(t *testing.T) {
	initTest()
	err := transpileTest(`int[] i = ["str"];`)
	expected := fmt.Errorf("test.scri:1:11: Cannot assign a value of type 'StrLiteral' to a var of type 'IntArray'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorArrayAppendWrongDataType(t *testing.T) {
	initTest()
	err := transpileTest(`
		int[] i = [42];
		i[] = "str";
	`)
	expected := fmt.Errorf("test.scri:3:11: Cannot assign a value of type 'StrLiteral' to array of type 'int-array'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorArrayIndexWrongDataType(t *testing.T) {
	initTest()
	err := transpileTest(`
		int[] i = [42];
		i["str"] = 43;
	`)
	expected := fmt.Errorf("test.scri:3:7: Array index is not the right type. Wanted 'IntLiteral'. Got 'StrLiteral'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleArray() {
	initTestForPrintMode()
	transpileTest(`
		# Declare
		# Empty array
		int[] i1 = [];
		# Array with value
		int[] i2 = [42];
		int[] i3 = [40, 41];

		# Assign
		# Empty array
		i1 = [];
		# Array with value
		i2 = [42];

		# Change
		i1[0] = 31;
		int one = 1;
		i2[one] = i2[0] + one;

		# Append
		i1[] = 32;
		int i44 = 44;
		i2[] = i44;

		# Print
		printLn(i1, i1[0]);

		# Array as return type
		func array() int[] {
			int[] tmpArray = [41, 42];
			return tmpArray;
		}
		int[] result = array();
	`)

	// Output:
	// #!/bin/bash
	//
	// # User script
	//
	// # Declare
	// # Empty array
	// i1=()
	// # Array with value
	// i2=(42)
	// i3=(40 41)
	// # Assign
	// # Empty array
	// i1=()
	// # Array with value
	// i2=(42)
	// # Change
	// i1[0]=31
	// one=1
	// i2[${one}]=$((${i2[0]} + ${one}))
	// # Append
	// i1+=(32)
	// i44=44
	// i2+=(${i44})
	// # Print
	// echo "${i1[@]} ${i1[0]}"
	// # Array as return type
	// # array() int[]
	// array () {
	// 	local tmpArray=(41 42)
	// 	tmpInts=${tmpArray[@]}
	// 	return
	// }
	//
	// tmpIndex=0
	// array
	// result=${tmpInts[@]}
}
