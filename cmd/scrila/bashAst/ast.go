package bashAst

type NodeType string

const (
	// Statements
	BashStmtNode        NodeType = "BashStmt"
	CommentNode         NodeType = "CommentStmt"
	FuncDeclarationNode NodeType = "FuncDeclarationStmt"
	IfStmtNode          NodeType = "IfStmt"
	ProgramNode         NodeType = "ProgramStmt"
	WhileStmtNode       NodeType = "WhileStmt"

	// Expressions
	AssignmentExprNode NodeType = "AssignmentExpr"
	BinaryCompExprNode NodeType = "BinaryCompExpr"
	BinaryOpExprNode   NodeType = "BinaryOpExpr"
	CallExprNode       NodeType = "CallExpr"
	ReturnExprNode     NodeType = "ReturnExpr"

	// Literals
	BoolLiteralNode NodeType = "BoolLiteral"
	IntLiteralNode  NodeType = "IntLiteral"
	StrLiteralNode  NodeType = "StrLiteral"
	VarLiteralNode  NodeType = "VarLiteral"

	VoidNode NodeType = "Void"
)
