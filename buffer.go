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

func (b *buffer) ParenthesesString(val string) {
	b.WriteByte(openParentheses)
	b.WriteString(val)
	b.WriteByte(closeParentheses)
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
		b.Table(t)
		if i+1 < len(tables) {
			b.Comma()
		}
	}
}

func (b *buffer) Expr(e *Expr) {
	if e.Table != "" {
		b.BackQuoteString(e.Table)
		b.Dot()
	}
	b.WriteString(e.Expr)
	if e.Alias != "" {
		b.WriteString(" AS ")
		b.BackQuoteString(e.Alias)
	}
}

func (b *buffer) Exprs(exprs []*Expr) {
	for i, e := range exprs {
		b.Expr(e)
		if i+1 < len(exprs) {
			b.Comma()
		}
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
		b.Field(f)
		if i+1 < len(fields) {
			b.Comma()
		}
	}
}

func (b *buffer) Conditions(conditions []whereCondition) {
	for _, c := range conditions {
		c.write(b)
	}
}

func (b *buffer) OrderSpecs(orderSpecs []*OrderSpec) {
	for i, spec := range orderSpecs {
		b.Field(spec.Field)
		if spec.OrderDirection != "" {
			b.Space()
			b.WriteString(string(spec.OrderDirection))
		}
		if i+1 < len(orderSpecs) {
			b.Comma()
		}
	}
}
