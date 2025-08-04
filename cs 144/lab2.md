



> 将一个 `WrappingInt32` 类型的 **32位序列号** `n` 转换成 **最接近给定 checkpoint 的 64位绝对序列号**。

------

## 🌟 背景知识

- `WrappingInt32`：TCP 中的 32-bit 序列号，它在 `2^32` 后会绕回（wrap around）。
- `unwrap(n, isn, checkpoint)` 的任务：从所有满足 `(x % 2^32) == n - isn` 的 `x` 中，选出距离 `checkpoint` 最近的那一个。

------

## ✅ 一步步解释

```
cpp


复制编辑
uint32_t offset = n - isn;
```

- `n` 和 `isn` 都是 `WrappingInt32` 类型，相减后会得出二者之间的 32 位距离（可能 wrap 过一次）。
- `offset` 表示 `n` 在 ISN 之后走了多少步。

例如：
 `n = WrappingInt32{100}`
 `isn = WrappingInt32{50}`
 `offset = 50`，表示 n 是从 isn 起的第 50 个。

------

```
cpp


复制编辑
uint64_t base = (checkpoint & 0xFFFFFFFF00000000ull);
```

- 把 `checkpoint` 向下对齐到 2^32 的整数倍。
- 等价于：`checkpoint / 2^32 * 2^32`。
- 作用是：构造一个候选 `x = base + offset`，使它和 `checkpoint` 属于同一个“2^32 块”。

例如：
 如果 `checkpoint = 0x0000000400000000 (2^34)`，
 则 `base = 0x0000000400000000`，保持在 checkpoint 相同的 2^32 区间起点。

------

```
cpp


复制编辑
uint64_t candidate = base + offset;
```

- 使用 base + offset 得出一个候选解。
- 这个 `candidate` 是我们首选的解（初始假设它是最接近 checkpoint 的那个）。

------

```
cpp


复制编辑
uint64_t best = candidate;
```

- 初始最优解设为 `candidate`。

------

```
cpp


复制编辑
if (candidate >= (1ul << 32)) {
    uint64_t cand2 = candidate - (1ul << 32);
    if (abs(static_cast<int64_t>(cand2 - checkpoint)) < abs(static_cast<int64_t>(best - checkpoint))) {
        best = cand2;
    }
}
```

### 👈 思路

- `candidate - 2^32` 也是可能的解，因为 wrap-around 意味着相同的 32 位值可能出现在多个 64 位位置。
- 只要结果仍为非负（即 `>= 0`），就合法。

### ✅ 作用

- 比较 `cand2` 和 `candidate` 谁更接近 `checkpoint`。
- 注意：必须转成 `int64_t` 来正确计算 `abs(...)`，防止无符号下出现负值变大整数。

------

```
cpp


复制编辑
uint64_t cand3 = candidate + (1ul << 32);
if (abs(static_cast<int64_t>(cand3 - checkpoint)) < abs(static_cast<int64_t>(best - checkpoint))) {
    best = cand3;
}
```

- 同理，`candidate + 2^32` 也可能是另一个 wrap-around 候选。
- 比较它和当前 best，若更接近 `checkpoint`，就更新。

------

```
cpp


复制编辑
return best;
```

- 返回最终距离 checkpoint 最近的 **绝对序列号**。

------

## 🧠 举个例子说明

假设：

- `isn = WrappingInt32{100}`
- `n = WrappingInt32{150}`
- 所以 offset = 150 - 100 = 50
- `checkpoint = 0x00000001FFFFFFFE (8589934590)`

步骤如下：

1. `base = checkpoint & 0xFFFFFFFF00000000 = 0x0000000100000000 = 4294967296`
2. `candidate = base + offset = 4294967296 + 50 = 4294967346`
3. 候选：
   - `cand1 = 4294967346`
   - `cand2 = 4294967346 - 2^32 = 50`
   - `cand3 = 4294967346 + 2^32 = 8589934642`
4. 计算与 checkpoint 的距离：
   - `abs(4294967346 - 8589934590) = 4294967244`
   - `abs(50 - 8589934590) = 很大`
   - `abs(8589934642 - 8589934590) = 52`

所以返回 `8589934642`（cand3）作为最优解。

------

## ✅ 总结

这个 `unwrap` 做了三件事：

1. **构造 offset**：计算目标值相对于 isn 的偏移。
2. **构造候选**：以 checkpoint 为参考，生成与之靠近的几个候选值。
3. **返回最优**：选出离 checkpoint 最近的那个。





## 1. std::move 的基本概念

### 定义和声明

```
// 在 <utility> 头文件中定义
template<typename T>
typename std::remove_reference<T>::type&& move(T&& t) noexcept;
```

### 核心作用

std::move 是一个类型转换函数，将左值转换为右值引用，从而启用移动语义。

## 2. std::move 的实现原理

```
template<typename T>
typename std::remove_reference<T>::type&& move(T&& t) noexcept {
    return static_cast<typename std::remove_reference<T>::type&&>(t);
}
```

### 实现解析

#### 1. 模板参数 T&&

```
template<typename T>
typename std::remove_reference<T>::type&& move(T&& t)
```

- T&& 是通用引用（Universal Reference）

- 可以接受任何类型的参数（左值或右值）

- 通过引用折叠规则确定最终类型



#### 2. std::remove_reference

