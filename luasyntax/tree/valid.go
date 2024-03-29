package tree

import (
	"github.com/GJSBRT/FiveMCompiler/luasyntax/token"
)

// When determining the validity of a node, there are three primary kinds of
// fields to consider:
//     - An interface: The underlying type is a child node.
//     - A Node: A child node.
//     - A Token: Not a node, and therefore not a child node. A node that embeds
//       a Token is not considered a token.
//
// A node is valid
//     - if its interface fields are non-nil.
//     - if its token fields have types acceptable for the node.
//     - even if its child nodes are not valid.

func isv(i interface{}) bool               { return i != nil }
func ist(tok Token, typ token.Type) bool   { return tok.Type == typ }
func ist2(tok Token, a, b token.Type) bool { return tok.Type == a || tok.Type == b }

func (f *File) IsValid() bool {
	return ist(f.EOFToken, token.EOF)
}

func (b *Block) IsValid() bool {
	if len(b.Seps) != len(b.Items) {
		return false
	}
	for _, stmt := range b.Items {
		if !isv(stmt) {
			return false
		}
	}
	for _, sep := range b.Seps {
		if !ist2(sep, token.SEMICOLON, token.INVALID) {
			return false
		}
	}
	return true
}

func (l *ExprList) IsValid() bool {
	if len(l.Items) == 0 || len(l.Seps) != len(l.Items)-1 {
		return false
	}
	for _, item := range l.Items {
		if !isv(item) {
			return false
		}
	}
	for _, sep := range l.Seps {
		if !ist(sep, token.COMMA) {
			return false
		}
	}
	return true
}

func (l *NameList) IsValid() bool {
	if len(l.Items) == 0 || len(l.Seps) != len(l.Items)-1 {
		return false
	}
	for _, item := range l.Items {
		if !ist(item, token.NAME) {
			return false
		}
	}
	for _, sep := range l.Seps {
		if !ist(sep, token.COMMA) {
			return false
		}
	}
	return true
}

func (e *NumberExpr) IsValid() bool {
	return e.NumberToken.Type.IsNumber()
}

func (e *StringExpr) IsValid() bool {
	return e.StringToken.Type.IsString()
}

func (e *NilExpr) IsValid() bool {
	return ist(e.NilToken, token.NIL)
}

func (e *BoolExpr) IsValid() bool {
	return e.BoolToken.Type.IsBool()
}

func (e *VarArgExpr) IsValid() bool {
	return ist(e.VarArgToken, token.VARARG)
}

func (e *UnopExpr) IsValid() bool {
	return e.UnopToken.Type.IsUnary() &&
		isv(e.Operand)
}

func (e *BinopExpr) IsValid() bool {
	return isv(e.Left) &&
		e.BinopToken.Type.IsBinary() &&
		isv(e.Right)
}

func (e *ParenExpr) IsValid() bool {
	return ist(e.LParenToken, token.LPAREN) &&
		isv(e.Value) &&
		ist(e.RParenToken, token.RPAREN)
}

func (e *VariableExpr) IsValid() bool {
	return ist(e.NameToken, token.NAME)
}

func (e *TableCtor) IsValid() bool {
	return ist(e.LBraceToken, token.LBRACE) &&
		ist(e.RBraceToken, token.RBRACE)
}

func (l *EntryList) IsValid() bool {
	if len(l.Seps) != len(l.Items) && len(l.Seps) != len(l.Items)-1 {
		return false
	}
	for _, entry := range l.Items {
		if !isv(entry) {
			return false
		}
	}
	for _, sep := range l.Seps {
		if !ist2(sep, token.COMMA, token.SEMICOLON) {
			return false
		}
	}
	return true
}

func (e *IndexEntry) IsValid() bool {
	return ist(e.LBrackToken, token.LBRACK) &&
		isv(e.Key) &&
		ist(e.RBrackToken, token.RBRACK) &&
		ist(e.AssignToken, token.ASSIGN) &&
		isv(e.Value)
}

func (e *FieldEntry) IsValid() bool {
	return ist(e.NameToken, token.NAME) &&
		ist(e.AssignToken, token.ASSIGN) &&
		isv(e.Value)
}

func (e *ValueEntry) IsValid() bool {
	return isv(e.Value)
}

func (e *FunctionExpr) IsValid() bool {
	if !(ist(e.FuncToken, token.FUNCTION) &&
		ist(e.LParenToken, token.LPAREN) &&
		ist(e.RParenToken, token.RPAREN) &&
		ist(e.EndToken, token.END)) {
		return false
	}
	if ist(e.VarArgToken, token.VARARG) {
		if e.Params != nil {
			return ist(e.VarArgSepToken, token.COMMA)
		}
		return ist(e.VarArgSepToken, token.INVALID)
	} else if ist(e.VarArgToken, token.INVALID) {
		return ist(e.VarArgSepToken, token.INVALID)
	}
	return false
}

