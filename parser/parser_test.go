package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct{
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, identifier string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let', got=%q", stmt.TokenLiteral())	
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", stmt)
		return false
	}

	if letStmt.Name.Value != identifier {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", identifier, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != identifier {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
		identifier, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

// This is a function that checks if the parser has any errors while parsing the program
// if there were any errors, it prints the errors and fails the test, otherwise, it returns
// without doing anything
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got=%q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not foobar. got=%s", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not foobar. got=%s", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf(
			"expected program to have 1 statement. got=%d",
			len(program.Statements),
		)
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
		program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf(
			"literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral(),
		)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input string
		operator string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements),
			)
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not *ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not *ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf(
				"exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator,
			)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, val int64) bool {
	il, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf(
			"expected expression to be *ast.IntegerLiteral, got=%T",
			exp,
		)
		return false
	}
	if il.Value != val {
		t.Errorf(
			"expected value to be %d, got=%d",
			il.Value, val,
		)
		return false
	}
	if il.TokenLiteral() != fmt.Sprintf("%d", val) {
		t.Errorf(
			"expected val string to be %s, got=%s",
			il.TokenLiteral(), val,
		)
		return false
	}
	return true
}