package parser

import (
	"ScriLa/cmd/scrila/ast"
	"ScriLa/cmd/scrila/lexer"
	"fmt"
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

func (self Parser) ProduceAST(sourceCode string) (ast.IProgram, error) {
	var err error
	self.tokens, err = self.lexer.Tokenize(sourceCode)
	if err != nil {
		return ast.NewProgram(), err
	}
	program := ast.NewProgram()

	for self.notEOF() {
		statement, err := self.parseStatement()
		if err != nil {
			return program, err
		}
		program.Body = append(program.GetBody(), statement)
	}

	return program, nil
}

func (self *Parser) notEOF() bool {
	if len(self.tokens) == 0 {
		return false
	}
	return self.tokens[0].TokenType != lexer.EndOfFile
}

func (self *Parser) parseStatement() (ast.IStatement, error) {
	var statement ast.IStatement
	var err error
	switch self.at().TokenType {
	case lexer.Const, lexer.BoolType, lexer.IntType, lexer.StrType, lexer.ObjType:
		statement, err = self.parseVarDeclaration()
		if err != nil {
			return ast.NewStatement(), err
		}
	case lexer.Function:
		return self.parseFunctionDeclaration()
	default:
		statement, err = self.parseExpr()
		if err != nil {
			return ast.NewStatement(), err
		}
	}

	_, err = self.expect(lexer.Semicolon, "Expressions must end with a semicolon.")
	return statement, err
}

// [const] [int|obj] IDENT = EXPR;
func (self *Parser) parseVarDeclaration() (ast.IStatement, error) {
	isConstant := self.at().TokenType == lexer.Const
	if isConstant {
		self.eat()
	}

	if !slices.Contains([]lexer.TokenType{lexer.ObjType, lexer.StrType, lexer.IntType, lexer.BoolType}, self.at().TokenType) {
		return ast.NewStatement(), fmt.Errorf("Variable type not given or supported. %s", self.at())
	}
	self.eat()

	// TODO Check if type matches with result of parseExpr()
	token, err := self.expect(lexer.Identifier, "Expected identifier name following [const] [int] keywords.")
	if err != nil {
		return ast.NewStatement(), err
	}
	identifier := token.Value
	_, err = self.expect(lexer.Equals, "Expected equals token following identifier in var declaration.")
	if err != nil {
		return ast.NewStatement(), err
	}
	expr, err := self.parseExpr()
	if err != nil {
		return ast.NewStatement(), err
	}
	declaration := ast.NewVarDeclaration(isConstant, identifier, expr)
	return declaration, nil
}

func (self *Parser) parseFunctionDeclaration() (ast.IStatement, error) {
	self.eat()
	token, err := self.expect(lexer.Identifier, "Expected function name following func keyword.")
	if err != nil {
		return ast.NewStatement(), err
	}
	name := token.Value
	args, err := self.parseArgs()
	if err != nil {
		return ast.NewStatement(), err
	}
	params := make([]string, 0)
	for _, arg := range args {
		if arg.GetKind() != ast.IdentifierNode {
			return ast.NewStatement(), fmt.Errorf("Inside function declaration expected parameters to be of type string.")
		}

		var i interface{} = arg
		identifier, _ := i.(ast.IIdentifier)
		params = append(params, identifier.GetSymbol())
	}
	_, err = self.expect(lexer.OpenBrace, "Expected function body following declaration.")
	if err != nil {
		return ast.NewStatement(), err
	}

	body := make([]ast.IStatement, 0)

	for self.notEOF() && self.at().TokenType != lexer.CloseBracket {
		statement, err := self.parseStatement()
		if err != nil {
			return ast.NewStatement(), err
		}
		body = append(body, statement)
	}

	_, err = self.expect(lexer.CloseBrace, "Closing brace expected inside function declaration.")
	return ast.NewFunctionDeclaration(name, params, body), err
}

func (self *Parser) parseExpr() (ast.IExpr, error) {
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

func (self *Parser) parseAssignmentExpr() (ast.IExpr, error) {
	left, err := self.parseObjectExpr()
	if err != nil {
		return ast.NewExpr(), err
	}

	if self.at().TokenType == lexer.Equals {
		self.eat() // Advance past equals

		value, err := self.parseAssignmentExpr() // This allows chaining e.g. x = y = 5; TODO Do not allow chaining
		if err != nil {
			return ast.NewExpr(), err
		}
		return ast.NewAssignmentExpr(left, value), nil
	}

	return left, nil
}

func (self *Parser) parseObjectExpr() (ast.IExpr, error) {
	// { Prop[] }

	if self.at().TokenType != lexer.OpenBrace {
		return self.parseAdditiveExpr()
	}
	self.eat() // Advance past open brace

	properties := make([]ast.IProperty, 0)
	for self.notEOF() && self.at().TokenType != lexer.CloseBrace {
		// { key: val, }

		token, err := self.expect(lexer.Identifier, "Object literal key expected.")
		if err != nil {
			return ast.NewExpr(), err
		}
		key := token.Value
		_, err = self.expect(lexer.Colon, "Missing colon following identifier in ObjectExpr.")
		if err != nil {
			return ast.NewExpr(), err
		}
		value, err := self.parseExpr()
		if err != nil {
			return ast.NewExpr(), err
		}
		_, err = self.expect(lexer.Comma, "Expected comma following Property.")
		if err != nil {
			return ast.NewExpr(), err
		}

		properties = append(properties, ast.NewProperty(key, value))
	}

	_, err := self.expect(lexer.CloseBrace, "Object literal missing closing brace.")
	return ast.NewObjectLiteral(properties), err
}

func (self *Parser) parseAdditiveExpr() (ast.IExpr, error) {
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

	left, err := self.parseMultiplicitaveExpr()
	if err != nil {
		return ast.NewExpr(), err
	}

	// Current token is an additive operator
	for slices.Contains(additiveOps, self.at().Value) {
		operator := self.eat().Value
		right, err := self.parseMultiplicitaveExpr()
		if err != nil {
			return ast.NewExpr(), err
		}
		left = ast.NewBinaryExpr(left, right, operator)
	}

	return left, nil
}

func (self *Parser) parseMultiplicitaveExpr() (ast.IExpr, error) {
	// Lefthand Prescidence (see func parseAdditiveExpr)
	left, err := self.parseCallMemberExpr()
	if err != nil {
		return ast.NewExpr(), err
	}

	// Current token is a multiplicitave operator
	for slices.Contains(multiplicitaveOps, self.at().Value) {
		operator := self.eat().Value
		right, err := self.parseCallMemberExpr()
		if err != nil {
			return ast.NewExpr(), err
		}
		left = ast.NewBinaryExpr(left, right, operator)
	}

	return left, nil
}

// foo.x()
func (self *Parser) parseCallMemberExpr() (ast.IExpr, error) {
	member, err := self.parseMemberExpr()
	if err != nil {
		return ast.NewExpr(), err
	}

	if self.at().TokenType == lexer.OpenParen {
		return self.parseCallExpr(member)
	}

	return member, nil
}

// foo()
func (self *Parser) parseCallExpr(caller ast.IExpr) (ast.IExpr, error) {
	var callExpr ast.IExpr
	args, err := self.parseArgs()
	if err != nil {
		return ast.NewExpr(), err
	}
	callExpr = ast.NewCallExpr(caller, args)

	// This allows chaining of function calls: e.g. foo()()
	if self.at().TokenType == lexer.OpenParen {
		var err error
		callExpr, err = self.parseCallExpr(callExpr)
		if err != nil {
			return ast.NewExpr(), err
		}
	}

	return callExpr, nil
}

// func add(a, b) {} <- a & b are parameters
// add(a, b) <- a & b are now args (when calling)
func (self *Parser) parseArgs() ([]ast.IExpr, error) {
	var args []ast.IExpr
	_, err := self.expect(lexer.OpenParen, "Expected open parenthesis")
	if err != nil {
		return args, err
	}
	if self.at().TokenType == lexer.CloseParen {
		args = make([]ast.IExpr, 0)
	} else {
		var err error
		args, err = self.parseArgumentsList()
		if err != nil {
			return args, err
		}
	}

	_, err = self.expect(lexer.CloseParen, "Missing closing parenthesis inside arguments list")
	return args, err
}

// foo(x = 5, v = "Bar")
// Set x to 5 and v to "Bar" in global scope and pass values afterwards
func (self *Parser) parseArgumentsList() ([]ast.IExpr, error) {
	args := []ast.IExpr{}
	expr, err := self.parseAssignmentExpr()
	if err != nil {
		return args, err
	}
	args = append(args, expr)

	for self.notEOF() && self.at().TokenType == lexer.Comma {
		self.eat()
		expr, err := self.parseAssignmentExpr()
		if err != nil {
			return args, err
		}
		args = append(args, expr)
	}

	return args, nil
}

func (self *Parser) parseMemberExpr() (ast.IExpr, error) {
	object, err := self.parsePrimaryExpr()
	if err != nil {
		return ast.NewExpr(), err
	}

	for self.at().TokenType == lexer.Dot || self.at().TokenType == lexer.OpenBracket {
		operator := self.eat()
		var property ast.IExpr
		var isComputed bool

		// Non-computed values aka "obj.expr"
		if operator.TokenType == lexer.Dot {
			isComputed = false
			// Get identifier
			property, err = self.parsePrimaryExpr()
			if err != nil {
				return ast.NewExpr(), err
			}

			if property.GetKind() != ast.IdentifierNode {
				return ast.NewExpr(), fmt.Errorf("Cannot use dot operator without right hand side being an identifier")
			}
		} else {
			isComputed = true
			// This allows chaining: obj[computedValue] e.g. obj1[obj2[getBar()]]
			property, err = self.parseExpr()
			if err != nil {
				return ast.NewExpr(), err
			}

			_, err = self.expect(lexer.CloseBracket, "Missing closing bracket in computed value.")
			if err != nil {
				return ast.NewExpr(), err
			}
		}

		object = ast.NewMemberExpr(object, property, isComputed)
	}

	return object, nil
}

func (self *Parser) parsePrimaryExpr() (ast.IExpr, error) {
	switch self.at().TokenType {
	case lexer.Identifier:
		return ast.NewIdentifier(self.eat().Value), nil
	case lexer.Int:
		strValue := self.eat().Value
		intValue, err := strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			return ast.NewExpr(), fmt.Errorf("Invalid Int '%s'", strValue)
		}
		return ast.NewIntLiteral(intValue), nil
	case lexer.Str:
		return ast.NewStrLiteral(self.eat().Value), nil
	case lexer.OpenParen:
		// Eat opening paren
		self.eat()
		value, err := self.parseExpr()
		if err != nil {
			return ast.NewExpr(), nil
		}
		// Eat closing paren
		_, err = self.expect(lexer.CloseParen, "Unexpexted token found inside parenthesised expression. Expected closing parenthesis.")
		return value, err
	default:
		return ast.NewExpr(), fmt.Errorf("Unexpected token '%s' ('%s') (Ln %d, Col %d) found during parsing\n", self.at().TokenType, self.at().Value, self.at().Ln, self.at().Col)
	}
}

func (self *Parser) at() *lexer.Token {
	return self.tokens[0]
}

func (self *Parser) expect(tokenType lexer.TokenType, errMsg string) (*lexer.Token, error) {
	prev := self.eat()
	if prev.TokenType != tokenType {
		return &lexer.Token{}, fmt.Errorf("Parser Error: %s\nExpected: %s\nGot: %s\n", errMsg, tokenType, prev)
	}
	return prev, nil
}

func (self *Parser) eat() *lexer.Token {
	var prev *lexer.Token
	prev, self.tokens = self.tokens[0], self.tokens[1:]
	return prev
}
