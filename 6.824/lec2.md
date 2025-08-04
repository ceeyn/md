

性能 -》分片 -》 容错（副本） -〉不一致问题

强一致性的一个体现：对于多台副本，执行请求队列顺序必须一致

master 存储2张表：1.文件到chunk的映射（持久化） 2.chunk到ip列表的映射（内存）

1：根据文件偏移对64mb取余就知道在哪个chunk

2:写文件的时候，要找到最后一个chunck，并且找到最新版本（master会存到log里），提交的时候，Primary和所有follower都写入成功才是成功，假如一个失败则会重发，直到写入成功。版本号是每个Primary一个版本

客户端向Primary和follower发数据是流式的，先发送到最近一个

Primary 通过租期避免脑裂，当master联系不上Primary时，会拒绝所有写请求，直到租期过期，然后重新制定Primary，这样可以避免脑裂，租期内两个Priamry都提供服务

<img src="/Users/moon/Library/Application Support/typora-user-images/image-20250411001456550.png" alt="image-20250411001456550" style="zoom:25%;" />

一个失败则会重发，直到写入成功，（部分写入失败，也会有部分成功，成功的追加在末尾，这样会导致乱序）



<img src="https://files.oaiusercontent.com/file-SETVv8iszKA9ZiD4f1pMWg?se=2025-04-11T13%3A11%3A42Z&sp=r&sv=2024-08-04&sr=b&rscc=max-age%3D299%2C%20immutable%2C%20private&rscd=attachment%3B%20filename%3Dimage.png&sig=USSDAD9it8dQb1t7/W56ARrbFeY8IMC6ZA9qYBNQbqs%3D" alt="已上传的图片" style="zoom:50%;" />



选择大块大小的好处包括：

1. **减少客户端与主节点的交互**。在同一个块内进行读写操作，仅需首次请求主节点获取位置信息；
   - 对于我们的典型负载（如顺序读写大文件）来说，减少效果尤为显著；
   - 即使是小规模随机读，客户端也能轻松缓存多个 TB 数据集的块位置。
2. **减少网络开销**。客户端在处理大块时往往会对同一个块进行多次操作，因此可以维持与 chunkserver 的**持久 TCP 连接**，节省重复建立连接的开销。
3. **减小主节点上的元数据体积**。由于每个文件需要的块数减少，主节点维护的块映射数量也减少，**元数据可以完全驻留在内存中**，从而带来更多性能优势（详见 2.6.1 节）。



### **3.2 数据流与控制流分离（Data Flow）**

GFS 将**数据流**和**控制流**解耦，提高网络利用率：

- 控制流：客户端 → 主副本 → 副本；
- 数据流：客户端 → chunkserver 链式传递（pipeline）；
- 这样可最大化带宽利用，减少高延迟链路冲突。

例如：客户端向副本 S1~S4 发送数据，会选择网络距离最近的顺序发送：

```
复制编辑
客户端 → S1 → S2 → S3 → S4
```

链式 TCP pipeline 可以实现较低延迟，例如 1MB 数据在 100Mbps 网络中可在约 80ms 内完成复制。



/data/bigfile       ← 正常文件
/snapshot/bigfile   ← 快照文件（引用相同的块）
用户对 /data/bigfile 发起快照，生成 /snapshot/bigfile

快照操作立即完成（因为两者引用相同块，未真正复制）

接下来用户向 /data/bigfile 写入数据块 42

主服务器检测到：块 42 被多个文件共享（引用数 > 1）

于是主服务器创建新块 42′（COW），将写操作落在 42′ 上

/data/bigfile 使用 42′，/snapshot/bigfile 仍然使用原始块 42





gfs

> GFS 是为 **Google 内部 append-heavy 的大规模数据处理场景** 优化设计的文件系统，拥有高吞吐、弱一致、高容错、快速恢复的特点。其通过 **Copy-on-Write 快照、Record Append、Shadow Master 和冗余副本** 等机制实现了可扩展性与高可用性。

master：拥有元信息（主机副本位置，块位置：master上线副本主动推送给master），

client读/写：问master拿到leader，follower位置信息，读任意都可以读，写数据先发，控制信号后发，发数据时是从最近的机器发，流式发，边写边往其他副本发，client通知leader写入，如果是append会维护队列，保证并发安全。普通写则会出现问题