```
// std::remove_reference 的实现
template<typename T>
struct remove_reference {
    using type = T;
};

template<typename T>
struct remove_reference<T&> {
    using type = T;
};

template<typename T>
struct remove_reference<T&&> {
    using type = T;
};
```





```
#include <iostream>
#include <utility>
#include <string>

class MoveableClass {
private:
    std::string data;
    
public:
    // 构造函数
    MoveableClass(const std::string& str) : data(str) {
        std::cout << "构造函数: " << str << std::endl;
    }
    
    // 拷贝构造函数
    MoveableClass(const MoveableClass& other) : data(other.data) {
        std::cout << "拷贝构造函数" << std::endl;
    }
    
    // 移动构造函数
    MoveableClass(MoveableClass&& other) noexcept : data(std::move(other.data)) {
        std::cout << "移动构造函数" << std::endl;
    }
    
    // 拷贝赋值运算符
    MoveableClass& operator=(const MoveableClass& other) {
        if (this != &other) {
            data = other.data;
            std::cout << "拷贝赋值运算符" << std::endl;
        }
        return *this;
    }
    
    // 移动赋值运算符
    MoveableClass& operator=(MoveableClass&& other) noexcept {
        if (this != &other) {
            data = std::move(other.data);
            std::cout << "移动赋值运算符" << std::endl;
        }
        return *this;
    }
    
    // 析构函数
    ~MoveableClass() {
        std::cout << "析构函数" << std::endl;
    }
};

int main() {
    std::cout << "=== 创建对象 ===" << std::endl;
    MoveableClass obj1("Hello");
    
    std::cout << "\n=== 拷贝构造 ===" << std::endl;
    MoveableClass obj2 = obj1;  // 拷贝构造
    
    std::cout << "\n=== 移动构造 ===" << std::endl;
    MoveableClass obj3 = std::move(obj1);  // 移动构造
    
    std::cout << "\n=== 移动赋值 ===" << std::endl;
    MoveableClass obj4("World");
    obj4 = std::move(obj2);  // 移动赋值
    
    return 0;
}
```





```
void TCPReceiver::segment_received(const TCPSegment &seg) {
    // 1. 取出 sn 查看是否超出窗口
    // 2。 使用 sn 转换为绝对序列号
    // 3. 将 data 存入 resemble 
    if (!_isn.has_value()) {
        if (seg.header().syn) {
            _isn = WrappingInt32(seg.header().seqno);
        } else {
            return;
        }
    }
    uint64_t absno = unwrap(seg.header().seqno, _isn.value(), _reassembler.stream_out().bytes_written());
    _reassembler.push_substring(seg.payload().copy(), (seg.header().syn==true)?0:(absno-1) , seg.header().fin);
    _ackno = wrap(_reassembler.stream_out().bytes_written() + (_reassembler.stream_out().input_ended()? 2 : 1), _isn.value());
}
```

**ack 是期望收到的下一个序列号**， 当包括 fin 时为+2

## 错误原因

### 1. FIN 标志位没有正确计入序列号

TCP 协议规定：

- SYN 标志占用 1 个序列号

- FIN 标志占用 1 个序列号

这段注释详细说明了 TCP 接收窗口（window size）的含义和计算方式。我们来逐句拆解讲解：

------

## 🧾 原文逐句解析：

```cpp
//! \brief The window size that should be sent to the peer
```

**这表示：**
 这个函数返回的“窗口大小（window size）”应该发送给对端（peer，也就是 TCP 的发送方），以告知当前接收缓冲区的剩余空间。

------

```cpp
//! Operationally: the capacity minus the number of bytes that the
//! TCPReceiver is holding in its byte stream (those that have been
//! reassembled, but not consumed).
```

**操作层面（Operationally）的解释：**

- 接收窗口 = 总容量（`_capacity`）
- 减去「已经被重组（reassembled）但尚未被消费（read）的字节数量」

这些字节目前正占用缓冲区空间。

🌰 举个例子：

- `capacity = 1000`
- 已组装好但还没被应用读取的字节数为 400
- 那么 window size = `1000 - 400 = 600`，发送方最多再发 600 字节。

------

```cpp
//! Formally: the difference between (a) the sequence number of
//! the first byte that falls after the window (and will not be
//! accepted by the receiver) and (b) the sequence number of the
//! beginning of the window (the ackno).
```

**形式上的（Formally）解释：**

窗口的大小可以通过 **TCP 序列号（sequence number）** 来衡量：

- `(a)`: 表示窗口“末尾”之后的第一个序号。再大的序号将会被拒收。
- `(b)`: 表示窗口的“开始”位置，也就是接收方的 `ackno`（期望接收的下一个字节序号）

所以：

```cpp
window size = (window_end_seqno) - (ackno)
```

这就是标准的 TCP 流量控制的“滑动窗口”模型。

------

## 🔁 两种描述的关系：

- **Operationally（操作层面）**：更偏工程实现，用 `capacity - buffer_size()` 计算
- **Formally（序列号层面）**：更偏协议定义，用 `(window_end_seqno - ackno)` 计算

但**两者等价**，因为：

```cpp
bytes in buffer = sequence numbers received - bytes consumed
```

------

## 📌 总结一句话：

> TCPReceiver 的 `window_size()` 表示「还能接收多少数据」，它是当前接收缓冲区的剩余容量。这个值会被发送到 TCP 的 ACK 报文中，让发送方知道“最多可以再发多少数据”。从协议上讲，它是窗口末尾与 `ackno` 之间的距离。
