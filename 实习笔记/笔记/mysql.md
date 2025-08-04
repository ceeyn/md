

# 日常开发常见MySQL性能优化策略及应用场景



## 1. 索引优化

**策略：**

- 为经常查询的列创建索引。
- 避免在索引列上使用函数，因为这会导致索引失效。
- 定期检查索引的效率，使用 `EXPLAIN` 命令分析查询。

**适用场景：**

- 数据库查询响应时间长。
- 需要快速检索大量数据。

### 真实场景应用示例

有一个电子商务平台的数据库，其中有一个 `orders` 表，存储了所有的订单信息。表结构可能如下：

- `order_id` (主键)
- `user_id` (用户ID)
- `order_date` (订单日期)
- `status` (订单状态)
- `total_amount` (订单总额)



**场景描述**

在该电子商务平台的日常运营中，经常需要根据用户的订单状态和日期进行查询，以生成报表或进行订单状态更新。但是，随着订单量的增加，这些查询的响应时间越来越长，影响了报表生成和订单处理的效率。



**优化措施**

1. **创建索引**：为 `status` 和 `order_date` 列创建索引，因为这些列经常作为查询条件。

```sql
    ALTER TABLE orders ADD INDEX idx_status (status);
    ALTER TABLE orders ADD INDEX idx_order_date (order_date);
```

1. **使用 `EXPLAIN` 分析**：通过 `EXPLAIN` 命令分析查询，确认索引是否被有效使用。

```sql
    EXPLAIN SELECT * FROM orders WHERE status = 'Shipped' AND order_date >= '2024-01-01';
```

1. **考虑复合索引**：如果查询经常同时基于 `status` 和 `order_date`，考虑创建一个复合索引。

```sql
    ALTER TABLE orders ADD INDEX idx_status_order_date (status, order_date);
```

1. **索引维护**：随着订单的不断增加和状态的变更，定期使用 `OPTIMIZE TABLE` 命令来优化表和索引，
   该命令的作用是整理和优化数据库表和索引，减少数据碎片，更新索引统计信息，回收未使用空间，并可能重新排序行以提高存储效率和查询性能。



```sql
    OPTIMIZE TABLE orders;
```



## 2. 查询优化

**策略：**

- 避免使用 `SELECT *`，只选择需要的列。
- 使用合适的 `JOIN` 类型和顺序。
- 减少子查询和复杂的嵌套查询。

**适用场景：**

- 查询结果集过大。
- 查询逻辑复杂，难以优化。



### 场景描述

运营一个视频分享网站，类似于YouTube，用户可以上传、观看、评价和分享视频。随着用户基数和视频内容的快速增长，视频推荐系统的查询效率成为了一个关键问题。用户在浏览网站时，希望快速获得个性化的视频推荐列表，而慢速的查询响应会影响用户体验。

### 问题表现

- 当用户请求个性化推荐时，系统需要从数百万甚至上亿的视频内容中筛选，但当前的查询逻辑导致响应时间过长。
- 视频内容更新频繁，用户评分和观看历史需要实时反映在推荐结果中，而现有的数据库查询无法满足实时性要求。

### 查询优化策略

1. **避免使用 `SELECT \*`**：

   - 原始查询可能如下，尝试获取所有视频的详细信息：

     ```sql
     SELECT * FROM videos WHERE video_id IN (推荐算法生成的视频ID列表);
     ```

     

   - 优化后的查询只选择必要的列，比如视频ID、标题和缩略图URL：

     ```sql
     SELECT video_id, title, thumbnail_url FROM videos WHERE video_id IN (推荐算法生成的视频ID列表);
     ```

     

2. **使用合适的 `JOIN` 类型和顺序**：

   - 如果推荐系统需要结合用户评分信息，原始的查询可能如下：

     ```sql
     SELECT * FROM videos JOIN ratings ON videos.video_id = ratings.video_id WHERE user_id = 用户ID;
     ```

     

   - 优化后的查询使用明确的连接类型，并只选择对推荐有用的列：

     ```sql
     SELECT videos.video_id, videos.title, AVG(ratings.score) as average_rating 
     FROM videos 
     LEFT JOIN ratings ON videos.video_id = ratings.video_id AND ratings.user_id = 用户ID 
     GROUP BY videos.video_id 
     HAVING COUNT(ratings.video_id) > 一定数量; -- 确保推荐的视频有足够的评分数据支撑
     ```

     

3. **减少子查询和复杂的嵌套查询**：

   - 假设推荐系统需要找出用户最近观看过的视频，原始查询可能包含子查询：

     ```sql
     SELECT * FROM videos WHERE video_id = (SELECT video_id FROM观看历史 WHERE user_id = 用户ID ORDER BY watch_date DESC LIMIT 1);
     ```

     

   - 优化后的查询使用连接代替子查询，提高效率：

     ```sql
     SELECT videos.* FROM videos 
     INNER JOIN 观看历史 ON videos.video_id = 观看历史.video_id 
     WHERE 观看历史.user_id = 用户ID 
     ORDER BY 观看历史.watch_date DESC LIMIT 1;
     ```

     



## 3. 数据库规范化

**策略：**

- 根据数据的逻辑关系进行规范化，减少数据冗余。
- 合理设计表结构，避免过多的表连接。

**适用场景：**

- 数据更新频繁，需要保持数据一致性。
- 数据库规模较大，需要减少数据冗余。



### 场景描述

考虑一个快速发展的外卖平台，类似于Uber Eats或Deliveroo，该平台需要处理大量的顾客订单、菜单项更新和送餐状态跟踪。随着业务的扩展，数据库的规模迅速增长，数据规范化成为提高效率和减少冗余的关键。

**问题表现**

- 顾客信息（如姓名、地址、电话号码）频繁更新。
- 餐厅菜单项（包括价格和菜品描述）经常变动。
- 订单数据量巨大，包含订单详情、状态和送餐信息。

**数据库规范化策略**

1. **分离实体数据**：
   - 将顾客信息存储在`customers`表中，字段包括`customer_id`, `name`, `address`, `phone_number`。
2. **分离菜单和订单数据**：
   - 餐厅菜单项存储在`menu_items`表中，字段包括`item_id`, `restaurant_id`, `item_name`, `price`。
   - 订单详情存储在`orders`表中，字段包括`order_id`, `customer_id`, `item_id`, `quantity`, `order_status`。
3. **设计事实表**：
   - 为订单和送餐状态创建事实表`order_details`，字段包括`order_id`, `restaurant_id`, `delivery_status`, `order_time`。



```sql
-- 创建顾客信息表
CREATE TABLE customers (
    customer_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    address VARCHAR(255),
    phone_number VARCHAR(20)
);

-- 创建菜单项表
CREATE TABLE menu_items (
    item_id INT AUTO_INCREMENT PRIMARY KEY,
    restaurant_id INT,
    item_name VARCHAR(100),
    price DECIMAL(10, 2)
);

-- 创建订单表
CREATE TABLE orders (
    order_id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT,
    item_id INT,
    quantity INT,
    order_status VARCHAR(50),
    order_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建订单详情表
CREATE TABLE order_details (
    order_id INT,
    restaurant_id INT,
    delivery_status VARCHAR(50),
    delivery_time TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(order_id),
    FOREIGN KEY (restaurant_id) REFERENCES menu_items(restaurant_id)
);
```



## 4. 缓存策略

**策略：**

- 使用应用层缓存或数据库内建的缓存机制。
- 对于不常变更的数据，使用缓存减少数据库访问。

**适用场景：**

- 数据读取频繁但更新不频繁。
- 需要减轻数据库的读取压力。



### 场景描述

一个社交媒体平台，如Twitter或Instagram，用户频繁查看和参与热门话题的讨论。热门话题通常包括最新的新闻、流行文化事件或用户高度关注的讨论。由于这些话题的页面被频繁访问，数据库需要处理大量的相同查询请求，这可能导致数据库读取压力增大，影响用户体验。

**问题表现**

- 用户访问热门话题页面时，数据库需要快速响应，提供最新的讨论和帖子。
- 【热门话题的更新频率相对较低，但访问频率非常高。】

**缓存策略**

1. **应用层缓存**：
   - 使用内存缓存系统（如Redis）来存储热门话题的帖子和讨论数据。
2. **缓存数据选择**：
   - 只缓存访问频率高且更新频率低的数据，如热门话题的帖子列表。
3. **缓存失效策略**：
   - 设定合理的缓存过期时间，或在后端服务中监听数据变更事件，以更新或失效缓存。



```go
// getPopularTopicPosts 尝试从缓存获取热门话题的帖子列表，如果缓存未命中，则从数据库获取并更新缓存
func getPopularTopicPosts(topicID string) ([]string, error) {
	// 构建缓存键名
	postsKey := "popular_topic_" + topicID + "_posts"

	// 尝试从Redis缓存中获取帖子列表
	postsJSON, err := cache.Get(ctx, postsKey).Result()
	if err == nil {
		// 缓存命中，反序列化JSON数据到切片
		var posts []string
		err = json.Unmarshal([]byte(postsJSON), &posts)
		if err == nil {
			fmt.Println("从缓存中获取热门话题帖子列表")
			return posts, nil // 返回帖子列表
		}
		// 如果反序列化失败，打印错误并返回空列表
		fmt.Println("反序列化缓存数据失败:", err)
		return nil, err
	}

	// 缓存未命中，模拟从数据库获取数据
	fmt.Println("缓存未命中，从数据库中获取热门话题帖子列表")
	// 这里应是数据库查询逻辑，此处使用模拟数据代替
	posts := []string{"帖子1", "帖子2", "帖子3"} // 假设的数据库查询结果

	// 序列化帖子列表为JSON字符串
	postsJSON, err = json.Marshal(posts)
	if err != nil {
		// 如果序列化失败，打印错误并返回错误
		fmt.Println("序列化帖子列表失败:", err)
		return nil, err
	}

	// 将序列化后的帖子列表存入缓存，并设置1小时的过期时间
	err = cache.SetEX(ctx, postsKey, time.Hour, string(postsJSON)).Err()
	if err != nil {
		// 如果缓存设置失败，打印错误并返回错误
		fmt.Println("设置缓存数据失败:", err)
		return nil, err
	}

	// 返回帖子列表
	return posts, nil
}
```



## 5. 并发控制

**策略：**

- 使用合适的事务隔离级别。
- 避免长事务，减少锁的竞争。

**适用场景：**

- 多用户环境下数据库操作冲突。
- 需要保证数据的一致性和完整性。



### 场景描述

考虑一个在线票务系统，类似于Eventbrite或Ticketmaster，用户可以浏览即将举行的活动并购买门票。此系统需要处理大量的并发请求，尤其是在热门事件的门票刚一开售时。在多用户同时操作的情况下，数据库必须确保每笔交易的完整性和一致性，避免超卖或数据冲突。

**问题表现**

- 在门票开售时，多个用户同时尝试购买同一场事件的门票，导致并发问题。
- 数据库事务处理不当可能导致某些门票被重复卖出（超卖）。

**并发控制策略**

1. **使用合适的事务隔离级别**：
   - 设置适宜的事务隔离级别以防止诸如脏读、不可重复读和幻读之类的问题。例如，可以设置为`REPEATABLE READ`或`SERIALIZABLE`。
2. **避免长事务**：
   - 确保事务尽可能短，减少锁定资源的时间，避免其他事务长时间等待锁释放。
3. **使用乐观锁或悲观锁**：
   - 对于更新操作，使用乐观锁（通过版本号或时间戳）或悲观锁（在事务开始时显式加锁）来控制并发。
4. **锁粒度控制**：
   - 尽量使用行级锁而不是表级锁，减少锁的范围，提高并发处理能力。
5. **事务超时设置**：
   - 设置合适的事务超时时间，避免事务长时间占用资源。

```go
// Ticket 代表一个门票的结构
type Ticket struct {
	EventID  int
	TicketID int
	Version  int // 用于乐观锁的版本号
}

// TicketService 处理票务服务的逻辑
type TicketService struct {
	DB *sql.DB
}

// PurchaseTicket 尝试购买门票
func (ts *TicketService) PurchaseTicket(eventID int, ticketID int) error {
	// 开始事务
	tx, err := ts.DB.Begin()
	if err != nil {
		return err
	}

	// 检查门票是否存在且未被售出
	var ticket Ticket
	err = tx.QueryRow("SELECT * FROM tickets WHERE event_id = ? AND ticket_id = ? FOR UPDATE", eventID, ticketID).Scan(&ticket.EventID, &ticket.TicketID, &ticket.Version)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 检查门票是否已被售出
	if ticket.Version == -1 { // 假设-1表示门票已售出
		tx.Rollback()
		return fmt.Errorf("票已售完")
	}

	// 更新门票状态为已售出
	// 使用乐观锁，通过版本号检查在事务开始后门票是否被修改过
	updatedRows, err := tx.Exec("UPDATE tickets SET version = -1 WHERE event_id = ? AND ticket_id = ? AND version = ?", eventID, ticketID, ticket.Version)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 检查是否更新了行
	if updatedRows == 0 {
		tx.Rollback()
		return fmt.Errorf("购票失败，门票可能已被其他人购买")
	}

	// 提交事务
	return tx.Commit()
}
```



## 6. 分区和分片

**策略：**

- 对大数据表进行分区，提高查询和维护效率。
- 在分布式系统中使用分片技术分散数据负载。

**适用场景：**

- 数据量巨大，单个表难以管理。
- 需要水平扩展数据库以应对高并发。



### 场景描述

