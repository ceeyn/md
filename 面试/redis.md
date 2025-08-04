string：set key val nx ex 10， setex key val， setnx， get，del， incry， decry

hash：hset，hget，hincryby，hexists，hgetall，hdel

list：lpush，lpop，lrange key start  stop

set：sadd，srem，sismember

sortedset：zadd，zrem，zrange key start stop，zrangebyscore key min max，zincryby key increment member，

TTL key  SCAN cursor













好的，我将更详细地解释 Redis 各个数据结构的用法，结合实际业务场景并提供具体的 Redis 命令操作和用例。

### 1. **String（字符串）**

**字符串是 Redis 中最基础的数据结构。它可以存储文本、整数、浮点数，甚至二进制数据。一个字符串的最大存储长度为 512 MB。**

#### **使用场景**

- **计数器**：可以用字符串类型存储访问计数、点击量等。
- **缓存**：将复杂对象的序列化结果（如 JSON、XML）存储在 Redis 中。
- **配置项**：可以用来存储全局配置，比如缓存的过期时间、开关等。

#### **示例**

- **存储和获取字符串**

```bash
SET user:name "Alice"  # 存储字符串
GET user:name           # 获取存储的字符串，返回 "Alice"
```

- **自增和自减**

```bash
SET page:visits 100      # 设置初始访问量为 100
INCR page:visits         # 访问一次页面，计数器加 1，返回 101
INCRBY page:visits 5     # 增加 5 次访问，返回 106
DECR page:visits         # 计数器减 1，返回 105
```

- **存储和获取整数、浮点数**

```bash
SET price:apple 1.99      # 设置苹果的价格为 1.99
GET price:apple           # 返回 "1.99"

SET balance 100           # 设置余额为 100
INCRBY balance 50         # 充值 50，返回 150
DECRBY balance 30         # 扣款 30，返回 120
```

#### **进阶用法**

- **设置键的过期时间**

```bash
SET session:token "abc123" EX 3600  # 设置 session，1小时后过期
TTL session:token                   # 查看过期时间，返回剩余秒数
```

### 2. **Hash（哈希）**

**Hash 是一个键值对的集合，适合存储对象。与 String 不同的是，Hash 可以在一个键中包含多个字段和值。每个字段和值本质上仍然是字符串。**

#### **使用场景**

- **用户信息存储**：存储用户资料，如用户名、年龄、邮箱等。
- **商品属性**：存储商品信息，如价格、库存、描述等。
- **配置项集合**：将一组相关配置项存储在一个哈希表中。

#### **示例**

- **存储和获取用户信息**

```bash
HSET user:1001 name "Alice"    # 设置用户 1001 的姓名为 Alice
HSET user:1001 age 30          # 设置用户 1001 的年龄为 30
HSET user:1001 email "alice@example.com"  # 设置用户 1001 的邮箱

HGET user:1001 name            # 获取用户 1001 的姓名，返回 "Alice"
HGETALL user:1001              # 获取用户 1001 的所有信息，返回 {"name": "Alice", "age": "30", "email": "alice@example.com"}
```

- **修改和增加字段**

```bash
HINCRBY user:1001 age 1        # 增加用户 1001 的年龄 1 岁，返回 31
HSETNX user:1001 email "new_email@example.com"  # 如果邮箱不存在，设置新邮箱
```

- **删除字段**

```bash
HDEL user:1001 email           # 删除用户 1001 的邮箱字段
HEXISTS user:1001 email        # 检查邮箱字段是否存在，返回 0（不存在）
```

### 3. **List（列表）**

**List 是一个双向链表，可以在列表的头部或尾部添加和删除元素。它特别适合用来实现队列（FIFO）和栈（LIFO）。**

#### **使用场景**

- **消息队列**：将任务、消息存储在列表中，按顺序处理。
- **社交动态**：存储用户的动态或操作历史。
- **评论系统**：按时间顺序存储评论。

#### **示例**

- **任务队列**

```bash
LPUSH tasks "task1"        # 将任务 task1 加入队列
LPUSH tasks "task2"        # 将任务 task2 加入队列（在队头）
RPUSH tasks "task3"        # 将任务 task3 加入队列尾部

LPOP tasks                 # 处理队列中最先加入的任务，返回 "task2"
RPOP tasks                 # 处理队列中最后加入的任务，返回 "task3"
```

- **获取列表中的元素**

```bash
LRANGE tasks 0 -1          # 获取列表中的所有任务，返回 ["task1"]
LINDEX tasks 0             # 获取第一个任务，返回 "task1"
LLEN tasks                 # 获取列表的长度，返回 1
```

- **修剪列表**

```bash
LTRIM tasks 0 2            # 只保留列表的前三个元素，超出部分删除
```

### 4. **Set（集合）**

**Set 是一个无序集合，集合中的元素是唯一的，不允许重复。Set 可以实现去重操作，适合用于需要快速查找的场景。**

#### **使用场景**

- **标签系统**：存储文章、博客的标签。
- **社交网络**：存储用户的好友、关注列表。
- **去重操作**：存储唯一的值，如访问 IP、用户 ID。

#### **示例**

- **存储用户的兴趣标签**

```bash
SADD user:1001:tags "redis"    # 添加 redis 标签
SADD user:1001:tags "nosql"    # 添加 nosql 标签
SADD user:1001:tags "database" # 添加 database 标签

SMEMBERS user:1001:tags        # 获取用户的所有标签，返回 ["redis", "nosql", "database"]
SISMEMBER user:1001:tags "nosql"  # 检查用户是否有 nosql 标签，返回 1（存在）
```

- **集合运算**

```bash
SADD set1 "a" "b" "c"           # 创建集合 set1，包含元素 a, b, c
SADD set2 "b" "c" "d"           # 创建集合 set2，包含元素 b, c, d

SDIFF set1 set2                 # 获取 set1 与 set2 的差集，返回 ["a"]
SINTER set1 set2                # 获取 set1 与 set2 的交集，返回 ["b", "c"]
SUNION set1 set2                # 获取 set1 与 set2 的并集，返回 ["a", "b", "c", "d"]
```

### 5. **Sorted Set（有序集合）**

**Sorted Set 是一个带有分数的集合。与 Set 不同的是，Sorted Set 中的元素是有序的，元素按照分数排序。它适合实现排行榜、优先级队列等场景。**

#### **使用场景**

