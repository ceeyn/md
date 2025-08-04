### 哈希分区

#### 节点取余分区

对特定数据采用hash算法得到一个整数，再通过整数对分区数取余就可以得到资源的存储路由。如redis的键或用户ID，再根据节点数量N使用公式：hash(key)%N计算出哈希值，用来决定数据映射到哪个分区节点。

1. 优点
   这种方式的突出优点就是简单，且常用于数据库的分库分表。如京东某系统中采用shardingjdbc，用这种方式进行分库分表路由处理。
2. 缺点
   当节点数量发生变化时，如扩容或收缩节点（没有遇到过），数据节点关系需要重新计算，会导致数据的重新迁移。所以扩容通常采用翻倍扩容，避免数据映射全部打乱而全部迁移，翻倍迁移只发生50%的数据迁移。如果不翻倍缩扩容，如某一台机器宕机，那么应该落在该机器的请求就无法得到正确的处理，这时需要将宕掉的服务器使用算法去除，此时候会有(N-1)/N的服务器的缓存数据需要重新进行计算；如果新增一台机器，会有N /(N+1)的服务器的缓存数据需要进行重新计算。对于系统而言，这通常是不可接受的颠簸（因为这意味着大量缓存的失效或者数据需要转移）。

## 一致性hash算法

一致性hash算法是因为节点数目发生改变时，尽可能少的数据迁移而出现的。比如扩容时，需要50%的数据迁移；但如果引入一种算法，可以减少数据的迁移量，所以就出现了一致性hash算法。将所有的存储节点排列在收尾相接的hash环上，每个key在计算Hash后会顺时针找到临接的存储节点存放。而当有节点加入或退出时，仅影响该节点在Hash环上顺时针相邻的后续节点。

1. 优点
   加入和删除节点只影响哈希环中顺时针方向的相邻的节点，对其他节点无影响。
2. 缺点
   数据的分布和节点的位置有关，因为这些节点不是均匀的分布在哈希环上的，所以数据在进行存储时达不到均匀分布的效果。所以，出现了**增加虚拟节点的方式来减少不均衡的现象。**









### ConsistentHashLoadBalance

一致性 hash 算法由麻省理工学院的 Karger 及其合作者于1997年提出的，算法提出之初是用于大规模缓存系统的负载均衡。它的工作过程是这样的，首先根据 ip 或者其他的信息为缓存节点生成一个 hash，并将这个 hash 投射到 [0, 232 - 1] 的圆环上。当有查询或写入请求时，则为缓存项的 key 生成一个 hash 值。然后查找第一个大于或等于该 hash 值的缓存节点，并到这个节点中查询或写入缓存项。如果当前节点挂了，则在下一次查询或写入缓存时，为缓存项查找另一个大于其 hash 值的缓存节点即可。大致效果如下图所示，每个缓存节点在圆环上占据一个位置。如果缓存项的 key 的 hash 值小于缓存节点 hash 值，则到该缓存节点中存储或读取缓存项。比如下面绿色点对应的缓存项将会被存储到 cache-2 节点中。**由于 cache-3 挂了，原本应该存到该节点中的缓存项最终会存储到 cache-4 节点中。**

