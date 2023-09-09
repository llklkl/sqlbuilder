package sqlbuilder

type Keyword string

const (
	Ignore     Keyword = "IGNORE"
	SqlCache   Keyword = "SQL_CACHE"
	SqlNoCache Keyword = "SQL_NO_CACHE"

	leftJoin  = "LEFT"
	rightJoin = "RIGHT"
	innerJoin = "INNER"
)

type OrderDirection string

const (
	Asc  OrderDirection = "ASC"
	Desc OrderDirection = "DESC"
)

const (
	space            = ' '
	singleQuote      = '\''
	openParentheses  = '('
	closeParentheses = ')'
	comma            = ','
	dot              = '.'
	equalMark        = '='
	questionMark     = '?'
	backQuote        = '`'
)

type Table struct {
	Table    string
	Alias    string
	Database string
}

type Expr struct {
	Expr  string
	Alias string
	Table string
}

type Field struct {
	Field string
	Alias string
	Table string
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