- **排行榜**：按分数、积分、时间排序。
- **优先级队列**：按任务优先级处理任务。
- **带权重的标签系统**：按使用频率排序标签。

#### **示例**

- **添加和获取排行榜**

```bash
ZADD leaderboard 100 "Alice"   # 添加 Alice，分数为 100
ZADD leaderboard 150 "Bob"     # 添加 Bob，分数为 150
ZADD leaderboard 120 "Charlie" # 添加 Charlie，分数为 120

ZRANGE leaderboard 0 -1 WITHSCORES  # 获取所有用户及其分数，按分数升序排列，返回 ["Alice", 100, "Charlie", 120, "Bob", 150]
ZRANK leaderboard "Charlie"         # 获取 Charlie 的排名，返回 1（从 0 开始计数）
```

- **倒序获取排名**

```bash
ZREVRANGE leaderboard 0 -1 WITHSCORES  # 获取所有用户及其分数，按分数降序排列，返回 ["Bob", 150, "Charlie", 120, "Alice", 100]
ZREVRANK leaderboard "Alice"           # 获取 Alice 的倒序排名，返回 2（从 0 开始计数）
```

- **按分数范围获取**

```bash
ZRANGEBYSCORE leaderboard 100 150  # 获取分数在 100 到 150 之间的用户，返回 ["Alice", "Charlie", "Bob"]
```

### 6. **Stream（流）**

**Stream 是一种强大的数据结构，适合处理消息队列、事件流等。它支持消费者组，可以实现消费和消息确认机制。**

#### **使用场景**

- **实时日志系统**：收集和处理日志数据。
- **消息队列**：异步处理任务，确保消息不丢失。
- **事件驱动系统**：处理用户行为

、传感器数据等实时事件。

#### **示例**

- **添加和读取消息**

```bash
XADD mystream * sensor-id 1234 temperature 19.8  # 添加传感器数据到流
XRANGE mystream - +                              # 获取所有流中的消息
```

- **消费者组**

```bash
XGROUP CREATE mystream mygroup $                # 创建消费者组
XREADGROUP GROUP mygroup consumer1 STREAMS mystream >  # 从消费者组读取消息
XACK mystream mygroup 1526569495631-0           # 确认消息已处理
```

### 7. **Bitmap（位图）**

**Bitmap 是一个位数组，用来表示二进制数据。每个位可以存储 0 或 1，它可以高效地存储布尔值。**

#### **使用场景**

- **签到系统**：记录用户每天是否签到。
- **活跃用户统计**：记录用户每天的活跃状态。
- **权限控制**：用位图来表示不同的权限位。

#### **示例**

- **用户签到**

```bash
SETBIT user:1001:sign 5 1    # 设置用户 1001 在第 5 天签到
GETBIT user:1001:sign 5      # 查询用户 1001 在第 5 天是否签到，返回 1
BITCOUNT user:1001:sign      # 统计用户 1001 签到的天数，返回 1
```

- **统计用户活跃天数**

```bash
SETBIT active:20230801 1 1   # 设置用户 1 在 2023-08-01 活跃
SETBIT active:20230802 2 1   # 设置用户 2 在 2023-08-02 活跃
BITCOUNT active:20230801     # 统计 2023-08-01 的活跃用户数量，返回 1
```

### 8. **HyperLogLog**

**HyperLogLog 是一种概率性数据结构，用于估算集合的基数（即不同元素的数量）。它占用极少的内存，但结果是近似的。**

#### **使用场景**

- **网站独立访问量统计**：统计不同 IP 的访问量。
- **去重统计**：统计不重复的元素数量，如用户 ID、设备 ID。

#### **示例**

- **统计独立访问用户数**

```bash
PFADD uv user1           # 添加访问用户 user1
PFADD uv user2           # 添加访问用户 user2
PFADD uv user1           # 添加重复访问用户 user1
PFCOUNT uv               # 获取独立访问用户数，返回 2（user1 和 user2）
```

### 9. **Geospatial Indexes（地理位置）**

**Geospatial Indexes 支持存储和查询地理位置数据，可以用于地理范围查询和计算距离。**

#### **使用场景**

- **附近的商店**：查询用户附近的商店位置。
- **位置服务**：根据用户位置提供个性化服务。

#### **示例**

- **存储和查询地理位置**

```bash
GEOADD places 13.361389 38.115556 "Palermo"   # 添加 Palermo 的经纬度
GEOADD places 15.087269 37.502669 "Catania"   # 添加 Catania 的经纬度
GEORADIUS places 15 37 200 km                # 查询200公里范围内的地点，返回 ["Catania"]
```

- **计算两地距离**

```bash
GEODIST places "Palermo" "Catania" km        # 计算 Palermo 和 Catania 的距离，返回 166.2742 km
```

---

通过详细的示例和使用场景，可以更好地理解 Redis 各个数据结构的作用和应用场景。根据业务需求选择合适的数据结构，可以有效提高系统性能和简化开发复杂度。





在 Redis 中，`hash` 数据结构类似于一个键值对的集合，具体到你提供的 Lua 脚本，`hash` 结构用于存储与某个秒杀活动相关的信息。以下是基于脚本中 `seckill:<voucherId>` 这个哈希表的具体结构。

### Hash 结构示例

假设 `voucherId` 为 `123`，则 Redis 中的哈希键名为 `seckill:123`。该哈希表的结构及其字段可以描述如下：

**Redis Key**: `seckill:123`

| **Field** | **Value Type**             | **Description**                            | **Example Value**                                    |
| --------- | -------------------------- | ------------------------------------------ | ---------------------------------------------------- |
| `begin`   | String (Integer as String) | 秒杀活动的开始时间戳（通常为 Unix 时间戳） | `"1693027200"` (对应于 UTC 时间 2023-08-26 00:00:00) |
| `end`     | String (Integer as String) | 秒杀活动的结束时间戳（通常为 Unix 时间戳） | `"1693113599"` (对应于 UTC 时间 2023-08-26 23:59:59) |
| `stock`   | String (Integer as String) | 当前剩余库存数量                           | `"100"` (表示当前库存数量为 100)                     |

### 示例：使用 Redis 命令查看和修改该 Hash 结构

#### 查看 Hash 结构中的所有字段和值

```bash
HGETALL seckill:123
```

**输出**:
```
1) "begin"
2) "1693027200"
3) "end"
4) "1693113599"
5) "stock"
6) "100"
```

#### 获取特定字段的值

获取秒杀活动的开始时间：

