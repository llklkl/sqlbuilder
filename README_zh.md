# sqlbuilder

一个支持链式调用的 DML SQL 简单语句构造工具。
支持生成 `SELECT`, `UPDATE`, `DELETE` 和 `INSERT` 简单的语句。

对于 `SELETE` 语句，内置了部分常用的条件过滤逻辑，诸如 `Ge`, `Eq`, `Ne` 等，
更多的条件可以通过 `Condition` 的方式自定义过滤逻辑。

提示:
+ **不要缓存链式调用过程中的任何中间结果，这样会导致最终生成 SQL 语句出错**
+ `SqlBuilder` 不是线程安全的，不能并发调用

## 一些特殊函数

> func T(args ...string) *Table

该函数用于定义一个 `Table`，用于 SQL 语句中定义表。

该函数会根据不同的参数个数，对参数进行不同解释：
+ 参数个数为 1 时，等价于 `func (table string) *Table`
+ 参数个数为 2 时，等价于 `func (table, alias string) *Table`
+ 参数个数为 3 时，等价于 `func (database, table, alias string) *Table`

> func F(args ...string) *Field

该函数用于定义一个 `Field`，通常用于条件过滤，或者 `SELETE` 查询字段。

该函数会根据不同的参数个数，对参数进行不同的解释：
+ 参数个数为 1 时, 等价于 `func (field string) *Field`
+ 参数个数为 2 时, 等价于 `func (table, field string) *Field`
+ 参数个数为 3 时, 等价于 `func (table, field, alias string) *Field`

> func E(args ...string) *Expr

该函数用于定义一个 `Expr`，在 `SELETE` 语句需要内置函数表达式时使用。

该函数会根据不同的参数个数，对参数进行不同的解释：
+ 参数个数为 1 时, 等价于 `func (expr string) *Expr`
+ 参数个数为 2 时, 等价于 `func (expr, alias string) *Expr`
