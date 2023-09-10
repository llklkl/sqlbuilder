package sqlbuilder

type SqlBuilder struct {
	buf  *buffer
	args []any
}

func New() *SqlBuilder {
	return &SqlBuilder{
		buf: nil,
	}
}

func (b *SqlBuilder) init() {
	b.buf = getBuffer()
}

func (b *SqlBuilder) Insert(kws ...Keyword) *insertBuilder {
	b.init()
	return (*insertBuilder)(b).init(kws)
}

func (b *SqlBuilder) Select(kws ...Keyword) *selectBuilder {
	b.init()
	return (*selectBuilder)(b).init(kws)
}

func (b *SqlBuilder) Delete(kws ...Keyword) *deleteBuilder {
	b.init()
	return (*deleteBuilder)(b).init(kws)
}

func (b *SqlBuilder) Update(kws ...Keyword) *updateBuilder {
	b.init()
	return (*updateBuilder)(b).init(kws)
}

type sqlBuilderBuild SqlBuilder

func (b *sqlBuilderBuild) Build() (string, []any) {
	sql := b.buf.String()
	releaseBuffer(b.buf)
	return sql, b.args
}