```bash
HGET seckill:123 begin
```

**输出**:
```
"1693027200"
```

#### 修改特定字段的值

减少库存：

```bash
HINCRBY seckill:123 stock -1
```

**输出**:
```
"99"  # 表示库存减少了1，当前库存为99
```

### 该 Hash 结构的使用场景

- **`begin` 和 `end` 字段**: 用于判断秒杀活动是否在进行中。根据脚本中的逻辑，在判断 `timestamp` 是否在 `begin` 和 `end` 之间，来决定是否允许秒杀。
- **`stock` 字段**: 用于记录当前秒杀活动的库存。当用户参与秒杀时，库存数量会减少，防止超卖。

这个哈希结构使得存储与秒杀活动相关的多个属性变得非常高效，所有信息集中在一个 Redis 键中，减少了数据查询的复杂性。



在 Redis 中，不同的数据结构对 `key` 和 `value` 的类型有不同的要求。下面我将详细列出 Redis 中每个主要数据结构的 `key` 和 `value` 类型。

### 1. **String（字符串）**

- **Key 类型**: `String`
  - `key` 是一个字符串，用来标识存储的值。它可以是任意的文本内容，比如用户 ID、配置项名称等。
  
- **Value 类型**: `String`
  - `value` 也是一个字符串。可以存储纯文本、数字（整型或浮点型）、二进制数据（如图片、音频）等。

**示例**:

```bash
SET mykey "Hello, World!"  # key 是 "mykey"，value 是 "Hello, World!"
GET mykey                  # 返回 "Hello, World!"
```

### 2. **Hash（哈希）**

- **Key 类型**: `String`
  - `key` 是哈希表的名称，标识存储的哈希表。

- **Field 类型**: `String`
  - `field` 是哈希表中某个字段的名称，相当于哈希表中的一个键。

- **Value 类型**: `String`
  - `value` 是哈希表中某个字段的值。与 `String` 类型类似，它也可以是纯文本、数字或其他类型的字符串。

**示例**:

```bash
HSET user:1001 name "Alice"    # key 是 "user:1001"，field 是 "name"，value 是 "Alice"
HSET user:1001 age "30"        # key 是 "user:1001"，field 是 "age"，value 是 "30"
HGET user:1001 name            # 返回 "Alice"
```

### 3. **List（列表）**

- **Key 类型**: `String`
  - `key` 标识列表的名称，表示存储的列表。

- **Value 类型**: `String`
  - `value` 是列表中的一个元素。每个元素是一个字符串，列表中的元素可以是文本、数字等。

**示例**:

```bash
LPUSH tasks "task1"       # key 是 "tasks"，value 是 "task1"
RPUSH tasks "task2"       # key 是 "tasks"，value 是 "task2"
LPOP tasks                # 返回 "task1"
```

### 4. **Set（集合）**

- **Key 类型**: `String`
  - `key` 标识集合的名称。

- **Value 类型**: `String`
  - `value` 是集合中的一个元素。每个元素是一个字符串，集合中的元素是唯一的。

**示例**:

```bash
SADD tags "redis"        # key 是 "tags"，value 是 "redis"
SADD tags "nosql"        # key 是 "tags"，value 是 "nosql"
SMEMBERS tags            # 返回 ["redis", "nosql"]
```

### 5. **Sorted Set（有序集合）**

- **Key 类型**: `String`
  - `key` 标识有序集合的名称。

- **Value 类型**: `String`
  - `value` 是有序集合中的一个元素。每个元素是一个字符串。

- **Score 类型**: `Double`
  - `score` 是一个双精度浮点数，用于排序。每个 `value` 与一个 `score` 相关联，`Sorted Set` 会根据 `score` 排序。

**示例**:

```bash
ZADD leaderboard 100 "Alice"  # key 是 "leaderboard"，score 是 100，value 是 "Alice"
ZADD leaderboard 150 "Bob"    # key 是 "leaderboard"，score 是 150，value 是 "Bob"
ZRANGE leaderboard 0 -1       # 返回 ["Alice", "Bob"]
```

### 6. **Stream（流）**

- **Key 类型**: `String`
  - `key` 标识流的名称。

- **Field 类型**: `String`
  - `field` 是流消息中的一个字段名称。

- **Value 类型**: `String`
  - `value` 是流消息中的一个字段值。

**示例**:

```bash
XADD mystream * sensor-id 1234 temperature 19.8  # key 是 "mystream"，field 是 "sensor-id" 和 "temperature"，value 是 "1234" 和 "19.8"
XRANGE mystream - +                              # 返回流中的所有消息
```

### 7. **Bitmap（位图）**

- **Key 类型**: `String`
  - `key` 是位图的名称，标识存储的位图。

- **Value 类型**: `Bit`
  - `value` 是位图中的一个位，值为 `0` 或 `1`。

**示例**:

```bash
SETBIT user:1001:login 5 1  # key 是 "user:1001:login"，offset 是 5，value 是 1，表示用户在第 5 天登录
GETBIT user:1001:login 5    # 返回 1
```

### 8. **HyperLogLog**

- **Key 类型**: `String`
  - `key` 标识 `HyperLogLog` 的名称。

- **Value 类型**: `String`
  - `value` 是 `HyperLogLog` 中存储的一个元素。

**示例**:

```bash
PFADD visitors "user1"  # key 是 "visitors"，value 是 "user1"
PFADD visitors "user2"  # key 是 "visitors"，value 是 "user2"
PFCOUNT visitors        # 返回 2，表示去重后的用户数
```

### 9. **Geospatial（地理位置）**

- **Key 类型**: `String`
  - `key` 标识地理位置集合的名称。

- **Value 类型**: `String`
  - `value` 是地理位置集合中的一个元素，通常是地理位置的名称或标识符。

- **Longitude and Latitude 类型**: `Double`
  - 经度和纬度是双精度浮点数，用来标识地理位置的坐标。

**示例**:

```bash
GEOADD cities 13.361389 38.115556 "Palermo"   # key 是 "cities"，经度是 13.361389，纬度是 38.115556，value 是 "Palermo"
GEORADIUS cities 15 37 200 km                # 返回 200 公里范围内的地点
```

### 总结
Redis 中不同的数据结构提供了丰富的 `key-value` 存储和管理方式。`key` 在所有情况下都是一个字符串，而 `value` 的类型取决于具体的数据结构和用途。理解每种数据结构的 `key` 和 `value` 的类型有助于选择合适的数据结构以提高 Redis 在不同应用场景中的效率和性能。





