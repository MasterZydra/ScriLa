package parser

import (
	"ScriLa/cmd/scrila/ast"
	"ScriLa/cmd/scrila/lexer"
	"fmt"

	"golang.org/x/exp/slices"
)

var additiveOps = []string{"+", "-"}

var multiplicitaveOps = []string{"*", "/"} // TODO Modulo %

var funcReturnTypes = []lexer.TokenType{lexer.BoolType, lexer.VoidType, lexer.IntType, lexer.StrType}

type Parser struct {
	lexer  *lexer.Lexer
	tokens []*lexer.Token
}

func NewParser() *Parser {
	return &Parser{lexer: lexer.NewLexer()}
}

var fileName string

func (self Parser) ProduceAST(sourceCode string, filename string) (ast.IProgram, error) {
	fileName = filename
	program := ast.NewProgram()

	var err error
	self.tokens, err = self.lexer.Tokenize(sourceCode, filename)
	if err != nil {
		return program, err
	}

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
	case lexer.Comment:
		return ast.NewComment(self.eat()), nil
	case lexer.Const, lexer.BoolType, lexer.IntType, lexer.StrType, lexer.ObjType:
		statement, err = self.parseVarDeclaration()
		if err != nil {
			return ast.NewEmptyStatement(), err
		}
	case lexer.Function:
		return self.parseFunctionDeclaration()
	case lexer.Return:
		self.eat()
		value, err := self.parseExpr()
		if err != nil {
			return ast.NewEmptyStatement(), err
		}
		statement = ast.NewReturnExpr(value)
	default:
		statement, err = self.parseExpr()
		if err != nil {
			return ast.NewEmptyStatement(), err
		}
	}

	_, err = self.expect(lexer.Semicolon, "Expression must end with a semicolon")
	return statement, err
}

// [const] [int|obj] IDENT = EXPR;
func (self *Parser) parseVarDeclaration() (ast.IStatement, error) {
	isConstant := self.at().TokenType == lexer.Const
	if isConstant {
		self.eat()
	}

	if !slices.Contains([]lexer.TokenType{lexer.ObjType, lexer.StrType, lexer.IntType, lexer.BoolType}, self.at().TokenType) {
		return ast.NewEmptyStatement(), fmt.Errorf("Variable type not given or supported. %s", self.at())
	}
	varType := self.eat().TokenType

	// TODO Check if type matches with result of parseExpr()
	token, err := self.expect(lexer.Identifier, "Expected identifier name following [const] [int] keywords.")
	if err != nil {
		return ast.NewEmptyStatement(), err
	}
	identifier := token.Value
	_, err = self.expect(lexer.Equals, "Expected equals token following identifier in var declaration.")
	if err != nil {
		return ast.NewEmptyStatement(), err
	}
	expr, err := self.parseExpr()
	if err != nil {
		return ast.NewEmptyStatement(), err
	}
	declaration := ast.NewVarDeclaration(varType, isConstant, identifier, expr)
	return declaration, nil
}

func (self *Parser) parseFunctionDeclaration() (ast.IStatement, error) {
	funcToken := self.eat()

	// Function name
	token, err := self.expect(lexer.Identifier, "Expected function name following func keyword.")
	if err != nil {
		return ast.NewEmptyStatement(), err
	}
	name := token.Value

	// Parameters
	params, err := self.parseParams()
	if err != nil {
		return ast.NewEmptyStatement(), err
	}

	// Return type
	returnType := self.eat()
	if returnType.TokenType == lexer.OpenBrace {
		return ast.NewEmptyStatement(), fmt.Errorf("%s:%d:%d: Return type is missing", fileName, returnType.Ln, returnType.Col)
	}
	if !slices.Contains(funcReturnTypes, returnType.TokenType) {
		return ast.NewEmptyStatement(), fmt.Errorf("%s:%d:%d: Unsupported return type '%s'", fileName, returnType.Ln, returnType.Col, returnType.Value)
	}

	// Body
	_, err = self.expect(lexer.OpenBrace, "Expected function body following declaration.")
	if err != nil {
		return ast.NewEmptyStatement(), err
	}

	body := make([]ast.IStatement, 0)

	for self.notEOF() && self.at().TokenType != lexer.CloseBracket {
		statement, err := self.parseStatement()
		if err != nil {
			return ast.NewEmptyStatement(), err
		}
		body = append(body, statement)
	}

	_, err = self.expect(lexer.CloseBrace, "Closing brace expected inside function declaration.")
	return ast.NewFunctionDeclaration(name, params, body, returnType.TokenType, funcToken.Ln, funcToken.Col), err
}

func (self *Parser) parseExpr() (ast.IExpr, error) {
	return self.parseAssignmentExpr()
}

// Orders of Precedence:
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
		return ast.NewEmptyExpr(), err
	}

	if self.at().TokenType == lexer.Equals {
		self.eat() // Advance past equals

		value, err := self.parseAssignmentExpr() // This allows chaining e.g. x = y = 5; TODO Do not allow chaining
		if err != nil {
			return ast.NewEmptyExpr(), err
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
			return ast.NewEmptyExpr(), err
		}
		_, err = self.expect(lexer.Colon, "Missing colon following identifier in ObjectExpr.")
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		value, err := self.parseExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		_, err = self.expect(lexer.Comma, "Expected comma following Property.")
		if err != nil {
			return ast.NewEmptyExpr(), err
		}

		properties = append(properties, ast.NewProperty(token.Value, value, token.Ln, token.Col))
	}

	_, err := self.expect(lexer.CloseBrace, "Object literal missing closing brace.")
	return ast.NewObjectLiteral(properties), err
}

