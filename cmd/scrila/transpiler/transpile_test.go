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
