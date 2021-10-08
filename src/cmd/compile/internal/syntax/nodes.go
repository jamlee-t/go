// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syntax

// JAMLEE: Node 是一个 interface。File 结构体实现了它。从定义上来看，File 结构体仅仅是存储
// 各种声明。Node，表达式是一个 node，语句也是一个 node。那node之间的关系是包含关系还是并列关系呢？
// ----------------------------------------------------------------------------
// Nodes

type Node interface {
	// Pos() returns the position associated with the node as follows:
	// 1) The position of a node representing a terminal syntax production
	//    (Name, BasicLit, etc.) is the position of the respective production
	//    in the source.
	// 2) The position of a node representing a non-terminal production
	//    (IndexExpr, IfStmt, etc.) is the position of a token uniquely
	//    associated with that production; usually the left-most one
	//    ('[' for IndexExpr, 'if' for IfStmt, etc.)
	Pos() Pos
	aNode()
}

// JAMLEE: node 是 Node 接口的实现
type node struct {
	// commented out for now since not yet used
	// doc  *Comment // nil means no comment(s) attached
	pos Pos
}

func (n *node) Pos() Pos { return n.pos }
func (*node) aNode()     {}

// JAMLEE: File 结构体表示的是最终结果吗？根据 dumper_test.go 看是可以的
// ----------------------------------------------------------------------------
// Files

// package PkgName; DeclList[0], DeclList[1], ...
type File struct {
	Pragma   Pragma
	PkgName  *Name
	DeclList []Decl
	Lines    uint
	node
}

// JAMLEE: 批量定义顶级声明，写在文件中的顶级的，不被其他内容包含
// ----------------------------------------------------------------------------
// Declarations

type (
	Decl interface {
		Node
		aDecl()
	}

	//              Path
	// LocalPkgName Path
	ImportDecl struct {
		Group        *Group // nil means not part of a group, JAMLEE: 将声明和group一起关联起来
		Pragma       Pragma
		LocalPkgName *Name // including "."; nil means no rename present， JAMLEE: 表示包的名字，表达式
		Path         *BasicLit // JAMLEE: 基本字面量表示，表达式。
		decl
	}

	// JAMLEE: ConstDecl 和  VarDecl 英国蕾西
	// NameList
	// NameList      = Values
	// NameList Type = Values
	ConstDecl struct {
		Group    *Group // nil means not part of a group, JAMLEE: 将常量和group一起关联起来
		Pragma   Pragma
		NameList []*Name
		Type     Expr // nil means no type
		Values   Expr // nil means no values
		decl
	}

	// Name Type
	TypeDecl struct {
		Group  *Group // nil means not part of a group, JAMLEE: 将类型和group一起关联起来
		Pragma Pragma
		Name   *Name
		Alias  bool
		Type   Expr
		decl
	}

	// JAMLEE: 这里类型为什么没有像 int 这种类型选择呢？因为 type 和 vaules 都用 Name 来表示的。这里负责处理所有 var 开头的语句。隐含 type 时，type必定为nil
	// NameList Type
	// NameList Type = Values
	// NameList      = Values
	VarDecl struct {
		Group    *Group // nil means not part of a group, JAMLEE: 将变量和group一起关联起来
		Pragma   Pragma
		NameList []*Name // JAMLEE: 存在一次声明多个同类型变量的情况
		Type     Expr // nil means no type, JAMLEE: 这里有很多类型可以选择例如，函数类型。类型用表达式 Name 表示, 也就是「字符串」。而不是 BasicLit
		              // 如果声明为 var func01 = func(int){}. 此时 type 为 nil (因为隐含了)， NameList 为 ["func01"]
		Values   Expr // nil means no values, 可以是 BasicLit, FuncLit
		decl
	}

	// JAMLEE: 函数不用关联 group
	// func          Name Type { Body }
	// func          Name Type
	// func Receiver Name Type { Body }
	// func Receiver Name Type
	FuncDecl struct {
		Pragma Pragma
		Recv   *Field // nil means regular function
		Name   *Name
		Type   *FuncType
		Body   *BlockStmt // nil means no body (forward declaration)
		decl
	}
)

