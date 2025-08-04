**2022.3.12**

![截屏2022-03-12 下午2.18.53](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-12 下午2.18.53.png)



![截屏2022-03-15 上午9.44.50_副本](/Users/haozhipeng/Downloads/我的笔记/images/截屏2022-03-15 上午9.44.50_副本.png)



![截屏2022-03-15 上午9.55.53](/Users/haozhipeng/Downloads/我的笔记/images/截屏2022-03-15 上午9.55.53.png)







![截屏2022-03-15 上午10.00.25](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-15 上午10.00.25.png)





![截屏2022-03-15 下午12.39.40](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-15 下午12.39.40.png)





![截屏2022-03-15 下午12.42.12](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-15 下午12.42.12.png)





![自定义Mybatis分析](/Users/haozhipeng/Downloads/我的笔记/images/自定义Mybatis分析.png)



<img src="/Users/haozhipeng/Downloads/我的笔记/images/非常重要的一张图-分析编写dao实现类Mybatis的执行过程.png" alt="非常重要的一张图-分析编写dao实现类Mybatis的执行过程" style="zoom: 200%;" />



![非常重要的一张图](/Users/haozhipeng/Downloads/我的笔记/images/非常重要的一张图.png)



![无标题](/Users/haozhipeng/Downloads/我的笔记/images/无标题.png)

![自定义mybatis开发流程图](/Users/haozhipeng/Downloads/我的笔记/images/自定义mybatis开发流程图.png)



![截屏2022-03-15 下午4.31.33](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-15 下午4.31.33.png)



查询条件是综合的查询条件，不仅包括用户查询条件还包括其它的查询条件（比如将用户购买商品信息也作为查询条件），这时可以使用包装对象传递输入参数。



**file协议   uri的获取方法之一**

![截屏2022-03-15 下午5.40.01](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-15 下午5.40.01.png)



![无标题](/Users/haozhipeng/Library/Application Support/typora-user-images/无标题.png)

**mybatis多表查询**

#### 一对一

##### 1、定义专门的 po 类作为输出类型，其中定义了 sql 查询结果集所有的字段。

```sql
实现查询账户信息时，也要查询账户所对应的用户信息。
SELECT 
account.*,
user.username,
user.address
FROM
account,
user
WHERE account.uid = user.id
```



##### 2、使用 resultMap，定义专门的 resultMap 用于映射一对一查询结果。

通过面向对象的(has a)关系可以得知，我们可以在 Account 类中加入一个 User 类的对象来代表这个账户

是哪个用户的。



#### 一对多

分析：

用户信息和他的账户信息为一对多关系，并且查询过程中如果用户没有账户信息，此时也要将用户信息

查询出来，我们想到了左外连接查询比较合适。

```sql
SELECT
u.*, acc.id id,
acc.uid,
acc.money
FROM
user u
LEFT JOIN account acc ON u.id = acc.uid
```

一对多不需要建表，只需要加外键。多对多需要另建一张表，不需要外键。

#### 多对多

需求：

实现查询所有对象并且加载它所分配的用户信息。

分析：

查询角色我们需要用到Role表，但角色分配的用户的信息我们并不能直接找到用户信息，而是要通过中

间表(USER_ROLE 表)才能关联到用户信息。

下面是实现的 SQL 语句：

```sql
SELECT u.*, r.id rid, r.role_name roleName, r.role_desc roleDesc from `user` u INNER JOIN user_role ur
on(u.id=ur.uid) INNER JOIN role r on(ur.rid=r.id) 
```



![截屏2022-03-17 下午1.41.55](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-17 下午1.41.55.png)





![截屏2022-03-17 下午2.22.36](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-17 下午2.22.36.png)





![截屏2022-03-17 下午3.18.06](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-17 下午3.18.06.png)









```java
@Results 注解
代替的是标签<resultMap>
该注解中可以使用单个@Result 注解，也可以使用@Result 集合
@Results（{@Result（），@Result（）}）或@Results（@Result（））
@Resutl 注解
代替了 <id>标签和<result>标签
@Result 中 属性介绍：
id 是否是主键字段
column 数据库的列名
property 需要装配的属性名
one 需要使用的@One 注解（@Result（one=@One）（）））
many 需要使用的@Many 注解（@Result（many=@many）（）））
@One 注解（一对一）
代替了<assocation>标签，是多表查询的关键，在注解中用来指定子查询返回单一对象。
@One 注解属性介绍：
select 指定用来多表查询的 sqlmapper
fetchType 会覆盖全局的配置参数 lazyLoadingEnabled。。
使用格式：
@Result(column=" ",property="",one=@One(select=""))
@Many 注解（多对一）
代替了<Collection>标签,是是多表查询的关键，在注解中用来指定子查询返回对象集合。
注意：聚集元素用来处理“一对多”的关系。需要指定映射的 Java 实体类的属性，属性的 javaType
（一般为 ArrayList）但是注解中可以不定义；
使用格式：
@Result(property="",column="",many=@Many(select=""))
```





Mapper 映射对象，映射器。

delegate 代理，代表

executor 执行人

**返回值是为了使用反射创建对象封装结果集。**

为什么实体类属性名要和表中列名一样：为了**使用反射封装属性值**，**根据列名反射创造出属性名的对象**，然后赋值。 不一致会怎么样呢？

**不一致时，返回值无法封装，此时需要定义resultmap。参数值需要修改成类的属性名称才能传递，**

\** #{user.username}它会先去找 user 对象，然后在 user 对象中找到 username 属性，并调用**

**getUsername()方法把值取出来。但是我们在 parameterType 属性上指定了实体类名称，所以可以省略 user.**

**而直接写 username。**

ognl 表达式：

它是 apache 提供的一种表达式语言，全称是：

Object Graphic Navigation Language 对象图导航语言

它是按照一定的语法格式来获取数据的。

语法格式就是使用 #{对象.对象}的方式

别名 **typeAliases**讲的是类型别名，定义后由全限定类名改为类名。





图中自定义dao使用内部代码是从获取预处理对象时开始，到得到结果集结束。**在executor中执行sql语句，在handler中封装结果（？？？），sqlsession实现最后结果，**没有结果集反射封装，关键词：executor

图中代理使用是从获取预处理对象时开始，到得到结果集结束。**使用mapper处理包装sql语句，执行sql语句并封装结果集在executor中，在sqlsession中通过代理实现最后结果。**并没有结果集反射封装，关键词：Mapper

executor是执行sql语句的地点。



