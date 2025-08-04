# DDL



```sql
show databases; 
select database(); # 查询当前使用数据库
create database if not exists itcast; # 
```

#### 创建表

```sql
mysql> create table tb_user (
    -> id int comment '编号’,
    '> name varchar(50) comment '名字'
    '> );
    '> select database();
    '> ;
    '> ;
    '> '
    -> ;
由'>变为-> 是因为'后面没有加' 因此mysql特别注意中英文标点的使用
  
# 表结构的查询
show create table tb_user; # 创建表语句
desc tb_user; # 表结构
  
  
```

Utf8 是三个字节，utf8-mb4是支持4个字节。



#### mysql的数据类型

在mysql中，一个中文汉字所占的字节数与编码格式有关：如果是GBK编码，则一个中文汉字占2个字节；如果是UTF8编码，则一个中文汉字占3个字节，而英文字母占1字节。

not null default *，这里的null default null

不能为null 默认为*
可以为null 默认为null

##### 数值类型

![截屏2022-05-01 11.36.08](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-05-01 11.36.08.png)



double（最长位数，几位小数）

##### 字符类型

![截屏2022-05-01 11.40.31](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-05-01 11.40.31.png)

##### 时间类型

![截屏2022-05-01 11.44.17](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-05-01 11.44.17.png)







插入数据

```sql

```





### 条件·分组·排序·分页

#### 条件查询[10种比较运算符，3种逻辑运算符]



![截屏2022-05-03 10.19.59](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-05-03 10.19.59.png)

```sql
条件列表指的是多个条件，多个条件组装起来通过逻辑运算符组装。
select * from emp where age = 18 or age = 20 or age = 40;
// 为了简化多个or连接 使用in
select * from emp where age in(18,20,40);
```





#### 聚合函数（分组查询一起使用）【5个】

![截屏2022-06-15 09.30.18](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-15 09.30.18.png)

#### 分组查询

![截屏2022-06-15 09.54.20](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-15 09.54.20.png)

```sql
查询年龄小于45员工，并根据工作地址分组，获取员工数量大于等于3的工作地址
select workaddress，count（*） where age < 45 group by workaddress having count（*） >= 3；
// select后面查询的就是要显示的 注意要显示什么就查询什么
```



#### 排序

![截屏2022-06-15 10.40.59](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-15 10.40.59.png)



#### 分页

![截屏2022-06-15 10.46.43](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-15 10.46.43.png)



#### sql执行顺序

![截屏2022-06-23 10.08.37](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-23 10.08.37.png)







#### dql小结

![截屏2022-06-23 10.10.51](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-23 10.10.51.png)





### DCL

#### 用户管理

![截屏2022-06-23 10.30.26](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-23 10.30.26.png)





#### 权限控制

![截屏2022-06-23 17.41.19](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-23 17.41.19.png)

![截屏2022-06-23 17.42.23](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-23 17.42.23.png)

### 函数

![截屏2022-06-23 18.02.12](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-23 18.02.12.png)





#### 数值函数

![截屏2022-06-23 18.18.08](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-23 18.18.08.png)



#### 日期函数

![截屏2022-06-23 18.24.38](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-23 18.24.38.png)

#### 流程函数

![截屏2022-06-23 18.40.33](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-23 18.40.33.png)





### 约束

![截屏2022-06-23 22.22.24](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-23 22.22.24.png)

![截屏2022-06-23 22.31.24](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-23 22.31.24.png)





#### 外键约束

![截屏2022-06-24 09.06.17](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-24 09.06.17.png)

### 多表关系

一对一经常做单表拆分

![截屏2022-06-24 09.16.34](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-24 09.16.34.png)

### 多表查询

![截屏2022-06-24 09.27.00](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-24 09.27.00.png)





#### 内连接

![截屏2022-06-24 09.28.20](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-24 09.28.20.png)





#### 外连接