有一个大型电商平台，类似于亚马逊或阿里巴巴，该平台每天产生数以百万计的订单。随着时间的推移，订单数据量迅速增长，单个订单表变得非常庞大，导致查询和维护效率低下。此外，在高流量时段，数据库需要处理大量的并发请求，这对数据库的性能提出了挑战。

**问题表现**

- 查询历史订单数据时，响应时间缓慢。
- 在促销或节日期间，数据库并发访问量剧增，导致性能瓶颈。

**分区和分片策略**

- 1. 对订单表进行分区

**技术**：MySQL内置的分区功能。

**策略详解**： 范围分区根据列值的范围将数据分散到不同的分区。在订单表的情况下，可以根据订单日期来创建分区。

**MySQL分区示例**：

```sql
CREATE TABLE orders (
    order_id INT NOT NULL,
    user_id INT NOT NULL,
    order_date DATE NOT NULL,
    -- 其他订单字段
) ENGINE=InnoDB
PARTITION BY RANGE (TO_DAYS(order_date)) (
    PARTITION p202301 VALUES LESS THAN (TO_DAYS('2023-02-01')),
    PARTITION p202302 VALUES LESS THAN (TO_DAYS('2023-03-01')),
    -- 为每个月创建一个分区
);
```

- 2. 使用分片技术

**技术**：分片中间件或自定义分片逻辑。

**策略详解**： 分片通常在应用层实现，根据分片键的值将数据路由到不同的数据库节点。

```go
// 根据用户ID进行哈希分片
func getShardKey(userID int) int {
    return userID % numShards // 假设有numShards个分片
}

// 获取对应分片的数据库连接
func getDBConnection(shardKey int) *sql.DB {
    // 根据分片键获取数据库连接
}
```

- 3. 实施智能路由

**技术**：应用层路由逻辑或使用代理服务器。

**策略详解**： 智能路由根据分片键将请求定向到正确的数据库节点，通常在应用层实现。

最常用的智能路由策略之一是**哈希分片**。这种策略通过哈希算法将数据均匀分布到不同的分片上，通常使用分片键（如用户ID或订单ID）的哈希值来决定数据应该存储在哪个分片。

```go
// Order 是订单数据的结构体
type Order struct {
	OrderID   int
	UserID    int
	OrderDate string
}

// DBShard 代表一个数据库分片
type DBShard struct {
	Database *sql.DB
	ShardKey int
}

// OrderService 结构体持有所有分片的数据库连接
type OrderService struct {
	Shards []*DBShard
}

// NewOrderService 初始化订单服务并创建所有分片的数据库连接
func NewOrderService(shardCount int) *OrderService {
	orderService := &OrderService{
		Shards: make([]*DBShard, shardCount),
	}
	// 这里应该添加数据库连接逻辑，创建每个分片的连接
	// 为了示例，只是简单地创建了分片
	for i := 0; i < shardCount; i++ {
		orderService.Shards[i] = &DBShard{ShardKey: i} // 实际上这里应创建并赋值 *sql.DB
	}
	return orderService
}

// GetOrderFromShard 根据订单ID的哈希值查询订单
func (os *OrderService) GetOrderFromShard(orderID int) (*Order, error) {
	// 计算订单ID的哈希值，确定分片键
	shardKey := hash(orderID)
	dbShard := os.Shards[shardKey]

	// 从选定的分片中查询订单
	// 这里使用了一个模拟的SQL查询，实际应用中应替换为真实的查询
	query := "SELECT order_id, user_id, order_date FROM orders WHERE order_id = ?"
	var order Order
	err := dbShard.Database.QueryRow(query, orderID).Scan(&order.OrderID, &order.UserID, &order.OrderDate)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// hash 函数根据订单ID生成哈希值，用于确定分片键
func hash(v int) int {
	// 使用简单的哈希函数，实际应用中可能需要更复杂的哈希算法
	return v % len(OrderService.Shards)
}
```



- 4. 维护分片的均衡

**技术**：监控系统和自动化脚本。

**策略详解**： 通过监控系统来跟踪每个分片的数据量和查询负载，使用自动化脚本来重新平衡数据。

维护分片均衡的最常用技术是**基于负载的自动分片重平衡**。

```go
// 定义Shard结构，代表数据库的一个分片
type Shard struct {
    ID     int    // 分片的唯一标识符
    Load   float64 // 当前分片的负载量，可以根据实际场景定义负载的具体含义，如查询率、数据量等
    MaxLoad float64 // 分片的最大负载阈值，超过这个值时需要进行负载均衡
}

// MigrateData方法实现将数据从当前分片迁移到目标分片的逻辑
func (s *Shard) MigrateData(target *Shard) {
    // 模拟数据迁移过程，这里简化为将一半的负载从当前分片迁移到目标分片
    // 实际情况可能涉及更复杂的数据迁移逻辑
    migrationLoad := s.Load / 2
    s.Load -= migrationLoad // 减少源分片的负载
    target.Load += migrationLoad // 增加目标分片的负载
}

// balanceShards函数负责检查所有分片的负载情况，并执行负载均衡
func balanceShards(shards []*Shard) {
    // 遍历分片列表，检查每个分片的负载
    for _, shard := range shards {
        // 如果分片的负载超过了预设的最大负载阈值
        if shard.Load > shard.MaxLoad {
            // 寻找负载最轻的分片以作为数据迁移的目标分片
            target := findLightestShard(shards)
            // 如果找到了较空闲的分片，则进行数据迁移
            if target != nil {
                shard.MigrateData(target)
            }
        }
    }
}

// findLightestShard函数遍历分片列表，找出负载最轻的分片
func findLightestShard(shards []*Shard) *Shard {
    // 假设第一个分片是负载最轻的分片
    lightest := shards[0]
    // 遍历所有分片，找出负载量最小的分片
    for _, shard := range shards {
        if shard.Load < lightest.Load {
            lightest = shard // 更新负载最轻的分片为当前分片
        }
    }
    // 返回负载最轻的分片
    return lightest
}
```



- 5. 分片和复制结合使用

**技术**：主从复制设置和故障转移机制。

**策略详解**： 每个分片都有一个或多个副本，主分片处理写入操作，从分片处理读取操作。在主分片故障时，从分片可以接管。

```sql
-- 在主分片上
CHANGE MASTER TO 
MASTER_HOST='master_host', 
MASTER_USER='master_user', 
MASTER_PASSWORD='master_password', 
MASTER_LOG_FILE='binlog_file', 
MASTER_LOG_POS=binlog_position;

-- 在从分片上
START SLAVE;
```



## 7. 定期维护

**策略：**

- 定期执行数据库的维护任务，如优化表、重建索引。
- 清理无用的数据和日志。

**适用场景：**

- 数据库长时间运行后性能下降。
- 需要保持数据库的长期健康和性能。



### 场景描述

公司内部日报系统用于记录员工每日的工作情况。随着时间推移，日报数据不断积累，可能会导致数据库性能下降。

**问题表现**

- 日报查看和提交操作变慢。
- 数据库存储空间紧张。

1. **清理旧数据**：
   - 定期删除或归档超过保留期限的日报数据。
2. **优化表**：
   - 定期执行`OPTIMIZE TABLE`，以整理数据文件和提高访问速度。
3. **重建索引**：
   - 定期检查索引，必要时重建以保持查询效率。
4. **监控数据库状态**：
   - 使用监控工具跟踪性能指标，及时发现问题。

```sql
-- 每月执行一次，清理一年前的日报数据
DELETE FROM daily_reports WHERE date < DATE_SUB(NOW(), INTERVAL 1 YEAR);

-- 优化日报表
OPTIMIZE TABLE daily_reports;
```



https://mp.weixin.qq.com/s?__biz=MzU2NjIzNDk5NQ%3D%3D&chksm=fcad713ccbdaf82a07e4b426abc493fa1a89da72f061e0d0a0f496a1f8bd583cd13f653912f8&idx=1&mid=2247510816&scene=27&sn=2a0945993af22696dd09312681c6cc46&utm_campaign=geek_search&utm_content=geek_search&utm_medium=geek_search&utm_source=geek_search&utm_term=geek_search#wechat_redirect

## 查询优化神器 – explain命令

