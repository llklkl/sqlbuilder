package sqlbuilder

type selectBuilder SqlBuilder

type selectBuilderExpr SqlBuilder

type selectBuilderTable SqlBuilder

type selectBuilderJoin SqlBuilder

type selectBuilderJoinSpec SqlBuilder

type selectBuilderWhere SqlBuilder

type selectBuilderGroup SqlBuilder

type selectBuilderOrder SqlBuilder

type selectBuilderLimit SqlBuilder

func (b *selectBuilder) init(kws []Keyword) *selectBuilder {
	b.buf.WriteString("SELECT")
	for _, kw := range kws {
		b.buf.WriteByte(space)
		b.buf.WriteString(string(kw))
	}
	return b
}

func (b *selectBuilder) Fields(fields ...string) *selectBuilderExpr {
	b.buf.Space()
	b.buf.BackQuoteStrings(fields)
	return (*selectBuilderExpr)(b)
}

func (b *selectBuilder) Expr(exprs ...*Expr) *selectBuilderExpr {
	b.buf.Space()
	b.buf.Exprs(exprs)

	return (*selectBuilderExpr)(b)
}

func (b *selectBuilderExpr) From(tables ...string) *selectBuilderTable {
	b.buf.Space()
	b.buf.WriteString("FROM")
	b.buf.Space()
	b.buf.BackQuoteStrings(tables)
	return (*selectBuilderTable)(b)
}

func (b *selectBuilderExpr) FromT(tables ...*Table) *selectBuilderTable {
	b.buf.Space()
	b.buf.WriteString("FROM")
	b.buf.Space()
	b.buf.Tables(tables)

	return (*selectBuilderTable)(b)
}

func (b *selectBuilderTable) LeftJoin(table *Table) *selectBuilderJoin {
	return (*selectBuilderJoin)(b).join(leftJoin, table)
}

func (b *selectBuilderTable) RightJoin(table *Table) *selectBuilderJoin {
	return (*selectBuilderJoin)(b).join(rightJoin, table)
}

func (b *selectBuilderTable) InnerJoin(table *Table) *selectBuilderJoin {
	return (*selectBuilderJoin)(b).join(innerJoin, table)
}

func (b *selectBuilderTable) Where(conditions ...whereCondition) *selectBuilderWhere {
	return (*selectBuilderWhere)(b).where(conditions)
}

func (b *selectBuilderJoin) join(typ string, table *Table) *selectBuilderJoin {
	b.buf.Space()
	b.buf.WriteString(typ)
	b.buf.Space()
	b.buf.WriteString("JOIN")
	b.buf.Space()
	b.buf.Table(table)
	return b
}

func (b *selectBuilderJoin) On(lhs, rhs *Field) *selectBuilderJoinSpec {
	b.buf.Space()
	b.buf.WriteString("ON")
	b.buf.Space()
	b.buf.Field(lhs)
	b.buf.Equal()
	b.buf.Field(rhs)
	return (*selectBuilderJoinSpec)(b)
}

func (b *selectBuilderJoin) Using(fields ...string) *selectBuilderJoinSpec {
	b.buf.Space()
	b.buf.WriteString("USING")
	b.buf.Space()
	b.buf.OpenParen()
	b.buf.BackQuoteStrings(fields)
	b.buf.CloseParen()
	return (*selectBuilderJoinSpec)(b)
}

func (b *selectBuilderJoinSpec) LeftJoin(table *Table) *selectBuilderJoin {
	return (*selectBuilderJoin)(b).join(leftJoin, table)
}

func (b *selectBuilderJoinSpec) RightJoin(table *Table) *selectBuilderJoin {
	return (*selectBuilderJoin)(b).join(rightJoin, table)
}

func (b *selectBuilderJoinSpec) InnerJoin(table *Table) *selectBuilderJoin {
	return (*selectBuilderJoin)(b).join(innerJoin, table)
}

func (b *selectBuilderJoinSpec) Where(conditions ...whereCondition) *selectBuilderWhere {
	return (*selectBuilderWhere)(b).where(conditions)
}

func (b *selectBuilderWhere) where(conditions []whereCondition) *selectBuilderWhere {
	b.buf.Space()
	b.buf.WriteString("WHERE")
	b.buf.Space()
	b.buf.Conditions(conditions)
	for _, c := range conditions {
		b.args = append(b.args, c.args()...)
	}
	return b
}

func (b *selectBuilderWhere) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}

func (b *selectBuilderWhere) OrderBy(orderSpecs ...*OrderSpec) *selectBuilderOrder {
	return (*selectBuilderOrder)(b).order(orderSpecs)
}

func (b *selectBuilderWhere) Limit(args ...any) *selectBuilderLimit {
	return (*selectBuilderLimit)(b).limit(args)
}

func (b *selectBuilderWhere) GroupBy(fields ...string) *selectBuilderGroup {
	b.buf.Space()
	b.buf.WriteString("GROUP BY")
	b.buf.Space()
	b.buf.BackQuoteStrings(fields)
	return (*selectBuilderGroup)(b)
}

func (b *selectBuilderGroup) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}

func (b *selectBuilderGroup) OrderBy(orderSpecs ...*OrderSpec) *selectBuilderOrder {
	return (*selectBuilderOrder)(b).order(orderSpecs)
}

func (b *selectBuilderGroup) Limit(args ...any) *selectBuilderLimit {
	return (*selectBuilderLimit)(b).limit(args)
}

func (b *selectBuilderOrder) order(orderSpecs []*OrderSpec) *selectBuilderOrder {
	b.buf.Space()
	b.buf.WriteString("ORDER BY")
	b.buf.Space()
	b.buf.OrderSpecs(orderSpecs)
	return b
}

func (b *selectBuilderOrder) Limit(args ...any) *selectBuilderLimit {
	return (*selectBuilderLimit)(b).limit(args)
}

func (b *selectBuilderOrder) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}

func (b *selectBuilderLimit) limit(args []any) *selectBuilderLimit {
	b.buf.Space()
	b.buf.WriteString("LIMIT")
	b.buf.Space()
	if len(args) == 1 {
		b.buf.Question()
		b.args = append(b.args, args...)
	} else if len(args) == 2 {
		b.buf.Question()
		b.buf.Comma()
		b.buf.Question()
		b.args = append(b.args, args...)
	}
	return b
}

func (b *selectBuilderLimit) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}
