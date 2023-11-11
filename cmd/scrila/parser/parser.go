package parser

import (
	"ScriLa/cmd/scrila/lexer"
	"ScriLa/cmd/scrila/scrilaAst"
	"fmt"

	"golang.org/x/exp/slices"
)

var additiveOps = []string{"+", "-"}

var multiplicitaveOps = []string{"*", "/"} // TODO Modulo %

var funcReturnTypes = []lexer.TokenType{lexer.BoolType, lexer.VoidType, lexer.IntType, lexer.StrType}

type Parser struct {
	lexer    *lexer.Lexer
	tokens   []*lexer.Token
	filename string
}

func NewParser() *Parser {
	return &Parser{lexer: lexer.NewLexer()}
}

func (self *Parser) ProduceAST(sourceCode string, filename string) (scrilaAst.IProgram, error) {
	self.filename = filename
	program := scrilaAst.NewProgram()

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

func (self *Parser) parseStatement() (scrilaAst.IStatement, error) {
	var statement scrilaAst.IStatement
	var err error
	switch self.at().TokenType {
	case lexer.Comment:
		return scrilaAst.NewComment(self.eat()), nil
	case lexer.Const, lexer.BoolType, lexer.IntType, lexer.StrType, lexer.ObjType:
		statement, err = self.parseVarDeclaration()
		if err != nil {
			return scrilaAst.NewEmptyStatement(), err
		}
	case lexer.If:
		return self.parseIfStatement(false)
	case lexer.While:
		return self.parserWhileStatement()
	case lexer.Function:
		return self.parseFunctionDeclaration()
	case lexer.Break, lexer.Continue:
		statement = scrilaAst.NewIdentifier(self.eat())
	case lexer.Return:
		statement, err = self.parseReturnExpr()
		if err != nil {
			return scrilaAst.NewEmptyStatement(), err
		}
	default:
		statement, err = self.parseExpr()
		if err != nil {
			return scrilaAst.NewEmptyStatement(), err
		}
	}

	_, err = self.expect(lexer.Semicolon, "Expression must end with a semicolon")
	return statement, err
}

// [const] [int|obj] IDENT = EXPR;
func (self *Parser) parseVarDeclaration() (scrilaAst.IStatement, error) {
	isConstant := self.at().TokenType == lexer.Const
	if isConstant {
		self.eat()
	}

	if !slices.Contains([]lexer.TokenType{lexer.ObjType, lexer.StrType, lexer.IntType, lexer.BoolType}, self.at().TokenType) {
		return scrilaAst.NewEmptyStatement(), fmt.Errorf("%s: Variable type '%s' not given or supported", self.getPos(self.at()), self.at().Value)
	}
	varType := self.eat().TokenType

	// TODO Check if type matches with result of parseExpr()
	token, err := self.expect(lexer.Identifier, "Expected identifier name following [const] [int] keywords")
	if err != nil {
		return scrilaAst.NewEmptyStatement(), err
	}
	identifier := token.Value
	_, err = self.expect(lexer.Equals, "Expected equals token following identifier in var declaration")
	if err != nil {
		return scrilaAst.NewEmptyStatement(), err
	}
	expr, err := self.parseExpr()
	if err != nil {
		return scrilaAst.NewEmptyStatement(), err
	}
	declaration := scrilaAst.NewVarDeclaration(varType, isConstant, identifier, expr, token.Ln, token.Col)
	return declaration, nil
}

func (self *Parser) parseReturnExpr() (scrilaAst.IStatement, error) {
	var value scrilaAst.IExpr

	returnToken := self.eat()
	isEmpty := self.at().TokenType == lexer.Semicolon
	if isEmpty {
		value = scrilaAst.NewEmptyExpr()
	} else {
		var err error
		value, err = self.parseExpr()
		if err != nil {
			return scrilaAst.NewEmptyStatement(), err
		}
	}

	return scrilaAst.NewReturnExpr(value, isEmpty, returnToken.Ln, returnToken.Col), nil
}

func (self *Parser) parseIfStatement(isElse bool) (scrilaAst.IStatement, error) {
	ifToken := self.eat()

	var condition scrilaAst.IExpr

	if !isElse || (isElse && self.at().TokenType == lexer.If) {
		// Else if
		if isElse && self.at().TokenType == lexer.If {
			self.eat()
		}

		// Condition wrapped in braces
		_, err := self.expect(lexer.OpenParen, "Expected condition wrapped in parentheses")
		if err != nil {
			return scrilaAst.NewEmptyStatement(), err
		}

		condition, err = self.parseBooleanExpr()
		if err != nil {
			return scrilaAst.NewEmptyStatement(), err
		}

		_, err = self.expect(lexer.CloseParen, "Expected closing parenthesis after condition")
		if err != nil {
			return scrilaAst.NewEmptyStatement(), err
		}
	}

	// Body
	_, err := self.expect(lexer.OpenBrace, "Expected block following condition")
	if err != nil {
		return scrilaAst.NewEmptyStatement(), err
	}

	body := make([]scrilaAst.IStatement, 0)

	for self.notEOF() && self.at().TokenType != lexer.CloseBracket {
		statement, err := self.parseStatement()
		if err != nil {
			return scrilaAst.NewEmptyStatement(), err
		}
		body = append(body, statement)
	}

	_, err = self.expect(lexer.CloseBrace, "Closing brace expected after if block")
	if err != nil {
		return scrilaAst.NewEmptyStatement(), err
	}

	// Else
	var elseBlock scrilaAst.IIfStatement
	if self.at().TokenType == lexer.Else {
		elseBlockStmt, err := self.parseIfStatement(true)
		if err != nil {
			return scrilaAst.NewEmptyStatement(), err
		}
		elseBlock = scrilaAst.ExprToIfStmt(elseBlockStmt)
	}

	return scrilaAst.NewIfStatement(condition, body, elseBlock, ifToken.Ln, ifToken.Col), nil
}

func (self *Parser) parserWhileStatement() (scrilaAst.IStatement, error) {
	whileToken := self.eat()

	// Condition wrapped in braces
	_, err := self.expect(lexer.OpenParen, "Expected condition wrapped in parentheses")
	if err != nil {
		return scrilaAst.NewEmptyStatement(), err
	}

	condition, err := self.parseBooleanExpr()
	if err != nil {
		return scrilaAst.NewEmptyStatement(), err
	}

	_, err = self.expect(lexer.CloseParen, "Expected closing parenthesis after condition")
	if err != nil {
		return scrilaAst.NewEmptyStatement(), err
	}

	// Body
	_, err = self.expect(lexer.OpenBrace, "Expected block following condition")
	if err != nil {
		return scrilaAst.NewEmptyStatement(), err
	}

	body := make([]scrilaAst.IStatement, 0)

	for self.notEOF() && self.at().TokenType != lexer.CloseBracket {
		statement, err := self.parseStatement()
		if err != nil {
			return scrilaAst.NewEmptyStatement(), err
		}
		body = append(body, statement)
	}

	_, err = self.expect(lexer.CloseBrace, "Closing brace expected after if block")
	if err != nil {
		return scrilaAst.NewEmptyStatement(), err
	}

	return scrilaAst.NewWhileStatement(condition, body, whileToken.Ln, whileToken.Col), nil
}

func (self *Parser) parseFunctionDeclaration() (scrilaAst.IStatement, error) {
	funcToken := self.eat()

	// Function name
	token, err := self.expect(lexer.Identifier, "Expected function name following func keyword")
	if err != nil {
		return scrilaAst.NewEmptyStatement(), err
	}
	name := token.Value

	// Parameters
	params, err := self.parseParams()
	if err != nil {
		return scrilaAst.NewEmptyStatement(), err
	}

	// Return type
	returnType := self.eat()
	if returnType.TokenType == lexer.OpenBrace {
		return scrilaAst.NewEmptyStatement(), fmt.Errorf("%s: Return type is missing", self.getPos(returnType))
	}
	if !slices.Contains(funcReturnTypes, returnType.TokenType) {
		return scrilaAst.NewEmptyStatement(), fmt.Errorf("%s: Unsupported return type '%s'", self.getPos(returnType), returnType.Value)
	}

	// Body
	_, err = self.expect(lexer.OpenBrace, "Expected function body following declaration")
	if err != nil {
		return scrilaAst.NewEmptyStatement(), err
	}

	body := make([]scrilaAst.IStatement, 0)

	for self.notEOF() && self.at().TokenType != lexer.CloseBracket {
		statement, err := self.parseStatement()
		if err != nil {
			return scrilaAst.NewEmptyStatement(), err
		}
		body = append(body, statement)
	}

	_, err = self.expect(lexer.CloseBrace, "Closing brace expected inside function declaration")
	return scrilaAst.NewFunctionDeclaration(name, params, body, returnType.TokenType, funcToken.Ln, funcToken.Col), err
}

func (self *Parser) parseExpr() (scrilaAst.IExpr, error) {
	return self.parseAssignmentExpr()
}

// Orders of Precedence:
// Processed from top to bottom
// Priority from bottom to top
// - AssignmentExpr
// - ObjectExpr
// - BooleanExpr
// - ComparisonExpr
// - AdditiveExpr
// - MultiplicitiveExpr
// - CallExpr
// - MemberExr
// - PrimaryExpr

// - LogicalExpr
// - Comparison
// - UnaryExpr

func (self *Parser) parseAssignmentExpr() (scrilaAst.IExpr, error) {
	left, err := self.parseObjectExpr()
	if err != nil {
		return scrilaAst.NewEmptyExpr(), err
	}

	if self.at().TokenType == lexer.Equals {
		self.eat() // Advance past equals

		value, err := self.parseAssignmentExpr() // This allows chaining e.g. x = y = 5; TODO Do not allow chaining
		if err != nil {
			return scrilaAst.NewEmptyExpr(), err
		}
		return scrilaAst.NewAssignmentExpr(left, value), nil
	}

	return left, nil
}

func (self *Parser) parseObjectExpr() (scrilaAst.IExpr, error) {
	// { Prop[] }

	if self.at().TokenType != lexer.OpenBrace {
		return self.parseBooleanExpr()
	}
	self.eat() // Advance past open brace

	properties := make([]scrilaAst.IProperty, 0)
	for self.notEOF() && self.at().TokenType != lexer.CloseBrace {
		// { key: val, }

		token, err := self.expect(lexer.Identifier, "Object literal key expected")
		if err != nil {
			return scrilaAst.NewEmptyExpr(), err
		}
		_, err = self.expect(lexer.Colon, "Missing colon following identifier in ObjectExpr")
		if err != nil {
			return scrilaAst.NewEmptyExpr(), err
		}
		value, err := self.parseExpr()
		if err != nil {
			return scrilaAst.NewEmptyExpr(), err
		}
		_, err = self.expect(lexer.Comma, "Expected comma following Property")
		if err != nil {
			return scrilaAst.NewEmptyExpr(), err
		}

		properties = append(properties, scrilaAst.NewProperty(token.Value, value, token.Ln, token.Col))
	}

	_, err := self.expect(lexer.CloseBrace, "Object literal missing closing brace")
	return scrilaAst.NewObjectLiteral(properties), err
}

func (self *Parser) parseBooleanExpr() (scrilaAst.IExpr, error) {
	left, err := self.parseComparisonExpr()
	if err != nil {
		return scrilaAst.NewEmptyExpr(), err
	}

	// Current token is an boolean operator
	for slices.Contains(lexer.BooleanOps, self.at().Value) {
		token := self.eat()
		right, err := self.parseComparisonExpr()
		if err != nil {
			return scrilaAst.NewEmptyExpr(), err
		}
		left = scrilaAst.NewBinaryExpr(left, right, token.Value, token.Ln, token.Col)
	}

	return left, nil
}

func (self *Parser) parseComparisonExpr() (scrilaAst.IExpr, error) {
	left, err := self.parseAdditiveExpr()
	if err != nil {
		return scrilaAst.NewEmptyExpr(), err
	}

	// Current token is an comparison operator
	for slices.Contains(lexer.ComparisonOps, self.at().Value) {
		token := self.eat()
		right, err := self.parseAdditiveExpr()
		if err != nil {
			return scrilaAst.NewEmptyExpr(), err
		}
		left = scrilaAst.NewBinaryExpr(left, right, token.Value, token.Ln, token.Col)
	}

	return left, nil
}

func (self *Parser) parseAdditiveExpr() (scrilaAst.IExpr, error) {
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
		return scrilaAst.NewEmptyExpr(), err
	}

	// Current token is an additive operator
	for slices.Contains(additiveOps, self.at().Value) {
		token := self.eat()
		right, err := self.parseMultiplicitaveExpr()
		if err != nil {
			return scrilaAst.NewEmptyExpr(), err
		}
		left = scrilaAst.NewBinaryExpr(left, right, token.Value, token.Ln, token.Col)
	}

	return left, nil
}

