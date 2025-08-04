1、直接安装

dmp包直接安装两个版本的jdk，比如我这里安装的1.8.0_181 、10.0.2。

2、配置环境

打开[环境变量](https://so.csdn.net/so/search?q=环境变量&spm=1001.2101.3001.7020)配置文件bash_profile

[vim](https://so.csdn.net/so/search?q=vim&spm=1001.2101.3001.7020) ~/.bash_profile 

```bash
# java



export JAVA_8_HOME="/Library/Java/JavaVirtualMachines/jdk1.8.0_181.jdk/Contents/Home"



export JAVA_10_HOME="/Library/Java/JavaVirtualMachines/jdk-10.0.2.jdk/Contents/Home"



 



alias jdk8="export JAVA_HOME=$JAVA_8_HOME"



alias jdk10="export JAVA_HOME=$JAVA_10_HOME"



 



export JAVA_HOME=$JAVA_8_HOME



export PATH="$JAVA_HOME:$PATH"
```


source ~/.bash_profile


3、任意切换java环境

控制台输入jdk8或者jdk10就会自动切换

```ruby
JeandeMBP:~ healerjean$ jdk8



JeandeMBP:~ healerjean$ java -version



java version "1.8.0_181"



Java(TM) SE Runtime Environment (build 1.8.0_181-b13)



Java HotSpot(TM) 64-Bit Server VM (build 25.181-b13, mixed mode)



 



JeandeMBP:~ healerjean$ jdk10



JeandeMBP:~ healerjean$ java -version



java version "10.0.2" 2018-07-17



Java(TM) SE Runtime Environment 18.3 (build 10.0.2+13)



Java HotSpot(TM) 64-Bit Server VM 18.3 (build 10.0.2+13, mixed mode)



 
```


4、删除jdk

输入 
sudo rm -fr /Library/Internet\ Plug-Ins/JavaAppletPlugin.plugin 
sudo rm -fr /Library/PreferencesPanes/JavaControlPanel.prefpane


查找当前版本 
输入：ls /Library/Java/JavaVirtualMachines/ 
输出：jdk1.8.0_181.jdk

sudo rm -rf /Library/Java/JavaVirtualMachines/jdk1.8.0_181.jdk


原文：https://blog.csdn.net/u012954706/article/details/82982667 
 