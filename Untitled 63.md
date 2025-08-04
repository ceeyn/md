

一

分析中纳入的受试者的人口统计学信息列于 SI 附录表 2 中。第一次静息扫描中**头部运动与年龄的皮尔逊相关性**为 r = −0.39，P = 2.4 × 10−15，第二次扫描中为 r = −0.39，P = 2.1 × 10−15。儿童的运动明显增加（第一次静息扫描：0.167 ± 0.096；第二次静息扫描：0.092 ± 0.055），而青少年（第一次静息扫描：0.108，与青少年（第一次静息扫描：0.108 ± 0.056；第二次静息扫描：0.059 ± 0.031；第一次静息扫描：t = 7.13，自由度 [df] = 376，P = 5 × 10−12；
第二次静息扫描：t = 7.13，自由度 = 376，P = 5 × 10−12）相比，儿童（第一次静息扫描：0.167 ± 0.096；第二次静息扫描：0.092 ± 0.055）的运动明显增加。与女性参与者（第一次静息扫描：0.121 ± 0.074；第二次静息扫描：0.067 ± 0.041；第一次静息扫描：t = 4.48，自由度 = 376，P = 9.9 × 10−6；

第二次静息扫描：t = 4.19，自由度 = 376，P = 3.4 × 10−5）相比，男性参与者（第一次静息扫描：0.160 ± 0.094；第二次静息扫描：0.088 ± 0.053）的运动有所增加。当将分析严格限制在性别（女性百分比：50%）并进行运动匹配以排除静息扫描之间头部运动差异引起的潜在混杂因素时，结果仍然一致。本研究仅使用第二次静息扫描数据来推导儿童组（n = 162；第二次静息扫描：absFD =
0.088 ± 0.052）和青少年组（n = 122；第二次静息扫描：absFD =
0.060 ± 0.027；t = −1.69，自由度 = 282，P = 0.09；附录S5）的梯度。在分别考察男性和女性参与者在儿童期和青少年期的梯度发展时，结果也一致（附录S6）



二

我们评估了运动（以不同体积的平均强度差异来衡量）和我们的离散度测量之间的残差相关性（见补充图 4）。我们确实发现运动和我们感兴趣的测量值之间存在微小的相关性。因此，为了确保我们的主要关联不会因这种运动相关性而受到干扰，在我们所有的模型中，运动都被始终作为不感兴趣的协变量。此外，我们通过系统地从数据集中移除高运动个体（REFRMS 的前 5％）并比较由此产生的与年龄相关的差异的 t 统计图（补充图 5）进行了额外的敏感性分析。我们发现，移除高运动个体不会影响与年龄相关的差异。



三