关于explain命令相信大家并不陌生，具体用法和字段含义可以参考官网[explain-output](http://dev.mysql.com/doc/refman/5.5/en/explain-output.html)，这里需要强调rows是核心指标，绝大部分rows小的语句执行一定很快（有例外，下面会讲到）。所以优化语句基本上都是在优化rows。

## 慢查询优化基本步骤

0.先运行看看是否真的很慢，注意设置SQL_NO_CACHE
1.where条件单表查，锁定最小返回记录表。这句话的意思是把查询语句的where都应用到表中返回的记录数最小的表开始查起，单表每个字段分别查询，看哪个字段的区分度最高
2.explain查看执行计划，是否与1预期一致（从锁定记录较少的表开始查询）
3**.order by limit 形式的sql语句让排序的表优先查**[**数据库会这样做，因此你需要建立覆盖索引的时候优先考虑order by**]
4.了解业务方使用场景
5.加索引时参照建索引的几大原则
6.观察结果，不符合预期继续从0分析

下面是对每个优化建议的详细解释，包括为什么这些建议是重要的，以及它们如何影响数据库查询性能。

**尽可能建立索引【order by优先，区分度高的优先，join，group by】**

```
在数据库优化的过程中，每一步的操作都有特定的目的和意义，针对你提出的几个问题，我将逐一进行详细解释。

### 0. 为什么要设置 `SQL_NO_CACHE`

`SQL_NO_CACHE` 是一个用于控制查询缓存的 SQL 提示，它告诉数据库服务器不要使用查询缓存，也不要将查询结果存储在缓存中。这在调试或优化 SQL 语句时非常重要，主要原因有：

- **避免误导性的优化结果**：查询缓存会导致后续相同的查询直接从缓存中读取结果，而不是重新执行 SQL 查询。如果你在调试或优化 SQL 查询时没有禁用缓存，可能会让你误以为查询非常快，而实际上只是缓存的结果。这样，你无法真实评估查询的性能表现。
- **精确分析查询时间**：通过禁用查询缓存，你可以确保每次查询都是真实的执行，从而更准确地分析 SQL 查询的响应时间，识别潜在的性能瓶颈。

因此，在优化 SQL 时设置 `SQL_NO_CACHE` 可以保证你看到的是实际查询执行的效果，而不是缓存的结果。

### 1. 为什么要看哪个字段的区分度最高

在查询优化中，选择性高的字段（即“区分度高”的字段）会显著影响查询效率。区分度指的是字段中不同值的数量占总记录数的比例。区分度越高，字段越容易缩小结果集范围，因此查询性能也会更好。

- **区分度高的字段有助于减少查询的返回结果**：当 `WHERE` 条件应用在区分度高的字段上时，可以更有效地缩小查询结果集，减少数据库的计算量。如果一个字段的区分度很低，几乎没有什么过滤效果，查询时需要遍历更多记录，性能会下降。
- **影响索引选择和使用**：通常数据库会优先选择高区分度的字段进行索引查找，因此在优化 SQL 查询时，我们需要首先评估每个字段的区分度，从区分度最高的字段开始优化。
  
你可以通过执行单字段的 `SELECT DISTINCT` 查询来评估哪个字段的区分度最高，从而选择合适的字段来进行查询优化。

### 2. 为什么要用 `EXPLAIN` 查看执行计划，是否与预期一致

`EXPLAIN` 是 SQL 中用于查看查询执行计划的工具，它能够告诉你数据库如何执行查询，具体包括使用的索引、扫描的表、返回的记录数估算等。通过 `EXPLAIN`，可以验证你的优化思路是否正确，比如：

- **查询是否从锁定记录较少的表开始**：根据 `EXPLAIN` 的输出，你可以看到 SQL 查询是如何在不同的表之间选择顺序和扫描路径。你可以通过 `rows` 列查看数据库预估要扫描的记录数，验证它是否从预期中返回记录最少的表开始。
- **是否正确利用索引**：你可以通过 `EXPLAIN` 查看查询是否使用了合适的索引，以及哪些字段的索引被使用。理想情况下，查询应该尽可能使用索引来减少全表扫描。

通过 `EXPLAIN` 工具，你可以识别查询过程中出现的瓶颈，如未使用索引、错误的扫描顺序等，并进行针对性优化。

### 3. 为什么 `ORDER BY LIMIT` 形式的 SQL 语句让排序的表优先查

在 SQL 查询中，如果使用 `ORDER BY` 和 `LIMIT`，查询通常会根据某个字段排序并只返回结果集的前几行。让排序的表优先查的原因是：

- **减少排序的计算量**：如果排序操作发生在查询结果集生成之前，查询就必须对所有结果进行排序，这会导致性能问题。通过让排序的表优先查，数据库可以更早地筛选出满足 `LIMIT` 条件的行，从而避免对整个结果集进行排序。
- **提高查询效率**：如果能够在查询初期就限制结果集的大小，数据库引擎可以减少不必要的扫描和计算，从而提升整体查询效率。特别是在大数据集上，这种策略能够极大地减少排序和返回结果的开销。

例如，当 SQL 查询语句有一个大表的 `ORDER BY` 和 `LIMIT` 时，优化器可以选择先查询返回的行较少的表，然后再进行排序，进而减少排序的工作量。

### 4. 了解业务方使用场景

了解业务场景对于 SQL 优化是非常重要的，因为优化的方向取决于具体的业务需求。例如：

- **查询模式**：某些场景可能更多地是读操作，另一些场景则涉及频繁的写操作。不同的查询模式下，索引策略和表的设计都会有所不同。
- **数据规模和增长情况**：在高并发场景或者大数据场景下，你可能需要考虑如何减少锁的竞争、提高并发度，以及如何避免全表扫描等性能问题。
- **响应时间要求**：有的业务对查询响应时间要求很高，这种情况下，你需要重点关注查询的延迟，并选择合适的缓存、索引策略来优化响应时间。

通过了解业务场景，可以更加有针对性地进行查询优化，从而避免盲目地进行索引或查询结构上的调整。

### 5. 加索引时参照建索引的几大原则

加索引是优化 SQL 性能的常用方法，但并不是所有情况都适合加索引，常见的索引设计原则有：

- **选择性高的字段加索引**：高选择性意味着可以有效减少查询的扫描行数，因此应该优先为选择性高的字段加索引。
- **频繁出现在 `WHERE`、`JOIN`、`GROUP BY` 和 `ORDER BY` 中的字段加索引**：这些操作往往涉及大量的数据处理，适合通过索引来提高效率。
- **避免为频繁更新的字段加索引**：因为每次更新操作都会涉及索引的更新，索引数量过多反而会拖慢插入和更新操作的效率。
- **复合索引的顺序要考虑查询使用场景**：复合索引的字段顺序要根据查询条件中的字段顺序来设计，优先将经常出现在 `WHERE` 和 `ORDER BY` 中的字段放在前面。

通过遵循这些原则，能够合理地为查询加索引，提升查询效率。

### 6. 为什么要观察结果，不符合预期继续从 0 分析

优化 SQL 查询是一个迭代的过程，查询的性能会受到多种因素的影响，比如数据量、索引设计、查询条件等。在进行优化后，观察实际查询的结果和性能表现是非常关键的一步：

- **评估优化效果**：通过观察查询执行的时间、返回的结果集、资源消耗情况，你可以确定优化是否达到了预期效果。如果没有达到，可能说明执行计划中有潜在问题。
- **不断调整和优化**：如果优化结果不符合预期，需要重新分析问题所在，可能需要调整索引、修改查询结构，或者是更深入地分析数据库执行计划中的瓶颈。

通过这种迭代的分析和调整，可以逐步找出影响查询性能的因素，并实现最佳的优化效果。

### 总结

- `SQL_NO_CACHE` 使得查询能够真实反映实际的执行时间，而不是从缓存中获取结果。
- 确定区分度最高的字段能够有效缩小查询范围，提高查询效率。
- `EXPLAIN` 用于查看执行计划，验证优化思路是否与预期一致。
- `ORDER BY LIMIT` 让排序的表优先查，能减少排序计算量，提升效率。
- 索引设计时需要遵循选择性、查询使用频率等原则。
- 优化 SQL 是一个迭代的过程，持续观察和调整至关重要。
```



### 3. JOIN 和 ORDER BY 都要添加索引，JOIN 的字段应该是相同的类型
- **解释**：
  - 在执行 `JOIN` 操作时，数据库会对两个表的相关字段进行匹配。**如果这两个字段没有索引，数据库需要对表进行全表扫描**，导致查询性能下降。因此，在 `JOIN` 的字段上创建索引可以显著提高查询效率。
  - 此外，`JOIN` 的字段应该是相同的数据类型，以避免数据库在匹配时进行类型转换，这也会导致性能问题。
  - 对 `ORDER BY` 使用的字段添加索引，可以加快排序操作，尤其是在大数据集上。

- **示例**：
  ```sql
  -- 假设有两个表：orders 和 customers
  SELECT orders.id, customers.name 
  FROM orders 
  JOIN customers ON orders.customer_id = customers.id 
  ORDER BY orders.date;
  
  -- 应为 orders.customer_id 和 customers.id 添加索引
  CREATE INDEX idx_orders_customer_id ON orders(customer_id);
  CREATE INDEX idx_customers_id ON customers(id);
  
  -- 应为 orders.date 添加索引
  CREATE INDEX idx_orders_date ON orders(date);
  ```

### 4. 避免 `SELECT *`，从数据库里读出越多的数据，查询就会变得越慢
- **解释**：
  - 使用 `SELECT *` 会将表中的所有列都返回，即使某些列在查询中不需要使用。这会增加数据传输的负担，尤其是在列数较多或行数较大的表上。此外，`SELECT *` 还会增加数据库的 I/O 操作，影响性能。
  - 更好的做法是只选择查询中需要的列，以减少传输的数据量和数据库的负载。

- **示例**：
  ```sql
  -- 避免使用 SELECT *
  SELECT id, name, date 
  FROM orders 
  WHERE status = 'completed';
  ```

### 5. `LIMIT` 分页会每次从表头开始扫描，注意添加索引条件缩小扫表范围
- **解释**：
  - 在分页查询中，使用 `LIMIT` 会让数据库从头开始扫描记录，直到找到满足条件的记录。这对于大表来说非常耗时。因此，最好在 `LIMIT` 查询中使用合适的索引条件，限制扫描的范围，减少不必要的记录扫描。
  - 在适当的列上添加索引，可以让查询从更靠近结果集的地方开始扫描，从而提高性能。

- **示例**：
  ```sql
  -- 优化分页查询
  SELECT id, name 
  FROM users 
  WHERE created_at > '2023-01-01' 
  ORDER BY id 
  LIMIT 100, 10;
  
  -- 应为 created_at 和 id 列添加索引
  CREATE INDEX idx_users_created_at ON users(created_at);
  CREATE INDEX idx_users_id ON users(id);
  ```

### 6. `COUNT(*)` 语句会按索引全表扫描，[[默认使用主键效率比指定二级普通索引低]]
- **解释**：
  - `COUNT(*)` 是一种常用的聚合查询，用于计算表中的记录数。默认情况下，MySQL 会使用主键进行全表扫描。如果表有大量记录，这会非常耗时。
  - 为了提高 `COUNT(*)` 的性能，可以考虑使用**覆盖索引（即只包含需要的字段的索引），使查询不必回表**，从而加快速度。

- **示例**：
  ```sql
  -- 创建一个覆盖索引以提高 COUNT(*) 的性能
  CREATE INDEX idx_users_status ON users(status);
  
  -- 使用 COUNT(*)，MySQL 可以通过 idx_users_status 索引完成统计
  SELECT COUNT(*) FROM users WHERE status = 'active';
  ```

### 7. 添加索引时，避免索引失效，索引遵循最左匹配原则
- **解释**：
  - 在创建复合索引时，MySQL 会根据索引的最左前缀进行匹配。如果查询中没有使用到索引的最左列，那么该索引就会失效。因此，创建索引时，必须根据查询的实际使用情况来设计索引列的顺序，保证最左边的列能够有效参与查询。

- **示例**：
  ```sql
  -- 复合索引
  CREATE INDEX idx_users_name_age ON users(name, age);
  
  -- 索引可以被利用，因为查询条件中使用了最左匹配列
  SELECT * FROM users WHERE name = 'John' AND age = 30;
  
  -- 索引失效，因为最左匹配列（name）没有在查询条件中使用
  SELECT * FROM users WHERE age = 30;
  ```

### 8. 每个数据表的数据不能超过 500 万，超过需要进行分库分表优化
- **解释**：
  - 随着数据量的增加，单个表的性能会逐渐下降。一般来说，当表中的数据量达到数百万甚至更多时，查询、插入、更新等操作的性能都会受到影响。为了保持高性能，可以将数据进行分库分表，将数据水平切分到多个表或数据库中，以减小单个表的大小和查询负载。
  - 分库分表可以基于某个字段（如用户 ID、订单 ID）进行分片，使每个分片的数据量保持在合理的范围内。

- **示例**：
  ```sql
  -- 水平分表示例
  CREATE TABLE users_0 LIKE users;
  CREATE TABLE users_1 LIKE users;
  
  -- 基于用户 ID 的最后一位来决定数据存储在哪个表中
  INSERT INTO users_0 SELECT * FROM users WHERE MOD(id, 2) = 0;
  INSERT INTO users_1 SELECT * FROM users WHERE MOD(id, 2) = 1;
  ```

### 9. SQL 进行 JOIN 连接时，需要遵循驱动表是小表，被驱动表是大表，并且在 ON 的条件里面添加索引
- **解释**：
  - 在联接查询中，MySQL 通常会选择记录较少的表作为驱动表（外表），从中获取数据后，再去驱动表（内表）中查找匹配的记录。这样做可以减少数据的匹配次数，优化查询性能。
  - 同时，`JOIN` 操作中的 `ON` 条件应该使用索引，这样可以避免全表扫描，提高查询效率。

- **示例**：
  ```sql
  -- 假设 orders 表较大，customers 表较小
  SELECT orders.id, customers.name 
  FROM orders 
  JOIN customers ON orders.customer_id = customers.id;
  
  -- 为 customer_id 添加索引
  CREATE INDEX idx_orders_customer_id ON orders(customer_id);
  ```

### 10. SQL 里面同时有 WHERE 和 ORDER BY 的时候，把这些字段建立联合索引，可以直接使用联合索引的默认排序，不需要通过 WHERE 查询出符合条件的行后，重新进行排序，可以大大提高查询效率
- **解释**：
  - 当查询中既有 `WHERE` 过滤条件，又有 `ORDER BY` 排序时，MySQL 通常会先执行 `WHERE` 条件过滤出符合条件的行，然后再对结果进行排序。如果在相关字段上创建联合索引，MySQL 可以直接利用索引来过滤和排序，从而避免额外的排序开销。
  - 这种优化特别有效于大数据集的查询，能够显著提高性能。

- **示例**：
  ```sql
  -- 假设查询需要根据 status 过滤，并根据 date 排序
  SELECT * FROM orders WHERE status = 'completed' ORDER BY date DESC;
  
  -- 创建联合索引以优化查询
  CREATE INDEX idx_orders_status_date ON orders(status, date);
  
  -- MySQL 可以直接利用联合索引完成查询和排序
  ```

### 总结
通过遵循这些优化建议，可以显著提高数据库查询的性能。这些优化不仅包括对索引的合理使用，还涉及到如何设计和优化数据库结构，以及如何根据具体的业务场景和数据特性进行调整。通过逐步应用这些优化技术，可以确保数据库在处理大数据量和复杂查询时仍然保持高效的性能。

````
### 慢查询优化基本步骤详细解释

#### 0. 先运行看看是否真的很慢，注意设置 SQL_NO_CACHE
- **目的**：首先要确认查询是否确实很慢。通过实际运行查询，获取查询的执行时间。
- **SQL_NO_CACHE**：通过设置 `SQL_NO_CACHE`，可以确保 MySQL 不会使用查询缓存，从而获得查询的真实执行时间。这是因为查询缓存可能会掩盖实际的查询性能问题。

```sql
SELECT SQL_NO_CACHE * FROM your_table WHERE your_condition;
```

#### 1. WHERE 条件单表查，锁定最小返回记录表
- **解释**：在进行复杂查询时，通常会涉及多个表的联接。在这种情况下，优化的关键是先从最小的可能返回记录的表开始查询。通过单表查询，逐个检查 `WHERE` 条件，查看每个表在应用 `WHERE` 条件后的返回记录数。
- **步骤**：
  1. 对每个表单独执行查询，检查每个 `WHERE` 条件的过滤效果。
  2. 识别出返回记录最少的表，因为这通常意味着查询的工作量最小。
  3. 优化查询顺序，使得最小返回记录的表优先参与联接。

```sql
-- 单表查询示例：
SELECT COUNT(*) FROM table1 WHERE condition1;
SELECT COUNT(*) FROM table2 WHERE condition2;
```

- **区分度最高的字段**：如果一个表的某个字段能够过滤出很少的记录（区分度高），优先使用这个字段进行筛选。

#### 2. EXPLAIN 查看执行计划，是否与1预期一致
- **目的**：使用 `EXPLAIN` 分析 MySQL 实际采用的执行计划，检查查询的执行顺序、使用的索引等是否符合预期。
- **步骤**：
  1. 运行 `EXPLAIN` 查看查询执行计划。
  2. 检查 `EXPLAIN` 输出中涉及到的表、访问类型、可能的键、使用的键等信息。
  3. 确保执行计划中优先访问了最小返回记录的表，并且使用了高区分度的字段。

```sql
EXPLAIN SELECT ... FROM table1 JOIN table2 ... WHERE ...;
```

- **注意**：如果 `EXPLAIN` 结果与预期不符，可能需要调整查询的结构或索引的设计。

#### 3. ORDER BY LIMIT 形式的 SQL 语句让排序的表优先查
- **解释**：在 `ORDER BY LIMIT` 形式的查询中，如果表数据量大且需要排序，通常会先从排序的表开始查，以减少数据量。
- **优化**：
  1. 【确保使用了适当的索引来支持 `ORDER BY` 的字段，这样可以避免全表扫描和全表排序。】
  2. 【如果是多表联接，尽量让 `ORDER BY` 涉及的表尽早在查询中出现，以减少排序的开销。】

```sql
SELECT * FROM table ORDER BY indexed_column LIMIT 10;
```

- **索引优化**：对于 `ORDER BY` 和 `LIMIT` 组合查询，确保索引能够覆盖排序字段和过滤字段，减少排序和扫描的工作量。

#### 4. 了解业务方使用场景
- **目的**：在进行优化时，理解业务场景非常重要。例如，查询的执行频率、数据的增长速度、哪些查询是关键路径等。
- **步骤**：
  1. 与业务团队沟通，了解哪些查询是最关键的，哪些数据是最常用的。
  2. 根据业务需求，判断是否可以通过优化数据库结构（如分区、分库分表）来提升性能。

- **业务场景的影响**：了解业务场景有助于判断是否需要进行额外的优化，比如是否可以通过缓存、数据归档、或者数据预计算来减少查询负担。

#### 5. 加索引时参照建索引的几大原则
- **建索引原则**：
  1. **选择性高的列**：对于数据选择性高（唯一值多）的列，索引会更加有效。
  2. **频繁参与 WHERE、JOIN、ORDER BY、GROUP BY 的列**：优先考虑为这些列建立索引。
  3. **联合索引最左前缀原则**：联合索引在匹配时，首先会匹配最左边的字段，所以应当根据查询条件中的字段顺序来设计联合索引。
  4. **覆盖索引**：索引的列应尽可能覆盖查询需要的所有列，避免回表操作。
  5. **避免冗余索引**：不要为同样的列重复创建多个索引，这会增加维护成本。

```sql
CREATE INDEX idx_column1 ON table_name(column1);
CREATE INDEX idx_multi_column ON table_name(column1, column2);
```

- **索引的成本**：虽然索引可以加快查询，但过多或不必要的索引会导致插入、更新操作的性能下降，因此需要权衡。

#### 6. 观察结果，不符合预期继续从0分析
- **目的**：通过实际的查询性能测试和分析，观察优化后的效果。如果结果不符合预期，可能需要重新审视查询的设计或执行计划。
- **步骤**：
  1. 运行优化后的查询，测量执行时间和资源消耗。
  2. 如果性能仍不理想，使用 `EXPLAIN` 或 `SHOW PROFILE` 分析查询的执行过程。
  3. 根据分析结果，判断是否需要进一步优化，例如调整索引、重新设计查询结构，或进行数据库配置调整。

- **迭代优化**：优化数据库查询通常是一个迭代过程，需要不断地测试、分析、调整，直到性能达到可接受的水平。

通过以上步骤，逐步优化 SQL 查询，可以显著提高查询性能，减少查询响应时间。在实际操作中，理解数据库的执行计划以及如何针对具体业务需求进行调整和优化是至关重要的。
````

下面是对给出的 SQL 语句的详细解释和分析，包括每个部分的功能和潜在的性能问题。

### SQL 语句

```sql
SELECT DISTINCT cert.emp_id 
FROM cm_log cl 
INNER JOIN (
    SELECT emp.id AS emp_id, emp_cert.id AS cert_id 
    FROM employee emp 
    LEFT JOIN emp_certificate emp_cert 
        ON emp.id = emp_cert.emp_id 
    WHERE emp.is_deleted = 0
) cert 
    ON (
        cl.ref_table = 'Employee' 
        AND cl.ref_oid = cert.emp_id
    ) 
    OR (
        cl.ref_table = 'EmpCertificate' 
        AND cl.ref_oid = cert.cert_id
    ) 
WHERE 
    cl.last_upd_date >= '2013-11-07 15:03:00' 
    AND cl.last_upd_date <= '2013-11-08 16:00:00';
```

### SQL 语句详细解释

#### 1. 内部子查询

```sql
SELECT emp.id AS emp_id, emp_cert.id AS cert_id 
FROM employee emp 
LEFT JOIN emp_certificate emp_cert 
    ON emp.id = emp_cert.emp_id 
WHERE emp.is_deleted = 0
```

- **目的：**
  - 该子查询从 `employee` 表中选取所有未被标记为删除的员工（`emp.is_deleted = 0`）。
  - 子查询将 `employee` 表与 `emp_certificate` 表进行左连接（`LEFT JOIN`），即使员工没有相关的证书记录，也会将该员工包含在结果集中。
  - 返回结果集包括两个字段：`emp_id`（员工 ID）和 `cert_id`（证书 ID）。如果员工没有相关证书，`cert_id` 将为 `NULL`。

- **潜在问题：**
  - 如果 `employee` 表较大，并且没有适当的索引，可能会导致全表扫描，影响查询性能。
  - 该子查询的结果集可能很大（包含所有未删除的员工及其证书），这会增加后续查询的计算负担。

#### 2. 主查询与子查询结果的联接

```sql
INNER JOIN (
    ... 子查询结果 ...
) cert 
ON (
    cl.ref_table = 'Employee' 
    AND cl.ref_oid = cert.emp_id
) 
OR (
    cl.ref_table = 'EmpCertificate' 
    AND cl.ref_oid = cert.cert_id
)
```

- **目的：**
  - 主查询将 `cm_log` 表与子查询结果（命名为 `cert`）进行联接。
  - 联接条件为：`cm_log` 表的 `ref_table` 字段与 `'Employee'` 匹配，并且 `ref_oid` 字段与子查询结果中的 `emp_id` 匹配；或者 `ref_table` 字段与 `'EmpCertificate'` 匹配，并且 `ref_oid` 字段与 `cert_id` 匹配。
  - 联接类型为 `INNER JOIN`，因此仅当 `cm_log` 表的记录符合其中一个条件时，才会包含在最终结果中。

- **潜在问题：**
  - `OR` 条件会导致 MySQL 不使用索引，而是进行全表扫描，严重影响性能。
  - 子查询结果集可能很大，联接操作可能涉及大量数据处理，增加计算开销。

#### 3. 查询过滤条件

```sql
WHERE 
    cl.last_upd_date >= '2013-11-07 15:03:00' 
    AND cl.last_upd_date <= '2013-11-08 16:00:00'
```

- **目的：**
  - 该过滤条件限制查询结果为 `cm_log` 表中 `last_upd_date` 字段在指定时间范围内的记录。
  - 通过使用时间范围条件，减少 `cm_log` 表中的查询结果集大小。

- **潜在问题：**
  - 如果 `cm_log` 表上 `last_upd_date` 字段没有适当的索引，可能会导致全表扫描。
  - 尽管这个时间范围的过滤条件有助于缩小结果集，但如果表非常大，依然可能需要扫描大量记录。

#### 4. `DISTINCT` 用法

```sql
SELECT DISTINCT cert.emp_id 
```

- **目的：**
  - 使用 `DISTINCT` 关键字，去除结果集中 `emp_id` 字段的重复值。
  - 这在某些情况下是必要的，例如，如果员工同时有多个证书，或者同一个员工在日志中有多条记录，`DISTINCT` 可以确保每个员工 ID 只出现一次。

- **潜在问题：**
  - `DISTINCT` 操作需要对结果集进行去重处理，这可能需要额外的内存和计算，特别是在结果集较大时。

### EXPLAIN 分析结合详细说明【6点】

#### EXPLAIN 结果（假设）

id 顺序

select_type 类型

```sql
+----+-------------+------------+-------+---------------------------------+-----------------------+---------+-------------------+-------+--------------------------------+
| id | select_type | table      | type  | possible_keys                   | key                   | key_len | ref               | rows  | Extra                          |
+----+-------------+------------+-------+---------------------------------+-----------------------+---------+-------------------+-------+--------------------------------+
|  1 | PRIMARY     | cl         | range | cm_log_cls_id,idx_last_upd_date | idx_last_upd_date     | 8       | NULL              |   379 | Using where; Using temporary   |
|  1 | PRIMARY     | <derived2> | ALL   | NULL                            | NULL                  | NULL    | NULL              | 63727 | Using where; Using join buffer |
|  2 | DERIVED     | emp        | ALL   | NULL                            | NULL                  | NULL    | NULL              | 13317 | Using where                    |
|  2 | DERIVED     | emp_cert   | ref   | emp_certificate_empid           | emp_certificate_empid | 4       | meituanorg.emp.id |     1 | Using index                    |
+----+-------------+------------+-------+---------------------------------+-----------------------+---------+-------------------+-------+--------------------------------+
```

#### 1. 主查询的执行顺序（`id = 1`）

- **表 (`table`)：** `cl` 表，即 `cm_log` 表。
- **访问类型 (`type`)：** `range`，表明使用了【范围扫描】，这是因为时间条件的过滤 (`last_upd_date`)。
- **索引 (`key`)：** `idx_last_upd_date`，表示使用了 `cm_log` 表上的 `last_upd_date` 索引。
- **读取的行数 (`rows`)：** 预计 379 行，表明索引范围扫描后，只需处理 379 条记录。
- **额外信息 (`Extra`)：** `Using where; Using temporary`，表示 MySQL 在过滤数据的过程中使用了临时表，这可能是因为后续的 `DISTINCT` 操作导致的。

#### 2. 主查询与子查询的联接（`id = 1`，派生表 `<derived2>`）

- **表 (`table`)：** `<derived2>`，表示子查询生成的派生表。
- **访问类型 (`type`)：** `ALL`，全表扫描，表示 MySQL 没有使用索引来优化联接，而是扫描了整个派生表。
- **读取的行数 (`rows`)：** 预计 63727 行，这表明派生表中包含了大量数据。
- **额外信息 (`Extra`)：** `Using where; Using join buffer`，表示 MySQL 在执行联接时使用了联接缓存，这通常意味着联接操作涉及大量数据，并且无法通过索引优化。

```
为什么 select_type 会是 PRIMARY？
PRIMARY 是指主查询，这里是整个查询的最外层查询。它与 DERIVED 类型的派生表进行联接，并在最终结果集中应用 WHERE 和 DISTINCT 等操作。具体来说：

外层查询（主查询）：
sql
复制代码
SELECT DISTINCT cert.emp_id 
FROM cm_log cl 
INNER JOIN (
    ... 子查询 ...
) cert 
ON (
    cl.ref_table = 'Employee' 
    AND cl.ref_oid = cert.emp_id
) 
OR (
    cl.ref_table = 'EmpCertificate' 
    AND cl.ref_oid = cert.cert_id
) 
WHERE 
    cl.last_upd_date >= '2013-11-07 15:03:00' 
    AND cl.last_upd_date <= '2013-11-08 16:00:00';
在这个查询中，最外层的 SELECT DISTINCT 语句负责从 cm_log 表中选择数据，并将其与子查询生成的派生表联接。因此，select_type 显示为 PRIMARY。
子查询（派生表）：
sql
复制代码
SELECT emp.id AS emp_id, emp_cert.id AS cert_id 
FROM employee emp 
LEFT JOIN emp_certificate emp_cert 
    ON emp.id = emp_cert.emp_id 
WHERE emp.is_deleted = 0
这个部分的 select_type 是 DERIVED，因为它是主查询的一部分，并生成了一个派生表（<derived2>），供主查询使用。

```

#### 3. 子查询对 `employee` 表的处理（`id = 2`）

- **表 (`table`)：** `emp` 表，即 `employee` 表。
- **访问类型 (`type`)：** `ALL`，全表扫描，表明在 `employee` 表上没有使用索引来优化查询。
- **读取的行数 (`rows`)：** 预计 13317 行，表示 `employee` 表中符合条件的记录数。
- **额外信息 (`Extra`)：** `Using where`，表示 MySQL 使用了 `WHERE` 过滤条件，但由于没有适当的索引，导致全表扫描，性能较低。



#### 4. 子查询对 `emp_certificate` 表的联接（`id = 2`）

- **表 (`table`)：** `emp_cert` 表，即 `emp_certificate` 表。
- **访问类型 (`type`)：** `ref`，使用了引用索引来进行联接，效率较高。
- **索引 (`key`)：** `emp_certificate_empid`，表明使用了 `emp_certificate_empid` 索引来加速联接操作。
- **读取的行数 (`rows`)：** 每次联接只需处理 1 行数据，表示联接的效率较高。
- **额外信息 (`Extra`)：** `Using index`，表示通过索引直接获取数据，无需回表，查询效率较高。

### 主要性能瓶颈和优化建议

1. **全

表扫描 (`employee` 表和派生表)：**
   - **问题：** `employee` 表和派生表都进行全表扫描，特别是派生表 `<derived2>` 返回了大量记录，这对性能有较大影响。
   - **优化建议：** 在 `employee` 表上创建索引，例如 `is_deleted` 字段，避免全表扫描。考虑优化子查询的条件，减少派生表的记录数量。

2. **使用 `OR` 进行联接：**
   - **问题：** 使用 `OR` 条件进行联接，导致 MySQL 无法使用索引，从而选择了全表扫描。
   - **优化建议：** 可以通过将查询拆分为两个独立的查询，分别处理 `Employee` 和 `EmpCertificate` 的条件，最后合并结果，从而避免使用 `OR`。

3. **使用临时表和 `DISTINCT`：**
   - **问题：** 使用临时表和 `DISTINCT` 可能会导致额外的内存和磁盘 I/O，特别是当结果集较大时。
   - **优化建议：** 如果可能，减少使用 `DISTINCT` 或者优化查询条件，确保去重前的结果集尽可能小。

通过以上分析，可以看出该 SQL 查询存在多个性能瓶颈，通过合理的索引设计和查询结构优化，可以显著提升查询效率。

在 MySQL 中，`EXPLAIN` 命令用于分析查询语句的执行计划，帮助我们理解 MySQL 将如何执行给定的查询。以下是对 `EXPLAIN` 结果中每个字段的详细解释，结合你提供的例子进行说明。

### `EXPLAIN` 字段详细解释

```sql
+----+-------------+------------+-------+---------------------------------+-----------------------+---------+-------------------+-------+--------------------------------+
| id | select_type | table      | type  | possible_keys                   | key                   | key_len | ref               | rows  | Extra                          |
+----+-------------+------------+-------+---------------------------------+-----------------------+---------+-------------------+-------+--------------------------------+
```

#### 1. `id`
- **解释：**  
  `id` 字段表示查询中执行操作的顺序，每个 `SELECT` 子句或联合查询的一部分都会被标记一个唯一的 `id`。`id` 数值越小，优先级越高。`id` 相同的行表示可以并行执行的操作，不同的行则表示按照顺序执行。

- **例子中的值：**  
  - `1`: 主查询（包含 `cm_log` 表和 `derived2` 子查询的联接）。
  - `2`: 子查询，生成派生表 `derived2`。

#### 2. `select_type`
- **解释：**  
  `select_type` 描述了每个 `SELECT` 子句的类型，标识了查询中各个 `SELECT` 的角色和作用。

- **常见的类型：**
  - `SIMPLE`: 简单查询，不包含子查询或联合查询。
  - `PRIMARY`: 最外层的 `SELECT` 查询。
  - `UNION`: 用于联合查询的第二个或后续的 `SELECT`。
  - `SUBQUERY`: 子查询中的第一个 `SELECT`。
  - `DERIVED`: 派生表的 `SELECT`（通常是子查询）。

- **例子中的值：**  
  - `PRIMARY`: 表示主查询部分，外层 `SELECT`。
  - `DERIVED`: 表示子查询，该子查询生成一个派生表供主查询使用。

#### 3. `table`
- **解释：**  
  `table` 字段显示查询中当前步骤操作的表名。对于派生表（子查询生成的临时表），会显示为 `derivedN`，其中 `N` 是子查询的 `id`。

- **例子中的值：**
  - `cl`: 表示 `cm_log` 表。
  - `<derived2>`: 表示子查询生成的派生表。

#### 4. `type`【系统、常量、相等、引用、范围、索引】
- **解释：**  
  `type` 字段表示 MySQL 选择的访问类型，显示 MySQL 如何查找满足查询条件的行。这是优化查询性能的重要指标。
- **访问类型从好到差的顺序：**
  - `system`: 表仅有一行（等同于 `const`）。
  - `const`: 表最多只有一行匹配，用主键或唯一索引查找。
  - `eq_ref`: **每个联接使用主键或唯一索引查找最多一行**。
  - `ref`: 使用非唯一索引查找可能多行。
  - `range`: 使用索引范围查找（通常用于 `BETWEEN`、`<`、`>` 操作）。
  - `index`: 全索引扫描。
  - `ALL`: 全表扫描。
- **例子中的值：**
  - `range`: 用于 `cm_log` 表，表示使用索引范围扫描，这在 `WHERE` 条件中指定了时间范围时较为常见。
  - `ALL`: 用于派生表 `<derived2>` 和 `emp` 表，表示全表扫描，效率较低。
  - `ref`: 用于 `emp_cert` 表，表示使用索引查找特定的 `emp_id`。

```
以下是对每个 MySQL 访问类型的详细解释，并附有具体的 SQL 示例来帮助理解这些类型在实际查询中的应用。

### 1. `system`
- **描述**：当表中只有一行数据时，MySQL 会使用 `system` 作为访问类型。这种情况非常罕见。
- **举例**：
  假设表 `config` 只有一行记录：
  ```sql
  SELECT * FROM config;
```
  在这种情况下，MySQL 知道表中只有一行数据，因此 `EXPLAIN` 结果会显示 `system`。

### 2. `const`
- **描述**：当查询使用主键或唯一索引并且条件能够锁定最多一行数据时，MySQL 会使用 `const`。这意味着 MySQL 可以直接访问所需的行，而无需扫描多个记录。
- **举例**：
  假设表 `users` 中有主键 `id`：
  ```sql
  SELECT * FROM users WHERE id = 1;
  ```
  因为 `id` 是主键，并且查询条件 `id = 1` 能够唯一标识一行，所以 MySQL 使用 `const` 访问类型。

### 3. `eq_ref`【连接列是一个表的主键】
- **描述**：这种类型用于联接操作中。当 MySQL 使用主键或唯一索引查找每个联接表中的一行时，会使用 `eq_ref`。这是最理想的联接类型之一。
- **举例**：
  假设有两张表 `orders` 和 `users`，其中 `orders` 表中的 `user_id` 是 `users` 表的主键：
  ```sql
  SELECT orders.id, users.name 
  FROM orders 
  JOIN users ON orders.user_id = users.id;，
  ```
  因为 `users.id` 是主键且 `orders.user_id` 是外键，因此 MySQL 对 `users` 表的访问类型为 `eq_ref`。

**唯一索引**（Unique Index）是一种数据库索引，它确保索引列中的所有值都是唯一的，即不允许两个不同的记录在该索引列上有相同的值。

### 特点和用途：

1. **唯一性保证**：
   - 唯一索引确保列中的每个值都是唯一的，防止在该列中插入重复的数据。
   - 当试图插入或更新数据时，如果新数据的值与唯一索引列中的已有值冲突，数据库会拒绝该操作，并返回错误。

2. **快速查找**：
   - 像其他索引一样，唯一索引可以加快查询速度。因为它强制列中数据的唯一性，数据库引擎可以更高效地查找特定的记录。

3. **允许空值**：
   - 大多数数据库允许唯一索引列包含 `NULL` 值。然而，不同数据库对 `NULL` 的处理有所不同：
     - **MySQL**：**在 MySQL 中，多个 `NULL` 值被视为不相等的，因此允许多个 `NULL` 值出现在唯一索引列中**。
     - **SQL Server**：SQL Server 允许唯一索引列中的一个 `NULL` 值。如果有第二个 `NULL` 值插入，SQL Server 会引发错误。
     - **Oracle**：类似于 MySQL，允许多个 `NULL` 值。

4. **主键与唯一索引**：
   - **主键是一种特殊的唯一索引。主键不仅要求唯一性，还不允许 `NULL` 值。此外，每个表只能有一个主键，但可以有多个唯一索引。**

### 示例

假设你有一个用户表，要求每个用户的电子邮件地址在数据库中是唯一的，那么你可以在 `email` 列上创建一个唯一索引。

#### 在 MySQL 中创建唯一索引的 SQL 语句：

```sql
CREATE UNIQUE INDEX idx_unique_email ON users(email);
```

#### 在创建表时直接定义唯一索引：

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255)
);
```

在这个示例中，`email` 列被定义为 `UNIQUE`，这意味着用户表中不允许有两个记录的 `email` 值相同。

### 应用场景

- **防止重复数据**：在需要保证数据列值唯一的场景下，例如用户的电子邮件地址、身份证号码、社保号等，可以使用唯一索引。
- **快速查找特定记录**：唯一索引不仅可以防止重复数据，还能提高查询速度，例如在按用户 ID 或邮箱地址查找记录时。

### 注意事项

- **性能影响**：虽然索引可以加快查询速度，但索引的创建和维护（如在插入、更新时）会占用额外的资源。因此，在创建索引时，需要平衡性能和数据完整性之间的关系。
- **索引覆盖**：如果一个查询仅访问了唯一索引的列，数据库可以直接从索引中获取结果，而不需要访问实际的数据表，进一步提高查询性能。

总结来说，唯一索引是保证数据库中某些列值唯一的重要工具，它不仅维护数据的完整性，还可以提高查询效率。

### 4. `ref`【只查找where id = 3的多行数据有索引】

- **描述**：当查询使用非唯一索引并且该索引可能返回多行结果时，MySQL 使用 `ref` 访问类型。这通常用于索引列上的联接或查找操作。
- **举例**：
  假设表 `users` 有一个非唯一索引 `index_city` 在 `city` 列上：
  ```sql
  SELECT * FROM users WHERE city = 'New York';
  ```
  因为 `city` 不是唯一索引，查询可能返回多行结果，所以 MySQL 使用 `ref` 作为访问类型。

### 5. `range`【只查找 where id > 2 and id < 3的多行数据】
- **描述**：`range` 访问类型用于在索引列上执行范围查询，如 `BETWEEN`、`<`、`>` 等操作。MySQL 会使用索引查找满足条件的一组连续记录。
- **举例**：
  假设表 `orders` 中有一个索引 `index_order_date` 在 `order_date` 列上：
  
  ```sql
  SELECT * FROM orders WHERE order_date BETWEEN '2023-01-01' AND '2023-01-31';
  ```
  MySQL 会使用 `range` 访问类型，在 `order_date` 的索引上进行范围查找。

### 6. `index`[查找id的所有数据]
- **描述**：`index` 访问类型表示 MySQL 需要扫描整个索引，而不是通过扫描数据行。这类似于 `ALL`，但它只扫描索引，不扫描表数据。这种类型通常在索引覆盖查询中出现。
- **举例**：
  假设表 `users` 有一个覆盖索引 `index_name_email` 在 `name` 和 `email` 列上：
  ```sql
  SELECT name, email FROM users;
  ```
  如果查询只涉及 `name` 和 `email`，MySQL 可以通过扫描 `index_name_email` 索引来满足查询，因此使用 `index` 作为访问类型。

### 7. `ALL`【查找全表】
- **描述**：`ALL` 访问类型表示 MySQL 需要进行全表扫描，即它必须扫描表中的所有记录。这通常是最慢的查询类型之一，除非表非常小。
- **举例**：
  如果我们在没有索引的表 `users` 上进行查询：
  ```sql
  SELECT * FROM users WHERE age > 30;
  ```
  如果 `age` 列上没有索引，MySQL 必须扫描 `users` 表的所有记录，以找到满足条件的行，因此会使用 `ALL` 访问类型。

### 总结
- `system` 和 `const` 代表了最简单、最轻量级的查询操作。
- `eq_ref` 和 `ref` 主要用于联接操作，前者是通过主键或唯一索引的精确联接，后者是通过非唯一索引的多行联接。
- `range` 和 `index` 表示 MySQL 在某些条件下对索引的有效利用，其中 `range` 用于范围查询，而 `index` 用于索引覆盖查询。
- `ALL` 则表示没有使用索引的全表扫描，是最不理想的情况，除非在处理非常小的表。 

### index和ref的区别

````
在 MySQL 中，`ref` 和 `index` 是两种不同的访问类型，它们都涉及使用索引，但它们的操作方式和适用场景不同。为了更好地理解这两者的区别，我将通过详细的解释和举例来说明。

### 1. `ref` 访问类型

**描述**:
- `ref` 访问类型表示 MySQL 使用非唯一索引或前缀索引来查找数据。它是通过索引列与查询条件进行匹配来定位数据行。
- `ref` 通常用于连接操作中的外键列或 `WHERE` 子句中带有非唯一索引的列。
- 由于索引列可能包含重复值，`ref` 访问类型可能会返回多行数据（即在索引中有多个条目匹配查询条件）。

**适用场景**:
- 查询条件匹配多个索引条目的情况，例如使用非唯一索引或多对多关系中的连接操作。
- 当查询中的条件不是唯一约束时（例如，查询条件匹配多条记录）。

**数据访问内容**:
- `ref` 访问类型根据查询条件在索引中查找匹配的条目，并返回所有匹配的表中行的数据。

**举例**:
假设有一个 `employees` 表，存储了员工的信息，包含以下字段：

```sql
CREATE TABLE employees (
    id INT PRIMARY KEY,
    name VARCHAR(255),
    department_id INT,
    salary DECIMAL(10, 2)
);
```

我们在 `department_id` 列上创建了一个非唯一索引：

```sql
CREATE INDEX idx_department_id ON employees(department_id);
```

现在，我们执行以下查询：

```sql
SELECT * FROM employees WHERE department_id = 3;
```

在这个查询中，`ref` 访问类型会使用 `idx_department_id` 索引在 `department_id` 列中查找等于 3 的条目。由于 `department_id` 可能不唯一，查询可能会返回多行数据，比如：

```
+----+------+--------------+--------+
| id | name | department_id| salary |
+----+------+--------------+--------+
|  2 | Bob  | 3            | 50000  |
|  7 | Alice| 3            | 60000  |
| 14 | John | 3            | 55000  |
+----+------+--------------+--------+
```

这里，`ref` 访问类型返回了 3 行数据，因为在 `department_id = 3` 的条件下，有 3 条记录匹配。

### 2. `index` 访问类型

**描述**:
- `index` 访问类型表示 MySQL 进行的是全索引扫描，这类似于全表扫描（`ALL`），但扫描的是整个索引而不是数据表。
- `index` 访问类型通常用于排序操作或者当索引覆盖查询时（即查询只涉及索引列，且不需要访问数据表的其他列）。

**适用场景**:
- 需要遍历索引中所有条目来执行排序或分组的操作。
- 查询中的列可以通过索引直接满足，不需要访问表中的其他列。

**数据访问内容**:
- `index` 访问类型扫描整个索引的所有条目，不依赖特定的查询条件。

**举例**:
假设我们对 `employees` 表执行以下查询：

```sql
SELECT department_id FROM employees ORDER BY department_id;
```

在这种情况下，MySQL 会使用 `index` 访问类型来扫描 `idx_department_id` 索引。因为我们需要对 `department_id` 列进行排序，而这个列已经有索引覆盖，所以 MySQL 会遍历整个索引，而不需要访问实际的表数据。

### **区别总结**:

- **`ref` 访问类型**：
  - 基于查询条件使用非唯一索引查找匹配的索引条目，可能返回多行数据。
  - 查询条件通常是等值匹配，如 `WHERE` 子句中的条件或连接操作中的外键。
  - 只扫描满足条件的索引条目。

- **`index` 访问类型**：
  - 扫描整个索引，不依赖查询条件。适用于排序、分组或索引覆盖查询。
  - 类似于全表扫描，但扫描的是索引数据而不是表数据。
  - 会遍历索引中的所有条目。

### 举例对比

- **使用 `ref` 访问类型的查询**：
  
  ```sql
  SELECT * FROM employees WHERE department_id = 3;
  ```
  - MySQL 使用 `ref` 访问类型，通过 `idx_department_id` 索引查找 `department_id = 3` 的所有条目，可能会返回多行数据（因为多个员工可以属于同一个部门）。

- **使用 `index` 访问类型的查询**：
  
  ```sql
  SELECT department_id FROM employees ORDER BY department_id;
  ```
  - MySQL 使用 `index` 访问类型遍历 `idx_department_id` 索引的所有条目，以获取 `department_id` 列的排序结果，不依赖具体的查询条件。

在优化查询时，选择 `ref` 或 `index` 取决于查询的需求：`ref` 更适合基于条件的精准查找，而 `index` 适合用于需要遍历或排序索引的场景。
````



### 什么叫只遍历索引

````
遍历索引和遍历整个表是两种不同的数据扫描方式，它们的区别主要在于效率、数据访问的内容以及适用的查询场景。

### 1. **遍历整个表（全表扫描）**

**描述**:
- 当 MySQL 使用 `ALL` 访问类型时，会对表中的每一行数据进行扫描。这意味着 MySQL 需要读取表中的所有记录，并根据查询条件（如 `WHERE` 子句）过滤出符合条件的记录。
- 全表扫描的效率通常较低，尤其是在表很大的情况下，因为它需要读取表中的所有数据行。

**适用场景**:
- 当没有合适的索引可用时，MySQL 可能会使用全表扫描。
- 当查询返回的大部分数据行时，全表扫描有时也是更高效的选择。

**数据访问内容**:
- 在全表扫描中，MySQL 直接访问表中的每一行数据。这意味着每一行的所有列都会被读取，无论查询是否实际需要这些列。

**举例**:
假设你有一个 `users` 表，包含 1 亿行数据，结构如下：

```sql
CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(255),
    age INT,
    email VARCHAR(255)
);
```

如果你执行以下查询：

```sql
SELECT * FROM users WHERE age = 30;
```

在没有索引的情况下，MySQL 会对 `users` 表进行全表扫描，检查每一行的 `age` 值，看看是否等于 30。

### 2. **遍历索引**

**描述**:
- 当 MySQL 使用 `index` 访问类型时，它会遍历整个索引结构，而不是直接扫描表中的每一行。索引是一种数据结构，用于加速对表中数据的查找。
- 索引的大小通常比表要小得多，并且索引中的数据是按排序顺序存储的，这使得遍历索引通常比全表扫描更快。

**适用场景**:
- 当查询只涉及索引覆盖的列时，MySQL 可以通过遍历索引直接获取所需的数据，而无需访问表中的实际行。
- 当执行 `ORDER BY` 或 `GROUP BY` 操作时，遍历索引可以避免对表进行额外排序。

**数据访问内容**:
- 在索引扫描中，MySQL 只访问索引中包含的列。如果查询只涉及这些列，MySQL 不需要再访问表中的数据行，这种情况称为“覆盖索引”。
- 如果查询需要访问索引中未包含的列，MySQL 可能还需要访问表中的数据行以获取这些列的值。

**举例**:
假设你在 `users` 表的 `age` 列上创建了一个索引：

```sql
CREATE INDEX idx_age ON users(age);
```

现在执行以下查询：

```sql
SELECT age FROM users ORDER BY age;
```

MySQL 可以使用 `index` 访问类型来遍历 `idx_age` 索引，而不需要扫描整个表。因为查询只涉及 `age` 列，MySQL 可以直接从索引中获取结果，而无需访问 `users` 表中的实际数据行。

### **总结**:

- **效率**: 索引遍历通常比全表扫描更高效，尤其是在表很大、但索引较小时。索引遍历可以避免读取不必要的数据。
  
- **数据访问**:
  - 全表扫描会读取表中的所有数据行，访问每行的所有列。
  - 索引遍历只读取索引覆盖的列，可以避免访问不需要的列，且索引的排序特性可以加速某些查询（如 `ORDER BY`）。

- **使用场景**:
  - 全表扫描通常在没有合适索引或查询需要访问大部分数据时使用。
  - 索引遍历在查询需要快速访问部分数据或按索引列排序时使用，并且可以使用覆盖索引优化查询性能。 

在优化 SQL 查询时，优先考虑索引遍历，因为它可以大大减少 I/O 操作，从而提升查询性能。
````

通过 `EXPLAIN` 语句分析查询的访问类型，开发者可以了解查询的执行情况，并针对性能问题进行优化。
```



#### 5. `possible_keys`
- **解释：**  
  `possible_keys` 列显示了 MySQL 在当前查询中可能使用的索引。这个字段反映了根据查询条件，哪些索引是 MySQL 可以考虑使用的。

- **例子中的值：**
  - `cm_log_cls_id, idx_last_upd_date`: 表示在 `cm_log` 表中，可能使用到的索引有 `cm_log_cls_id` 和 `idx_last_upd_date`。
  - `NULL`: 表示在 `emp` 表和派生表 `<derived2>` 中，MySQL 认为没有可用的索引。

#### 6. `key`
- **解释：**  
  `key` 列表示 MySQL 实际上选择使用的索引。如果没有选择使用索引，则显示 `NULL`。

- **例子中的值：**
  - `idx_last_upd_date`: 用于 `cm_log` 表，表示实际使用了 `idx_last_upd_date` 索引来执行范围扫描。
  - `NULL`: 用于派生表 `<derived2>` 和 `emp` 表，表示未使用索引。
  - `emp_certificate_empid`: 用于 `emp_cert` 表，表示 MySQL 使用了 `emp_certificate_empid` 索引来进行查询。

#### 7. `key_len`
- **解释：**  
  `key_len` 列表示 MySQL 使用的索引的长度（字节数）。这是通过索引字段计算得出的，显示了 MySQL 实际使用了多少索引前缀来进行查询。

- **例子中的值：**
  - `8`: 用于 `cm_log` 表，表示使用了 `idx_last_upd_date` 索引的 8 字节长度，这通常是一个日期时间字段的字节数。
  - `4`: 用于 `emp_cert` 表，表示使用了 `emp_certificate_empid` 索引的 4 字节长度，可能是一个整数类型字段。

#### 8. `ref`
- **解释：**  
  `ref` 列显示的是查询中哪个字段或常量与 `key` 所指的索引一起被用于查找行。如果是常量查找（`const`），则表示为 `const`。

- **例子中的值：**
  - `NULL`: 用于 `cm_log` 和 `<derived2>`，表示没有特定的引用。
  - `meituanorg.emp.id`: 用于 `emp_cert` 表，表示 `emp_certificate_empid` 索引引用了 `employee` 表中的 `id` 字段。

#### 9. `rows`
- **解释：**  
  `rows` 列是 MySQL 估计需要读取的行数，用于执行查询。这是一个重要的指标，显示了查询扫描的行数多少，从而反映了查询的潜在成本。

- **例子中的值：**
  - `379`: 对 `cm_log` 表的索引扫描返回了 379 行。
  - `63727`: 派生表 `<derived2>` 返回了 63727 行，表明子查询生成了大量中间数据。
  - `13317`: 对 `emp` 表的全表扫描返回了 13317 行。
  - `1`: 对 `emp_cert` 表的索引查找，每次只返回 1 行。

#### 10. `Extra`
- **解释：**  
  `Extra` 列包含了 MySQL 查询计划的额外信息，描述了查询过程中使用的特定操作。

- **常见的值：**
  - `Using where`: 表示查询使用了 `WHERE` 过滤条件。
  - `Using index`: 表示查询只通过索引用到的数据，未回表。
  - `Using temporary`: 表示查询使用了临时表。
  - `Using filesort`: 表示 MySQL 需要额外的文件排序来处理 `ORDER BY` 操作。

- **例子中的值：**
  - `Using where; Using temporary`: 表示 `cm_log` 表的查询使用了 `WHERE` 条件过滤，并且使用了临时表。
  - `Using where; Using join buffer`: 表示派生表 `<derived2>` 使用了 `WHERE` 条件过滤，并且 MySQL 使用了联接缓存（Join Buffer）来加速联接操作。
  - `Using where`: 用于 `emp` 表，表示使用了 `WHERE` 过滤条件。
  - `Using index`: 用于 `emp_cert` 表，表示查询使用了索引，并未需要回表查询其他字段。

### 总结

通过上述 `EXPLAIN` 字段的解释，可以看出 MySQL 如何执行查询，并且可以识别出潜在的性能瓶颈。了解这些字段及其意义，有助于我们优化查询，减少不必要的全表扫描，合理使用索引，从而提高查询性能。





简述一下执行计划，首先mysql根据idx_last_upd_date索引扫描cm_log表获得379条记录；然后查表扫描了63727条记录，分为两部分，derived表示构造表，也就是不存在的表，可以简单理解成是一个语句形成的结果集，后面的数字表示语句的ID。derived2表示的是ID = 2的查询构造了虚拟表，并且返回了63727条记录。我们再来看看ID = 2的语句究竟做了写什么返回了这么大量的数据，首先全表扫描employee表13317条记录，然后根据索引emp_certificate_empid关联emp_certificate表，rows = 1表示，每个关联都只锁定了一条记录，效率比较高。获得后，再和cm_log的379条记录根据规则关联。从执行过程上可以看出返回了太多的数据，返回的数据绝大部分cm_log都用不到，因为cm_log只锁定了379条记录。

如何优化呢？可以看到我们在运行完后还是要和cm_log做join,那么我们能不能之前和cm_log做join呢？仔细分析语句不难发现，其基本思想是如果cm_log的ref_table是EmpCertificate就关联emp_certificate表，如果ref_table是Employee就关联employee表，我们完全可以拆成两部分，并用union连接起来，注意这里用union，而不用union all是因为原语句有“distinct”来得到唯一的记录，而union恰好具备了这种功能。如果原语句中没有distinct不需要去重，我们就可以直接使用union all了，因为使用union需要去重的动作，会影响SQL性能。
优化过的语句如下



MySQL

| 12345678910111213141516171819202122232425262728 | **select**  emp.id  **from**  cm_log cl  **inner join**  employee emp     **on** cl.ref_table = 'Employee'     **and** cl.ref_oid = emp.id  **where**  cl.last_upd_date >='2013-11-07 15:03:00'   **and** cl.last_upd_date<='2013-11-08 16:00:00'   **and** emp.is_deleted = 0  **union** **select**  emp.id  **from**  cm_log cl  **inner join**  emp_certificate ec     **on** cl.ref_table = 'EmpCertificate'     **and** cl.ref_oid = ec.id  **inner join**  employee emp     **on** emp.id = ec.emp_id  **where**  cl.last_upd_date >='2013-11-07 15:03:00'   **and** cl.last_upd_date<='2013-11-08 16:00:00'   **and** emp.is_deleted = 0 |
| ----------------------------------------------- | ------------------------------------------------------------ |
|                                                 |                                                              |

4.不需要了解业务场景，只需要改造的语句和改造之前的语句保持结果一致

5.现有索引可以满足，不需要建索引

6.用改造后的语句实验一下，只需要10ms 降低了近200倍！

| 12345678910 | +----+--------------+------------+--------+---------------------------------+-------------------+---------+-----------------------+------+-------------+ \| id \| select_type \| table   \| type  \| possible_keys          \| key        \| key_len \| ref          \| rows \| Extra    \| +----+--------------+------------+--------+---------------------------------+-------------------+---------+-----------------------+------+-------------+ \| 1 \| PRIMARY   \| cl     \| range \| cm_log_cls_id,idx_last_upd_date \| idx_last_upd_date \| 8    \| **NULL**         \| 379 \| Using where \| \| 1 \| PRIMARY   \| emp    \| eq_ref \| PRIMARY             \| PRIMARY      \| 4    \| meituanorg.cl.ref_oid \|  1 \| Using where \| \| 2 \| UNION    \| cl     \| range \| cm_log_cls_id,idx_last_upd_date \| idx_last_upd_date \| 8    \| **NULL**         \| 379 \| Using where \| \| 2 \| UNION    \| ec     \| eq_ref \| PRIMARY,emp_certificate_empid  \| PRIMARY      \| 4    \| meituanorg.cl.ref_oid \|  1 \|       \| \| 2 \| UNION    \| emp    \| eq_ref \| PRIMARY             \| PRIMARY      \| 4    \| meituanorg.ec.emp_id \|  1 \| Using where \| \| **NULL** \| UNION RESULT \| &lt;union1,2&gt; \| ALL  \| **NULL**              \| **NULL**       \| **NULL**  \| **NULL**         \| **NULL** \|       \| +----+--------------+------------+--------+---------------------------------+-------------------+---------+-----------------------+------+-------------+ |
| ----------- | ------------------------------------------------------------ |
|             |                                                              |



### 明确应用场景

举这个例子的目的在于颠覆我们对列的区分度的认知，一般上我们认为区分度越高的列，越容易锁定更少的记录，但在一些特殊的情况下，这种理论是有局限性的



MySQL

| 1234567891011 | **select**  *  **from**  stage_poi sp  **where**  sp.accurate_result=1   **and** (    sp.sync_status=0     **or** sp.sync_status=2     **or** sp.sync_status=4  ); |
| ------------- | ------------------------------------------------------------ |
|               |                                                              |

0.先看看运行多长时间,951条数据6.22秒，真的很慢

| 1    | 951 rows **in** set (6.22 sec) |
| ---- | ------------------------------ |
|      |                                |

1.先explain，rows达到了361万，type = ALL表明是全表扫描

| 12345 | +----+-------------+-------+------+---------------+------+---------+------+---------+-------------+ \| id \| select_type \| table \| type \| possible_keys \| key \| key_len \| ref \| rows  \| Extra    \| +----+-------------+-------+------+---------------+------+---------+------+---------+-------------+ \| 1 \| SIMPLE   \| sp  \| ALL \| **NULL**     \| **NULL** \| **NULL**  \| **NULL** \| 3613155 \| Using where \| +----+-------------+-------+------+---------------+------+---------+------+---------+-------------+ |
| ----- | ------------------------------------------------------------ |
|       |                                                              |

2.所有字段都应用查询返回记录数，因为是单表查询 0已经做过了951条

3.让explain的rows 尽量逼近951

看一下accurate_result = 1的记录数

| 12345678 | select count(*),accurate_result from stage_poi group by accurate_result; +----------+-----------------+ \| count(*) \| accurate_result \| +----------+-----------------+ \|   1023 \|       -1 \| \| 2114655 \|        0 \| \|  972815 \|        1 \| +----------+-----------------+ |
| -------- | ------------------------------------------------------------ |
|          |                                                              |

我们看到accurate_result这个字段的区分度非常低，整个表只有-1,0,1三个值，加上索引也无法锁定特别少量的数据

再看一下sync_status字段的情况

| 1234567 | select count(*),sync_status from stage_poi group by sync_status; +----------+-------------+ \| count(*) \| sync_status \| +----------+-------------+ \|   3080 \|      0 \| \| 3085413 \|      3 \| +----------+-------------+ |
| ------- | ------------------------------------------------------------ |
|         |                                                              |

同样的区分度也很低，根据理论，也不适合建立索引

问题分析到这，好像得出了这个表无法优化的结论，两个列的区分度都很低，即便加上索引也只能适应这种情况，很难做普遍性的优化，比如当sync_status 0、3分布的很平均，那么锁定记录也是百万级别的

4.找业务方去沟通，看看使用场景。业务方是这么来使用这个SQL语句的，每隔五分钟会扫描符合条件的数据，处理完成后把sync_status这个字段变成1,五分钟符合条件的记录数并不会太多，1000个左右。了解了业务方的使用场景后，优化这个SQL就变得简单了，因为业务方保证了数据的不平衡，如果加上索引可以过滤掉绝大部分不需要的数据

5.根据建立索引规则，使用如下语句建立索引

| 1    | alter table stage_poi add index idx_acc_status(accurate_result,sync_status); |
| ---- | ------------------------------------------------------------ |
|      |                                                              |

6.观察预期结果,发现只需要200ms，快了30多倍。

| 1    | 952 rows **in** set (0.20 sec) |
| ---- | ------------------------------ |
|      |                                |

我们再来回顾一下分析问题的过程，单表查询相对来说比较好优化，大部分时候只需要把where条件里面的字段依照规则加上索引就好，如果只是这种“无脑”优化的话，显然一些区分度非常低的列，不应该加索引的列也会被加上索引，这样会对插入、更新性能造成严重的影响，同时也有可能影响其它的查询语句。所以我们第4步调差SQL的使用场景非常关键，我们只有知道这个业务场景，才能更好地辅助我们更好的分析和优化查询语句。

### 无法优化的语句





MySQL

| 12345678910111213141516171819202122232425262728293031323334353637 | **select**  c.id,  c.name,  c.position,  c.sex,  c.phone,  c.office_phone,  c.feature_info,  c.birthday,  c.creator_id,  c.is_keyperson,  c.giveup_reason,  c.**status**,  c.data_source,  from_unixtime(c.created_time) **as** created_time,  from_unixtime(c.last_modified) **as** last_modified,  c.last_modified_user_id  **from**  contact c  **inner join**  contact_branch cb     **on** c.id = cb.contact_id  **inner join**  branch_user bu     **on** cb.branch_id = bu.branch_id     **and** bu.**status** **in** (     1,    2)   **inner join**    org_emp_info oei      **on** oei.data_id = bu.user_id      **and** oei.node_left >= 2875      **and** oei.node_right <= 10802      **and** oei.org_category = - 1   **order by**    c.created_time **desc** **limit** 0 ,    10; |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
|                                                              |                                                              |

还是几个步骤

0.先看语句运行多长时间，10条记录用了13秒，已经不可忍受

| 1    | 10 rows **in** set (13.06 sec) |
| ---- | ------------------------------ |
|      |                                |

1.explain

| 12345678 | +----+-------------+-------+--------+-------------------------------------+-------------------------+---------+--------------------------+------+----------------------------------------------+ \| id \| select_type \| table \| type  \| possible_keys            \| key           \| key_len \| ref           \| rows \| Extra                    \| +----+-------------+-------+--------+-------------------------------------+-------------------------+---------+--------------------------+------+----------------------------------------------+ \| 1 \| SIMPLE   \| oei  \| ref  \| idx_category_left_right,idx_data_id \| idx_category_left_right \| 5    \| const          \| 8849 \| Using where; Using temporary; Using filesort \| \| 1 \| SIMPLE   \| bu  \| ref  \| PRIMARY,idx_userid_status      \| idx_userid_status    \| 4    \| meituancrm.oei.data_id  \|  76 \| Using where; Using index           \| \| 1 \| SIMPLE   \| cb  \| ref  \| idx_branch_id,idx_contact_branch_id \| idx_branch_id      \| 4    \| meituancrm.bu.branch_id \|  1 \|                       \| \| 1 \| SIMPLE   \| c   \| eq_ref \| PRIMARY               \| PRIMARY         \| 108   \| meituancrm.cb.contact_id \|  1 \|                       \| +----+-------------+-------+--------+-------------------------------------+-------------------------+---------+--------------------------+------+----------------------------------------------+ |
| -------- | ------------------------------------------------------------ |
|          |                                                              |

从执行计划上看，mysql先查org_emp_info表扫描8849记录，再用索引idx_userid_status关联branch_user表，再用索引idx_branch_id关联contact_branch表，最后主键关联contact表。

rows返回的都非常少，看不到有什么异常情况。我们在看一下语句，发现后面有order by + limit组合，会不会是排序量太大搞的？于是我们简化SQL，去掉后面的order by 和 limit，看看到底用了多少记录来排序



MySQL

| 12345678910111213141516171819202122232425 | **select**  **count**(*) **from**  contact c  **inner join**  contact_branch cb     **on** c.id = cb.contact_id  **inner join**  branch_user bu     **on** cb.branch_id = bu.branch_id     **and** bu.**status** **in** (     1,    2)   **inner join**    org_emp_info oei      **on** oei.data_id = bu.user_id      **and** oei.node_left >= 2875      **and** oei.node_right <= 10802      **and** oei.org_category = - 1  +----------+ \| count(*) \| +----------+ \|  778878 \| +----------+ 1 **row** **in** **set** (5.19 sec) |
| ----------------------------------------- | ------------------------------------------------------------ |
|                                           |                                                              |

发现排序之前居然锁定了778878条记录，如果针对70万的结果集排序，将是灾难性的，怪不得这么慢，那我们能不能换个思路，先根据contact的created_time排序，再来join会不会比较快呢？

于是改造成下面的语句，也可以用straight_join来优化



MySQL

| 12345678910111213141516171819202122232425262728293031323334353637383940414243 | **select** c.id, c.name, c.position, c.sex, c.phone, c.office_phone, c.feature_info, c.birthday, c.creator_id, c.is_keyperson, c.giveup_reason, c.**status**, c.data_source, from_unixtime(c.created_time) **as** created_time, from_unixtime(c.last_modified) **as** last_modified, c.last_modified_user_id **from** contact c **where** **exists** ( **select** 1 **from** contact_branch cb **inner join** branch_user bu **on** cb.branch_id = bu.branch_id **and** bu.**status** **in** ( 1, 2) **inner join** org_emp_info oei **on** oei.data_id = bu.user_id **and** oei.node_left >= 2875 **and** oei.node_right <= 10802 **and** oei.org_category = - 1 **where** c.id = cb.contact_id ) **order by** c.created_time **desc** **limit** 0 , 10; |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
|                                                              |                                                              |

验证一下效果 预计在1ms内，提升了13000多倍！

| 1    | 10 rows **in** set (0.00 sec) |
| ---- | ----------------------------- |
|      |                               |

本以为至此大工告成，但我们在前面的分析中漏了一个细节，先排序再join和先join再排序理论上开销是一样的，为何提升这么多是因为有一个limit！大致执行过程是：mysql先按索引排序得到前10条记录，然后再去join过滤，当发现不够10条的时候，再次去10条，再次join，这显然在内层join过滤的数据非常多的时候，将是灾难的，极端情况，内层一条数据都找不到，mysql还傻乎乎的每次取10条，几乎遍历了这个数据表！

用不同参数的SQL试验下



MySQL

| 1234567891011121314151617181920212223242526272829303132333435363738394041424344 | **select**  **sql_no_cache**  c.id,  c.name,  c.position,  c.sex,  c.phone,  c.office_phone,  c.feature_info,  c.birthday,  c.creator_id,  c.is_keyperson,  c.giveup_reason,  c.**status**,  c.data_source,  from_unixtime(c.created_time) **as** created_time,  from_unixtime(c.last_modified) **as** last_modified,  c.last_modified_user_id   **from**  contact c   **where**  **exists** (    **select**     1        **from**     contact_branch cb         **inner join**     branch_user bu                  **on** cb.branch_id = bu.branch_id                  **and** bu.**status** **in** (        1,       2)             **inner join**       org_emp_info oei                      **on** oei.data_id = bu.user_id                      **and** oei.node_left >= 2875                      **and** oei.node_right <= 2875                      **and** oei.org_category = - 1             **where**       c.id = cb.contact_id          )      **order by**    c.created_time **desc** **limit** 0 ,    10; Empty **set** (2 min 18.99 sec) |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
|                                                              |                                                              |

2 min 18.99 sec！比之前的情况还糟糕很多。由于mysql的nested loop机制，遇到这种情况，基本是无法优化的。这条语句最终也只能交给应用系统去优化自己的逻辑了。

通过这个例子我们可以看到，并不是所有语句都能优化，而往往我们优化时，由于SQL用例回归时落掉一些极端情况，会造成比原来还严重的后果。所以，第一：不要指望所有语句都能通过SQL优化，第二：不要过于自信，只针对具体case来优化，而忽略了更复杂的情况。

慢查询的案例就分析到这儿，以上只是一些比较典型的案例。我们在优化过程中遇到过超过1000行，涉及到16个表join的“垃圾SQL”，也遇到过线上线下数据库差异导致应用直接被慢查询拖死，也遇到过varchar等值比较没有写单引号，还遇到过笛卡尔积查询直接把从库搞死。再多的案例其实也只是一些经验的积累，如果我们熟悉查询优化器、索引的内部原理，那么分析这些案例就变得特别简单了。

```





首先来看这个面试题：
已知表t是innodb引擘，有主键：id（int类型) ，下面3条语句是否加锁？加锁的话，是什么锁？

