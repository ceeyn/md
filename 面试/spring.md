



![image-20241204162115085](/Users/haozhipeng/Library/Application Support/typora-user-images/image-20241204162115085.png)





# 自动配置

https://mp.weixin.qq.com/s?__biz=MzAxODcyNjEzNQ%3D%3D&chksm=9bd0b8e2aca731f4bab79c636fabed11b9a22a25fc90bd3808409a3a614189c5ab13d9c067d6&idx=1&mid=2247488890&scene=27&sn=b6fb5b2e67629ee348dd901eb7076e8e&utm_campaign=geek_search&utm_content=geek_search&utm_medium=geek_search&utm_source=geek_search&utm_term=geek_search#wechat_redirect

1.创建autoconfig moudle，底下新建spring.factories文件导入自动配置类，自动配置类中绑定一个configrationProperties类，该类属性绑定application.yaml的配置，自动配置类通过conditional注解选择导入一些bean

2.创建starter moudle，只有pom文件，依赖autoconfig moudle

3.业务项目只依赖于starter moudle，application.yaml文件的属性被autoconfig moudle读入





## 一。自定义注解扫描

`@RpcComponentScan` 之所以能被扫描到，是因为 Spring 容器在启动时会扫描所有的配置类，并处理这些配置类上的注解。下面详细解释 `@RpcComponentScan` 是如何被扫描到的过程。

### 1. `@Configuration`注解和配置类

首先，Spring 的配置类（如 `AppConfig`）通常会被标注为 `@Configuration`。Spring 容器在启动时会扫描这些配置类，解析其上的注解。

```
java
复制代码
@Configuration
@RpcComponentScan(basePackages = "com.example.rpc")
public class AppConfig {
    // 配置类内容
}
```

### 2. `@Import`注解

`@RpcComponentScan` 注解上使用了 `@Import` 注解，这使得 Spring 容器会在处理 `@RpcComponentScan` 时，导入指定的类（如 `RpcBeanDefinitionRegistrar`）。

```
java
复制代码
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
@Documented
@Inherited
@Import(RpcBeanDefinitionRegistrar.class)
public @interface RpcComponentScan {
    @AliasFor("basePackages")
    String[] value() default {};
    @AliasFor("value")
    String[] basePackages() default {};
}
```

### 3. Spring 容器启动过程

#### 3.1 扫描和解析配置类

Spring 容器启动时，使用 `ConfigurationClassPostProcessor` 处理所有标注了 `@Configuration` 的类。这是通过 Spring 的组件扫描机制或显式注册（如通过 `AnnotationConfigApplicationContext.register` 方法）来实现的。

```
java
复制代码
AnnotationConfigApplicationContext context = new AnnotationConfigApplicationContext(AppConfig.class);
```

#### 3.2 处理 `@Import` 注解

在处理配置类时，`ConfigurationClassPostProcessor` 会扫描所有注解，并处理 `@Import` 注解。对于 `@RpcComponentScan` 注解，由于它包含 `@Import(RpcBeanDefinitionRegistrar.class)`，Spring 会实例化 `RpcBeanDefinitionRegistrar` 并调用它的 `registerBeanDefinitions` 方法。

## 二。beandefinition作用

### 2. `BeanDefinition` 对象注册过程

1. **扫描和解析**：使用 `RpcClassPathBeanDefinitionScanner` 扫描指定的包路径，找到所有标注了 `@RpcService` 注解的类。
2. **生成和注册 `BeanDefinition`**：对于每个找到的类，生成对应的 `BeanDefinition` 对象，并将其注册到 `BeanDefinitionRegistry` 中。

### 3. Bean 实例的创建

bean 实例的创建发生在 Spring 容器启动并完成 `BeanDefinition` 注册之后。当 Spring 容器需要使用某个 bean 时，会根据对应的 `BeanDefinition` 创建该 bean 的实例。

`BeanDefinition` 是 Spring 框架中的一个核心接口，用于描述一个 bean 的定义。它是 Spring IoC 容器中用于存储和管理 bean 元数据的基础结构。`BeanDefinition` 包含了 bean 的配置信息，例如 bean 的类型、作用域、依赖关系、初始化方法、销毁方法等。

### `BeanDefinition` 的作用

1. **描述 bean 的元数据**： `BeanDefinition` 对象包含了 bean 的所有配置信息，包括类名、作用域、初始化方法、销毁方法、属性值、构造函数参数等。这些信息用于告诉 Spring 容器如何实例化和管理 bean。
2. **支持 IoC 容器的操作**： Spring 容器在启动时，会根据配置文件（XML、Java 注解或 Java 配置类）创建相应的 `BeanDefinition` 对象。这些对象存储在 `BeanDefinitionRegistry` 中，容器可以根据这些定义来创建、配置和管理 bean 实例。
3. **定义 bean 的作用域**： `BeanDefinition` 可以指定 bean 的作用域，例如单例（singleton）、原型（prototype）、请求（request）、会话（session）等。作用域定义了 bean 的生命周期和可见范围。
4. **管理 bean 的依赖关系**： `BeanDefinition` 可以描述 bean 之间的依赖关系，Spring 容器会根据这些依赖关系进行依赖注入（Dependency Injection），确保在实例化某个 bean 之前，先实例化并注入它所依赖的其他 bean。

### `BeanDefinition` 的主要属性

- **beanClassName**：bean 的类名，Spring 容器通过反射机制根据类名创建 bean 实例。
- **scope**：bean 的作用域，例如单例（singleton）或原型（prototype）。
- **propertyValues**：bean 的属性值集合，包含了需要注入到 bean 中的属性及其值。
- **constructorArgumentValues**：构造函数参数值，用于在实例化 bean 时传递构造函数参数。
- **initMethodName**：初始化方法名，在 bean 实例化之后调用。
- **destroyMethodName**：销毁方法名，在 bean 销毁之前调用。
- **lazyInit**：是否延迟初始化，标记为 true 时，只有在第一次使用时才会实例化 bean。
- **dependsOn**：该 bean 所依赖的其他 bean 的名称列表。
- **primary**：是否是主要候选者，在自动装配过程中，如果有多个 bean 可用，将优先选择标记为 primary 的 bean。

### 示例：使用 `BeanDefinition`

以下是一个使用 `BeanDefinition` 的示例，展示了如何在代码中创建和注册一个 bean 定义：

```
java
复制代码
import org.springframework.beans.factory.support.DefaultListableBeanFactory;
import org.springframework.beans.factory.support.GenericBeanDefinition;

public class BeanDefinitionExample {
    public static void main(String[] args) {
        // 创建一个 BeanFactory
        DefaultListableBeanFactory beanFactory = new DefaultListableBeanFactory();

        // 创建一个 BeanDefinition 对象
        GenericBeanDefinition beanDefinition = new GenericBeanDefinition();
        beanDefinition.setBeanClass(MyBean.class);
        beanDefinition.setScope("singleton");

        // 设置属性值
        beanDefinition.getPropertyValues().add("name", "Spring Bean");

        // 注册 BeanDefinition
        beanFactory.registerBeanDefinition("myBean", beanDefinition);

        // 获取并使用 bean
        MyBean myBean = (MyBean) beanFactory.getBean("myBean");
        System.out.println(myBean.getName());
    }
}

class MyBean {
    private String name;

    public void setName(String name) {
        this.name = name;
    }

    public String getName() {
        return name;
    }
}
```

在这个示例中，我们：

1. 创建了一个 `DefaultListableBeanFactory`，它是 Spring 容器的一种实现。
2. 创建了一个 `GenericBeanDefinition`，用于定义 bean 的元数据。
3. 设置了 bean 的类名、作用域和属性值。
4. 将 `BeanDefinition` 注册到 `BeanFactory` 中。
5. 从 `BeanFactory` 中获取并使用 bean。

