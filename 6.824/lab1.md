

https://blog.josejg.com/debugging-pretty/



> 你在使用 `go build -buildmode=plugin ../mrapps/wc.go` 时遇到的这个错误：
>
> ```
> cpp
> 
> 
> 复制编辑
> "../mr" is relative, but relative import paths are not supported in module mode
> ```
>
> 是因为 **Go 的 module 模式**（也就是启用了 `go.mod` 的项目）**不允许使用相对路径导入包**。
>
> ------
>
> ### ✅ 错误原因解析：
>
> 你可能在 `../mrapps/wc.go` 中写了类似这样的 import：
>
> ```
> go
> 
> 
> 复制编辑
> import "../mr"
> ```
>
> 这是合法的 **GOPATH 模式** 写法，但在 **module 模式** 中会报错。
>
> ------
>
> ### ✅ 解决方法一：修改 `import` 为模块路径
>
> 1. 查看 `go.mod` 中的 module 名称，比如：
>
> ```
> go
> 
> 
> 复制编辑
> module github.com/yourname/6.824
> ```
>
> 1. 把 `import "../mr"` 改为：
>
> ```
> go
> 
> 
> 复制编辑
> import "github.com/yourname/6.824/mr"
> ```
>
> 1. 再重新编译：
>
> ```
> bash
> 
> 
> 复制编辑
> go build -buildmode=plugin ../mrapps/wc.go
> ```
>
> ------
>
> ### ✅ 解决方法二：使用 GOPATH 模式（不推荐）
>
> 如果你不想用 `go.mod`，可以切换到老的 GOPATH 模式：
>
> ```
> bash
> 
> 
> 复制编辑
> export GO111MODULE=off
> ```
>
> 然后再执行：
>
> ```
> bash
> 
> 
> 复制编辑
> go build -buildmode=plugin ../mrapps/wc.go
> ```
>
> > ⚠️ 这种方式不推荐使用，因为现在 Go 的主流开发方式是基于 module 的。
>
> ------

### 

Go 项目**默认建议**一个项目一个模块（就是顶层有一个 `go.mod`），
 但**在以下场景中，确实需要多个模块：**

### 📌 多模块场景举例：

| 场景                         | 举例                                              |
| ---------------------------- | ------------------------------------------------- |
| 拆分为多个独立可发布组件     | `core/`, `client/`, `plugin/` 各有自己的 `go.mod` |
| 插件开发（plugin buildmode） | 每个 `.so` 插件目录可以是一个独立模块             |
| 版本隔离                     | 某些子模块依赖不同版本的库                        |
| 测试多个实验性组件           | 不希望主项目的 `go.mod` 被污染                    |





📦 模块（Module）
│
├── 📁 包 package main
│     └── main.go
│
├── 📁 包 package utils
│     └── math.go
│
└── 📁 包 package service
      └── api.go



## 🧱 `go mod init` 是什么？

### ✅ 作用：**初始化当前目录为一个 Go 模块项目**

### 🛠 使用方式：

```
bash


复制编辑
go mod init <模块名>
```

### 📝 会做什么？

1. 生成一个 `go.mod` 文件（模块的配置文件）
2. 指定模块路径（如 `github.com/yourname/project`）
3. 从你当前项目的 `import` 路径中识别出依赖（只要能找到）

### 💡 举个例子：

假如你的项目目录是：

```
bash


复制编辑
/Users/moon/GolandProjects/hello
```

你在这个目录下执行：

```
bash


复制编辑
go mod init hello
```

就会生成一个 `go.mod` 文件，内容类似于：

```
go


复制编辑
module hello

go 1.20
```

------

## 🧹 `go mod tidy` 是什么？

### ✅ 作用：**自动整理模块依赖**

### 🚀 会做哪些事情？

1. **添加** 代码中用到但 `go.mod` 没有的依赖
2. **移除** `go.mod` 中存在但代码没用到的依赖
3. **生成或更新 `go.sum`**（校验依赖哈希用）

### 🧪 什么时候用？

- 引入新包后
- 删除代码或包后
- 拉别人的代码后第一次构建
- 执行完 `go get` 等命令后

### 📌 命令：

```
bash


复制编辑
go mod tidy
```





