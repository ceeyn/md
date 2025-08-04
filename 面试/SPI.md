

https://www.infoq.cn/article/4JPR7KBVAHKSEpOPDRht?utm_campaign=geek_search&utm_content=geek_search&utm_medium=geek_search&utm_source=geek_search&utm_term=geek_search

https://www.infoq.cn/article/iaQ18w7uvbHavtg5cA42?utm_campaign=geek_search&utm_content=geek_search&utm_medium=geek_search&utm_source=geek_search&utm_term=geek_search

https://xie.infoq.cn/article/69036229e43b3d811e2f2cae6?utm_campaign=geek_search&utm_content=geek_search&utm_medium=geek_search&utm_source=geek_search&utm_term=geek_search



https://blog.csdn.net/xmtblog/article/details/118947643





![双亲委派](/Users/haozhipeng/Downloads/我的笔记/images/双亲委派.png)



双亲委派模型缺点：父类无
法加载子类加载器加载路径
下的类，在spi服务时需要打
破双亲委派

jar包的meta-inf的service下以接口名称定义的文件，
里面有很多实现类，spi就是根据接口查找服务实现

driverManager执行class.forName(驱动)加载驱动类，但

driverManager是bootstrap加载器加载的，父加载器

加载不了子加载器的目录下的类，只能通过线程上下文

加载器【线程私有，一开始赋值为app，类似于thread

local】加载