通过 `BeanDefinition`，我们可以灵活地定义和管理 Spring 容器中的 bean，从而实现复杂的依赖注入和配置管理。



## 三。bean的创建和beandefination创建不同

1. **Bean 初始化后处理**：
   - 当 Spring 容器初始化一个 bean 后，`postProcessAfterInitialization` 方法会被调用。
   - 如果该 bean 被 `@RpcService` 注解标注，则进行额外的处理，如服务注册和本地缓存。
2. **服务注册**：
   - 根据 `@RpcService` 注解的信息创建 `ServiceInfo` 对象。
   - 将 `ServiceInfo` 对象注册到服务注册中心（如 Zookeeper）。
3. **本地缓存**：
   - 将服务实例添加到本地缓存中，以便在 RPC 调用时快速查找。
4. **服务启动**：
   - 在 Spring Boot 应用启动后，通过 `CommandLineRunner` 接口的 `run` 方法启动 RPC 服务器。

### 总结

带有 `@RpcService` 注解的 bean 的创建时机取决于 bean 的作用域和配置：

- **单例作用域**：在容器启动时立即创建。
- **原型作用域**：在每次请求时创建。
- **延迟加载**：在首次请求时创建。



Spring 框架是 Java 开发中非常流行的框架，它运用了多种经典的设计模式来实现其强大的功能和灵活性。以下是 Spring 框架中常见的设计模式及其应用场景：

### 1. **单例模式（Singleton Pattern）**

#### 应用场景：
- 在 Spring 中，默认情况下，所有的 `Bean` 都是单例的，即在 Spring 容器中，某个 `Bean` 的实例在整个应用程序中只有一个。这是通过 `@Scope("singleton")` 或者默认的 `@Scope` 注解来实现的。

#### 作用：
- **节省资源**：避免多次创建相同的对象，节省内存开销。
- **全局访问点**：提供全局访问该实例的方式，方便管理。

#### 示例：
```java
@Component
public class MyService {
    // 默认情况下，Spring 容器中只会有一个 MyService 实例
}
```

### 2. **工厂模式（Factory Pattern）**

#### 应用场景：
- Spring 中的 `BeanFactory` 和 `ApplicationContext` 就是典型的工厂模式应用，它们通过配置文件或注解，负责创建和管理 `Bean` 的实例。

#### 作用：
- **解耦对象的创建过程**：通过工厂类来创建对象，避免在代码中直接 `new` 对象。
- **提供统一的创建接口**：可以通过工厂模式灵活地管理对象的创建逻辑。

#### 示例：
```java
ApplicationContext context = new ClassPathXmlApplicationContext("applicationContext.xml");
MyService myService = context.getBean("myService", MyService.class);
```

### 3. **代理模式（Proxy Pattern）**

#### 应用场景：
- Spring AOP（面向切面编程）是代理模式的经典应用。Spring AOP 通过动态代理（JDK 动态代理或 CGLIB 代理）为目标对象创建代理对象，来增强目标对象的方法执行，如添加事务管理、日志记录、权限控制等。

#### 作用：
- **在不修改目标对象的情况下增强功能**：代理对象可以在方法执行前后添加额外的逻辑。
- **控制访问**：代理模式可以控制对目标对象的访问，如在执行某些操作之前进行权限检查。

#### 示例：
```java
@Aspect
@Component
public class LoggingAspect {

    @Before("execution(* com.example.service.*.*(..))")
    public void logBefore(JoinPoint joinPoint) {
        System.out.println("Logging before method: " + joinPoint.getSignature().getName());
    }
}
```

### 4. **模板方法模式（Template Method Pattern）**

#### 应用场景：
- Spring 提供的 `JdbcTemplate`、`RestTemplate` 等都是模板方法模式的应用。这些模板类定义了一个操作的骨架，而将一些具体的步骤留给子类或回调函数来实现。

#### 作用：
- **复用代码**：将固定的操作步骤封装在模板方法中，避免重复代码。
- **易于扩展**：开发者只需实现具体步骤的逻辑，而不需要关心整体流程的控制。

#### 示例：
```java
JdbcTemplate jdbcTemplate = new JdbcTemplate(dataSource);
String sql = "SELECT COUNT(*) FROM users";
int count = jdbcTemplate.queryForObject(sql, Integer.class);
```

### 5. **观察者模式（Observer Pattern）**

#### 应用场景：
- Spring 的事件驱动模型（`ApplicationEvent` 和 `ApplicationListener`）是观察者模式的应用。Spring 容器中某个事件发生时，所有监听该事件的监听器都会收到通知并作出响应。

#### 作用：
- **松耦合**：对象之间通过事件进行通信，减少对象之间的直接依赖。
- **扩展性好**：可以轻松添加或移除监听器，而不需要修改事件发布者的代码。

#### 示例：
```java
@Component
public class MyEventListener implements ApplicationListener<ContextRefreshedEvent> {
    @Override
    public void onApplicationEvent(ContextRefreshedEvent event) {
        System.out.println("Context refreshed event received.");
    }
}
```

### 6. **策略模式（Strategy Pattern）**

#### 应用场景：
- Spring 中的策略模式应用非常广泛。例如，Spring 的 `TaskExecutor` 可以使用不同的策略执行任务（如线程池、异步执行等）；`DataSource` 可以使用不同的策略管理数据库连接。

#### 作用：
- **解耦策略的定义与使用**：可以在运行时选择不同的策略实现。
- **提高系统的灵活性**：通过策略模式，可以动态地更改业务逻辑而不影响客户端代码。

#### 示例：
```java
@Configuration
public class AppConfig {

    @Bean
    public TaskExecutor taskExecutor() {
        return new ThreadPoolTaskExecutor(); // 使用线程池策略执行任务
    }
}
```

### 7. **责任链模式（Chain of Responsibility Pattern）**

#### 应用场景：
- Spring 中的 `Interceptor`（拦截器）和 `Filter`（过滤器）链条就是责任链模式的典型应用。多个拦截器或过滤器可以依次处理请求，直到找到能够处理请求的对象。

#### 作用：
- **解耦请求的发送者和处理者**：请求可以在处理链中传递，直到被处理为止。
- **灵活的请求处理**：可以动态地添加、删除或重排序处理链中的处理者。

#### 示例：
```java
public class MyInterceptor implements HandlerInterceptor {
    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) {
        System.out.println("Pre Handle method is Calling");
        return true;
    }

    @Override
    public void postHandle(HttpServletRequest request, HttpServletResponse response, Object handler, ModelAndView modelAndView) {
        System.out.println("Post Handle method is Calling");
    }

    @Override
    public void afterCompletion(HttpServletRequest request, HttpServletResponse response, Object handler, Exception exception) {
        System.out.println("Request and Response is completed");
    }
}
```

### 8. **装饰者模式（Decorator Pattern）**

#### 应用场景：
- Spring AOP 也是装饰者模式的一种应用。通过装饰者模式，Spring AOP 可以在不修改原始 `Bean` 的情况下，动态地为其添加新的功能（如日志记录、事务管理等）。

#### 作用：
- **动态扩展对象的功能**：通过装饰者模式，可以在不修改原始类的情况下，添加额外的功能。
- **提高代码的可复用性**：可以通过组合多个装饰者来实现不同的功能扩展。

#### 示例：
```java
@Component
public class MyServiceDecorator implements MyService {
    private final MyService delegate;

    @Autowired
    public MyServiceDecorator(MyService delegate) {
        this.delegate = delegate;
    }

    @Override
    public void performOperation() {
        System.out.println("Before operation");
        delegate.performOperation();
        System.out.println("After operation");
    }
}
```