```
// 1。将每个文件新建线程调用map处理
for index, fileName := range m.fileNames {
    go func(fileName string) {
       // 1.1 提取content
       file, err := os.Open(fileName)
       if err != nil {
          log.Fatalf("can't open %v", fileName)
       }
       content, err := ioutil.ReadAll(file)
       if err != nil {
          log.Fatalf("can't read %v", fileName)
       }
       // 1.2 调用map map中完成排序输出到对应文件

    }(fileName)
}
// 2.等待所有map结束后，调用reduce-Y(k, file)处理

// 3.等待所有reduce结束后，返回
// Your code here.
```





2025/03/29 09:18:21 rpc: gob error encoding body: gob: type mr.TaskReply has no exported fields,reading body unexpected EOF

**你的 `mr.TaskReply` 结构体没有任何**“导出字段”**，即所有字段都是小写开头的私有字段，gob 编码器无法序列化它们。**

需要确保 `TaskReply`（以及其它参与 RPC 的 struct 类型）的字段是导出的：字段名首字母必须大写。







在 Go 中，**要想在其他包中访问变量/函数/结构体字段，它们必须以大写字母开头**。

### 🔧 你当前的方法定义：

```
go


复制编辑
func (m *Master) MapWorkFinished(args *ExampleArgs, reply *TaskReply) {
```

你少了一个返回值 `error`，**不满足第 4 条规则**，因此不能被注册为 RPC 方法！

####  一个方法能作为 RPC 方法，必须满足：

1. **方法是导出（首字母大写）的**
2. **接收者是导出类型的指针类型**（如 `*Master`）
3. **方法有两个参数**，且都是导出类型指针：
   - 第一个参数是 `args *SomeArgs`
   - 第二个参数是 `reply *SomeReply`
4. **方法返回类型为 `error`**



### 问题：reduce parallelism test: FAIL

所以整体执行顺序如下：

1. 第一行 worker（Worker A）启动，并在后台运行
2. 第二行 worker（Worker B）启动，并阻塞脚本，直到退出
3. 等 B 退出后，脚本才继续执行下面这行：

```
NT=$(cat mr-out* | grep '^[a-zA-Z] 2' | wc -l | sed 's/ //g')
```



time.Sleep(1 * time.Second) 去掉可以通过，可能是时间太长执行不完导致的失败



```
两个 reduce worker 几乎同时跑：

Worker A：pid=90280

Worker B：pid=90281

当 A 执行 nparallel("reduce")：

创建文件 mr-worker-reduce-90280

看到两个文件（自己 + B）

检查两个 pid 都 alive → ret = 2

同理 B 也看到两个 worker → ret = 2

func nparallel(phase string) int {
    // create a file so that other workers will see that
    // we're running at the same time as them.
    pid := os.Getpid()
    myfilename := fmt.Sprintf("mr-worker-%s-%d", phase, pid)
    err := ioutil.WriteFile(myfilename, []byte("x"), 0666)
    if err != nil {
       panic(err)
    }

    // are any other workers running?
    // find their PIDs by scanning directory for mr-worker-XXX files.
    dd, err := os.Open(".")
    if err != nil {
       panic(err)
    }
    names, err := dd.Readdirnames(1000000)
    if err != nil {
       panic(err)
    }
    ret := 0
    for _, name := range names {
       var xpid int
       pat := fmt.Sprintf("mr-worker-%s-%%d", phase)
       n, err := fmt.Sscanf(name, pat, &xpid)
       if n == 1 && err == nil {
          err := syscall.Kill(xpid, 0)
          if err == nil {
             // if err == nil, xpid is alive.
             ret += 1
          }
       }
    }
    dd.Close()

    time.Sleep(1 * time.Second)

    err = os.Remove(myfilename)
    if err != nil {
       panic(err)
    }

    return ret
}
```





### Map 的随机会带来不确定性