![img](https://cn.dubbo.apache.org/imgs/dev/consistent-hash.jpg)

下面来看看一致性 hash 在 Dubbo 中的应用。我们把上图的缓存节点替换成 Dubbo 的服务提供者，于是得到了下图：

![img](https://cn.dubbo.apache.org/imgs/dev/consistent-hash-invoker.jpg)

这里相同颜色的节点均属于同一个服务提供者，比如 Invoker1-1，Invoker1-2，……, Invoker1-160。这样做的目的是通过引入虚拟节点，让 Invoker 在圆环上分散开来，避免数据倾斜问题。所谓数据倾斜是指，由于节点不够分散，导致大量请求落到了同一个节点上，而其他节点只会接收到了少量请求的情况。比如：

![img](https://cn.dubbo.apache.org/imgs/dev/consistent-hash-data-incline.jpg)

如上，由于 Invoker-1 和 Invoker-2 在圆环上分布不均，导致系统中75%的请求都会落到 Invoker-1 上，只有 25% 的请求会落到 Invoker-2 上。解决这个问题办法是引入虚拟节点，通过虚拟节点均衡各个节点的请求量

抱歉之前的解释没有完全展示如何生成 160 个虚拟节点。现在我会详细列出如何从一个 IP 地址生成 160 个虚拟节点，并解释每个虚拟节点的哈希值是如何计算的。

### 步骤 1：生成虚拟节点

首先，回顾一下代码中的生成虚拟节点的逻辑：

```java
public ConsistentHashSelector(List<ServiceInfo> invokers, int replicaNumber, int identityHashCode) {
    this.virtualInvokers = new TreeMap<>();
    this.identityHashCode = identityHashCode;

    for (ServiceInfo invoker : invokers) {
        String address = invoker.getAddress();
        for (int i = 0; i < replicaNumber / 4; i++) {
            byte[] digest = md5(address + i); // 对 address + i 进行 MD5 运算，得到一个长度为 16 字节的数组
            for (int h = 0; h < 4; h++) {    // 每 4 字节生成一个虚拟节点
                long m = hash(digest, h);     // 计算哈希值
                virtualInvokers.put(m, invoker); // 将哈希值与节点地址绑定
            }
        }
    }
}
```

### 具体操作步骤：

#### 假设：
- 节点 A 的 IP 地址是 `"192.168.0.1"`
- `replicaNumber = 160`

每个物理节点会生成 160 个虚拟节点。生成过程如下：

1. **MD5 哈希**：  
   首先，将 IP 地址和递增的整数 `i` 进行字符串拼接，然后通过 `MD5` 算法生成哈希值。`MD5` 生成一个 16 字节的字节数组。

2. **生成虚拟节点**：  
   每个 16 字节的 `digest` 被分成 4 部分，每部分 4 字节，共生成 4 个虚拟节点。重复这一过程 `160 / 4 = 40` 次，以生成 160 个虚拟节点。

### 示例：生成第一个虚拟节点组

#### 假设：
- `address = "192.168.0.1"`
- `i = 0`

拼接后字符串为 `"192.168.0.10"`。对其进行 MD5 运算。

1. **MD5 结果**（假设结果为）：

```java
digest = md5("192.168.0.10");
```

假设 `digest` 得到的 16 字节数组如下：

```
digest = [0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88]
```

2. **生成虚拟节点哈希值**：

我们将 `digest` 分成 4 个部分，每部分 4 字节：

- **第一部分**（`h = 0`）：

```java
long m0 = (((long) (digest[3] & 0xFF) << 24)
          | ((long) (digest[2] & 0xFF) << 16)
          | ((long) (digest[1] & 0xFF) << 8)
          | (digest[0] & 0xFF)) & 0xFFFFFFFFL;
```

计算结果：

```
m0 = (0x78 << 24) | (0x56 << 16) | (0x34 << 8) | 0x12 = 0x78563412
```

- **第二部分**（`h = 1`）：

```java
long m1 = (((long) (digest[7] & 0xFF) << 24)
          | ((long) (digest[6] & 0xFF) << 16)
          | ((long) (digest[5] & 0xFF) << 8)
          | (digest[4] & 0xFF)) & 0xFFFFFFFFL;
```

计算结果：

```
m1 = (0xF0 << 24) | (0xDE << 16) | (0xBC << 8) | 0x9A = 0xF0DEBC9A
```

- **第三部分**（`h = 2`）：

```java
long m2 = (((long) (digest[11] & 0xFF) << 24)
          | ((long) (digest[10] & 0xFF) << 16)
          | ((long) (digest[9] & 0xFF) << 8)
          | (digest[8] & 0xFF)) & 0xFFFFFFFFL;
```

计算结果：

```
m2 = (0x44 << 24) | (0x33 << 16) | (0x22 << 8) | 0x11 = 0x44332211
```

- **第四部分**（`h = 3`）：

```java
long m3 = (((long) (digest[15] & 0xFF) << 24)
          | ((long) (digest[14] & 0xFF) << 16)
          | ((long) (digest[13] & 0xFF) << 8)
          | (digest[12] & 0xFF)) & 0xFFFFFFFFL;
```

计算结果：

```
m3 = (0x88 << 24) | (0x77 << 16) | (0x66 << 8) | 0x55 = 0x88776655
```

这四个哈希值分别是 `m0, m1, m2, m3`，对应四个虚拟节点的位置。

### 重复上述步骤：

通过增加 `i` 的值（从 0 到 39），并且每次对 `address + i` 进行 MD5 哈希运算，再从中生成 4 个虚拟节点的哈希值。这样，我们最终为 `"192.168.0.1"` 生成了 160 个虚拟节点的哈希值，存储在 `TreeMap` 中。

### 步骤 3：查询哈希环上的节点

当我们对一个请求进行哈希运算并得到一个哈希值时，例如得到哈希值 `0x44332211`，我们将会查找哈希环上第一个大于等于该哈希值的虚拟节点：

```java
long hash = 0x44332211L;
Map.Entry<Long, ServiceInfo> entry = virtualInvokers.ceilingEntry(hash);
```

如果 `ceilingEntry` 找到的节点属于 `A`，则该请求将被路由到物理节点 `A`。

### 总结

在这个例子中，我们详细展示了如何从 IP 地址生成 160 个虚拟节点，并解释了每个虚拟节点的哈希值是如何计算的。通过将这些虚拟节点均匀分布在哈希环上，我们可以有效地实现负载均衡，并减少数据迁移的负担。







### 1. 背景

我们有三个物理节点 `A`、`B` 和 `C`，它们分别通过一致性哈希算法映射到 160 个虚拟节点，并均匀分布在哈希环上。

我们生成了 160 个虚拟节点，对于每个物理节点，具体步骤前面已经详细解释过。现在我们要看一下一个请求（`RpcRequest`）是如何通过哈希环找到对应的物理节点。

### 2. 请求的映射过程

假设我们有一个 `RpcRequest`，方法是 `getUser`，参数是 `"user123"`，我们需要将该请求映射到一个物理节点上。这个过程如下：

#### 2.1 生成请求的唯一键

首先，我们生成请求的唯一键（`selectKey`），它包含了服务的全限定类名、方法名以及参数信息：

```
java
复制代码
RpcRequest request = new RpcRequest("com.example.Service", "getUser", new Object[]{"user123"});
String selectKey = request.getServiceName() + "." + request.getMethod() + Arrays.toString(request.getParameterValues());
```

生成的 `selectKey` 可能是：

```
arduino
复制代码
"com.example.Service.getUser[user123]"
```

#### 2.2 对 `selectKey` 进行哈希计算

接下来，我们对 `selectKey` 进行 MD5 哈希计算，得到一个哈希值：

```
java
复制代码
byte[] digest = md5(selectKey);
long hash = hash(digest, 0);
```

假设计算得到的哈希值是 `0x44332211`（十六进制表示）。

#### 2.3 在哈希环上找到对应的虚拟节点

我们现在需要在哈希环上找到第一个大于或等于这个哈希值的虚拟节点：

```
java
复制代码
Map.Entry<Long, ServiceInfo> entry = virtualInvokers.ceilingEntry(hash);
```

- `ceilingEntry(hash)` 方法会在 `TreeMap` 中找到第一个键值（虚拟节点的哈希值）大于或等于 `0x44332211` 的条目。

假设 `ceilingEntry` 找到的虚拟节点属于物理节点 `B`。

#### 2.4 获取物理节点信息

通过 `entry.getValue()`，我们可以获取这个虚拟节点所对应的物理节点信息：

```
java
复制代码
ServiceInfo selectedNode = entry.getValue();
```

如果找到的虚拟节点属于 `B`，那么这个请求就会被路由到物理节点 `B`。

#### 2.5 如果哈希值超出环的最大值

如果生成的哈希值比哈希环上所有虚拟节点的哈希值都要大，`ceilingEntry(hash)` 将返回 `null`，这意味着我们需要回到环的起点。这种情况下，我们选择 `TreeMap` 中的第一个条目：

```
java
复制代码
if (entry == null) {
    entry = virtualInvokers.firstEntry();
}
ServiceInfo selectedNode = entry.getValue();
```

这样，即使哈希值超出环的最大值，我们也能找到最接近的节点。

### 3. 总结请求映射过程

- **步骤 1**: 生成请求的唯一键 `selectKey`。
- **步骤 2**: 对 `selectKey` 进行 MD5 哈希计算，得到哈希值 `hash`。
- **步骤 3**: 在哈希环上找到第一个大于或等于 `hash` 的虚拟节点。
- **步骤 4**: 获取该虚拟节点对应的物理节点，并将请求路由到该节点。
- **步骤 5**: 如果没有找到合适的虚拟节点，选择哈希环上的第一个虚拟节点。

通过这种方式，请求被映射到哈希环上的某个虚拟节点，而虚拟节点最终对应的是一个物理节点。这种机制确保了即使节点发生变化，数据迁移量也最小，并且请求始终可以路由到正确的节点上。



