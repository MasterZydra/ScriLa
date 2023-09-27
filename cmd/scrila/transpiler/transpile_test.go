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

func ExampleIntDeclaration() {
	setTestMode()
	transpileTest(`
		int i = 42;
		print(i);
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// i=42
	// echo "$i"
}

func ExampleIntAssignment() {
	setTestMode()
	transpileTest(`
		int i = 42;
		i = 101;
		print(i);
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// i=42
	// i=101
	// echo "$i"
}

func TestIntDeclarationWithMissingSemicolon(t *testing.T) {
	err := transpileTest(`
		int i = 42
		print(i);
	`)
	expected := fmt.Errorf("Parser Error: Expressions must end with a semicolon.")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestIntAssignmentWithMissingDeclaration(t *testing.T) {
	err := transpileTest(`
		i = 42;
		print(i);
	`)
	expected := fmt.Errorf("Cannot resolve variable 'i' as it does not exist.")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}
