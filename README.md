# sqlbuilder

[![Coverage Status](https://coveralls.io/repos/github/llklkl/sqlbuilder/badge.svg?branch=main)](https://coveralls.io/github/llklkl/sqlbuilder?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/llklkl/sqlbuilder)](https://goreportcard.com/report/github.com/llklkl/sqlbuilder)

[简体中文](README_zh.md)

A DML SQL simple statement construction tool, which supports chained calls.
Supports generating `SELECT`, `UPDATE`, `DELETE` and `INSERT` simple statements.

**hint**:

+ **Do not cache any intermediate results in the chain call process, as this will cause errors in the final generated
  SQL statement**
+ `SqlBuilder` is not thread-safe and cannot be called concurrently

## Contents

<!-- TOC -->
* [sqlbuilder](#sqlbuilder)
  * [Contents](#contents)
  * [Install](#install)
  * [Usage](#usage)
    * [insert statement](#insert-statement)
      * [insert a single piece of data](#insert-a-single-piece-of-data)
      * [insert multiple pieces of data](#insert-multiple-pieces-of-data)
    * [select statement](#select-statement)
    * [select statement](#select-statement-1)
    * [delete statement](#delete-statement)
    * [construct where condition](#construct-where-condition)
  * [Some special functions](#some-special-functions)
    * [func T(args ...string) *Table](#func-targs-string-table)
    * [func F(args ...string) *Field](#func-fargs-string-field)
    * [func E(args ...string) *Expr](#func-eargs-string-expr)
    * [func O(field any, direction OrderDirection) *OrderSpec](#func-ofield-any-direction-orderdirection-orderspec)
  * [License](#license)
<!-- TOC -->

## Install

```shell
go get github.com/llklkl/sqlbuilder@latest
```

## Usage

### insert statement

#### insert a single piece of data

the original sql:

```sql
INSERT INTO `demo` (`name`, `age`)
VALUES (?, ?)
```

construct using `sqlbuilder`:

```go
package main

import (
	"fmt"

	sb "github.com/llklkl/sqlbuilder"
)

func main() {
	sql, args := sb.New().Insert().Into("demo").
		Fields("name", "age").
		Values("alice", 20).Build()
	fmt.Println(sql)
	fmt.Println(args)
}
```

#### insert multiple pieces of data

the original sql:

```sql
INSERT INTO `demo` (`name`, `age`)
VALUES (?, ?),
       (?, ?),
       (?, ?) ON DUPLICATE KEY
UPDATE `name`=?,`age`=`age`+1
```

construct using `sqlbuilder`:

```go
package main

import (
	"fmt"

	sb "github.com/llklkl/sqlbuilder"
)

type Student struct {
	Name string
	Age  int
}

func main() {
	students := []*Student{
		{Name: "alice", Age: 19},
		{Name: "bob", Age: 20},
		{Name: "carol", Age: 21},
	}
	sql, args := sb.New().Insert().Into("demo").
		Fields("name", "age").
		Bulk(len(students), func(index int) []any {
			return []any{students[index].Name, students[index].Age}
		}).
		OnDuplicate(
			sb.Set(sb.F("name"), "duplicate"),
			sb.Value("`age`=`age`+1"),
		).Build()
	fmt.Println(sql)
	fmt.Println(args)
}

```

### select statement

the original sql:

```sql
SELECT `c`.`class_name`, `s`.`name`, `s`.`score`
FROM `t_student` AS `s`
         RIGHT JOIN `t_class` AS `c` USING (`class_id`)
WHERE `c`.`class_name` = ?
  AND `s`.`score` >= ?
ORDER BY `s`.`name` ASC LIMIT ?,?
```

construct using `sqlbuilder`:

```go
package main

import (
	"fmt"

	sb "github.com/llklkl/sqlbuilder"
)

func main() {
	sql, args := sb.New().Select().
		Field(
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
		OrderBy(sb.O(sb.F("s", "name"), sb.Asc)).
		LimitOffset(0, 10).Build()
	fmt.Println(sql)
	fmt.Println(args)
}

```

### select statement

the original sql:

```sql
UPDATE `demo`
SET `name`=?,
    `age`=?
WHERE `name` = ? LIMIT ?
```

construct using `sqlbuilder`:

```go
package main

import (
	"fmt"

	sb "github.com/llklkl/sqlbuilder"
)

func main() {
	sql, args := sb.New().Update().Table("demo").
		Set(
			sb.Set(sb.F("name"), "alice"),
			sb.Set(sb.F("age"), 22),
		).Where(sb.Eq(sb.F("name"), "bob")).
		Limit(5).Build()
	fmt.Println(sql)
	fmt.Println(args)
}

```

### delete statement

the original sql:

```sql
DELETE
FROM `demo`
WHERE `age` >= ? ORDER BY `name` DESC LIMIT ?
```

construct using `sqlbuilder`:

```go
package main

import (
	"fmt"

	sb "github.com/llklkl/sqlbuilder"
)

func main() {
	sql, args := sb.New().Delete().From("demo").
		Where(sb.Ge(sb.F("age"), 20)).
		Order(sb.O(sb.F("name"), sb.Desc)).
		Limit(10).Build()
	fmt.Println(sql)
	fmt.Println(args)
}

```

### construct where condition

The `Where` method in the `SELECT` statement connects multiple conditions in an `AND` manner by default. Multiple
conditions can be nested through the `And`, `Or` methods.

Currently, the following where conditions are supported:

+ And: Multiple where conditions can be nested and connected with `AND`
+ Or: Multiple where conditions can be nested and connected using `OR`
+ Lt
+ Le
+ Eq
+ Gt
+ Ge
+ Ne
+ Between And
+ Like
+ IsNull
+ NotNull
+ In
+ Not In
+ Exists: supports subquery statement
+ Not Exists: supports subquery statement
+ Condition: supports customizing arbitrary conditions. For example, `Condition("file_sha=UNHEX(?)", fileSha)` defines a
  condition of `file_sha=UNHEX(?)`
+ ...

## Some special functions

### func T(args ...string) *Table

This function defines a `Table`, which is used to define the table in SQL statements.

This function will interpret the parameters differently depending on the number of parameters:

+ When the number of parameters is 1, it is equivalent to `func (table string) *Table`
+ When the number of parameters is 2, it is equivalent to `func (table, alias string) *Table`
+ When the number of parameters is 3, it is equivalent to `func (database, table, alias string) *Table`

### func F(args ...string) *Field

This function defines a `Field`, usually used for conditional filtering, or `SELECT` query fields.

This function will interpret the parameters differently depending on the number of parameters:

+ When the number of parameters is 1, it is equivalent to `func (field string) *Field`
+ When the number of parameters is 2, it is equivalent to `func (table, field string) *Field`
+ When the number of parameters is 3, it is equivalent to `func (table, field, alias string) *Field`

### func E(args ...string) *Expr

This function defines an `Expr`, which is used when the `SELECT` statement requires a built-in function expression.

This function will interpret the parameters differently depending on the number of parameters:

+ When the number of parameters is 1, it is equivalent to `func (expr string) *Expr`
+ When the number of parameters is 2, it is equivalent to `func (expr, alias string) *Expr`

### func O(field any, direction OrderDirection) *OrderSpec

This function defines an OrderSpec for Select...Order to specify the sorting field.

## License

[MIT](https://github.com/sunyctf/ChineseREADME/blob/main/LICENSE) © llklkl

