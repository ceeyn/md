

1.锁不是锁为了保持某个原子性操作，而是为了锁某种不变性，例如银行转账两个账户总量相等

2.耗时操作例如 rpc 应释放锁，在准备req 和处理 resp 时获取锁

3.很多并发问题都源于先检查-再操作，然而检查的结果具有不可靠性，有可能在你检查完毕到操作的时候，检查的结果已经不可靠

4.不要同时使用两种同步机制保证同步

5.封装使得对程序正确性分析更为可能，破坏约束条件更难









### 3.4 volatile发布不可变类型

**请求进来** → Servlet 读取请求中的数字 `i`。

**检查缓存**：

- 调用 `cache.getFactors(i)`；如果返回 `null`，说明缓存未命中。
- 如果不是 `null`，就说明命中了缓存，直接把这个数组（的副本）写回响应。

**计算因数分解**：

- 如果缓存未命中，则执行 `factor(i)` 进行计算。
- 创建新的不可变缓存对象 `new OneValueCache(i, factors)`。
- 用 `cache = ...` 更新全局 `cache`。

**返回结果**：

- 把因数分解结果通过 `encodeIntoResponse` 写到响应里。

如果 `OneValueCache` 中的字段（如 `lastNumber` 和 `lastFactors`）没有被声明为 `final`，则在多线程环境下可能出现安全发布的问题。具体来说，假设在构造函数中先给 `lastNumber` 赋值，然后再给 `lastFactors` 赋值，但由于重排序或内存可见性问题，另一个线程可能在看到 `lastNumber` 的更新后，却没有看到 `lastFactors` 的更新（例如 `lastFactors` 仍然为 `null` 或不完整），从而导致以下问题：

**final 作用：**不可修改 + 内存屏障：指令不会重排序，另一个线程 B 读取 cache 时，由于 final 字段的初始化屏障，所有字段值已正确可见
// 不会看到lastNumber 初始化但lastFactors未初始化，但OneValueCache已经是新的的情况

【**final 的初始化屏障**： JMM 强制 final 字段的初始化必须在构造函数完成前执行，禁止重排序； 

**volatile 的引用可见性**： 线程 B 通过 volatile 读取到新引用时，所有 final 字段已正确初始化  **只有final 会 读到旧的引用，**】

 

```java
// 3.4 不变性的实践
public class OneValueCache {
    private final BigInteger lastNumber;

    private final BigInteger[] lastFactors;

    public OneValueCache(BigInteger lastNumber, BigInteger[] lastFactors) {
        this.lastNumber = lastNumber;
        this.lastFactors = Arrays.copyOf(lastFactors, lastFactors.length);
    }

    public BigInteger[] getLastFactors(BigInteger i) {
        if (lastNumber == null || !lastNumber.equals(i)) {
            return null;
        } else {
            return Arrays.copyOf(lastFactors, lastFactors.length);
        }
    }

}

// 3.4 不变性的实践
public class VolatileCachedFactorizer implements Servlet {
    private volatile OneValueCache cache = new OneValueCache(null, null);

    @Override
    public void init(ServletConfig servletConfig) throws ServletException {

    }

    @Override
    public ServletConfig getServletConfig() {
        return null;
    }

    @Override
    public void service(ServletRequest servletRequest, ServletResponse servletResponse) throws ServletException, IOException {
        BigInteger number = (BigInteger) servletRequest.getAttribute("number");
        BigInteger[] lastFactors = cache.getLastFactors(number);
        if (lastFactors == null) {
            lastFactors = factors(number);
            cache = new OneValueCache(number, lastFactors);
        }
        encodeIntoResponse(servletResponse, cache);
    }
    private void encodeIntoResponse(ServletResponse servletResponse, OneValueCache cache) {

    }

    @Override
    public String getServletInfo() {
        return null;
    }

    @Override
    public void destroy() {

    }

    private BigInteger[] factors(BigInteger i) {
        return null;
    }
}


```