func (self *Parser) parseMultiplicitaveExpr() (scrilaAst.IExpr, error) {
	// Lefthand Precedence (see func parseAdditiveExpr)
	left, err := self.parseCallMemberExpr()
	if err != nil {
		return scrilaAst.NewEmptyExpr(), err
	}

	// Current token is a multiplicitave operator
	for slices.Contains(multiplicitaveOps, self.at().Value) {
		token := self.eat()
		right, err := self.parseCallMemberExpr()
		if err != nil {
			return scrilaAst.NewEmptyExpr(), err
		}
		left = scrilaAst.NewBinaryExpr(left, right, token.Value, token.Ln, token.Col)
	}

	return left, nil
}

// foo.x()
func (self *Parser) parseCallMemberExpr() (scrilaAst.IExpr, error) {
	member, err := self.parseMemberExpr()
	if err != nil {
		return scrilaAst.NewEmptyExpr(), err
	}

	if self.at().TokenType == lexer.OpenParen {
		return self.parseCallExpr(member)
	}

	return member, nil
}

// foo()
func (self *Parser) parseCallExpr(caller scrilaAst.IExpr) (scrilaAst.IExpr, error) {
	var callExpr scrilaAst.IExpr
	args, err := self.parseArgs()
	if err != nil {
		return scrilaAst.NewEmptyExpr(), err
	}
	callExpr = scrilaAst.NewCallExpr(caller, args)

	// This allows chaining of function calls: e.g. foo()()
	// TODO Do not allow chaining
	if self.at().TokenType == lexer.OpenParen {
		var err error
		callExpr, err = self.parseCallExpr(callExpr)
		if err != nil {
			return scrilaAst.NewEmptyExpr(), err
		}
	}

	return callExpr, nil
}

