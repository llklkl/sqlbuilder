package sqlbuilder

type deleteBuilder SqlBuilder

type deleteBuilderTable SqlBuilder

type deleteBuilderWhere SqlBuilder

type deleteBuilderOrder SqlBuilder

type deleteBuilderLimit SqlBuilder

func (b *deleteBuilder) init(kws []Keyword) *deleteBuilder {
	b.buf.WriteString("DELETE")
	for _, kw := range kws {
		b.buf.Space()
		b.buf.WriteString(string(kw))
	}
	return b
}

func (b *deleteBuilder) From(table string) *deleteBuilderTable {
	b.buf.Space()
	b.buf.WriteString("FROM")
	b.buf.Space()
	b.buf.BackQuoteString(table)
	return (*deleteBuilderTable)(b)
}

func (b *deleteBuilder) FromT(table *Table) *deleteBuilderTable {
	b.buf.Space()
	b.buf.Table(table)
	return (*deleteBuilderTable)(b)
}

func (b *deleteBuilderTable) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}

func (b *deleteBuilderTable) Where(conditions ...whereCondition) *deleteBuilderWhere {
	b.buf.Space()
	b.buf.WriteString("WHERE")
	b.buf.Space()
	b.buf.Conditions(conditions)
	for i := range conditions {
		b.args = append(b.args, conditions[i].args()...)
	}
	return (*deleteBuilderWhere)(b)
}

func (b *deleteBuilderWhere) Order(orderSpecs ...*OrderSpec) *deleteBuilderOrder {
	return (*deleteBuilderOrder)(b).order(orderSpecs)
}

func (b *deleteBuilderWhere) Limit(limit any) *deleteBuilderLimit {
	return (*deleteBuilderLimit)(b).limit(limit)
}

func (b *deleteBuilderWhere) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}

func (b *deleteBuilderOrder) order(orderSpecs []*OrderSpec) *deleteBuilderOrder {
	b.buf.Space()
	b.buf.WriteString("ORDER BY")
	b.buf.Space()
	b.buf.OrderSpecs(orderSpecs)
	return b
}

func (b *deleteBuilderOrder) Limit(limit any) *deleteBuilderLimit {
	return (*deleteBuilderLimit)(b).limit(limit)
}

func (b *deleteBuilderOrder) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}

func (b *deleteBuilderLimit) limit(limit any) *deleteBuilderLimit {
	b.buf.Space()
	b.buf.WriteString("LIMIT")
	b.buf.Space()
	b.buf.Question()
	b.args = append(b.args, limit)
	return b
}

func (b *deleteBuilderLimit) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}
