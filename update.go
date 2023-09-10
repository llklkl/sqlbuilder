package sqlbuilder

type updateBuilder SqlBuilder

type updateBuilderTable SqlBuilder

type updateBuilderSet SqlBuilder

type updateBuilderWhere SqlBuilder

type updateBuilderOrder SqlBuilder

type updateBuilderLimit SqlBuilder

func (b *updateBuilder) init(kws []Keyword) *updateBuilder {
	b.buf.WriteString("UPDATE")
	for _, kw := range kws {
		b.buf.Space()
		b.buf.WriteString(string(kw))
	}
	return b
}

func (b *updateBuilder) Table(table string) *updateBuilderTable {
	b.buf.Space()
	b.buf.BackQuoteString(table)
	return (*updateBuilderTable)(b)
}

func (b *updateBuilder) TableT(table *Table) *updateBuilderTable {
	b.buf.Space()
	b.buf.Table(table)
	return (*updateBuilderTable)(b)
}

func (b *updateBuilderTable) Set(vps ...valueUpdater) *updateBuilderSet {
	b.buf.Space()
	b.buf.WriteString("SET")
	b.buf.Space()
	b.buf.ValueUpdater(vps)
	for i := range vps {
		b.args = append(b.args, vps[i].args()...)
	}
	return (*updateBuilderSet)(b)
}

func (b *updateBuilderSet) Where(conditions ...whereCondition) *updateBuilderWhere {
	b.buf.Space()
	b.buf.WriteString("WHERE")
	b.buf.Space()
	b.buf.Conditions(conditions)
	for i := range conditions {
		b.args = append(b.args, conditions[i].args()...)
	}
	return (*updateBuilderWhere)(b)
}

func (b *updateBuilderSet) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}

func (b *updateBuilderWhere) Order(orderSpecs ...*OrderSpec) *updateBuilderOrder {
	return (*updateBuilderOrder)(b).order(orderSpecs)
}

func (b *updateBuilderWhere) Limit(limit any) *updateBuilderLimit {
	return (*updateBuilderLimit)(b).limit(limit)
}

func (b *updateBuilderWhere) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}

func (b *updateBuilderOrder) order(orderSpecs []*OrderSpec) *updateBuilderOrder {
	b.buf.Space()
	b.buf.WriteString("ORDER BY")
	b.buf.Space()
	b.buf.OrderSpecs(orderSpecs)
	return b
}

func (b *updateBuilderOrder) Limit(limit any) *updateBuilderLimit {
	return (*updateBuilderLimit)(b).limit(limit)
}

func (b *updateBuilderOrder) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}

func (b *updateBuilderLimit) limit(limit any) *updateBuilderLimit {
	b.buf.Space()
	b.buf.WriteString("LIMIT")
	b.buf.Space()
	b.buf.Question()
	b.args = append(b.args, limit)
	return b
}

func (b *updateBuilderLimit) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}
