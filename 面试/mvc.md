

5.6 @ResponseBody注解的作用
作用： 该注解用于将Controller的方法返回的对象，通过适当的HttpMessageConverter转换为指定格式后，写入到Response对象bod数据区。使用时机：返回的数据不是html标签的页面，而是其他某种格式的数据时（如json、xml等）使用；





3：拦截器于过滤器对比
两者都是AOP编程思想的实现，都能够实现权限控制和日志记录等问题的处理，但是两者粒度不同拦截对象不一样

适用范围不同：Filter是servlet的规范，只能用于web程序，但是拦截器可以用于application等程序。

规范不同：Filter是servlet的规范。但是Interceptor是spring容器支撑，有spring框架支持。

使用资源不一样：spring的拦截器由于依赖spring，也是spring的一个组件，因此能够在拦截器中使用spring的任何资源和对象。例如service对象，数据源，事务管理等，通过ioc注入拦截器即可，而filter不能

粒度不同：Filter只能在servlet的前后起作用，而拦截器能在方法前后异常前后执行，更加灵活，粒度更小，spring框架程序优先使用拦截器。


<img src="/Users/haozhipeng/Library/Application Support/typora-user-images/image-20240923162833456.png" alt="image-20240923162833456" style="zoom:50%;" />



### 拦截器

**Spring MVC 拦截器**（`HandlerInterceptor`）是 Spring MVC 框架中的一种机制，用于在处理 HTTP 请求的过程中拦截并执行特定的逻辑。它在请求到达控制器（Controller）之前、处理请求之后，以及响应生成之后的不同阶段进行拦截，允许开发者在这些阶段添加自定义的处理逻辑。

拦截器在 Spring MVC 中常用于身份验证、权限检查、记录日志、数据预处理等功能。它类似于 Servlet 的过滤器（Filter），但拦截器更细粒度且与 Spring MVC 深度集成。

### 拦截器的三大方法

Spring MVC 的拦截器基于 **`HandlerInterceptor`** 接口进行实现，该接口提供了三个核心方法：

1. **`preHandle()`**：在请求到达控制器之前执行（可以在此进行权限验证、身份验证等）。
2. **`postHandle()`**：在控制器处理完请求，但还未生成视图时执行（可以修改 ModelAndView 数据）。
3. **`afterCompletion()`**：在整个请求完成之后执行，通常用于资源清理或异常处理。

````


#### HandlerInterceptor 的方法签名：
```java
public interface HandlerInterceptor {
    // 在处理请求之前调用
    boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) throws Exception;

    // 在处理请求之后，生成视图之前调用
    void postHandle(HttpServletRequest request, HttpServletResponse response, Object handler, ModelAndView modelAndView) throws Exception;

    // 请求处理完成后调用，通常用于清理资源
    void afterCompletion(HttpServletRequest request, HttpServletResponse response, Object handler, Exception ex) throws Exception;
}
```

### 拦截器执行的顺序

- **preHandle()**：此方法在控制器方法执行之前调用。可以在这里进行身份验证、权限检查、日志记录等操作。返回值是一个 `boolean` 值：
  - 如果返回 `true`，请求会继续执行，控制器方法会被调用。
  - 如果返回 `false`，请求会被中断，控制器方法不会执行，通常会直接返回一个响应（如重定向到登录页面）。

- **postHandle()**：此方法在控制器方法执行之后，且视图渲染之前调用。可以用来对 `ModelAndView` 进行处理，如在返回给前端之前修改模型数据，添加全局数据，或改变视图名称。

- **afterCompletion()**：此方法在整个请求处理完成后（包括视图渲染完成）调用，通常用于清理资源，或记录请求的结束时间。这个方法也会在请求处理过程中出现异常时被调用，可以用于异常日志记录等。

### 拦截器的使用场景

- **身份验证与授权**：在用户访问某些受保护的资源时，拦截器可以检查用户是否已经登录，或者是否具有访问该资源的权限。
- **日志记录**：拦截器可以记录请求的开始时间、结束时间、请求内容、响应内容等。
- **全局数据处理**：在 `postHandle` 中，可以在响应之前统一处理全局数据，传递给前端视图。
- **异常处理**：在 `afterCompletion` 方法中，可以处理在请求处理过程中发生的异常。

