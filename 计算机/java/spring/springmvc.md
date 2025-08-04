**创建springmvc项目**



![截屏2022-03-10 下午4.43.45](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 下午4.43.45.png)





**servlet在第一次发请求时创建**



#### 一.点击请求

![截屏2022-03-10 下午4.45.57](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 下午4.45.57.png)

#### **二.方法执行**

![截屏2022-03-10 下午4.47.04](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 下午4.47.04.png)



#### 三.扫描不到springmvc.xml，注解扫描不到，因此配置web.xml

![截屏2022-03-10 下午4.47.58](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 下午4.47.58.png)



#### 四.返回字符串表示jsp的名字，使用视图解析器找到页面

![截屏2022-03-10 下午4.50.52](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 下午4.50.52.png)

### 流程总结

![截屏2022-03-10 下午4.58.12](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-03-10 下午4.58.12.png)







# 注解

|      | 作用在类上                             | 作用在方法上                           | 作用在方法参数上 | 作用在方法的返回值上 |
| ---- | :------------------------------------- | -------------------------------------- | ---------------- | -------------------- |
|      | @RequestMapping（value=“”，method=“”） | @RequestMapping（value=“”，method=“”） | @RequestParam    | @ResponseBody        |
|      | @SessionAttributes                     |                                        | @RequestBody     |                      |
|      |                                        |                                        | @PathVariable    |                      |
|      |                                        |                                        | @RequestHeader   |                      |
|      |                                        |                                        | @CookieValue     |                      |
|      |                                        | @ModelAttribute                        | @ModelAttribute  |                      |



# 组件



##类型转换器（conveter）

实现Converter的接口





## 异常处理器（HandlerExceptionResolver）





## 拦截器（HandlerInterceptor）









