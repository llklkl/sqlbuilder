package sqlbuilder

import (
	"testing"
)

func Test_buffer_AnyField(t *testing.T) {
	type args struct {
		field any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				"string_field",
			},
			want: "`string_field`",
		},
		{
			name: "",
			args: args{
				F("field"),
			},
			want: "`field`",
		},
		{
			name: "",
			args: args{
				F("table", "field"),
			},
			want: "`table`.`field`",
		},
		{
			name: "",
			args: args{
				F("table", "field", "alias"),
			},
			want: "`table`.`field` AS `alias`",
		},
		{
			name: "",
			args: args{
				E("expr"),
			},
			want: "expr",
		},
		{
			name: "",
			args: args{
				E("expr", "alias"),
			},
			want: "expr AS `alias`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := getBuffer()
			b.AnyField(tt.args.field)
			got := b.String()
			if got != tt.want {
				t.Errorf("AnyField got=%v, want=%v", got, tt.want)
			}
			releaseBuffer(b)
		})
	}
}

func Test_buffer_AnyFields(t *testing.T) {
	type args struct {
		field []any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				[]any{"string_field", F("field"), E("expr")},
			},
			want: "`string_field`,`field`,expr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := getBuffer()
			b.AnyFields(tt.args.field)
			got := b.String()
			if got != tt.want {
				t.Errorf("AnyField got=%v, want=%v", got, tt.want)
			}
			releaseBuffer(b)
		})
	}
}

func Test_buffer_BackQuoteString(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				"",
			},
			want: "``",
		},
		{
			name: "",
			args: args{
				"name",
			},
			want: "`name`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := getBuffer()
			b.BackQuoteString(tt.args.val)
			got := b.String()
			if got != tt.want {
				t.Errorf("BackQuoteString got=%v, want=%v", got, tt.want)
			}
			releaseBuffer(b)
		})
	}
}

func Test_buffer_BackQuoteStrings(t *testing.T) {
	type args struct {
		vals []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				[]string{"name", "name2"},
			},
			want: "`name`,`name2`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := getBuffer()
			b.BackQuoteStrings(tt.args.vals)
			got := b.String()
			if got != tt.want {
				t.Errorf("BackQuoteStrings got=%v, want=%v", got, tt.want)
			}
			releaseBuffer(b)
		})
	}
}

func Test_buffer_OrderSpecs(t *testing.T) {
	type args struct {
		orderSpecs []*OrderSpec
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				[]*OrderSpec{
					O("name", Asc),
				},
			},
			want: "`name` ASC",
		},
		{
			name: "",
			args: args{
				[]*OrderSpec{
					O("name", Asc),
					O("age", Desc),
				},
			},
			want: "`name` ASC,`age` DESC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := getBuffer()
			b.OrderSpecs(tt.args.orderSpecs)
			got := b.String()
			if got != tt.want {
				t.Errorf("OrderSpecs got=%v, want=%v", got, tt.want)
			}
			releaseBuffer(b)
		})
	}
}

func Test_buffer_Table(t *testing.T) {
	type args struct {
		t *Table
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				t: T("table"),
			},
			want: "`table`",
		},
		{
			name: "",
			args: args{
				t: T("table", "alias"),
			},
			want: "`table` AS `alias`",
		},
		{
			name: "",
			args: args{
				t: T("database", "table", "alias"),
			},
			want: "`database`.`table` AS `alias`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := getBuffer()
			b.Table(tt.args.t)
			got := b.String()
			if got != tt.want {
				t.Errorf("Table got=%v, want=%v", got, tt.want)
			}
			releaseBuffer(b)
		})
	}
}

func Test_buffer_ValueUpdater(t *testing.T) {
	type args struct {
		vps []valueUpdater
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				vps: []valueUpdater{
					Set("field", "new_value"),
					Value("a=1"),
				},
			},
			want: "`field`=?,a=1",
		},
		{
			name: "",
			args: args{
				vps: []valueUpdater{
					Value("a=1"),
				},
			},
			want: "a=1",
		},
		{
			name: "",
			args: args{
				vps: []valueUpdater{
					Set("field", "new_value"),
				},
			},
			want: "`field`=?",
		},
		{
			name: "",
			args: args{
				vps: []valueUpdater{
					Set("field", "new_value"),
					Set("field2", "new_value"),
				},
			},
			want: "`field`=?,`field2`=?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := getBuffer()
			b.ValueUpdater(tt.args.vps)
			got := b.String()
			if got != tt.want {
				t.Errorf("ValueUpdater got=%v, want=%v", got, tt.want)
			}
			releaseBuffer(b)
		})
	}
}