### 拦截器的实现与配置

#### 1. **创建拦截器**

一个拦截器通常是实现了 `HandlerInterceptor` 接口的类。下面是一个简单的拦截器示例，用于检查用户是否登录：

```java
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.ModelAndView;

public class LoginInterceptor implements HandlerInterceptor {

    // 请求处理前的拦截逻辑
    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) throws Exception {
        // 检查用户是否登录，假设我们通过 session 中的 "user" 属性判断登录状态
        Object user = request.getSession().getAttribute("user");
        if (user == null) {
            // 如果用户未登录，重定向到登录页面，并拦截请求
            response.sendRedirect("/login");
            return false;
        }
        // 如果用户已登录，放行请求
        return true;
    }

    // 请求处理后视图渲染前的处理
    @Override
    public void postHandle(HttpServletRequest request, HttpServletResponse response, Object handler, ModelAndView modelAndView) throws Exception {
        // 可以在这里对 ModelAndView 进行修改
    }

    // 请求完成后的处理（包括异常情况）
    @Override
    public void afterCompletion(HttpServletRequest request, HttpServletResponse response, Object handler, Exception ex) throws Exception {
        // 用于资源清理或异常处理
    }
}
```

#### 2. **注册拦截器**

拦截器需要通过 Spring 配置来进行注册和应用。你可以使用 Java 配置或 XML 配置来注册拦截器，并指定拦截的路径模式。

##### 使用 Java 配置

在 Spring Boot 或 Spring 的 Java 配置模式下，可以通过实现 `WebMvcConfigurer` 来添加拦截器：

```java
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

@Configuration
public class WebConfig implements WebMvcConfigurer {

    @Autowired
    private LoginInterceptor loginInterceptor;

    @Override
    public void addInterceptors(InterceptorRegistry registry) {
        // 注册拦截器，并指定拦截的路径模式
        registry.addInterceptor(loginInterceptor)
                .addPathPatterns("/profile/**", "/admin/**") // 拦截的路径
                .excludePathPatterns("/login", "/register"); // 排除不需要拦截的路径
    }
}
```

在这个配置中，`addInterceptor()` 方法将 `LoginInterceptor` 注册为拦截器，并通过 `addPathPatterns()` 来指定哪些路径需要被拦截。通过 `excludePathPatterns()` 可以排除某些路径不受拦截器的影响（比如登录和注册页面不需要拦截）。

##### 使用 XML 配置

如果使用 XML 配置拦截器，可以这样进行配置：

```xml
<mvc:interceptors>
    <mvc:interceptor>
        <mvc:mapping path="/profile/**"/>
        <mvc:mapping path="/admin/**"/>
        <mvc:exclude-mapping path="/login"/>
        <mvc:exclude-mapping path="/register"/>
        <bean class="com.example.LoginInterceptor"/>
    </mvc:interceptor>
</mvc:interceptors>
```

在 XML 配置中，通过 `<mvc:interceptor>` 标签来定义拦截器，并指定拦截和排除的路径。

### 3. **拦截器与过滤器的区别**

虽然拦截器和 Servlet 过滤器（`Filter`）的功能类似，都是用于拦截 HTTP 请求和响应，但它们之间有以下不同点：

- **拦截器是 Spring MVC 特有的**，与 Spring MVC 的处理流程（如控制器和视图）紧密结合。
- **过滤器是 Servlet 规范的一部分**，拦截的是整个请求和响应的生命周期，通常用于更底层的处理。
- **拦截器可以访问 Spring MVC 的特定上下文**（如 `Handler` 和 `ModelAndView`），而过滤器只能操作请求和响应对象。
- 拦截器处理的是经过调度器分发（DispatcherServlet）之后的请求，而过滤器可以在任何请求处理之前进行拦截。

### 4. **拦截器执行顺序**

如果有多个拦截器，它们的执行顺序与注册顺序有关。在多个拦截器中：
- **`preHandle()`** 方法按注册顺序执行。
- **`postHandle()` 和 `afterCompletion()`** 方法按相反顺序执行。

例如，假设有两个拦截器 A 和 B，注册顺序是 A -> B：
- `preHandle()` 方法的执行顺序是：A -> B。
- `postHandle()` 和 `afterCompletion()` 方法的执行顺序是：B -> A。