```go
func Shuffle(kvs []KeyValue, x int) {
	log.Printf("begin Shuffle: %v", x)
	if len(kvs) == 0 {
		return
	}
	sort.Sort(ByKey(kvs))

	buckets := make(map[int][]KeyValue)
	for _, kv := range kvs {
		y := ihash(kv.Key) % 10
		buckets[y] = append(buckets[y], kv)
	}

	for y := 0; y < 10; y++ {
		// 1. 写临时文件
		tempName := fmt.Sprintf("mr-%d-%d-tmp", x, y)
		file, err := os.Create(tempName)
		if err != nil {
			log.Fatalf("创建临时文件失败: %v", err)
		}
		enc := json.NewEncoder(file)
		for _, kv := range buckets[y] {
			if err := enc.Encode(&kv); err != nil {
				log.Fatalf("写入临时文件失败: %v", err)
			}
		}
		file.Close()

		// 2. 重命名为正式文件
		finalName := fmt.Sprintf("mr-%d%d", x, y)
		err = os.Rename(tempName, finalName)
		if err != nil {
			log.Fatalf("重命名文件失败: %v", err)
		}
		log.Printf("写入中间文件: %s", finalName)
	}
}
```

这个问题非常关键，咱们一起来拆开解释：

------

## 🎯 你问的是：

### 为什么：

```go
for y := 0; y < 10; y++ { ... }
```

是对的，而：

```go
for y, kvList := range buckets { ... }
```

**是错的**？

------

## 🔍 核心区别在于：

### ✅ 第一种：`for y := 0; y < 10; y++` 是 **固定顺序、完整、无遗漏的遍历**

你显式指定了 **reduce ID 从 0 到 9**（共 10 个桶），即使某些桶是空的，也不会跳过。

它的含义是：

> 我不管 `buckets[y]` 有没有数据，我都遍历 0 到 9，每个中间文件我都确保创建、写入（哪怕是空）。

这在 MapReduce 中是非常重要的，因为：

- 每个 reduce 任务会读取所有 map 生成的 `mr-X-Y` 文件
- 如果你漏写某个 reduce 的中间文件，比如 `mr-3-7`，Reduce 3 读不到，就挂了

------

### ❌ 第二种：`for y, kvList := range buckets` 是 **不稳定、可能遗漏、顺序随机的遍历**

这是 Go 对 `map[int][]KeyValue` 的遍历，Go 官方文档明确说明：

> **Map 的遍历顺序是随机的，并且在每次运行时都不同。**

#### 更严重的是：

如果某个 `buckets[y]` 是空的，**它根本不会出现在遍历结果中**！

比如：

```go
buckets := map[int][]KeyValue{
  0: {...},
  2: {...},
  5: {...},
}
```

上面遍历中 **1、3、4、6、7、8、9** 全部会被跳过！

❗那就意味着你不会生成 `mr-x-1`, `mr-x-3`, ..., 这些 reduce 的输入文件被 **意外遗漏了**，导致 reduce 阶段无法执行。

------

## ✅ 举个实际例子

```go
// 有 10 个 reduce bucket，理论上应写 10 个文件
buckets := map[int][]KeyValue{}
buckets[0] = [...] // 有数据
buckets[5] = [...] // 有数据
// 其余 1~4, 6~9 没有数据，默认不会出现在 map 中
```

**错误做法：**

```go
for y, kvList := range buckets { // 只写 0、5
	file := os.Create(fmt.Sprintf("mr-%d-%d", x, y))
	...
}
```

**正确做法：**

```go
for y := 0; y < 10; y++ {
	kvList := buckets[y] // 空也 OK
	file := os.Create(fmt.Sprintf("mr-%d-%d", x, y))
	...
}
```

✅ **无论有没有值，我都生成这个中间文件**，这是对 Reduce 最友好的方式！

------

## ✅ 最终结论

| 方式                           | 是否稳定顺序 | 会不会漏写 | 是否推荐 |
| ------------------------------ | ------------ | ---------- | -------- |
| `for y := 0; y < NReduce; y++` | ✅ 是         | ❌ 不会     | ✅ 推荐   |
| `for y, _ := range buckets`    | ❌ 否         | ✅ 可能会漏 | ❌ 不推荐 |

------

## 🧠 补充建议

如果你硬要用 `map[int][]KeyValue` 结构，也建议这么遍历：

```go
for y := 0; y < NReduce; y++ {
	kvList := buckets[y] // 没有就默认是 nil
	// 一样写文件即可
}
```

这样既保证顺序，又不会遗漏 reduce 所需的输入文件。

------