在高并发场景下，限流是非常重要的技术，用于控制进入系统的请求数量，避免系统过载。Lua 脚本可以在 Redis 中实现限流策略，比如使用令牌桶（Token Bucket）或漏斗算法（Leaky Bucket）等。这些算法可以确保即使在高并发的情况下，系统也能稳定地处理请求，而不会因为瞬间的流量高峰导致系统崩溃或响应延迟过长。

### **Lua 限流的实现**

**1. 令牌桶（Token Bucket）算法**

令牌桶算法是常用的限流策略之一。该算法的核心思想是：

- 系统中存在一个“桶”，桶中存储了一定数量的“令牌”。
- 每次请求到来时，必须从桶中取出一个令牌才能执行。如果桶中没有令牌，则请求被拒绝或延迟处理。
- 令牌以固定的速率（如每秒）被加入到桶中，桶的容量是有限的，超过容量的令牌会被丢弃。

**Lua 实现令牌桶限流**：

```lua
-- 参数列表
local key = KEYS[1]            -- 限流的key
local rate = tonumber(ARGV[1]) -- 令牌发放速率，每秒的令牌数量
local capacity = tonumber(ARGV[2]) -- 桶的最大容量
local now = tonumber(ARGV[3])  -- 当前时间戳（秒）
local requested = tonumber(ARGV[4]) -- 请求的令牌数量

-- 获取当前的令牌数量和上次获取的时间
local tokens = tonumber(redis.call('get', key .. ':tokens')) or capacity
local last_refill = tonumber(redis.call('get', key .. ':timestamp')) or 0

-- 计算时间差，生成新令牌
local delta = math.max(0, now - last_refill)
local new_tokens = math.min(capacity, tokens + delta * rate)

-- 判断是否可以发放令牌
if new_tokens < requested then
    -- 令牌不足，拒绝请求
    return 0
else
    -- 令牌充足，扣减令牌并允许请求
    redis.call('set', key .. ':tokens', new_tokens - requested)
    redis.call('set', key .. ':timestamp', now)
    return 1
end
```

**使用示例**:

```bash
-- 初始化一个令牌桶限流，每秒发放5个令牌，桶容量为10
local rate = 5
local capacity = 10
local now = os.time()  -- 当前时间戳
local requested = 1    -- 每次请求1个令牌

-- 调用限流Lua脚本
redis-cli --eval limit.lua myservice , rate capacity now requested
```

在这个例子中，每当一个请求到达时，脚本会检查令牌桶中是否有足够的令牌。如果有，允许请求并减少桶中的令牌数量。如果没有足够的令牌，则拒绝请求。

**2. 漏斗算法（Leaky Bucket）**

漏斗算法与令牌桶算法类似，但它是按固定速率排放请求的。可以理解为请求进入一个漏斗，漏斗以固定速率排放请求，超出容量的请求会被丢弃。

**Lua 实现漏斗算法限流**：

```lua
-- 参数列表
local key = KEYS[1]            -- 漏斗的key
local capacity = tonumber(ARGV[1]) -- 漏斗的容量
local leak_rate = tonumber(ARGV[2]) -- 漏斗的漏出速率，多少时间排出一个请求
local now = tonumber(ARGV[3])  -- 当前时间戳（秒）

-- 获取当前漏斗中的水量和上次漏水时间
local water = tonumber(redis.call('get', key .. ':water')) or 0
local last_time = tonumber(redis.call('get', key .. ':time')) or now

-- 计算经过时间漏掉的水量
local delta = math.max(0, now - last_time)
local new_water = math.max(0, water - delta * leak_rate)

-- 判断是否可以放入新的请求
if new_water + 1 > capacity then
    -- 漏斗已满，拒绝请求
    return 0
else
    -- 漏斗未满，接受请求并增加水量
    redis.call('set', key .. ':water', new_water + 1)
    redis.call('set', key .. ':time', now)
    return 1
end
```

**使用示例**:

```bash
-- 初始化一个漏斗限流，漏斗容量为10，漏速为1秒1个请求
local capacity = 10
local leak_rate = 1
local now = os.time()

-- 调用限流Lua脚本
redis-cli --eval leakybucket.lua myservice , capacity leak_rate now
```

**总结：**

通过 Lua 脚本在 Redis 中实现限流，可以有效地控制高并发情况下的请求数量，防止系统因流量过大而崩溃。无论是令牌桶算法还是漏斗算法，它们都能很好地管理和控制请求流量，确保系统的稳定性。Lua 脚本的执行速度非常快，足以满足大多数实时性要求较高的应用场景。





漏斗算法（Leaky Bucket）的核心思想是将请求比喻为水滴，将它们倒入一个“漏斗”中，然后以固定的速率从漏斗底部流出（处理请求）。漏斗有一定的容量，**表示它能容纳的最大请求数。如果请求流入速度太快且超过了漏斗的容量，那么多余的请求将被丢弃（即请求被拒绝或延迟处理）。**

### 举例说明

假设有一个漏斗，容量为 10 个请求，漏斗以每秒处理 1 个请求的速率排放（处理）请求。

#### **场景 1：请求速率低于漏斗排放速率**

- **初始状态**: 漏斗是空的。
- **第 1 秒**: 2 个请求到达，这些请求都能进入漏斗，因为漏斗还有足够的空间（容量为 10，只有 2 个请求）。
  - 漏斗状态：2 个请求在漏斗中。
- **第 2 秒**: 漏斗排放 1 个请求，现在漏斗中剩下 1 个请求。又有 2 个请求到达，现在漏斗中有 3 个请求。
  - 漏斗状态：3 个请求在漏斗中。

在这种情况下，漏斗可以正常处理所有请求，没有请求被丢弃，因为请求的速率没有超过漏斗的排放速率和容量。

#### **场景 2：请求速率高于漏斗排放速率，但不超过漏斗容量**

- **初始状态**: 漏斗是空的。
- **第 1 秒**: 5 个请求到达，漏斗可以容纳这些请求。
  - 漏斗状态：5 个请求在漏斗中。
- **第 2 秒**: 漏斗排放 1 个请求，现在漏斗中剩下 4 个请求。又有 5 个请求到达，漏斗中现在有 9 个请求。
  - 漏斗状态：9 个请求在漏斗中。

