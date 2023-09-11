# sqlbuilder

一个支持链式调用的 DML SQL 简单语句构造工具。
支持生成 `SELECT`, `UPDATE`, `DELETE` 和 `INSERT` 简单的语句。

提示:

+ **不要缓存链式调用过程中的任何中间结果，这样会导致最终生成 SQL 语句出错**
+ `SqlBuilder` 不是线程安全的，不能并发调用

## 目录

<!-- TOC -->
* [sqlbuilder](#sqlbuilder)
  * [目录](#目录)
  * [安装](#安装)
  * [使用方法](#使用方法)
    * [insert 语句](#insert-语句)
      * [插入单条数据](#插入单条数据)
      * [插入多条数据](#插入多条数据)
    * [select 语句](#select-语句)
    * [update 语句](#update-语句)
    * [delete 语句](#delete-语句)
    * [构造 where 条件](#构造-where-条件)
  * [一些特殊函数](#一些特殊函数)
    * [func T(args ...string) *Table](#func-targs-string-table)
    * [func F(args ...string) *Field](#func-fargs-string-field)
    * [func E(args ...string) *Expr](#func-eargs-string-expr)
  * [开源协议](#开源协议)
<!-- TOC -->

## 安装

```shell
go get github.com/llklkl/sqlbuilder@latest
```

## 使用方法

### insert 语句

#### 插入单条数据

原始sql：

```sql
INSERT INTO `demo` (`name`, `age`)
VALUES (?, ?)
```

使用 `sqlbuilder` 构造：

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

#### 插入多条数据
原始sql：

```sql
INSERT INTO `demo` (`name`, `age`)
VALUES (?, ?),
       (?, ?),
       (?, ?) ON DUPLICATE KEY
UPDATE `name`=?,`age`=`age`+1
```

使用 `sqlbuilder` 构造：

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

### select 语句

原始sql：

```sql
SELECT `c`.`class_name`, `s`.`name`, `s`.`score`
FROM `t_student` AS `s`
         RIGHT JOIN `t_class` AS `c` USING (`class_id`)
WHERE `c`.`class_name` = ?
  AND `s`.`score` >= ?
ORDER BY `s`.`name` ASC LIMIT ?,?
```

使用 `sqlbuilder` 构造：

```go
package main

import (
	"fmt"

	sb "github.com/llklkl/sqlbuilder"
)

func main() {
	sql, args := sb.New().Select().
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
	fmt.Println(sql)
	fmt.Println(args)
}

```

### update 语句

原始sql：

```sql
UPDATE `demo`
SET `name`=?,
    `age`=?
WHERE `name` = ? LIMIT ?
```

使用 `sqlbuilder` 构造：

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

### delete 语句

原始sql：

```sql
DELETE
FROM `demo`
WHERE `age` >= ? ORDER BY `name` DESC LIMIT ?
```

使用 `sqlbuilder` 构造：

```go
package main

import (
	"fmt"

	sb "github.com/llklkl/sqlbuilder"
)

func main() {
	sql, args := sb.New().Delete().From("demo").
		Where(sb.Ge(sb.F("age"), 20)).
		Order(sb.Order(sb.F("name"), sb.Desc)).
		Limit(10).Build()
	fmt.Println(sql)
	fmt.Println(args)
}

```

### 构造 where 条件

`SELECT` 语句中的 `Where` 方法默认以 `AND` 的方式连接多个条件。多个条件可以通过 `And`, `Or` 方法嵌套。

目前支持构造以下的 where 条件：

+ And: 可以嵌套多个 where 条件，并用 `AND` 连接
+ Or: 可以嵌套多个 where 条件，并用 `OR` 连接
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
+ Exists: 支持加入一个条子查询语句
+ Not Exists： 支持加入一条子查询语句
+ Condition: 支持自定义任意条件。如，`Condition("file_sha=UNHEX(?)", fileSha)`定义一个`file_sha=UNHEX(?)`的条件
+ ...

## 一些特殊函数

### func T(args ...string) *Table

该函数用于定义一个 `Table`，用于 SQL 语句中定义表。

该函数会根据不同的参数个数，对参数进行不同解释：

+ 参数个数为 1 时，等价于 `func (table string) *Table`
+ 参数个数为 2 时，等价于 `func (table, alias string) *Table`
+ 参数个数为 3 时，等价于 `func (database, table, alias string) *Table`

### func F(args ...string) *Field

该函数用于定义一个 `Field`，通常用于条件过滤，或者 `SELECT` 查询字段。

该函数会根据不同的参数个数，对参数进行不同的解释：

+ 参数个数为 1 时, 等价于 `func (field string) *Field`
+ 参数个数为 2 时, 等价于 `func (table, field string) *Field`
+ 参数个数为 3 时, 等价于 `func (table, field, alias string) *Field`

### func E(args ...string) *Expr

该函数用于定义一个 `Expr`，在 `SELECT` 语句需要内置函数表达式时使用。

该函数会根据不同的参数个数，对参数进行不同的解释：

+ 参数个数为 1 时, 等价于 `func (expr string) *Expr`
+ 参数个数为 2 时, 等价于 `func (expr, alias string) *Expr`

## 开源协议

[MIT](https://github.com/sunyctf/ChineseREADME/blob/main/LICENSE) © llklkl