如你愿意，我可以再帮你写一个健壮的 `Shuffle()` 版本，包含：

- 临时文件 + rename
- 保证输出所有中间文件
- 错误处理

要不要？😎

缺失的是 263

```
send id: 3, msg: {true {Append 3 x 3 14 y 52 1878463325490364172} 263}
```

3 给 2 发了到 265 前所有日志

```
3 finish sendInstallSnapshotRpc to 2.....nextIndex:266, matchIndex:265 
```



2 发送快照以后是 last 265   

下次发心跳开始args.preIndex: 265



```
2025/05/07 22:37:55 AppendEntries,args.leaderId:3,args.epoch:1，args.preIndex: 265, args.preTerm: 1, rf.me：2,rf.Term：1,rf.voteFor：1,rf.state：0, rf.log：[{0 20516}], args.log:[{1 {Append 4 x 4 31 y 57 2391327139425014098}} {1 {Append 1 x 1 29 y 54 1954610769135245984}} {1 {Append 0 x 0 22 y 53 2668356850362338520}} {1 {Append 2 x 2 14 y 48 2093277762212500425}} {1 {Append 3 x 3 15 y 53 1878463325490364172}} {1 {Get 3  54 1878463325490364172}} {1 {Get 2  49 2093277762212500425}} {1 {Get 1  55 1954610769135245984}} {1 {Get 0  54 2668356850362338520}}]
```



```
2025/05/08 11:50:21 2 begin CupLogExceedMaxSizeAndSaveSnapShot.....rf.LastIncludedIndex:310,rf.ApplyId:319, rf.log:[{0 20516} {1 {Append 2 x 2 5 y 7 4135897827882552221}} {1 {Append 3 x 3 2 y 6 2733491079772594813}} {1 {Append 0 x 0 6 y 10 4293905909427359528}} {1 {Get 4  8 1430075913155013695}} {1 {Get 1  8 1808890630129367958}} {1 {Append 2 x 2 6 y 8 4135897827882552221}} {1 {Get 1  9 1808890630129367958}} {1 {Append 3 x 3 3 y 7 2733491079772594813}} {1 {Append 4 x 4 2 y 9 1430075913155013695}} {1 {Get 2  9 4135897827882552221}}], kvs:map[0:x 0 0 yx 0 1 yx 0 2 yx 0 3 yx 0 4 yx 0 5 yx 0 6 y 1:x 1 0 yx 1 1 y 2:x 2 0 yx 2 1 yx 2 2 yx 2 3 yx 2 4 yx 2 5 yx 2 6 y 3:x 3 0 yx 3 1 yx 3 2 y 4:x 4 0 yx 4 1 y]
2025/05/08 11:50:21 2 getRelativeLogIndex rf.LastIncludedIndex:310
2025/05/08 11:50:21 2 getRelativeLogIndex rf.LastIncludedIndex:310
2025/05/08 11:50:21 3 getRelativeLogIndex rf.LastIncludedIndex:308
2025/05/08 11:50:21 before persist....rf.me:2, rf.CurrentTerm:1, rf.votedFor:2, rf.log:[{0 20516} {1 {Get 2  9 4135897827882552221}}]
2025/05/08 11:50:21 0 begin, apply: {Get 4  8 1430075913155013695}
2025/05/08 11:50:21 0 apply, get: x 4 0 yx 4 1 y
2025/05/08 11:50:21 0 begin getState
2025/05/08 11:50:21 GetState release lock
2025/05/08 11:50:21 {true {Get 4  8 1430075913155013695} 314} chan don't exist
2025/05/08 11:50:21 3 begin CupLogExceedMaxSizeAndSaveSnapShot.....rf.LastIncludedIndex:308,rf.ApplyId:314, rf.log:[{0 20516} {1 {Get 1  7 1808890630129367958}} {1 {Append 0 x 0 5 y 9 4293905909427359528}} {1 {Append 2 x 2 5 y 7 4135897827882552221}} {1 {Append 3 x 3 2 y 6 2733491079772594813}} {1 {Append 0 x 0 6 y 10 4293905909427359528}} {1 {Get 4  8 1430075913155013695}} {1 {Get 1  8 1808890630129367958}} {1 {Append 2 x 2 6 y 8 4135897827882552221}} {1 {Get 1  9 1808890630129367958}} {1 {Append 3 x 3 3 y 7 2733491079772594813}} {1 {Append 4 x 4 2 y 9 1430075913155013695}}], kvs:map[0:x 0 0 yx 0 1 yx 0 2 yx 0 3 yx 0 4 yx 0 5 y 1:x 1 0 yx 1 1 y 2:x 2 0 yx 2 1 yx 2 2 yx 2 3 yx 2 4 yx 2 5 y 3:x 3 0 yx 3 1 yx 3 2 y 4:x 4 0 yx 4 1 y]
2025/05/08 11:50:21 3 getRelativeLogIndex rf.LastIncludedIndex:308
2025/05/08 11:50:21 id: 3, msg: {true {Get 1  8 1808890630129367958} 315}
2025/05/08 11:50:21 send id: 3, msg: {true {Get 1  8 1808890630129367958} 315}
2025/05/08 11:50:21 finish persist....
2025/05/08 11:50:21 2 finish CupLogExceedMaxSizeAndSaveSnapShot.....log:[{0 20516} {1 {Get 2  9 4135897827882552221}}], rf.LastIncludedIndex:319, rf.LastIncludedTerm:1
2025/05/08 11:50:21 msg 内容：{true {Append 3 x 3 3 y 7 2733491079772594813} 318}
2025/05/08 11:50:21 2 begin, apply: {Append 3 x 3 3 y 7 2733491079772594813}
```