snapshot：master先收回lease（写的权限），然后复制一些指针的元数据信息，然后等到有写请求来时发现chunk有多个指针，则进行真正数据的复制生成snapshot，写请求写入原文件。

showdmater：开始时从master获取log文件进行同步

checksum：chunck64mb，被分成若干块64kb，每块一个checksum 32b，每次读前先计算checksum比对没问题后才读，append写的话只用修改最后一块的checksum



IO concurrency io 并发

我有一个进程 启动了多个线程，可以执行在网络上不同服务器远程调用，一个写磁盘，一个写网卡等等

事件编程 一个线程 一个while

维护一张不同客户端对应的进度表，每次收到请求，查看是哪个客户端，执行相应的进度

将代码切割成一个个小块（事件），由while决定激活相应的事件



a进程1，2，3； b进程1，2

假如单核cpu，操作系统调度时决定从a1，切换到b1，这样进行调度，而routine则是在线程内运行的，对操作系统透明的。



内部锁的一个问题是，两个内部锁对象可能互相死锁

```go
	var done sync.WaitGroup
	for _, u := range urls {
		done.Add(1)
		u2 := u
		go func() {
			defer done.Done()
			ConcurrentMutex(u2, fetcher, f)
		}()
		//go func(u string) {
		//	defer done.Done()
		//	ConcurrentMutex(u, fetcher, f)
		//}(u)
	}
	1.done.Wait() 对于不同线程done都是一个，而done没有任何同步措施，这也就说明了done内部必定有锁。
	2.u 并不是每次都新创建，而是同一个变量被反复赋值。所以 所有 goroutine 访问的其实是“同一个 u”，
```

> **函数参数传递是值传递**，但对于引用类型（如 map、slice、指针、chan 等），传递的是“引用的值”。

闭包（closure）的通用行为，**捕获变量引用**而不是其值。

语言设计者这样做，是为了节省内存和支持动态访问变量 —— 但也带来了副作用：**共享变量在异步场景下可能出错**。



##### goroutine 延迟执行，等开始执行时，u 已经变了

- goroutine 是异步的，可能在 `for` 循环执行完之后才开始执行。
- 到那时，`u` 的值已经变成了循环的最后一个元素（比如 `"c"`）。

所以你看到的现象就是：

> 所有 goroutine 访问到的 `u` 值都是最后一次循环赋的值。



##### 逃逸

逃逸 = **变量在函数作用域外还会被引用**，就不能放在栈上了，只能放在堆上。

func foo() {
  x := 100   ← 栈上

  return &x  ← x 要逃逸到堆上
}

### 📦 举个简单例子：

```
go


复制编辑
func makeAdder(x int) func(int) int {
    return func(y int) int {
        return x + y
    }
}
```

🔍 上面这段代码中，变量 `x` 是 `makeAdder` 函数的局部变量，但是它被返回的闭包引用了。

> 即使 `makeAdder` 已经执行完了，返回的闭包还要用到 `x`！

➡️ **所以 `x` 必须分配在堆上**，确保闭包外部还可以访问它。



## 插件

插件（`.so`）是用 Go 写的可执行模块，主程序通过 `plugin.Open()` 在运行时动态加载，并通过 `Lookup()` 找到函数/变量，从而实现模块热插拔、功能解耦。



### 1️⃣ 插件：`wc.go`

```
go


复制编辑
package main

import "strings"

func Map(filename string, contents string) []KeyValue {
    var kva []KeyValue
    words := strings.Fields(contents)
    for _, w := range words {
        kva = append(kva, KeyValue{w, "1"})
    }
    return kva
}

func Reduce(key string, values []string) string {
    return strconv.Itoa(len(values))
}
```

> 然后我们执行命令生成插件：

```
bash


复制编辑
go build -buildmode=plugin -o wc.so wc.go
```

------

### 2️⃣ 主程序：动态加载插件

```
go


复制编辑
package main

import (
    "plugin"
    "fmt"
)

func main() {
    p, err := plugin.Open("wc.so") // 动态加载插件
    if err != nil {
        panic(err)
    }

    // 查找函数
    mapFunc, err := p.Lookup("Map")
    if err != nil {
        panic(err)
    }

    // 将函数转换成实际类型然后调用
    result := mapFunc.(func(string, string) []KeyValue)("file.txt", "hello world hello")

    fmt.Println(result)
}
```