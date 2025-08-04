# **1、查询ip**

[github](https://so.csdn.net/so/search?q=github&spm=1001.2101.3001.7020)网址查询：[The world's leading software development platform · GitHub](https://link.zhihu.com/?target=https%3A//github.com.ipaddress.com/)

![img](https://img-blog.csdnimg.cn/b5519bce1ed04ffbbd80fdd1f96f730b.png?x-oss-process=image/watermark,type_d3F5LXplbmhlaQ,shadow_50,text_Q1NETiBA5aSp5L2R5pyo5p6r,size_20,color_FFFFFF,t_70,g_se,x_16)

 

github域名查询：[▷ github.global.ssl.fastly.net Website statistics and traffic analysis](https://link.zhihu.com/?target=https%3A//fastly.net.ipaddress.com/github.global.ssl.fastly.net)

![img](https://img-blog.csdnimg.cn/d23bd5e0979f4e32a6f937932b52a397.png?x-oss-process=image/watermark,type_d3F5LXplbmhlaQ,shadow_50,text_Q1NETiBA5aSp5L2R5pyo5p6r,size_20,color_FFFFFF,t_70,g_se,x_16)

 

github静态资源ip：[assets-cdn.Github.com ᐅ Learn more about Github](https://ipaddress.com/website/assets-cdn.github.com)

![img](https://img-blog.csdnimg.cn/3282508d2c934881a0b8386c101634b5.png?x-oss-process=image/watermark,type_d3F5LXplbmhlaQ,shadow_50,text_Q1NETiBA5aSp5L2R5pyo5p6r,size_20,color_FFFFFF,t_70,g_se,x_16)





将上面的ip地址都记下来，格式：

```python
140.82.114.3  github.com



199.232.69.194  github.global.ssl.fastly.net



185.199.111.153 assets-cdn.github.com 



185.199.110.153 assets-cdn.github.com 



185.199.108.153 assets-cdn.github.com 
```

# **2、找hosts文件**

1.新建一个访达窗口，同时按住shift command G三个键，进入前往文件夹页面

2.在输入框内输入/etc/hosts 

3.找到hosts文件夹

4.由于hosts文件夹不可编辑，所以复制一份hosts文件先保存到本地桌面

5.在新的hosts文件夹里输入上面的ip格式，然后点击保存。

6.将/etc/hosts原来的文件删除，删除的时候需要输入你的电脑开机密码

7.再将修改后的保存到桌面的hosts文件拖拽到/etc文件夹下，也需要你输入开机密码，文件名是hosts，不要修改

9.然后在浏览器中刷新一下之前没有打开的GitHub页面，或者重新打开

至此，mac电脑无法访问github的问题终于得到了解决