### 9. **工厂方法模式（Factory Method Pattern）**

#### 应用场景：
- 在 Spring 中，工厂方法模式用于创建 `Bean` 实例，例如 `FactoryBean` 接口允许开发者定制 `Bean` 的创建过程。

#### 作用：
- **控制对象的创建**：通过工厂方法，可以根据需要定制对象的创建过程。
- **解耦创建逻辑与使用逻辑**：调用者只需调用工厂方法获得实例，而无需关心实例的创建细节。

#### 示例：
```java
public class MyBeanFactory implements FactoryBean<MyService> {
    @Override
    public MyService getObject() throws Exception {
        return new MyServiceImpl();
    }

    @Override
    public Class<?> getObjectType() {
        return MyService.class;
    }

    @Override
    public boolean isSingleton() {
        return true;
    }
}
```

### 10. **桥接模式（Bridge Pattern）**

#### 应用场景：
- 在 Spring 中，桥接模式通常用于将抽象和实现解耦。例如，Spring 的 `Resource` 抽象了不同类型的资源（如文件、URL、classpath 资源等），通过桥接模式实现了对多种资源类型的统一处理。

#### 作用：
- **解耦抽象与实现**：允许它们独立变化，从而提高系统的扩展性。
- **提高代码的灵活性**：可以在不同的实现之间切换，而不影响抽象部分。

#### 示例：
```java
public class MyService {
    private Resource resource;

    public MyService(Resource resource) {
        this.resource = resource;
    }

    public void performOperation() {
        InputStream inputStream = resource.getInputStream();
        // 处理资源
    }
}
```

### 总结

Spring 框架运用了大量的设计模式，使得它在处理各种复杂问题时能够保持灵活性和高效性。通过理解这些设计模式及其在 Spring 中的应用，开发者可以更好地理解 Spring 的设计思想，并在自己的项目中应用这些模式来提高代码的质量和可维护性。







### 观察者模式（Observer Pattern）

#### 概述

观察者模式是一种行为设计模式，定义了一种**一对多**的关系，让多个观察者对象同时监听某一个主题对象。当这个主题对象的状态发生变化时，会通知所有的观察者对象，使得它们能够自动更新。这种模式常用于事件处理系统中，例如 GUI 框架、订阅/发布系统、事件驱动架构等。

在 Java 中，观察者模式是通过 `java.util.Observable` 类和 `java.util.Observer` 接口来实现的。而在 Spring 框架中，观察者模式的应用主要体现在**事件驱动模型**上，通过 `ApplicationEvent` 和 `ApplicationListener` 来实现事件的发布和监听。

#### 观察者模式的关键组成部分

1. **Subject（主题/被观察者）**：
   - 被观察者维护了一组观察者。当状态发生改变时，通知所有注册的观察者。

2. **Observer（观察者）**：
   - 观察者是对主题感兴趣的对象。当主题发生变化时，观察者会接收到通知并进行相应的操作。

3. **Event（事件）**：
   - 事件是特定的消息或数据，当某些条件满足时被发布。例如，Spring 中的 `ApplicationEvent` 就是一个事件的抽象。

4. **EventListener（事件监听器）**：
   - 事件监听器是观察者的一种具体实现。它注册到特定的事件类型，当该事件发生时，它会被通知。

### Spring 中的观察者模式应用

在 Spring 中，观察者模式被用于实现事件驱动模型。Spring 提供了一套机制，允许你在 Spring 容器中的 Bean 之间进行松耦合的通信，这套机制就是基于 `ApplicationEvent` 和 `ApplicationListener` 实现的。

#### 1. **ApplicationEvent（事件）**

`ApplicationEvent` 是 Spring 的事件基类，所有的事件都需要继承这个类。它表示应用程序中发生的某个事件，可以携带相关的消息或数据。Spring 提供了一些内置事件，比如 `ContextRefreshedEvent`、`ContextClosedEvent` 等，这些事件会在特定的时机触发。

#### 2. **ApplicationListener（事件监听器）**

`ApplicationListener` 是一个接口，所有希望监听特定事件的 Bean 需要实现这个接口。实现这个接口的 Bean 会自动注册到 Spring 的事件机制中，当相应的事件发生时，Spring 会调用该监听器的 `onApplicationEvent` 方法。

#### 3. **事件发布**

在 Spring 中，可以通过 `ApplicationEventPublisher` 或 `ApplicationContext` 来发布事件。Spring 容器会将事件发布给所有相关的 `ApplicationListener`，以通知它们进行处理。

### 观察者模式的应用场景

Spring 框架中，观察者模式主要用于事件驱动模型，它允许 Spring 容器中不同的组件以松耦合的方式进行通信。具体应用场景包括：

- **容器刷新事件**：当 Spring 容器完成刷新（初始化或刷新配置）时，发布 `ContextRefreshedEvent`，通知所有相关的监听器。例如，可以在容器刷新后进行一些资源加载或初始化操作。
  
- **容器关闭事件**：当 Spring 容器关闭时，发布 `ContextClosedEvent`，通知所有监听器执行相应的清理操作。

- **自定义事件**：用户可以定义自己的事件并发布，当某些特定条件发生时，通知感兴趣的监听器。例如，当用户注册成功后，发布一个 `UserRegisteredEvent`，通知其他组件执行相应的逻辑，如发送欢迎邮件或记录日志。

### 观察者模式的作用

1. **松耦合**：
   - 观察者模式允许对象之间通过事件进行通信，而不需要直接依赖对方。事件发布者只负责发布事件，具体谁来处理事件是由 Spring 容器来管理的。这种松耦合的设计减少了对象之间的依赖关系，提高了系统的可维护性。

2. **扩展性好**：
   - 通过观察者模式，可以很方便地添加或移除事件监听器，而不需要修改事件发布者的代码。这样，系统可以根据需要灵活地增加新功能，而不影响现有的功能。

### 代码示例

下面是一个简单的代码示例，展示了如何在 Spring 中使用观察者模式来处理容器刷新事件。

```java
import org.springframework.context.ApplicationListener;
import org.springframework.context.event.ContextRefreshedEvent;
import org.springframework.stereotype.Component;

@Component
public class MyEventListener implements ApplicationListener<ContextRefreshedEvent> {

    @Override
    public void onApplicationEvent(ContextRefreshedEvent event) {
        System.out.println("Context refreshed event received.");
        // 可以在这里添加自定义逻辑，比如初始化资源等
    }
}
```

#### 解释：

- **@Component**：这个注解将 `MyEventListener` 注册为 Spring 容器中的一个 Bean。Spring 会自动扫描并将其添加到应用上下文中。

- **ApplicationListener<ContextRefreshedEvent>**：`MyEventListener` 实现了 `ApplicationListener` 接口，并指定了它感兴趣的事件类型为 `ContextRefreshedEvent`。这意味着当 Spring 容器触发 `ContextRefreshedEvent` 事件时，`MyEventListener` 会被通知并执行 `onApplicationEvent` 方法。

- **onApplicationEvent**：当容器刷新事件发生时，Spring 会调用这个方法。在这个方法中，你可以添加任何你需要的逻辑，比如记录日志、初始化资源、启动任务等。

### 发布自定义事件

除了监听 Spring 内置的事件外，开发者还可以自定义事件并在应用中发布和处理。例如：

#### 1. 自定义事件类

```java
import org.springframework.context.ApplicationEvent;

public class UserRegisteredEvent extends ApplicationEvent {

    private String username;

    public UserRegisteredEvent(Object source, String username) {
        super(source);
        this.username = username;
    }

    public String getUsername() {
        return username;
    }
}
```