// JAMLEE: decl 也是 一个 node
type decl struct{ node }

func (*decl) aDecl() {}

// JAMLEE: 所有的声明属于一个 group，指向同一个 group 节点
// All declarations belonging to the same group point to the same Group node.
type Group struct {
	dummy int // not empty so we are guaranteed different Group instances
}

// JAMLEE: 表达式的所有定义
// ----------------------------------------------------------------------------
// Expressions

type (
	Expr interface {
		Node
		aExpr()
	}

	// Placeholder for an expression that failed to parse
	// correctly and where we can't provide a better node.
	BadExpr struct {
		expr
	}

	// JAMLEE: 表示例如函数名，变量名。例如: var test int32, 其中 VarDecl 中 Type 是 "int32", NameList 是 ["test"]
	// Value
	Name struct {
		Value string
		expr
	}

	// JAMLEE: 在语法树中表示基本字面量数值。还可以表示引入包例如 "fmt" 的值
	// Value
	BasicLit struct {
		Value string  // JAMLEE: 字面量统一用字符串表示
		Kind  LitKind // JAMLEE: 表示基本类型赋值，但是这里没有类型区分 int64 和 int32 这种区别。
		Bad   bool // true means the literal Value has syntax errors
		expr
	}

	// Type { ElemList[0], ElemList[1], ... }
	CompositeLit struct {
		Type     Expr // nil means no literal type
		ElemList []Expr
		NKeys    int // number of elements with keys
		Rbrace   Pos
		expr
	}

	// Key: Value
	KeyValueExpr struct {
		Key, Value Expr
		expr
	}

	// JAMLEE: 匿名函数类型，可以赋值给变量
	// func Type { Body }
	FuncLit struct {
		Type *FuncType
		Body *BlockStmt
		expr
	}

	// (X)
	ParenExpr struct {
		X Expr
		expr
	}

	// X.Sel
	SelectorExpr struct {
		X   Expr
		Sel *Name
		expr
	}

	// X[Index]
	IndexExpr struct {
		X     Expr
		Index Expr
		expr
	}

	// X[Index[0] : Index[1] : Index[2]]
	SliceExpr struct {
		X     Expr
		Index [3]Expr
		// Full indicates whether this is a simple or full slice expression.
		// In a valid AST, this is equivalent to Index[2] != nil.
		// TODO(mdempsky): This is only needed to report the "3-index
		// slice of string" error when Index[2] is missing.
		Full bool
		expr
	}

	// X.(Type)
	AssertExpr struct {
		X    Expr
		Type Expr
		expr
	}

	// X.(type)
	// Lhs := X.(type)
	TypeSwitchGuard struct {
		Lhs *Name // nil means no Lhs :=
		X   Expr  // X.(type)
		expr
	}

	Operation struct {
		Op   Operator
		X, Y Expr // Y == nil means unary expression
		expr
	}

	// Fun(ArgList[0], ArgList[1], ...)
	CallExpr struct {
		Fun     Expr
		ArgList []Expr // nil means no arguments
		HasDots bool   // last argument is followed by ...
		expr
	}

	// ElemList[0], ElemList[1], ...
	ListExpr struct {
		ElemList []Expr
		expr
	}

	// [Len]Elem
	ArrayType struct {
		// TODO(gri) consider using Name{"..."} instead of nil (permits attaching of comments)
		Len  Expr // nil means Len is ...
		Elem Expr
		expr
	}

	// []Elem
	SliceType struct {
		Elem Expr
		expr
	}

	// ...Elem
	DotsType struct {
		Elem Expr
		expr
	}

	// struct { FieldList[0] TagList[0]; FieldList[1] TagList[1]; ... }
	StructType struct {
		FieldList []*Field
		TagList   []*BasicLit // i >= len(TagList) || TagList[i] == nil means no tag for field i
		expr
	}

	// Name Type
	//      Type
	Field struct {
		Name *Name // nil means anonymous field/parameter (structs/parameters), or embedded interface (interfaces)
		Type Expr  // field names declared in a list share the same Type (identical pointers)
		node
	}

	// interface { MethodList[0]; MethodList[1]; ... }
	InterfaceType struct {
		MethodList []*Field
		expr
	}

	// JAMLEE: 表示在例如变量声明中所显示的赋值函数类型
	FuncType struct {
		ParamList  []*Field
		ResultList []*Field
		expr
	}

	// map[Key]Value
	MapType struct {
		Key, Value Expr
		expr
	}

	//   chan Elem
	// <-chan Elem
	// chan<- Elem
	ChanType struct {
		Dir  ChanDir // 0 means no direction
		Elem Expr
		expr
	}
)