2 给 0,1，3发了

```
msg 内容：{true {Append 3 x 3 3 y 7 2733491079772594813} 318}
```



```
2025/05/08 11:50:21 4 begin InstallSnapshotRpcHandler.....snapShotReq.LastIncludedIndex:319
2025/05/08 11:50:21 4 getRelativeLogIndex rf.LastIncludedIndex:308
2025/05/08 11:50:21 before persist....rf.me:4, rf.CurrentTerm:1, rf.votedFor:2, rf.log:[{0 20516}]
2025/05/08 11:50:21 finish persist....
2025/05/08 11:50:21 finish InstallSnapshotRpcHandler.....rf.log:[{0 20516}]
2025/05/08 11:50:21 receve snapShot
2025/05/08 11:50:21 4 begin readSnapShot....
2025/05/08 11:50:21 2 finish sendInstallSnapshotRpc to 4.....nextIndex:320, matchIndex:319
2025/05/08 11:50:21 4 finish readSnapShot...., kv.kvs:map[0:x 0 0 yx 0 1 yx 0 2 yx 0 3 yx 0 4 yx 0 5 yx 0 6 y 1:x 1 0 yx 1 1 y 2:x 2 0 yx 2 1 yx 2 2 yx 2 3 yx 2 4 yx 2 5 yx 2 6 y 3:x 3 0 yx 3 1 yx 3 2 y 4:x 4 0 yx 4 1 y]
```



177 的快照不应该包括 179 的内容

sss

msg 内容：{true {Append 1 x 1 15 y 32 1950430866193101281} 174}

3 begin start putappend: {Append 1 x 1 15 y 32 1950430866193101281} 174

2025/05/08 17:41:15 3 finish CupLogExceedMaxSizeAndSaveSnapShot.....log:[{0 20516} {59 {Append 1 x 1 15 y 32 1950430866193101281}} {59 {Append 2 x 2 14 y 29 2091085765881772888}} {59 {Get 0  32 777403258639395625}} {59 {Get 3  31 3698835677464332326}} {62 {Get 4  32 3073720048189741165}} {62 {Append 1 x 1 15 y 32 1950430866193101281}} {62 {Append 2 x 2 14 y 29 2091085765881772888}}], rf.LastIncludedIndex:173, rf.LastIncludedTerm:59

【两个{Append 1 x 1 15 y 32 1950430866193101281} 】

```
2025/05/08 17:41:14 msg 内容：{true {Get 1  31 1950430866193101281} 167}
2025/05/08 17:41:14 4 begin, apply: {Get 1  31 1950430866193101281}
2025/05/08 17:41:14 4 apply, get: x 1 0 yx 1 1 yx 1 2 yx 1 3 yx 1 4 yx 1 5 yx 1 6 yx 1 7 yx 1 8 yx 1 9 yx 1 10 yx 1 11 yx 1 12 yx 1 13 yx 1 14 y
```



