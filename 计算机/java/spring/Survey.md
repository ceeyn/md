1.**When allowCredentials is true, allowedOrigins cannot contain the special value "*" since that cannot be set on the "Access-Control-Allow-Origin" response header. To allow credentials to a set of origins, list them explicitly or consider using "allowedOrig**

https://blog.csdn.net/qq_33532713/article/details/122898582





2.[关于@JsonIgnore的理解](https://www.cnblogs.com/czrb/p/9985196.html)

首先：@JsonIgnore是一个能够在后端发送给前端数据的时候对后端发送出的json字符串能够发挥作用的一个注解

如果比如user表，你在password这个属性上添加@JsonIgnore，就会set和get方法都不会出现password，就是这个属性不会被作用，

但是在注解get上，就会发现前端不会显示password的数据，注解在set上，后端就拿不到前端传过来的password的数据，





3.**每种不同的收参方式，也意味着前端要有对应的不同的传参方式。常见的方式有：Path["/"]，Query，Body，Form-data等。**

https://blog.csdn.net/weixin_39606048/article/details/111368220

pathvariable "/"

**requestParam "?"  以？拼接的默认不需要@requestParam注解去获取，只要方法参数名和url参数名一致即可接收**



4.**Request processing failed; nested exception is org.springframework.dao.DataIntegrityViolationException:** 

**Error updating database.  Cause: java.sql.SQLException: Field 'id' doesn't have a default value**



原因 ： int型id没有写atuo_increased

http://t.zoukankan.com/mark5-p-14268122.html

```
@TableId(value = "id", type = IdType.AUTO)
private Integer id;
auto自增id

@TableId(value = "user_id", type = IdType.ASSIGN_ID)
private String id;
字符串或长数字类型id

```







5.javax.servlet.ServletException: Could not resolve view with name 'employee' in servlet with name 'dispatcherServlet'
	at org.springframework.web.servlet.DispatcherServlet.render(DispatcherServlet.java:1373) ~[spring-webmvc-5.3.6.jar:5.3.6]
	at org.springframework.web.servlet.DispatcherServlet.processDispatchResult(DispatcherServlet.java:1138) ~[spring-webmvc-5.3.6.jar:5.3.6]
	at org.springframework.web.servlet.DispatcherServlet.doDispatch(DispatcherServlet.java:1077) ~[spring-webmvc-5.3.6.jar:5.3.6]
	at org.springframework.web.servlet.DispatcherServlet.doService(DispatcherServlet.java:962) ~[spring-webmvc-5.3.6.jar:5.3.6]
	at org.springframework.web.servlet.FrameworkServlet.processRequest(FrameworkServlet.java:1006) ~[spring-webmvc-5.3.6.jar:5.3.6]
	at org.springframework.web.servlet.FrameworkServlet.doPost(FrameworkServlet.java:909) ~[spring-webmvc-5.3.6.jar:5.3.6]

**解决**：@ResponseBody这个注解一般是作用在方法上的，加上该注解表示该方法的返回结果直接写到Http response Body中，在RequestMapping中 return返回值默认解析为跳转路径，如果你此时想让Controller返回一个字符串或者对象到前台。





**6.org.apache.tomcat.util.http.fileupload.impl.SizeLimitExceededException: the request was rejected because its size (184783270) exceeds the configured maximum (10485760)**

**报错的原因是：** springBoot项目自带的tomcat对上传的文件大小有默认的限制，SpringBoot官方文档中展示：每个文件的配置最大为1Mb，单次请求的文件的总数不能大于10Mb。

```properties
# maxFileSize 单个数据大小
spring.servlet.multipart.maxFileSize=10MB
# maxRequestSize 是总数据大小
spring.servlet.multipart.maxRequestSize=100MB
```





7.Resolved [org.springframework.http.converter.HttpMessageNotReadableException: JSON parse error: Cannot deserialize value of type `java.lang.Integer` from String "15340663845": Overflow: numeric value (15340663845) out of range of `java.lang.Integer` (-2147483648 -2147483647); nested exception is com.fasterxml.jackson.databind.exc.InvalidFormatException: Cannot deserialize value of type `java.lang.Integer` from String "15340663845": Overflow: numeric value (15340663845) out of range of `java.lang.Integer` (-2147483648 -2147483647)<EOL> at [Source: (org.springframework.util.StreamUtils$NonClosingInputStream); line: 7, column: 14] (through reference chain: com.tyut.survey.model.entity.Video["workId"])]





8.Resolved [org.springframework.web.HttpRequestMethodNotSupportedException: Request method 'DELETE' not supported]![截屏2022-11-12 17.34.53](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-12 17.34.53.png)



9.Optional int parameter 'status' is present but cannot be translated into a null value due to being declared as a primitive type. Consider declaring it as object wrapper for the corresponding primitive type![截屏2022-11-12 20.07.01](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-12 20.07.01.png)





nohup java -jar survey-0.0.1-SNAPSHOT.jar &

netstat -tln命令查看了端口状态

 put /Users/haozhipeng/Desktop/survey-0.0.1-SNAPSHOT.jar /home/hao/communication/

ssh root@152.136.105.72





**一个类A调用另一个类B的实例的时候采用new方法新建，但被调用类B里面@Value和@Autowired注解都没用了。** 导致报空指针异常。

**@Autowired相当于setter，在注入之前，对象已经实例化，是在这个接口注解的时候实例化的；
而new只是实例化一个对象，而且new的对象不能调用注入的其他类**

上面的**private MyFileLister myFileLister = new MyFileLister();**应该改成

```css
@Autowried
MyFileLister myFileLister
```



# SpringBoot文件上传同时，接收复杂参数

https://blog.csdn.net/lzhfdxhxm/article/details/127884983







**BeanUtils.copyProperties方法简单来说就是将两个字段相同的对象进行属性值的复制。如果 两个对象之间存在名称不相同的属性，则 BeanUtils 不对这些属性进行处理，需要程序手动处理。**



.zip 不能用tar解压





Mysql 不需要密码就能登陆

```sql
UPDATE user SET Password = PASSWORD(``'root'``) WHERE user = ``'root'``;
```







jar:    **ClassPath对应着Jar包的根目录，对应着编译后的target的classes目录**

war:    **classpath : 指的是打成war包以后的web-info 文件夹下面的classes 文件夹里面的路径**<img src="../images/截屏2023-02-23 下午12.47.10.png" alt="截屏2023-02-23 下午12.47.10" style="zoom:50%;" />

### Java项目中加载资源文件通常有两种方式:

**1、Class.getResource(String path)**

**2、Class.getClassLoader().getResource(String path)**

<img src="../images/截屏2023-02-23 下午12.49.23.png" alt="截屏2023-02-23 下午12.49.23" style="zoom:50%;" />

<img src="../images/截屏2023-02-23 下午12.49.46.png" alt="截屏2023-02-23 下午12.49.46" style="zoom:50%;" />

###  批量文件上传不能用swagger



### 字符串操作

中文字符在string每个也算一个

https://blog.csdn.net/m0_58761900/article/details/125014074





