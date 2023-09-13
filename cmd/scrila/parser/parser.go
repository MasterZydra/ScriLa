package parser

import (
	"ScriLa/cmd/scrila/ast"
	"ScriLa/cmd/scrila/lexer"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/exp/slices"
)

var additiveOps = []string{"+", "-"}

var multiplicitaveOps = []string{"*", "/"} // Modulo %

type Parser struct {
	lexer  *lexer.Lexer
	tokens []*lexer.Token
}

func NewParser() *Parser {
	return &Parser{
		lexer: lexer.NewLexer(),
	}
}

func (self Parser) ProduceAST(sourceCode string) ast.IProgram {
	self.tokens = self.lexer.Tokenize(sourceCode)
	program := ast.NewProgram()

	for self.notEOF() {
		program.Body = append(program.GetBody(), self.parseStatement())
	}

	return program
}

func (self *Parser) notEOF() bool {
	return self.tokens[0].TokenType != lexer.EndOfFile
}

func (self *Parser) parseStatement() ast.IStatement {
	switch self.at().TokenType {
	case lexer.Const, lexer.IntType, lexer.ObjType:
		return self.parseVarDeclaration()
	default:
		return self.parseExpr()
	}
}

// [const] [int|obj] IDENT = EXPR;
func (self *Parser) parseVarDeclaration() ast.IStatement {
	isConstant := self.at().TokenType == lexer.Const
	if isConstant {
		self.eat()
	}

	if self.at().TokenType != lexer.ObjType && self.at().TokenType != lexer.IntType {
		fmt.Println("Variable type not given or supported.", self.at())
		os.Exit(1)
	}
	self.eat()

	// TODO Check if type matches with result of parseExpr()
	identifier := self.expect(lexer.Identifier, "Expected identifier name following [const] [int] keywords.").Value
	self.expect(lexer.Equals, "Expected equals token following identifier in var declaration.")
	declaration := ast.NewVarDeclaration(isConstant, identifier, self.parseExpr())
	self.expect(lexer.Semicolon, "Variable declaration statement must end with semicolon.")
	return declaration
}

func (self *Parser) parseExpr() ast.IExpr {
	return self.parseAssignmentExpr()
}

// Orders of Prescidence:
// Processed from top to bottom
// Priority from bottom to top
// - AssignmentExpr
// - ObjectExpr
// - AdditiveExpr
// - MultiplicitiveExpr
// - CallExpr
// - MemberExr
// - PrimaryExpr

// - LogicalExpr
// - Comparison
// - UnaryExpr

func (self *Parser) parseAssignmentExpr() ast.IExpr {
	left := self.parseObjectExpr()

	if self.at().TokenType == lexer.Equals {
		self.eat() // Advance past equals

		value := self.parseAssignmentExpr() // This allows chaining e.g. x = y = 5;
		self.expect(lexer.Semicolon, "Variable assignment expr must end with semicolon.")
		return ast.NewAssignmentExpr(left, value)
	}

	return left
}

func (self *Parser) parseObjectExpr() ast.IExpr {
	// { Prop[] }

	if self.at().TokenType != lexer.OpenBrace {
		return self.parseAdditiveExpr()
	}
	self.eat() // Advance past open brace

	properties := make([]ast.IProperty, 0)
	for self.notEOF() && self.at().TokenType != lexer.CloseBrace {
		// { key: val, }

		key := self.expect(lexer.Identifier, "Object literal key expected.").Value
		self.expect(lexer.Colon, "Missing colon following identifier in ObjectExpr.")
		value := self.parseExpr()
		self.expect(lexer.Comma, "Expected comma following Property.")

		properties = append(properties, ast.NewProperty(key, value))
	}

	self.expect(lexer.CloseBrace, "Object literal missing closing brace.")

	return ast.NewObjectLiteral(properties)
}

func (self *Parser) parseAdditiveExpr() ast.IExpr {
	// Lefthand Prescidence
	//
	//      10 + 5 - 6      10 + (5 - 6)        10 * 5 - 6     10 + 5 * 6
	//
	//           o               o                   o              o
	//          /|\             /|\                 /|\            /|\
	//         / | \           / | \               / | \          / | \
	//        /  -  6        10  +  \             /  -  6       10  +  \
	//       o                       o           o                      o
	//      /|\                     /|\         /|\                    /|\
	//     / | \                   / | \       / | \                  / | \
	//   10  +  5                 5  -  6     10  *  5                5  *  6

	left := self.parseMultiplicitaveExpr()

	// Current token is an additive operator
	for slices.Contains(additiveOps, self.at().Value) {
		operator := self.eat().Value
		right := self.parseMultiplicitaveExpr()
		left = ast.NewBinaryExpr(left, right, operator)
	}

	return left
}