1. select * from t where id=X;
2. begin;select * from t where id=X;
3. begin;select * from t where id=X for update;

这里其实有很多依赖条件，并不能一开始就给出一个很确定的答复。我们一层层展开来分析。

# 1 MySQL有哪些锁？

![8275bf1ffeb04bbc9da3ca4cd89268de.png](/Users/giffinhao/Downloads/笔记/pic/8275bf1ffeb04bbc9da3ca4cd89268de.png.webp)

首先要知道MySQL有哪些锁，如上图所示，至少有12类锁（其中**自增锁**是事务向包含了AUTO_INCREMENT列的表中新增数据时会持有，**predicate locks for spatial index** 为空间索引专用，本文不讨论这2类锁）。

锁按**粒度**可分为分为**全局，表级，行级**3类。

## **1.1 全局锁**

对整个数据库实例加锁。
加锁表现：数据库处于只读状态，阻塞对数据的所有DML/DDL
加锁方式：**Flush tables with read lock** 释放锁：**unlock tables**(发生异常时会自动释放)
作用场景：全局锁主要用于做数据库实例的逻辑备份，与设置数据库只读命令**set global readonly=true**相比，全局锁在发生异常时会自动释放

## 1.2 表锁

对操作的整张表加锁， 锁定颗粒度大，资源消耗少，不会出现死锁，但会导致写入并发度低。具体又分为3类：
**1）显式表锁：**分为共享锁（S）和排他锁（X）
显示加锁方式：**lock tables ... read/write**
释放锁：**unlock tables**(连接中断也会自动释放)
**2）Metadata-Lock（元数据锁）**：MySQL5.5版本开始引入，主要功能是并发条件下，防止session1的查询事务未结束的情况下，session2对表结构进行修改，保护元数据的一致性。在session1持有 metadata-lock的情况下，session2处于等待状态：show processlist可见**Waiting for table metadata lock**。其具体加锁机制如下：

