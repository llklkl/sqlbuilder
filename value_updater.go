package sqlbuilder

type _valueUpdater interface {
	valueUpdater()
}

type valueUpdater interface {
	_valueUpdater
	args() []any
	write(*buffer)
}

func Set(field *Field, args ...any) *SetValuer {
	return &SetValuer{
		Field: field,
		Args:  args,
	}
}

func Value(expr string, args ...any) *Valuer {
	return &Valuer{
		Expr: expr,
		Args: args,
	}
}

type SetValuer struct {
	_valueUpdater
	Field *Field
	Args  []any
}

func (v *SetValuer) write(buf *buffer) {
	buf.Field(v.Field)
	buf.Equal()
	buf.Question()
}

func (v *SetValuer) args() []any {
	return v.Args
}

type Valuer struct {
	_valueUpdater
	Expr string
	Args []any
}

func (v *Valuer) write(buf *buffer) {
	buf.WriteString(v.Expr)
}

func (v *Valuer) args() []any {
	return v.Args
}
