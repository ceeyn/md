

### 锁是为了保护不变量的属性

例如，不是为了保护bob +和 alice -的原子性，而是为了维持 bob+alice 的总数不变的这个不变量。当 lock 和 unLock 的时候实际上破坏了这个不变量。

<img src="/Users/moon/Library/Application Support/typora-user-images/image-20250419214451331.png" alt="image-20250419214451331" style="zoom:33%;" />



<img src="/Users/moon/Library/Application Support/typora-user-images/image-20250419221617611.png" alt="image-20250419221617611" style="zoom: 33%;" />

Chan 死锁

<img src="/Users/moon/Library/Application Support/typora-user-images/image-20250419232655926.png" alt="image-20250419232655926" style="zoom:33%;" />



Chan 实现 waitGroup

<img src="/Users/moon/Library/Application Support/typora-user-images/image-20250419232531083.png" alt="image-20250419232531083" style="zoom:50%;" />

waitGroup

<img src="/Users/moon/Library/Application Support/typora-user-images/image-20250419233523855.png" alt="image-20250419233523855" style="zoom:33%;" />

#### 死锁的情况

<img src="/Users/moon/Library/Application Support/typora-user-images/image-20250420092812508.png" alt="image-20250420092812508" style="zoom:50%;" />



<img src="/Users/moon/Library/Application Support/typora-user-images/image-20250420092926758.png" alt="image-20250420092926758" style="zoom:50%;" />

```
S0 lock0 send s1

S1 lock1 send s0

S0 need lock0 to reply

S1 need lock1 to reply

死锁了。。。。。
```



 rpc调用的时候【等待其他人回复的时候】，不应该持有锁

准备 rpc 的发送参数的时候锁定，发送之前解锁

<img src="/Users/moon/Library/Application Support/typora-user-images/image-20250420092415425.png" alt="image-20250420092415425" style="zoom:50%;" />





<img src="/Users/moon/Library/Application Support/typora-user-images/image-20250420120754849.png" alt="image-20250420120754849" style="zoom: 33%;" />



<img src="/Users/moon/Library/Application Support/typora-user-images/image-20250420120847809.png" alt="image-20250420120847809" style="zoom:33%;" />



Leader1:3

Leader2:1

other：0

最后 3 成为 leader 了

## 🔒 同步（Synchronization）

### 🧩 初始化（Initialization）

- 程序初始化只在一个 goroutine 中运行。
- 该 goroutine **可以创建其它 goroutine**，这些 goroutine 会并发运行。
- 如果包 `p` 导入了包 `q`，那么：
  - **`q` 包中的 `init` 函数会先执行完，才会开始执行 `p` 包的 `init` 函数**。
- 所有包的 `init` 函数执行完毕之后，才会开始执行 `main.main()` 函数。
- ✅ 这部分初始化是**串行同步**完成的。

------

### 🌀 goroutine 创建（Goroutine creation）

通过 `go` 关键字启动的新 goroutine，会**在当前 goroutine 中创建并同步启动**。

例如这个程序：

```
go


复制编辑
var a string

func f() {
	print(a)
}

func hello() {
	a = "hello, world"
	go f()
}
```

调用 `hello()` 时，它**保证 `f()` 会在某个未来时间点打印出 `"hello, world"`**（即使是在 `hello()` 返回之后）。

------

### 🔚 goroutine 销毁（Goroutine destruction）

goroutine 的退出**不保证与程序中其它事件的同步**。

例如这个程序：

```
go


复制编辑
var a string

func hello() {
	go func() { a = "hello" }()
	print(a)
}
```

此处 `a = "hello"` 的赋值操作没有任何同步保障，所以：

- 其它 goroutine（比如主线程）**可能观察不到这个赋值结果**。
- 编译器甚至可能**优化掉整个 `go` 语句**，因为它“看上去没影响”。

✅ 正确做法是使用**同步机制**（例如锁、channel）来确保顺序关系。