- DML->先加MDL 读锁（SHARED_READ，SHARED_WRITE）
- DDL->先加MDL 写锁（EXCLUSIVE）
- 读锁之间兼容
- 读写锁之间、写锁之间互斥

**3）Intention Locks（意向锁）**：**意向锁**为表锁（表示为IS或者IX），由存储引擎自己维护，用户无法干预。下面举一个例子说明其功能：
假设有2个事务：T1和T2
T1: 锁住表中的一行，只能读不能写（行级读锁）。
T2: 申请整个表的写锁（表级写锁）。
如T2申请成功，则能任意修改表中的一行，但这与T1持有的行锁是冲突的。故数据库应识别这种冲突，让T2的锁申请被阻塞，直到T1释放行锁。
有2种方法可以实现冲突检测：

1. 判断表是否已被其它事务用表锁锁住。
2. 判断表中的每一行是否已被行锁锁住。

其中2需要遍历整个表，**效率太低**。因此innodb使用意向锁来解决这个问题：
T1需要先申请表的**意向共享锁（IS）**，成功后再申请某一行的**记录锁S**。
在意向锁存在的情况下，上面的判断可以改为：
T2发现表上有意向共享锁IS，因此申请表的写锁被阻塞。

##  **1.3 行锁**

InnoDB引擘支持行级别锁，行锁粒度小，并发度高，但加锁开销大，也可能会出现死锁。
加锁机制：innodb行锁锁住的是索引页，回表时，主键的聚簇索引也会加上锁。