### 5. **常见的拦截器应用场景**

- **身份验证**：在 `preHandle()` 中检查用户的登录状态，如果未登录，则重定向到登录页面。
- **权限控制**：在 `preHandle()` 中判断用户是否有权限访问某些特定的资源。
- **日志记录**：可以记录请求的开始时间，在 `afterCompletion()` 中记录请求的结束时间，以计算请求的处理时间。
- **全局数据传递**：在 `postHandle()` 中添加一些全局的数据到 `ModelAndView`，如当前用户的角色信息、公共数据等。
- **异常处理**：在 `afterCompletion()` 中捕获并处理请求过程中抛出的异常，统一进行日志记录或返回定制的错误页面。

### 6. **总结**

- Spring MVC 拦截器是用于在 HTTP 请求的不同阶段进行拦截和处理的工具。
- 通过实现 `
````



### 6.2 Spring MVC怎么样设定重定向和转发的？

（1）转发：在返回值前面加"forward:“，譬如"forward:user.do?name=method4” **301**
（2）重定向：在返回值前面加"redirect:“，譬如"redirect:www.baidu.com” **302**

````
在 Spring MVC 中，**重定向**（redirect）和**转发**（forward）是两种常见的请求处理方式。它们用于控制请求的流转路径，但两者之间有明显的区别，具体表现在请求的方式、服务器与客户端的交互、数据的传递等方面。接下来我将详细解释 Spring MVC 中如何设置转发和重定向，并探讨它们的原理及使用场景。

---

### 1. **转发（Forward）**

#### 1.1 什么是转发？

**转发** 是指服务器在处理客户端的请求时，将请求转交给另一个内部资源（如另一个控制器或 JSP 页面）继续处理，但客户端（浏览器）并不知道转发的发生。**URL 不会改变**，整个请求过程只在服务器内部完成。转发发生在服务器端，它避免了客户端与服务器之间的再次通信。

#### 1.2 Spring MVC 中的转发设置

在 Spring MVC 中，通过在控制器返回值前加上 `forward:` 来实现转发操作。

**示例**：
```java
@Controller
public class UserController {
    @RequestMapping("/testForward")
    public String testForward() {
        // 转发请求到 /user.do 并带上参数
        return "forward:/user.do?name=method4";
    }
}
```

- **`forward:`**：告知 Spring MVC 将当前请求转发到指定路径 `/user.do`。
- **参数传递**：如果有参数，它们可以继续通过请求传递，因为转发的请求是同一个请求。

#### 1.3 转发的流程

1. **客户端发起请求**：客户端向服务器发送 HTTP 请求，例如请求 `/testForward`。
2. **控制器处理请求**：Spring MVC 的 `DispatcherServlet` 调用控制器的处理方法，返回 `forward:/user.do`。
3. **服务器内部转发**：服务器通过 `RequestDispatcher.forward()` 方法将请求内部转发到 `/user.do`，而不会通知客户端，也不会改变 URL。
4. **继续处理请求**：转发后的新目标处理请求，处理逻辑完成后返回结果，响应给客户端。
   
在整个转发流程中，浏览器端的 URL 不会发生变化，它仍然是最初的请求 URL（如 `/testForward`）。

#### 1.4 转发的特点

- **请求不变**：转发不会创建新的请求，原始请求的数据（如请求头、请求体、请求参数）都会保留并传递给目标资源。
- **服务器端行为**：转发只在服务器内部发生，浏览器并不会察觉，因此客户端的 URL 不会改变。
- **适用场景**：转发通常用于服务器内部的资源跳转，例如在处理完某个逻辑后，直接将请求转发给 JSP 页面进行渲染，或者转发给另一个控制器继续处理。

---

### 2. **重定向（Redirect）**

#### 2.1 什么是重定向？

**重定向** 是指服务器在处理客户端请求时，**告诉客户端（浏览器）重新发送请求到另一个 URL**。重定向是一种客户端行为，**URL 会发生改变**，浏览器接收到重定向响应后，会发起新的 HTTP 请求到新的 URL。重定向适合需要让客户端感知跳转的场景。

#### 2.2 Spring MVC 中的重定向设置

在 Spring MVC 中，通过在控制器返回值前加上 `redirect:` 来实现重定向操作。

**示例**：
```java
@Controller
public class UserController {
    @RequestMapping("/testRedirect")
    public String testRedirect() {
        // 重定向到百度首页
        return "redirect:https://www.baidu.com";
    }
}
```

- **`redirect:`**：告知 Spring MVC 返回重定向响应，通知浏览器请求另一个 URL。

#### 2.3 重定向的流程

1. **客户端发起请求**：客户端请求 `/testRedirect`。
2. **控制器处理请求**：Spring MVC 的 `DispatcherServlet` 调用控制器方法，返回 `redirect:https://www.baidu.com`。
3. **返回重定向响应**：服务器返回 HTTP 状态码 302（Found）和 `Location` 头部，指示客户端应该重新发送请求到 `https://www.baidu.com`。
4. **浏览器重新发送请求**：客户端（浏览器）接收到 302 响应后，自动向 `Location` 头部指定的新 URL 发送新的 HTTP 请求。
5. **新请求处理**：服务器处理新的请求，返回结果给客户端。

