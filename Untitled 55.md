# 语言网络侧化体现在皮层的宏观功能组织中

其次，从 rs-fMRI 数据中计算出另外 3 个特征。对于 4 次 rs-fMRI 扫描中的每一次、每个个体和 18 个语言区域中的每一个，通过对位于该区域体积内的所有体素的 BOLD fMRI 时间序列取平均值来计算单独的 BOLD rs-fMRI 时间序列。然后为 995 个个体和扫描中的每一个计算内在连接矩阵。内在连接矩阵非对角线元素是区域对的 rs-fMRI 时间序列之间的 Pearson 相关系数。对于每个个体，4 个连接矩阵在取平均值之前进行*z变换，并使用双曲正切函数进行**r*变换。对 4 次扫描取平均值以提高信噪比和生成个体功能连接矩阵的可靠性[133](https://www.nature.com/articles/s41467-023-39131-y#ref-CR133)。对于每个人和每个区域，计算每个半球的强度或中心度。强度计算为一个区域与所有 18 个其他区域之间存在的相关性的总和。然后对同一半球的 18 个区域的强度值进行平均，并将得到的左、右平均强度值相加。还计算了左减右差值。最后，通过对 18 个区域对中同伦区域之间*z*变换内在相关系数进行平均，估计出每个个体的半球间连接强度。

### 连接嵌入

对于每个参与者，我们获得了前 3 个功能梯度的值。梯度反映了参与者的连接矩阵，通过 Margulies 及其同事[31的方法降低了其维数。功能梯度反映了 Mesulam ](https://www.nature.com/articles/s41467-023-39131-y#ref-CR31)[134](https://www.nature.com/articles/s41467-023-39131-y#ref-CR134)所描述的感觉整合流方面的皮质拓扑组织。梯度是使用 Python [135](https://www.nature.com/articles/s41467-023-39131-y#ref-CR135)（Python 版本：3.8.10）和*BrainSpace*库[136](https://www.nature.com/articles/s41467-023-39131-y#ref-CR136) （Python 库版本：0.1.3）计算的。在区域和顶点级别计算的梯度显示出相似的性能[136](https://www.nature.com/articles/s41467-023-39131-y#ref-CR136)。

为每个个体生成了整个皮层（*即*0.384 个 AICHA 大脑区域，与语言连接矩阵相同的过程）的 4 次扫描的平均区域级功能连接矩阵。与之前的工作一致，保留每个区域前 10％ 的连接，并将矩阵中的其他元素设置为 0 以强制稀疏性[31、49。计算矩阵任意两行之间的归一化角度距离以获得对称相似性矩阵。在相似性矩阵上实施扩散图嵌入 50、137、138](https://www.nature.com/articles/s41467-023-39131-y#ref-CR50)[以](https://www.nature.com/articles/s41467-023-39131-y#ref-CR31)[导出](https://www.nature.com/articles/s41467-023-39131-y#ref-CR49)前3个梯度。请注意，使用 Procrustes 旋转（N次迭代[= ](https://www.nature.com/articles/s41467-023-39131-y#ref-CR137)[10](https://www.nature.com/articles/s41467-023-39131-y#ref-CR138) ）将个体级梯度与相应的组级梯度对齐。此对齐过程用于提高个体级梯度与先前文献中的梯度的相似性。在个体层面对整个大脑进行最小-最大归一化（0-100）[42](https://www.nature.com/articles/s41467-023-39131-y#ref-CR42)。

为了使后续分析局限于大规模网络大脑组织，根据 Yeo 及其同事[139](https://www.nature.com/articles/s41467-023-39131-y#ref-CR139)描述的 7 个典型网络，对每个参与者的梯度值进行了平均。在平均步骤之前，每个 AICHA 区域已根据其与给定网络的空间重叠被分配到 7 个典型网络之一。然后计算每个参与者和区域的梯度不对称性。对于给定网络，梯度不对称对应于左半球的标准化梯度值减去右半球的梯度值之间的差值。





- 计算每个ROI的平均时间序列。
- 清理内存，记录计算平均时间序列的时间。

```
r
复制代码
  #.............................................................................
  # Compute correlation matrix
  c_1LR = cor(ts_ROI_1LR)
  c_2LR = cor(ts_ROI_2LR)
  c_1RL = cor(ts_ROI_1RL)
  c_2RL = cor(ts_ROI_2RL)
  ts_ROI_1LR = ts_ROI_2LR = ts_ROI_1RL= ts_ROI_2RL = NULL
  gc()
```

- 计算每个时间序列的相关矩阵。

```
r
复制代码
  #.............................................................................
  # Compute the average matrix of the 4 different rfMRI scans...................
  concat_cor = fisherz2r(((fisherz(c_1LR) + fisherz(c_2LR) +
                             fisherz(c_1RL) + fisherz(c_2RL))/4))
```



#### 导入库

```
python
复制代码
import matplotlib as mpl
import matplotlib.pyplot as plt
from nilearn import plotting
import numpy as np
import glob as gb
from brainspace.gradient import GradientMaps
```

这段代码导入了一些用于数据处理、绘图和梯度映射的库：

- `matplotlib` 和 `mpl` 用于绘图。
- `nilearn` 用于脑图像数据的处理和绘图。
- `numpy` 用于数值计算。
- `glob` 用于文件路径模式匹配。
- `brainspace.gradient` 模块中的 `GradientMaps` 用于计算梯度映射。

#### 读取并处理平均功能连接矩阵

```
python
复制代码
average_matrix = np.genfromtxt('/Data/averageData/groupLevel_correlationMatrix.txt', delimiter='\t')
np.fill_diagonal(average_matrix, 0)
```

这段代码读取存储在 `'/Data/averageData/groupLevel_correlationMatrix.txt'` 文件中的组平均功能连接矩阵，并将对角线元素设置为0。这是因为对角线上的元素表示自相关，通常为1，在功能连接分析中并不关心。

#### 初始化和拟合梯度映射

```
python
复制代码
gm = GradientMaps(kernel='normalized_angle', approach='dm', n_components=10, random_state=0)
grad_m = gm.fit(average_matrix, sparsity=0.9)
np.savetxt("/Data/groupLevel_gradient.csv", grad_m.gradients_[:, [0, 1, 2]], delimiter=',')
```

这段代码初始化 `GradientMaps` 对象，设置内核为 `normalized_angle`，方法为 `dm` (Diffusion Maps)，并指定要计算的前10个梯度组件。`random_state` 确保结果可重复。然后，该对象使用组平均功能连接矩阵进行拟合，`sparsity=0.9` 意味着在计算过程中只保留相似性矩阵中最高的90%连接。拟合完成后，将前3个梯度保存到 `"/Data/groupLevel_gradient.csv"` 文件中。

#### 处理和对齐个体功能连接矩阵

```
python
复制代码
path_matrices = gb.glob("/correlationMatrix_Individual/*.csv")
for s in range(len(path_matrices)):
    print(s)
    individual_matrix = np.genfromtxt(path_matrices[s], delimiter=',', skip_header=1)
    individual_matrix = np.delete(individual_matrix, 0, 1)
    np.fill_diagonal(individual_matrix, 0)
    gm_sujet = GradientMaps(kernel='normalized_angle', approach='dm', n_components=10, random_state=0, alignment='procrustes')
    chemin = path_matrices[s].replace('correlationMatrix_Individual/', '')
    chemin = chemin.replace('_correlationMatrix_AICHA_4scansAverage.csv', '')
    grad_individual = gm_sujet.fit([average_matrix, individual_matrix], sparsity=0.9)
    np.savetxt((chemin + '_gradient_aligned.csv'), grad_individual.aligned_[1][:, [0, 1, 2]], delimiter=',')
```

这段代码对每个个体的功能连接矩阵进行处理和对齐：

1. 使用 `glob` 库获取所有个体功能连接矩阵的文件路径。
2. 对于每个文件：
   - 读取功能连接矩阵，并删除第一列（通常是索引列）。
   - 将对角线元素设置为0。
   - 初始化一个新的 `GradientMaps` 对象，使用 `procrustes` 对齐方法。
   - 从文件路径中提取个体ID。
   - 使用组平均功能连接矩阵和个体功能连接矩阵进行拟合和对齐。
   - 将对齐后的梯度保存到新的CSV文件中。

### 

# 自闭症内在功能组织的发散不对称

### 功能连接组梯度

将预处理的时间序列进行分区后，我们得到了时间序列*分区的数组。我们首先使用时间序列计算分区之间的皮尔逊相关性，并使用 Fisher z 变换将*r*值转换为*z*值。这将为每个个体生成 360*360 的功能连接 (FC) 矩阵。然后，为了计算功能连接组梯度，我们使用非线性流形学习算法对 FC 矩阵进行降维。与功能梯度不对称框架 [ [31](https://www.nature.com/articles/s41380-023-02220-x#ref-CR31) ] 一致，我们将每个个体梯度与模板梯度 (即左-左组级梯度) 进行对齐，并使用 Procrustes 旋转使个体梯度具有可比性 [ [30](https://www.nature.com/articles/s41380-023-02220-x#ref-CR30) ]。为了在年轻人中获得无年龄或性别偏见的无偏左-左组级梯度模板，我们使用了人类功能连接组项目 S1200 版本 (HCP S1200) 的数据。这之前已经做过了 [ [31](https://www.nature.com/articles/s41380-023-02220-x#ref-CR31) ]。简而言之，我们对 1104 名受试者的 HCP S1200 FC 矩阵取平均值，并根据平均左-左 FC 矩阵计算组级梯度。第一个特征向量反映单峰-跨峰梯度 (G1)、感觉-视觉梯度 (G2) 和多需求梯度 (G3)，分别解释总方差的 24.1%、18.4% 和 15.1%。

梯度分析是在 BrainSpace [ [30](https://www.nature.com/articles/s41380-023-02220-x#ref-CR30) ] 中执行的，这是一个用于大脑降维的 Matlab/python 工具箱（https://brainspace.readthedocs.io/en/latest/pages/install.html）。梯度是连接组的低维特征向量，沿着梯度，通过许多超阈值边或少量非常强的边紧密相连的皮质节点靠得更近。同样，连接性较小的节点相距较远。这反映了功能连接配置文件的相似性/相异性，可以解释为以前三个梯度构建的公共坐标空间 [ [33](https://www.nature.com/articles/s41380-023-02220-x#ref-CR33) ] 形式描述的区域之间的功能整合和分离。这种方法属于图拉普拉斯算子家族，其名称源自扩散图嵌入中点之间的欧几里得距离的等价性 [ [32](https://www.nature.com/articles/s41380-023-02220-x#ref-CR32) , [88](https://www.nature.com/articles/s41380-023-02220-x#ref-CR88) ]。它由单个参数 α 控制，该参数反映了采样点密度对流形的影响（*α* = 0，影响最大；*α = 1，无影响）。在前人研究 [* [32](https://www.nature.com/articles/s41380-023-02220-x#ref-CR32) ] 的基础上，我们遵循建议并设置*α* = 0.5，这个选择保留了嵌入空间中数据点之间的全局关系，并且被认为对协方差矩阵中的噪声具有相对的鲁棒性。FC 矩阵中前 10% 的值被用作进入计算的阈值，这与之前的研究 [ [4](https://www.nature.com/articles/s41380-023-02220-x#ref-CR4)，[31](https://www.nature.com/articles/s41380-023-02220-x#ref-CR31)，[32](https://www.nature.com/articles/s41380-023-02220-x#ref-CR32) ] 一致。

### 不对称指数

为了量化左右半球的差异，我们选择左右作为不对称指数（AI）[ [31](https://www.nature.com/articles/s41380-023-02220-x#ref-CR31) ]。我们没有选择标准化的AI，即（左-右）/（左+右），因为梯度方差（标准化角度）既有负值也有正值[ [14](https://www.nature.com/articles/s41380-023-02220-x#ref-CR14) ]，而标准化的AI会夸大差异值或导致分母不连续[ [89](https://www.nature.com/articles/s41380-023-02220-x#ref-CR89) ]。标准化的AI与相关系数大于0.9的非标准化AI非常相似[ [31](https://www.nature.com/articles/s41380-023-02220-x#ref-CR31) ]。对于半球内模式，AI是使用从左到左的连接组梯度减去从右到右的连接组梯度来计算的。正的AI分数表示半球特征向左占主导地位，而负的AI分数表示向右占主导地位。对于半球间模式，我们使用从左到右的连接组梯度减去从右到左的连接组梯度来计算AI。我们在图中 Cohen's d 分数中添加了一个“减号”，以便于查看侧化方向（即向左或向右）。





#### 处理每个参与者的功能连接矩阵

```
python
复制代码
path = '../../data/data_autism/1_fc/ABIDE-I/'
for i in range(len(dir_num)):
  file = path+dir_num[i]+'/hcp_processed/'+file_name[i]+'_func.dtseries.nii'  
  if os.path.exists(file):
    img = nib.load(file).get_fdata()
    print('executing sub.'+dir_num[i]+'......')
    img_mmp = np.zeros((img.shape[0],360))
    img_cortex=img[:,0:18722]
    for m in range(img.shape[0]):
      for n in range(360):
        img_mmp[m][n]=np.mean(img_cortex[m][mmp_lr_clean==n+1])
    corr_matrix = np.corrcoef(img_mmp.T)
    corr_matrix[corr_matrix>0.99999] = 1
    fc = np.arctanh(corr_matrix)
    fc[fc == np.inf] = 0
    np.savetxt('../results/fc/full/'+dir_num[i]+'.csv', fc, delimiter=',')
    fc_LL = fc[0:180,0:180]
    fc_RR = fc[180:360,180:360]
    fc_LR = fc[0:180,180:360]
    fc_RL = fc[180:360,0:180]
    fc_LLRR = fc_LL - fc_RR
    fc_LRRL = fc_LR - fc_RL
    np.savetxt('../results/fc/LL/'+dir_num[i]+'.csv', fc_LL, delimiter=',')
    np.savetxt('../results/fc/RR/'+dir_num[i]+'.csv', fc_RR, delimiter=',')
    np.savetxt('../results/fc/LR/'+dir_num[i]+'.csv', fc_LR, delimiter=',')
    np.savetxt('../results/fc/RL/'+dir_num[i]+'.csv', fc_RL, delimiter=',')
    np.savetxt('../results/fc/intra/'+dir_num[i]+'.csv', fc_LLRR, delimiter=',')
    np.savetxt('../results/fc/inter/'+dir_num[i]+'.csv', fc_LRRL, delimiter=',')
```

这段代码处理每个参与者的功能连接矩阵：

1. 加载功能数据并提取皮层数据。
2. 计算每个节点的平均值。
3. 计算相关矩阵，并将大于 0.99999 的值设置为 1。
4. 对相关矩阵进行 Fisher Z 变换，并将无穷大值设置为 0。
5. 保存完整功能连接矩阵以及各个子矩阵（LL、RR、LR、RL、intra、inter）。

#### 处理 ABIDE II 数据

```
python
复制代码
df = pd.read_csv('../data/abideII_clean.csv')
dir_num = np.array(df['ID']).astype(str)
path = '../../data/data_autism/1_fc/ABIDE-II/'
for i in range(len(dir_num)):
  file_name = '%.3s'%np.array(df['Site'])[i] +'_'+dir_num[i]
  file = path+dir_num[i]+'/hcp_processed/'+file_name+'_func.dtseries.nii'  
  if os.path.exists(file):
    img = nib.load(file).get_fdata()
    print('executing sub.'+dir_num[i]+'......')
    img_mmp = np.zeros((img.shape[0],360))
    img_cortex=img[:,0:18722]
    for m in range(img.shape[0]):
      for n in range(360):
        img_mmp[m][n]=np.mean(img_cortex[m][mmp_lr_clean==n+1])
    corr_matrix = np.corrcoef(img_mmp.T)
    corr_matrix[corr_matrix>0.99999] = 1
    fc = np.arctanh(corr_matrix)
    fc[fc == np.inf] = 0
    np.savetxt('../results/fc/full/'+dir_num[i]+'.csv', fc, delimiter=',')
    fc_LL = fc[0:180,0:180]
    fc_RR = fc[180:360,180:360]
    fc_LR = fc[0:180,180:360]
    fc_RL = fc[180:360,0:180]
    fc_LLRR = fc_LL - fc_RR
    fc_LRRL = fc_LR - fc_RL
    np.savetxt('../results/fc/LL/'+dir_num[i]+'.csv', fc_LL, delimiter=',')
    np.savetxt('../results/fc/RR/'+dir_num[i]+'.csv', fc_RR, delimiter=',')
    np.savetxt('../results/fc/LR/'+dir_num[i]+'.csv', fc_LR, delimiter=',')
    np.savetxt('../results/fc/RL/'+dir_num[i]+'.csv', fc_RL, delimiter=',')
    np.savetxt('../results/fc/intra/'+dir_num[i]+'.csv', fc_LLRR, delimiter=',')
    np.savetxt('../results/fc/inter/'+dir_num[i]+'.csv', fc_LRRL, delimiter=',')
```

这段代码与处理 ABIDE I 数据的逻辑相同，只是数据源和路径不同。处理流程依旧包括加载数据、计算功能连接矩阵、进行 Fisher Z 变换和稀疏化处理，并保存结果。



- 初始化变量，用于计算左脑功能连接矩阵的平均矩阵。`n` 是文件数量，`matrix_fc_LL` 用于存储每个文件的矩阵，`total_fc_LL` 是用于累加所有矩阵的数组。

```
python
复制代码
for i in range(n):
  matrix_fc_LL[i] = np.array(pd.read_csv(path+'LL/'+path_list[i], header=None))
  total_fc_LL += matrix_fc_LL[i]
```

- 遍历每个文件，读取左脑功能连接矩阵，并累加到 `total_fc_LL` 中。

```
python
复制代码
mean_fc_LL = total_fc_LL/n
np.savetxt(path+'LL_groupmean.csv', mean_fc_LL, delimiter = ',')
```

- 计算左脑功能连接矩阵的平均矩阵，并保存为 `LL_groupmean.csv` 文件。

```
python
复制代码
### group gradients
gm = GradientMaps(approach='dm', kernel='normalized_angle',n_components=10,random_state=0)
LL = np.array(pd.read_csv('../data/LL_groupFC_HCP.csv', header=None))
gm.fit(LL)
group_grad_LL = gm.gradients_

path_add = '../results/grad/'

np.savetxt('../results/grad/group_grad_LL.csv', 
             group_grad_LL, delimiter = ',')
np.savetxt('../results/grad/group_grad_LL_lambdas.csv', 
             gm.lambdas_, delimiter = ',')
```

- 初始化 `GradientMaps` 对象，读取左脑功能连接矩阵的组级数据，并计算梯度。将计算的梯度和特征值保存到文件中。

```
python
复制代码
# individual gradients
for i in path_list:
  align = GradientMaps(n_components=10, random_state=0, approach='dm', 
                       kernel='normalized_angle', alignment='procrustes')  
  fc_LL = np.array(pd.read_csv('../results/fc/LL/'+i, header=None))
  align.fit(fc_LL,reference=group_grad_LL)
  grad_LL = align.aligned_
  np.savetxt(path_add+'LL/'+i, grad_LL, delimiter = ',')
  align = GradientMaps(n_components=10, random_state=0, approach='dm', 
                       kernel='normalized_angle', alignment='procrustes')
  fc_RR = np.array(pd.read_csv('../results/fc/RR/'+i, header=None))
  align.fit(fc_RR,reference=group_grad_LL)
  grad_RR = align.aligned_
  np.savetxt(path_add+'RR/'+i, grad_RR, delimiter = ',')
  align = GradientMaps(n_components=10, random_state=0, approach='dm', 
                       kernel='normalized_angle', alignment='procrustes')
  fc_LR = np.array(pd.read_csv('../results/fc/LR/'+i, header=None))
  align.fit(fc_LR,reference=group_grad_LL)
  grad_LR = align.aligned_
  np.savetxt(path_add+'LR/'+i, grad_LR, delimiter = ',')
  align = GradientMaps(n_components=10, random_state=0, approach='dm', 
                       kernel='normalized_angle', alignment='procrustes')
  fc_RL = np.array(pd.read_csv('../results/fc/RL/'+i, header=None))
  align.fit(fc_RL,reference=group_grad_LL)
  grad_RL = align.aligned_
  np.savetxt(path_add+'RL/'+i, grad_RL, delimiter = ',')
  print('finish   ' + i)
```

- 对每个个体的功能连接矩阵计算梯度，并与组级梯度对齐。将计算的梯度保存到相应文件夹中。

```
python
复制代码
# correct individual gradient if the correlation is negative
for dir in path_list:
  # LL
  df = np.array(pd.read_csv(path_add+'LL/'+dir,header=None))
  r = [None] * 10
  corrected = [None]*10
  for i in range(10):
    r[i] = stats.pearsonr(group_grad_LL[:,i],df[:,i])    
    if r[i][0] > 0:
      corrected[i]=df[:,i]
    else:
      corrected[i]=-1*df[:,i]
  corrected_ll = np.array(corrected).T
  np.savetxt(path_add+'LL/'+dir, corrected_ll, delimiter = ',')
  
  # RR
  df = np.array(pd.read_csv(path_add+'RR/'+dir,header=None))
  r = [None] * 10
  corrected = [None]*10
  for i in range(10):
    r[i] = stats.pearsonr(group_grad_LL[:,i],df[:,i])    
    if r[i][0] > 0:
      corrected[i]=df[:,i]
    else:
      corrected[i]=-1*df[:,i]
  corrected_rr = np.array(corrected).T
  np.savetxt(path_add+'RR/'+dir, corrected_rr, delimiter = ',')
    
  # LR
  df = np.array(pd.read_csv(path_add+'LR/'+dir,header=None))
  r = [None] * 10   
  corrected = [None]*10
  for i in range(10):
    r[i] = stats.pearsonr(group_grad_LL[:,i],df[:,i])    
    if r[i][0] > 0:
      corrected[i]=df[:,i]
    else:
      corrected[i]=-1*df[:,i]
  corrected_lr = np.array(corrected).T
  np.savetxt(path_add+'LR/'+dir, corrected_lr, delimiter = ',')
    
  # RL
  df = np.array(pd.read_csv(path_add+'RL/'+dir,header=None))
  r = [None] * 10
  corrected = [None]*10
  for i in range(10):
    r[i] = stats.pearsonr(group_grad_LL[:,i],df[:,i])    
    if r[i][0] > 0:
      corrected[i]=df[:,i]
    else:
      corrected[i]=-1*df[:,i]
  corrected_rl = np.array(corrected).T
  np.savetxt(path_add+'RL/'+dir, corrected_rl, delimiter = ',')
```

- 对每个个体的梯度进行校正，如果梯度与组级梯度的相关性为负，则反转梯度。将校正后的梯度保存到文件中。

```
python
复制代码
  # RR-LL, RL-RL
  AI_llrr = corrected_ll - corrected_rr
  AI_lrrl = corrected_lr - corrected_rl
  np.savetxt(path_add+'intra/'+dir, AI_llrr, delimiter = ',')
  np.savetxt(path_add+'inter/'+dir, AI_lrrl, delimiter = ',')
  print('finish   ' + dir)
```

- 计算左右脑之间的差异（RR-LL 和 RL-RL），并将结果保存到文件中。

```
python
复制代码
# ca network
ca_l = np.array(pd.read_csv('../data/ca_glasser_network.csv',header=None))[:,0][:180]
ca_r = np.array(pd.read_csv('../data/ca_glasser_network.csv',header=None))[:,0][180:]

for n in range(len(path_list)):
  ll = np.array(pd.read_csv(path_add+'LL/'+path_list[n],header=None))
  rr = np.array(pd.read_csv(path_add+'RR/'+path_list[n],header=None))
  lr = np.array(pd.read_csv(path_add+'LR/'+path_list[n],header=None))
  rl = np.array(pd.read_csv(path_add+'RL/'+path_list[n],header=None))
  intra = [None] * 3  
  for i in range(3):
    intra[i] = [np.mean(ll[:,i][ca_l==1])-np.mean(rr[:,i][ca_r==1]),
                np.mean(ll[:,i][ca_l==2])-np.mean(rr[:,i][ca_r==2]),
                np.mean(ll[:,i][ca_l==3])-np.mean(rr[:,i][ca_r==3]),
                np.mean(ll[:,i][ca_l==4])-np.mean(rr[:,i][ca_r==4]),
                np.mean(ll[:,i][ca_l==5])-np.mean(rr[:,i][ca_r==5]),
                np.mean(ll[:,i][ca_l==6])-np.mean(rr[:,i][ca_r==6]),
                np.mean(ll[:,i][ca_l==7])-np.mean(rr[:,i][ca_r==7]),
                np.mean(ll[:,i][ca_l==8])-np.mean(rr[:,i][ca_r==8]),
                np.mean(ll[:,i][ca_l==9])-np.mean(rr[:,i][ca_r==9]),
                np.mean(ll[:,i][ca_l==10])-np.mean(rr[:,i][ca_r==10]),
                np.mean(ll[:,i][ca_l==11])-np.mean(rr[:,i][ca_r==11]),
                np.mean(ll[:,i][ca_l==12])-np.mean(rr[:,i][ca_r==12])]
  np.savetxt(path_add+'network/intra/'+path_list[n], np.array(intra).T, delimiter = ',')
```

- 读取 Glasser 网络数据，并计算左右脑内网的差异，将结果保存到文件中。

```
python
复制代码
  inter = [None] * 3  
  for i in range(3):
    inter[i] = [np.mean(lr[:,i][ca_l==1])-np.mean(rl[:,i][ca_r==1]),
                np.mean(lr[:,i][ca_l==2])-np.mean(rl[:,i][ca_r==2]),
                np.mean(lr[:,i][ca_l==3])-np.mean(rl[:,i][ca_r==3]),
                np.mean(lr[:,i][ca_l==4])-np.mean(rl[:,i][ca_r==4]),
                np.mean(lr[:,i][ca_l==5])-np.mean(rl[:,i][ca_r==5]),
                np.mean(lr[:,i][ca_l==6])-np.mean(rl[:,i][ca_r==6]),
                np.mean(lr[:,i][ca_l==7])-np.mean(rl[:,i][ca_r==7]),
                np.mean(lr[:,i][ca_l==8])-np.mean(rl[:,i][ca_r==8]),
                np.mean(lr[:,i][ca_l==9])-np.mean(rl[:,i][ca_r==9]),
                np.mean(lr[:,i][ca_l==10])-np.mean(rl[:,i][ca_r==10]),
                np.mean(lr[:,i][ca_l==11])-np.mean(rl[:,i][ca_r==11]),
                np.mean(lr[:,i][ca_l==12])-np.mean(rl[:,i][ca_r==12])]
  
  np.savetxt(path_add+'network/inter/'+path_list[n], np.array(inter).T, delimiter = ',')
  print('finish......'+path_list[n])
```

- 计算左右脑之间的跨网差异，并将结果保存到文件中。输出处理完成的文件名。





请写成学术范的中文，可以替换专业名词，将逻辑调整的更为通顺，记得保留所有文献引用

请你扮演一个专业的论文评审专家，从图中得出一些结论，一篇学术论文中的一个段落

如果我希望将论文发表在「X」会议/期刊上，请按照「X」文章的风格，对上面的内容进行润色。



请你扮演一个专业的论文评审专家,如果我希望将论文发表在nature会议上，请按照nature文章的风格对上面的内容进行润色,请帮助我提高论文的学术严谨性，纠正任何语法错误，改进句子结构以符合学术标准，并在需要时使文本更加正式,给出修改后的英文和相应的中文翻译



## 1.英文润色

润色学术论文段落，确保其遵循学术规范，同时提升文章的拼写、语法、清晰度、简洁性以及整体阅读性。此外，可要求AI在Markdown表格中明确列出所有更改及其理由。中文指令：这是一篇学术论文中的一个段落。请重新润色，以符合学术风格，并提升其拼写、语法、清晰度、简洁性及整体可读性。如有必要，请重写句子，并在Markdown表格中列出所有更改及其原因。英文指令：Here is a paragraph from an academic paper. Please polish it to conform to academic standards, enhancing its spelling, grammar, clarity, conciseness, and overall readability. If needed, sentences should be rewritten. Additionally, list all modifications and their reasons in a Markdown table.

## 2.中文润色

改善中文学术论文的文本质量，提高其语言表达和结构的清晰度和简洁性。只提供经过修改的文本，不需附加修改原因。中文指令：作为中文学术论文写作改进专家，你的任务是增强提供文本的拼写、语法、清晰度、简洁性及整体可读性，同时简化长句，减少重复内容，并给出改进建议。请仅提供修改后的文本版本，无需解释更改原因。英文指令：As an expert in improving Chinese academic writing, your task is to enhance the spelling, grammar, clarity, conciseness, and overall readability of the provided text. Simplify complex sentences, minimize repetition, and offer improvement suggestions. Only the revised text should be provided without explanations for the changes.

## 3.SCI论文润色

在准备提交SCI论文时，逐段进行润色，以提高文章的学术严谨性，包括修正语法错误和改善句子结构，必要时提升文本的正式程度。将所有更改放入Markdown表格，包括原句、修改后的句子及其修改理由，并重新编写整个段落。中文指令：我正准备提交我的SCI论文，需要帮助润色每个段落。请帮助我提高论文的学术严谨性，纠正任何语法错误，改进句子结构以符合学术标准，并在需要时使文本更加正式。对每个需改进的段落，将修改后的句子及其原因列入Markdown表格，并重写整个段落。

英文指令：I am preparing my SCI paper for submission and need assistance in polishing each paragraph to enhance its academic rigor. Please correct any grammatical errors, improve sentence structure to meet academic standards, and formalize the text as necessary. For each paragraph that requires improvement, include the revised sentences and their reasons in a Markdown table and rewrite the entire paragraph.

## 4.期刊/会议风格润色

根据指定期刊或会议的要求，调整文本以符合其发表标准。这种润色旨在确保文本的风格、术语使用和格式均符合特定学术圈的标准。中文指令：若我计划在某XXX会议/期刊发表论文，请依照XXX文章的风格对文本进行润色，确保其符合发表要求。

英文指令：If I intend to publish my paper in a specific XXX conference/journal, please adjust the content to align with the style requirements of XXX publications, ensuring it meets the standards for submission.

## 5.结构和逻辑润色

分析并提升英文段落的逻辑性和句子间的连贯性，旨在增强文本的整体质量和易读性。此项润色专注于提升内容的逻辑流畅性，确保每个论点紧密相连，清晰表达。中文指令：作为一位专注于***研究领域的学者，我正在修订手稿以提交至***期刊。请分析并提升文本中每个段落的逻辑连贯性，改善句子间的流畅性，并给出具体的改进建议，以提高文本的整体质量和易读性。

英文指令：As a researcher in *** field, I am revising my manuscript for submission to *** journal. Please analyze and enhance the logical coherence and fluidity among sentences in each paragraph, providing specific recommendations to improve the overall quality and readability of the text.

## 6.错误纠正

如果发现前述回答中的错误或误解，根据新的信息重新进行回答。这要求AI准确理解新增的内容，并据此调整其回答，确保信息的准确性和相关性。中文指令：若先前的回答存在误解或错误，请依据我提供的新信息重新回答问题，确保内容准确无误。

英文指令：If there are misunderstandings or errors in the previous answers, please re-answer the question based on the new information I provide, ensuring the accuracy and relevance of the content.

## 7.语法检查

细致地检查文本的语法和拼写，确保没有错误。这个过程专注于识别和纠正文本中的语法和拼写问题，而不进行内容上的润色。中文指令：请帮我检查文本的语法和拼写是否准确，不需要进行内容润色。如果未发现错误，告诉我文本无误。若有发现错误，请在Markdown表格中列出，同时突出显示修正部分。

英文指令：Please check the grammar and spelling of the text for accuracy, without polishing the content. If no errors are found, inform me that the text is correct. If errors are identified, list them in a Markdown table and highlight the corrections.

## 8.语法矫正

确保文本在语法和拼写上无误，如果发现错误，则在Markdown表格中详细列出并更正，同时突出关键的修正词。中文指令：作为一名研究者，我正在修订我的手稿准备投稿。请确保文本的语法和拼写无误。如发现错误，请在Markdown表格中详细列出并更正，同时突出显示关键修正词

英文指令：As a researcher, I am revising my manuscript for submission. Please ensure that the text's grammar and spelling are correct. If errors are found, detail them in a Markdown table, correct them, and highlight the key corrections.

## 9.修改建议

作为领域专家，提供针对文本的具体修改建议，而不进行全面修改。这要求AI具体指出文本中哪些部分需要改进，并提出具体的修改建议。中文指令：作为领域专家，你认为文本中哪些部分需要修改，请具体指出并提出修改建议，但不需要全文修改

英文指令：As an expert in the field, identify specific areas of the text that require modification and provide targeted suggestions, without needing to revise the entire text.

## 10.封装基本事实/原理/背景

在润色时考虑和整合相关的基本事实、原理或背景信息，确保文本的准确性和深度。中文指令：为了帮助更好地润色文本，我需要你了解并记住特定的原理或背景信息。在润色时，请考虑这些信息，确保文本的准确性和深度

英文指令：To better polish the text, I need you to understand and remember specific principles or background information. When polishing, consider this information to ensure the text's accuracy and depth.

## 11.批判性分析

对文本进行批判性分析，识别并强化论点的有效性，确保每个论据都有充分的证据支持，并对可能的反驳进行预防和应对。中文指令：对文本进行深入的批判性分析，加强论点的说服力，确保所有论据都有充分的证据支持，并对潜在的反驳做出预防和回应。

英文指令：Conduct a critical analysis of the text to strengthen the persuasiveness of arguments, ensuring each claim is well-supported by evidence and adequately addresses potential counterarguments.

## 12.数据验证和引用

检查核实文中提到的数据和引用的准确性，确保所有引用都正确无误，并符合相关学术领域的引用标准。中文指令：验证文本中的数据和引用的准确性，确保每项引用都准确无误，并遵循该学术领域的引用规范。

英文指令：Verify the accuracy of data and citations mentioned in the text, ensuring all references are correct and adhere to the citation standards of the relevant academic field.

## 13.内容丰富和深度

拓展在保持原有论点基础上，进一步丰富内容，提供更多背景信息、案例研究或相关理论，以增强论文的深度和广度。中文指令：在不改变原有论点的基础上，丰富文本内容，添加更多背景信息、案例分析或相关理论，增加论文的深度和广度。

英文指令：Without altering the original arguments, enrich the text with additional background information, case studies, or relevant theories to enhance the depth and breadth of the paper.

﻿



[梯度不对称与年龄] 自闭症内在功能组织的发散不对称

假如你是一个java专家，基于后台开发知识，详细解释代码