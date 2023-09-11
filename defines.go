package sqlbuilder

type Keyword string

const (
	Ignore     Keyword = "IGNORE"
	SqlCache   Keyword = "SQL_CACHE"
	SqlNoCache Keyword = "SQL_NO_CACHE"

	leftJoin  Keyword = "LEFT JOIN"
	rightJoin Keyword = "RIGHT JOIN"
	innerJoin Keyword = "INNER JOIN"
)

type OrderDirection string

const (
	Asc  OrderDirection = "ASC"
	Desc OrderDirection = "DESC"
)

type ConditionOperator string

const (
	AndOperator ConditionOperator = "AND"
	OrOperator  ConditionOperator = "OR"

	EqOperator ConditionOperator = "="
	LtOperator ConditionOperator = "<"
	LeOperator ConditionOperator = "<="
	GtOperator ConditionOperator = ">"
	GeOperator ConditionOperator = ">="
	NeOperator ConditionOperator = "!="

	BetweenOperator   ConditionOperator = "BETWEEN"
	ExistsOperator    ConditionOperator = "EXISTS"
	NotExistsOperator ConditionOperator = "NOT EXISTS"
	LikeOperator      ConditionOperator = "LIKE"

	InOperator    ConditionOperator = "IN"
	NotInOperator ConditionOperator = "NOT IN"

	IsNullOperator  ConditionOperator = "IS NULL"
	NotNullOperator ConditionOperator = "IS NOT NULL"
)

const (
	space            = ' '
	openParentheses  = '('
	closeParentheses = ')'
	comma            = ','
	dot              = '.'
	equalMark        = '='
	questionMark     = '?'
	backQuote        = '`'
)

type _selectField interface {
	selectField()
}

type Table struct {
	Table    string
	Alias    string
	Database string
}

// T Specify a table. Different numbers of parameters will have different effects.
//
// When the number of parameters is,
//
// 1: func (table string) *Table
//
// 2: func (table, alias string) *Table
//
// 3. func (database, table, alias string) *Table
func T(args ...string) *Table {
	table := &Table{}
	switch len(args) {
	case 1:
		table.Table = args[0]
	case 2:
		table.Table = args[0]
		table.Alias = args[1]
	case 3:
		table.Database = args[0]
		table.Table = args[1]
		table.Alias = args[2]
	}
	return table
}

type Expr struct {
	_selectField
	Expr  string
	Alias string
}

// E Specify an Expression. Different numbers of parameters will have different effects.
//
// When the number of parameters is,
//
// 1: func (expr string) *Expr
//
// 2: func (expr, alias string) *Expr
func E(args ...string) *Expr {
	expr := &Expr{}
	switch len(args) {
	case 1:
		expr.Expr = args[0]
	case 2:
		expr.Expr = args[0]
		expr.Alias = args[1]
	}
	return expr
}

type Field struct {
	_selectField
	Field string
	Alias string
	Table string
}

// F Specify a field. Different numbers of parameters will have different effects.
//
// When the number of parameters is,
//
// 1: func(field string) *Field
//
// 2: func(table, field string) *Field
//
// 3: func(table, field, alias string) *Field
func F(args ...string) *Field {
	f := &Field{}
	switch len(args) {
	case 1:
		f.Field = args[0]
	case 2:
		f.Table = args[0]
		f.Field = args[1]
	case 3:
		f.Table = args[0]
		f.Field = args[1]
		f.Alias = args[2]
	}
	return f
}

type OrderSpec struct {
	Field          *Field
	OrderDirection OrderDirection
}

func Order(f *Field, direction OrderDirection) *OrderSpec {
	return &OrderSpec{
		Field:          f,
		OrderDirection: direction,
	}
}
