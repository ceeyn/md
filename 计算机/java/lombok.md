# 前言

`Lombok`是一款Java开发插件，使得Java开发者可以通过其定义的一些注解来消除业务工程中冗长和繁琐的代码，尤其对于简单的Java模型对象（[POJO](https://so.csdn.net/so/search?q=POJO&spm=1001.2101.3001.7020)）。在开发环境中使用Lombok插件后，Java开发人员可以节省出重复构建，诸如hashCode和equals这样的方法以及各种业务对象模型的accessor和ToString等方法的大量时间。对于这些方法，它能够在编译源代码期间自动帮我们生成这些方法，并没有如反射那样降低程序的性能。

它所有的增强都是通过注解实现，所以了解其使用主要了解一下注解即可

注解列表

Maven下的依赖

```xml
<dependency>
    <groupId>org.projectlombok</groupId>
    <artifactId>lombok</artifactId>
    <scope>provided</scope>
    <version>1.18.4</version>
</dependency>
123456
```

Gradle下 的依赖(kts)

```kotlin
dependencies {
    //*Lombok*
    // <p>注解处理将不再compile classpath中，需要手动添加到annotation processor path</p>
    compileOnly("org.projectlombok:lombok:1.18.4")
    annotationProcessor("org.projectlombok:lombok:1.18.4")
    testCompileOnly("org.projectlombok:lombok:1.18.4")
    testAnnotationProcessor("org.projectlombok:lombok:1.18.4")
}
12345678
```

先介绍这一波最常用的注解：

# @NoArgsConstructor/@RequiredArgsConstructor/@AllArgsConstructor

这三个注解都是用在类上的，第一个和第三个都很好理解，就是为该类产生无参的构造方法和包含所有参数的构造方法，第二个注解则使用类中所有带有`@NonNull`注解的或者带有`final`修饰的成员变量生成对应的构造方法，当然，和前面几个注解一样，成员变量都是非静态的。

另外，**如果类中含有final修饰的成员变量**，是无法使用`@NoArgsConstructor`注解的。

三个注解都可以指定生成的构造方法的访问权限，还可指定生成一个静态方法

使用案例：

```java
@AllArgsConstructor
public class Demo {
    private String name;
    private int age;
}

@AllArgsConstructor
class Parent {
    private Integer id;
}
12345678910
```

编译后的两个class文件如下：

```

```



```java
public class Demo {
    private String name;
    private int age;

    public Demo(String name, int age) {
        this.name = name;
        this.age = age;
    }
}

//第二个类
class Parent {
    private Integer id;

    public Parent(Integer id) {
        this.id = id;
    }
}
123456789101112131415161718
```

由此课件，此注解并不会把父类的属性id拿到Demo的构造器里面去，这是需要注意的地方。并且它也没有默认的构造器了

```java
@AllArgsConstructor(access = AccessLevel.PROTECTED, staticName = "test")
public class Demo {
    private final int finalVal = 10;
    private String name;
    private int age;
}
123456
```

生成如下：

```java
public class Demo {
    private final int finalVal = 10;
    private String name;
    private int age;

    private Demo(String name, int age) {
        this.name = name;
        this.age = age;
    }

    protected static Demo test(String name, int age) {
        return new Demo(name, age);
    }
}
1234567891011121314
```

看出来的效果为：可以指定生成的构造器的访问权限。但是，但是如果指定了一个静态方法，那么构造器会自动会被private，只通过静态方法对外提供反问，并且我们发现final的属性值，是不会放进构造函数里面的。

NoArgsConstructor的使用方式同上，RequiredArgsConstructor看看效果：

```java
@RequiredArgsConstructor
public class Demo {
    private final int finalVal = 10;

    @NonNull
    private String name;
    @NonNull
    private int age;
}
123456789
```

编译后：

```java
public class Demo {
    private final int finalVal = 10;
    @NonNull
    private String name;
    @NonNull
    private int age;

    public Demo(@NonNull String name, @NonNull int age) {
        if (name == null) {
            throw new NullPointerException("name is marked @NonNull but is null");
        } else {
            this.name = name;
            this.age = age;
        }
    }
}
12345678910111213141516
```

解释：该注解会识别`@nonNull`字段，然后以该字段为元素产生一个构造函数。备注：如果所有字段都没有@nonNull注解，那效果同`NoArgsConstructor`

# @Builder

`@Builder`提供了一种比较推崇的构建值对象的方式
非常推荐的一种构建值对象的方式。缺点就是父类的属性不能产于builder

```java
@Builder
public class Demo {
    private final int finalVal = 10;

    private String name;
    private int age;
}
1234567
```

编译后：

```java
public class Demo {
    private final int finalVal = 10;
    private String name;
    private int age;

    Demo(String name, int age) {
        this.name = name;
        this.age = age;
    }

    public static Demo.DemoBuilder builder() {
        return new Demo.DemoBuilder();
    }

    public static class DemoBuilder {
        private String name;
        private int age;

        DemoBuilder() {
        }

        public Demo.DemoBuilder name(String name) {
            this.name = name;
            return this;
        }

        public Demo.DemoBuilder age(int age) {
            this.age = age;
            return this;
        }

        public Demo build() {
            return new Demo(this.name, this.age);
        }

        public String toString() {
            String var10000 = this.name;
            return this.age;
        }
    }
}
1234567891011121314151617181920212223242526272829303132333435363738394041
```

因此我们构造一个对象就可以优雅的这么来：

```java
 public static void main(String[] args) {
        Demo demo = Demo.builder().name("aa").age(10).build();
        System.out.println(demo); 
    }
1234
```

里面有一些自定义参数，我表示，完全没有必要去自定义。

# @Cleanup

```
@Cleanup`能够`自动释放资源
```

这个注解用在变量前面，可以保证此变量代表的资源会被自动关闭，默认是调用资源的close()方法。如果该资源有其它关闭方法，可使用@Cleanup(“methodName”)来指定要调用的方法，就用输入输出流来举个例子吧：

```java
public static void main(String[] args) throws Exception {
        @Cleanup InputStream in = new FileInputStream(args[0]);
        @Cleanup OutputStream out = new FileOutputStream(args[1]);
        byte[] b = new byte[1024];
        while (true) {
            int r = in.read(b);
            if (r == -1) break;
            out.write(b, 0, r);
        }
    }
12345678910
```

编译后：

```java
public static void main(String[] args) throws Exception {
        FileInputStream in = new FileInputStream(args[0]);

        try {
            FileOutputStream out = new FileOutputStream(args[1]);

            try {
                byte[] b = new byte[1024];

                while(true) {
                    int r = in.read(b);
                    if (r == -1) {
                        return;
                    }

                    out.write(b, 0, r);
                }
            } finally {
                if (Collections.singletonList(out).get(0) != null) {
                    out.close();
                }

            }
        } finally {
            if (Collections.singletonList(in).get(0) != null) {
                in.close();
            }

        }
    }
123456789101112131415161718192021222324252627282930
```

就这么简单的一个注解，就实现了优雅的关流操作哟。

# @Data

`@Data` 强悍的组合功能包

相当于注解集合。效果等同于**@Getter + @Setter + @ToString + @EqualsAndHashCode + @RequiredArgsConstructor** 由于生成的代码篇幅太长，这里就不给demo了，反正效果同上5个注解的效果，强悍

需要注意的是，这里**不包括@NoArgsConstructor和@AllArgsConstructor**

# @Value

`@Value`注解和@Data类似，区别在于它会把所有成员变量默认定义为`private final`修饰，并且**不会生成set方法**。

所以@Value更适合`只读性`更强的类，所以特殊情况下，还是可以使用的。

# @ToString/@EqualsAndHashCode

这两个注解也比较好理解，就是生成`toString`，`equals`和`hashcode`方法，同时后者还会生成一个canEqual方法，用于判断某个对象是否是当前类的实例。，生成方法时只会使用类中的非静态成员变量，这些都比较好理解。毕竟静态的东西并不属于对象本身

```java
@ToString
public class Demo {
    private final int finalVal = 10;

    private transient String name = "aa";
    private int age;

}


 public static void main(String[] args) throws Exception {
        Demo demo = new Demo();
        System.out.println(demo); //Demo(finalVal=10, age=0)
    }
1234567891011121314
```

我们发现静态字段它是不输出的。
有些关键的属性，可以控制toString的输出，我们可以了解一下：

```java
@ToString(
        includeFieldNames = true, //是否使用字段名
        exclude = {"name"}, //排除某些字段
        of = {"age"}, //只使用某些字段
        callSuper = true //是否让父类字段也参与 默认false
)
123456
```

> 备注：大多数情况下，使用默认的即可，毕竟大多数情况都是POJO

`@Generated`：暂时貌似没什么用

# @Getter/@Setter

这一对注解从名字上就很好理解，用在成员变量上面或者类上面，相当于为成员变量生成对应的`get`和`set`方法，同时还可以为生成的方法指定访问修饰符，当然，`默认为public`

这两个注解直接用在类上，可以为此类里的所有非静态成员变量生成对应的get和set方法。**如果是final变量，那就只会有get方法**

```java
@Getter
@Setter
public class Demo {
    private final int finalVal = 10;

    private String name;
    private int age;

}
123456789
```

编译后：

```java
public class Demo {
    private final int finalVal = 10;
    private String name;
    private int age;

    public Demo() {
    }

    public int getFinalVal() {
        Objects.requireNonNull(this);
        return 10;
    }

    public String getName() {
        return this.name;
    }

    public int getAge() {
        return this.age;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void setAge(int age) {
        this.age = age;
    }
}
1234567891011121314151617181920212223242526272829
```

# @NonNull

这个注解可以用在成员方法或者构造方法的参数前面，会自动产生一个关于此参数的`非空检查`，如果参数为空，则抛出一个空指针异常。

```java
//成员方法参数加上@NonNull注解
public String getName(@NonNull Person p){
    return p.getName();
}
1234
```

编译后：

```java
public String getName(@NonNull Person p){
    if(p==null){
        throw new NullPointerException("person");
    }
    return p.getName();
}
123456
```

`@Singular` 默认值 暂时也没太大用处

# @SneakyThrows

这个注解用在方法上，可以将方法中的代码用`try-catch`语句包裹起来，捕获异常并在catch中用[Lombok](https://so.csdn.net/so/search?q=Lombok&spm=1001.2101.3001.7020).sneakyThrow(e)把异常抛出，可以使用@SneakyThrows(Exception.class)的形式指定抛出哪种异常

```java
 @SneakyThrows(UnsupportedEncodingException.class)
    public String utf8ToString(byte[] bytes) {
        return new String(bytes, "UTF-8");
    }
1234
```

编译后：

```java
@SneakyThrows(UnsupportedEncodingException.class)
    public String utf8ToString(byte[] bytes) {
        try{
            return new String(bytes, "UTF-8");
        }catch(UnsupportedEncodingException uee){
            throw Lombok.sneakyThrow(uee);
        }
    }
12345678
```

这里有必要贴出来Lombok.sneakyThrow的代码：

```java
 public static RuntimeException sneakyThrow(Throwable t) {
        if (t == null) {
            throw new NullPointerException("t");
        } else {
            return (RuntimeException)sneakyThrow0(t);
        }
    }

    private static <T extends Throwable> T sneakyThrow0(Throwable t) throws T {
        throw t;
    }
1234567891011
```

其实就是转化为了`RuntimeException`，其实我想说，这个注解也没大用。毕竟我们碰到异常，是希望自己处理的

# @Synchronized

这个注解用在类方法或者实例方法上，效果和`synchronized`关键字相同，区别在于`锁对象不同`，对于类方法和实例方法，synchronized关键字的锁对象分别是类的class对象和this对象，而@Synchronized得锁对象分别是私有静态final对象LOCK`和`私有final对象lock`，当然，也可以自己指定锁对象

```java
@Synchronized
    public static void hello() {
        System.out.println("world");
    }

    @Synchronized
    public int answerToLife() {
        return 42;
    }

    @Synchronized("readLock")
    public void foo() {
        System.out.println("bar");
    }
1234567891011121314
```

编译后：

```java
public static void hello() {
        Object var0 = $LOCK;
        synchronized($LOCK) {
            System.out.println("world");
        }
    }

    public int answerToLife() {
        Object var1 = this.$lock;
        synchronized(this.$lock) {
            return 42;
        }
    }

    public void foo() {
        Object var1 = this.readLock;
        synchronized(this.readLock) {
            System.out.println("bar");
        }
    }
1234567891011121314151617181920
```

我只能说，这个注解也挺鸡肋的。

# @Val

@Val 很强的类型推断 var注解，在Java10之后就不能使用了

```java
class Parent {
    //private static final val set = new HashSet<String>(); //编译不通过

    public static void main(String[] args) {
        val set = new HashSet<String>();
        set.add("aa");
        System.out.println(set); //[aa]
    }

}
12345678910
```

编译后：

```java
class Parent {
    Parent() {
    }

    public static void main(String[] args) {
        HashSet<String> set = new HashSet();
        set.add("aa");
        System.out.println(set);
    }
}
12345678910
```

这个和Java10里的Var很像，强大的类型推断。并且不能使用在全局变量上，只能使用在局部变量的定义中。

# @Log、CommonsLog、Slf4j、XSlf4j、Log4j、Log4j2等日志注解

这个注解用在类上，可以省去从日志工厂生成日志对象这一步，直接进行日志记录，具体注解根据日志工具的不同而不同，同时，可以在注解中使用topic来指定生成log对象时的类名。不同的日志注解总结如下(上面是注解，下面是实际作用)：

```java
@CommonsLog
private static final org.apache.commons.logging.Log log = org.apache.commons.logging.LogFactory.getLog(LogExample.class);
@JBossLog
private static final org.jboss.logging.Logger log = org.jboss.logging.Logger.getLogger(LogExample.class);
@Log
private static final java.util.logging.Logger log = java.util.logging.Logger.getLogger(LogExample.class.getName());
@Log4j
private static final org.apache.log4j.Logger log = org.apache.log4j.Logger.getLogger(LogExample.class);
@Log4j2
private static final org.apache.logging.log4j.Logger log = org.apache.logging.log4j.LogManager.getLogger(LogExample.class);
@Slf4j
private static final org.slf4j.Logger log = org.slf4j.LoggerFactory.getLogger(LogExample.class);
@XSlf4j
private static final org.slf4j.ext.XLogger log = org.slf4j.ext.XLoggerFactory.getXLogger(LogExample.class);
1234567891011121314
```

这个注解还是非常有用的，特别是Slf4j这个，在平时开发中挺有用的

```java
@Slf4j
class Parent {
}
123
```

编译后：

```java
class Parent {
    private static final Logger log = LoggerFactory.getLogger(Parent.class);

    Parent() {
    }
}
123456
```

也可topic的名称：

```java
@Slf4j
@CommonsLog(topic = "commonLog")
class Parent {
}
1234
```

编译后：

```java
class Parent {
    private static final Logger log = LoggerFactory.getLogger("commonLog");

    Parent() {
    }
}
123456
```

lombok中有experimental的包：
实验性因为：

我们可能想将这些特性和更完全的性质支持概念融为一体(普通话：这些性能还在研究)
新特性-需要社区反馈

# @Accessors

`@Accessors` 一个为`getter`和`setter`设计的更流畅的注解
这个注解要搭配`@Getter`与`@Setter`使用，用来修改默认的setter与getter方法的形式。所以单独使用是没有意义的

```java
@Accessors(fluent = true)
@Getter
@Setter
public class Demo extends Parent {
    private final int finalVal = 10;

    private String name;
    private int age;

}
12345678910
```

编译后：

```java
public class Demo extends Parent {
    private final int finalVal = 10;
    private String name;
    private int age;

    public Demo() {
    }

    public int finalVal() {
        Objects.requireNonNull(this);
        return 10;
    }

    public String name() {
        return this.name;
    }

    public int age() {
        return this.age;
    }

    public Demo name(String name) {
        this.name = name;
        return this;
    }

    public Demo age(int age) {
        this.age = age;
        return this;
    }
}
12345678910111213141516171819202122232425262728293031
```

它的三个参数解释：

- `chain` 链式的形式 这个特别好用，方法连缀越来越方便了
- `fluent` 流式的形式（若无显示指定chain的值，也会把chain设置为true）
- `prefix` 生成指定前缀的属性的getter与setter方法，并且生成的getter与setter方法时会去除前缀

```java
@Accessors(prefix = "xxx")
@Getter
@Setter
public class Demo extends Parent {
    private final int finalVal = 10;

    private String xxxName;
    private int age;

}
12345678910
```

编译后：

```java
public class Demo extends Parent {
    private final int finalVal = 10;
    private String xxxName;
    private int age;

    public Demo() {
    }

    public String getName() {
        return this.xxxName;
    }

    public void setName(String xxxName) {
        this.xxxName = xxxName;
    }
}
12345678910111213141516
```

我们发现prefix可以在生成get/set的时候，去掉xxx等prefix前缀，达到很好的一致性。但是，但是需要注意，因为此处age没有匹配上xxx前缀，所有根本就不给生成，所以使用的时候一定要注意。

属性名没有一个以其中的一个前缀开头，则属性会被lombok完全忽略掉，并且会产生一个警告。

# @Delegate

`@Delegate`注释的属性，会把这个属性对象的公有非静态方法合到当前类

代理模式，把字段的方法代理给类，默认代理所有方法。注意：公共 非静态方法

```java
public class Demo extends Parent {
    private final int finalVal = 10;

    @Delegate
    private String xxxName;
    private int age;

}
12345678
```

编译后：把String类的公共 非静态方法全拿来了 个人觉得很鸡肋有木有

```java
public class Demo extends Parent {
    private final int finalVal = 10;
    private String xxxName;
    private int age;

    public Demo() {
    }

    public int length() {
        return this.xxxName.length();
    }

    public boolean isEmpty() {
        return this.xxxName.isEmpty();
    }

    public char charAt(int index) {
        return this.xxxName.charAt(index);
    }

    public int codePointAt(int index) {
        return this.xxxName.codePointAt(index);
    }
1234567891011121314151617181920212223
```

> 备注：它不能用于基本数据类型字段比如int，只能用在包装类型比如Integer

参数们：

- types：指定代理的方法
- excludes：和types相反

@NonFinal 设置不为Final，@FieldDefaults和@Value也有这功能
@SuperBuilder 本以为它是支持到了父类属性的builder构建，但其实，我们还是等等吧 目前还不好使
@UtilityClass 工具类 会把所有字段方法static掉，没啥用
@Wither 生成withXXX方法，返回类实例 没啥用，因为还有bug

# @Builder和@NoArgsConstructor一起使用冲突问题

当我们这么使用时候：

编译报错：

```bash
Error:(17, 1) java: 无法将类 com.sayabc.groupclass.dtos.appoint.TeaPoolLogicalDelDto中的构造器 TeaPoolLogicalDelDto应用到给定类型;
1
  需要: 没有参数
  找到: java.lang.Long,java.lang.Long,java.lang.Long,java.lang.Integer
  原因: 实际参数列表和形式参数列表长度不同
123
```

其实原因很简单，自己点进去看编译后的源码一看便知。
只使用`@Builder`会自动创建全参构造器。而添加上`@NoArgsConstructor`后就不会自动产生全参构造器

两种解决方式：

- 去掉@NoArgsConstructor
- 添加@AllArgsConstructor（建议使用这种，毕竟无参构造最好保证是有的）

but，枚举值建议这样来就行了，不要加@NoArgsConstructor

我认为这也是Lombok的一个bug，希望在后续版本中能够修复

@builder注解影响设置默认值的问题
例子如下，本来我是想给age字段直接赋一个默认值的：
没有使用lombok，我们这么写：

```java
  public static void main(String[] args) {
        Demo demo = new Demo();
        System.out.println(demo); //Demo{id=null, age=10}
    }

    private static class Demo {
        private Integer id;
        private Integer age = 10; //放置默认值年龄

        //省略手动书写的get、set、方法和toString方法

        @Override
        public String toString() {
            return "Demo{" +
                    "id=" + id +
                    ", age=" + age +
                    '}';
        }
    }
12345678910111213141516171819
```

我们发现，这样运行没有问题，默认值也生效了。但是，但是我们用了强大的lombok，我们怎么可能还愿意手写get/set呢？关键是，我们一般情况下还会用到它的@buider注解：

```java
 public static void main(String[] args) {
        Demo demo = new Demo();
        System.out.println(demo); //Demo{id=null, age=10}

        //采用builder构建  这是我们使用最多的场景吧
        Demo demo2 = Demo.builder().build();
        System.out.println(demo2); //PeriodAddReq.Demo(id=null, age=null)
    }

    @Getter
    @Setter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    @ToString
    private static class Demo {
        private Integer id;
        private Integer age = 10; //放置默认值年龄
    }
12345678910111213141516171819
```

代码简洁了不少。但是我们却发现一个问题。new出来的对象默认值仍然没有问题，但是buider构建出来的demo2对象，默认值却没有设置进去。这是一个非常隐晦的问题，一不小心，就可能留下一个惊天大坑，所以需要注意
其实在执行编译的时候，idea开发工具已经警告我们了：

```
Warning:(51, 25) java: @Builder will ignore the initializing expression entirely. If you want the initializing expression to serve as default, add @Builder.Default. If it is not supposed to be settable during building, make the field final.
1
```

## 方案一：

从它的建议可以看出，把字段标为final就ok了（亲测好用）。但很显然，绝大多数我们并不希望他是final的字段。

因此我们采用第二个方案：

```java
@Getter
    @Setter
    @Builder
    @NoArgsConstructor
    @AllArgsConstructor
    @ToString
    private static class Demo {
        private Integer id;
        @Builder.Default
        private Integer age = 10; //放置默认值年龄
    }
1234567891011
```

lombok考虑到了这种现象，因此我们只需要在需要设置默认值的字段上面加上 @Builder.Default注解就ok了。

```java
  public static void main(String[] args) {
        Demo demo = new Demo();
        System.out.println(demo); //PeriodAddReq.Demo(id=null, age=null)

        //采用builder构建  这是我们使用最多的场景吧
        Demo demo2 = Demo.builder().build();
        System.out.println(demo2); //PeriodAddReq.Demo(id=null, age=10)
    }
12345678
```

但是我们坑爹的发现：builder默认值没问题了，但是new出来又有问题了。见鬼啊，

我认为这是lombok的一个大bug，希望后续版本中能够修复

但是我们不能因为有这么一个问题，咱们就不使用它了。本文主要提醒读者，在使用的时候留心这个问题即可。

> 备注：@Builder.Default会使得使用@NoArgsConstructor生成的无参构造没有默认值，自己显示写出来的也不会给你设置默认值的，需要注意。

2019年1.18日补充内容：Lombok 1.18.4版本
上面已经指出了Lombok设置默认值的bug，果不其然。官方在1.18.4这个版本修复了这个bug。各位要有版本意识：这个版本级以上版本是好用的，比这版本低的都不行。

用这个版本运行上面例子，默认值没有问题了。

```
Main.Demo(id=null, age=10)
Main.Demo(id=null, age=10)
12
```

我们不用自动生成空构造，显示书写出来呢？如下：

```java
    @Getter
    @Setter
    @Builder
    @AllArgsConstructor
    @ToString
    private static class Demo {
        private Integer id;
        @Builder.Default
        private Integer age = 10; //放置默认值年龄Default

        //显示书写出空构造
        public Demo() {
        }
    }
1234567891011121314
```

我们发现手动书写出来的空构造，默认值是不生效的。这点需要特别注意。

这个就不说是Lombok的bug了，因为既然你都使用Lombok了，为何还自己写空构造呢？不是作死吗？

# Lombok背后的自定义注解原理

作为一个Java开发者来说光了解插件或者技术框架的用法只是做到了“知其然而不知其所以然”，如果真正掌握其背后的技术原理，看明白源码设计理念才能真正做到“知其然知其所以然”。好了，话不多说下面进入本章节的正题，看下Lombok背后注解的深入原理。

可能熟悉Java自定义注解的同学已经猜到，Lombok这款插件正是依靠可插件化的Java自定义注解处理API（JSR 269: Pluggable Annotation Processing API）来实现在Javac编译阶段利用“Annotation Processor”对自定义的注解进行预处理后生成真正在JVM上面执行的“Class文件”。有兴趣的同学反编译带有Lombok注解的类文件也就一目了然了。其大致执行原理图如下：

从上面的Lombok执行的流程图中可以看出，在Javac 解析成AST抽象语法树之后, Lombok 根据自己编写的注解处理器，动态地修改 AST，增加新的节点（即Lombok自定义注解所需要生成的代码），最终通过分析生成JVM可执行的字节码Class文件。使用Annotation Processing自定义注解是在编译阶段进行修改，而JDK的[反射](https://so.csdn.net/so/search?q=反射&spm=1001.2101.3001.7020)技术是在运行时动态修改，两者相比，反射虽然更加灵活一些但是带来的性能损耗更加大。

需要更加深入理解Lombok插件的细节，自己查阅其源代码是必比可少的。

`AnnotationProcessor`这个类是Lombok自定义注解处理的入口。该类有两个比较重要的方法一个是init方法，另外一个是process方法。在init方法中，先用来做参数的初始化，将AnnotationProcessor类中定义的内部类（JavacDescriptor、EcjDescriptor）先注册到`ProcessorDescriptor`类型定义的列表中。其中，内部静态类—JavacDescriptor在其加载的时候就将 `lombok.javac.apt.LombokProcessor`这个类进行对象实例化并注册。在 LombokProcessor处理器中，其中的process方法会根据优先级来分别运行相应的handler处理类。Lombok中的多个自定义注解都分别有对应的handler处理类.

在Lombok中对于其自定义注解进行实际的替换、修改和处理的正是这些handler类。对于其实现的细节可以具体参考其中的代码。