![51e121fb018d4aba9ae85e43d2a190ff.png](/Users/giffinhao/Downloads/笔记/pic/51e121fb018d4aba9ae85e43d2a190ff.png.webp)

行锁具体类别如上图所示，包括：**Record lock/Gap Locks/Next-Key Locks**，每类又可分为**共享锁（S）**或者**排它锁（X）**，一共2*3=6类，最后还有1类插入意向锁：
**Record lock（记录锁）：**最简单的行锁，仅仅锁住一行。记录锁永远都是加在索引上的，即使一个表没有索引，InnoDB也会隐式的创建一个索引，并使用这个索引实施记录锁。
**Gap Locks（间隙锁）：**加在两个索引值之间的锁，或者加在第一个索引值之前，或最后一个索引值之后的间隙。使用间隙锁锁住的是一个区间，而不仅仅是这个区间中的每一条数据。间隙锁只阻止其他事务插入到间隙中，不阻止其它事务在同一个间隙上获得间隙锁，所以 gap x lock 和 gap s lock 有相同的作用。它是一个**左开右开**区间：如（1，3）
**Next-Key Locks：\**记录锁\**和\**间隙锁\**的组合，它指的是加在某条记录以及这条记录前面间隙上的锁。它是一个左开右闭**区间：如（1，3】
**Insert Intention（插入意向锁）**：该锁只会出现在insert操作执行前（并不是所有insert操作都会出现），目的是为了提高并发插入能力。它在插入一行记录操作之前设置一种特殊的间隙锁，多个事务在相同的索引间隙插入时，如果不是插入间隙中相同的位置就不需要互相等待。

