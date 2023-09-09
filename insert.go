package sqlbuilder

type insertBuilder SqlBuilder

type insertBuilderFields SqlBuilder

type insertBuilderTable SqlBuilder

type insertBuilderValues SqlBuilder

func (b *insertBuilder) init(kws []Keyword) *insertBuilder {
	b.buf.WriteString("INSERT")
	for _, kw := range kws {
		b.buf.Space()
		b.buf.WriteString(string(kw))
	}
	return b
}

func (b *insertBuilder) Into(tableName string) *insertBuilderTable {
	b.buf.Space()
	b.buf.WriteString("INTO ")
	b.buf.BackQuoteString(tableName)
	return (*insertBuilderTable)(b)
}

func (b *insertBuilderTable) Fields(fields ...string) *insertBuilderFields {
	b.buf.Space()
	b.buf.OpenParen()
	b.buf.BackQuoteStrings(fields)
	b.buf.CloseParen()
	return (*insertBuilderFields)(b)
}

func (b *insertBuilderFields) Values(args ...any) *insertBuilderValues {
	b.buf.Space()
	b.buf.WriteString("VALUES ")

	b.buf.WriteString(QuestionMarks(len(args)))
	b.args = append(b.args, args...)

	return (*insertBuilderValues)(b)
}

func (b *insertBuilderFields) ValuesFn(n int, argf func(index int) []any) *insertBuilderValues {
	b.buf.Space()
	b.buf.WriteString("VALUES ")

	numsOfArgs := -1
	for i := 0; i < n; i++ {
		args := argf(i)
		// FIXME what if length of args if different
		if numsOfArgs == -1 {
			numsOfArgs = len(args)
			b.args = make([]any, 0, n*numsOfArgs)
		}
		b.args = append(b.args, args...)
	}

	qs := QuestionMarks(numsOfArgs)
	for i := 0; i < n; i++ {
		b.buf.WriteString(qs)
		if i+1 < n {
			b.buf.WriteByte(',')
		}
	}

	return (*insertBuilderValues)(b)
}

func (b *insertBuilderValues) Build() (string, []any) {
	return (*sqlBuilderBuild)(b).Build()
}