即使请求速率较高，但只要不超过漏斗的容量，所有请求仍然可以进入漏斗并排队处理。

#### **场景 3：请求速率高于漏斗排放速率，且超过漏斗容量**

- **初始状态**: 漏斗是空的。
- **第 1 秒**: 12 个请求到达，但漏斗只能容纳 10 个请求，因此前 10 个请求进入漏斗，另外 2 个请求被丢弃。
  - 漏斗状态：10 个请求在漏斗中，漏斗已满。
- **第 2 秒**: 漏斗排放 1 个请求，现在漏斗中剩下 9 个请求。又有 12 个请求到达，但漏斗只能容纳 10 个请求，因此仅能接收 1 个新请求，剩下 11 个请求被丢弃。
  - 漏斗状态：10 个请求在漏斗中，漏斗再次满了。

在这种情况下，因请求到达速率过高，超过了漏斗的容量，导致部分请求被丢弃。这意味着系统只能处理漏斗容量内的请求，多余的请求将被拒绝或丢失。

### 总结

- **漏斗容量**: 指的是系统可以在短时间内处理或排队的最大请求数。如果请求数量超过了这个容量，那么漏斗就会“溢出”，多余的请求无法进入漏斗，被视为丢弃或拒绝。
  
- **请求排放速率**: 是漏斗处理请求的速率。即使请求大量涌入，漏斗也只能以固定速率处理请求。

- **超出容量**: 当请求到达的速度过快，超过漏斗的容纳能力时，多余的请求将被丢弃或延迟处理。这种机制可以防止系统过载，从而保持系统的稳定性和响应速度。



令牌桶算法（Token Bucket）是一种常用的流量控制算法，广泛应用于网络流量整形和限流。它的基本思想是：系统会以固定速率生成“令牌”放入一个“桶”中，每当有一个请求到达时，必须从桶中取出一个令牌才能被处理。如果桶中没有令牌了，说明当前系统处理能力已经达到上限，新的请求要么被拒绝，要么被延迟处理。

### 令牌桶算法的关键点

1. **令牌生成速率**: 令牌以固定的速率生成，例如每秒生成 5 个令牌。
2. **桶的容量**: 桶的容量是固定的，表示桶中最多能存放的令牌数量。例如桶的容量为 10 个令牌。
3. **请求处理**: 每个请求到来时，必须从桶中拿到一个令牌才能被处理。如果桶中没有足够的令牌，请求要么被拒绝，要么被延迟处理。

### 举例说明

假设我们有一个令牌桶，设定以下参数：
- **令牌生成速率**: 每秒 5 个令牌。
- **桶的容量**: 10 个令牌。
- **请求情况**: 系统接收到的请求频率可能是每秒 3 次、10 次或者更多。

#### **场景 1：请求速率低于令牌生成速率**

- **初始状态**: 桶中有 10 个令牌（已满）。
- **第 1 秒**: 有 3 个请求到达。每个请求从桶中拿走一个令牌，成功处理了 3 个请求，桶中剩余 7 个令牌。与此同时，系统以每秒 5 个令牌的速率生成新令牌，1 秒后桶中有 7 + 5 = 10 个令牌（再次满了）。
- **第 2 秒**: 再次有 3 个请求到达，每个请求成功处理，桶中剩余 7 个令牌，随后生成新令牌，桶中再次满了。

在这种情况下，请求速率较低，系统始终能及时处理所有请求，没有请求被拒绝。

#### **场景 2：请求速率高于令牌生成速率，但低于令牌生成速率与桶容量的总和**

- **初始状态**: 桶中有 10 个令牌（已满）。
- **第 1 秒**: 有 8 个请求到达，成功处理了 8 个请求，桶中剩余 2 个令牌。系统以每秒 5 个令牌的速率生成新令牌，1 秒后桶中有 2 + 5 = 7 个令牌。
- **第 2 秒**: 再次有 8 个请求到达，成功处理了 7 个请求，但桶中只有 7 个令牌，因此第 8 个请求被拒绝或延迟处理，桶中剩余 0 个令牌。1 秒后生成 5 个新令牌，桶中有 5 个令牌。

在这种情况下，请求速率略高于令牌生成速率，但桶容量能够在一定程度上缓冲请求的突发增长，因此大多数请求可以被处理，少部分请求可能会被拒绝或延迟处理。

#### **场景 3：请求速率远高于令牌生成速率**

- **初始状态**: 桶中有 10 个令牌（已满）。
- **第 1 秒**: 有 15 个请求到达，成功处理了 10 个请求，桶中剩余 0 个令牌。系统以每秒 5 个令牌的速率生成新令牌，1 秒后桶中有 0 + 5 = 5 个令牌。
- **第 2 秒**: 再次有 15 个请求到达，成功处理了 5 个请求，桶中令牌再次耗尽。桶中剩余 0 个令牌，1 秒后生成 5 个新令牌，桶中有 5 个令牌。

在这种情况下，由于请求速率远高于令牌生成速率，大量请求会被拒绝或延迟处理，只有部分请求能够被及时处理。

### **令牌桶算法的优势**

1. **灵活性**: 令牌桶算法允许突发流量的存在，只要桶中有足够的令牌，突发请求可以被立即处理。这使得系统能更好地应对不均匀的流量。
2. **精细控制**: 通过调节令牌的生成速率和桶的容量，可以精细控制系统的处理能力，确保在高负载下系统仍然稳定运行。
3. **避免过载**: 当请求速率持续高于处理能力时，多余的请求将被拒绝，从而避免系统过载。

### **Lua 实现令牌桶算法的示例**

以下是一个用 Lua 脚本在 Redis 中实现令牌桶限流的简单示例：

```lua
-- 参数列表
local key = KEYS[1]             -- 限流的key，唯一标识一个限流器
local rate = tonumber(ARGV[1])  -- 令牌生成速率，每秒生成的令牌数
local capacity = tonumber(ARGV[2]) -- 桶的最大容量
local now = tonumber(ARGV[3])   -- 当前时间戳（秒）
local requested = tonumber(ARGV[4]) -- 本次请求消耗的令牌数

-- 获取当前桶中剩余的令牌数
local tokens = tonumber(redis.call('get', key .. ':tokens')) or capacity
local last_time = tonumber(redis.call('get', key .. ':timestamp')) or 0

-- 计算从上次生成令牌到现在的时间差，增加相应数量的令牌
local delta = math.max(0, now - last_time)
local new_tokens = math.min(capacity, tokens + delta * rate)

-- 判断桶中是否有足够的令牌来处理当前请求
if new_tokens < requested then
    -- 令牌不足，拒绝请求
    return 0
else
    -- 令牌充足，扣减令牌并允许请求
    redis.call('set', key .. ':tokens', new_tokens - requested)
    redis.call('set', key .. ':timestamp', now)
    return 1
end
```

