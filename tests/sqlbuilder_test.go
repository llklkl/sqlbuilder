package tests

import (
	"reflect"
	"testing"

	"github.com/pingcap/tidb/parser"
	_ "github.com/pingcap/tidb/parser/test_driver"

	. "github.com/llklkl/sqlbuilder"
)

func sqlCheck(t *testing.T, sql string) {
	parse := parser.New()
	_, _, err := parse.Parse(sql, "", "")
	if err != nil {
		t.Errorf("bad sql [%s] err=%v", sql, err)
	}
}

func TestSqlBuilder_Insert(t *testing.T) {
	tests := []struct {
		name     string
		workFn   func() (string, []any)
		wantSql  string
		wantArgs []any
	}{
		{
			name: "insert",
			workFn: func() (string, []any) {
				return New().Insert().Into("demo").
					Fields("name", "age").
					Values("alice", 20).Build()
			},
			wantSql:  "INSERT INTO `demo` (`name`,`age`) VALUES (?,?)",
			wantArgs: []any{"alice", 20},
		},
		{
			name: "bulk insert",
			workFn: func() (string, []any) {
				ages := []int{19, 20, 21}
				names := []string{"alice", "bob", "carol"}
				return New().Insert().Into("demo").
					Fields("name", "age").
					ValuesFn(len(ages), func(index int) []any {
						return []any{names[index], ages[index]}
					}).Build()
			},
			wantSql:  "INSERT INTO `demo` (`name`,`age`) VALUES (?,?),(?,?),(?,?)",
			wantArgs: []any{"alice", 19, "bob", 20, "carol", 21},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := tt.workFn()
			if sql != tt.wantSql {
				t.Errorf("Insert sql got = %v, want %v", sql, tt.wantSql)
			}
			if !reflect.DeepEqual(args, tt.wantArgs) {
				t.Errorf("Insert args got1 = %v, want %v", args, tt.wantArgs)
			}
			sqlCheck(t, sql)
		})
	}
}

func TestSqlBuilder_Select(t *testing.T) {
	tests := []struct {
		name     string
		workFn   func() (string, []any)
		wantSql  string
		wantArgs []any
	}{
		{
			name: "select",
			workFn: func() (string, []any) {
				return New().Select().Expr(&Expr{Expr: "name", Table: "d"}, &Expr{Expr: "age", Table: "d1"}).
					FromT(&Table{Table: "demo", Alias: "d"}).
					LeftJoin(&Table{Table: "demo1", Alias: "d1"}).On(&Field{Field: "name", Table: "d"}, &Field{Field: "name", Table: "d1"}).
					RightJoin(&Table{Table: "demo2"}).Using("age").
					Where(Condition("name=?", 1)).
					GroupBy("name").
					OrderBy(Order(&Field{Field: "age"}, Asc)).
					Limit(1, 2).Build()
			},
			wantSql:  "SELECT `d`.name,`d1`.age FROM `demo` AS `d` LEFT JOIN `demo1` AS `d1` ON `d`.`name`=`d1`.`name` RIGHT JOIN `demo2` USING (`age`) WHERE name=? GROUP BY `name` ORDER BY `age` ASC LIMIT ?,?",
			wantArgs: []any{1, 1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := tt.workFn()
			if sql != tt.wantSql {
				t.Errorf("Select sql got = %v, want %v", sql, tt.wantSql)
			}
			if !reflect.DeepEqual(args, tt.wantArgs) {
				t.Errorf("Selete args got1 = %v, want %v", args, tt.wantArgs)
			}
			sqlCheck(t, sql)
		})
	}
}
