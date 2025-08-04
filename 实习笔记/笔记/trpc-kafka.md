trpc-go kafka使用及原理浅析 

### kafka介绍

在本文开始之前，我们首先要思考如下几个问题？



1. 什么是kafka？



首先介绍一下什么是kafka，以下是官网的定义：Everything you need to know about Kafka in 10 minutes :[Apache Kafka](https://kafka.apache.org/intro) 简单来说，kafka是一个开源的、分布式的消息系统，使用Scala编写，它以可水平扩展和高吞吐率而被广泛使用。可以把kafka简单的理解成为一个队列，一头往队列里插元素（producer），另一头取元素（consumer） 2. kafka能干嘛？ 在实际业务开发过程中，在系统设计之初的时候，可能考虑的不是很全面，在后续更新迭代的过程中，发现要增加一些新的功能，这个时候可能就要增加新的接口，增加新的数据库表等，而且随着业务的不断扩展，需要加的东西可能会越来越多，很多时候不同的业务之间用到的数据是一样的，如果不同业务之间都这样干的时候，会造成系统的膨胀和强耦合性，此时就可以使用kafka来解决这个问题了，不同的业务方可以订阅相同的或者不同的topic，消费数据，至于拿到数据之后干嘛，灵活性就很大了，就好比定义了一个接口，可以有很多种实现，实现了高可扩展性和接触了相互之间的耦合。除了上述场景之外，kafka还提功如下能力：



- 顺序保证



在大多使用场景下，数据处理的顺序都很重要。大部分消息队列本来就是排序的，并且能保证数据会按照特定的顺序来处理。Kafka保证一个Partition内的消息的有序性。



- 缓冲



在任何重要的系统中，都会有需要不同的处理时间的元素。例如，加载一张图片比应用过滤器花费更少的时间。消息队列通过一个缓冲层来帮助任务最高效率的执行———写入队列的处理会尽可能的快速。该缓冲有助于控制和优化数据流经过系统的速



- 灵活性和削峰



在访问量剧增的情况下，应用仍然需要继续发挥作用，但是这样的突发流量并不常见；如果为以能处理这类峰值访问为标准来投入资源随时待命无疑是巨大的浪费。使用消息队列能够使关键组件顶住突发的访问压力，而不会因为突发的超负荷的请求而完全崩溃。



- 高性能



Kafka虽然是基于磁盘做的数据存储，但却具有高性能、高吞吐、低延时的特点，其吞吐量动辄几万、几十上百万。一般印象中，磁盘的速度是比较慢的，但是由于kafka是顺序读写的，而磁盘的顺序读写性能却很高，一般而言要高出磁盘随机读写三个数量级，除了顺序读写，Page Cache、Zero Copy等技术也是kafka快的原因，具体可参考[Kafka为什么吞吐量大、速度快](https://zhuanlan.zhihu.com/p/120967989) 3. 什么时候使用kafka？ 虽然前面介绍了kafka的各种能力，那为什么要选择使用kafka，或者说在什么情况下选用kafka呢？一般情况下，当你的程序中有如下诉求的时候，kafka是个比较好的选择。



- 异步通信：将业务中属于非核心或不重要的流程部分，使用消息异步通知的方式发给目标系统，这样主业务流程无需同步等待其他系统的处理结果，从而达到系统快速响应的目的。
- 消息系统：Kafka被当作传统消息中间件的替代品。与大多数消息系统相比，Kafka具有更好的吞吐量，内置的分区，多副本和容错性，这使其成为大规模消息处理应用程序的良好解决方案。
- 错峰流控与流量削峰：在电子商务系统或大型网站中，上下游系统处理能力存在差异，处理能力高的上游系统的突发流量可能会对处理能力低的某些下游系统造成冲击，需要提高系统的可用性的同时降低系统实现的复杂性。电商大促销等流量洪流突然来袭时，可以通过队列服务堆积缓存订单等信息，在下游系统有能力处理消息的时候再处理，避免下游订阅系统因突发流量崩溃。消息队列提供亿级消息堆积能力，3天的默认保留时长，消息消费系统可以错峰进行消息处理。另外，在商品秒杀、抢购等流量短时间内暴增场景中，为了防止后端应用被压垮，可在前后端系统间使用Kafka消息队列传递请求。
- 指标：用Kafka采集应用程序和服务器健康相关的指标，如CPU占用率、IO、内存、连接数、TPS、QPS等，然后将指标信息进行处理，从而构建一个具有监控仪表盘、曲线图等可视化监控系统。例如，司内天机阁二代所基于的OpenTelemetry就是这样做的。



### trpc-go使用kafka

下面进入本文的主题，如何在trpc-go中使用kafka？ trpc-go官方给我们提供了trpc-database/kafka框架，基于Shopify/sarama框架二次封装，可以非常方便的配合 trpc 使用，本文主要介绍我们作为消费方（consumer）如何接入以及使用kafka。 在trpc-database/kafka中，consumer是作为一个service存在的，service是什么？service是trpc的一个服务，在mian函数的trpc.NewServer()里面会对配置文件进行处理，然后使用s.Serve()方法启动service，也就是说我们的每一个kafka消费者其实单独起了个goroutine去消费消息。所以在使用trpc-go中使用kafka，首先需要在trpc_go.yaml文件类似普通的trpc服务那样配置service信息，如下：

```javascript
server:
  service:
   - name: trpc.app.server.consumer1 # service 名，保证唯一即可，任意取
    address: ip1:port1,ip2:port2?topics=topic1,topic2&amp;amp;group=xxx&amp;amp;version=x.x.x.x # 注意书写格式，不然框架会解析失败
    protocol: kafka # 协议填kafka
    timeout: 1000
```

可以配置多个consumer service，也可以在单个service下配置多个topic。

在配置文件中配置好service信息后，我们有两种方式来消费service，一种是实现kafka/service_desc的KafkaConsumer接口的Handle方法。

```javascript
// git.code.oa.com/trpc-go/trpc-database/kafka@v0.1.9/service_desc.go
// 注意在0.1.9版本之后，实现的接口以及Handle函数的签名发生了变化，推荐使用0.1.9之后的版本
// KafkaConsumer 消费者接口
type KafkaConsumer interface {
  Handle(ctx context.Context, msg *sarama.ConsumerMessage) error
}
```

实现KafkaConsumer接口

```javascript
// Consumer 回调对象
type Consumer1 struct {
}

// Handle 回调方法
// 自定义消费逻辑，do anything what you want
func (Consumer1) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
  log.Infof("get kafka message: %+v", msg)
  // ...  消费业务逻辑
  // 返回nil才会确认消费成功
  return nil
}
```

注册consumer service

```javascript
// main.go
func main() {
  s := trpc.NewServer()
 // trpc.app.server.consumer1 service 名与yaml配置文件中保持一致
  kafka.RegisterKafkaConsumerService(s.Service("trpc.app.server.consumer1"),&amp;amp;Consumer1{})
 
 // 可以注册多个service，
 // kafka.RegisterKafkaConsumerService(s.Service("trpc.app.server.consumer2"),&amp;amp;Consumer{})
 
 // 启动服务
  err := s.Serve()
  if err != nil {
    panic(err)
  }
}
```

另一种方法是直接注册handle函数，而不是通过实现接口。

```javascript
func handleXXX(ctx context.Context, msg *sarama.ConsumerMessage) error {
  log.Info(msg)
  // ...  消费业务逻辑
  // 返回nil才会确认消费成功
  return nil
}
```

注册handle

```javascript
// main.go
func main() {
  s := trpc.NewServer()
 // trpc.app.server.consumer1 service 名与yaml配置文件中保持一致
  kafka.RegisterKafkaHandlerService(s.Service("trpc.app.server.consumer1"), handleXXX)
 
 // 可以注册多个handle
 // kafka.RegisterKafkaHandlerService(s.Service("trpc.app.server.consumer2"), handleXXX)

 // 启动服务
  err := s.Serve()
  if err != nil {
    panic(err)
  }
}
```

上述的方式一次只能消费一个topic，如果想一次消费多个topic，在main函数中注册批量消费的handler即可 hanler函数实现

```javascript
// Handle 回调方法
func Handle(ctx context.Context, msgArray []*sarama.ConsumerMessage) error {
  // 返回nil才会确认消费成功
  log.Infof("len(msgArray) = %d", len(msgArray))
  for _, v := range msgArray {
    log.Infof("[consume][topic]%v\t[partition]%v\t[offset]%v\t[key]%v\t[value]%v\n",
      v.Topic, v.Partition, v.Offset, string(v.Key), string(v.Value))
  // ... 消费业务逻辑
  }
  // 返回nil才会确认消费成功
  return nil
}
```

并在main函数中注册

```javascript
func main() {

  s := trpc.NewServer()
   // trpc.app.server.consumer1 service 名与yaml配置文件中保持一致
  kafka.RegisterBatchHandlerService(s.Service("trpc.app.server.consumer1"), Handle)
  err := s.Serve()
  if err != nil {
    panic(err)
  }
}
```

### 原理浅析

前面介绍了如何在trpc-go框架中使用kafka，可以看到，接入使用还是比较简单的，前面说到consumer是作为trpc-go中一个service而存在的，但又不同于普通的trpc service提供接口供外部使用，实现KafkaConsumer接口的Handle方法类似于trpc服务的接口，但是显然，这个“接口”不是给外部调用的，trpc-database/kafka框架是如何做到当生产者生产消息的时候，consumer能够持续不断的消费来自producer的消息呢？换句话说我们实现的Handle方法在什么时候会被调用呢？ 回顾一下trpc-go框架启动服务的流程，大致如下：

![image-20210718143502635](/Users/giffinhao/Downloads/笔记/pic/image-20210718143502635.png) 第6步才是真正启动一个service的步骤，我们进入源码里看看

```javascript
// Serve 启动服务
func (s *service) Serve() (err error) {

  pid := os.Getpid()

  // 确保正常监听之后才能启动服务注册
  if err = s.opts.Transport.ListenAndServe(s.ctx, s.opts.ServeOptions...); err != nil {
    log.Errorf("process:%d service:%s ListenAndServe fail:%v", pid, s.opts.ServiceName, err)
    return err
  }

  if s.opts.Registry != nil {
    err = s.opts.Registry.Register(s.opts.ServiceName, registry.WithAddress(s.opts.Address))
    if err != nil {
      // 有注册失败，关闭service，并返回给上层错误
      log.Errorf("process:%d, service:%s register fail:%v", pid, s.opts.ServiceName, err)
      return err
    }
  }
  // ....
}
```

在启动服务注册之前会对服务进行监听，这个地方很关键，当我们在main.go文件中引入trpc-database/kafka依赖的时候，会执行kafka包的init函数，把kafka transport注册到trpc-go的transport工厂里面，回到服务监听，看看kafka的ListenAndServe的实现:

```javascript
// ListenAndServe 启动监听，如果监听失败则返回错误
func (s *ServerTransport) ListenAndServe(ctx context.Context, opts ...transport.ListenServeOption) (err error) {

  kafkalsopts := &amp;amp;transport.ListenServeOptions{}
  for _, opt := range opts {
    opt(kafkalsopts)
  }

  kafkaUserConfig, err := parseAddress(kafkalsopts.Address)
  if err != nil {
    return err
  }

  config := sarama.NewConfig()
  config.Version = kafkaUserConfig.version
  config.Consumer.Group.Rebalance.Strategy = kafkaUserConfig.strategy
  config.Consumer.Offsets.Initial = kafkaUserConfig.initial
  config.ClientID = kafkaUserConfig.group

  config.Consumer.Fetch.Default = int32(kafkaUserConfig.fetchDefault)
  config.Consumer.Fetch.Max = int32(kafkaUserConfig.fetchMax)
  config.Consumer.MaxWaitTime = kafkaUserConfig.maxWaitTime

  config.Metadata.Full = false                 //禁止拉取所有元数据
  config.Metadata.Retry.Max = 1                 //元数据更新重次次数
  config.Metadata.Retry.Backoff = time.Second          //元数据更新等待时间
  config.Consumer.Offsets.AutoCommit.Interval = 3 * time.Second //定时多久一次提交消费进度
  kafkaUserConfig.scramClient.config(config)

  // 连接broker，失败会返回错误
  consumerGroup, err := sarama.NewConsumerGroup(kafkaUserConfig.brokers, kafkaUserConfig.group, config)
  if err != nil {
    return err
  }
  go func() {
    for {
      if consumerGroup != nil {
        if kafkaUserConfig.batchConsumeCount > 0 {
          consumer := &amp;amp;batchConsumer{opts: kafkalsopts, ctx: ctx,
            maxNum: kafkaUserConfig.batchConsumeCount, flushInterval: kafkaUserConfig.batchFlush}
          err = consumerGroup.Consume(ctx, kafkaUserConfig.topics, consumer)
        } else {
          err = consumerGroup.Consume(ctx, kafkaUserConfig.topics, &amp;amp;consumer{opts: kafkalsopts, ctx: ctx})
        }
      }
      select {
      case <-ctx.Done():
        log.ErrorContextf(ctx, "kafka server transport: context done:%v, close", ctx.Err())
        return
      default:
      }

      time.Sleep(time.Second * 3)
      if err == nil {
        continue
      }

      log.ErrorContextf(ctx, "kafka server transport: consume fail:%v, reconnect", err)
      if consumerGroup != nil {
        consumerGroup.Close()
        consumerGroup = nil
      }

      //重新连接broker，失败会返回错误
      consumerGroup, err = sarama.NewConsumerGroup(kafkaUserConfig.brokers, kafkaUserConfig.group, config)
      if err != nil {
        log.ErrorContextf(ctx, "kafka server transport: consume reconnect fail:%v", err)
      }
    }
  }()

  return nil
}
```

在初次启动consumer service的时候，会在后台启动一个协程死循环不断的轮询去消费来自producer的消息，并且为了避免过度频繁拉取kafka消息，会在每次拉取后等3s再去拉取。到这我们拿到了消费的消息，但是怎么找到我们实现的Handle方法并去调用呢？继续往下看，消费kafka消息逻辑是在43行调用Consume方法，那么这个方法具体干了什么事呢？我们在这个方法中传入了consumer对象，该结构体实现了ConsumerGroupHandler接口，该接口包含Setup、Cleanup、ConsumeClaim三个方法，前两个方法分别是ConsumeClaim的前置和后置操作，重点在于这个ConsumeClaim方法，我们来看一下这个方法的代码：

```javascript
// git.code.oa.com/trpc-go/trpc-database/kafka@v0.1.9/transport.go
// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (s *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
  for message := range claim.Messages() {

    select {
    case <-s.ctx.Done():
      return errors.New("consumer service close")
    default:
    }

    //生成新的空的通用消息结构数据，并保存到ctx里面
    ctx, msg := codec.WithNewMessage(context.Background())

    //填充被调方，自己
    msg.WithCompressType(codec.CompressTypeNoop) //不解压缩
    msg.WithServerReqHead(message)

    _, err := s.opts.Handler.Handle(ctx, nil)
    if err != nil || msg.ServerRspErr() != nil {
      msgInfo := ""
      if message != nil {
        msgInfo = fmt.Sprintf("%+v:%+v:%+v", message.Topic, message.Partition, message.Offset)
      }
      if s.maxRetry == 0 || s.retryNum < s.maxRetry {
        s.retryNum++
        log.ErrorContextf(ctx, "kafka consumer handle fail:%v %v, retry: %+v, msg: %+v",
          err, msg.ServerRspErr(), s.retryNum, msgInfo)
        return nil
      }
    }

    //确认消息消费成功
    s.retryNum = 0
    session.MarkMessage(message, "")
  }

  return nil
}
```

注意第19行，这里调用了Handler.Handle方法，其实调用的是service的Handle方法

```javascript
// git.code.oa.com/trpc-go/trpc-go@v0.6.6/server/service.go
// Handle server transport收到请求包后调用此函数
func (s *service) Handle(ctx context.Context, reqbuf []byte) (rspbuf []byte, err error) {
  // ....

  msg := codec.Message(ctx)
  reqbodybuf, err := s.decode(ctx, msg, reqbuf)
  if err != nil {
    return s.encode(ctx, msg, nil, err)
  }
  // 已经有错误了（通常是过载错误），在解析完包头拿到 RequestID 后立刻返回客户端。
  if err := msg.ServerRspErr(); err != nil {
    return s.encode(ctx, msg, nil, err)
  }
  rspbody, err := s.handle(ctx, msg, reqbodybuf)
 
 // ....
}
```

继续看15行的s.handle方法

```javascript
func (s *service) handle(ctx context.Context, msg codec.Msg, reqbodybuf []byte) (interface{}, error) {
  // ...
  handler, ok := s.handlers[msg.ServerRPCName()]
  // ...
  rspbody, err := handler(ctx, s.filterFunc(msg, reqbodybuf))
  // ...
}

//
// service Service实现
type service struct {
  ctx      context.Context  // service关闭
  cancel     context.CancelFunc // service关闭
  opts      *Options      // service选项
  handlers    map[string]Handler // rpcname => handler
  streamHandlers map[string]StreamHandler
}
```

可以看到从service的handlers工厂map中根据ServerRPCName取出对应的handler，那这个handler是怎么来的？或者说是在什么时候塞到s.handlers里面的呢？回顾一下，我们采用kafka.RegisterKafkaHandlerService或者kafka.RegisterKafkaConsumerService注册的时候都干了些啥呢？其实这两种方式没啥区别，只是传入的参数不一样，一个传入handle函数由代码的handler帮我们处理了，一个传入自己实现了handle方法的Consumer对象。

```javascript
// git.code.oa.com/trpc-go/trpc-database/kafka@v0.1.9/service_desc.go
// RegisterKafkaConsumerService 注册service
func RegisterKafkaConsumerService(s server.Service, svr KafkaConsumer) {
  s.Register(&amp;amp;KafkaConsumerServiceDesc, svr)
}

// KafkaConsumerServiceDesc 对应的descriptor
var KafkaConsumerServiceDesc = server.ServiceDesc{
  ServiceName: fmt.Sprintf("trpc.kafka.consumer.service"),
  HandlerType: ((*KafkaConsumer)(nil)),
  Methods: []server.Method{{
    Name: "/trpc.kafka.consumer.service/handle",
    Func: KafkaConsumerHandle,
  }},
}

// KafkaConsumerHandle consumer service handler wrapper
func KafkaConsumerHandle(svr interface{}, ctx context.Context, f server.FilterFunc) (rspbody interface{}, err error) {
  filters, err := f(nil)
  // ...

  handleFunc := func(ctx context.Context, reqbody interface{}, rspbody interface{}) error {
    msg := codec.Message(ctx)
    m, ok := msg.ServerReqHead().(*sarama.ConsumerMessage)
    if !ok {
      return errs.NewFrameError(errs.RetServerDecodeFail, "kafka consumer handler: message type invalid")
    }
    return svr.(KafkaConsumer).Handle(ctx, m)
  }

  err = filters.Handle(ctx, nil, nil, handleFunc)
  // ...
}

// git.code.oa.com/trpc-go/trpc-go@v0.6.6/server/service.go
// Register 把service业务实现接口注册到server里面
func (s *service) Register(serviceDesc interface{}, serviceImpl interface{}) error {
  // ....
  for _, method := range desc.Methods {
    h := method.Func
    s.handlers[method.Name] = func(ctx context.Context, f FilterFunc) (rsp interface{}, err error) {
      return h(serviceImpl, ctx, f)
    }
    // ....
}
```

这就回答了我们刚才的问题，我们在服务注册的时候，会把KafkaConsumerHandle函数塞到s.handlers里面，s.handlers[msg.ServerRPCName()]实际上就是KafkaConsumerHandle函数，这个函数会调用filters.Handle(ctx, nil, nil, handleFunc)来执行我们实现的消费逻辑，到这里整个链条就很清楚了。

<img src="/Users/giffinhao/Downloads/笔记/pic/image-20210718184210212.png" alt="image-20210718184210212" style="zoom:50%;" />

关键是第3步起一个后台死循环不间断的轮询去获取最新消息消费。

### 小结

trpc-go和trpc-database/kafka中有很多巧妙的设计，当我们使用一个框架的时候，使用起来很感觉非常简便快捷的时候，这个时候就要思考一下为什么如此好用，只需如此这般就能够满足我们的需求，通过不断模仿学习继而去学习思考其中的优秀设计，在自己写代码的时候，不管是一个微服务还是SDK，在设计的时候也应该考虑提供足够简便方式供调用者使用。





## service_desc.go

### KafkaConsumer 接口

```
go
复制代码
// KafkaConsumer 消费者接口
type KafkaConsumer interface {
	Handle(ctx context.Context, msg *sarama.ConsumerMessage) error
}
```

定义了一个 Kafka 消费者接口，包含一个 `Handle` 方法，用于处理单条 Kafka 消息。

### kafkaHandler 类型及其 Handle 方法

```
go
复制代码
type kafkaHandler func(ctx context.Context, msg *sarama.ConsumerMessage) error

// Handle 主处理函数
func (h kafkaHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	return h(ctx, msg)
}
```

定义了一个 `kafkaHandler` 类型，它是一个函数类型，包含 `Handle` 方法，用于处理单条 Kafka 消息。

### KafkaConsumerServiceDesc 描述符

```
go
复制代码
// KafkaConsumerServiceDesc 对应的 descriptor
var KafkaConsumerServiceDesc = server.ServiceDesc{
	ServiceName: "trpc.kafka.consumer.service",
	HandlerType: ((*KafkaConsumer)(nil)),
	Methods: []server.Method{{
		Name: "/trpc.kafka.consumer.service/handle",
		Func: KafkaConsumerHandle,
	}},
}
```

定义了一个 Kafka 消费者服务的描述符 `KafkaConsumerServiceDesc`，包含服务名称和方法。

### KafkaConsumerHandle 函数

```
go
复制代码
// KafkaConsumerHandle consumer service handler wrapper
func KafkaConsumerHandle(svr interface{}, ctx context.Context, f server.FilterFunc) (rspbody interface{}, err error) {
	filters, err := f(nil)
	if err != nil {
		return nil, err
	}
	handleFunc := func(ctx context.Context, reqbody interface{}, rspbody interface{}) error {
		msg := codec.Message(ctx)
		m, ok := msg.ServerReqHead().(*sarama.ConsumerMessage)
		if !ok {
			return errs.NewFrameError(errs.RetServerDecodeFail, "kafka consumer handler: message type invalid")
		}
		return svr.(KafkaConsumer).Handle(ctx, m)
	}
	if err := filters.Handle(ctx, nil, nil, handleFunc); err != nil {
		return nil, err
	}
	return nil, nil
}
```

定义了一个 Kafka 消费者服务的处理函数 `KafkaConsumerHandle`，该函数包裹了实际的处理逻辑：

1. 获取过滤器。
2. 创建一个处理函数，用于处理 Kafka 消息。
3. 调用过滤器链处理消息。

### RegisterKafkaConsumerService 函数

```
go
复制代码
// RegisterKafkaConsumerService 注册 service
func RegisterKafkaConsumerService(s server.Service, svr KafkaConsumer) {
	_ = s.Register(&KafkaConsumerServiceDesc, svr)
}
```

定义了一个函数 `RegisterKafkaConsumerService`，用于注册 Kafka 消费者服务。

### RegisterKafkaHandlerService 函数

```
go
复制代码
// RegisterKafkaHandlerService 注册 handle
func RegisterKafkaHandlerService(s server.Service,
	handle func(ctx context.Context, msg *sarama.ConsumerMessage) error,
) {
	_ = s.Register(&KafkaConsumerServiceDesc, kafkaHandler(handle))
}
```

定义了一个函数 `RegisterKafkaHandlerService`，用于注册 Kafka 消费者处理函数。

### BatchConsumer 接口

```
go
复制代码
// BatchConsumer 批量消费者
type BatchConsumer interface {
	// Handle 接收到消息时的回调函数
	Handle(ctx context.Context, msgArray []*sarama.ConsumerMessage) error
}
```

定义了一个批量消费者接口，包含一个 `Handle` 方法，用于处理一批 Kafka 消息。

### batchHandler 类型及其 Handle 方法

```
go
复制代码
type batchHandler func(ctx context.Context, msgArray []*sarama.ConsumerMessage) error

// Handle handle
func (h batchHandler) Handle(ctx context.Context, msgArray []*sarama.ConsumerMessage) error {
	return h(ctx, msgArray)
}
```

定义了一个 `batchHandler` 类型，它是一个函数类型，包含 `Handle` 方法，用于处理一批 Kafka 消息。

### BatchConsumerServiceDesc 描述符

```
go
复制代码
// BatchConsumerServiceDesc descriptor for server.RegisterService
var BatchConsumerServiceDesc = server.ServiceDesc{
	ServiceName: "trpc.kafka.consumer.service",
	HandlerType: ((*BatchConsumer)(nil)),
	Methods: []server.Method{
		{
			Name: "/trpc.kafka.consumer.service/handle",
			Func: BatchConsumerHandle,
		},
	},
}
```

定义了一个批量消费者服务的描述符 `BatchConsumerServiceDesc`，包含服务名称和方法。

### BatchConsumerHandle 函数

```
go
复制代码
// BatchConsumerHandle batch consumer service handler wrapper
func BatchConsumerHandle(svr interface{}, ctx context.Context, f server.FilterFunc) (rspbody interface{}, err error) {
	filters, err := f(nil)
	if err != nil {
		return nil, err
	}

	handleFunc := func(ctx context.Context, reqbody interface{}, rspbody interface{}) error {
		msg := codec.Message(ctx)
		msgs, ok := msg.ServerReqHead().([]*sarama.ConsumerMessage)
		if !ok {
			return errs.NewFrameError(errs.RetServerDecodeFail, "kafka consumer handler: message type invalid")
		}
		return svr.(BatchConsumer).Handle(ctx, msgs)
	}

	err = filters.Handle(ctx, nil, nil, handleFunc)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
```

定义了一个批量消费者服务的处理函数 `BatchConsumerHandle`，该函数包裹了实际的处理逻辑：

1. 获取过滤器。
2. 创建一个处理函数，用于处理一批 Kafka 消息。
3. 调用过滤器链处理消息。

### RegisterBatchHandlerService 函数

```
go
复制代码
// RegisterBatchHandlerService register consumer function
func RegisterBatchHandlerService(s server.Service,
	handle func(ctx context.Context, msgArray []*sarama.ConsumerMessage) error,
) {
	_ = s.Register(&BatchConsumerServiceDesc, batchHandler(handle))
}
```

定义了一个函数 `RegisterBatchHandlerService`，用于注册批量 Kafka 消费者处理函数。

### 总结

这段代码实现了 Kafka 消费者服务和批量消费者服务的注册和处理机制。具体步骤如下：

1. **KafkaConsumer 接口和实现**：定义 Kafka 消费者接口 `KafkaConsumer` 和处理函数 `kafkaHandler`。
2. **服务描述符**：定义 Kafka 消费者服务描述符 `KafkaConsumerServiceDesc` 和批量消费者服务描述符 `BatchConsumerServiceDesc`。
3. **处理函数**：实现 Kafka 消费者服务处理函数 `KafkaConsumerHandle` 和批量消费者服务处理函数 `BatchConsumerHandle`。
4. **服务注册**：提供 `RegisterKafkaConsumerService` 和 `RegisterKafkaHandlerService` 函数用于注册 Kafka 消费者服务，提供 `RegisterBatchHandlerService` 用于注册批量 Kafka 消费者服务。

通过这些步骤，Kafka 消费者服务和批量消费者服务可以与 `trpc-go` 框架集成，处理从 Kafka 消费到的消息。



`KafkaConsumerServiceDesc`有什么用



ChatGPT

`KafkaConsumerServiceDesc` 在 `trpc-go` 框架中用于描述和注册 Kafka 消费者服务。它是一个服务描述符，包含了服务的基本信息、处理类型和方法等。在 `trpc-go` 框架中，服务描述符是用于定义服务的结构和行为的核心组件。

### `KafkaConsumerServiceDesc` 的定义

```
go
复制代码
var KafkaConsumerServiceDesc = server.ServiceDesc{
	ServiceName: "trpc.kafka.consumer.service",
	HandlerType: ((*KafkaConsumer)(nil)),
	Methods: []server.Method{{
		Name: "/trpc.kafka.consumer.service/handle",
		Func: KafkaConsumerHandle,
	}},
}
```

### `KafkaConsumerServiceDesc` 的组成部分

1. **ServiceName**：服务名称
   - 这是服务的唯一标识符，在注册服务时使用。
2. **HandlerType**：处理类型
   - 这是一个接口类型，用于指定服务的处理方法。
   - `HandlerType: ((*KafkaConsumer)(nil))` 表示这个服务将使用 `KafkaConsumer` 接口来处理请求。
3. **Methods**：服务方法
   - 这是一个方法列表，定义了服务的具体方法和处理函数。
   - 每个方法都有一个名称和对应的处理函数。
   - 在这个例子中，只有一个方法，名称为 `"/trpc.kafka.consumer.service/handle"`，对应的处理函数为 `KafkaConsumerHandle`。

### `KafkaConsumerServiceDesc` 的作用

1. **服务注册**

   - 在 `trpc-go` 框架中，服务描述符用于注册服务。通过调用 `RegisterKafkaConsumerService` 函数，将 `KafkaConsumerServiceDesc` 注册到服务器实例中。

   - ```
     RegisterKafkaConsumerService
     ```

      函数：

     ```
     go
     复制代码
     func RegisterKafkaConsumerService(s server.Service, svr KafkaConsumer) {
         _ = s.Register(&KafkaConsumerServiceDesc, svr)
     }
     ```

   - 这会将 `KafkaConsumerServiceDesc` 描述的服务注册到 `trpc-go` 服务器实例 `s`，并将实现了 `KafkaConsumer` 接口的 `svr` 绑定到服务上。

2. **服务处理**

   - 服务描述符定义了服务的处理逻辑。当服务器接收到请求时，它会根据服务描述符找到对应的处理函数并调用它。
   - 在这个例子中，当有 Kafka 消息到达时，服务器会调用 `KafkaConsumerHandle` 函数来处理消息。

3. **统一管理**

   - 使用服务描述符，可以统一管理服务的结构和行为。这样可以使代码更加模块化、可维护。
   - 服务描述符将服务的名称、接口类型、方法和处理函数集中在一起，提供了清晰的服务定义。

### 示例：服务注册和处理流程

1. **定义 Kafka 消费者服务描述符**

   ```
   go
   复制代码
   var KafkaConsumerServiceDesc = server.ServiceDesc{
       ServiceName: "trpc.kafka.consumer.service",
       HandlerType: ((*KafkaConsumer)(nil)),
       Methods: []server.Method{{
           Name: "/trpc.kafka.consumer.service/handle",
           Func: KafkaConsumerHandle,
       }},
   }
   ```

2. **实现 KafkaConsumer 接口**

   ```
   go
   复制代码
   type Consumer1 struct {}
   
   func (Consumer1) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
       log.Infof("get kafka message: %+v", msg)
       // 消费业务逻辑
       return nil
   }
   ```

3. **注册 Kafka 消费者服务**

   ```
   go
   复制代码
   func main() {
       s := trpc.NewServer()
       kafka.RegisterKafkaConsumerService(s.Service("trpc.app.server.consumer1"), &Consumer1{})
       err := s.Serve()
       if err != nil {
           panic(err)
       }
   }
   ```

4. **处理 Kafka 消息**

   ```
   go
   复制代码
   func KafkaConsumerHandle(svr interface{}, ctx context.Context, f server.FilterFunc) (rspbody interface{}, err error) {
       filters, err := f(nil)
       if err != nil {
           return nil, err
       }
       handleFunc := func(ctx context.Context, reqbody interface{}, rspbody interface{}) error {
           msg := codec.Message(ctx)
           m, ok := msg.ServerReqHead().(*sarama.ConsumerMessage)
           if !ok {
               return errs.NewFrameError(errs.RetServerDecodeFail, "kafka consumer handler: message type invalid")
           }
           return svr.(KafkaConsumer).Handle(ctx, m)
       }
       if err := filters.Handle(ctx, nil, nil, handleFunc); err != nil {
           return nil, err
       }
       return nil, nil
   }
   ```

### 总结

`KafkaConsumerServiceDesc` 是在 `trpc-go` 框架中描述和注册 Kafka 消费者服务的关键组件。它定义了服务的基本信息、处理类型和方法，并用于将实现了 `KafkaConsumer` 接口的对象注册到 `trpc-go` 服务器中，从而使得服务器能够正确处理 Kafka 消息。通过这种方式，可以在 `trpc-go` 框架中方便地集成 Kafka 消费者服务。





`ConsumerMessage` 是 Kafka 消费者从 Kafka 主题中读取的消息。它封装了 Kafka 消息的各种属性和元数据，使得开发者可以方便地访问和处理这些信息。下面是对 `ConsumerMessage` 结构体各个字段的详细解释：

### `ConsumerMessage` 结构体

```
go
复制代码
// ConsumerMessage encapsulates a Kafka message returned by the consumer.
type ConsumerMessage struct {
	Headers        []*RecordHeader // 仅在 Kafka 版本 0.11+ 中设置，表示消息的头部信息
	Timestamp      time.Time       // 仅在 Kafka 版本 0.10+ 中设置，表示消息的时间戳（内部消息的时间戳）
	BlockTimestamp time.Time       // 仅在 Kafka 版本 0.10+ 中设置，表示外部（压缩）块的时间戳

	Key, Value []byte              // 消息的键和值
	Topic      string              // 消息所属的主题
	Partition  int32               // 消息所属的分区
	Offset     int64               // 消息在分区中的偏移量
}
```

### 字段解释

1. **Headers**：
   - 类型：`[]*RecordHeader`
   - 解释：消息的头部信息，包含键值对。这个字段只有在 Kafka 版本 0.11 及以上版本才会被设置。头部信息可以用于存储一些元数据，例如消息的类型、处理状态等。
2. **Timestamp**：
   - 类型：`time.Time`
   - 解释：消息的时间戳，表示消息创建的时间。这个字段只有在 Kafka 版本 0.10 及以上版本才会被设置。
3. **BlockTimestamp**：
   - 类型：`time.Time`
   - 解释：外部（压缩）块的时间戳。这个字段也是在 Kafka 版本 0.10 及以上版本才会被设置。对于压缩的消息集合，这个时间戳表示整个消息块的时间戳。
4. **Key**：
   - 类型：`[]byte`
   - 解释：消息的键。通常用于分区键，以确保具有相同键的消息被发送到同一个分区。
5. **Value**：
   - 类型：`[]byte`
   - 解释：消息的值。实际的消息内容。
6. **Topic**：
   - 类型：`string`
   - 解释：消息所属的 Kafka 主题的名称。
7. **Partition**：
   - 类型：`int32`
   - 解释：消息所属的分区。Kafka 主题是分区的，消息在特定的分区内按顺序排列。
8. **Offset**：
   - 类型：`int64`
   - 解释：消息在分区中的偏移量。每个消息在其分区内都有一个唯一的偏移量，用于标识消息在分区中的位置。