2025/05/08 17:41:15 4 begin InstallSnapshotRpcHandler.....snapShotReq.LastIncludedIndex:177
2025/05/08 17:41:15 4 getRelativeLogIndex rf.LastIncludedIndex:163
2025/05/08 17:41:15 before persist....rf.me:4, rf.CurrentTerm:64, rf.votedFor:3, rf.log:[{0 20516}]
2025/05/08 17:41:15 finish persist....
2025/05/08 17:41:15 finish InstallSnapshotRpcHandler.....rf.log:[{0 20516}]
2025/05/08 17:41:15 3 finish sendInstallSnapshotRpc to 4.....nextIndex:178, matchIndex:177
2025/05/08 17:41:15 receve snapShot
2025/05/08 17:41:15 4 begin readSnapShot....
2025/05/08 17:41:15 4 finish readSnapShot...., kv.kvs:map[0:x 0 0 yx 0 1 yx 0 2 yx 0 3 yx 0 4 yx 0 5 yx 0 6 yx 0 7 yx 0 8 yx 0 9 yx 0 10 yx 0 11 yx 0 12 yx 0 13 yx 0 14 yx 0 15 yx 0 16 yx 0 17 y 1:x 1 0 yx 1 1 yx 1 2 yx 1 3 yx 1 4 yx 1 5 yx 1 6 yx 1 7 yx 1 8 yx 1 9 yx 1 10 yx 1 11 yx 1 12 yx 1 13 yx 1 14 yx 1 15 y 2:x 2 0 yx 2 1 yx 2 2 yx 2 3 yx 2 4 yx 2 5 yx 2 6 yx 2 7 yx 2 8 yx 2 9 yx 2 10 yx 2 11 yx 2 12 yx 2 13 yx 2 14 y 3:x 3 0 yx 3 1 yx 3 2 yx 3 3 yx 3 4 yx 3 5 yx 3 6 yx 3 7 yx 3 8 yx 3 9 yx 3 10 yx 3 11 yx 3 12 yx 3 13 yx 3 14 yx 3 15 yx 3 16 yx 3 17 yx 3 18 y 4:x 4 0 yx 4 1 yx 4 2 yx 4 3 yx 4 4 yx 4 5 yx 4 6 yx 4 7 yx 4 8 yx 4 9 yx 4 10 yx 4 11 yx 4 12 yx 4 13 y]





2025/05/08 17:41:15 AppendEntries,args.leaderId:3,args.epoch:64，args.preIndex: 177, args.preTerm: 59, rf.me：0,rf.Term：64,rf.voteFor：3,rf.state：0, rf.log：[{0 20516}], args.log:[{62 {Get 4  32 3073720048189741165}} {62 {Append 1 x 1 15 y 32 1950430866193101281}} {62 {Append 2 x 2 14 y 29 2091085765881772888}} {62 {Get 4  33 3073720048189741165}} {62 {Get 1  33 1950430866193101281}} {62 {Append 0 x 0 18 y 33 777403258639395625}} {62 {Append 3 x 3 19 y 32 3698835677464332326}}]
2025/05/08 17:41:15 0 getAbsLogIndex rf.LastIncludedIndex:177
2025/05/08 17:41:15 0 getRelativeLogIndex rf.LastIncludedIndex:177
2025/05/08 17:41:15 before leaderId: 3, rf.me: 0, LeaderCommit: 184, rf.log:[{0 20516}]
2025/05/08 17:41:15 after leaderId: 3, rf.me: 0, LeaderCommit: 184, rf.log:[{0 20516} {62 {Get 4  32 3073720048189741165}} {62 {Append 1 x 1 15 y 32 1950430866193101281}} {62 {Append 2 x 2 14 y 29 2091085765881772888}} {62 {Get 4  33 3073720048189741165}} {62 {Get 1  33 1950430866193101281}} {62 {Append 0 x 0 18 y 33 777403258639395625}} {62 {Append 3 x 3 19 y 32 3698835677464332326}}]
2025/05/08 17:41:15 AppendEntries,args.leaderId:3,args.epoch:64，args.preIndex: 177, args.preTerm: 59, rf.me：4,rf.Term：64,rf.voteFor：3,rf.state：0, rf.log：[{0 20516}], args.log:[{62 {Get 4  32 3073720048189741165}} {62 {Append 1 x 1 15 y 32 1950430866193101281}} {62 {Append 2 x 2 14 y 29 2091085765881772888}} {62 {Get 4  33 3073720048189741165}} {62 {Get 1  33 1950430866193101281}} {62 {Append 0 x 0 18 y 33 777403258639395625}} {62 {Append 3 x 3 19 y 32 3698835677464332326}}]



