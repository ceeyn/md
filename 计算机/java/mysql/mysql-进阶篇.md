![截屏2022-06-26 21.20.45](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 21.20.45.png)

![截屏2022-06-26 21.21.17](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 21.21.17.png)





![截屏2022-06-26 21.57.33](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 21.57.33.png)





![截屏2022-06-26 22.09.35](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-26 22.09.35.png)





#### myisam特点

![截屏2022-06-27 09.38.22](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-27 09.38.22.png)









![截屏2022-06-27 09.40.05](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-27 09.40.05.png)



##### memory特点

![截屏2022-06-27 09.39.10](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-27 09.39.10.png)



![截屏2022-06-27 09.38.01](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-27 09.38.01.png)









### 索引

![截屏2022-06-28 09.44.09](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 09.44.09.png)



![截屏2022-06-28 09.46.55](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 09.46.55.png)





![截屏2022-06-28 09.47.20](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 09.47.20.png)





![截屏2022-06-28 10.00.21](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 10.00.21.png)



#### 哈希索引

![截屏2022-06-28 11.19.17](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 11.19.17.png)



![截屏2022-06-28 11.22.07](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 11.22.07.png)



![截屏2022-06-28 11.23.29](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 11.23.29.png)



![截屏2022-06-28 11.25.26](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 11.25.26.png)





![截屏2022-06-28 11.30.24](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 11.30.24.png)



 ![截屏2022-06-28 11.31.33](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 11.31.33.png)





![截屏2022-06-28 11.50.55](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 11.50.55.png)





![截屏2022-06-28 11.58.22](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 11.58.22.png)



##### sql执行频率查看

![截屏2022-06-28 12.09.20](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 12.09.20.png)

![截屏2022-06-28 12.16.26](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 12.16.26.png)



##### show profiles

![截屏2022-06-28 14.19.40](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 14.19.40.png)

![截屏2022-06-28 14.20.16](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 14.20.16.png)

#### explain

![截屏2022-06-28 14.50.33](data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 2342 1252"></svg>)

![截屏2022-06-28 14.51.12](data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 2430 1234"></svg>)

```sql
null 用不到表
system 用系统表
const 用主键或者唯一索引进行查询
eq_ref ref 用其它索引查询
index 对所有索引进行遍历
all 对全表进行扫描
```





![截屏2022-06-28 16.53.12](data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 2440 1116"></svg>)



##### 最左前缀法则

![截屏2022-06-28 16.59.21](data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 2392 1204"></svg>)



![截屏2022-06-28 17.09.26](data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 2368 648"></svg>)



![截屏2022-06-28 17.13.01](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 17.13.01.png)

![截屏2022-06-28 17.15.26](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 17.15.26.png)

![截屏2022-06-28 17.18.04](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 17.18.04.png)

![截屏2022-06-28 17.58.58](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 17.58.58.png)

![截屏2022-06-28 21.24.57](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 21.24.57.png)



##### 覆盖索引

![截屏2022-06-28 21.30.00](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 21.30.00.png)



![截屏2022-06-28 21.33.04](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 21.33.04.png)

##### 前缀索引

![截屏2022-06-28 21.41.07](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 21.41.07.png)

 ![截屏2022-06-28 21.43.39](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 21.43.39.png)





![截屏2022-06-28 21.55.17](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 21.55.17.png)

![截屏2022-06-28 21.57.33](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 22.01.57.png)







![截屏2022-06-28 22.09.11](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 22.09.11.png)





![截屏2022-06-28 22.08.11](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-28 22.08.11.png)





### SQL优化



![截屏2022-06-29 09.08.24](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 09.08.24.png)

![截屏2022-06-29 09.11.46](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 09.11.46.png)





 ![截屏2022-06-29 10.03.08](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 10.03.08.png)

![截屏2022-06-29 15.59.07](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 15.59.07.png)

![截屏2022-06-29 16.02.10](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 16.02.10.png)



#### order by优化

![截屏2022-06-29 16.04.24](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 16.04.24.png)







![截屏2022-06-29 16.09.22](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 16.09.22.png)

```sql
先按age升序 再按phone升序 因此会反向排
```







![截屏2022-06-29 16.33.21](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 16.33.21.png)

![截屏2022-06-29 16.35.45](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 16.35.45.png)

![截屏2022-06-29 16.36.32](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 16.36.32.png)

 

#### group by优化

![截屏2022-06-29 16.43.56](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 16.43.56.png)

```
对条件子句按最左法则判断，对select按索引覆盖判断
```

 ![截屏2022-06-29 16.46.53](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 16.46.53.png)



#### Limit 优化

![截屏2022-06-29 16.55.29](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 16.55.29.png)



#### Count 优化

![截屏2022-06-29 16.57.26](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 16.57.26.png)

![截屏2022-06-29 17.00.17](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 17.00.17.png)

![截屏2022-06-29 17.02.53](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 17.02.53.png)

![截屏2022-06-29 17.11.12](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 17.11.12.png)

![截屏2022-06-29 17.16.27](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-29 17.16.27.png)





## 视图

![截屏2022-06-30 09.47.29](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 09.47.29.png)



#### cascade

![截屏2022-06-30 10.03.23](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 10.03.23.png)

![截屏2022-06-30 10.05.45](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 10.05.45.png)

```
cascade会对上层即使没有with check option也做检查
```



#### local

![截屏2022-06-30 11.49.38](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 11.49.38.png)

![截屏2022-06-30 11.49.52](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 11.49.52.png)

```
local会对上层带有check option做检查 不带则不检查
```





![截屏2022-06-30 11.54.15](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 11.54.15.png)

![截屏2022-06-30 11.56.56](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 11.56.56.png)



### 存储过程

![截屏2022-06-30 12.11.00](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 12.11.00.png)

![截屏2022-06-30 12.12.07](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 12.12.07.png)



![截屏2022-06-30 14.37.07](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 14.37.07.png)

![截屏2022-06-30 14.36.49](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 14.36.49.png)

 

#### 系统变量

![截屏2022-06-30 14.49.39](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 14.49.39.png)



#### 用户变量

![截屏2022-06-30 15.35.49](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 15.35.49.png)

 

#### 局部变量

![截屏2022-06-30 16.00.49](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 16.00.49.png)

![截屏2022-06-30 17.16.48](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 17.16.48.png)

##### case

![截屏2022-06-30 18.10.07](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-06-30 18.10.07.png)



![截屏2022-07-02 17.10.03](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-02 17.10.03.png)

![截屏2022-07-02 17.09.31](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-02 17.09.31.png)

![截屏2022-07-02 17.22.16](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-02 17.22.16.png)



### 游标

****游标的引出****

![截屏2022-07-07 17.58.54](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-07 17.58.54.png)

![截屏2022-07-07 17.47.23](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-07 17.47.23.png)

#### 条件处理程序

![截屏2022-07-07 18.04.21](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-07 18.04.21.png)





### 存储函数

![截屏2022-07-07 21.12.14](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-07 21.12.14.png)



## 触发器

![截屏2022-07-07 21.16.57](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-07 21.16.57.png)





![截屏2022-07-07 21.19.09](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-07 21.19.09.png)





![截屏2022-07-07 21.33.40](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-07 21.33.40.png)



![截屏2022-07-07 21.35.30](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-07 21.35.30.png)





![截屏2022-07-07 21.35.51](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-07 21.35.51.png)



#### 全局锁

![截屏2022-07-08 16.26.25](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-08 16.26.25.png)

![截屏2022-07-08 16.31.05](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-08 16.31.05.png)

 ![截屏2022-07-08 16.34.45](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-07-08 16.34.45.png)