### **总结**

令牌桶算法通过控制令牌生成速率和桶的容量，实现了对系统请求处理能力的精确控制。它允许系统应对突发请求，但也能有效防止系统在持续高负载下被过载。通过 Lua 脚本将令牌桶算法集成到 Redis 中，可以在高并发环境下为应用程序提供稳定的限流机制。





# redission

### 1.hash 大key：锁名称  小key：线程名称   val：重入次数

### 2.超时释放【看门狗机制】

### 3.可重试【while+失败休息+订阅释放锁的消息后回调函数中进行重试】

```
在 Redisson 的可重入锁（Reentrant Lock）实现中，锁的核心机制借助 Redis 的 `hash` 结构和 Lua 脚本来实现分布式锁的功能。相比于本地锁（如 `synchronized` 或 `Lock`），Redisson 分布式锁的实现允许多线程进程在不同节点上通过 Redis 来协调共享资源的访问。

### 1. 可重入锁的本地实现原理
在本地锁中，`Lock` 和 `synchronized` 通过状态变量来实现可重入锁功能。

- **Lock锁：**
  `Lock` 使用一个 `volatile` 的 `state` 变量来表示锁的状态，`state` 初始为 `0`，表示没有任何线程持有锁。如果某个线程第一次获取锁，`state` 变为 `1`，如果该线程再次获取锁，`state` 递增，表明该线程多次重入锁。每当线程释放锁时，`state` 减少，当 `state` 为 `0` 时，锁完全释放。
  
- **synchronized：**
  在 JVM 的底层实现中，`synchronized` 也是通过计数器实现可重入锁。当一个线程重入锁时，计数器递增；当线程退出同步块时，计数器递减，直到归零。

### 2. Redisson 可重入锁的分布式实现
在分布式系统中，锁的状态无法依赖单个节点的内存，因此需要借助 Redis 来存储锁的状态，并通过 Lua 脚本来保证操作的原子性。Redisson 通过 Redis 的 `hash` 结构来实现锁的重入功能。

- **锁的结构：**
  在 Redis 中，Redisson 通过 `hash` 数据结构来存储锁的状态。`hash` 中有一个大 `key` 作为锁的名称，每个线程获取锁时会用一个独特的标识（由 `id` 和 `threadId` 组成）作为 `hash` 的小 `key`，表示锁的拥有者。
  
- **Lua 脚本的关键逻辑：**

  ```lua
  if (redis.call('exists', KEYS[1]) == 0) then 
      redis.call('hset', KEYS[1], ARGV[2], 1); 
      redis.call('pexpire', KEYS[1], ARGV[1]); 
      return nil; 
  end; 
  if (redis.call('hexists', KEYS[1], ARGV[2]) == 1) then 
      redis.call('hincrby', KEYS[1], ARGV[2], 1); 
      redis.call('pexpire', KEYS[1], ARGV[1]); 
      return nil; 
  end; 
  return redis.call('pttl', KEYS[1]);
```

- **参数解释：**
  - `KEYS[1]`: 锁的名称（大 `key`），表示这把锁在 Redis 中的名称。
  - `ARGV[1]`: 锁的失效时间，单位为毫秒，确保即使锁持有者异常退出，锁也能被自动释放。
  - `ARGV[2]`: 由 `id + ":" + threadId` 组成的字符串，作为 `hash` 结构中的小 `key`，用来区分不同线程持有同一锁的情况。

### 3. 执行流程
1. **锁不存在时：**
   当 `redis.call('exists', KEYS[1]) == 0`，表示锁不存在。此时，使用 `hset` 命令在 Redis 中创建一个新的 `hash` 结构，存储线程的唯一标识 `id:threadId`，并将该值设置为 `1`，表示该线程第一次获取了这把锁。接着，使用 `pexpire` 设置锁的过期时间，以防持有锁的线程异常退出时锁永不释放。

2. **锁已经存在且属于当前线程：**
   如果锁已经存在，但 `redis.call('hexists', KEYS[1], ARGV[2]) == 1` 表示当前锁是由同一线程持有的，则执行 `hincrby` 命令将该锁的计数器加 `1`，表示同一线程的可重入操作。同时，重新设置锁的过期时间，确保锁不会因过期而自动释放。

3. **锁已经存在且不属于当前线程：**
   如果锁存在但不属于当前线程，直接返回锁的剩余生存时间 `pttl`，这意味着当前线程未能获取锁。

4. **自旋等待：**
   如果线程未能抢到锁，它将进入 `while(true)` 循环，不断尝试获取锁，直到成功。

### 4. 举例说明
假设有两个线程 `A` 和 `B` 试图获取一把名为 `myLock` 的分布式锁。以下是具体的执行过程：

1. **线程 A 获取锁：**
   - `A` 调用 `tryLock()`，Redis 中 `myLock` 不存在，`A` 通过 `hset` 创建 `myLock`，并设置 `id:threadIdA` 为 `1`，表示 `A` 持有了锁。
   - `A` 开始执行业务逻辑，锁的过期时间为 `3000` 毫秒。

2. **线程 A 重入锁：**
   - 在业务逻辑中，`A` 再次调用 `tryLock()`，此时锁已经存在且属于 `A`，`hexists` 判断 `myLock` 中的 `id:threadIdA` 存在。于是 `hincrby` 将计数器增加到 `2`，同时重置锁的过期时间。

3. **线程 B 尝试获取锁：**
   - 在 `A` 持有锁期间，`B` 调用 `tryLock()`，此时 `myLock` 存在且 `id:threadIdA` 存在于 `hash` 中。因此，`B` 未能获取锁，Lua 脚本返回 `myLock` 的剩余生存时间，表示锁当前被 `A` 持有。

4. **线程 A 释放锁：**
   - 当 `A` 结束业务逻辑时，`A` 调用 `unlock()`，锁的计数器减 `1`，直到计数器归零，锁完全释放。此时，`B` 可以继续尝试获取锁。

### 总结
Redisson 的可重入锁借助 Redis 的 `hash` 结构和 Lua 脚本实现了分布式环境下的锁重入机制。通过 `id:threadId` 来区分不同线程持有同一锁，避免了锁的误释放问题，同时保证了锁的原子性和可重入性。
```