func (e *FieldExpr) IsValid() bool {
	return isv(e.Value) &&
		ist(e.DotToken, token.DOT) &&
		ist(e.NameToken, token.NAME)
}

func (e *IndexExpr) IsValid() bool {
	return isv(e.Value) &&
		ist(e.LBrackToken, token.LBRACK) &&
		isv(e.Index) &&
		ist(e.RBrackToken, token.RBRACK)
}

func (e *MethodExpr) IsValid() bool {
	return isv(e.Value) &&
		ist(e.ColonToken, token.COLON) &&
		ist(e.NameToken, token.NAME) &&
		isv(e.Args)
}

func (e *CallExpr) IsValid() bool {
	return isv(e.Value) &&
		isv(e.Args)
}

func (c *ListArgs) IsValid() bool {
	return ist(c.LParenToken, token.LPAREN) &&
		ist(c.RParenToken, token.RPAREN)
}

func (c *TableArg) IsValid() bool {
	return true
}

func (c *StringArg) IsValid() bool {
	return true
}

func (s *DoStmt) IsValid() bool {
	return ist(s.DoToken, token.DO) &&
		ist(s.EndToken, token.END)
}

func (s *AssignStmt) IsValid() bool {
	return ist(s.AssignToken, token.ASSIGN)
}

func (s *CallStmt) IsValid() bool {
	return isv(s.Call)
}

func (s *IfStmt) IsValid() bool {
	return ist(s.IfToken, token.IF) &&
		isv(s.Cond) &&
		ist(s.ThenToken, token.THEN) &&
		ist(s.EndToken, token.END)
}

func (c *ElseIfClause) IsValid() bool {
	return ist(c.ElseIfToken, token.ELSEIF) &&
		isv(c.Cond) &&
		ist(c.ThenToken, token.THEN)
}

func (c *ElseClause) IsValid() bool {
	return ist(c.ElseToken, token.ELSE)
}

func (s *NumericForStmt) IsValid() bool {
	if !(ist(s.ForToken, token.FOR) &&
		ist(s.NameToken, token.NAME) &&
		ist(s.AssignToken, token.ASSIGN) &&
		isv(s.Min) &&
		ist(s.MaxSepToken, token.COMMA) &&
		isv(s.Max) &&
		ist(s.DoToken, token.DO) &&
		ist(s.EndToken, token.END)) {
		return false
	}
	if ist(s.StepSepToken, token.COMMA) {
		return isv(s.Step)
	} else if ist(s.StepSepToken, token.INVALID) {
		return !isv(s.Step)
	}
	return false
}

func (s *GenericForStmt) IsValid() bool {
	return ist(s.ForToken, token.FOR) &&
		ist(s.InToken, token.IN) &&
		ist(s.DoToken, token.DO) &&
		ist(s.EndToken, token.END)
}

func (s *WhileStmt) IsValid() bool {
	return ist(s.WhileToken, token.WHILE) &&
		isv(s.Cond) &&
		ist(s.DoToken, token.DO) &&
		ist(s.EndToken, token.END)
}

func (s *RepeatStmt) IsValid() bool {
	return ist(s.RepeatToken, token.REPEAT) &&
		ist(s.UntilToken, token.UNTIL) &&
		isv(s.Cond)
}

func (s *LocalVarStmt) IsValid() bool {
	if !ist(s.LocalToken, token.LOCAL) {
		return false
	}
	if ist(s.AssignToken, token.ASSIGN) {
		return s.Values != nil
	} else if ist(s.AssignToken, token.INVALID) {
		return s.Values == nil
	}
	return false
}

func (s *LocalFunctionStmt) IsValid() bool {
	return ist(s.LocalToken, token.LOCAL) &&
		ist(s.NameToken, token.NAME)
}

func (s *FunctionStmt) IsValid() bool {
	return true
}

func (l *FuncNameList) IsValid() bool {
	if len(l.Items) == 0 || len(l.Seps) != len(l.Items)-1 {
		return false
	}
	for _, item := range l.Items {
		if !ist(item, token.NAME) {
			return false
		}
	}
	for _, sep := range l.Seps {
		if !ist(sep, token.DOT) {
			return false
		}
	}
	if ist(l.ColonToken, token.COLON) {
		return ist(l.MethodToken, token.NAME)
	} else if ist(l.ColonToken, token.INVALID) {
		return ist(l.MethodToken, token.INVALID)
	}
	return false
}

func (s *BreakStmt) IsValid() bool {
	return ist(s.BreakToken, token.BREAK)
}

func (s *ReturnStmt) IsValid() bool {
	return ist(s.ReturnToken, token.RETURN)
}