rs-fMRI 数据的预处理包括以下步骤：（*i*）去除前 15 个时间点；（*ii*）切片时间校正；（*iii*）头部运动校正并排除头部运动过度的新生儿（最大头部运动> 5 mm 或> 5°，或平均帧位移> 1 mm）；（*iv*）将功能图像标准化为定制的组级模板。具体来说，我们首先将功能图像标准化为单个 T2w 结构图像，然后将单个 T2w 图像注册到公开的 37 周大脑模板。45[然后](https://www.cell.com/iscience/fulltext/S2589-0042(24)00202-5?_returnURL=https%3A%2F%2Flinkinghub.elsevier.com%2Fretrieve%2Fpii%2FS2589004224002025%3Fshowall%3Dtrue#)我们使用组平均第一个注册的 T2w 图像作为组级模板重新注册单个 T2w 图像。最后，使用第二次注册的变换参数将功能图像重新标准化为组级模板，并重新采样为 3 mm 各向同性分辨率；（*v*）平滑时间序列信号（FWHM = 4 mm）；（*vi*）去除线性漂移趋势； ( *vii* ) 通过线性回归去除虚假方差。干扰回归量包括 Friston 的 24 个头部运动参数、白质、脑脊液信号和全脑信号；( *viii* ) 时间带通滤波（0.01–0.08 Hz）；( *ix* ) 通过线性回归去除体素特定的头部运动不良时间点（FD > 0.2 mm），并排除不良时间点过多（>50%）的新生儿



四

[使用 fMRIPrep 20.2.1 ( Esteban et al., 2019](https://www.frontiersin.org/journals/aging-neuroscience/articles/10.3389/fnagi.2024.1331574/full#ref18) )对 fMRI 图像进行预处理。整体预处理包括以下步骤：从整个系列中删除前 10 个卷；切片时间校正；头部运动校正；与 T1w 图像配准；提取生理噪声回归量；估计几个混杂参数和时间序列；使用半峰全宽 (FWHM) 为 6 毫米的高斯平滑核进行空间平滑；并回归出混杂变量。脑结构异常、平均帧位移 (FD) 参数 > 0.3 毫米、或最大头部运动 > 1.5 毫米或 1.5 度、和/或超过 10% 的帧 FD > 0.5 毫米的参与者被排除在外。详细信息请参见[补充材料](https://www.frontiersin.org/journals/aging-neuroscience/articles/10.3389/fnagi.2024.1331574/full#SM1)。



五

对原生空间中去噪时间序列进行图像质量评估，以识别和排除配准不成功、残留噪声（逐帧位移 > 0.50 毫米，加上去噪时间序列显示 DVARS > 1，[Power 等人，2012 年](javascript:;)）、时间信噪比 < 50，或保留的 BOLD 样成分少于 10 个的参与者（见[补充图 1](https://oup.silverchair-cdn.com/oup/backfile/Content_public/Journal/cercor/33/1/10.1093_cercor_bhac056/1/setton_age_differences_in_the_functional_architecture_of_the_human_brain_suppl_bhac056.pdf?Expires=1748840194&Signature=SeTo6FHVG2z68iT8tvJcfOGaubTvAEedwHgFMwJkt9aIjFZoV0sqSXdZY-8tZtRsPpqB1kLHM8CSER-FgswDqB4vNRI7ZcRZm8grMiSZEhMN-ss4VG4SubUWlvK~Q4D-X1V0fxyjYeQ4P4J4qucUW58I39RQJyRQYZg4o8Acw~MBNK0JmFvlgvJ~EKnF24VZz6cCoO4IWmCJLSDuE60P00pYAxjTEezoJY7a-afHHFZDbJnJTVr2HrwXkhmUtsVVvaQUTXpRlGukfSqJkLpWfHH8sqWj8PXIkfmpnEz~W~X7FbDxc6FL2r4j-62quq2aje~nAnpSK-gldJqD4tiQ1Q__&Key-Pair-Id=APKAIE5G5CRDK6RD3PGA)为组时间信噪比图。





六（nki-rs）

此外，我们以横断面寿命样本验证了我们的研究结果。样本选取了313名健康参与者（214名女性，年龄：6-85岁，42.2±22.4岁），这些参与者来自NKI-RS数据集，未诊断任何精神或[神经系统疾病](https://www.sciencedirect.com/topics/neuroscience/neurological-disorder)，并通过了头部运动标准（平均框架位移<0.25毫米）的质量控制。NKI-RS数据由内森·克莱恩研究所的西门子TrioTim 3 Tesla扫描仪采用多波段序列采集。静息态fMRI采集时，多波段因子为4，各向同性分辨率为3毫米，重复时间为0.645秒，持续时间为9.7分钟，每次运行的vol值为900。使用连接组计算系统 ( [Xu et al., 2015](https://www.sciencedirect.com/science/article/pii/S1053811923002057?via%3Dihub#bib0091) ) 进行预处理，包括丢弃前五个时间点、压缩时间脉冲 (AFNI 3dDespike)、切片时间校正、运动校正、4D 全局平均强度归一化、干扰回归 (Friston's 24 模型、脑脊液和白质)、线性和二次去趋势以及带通滤波 (0.01–0.1 Hz)。重复预处理步骤，分别包括 GSR 和非 GSR。然后将预处理后的数据投影到 fsaverage 表面上，类似于将 HCP 数据集下采样到 fsaverage4 表面。数据集和预处理的详细信息在我们之前的研究中进行了描述 ( [Nenning et al., 2020](https://www.sciencedirect.com/science/article/pii/S1053811923002057?via%3Dihub#bib0060) ; [Nooner et al., 2012](https://www.sciencedirect.com/science/article/pii/S1053811923002057?via%3Dihub#bib0061) )。我们使用皮尔逊相关性来建立 GSR 和非 GSR 数据的单独连接矩阵。



七（nki-rs）

我们采用了横断面寿命样本，其中包括 313 名健康参与者（214 名女性，年龄：6-85 岁，42.2 ± 22.4 岁），以评估我们的方法对全脑关联研究的实际意义。参与者选自 NKI-RS [16](https://www.nature.com/articles/s42003-024-06401-4#ref-CR16)，他们未被诊断出任何精神或神经疾病，并通过了头部运动标准的质量控制（平均框架位移 < 0.25 毫米）。NKI-RS 数据是在 Nathan Kline 研究所使用西门子 TrioTim 3 Tesla 扫描仪采集的。Nathan Kline 研究所获得了机构审查委员会的批准，并获得了所有研究参与者的书面知情同意书。遵守与人类研究参与者相关的所有道德法规。以 4 的多波段因子、3 毫米各向同性分辨率和 0.645 秒的重复时间采集静息态 fMRI，持续时间为 9.7 分钟，每次运行产生 900 个体积。使用 Connectome Computational System [24](https://www.nature.com/articles/s42003-024-06401-4#ref-CR24)进行预处理，包括丢弃前五个时间点、压缩时间脉冲、切片时间校正、运动校正、四维全局平均强度归一化、干扰回归（Friston 24 模型，脑脊液和白质）、线性和二次去趋势、带通滤波（0.01-0.1 Hz）以及全局信号回归。然后将预处理后的数据投影到 32k fsLR 曲面模板上，每个半球有 32,492 个顶点



八（nki-rs）

#### 2.1.2 . NKI-RS寿命

NKI数据集来自公开共享的增强型Nathan Kline Institute-Rockland Sample数据存储库（http://fcon_1000.projects.nitrc.org/indi/enhanced/）。在本研究中，我们选择了312名参与者（214名女性，年龄：6-85岁，42.2±22.4岁），他们没有任何精神或[神经疾病](https://www.sciencedirect.com/topics/neuroscience/neurological-disorder)的诊断，并通过了头部运动标准（平均帧位移<0.25mm）的[质量控制](https://www.sciencedirect.com/topics/agricultural-and-biological-sciences/quality-control)。对于每个个体，本研究包括高分辨率T1加权扫描（1毫米各向同性分辨率，TR = 1900毫秒，TE = 2.52毫秒，翻转角= 9°）和10分钟rs-fMRI采集（3毫米各向同性分辨率，TR = 645毫秒，TE = 30毫秒，翻转角= 60°）。数据集和序列的细节在其他地方描述（[Nooner 等人，2012](https://www.sciencedirect.com/science/article/pii/S1053811920307187?via%3Dihub#bib0048)）。

NKI 数据集使用连接组计算系统 (CCS: https://github.com/zuoxinian/CCS ) 进行预处理 ( [Xu et al., 2015](https://www.sciencedirect.com/science/article/pii/S1053811920307187?via%3Dihub#bib0061) )。简而言之，结构 MRI 预处理包括空间去噪 (CAT12)、[脑](https://www.sciencedirect.com/topics/agricultural-and-biological-sciences/protocerebrum)提取、分割和表面重建。在 FreeSurfer 中计算从原生表面到 fsaverage 空间的变形球面。然后，我们生成标准表面网格，其中每个节点与原生空间中高分辨率 fsaverage 表面上的节点一一对应。功能图像预处理包括丢弃前五个时间点、压缩时间脉冲 (AFNI 3dDespike)、切片时间校正、运动校正、4D 全局平均强度归一化、干扰回归（Friston's 24 模型，[脑脊液](https://www.sciencedirect.com/topics/agricultural-and-biological-sciences/cerebrospinal-fluid)和白质）、线性和二次去趋势以及带通滤波（0.01–0.1 Hz）。然后，使用基于边界的配准方法将预处理后的数据配准到解剖空间，并投影到原生空间中的高分辨率表面（即 fsaverage）。之后，rs-fMRI 数据在 fsaverage 表面（FWHM = 6 mm）上进行空间平滑处理。最后，与 HCP 数据集类似，将数据下采样至 fsaverage4 表面，以便进行后续分析。