func (self *Parser) parseParams() ([]*scrilaAst.Parameter, error) {
	var args []*scrilaAst.Parameter
	_, err := self.expect(lexer.OpenParen, "Expected open parenthesis")
	if err != nil {
		return args, err
	}
	if self.at().TokenType == lexer.CloseParen {
		args = make([]*scrilaAst.Parameter, 0)
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

func (self *Parser) parseParametersList() ([]*scrilaAst.Parameter, error) {
	params := make([]*scrilaAst.Parameter, 0)

	if !slices.Contains([]lexer.TokenType{lexer.StrType, lexer.BoolType, lexer.IntType, lexer.ObjType}, self.at().TokenType) {
		return params, fmt.Errorf("%s: Expected param type but got %s '%s'", self.getPos(self.at()), self.at().TokenType, self.at().Value)
	}

	for self.notEOF() && slices.Contains([]lexer.TokenType{lexer.StrType, lexer.BoolType, lexer.IntType, lexer.ObjType}, self.at().TokenType) {
		paramType := self.eat().TokenType
		ident, err := self.expect(lexer.Identifier, "parseParametersList: Expected identifier following param type")
		if err != nil {
			return params, err
		}
		params = append(params, scrilaAst.NewParameter(ident.Value, paramType))

		if self.at().TokenType == lexer.Comma {
			self.eat()
		} else if self.at().TokenType == lexer.CloseParen {
			return params, nil
		}
	}
	return params, fmt.Errorf("%s: Unexpected token '%s' in parameter list", self.getPos(self.at()), self.at().Value)
}

// func add(a, b) {} <- a & b are parameters
// add(a, b) <- a & b are now args (when calling)
func (self *Parser) parseArgs() ([]scrilaAst.IExpr, error) {
	var args []scrilaAst.IExpr
	_, err := self.expect(lexer.OpenParen, "Expected open parenthesis")
	if err != nil {
		return args, err
	}
	if self.at().TokenType == lexer.CloseParen {
		args = make([]scrilaAst.IExpr, 0)
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
func (self *Parser) parseArgumentsList() ([]scrilaAst.IExpr, error) {
	args := []scrilaAst.IExpr{}
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

func (self *Parser) parseMemberExpr() (scrilaAst.IExpr, error) {
	object, err := self.parsePrimaryExpr()
	if err != nil {
		return scrilaAst.NewEmptyExpr(), err
	}

	for self.at().TokenType == lexer.Dot || self.at().TokenType == lexer.OpenBracket {
		operator := self.eat()
		var property scrilaAst.IExpr
		var isComputed bool

		// Non-computed values aka "obj.expr"
		if operator.TokenType == lexer.Dot {
			isComputed = false
			// Get identifier
			property, err = self.parsePrimaryExpr()
			if err != nil {
				return scrilaAst.NewEmptyExpr(), err
			}

			if property.GetKind() != scrilaAst.IdentifierNode {
				return scrilaAst.NewEmptyExpr(), fmt.Errorf("%s: Cannot use dot operator without right hand side being an identifier", self.getPosExpr(property))
			}
		} else {
			isComputed = true
			// This allows chaining: obj[computedValue] e.g. obj1[obj2[getBar()]]
			property, err = self.parseExpr()
			if err != nil {
				return scrilaAst.NewEmptyExpr(), err
			}

			_, err = self.expect(lexer.CloseBracket, "Missing closing bracket in computed value")
			if err != nil {
				return scrilaAst.NewEmptyExpr(), err
			}
		}

		object = scrilaAst.NewMemberExpr(object, property, isComputed)
	}

	return object, nil
}

func (self *Parser) parsePrimaryExpr() (scrilaAst.IExpr, error) {
	switch self.at().TokenType {
	case lexer.Identifier:
		return scrilaAst.NewIdentifier(self.eat()), nil
	case lexer.Int:
		return scrilaAst.NewIntLiteral(self.eat())
	case lexer.Str:
		return scrilaAst.NewStrLiteral(self.eat()), nil
	case lexer.OpenParen:
		// Eat opening paren
		self.eat()
		value, err := self.parseExpr()
		if err != nil {
			return scrilaAst.NewEmptyExpr(), nil
		}
		// Eat closing paren
		_, err = self.expect(lexer.CloseParen, "Unexpexted token found inside parenthesised expression. Expected closing parenthesis")
		return value, err
	default:
		return scrilaAst.NewEmptyExpr(), fmt.Errorf("%s: Unexpected token '%s' ('%s') found during parsing", self.getPos(self.at()), self.at().TokenType, self.at().Value)
	}
}

func (self *Parser) at() *lexer.Token {
	return self.tokens[0]
}

func (self *Parser) expect(tokenType lexer.TokenType, errMsg string) (*lexer.Token, error) {
	prev := self.eat()
	if prev.TokenType != tokenType {
		return &lexer.Token{}, fmt.Errorf("%s: %s\nExpected: %s\nGot: %s", self.getPos(prev), errMsg, tokenType, prev)
	}
	return prev, nil
}

func (self *Parser) eat() *lexer.Token {
	var prev *lexer.Token
	prev, self.tokens = self.tokens[0], self.tokens[1:]
	return prev
}

func (self *Parser) getPos(token *lexer.Token) string {
	return fmt.Sprintf("%s:%d:%d", self.filename, token.Ln, token.Col)
}

func (self *Parser) getPosExpr(expr scrilaAst.IExpr) string {
	return fmt.Sprintf("%s:%d:%d", self.filename, expr.GetLn(), expr.GetCol())
}
