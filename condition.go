package sqlbuilder

type _whereCondition interface {
	whereCondition()
}

type whereCondition interface {
	_whereCondition
	args() []any
	write(*buffer)
}

func And(conditions ...whereCondition) *BoolCondition {
	return &BoolCondition{
		Op:         AndOperator,
		Conditions: conditions,
	}
}

func Or(conditions ...whereCondition) *BoolCondition {
	return &BoolCondition{
		Op:         OrOperator,
		Conditions: conditions,
	}
}

func Lt(f *Field, arg any) *BinaryCondition {
	return &BinaryCondition{
		Field: f,
		Op:    LtOperator,
		Args:  []any{arg},
	}
}

func Le(f *Field, arg any) *BinaryCondition {
	return &BinaryCondition{
		Field: f,
		Op:    LeOperator,
		Args:  []any{arg},
	}
}

func Eq(f *Field, arg any) *BinaryCondition {
	return &BinaryCondition{
		Field: f,
		Op:    EqOperator,
		Args:  []any{arg},
	}
}

func Ge(f *Field, arg any) *BinaryCondition {
	return &BinaryCondition{
		Field: f,
		Op:    GeOperator,
		Args:  []any{arg},
	}
}

func Gt(f *Field, arg any) *BinaryCondition {
	return &BinaryCondition{
		Field: f,
		Op:    GtOperator,
		Args:  []any{arg},
	}
}

func Ne(f *Field, arg any) *BinaryCondition {
	return &BinaryCondition{
		Field: f,
		Op:    NeOperator,
		Args:  []any{arg},
	}
}

func Between(f *Field, ge any, le any) *BinaryCondition {
	return &BinaryCondition{
		Field: f,
		Op:    BetweenOperator,
		Args:  []any{ge, le},
	}
}

func Like(f *Field, args ...any) *BinaryCondition {
	return &BinaryCondition{
		Field: f,
		Op:    LikeOperator,
		Args:  args,
	}
}

func IsNull(f *Field) *UnaryCondition {
	return &UnaryCondition{
		Field: f,
		Op:    IsNullOperator,
	}
}

func NotNull(f *Field) *UnaryCondition {
	return &UnaryCondition{
		Field: f,
		Op:    NotNullOperator,
	}
}

func In(f *Field, args ...any) *InCondition {
	return &InCondition{
		Field: f,
		Op:    InOperator,
		Args:  args,
	}
}

func NotIn(f *Field, args ...any) *InCondition {
	return &InCondition{
		Field: f,
		Op:    NotInOperator,
		Args:  args,
	}
}

func Exists(subQuerySql string) *otherCondition {
	return &otherCondition{
		Op:   ExistsOperator,
		Args: []any{subQuerySql},
	}
}

func NotExists(subQuerySql string) *otherCondition {
	return &otherCondition{
		Op:   NotExistsOperator,
		Args: []any{subQuerySql},
	}
}

func Condition(expr string, args ...any) *AnyCondition {
	return &AnyCondition{Expr: expr, Args: args}
}

type BoolCondition struct {
	_whereCondition
	Op         ConditionOperator
	Conditions []whereCondition
}

func (c *BoolCondition) args() []any {
	list := make([]any, 0)
	for _, cd := range c.Conditions {
		list = append(list, cd.args()...)
	}
	return list
}

func (c *BoolCondition) write(buf *buffer) {
	buf.OpenParen()
	for i, cd := range c.Conditions {
		if i > 0 {
			buf.Space()
			buf.WriteString(string(c.Op))
			buf.Space()
		}
		cd.write(buf)
	}
	buf.CloseParen()
}

type UnaryCondition struct {
	_whereCondition
	Field *Field
	Op    ConditionOperator
}

func (c *UnaryCondition) write(buf *buffer) {
	buf.Field(c.Field)
	buf.Space()
	buf.WriteString(string(c.Op))
}

func (c *UnaryCondition) args() []any {
	return nil
}

type BinaryCondition struct {
	_whereCondition
	Field *Field
	Op    ConditionOperator
	Args  []any
}

func (c *BinaryCondition) write(buf *buffer) {
	switch c.Op {
	case BetweenOperator:
		buf.Field(c.Field)
		buf.Space()
		buf.WriteString(string(BetweenOperator))
		buf.Space()
		buf.Question()
		buf.Space()
		buf.WriteString(string(AndOperator))
		buf.Space()
		buf.Question()
	default:
		buf.Field(c.Field)
		buf.Space()
		buf.WriteString(string(c.Op))
		buf.Space()
		buf.Question()
	}
}

func (c *BinaryCondition) args() []any {
	return c.Args
}

type InCondition struct {
	_whereCondition
	Field *Field
	Op    ConditionOperator
	Args  []any
}

func (c *InCondition) write(buf *buffer) {
	buf.Field(c.Field)
	buf.Space()
	buf.WriteString(string(c.Op))
	buf.Space()
	buf.WriteString(QuestionMarks(len(c.Args)))
}

func (c *InCondition) args() []any {
	return c.Args
}

type otherCondition struct {
	_whereCondition
	Field *Field
	Op    ConditionOperator
	Args  []any
}

func (c *otherCondition) write(buf *buffer) {
	switch c.Op {
	case ExistsOperator, NotExistsOperator:
		buf.WriteString(string(c.Op))
		buf.Space()
		buf.OpenParen()
		buf.WriteString(c.Args[0].(string))
		buf.CloseParen()
	}
}

func (c *otherCondition) args() []any {
	switch c.Op {
	case ExistsOperator, NotExistsOperator:
		return nil
	default:
		return nil
	}
}

type AnyCondition struct {
	_whereCondition
	Expr string
	Args []any
}

func (c *AnyCondition) args() []any {
	return c.Args
}

func (c *AnyCondition) write(buf *buffer) {
	buf.WriteString(c.Expr)
}
