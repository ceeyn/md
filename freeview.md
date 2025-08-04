.dcm是一个个【一层层】的图像

.nii是整个大脑

**dcm2nii 环境变量配置**

```shell
# micron
export dcm2nii = /Applications/MRIcron.app/Contents/Resources
export PATH = ${dcm2nii}:${PATH}
```



```shell
for subj in `ls ./test`
do
	mkdir -p ./${subj}nii
	dcm2niix -o ./${subj}nii/Users/haozhipeng/downloads/test/${subj}
done
	
```

![截屏2022-11-02 10.26.30](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-02 10.26.30.png)

## FreeSufer

分割皮质和白质

把勾回”inflation“至平滑表面，”顶点“ vortex

对每个半球进行皮层重构后生成完整大脑皮层图， 共有163842个顶点

**freesufer配置**

![截屏2022-11-01 13.58.49](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-01 13.58.49.png)



```bash
# 皮层分割重建
recon -all
```

![截屏2022-11-01 14.01.46](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-01 14.01.46.png)





```shell
export SUBJECTS_DIR = /Users/haozhipeng/downloads/recon
# -s 表示被试名
# -i 表示输入文件
# -qcahe 表示大脑配准 fsaverage

for subj in `ls ./test`
do
	recon-all -s $subj -i ./test/$subj/*.nii -all -qcahe
done
```



![截屏2022-11-02 10.27.43](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-02 10.27.43.png)



![截屏2022-11-01 14.40.34](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-01 14.40.34.png)

**.ctab 用来给annotation上色的**

**.annot 根据具体功能对大脑进行分割** 

**.label annot中具体的一个个分割** **不同的脑区**，可以直接获取这个分区皮层厚度值进行组间比较

.annot是不同的.label加上代表颜色的.ctab组成的

![截屏2022-11-01 14.42.30](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-01 14.42.30.png)





**Lh.pial 根据软脑膜绘制出的底板**

皮层厚度或者皮层表面信号的选这个 

![截屏2022-11-01 15.03.41](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-01 15.03.41.png)



 lh.white 根据白质和灰质交界绘制出的底板

**不是在皮层表面，而是在灰质下选用这个做模板比较好**

![](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-01 15.04.19.png)



lh.inflated 大脑膨胀后获得的脑模板信息

信号刚好是在沟沟壑壑中，用这个

![截屏2022-11-01 14.50.52](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-01 15.04.40.png)



**Annotation选用以后**

![截屏2022-11-01 15.08.11](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-01 15.08.11.png)







![截屏2022-11-02 10.27.15](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-02 10.27.15.png)

.area面积

.fsaverage. 标准脑空间下

.fwhm 平滑后 去掉噪音 数字越来越高代表越模糊 1mm常使用10-15的结果



![截屏2022-11-01 15.19.17](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-01 15.19.17.png)

.sulc.皮层沟宽度

.thickness.皮层厚度

.volume.皮层体积

w-g 白质减灰质信号差值





![截屏2022-11-01 15.22.27](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-01 15.22.27.png)

mgz相当于nii 更接近原始图像

aparc 皮层分割

aseg 皮质下分割



![截屏2022-11-01 15.27.41](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-01 15.27.41.png)





```shell
# origz 原始图像 surf 勾勒皮层和灰白质交界 
freeview -v mri/orig.mgz -f surf/lh.white surf/lh.pial
```

![截屏2022-11-01 16.15.21](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-01 16.15.21.png)



```shell
mris_expand lh.white -1 lh.white_1
```

![截屏2022-11-01 16.28.04](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-01 16.28.04.png)



```shell
# 脑膜层绿色 灰白质中间百分之五十是黄色 灰白质交界处是红色 红色以里1mm是蓝色
freeview -v mri/orig.mgz -f surf/lh.white  surf/lh.white:edgecolor='255.0.0' -f surf/lh.graymid
-f surf/lh.innerwhite:edgecolor=blue -f surf/lh.pial:edgecolor=green -f surf/rh.white
```

![截屏2022-11-01 17.21.55](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-11-01 17.21.55.png)