func (self *Parser) parseMultiplicitaveExpr() ast.IExpr {
	// Lefthand Prescidence (see func parseAdditiveExpr)
	left := self.parseCallMemberExpr()

	// Current token is a multiplicitave operator
	for slices.Contains(multiplicitaveOps, self.at().Value) {
		operator := self.eat().Value
		right := self.parseCallMemberExpr()
		left = ast.NewBinaryExpr(left, right, operator)
	}

	return left
}

// foo.x()
func (self *Parser) parseCallMemberExpr() ast.IExpr {
	member := self.parseMemberExpr()

	if self.at().TokenType == lexer.OpenParen {
		return self.parseCallExpr(member)
	}

	return member
}

// foo()
func (self *Parser) parseCallExpr(caller ast.IExpr) ast.IExpr {
	var callExpr ast.IExpr
	callExpr = ast.NewCallExpr(caller, self.parseArgs())

	// This allows chaining of function calls: e.g. foo()()
	if self.at().TokenType == lexer.OpenParen {
		callExpr = self.parseCallExpr(callExpr)
	}

	self.expect(lexer.Semicolon, "Function calls must end with a semicolon.")

	return callExpr
}

// func add(a, b) {} <- a & b are parameters
// add(a, b) <- a & b are now args (when calling)
func (self *Parser) parseArgs() []ast.IExpr {
	self.expect(lexer.OpenParen, "Expected open parenthesis")
	var args []ast.IExpr
	if self.at().TokenType == lexer.CloseParen {
		args = make([]ast.IExpr, 0)
	} else {
		args = self.parseArgumentsList()
	}

	self.expect(lexer.CloseParen, "Missing closing parenthesis inside arguments list")
	return args
}

// foo(x = 5, v = "Bar")
// Set x to 5 and v to "Bar" in global scope and pass values afterwards
func (self *Parser) parseArgumentsList() []ast.IExpr {
	args := []ast.IExpr{self.parseAssignmentExpr()}

	for self.notEOF() && self.at().TokenType == lexer.Comma {
		self.eat()
		args = append(args, self.parseAssignmentExpr())
	}

	return args
}

func (self *Parser) parseMemberExpr() ast.IExpr {
	object := self.parsePrimaryExpr()

	for self.at().TokenType == lexer.Dot || self.at().TokenType == lexer.OpenBracket {
		operator := self.eat()
		var property ast.IExpr
		var isComputed bool

		// Non-computed values aka "obj.expr"
		if operator.TokenType == lexer.Dot {
			isComputed = false
			// Get identifier
			property = self.parsePrimaryExpr()

			if property.GetKind() != ast.IdentifierNode {
				fmt.Println("Cannot use dot operator without right hand side being an identifier")
				os.Exit(1)
			}
		} else {
			isComputed = true
			// This allows chaining: obj[computedValue] e.g. obj1[obj2[getBar()]]
			property = self.parseExpr()

			self.expect(lexer.CloseBracket, "Missing closing bracket in computed value.")
		}

		object = ast.NewMemberExpr(object, property, isComputed)
	}

	return object
}

func (self *Parser) parsePrimaryExpr() ast.IExpr {
	switch self.at().TokenType {
	case lexer.Identifier:
		return ast.NewIdentifier(self.eat().Value)
	case lexer.Int:
		strValue := self.eat().Value
		intValue, err := strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			fmt.Println("Invalid Int '" + strValue + "'")
			os.Exit(1)
		}
		return ast.NewIntLiteral(intValue)
	case lexer.OpenParen:
		// Eat opening paren
		self.eat()
		value := self.parseExpr()
		// Eat closing paren
		self.expect(lexer.CloseParen, "Unexpexted token found inside parenthesised expression. Expected closing parenthesis.")
		return value
	default:
		fmt.Println("Unexpected token found during parsing!", self.at())
		os.Exit(1)
		return &ast.Expr{}
	}
}

func (self *Parser) at() *lexer.Token {
	return self.tokens[0]
}

func (self *Parser) expect(tokenType lexer.TokenType, errMsg string) *lexer.Token {
	prev := self.eat()
	if prev.TokenType != tokenType {
		fmt.Printf("\nParser Error: %s\nExpected: %s\nGot: %s\n", errMsg, tokenType, prev)
		os.Exit(1)
	}
	return prev
}

func (self *Parser) eat() *lexer.Token {
	var prev *lexer.Token
	prev, self.tokens = self.tokens[0], self.tokens[1:]
	return prev
}
