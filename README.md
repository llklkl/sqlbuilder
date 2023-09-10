# sqlbuilder

A DML SQL simple statement construction tool, which supports chained calls.
Supports generating `SELECT`, `UPDATE`, `DELETE` and `INSERT` simple statements.

For the `SELETE` statement, some commonly used conditional filtering logic is built-in,
such as `Ge`, `Eq`, `Ne`, etc. More conditions can be customized through `Condition`.

**hint**:
+ **Do not cache any intermediate results in the chain call process, as this will cause errors in the final generated SQL statement**
+ `SqlBuilder` is not thread-safe and cannot be called concurrently

## Some special functions

> func T(args ...string) *Table

This function defines a `Table`, which is used to define the table in SQL statements.

This function will interpret the parameters differently depending on the number of parameters:
+ When the number of parameters is 1, it is equivalent to `func (table string) *Table`
+ When the number of parameters is 2, it is equivalent to `func (table, alias string) *Table`
+ When the number of parameters is 3, it is equivalent to `func (database, table, alias string) *Table`

> func F(args ...string) *Field

This function defines a `Field`, usually used for conditional filtering, or `SELETE` query fields.

This function will interpret the parameters differently depending on the number of parameters:
+ When the number of parameters is 1, it is equivalent to `func (field string) *Field`
+ When the number of parameters is 2, it is equivalent to `func (table, field string) *Field`
+ When the number of parameters is 3, it is equivalent to `func (table, field, alias string) *Field`

> func E(args ...string) *Expr

This function defines an `Expr`, which is used when the `SELETE` statement requires a built-in function expression.

This function will interpret the parameters differently depending on the number of parameters:
+ When the number of parameters is 1, it is equivalent to `func (expr string) *Expr`
+ When the number of parameters is 2, it is equivalent to `func (expr, alias string) *Expr`


