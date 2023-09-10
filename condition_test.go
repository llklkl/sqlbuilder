package sqlbuilder

import (
	"testing"
)

func TestCondition_multiCondition(t *testing.T) {
	tests := []struct {
		name       string
		conditions []whereCondition
		want       string
	}{
		{
			name: "",
			conditions: []whereCondition{
				Lt(F("lt"), 1),
				Le(F("le"), 2),
				Eq(F("eq"), 3),
				Ge(F("ge"), 4),
				Gt(F("gt"), 5),
				Ne(F("ne"), 6),
				Between(F("be"), 7, 8),
				Condition("any_condition=?", 9),
			},
			want: "`lt` < ? AND `le` <= ? AND `eq` = ? AND `ge` >= ? AND `gt` > ? AND `ne` != ? AND `be` BETWEEN ? AND ? AND any_condition=?",
		},
		{
			name: "",
			conditions: []whereCondition{
				In(F("name"), "alice", "bob"),
				NotIn(F("name"), "carol"),
				Exists("select 1"),
				NotExists("select 2"),
			},
			want: "`name` IN (?,?) AND `name` NOT IN (?) AND EXISTS (select 1) AND NOT EXISTS (select 2)",
		},
		{
			name: "",
			conditions: []whereCondition{
				And(
					Ge(F("age"), 20),
					Or(
						Eq(F("name"), "alice"),
						Eq(F("name"), "bob"),
					),
				),
				In(F("name"), "carol"),
				Or(
					Eq(F("name"), "dave"),
					Gt(F("age"), "100"),
				),
			},
			want: "(`age` >= ? AND (`name` = ? OR `name` = ?)) AND `name` IN (?) AND (`name` = ? OR `age` > ?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := newBuffer(1024)
			buf.Conditions(tt.conditions)
			got := buf.String()
			if got != tt.want {
				t.Errorf("multiCondition got=%v, want=%v", got, tt.want)
			}
		})
	}
}