```
025/05/08 17:41:15 0 before apply ApplyIdid1: 178, rf.CommitId: 184, rf.getAbsLogIndex(len(rf.log)):185,relativeIndex:3
2025/05/08 17:41:15 4 begin, apply: {Append 1 x 1 15 y 32 1950430866193101281}
2025/05/08 17:41:15 4 apply, append: x 1 0 yx 1 1 yx 1 2 yx 1 3 yx 1 4 yx 1 5 yx 1 6 yx 1 7 yx 1 8 yx 1 9 yx 1 10 yx 1 11 yx 1 12 yx 1 13 yx 1 14 yx 1 15 yx 1 15 y
```



{Append 1 x 1 15 y 32 1950430866193101281}  179







```
2025/05/08 17:41:18 4 begin CupLogExceedMaxSizeAndSaveSnapShot.....rf.LastIncludedIndex:0,rf.ApplyId:1, rf.log:[{0 20516} {62 {Get 4  32 3073720048189741165}} {62 {Append 1 x 1 15 y 32 1950430866193101281}} {62 {Append 2 x 2 14 y 29 2091085765881772888}} {62 {Get 4  33 3073720048189741165}} {62 {Get 1  33 1950430866193101281}} {62 {Append 0 x 0 18 y 33 777403258639395625}} {62 {Append 3 x 3 19 y 32 3698835677464332326}} {62 {Append 0 x 0 18 y 33 777403258639395625}} {62 {Append 3 x 3 19 y 32 3698835677464332326}} {65 {Get 0  10 883003686767608302}}], kvs:map[0:x 0 0 yx 0 1 yx 0 2 yx 0 3 yx 0 4 yx 0 5 yx 0 6 yx 0 7 yx 0 8 yx 0 9 yx 0 10 yx 0 11 yx 0 12 yx 0 13 yx 0 14 yx 0 15 yx 0 16 yx 0 17 y 1:x 1 0 yx 1 1 yx 1 2 yx 1 3 yx 1 4 yx 1 5 yx 1 6 yx 1 7 yx 1 8 yx 1 9 yx 1 10 yx 1 11 yx 1 12 yx 1 13 yx 1 14 yx 1 15 y 2:x 2 0 yx 2 1 yx 2 2 yx 2 3 yx 2 4 yx 2 5 yx 2 6 yx 2 7 yx 2 8 yx 2 9 yx 2 10 yx 2 11 yx 2 12 yx 2 13 yx 2 14 y 3:x 3 0 yx 3 1 yx 3 2 yx 3 3 yx 3 4 yx 3 5 yx 3 6 yx 3 7 yx 3 8 yx 3 9 yx 3 10 yx 3 11 yx 3 12 yx 3 13 yx 3 14 yx 3 15 yx 3 16 yx 3 17 yx 3 18 y 4:x 4 0 yx 4 1 yx 4 2 yx 4 3 yx 4 4 yx 4 5 yx 4 6 yx 4 7 yx 4 8 yx 4 9 yx 4 10 yx 4 11 yx 4 12 yx 4 13 y]
2025/05/08 17:41:18 4 getRelativeLogIndex rf.LastIncludedIndex:0
```



刚开始日志提交了两个重复的 append，假如 server 挂了重新恢复的时候，维护的每个 cli最后一个的值没了，这个时候重新执行就会执行成功。





2 正常之前

```
2025/05/09 02:30:53 msg 内容：{true {Get 4  7 402244904228922384} 38}
2025/05/09 02:30:53 2 begin, apply: {Get 4  7 402244904228922384}
2025/05/09 02:30:53 2 apply, get: x 4 0 yx 4 1 yx 4 2 yx 4 3 y
```





