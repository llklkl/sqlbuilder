package sqlbuilder

import (
	"bytes"
	"strings"
	"sync"
)

const (
	defaultBufferSize = 1024
	maxBufferSize     = 64 * 1024
)

var (
	bufferPool *sync.Pool
)

func init() {
	bufferPool = &sync.Pool{
		New: func() any {
			return newBuffer(defaultBufferSize)
		},
	}
}

func getBuffer() *buffer {
	return bufferPool.Get().(*buffer)
}

func releaseBuffer(buf *buffer) {
	if buf.Cap() > maxBufferSize {
		return
	}
	buf.Reset()
	bufferPool.Put(buf)
}

type buffer struct {
	*bytes.Buffer
}

func newBuffer(length int) *buffer {
	return &buffer{
		Buffer: bytes.NewBuffer(make([]byte, 0, length)),
	}
}

func (b *buffer) Space() {
	b.WriteByte(space)
}

func (b *buffer) Question() {
	b.WriteByte(questionMark)
}

func (b *buffer) Dot() {
	b.WriteByte(dot)
}

func (b *buffer) Comma() {
	b.WriteByte(comma)
}

func (b *buffer) Equal() {
	b.WriteByte(equalMark)
}

func (b *buffer) OpenParen() {
	b.WriteByte(openParentheses)
}

func (b *buffer) CloseParen() {
	b.WriteByte(closeParentheses)
}

func (b *buffer) BackQuoteString(val string) {
	b.WriteByte(backQuote)
	b.WriteString(val)
	b.WriteByte(backQuote)
}

func (b *buffer) BackQuoteStrings(vals []string) {
	b.WriteByte(backQuote)
	b.WriteString(strings.Join(vals, "`,`"))
	b.WriteByte(backQuote)
}

func (b *buffer) Table(t *Table) {
	if t.Database != "" {
		b.BackQuoteString(t.Database)
		b.Dot()
	}
	b.BackQuoteString(t.Table)
	if t.Alias != "" {
		b.WriteString(" AS ")
		b.BackQuoteString(t.Alias)
	}
}

func (b *buffer) Tables(tables []*Table) {
	for i, t := range tables {
		if i > 0 {
			b.Comma()
		}
		b.Table(t)
	}
}

func (b *buffer) Expr(e *Expr) {
	b.WriteString(e.Expr)
	if e.Alias != "" {
		b.WriteString(" AS ")
		b.BackQuoteString(e.Alias)
	}
}

func (b *buffer) Exprs(exprs []*Expr) {
	for i, e := range exprs {
		if i > 0 {
			b.Comma()
		}
		b.Expr(e)
	}
}

func (b *buffer) Field(f *Field) {
	if f.Table != "" {
		b.BackQuoteString(f.Table)
		b.Dot()
	}
	b.BackQuoteString(f.Field)
	if f.Alias != "" {
		b.WriteString(" AS ")
		b.BackQuoteString(f.Alias)
	}
}

func (b *buffer) Fields(fields []*Field) {
	for i, f := range fields {
		if i > 0 {
			b.Comma()
		}
		b.Field(f)
	}
}

func (b *buffer) SelectField(fields []_selectField) {
	for i := range fields {
		if i > 0 {
			b.Comma()
		}
		switch v := fields[i].(type) {
		case *Field:
			b.Field(v)
		case *Expr:
			b.Expr(v)
		}
	}
}

func (b *buffer) Conditions(conditions []whereCondition) {
	for i := range conditions {
		if i > 0 {
			b.Space()
			b.WriteString("AND")
			b.Space()
		}
		conditions[i].write(b)
	}
}

func (b *buffer) ValueUpdater(vps []valueUpdater) {
	for i := range vps {
		if i > 0 {
			b.Comma()
		}
		vps[i].write(b)
	}
}

func (b *buffer) OrderSpecs(orderSpecs []*OrderSpec) {
	for i, spec := range orderSpecs {
		if i > 0 {
			b.Comma()
		}
		b.Field(spec.Field)
		if spec.OrderDirection != "" {
			b.Space()
			b.WriteString(string(spec.OrderDirection))
		}
	}
}
