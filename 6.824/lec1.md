https://mit-public-courses-cn-translatio.gitbook.io/mit6-824/lecture-07-raft2/7.3-hui-fu-jia-su-backup-acceleration



Avilable

hard consistency

performance







MapDuce  分发者

worker：

​	计算：cal(k, v) : case k -> for (执行v的计算) 输入是文件【GFS的机器上运行worker线程，避免网络开销】,输出也是文件

worker：

​	汇总：sum(k,v) 从各台服务器上汇总完结果，发送到GFS中 网络开销比较严重

GFS 将大文件拆成64kb进行存储



map、reduce任务：失败重跑

map结果写入本地磁盘

shufile：将map partition结果发送给reduce，



reduce结果写入中间文件，中间文件写入完写入到gfs中，整个过程原子性的，不支持两阶段提交（除了最终输出文件外，最多写入一个其他文件，但不保证整个文件最后是正确的）

其他机器（即其他 Map Worker）上处理出来的 `"baidu.com"` 和 `"amazon.com"`，**只要它们的 key 经过相同的 partition 函数（默认是 `hash(key) % R`）得出的结果为 0，它们就会被分到 Reduce 分区 0，由同一个 Reduce Worker 来处理。**



------

## ✅ 举个完整例子

假设：

- 我们有 3 台 Map Worker：`MapWorker0`, `MapWorker1`, `MapWorker2`
- 我们有 2 个 Reduce Worker（分区 0 和 1）
- partition 函数是：`hash(key) % 2`

**MapWorker0 的输出：**

```
"baidu.com"   → hash("baidu.com") % 2 = 0 → 放入本地 mapWorker0-part-0
"google.com"  → hash("google.com") % 2 = 1 → 放入 mapWorker0-part-1
```

**MapWorker1 的输出：**

```
"amazon.com"  → hash("amazon.com") % 2 = 0 → 放入 mapWorker1-part-0
"baidu.com"   → hash("baidu.com") % 2 = 0 → 放入 mapWorker1-part-0
```

**MapWorker2 的输出：**

```
"google.com"  → hash("google.com") % 2 = 1 → 放入 mapWorker2-part-1
"baidu.com"   → hash("baidu.com") % 2 = 0 → 放入 mapWorker2-part-0
```

------

## 🧠 Reduce 处理阶段：

Reduce Worker 会根据分区编号分别：

- **Reduce 0**：
  - 拉取所有 map worker 中的 `part-0`：
    - `mapWorker0-part-0`
    - `mapWorker1-part-0`
    - `mapWorker2-part-0`
  - 合并所有 `"baidu.com"`、`"amazon.com"` 的 key/value，聚合统计等。
- **Reduce 1**：
  - 拉取所有 `part-1` 的内容（比如 `"google.com"`）

------

## ✅ 所以结论：

> 不管 `"baidu.com"` 是在哪个 Map Worker 被处理出来的，**只要 hash 值对应的是 Reduce 分区 0，它最终就会归属于 Reduce Worker 0**，被统一归并处理。

------

## 🔁 这就是 MapReduce 中的「**Shuffle 阶段**」核心：

- Map Worker 本地将输出分 R 份（每份给一个 Reduce 分区）
- Reduce Worker 根据 key 分区编号，**从所有 Map Worker 拉取属于自己的那一份数据**

------

如果你还想了解：

- Shuffle 过程中如何排序
- Reduce 如何并行处理大量 key
- 如何防止 Reduce 数据倾斜（如一个 key 非常大）

我也可以帮你讲解 😎