### 为什么使用看门狗机制？

在分布式锁中，使用看门狗（WatchDog）机制而不是设置锁永不过期，主要是为了提高系统的可靠性和避免一些潜在的风险：

#### 1. **防止锁被长时间占用**

如果将锁的过期时间设置为无限长，可能会导致以下问题：

- **线程宕机或崩溃**：如果持有锁的线程因故障或崩溃而无法释放锁，锁将永远被占用，其他线程将无法获取锁。这会导致系统停滞不前，无法正常执行任务。
- **资源浪费**：即使锁的持有者已不再需要该锁，锁依然存在，可能会占用 Redis 的存储资源，影响系统性能。

看门狗机制的核心思想是确保锁的过期时间足够长，并在锁被持有期间定期续约。这样即使持有锁的线程发生故障，锁也会在一定时间后自动释放，从而避免锁的长期占用问题。

#### 2. **防止资源泄漏**

看门狗机制通过定期续约锁的过期时间，确保在持有锁的过程中，锁不会因过期而被意外释放。这样可以有效防止锁的资源泄漏问题。

- **自动续约**：看门狗机制会周期性地续约锁的过期时间，直到持有锁的线程显式地释放锁。这样可以确保锁在使用过程中不会因为超时而被释放。
- **失效处理**：如果持有锁的线程发生故障，看门狗机制会在锁的过期时间后自动释放锁，避免锁被长时间占用。

#### 3. **提高系统的鲁棒性**

看门狗机制可以提高系统的鲁棒性，确保即使在网络延迟或线程故障的情况下，系统依然能够正常工作。

- **网络延迟**：在高延迟或不稳定的网络环境中，看门狗机制可以确保锁的有效性，不会因为网络波动导致锁的过期。
- **线程崩溃**：如果线程崩溃或被终止，看门狗机制会确保锁在一段时间后被自动释放，避免对系统的长时间锁定。

#### 4. **灵活性和配置**

看门狗机制提供了灵活的配置选项，使得系统可以根据具体需求调整锁的过期时间和续约策略。

- **自定义过期时间**：可以根据业务需求设置不同的锁过期时间。
- **动态调整**：看门狗机制允许动态调整锁的过期时间和续约策略，以适应不同的业务场景。



```
`subscribeFuture` 是一个 `RFuture<RedissonLockEntry>` 类型的对象。它表示一个异步操作的结果，用于在 Redisson 分布式锁的上下文中处理锁释放事件的订阅。

### 详细解释

#### 1. **类型定义**

- **`RFuture<T>`** 是 Redisson 提供的一个接口，表示异步操作的结果。在 Redisson 中，`RFuture<T>` 用于处理与 Redis 交互的异步操作，`T` 表示操作的结果类型。

- **`RedissonLockEntry`** 是一个特定于 Redisson 锁的类，它可能包含有关锁的详细信息，例如锁的持有者或锁的状态。

#### 2. **功能**

- `subscribeFuture` 是通过 `this.subscribe()` 方法获得的，该方法用于订阅锁的释放事件。当锁被释放时，该订阅将会通知 `subscribeFuture`。

- `subscribeFuture.onComplete((r, ex) -> { ... })` 是一个回调函数，当 `subscribeFuture` 完成时，这个回调函数会被调用。`r` 是订阅操作的结果（即 `RedissonLockEntry`），而 `ex` 是可能发生的异常。

  - **`ex != null`**：如果订阅操作失败（即 `ex` 不为空），则通过 `result.tryFailure(ex)` 将结果标记为失败，并传递异常。

  - **`if (futureRef.get() != null)`**：如果有一个定时任务正在等待订阅完成，则取消该定时任务。

  - **`time.addAndGet(-elapsed)`**：计算并更新剩余的等待时间。

  - **`this.tryAcquireAsync(time, permits, subscribeFuture, result, ttl, timeUnit)`**：在订阅成功的情况下，重新尝试获取锁（递归调用 `tryAcquireAsync` 方法）。

#### 3. **订阅机制**

`subscribeFuture` 的主要作用是处理锁释放的通知。在获取锁失败的情况下，通过订阅机制可以在锁释放时尝试重新获取锁。这种机制确保了在锁被释放时，系统能够及时响应并尝试重新获取锁，从而提高了锁的获取成功率。

### 总结

`subscribeFuture` 是一个 `RFuture<RedissonLockEntry>` 类型的对象，用于表示锁释放事件的异步订阅操作。它在获取锁失败时通过订阅机制监听锁释放事件，并在事件发生时尝试重新获取锁。
```



### 所有兜底都可以是较短的缓存过期时间







![无标题-2024-09-07-1323](/Users/haozhipeng/Downloads/我的笔记/images/无标题-2024-09-07-1323.png)





**Redis** 是一个基于内存的高性能键值数据库，广泛用于缓存、会话管理和实时数据存储等场景。Redis 内部使用了多种数据结构来存储数据，但所有这些数据最终都由一个全局哈希表来管理。Redis 的核心数据存储机制依赖于一个 **全局哈希表**，它负责存储和索引所有键值对数据。

### 1. 全局哈希表的概念

Redis 的全局哈希表是用于存储所有键值对的核心数据结构。它是一个由键值对（key-value pairs）组成的大型哈希表，Redis 实例中的所有键值对都存储在这个哈希表中。Redis 使用哈希表的设计是为了提供 O(1) 时间复杂度的查找、插入和删除操作。

### 2. Redis 全局哈希表的结构

Redis 的全局哈希表实际上是由两个哈希表组成的，分别是 `ht[0]` 和 `ht[1]`，它们实现了动态扩展和渐进式 rehashing（再哈希）的机制：

- **ht[0]**：主哈希表，用于存储当前所有的键值对。
- **ht[1]**：备用哈希表，仅在扩容或缩容时使用，用来迁移数据。当扩容或缩容完成后，`ht[1]` 被清空或丢弃。

#### 哈希表内部结构

哈希表本质上是一个包含许多 **槽位（bucket）** 的数组，每个槽位存储一个链表，链表中每个节点包含一个键值对。Redis 的哈希表通过一个哈希函数将键映射到数组中的某个槽位中：

