https://www.cnblogs.com/tera/p/13976714.html



## [JNI-从jvm源码分析Thread.interrupt的系统级别线程打断原理](https://www.cnblogs.com/tera/p/13976714.html)

2020-11-15 14:50 [tera](https://www.cnblogs.com/tera) 阅读(1346) 评论(2) [编辑](https://i.cnblogs.com/EditPosts.aspx?postid=13976714) [收藏](javascript:void(0)) [举报](javascript:void(0))

### **前言**

在java编程中，我们经常会调用Thread.sleep()方法使得线程停止运行一段时间，而Thread类中也提供了interrupt方法供我们去主动打断一个线程。那么线程挂起和打断的本质究竟是什么，本文就此问题作一个探究。

本文主要分为以下几个部分

1.interrupt的使用特点

2.jvm层面上interrupt方法的本质

3.ParkEvent对象的本质

4.Park()对象的本质

5.利用jni实现一个可以被打断的MyThread类

## 1.interrupt的使用特点

我们先看2个线程打断的示例

首先是可打断的情况：

```java
@Test
public void interruptedTest() throws InterruptedException {
    Thread sleep = new Thread(() -> {
        try {
            log.info("sleep thread start");
            TimeUnit.SECONDS.sleep(1);
            log.info("sleep thread end");
        } catch (InterruptedException e) {
            log.info("sleep thread interrupted");
        }
    }, "sleep_thread");
    sleep.start();

    TimeUnit.MILLISECONDS.sleep(100);
    log.info("ready to interrupt sleep");
    sleep.interrupt();
}
```

我们创建了一个“sleep”线程，其中调用了会抛出InterruptedException异常的sleep方法。“sleep”线程启动100毫秒后，主线程调用其打断方法，此时输出如下：

```armasm
09:50:39.312 [sleep_thread] INFO cn.tera.thread.ThreadTest - sleep thread start
09:50:39.412 [main] INFO cn.tera.thread.ThreadTest - ready to interrupt sleep
09:50:39.412 [sleep_thread] INFO cn.tera.thread.ThreadTest - sleep thread interrupted
```

可以看到“sleep”线程被打断后，抛出了InterruptedException异常，并直接进入了catch的逻辑。

接着我们看一个不可打断的情况：

```java
@Test
public void normalTest() throws InterruptedException {
    Thread normal = new Thread(() -> {
        log.info("normal thread start");
        int i = 0;
        while (true) {
            i++;
        }
    }, "normal_thread");
    normal.start();
    TimeUnit.MILLISECONDS.sleep(100);
    log.info("ready to interrupt normal");
    normal.interrupt();
}
```

我们创建了一个“normal”线程，其中是一个死循环对i++，此时输出如下：

```armasm
10:09:20.237 [normal_thread] INFO cn.tera.thread.ThreadTest - normal thread start
10:09:20.338 [main] INFO cn.tera.thread.ThreadTest - ready to interrupt normal
```

可以看到“normal”线程被打断后，并不会抛出异常，且会继续执行业务流程。

所以打断线程并非是任何时候都会生效的，那么我们就需要探究下interrupt究竟做了什么。

## 2.jvm层面上interrupt方法的本质

#### **Thread.java**

查看interrupt方法，其中的interrupt0()正是打断的主要方法

```java
public void interrupt() {
    if (this != Thread.currentThread())
        checkAccess();

    synchronized (blockerLock) {
        Interruptible b = blocker;
        if (b != null) {
            //打断的主要方法，该方法的主要作用是设置一个打断标记
            interrupt0();
            b.interrupt(this);
            return;
        }
    }
    interrupt0();
}
```

查看interrupt0()方法：

```java
private native void interrupt0();
```

因为interrupt0()是一个本地方法，所以要了解其的究竟做了什么，我们就需要深入到jvm中看源码。其中涉及到了jni相关的知识，有兴趣的同学可以参看我之前写的jni基础应用的文章。
[JNI-从jvm源码分析Thread.start的调用与Thread.run的回调](https://www.cnblogs.com/tera/p/13937611.html)

首先我们还是需要下载open-jdk的源码，包括jdk和hotspot（jvm）

下载地址：http://hg.openjdk.java.net/jdk8

因为C和C++的代码对于java程序员来说比较晦涩难懂，所以在下方展示源码的时候我只会贴出我们关心的重点代码，其余的部分就省略了。

**查看Thread.c：jdk源码目录src/java.base/share/native/libjava**

找到如下代码：

```c++
static JNINativeMethod methods[] = {
    ...
    {"interrupt0",       "()V",        (void *)&JVM_Interrupt}
    ...
};
```

可以看到interrupt0对应的jvm方法是JVM_Interrupt

**查看jvm.cpp，hotspot目录src/share/vm/prims**

可以找到JVM_Interrupt方法的实现，这个方法挺简单的：

```c++
JVM_ENTRY(void, JVM_Interrupt(JNIEnv* env, jobject jthread))
  JVMWrapper("JVM_Interrupt");
  ...
  if (thr != NULL) {
    //执行线程打断操作
    Thread::interrupt(thr);
  }
JVM_END
```

**查看thread.cpp，hotspot目录src/share/vm/runtime**

找到interrupt方法：

```c++
void Thread::interrupt(Thread* thread) {
  //执行os层面的打断
  os::interrupt(thread);
}
```

**查看os_posix.cpp，hotspot目录src/os/posix/vm**

找到interrupt方法，这个方法正是打断的重点：

```c++
void os::interrupt(Thread* thread) {
  ...
  //获得c++线程对应的系统线程
  OSThread* osthread = thread->osthread();
  //如果系统线程的打断标记是false，意味着还未被打断
  if (!osthread->interrupted()) {
    //将系统线程的打断标记设为true
    osthread->set_interrupted(true);
    //这个涉及到内存屏障，本文不展开
    OrderAccess::fence();
    //这里获取一个_SleepEvent，并调用其unpark()方法
    ParkEvent * const slp = thread->_SleepEvent ;
    if (slp != NULL) slp->unpark() ;
  }

  //这里依据JSR166标准，即使打断标记为true，依然要调用下面的2个unpark
  if (thread->is_Java_thread())
    //如果是一个java线程，这里获取一个parker对象，并调用其unpark()方法
    ((JavaThread*)thread)->parker()->unpark();

  ParkEvent * ev = thread->_ParkEvent ;
  //这里获取一个_ParkEvent，并调用其unpark()方法
  if (ev != NULL) ev->unpark() ;
}
```

这个方法中，首先判断线程的打断标志，如果为false，则将其设置为true

并且调用了3个对象的unpark()方法，一会儿介绍着3个对象的作用。

**总而言之，线程打断的本质做了2件事情**

**1.将线程的打断标志设置为true**

**2.调用3个对象的unpark方法唤醒线程**

## 3.ParkEvent对象的本质

在前面我们看到线程在调用interrupt方法的最底层其实是调用了thread中3个对象的unpark()方法，那么这3个对象究竟代表了什么呢，我们继续探究。

首先我们先看**SleepEvent**和**ParkEvent**对象，这2个对象的类型是相同的

**查看thread.cpp，hotspot目录src/share/vm/runtime**

找到SleepEvent和ParkEvent的定义，jvm已经给我们注释了，ParkEven是供synchronized()使用，SleepEvent是供Thread.sleep使用：

```c++
ParkEvent * _ParkEvent;    // for synchronized()
ParkEvent * _SleepEvent;   // for Thread.sleep
```

**查看park.hpp，hotspot目录src/share/vm/runtime**

在头文件中能找到ParkEvent类的定义，继承自**os::PlatformEvent**，是一个和系统相关的的PlatformEvent：

```c++
class ParkEvent : public os::PlatformEvent {
  ...
}
```

**查看os_linux.hpp，hotspot目录src/os/linux/vm**

以linux系统为例，在头文件中可以看到PlatformEvent的具体定义，我们只关注其中的重点：

首先是2个私有对象，一个**pthread_mutex_t操作系统级别的信号量**，一个**pthread_cond_t操作系统级别的条件变量**，这2个变量是一个数组，长度都是1，这些在后面会看到是如何使用的

其次是定义了3个方法，park()、unpark()、park(jlong millis)，控制线程的挂起和继续执行

```c++
class PlatformEvent : public CHeapObj<mtInternal> {
 private:
  ...
  pthread_mutex_t _mutex[1];
  pthread_cond_t  _cond[1];
  ...
  void park();
  void unpark();
  int  park(jlong millis); // relative timed-wait only
  ...
};
```

**查看os_linux.cpp，hotspot目录src/os/linux/vm**

接着我们就需要去看park和unpark方法的具体实现，并看看2个私有变量是如何被使用的

先看**park()**方法，这里我们主要关注3个系统底层方法的调用

**pthread_mutex_lock(_mutex)：锁住信号量**

**status = pthread_cond_wait(_cond, _mutex)：释放信号量，并在条件变量上等待**

**status = pthread_mutex_unlock(_mutex)：释放信号量**

```c++
void os::PlatformEvent::park() { 
    ...
    //锁住信号量
    int status = pthread_mutex_lock(_mutex);
    while (_Event < 0) {
      //释放信号量，并在条件变量上等待
      status = pthread_cond_wait(_cond, _mutex);
    }
    //释放信号量
    status = pthread_mutex_unlock(_mutex);
}
```

这个方法其实非常好理解，就相当于：

```java
synchronize(obj){
  obj.wait();
}
```

或者：

```java
ReentrantLock lock = new ReentrantLock();
Condition condition = lock.newCondition();
lock.lock();
condition.wait();
lock.unlock();
```

park(jlong millis)方法就不展示了，区别只是调用一个接受时间参数的等待方法。

### **所以park()方法底层其实是调用系统层面的锁和条件等待去挂起线程的**

接着我们看unpark()方法，其中最重要的方法当然是

**pthread_cond_signal(_cond)：唤醒条件变量**

```c++
void os::PlatformEvent::unpark() {
  ...
  if (AnyWaiters != 0) {
    //唤醒条件变量
    status = pthread_cond_signal(_cond);
  }
  ...
}
```

### **所以unpark()方法底层其实是调用系统层面的唤醒条件变量达到唤醒线程的目的**

## 4.Park()对象的本质

看完了2个ParkEvent对象的本质，那么接着我们还剩一个park()对象

**查看thread.hpp，hotspot目录src/share/vm/runtime**

park()对象的定义如下：

```c++
public:
  Parker*     parker() { return _parker; }
```

**查看park.hpp，hotspot目录src/share/vm/runtime**

可以看到，它是继承自**os::PlatformParker，和ParkEvent不同**，下面可以看到，等待变量的数组长度变为了2，其中一个给相对时间使用，一个给绝对时间使用

```c++
class Parker : public os::PlatformParker {
    pthread_mutex_t _mutex[1];
    pthread_cond_t  _cond[2]; // one for relative times and one for abs.
}
```

**查看os_linux.cpp，hotspot目录src/os/linux/vm**

还是先看park方法的实现，这个方法其实是对ParkEvent中的park方法的改良版，不过总体的逻辑还是没有变

最终还是调用**pthread_cond_wait方法挂起线程**

```c++
void Parker::park(bool isAbsolute, jlong time) {
  ...
  if (time == 0) {
    //这里是直接长时间等待
    _cur_index = REL_INDEX; 
    status = pthread_cond_wait(&_cond[_cur_index], _mutex);
  } else {
    //这里会根据时间是否是绝对时间，分别等待在不同的条件上
    _cur_index = isAbsolute ? ABS_INDEX : REL_INDEX;
    status = pthread_cond_timedwait(&_cond[_cur_index], _mutex, &absTime);
  }
  ...
}
```

最后看一下unpark方法，这里需要先获取一个正确的等待对象，然后通知即可：

```c++
void Parker::unpark() {
  int status = pthread_mutex_lock(_mutex);
  ...
  //因为在等待的时候会有2个等待对象，所以需要先获取正确的索引
  int index = _cur_index;
  ...
  status = pthread_mutex_unlock(_mutex);
  if (s < 1 && index != -1) {
    //唤醒线程
    status = pthread_cond_signal(&_cond[index]);
  }
  ...
}
```

## 5.利用jni实现一个可以被打断的MyThread类

结合上一篇文章，我们利用jni实现一个自己可以被打断的简易MyThread类

对于jni的基础使用和Thread在jvm级别的本质可以参看上一篇文章，对下面每一步的意义都作了详细的解释
[JNI-从jvm源码分析Thread.start的调用与Thread.run的回调](https://www.cnblogs.com/tera/p/13937611.html)

首先定义**MyThread.java**

```java
import java.util.concurrent.TimeUnit;
import java.time.LocalDateTime;

public class MyThread {

    static {
        //设置查找路径为当前项目路径
        System.setProperty("java.library.path", ".");
        //加载动态库的名称
        System.loadLibrary("MyThread");
    }

    public native void startAndPark();

    public native void interrupt();

    public static void main(String[] args) throws InterruptedException {
        MyThread thread = new MyThread();
        //启动线程打印一段文字，并睡眠
        thread.startAndPark();
        //1秒后主线程打断子线程
        TimeUnit.MILLISECONDS.sleep(1000);
        System.out.println(LocalDateTime.now() + "：Main---准备打断线程");
        //打断子线程
        thread.interrupt();
        System.out.println(LocalDateTime.now() + "：Main---打断完成");
    }
}
```

执行命令编译MyThread.class文件并生成MyThread.h头文件

```mipsasm
javac -h . MyThread.java
```

创建**MyThread.c**文件

当java代码调用startAndPark()方法的时候，创建了一个系统级别的线程，并调用pthread_cond_wait进行休眠

当java代码调用interrupt()方法的时候，会唤醒休眠中的线程

```c++
#include <pthread.h>
#include <stdio.h>
#include "MyThread.h"
#include "time.h"

pthread_t pid;
pthread_mutex_t _mutex = PTHREAD_MUTEX_INITIALIZER;
pthread_cond_t  _cond = PTHREAD_COND_INITIALIZER; 

//打印时间
void printTime(){
    char strTm[50] = { 0 };
	  time_t currentTm;
	  time(&currentTm);
	  strftime(strTm, sizeof(strTm), "%x %X", localtime(&currentTm));
	  puts(strTm);
}

//子线程执行的方法
void* thread_entity(void* arg){
    printTime();
    printf("MyThread---启动\n");
    printTime();
    printf("MyThread---准备休眠\n");
    //阻塞线程，等待唤醒
    pthread_cond_wait(&_cond, &_mutex);
    printTime();
    printf("MyThread---休眠被打断\n");
}
//对应MyThread中的startAndPark方法
JNIEXPORT void JNICALL Java_MyThread_startAndPark(JNIEnv *env, jobject c1){
    //创建一个子线程
    pthread_create(&pid, NULL, thread_entity, NULL);
}
//对应MyThread中的interrupt方法
JNIEXPORT void JNICALL Java_MyThread_interrupt(JNIEnv *env, jobject c1){
    //唤醒线程
    pthread_cond_signal(&_cond);
}
```

执行命令创建动态链接库

```bash
gcc -dynamiclib -I /Library/Java/JavaVirtualMachines/jdk1.8.0_241.jdk/Contents/Home/include MyThread.c -o libMyThread.jnilib 
```

执行java的main方法，得到结果

子线程启动后进入睡眠，主线程1秒钟后打断子线程，完全符合我们的预期

```lua
2020/11/13 19时42分57秒
MyThread---启动
2020/11/13 19时42分57秒
MyThread---准备休眠
2020-11-13T19:42:58.891：Main---准备打断线程
2020/11/13 19时42分58秒
MyThread---休眠被打断
2020-11-13T19:42:58.891：Main---打断完成
```

最后总结一下本文的内容

1.线程打断的本质做了2件事情：设置线程的打断标记，并调用线程3个Park对象的unpark()方法唤醒线程

2.线程挂起的本质是调用系统级别的**pthread_cond_wait**方法，使得等待在一个条件变量上

3.线程唤醒的本质是调用系统级别的**pthread_cond_signal**方法，唤醒等待的线程

4.通过实现一个自己的可以打断的线程类更好地理解线程打断的本质

