package ast

type NodeType string

const (
	// Statements
	StatementNode           NodeType = "Statement"
	CommentNode             NodeType = "Comment"
	ProgramNode             NodeType = "Program"
	VarDeclarationNode      NodeType = "VarDeclaration"
	FunctionDeclarationNode NodeType = "FunctionDeclaration"

	// Expressions
	ExprNode           NodeType = "Expr"
	AssignmentExprNode NodeType = "AssignmentExpr"
	BinaryExprNode     NodeType = "BinaryExpr"
	UnaryExprNode      NodeType = "UnaryExpr"
	CallExprNode       NodeType = "CallExpr"
	MemberExprNode     NodeType = "MemberExpr"
	ReturnExprNode     NodeType = "ReturnExpr"

	// Literals
	PropertyNode      NodeType = "Property"
	ObjectLiteralNode NodeType = "ObjectLiteral"
	IdentifierNode    NodeType = "Identifier"
	IntLiteralNode    NodeType = "IntLiteral"
	StrLiteralNode    NodeType = "StrLiteral"
)
