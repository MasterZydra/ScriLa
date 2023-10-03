package transpiler

import (
	"ScriLa/cmd/scrila/parser"
	"fmt"
	"strings"
	"testing"
)

func setTestMode() {
	testMode = true
}

func transpileTest(code string) error {
	parser := parser.NewParser()
	env := NewEnvironment(nil)

	program, err := parser.ProduceAST(code)
	if err != nil {
		return err
	}
	return Transpile(program, env, "")
}

func ExmaplePrint() {
	setTestMode()
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
	setTestMode()
	transpileTest(`
		printLn(42, "str", true, false, null);
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// echo "42 str true false null"
}

func ExamplePrintVariables() {
	setTestMode()
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
	// b=false
	// echo "${i} ${s} ${b}"
}

func ExampleIntDeclaration() {
	setTestMode()
	transpileTest(`
		int i = 42;
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// i=42
}

func ExampleIntAssignment() {
	setTestMode()
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

func TestIntDeclarationWithMissingSemicolon(t *testing.T) {
	err := transpileTest(`
		int i = 42
	`)
	expected := fmt.Errorf("Parser Error: Expressions must end with a semicolon.")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestIntAssignmentWithMissingDeclaration(t *testing.T) {
	err := transpileTest(`
		i = 42;
	`)
	expected := fmt.Errorf("Cannot resolve variable 'i' as it does not exist.")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleIntAssignmentBinaryExpr() {
	setTestMode()
	transpileTest(`
		int i = 42 * 2;
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// i=$((42 * 2))
}

func ExampleIntAssignmentBinaryExprWithVar() {
	setTestMode()
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
	setTestMode()
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
	setTestMode()
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

func TestAssignDifferentVarTypes(t *testing.T) {
	err := transpileTest(`
		int i = 123;
		str s = "str";
		s = i;
	`)
	expected := fmt.Errorf("Cannot assign a value of type 'IntType' to a var of type 'StrType'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestDeclareDifferentVarTypes(t *testing.T) {
	err := transpileTest(`
		int i = 123;
		str s = i;
	`)
	expected := fmt.Errorf("Cannot assign a value of type 'IntType' to a var of type 'StrType'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestDeclareDifferentType(t *testing.T) {
	err := transpileTest(`
		int i = "123";
	`)
	expected := fmt.Errorf("Cannot assign a value of type 'StrType' to a var of type 'IntType'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestAssignDifferentType(t *testing.T) {
	err := transpileTest(`
		int i = 123;
		i = "456";
	`)
	expected := fmt.Errorf("Cannot assign a value of type 'StrType' to a var of type 'IntType'")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleComment() {
	setTestMode()
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

func ExampleFuncDeclarationWithCall() {
	setTestMode()
	transpileTest(`
		# Function without params
		func funcWithoutParams() {
			str str1 = "Test";
			printLn(str1);
		}

		# Function with params
		func funcWithParams(int a, str s) {
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
		func add(int a, int b) {
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
	// 	return $((${a} + ${b}))
	// }
	//
	// add ${i} 321
	// sum=$?
	// add 123 321
	// echo "$?"
}

func ExampleObject() {
	setTestMode()
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

func TestObjectWithMissingComma(t *testing.T) {
	err := transpileTest(`
		int o = { p1: 123 };
	`)
	expected := fmt.Errorf("Parser Error: Expected comma following Property.")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestObjectWithMissingColon(t *testing.T) {
	err := transpileTest(`
		int o = { p1 };
	`)
	expected := fmt.Errorf("Parser Error: Missing colon following identifier in ObjectExpr.")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestObjectWithMissingValue(t *testing.T) {
	err := transpileTest(`
		int o = { p1: , };
	`)
	expected := fmt.Errorf("parsePrimaryExpr: Unexpected token 'Comma'")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestInputWithoutPrompt(t *testing.T) {
	err := transpileTest(`
		input();
	`)
	expected := fmt.Errorf("Expected syntax: input(str prompt)")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleInput() {
	setTestMode()
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

func TestSleepWithoutSeconds(t *testing.T) {
	err := transpileTest(`
		sleep();
	`)
	expected := fmt.Errorf("Expected syntax: sleep(int seconds)")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleSleep() {
	setTestMode()
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
