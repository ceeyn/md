



### 一个订单表，有order_id,user_id,order_price，需要查询出消费总金额第二大的所有用户，用户可能有多个

```sql
SELECT user_id, total_amount
FROM (
    SELECT user_id, SUM(order_price) AS total_amount
    FROM orders
    GROUP BY user_id
) AS user_totals
WHERE total_amount = (
    SELECT MAX(total_amount) 
    FROM (
        SELECT SUM(order_price) AS total_amount
        FROM orders
        GROUP BY user_id
        HAVING SUM(order_price) < (
            SELECT MAX(total_amount)
            FROM (
                SELECT SUM(order_price) AS total_amount
                FROM orders
                GROUP BY user_id
            ) AS temp1
        )
    ) AS temp2
);

```

### 解释：

1. **第一层子查询 (`temp1`)**：计算每个用户的总消费金额，并找出最高的消费金额。
2. **第二层子查询 (`temp2`)**：找出总消费金额小于最高消费金额的用户中的最大值（即第二高消费金额）。
3. **主查询 (`user_totals`)**：找出总消费金额等于第二高消费金额的所有用户。

这样，查询结果将返回所有消费总金额第二大的用户及其总金额。



### With

`WITH` 是 SQL 中的一个关键字，用于定义 **CTE（Common Table Expression）**，也叫**公用表表达式**。CTE 允许你在查询中定义一个临时结果集，供后续的查询使用。`WITH` 后面跟的是 CTE 的名称（在你的例子中是 `user_totals`），并且 `AS` 后跟的是一个子查询。

简单来说，`WITH user_totals AS` 是在定义一个临时的查询结果表，叫做 `user_totals`，这个表只在当前的 SQL 语句中有效。

````
`WITH` 是 SQL 中的一个关键字，用于定义 **CTE（Common Table Expression）**，也叫**公用表表达式**。CTE 允许你在查询中定义一个临时结果集，供后续的查询使用。`WITH` 后面跟的是 CTE 的名称（在你的例子中是 `user_totals`），并且 `AS` 后跟的是一个子查询。

简单来说，`WITH user_totals AS` 是在定义一个临时的查询结果表，叫做 `user_totals`，这个表只在当前的 SQL 语句中有效。

### 详细解释：
```sql
WITH user_totals AS (
    SELECT user_id, SUM(order_price) AS total_amount
    FROM orders
    GROUP BY user_id
)
```

- **`WITH`**：引入一个 CTE。
- **`user_totals`**：这是 CTE 的名字，表示在接下来的查询中，你可以把 `user_totals` 当作一个表来使用。
- **`AS`**：后面跟的是一个子查询，该子查询的结果会被存储在 `user_totals` 这个临时表中。

在这个例子中，`user_totals` 是一个包含 `user_id` 和每个用户的 `total_amount`（消费总金额）的临时表。这个临时表接下来可以在其他查询中直接引用。

### 这个 CTE 相当于：
你可以把它理解成把这个查询的结果先存储在一个虚拟表 `user_totals` 里，后续的查询再基于这个虚拟表来做进一步的处理。比如在这个例子中，后续的查询使用 `ranked_totals` 来为每个用户的消费总金额进行排名。

CTE 的好处是：
1. **简洁性**：可以让复杂查询变得更清晰，逻辑分明。
2. **可读性**：便于理解各个步骤的含义，而不是嵌套多个子查询。
3. **重用**：在主查询中可以重用临时结果集，而无需重复编写相同的查询逻辑。

如果不使用 `WITH` 关键字，等效的查询会变得冗长复杂。
````





