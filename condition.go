package sqlbuilder

type whereCondition interface {
	args() []any
	write(*buffer)
}

func And(cds ...whereCondition) *AndCondition {
	return &AndCondition{Conditions: cds}
}

func Or(cds ...whereCondition) *OrCondition {
	return &OrCondition{Conditions: cds}
}

func Condition(expr string, args ...any) *AnyCondition {
	return &AnyCondition{Expr: expr, Args: args}
}

type AndCondition struct {
	Conditions []whereCondition
}

func (c *AndCondition) args() []any {
	list := make([]any, 0)
	for _, cd := range c.Conditions {
		list = append(list, cd.args()...)
	}
	return list
}

func (c *AndCondition) write(buf *buffer) {
	buf.OpenParen()
	for _, cd := range c.Conditions {
		cd.write(buf)
	}
	buf.CloseParen()
}

type OrCondition struct {
	Conditions []whereCondition
}

func (c *OrCondition) args() []any {
	list := make([]any, 0)
	for _, cd := range c.Conditions {
		list = append(list, cd.args()...)
	}
	return list
}

func (c *OrCondition) write(buf *buffer) {
	buf.OpenParen()
	for _, cd := range c.Conditions {
		cd.write(buf)
	}
	buf.CloseParen()
}

type AnyCondition struct {
	Expr string
	Args []any
}

func (c *AnyCondition) args() []any {
	return c.Args
}

func (c *AnyCondition) write(buf *buffer) {
	buf.WriteString(c.Expr)
}
