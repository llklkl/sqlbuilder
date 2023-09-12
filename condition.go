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

func Lt(field any, arg any) *BinaryCondition {
	return &BinaryCondition{
		Field: field,
		Op:    LtOperator,
		Args:  []any{arg},
	}
}

func Eq(field any, arg any) *BinaryCondition {
	return &BinaryCondition{
		Field: field,
		Op:    EqOperator,
		Args:  []any{arg},
	}
}

func Le(field any, arg any) *BinaryCondition {
	return &BinaryCondition{
		Field: field,
		Op:    LeOperator,
		Args:  []any{arg},
	}
}

func Ge(field any, arg any) *BinaryCondition {
	return &BinaryCondition{
		Field: field,
		Op:    GeOperator,
		Args:  []any{arg},
	}
}

func Gt(field any, arg any) *BinaryCondition {
	return &BinaryCondition{
		Field: field,
		Op:    GtOperator,
		Args:  []any{arg},
	}
}

func Ne(field any, arg any) *BinaryCondition {
	return &BinaryCondition{
		Field: field,
		Op:    NeOperator,
		Args:  []any{arg},
	}
}

func Between(field any, ge any, le any) *BinaryCondition {
	return &BinaryCondition{
		Field: field,
		Op:    BetweenOperator,
		Args:  []any{ge, le},
	}
}

func Like(field any, args ...any) *BinaryCondition {
	return &BinaryCondition{
		Field: field,
		Op:    LikeOperator,
		Args:  args,
	}
}

func IsNull(field any) *UnaryCondition {
	return &UnaryCondition{
		Field: field,
		Op:    IsNullOperator,
	}
}

func NotNull(field any) *UnaryCondition {
	return &UnaryCondition{
		Field: field,
		Op:    NotNullOperator,
	}
}

func In(field any, args ...any) *InCondition {
	return &InCondition{
		Field: field,
		Op:    InOperator,
		Args:  args,
	}
}

func NotIn(field any, args ...any) *InCondition {
	return &InCondition{
		Field: field,
		Op:    NotInOperator,
		Args:  args,
	}
}

func Exists(subQuerySql string, args ...any) *SubqueryCondition {
	return &SubqueryCondition{
		Subquery: subQuerySql,
		Op:       ExistsOperator,
		Args:     args,
	}
}

func NotExists(subQuerySql string, args ...any) *SubqueryCondition {
	return &SubqueryCondition{
		Subquery: subQuerySql,
		Op:       NotExistsOperator,
		Args:     args,
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
	Field any
	Op    ConditionOperator
}

func (c *UnaryCondition) write(buf *buffer) {
	buf.AnyField(c.Field)
	buf.Space()
	buf.WriteString(string(c.Op))
}

func (c *UnaryCondition) args() []any {
	return nil
}

type BinaryCondition struct {
	_whereCondition
	Field any
	Op    ConditionOperator
	Args  []any
}

func (c *BinaryCondition) write(buf *buffer) {
	switch c.Op {
	case BetweenOperator:
		buf.AnyField(c.Field)
		buf.Space()
		buf.WriteString(string(BetweenOperator))
		buf.Space()
		buf.Question()
		buf.Space()
		buf.WriteString(string(AndOperator))
		buf.Space()
		buf.Question()
	default:
		buf.AnyField(c.Field)
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
	Field any
	Op    ConditionOperator
	Args  []any
}

func (c *InCondition) write(buf *buffer) {
	buf.AnyField(c.Field)
	buf.Space()
	buf.WriteString(string(c.Op))
	buf.Space()
	buf.WriteString(QuestionMarks(len(c.Args)))
}

func (c *InCondition) args() []any {
	return c.Args
}

type SubqueryCondition struct {
	_whereCondition
	Field    any
	Subquery string
	Op       ConditionOperator
	Args     []any
}

func (c *SubqueryCondition) write(buf *buffer) {
	if c.Field != nil {
		buf.AnyField(c.Field)
		buf.Space()
	}
	buf.WriteString(string(c.Op))
	buf.Space()
	buf.OpenParen()
	buf.WriteString(c.Subquery)
	buf.CloseParen()
}

func (c *SubqueryCondition) args() []any {
	return c.Args
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
