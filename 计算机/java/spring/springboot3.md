full:true 有判断 单实例

lite:false 无判断 多实例

## 1.@SpringBootApplication  ==
	@SpringBootConfiguration 
	@EnableAutoConfiguration 	
	@ComponentScan

```java
@SpringBootApplication
等同于
@SpringBootConfiguration 
@EnableAutoConfiguration // @AutoConfigurationPackage + @Import(AutoConfigurationImportSelector.class)
@ComponentScan("com.atguigu.boot")
```

### 其中@EnableAutoConfiguration

#### 1.@AutoConfigurationPackag

```java
@Import(AutoConfigurationPackages.Registrar.class)  //给容器中导入一个组件
public @interface AutoConfigurationPackage {}

//利用Registrar给容器中导入一系列组件
//将指定的一个包下的所有组件导入进来？MainApplication 所在包下。

```



#### 2.@Import(AutoConfigurationImportSelector.class)

```java
1、利用getAutoConfigurationEntry(annotationMetadata);给容器中批量导入一些组件
2、调用List<String> configurations = getCandidateConfigurations(annotationMetadata, attributes)获取到所有需要导入到容器中的配置类
3、利用工厂加载 Map<String, List<String>> loadSpringFactories(@Nullable ClassLoader classLoader)；得到所有的组件
4、从META-INF/spring.factories位置来加载一个文件。
	默认扫描我们当前系统里面所有META-INF/spring.factories位置的文件
    spring-boot-autoconfigure-2.3.4.RELEASE.jar包里面也有META-INF/spring.factories
    
```







@responsebody + @controller = @restcontroller



@Component + @ConfigurationProperties

```java
@Component
@ConfigurationProperties(prefix = "mycar") // 标注在类上 通过application.properties里的mycar前缀绑定， 要想注入， car必须有get set方法
public class Car 
```

@EnableConfigurationProperties + @ConfigurationProperties

```java
//@EnableConfigurationProperties(Car.class)
//1、开启Car配置绑定功能
//2、把这个Car这个组件自动注册到容器中
public class MyConfig {

```





##@SpringBootApplication() 等于三个相加

![截屏2022-01-26 上午10.23.16](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-01-26 上午10.23.16.png)

### 1、@SpringBootConfiguration

@Configuration。代表当前是一个配置类

### 2、@ComponentScan

指定扫描哪些，Spring注解；

### 3、@EnableAutoConfiguration

```java
@AutoConfigurationPackage
@Import(AutoConfigurationImportSelector.class)
public @interface EnableAutoConfiguration {}
```

#### 1、@AutoConfigurationPackage

自动配置包？指定了默认的包规则

```java
@Import(AutoConfigurationPackages.Registrar.class)  //给容器中导入一个组件
public @interface AutoConfigurationPackage {}

//利用Registrar给容器中导入一系列组件
//将指定的一个包下的所有组件导入进来？MainApplication 所在包下。

```

![截屏2022-01-26 下午3.14.16](/Users/haozhipeng/Downloads/我的笔记/images/截屏2022-01-26 下午3.14.16.png)

2.