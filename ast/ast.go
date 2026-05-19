package ast

import (
	"interpreter/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out strings.Builder

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct {
	Token token.Token
	Name *Identifier
	Value Expression
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) statementNode() {}

// Return something like "let x = 5;" for a let statement.
func (ls *LetStatement) String() string {
	var out strings.Builder
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// Satisfies Expression interface.
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

// ReturnStatement should start with 'return' keyword followed by
// a return value which should be an expression. The return value is optional, so it can be nil.
type ReturnStatement struct {
	Token token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
// Return something like "return 5;".
func (rs *ReturnStatement) String() string {
	var out strings.Builder
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatements are expressions which are allowed to
// behave somewhat like pseudo-statements. They are not actually statements
// But are expressions sharing some similarities with statements.
// e.g let x = 5; is a full-blown statement.
// x + 10; is an expression statement.
type ExpressionStatement struct {
	Token token.Token // The first token of the expression
	Expression Expression
}

// To implement the Statement interface.
func (es *ExpressionStatement) statementNode() {}
// To implement the Node interface.
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// Satisfies Node interface via TokenLiteral() and String() methods.
// Satisfies Expression interface via Node interface and 
// expressionNode() method. IntegerLiteral refers to such thing as
// "5;". That singular thing is an expression.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string { 
	return il.Token.Literal 
}

// A prefix expression consists of an operator and an expression
// to the right of the operator. e.g "-5", "!foobar". 
// It satisfies the Node interface via the String() and TokenLiteral()
// methods. It satisfies the Expression interface via the Node interface
// and expressionNode() method.
type PrefixExpression struct {
	Operator string
	Right Expression
}

var _ Expression

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string {
	return pe.Operator
}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Operator
}