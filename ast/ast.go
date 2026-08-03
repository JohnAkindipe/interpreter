package ast

import (
	"fmt"
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
	Name  *Identifier
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
	Token       token.Token
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
	Token      token.Token // The first token of the expression
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

func (il *IntegerLiteral) expressionNode()      {}
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
	Token    token.Token
	Operator string
	Right    Expression
}

func (ie *PrefixExpression) expressionNode() {}
func (ie *PrefixExpression) String() string {
	var out strings.Builder
	out.WriteString("(")
	out.WriteString(ie.Operator)
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
func (ie *PrefixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

type InfixExpression struct {
	Left     Expression
	Token    token.Token
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) String() string {
	var out strings.Builder
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// A boolean AST node. Satisfies the expression interface
type Boolean struct {
	Token token.Token // first token of the boolean node
	Value bool        // value, either "true" or "false"
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // else block
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out strings.Builder
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out strings.Builder
	out.WriteString("{")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString("}")
	return out.String()
}

var _ Expression = (*FunctionLiteral)(nil)

type FunctionLiteral struct {
	Token      token.Token // i.e. "fn"
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out strings.Builder
	params := []string{}

	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.Token.Literal)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

var _ Expression = (*ParameterExpression)(nil)

type ParameterExpression struct {
	Token token.Token // i.e. "("
	Body  []Expression
}

func (pe *ParameterExpression) expressionNode() {}
func (pe *ParameterExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *ParameterExpression) String() string {
	var out strings.Builder
	out.WriteString("(")
	for i, v := range pe.Body {
		out.WriteString(v.String())
		if i < len(pe.Body)-2 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")
	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier (for named fn) or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// returns somn like: [getName]((1 + 2), (3 * 4))
func (ce *CallExpression) String() string {
	var out strings.Builder
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