#### 2. 发布自定义事件

```java
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.stereotype.Service;

@Service
public class UserService {

    @Autowired
    private ApplicationEventPublisher eventPublisher;

    public void registerUser(String username) {
        // 业务逻辑，如保存用户信息到数据库
        System.out.println("User registered: " + username);

        // 发布事件
        UserRegisteredEvent event = new UserRegisteredEvent(this, username);
        eventPublisher.publishEvent(event);
    }
}
```

#### 3. 监听自定义事件

```java
import org.springframework.context.ApplicationListener;
import org.springframework.stereotype.Component;

@Component
public class UserRegisteredListener implements ApplicationListener<UserRegisteredEvent> {

    @Override
    public void onApplicationEvent(UserRegisteredEvent event) {
        System.out.println("User registered event received for user: " + event.getUsername());
        // 可以在这里添加用户注册后的处理逻辑，比如发送欢迎邮件等
    }
}
```

### 总结

- **观察者模式**在 Spring 中通过 `ApplicationEvent` 和 `ApplicationListener` 实现，使得组件之间可以通过事件进行通信，解耦了事件发布者和监听者之间的关系。
- **松耦合**：通过观察者模式，可以减少对象之间的直接依赖，提高系统的灵活性和可维护性。
- **扩展性**：通过添加或移除监听器，可以轻松扩展系统功能，而不需要修改事件发布者的代码。

这种设计模式在实际开发中非常有用，尤其是在需要处理多种异步事件的复杂系统中。理解和掌握 Spring 的事件驱动机制，可以帮助开发者构建更加灵活和健壮的应用程序。





为了更好地理解 Spring Boot 自动配置的工作流程，我将通过一个具体的例子，按流程详细说明自动配置是如何工作的。这个例子将以 Spring Boot 自动配置 `DataSource`（数据源）为例，演示从项目启动到自动配置生效的全过程。

### 场景：自动配置 DataSource

假设我们创建了一个简单的 Spring Boot 项目，并添加了与数据库相关的依赖，如 `spring-boot-starter-data-jpa`。我们希望 Spring Boot 能自动配置一个 `DataSource`，用于数据库连接。

### 1. **引入依赖**

首先，在项目的 `pom.xml` 中引入 Spring Boot 的 JPA 依赖包：

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-data-jpa</artifactId>
</dependency>
```

这个依赖包括了 `spring-boot-starter-jdbc`，其中包含了 JDBC 和 HikariCP（默认的数据源连接池）等相关的类。

### 2. **Spring Boot 启动应用**

当我们启动 Spring Boot 应用时，Spring Boot 会扫描项目中的所有依赖和类路径，并加载应用上下文。在启动过程中，Spring Boot 会自动进行以下操作：

1. **加载 `@SpringBootApplication` 注解**：

   项目的主启动类通常使用 `@SpringBootApplication` 注解，这个注解实际上是 `@SpringBootConfiguration`、`@EnableAutoConfiguration` 和 `@ComponentScan` 的组合。重点是 `@EnableAutoConfiguration`，它会触发 Spring Boot 的自动配置机制。

   ```java
   @SpringBootApplication
   public class MyApplication {
       public static void main(String[] args) {
           SpringApplication.run(MyApplication.class, args);
       }
   }
   ```

2. **扫描 `spring.factories` 文件**：

   Spring Boot 会在类路径下的 `META-INF/spring.factories` 文件中查找自动配置类。这些配置类包含了各种可能需要自动配置的组件，比如 `DataSourceAutoConfiguration`。

   在 `spring-boot-autoconfigure` JAR 包中，`spring.factories` 文件中可能包含如下配置：

   ```properties
   org.springframework.boot.autoconfigure.EnableAutoConfiguration=\
   org.springframework.boot.autoconfigure.jdbc.DataSourceAutoConfiguration,\
   ...
   ```

   这意味着，当 `@EnableAutoConfiguration` 被触发时，Spring Boot 会加载并执行 `DataSourceAutoConfiguration` 类。

### 3. **加载自动配置类**

`DataSourceAutoConfiguration` 是一个 Spring 配置类，负责自动配置数据源。它可能像这样定义：

```java
@Configuration
@ConditionalOnClass(DataSource.class)  // 仅在类路径中存在 DataSource 时生效
@EnableConfigurationProperties(DataSourceProperties.class)
@Import({DataSourceInitializationConfiguration.class})
public class DataSourceAutoConfiguration {

    @Bean
    @ConditionalOnMissingBean
    public DataSource dataSource(DataSourceProperties properties) {
        // 使用 DataSourceProperties 构建并返回一个 DataSource Bean
        return properties.initializeDataSourceBuilder().build();
    }
}
```

这个配置类做了以下事情：

1. **条件判断**：
   - `@ConditionalOnClass(DataSource.class)`：仅当类路径中存在 `DataSource` 类时，才会激活这个配置类。
   - `@ConditionalOnMissingBean`：如果 Spring 上下文中没有其他的 `DataSource` Bean，则创建并注册一个默认的 `DataSource` Bean。

2. **加载配置属性**：
   - `@EnableConfigurationProperties(DataSourceProperties.class)`：启用 `DataSourceProperties`，用于从 `application.properties` 或 `application.yml` 中加载数据源的相关配置（如 URL、用户名、密码等）。

### 4. **初始化数据源**

如果满足条件，Spring Boot 会调用 `dataSource()` 方法，初始化并注册一个 `DataSource` Bean。此时，Spring Boot 会检查 `application.properties` 或 `application.yml` 文件中的配置项：

```yaml
spring:
  datasource:
    url: jdbc:mysql://localhost:3306/mydb
    username: root
    password: password
    driver-class-name: com.mysql.cj.jdbc.Driver
```

这些配置会被加载到 `DataSourceProperties` 中，并用于初始化 `DataSource`。

### 5. **注册 Bean 到 Spring 上下文**

当 `DataSourceAutoConfiguration` 完成后，Spring Boot 会将配置好的 `DataSource` Bean 注册到 Spring 应用上下文中。这个 `DataSource` Bean 现在可以在应用的任何地方通过依赖注入的方式被使用：

```java
@Service
public class MyService {

    @Autowired
    private DataSource dataSource;

