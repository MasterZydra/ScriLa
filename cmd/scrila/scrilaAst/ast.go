package scrilaAst

type NodeType string

const (
	// Statements
	StatementNode           NodeType = "Statement"
	CommentNode             NodeType = "Comment"
	ProgramNode             NodeType = "Program"
	VarDeclarationNode      NodeType = "VarDeclaration"
	FunctionDeclarationNode NodeType = "FunctionDeclaration"
	IfStatementNode         NodeType = "IfStmt"
	WhileStatementNode      NodeType = "WhileLoop"

	// Expressions
	ExprNode           NodeType = "Expr"
	AssignmentExprNode NodeType = "AssignmentExpr"
	BinaryExprNode     NodeType = "BinaryExpr"
	UnaryExprNode      NodeType = "UnaryExpr"
	CallExprNode       NodeType = "CallExpr"
	MemberExprNode     NodeType = "MemberExpr"
	ReturnExprNode     NodeType = "ReturnExpr"
	BreakExprNode      NodeType = "BreakExpr"
	ContinueExprNode   NodeType = "ContinueExpr"

	// Literals
	PropertyNode      NodeType = "Property"
	ObjectLiteralNode NodeType = "ObjectLiteral"
	IdentifierNode    NodeType = "Identifier"
	IntLiteralNode    NodeType = "IntLiteral"
	StrLiteralNode    NodeType = "StrLiteral"
	BoolLiteralNode   NodeType = "BoolLiteral"

	VoidNode NodeType = "Void"
)
