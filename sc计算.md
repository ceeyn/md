

	bash  CBIG_DiffProc_batch_tractography.sh --subj_list /home/lt292794642/input_dir/sub_1.txt --input_dir /home/lt292794642/input_dir --output_dir /home/lt292794642/output_dir --py_env AMICO --mask_output_dir /home/lt292794642/output_dir/mask/
	# subj_list后面接的是存有被试信息的txt
	# 需要放置如下的目录结构
	# --py_env需要conda配置环境
		

<img src="/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2023-06-20 下午5.57.33.png" alt="截屏2023-06-20 下午5.57.33" style="zoom:50%;" />



<img src="/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2023-06-20 下午5.56.51.png" alt="截屏2023-06-20 下午5.56.51" style="zoom:50%;" />



<img src="/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2023-06-20 下午5.57.21.png" alt="截屏2023-06-20 下午5.57.21" style="zoom:50%;" />

#### -. --py_env 先安装conda 再配置conda的env

![截屏2023-06-20 下午2.00.36](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2023-06-20 下午2.00.36.png)



#### 二、requirements.txt和environment.yml

###### 2. requirements.txt

###### 2.1 生成和使用命令

​    requirements.txt的生成（开发者写的）用pip freeze命令，安装时使用也需要用pip命令，pip生成的requirements.txt用conda install无法识别。如下例所示：

```bash
pip freeze > requirements.txt # 生成requirements.txt



 



pip install -r requirements.txt # 从requirements.txt安装依赖
```



###### 2.2 内容

​    以下为一个（我正在鼓捣的一个包的）requirements.txt示例，当然这里并没有包含requirements.txt所有可能的语法要素（一般的像我这样的菜鸟也管不了这些），知道以上两个命令在大部分情况下足以生活自理了^-^。如果用"=="的话是指定了一个特定版本的包，而用“>=”则表示只要不低于这个版本就可以了，简明易懂。至于带"-e"选项的那两行我也不懂（待查阅学习和补充）。。。^-^

```bash
gym>=0.14.0



jupyter>=1.0.0



numpy>=1.16.4



pandas>=0.24.2



scipy>=1.3.0



scikit-learn>=0.21.2



matplotlib>=3.1.0



-e git+https://github.com/ntasfi/PyGame-Learning-Environment.git#egg=ple



-e git+https://github.com/lusob/gym-ple.git#egg=gym-ple



h5py>=2.9.0



pygame>=1.9.6



tqdm>=4.32.1
```

​    注意，“pip freeze”命令因为是提取当前环境的信息，因此所生成的requirements.txt应该都是"=="，">="是（确信对应包只要不低于这个版本即可而）手动编辑修改的（我瞎猜的，待确认）。

###### 3. environment.yml

​    注：关于yml or yaml?, 参见本文最后的解释。

​    environment.yml是用conda命令将环境信息导出备份的文件。

​    创建命令如下：

```bash
conda env export > environment.yml
```

​    软件安装时则执行以下命令就可以恢复其运行环境和依赖包：

```bash
conda env create -f environment.yml
```

​    注1：.yml文件移植过来的环境只是安装了你原来环境里用conda install等命令直接安装的包，你用pip之类装的东西没有移植过来，需要你重新安装。--待确认。

​    注2：environment.yml中包含该文件创建时所在的虚拟环境名称，不需要先执行"conda env create"创建并进入虚拟环境，直接在base环境下执行就会自动创建虚拟环境以及安装其中的依赖包（这个是与pip install -r requirements.txt不同的）。当然这就要求你的当前环境中没有同名的虚拟环境。如果暗装者不想使用environment.yml中内置的虚拟环境名(在environment.yml的第一行)，可以使用-n选项来指定新的虚拟环境名，如下所示：

```bash
conda env create -f environment.yml -n new_env_name
```



#### 三、yml和txt错误

![截屏2023-06-20 下午2.04.43](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2023-06-20 下午2.04.43.png)

**解决方法**

![截屏2023-06-20 下午2.06.10](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2023-06-20 下午2.06.10.png)