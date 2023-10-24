package transpiler

import (
	"fmt"
	"strings"
	"testing"
)

// -------- While --------

func TestErrorWhileWithoutOpenParen(t *testing.T) {
	initTest()
	err := transpileTest(`while true {}`)
	expected := fmt.Errorf("test.scri:1:7: Expected condition wrapped in parentheses")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorWhileWithoutOpenBrace(t *testing.T) {
	initTest()
	err := transpileTest(`
		while (true)
		printLn("str");
	`)
	expected := fmt.Errorf("test.scri:3:3: Expected block following condition")
	if !strings.HasPrefix(err.Error(), expected.Error()) {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorWhileWithWrongBinaryExprType(t *testing.T) {
	initTest()
	err := transpileTest(`
		while (1 + 1) {
			printLn("str");
		}
	`)
	expected := fmt.Errorf("test.scri:2:12: Condition is no boolean expression. Got int")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func TestErrorWhileWithWrongVarType(t *testing.T) {
	initTest()
	err := transpileTest(`
		int i = 42;
		while (i) {
			printLn("str");
		}
	`)
	expected := fmt.Errorf("test.scri:3:10: Condition is not of type bool. Got IntType")
	if err.Error() != expected.Error() {
		t.Errorf("Expected: \"%s\", Got: \"%s\"", expected, err)
	}
}

func ExampleWhile() {
	initTestForPrintMode()
	transpileTest(`
		while (true && false) {
			printLn("true");
		}
		while (true || false) {
		}
		bool b = true;
		while (b) {
			printLn("true");
		}
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	//
	// # User script
	// while [[ "true" == "true" ]] && [[ "false" == "true" ]]
	// do
	// 	echo "true"
	// done
	// while [[ "true" == "true" ]] || [[ "false" == "true" ]]
	// do
	// 	:
	// done
	// b="true"
	// while [[ "${b}" == "true" ]]
	// do
	// 	echo "true"
	// done
}

// -------- If --------

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
		} else if (true && false) {
			printLn("true");
		}
		if (true || false) {
			printLn("true");
		} else {
			printLn("false");
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
	// if [[ "true" == "true" ]]
	// then
	// 	echo "true"
	// elif [[ "true" == "true" ]] && [[ "false" == "true" ]]
	// then
	// 	echo "true"
	// fi
	// if [[ "true" == "true" ]] || [[ "false" == "true" ]]
	// then
	// 	echo "true"
	// else
	// 	echo "false"
	// fi
	// b="true"
	// if [[ "${b}" == "true" ]]
	// then
	// 	echo "true"
	// fi
}

func ExampleIfComparisons() {
	initTestForPrintMode()
	transpileTest(`
		# Integer comparison
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
		# String comparison
		str s = "str";
		if (s > "ttr") {
			printLn(true);
		}
		if (s < "ttr") {
			printLn(true);
		}
		if ("str" == s) {
			printLn(true);
		}
		if (s != "str") {
		}
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	//
	// # User script
	// # Integer comparison
	// i=123
	// if [[ ${i} -gt 122 ]]
	// then
	// 	echo "true"
	// fi
	// if [[ ${i} -lt 124 ]]
	// then
	// 	echo "true"
	// fi
	// if [[ ${i} -ge 122 ]]
	// then
	// 	echo "true"
	// fi
	// if [[ 122 -le ${i} ]]
	// then
	// 	echo "true"
	// fi
	// if [[ 123 -eq ${i} ]]
	// then
	// 	echo "true"
	// fi
	// if [[ ${i} -ne 321 ]]
	// then
	// 	echo "true"
	// fi
	// # String comparison
	// s="str"
	// if [[ "${s}" > "ttr" ]]
	// then
	// 	echo "true"
	// fi
	// if [[ "${s}" < "ttr" ]]
	// then
	// 	echo "true"
	// fi
	// if [[ "str" == "${s}" ]]
	// then
	// 	echo "true"
	// fi
	// if [[ "${s}" != "str" ]]
	// then
	// 	:
	// fi
}
