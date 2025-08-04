![截屏2022-11-24 10.23.32](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 10.23.32.png)

![截屏2022-11-24 10.46.28](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 10.46.28.png)

![截屏2022-11-24 10.48.34](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 10.48.34.png)

opt:Optional application software packages

etc：Editable Text Configuration

. 当前目录

pwd 显示当前路径

![截屏2022-11-24 11.05.52](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 11.05.52.png)

![截屏2022-11-24 11.11.10](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 11.11.10.png)

![截屏2022-11-24 11.16.01](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 11.16.01.png)

![截屏2022-11-24 11.20.24](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 11.20.24.png)

![截屏2022-11-24 11.24.29](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 11.24.29.png)

![截屏2022-11-24 11.26.21](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 11.26.21.png)

![截屏2022-11-24 11.32.59](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 11.32.59.png)

![截屏2022-11-24 15.20.45](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 15.20.45.png)

![截屏2022-11-24 16.00.25](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 16.00.25.png)

![截屏2022-11-24 16.10.31](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 16.10.31.png)

![截屏2022-11-24 16.14.03](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 16.14.03.png)

![截屏2022-11-24 16.17.58](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 16.17.58.png)

![截屏2022-11-24 16.20.26](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 16.20.26.png)



![截屏2022-11-24 16.31.59](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 16.31.59.png)

![截屏2022-11-24 16.39.25](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-24 16.39.25.png)





```shell
$ 1.A=$(命令)   2.A=`反引号` # sh中运行命令,返回给A

# sh获取命令行参数
$7
$*
$@
$#
$$
$?
$!
# 运算操作
$((表达式))  
$[] expr /* # 中括号前后可加空格

```



case1

```bash
#!/bin/bash
# 测试中括号左右的空格
echo $[ 1 + 3 ]


# 测试 -e 必须后面写参数，if和【】之间有空格，echo和输出有空格
if [ -e /Users/haozhipeng/Desktop/截屏2022-12-23\ 12.55.58.png ]
then echo "存在"
fi


# 测试case
case $1 in
"1")
echo "周一"
;;
"2")
echo "周二"
;;
esac


# 测试for
for  i in "$*"
do
 echo "num is $i"
done
echo "=================="
for j in "$@"
do
 echo "num is $j"
done

# 测试while,1.可以不定义变量；2.$是取变量的值，当需要赋值，如等号左边时不应该加$
i=0
while [ $i -le $1 ]
do
        sum=$[ $sum+$i ]
        i=$[ $i+1 ]
done
echo $sum
~                 

# ""用于输出连续内容，
if [ -d "$BACKUP/$dir" ] #正确
if [ -d $BACKUP/$dir ] #错误
# 
```



case2

```sh
#!/bin/bash

hostname=localhost
# 执行命令用$()
dir=$(date)
echo $dir
username=root
passwd=root
BACKUP=/users/haozhipeng/db/
database=video_analyzing
#echo $[ -d $BACKUP/$dir ]
# ！和-d之间有空格
if [ ! -d "$BACKUP/$dir/" ]
then
        mkdir -p "$BACKUP/$dir/"
fi
# gzip 后不能加“”，一般都不加“”，遇到条件表达式再加
mysqldump -u$username -p$passwd --host=$hostname -q -R --databases $database | gzip > $BACKUP/$dir/$dir.sql.gz
echo "备份开始========"
echo "备份结束========"
```

**1、蓝色是目录**

**2、白色是一般性文件，如文本文件、配置文件、源码文件等**

**3、绿色是可执行文件**

**4、黄色是设备文件**

**5、红色是压缩文件**

**6、红色是闪烁代表连接文件有问题**

**7、灰色是其他文件**

**8、浅蓝色是链接文件**




# 指令



### wc指令

```shell
$ wc testfile           # testfile文件的统计信息  
3 92 598 testfile       # testfile文件的行数为3、单词数92、字节数598 

```



### 磁盘相关指令

```shell
du -hac
lsblk -f

```

###  文本操作指令

```shell
cut -d '/' -f 3 // d以自定义/分割，并取第三个
//当重复的行并不相邻时，uniq 命令是不起作用的，搭配sort使用。uniq删除重复行，uniq -c 统计重复行出现次数
sort | uniq -c
sort -nr

#awk
awk [-F|-f|-v] ‘BEGIN{} //{command1; command2} END{}’ file

 [-F|-f|-v]   大参数，-F指定分隔符，-f调用脚本，-v定义变量 var=value

'  '          引用代码块

BEGIN   初始化代码块，在对每一行进行处理之前，初始化代码，主要是引用全局变量，设置FS分隔符

//           匹配代码块，可以是字符串或正则表达式

{}           命令代码块，包含一条或多条命令

；          多条命令使用分号分隔

END      结尾代码块，在对每一行进行处理之后再执行的代码块，主要是进行最终计算或输出结尾摘要信息
```

### 包管理

```shell
rpm -qa | grep mysql
rpm -e mysql
rpm -ivh mysql


```

### 服务

```shell
ls /usr/lib/systemd/system
systemctl status 服务
firewall-cmd --permanent --add-port=端口号/tcp

```

### 网络

```
netstat -anp |grep a所有 n是IP地址显示 p显示进程
netstat -tunlp 
-t或–tcp：显示TCP传输协议的连线状况；
-u或–udp：显示UDP传输协议的连线状况；
-n或–numeric：直接使用ip地址，而不通过域名服务器；
-l或–listening：显示监控中的服务器的Socket
-p或–programs：显示正在使用Socket的程序识别码和程序名称；
```

### 进程

```shell
ps -aux | grep  a所有 u用户 x后台信息
ps -ef 全格式显示
kill -9 
killall 名称

```

### 解压

```shell
tar -zxvf
tar -zcvf
tar -cvf
tar -xvf
```