    public void doSomething() {
        // 使用数据源进行数据库操作
        Connection connection = dataSource.getConnection();
        // ...
    }
}
```

### 6. **应用启动完成**

当所有的自动配置类都被加载和执行后，Spring Boot 完成了整个应用上下文的构建，并启动应用。在这个例子中，数据源的配置完全由 Spring Boot 自动完成，开发者无需手动配置大量的 `Bean` 或 XML 文件。

### 7. **覆盖默认配置**

如果你需要自定义数据源配置，Spring Boot 提供了多种方式覆盖默认的自动配置：

1. **通过配置文件覆盖**：直接在 `application.properties` 或 `application.yml` 文件中覆盖默认的配置。

2. **自定义 Bean**：你也可以通过定义一个 `@Bean` 来覆盖自动配置的 `DataSource` Bean：

   ```java
   @Configuration
   public class MyDataSourceConfig {
   
       @Bean
       public DataSource customDataSource() {
           // 自定义数据源
           return new HikariDataSource();
       }
   }
   ```

   由于 Spring Boot 的 `@ConditionalOnMissingBean` 注解，自动配置的 `DataSource` Bean 将不会被创建，而是使用你自定义的 Bean。

### 总结

通过这个例子，我们可以看到 Spring Boot 自动配置的整个流程：

1. **引入依赖**：通过 `pom.xml` 引入所需的依赖，如 `spring-boot-starter-data-jpa`。
2. **启动应用**：`@EnableAutoConfiguration` 触发自动配置机制。
3. **扫描并加载自动配置类**：Spring Boot 从 `spring.factories` 文件中读取自动配置类，并根据条件判断是否加载配置。
4. **初始化 Bean**：根据配置类的逻辑，Spring Boot 自动初始化并注册 `Bean`，如 `DataSource`。
5. **使用 Bean**：Spring 上下文中注册的 `Bean` 可以通过依赖注入在应用的其他地方被使用。
6. **覆盖配置**：开发者可以通过配置文件或自定义 `Bean` 来覆盖默认的自动配置。

Spring Boot 的自动配置极大地简化了应用开发，减少了繁琐的配置工作，使得开发者可以更专注于业务逻辑的实现。





在 Spring Boot 项目中，`spring-boot-configuration-processor` 和 `spring-boot-autoconfigure-processor` 是两个非常重要的依赖项，它们用于处理 Spring Boot 的自动配置和自定义配置属性。下面将详细解释它们的作用、工作原理和使用场景。

### 1. **spring-boot-configuration-processor**

#### 作用：
`spring-boot-configuration-processor` 是 Spring Boot 提供的一个注解处理器，它用于生成帮助开发者编写配置类的元数据文件。在使用 Spring Boot 的自定义配置属性时，这个处理器能够帮助生成配置提示（比如在 IDE 中自动补全）和验证配置值。

#### 工作原理：
- 当你在 Spring Boot 项目中定义了一个配置属性类（通常使用 `@ConfigurationProperties` 注解）时，`spring-boot-configuration-processor` 会扫描这些类，并生成相应的元数据文件。这些文件位于 `META-INF/spring-configuration-metadata.json`，它们描述了配置属性的名称、类型、默认值等信息。
- 这个元数据文件被 IDE 读取后，可以在 `application.properties` 或 `application.yml` 中提供自动补全和提示功能。

#### 示例：
假设你定义了一个自定义配置类：

```java
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

@Component
@ConfigurationProperties(prefix = "app")
public class MyAppProperties {
    
    private String name;
    private int timeout;
    
    // Getters and Setters

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public int getTimeout() {
        return timeout;
    }

    public void setTimeout(int timeout) {
        this.timeout = timeout;
    }
}
```

在这个例子中，`MyAppProperties` 类定义了两个配置属性 `app.name` 和 `app.timeout`。`spring-boot-configuration-processor` 会扫描这个类，并生成对应的元数据文件。

生成的 `META-INF/spring-configuration-metadata.json` 文件可能包含如下内容：

```json
{
  "groups": [
    {
      "name": "app",
      "type": "com.example.MyAppProperties",
      "sourceType": "com.example.MyAppProperties"
    }
  ],
  "properties": [
    {
      "name": "app.name",
      "type": "java.lang.String",
      "description": "Name of the application",
      "sourceType": "com.example.MyAppProperties"
    },
    {
      "name": "app.timeout",
      "type": "java.lang.Integer",
      "description": "Timeout for the application",
      "sourceType": "com.example.MyAppProperties"
    }
  ]
}
```

这个元数据文件使得 IDE 可以识别和提供 `app.name` 和 `app.timeout` 的自动补全提示。

#### 使用场景：
- 在 Spring Boot 项目中使用自定义配置属性时，建议引入 `spring-boot-configuration-processor` 依赖，尤其是在开发公共库或框架时，它能够帮助其他开发者更好地使用和配置这些属性。

### 2. **spring-boot-autoconfigure-processor**

#### 作用：
`spring-boot-autoconfigure-processor` 是 Spring Boot 提供的另一个注解处理器，它用于优化和加速 Spring Boot 自动配置类的加载过程。具体来说，它帮助生成自动配置的元数据文件，Spring Boot 可以在启动时根据这些元数据更快地确定哪些自动配置类需要加载。

#### 工作原理：
- 当 Spring Boot 应用启动时，`@EnableAutoConfiguration` 注解会触发自动配置类的加载。Spring Boot 会扫描类路径上的 `META-INF/spring.factories` 文件，找到所有的自动配置类。
- `spring-boot-autoconfigure-processor` 处理器生成的元数据文件（如 `META-INF/spring-autoconfigure-metadata.properties`）可以帮助 Spring Boot 在启动时更高效地决定哪些自动配置类需要加载，哪些类可以跳过，从而提高启动速度。

#### 示例：
如果你开发一个包含自动配置类的库（如 `spring-boot-starter` 类库），你可能会编写如下的自动配置类：

```java
@Configuration
@ConditionalOnClass(MyService.class)
public class MyServiceAutoConfiguration {
    
    @Bean
    @ConditionalOnMissingBean
    public MyService myService() {
        return new MyServiceImpl();
    }
}
```

`spring-boot-autoconfigure-processor` 会在编译时生成一个 `META-INF/spring-autoconfigure-metadata.properties` 文件，该文件描述了自动配置类和条件注解的信息，使得 Spring Boot 在启动时更快地做出配置决策。

生成的 `META-INF/spring-autoconfigure-metadata.properties` 文件可能包含如下内容：

```properties
# Automatically generated
org.springframework.boot.autoconfigure.EnableAutoConfiguration=\
com.example.MyServiceAutoConfiguration
```

#### 使用场景：
- 当你开发一个 Spring Boot 的 Starter 包或其他包含自动配置类的库时，建议引入 `spring-boot-autoconfigure-processor` 依赖，以优化自动配置的处理过程。
- 对于普通的 Spring Boot 应用项目，通常不需要手动添加这个依赖，除非你开发了自定义的自动配置类。

### 总结

- **`spring-boot-configuration-processor`**：用于生成配置类的元数据文件，帮助 IDE 提供配置属性的自动补全和提示功能。适用于任何自定义配置类的 Spring Boot 项目，特别是公共库开发中。
  
- **`spring-boot-autoconfigure-processor`**：用于优化 Spring Boot 自动配置类的处理，生成自动配置的元数据文件，提升应用启动速度。适用于开发包含自动配置类的 Spring Boot Starter 包或其他框架扩展。

这两个依赖项都可以在编译时自动生成元数据文件，使得 Spring Boot 应用更加智能和高效，提升了开发体验和应用性能。





好，接下来我将更加详细地说明 Spring 解决循环依赖的机制、无法处理的情况，以及如何在遇到这种问题时进行解决。首先，我们需要了解什么是循环依赖，以及 Spring 在创建 Bean 时的生命周期，最后再深度剖析不同注入方式下循环依赖的解决和限制。

### 1. 循环依赖的定义

**循环依赖**，也称为循环引用，发生在两个或多个 Bean 相互依赖时。例如：

- `Bean A` 依赖于 `Bean B`，而 `Bean B` 依赖于 `Bean A`，这就形成了一个循环。
  

这种情况下，Spring 需要在实例化其中一个 Bean 时找到另一个 Bean，但是由于它们互相依赖而都未完成创建，导致依赖注入无法完成。

### 2. Spring Bean 的创建过程

为了更好地理解循环依赖问题，首先要了解 Spring 在创建 Bean 时的整个生命周期：

1. **实例化 Bean**：Spring 调用类的构造器，创建 Bean 的实例（这时没有注入任何依赖）。
2. **将未完成的 Bean 放入三级缓存**：Spring 有一个三级缓存机制用于解决循环依赖，缓存包含：
    - 一级缓存（`singletonObjects`）：存储**已经完全初始化的 Bean。**
    - 二级缓存（`earlySingletonObjects`）：存储**已实例化但未完全初始化的 Bean。**
    - 三级缓存（`singletonFactories`）：存储 **BeanFactory 的代理工厂，用于创建早期 Bean 的引用。**
3. **依赖注入**：通过 Setter 方法或构造器，将**依赖项注入到 Bean 中。**
4. **初始化完成**：在 Bean 初始化完成后，**Spring 会将其从二级缓存中移入一级缓存**。

Spring 解决 Setter 注入中的循环依赖的关键步骤是在第 2 步时，提前暴露出未完成初始化的 Bean，使其能够被其他 Bean 注入。

### 3. Spring 如何解决 Setter 注入的循环依赖

Setter 注入指的是通过 Setter 方法将依赖注入到 Bean 中，**这种方式比较容易解决循环依赖。Spring 通过提前暴露一个 Bean 的引用，即使该 Bean 尚未完全初始化，也能够将其注入到其他 Bean 中**。

#### 3.1 Setter 注入的工作机制

假设有以下两个类，分别为 `A` 和 `B`，它们互相依赖：

```java
@Component
public class A {
    private B b;

