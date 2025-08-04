# IOC



## 耦合

编译期依赖关系需要避免，A a = new A();





## **工厂模式解耦**



**工厂：创建JavaBean，**

**javabean：bean是可重用组件的意思，组件是组成部分的意思，service，dao都可以称之为组件。**

**可重用组件就是可以反复使用的组件，一个service可以被好几个servelet使用，一个dao可以被好几个service使用，就是可重用组件，也就是bean**





在实际开发中我们可以把三层的对象都使用配置文件配置起来，当启动服务器应用加载的时候，让一个类中的

方法通过读取配置文件，把这些对象创建出来**并存起来**。在接下来的使用的时候，直接拿过来用就好了。

那么，这个读取配置文件，创建和获取三层对象的类就是工厂。







## **Inversion Of Control**

上一小节解耦的思路有 2 个问题：

**1**、**存哪去？**

分析：由于我们是很多对象，肯定要找个集合来存。这时候有 Map 和 List 供选择。

到底选 Map 还是 List 就看我们有没有查找需求。有查找需求，选 Map。

所以我们的答案就是

在应用加载时，创建一个 Map，用于存放三层对象。

我们把这个 map 称之为**容器**。 

2、还是没解释什么是工厂？

**工厂就是负责给我们从容器中获取指定对象的类。这时候我们获取对象的方式发生了改变。**

原来：

我们在获取对象时，都是**采用 new 的方式。是主动的。**

![ioc](/Users/haozhipeng/Downloads/读书随笔/我的笔记/images/ioc.png)

现在：

**我们获取对象时，同时跟工厂要，有工厂为我们查找或者创建对象。是被动的。**

**这种被动接收的方式获取对象的思想就是控制反转，它是 spring 框架的核心之一。**



**========================== 2022 2.27**

1、**spring的maven配置** 



![截屏2022-02-27 下午7.10.19](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-02-27 下午7.10.19.png)



2、**resources下放置xml文件**

#### xml引入

![截屏2022-03-10 上午10.26.11](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 上午10.26.11.png)





##### ==========================2022.3.6

**构造函数注入**

![截屏2022-03-06 下午7.25.30](/Users/haozhipeng/Downloads/我的笔记/images/截屏2022-03-06 下午7.25.30.png)



**set方法注入**

![截屏2022-03-06 下午7.31.42](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-06 下午7.31.42.png)

弊端在于没写set方法则无法注入





3、**set方法注入**：**property**



**=========集合**



![截屏2022-02-27 下午7.44.35](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-02-27 下午7.44.35.png)





#### ==============================注解引入

![截屏2022-03-10 上午10.24.58](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 上午10.24.58.png)

```xml
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns:context="http://www.springframework.org/schema/context"
    xsi:schemaLocation="http://www.springframework.org/schema/beans
        https://www.springframework.org/schema/beans/spring-beans.xsd
        http://www.springframework.org/schema/context
        https://www.springframework.org/schema/context/spring-context.xsd"> 
  //替换
  
<?xml version="1.0" encoding="UTF-8"?>
<beans xmlns="http://www.springframework.org/schema/beans"
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
       xsi:schemaLocation="http://www.springframework.org/schema/beans
http://www.springframework.org/schema/beans/spring-beans.xsd">
  
  
  
  
```



##### =============2.28  注解继续



![截屏2022-02-28 上午11.59.42](/Users/haozhipeng/Downloads/我的笔记/images/截屏2022-02-28 上午11.59.42.png)





![截屏2022-02-28 下午12.04.29](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-02-28 下午12.04.29.png)





##### @Autowired原理



![截屏2022-02-28 下午12.06.47](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-02-28 下午12.06.47.png)



**如果有两个 则需改bean的id 使其匹配成功**



![截屏2022-02-28 下午12.08.59](/Users/haozhipeng/Downloads/我的笔记/images/截屏2022-02-28 下午12.08.59.png)



![截屏2022-02-28 下午1.43.18](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-02-28 下午1.43.18.png)







![截屏2022-02-28 下午1.41.51](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-02-28 下午1.41.51.png)

什么的el表达式就去制定地方获取





##### ================================2022.3.1

![截屏2022-03-01 下午1.19.02](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-01 下午1.19.02.png)



#### =================================2022.3.7

![截屏2022-03-07 下午3.43.34](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-07 下午3.43.34.png)





![截屏2022-03-07 下午4.29.48](/Users/haozhipeng/Downloads/我的笔记/images/截屏2022-03-07 下午4.29.48.png)



**diagrams**

![截屏2022-03-07 下午5.45.21](/Users/haozhipeng/Downloads/我的笔记/images/截屏2022-03-07 下午5.45.21.png)





#### **@Import**

![截屏2022-03-07 下午6.13.48](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-07 下午6.13.48.png)