**TIPS:**

1.不存在unlock tables … read/write，只有unlock tables
2.If a session begins a transaction, an implicit UNLOCK TABLES is performed

# 2 锁的兼容情况

引入意向锁后，表锁之间的兼容性情况如下表：

![66ec224bc89f4cb4a173caf4bec23b8f.png](/Users/giffinhao/Downloads/笔记/pic/66ec224bc89f4cb4a173caf4bec23b8f.png.webp)

总结：

1. 意向锁之间都兼容
2. X,IX和其它都不兼容（除了1）
3. S,IS和其它都兼容（除了1,2）

# 3 锁信息查看方式

- MySQL 5.6.16版本之前，需要建立一张特殊表innodb_lock_monitor，然后使用**show engine innodb status**查看

CREATE TABLE innodb_lock_monitor (a INT) ENGINE=INNODB;

DROP TABLE innodb_lock_monitor;

- MySQL 5.6.16版本之后，修改系统参数innodb_status_output后，使用**show engine innodb status**查看

set GLOBAL innodb_status_output=ON;

set GLOBAL innodb_status_output_locks=ON;

每15秒输出一次INNODB运行状态信息到错误日志

- MySQL5.7版本之后

可以通过information_schema.innodb_locks查看事务的锁情况，但只能看到阻塞事务的锁；如果事务并未被阻塞，则在该表中看不到该事务的锁情况