    @Autowired
    public void setB(B b) {
        this.b = b;
    }
}

@Component
public class B {
    private A a;

    @Autowired
    public void setA(A a) {
        this.a = a;
    }
}
```

在这种情况下，Spring 创建 `A` 时，会经历以下步骤：

1. **实例化 `A`**：Spring 调用 `A` 的默认构造器，创建 `A` 的实例。此时，`A` 中的 `b` 属性尚未注入。
2. **暴露 `A` 的早期引用**：Spring 会将 `A` 的早期引用放入三级缓存（`singletonFactories`）中，方便后续的依赖注入。如果其他 Bean 需要 `A`，可以从缓存中获取到未完全初始化的 `A`。
3. **创建 `B` 并注入 `A`**：接着，Spring 开始创建 `B`。在 `B` 的依赖注入阶段，`B` 需要 `A`，Spring 会从三级缓存中获取到 `A` 的早期引用并注入 `B` 中。
4. **完成 `B` 的创建**：`B` 被完全初始化，并放入一级缓存中。
5. **注入 `B` 到 `A`**：最后，Spring 将 `B` 注入到 `A` 的 `b` 属性中，完成 `A` 的初始化，并将 `A` 放入一级缓存中。

通过这种机制，Spring 能够打破 Setter 注入的循环依赖问题。

### 4. Spring 无法解决的循环依赖：构造器注入

构造器注入是通过构造函数将依赖传递给 Bean，在对象创建时立即进行依赖注入。在构造器注入的场景下，Spring 无法提前暴露 Bean 的早期引用，因此不能解决构造器注入的循环依赖。

#### 4.1 构造器注入中的循环依赖示例

```java
@Component
public class A {
    private B b;

    @Autowired
    public A(B b) {
        this.b = b;
    }
}

@Component
public class B {
    private A a;

    @Autowired
    public B(A a) {
        this.a = a;
    }
}
```

在这个示例中，`A` 和 `B` 都使用构造器注入，导致 Spring 无法完成 Bean 的初始化。流程如下：

1. Spring 尝试实例化 `A`，但由于 `A` 的构造器依赖于 `B`，因此 Spring 需要先创建 `B`。
2. 然而，在创建 `B` 时，`B` 的构造器又依赖于 `A`。
3. 由于两者都需要对方的实例才能完成构造，Spring 无法进行下一步，导致循环依赖问题。

#### 4.2 构造器注入不能解决循环依赖的原因

构造器注入的循环依赖无法解决的原因主要在于：
- **构造器注入时 Bean 必须完全初始化**：构造器注入要求所有依赖在构造方法调用时必须存在，因此 Bean 必须完全初始化，这与 Spring 通过缓存暴露早期 Bean 的机制冲突。
- **无法提前暴露 Bean**：在构造器注入过程中，Spring 无法将一个未初始化的 Bean 暴露出来，因为构造器本身需要依赖对象的完整实例。

### 5. 三级缓存的作用与机制

为了解决 Setter 注入时的循环依赖问题，Spring 引入了三级缓存：

1. **一级缓存（`singletonObjects`）**：存储完全初始化的单例对象。
2. **二级缓存（`earlySingletonObjects`）**：存储早期但部分初始化的 Bean。
3. **三级缓存（`singletonFactories`）**：存储 Bean 工厂的引用，允许创建早期 Bean 的代理对象。

### 6. 如何解决构造器注入的循环依赖

虽然 Spring 不能自动解决构造器注入的循环依赖，但我们可以通过以下几种方式来手动避免或解决这个问题：

#### 6.1 使用 Setter 注入或 Field 注入

最简单的解决办法是将构造器注入改为 Setter 注入或 Field 注入。这种方法利用 Spring 内置的机制来处理循环依赖。

```java
@Component
public class A {
    private B b;

    @Autowired
    public void setB(B b) {
        this.b = b;
    }
}

@Component
public class B {
    private A a;

    @Autowired
    public void setA(A a) {
        this.a = a;
    }
}
```

#### 6.2 使用 `@Lazy` 延迟加载

使用 `@Lazy` 注解，可以让 Spring 在真正需要某个 Bean 时再去创建它，从而打破构造器注入的循环依赖。`@Lazy` 会使被依赖的 Bean 延迟初始化，避免在 Bean 创建时立刻需要它的依赖项。

```java
@Component
public class A {
    private B b;

    @Autowired
    public A(@Lazy B b) {
        this.b = b;
    }
}

@Component
public class B {
    private A a;