2025/05/09 02:30:53 AppendEntries,args.leaderId:3,args.epoch:1，args.preIndex: 38, args.preTerm: 1, rf.me：2,rf.Term：1,rf.voteFor：3,rf.state：0, rf.log：[{0 20516} {1 {Append 0 x 0 4 y 7 1123720622585023599}} {1 {Get 4  6 402244904228922384}} {1 {Append 2 x 2 0 y 5 1228272264126474379}} {1 {Append 1 x 1 2 y 7 2016228440974134146}} {1 {Append 0 x 0 5 y 8 1123720622585023599}} {1 {Append 3 x 3 4 y 6 1098643115262335781}} {1 {Get 4  7 402244904228922384}}], args.log:[{1 {Get 2  6 1228272264126474379}} {1 {Append 0 x 0 6 y 9 1123720622585023599}} {1 {Append 3 x 3 5 y 7 1098643115262335781}} {1 {Append 4 x 4 4 y 8 402244904228922384}} {1 {Append 1 x 1 3 y 8 2016228440974134146}}]

2025/05/09 02:30:53 3sendRequestAppendEntries to 2  after rf.matchIndex: 43, rf.nextIndex: 44



中间 leader 发了一次心跳，这个时候 2 日志里还有

```
2025/05/09 02:30:53 AppendEntries,args.leaderId:3,args.epoch:1，args.preIndex: 43, args.preTerm: 1, rf.me：2,rf.Term：1,rf.voteFor：3,rf.state：0, rf.log：[{0 20516} {1 {Append 0 x 0 5 y 8 1123720622585023599}} {1 {Append 3 x 3 4 y 6 1098643115262335781}} {1 {Get 4  7 402244904228922384}} {1 {Get 2  6 1228272264126474379}} {1 {Append 0 x 0 6 y 9 1123720622585023599}} {1 {Append 3 x 3 5 y 7 1098643115262335781}} {1 {Append 4 x 4 4 y 8 402244904228922384}} {1 {Append 1 x 1 3 y 8 2016228440974134146}}], args.log:[]
```

根本原因是因为 2 落后太多触发快照，然后快照解码失败了，只更新了LastIncludedIndex，没有更新完 kvs

```
025/05/09 02:30:53 2 begin InstallSnapshotRpcHandler.....snapShotReq.LastIncludedIndex:44
2025/05/09 02:30:53 2 getRelativeLogIndex rf.LastIncludedIndex:35
2025/05/09 02:30:53 before persist....rf.me:2, rf.CurrentTerm:1, rf.votedFor:3, rf.log:[{0 20516}]
2025/05/09 02:30:53 finish persist....
2025/05/09 02:30:53 finish InstallSnapshotRpcHandler.....rf.log:[{0 20516}]
2025/05/09 02:30:53 receve snapShot
2025/05/09 02:30:53 2 begin readSnapShot....
2025/05/09 02:30:53 readSnapShot decode error....
```



msg 内容：{true {Append 4 x 4 4 y 8 402244904228922384} 42}

这个时候 2 的LastIncludedIndex 明显变大了

```
2025/05/09 02:30:53 AppendEntries,args.leaderId:3,args.epoch:1，args.preIndex: 44, args.preTerm: 1, rf.me：2,rf.Term：1,rf.voteFor：3,rf.state：0, rf.log：[{0 20516}], args.log:[{1 {Append 4 x 4 5 y 9 402244904228922384}} {1 {Append 2 x 2 1 y 7 1228272264126474379}} {1 {Get 0  10 1123720622585023599}} {1 {Get 2  8 1228272264126474379}} {1 {Get 3  9 1098643115262335781}} {1 {Append 0 x 0 7 y 11 1123720622585023599}} {1 {Append 1 x 1 4 y 9 2016228440974134146}} {1 {Append 4 x 4 6 y 10 402244904228922384}} {1 {Append 3 x 3 7 y 10 1098643115262335781}} {1 {Append 0 x 0 8 y 12 1123720622585023599}} {1 {Get 2  9 1228272264126474379}} {1 {Append 1 x 1 5 y 10 2016228440974134146}}]
2025/05/09 02:30:53 2 getAbsLogIndex rf.LastIncludedIndex:44
```