=======@Configuration+ComponentScan





#### @PropertySource

![截屏2022-03-07 下午6.25.19](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-07 下午6.25.19.png)



**============================2022.3.9**



![截屏2022-03-09 下午8.03.51](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-09 下午8.03.51.png)



要成功一起成功，要失败一起失败



**连接池和线程池** 

应用初始化时会有一个线程池，每次用的时候从池中获取，不用的时候还回池中。





![截屏2022-03-09 下午9.32.29](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-09 下午9.32.29.png)



![截屏2022-03-09 下午9.33.05](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-09 下午9.33.05.png)

![截屏2022-03-09 下午9.34.09](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-09 下午9.34.09.png)







**基于子类的动态代理**

![截屏2022-03-09 下午9.38.14](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-09 下午9.38.14.png)



![截屏2022-03-09 下午9.45.54](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-09 下午9.45.54.png)





# 创建bean的方式：

| 创造方式                                                     | xml配置                                                      | 注解配置                                                     |
| :----------------------------------------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| **默认无参构造方法创造bean**                                 | <bean id=*"accountService"* class=*"com.itheima.service.impl.AccountServiceImpl"*/> | **@Component**value：指定 bean 的 id。如果不指定 value 属性，默认 bean 的 id 是当前类的类名。首字母小写。 |
| **静态工厂创造bean**                                         | <bean id=*"accountService"*class=*"com.itheima.factory.StaticFactory"* |                                                              |
| **动态工厂创造bean**(先把工厂的创建交给 spring 来管理。然后在使用工厂的 bean 来调用里面的方法)即**要创建的对象在方法的返回值上** | <bean id=*"instancFactory"* class=*"com.itheima.factory.InstanceFactory"*></bean> <bean id=*"accountService"  *factory-bean=*"instancFactory"*  factory-method=*"createAccountService"*></bean> | @Bean                                                        |



# 依赖注入（DI）

| 注入方式             | xml配置                                                      | 注解配置                                               |
| -------------------- | ------------------------------------------------------------ | ------------------------------------------------------ |
| **构造函数注入**     | <bean id=*"accountService"* class=*"com.itheima.service.impl.AccountServiceImpl"*> <constructor-arg name=*"name"* value=*"**张三**"*> </constructor-arg> |                                                        |
| **set** **方法注入** | <bean id=*"accountService"* class=*"com.itheima.service.impl.AccountServiceImpl"*> <property name=*"name"* value=*"test"*></property><property name=*"age"* value=*"21"*></property></bean> |                                                        |
| 属性直接注入         |                                                              | **@Autowired** **@Qualifier** **@Resource** **@Value** |

#### 动态代理  是代理（增强）某个类中的所有方法

![截屏2022-03-10 下午12.45.39](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 下午12.45.39.png)



![截屏2022-03-10 下午12.46.47](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 下午12.46.47.png)



![截屏2022-03-10 下午12.47.56](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 下午12.53.19.png)





# AOP

**连接点之所以都是业务层的所有方法，是因为业务层很多事物控制是重复的，将重复的事务控制抽取出来给业务层方法加强，从而形成完整的业务逻辑。所以业务层是连接点，aop就是解决业务层的事务控制的**

**业务层所有的方法都叫连接点，被增强的方法才叫切入点**

**![截屏2022-03-09 下午10.09.47](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-09 下午10.09.47.png)**

**仅有test没有被增强，所以test是连接点，不是切入点**

![截屏2022-03-09 下午10.10.39](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-09 下午10.10.39.png)

**抽取的公共代码就是通知**

![截屏2022-03-09 下午10.13.43](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-09 下午10.13.43.png)

**抽取公共代码制作的类既是通知类也是切面类**



**===================2022.3.10**

#### aop的引入

![截屏2022-03-10 上午10.23.44](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 上午10.23.44.png)





**切面就是通知（公共代码）+切入点**

![截屏2022-03-10 上午10.19.01](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 上午10.19.01.png)





**切入表达式写法**

![截屏2022-03-10 上午10.22.11](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 上午10.22.11.png)

![截屏2022-03-10 上午10.44.25](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 上午10.44.25.png)



![截屏2022-03-10 上午10.46.24](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 上午10.46.24.png)

```java
//实际开发中 切入表达式写法：切到业务层实现类下所有方法
* com.yunding.service.impl.*.*(..)
```



**配置切入点表达式**

**可以写在外面，但必须在切面之前**

![截屏2022-03-10 上午11.59.26](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 上午11.59.26.png)



**环绕通知**

![截屏2022-03-10 下午12.07.00](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 下午12.07.00.png)

![截屏2022-03-10 下午12.06.06](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 下午12.06.06.png)







#### 基于注解aop

![截屏2022-03-10 下午12.41.14](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 下午12.41.14.png)