    @Autowired
    public B(@Lazy A a) {
        this.a = a;
    }
}
```

#### 6.3 拆分依赖关系

另一个解决构造器注入循环依赖的方法是重构代码，拆分依赖关系，避免互相依赖。例如，利用中间层将互相依赖的 Bean 分开。

### 7. 总结

- **Spring 能解决的循环依赖**：Spring 可以自动处理通过 **Setter 注入** 和 **Field 注入** 产生的循环依赖。它通过三级缓存机制，在依赖注入时提前暴露部分初始化的 Bean 来解决问题。
- **Spring 无法解决的循环依赖**：Spring 无法自动解决 **构造器注入** 中的循环依赖，因为构造器注入要求依赖项在 Bean 创建时已经完全可用。
- **解决方案**：可以通过将构造器注入改为 Setter 注入、使用 `@Lazy` 注解、或重构依赖关系来解决构造器注入引发的循环依赖问题。

Spring 通过灵活的依赖注入机制和缓存机制，在大多数场景下能够很好地解决循环依赖问题，但在设计系统时，尽量避免复杂的循环依赖关系，尤其是在使用构造器注入时要特别小心。



确实，`A` 的引用不会直接放入二级缓存，而是通过 **三级缓存** 提前暴露。我们需要明确三级缓存和二级缓存的区别以及它们的具体作用。在 Spring 的三级缓存机制中，首先是通过 **三级缓存**（`singletonFactories`）暴露 `A` 的早期引用，然后当某些条件满足时，`A` 才会从三级缓存移到二级缓存中。让我们详细解释这个过程。

### 1. 三级缓存与二级缓存的关系

- **三级缓存（`singletonFactories`）**：存放的是 `ObjectFactory`，用于在需要时创建 Bean 的代理对象或早期引用。它是最早的一个缓存，用于延迟创建代理对象，确保 AOP 和其他代理机制的正常工作。
  
- **二级缓存（`earlySingletonObjects`）**：存放已经创建但尚未完全初始化的 **早期引用**。在代理场景中，三级缓存中的 `ObjectFactory` 会根据需要将代理对象或早期引用转移到二级缓存中。

因此，三级缓存实际上是用来动态生成对象（或代理）的工厂，而二级缓存是用于存储这个生成的早期引用（可以是代理对象，也可以是原始对象）。当代理不涉及时，二级缓存可能会直接保存对象引用，而不需要三级缓存的 `ObjectFactory`。

### 2. 创建过程中的缓存工作机制

我们再以一个 **`A` 依赖 `B`，`B` 依赖 `A`** 的例子来详细说明 `A` 是如何逐步放入三级缓存、二级缓存以及一级缓存的。

#### Step 1：开始创建 `A`

当 Spring 容器开始创建 `A` 时，它首先会检查 **一级缓存（`singletonObjects`）**，看看 `A` 是否已经完全初始化并可以直接使用。如果 `A` 还未创建，则一级缓存中不会有它的实例。

```java
Object singleton = this.singletonObjects.get("A");  // 从一级缓存获取A，结果为null
```

接着，Spring 会检查 **二级缓存（`earlySingletonObjects`）**。二级缓存存储的是已经创建但未完全初始化的 Bean 的引用（早期引用），如果此时 `A` 还未暴露到二级缓存中，结果也为 `null`。

```java
Object earlySingleton = this.earlySingletonObjects.get("A");  // 从二级缓存获取A，结果为null
```

最后，Spring 还会检查 **三级缓存（`singletonFactories`）**，三级缓存中存放的是 `ObjectFactory`，可以动态生成 Bean 的引用或代理。如果 `A` 还没有放入三级缓存，结果也为 `null`。

```java
ObjectFactory<?> singletonFactory = this.singletonFactories.get("A");  // 从三级缓存获取A的工厂，结果为null
```

由于 `A` 还未创建，Spring 开始实例化 `A`。

```java
A a = new A();  // 开始创建A
```

#### Step 2：提前暴露 `A` 到三级缓存

为了避免后续创建 `B` 时发生循环依赖，Spring 会 **提前暴露 `A` 的引用**。但是此时，`A` 还没有完全初始化，所以不会直接放入二级缓存或一级缓存。

而是通过一个 **`ObjectFactory`** 将 `A` 的引用放入 **三级缓存**：

```java
this.singletonFactories.put("A", () -> getEarlyBeanReference(a));
```

`getEarlyBeanReference(a)` 是一个方法，用于在需要时返回 `A` 的早期引用（或代理对象）。这一步非常关键，它允许 Spring 以后通过这个工厂来获取 `A` 的引用，解决循环依赖问题。



#### Step 3：开始创建 `B`

接下来，Spring 发现 `A` 依赖于 `B`，于是开始创建 `B`。在创建 `B` 的过程中，`B` 依赖于 `A`，所以需要从缓存中获取 `A`。

此时，Spring 检查缓存的顺序如下：

1. **一级缓存** 中没有 `A`，因为 `A` 尚未完成初始化。
2. **二级缓存** 中也没有 `A`，因为 `A` 的引用还没有进入二级缓存。
3. **三级缓存** 中存在 `A` 的 `ObjectFactory`。

Spring 使用 `ObjectFactory` 从 **三级缓存** 中获取 `A` 的早期引用：

```java
A earlyA = this.singletonFactories.get("A").getObject();  // 获取A的早期引用
```

这个早期引用可以是 `A` 的原始对象，也可以是 `A` 的代理对象（如果涉及 AOP 代理等机制）。



#### Step 4：将 `A` 的引用放入二级缓存

当 `B` 获取到 `A` 的早期引用后，Spring 将 `A` 放入 **二级缓存**，以便后续可以更快地使用 `A`，而不需要再次通过 `ObjectFactory` 从三级缓存中获取。

```java
this.earlySingletonObjects.put("A", earlyA);  // 将A的早期引用放入二级缓存
```

此时，`A` 的早期引用已经进入了二级缓存，其他依赖 `A` 的 Bean 可以直接从二级缓存中获取 `A` 的引用。



#### Step 5：完成 `B` 的创建并放入一级缓存

`B` 获取到 `A` 的早期引用后，完成了 `B` 的创建和初始化。`B` 完全创建后，Spring 将其放入 **一级缓存** 中，表示 `B` 已经完全初始化。

```java
this.singletonObjects.put("B", b);  // 将B放入一级缓存
this.singletonFactories.remove("B");  // 从三级缓存移除B的工厂
```



#### Step 6：完成 `A` 的创建并放入一级缓存

现在返回到 `A` 的创建过程。在 `B` 创建完成后，`A` 也继续完成初始化。`A` 完成初始化后，将其从 **二级缓存** 和 **三级缓存** 中移除，并最终放入 **一级缓存**：

```java
this.singletonObjects.put("A", a);  // 将A放入一级缓存
this.earlySingletonObjects.remove("A");  // 从二级缓存移除A
this.singletonFactories.remove("A");  // 从三级缓存移除A的工厂
```

至此，`A` 和 `B` 的循环依赖问题得到了成功解决。



### 3. 为什么 `A` 不直接放入二级缓存？

### 核心原因：三级缓存提供了 **代理对象的能力**。而二级缓存仅仅存放的是 **早期引用**，这可能是原始对象，也可能是代理对象。

1. **如果没有代理**：
   - Spring 可以直接将 `A` 的早期引用放入二级缓存，这种情况较为简单。
   - 在这种情况下，`A` 通过二级缓存暴露，能够解决大多数简单的循环依赖问题。

2. **如果存在代理（如 AOP）**：
   - 在代理对象场景下，Spring 不能直接将原始对象 `A` 放入二级缓存。
   - **代理对象的生成需要在 Bean 完全初始化后进行，所以在此之前，Spring 只能通过 `ObjectFactory` 来延迟生成代理对象。**
   - 三级缓存的 `ObjectFactory` 提供了延迟生成代理对象的能力，这样当 `B` 需要 `A` 时，它获取到的是 **代理对象** 而不是原始对象。

因此，在 Spring 的三级缓存设计中，**先通过三级缓存暴露 `ObjectFactory`，然后再将生成的早期引用（可以是代理对象）放入二级缓存**，确保代理机制和依赖注入的正确性。

### 4. 总结：三级缓存的必要性

三级缓存的存在是为了应对复杂的 Bean 创建场景，特别是当 Bean 涉及 **代理** 或 **AOP** 时，单纯依赖二级缓存无法提供动态代理功能：

- **三级缓存（`singletonFactories`）** 存储的是 `ObjectFactory`，它可以在需要时生成代理对象或早期引用。
- **二级缓存（`earlySingletonObjects`）** 仅存储提前暴露的原始 Bean 实例或者代理对象（如果已经生成）。
- **一级缓存（`singletonObjects`）** 存储完全初始化的单例 Bean。

通过三级缓存，Spring 能够在 Bean 尚未完全初始化时，通过 `ObjectFactory` 提供代理对象，确保 AOP 等功能正常运作。这是二级缓存无法实现的功能。



### 三级缓存如何创建代理对象    getEarlyBeanFactories，相当于在初始化前提前执行动态代理

````
在 Spring 的三级缓存机制中，`ObjectFactory` 提供了延迟生成代理对象的能力。当 `B` 需要 `A` 时，通过三级缓存中的 `ObjectFactory` 可以确保 `B` 注入的是 **代理对象**，而不是 `A` 的原始对象。这种机制主要用于解决 **AOP 代理** 和 **循环依赖** 并存的场景。

我将详细解释 Spring 是如何通过 `ObjectFactory` 实现代理对象的延迟生成，并在需要时注入代理对象，而不是原始对象。为了更清楚理解，我们首先回顾一下相关背景，然后逐步解析各个环节。

### 1. 背景：代理对象与原始对象

在 Spring 中，代理对象通常用于支持 **AOP（Aspect-Oriented Programming，面向切面编程）** 功能，例如：

- **事务管理**：通过 `@Transactional` 注解，Spring 会为某些方法创建代理对象，用于在方法执行前后进行事务的管理。
- **日志记录**：通过 AOP 代理对象，可以在目标方法执行前后进行日志记录。

这些功能依赖于 **代理对象** 的存在，而非直接操作 Bean 的原始对象。当涉及到 AOP 代理时，Spring 不仅需要创建 Bean，还需要创建该 Bean 的代理对象。这个代理对象会在一些方法调用时执行额外的逻辑，例如事务的开启和提交。

### 2. 三级缓存的核心作用

**三级缓存** 的核心是通过 `ObjectFactory` 延迟生成代理对象。三级缓存中的 `ObjectFactory` 提供了一个方法，当某个依赖 `A` 的 Bean（例如 `B`）需要 `A` 时，Spring 通过调用 `ObjectFactory.getObject()` 来返回代理对象，而不是原始的 `A` 对象。

### 3. 代理对象生成的过程

Spring 是如何通过三级缓存中的 `ObjectFactory` 生成代理对象的呢？关键在于 **`getEarlyBeanReference()`** 方法的实现。我们先看下其主要步骤：

#### Step 1：创建 Bean `A`

假设 `A` 是一个通过 AOP 增强的 Bean，当 Spring 开始创建 `A` 时，它并不会立刻生成代理对象，而是先创建 `A` 的原始对象。

```java
A a = new A();  // 创建 A 的原始对象
```

此时 `A` 还没有完全初始化，并且没有进行代理对象的创建。

#### Step 2：将 `A` 的 `ObjectFactory` 放入三级缓存

为了处理依赖注入中的循环依赖问题，Spring 会提前将 `A` 的引用暴露出来，但并不会直接暴露原始的 `A`。相反，它会将一个 **`ObjectFactory`** 放入三级缓存 `singletonFactories` 中。这个 `ObjectFactory` 可以在以后根据需要生成 `A` 的代理对象或早期引用。

```java
this.singletonFactories.put("A", () -> getEarlyBeanReference(a, "A"));
```

`getEarlyBeanReference(a, "A")` 是一个方法，它决定返回的是原始对象还是代理对象。

#### Step 3：`getEarlyBeanReference()` 方法

`getEarlyBeanReference()` 方法是 Spring 中实现代理对象生成的关键。其主要逻辑如下：

```java
protected Object getEarlyBeanReference(Object bean, String beanName) {
    // 检查当前 Bean 是否需要被代理
    if (bean instanceof Advised) {  // 判断是否需要代理
        return bean;  // 如果不需要代理，返回原始对象
    }

    // 使用 AOP 的相关机制，创建代理对象
    return createProxy(bean, beanName);  // 创建代理对象
}
```

- **检查是否需要代理**：Spring 会首先检查 `A` 是否配置了 AOP 或其他代理增强。如果不需要代理，`getEarlyBeanReference()` 方法直接返回原始的 `A` 对象。
- **创建代理对象**：如果 `A` 需要代理（例如 `@Transactional` 注解），`getEarlyBeanReference()` 会调用代理工厂方法 `createProxy()` 来为 `A` 创建代理对象。

#### Step 4：依赖注入时调用 `ObjectFactory.getObject()`

当 `B` 依赖 `A` 时，Spring 会尝试通过三级缓存中的 `ObjectFactory` 获取 `A` 的引用。在这个过程中，Spring 调用 `ObjectFactory.getObject()` 方法，这个方法的实现如下：

```java
ObjectFactory<?> singletonFactory = this.singletonFactories.get("A");
Object earlyA = singletonFactory.getObject();  // 获取 A 的代理对象或原始对象
```

在 `ObjectFactory.getObject()` 方法中，Spring 调用的是 `getEarlyBeanReference(a, "A")`，从而决定返回 `A` 的原始对象还是代理对象。

- 如果 `A` 需要代理，`ObjectFactory.getObject()` 返回的是 `A` 的代理对象。
- 如果 `A` 不需要代理，`ObjectFactory.getObject()` 返回的是 `A` 的原始对象。

通过这种机制，Spring 可以根据 Bean 的实际需求，灵活地在依赖注入过程中返回原始对象或代理对象。

#### Step 5：完成 `A` 的创建和放入一级缓存

在依赖注入完成后，Spring 会继续对 `A` 进行初始化，并最终将 `A`（无论是代理对象还是原始对象）放入一级缓存 `singletonObjects` 中：

```java
this.singletonObjects.put("A", a);
this.singletonFactories.remove("A");
```

此时，`A` 已经完全初始化，三级缓存中的 `ObjectFactory` 也被清除。`A` 的代理对象或原始对象将保存在一级缓存中，供后续使用。

### 4. 代理对象的生成：AOP 与 `AdvisedSupport`

当 Spring 确定 `A` 需要代理时，它使用 **AOP 代理工厂** 来生成代理对象。Spring 的 AOP 代理对象生成通常分为以下几个步骤：

#### Step 1：获取 `AdvisedSupport`

Spring 使用 `AdvisedSupport` 类来保存代理对象的元数据（如目标对象、切面等）。Spring 通过这个类来管理代理对象的创建和行为。

#### Step 2：使用代理工厂创建代理对象

Spring 有两种常见的代理方式：
- **JDK 动态代理**：用于接口代理。如果目标对象实现了接口，Spring 使用 JDK 动态代理来创建代理对象。
- **CGLIB 代理**：用于类代理。如果目标对象没有实现接口，Spring 使用 CGLIB 生成目标对象的子类，并通过子类来代理目标类。

Spring 在创建代理对象时，会根据目标对象的类型和配置，选择合适的代理方式。

```java
// 判断是否使用 JDK 动态代理或 CGLIB
if (bean instanceof SomeInterface) {
    return ProxyFactory.createJdkDynamicProxy(bean);  // 使用JDK动态代理
} else {
    return ProxyFactory.createCglibProxy(bean);  // 使用CGLIB代理
}
```

#### Step 3：返回代理对象

创建完成后，代理对象会被返回给 `ObjectFactory`，并通过三级缓存中的 `getObject()` 方法被注入到依赖它的 Bean 中。通过这种方式，Spring 确保了在依赖注入过程中，如果某个 Bean 需要代理增强，其依赖的 Bean 获取到的是代理对象而不是原始对象。

### 5. 三级缓存与 AOP 的关系

通过三级缓存中的 `ObjectFactory`，Spring 可以延迟生成代理对象。这是因为代理对象的生成通常依赖于 Bean 的完整生命周期，只有在 Bean 完全初始化后，代理对象才能正确生成。

三级缓存的 `ObjectFactory` 通过 `getEarlyBeanReference()` 方法，在 Bean 的早期创建阶段暴露代理对象，使得即便在循环依赖场景下，Spring 依然可以通过 `ObjectFactory` 生成代理对象并完成依赖注入。

### 6. 总结

- **三级缓存中的 `ObjectFactory` 提供了延迟生成代理对象的能力**。在依赖注入过程中，当 `B` 依赖 `A` 时，Spring 通过三级缓存中的 `ObjectFactory`，可以灵活决定返回原始对象还是代理对象。
  
- **代理对象的生成过程** 是通过 `getEarlyBeanReference()` 方法实现的，这个方法会根据 Bean 的配置（如 AOP 增强）来决定是否返回代理对象。

- **最终的代理对象** 或原始对象在依赖注入完成后，会被放入一级缓存 `singletonObjects` 中，供后续使用。这保证了代理机制与依赖注入的正常工作，即使在复杂的循环依赖场景下。

通过这种三级缓存机制，Spring 成功解决了代理对象的延迟生成问题，并确保在循环依赖的情况下，注入的是代理对象而不是原始对象。这也是三级缓存相较于二级缓存的重要优势之一。
````