```c
typedef struct dictEntry {
    void *key;
    void *value;
    struct dictEntry *next;  // 用于处理哈希冲突
} dictEntry;
```

每个 `dictEntry` 结构包含三个部分：
- **key**：表示 Redis 中的键。
- **value**：表示与该键相关联的值（可以是简单值或复杂数据类型）。
- **next**：用于在哈希冲突时链接到链表中的下一个节点。

### 3. 哈希表的扩展与缩容

Redis 中的哈希表是动态扩展和缩容的，目的是保持哈希表的负载因子在合理范围内。负载因子定义为：
```
负载因子 = 哈希表中存储的键值对数量 / 哈希表的槽位数量
```
Redis 通过调整哈希表的大小，确保负载因子不会过高，以避免哈希冲突过多而影响性能。

#### 3.1 扩容

当 Redis 中存储的键值对数量超过了当前哈希表的容量时，Redis 会触发 **扩容** 操作。扩容的目的是增加哈希表的槽位数量，从而减少哈希冲突。扩容的过程如下：

1. Redis 创建一个新的、更大的哈希表 `ht[1]`，容量通常是原哈希表的两倍。
2. Redis 将旧哈希表 `ht[0]` 中的键值对逐步迁移到新哈希表 `ht[1]` 中，这个过程叫做 **渐进式 rehashing**。
3. 每次对 Redis 进行插入、删除、查询等操作时，Redis 会将一部分键值对从 `ht[0]` 迁移到 `ht[1]` 中。
4. 当所有键值对都迁移完毕后，Redis 丢弃旧哈希表，将 `ht[1]` 设置为新的 `ht[0]`。

这种 **渐进式 rehashing** 可以避免在扩容过程中因一次性迁移大量数据而引发性能瓶颈。每次操作时只迁移一部分数据，分散了迁移开销。

#### 3.2 缩容

当 Redis 中存储的键值对数量减少时（例如大量删除操作），Redis 会触发 **缩容** 操作，以释放内存。缩容的过程与扩容类似，Redis 会创建一个较小的哈希表，并逐步将旧哈希表中的键值对迁移到新哈希表中。

缩容通常在负载因子过低时触发，即哈希表的槽位数量远超实际存储的键值对数量。

### 4. 渐进式 Rehashing 机制

Redis 的 **渐进式 rehashing** 是哈希表动态扩容和缩容的关键机制。当哈希表需要扩容或缩容时，Redis 不会一次性将所有数据迁移到新的哈希表中，而是采用渐进式的方式，分批次逐步完成数据迁移。

具体过程如下：
1. Redis 在执行插入、查询或删除操作时，会同时处理部分数据迁移。每次数据操作时，都会从旧的哈希表 `ht[0]` 中取出一部分键值对，迁移到新的哈希表 `ht[1]`。
2. Redis 使用 `rehashidx` 变量跟踪迁移进度。每次数据迁移时，`rehashidx` 会递增，表示下一个要迁移的槽位。
3. 当所有键值对迁移完毕后，Redis 将 `ht[1]` 作为新的主哈希表 `ht[0]`，旧哈希表 `ht[0]` 被丢弃或释放内存。

这种渐进式迁移机制避免了 Redis 在扩容或缩容时对性能的影响。数据操作和迁移任务是并行进行的，每次操作只需处理少量数据，保证了 Redis 的高并发性能。

### 5. 哈希冲突及其处理

在哈希表中，**哈希冲突** 是指两个不同的键通过哈希函数计算后，映射到了相同的槽位。这种情况下，Redis 采用 **链地址法** 解决冲突问题。

具体来说，哈希表的每个槽位都是一个链表。当多个键映射到同一个槽位时，Redis 会将它们链接成一个链表，每个链表节点保存一个键值对。Redis 在执行查询、插入或删除操作时，首先通过哈希值找到对应的槽位，然后在链表中遍历键值对，找到所需的键。

### 6. Redis 全局哈希表的多种数据结构

虽然 Redis 中所有数据都存储在全局哈希表中，但 Redis 支持多种不同的数据类型，每个数据类型都通过键值对存储在哈希表中。以下是 Redis 支持的主要数据类型及其存储方式：

- **字符串（String）**：Redis 最基本的数据类型，每个字符串键都映射到一个字符串值。
- **列表（List）**：列表是一个双向链表，支持从头部和尾部插入和删除元素。每个列表通过一个键映射到一个双向链表。
- **集合（Set）**：集合是一个无序的、唯一的元素集合。Redis 使用哈希表实现集合，每个集合通过一个键映射到一个哈希表。
- **有序集合（Sorted Set）**：有序集合类似于集合，但每个元素都有一个关联的分数。Redis 使用跳表（Skip List）和哈希表实现有序集合。
- **哈希表（Hash）**：哈希表用于存储键值对集合，每个 Redis 哈希通过一个键映射到一个内部的哈希表。

### 7. Redis 的高效性与可扩展性

Redis 全局哈希表的设计使得 Redis 能够在高并发场景下提供极高的性能。其 O(1) 时间复杂度的操作加上动态扩容、渐进式 rehashing 和链地址法的冲突处理，使 Redis 能够应对大规模的键值对存储和查询需求。

此外，Redis 还通过数据持久化、主从复制、集群模式等机制，确保数据的高可用性和可扩展性。

### 8. 总结

Redis 的全局哈希表是其数据存储的核心结构，负责管理所有键值对的数据。通过双哈希表和渐进式 rehashing，Redis 能够在高效地管理哈希表大小的同时，避免性能瓶颈。链地址法用于解决哈希冲突，而全局哈希表则管理各种数据类型（如字符串、列表、集合、哈希表等），使 Redis 能够高效、灵活地处理各种存储需求。

Redis 的这一全局哈希表机制，是其高性能和低延迟的基础，同时支持各种复杂的数据操作场景。



<img src="https://files.xyhelper.cn/file-5scloDsMIwJOh46Gwwsds5B8?se=2024-09-23T14%3A46%3A10Z&sp=r&sv=2024-08-04&sr=b&rscc=max-age%3D299%2C%20immutable%2C%20private&rscd=attachment%3B%20filename%3Dimage.png&sig=E82skSygllBbkeGPhiMyfNTgEsz7/cM5gGLbviDILjU%3D" alt="已上传的图片" style="zoom:50%;" />

<img src="/Users/moon/Library/Application Support/typora-user-images/image-20240923224314305.png" alt="image-20240923224314305" style="zoom:50%;" />