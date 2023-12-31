package sqlbuilder

type insertBuilder SqlBuilder

type insertBuilderFields SqlBuilder

type insertBuilderTable SqlBuilder

type insertBuilderValues SqlBuilder

type insertBuilderSelect SqlBuilder

type insertBuilderDup SqlBuilder

func (b *insertBuilder) init(kws []Keyword) *insertBuilder {
	b.buf.WriteString("INSERT")
	for _, kw := range kws {
		b.buf.Space()
		b.buf.WriteString(string(kw))
	}
	return b
}

func (b *insertBuilder) Into(table string) *insertBuilderTable {
	b.buf.Space()
	b.buf.WriteString("INTO")
	b.buf.Space()
	b.buf.BackQuoteString(table)
	return (*insertBuilderTable)(b)
}

func (b *insertBuilder) IntoT(table *Table) *insertBuilderTable {
	b.buf.Space()
	b.buf.WriteString("INTO")
	b.buf.Space()
	b.buf.Table(table)
	return (*insertBuilderTable)(b)
}

func (b *insertBuilderTable) Fields(fields ...string) *insertBuilderFields {
	b.buf.Space()
	b.buf.OpenParen()
	b.buf.BackQuoteStrings(fields)
	b.buf.CloseParen()
	return (*insertBuilderFields)(b)
}

func (b *insertBuilderTable) Select(subquery string, args ...any) *insertBuilderSelect {
	return (*insertBuilderSelect)(b).selectSub(subquery, args)
}

func (b *insertBuilderFields) Values(args ...any) *insertBuilderValues {
	b.buf.Space()
	b.buf.WriteString("VALUES")
	b.buf.Space()

	b.buf.WriteString(QuestionMarks(len(args)))
	b.args = append(b.args, args...)

	return (*insertBuilderValues)(b)
}

func (b *insertBuilderFields) Bulk(n int, argf func(index int) []any) *insertBuilderValues {
	b.buf.Space()
	b.buf.WriteString("VALUES")
	b.buf.Space()

	numsOfArgs := -1
	for i := 0; i < n; i++ {
		args := argf(i)
		if numsOfArgs == -1 {
			numsOfArgs = len(args)
			b.args = make([]any, 0, n*numsOfArgs)
		}
		b.args = append(b.args, args...)
	}

	qs := QuestionMarks(numsOfArgs)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.buf.WriteByte(',')
		}
		b.buf.WriteString(qs)
	}

	return (*insertBuilderValues)(b)
}

func (b *insertBuilderFields) Select(subquery string, args ...any) *insertBuilderSelect {
	return (*insertBuilderSelect)(b).selectSub(subquery, args)
}

func (b *insertBuilderValues) OnDuplicate(vps ...valueUpdater) *insertBuilderDup {
	b.buf.Space()
	b.buf.WriteString("ON DUPLICATE KEY UPDATE")
	b.buf.Space()
	b.buf.ValueUpdater(vps)
	for i := range vps {
		b.args = append(b.args, vps[i].args()...)
	}
	return (*insertBuilderDup)(b)
}

func (b *insertBuilderValues) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}

func (b *insertBuilderSelect) selectSub(subquery string, args []any) *insertBuilderSelect {
	b.buf.Space()
	b.buf.WriteString(subquery)
	b.args = append(b.args, args...)
	return b
}

func (b *insertBuilderSelect) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}

func (b *insertBuilderDup) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}
