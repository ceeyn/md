猕猴数据。
这项研究中的所有数据集都来自公开的来源。
我们主要分析中使用的数据来自最近成立的Prime-DE(https://fcon_1000.pro-jects.nitrc.org/indi/indiPRIME.html)38.的一个队列(加州大学戴维斯分校)。
完整的数据集包括19只恒河猴(猕猴，均为雌性，年龄±SD=20.38±0.93岁，体重=9.70±1.58 kg)，使用西门子Skyra 3T四通道翻转线圈进行扫描。
所有动物都在麻醉下进行扫描。
简而言之，猕猴注射氯胺酮(10 mg/kg)、右美托咪定(0.01 mg/kg)和丁丙诺啡(0.01 mg/kg)。
麻醉维持用1-2%异氟烷酸酯。
有关扫描和麻醉方案的详细信息，请访问(https://fcon_1000.projects.nitrc.org/indi/PRIME/ucdavis.html).。
神经成像实验和相关程序是在加州国家灵长类研究中心(CNPRC)根据加州大学戴维斯机构动物护理和使用委员会107批准的协议进行的。
静息状态下采集fMRI数据，分辨率1.4×1.4×1.4 mm，tr=1.6 S，麻醉6.67min(2 5 0体积)。
扫描过程中未使用造影剂。
结构数据(T1w和T2w)以0.3×0.6×0.6 mm的分辨率采集，内插生成0.3 mm的各向同性分辨率(T1w：TR=2500ms，TE=3.65ms，TI=1100ms，翻转角=7度，FOV=154 mm；T2w：TR=3000ms，TE=307ms)。为了评估猕猴发现的稳健性，我们使用参考文献41中公开的RsFC数据评估了两个附加的猕猴样本。这些样本的年龄范围、性别和采集参数不同。
纽卡斯尔数据集(Awake108,109)包括10只恒河猴(8只雄性，平均年龄±SD=8.28±2.33，体重=11.76±3.38)，在Vertical Bruker 4.7T灵长类扫描仪上扫描。
FMRI扫描分辨率为1.2×1.2×1.2 mm，TR值为2000ms，扫描时间为8.33min(250卷×2次扫描)。
扫描过程中不使用造影剂。在牛津数据集(麻醉)的情况下，我们包括19只恒河猴，与之前的工作41一样，进行了预处理和表面重建(所有男性，年龄=4.01±0.98岁，体重=6.61±2.04公斤)。
这些猕猴是在带有4通道线圈110的3T上进行扫描的。
采集静息功能磁共振成像(rS-fMRI)数据，分辨率2倍，tr=2000ms，53.3min(1600个体积)。
扫描过程中未使用造影剂。



Macaque data. 

All datasets in this study were from openly available sources. Themacaque data used for our main analyses stemmed from one cohort (University ofCalifornia, Davis) of the recently established PRIME-DE (https://fcon_1000.pro-jects.nitrc.org/indi/indiPRIME.html)38. The full dataset consisted of 19 rhesusmacaque monkeys (macaca mulatta, all female, age ± SD=20.38 ± 0.93 years,weight=9.70 ± 1.58 kg) scanned on a Siemens Skyra 3T with 4‐channel clamshellcoil. All the animals were scanned under anesthesia. In brief, the macaques weresedated with injection of ketamine (10 mg/kg), dexmedetomidine (0.01 mg/kg),and buprenorphine (0.01 mg/kg). The anesthesia was maintained with isofluraneat 1–2%. The details of the scan and anesthesia protocol can be found at (https://fcon_1000.projects.nitrc.org/indi/PRIME/ucdavis.html). The neuroimagingexperiments and associated procedures were performed at the California NationalPrimate Research Center (CNPRC) under protocols approved by the University ofCalifornia, Davis Institutional Animal Care and Use Committee107. The resting-state fMRI data were collected with 1.4 × 1.4 × 1.4 mm resolution, TR=1.6 s,6.67 min (250 volumes) under anesthesia. No contrast-agent was used duringthe scans. Structural data (T1w and T2w) were acquired with 0.3 × 0.6 × 0.6 mmresolution, with interpolation on to generate 0.3 mm isotropic resolution(T1w: TR=2500 ms, TE=3.65 ms, TI=1100 ms,flip angle=7degrees,FOV=154 mm; T2w: TR=3000 ms, TE=307 ms).To evaluate robustness of the macaquefindings, we evaluated two additionalmacaque samples with rsFC data openly available in PRIME-DE (https://fcon_1000.projects.nitrc.org/indi/indiPRIME.html)38, previously used in ref.41.Thesamples varied in age-range, sex, and acquisition parameters. The Newcastle dataset(awake108,109) consisted of 10 rhesus macaques (8 males, age mean ± SD=8.28 ± 2.33,weight=11.76 ± 3.38) scanned on a Vertical Bruker 4.7 T primate scanner. The fMRIsession was acquired with 1.2 × 1.2 × 1.2 mm resolution, TR=2000 ms, 8.33-min perscan (250 volumes × 2 scan) per animal. No contrast-agent was used during the scans.In the case of the Oxford dataset (anesthetized), we included nineteen rhesus macaqueswith preprocessing and surface reconstruction as in previous work41(all males,age=4.01 ± 0.98 years, weight=6.61 ± 2.04 kg). The macaques were scanned on a 3Twith a 4-channel coil110. Resting-state fMRI (rs-fMRI) data were collected with 2 mmisotropic resolution, TR=2000 ms, 53.3 min (1600 volumes). No contrast-agent wasused during the scans





SC-FC耦合是可遗传的，并且不同于FC或SC的遗传性。接下来，我们使用最近开发的建模方法评估了SC-FC耦合的遗传性，该方法考虑了**所讨论的成像生物标志物的测量误差水平26。具体而言，我们设计了一个线性混合效应（LME）模型，独立估计了总体表型变异的组内和组间主观变异（代表不稳定的、短暂的成分和测量误差）。遗传性被定义为归因于基因的组间变异的比例。除了年龄、性别和利手性外，我们还将SC和FC节点强度（每行的L1范数）作为固定效应协变量包括在模型中**。总体而言，SC-FC耦合具有很高的遗传性，尤其是在亚皮质、小脑/脑干区和视觉网络中，其遗传性显著高于其他网络（中位数遗传性分别为0.78±0.16、0.70±0.22和0.57±0.20，见图5a、b）。SC-FC耦合强度与其遗传性之间存在较弱的相关性（Pearson相关系数r=0.114，p=0.140，见图5j），这表明SC-FC耦合的遗传性与其大小无关。为了与耦合遗传性进行比较，我们**计算了每种模态的区域节点强度的遗传性，见图5d、g，其中年龄、性别和利手性作为协变量。**FC的遗传性水平与SC-FC耦合相似，而SC的整体遗传性水平较低。SC-FC耦合遗传性与SC或FC的遗传性并没有明显的关联，这可由SC-FC耦合和FC遗传性之间的中等正相关性（Pearson相关系数r=0.309，p=0）以及SC-FC耦合和SC遗传性之间的无相关性（Pearson相关系数r=0.021，p=0.392）证实，见（图5k、l）。此外，FC节点强度的遗传性与SC节点强度的遗传性也没有相关性（Pearson相关系数r=0.086，p=0.089）。SC-FC耦合、FC和SC节点强度遗传性模型中每个成分（遗传效应、共同环境效应、独特环境效应和组内主观测量误差）解释的方差在补充图13、14和15中显示。



量化SC-FC耦合的遗传性。我们采用**线性混合效应（LME）模型来区分组间和组内的变异**65。这种LME方法最近被用于HCP数据，以量化**与组间成分相关的功能连接指纹的遗传性**，同时消除了单个主体观测中的短暂变化的影响26。这种方法允许考察基因关系与表型相似性之间的关联，同时考虑兄弟姐妹之间的共享环境。具体而言，我们将模型表示为：yij = xijβ + γi + εij，其中i = 1,2，...，n，j = 1,2，...，mi。mi是主体i的重复测量总数。变量yij是主体i的表型测量值，xij包含所有的q个协变量，向量β也是长度为q的未知固定总体水平效应。标量γ表示个体相对于总体平均值的偏差，εij表示yij的组内测量误差（短暂成分），并假定它与随机效应无关且在重复测量之间是独立的