------

### 📬 Channel 通信（Channel communication）

Go 中最常用的同步手段就是 **channel 通信**。

- 每次 `channel` 的发送操作 `chan <- x`，会与一个接收 `<-chan` 匹配；
- 通常在**不同的 goroutine 中发生配对**。
- **发送一定在接收之前完成。**

#### 例子 1：有缓冲 channel（Buffered）

```
go


复制编辑
var c = make(chan int, 10)
var a string

func f() {
	a = "hello, world"
	c <- 0
}

func main() {
	go f()
	<-c
	print(a)
}
```

**✅ 一定会打印 "hello, world"**，因为：

1. `a = "hello, world"` 在 `c <- 0` 之前；
2. `c <- 0` 和 `<-c` 是同步点；
3. `<-c` 之后是 `print(a)`；
4. 所以 `a` 一定被赋值了。

> 替换 `c <- 0` 为 `close(c)`，行为也是一样的，因为关闭 channel 的操作也会被同步。

------

#### 例子 2：无缓冲 channel（Unbuffered）

```
go


复制编辑
var c = make(chan int)
var a string

func f() {
	a = "hello, world"
	<-c
}

func main() {
	go f()
	c <- 0
	print(a)
}
```

**✅ 同样会打印 "hello, world"**，因为：

- `a = ...` 发生在 `<-c` 之前；
- `<-c` 与 `c <-` 是同步点；
- 所以 `print(a)` 发生在写 `a` 之后。

> ⚠️ 若是改成 `make(chan int, 1)`（带缓冲），则无法保证 `a` 已经赋值，因为 send 不会阻塞。

------

#### 



## 📦 高级同步规则：缓冲 channel 的 “第 K 次接收” 同步语义

> **Go 保证**：第 `k` 次 `<-ch` 操作一定发生在第 `k + C` 次 `ch <-` **之前完成**（其中 `C` 是 channel 容量）。

这条规则是为了对缓冲 channel 的行为提供更可控的顺序保证。

### ✅ 举例说明

假设我们有一个缓冲区大小为 3 的 channel：

```
go


复制编辑
ch := make(chan int, 3)
```

我们发送了 5 个消息：

```
go


复制编辑
ch <- 1 // 第1次send
ch <- 2 // 第2次send
ch <- 3 // 第3次send （此时buffer满）

// 到这里都没阻塞

// 第4次send 需要等待前面的 receive，因为缓冲区已满
```

这时我们：

```
go


复制编辑
<-ch // 第1次recv，释放一个槽位
```

于是：

```
go


复制编辑
ch <- 4 // 第4次send 可以继续执行
```

### 🧠 高级语义：

- 第 `1` 次 `<-ch` 一定完成于第 `4` 次 `ch <-` 之前；
- 第 `2` 次 `<-ch` 一定完成于第 `5` 次 `ch <-` 之前；

换句话说：

> **第 k 次 recv 发生在第 k+C 次 send 之前完成。**

这种语义允许我们将 channel 当作一种**有限资源信号量**（如并发池）使用。

------

## 🎯 实际应用：用缓冲 channel 限制并发（信号量）

```
go


复制编辑
limit := make(chan struct{}, 3) // 容量为3，表示最大允许3个并发任务

for _, job := range jobs {
    limit <- struct{}{} // 获取一个信号，占用一个槽位（可能会阻塞）
    go func(j Job) {
        process(j)
        <-limit // 完成任务后释放槽位
    }(job)
}
```

### ✅ 保证：

- 任意时刻最多只有 `3` 个 goroutine 正在处理任务；
- 即使你同时启动了 1000 个 goroutine，也只会有 3 个活跃。

------

## 



## 🔐 锁（Locks）

Go 标准库 `sync` 提供了两种锁类型：`sync.Mutex` 和 `sync.RWMutex`。

### ✅ 锁的同步保证：