func (self *Parser) parseAdditiveExpr() (ast.IExpr, error) {
	// Lefthand Precedence
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
		return ast.NewEmptyExpr(), err
	}

	// Current token is an additive operator
	for slices.Contains(additiveOps, self.at().Value) {
		token := self.eat()
		right, err := self.parseMultiplicitaveExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		left = ast.NewBinaryExpr(left, right, token.Value, token.Ln, token.Col)
	}

	return left, nil
}

func (self *Parser) parseMultiplicitaveExpr() (ast.IExpr, error) {
	// Lefthand Precedence (see func parseAdditiveExpr)
	left, err := self.parseCallMemberExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
	}

	// Current token is a multiplicitave operator
	for slices.Contains(multiplicitaveOps, self.at().Value) {
		token := self.eat()
		right, err := self.parseCallMemberExpr()
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
		left = ast.NewBinaryExpr(left, right, token.Value, token.Ln, token.Col)
	}

	return left, nil
}

// foo.x()
func (self *Parser) parseCallMemberExpr() (ast.IExpr, error) {
	member, err := self.parseMemberExpr()
	if err != nil {
		return ast.NewEmptyExpr(), err
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
		return ast.NewEmptyExpr(), err
	}
	callExpr = ast.NewCallExpr(caller, args)

	// This allows chaining of function calls: e.g. foo()()
	// TODO Do not allow chaining
	if self.at().TokenType == lexer.OpenParen {
		var err error
		callExpr, err = self.parseCallExpr(callExpr)
		if err != nil {
			return ast.NewEmptyExpr(), err
		}
	}

	return callExpr, nil
}

func (self *Parser) parseParams() ([]*ast.Parameter, error) {
	var args []*ast.Parameter
	_, err := self.expect(lexer.OpenParen, "Expected open parenthesis")
	if err != nil {
		return args, err
	}
	if self.at().TokenType == lexer.CloseParen {
		args = make([]*ast.Parameter, 0)
	} else {
		var err error
		args, err = self.parseParametersList()
		if err != nil {
			return args, err
		}
	}

	_, err = self.expect(lexer.CloseParen, "Missing closing parenthesis inside arguments list")
	return args, err
}

func (self *Parser) parseParametersList() ([]*ast.Parameter, error) {
	params := make([]*ast.Parameter, 0)

	if !slices.Contains([]lexer.TokenType{lexer.StrType, lexer.BoolType, lexer.IntType, lexer.ObjType}, self.at().TokenType) {
		return params, fmt.Errorf("parseParametersList: Expected param type but got: %s", self.at())
	}

	for self.notEOF() && slices.Contains([]lexer.TokenType{lexer.StrType, lexer.BoolType, lexer.IntType, lexer.ObjType}, self.at().TokenType) {
		paramType := self.eat().TokenType
		ident, err := self.expect(lexer.Identifier, "parseParametersList: Expected identifier following param type.")
		if err != nil {
			return params, err
		}
		params = append(params, ast.NewParameter(ident.Value, paramType))

		if self.at().TokenType == lexer.Comma {
			self.eat()
		} else if self.at().TokenType == lexer.CloseParen {
			return params, nil
		} else {
			return params, fmt.Errorf("parseParametersList: Unexpected token: %s", self.at())
		}
	}
	return params, fmt.Errorf("parseParametersList: Unexpected token: %s", self.at())
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
		return ast.NewEmptyExpr(), err
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
				return ast.NewEmptyExpr(), err
			}

			if property.GetKind() != ast.IdentifierNode {
				return ast.NewEmptyExpr(), fmt.Errorf("Cannot use dot operator without right hand side being an identifier")
			}
		} else {
			isComputed = true
			// This allows chaining: obj[computedValue] e.g. obj1[obj2[getBar()]]
			property, err = self.parseExpr()
			if err != nil {
				return ast.NewEmptyExpr(), err
			}

			_, err = self.expect(lexer.CloseBracket, "Missing closing bracket in computed value.")
			if err != nil {
				return ast.NewEmptyExpr(), err
			}
		}

		object = ast.NewMemberExpr(object, property, isComputed)
	}

	return object, nil
}

func (self *Parser) parsePrimaryExpr() (ast.IExpr, error) {
	switch self.at().TokenType {
	case lexer.Identifier:
		return ast.NewIdentifier(self.eat()), nil
	case lexer.Int:
		return ast.NewIntLiteral(self.eat())
	case lexer.Str:
		return ast.NewStrLiteral(self.eat()), nil
	case lexer.OpenParen:
		// Eat opening paren
		self.eat()
		value, err := self.parseExpr()
		if err != nil {
			return ast.NewEmptyExpr(), nil
		}
		// Eat closing paren
		_, err = self.expect(lexer.CloseParen, "Unexpexted token found inside parenthesised expression. Expected closing parenthesis.")
		return value, err
	default:
		return ast.NewEmptyExpr(), fmt.Errorf("%s:%d:%d: Unexpected token '%s' ('%s') found during parsing", fileName, self.at().Ln, self.at().Col, self.at().TokenType, self.at().Value)
	}
}

func (self *Parser) at() *lexer.Token {
	return self.tokens[0]
}

func (self *Parser) expect(tokenType lexer.TokenType, errMsg string) (*lexer.Token, error) {
	prev := self.eat()
	if prev.TokenType != tokenType {
		return &lexer.Token{}, fmt.Errorf("%s:%d:%d: %s\nExpected: %s\nGot: %s", fileName, prev.Ln, prev.Col, errMsg, tokenType, prev)
	}
	return prev, nil
}

func (self *Parser) eat() *lexer.Token {
	var prev *lexer.Token
	prev, self.tokens = self.tokens[0], self.tokens[1:]
	return prev
}
