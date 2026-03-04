package ast

import (
	"bytes"
	"modern-basic/token"
)

// Node is the base AST interface.
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement represents AST nodes that do not produce a value directly.
type Statement interface {
	Node
	statementNode()
}

// Expression represents AST nodes that produce a value.
type Expression interface {
	Node
	expressionNode()
}

// Program is the AST root containing all top-level statements.
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the first statement token literal or empty if program is empty.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

// String renders the full program in a deterministic form for testing.
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// LetStatement represents `let <ident> = <expression>;`.
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral returns the source token literal.
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// AssignStatement represents `<ident> = <expression>;`.
type AssignStatement struct {
	Token token.Token // token.IDENT
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode() {}

// TokenLiteral returns the source token literal.
func (as *AssignStatement) TokenLiteral() string {
	return as.Token.Literal
}

// String renders an assignment statement.
func (as *AssignStatement) String() string {
	var out bytes.Buffer

	if as.Name != nil {
		out.WriteString(as.Name.String())
	}
	out.WriteString(" = ")
	if as.Value != nil {
		out.WriteString(as.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

// String renders a let statement.
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral())
	out.WriteString(" ")
	if ls.Name != nil {
		out.WriteString(ls.Name.String())
	}
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

// ReturnStatement represents `return <expression>;`.
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral returns the source token literal.
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// String renders a return statement.
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral())
	out.WriteString(" ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

// ExpressionStatement wraps a bare expression used as a statement.
type ExpressionStatement struct {
	Token      token.Token // first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral returns the source token literal.
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// String renders the wrapped expression.
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

// Identifier represents a variable name.
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral returns the source token literal.
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// String renders the identifier name.
func (i *Identifier) String() string {
	return i.Value
}

// IntegerLiteral represents integer values.
type IntegerLiteral struct {
	Token token.Token // token.INT
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral returns the source token literal.
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// String renders the literal as it appeared in source.
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// StringLiteral represents string values.
type StringLiteral struct {
	Token token.Token // token.STRING
	Value string
}

func (sl *StringLiteral) expressionNode() {}

// TokenLiteral returns the source token literal.
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

// String renders the literal as it appeared in source.
func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

// Boolean represents true/false expressions.
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

// TokenLiteral returns the source token literal.
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// String renders the literal as it appeared in source.
func (b *Boolean) String() string {
	return b.Token.Literal
}

// PrefixExpression represents prefix operators like !x or -5.
type PrefixExpression struct {
	Token    token.Token // prefix token, e.g. ! or -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral returns the source token literal.
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

// String renders a parenthesized prefix expression.
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	if pe.Right != nil {
		out.WriteString(pe.Right.String())
	}
	out.WriteString(")")

	return out.String()
}

// InfixExpression represents binary operators like x + y.
type InfixExpression struct {
	Token    token.Token // infix token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral returns the source token literal.
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String renders a parenthesized infix expression.
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	if ie.Left != nil {
		out.WriteString(ie.Left.String())
	}
	out.WriteString(" ")
	out.WriteString(ie.Operator)
	out.WriteString(" ")
	if ie.Right != nil {
		out.WriteString(ie.Right.String())
	}
	out.WriteString(")")

	return out.String()
}

// BlockStatement represents a block of statements surrounded by braces.
type BlockStatement struct {
	Token      token.Token // token.LBRACE
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral returns the source token literal.
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// String renders all statements in order.
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// IfBranch stores one conditional branch in an if/elseif chain.
type IfBranch struct {
	Condition   Expression
	Consequence *BlockStatement
}

// IfExpression represents if/elseif/else control flow.
type IfExpression struct {
	Token       token.Token // token.IF
	Branches    []IfBranch
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral returns the source token literal.
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String renders the conditional expression.
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	for i, branch := range ie.Branches {
		if i == 0 {
			out.WriteString("if")
		} else {
			out.WriteString("elseif")
		}
		out.WriteString("(")
		if branch.Condition != nil {
			out.WriteString(branch.Condition.String())
		}
		out.WriteString(")")
		out.WriteString("{")
		if branch.Consequence != nil {
			out.WriteString(branch.Consequence.String())
		}
		out.WriteString("}")
	}

	if ie.Alternative != nil {
		out.WriteString("else")
		out.WriteString("{")
		out.WriteString(ie.Alternative.String())
		out.WriteString("}")
	}

	return out.String()
}

// WhileExpression represents a while loop.
type WhileExpression struct {
	Token     token.Token // token.WHILE
	Condition Expression
	Body      *BlockStatement
}

func (we *WhileExpression) expressionNode() {}

// TokenLiteral returns the source token literal.
func (we *WhileExpression) TokenLiteral() string {
	return we.Token.Literal
}

// String renders the while loop.
func (we *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString("while(")
	if we.Condition != nil {
		out.WriteString(we.Condition.String())
	}
	out.WriteString("){")
	if we.Body != nil {
		out.WriteString(we.Body.String())
	}
	out.WriteString("}")

	return out.String()
}

// ForExpression represents a C-style for loop.
type ForExpression struct {
	Token     token.Token // token.FOR
	Init      Statement
	Condition Expression
	Post      Statement
	Body      *BlockStatement
}

func (fe *ForExpression) expressionNode() {}

// TokenLiteral returns the source token literal.
func (fe *ForExpression) TokenLiteral() string {
	return fe.Token.Literal
}

// String renders the for loop.
func (fe *ForExpression) String() string {
	var out bytes.Buffer

	out.WriteString("for(")
	if fe.Init != nil {
		out.WriteString(fe.Init.String())
	}
	out.WriteString(";")
	if fe.Condition != nil {
		out.WriteString(fe.Condition.String())
	}
	out.WriteString(";")
	if fe.Post != nil {
		out.WriteString(fe.Post.String())
	}
	out.WriteString("){")
	if fe.Body != nil {
		out.WriteString(fe.Body.String())
	}
	out.WriteString("}")

	return out.String()
}

// FunctionLiteral represents `fn(x, y) { ... }`.
type FunctionLiteral struct {
	Token      token.Token // token.FUNCTION
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral returns the source token literal.
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

// String renders the function literal.
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	for i, p := range fl.Parameters {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(p.String())
	}
	out.WriteString(")")
	out.WriteString("{")
	if fl.Body != nil {
		out.WriteString(fl.Body.String())
	}
	out.WriteString("}")

	return out.String()
}

// CallExpression represents function invocation.
type CallExpression struct {
	Token     token.Token // token.LPAREN
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// TokenLiteral returns the source token literal.
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

// String renders the call expression.
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	if ce.Function != nil {
		out.WriteString(ce.Function.String())
	}
	out.WriteString("(")
	for i, arg := range ce.Arguments {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(arg.String())
	}
	out.WriteString(")")

	return out.String()
}

// DotExpression represents dotted access: left.right.
type DotExpression struct {
	Token token.Token // token.DOT
	Left  Expression
	Right *Identifier
}

func (de *DotExpression) expressionNode() {}

// TokenLiteral returns the source token literal.
func (de *DotExpression) TokenLiteral() string {
	return de.Token.Literal
}

// String renders a dotted expression.
func (de *DotExpression) String() string {
	var out bytes.Buffer

	if de.Left != nil {
		out.WriteString(de.Left.String())
	}
	out.WriteString(".")
	if de.Right != nil {
		out.WriteString(de.Right.String())
	}

	return out.String()
}

// ArrayLiteral represents list literals like [1, 2, 3].
type ArrayLiteral struct {
	Token    token.Token // token.LBRACKET
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

// TokenLiteral returns the source token literal.
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

// String renders the array literal.
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("[")
	for i, el := range al.Elements {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(el.String())
	}
	out.WriteString("]")

	return out.String()
}

// IndexExpression represents indexing into arrays/hashes: left[index].
type IndexExpression struct {
	Token token.Token // token.LBRACKET
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

// TokenLiteral returns the source token literal.
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String renders the index expression.
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	if ie.Left != nil {
		out.WriteString(ie.Left.String())
	}
	out.WriteString("[")
	if ie.Index != nil {
		out.WriteString(ie.Index.String())
	}
	out.WriteString("]")
	out.WriteString(")")

	return out.String()
}

// HashLiteral represents dictionary literals: {key: value}.
type HashLiteral struct {
	Token token.Token // token.LBRACE
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}

// TokenLiteral returns the source token literal.
func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

// String renders the hash literal.
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	i := 0
	for key, value := range hl.Pairs {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(key.String())
		out.WriteString(": ")
		out.WriteString(value.String())
		i++
	}
	out.WriteString("}")

	return out.String()
}