// JAMLEE: expr 继承了 node。所以表达式是 1 个 node
type expr struct{ node }

func (*expr) aExpr() {}

type ChanDir uint

const (
	_ ChanDir = iota
	SendOnly
	RecvOnly
)

// JAMLEE: 声明所有的语句。分为 Stmt 和 SimpleStmt  两种。每一个语句也是一个 node
// ----------------------------------------------------------------------------
// Statements

type (
	Stmt interface {
		Node
		aStmt()
	}

	SimpleStmt interface {
		Stmt
		aSimpleStmt()
	}

	EmptyStmt struct {
		simpleStmt
	}

	LabeledStmt struct {
		Label *Name
		Stmt  Stmt
		stmt
	}

	BlockStmt struct {
		List   []Stmt
		Rbrace Pos
		stmt
	}

	ExprStmt struct {
		X Expr
		simpleStmt
	}

	SendStmt struct {
		Chan, Value Expr // Chan <- Value
		simpleStmt
	}

	DeclStmt struct {
		DeclList []Decl
		stmt
	}

	AssignStmt struct {
		Op       Operator // 0 means no operation
		Lhs, Rhs Expr     // Rhs == ImplicitOne means Lhs++ (Op == Add) or Lhs-- (Op == Sub)
		simpleStmt
	}

	BranchStmt struct {
		Tok   token // Break, Continue, Fallthrough, or Goto
		Label *Name
		// Target is the continuation of the control flow after executing
		// the branch; it is computed by the parser if CheckBranches is set.
		// Target is a *LabeledStmt for gotos, and a *SwitchStmt, *SelectStmt,
		// or *ForStmt for breaks and continues, depending on the context of
		// the branch. Target is not set for fallthroughs.
		Target Stmt
		stmt
	}

	CallStmt struct {
		Tok  token // Go or Defer
		Call *CallExpr
		stmt
	}

	ReturnStmt struct {
		Results Expr // nil means no explicit return values
		stmt
	}

	IfStmt struct {
		Init SimpleStmt
		Cond Expr
		Then *BlockStmt
		Else Stmt // either nil, *IfStmt, or *BlockStmt
		stmt
	}

	ForStmt struct {
		Init SimpleStmt // incl. *RangeClause
		Cond Expr
		Post SimpleStmt
		Body *BlockStmt
		stmt
	}

	SwitchStmt struct {
		Init   SimpleStmt
		Tag    Expr // incl. *TypeSwitchGuard
		Body   []*CaseClause
		Rbrace Pos
		stmt
	}

	SelectStmt struct {
		Body   []*CommClause
		Rbrace Pos
		stmt
	}
)

type (
	RangeClause struct {
		Lhs Expr // nil means no Lhs = or Lhs :=
		Def bool // means :=
		X   Expr // range X
		simpleStmt
	}

	CaseClause struct {
		Cases Expr // nil means default clause
		Body  []Stmt
		Colon Pos
		node
	}

	CommClause struct {
		Comm  SimpleStmt // send or receive stmt; nil means default clause
		Body  []Stmt
		Colon Pos
		node
	}
)

// JAMLEE: 语句也是 1 个 node
type stmt struct{ node }

func (stmt) aStmt() {}

type simpleStmt struct {
	stmt
}

func (simpleStmt) aSimpleStmt() {}

// ----------------------------------------------------------------------------
// Comments

// TODO(gri) Consider renaming to CommentPos, CommentPlacement, etc.
//           Kind = Above doesn't make much sense.
type CommentKind uint

const (
	Above CommentKind = iota
	Below
	Left
	Right
)

type Comment struct {
	Kind CommentKind
	Text string
	Next *Comment
}
