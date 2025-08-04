





> 面对大量用户访问，高并发请求，海量数据，可以使用高性能的服务器，大型的存储资源，高效率的编程语言， 但是当单机容量达到极限时，需要考虑的就是进行业务拆分和分布式部署，来避免整个系统的三高（高并发、高性能、高可用）问题。
>
> 当进行应用拆分，部署到不同的机器上，实现大规模的分布式系统后，紧接着要解决的就是对流量进行合理的分发，对压力过大时系统整体可用的保证。于是负载均衡、动态路由、心跳检查、熔断、限流等策略应运而生。

## 1. 负载均衡

>  负载均衡，英文名称为Load Balance，其含义就是指将负载（工作任务）进行平衡、分摊到多个操作单元上进行运行，例如FTP服务器、Web服务器、企业核心应用服务器和其它主要任务服务器等，从而协同完成工作任务。

任务是平衡、分摊网络请求，既然涉及到网络，而网络又是分层次的，所以针对[OSI](https://baike.baidu.com/item/网络七层协议) 中涉及的层级，都有各自的负载均衡策略。

#### 1.1 网络层面负载均衡分类

##### 1. 二层负载均衡（MAC）

根据OSI模型分的二层负载，一般是用虚拟mac地址方式，外部对虚拟MAC地址请求，负载均衡接收后分配后端实际的MAC地址响应.

##### 2. 三层负载均衡（IP）

一般采用虚拟IP地址方式，外部对虚拟的ip地址请求，负载均衡接收后分配后端实际的IP地址响应. (即一个ip对一个ip的转发, 端口全放开)

##### 3. 四层负载均衡（TCP）

在三次负载均衡的基础上，即从第四层"传输层"开始, 使用"ip+port"接收请求，再转发到对应的机器。这里不需要理解具体的包内容。对于这层负载均衡实现较好的软件包括：

- F5：硬件负载均衡器，功能很好，但是成本很高。
- lvs：重量级的四层负载软件
- nginx：轻量级的四层负载软件，带缓存功能，正则表达式较灵活
- haproxy：模拟四层转发，较灵活

##### 4. 七层负载均衡（应用层）

通过虚拟的URL或主机名接收请求，然后再分配到真实的服务器。

上面提到的是网络部署相关的。我们在实际的开发中，相对关心的其实框架提供的负载均衡策略，也就是在微服务中客户端在进行服务发现之后的一个重要组件，用于在服务返回健康节点列表后的平衡选择策略。

#### 1.2 常见的负载均衡算法

##### 轮询法（Round Robin）

> 将请求按顺序轮流地分配到后端服务器上，它均衡地对待后端的每一台服务器，而不关心服务器实际的连接数和当前的系统负载。
> 轮询法又有两个常见的改进方法
>
> 1. 加权轮询法： 根据服务器的不同处理能力，给每个服务器分配不同的权值，使其能够接受相应权值数的服务请求。
> 2. 平滑加权轮询法： 类似于加权轮询法，但生成的序列更加均匀。
>
> 
> 举个例子
> ![enter image description here](http://km.woa.com/files/photos/pictures/202010/1603954657_51_w808_h506.png =350x)

##### 随机法（Random）

> 通过系统的随机算法，根据后端服务器的列表大小值来随机选取其中的一台服务器进行访问。

##### 加权随机法（Weight Random）

> 与加权轮询法一样，加权随机法也根据后端机器的配置，系统的负载分配不同的权重。不同的是，它是按照权重随机请求后端服务器，而非顺序。

##### 最小连接法（Least Connections）

> 根据后端服务器当前的连接情况，动态地选取其中当前积压连接数最少的一台服务器来处理当前的请求，比较适用于长连接。

##### 源地址哈希法（Hash）

> 通过哈希运算获取下标在进行路由，如果机器数量不固定，还需要采用一致性哈希算法。

#### 1.3 tRPC-go中的实现

服务调用链路图

![f75a4530f1b549a0b6ae6de57abd5abb.png](/Users/giffinhao/Downloads/笔记/pic/f75a4530f1b549a0b6ae6de57abd5abb.png.webp)



可以了解到， 负载均衡是在主调方，也就是客户端的策略，于是可以想到，策略的具体代码是在`trpc-go-cmdline`生成的client、以及对应pb协议的代码中。

![image-20240821160228449](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821160228449.png)

![image-20240821160250309](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821160250309.png)

![1e0b60400c094cb2b13187c731ee98cc.png](/Users/giffinhao/Downloads/笔记/pic/1e0b60400c094cb2b13187c731ee98cc.png.webp)

![image-20240821160408219](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821160408219.png)

![image-20240821160454801](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821160454801.png)



![image-20240821160514006](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821160514006.png)

为了说明加权轮询算法的运行过程，我们以三个服务器节点为例，它们的权重分别是3、2、1。我们来详细演示这个算法是如何在这些服务器之间分配请求的。

### 服务器设置
- **Server A**: 权重 = 3
- **Server B**: 权重 = 2
- **Server C**: 权重 = 1

每个服务器的 `effectiveWeight` 初始值就是它的权重，而 `currentWeight` 初始为 0。

### 初始状态
- **Server A**: `currentWeight = 0`, `effectiveWeight = 3`
- **Server B**: `currentWeight = 0`, `effectiveWeight = 2`
- **Server C**: `currentWeight = 0`, `effectiveWeight = 1`

### 第一次选择
- **Server A**: `currentWeight = 0 + 3 = 3`
- **Server B**: `currentWeight = 0 + 2 = 2`
- **Server C**: `currentWeight = 0 + 1 = 1`

在这一轮中，`Server A` 的 `currentWeight` 最大，因此选择 `Server A` 来处理请求。

然后，更新选中的 `Server A` 的 `currentWeight`：
- `Server A`: `currentWeight = 3 - (3 + 2 + 1) = 3 - 6 = -3`

### 第二次选择
- **Server A**: `currentWeight = -3 + 3 = 0`
- **Server B**: `currentWeight = 2 + 2 = 4`
- **Server C**: `currentWeight = 1 + 1 = 2`

在这一轮中，`Server B` 的 `currentWeight` 最大，因此选择 `Server B` 来处理请求。

然后，更新选中的 `Server B` 的 `currentWeight`：
- `Server B`: `currentWeight = 4 - 6 = -2`

### 第三次选择
- **Server A**: `currentWeight = 0 + 3 = 3`
- **Server B**: `currentWeight = -2 + 2 = 0`
- **Server C**: `currentWeight = 2 + 1 = 3`

在这一轮中，`Server A` 和 `Server C` 的 `currentWeight` 都是 3，但由于 `Server A` 排在前面，它将被选择来处理请求。

然后，更新选中的 `Server A` 的 `currentWeight`：
- `Server A`: `currentWeight = 3 - 6 = -3`

### 第四次选择
- **Server A**: `currentWeight = -3 + 3 = 0`
- **Server B**: `currentWeight = 0 + 2 = 2`
- **Server C**: `currentWeight = 3 + 1 = 4`

在这一轮中，`Server C` 的 `currentWeight` 最大，因此选择 `Server C` 来处理请求。

然后，更新选中的 `Server C` 的 `currentWeight`：
- `Server C`: `currentWeight = 4 - 6 = -2`

### 第五次选择
- **Server A**: `currentWeight = 0 + 3 = 3`
- **Server B**: `currentWeight = 2 + 2 = 4`
- **Server C**: `currentWeight = -2 + 1 = -1`

在这一轮中，`Server B` 的 `currentWeight` 最大，因此选择 `Server B` 来处理请求。

然后，更新选中的 `Server B` 的 `currentWeight`：
- `Server B`: `currentWeight = 4 - 6 = -2`

### 第六次选择
- **Server A**: `currentWeight = 3 + 3 = 6`
- **Server B**: `currentWeight = -2 + 2 = 0`
- **Server C**: `currentWeight = -1 + 1 = 0`

在这一轮中，`Server A` 的 `currentWeight` 最大，因此选择 `Server A` 来处理请求。

然后，更新选中的 `Server A` 的 `currentWeight`：
- `Server A`: `currentWeight = 6 - 6 = 0`

### 循环开始
到此，第六次选择结束，负载均衡器会继续进行下一轮选择，从头开始计算 `currentWeight`，这个过程将不断循环。

### 总结

在这个过程中：
- **Server A** 被选中的次数最多（共 3 次），因为它的权重最高。
- **Server B** 被选中的次数居中（共 2 次），因为它的权重中等。
- **Server C** 被选中的次数最少（共 1 次），因为它的权重最低。

这个算法有效地根据权重分配了请求，同时确保了每个服务器都有机会处理请求，达到负载均衡的效果。



![image-20240821161013542](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821161013542.png)

![image-20240821161118765](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821161118765.png)

#### .1 常见熔断策略

> 熔断的策略就是衡量资源是否处于稳定的状态的方式。业内使用较多的熔断中间件包括 `Hystrix`，`Sentinel` 都涵盖了以下几种策略

##### 1. 平均响应时间

比如：当 1s 内持续进入 5 个请求，对应时刻的平均响应时间（秒级）均超过阈值（count，以 ms 为单位），那么在接下的时间窗口之内，对这个方法的调用都会自动地熔断。

##### 2. 异常比例

比如: 当资源的每秒请求量 >= 5，并且每秒异常总数占通过量的比值超过阈值之后，资源进入熔断状态，即在接下的时间窗口之内，对这个方法的调用都会自动地返回。

##### 3. 异常数

比如： 当资源近 1 分钟的异常数目超过阈值之后会进行熔断。



![image-20240821180022196](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821180022196.png)







#### 常见的限流策略

##### 1. 令牌桶算法

令牌桶算法的原理是系统以恒定的速率产生令牌，然后把令牌放到令牌桶中，令牌桶有一个容量，当令牌桶满了的时候，再向其中放令牌，那么多余的令牌会被丢弃；当想要处理一个请求的时候，需要从令牌桶中取出一个令牌，如果此时令牌桶中没有令牌，那么则拒绝该请求。

生成令牌的速度是恒定的，而请求去拿令牌是没有速度限制的。这意味，面对瞬时大流量，该算法可以在短时间内请求拿到大量令牌，而且拿令牌的过程并不是消耗很大的事情。



##### 2. 漏桶算法

把请求比作是水，水来了都先放进桶里，并以限定的速度出水，当水来得过猛而出水不够快时就会导致水直接溢出，即拒绝服务。

漏桶的出水速度是恒定的，那么意味着如果瞬时大流量的话，将有大部分请求被丢弃掉（也就是所谓的溢出）。



##### 3. 滑动窗口



![image-20240821180523066](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821180523066.png)





![image-20240821180551032](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821180551032.png)

![image-20240821180609464](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821180609464.png)

![image-20240821180639976](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821180639976.png)

![image-20240821180719609](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821180719609.png)

![image-20240821180743224](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821180743224.png)







![image-20240821180451384](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821180451384.png)

设定的单位时间就是一个窗口，窗口可以分割多个更小的时间单元，随着时间的推移，窗口会向右移动。窗口内的请求会被记录，超过既定的阀值，就开始拒绝服务。

![image-20240821180326472](/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821180326472.png)

<img src="/Users/giffinhao/Library/Application Support/typora-user-images/image-20240821180831914.png" alt="image-20240821180831914" style="zoom:50%;" />