

# 原文对应 3页到6页开头第一段  14页到16页中间部分

## 原文

**Results**

**Formulation of the histology-based microstructure profile covariance**

**analysis**

We modelled cortico–cortical microstructural similarity across a 100-μm resolution Merker

stained 3D histological reconstruction of an entire post mortem human brain (BigBrain;

https://bigbrain.loris.ca/main.php [35]) (Fig 1A). Staining intensity profiles, representing neu

ronal density and soma size by cortical depth, were generated along 160,000 surface points

(henceforth, vertices) for each hemisphere (Fig 1B). Profile residuals, obtained after correcting

intensity profile data for the *y* coordinate to account for measurable shifts in intensity in ante

rior-to-posterior direction (S1 Fig), were averaged within 1,012 equally sized, spatially contigu

ous nodes [38]. Pairwise correlations of nodal intensity profiles, covaried for average intensity

profile, were thresholded at 0, and positive edges were log transformed to produce a micro

structure profile covariance (MPCHIST) matrix; in other words, MPCHIST captures cytoarchi

tectural similarity between cortical areas (see S2 Fig for distribution of values).

The pipeline was optimised with respect to the number of intracortical surfaces based on

matrix stability (see Methods and S3 Fig). While microstructural similarity had a small but sig

nificant statistical relationship with spatial proximity (adjusted R2 = 0.02, *P* < 0.001), similar

findings were obtained after correcting for geodesic distance.

**The principal gradient of microstructural similarity reflects sensory–fugal**

**neurostructural variation**

Diffusion map embedding, a nonlinear dimensionality reduction algorithm [39] recently

applied to identify an intrinsic functional segregation of cortical regions based on resting-state

functional MRI [28], was applied to the histology-based MPCHIST matrix (Fig 2A). The relative

positioning of nodes in this embedding space informs on (dis)similarity of their covariance

patterns. The first principal gradient (G1HIST), accounting for 14.5% of variance, was anchored

on one end by primary sensory and motor areas and on the other end by transmodal associa

tion and limbic cortices (Fig 2B; see S4 Fig for the second gradient and S5 Fig for results on

inflated cortical surfaces). G1HIST depicted the most distinguishable transition in the shape of

microstructure profiles (Fig 2B, right). Regions of the prefrontal cortex (green) expressed an

intracortical profile that was closest to the cortex-wide average. Extending outward from the

centre of G1HIST, sensory and motor regions (blue-purple) exhibited heightened cellular den

sity around the midsurface, whereas paralimbic cortex (red) displayed specifically enhanced

density near the cortical borders. For further validation of the biological basis of G1HIST, we

mapped independent atlases of laminar differentiation [40] and cytoarchitectural class [13,41]

onto the BigBrain midsurface (Fig 2C).

![截屏2022-10-20 16.36.44](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-10-20 16.36.44.png)

**Fig 1. Histology-based MPC****HIST** **analysis.** (A) Pial (purple) and WM (yellow) surfaces displayed against a sagittal

slice of the BigBrain (left) and with the midsurface (blue) in magnified view (right). 

(B) Mean and SD in residual intensity at each node are displayed on the cortex (left). Cortex-wide intensity profiles were calculated by systematic intensity sampling across intracortical surfaces (rows) and nodes (columns). 

(C) The MPCHIST matrix depicts node wise partial correlations in intensity profiles, controlling for the average intensity profile. Exemplary patterns of microstructural similarity from S1, ACC, V1, and the temporal pole. Seed nodes are shown in white. Histological data is openly available as part of the BigBrain repository (https://bigbrain.loris.ca/main.php). ACC, anterior cingulate

cortex; HIST, histology-based; MPC, microstructure profile covariance; S1, primary somatosensory; V1, primary visual; WM, white matter.

![截屏2022-10-20 16.37.17](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-10-20 16.37.17.png)

**Fig 2. The G1****HIST** **of the histology-based MPC****HIST.** 

(A) Identification: the MPCHIST matrix was transformed into an affinity matrix, which captures similarities in MPCHIST between nodes; this affinity matrix was subjected to diffusion map embedding, a nonlinear compression algorithm that sorts nodes based on MPCHIST similarity. 

(B) Variance explained by embedding components (left). The first component, G1HIST, describes a gradual transition from primary sensory and motor

(blue) to transmodal and limbic areas (red), corresponding to changes in intensity profiles, illustrated with the mean residual intensity profile across 10 discrete bins of the gradient (right).

 (C) Spatial associations between G1HIST and levels of laminar differentiation (left; [40]) and cytoarchitectural taxonomy (right; [13,41]), ordered by median. Histological data is openly available as part of the BigBrain repository (https://bigbrain.loris.ca/main.php). G1, first principal gradient; HIST, histology-based;

MPC, microstructure profile covariance.





Multiple linear regression analyses showed that levels of laminar differentiation and

cytoarchitectural taxonomy each accounted for 17% of variance in G1HIST (S2–S3 Tables).

Strongest predictors were idiotypic (β = 0.06, *P* < 0.001) and paralimbic (β = 0.05, *P* < 0.001)

laminar differentiation levels and limbic (β = 0.07, *P* < 0.001) as well as motor (β = 0.06,

*P* < 0.001) classes in the cytoarchitectural model, demonstrating the cytoarchitectural distinc

tiveness of regions at the extremes of G1HIST.



**Methods**

**Histology-based MPC**

**Histological data acquisition and preprocessing.** An ultra-high–resolution Merker

stained 3D volumetric histological reconstruction of a post mortem human brain from a

65-year-old male was obtained from the open-access BigBrain repository on February 2, 2018

(https://bigbrain.loris.ca/main.php [35]). The post mortem brain was paraffin-embedded, cor

onally sliced into 7,400 20-μm sections, silver-stained for cell bodies [98], and digitised. Man

ual inspection for artefacts (i.e., rips, tears, shears, and stain crystallisation) was followed by

automatic repair procedures, involving nonlinear alignment to a post mortem MRI, intensity

normalisation, and block averaging [99]. 3D reconstruction was implemented with a succes

sive coarse-to-fine hierarchical procedure [100]. We downloaded the 3D volume at four reso

lutions, with 100-, 200-, 300-, and 400-μm isovoxel size. We primarily analysed 100-μm data

and used 200-, 300-, and 400-μm data to assess consistency of findings across spatial scales.

Computations were performed on inverted images, on which staining intensity reflects cellular

density and soma size. Geometric meshes approximating the outer and inner cortical interface

(i.e., the GM/CSF boundary and the GM/WM boundary) with 163,842 matched vertices per

hemisphere were also available [101].

**Histology-based MPC analysis.** 1) Surface sampling. We systematically constructed 10–

100 equivolumetric surfaces in steps of 1 between the outer and inner cortical surfaces [45].

The equivolumetric model compensates for cortical folding by varying the Euclidean distance

ρ between pairs of intracortical surfaces throughout the cortex to preserve the fractional vol

ume between surfaces [44]. ρ was calculated as follows for each surface:

![截屏2022-10-20 16.40.12](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-10-20 16.40.12.png)

in which α represents fraction of the total volume of the segment accounted for by the surface,

while Aout and Ain represent the surface area of the outer and inner cortical surfaces, respec

tively. Next, vertex-wise microstructure profiles were estimated by sampling intensities along

linked vertices from the outer to the inner surface across the whole cortex. In line with previ

ous work [85], layer 1 was approximated as the top 10% of surfaces and removed from the

analysis due to little inter-regional variability. Note, however, that findings were nevertheless

virtually identical when keeping the top 10% of surfaces. To reduce the impact of partial vol

ume effects, the deepest surface was also removed. Surface-based linear models, implemented

via SurfStat for Matlab (http://mica-mni.github.io/surfstat) [102], were used to account for an

anterior–posterior increase in intensity values across the BigBrain due to coronal slicing and

reconstruction [35], whereby standardised residuals from a simple linear model of surface

wide intensity values predicted by the midsurface *y* coordinate were used in further analyses.



\2) MPC matrix construction. Cortical vertices were parcellated into 1,012 spatially contigu

ous cortical ‘nodes’ of approximately 1.5 cm2 surface area, excluding outlier vertices with

median intensities more than three scaled median absolute deviations away from the node

median intensity. The parcellation scheme preserves the boundaries of the Desikan Killany

atlas [38] and was transformed from conte69 surface to the BigBrain midsurface via nearest

neighbour interpolation. Nodal intensity profiles underwent pairwise Pearson product–

moment correlations, controlling for the average whole-cortex intensity profile. MPCHIST for a

given pair of nodes *i* and *j* was thus

![截屏2022-10-20 16.40.53](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-10-20 16.40.53.png)

in which *r**ij* is the Pearson product–moment correlation coefficient of the BigBrain intensity

profiles at nodes *i* and *j*, *r**ic* the correlation coefficient of the intensity profile at node *i* with the

average intensity profile across the entire cortex, and *r**jc* the Pearson correlation of the intensity

profile at node *j* with the average intensity profile across the whole brain. The MPC matrix was

thresholded above zero, and remaining MPC values were log-transformed to produce a sym

metric 1,012 X 1,012 MPCHIST matrix. The in-house developed code for MPC construction is

available online (https://github.com/MICA-MNI/micaopen/tree/master/MPC).

\3) Parameter estimation. The optimal surface number was determined based on the stabil

ity of the MPC matrix. This procedure involved (repeatedly and randomly) dividing the vertex

intensity profiles within each node into two groups and constructing two MPC matrices, then

calculating the Euclidean distance between them. The procedure was repeated 1,000 times.

Although the MPC matrix instability was robust to variations in surface number, the 18-sur

face solution exhibited a noticeable local minimum MPC instability in the studied range (10–

100 surfaces) and was used in subsequent analyses (S2 Fig). Notably, the MPC gradient was

similar using two finer grained solutions (i.e., 54 and 91 surfaces), in which local minima were

observed as well. More details on the origins of the stability statistic in clustering algorithms

may be found elsewhere [103].

\4) Relation to spatial proximity. To determine whether MPCHIST was not purely driven by

spatial proximity, we correlated MPCHIST strength with geodesic distance for all node pairs.

The latter was calculated using the Fast Marching Toolbox between all pairs of vertices, then

averaged by node (https://github.com/gpeyre/matlab-toolboxes/tree/master/).

**Histology-based MPC gradient mapping.** In line with previous studies [28,104], the

MPCHIST matrix was proportionally thresholded at 90% per row and converted into a normal

ised angle matrix. Diffusion map embedding [39], a nonlinear manifold learning technique,

identified principal gradient components, explaining MPCHIST variance in descending order

(each of 1 × 1,012). In brief, the algorithm estimates a low-dimensional embedding from a

high-dimensional affinity matrix. In this space, cortical nodes that are strongly interconnected

by either many suprathreshold edges or few very strong edges are closer together, whereas

nodes with little or no intercovariance are farther apart. The name of this approach, which

belongs to the family of graph Laplacians, derives from the equivalence of the Euclidean dis

tance between points in the embedded space and the diffusion distance between probability

distributions centred at those points. Compared to other nonlinear manifold learning tech

niques, the algorithm is relatively robust to noise and computationally inexpensive [105,106].

Notably, it is controlled by a single parameter α, which controls the influence of the density of

sampling points on the manifold (α = 0, maximal influence; α = 1, no influence). In this and

previous studies [28,104], we followed recommendations and set α = 0.5, a choice that retains

the global relations between data points in the embedded space and has been suggested to be

relatively robust to noise in the covariance matrix. Gradients were mapped onto BigBrain mid

surface visualised using SurfStat (http://mica-mni.github.io/surfstat) [102], and we assessed

the amount of MPCHIST variance explained. To show how the principal gradient in MPCHIST

(G1HIST) relates to systematic variations in microstructure, we calculated and plotted the mean

microstructure profiles within ten equally sized discrete bins of G1HIST.

**Relation of G1****HIST** **to laminar differentiation and cytoarchitectural taxonomy.** We

evaluated correspondence of G1HIST to atlas information on laminar differentiation and

cytoarchitectural class. To this end, each cortical node was assigned to one of four levels of

laminar differentiation (i.e., idiotypic, unimodal, heteromodal, or paralimbic) derived from a

seminal model of Mesulam, which was built on the integration of neuroanatomical,

electrophysiological, and behavioural studies in human and nonhuman primates [40] and one

of the seven Von-Economo/Koskinas cytoarchitectural classes (i.e., primary sensory, second

ary sensory, motor, association 1, association 2, limbic, or insular) [13,41]. In the case of lami

nar differentiation maps, assignment was done manually; in the case of cytoarchitectural

classes, we mapped previously published Von Economo/Koskinas classes [91] to the BigBrain

midsurface with nearest neighbour interpolation and assigned nodes to the cytoarchitectural

class most often represented by the underlying vertices. Finally, we estimated the contribution

of level of laminar differentiation (D, a categorical variable) and cytoarchitectural class (C, a

categorical variable) to the principal gradient G1HIST of the MPCHIST within two separate mul

tiple regression models:

![截屏2022-10-20 16.41.50](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-10-20 16.41.50.png)

We evaluated model fit via adjusted R2 statistics and unique variances explained by each

predictor (β).



## 翻译

## 结果

### 基于组织学的微观结构轮廓协方差的制定分析



我们在 100-μm 分辨率下对整个死后人脑 (BigBrain;https://bigbrain.loris.ca/main.php [35]）（图 1A）。沿 160,000 个表面点生成**染色强度分布**，**通过皮层深度表示神经元密度和体细胞大小**（以下称为顶点）每个半球（图 1B）。在校正 y 坐标的强度剖面数据以考虑前后方向的可测量强度变化后获得的剖面残差（S1 图）在 1,012 个相同大小的空间连续节点内进行平均 [38]。节点强度分布的成对相关性，平均强度分布的协变，阈值为 0，并且正边缘被对数转换以产生微观结构分布协方差 (MPCHIST) 矩阵； 换句话说，MPCHIST 捕获了皮质区域之间的细胞结构相似性

**根据基质稳定性，管道对皮质内表面的数量进行了优化**（见方法和S3图）。而微观结构的相似性有一个小但显著的统计数据 空间与空间接近度的物理关系（校正后的R2 = 0.02，P < 0.001），校正测地线距离后也得到了类似的结果。

![截屏2022-10-20 16.36.23](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-10-20 16.36.23.png)



**pial 软膜**

图 1. 基于组织学的 MPCHIST 分析。 (A) Pial (紫色) 和 WM (黄色) 表面显示在 BigBrain (左) 的矢状切片和放大视图 (右) 中的中间表面 (蓝色) 上。 (B) 每个节点的残余强度的平均值和 SD 显示在皮层 (左)。 通过跨皮质内表面（行）和节点（列）的系统强度采样计算皮质范围的强度分布。 (C) MPCHIST 矩阵描述了强度分布中的节点部分相关性，控制了平均强度分布。 来自 S1、ACC、V1 和颞极的微观结构相似性的示例性模式。 种子节点以白色显示。 ACC，前扣带皮层； HIST，基于组织学； MPC，微观结构轮廓协方差； S1，初级体感； V1，初级视觉； WM，白质。



### 微观结构相似性的主要梯度反映了感觉-模糊的神经结构的变化

**扩散图嵌入，一种非线性降维算法最近应用于识别基于静息状态的皮质区域的内在功能分离**
功能 MRI [28] 被应用于基于组织学的 MPCHIST 矩阵（图 2A）。亲戚节点在此嵌入空间中的定位告知它们协方差的（不）相似性模式。**第一个主梯度（G1HIST）**，占方差的 14.5%，**被锚定一端是初级感觉和运动区域，另一端是跨模式关联和边缘皮质**（图 2B；第二个梯度见 S4 图和 S5 图膨胀的皮质表面）。 **G1HIST 描绘了形状最明显的过渡微观结构轮廓**（图 2B，右）。前额叶皮层（绿色）区域表达最接近全皮层平均值的皮层内轮廓。从向外延伸G1HIST 的中心，感觉和运动区域（蓝紫色）在中表面周围表现出更高的细胞密度，而旁边缘皮层（红色）表现出特别增强皮质边缘附近的密度。为了进一步验证 G1HIST 的生物学基础，我们映射独立的层状分化图谱 [40] 和细胞结构类别 [13,41]到 BigBrain 中表面

**矩阵归一化是为了方便处理**

**基于组织学的 MPCHIST 的 G1HIST。**

 (A) 识别：将 MPCHIST 矩阵转换为亲和矩阵，该矩阵捕获了相似性节点之间的 MPCHIST；该亲和矩阵经过扩散图嵌入，这是一种基于 MPCHIST 对节点进行排序的非线性压缩算法相似。

 (B) 嵌入组件解释的方差（左）。第一个组成部分，G1HIST，描述了从初级感觉和运动的逐渐过渡（蓝色）到跨模式和边缘区域（红色），对应于强度分布的变化，用 10 个离散箱的平均残余强度分布来说明渐变（右）。

 (C) **G1HIST 与层状分化水平（左；[40]）和细胞结构分类（右；[13,41]）之间的空间关联**，按顺序排列中位数。组织学数据作为 BigBrain 存储库 (https://bigbrain.loris.ca/main.php) 的一部分公开提供。 **G1，第一主梯度； HIST，基于组织学；MPC，微观结构轮廓协方差。**

![截屏2022-10-16 16.44.18](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-10-16 16.44.18.png)

herteromodal 异模式

unimodal 单峰

residual intensity profiles 剩余强度分布

prinicpal componet 主成分

variance explained 方差解释

narmalised angle matrix 归一化角矩阵

gradient ordered matrix  梯度有序矩阵

laminar differentiation 层流分化

cytoarchitectural 细胞结构





### 基于组织学的MPC

组织学数据采集和预处理。超高分辨率染色 3D 体积组织学重建死后人脑2018 年 2 月 2 日从开放获取的 BigBrain 存储库中获得 65 岁男性（https://bigbrain.loris.ca/main.php [35]）。死后的大脑被石蜡包埋，冠状切片成 7,400 个 20 微米的切片，细胞体银染 [98]，然后数字化。人工检查人工制品（即裂口、撕裂、剪切和污点结晶）之后是自动修复程序，包括与验尸 MRI 的非线性对齐、强度
归一化和块平均[99]。 3D 重建是通过连续的从粗到细的分层过程实现的 [100]。我们下载了四种分辨率的 3D 体积，100-、200-、300- 和 400-μm 等体素大小。我们主要分析了 100 微米的数据，并使用 200、300 和 400 微米的数据来评估跨空间尺度的发现的一致性。
计算是在倒置图像上进行的，**染色强度反映了细胞密度和体细胞大小**。近似外层和内层皮质界面的几何网格（即 GM/CSF 边界和 GM/WM 边界），每个有 163,842 个匹配的顶点半球也可用，

**基于组织学的MPC分析**

1)表面采样1) 表面采样。 我们系统地构建了 10-100 个等体积表面，在外皮层和内皮表面之间以 1 为步长 。等体积模型通过改变整个**皮层的成对皮层内表面之间的欧几里德距离 ρ 来补偿皮层折叠，以保持表面之间的分数体积**。 每个表面的 ρ 计算如下

其中 α 表示由表面占的段总体积的分数，

**![截屏2022-10-18 14.31.16](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-10-18 14.31.16.png)**
**而 Aout 和 Ain 分别代表外层和内层皮质表面的表面积。**接下来，通过沿采样强度估计顶点方向的微观结构轮廓连接整个皮层的外表面到内表面的顶点。根据之前的工作 [85]，第 1 层近似为表面的前 10%，并从由于区域间差异很小，因此进行了分析。然而，请注意，调查结果仍然保持前 10% 的表面时几乎相同。为了减少部分体积效应的影响，还去除了最深的表面。基于表面的线性模型，已实现
通过 Matlab 的 SurfStat (http://mica-mni.github.io/surfstat) [102]，用于解释由于冠状切片和重建[35]，从而在进一步分析中使用来自中表面y坐标预测的表面宽强度值的简单线性模型的标准化残差



2）MPC矩阵构造。 **皮质顶点被分割成 1,012 个空间连续的皮质“节点”**，表面积约为 1.5 cm2，不包括异常顶点，远离节点的中值强度超过三个缩放中值绝对偏差中强度。 分区方案保留了 Desikan Killany 的边界，并通过最近的从 conte69 表面转换为 BigBrain 中表面
邻域插值。 节点强度分布经历了成对的 Pearson 乘积——矩相关性，控制平均全皮层强度分布。 

**其中 rij 是节点 i 和 j 处 BigBrain 强度分布的 Pearson 积矩相关系数，ric 是节点 i 处的强度分布与整个皮层的平均强度分布的相关系数，rjc 是 节点 j 的强度分布与整个大脑的平均强度分布。 MPC 矩阵的阈值高于零，剩余的 MPC 值经过对数变换以生成对称的 1,012 X 1,012 MPCHIST 矩阵**

。内部开发的 MPC 构建代码是

可在线获取![截屏2022-10-18 14.41.22](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-10-18 14.41.22.png)

基于 MPC 矩阵的稳定性确定最佳表面数。 此过程涉及（重复和随机）划分顶点，**每个节点内的强度分布分为两组并构造两个 MPC 矩阵，然后计算它们之间的欧几里得距离。 该过程重复了 1,000 次。尽管 MPC 矩阵不稳定性对表面数的变化具有鲁棒性**，但表面解在研究范围内表现出明显的局部最小 MPC 不稳定性（10-100 个表面）并用于后续分析（S2 图）。 值得注意的是，**MPC 梯度为使用两个更细粒度的解决方案（即 54 和 91 表面）类似**，其中局部最小值是也观察到。 有关聚类算法中稳定性统计量起源的更多详细信息，可能在其他地方找到

**与空间接近度的关系。 确定 MPCHIST 是否不是纯粹由空间接近度，我们将 MPCHIST 强度与所有节点对的测地距离相关联。后者是使用所有顶点对之间的快速行进工具箱计算的，然后按节点平均**

【协方差】**协方差**（Covariance）在[概率论](https://baike.baidu.com/item/概率论/829122?fromModule=lemma_inlink)和[统计学](https://baike.baidu.com/item/统计学/1175?fromModule=lemma_inlink)中用于衡量两个变量的总体[误差](https://baike.baidu.com/item/误差/738024?fromModule=lemma_inlink)。而[方差](https://baike.baidu.com/item/方差/3108412?fromModule=lemma_inlink)是协方差的一种特殊情况，即当两个变量是相同的情况。

【协方差矩阵】分别为*m*与*n*个标量元素的列向量随机变量*X*与*Y*，这两个变量之间的协方差定义为*m*×*n*矩阵.其中X包含变量X1.X2......Xm，Y包含变量Y1.Y2......Yn，假设X1的期望值为μ1，Y2的期望值为v2，那么在[协方差矩阵](https://baike.baidu.com/item/协方差矩阵/9822183?fromModule=lemma_inlink)中（1,2）的元素就是X1和Y2的协方差。

两个向量变量的协方差Cov(*X*,*Y*)与Cov(*Y*,*X*)互为转置矩阵。

协方差有时也称为是两个随机变量之间“线性独立性”的度量，但是这个含义与[线性代数](https://baike.baidu.com/item/线性代数/800?fromModule=lemma_inlink)中严格的线性独立性不同。

### 基于组织学的 MPC 梯度映射

根据之前的研究 [28,104]，MPCHIST 矩阵在每行 90% 处按比例阈值化，并转换为正常角度矩阵。扩散图嵌入识别主梯度分量，按降序解释 MPCHIST 方差（每个 1 × 1,012）。简而言之，**该算法从一个高维亲和矩阵。在这个空间中，紧密相连的皮质节点，通过许多超阈值边缘或少数非常强的边缘更靠近在一起，而具有很少或没有协方差的节点相距较远。**这种方法的名称，它
属于图拉普拉斯算子家族，**源自嵌入空间中点之间的欧几里得距离与概率之间的扩散距离的等价性**，以这些点为中心的分布。与其他非线性流形学习技术相比，该算法对噪声相对鲁棒且计算成本低廉[105,106]。值得注意的是，它由单个参数 α 控制，该参数控制密度的影响
流形上的采样点（α = 0，最大影响；α = 1，无影响）。在这和以前的研究 [28,104]，**我们遵循建议并设置 α = 0.5，这是一个保留嵌入空间中数据点之间的全局关系，并被建议为，对协方差矩阵中的噪声相对鲁棒**。 梯度被映射到使用 SurfStat (http://mica-mni.github.io/surfstat) [102] 可视化的 BigBrain 中表面上，我们评估了解释的 MPCHIST 方差量。 显示 MPCHIST 中的主梯度如何（G1HIST）与微观结构的系统变化有关，我们计算并绘制了平均值G1HIST 的十个相同大小的离散箱内的微观结构轮廓。



### G1HIST与层流分化和细胞结构分类的关系

我们评估了 **G1HIST 与关于层流分化的图谱信息的对应关系和细胞结构类。**为此，每个皮层节点被分配到四个层次之一
从一个Mesulam 的开创性模型，建立在神经解剖学的整合之上，人类和非人类灵长类动物的电生理和行为研究 [40] 和一项
七个 Von-Economo/Koskinas 细胞结构类别（即初级感觉、次级感觉、运动、关联 1、关联 2、边缘或岛状）[13,41]。**在层流分化图的情况下，分配是手动完成的；****在细胞结构的情况下类，我们将之前发布的 Von Economo/Koskinas 类 [91] 映射到 BigBrain**
**具有最近邻插值并将节点分配给细胞结构的中间表面，最常由底层顶点表示的类。**最后，我们估计了贡献层状分化水平（D，分类变量）和细胞结构类别（C，a分类变量）到两个独立的多重回归模型中 MPCHIST 的主梯度 G1HIST

![截屏2022-10-17 10.59.54](/Users/haozhipeng/Library/Application Support/typora-user-images/截屏2022-10-17 10.59.54.png)

我们通过调整后的 R2 统计量和每个解释的独特方差来评估模型拟合预测器（β）。

