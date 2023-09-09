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
	tokens []lexer.Token
}

func NewParser() *Parser {
	return &Parser{}
}

func (self Parser) ProduceAST(sourceCode string) ast.IProgram {
	self.tokens = lexer.Tokenize(sourceCode)
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
	// Skip for now
	return self.parseExpr()
}

func (self *Parser) parseExpr() ast.IExpr {
	return self.parseAdditiveExpr()
}

// Orders of Prescidence:
// Processed from top to bottom
// Priority from bottom to top
// - AssignmentExpr
// - MemberExr
// - FunctionCall
// - LogicalExpr
// - Comparison
// - AdditiveExpr
// - MultiplicitiveExpr
// - UnaryExpr
// - PrimaryExpr

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
	left := self.parsePrimaryExpr()

	// Current token is a multiplicitave operator
	for slices.Contains(multiplicitaveOps, self.at().Value) {
		operator := self.eat().Value
		right := self.parsePrimaryExpr()
		left = ast.NewBinaryExpr(left, right, operator)
	}

	return left
}

func (self *Parser) parsePrimaryExpr() ast.IExpr {
	switch self.at().TokenType {
	case lexer.Identifier:
		return ast.NewIdentifier(self.eat().Value)
	case lexer.NullType:
		self.eat() // Advance post null keyword
		return ast.NewNullLiteral()
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

func (self *Parser) at() lexer.Token {
	return self.tokens[0]
}

func (self *Parser) expect(tokenType lexer.TokenType, errMsg string) lexer.Token {
	var prev lexer.Token
	prev, self.tokens = self.tokens[0], self.tokens[1:]
	if prev.TokenType != tokenType {
		fmt.Printf("\nParser Error: %s\nExpected: %s\n", errMsg, prev)
		os.Exit(1)
	}
	return prev
}

func (self *Parser) eat() lexer.Token {
	var prev lexer.Token
	prev, self.tokens = self.tokens[0], self.tokens[1:]
	return prev
}