对于任意 `sync.Mutex` 或 `sync.RWMutex` 类型的变量 `l`，如果 `n < m`，那么：

> 第 n 次调用 `l.Unlock()` **先于** 第 m 次调用 `l.Lock()` 的返回。

------

### ✅ 例子说明：

```
go


复制编辑
var l sync.Mutex
var a string

func f() {
	a = "hello, world"
	l.Unlock()
}

func main() {
	l.Lock()       // 🔐 获得锁
	go f()         // 🚀 启动 goroutine 写入 a，然后释放锁
	l.Lock()       // 🚪 等待 f() 解锁后继续
	print(a)       // ✅ 打印结果
}
```

上述代码一定会打印 `"hello, world"`，因为 `f()` 中的 `l.Unlock()` 与 `main` 中的第二次 `l.Lock()` 建立了同步顺序，`print(a)` 肯定发生在 `a = "hello, world"` 之后。

------

### 📖 关于 `RWMutex`

- 对 `sync.RWMutex` 使用 `l.RLock()` 时：
  - 存在某个 `l.Unlock()` 的调用（假设是第 n 次），它 **先于** `l.RLock()` 返回。
  - 对应的 `l.RUnlock()` **先于** 第 n+1 次对 `l.Lock()` 的返回。

------

## 🔁 TryLock 行为说明：

- `l.TryLock()` / `l.TryRLock()` 成功时，等价于 `Lock()`。
- **失败时**（即锁被占用），不产生任何同步效果。
- **注意**：Go 的内存模型允许 `TryLock()` 在锁是空闲的情况下返回 `false` —— 这是语言层的保守设计，避免依赖未同步行为。

------

## 🧩 Once（只执行一次的初始化）

Go 的 `sync.Once` 提供线程安全的一次性初始化机制。

```
go


复制编辑
var once sync.Once
once.Do(f) // 多个 goroutine 并发调用，只会执行一次 f()
```

- 所有对 `once.Do(f)` 的调用，都会在 `f()` 完成之后才返回。
- 如果某 goroutine 在执行 `f()`，其他并发调用将会 **阻塞等待**。

------

### ✅ Once 的例子

```
go


复制编辑
var a string
var once sync.Once

func setup() {
	a = "hello, world"
}

func doprint() {
	once.Do(setup)
	print(a)
}

func twoprint() {
	go doprint()
	go doprint()
}
```

调用 `twoprint()` 会保证：

- `setup()` 只被调用一次；
- 两个 `print(a)` 都打印 `"hello, world"`；
- 并且在 `setup()` 执行完成后才执行 `print(a)`。

------

## ⚛️ Atomic Values（原子操作）

Go 的 `sync/atomic` 包提供了原子操作函数，用于在并发环境下实现无锁安全访问。

### ✅ 同步保证：

> 如果原子操作 A 的结果被 B 观察到了，那么 A 就**先于** B 被执行（synchronized-before）。

- 所有原子操作在程序中表现为**顺序一致（sequentially consistent）**。
- 对应于 C++ 的 `std::atomic`，Java 的 `volatile`。

------

## 🧹 Finalizers（终结器）

- Go 的 `runtime.SetFinalizer(x, f)` 用于注册对象 `x` 的终结函数 `f(x)`。
- 保证：`SetFinalizer(x, f)` 调用 **先于** `f(x)` 被执行。





你的总结非常到位，下面我将**更深入地解释这些错误模式为什么会出现**，以及**为什么必须使用显式同步手段（如 channel、`sync.Mutex`、`sync.WaitGroup`）来避免问题**。

------

## ❌ 错误同步行为详解

### 📌 原则复述

> 在并发程序中，如果**两个 goroutine 同时访问一个变量，其中至少有一个是写操作**，**且没有使用同步手段（锁、channel、原子操作等）来保护访问**，就会导致“数据竞争（race condition）”。

数据竞争的结果是**未定义的行为**，即你**无法预期程序会输出什么**，也不能依赖特定的值是否已经“写入成功”。

