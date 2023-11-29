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
	ArrayLiteralNode  NodeType = "Array"
	IntLiteralNode    NodeType = "IntLiteral"  // Also data type
	StrLiteralNode    NodeType = "StrLiteral"  // Also data type
	BoolLiteralNode   NodeType = "BoolLiteral" // Also data type

	// Data types
	VoidNode      NodeType = "Void"
	BoolArrayNode NodeType = "BoolArray"
	IntArrayNode  NodeType = "IntArray"
	StrArrayNode  NodeType = "StrArray"
)