#### 2.4 重定向的特点

- **请求发生变化**：重定向会发起两个请求，第一次是对原始 URL 的请求，第二次是对重定向目标 URL 的请求。原始请求的数据不会自动传递给新的请求。
- **客户端感知跳转**：浏览器通过 HTTP 302 状态码感知重定向，并主动向新的 URL 发起请求，因此客户端的 URL 会改变。
- **适用场景**：重定向通常用于一些需要跳转页面的场景，如用户登录成功后，重定向到首页，或提交表单后重定向到新的页面以防止表单重复提交。

#### 2.5 重定向的参数传递

重定向时，如果需要传递参数，通常通过 **URL 拼接** 或 **Session** 来实现。

**通过 URL 拼接参数**：
```java
@Controller
public class UserController {
    @RequestMapping("/testRedirectWithParam")
    public String testRedirectWithParam() {
        // 重定向并附带参数
        return "redirect:/userProfile?userId=123";
    }
}
```

通过这种方式，参数会附加在重定向的 URL 上，浏览器会发送带参数的请求。

---

### 3. **转发与重定向的区别**

| **特性**        | **转发（Forward）**                         | **重定向（Redirect）**                        |
|-----------------|--------------------------------------------|---------------------------------------------|
| **发生位置**    | 服务器内部行为，浏览器不可见               | 服务器通知浏览器跳转，浏览器发起新的请求     |
| **URL 变化**    | 浏览器的 URL 不变                           | 浏览器的 URL 改变                           |
| **请求次数**    | 一次请求                                    | 两次请求（原始请求 + 新请求）                |
| **请求数据**    | 保留原始请求的数据（参数、请求头等）         | 原始请求数据不会自动传递                    |
| **适用场景**    | 服务器内部资源跳转，不需要客户端感知的跳转 | 登录后跳转首页、表单提交后跳转新页面等需要客户端感知的跳转 |
| **性能开销**    | 仅在服务器内部跳转，性能开销较小            | 需要客户端发起新请求，存在一定的网络开销    |
| **状态码**      | 无状态码变化（同一个请求）                  | HTTP 302（重定向）                          |

---

### 4. **转发与重定向的使用场景**

#### 4.1 转发的使用场景

- **页面跳转**：在请求完成后，服务器需要将结果转发给 JSP 页面渲染，而不需要客户端感知。
- **逻辑处理**：请求进入某个控制器后，需要转发给另一个控制器或服务来完成处理。例如，模块间的数据传递。

#### 4.2 重定向的使用场景

- **登录重定向**：用户登录成功后，需要跳转到用户的个人中心或者首页。
- **表单提交后跳转**：为了避免表单重复提交，通常会在表单提交成功后进行重定向。
- **外部资源跳转**：需要将请求转发到外部网站，通常需要使用重定向。
  
---

### 5. **总结**

- **转发**：只在服务器内部完成跳转，客户端不会察觉，原始请求的数据会保留，适用于服务器内部的页面或控制器跳转。
- **重定向**：服务器通知客户端重新发起请求，URL 发生变化，适用于客户端可感知的跳转，如表单提交后跳转、登录后跳转等。

在 Spring MVC 中，转发和重定向的设置非常简单，通过 `forward:` 和 `redirect:` 前缀可以轻松实现。这两者的选择应根据实际的业务场景来决定。
````

