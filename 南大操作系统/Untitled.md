Lec1



<img src="/Users/moon/Library/Application Support/typora-user-images/image-20250510111411656.png" alt="image-20250510111411656" style="zoom:33%;" />

不能让程序访问硬件，如果可以直接用程序访问硬件的话，电脑随时有 crash的风险，必须通过程序申请操作系统特权进入内核态来访问硬件。







调试，找 bug

1.第一个错的地方往往是一切问题的开始，log n 二分查找找到第一个错的地方

2.需求-》设计-》代码-》执行【error 是不被观察到的错：状态机中间某个状态错误但最后结果正确，failure 是被观察到的错（可以用断言将 failure 尽可能提前）】



