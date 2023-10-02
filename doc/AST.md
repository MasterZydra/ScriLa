# AST Nodes

```mermaid
classDiagram
  Statement <|-- Comment
  Statement <|-- Program
  Statement <|-- VarDeclaration
  Statement <|-- FunctionDeclaration
  Statement <|-- Expr
  Expr <|-- AssignmentExpr
  Expr <|-- BinaryExpr
  Expr <|-- CallExpr
  Expr <|-- ReturnExpr
  Expr <|-- MemberExpr
  Expr <|-- Identifier
  Expr <|-- IntLiteral
  Expr <|-- StrLiteral
  Expr <|-- Property
  Expr <|-- ObjectLiteral
```