------

## 🧨 案例 1：无同步写入和读取

```go
var a, b int

func f() {
	a = 1
	b = 2
}

func g() {
	print(b)
	print(a)
}

func main() {
	go f()
	g()
```

### ❗ 为什么可能输出 `2 0`？

- `f()` 在 goroutine 中执行，但 **main 线程可能在 `f()` 完成前就开始执行 `g()`**；
- `g()` 中先读 `b`（读到 2），再读 `a`（读到初始值 0）；
- 原因是 **goroutine 与主线程之间没有任何同步机制**；
- 所以，`g()` 看到了 `b = 2`，但**不能保证**也能看到 `a = 1`。

### ✅ 正确写法

```go
var a, b int
var wg sync.WaitGroup

func f() {
	a = 1
	b = 2
	wg.Done()
}

func g() {
	wg.Wait()
	print(b)
	print(a)
}

func main() {
	wg.Add(1)
	go f()
	g()
}
```

------

## 🧨 案例 2：双重检查锁定（Double-Checked Locking）

```
go


复制编辑
var a string
var done bool

func setup() {
	a = "hello, world"
	done = true
}

func doprint() {
	if !done {
		once.Do(setup)
	}
	print(a)
}
```

### ❗ 问题在哪？

即使你看到 `done == true`，也**不能保证看到 `a = "hello, world"`**。 这是因为：

- `done` 是一个布尔值；
- `a` 是一个字符串；
- 它们之间**没有同步边（happens-before）**；
- 在并发下，CPU 和编译器可能会**乱序执行**（写入 a 在 done 之后或重排序）；
- 另一个 goroutine **可能先看到 `done = true`，但此时 `a` 还没设置好**！

### ✅ 正确做法

```
go


复制编辑
var a string
var once sync.Once

func setup() {
	a = "hello, world"
}

func doprint() {
	once.Do(setup)
	print(a)
}
```

------

## 🧨 案例 3：忙等（Busy Waiting）

```
go


复制编辑
var a string
var done bool

func setup() {
	a = "hello, world"
	done = true
}

func main() {
	go setup()
	for !done {
	}
	print(a)
}
```

### ❗ 典型错误点

- 虽然主线程在循环等 `done == true`；
- 但 Go 编译器和 CPU 优化器可能让它 **永远看不到最新的 done 值**；
- 更重要的是，即使 `done == true`，你也不能**指望看到 `a = "hello, world"`**，因为缺乏同步事件。

### ✅ 正确做法

```
go


复制编辑
var a string
var ch = make(chan struct{})

func setup() {
	a = "hello, world"
	close(ch)
}

func main() {
	go setup()
	<-ch  // 阻塞直到 setup 完成
	print(a)
}
```

------

## 🧨 案例 4：对象赋值的同步问题

```
go


复制编辑
type T struct {
	msg string
}

var g *T

func setup() {
	t := new(T)
	t.msg = "hello, world"
	g = t
}

func main() {
	go setup()
	for g == nil {
	}
	print(g.msg)
}
```

### ❗ 为什么出错？

- `main` 线程看到 `g != nil`，但 `t.msg` 的写入不一定已经“同步可见”；
- `g = t` 和 `t.msg = ...` 虽然在 `setup()` 中顺序执行，但另一个线程未同步看到；
- 所以可能会打印空字符串或发生 **“部分初始化”错误**（可能看到 g 是非空指针，但 msg 还是空）。

### ✅ 正确做法

```
go


复制编辑
type T struct {
	msg string
}

var g *T
var ch = make(chan struct{})

func setup() {
	t := new(T)
	t.msg = "hello, world"
	g = t
	close(ch)
}

func main() {
	go setup()
	<-ch
	print(g.msg)
}
```





<img src="/Users/moon/Library/Application Support/typora-user-images/image-20250419211135114.png" alt="image-20250419211135114" style="zoom:50%;" />