![截屏2022-06-24 09.42.17](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-24 09.42.17.png)

### 左外连接

![截屏2022-11-11 17.01.57](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-11 17.01.57.png)



### 右外连接

![截屏2022-11-11 17.06.50](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-11 17.06.50.png)





```
左右外连接和直接查询左表右表所有数据不同的是，左右外连接可以在结果中明确指出两张表交集
```



#### 自连接

![截屏2022-11-11 17.16.54](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-11 17.16.54.png)



![截屏2022-06-25 08.19.20](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-25 08.19.20.png)

![截屏2022-11-11 17.10.38](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-11 17.10.38.png)

### 自连接外查询

![截屏2022-11-11 17.17.59](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-11 17.17.59.png)



#### 联合查询

![截屏2022-06-25 08.26.08](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-25 08.26.08.png)



![截屏2022-11-11 17.21.20](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-11 17.21.20.png)



![截屏2022-11-11 17.20.58](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-11 17.20.58.png)



### 子查询

![截屏2022-06-25 08.27.50](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-25 08.27.50.png)

#### 标量子查询

![截屏2022-06-25 08.31.42](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-25 08.31.42.png)



#### 列子查询



![截屏2022-06-25 08.40.10](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-25 08.40.10.png)

#### 行子查询

![截屏2022-06-25 08.43.06](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-25 08.43.06.png)

#### 表子查询

![截屏2022-06-25 08.49.03](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-25 08.49.03.png)

##### 练习

```
多表查询时，以思考集合的方式去写sql语句
```

![截屏2022-06-26 09.49.46](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 09.49.46.png)

```
n张表 n-1个连接条件
隐式内连接多个连接条件用and连接 并且查询条件也拥and连接
```

![截屏2022-06-26 10.13.45](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 10.13.45.png)

 



![截屏2022-06-26 10.25.47](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 10.25.47.png)



？？？？？ 

![截屏2022-06-26 11.06.21](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 11.06.21.png)

```
涉及到几张表
连接条件？(n张表 n-1个连接条件，内连接、外连接、自连接)
拆分sql语句(子查询)
```



### 事务



![截屏2022-06-26 16.27.27](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 16.27.27.png)





![截屏2022-06-26 16.31.53](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 16.31.53.png)



#### 并发事务问题

##### 脏读

![截屏2022-06-26 17.14.57](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 17.14.57.png)







##### 不可重复读

![截屏2022-06-26 17.15.32](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 17.15.32.png)



##### 幻读

![截屏2022-06-26 17.16.12](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 17.16.12.png)



#### 事务隔离级别

![截屏2022-06-26 17.29.08](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 17.29.08.png)



##### read uncommited



![截屏2022-06-26 17.26.27](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 17.26.27.png)



##### read commited

![截屏2022-06-26 17.28.09](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 17.28.09.png)

 

##### repeatable read

![截屏2022-06-26 17.34.05](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 17.34.05.png)



##### serializable

![截屏2022-06-26 17.38.08](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 17.38.08.png)





基础篇总结

![截屏2022-06-26 17.43.29](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 17.43.29.png)



#### explain

![截屏2022-06-28 14.50.33](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 14.50.33.png)

![截屏2022-06-28 14.51.12](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 14.51.12.png)

```sql
null 用不到表
system 用系统表
const 用主键或者唯一索引进行查询
eq_ref ref 用其它索引查询
index 对所有索引进行遍历
all 对全表进行扫描
```





![截屏2022-06-28 16.53.12](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 16.53.12.png)



##### 最左前缀法则

![截屏2022-06-28 16.59.21](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 16.59.21.png)



![截屏2022-06-28 17.09.26](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 17.09.26.png)



![截屏2022-06-28 17.13.01](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 17.13.01.png)

![截屏2022-06-28 17.15.26](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 17.15.26.png)

![截屏2022-06-28 17.18.04](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 17.18.04.png)

