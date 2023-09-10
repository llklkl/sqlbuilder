package tests

import (
	"reflect"
	"testing"

	"github.com/pingcap/tidb/parser"
	_ "github.com/pingcap/tidb/parser/test_driver"

	sb "github.com/llklkl/sqlbuilder"
)

func sqlCheck(t *testing.T, sql string) {
	parse := parser.New()
	_, _, err := parse.Parse(sql, "", "")
	if err != nil {
		t.Errorf("bad sql [%s], err=[%v]", sql, err)
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
				return sb.New().Insert().Into("demo").
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
				return sb.New().Insert().Into("demo").
					Fields("name", "age").
					ValuesFn(len(ages), func(index int) []any {
						return []any{names[index], ages[index]}
					}).Build()
			},
			wantSql:  "INSERT INTO `demo` (`name`,`age`) VALUES (?,?),(?,?),(?,?)",
			wantArgs: []any{"alice", 19, "bob", 20, "carol", 21},
		},
		{
			name: "insert, duplicate",
			workFn: func() (string, []any) {
				ages := []int{19, 20, 21}
				names := []string{"alice", "bob", "carol"}
				return sb.New().Insert().Into("demo").
					Fields("name", "age").
					ValuesFn(len(ages), func(index int) []any {
						return []any{names[index], ages[index]}
					}).
					OnDuplicate(
						sb.Set(sb.F("name"), "duplicate"),
						sb.Value("`age`=`age`+1"),
					).Build()
			},
			wantSql:  "INSERT INTO `demo` (`name`,`age`) VALUES (?,?),(?,?),(?,?) ON DUPLICATE KEY UPDATE `name`=?,`age`=`age`+1",
			wantArgs: []any{"alice", 19, "bob", 20, "carol", 21, "duplicate"},
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
				return sb.New().Select().
					FieldT(
						sb.F("c", "class_name"),
						sb.F("s", "name"),
						sb.F("s", "score"),
					).
					FromT(sb.T("t_student", "s")).
					RightJoin(sb.T("t_class", "c")).Using("class_id").
					Where(
						sb.Eq(sb.F("c", "class_name"), "class1"),
						sb.Ge(sb.F("s", "score"), 85),
					).
					OrderBy(sb.Order(sb.F("s", "name"), sb.Asc)).
					LimitOffset(0, 10).Build()
			},
			wantSql:  "SELECT `c`.`class_name`,`s`.`name`,`s`.`score` FROM `t_student` AS `s` RIGHT JOIN `t_class` AS `c` USING (`class_id`) WHERE `c`.`class_name` = ? AND `s`.`score` >= ? ORDER BY `s`.`name` ASC LIMIT ?,?",
			wantArgs: []any{"class1", 85, 10, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := tt.workFn()
			if sql != tt.wantSql {
				t.Errorf("Select sql got = %v, want %v", sql, tt.wantSql)
			}
			if !reflect.DeepEqual(args, tt.wantArgs) {
				t.Errorf("SELECT args got1 = %v, want %v", args, tt.wantArgs)
			}
			sqlCheck(t, sql)
		})
	}
}

func TestSqlBuilder_Delete(t *testing.T) {
	tests := []struct {
		name     string
		workFn   func() (string, []any)
		wantSql  string
		wantArgs []any
	}{
		{
			name: "DELETE",
			workFn: func() (string, []any) {
				return sb.New().Delete().From("demo").
					Where(sb.Ge(sb.F("age"), 20)).
					Order(sb.Order(sb.F("name"), sb.Desc)).
					Limit(10).Build()
			},
			wantSql:  "DELETE FROM `demo` WHERE `age` >= ? ORDER BY `name` DESC LIMIT ?",
			wantArgs: []any{20, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := tt.workFn()
			if sql != tt.wantSql {
				t.Errorf("Delete sql got = %v, want %v", sql, tt.wantSql)
			}
			if !reflect.DeepEqual(args, tt.wantArgs) {
				t.Errorf("Delete args got1 = %v, want %v", args, tt.wantArgs)
			}
			sqlCheck(t, sql)
		})
	}
}

func TestSqlBuilder_Update(t *testing.T) {
	tests := []struct {
		name     string
		workFn   func() (string, []any)
		wantSql  string
		wantArgs []any
	}{
		{
			name: "Update",
			workFn: func() (string, []any) {
				return sb.New().Update().Table("demo").
					Set(
						sb.Set(sb.F("name"), "alice"),
						sb.Set(sb.F("age"), 22),
					).Where(sb.Eq(sb.F("name"), "bob")).
					Limit(5).Build()
			},
			wantSql:  "UPDATE `demo` SET `name`=?,`age`=? WHERE `name` = ? LIMIT ?",
			wantArgs: []any{"alice", 22, "bob", 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := tt.workFn()
			if sql != tt.wantSql {
				t.Errorf("Update sql got = %v, want %v", sql, tt.wantSql)
			}
			if !reflect.DeepEqual(args, tt.wantArgs) {
				t.Errorf("Update args got1 = %v, want %v", args, tt.wantArgs)
			}
			sqlCheck(t, sql)
		})
	}
}