- MySQL8.0

删除information_schema.innodb_locks，添加performance_schema.data_locks，即使事务并未被阻塞，依然可以看到事务所持有的锁，同时通过performance_schema.table_handles、performance_schema.metadata_locks可以非常方便的看到元数据锁等表锁。

# 4 测试环境搭建

## 4.1 建立测试表

该表包含一个主键，一个唯一键和一个非唯一键：

CREATE TABLE `t` (

`id` int(11) NOT NULL,

`a` int(11) DEFAULT NULL,

`b` int(11) DEFAULT NULL,

`c` varchar(10),

PRIMARY KEY (`id`),

unique KEY `a` (`a`),

key b(b))

ENGINE=InnoDB;

## 4.2 写入测试数据

insert into t values(1,10,100,'a'),(3,30,300,'c'),(5,50,500,'e');

# 5 记录存在时的加锁

对于innodb引擘来说，加锁的2个决定因素：

1）当前的**事务隔离级别**
2）当前**记录是否存在**

假设id为3的记录存在，则在不同的4个隔离级别下3个语句的加锁情况汇总如下表(select 3表示 select * from t where id=3)：

| **隔离级别** | **select 3** | **begin;select 3**              | **begin;select 3 for update**    |
| ------------ | ------------ | ------------------------------- | -------------------------------- |
| **RU**       | 无           | SHARED_READ                     | SHARED_WRITE IX X,REC_NOT_GAP：3 |
| **RC**       | 无           | SHARED_READ                     | SHARED_WRITE IX X,REC_NOT_GAP：3 |
| **RR**       | 无           | SHARED_READ                     | SHARED_WRITE IX X,REC_NOT_GAP：3 |
| **Serial**   | 无           | SHARED_READ IS S,REC_NOT_GAP：3 | SHARED_WRITE IX X,REC_NOT_GAP：3 |

分析：

1. 使用以下语句在4种隔离级别之间切换：
   set global transaction_isolation='READ-UNCOMMITTED';
   set global transaction_isolation='READ-COMMITTED';
   set global transaction_isolation='REPEATABLE-READ';
   set global transaction_isolation='Serializable';
2. 对于auto commit=true，select 没有显式开启事务（begin）的语句，元数据锁和行锁都不加，是真的“**读不加锁**”
3. 对于begin; select ... where id=3这种只读事务，会加**元数据锁SHARED_READ**，防止事务执行期间表结构变化，查询**performance_schema.metadata_locks**表可见此锁：

![4b66ac8cd6b8462f8226bfe0aa667367.png](/Users/giffinhao/Downloads/笔记/pic/4b66ac8cd6b8462f8226bfe0aa667367.png.webp)

1. 对于begin; select ... where id=3这种只读事务，MySQL在RC和RR隔离级别下，使用MVCC快照读，不加行锁，但在Serial隔离级别下，读写互斥，会加**意向共享锁（表锁）**和**共享记录锁（行锁）**
2. 对于begin; select ... where id=3 for update，会加**元数据锁SHARED_WRITE**
3. 对于begin; select ... where id=3 or update，4种隔离级别都会加**意向排它锁（表锁）**和**排它记录锁（行锁）,**查询**performance_schema.data_locks**可见此2类锁

![f9f0089c36514b828e896769a9f6f6ad.png](/Users/giffinhao/Downloads/笔记/pic/f9f0089c36514b828e896769a9f6f6ad.png.webp)

# 6 记录不存在时的加锁

| **隔离级别** | **select 2** | **begin;select 2**      | **begin;select 2 for update** |
| ------------ | ------------ | ----------------------- | ----------------------------- |
| RU           | 无           | SHARED_READ             | SHARED_WRITE IX               |
| RC           | 无           | SHARED_READ             | SHARED_WRITE IX               |
| RR           | 无           | SHARED_READ             | SHARED_WRITE IX X,GAP：3      |
| Serial       | 无           | SHARED_READ IS S,GAP：3 | SHARED_WRITE IX X,GAP：3      |

分析：

1. 当记录不存在的时候，RU和RC隔离级别只有意向锁，没有行锁了
2. RR，Serial隔离级别下，记录锁变成了**Gap Locks（间隙锁），**可以防止**幻读，**lock_data为3的GAP lock锁住区间（1，3），此时ID=2的记录插入会被阻塞。

![e8e4c6c5deff458896d6f0b2629727c4.png](/Users/giffinhao/Downloads/笔记/pic/e8e4c6c5deff458896d6f0b2629727c4.png.webp)

# 1 构造测试环境

该表包含一个主键，一个唯一键和一个非唯一键，有3条测试记录：
CREATE TABLE `t` (
`id` int(11) NOT NULL,
`a` int(11) DEFAULT NULL,
`b` int(11) DEFAULT NULL,
`c` varchar(10),
PRIMARY KEY (`id`),
unique KEY `a` (`a`),
key b(b))
ENGINE=InnoDB;

insert into t values(1,10,100,'a'),(3,30,300,'c'),(5,50,500,'e');

# 2 主键范围读取

## 2.1 RR隔离级别

begin;
select * from t where id>1 and id<7 for update;

![343055ea1cef4b379127ce2ec7bf4684.jpeg](/Users/giffinhao/Downloads/笔记/pic/343055ea1cef4b379127ce2ec7bf4684.jpeg.webp)

![869e0ad7098b4a67948de7c38ed22916.jpeg](/Users/giffinhao/Downloads/笔记/pic/869e0ad7098b4a67948de7c38ed22916.jpeg.webp)

1. 原则1：innodb行锁锁住的是索引页
2. 原则2：索引查找过程中访问到的对象会加锁
3. 原则3：RR隔离级别为了防止幻读，存在间隙锁（GAP LOCK）
4. 原则4：加锁的基本单位是 next-key lock，next-key lock 是前开后闭区间
5. 所以加了3个X锁（锁定记录本身和之前的区间，等于间隙锁+行锁），分别锁定(1,3】，(3,5】，(5,+∞】区间
   **说明：
   **1）InnoDB 给每个索引加了一个不存在的最大值 supremum，代表+∞
   2）幻读：当某个事务在读取某个范围内的记录时，另一个事务又在该范围内插入了新的记录，当之前的事务再次读取该范围的记录时，会产生幻读。

# 3 唯一索引等值查询

## 3.1 RR隔离级别

begin;
select * from t where a=30 for update;

![be909e91265e41819a1d2534300f4642.jpeg](/Users/giffinhao/Downloads/笔记/pic/be909e91265e41819a1d2534300f4642.jpeg.webp)



1. 原则1：**innodb行锁锁住的是索引页，回表时，主键的聚簇索引也会加上锁。**
2. 原则2：二级索引（非聚簇索引）的叶子节点包含了引用行的主键值。
3. 原则3：索引上的等值查询，给唯一索引加锁的时候，next-key lock 退化为行锁。
4. 所以加了2个记录锁，记录锁30，3代表锁定唯一索引a上的（id=3,a=30）这条记录，记录锁3代表锁定了主键上的id=3这条记录

## 3.2 RC隔离级别

begin;
select * from t where a=30 for update;

![be909e91265e41819a1d2534300f4642.jpeg](/Users/giffinhao/Downloads/笔记/pic/be909e91265e41819a1d2534300f4642.jpeg.webp)

1. 对于该条语句，RC隔离级别下加锁完全一样

# 4 非唯一索引等值查询

## 4.1 RR隔离级别

begin;
select * from t where b=300 for update;

![59ed5feb79d14e8eb4de411fe450c598.jpeg](/Users/giffinhao/Downloads/笔记/pic/59ed5feb79d14e8eb4de411fe450c598.jpeg.webp)

1. 原则：索引上的等值查询，向右遍历时且最后一个值不满足等值条件的时候，next-key lock 退化为间隙锁。
2. 所以对于非唯一索引b，锁定了((b=100,id=1),(b=300,id=3)】区间和((b=300,id=3),(b=500,id=5))区间和主键上的id=3

begin;
select * from t where b=400 for update;

![59ed5feb79d14e8eb4de411fe450c598.jpeg](/Users/giffinhao/Downloads/笔记/pic/59ed5feb79d14e8eb4de411fe450c598.jpeg.webp)

1. 可以看到，查询的值b=400不存在，但加锁情况和b=300值存在的时候是一样的

## 4.2 RC隔离级别

begin;
select * from t where b=300 for update;

![7ebf12a9abd84c9f9865664edfbe4521.jpeg](/Users/giffinhao/Downloads/笔记/pic/7ebf12a9abd84c9f9865664edfbe4521.jpeg.webp)

1. 原则：读提交隔离级别 (read-committed) 只有行锁，没有间隙锁
2. 所以只锁定了锁引b上的(b=300,id=3)和主键上的id=3

begin;
select * from t where b=400 for update;

![4aa58f8b308c42e5b4a3440e9c47d03d.jpeg](/Users/giffinhao/Downloads/笔记/pic/4aa58f8b308c42e5b4a3440e9c47d03d.jpeg.webp)

1. 因为RC隔离级别没有间隙锁，所以b=400值不存在的时候，只有IX意向排它锁。

# 5 非唯一索引加覆盖索引

## 5.1 RR隔离级别

select id from t where b=300 lock in share mode;

![7bbad7602ca5454cba890b49696173bd.jpeg](/Users/giffinhao/Downloads/笔记/pic/7bbad7602ca5454cba890b49696173bd.jpeg.webp)

1. 原则：如果一个索引包含所有需要查询的字段的值，就是覆盖索引，对于二级索引来说，可以避免对主键索引的查询（回表）
2. 因为二级索引b包括(b,id)，所以主键索引上无锁
3. 因为是lock in share mode所以加的是共享锁（S）和共享意向锁（IS）

# 6 无索引

begin;
select * from t where c='aa' for update;

![85cd5db70d32469d880a187ec9a25e42.jpeg](/Users/giffinhao/Downloads/笔记/pic/85cd5db70d32469d880a187ec9a25e42.jpeg.webp)

1. 没有索引的时候，要全表扫描，有主键就扫主键
2. 所以锁定范围：(-∞,1]、(1,3]、(3,5]、(5,+supremum]，可以看出来把整张表都锁住了，所以对于实时业务一定要避免非索引查询
