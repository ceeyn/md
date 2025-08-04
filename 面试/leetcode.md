### 优先写主要逻辑，变量名定义清晰

n/k向上取整转为向下取整,（n+（k-1））/ k

向下取整转为向上取整，（n-（k-1））/k

自己定义 链表、树1，写在类外面，2.构造方法里只设值

default变量在包内都可见

内部类的私有对外部类可见

在字符串对齐问题中，`i - j`、`i - j + 1` 和 `i - j - 1` 的选择取决于你在处理单词或空格时的具体需求。它们通常用于控制循环的次数或确定单词与空格的分布情况。让我们根据你的代码中的情况来解释每种表达式的具体使用场景。

假设 `i` 和 `j` 分别是数组的下标，并且 `[j, i]` 表示一个区间，其中 `j <= i`。在这种情况下，`i-j`、`i-j+1` 和 `i-j-1` 通常用于计算该区间的长度、范围或相关的操作。让我们逐一分析它们的含义：

### 1. **`i - j` 表示什么？**

`i - j` 表示区间 `[j, i]` 的**两个端点之间的差值**，不包含区间的长度。**它只计算从下标 `j` 到下标 `i` 之间有多少步。**

- **含义**：`i - j` 代表在数组中从位置 `j` 到位置 `i` 有多少个位置的差值。
- **注意**：`i - j` 的结果是**非负数**，因为假设 `j <= i`。

**示例：**
```python
i = 7
j = 3
print(i - j)  # 输出 4，表示从 j 到 i 之间的差值
```

### 2. **`i - j + 1` 表示什么？**

`i - j + 1` 通常用于计算**区间 `[j, i]` 的实际长度**。因为区间 `[j, i]` 包含 `i` 和 `j` 两个端点，所以在计算区间长度时，需要在差值的基础上加 1。

- **含义**：`i - j + 1` 表示从 `j` 到 `i` 的区间中元素的**个数**。这个公式正确地计算了包含两个端点的区间长度。
- **区间长度公式**：`区间长度 = i - j + 1`

**示例：**
```python
i = 7
j = 3
print(i - j + 1)  # 输出 5，表示区间 [3, 7] 的长度（包含5个元素）
```

### 3. **`i - j - 1` 表示什么？**

`i - j - 1` 通常用于表示**区间 `（j, i）` 之间的内部元素个数**，**即不包含端点 `j` 和 `i` 的元素个数。如果你需要访问区间 `[j+1, i-1]` 之间的内容，`i - j - 1` 就表示该区间内部元素的数量。**

- **含义**：`i - j - 1` 表示区间 `[j, i]` 中，去掉两个端点后，区间内部剩余的元素个数。
- **注意**：当 `i == j` 时，`i - j - 1` 的结果是负数，表示区间中没有内部元素。

**示例：**
```python
i = 7
j = 3
print(i - j - 1)  # 输出 3，表示区间 (3, 7) 内部元素的数量
```

### 总结：
- **`i - j`**：表示区间 `[j, i]` **两个端点之间的差值（不含端点数）。**
- **`i - j + 1`**：表示区间 `[j, i]` 的长度**，包括端点在内的元素个数。**
- **`i - j - 1`**：表示区间 `[j, i]` 内部的元素个数，不包含端点的元素数量。

这些公式常用于数组或区间操作中，用于索引、长度计算或遍历特定区间时确定范围。

### `ArrayDeque` 不允许存储 `null` 元素，这是一个限制。你可以用其他的数据结构，如 `LinkedList`，它是允许存储 `null` 的。

我感觉摩尔投票是动态规划的思想，他们都有一个共同点：直到某一更好的状态出现之前，答案取当前最好的状态。

从本题看，给定一个数组，如[1,5,0,5,5]，在1或0没被抵消完之前，1或0就是当前的众数，直到5出现的次数发生正增长。

然后举一个动态规划的例子，比如买卖股票的问题，给定一个每日股票价格数列，如[1,5,0,5,5]，只允许执行一次交易，求最大利润的买入时机。在买入价0元出现之前，最佳买入时机一直是1元，直到0元买入5元卖出的交易窗口出现。





**首先比较两个字符串的长度，如果长度不同，则较长的字符串更大。如果长度相同，则逐位比较字符，直到找到不同的字符为止。这种方法适用于比较较大的数字字符串**

# hot 100

### 模拟

[15. 三数之和](https://leetcode.cn/problems/3sum/)【1 i 从 0 到n-2，if (i+j+k == 0) j++，k--，然后跳过j，k重复的元素，否则只是移动j，k】

[400. 第 N 位数字](https://leetcode.cn/problems/nth-digit/)【1.找到digit位，找到start开始，2.找到在哪个数字start+(n-1)/digit 3.找到在数字具体哪位n-1 % digit】

**JZ61** **扑克牌顺子**【1.排序后，不能有重复，间隔数 <= 空格数，2.set判断不能有重复，max-min《= 4】

**JZ21** **调整数组顺序使奇数位于偶数前面(一)**【两个指针，一个遍历数组，另一个指向奇数所在位置，找到奇数后，将元素移动到奇数所在位置】

[7. 整数反转](https://leetcode.cn/problems/reverse-integer/)【ans < Integer.MIN_VALUE / 10-cur / 10 || ans > Integer.MAX_VALUE / 10-cur / 10说明溢出】

**JZ39** **数组中出现次数超过一半的数字**【摩尔投票，当前元素等于cond，cnt++，不等于cnt--，若cnt == 0 说明没有候选者】

**JZ43** **整数中1出现的次数（从1到n整数中1出现的次数）**【cur % 10 拿到每一位，是1的话++】

**JZ45** **把数组排成最小的数**【快速排序，都是大于，除了大于等于不移动】【最小，排序规则 i + j  > j + i，则j排在i前】

[165. 比较版本号](https://leetcode.cn/problems/compare-version-numbers/)【1.split("\\\\.")分割，while (i<m || j < n)进行对比】

[43. 字符串相乘](https://leetcode.cn/problems/multiply-strings/)【a * b，一共 m+n位，第i位，第j位乘积对应i+j+1，进位i+j】

**NC359** **大数相减**【num1的i，num2的j，进位carry，while(i < num1 || j < num2 || carry != 0), 当前位 num1[i]+num2[j]+carry,若小于0，则+10，carry=-1，若大于0，则carry = 0】

[470. 用 Rand7() 实现 Rand10()](https://leetcode.cn/problems/implement-rand10-using-rand7/)【1.利用一个小范围随机数，得到一个大范围等概率随机数：rand xy =（ rand x - 1）* y + rand y [✖️谁就减谁+rand谁]2. 利用一个大范围随机数，得到一个小范围等概率随机数rand x = rand x的倍数 % x + 1，3. 拒绝采样】

[394. 字符串解码](https://leetcode.cn/problems/decode-string/)【1.栈中保存《左括号前的字符，左括号前的数字》，定义cur为当前所有不在括号内的字符组合，遇到左括号压栈，遇到右括号弹出栈顶，cur = 左前字符 + 左前数字 * cur，最后返回cur；2. Map（左括号下标，对应右括号下标），cal（s，开始下标，结束下标，map）计算开始下标到结束下标的结果，遇到左括号时递归调用】

[224. 基本计算器](https://leetcode.cn/problems/basic-calculator/)【1.Map（左括号下标，对应右括号下标），cal（s，开始下标，结束下标，map）计算开始下标到结束下标的结果，遇到左括号时递归调用，将s看做，sign，num，【符号·数字】，每次遇到符号时或者 i 为 r时，将sign num入栈,或者 sign 为* /，将 st 上一个元素 * num，然后 sign = cur，num = 0】

[179. 最大数](https://leetcode.cn/problems/largest-number/)【若干个数组成最大数，转换成string比较字典序 i+j 》 j+ i，则i在j前面，然后去掉前导0】

[498. 对角线遍历](https://leetcode.cn/problems/diagonal-traverse/)【1.利用x+y = k， k是对角线的序号，以及 0=<x < m;  0=<y < n来构造，k为偶数右上，k为奇数左下；2.goingUp = true右上走, 若先到达最右边，x++，先到达上边 y++，左下颠倒一下】

[384. 打乱数组](https://leetcode.cn/problems/shuffle-an-array/)【1.洗牌算法，i和后面i+random.nextInt(n-i)交换位置，random.nextInt(n)返回0-n-1的一个数】

[287. 寻找重复数](https://leetcode.cn/problems/find-the-duplicate-number/)【1.将值当成next的索引，因此 target 这个位置一定有起码两条指向它的边，因此整张图一定存在环，且我们要找到的 target 就是这个环的入口,2.二分答案，求最小/第一个思路，l不符合，r符合，小于mid元素个数 > mid,说明在左边，否则说明在右边】

[678. 有效的括号字符串](https://leetcode.cn/problems/valid-parenthesis-string/)【1.*一个栈，左括号一个栈，遇到右括号，优先用左括号匹配，最后如果左括号栈不为空，则检查星号栈，每个星号位置必须大于左括号】

[611. 有效三角形的个数](https://leetcode.cn/problems/valid-triangle-number/)【a+b > c, i倒序枚举c，j，k枚举a，b，若j+k>i,则k和j之间所有都符合，ans += k-j】

[134. 加油站](https://leetcode.cn/problems/gas-station/)【局部失败意味着全局失败：当你从一个加油站开始出发，途中发现汽油不足以到达下一个加油站，这意味着从这些加油站出发就不可能绕一圈回到这个加油站】

[71. 简化路径](https://leetcode.cn/problems/simplify-path/)【遇到..栈不空则弹出栈顶，for 从栈底遍历栈】

[168. Excel 表列名称](https://leetcode.cn/problems/excel-sheet-column-title/)【转为26进制，不断除k取余】

[1047. 删除字符串中的所有相邻重复项](https://leetcode.cn/problems/remove-all-adjacent-duplicates-in-string/)【1.使用栈2.双指针：fast探测，slow模拟栈顶指针】

[316. 去除重复字母](https://leetcode.cn/problems/remove-duplicate-letters/)【维护递增单调栈，如果栈中相邻的元素字典序更大，那么我们选择丢弃相邻的栈中的元素。】

* [402. 移掉 K 位数字](https://leetcode.cn/problems/remove-k-digits/)【维护递增单调栈】

[395. 至少有 K 个重复字符的最长子串](https://leetcode.cn/problems/longest-substring-with-at-least-k-repeating-characters/)【1.对字符种类数量进行滑窗，2.递归】

**NC89** **字符串变形**【1.桟，2.两次翻转】

**字符串出现次数的TopK问题**【1.hashmap 统计前 k 个，2.priorityqueue 统计】

**NC106** **三个数的最大乘积**【1.找到 max3，max2，max1，min1，min2，最大值是 Max（max1*min1*min2，max1max2max3）】

**阶乘末尾0的数量**【每次ans += n / 5】

**NC142** **最长重复子串**【从大到小枚举长度，然后枚举起点 i，如果 i ！=i+len，i 之前所有起点都不行】

### hash

1.[128. 最长连续序列](https://leetcode.cn/problems/longest-consecutive-sequence/)【使用hash，如果不是起始第一个元素则continue，若是第一个元素则向后拓展求最大长度】

2.[41. 缺失的第一个正数](https://leetcode.cn/problems/first-missing-positive/)【原地哈希，不相等一直交换到正确】 x 应当出现在数组中的 x−1 的位置，因此交换 nums[i] 和 nums[x−1]，这样 x 就出现在了正确的位置。在完成交换后，新的 nums[i] 可能还在 [1,N] 的范围内，我们需要继续进行交换操作

3.**JZ75** **字符流中第一个不重复的字符**【hash表记录每个元素出现次数，pq是第一个不重复元素的队列】

4.[49. 字母异位词分组](https://leetcode.cn/problems/group-anagrams/)【每个字符排序后hash分组】

5.[128. 最长连续序列](https://leetcode.cn/problems/longest-consecutive-sequence/)【通过set判断i+1是否存在】

[349. 两个数组的交集](https://leetcode.cn/problems/intersection-of-two-arrays/)【两个hash，hash1记录数组1，hash2记录hash1有且数组2有的元素，即答案】

### 滑动窗[最大最长，个数至少，外层更新]

1.438找到字符的所有字母异位词

2.**JZ74** **和为S的连续正数序列**【滑窗枚举右端点，while内更新答案，（r+l）*（r-l+1）/2 等于tar就是答案】

[76. 最小覆盖子串](https://leetcode.cn/problems/minimum-window-substring/)【need map，cur map，cur map中符合条件的个数count，入出窗口只有是need中包含的字符才加入cur】

[1004. 最大连续1的个数 III](https://leetcode.cn/problems/max-consecutive-ones-iii/)

### 前缀和

1.560和为k的子数组（前缀和+hash，// 先判断，后更新）

**NC125** **和为K的连续子数组**【map 的 key 是 sum，val 是 sum 第一次出现的下标，刚开始放入 0，-1；每次先更新 sum，再更新 ans，最后更新 map】

2.[53. 最大子数组和](https://leetcode.cn/problems/maximum-subarray/)（1.买卖股票，维护前缀和最小值2.dp）

3.**JZ85** **连续子数组的最大和(二)**【同上面做法，维护前缀和最小值, 初始化sum=0, gl = gr = 0, 上次最小值的位置为ll = 0，大于等于ans时，记录gr = i， gl = ll， 每次更新最小值记录ll = i+1 】

### 贪心

1.[763. 划分字母区间](https://leetcode.cn/problems/partition-labels/)【如果不要求合并后的具体区间的话，可以维护当前区间的end这种做法】

2.[55. 跳跃游戏](https://leetcode.cn/problems/jump-game/)【如果不要求合并后的具体区间的话，可以维护当前区间的end这种做法，// 先判断答案(i > end) 再进行更新】

3.[45. 跳跃游戏 II](https://leetcode.cn/problems/jump-game-ii/)【int curRight = 0; // 已建造的桥的右端点，int nextRight = 0; // 下一座桥的右端点的最大值，先判断后更新，超过 curRight 时就ans++，更新nextRight】

4.[452. 用最少数量的箭引爆气球](https://leetcode.cn/problems/minimum-number-of-arrows-to-burst-balloons/)【1.右端点从小到大，维护右端点最小值，当当前区间左端点大于右端点最小值时，ans++2.合并区间求交集，答案是交集的数量】

**主持人调度（二）**【将活动分为 starts，ends 数组，对两个数组分别排序，双指针指向两个数组，如果当前 start[i] < start[j]则 ans++，否则 j++】

### 动态规划

2.[53. 最大子数组和](https://leetcode.cn/problems/maximum-subarray/)（1.curMax = Math.max(curMax,0)+nums[i]，ans = Math.max(curMax,ans)）

[918. 环形子数组的最大和](https://leetcode.cn/problems/maximum-sum-circular-subarray/)【维护gMax，gMin，curMax，curMin，sum，sum == gMin？gMax ：Math.max(gMax,sum-gMin)】

3.[LCR 168. 丑数](https://leetcode.cn/problems/chou-shu-lcof/)【优先队列队头元素是最小丑数，每次取出队头x2，x3，x5利用hashset去重后再放入queue中】【三个丑数序列每次x2，x3，x5，是1x2，1x3,1x5，合并这三个序列，三个指针最小的放入数组并后移一个位置】

4.[LCR 187. 破冰游戏](https://leetcode.cn/problems/yuan-quan-zhong-zui-hou-sheng-xia-de-shu-zi-lcof/)【约瑟夫环，看起点的变化，fn = （fn-1+m） % n，fn = f·n-1，删除第m个后起点变成m，fn-1起点是0，f·n-1 = fn-1 + m = fn】

**NC67** **汉诺塔问题**【1.dfs(left,mid,right)将left借助mid移动到right上，将前n-1个节点借助right移动到mid，将第n个节点借助mid移动到right，再将n-1个节点借助left移动到right】

[718. 最长重复子数组](https://leetcode.cn/problems/maximum-length-of-repeated-subarray/)【dfs(i,j,pre_select)代表以nums[i],nums[j]结尾的最大长度，上一个i+1，j+1，select=1的情况下，当前若nums[i] == nums[j]，则必须选】

5.[LCR 126. 斐波那契数](https://leetcode.cn/problems/fei-bo-na-qi-shu-lie-lcof/)【滚动数组，fn = f n-1 + f n-2 ，p =q，q =r ， r = p + q】

### 矩阵

1.[54. 螺旋矩阵](https://leetcode.cn/problems/spiral-matrix/)【t，b，l，r，尽可能遍历到底，遍历完第一行相应t++，例如i从l到r，t++】

2.[48. 旋转图像](https://leetcode.cn/problems/rotate-image/)【上下：行到m/2，i行变成m-1-i，对角线：列到i，ij变ji】

[240. 搜索二维矩阵 II](https://leetcode.cn/problems/search-a-2d-matrix-ii/)【排除法双指针，每次比较右上角元素和 tar，如果大于则减少一列 j--，小于则增加一行 i++】

3.**JZ58** **左旋转字符串**【整体翻转，反转前i个，反转剩余部分】



### 链表

* *[92. 反转链表 II](https://leetcode.cn/problems/reverse-linked-list-ii/)【要设 dummy，p0指向上一段末尾，反转结束后pre是这一段末尾，cur是下一段开始】

1.[138. 随机链表的复制](https://leetcode.cn/problems/copy-list-with-random-pointer/)【map<原来节点，新节点>】

2.[25. K 个一组翻转链表](https://leetcode.cn/problems/reverse-nodes-in-k-group/)【得到sz，p0 = dummy,代表上一段最后一个节点， pre = null， 反转完k个后pre指向这段最后一个节点，pre指向这一段最后一个节点，cur指向下一段第一个，p0.next.next = cur, p0.next = pre, p0 = p0.next】

3.[82. 删除排序链表中的重复元素 II](https://leetcode.cn/problems/remove-duplicates-from-sorted-list-ii/)【如果cur.next = cur.next.next，套一个循环一直删除cur.next，不等于的话cur = cur.next】

4.[146. LRU 缓存](https://leetcode.cn/problems/lru-cache/)【三个变量：map[k->Node]，cap，dummy；三个函数getNode(key),putFront(Node),remove(Node),get抽书放到最上面，put 有这本书则getNode然后更新节点，没有这本书则放到最上面，如果超出最大容量，则删除末尾,因此需要在node里存key，为了在hash中删除最后节点】

5.[61. 旋转链表](https://leetcode.cn/problems/rotate-list/)【旋转k个，快慢指针找到倒数第k个，然后把链表末尾衔接表头，倒数第k个新表头】

6.[148. 排序链表](https://leetcode.cn/problems/sort-list/)【快慢指针找到中点+归并排序】

[24. 两两交换链表中的节点](https://leetcode.cn/problems/swap-nodes-in-pairs/)【node0，node1，node2，node3,0—>2,2->1,1->3,0 = 1,1 =3】

[460. LFU 缓存](https://leetcode.cn/problems/lfu-cache/)【1.四个变量：map[freq->Node] 多摞书，map[k->Node],minFreq,cap，四个函数getNode(key),putFront(freq,Node),remove(Node)，newList()】

[LCR 171. 训练计划 V](https://leetcode.cn/problems/liang-ge-lian-biao-de-di-yi-ge-gong-gong-jie-dian-lcof/)【l1走到headA的结尾则，l1 = headB，l2同理，直到相遇】

### 树

1.[101. 对称二叉树](https://leetcode.cn/problems/symmetric-tree/)【左子树和右子树必须是相同的树，相同定义：t1的左等于t2的右，t2的右等于t1的左】

2.[543. 二叉树的直径](https://leetcode.cn/problems/diameter-of-binary-tree/)【最大链长】

3.[437. 路径总和 III](https://leetcode.cn/problems/path-sum-iii/)【记得用long：1.前缀和+哈希2.dfs(t1,t2)--------->do(t1,tar)函数实际以t1为根的树做操作，如果t1为根往下和为tar即t1.val == tar则个数加1，dfs(t1.left,tar)，dfs(t1.right,tar)在这棵树上游走】

4.[236. 二叉树的最近公共祖先](https://leetcode.cn/problems/lowest-common-ancestor-of-a-binary-tree/)【dfs(TreeNode root,TreeNode p, TreeNode q) 在root为根的树上找p、q的最近祖先】

5.[124. 二叉树中的最大路径和](https://leetcode.cn/problems/binary-tree-maximum-path-sum/)【最大链，直径更新全局最小】

**[129. 求根节点到叶节点数字之和](https://leetcode.cn/problems/sum-root-to-leaf-numbers/)【dfs(root,sum) sum是从root到某一个叶节点的和】

[LCR 143. 子结构判断](https://leetcode.cn/problems/shu-de-zi-jie-gou-lcof/)【dfs(t1,t2)--------->do(t1,t2)函数实际以t1，t2为根的树做操作，dfs(t1.left,t2)，dfs(t1.right,t2)在这棵树上游走】

* [LCR 152. 验证二叉搜索树的后序遍历序列](https://leetcode.cn/problems/er-cha-sou-suo-shu-de-hou-xu-bian-li-xu-lie-lcof/)【dfs(l,r)判断是否是后序序列，r是根，从左到右找到左子树的边界j，检查剩余[j,r-1]是否都大于r,然后递归左dfs(l,j),右dfs(j+1,r-1)】

JZ36 二叉搜索树与双向链表【全局pre，head,每次遍历cur时，pre.right = cur,cur.left = pre, pre = cur,最后返回head】

**JZ8** **二叉树的下一个结点**【分情况讨论：1.当前节点的右子树最左下，2.当前节点是父节点的左子树，3.当前节点是父节点右子树】

[230. 二叉搜索树中第 K 小的元素](https://leetcode.cn/problems/kth-smallest-element-in-a-bst/)【中序位置k--，如果k==0，就是该元素】

[98. 验证二叉搜索树](https://leetcode.cn/problems/validate-binary-search-tree/)【1.前序：参数传入min，max作为约束，判断每个节点是否符合；2.中序：记录前缀，遍历完左子树，判断当前val是否大于pre，然后设置pre为val，然后遍历右子树，3.后序：记录子树的min，max，要求当前节点val大于左子树最大值，小于右子树最小值】

* [173. 二叉搜索树迭代器](https://leetcode.cn/problems/binary-search-tree-iterator/)【双色遍历法】

**NC11** **将升序数组转化为平衡二叉搜索树**【dfs(l,r),mid做根】

* **找到搜索二叉树中两个错误的节点**【pre是前一个节点，第一次出现前一个节点大于后一个节点时，前一个节点是较大那个，第二次出现前一个节点大于后一个节点时，后一个节点是较小那个】

**NC60** **判断一棵二叉树是否为搜索二叉树和完全二叉树**【1.判断是否是完全二叉树：层序遍历的过程中遇到第一个空节点之后不应该再出现非空节点】

* [662. 二叉树最大宽度](https://leetcode.cn/problems/maximum-width-of-binary-tree/)【记录节点序号，左 2 * index + 1，右 2 *index+2，找到每层第一个节点，然后每层当前index - 最小index+1更新答案，bfs，dfs都可以，dfs(TreeNode root, int depth,int index)先序遍历root的深度为depth，编号为index，得到宽度，其中每一层最左边节点的编号用HashMap<Integer, Integer>记录】

[199. 二叉树的右视图](https://leetcode.cn/problems/binary-tree-right-side-view/)【递归右子树，再递归左子树，当depth == ans.size()时，就是每层第一个】

* [297. 二叉树的序列化与反序列化](https://leetcode.cn/problems/serialize-and-deserialize-binary-tree/)【1.dfs：先序遍历将字符串序列化成树:serialize(root) ,反序列化TreeNode deserialize(String data)时使用队列，先把字符串都加进队列中，dese(pq)每次从队头弹出，当做当前节点，左节点等于dese(pq)的返回值】

### dfs

1.[207. 课程表](https://leetcode.cn/problems/course-schedule/)【dfs判断是否有环，0 没有访问过 1 当前节点正在访问 2 当前节点访问完毕 若出现两次1则说明有环】

2.[210. 课程表 II](https://leetcode.cn/problems/course-schedule-ii/)【1.dfs(i)判断环之后，将i加入ans，最后输出逆序的ans即为拓扑序列；2.bfs 维护领接表，入度表，刚开始把入度为0加入队列，每次出队元素将后继节点入度--，继续过程】

[679. 24 点游戏](https://leetcode.cn/problems/24-game/)【dfs(int[]nums)以nums运算能否得到24点，最小情况len(nums) == 1，每次从nums选i,j做运算将答案curRes放入next，递归到dfs(int[]next)】

**NC138** **矩阵最长递增路径**【往上，下，左，右四个方向移动，dfs（i，j）表示以 i，j 为结尾的最大长度，可以从上，下，左右转移来,不能用动态规划，动态规划最多处理两个方向】

### 二分

1.[33. 搜索旋转排序数组](https://leetcode.cn/problems/search-in-rotated-sorted-array/)【红蓝染色，红色在目标值左边，蓝色在目标值右边，和最后一个元素比(分为大于最后一个，小于等于最后一个)，check】

2.[153. 寻找旋转排序数组中的最小值](https://leetcode.cn/problems/find-minimum-in-rotated-sorted-array/)【红蓝染色，红色在最小值左边，蓝色在最小值右边，和最后一个元素比，check】

3.[4. 寻找两个正序数组的中位数](https://leetcode.cn/problems/median-of-two-sorted-arrays/)【奇数中位数 n+1 / 2,偶数中位数 (（n+1)/2 + (n+2)/2)/2,   findk（nums1，start1，nums2，start2，k）从nums1的start1开始，nums2的start2开始找第k个数，拿出来两个的k/2的数（对应start1+k/2-1）比较，若nums1大，淘汰nums2前k/2个数】

4.[69. x 的平方根 ](https://leetcode.cn/problems/sqrtx/)【求满足题目条件的最大值，左边是满足条件元素，check(mid)为ture更新左边，最后返回左边】



### 单调栈

1.[84. 柱状图中最大的矩形](https://leetcode.cn/problems/largest-rectangle-in-histogram/)【三步走：1.栈顶是满足答案条件的元素 2.新元素入栈时，对于新入栈元素，及时去掉无用元素（while pop）3.新元素入栈并更新答案，单调栈1求左边小于当前高度的下标i，单调栈2求右边小于当前高度的下标j，当前高度的最大面积h * (j-i+1)】

2.[85. 最大矩形](https://leetcode.cn/problems/maximal-rectangle/)【计算每一行为底部的高度前缀和数组，然后对每行进行84的求法】



### 栈

1.*[32. 最长有效括号](https://leetcode.cn/problems/longest-valid-parentheses/)【栈顶是左括号或者尚未匹配上的右括号，刚开始设置栈的初始基准位置-1】

2.栈实现队列【输入栈、输出栈，输入栈一直入，输出时，若输出栈为空，将输入栈所有放到输出栈，就是把先入先出-》后入先出】

3.**JZ30** **包含min函数的栈**【两个栈，一个正常输入输出，一个最小栈，栈顶永远是最小元素】

4.**JZ31** **栈的压入、弹出序列**【遍历弹出序列，如果栈顶不是弹出元素，则一直入栈，最后判断栈为空说明符合】

5.[71. 简化路径](https://leetcode.cn/problems/simplify-path/)【split("/")分割，压入栈，“..”弹出栈顶，最后拼接】

**NC115** **栈和排序**【输入某个数字序列最大的输出序列，记录每个元素后面最大的元素是什么，某个元素入栈时，如果该元素后面没有更大的元素，则输出，否则入栈】



### 单调队列

[239. 滑动窗口最大值](https://leetcode.cn/problems/sliding-window-maximum/)【1.滑动窗口，入时及时去掉队列中无用元素，出时判断队首是否还在k个元素内】



```
// i指的都是被选元素下标
// 选或不选：dfs(i)指的是被选元素中第i个选或不选
// 枚举选哪个：dfs(i)指的是当前选的元素一定是大于等于i的一个
```



```
// 选或不选
// 当前操作：枚举第i个数选或不选
// 子问题：从下标>=i的元素中构造子集
// 下一个子问题：从下标>=i+1的数字中构造子集
```



```
 // 枚举选哪个[基于答案角度一定选]
 // 当前操作：第i个位置，枚举一个下标j>=i的数字，加入path
 // 子问题：第i个位置，从下标>=i的数字中构造子集
 // 下一个子问题：第j个位置，从下标>=j+1的数字中构造子集
 // j是待候选元素的下标，只是刚好在子集问题中可以代表i的下标，i才是代表构造答案的哪一位
```



### 选或不选 

### 完全背包【】

### 01背包【[416. 分割等和子集](https://leetcode.cn/problems/partition-equal-subset-sum/)】】

**NC17** **最长回文子串**【1.dfs(i,j)选或不选返回长度，i == j 返回1，全局gl，gr，分情况i和j处字符相等： j-i+1 = dfs(i+1,j-1),更新gl，gr，不等，dfs(i+1,j),dfs(i,j-1)更新gl，gr】

[213. 打家劫舍 II](https://leetcode.cn/problems/house-robber-ii/)【分为偷nums[0]，不偷nums[0]讨论，偷nums[0]，则rob(2,n-2),不偷rob(1,n-1)】

*[221. 最大正方形](https://leetcode.cn/problems/maximal-square/)【dfs(i,j)代表以i为结尾，以j为结尾的最大面积，dfs(i,j) = Math.min(dfs(i-1,j),dfs(i,j-1),dfs(i-1,j-1))+1】

**NC138** **矩阵最长递增路径**【往上，下，左，右四个方向移动，dfs（i，j）表示以 i，j 为结尾的最大长度，可以从上，下，左右转移来,不能用动态规划，动态规划最多处理两个方向】

### 枚举选哪个

【[300. 最长递增子序列](https://leetcode.cn/problems/longest-increasing-subsequence/)】【dfs(i)代表以i为结尾的最长子序列长度，一定选】

[673. 最长递增子序列的个数](https://leetcode.cn/problems/number-of-longest-increasing-subsequence/)【dfs(i)代表以i为结尾的最长子序列长度，最长子序列个数，维护全局最长子序列长度，最长子序列个数，进行更新】

**JZ46** **把数字翻译成字符串**【划分型dp，dfs(i)代表子数组以i为结尾，枚举划分的左端点j，下一个子问题就是dfs(j-1)】

[1884. 鸡蛋掉落-两枚鸡蛋](https://leetcode.cn/problems/egg-drop-with-2-eggs-and-n-floors/)【dfs(i)代表n层的最小次数，枚举j从1到i，dfs(i) = Math.min(Math.max(k,dfs (j-k)+1)),在第k层扔，分为两种情况：1.在第k层碎了，需要从1扔到k一共k次；2，在第k层没碎，剩下等价于dfs（i-k)子问题】最优解：如果已知答案（操作次数），*n* 最大可以是多少？假设最优次数x，第一次在第x层扔最好，第二次 x+x-1。。。。。一共要x+x-1+。。。 = n，解出x

[887. 鸡蛋掉落](https://leetcode.cn/problems/super-egg-drop/)【dfs(i,j)代表i次机会，j个鸡蛋能到达的最大楼层，分为碎没碎，碎了相当于dfs(i-1,j-1)的楼层，没碎则相当于从碎了那一层+1开始为0层，又可以扔：为dfs(i-1,j-1)+1+dfs(i-1,j)若楼层固定，则i次机会是最少的，dfs(i,j) = max(dfs(i-1,j-1),dfs(i-1,j-1)+1+dfs(i-1,j))】

[53. 最大子数组和](https://leetcode.cn/problems/maximum-subarray/)【dfs(int i)代表以i为结尾一定选的最大子数组和，不是】



### 以i为结尾一定选【子数组、子串问题】

###### 【53最大子数组和】【dfs(int i)代表以i为结尾一定选的最大子数组和，不是大于等于i一定选一个】

[5. 最长回文子串](https://leetcode.cn/problems/longest-palindromic-substring/)【boolean dp[i] [j] 代表ij是否是回文子串，先枚举j再枚举i，若s.charat(i) == s.charAt(j) 并且【 dp[i+1] [j-1]为true，或者j-i <= 2】，否则为false】

**连续子数组的最大和(二)**【记录curL，curR代表当前区间的l和r，若dp[i] = dp[i-1] + num[i]，则curR =i；若不是则curL=curR=i，gl和gr在dp[i]比ans更大时更新】

**NC127** **最长公共子串**【dfs（i，j）代表以 i 为结尾，以 j 为结尾一定选的最长长度，维护全局最长长度，最长的下标 i，最后返回 substring(i-ans+1,i+1)】

[230. 二叉搜索树中第 K 小的元素](https://leetcode.cn/problems/kth-smallest-element-in-a-bst/)

在当前代码中，`k` 是一个整型值，在 `dfs` 递归调用中不断地减小，但由于 Java 是按值传递参数，所以在递归过程中，`k` 的变化不会影响其他递归调用。这导致在每次递归中，`k` 的值始终是初始的，而不能正确地进行计数。

为了解决这个问题，我们需要使用一个类变量或使用引用来记录 `k` 的变化。我建议改用一个数组来保存 `k`，这样它的值可以在递归的过程中保持一致。

```
public class Solution {
    /**
     * 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
     *
     * 
     * @param proot TreeNode类 
     * @param k int整型 
     * @return int整型
     */
    int tar;
    public int KthNode (TreeNode proot, int k) {
        // write code here
        dfs(proot,k);
        return tar;
    }
    public void dfs(TreeNode root, int k) {
        if (root == null) return;
        dfs(root.left,k);
        k--;
        if (k == 0) {
            tar = root.val;
        }
        dfs(root.right,k);
    }
}
```

输入格式：

[a,b,c,d]

[java,c++,java,js]

### 排序

* [215. 数组中的第K个最大元素](https://leetcode.cn/problems/kth-largest-element-in-an-array/)【快速选择排序：int quickSort返回第k个元素，分为小于等于大于的情况，【都return】，【nextInt(r-l+1)生成[l,r]的元素当作枢纽，然后swap(l，位置)，最后将p放到l位置】】

归并排序[LCR 170. 交易逆序对的总数](https://leetcode.cn/problems/shu-zu-zhong-de-ni-xu-dui-lcof/)【归并排序 ，填充以及比较的都是temp，这个是原数组，合并后的直接填入nums,每次合并当nums[i] > nums[j]都有mid-i+1个逆序】

**JZ45** **把数组排成最小的数**【快速排序，quicksort 继续进行的规则是 l <= r(只有一个数也应该继续),partition 循坏继续进行的规则是 l < r（l 到 r 是需要 par 的区间，只有一个元素则不需 par）,大于等于不移动】【最小，排序规则 i + j  > j + i，则j排在i前】

**JZ41** **数据流中的中位数**【左大右小，左多右少，维护这个性质】

[164. 最大间距](https://leetcode.cn/problems/maximum-gap/)【桶排序，将值域max-min分为 n-1 份，  1.求桶大小 （max-min + n-2） / n-1 向下取整，2.求桶个数：（max-min）/d+1,2 3.将每个元素 ele放入 ele - min / d 的桶内，桶内元素各自排序】

### 位运算

**JZ65** **不用加减乘除做加法**【当前位的和 ^, 进位 &， a是当前和，b是进位，一直运算直到b为0】

gcd【a，b，a = b， b = a % b，直到b等于0，返回a】

**JZ16** **数值的整数次方**【将指数拆成二进制，每位是1的话base * base】

**JZ56** **数组中只出现一次的两个数字**【异或和找出两个数字的和，然后lowbit找出某位两个数字1个是1,1个是0，按这位分组异或和】

**JZ64** **求1+2+3+...+n**【利用短路运算特性，&&， a && b，a不成立b不执行，a||b，a成立b不执行，(n > 0) && (n += Sum_Solution(n-1)) > 0】



### \ + 需要转义的字符例如：‘，“，/

### **1. 转义字符分类**

#### **(1) 通用转义字符**

常用于表示不可直接书写的特殊字符。

| 转义字符 | 含义   | 示例代码                                   | 输出               |
| -------- | ------ | ------------------------------------------ | ------------------ |
| `\'`     | 单引号 | `System.out.println('\'');`                | `'`                |
| `\"`     | 双引号 | `System.out.println("\"Hello\"");`         | `"Hello"`          |
| `\\`     | 反斜杠 | `System.out.println("C:\\Program Files");` | `C:\Program Files` |



### **1. 正则表达式中的特殊字符**

正则表达式中有以下特殊字符需要转义：

| 特殊字符 | 描述                           | 转义方式    |
| -------- | ------------------------------ | ----------- |
| `.`      | 匹配任意单个字符（除了换行符） | `\\.`       |
| `*`      | 匹配前面的子表达式零次或多次   | `\\*`       |
| `+`      | 匹配前面的子表达式一次或多次   | `\\+`       |
| `?`      | 匹配前面的子表达式零次或一次   | `\\?`       |
| `^`      | 匹配输入的开始                 | `\\^`       |
| `$`      | 匹配输入的结尾                 | `\\$`       |
| `{` `}`  | 定义量词                       | `\\{` `\\}` |
| `[` `]`  | 定义字符类                     | `\\[` `\\]` |
| `(` `)`  | 定义分组                       | `\\(` `\\)` |
| `        | `                              | 表示“或”    |
| `\`      | 转义字符本身，需要双转义       | `\\\\`      |









【搜索提示词（sug词）评测规则】

一、sug提示词的定义与用途： 搜索提示词（sug词）是指针对用户搜索词的联想或补全，用于引导用户快速找到感兴趣的内容并进行点击。

二、评测维度： 在对sug提示词进行质量评测时，需从以下维度进行严格评测：

1. 单rank维度

- 字面质量问题（错别字、冗余、歧义、语义不明、特殊情况、口语化、倒装、大小写敏感、缺词漏字）
- 风控问题（涉政问题、违法违规、明确色情、映射色情、不良导向、谣言问题）
- 相关性（强相关、弱相关、不相关）

1. 整页维度（召回多个结果的考虑）

- 多样性（避免重复问题）
- 排序问题（确保最优质的内容优先呈现）

1. 时效性

- sug排序关注
- 搜索发现单rank关注

四、sug业务与其他内增业务问题处理优先级： 除sug业务外，其他内增业务所有评测维度出现问题均需多选，sug业务单选且需注意问题优先级。

- 搜索场景：
  - 增长词同时命中名人（大V、明星、历史人物等）和普通事件或普通用户，优先考虑名人和事件。
  - 示例：q=振兴，rs=振新中华，优先字面质量纠错处理，而非精准普通用户不相关（除非明确指向普通用户如振新中华718）。
- 非搜索场景：
  - 常见字面和用户名词可参考站内指向。
  - 示例：搜发=振新中华，站内历史相关可作字面质量问题处理。

五、问题优先级选择规则： 当出现多个评测维度的问题时，需按照以下优先级进行问题标记：

优先级最高：风控维度问题（涉政、违规、色情、不良导向、谣言）> 字面质量问题 > 相关性问题 > 重复性问题 > 排序问题 > 多样性问题

六、评测流程：

1. 分析用户query（q词）
2. 判断sug词和q词的相关性
3. 评估质量问题、排序问题、重复问题等
4. 根据规则中的优先级选择最严重的问题进行标记与反馈

## sug查询规则说明与应用示例

### 一、问题类型与相关性标准

#### 1. 包含q词可视为强相关

- 规则：当`sug`（suggestion）中包含用户输入的查询词（`q词`），或对`q词`进行了合理的谐音、拼音、数字、多音字延展，只要整体语义通顺、有实际意义且符合用户真实需求，即判定为强相关。
- 示例：
  - 正确示例：
    - `q=天津`，`sug=天津旅游攻略`
    - **说明**：`天津旅游攻略`包含`天津`（q词），且有明确实际意义，强相关。
    - **示例拓展**：`q=98`，`sug=酒吧视频`
    - **说明**：98谐音为“酒吧”，谐音明确且匹配实际语义，也视作强相关。
- 注意事项：
  - 谐音、拼音、多音字、数字扩展均视作正常延展，前提是语义完整明确。
  - 若出现明显的语义混乱或无意义则判定为不相关。

#### 2. 分词错误

- 规则：当`sug`中明显出现分词错误导致与`q`无明确关系时，应判为**不相关**。
- **反例**：
  - `q=海水`，`sug=青海水车`
    - **错误原因**：正确分词为`青海 水车`，而非`青 海水 车`。
- **优化建议**：
  - 分词出错时需优先纠正原意再进行相关性判断；明确断句、断词以确保准确性。

#### 3. 命中`q词`核心词即为强相关

- 规则：如`sug`命中了`q`词中的核心词，即视为强相关。
  - 注意**核心词数量**要求：
    - 若q词包含多个核心关键词：
      - 命中1个核心词可视为相关；
      - 3个及以上核心词命中2个及以上即可视为强相关； 若不足则为不相关或弱相关（视上下文确定）。
- **应用示例**：
  - `q=李白出装铭文推荐`，`sug=李白玩法教学`
    - **说明**：明确命中核心关键词`李白`，语义相关且有明确需求，视为强相关。
  - `q=赵云李白对局技巧`，`sug=王者荣耀攻略`
    - 若无具体命中词，判定弱相关或不相关。

### 二、关联匹配模式

#### 1. 联想匹配与宽泛认知匹配

- 规则：无核心词的宽泛查询句，`sug`命中外接关键成分或匹配字符数≥50%即可判定为强相关。
  - 中文、英文、数字混合整体视作一个整体字符。
- **应用示例**：
  - `q=喝表`，`sug=喝水时间表`
    - 有宽泛联想匹配（喝→喝水，表→时间表），判定强相关。
- **补充优化**：
  - 增加更清晰的字数计算与比例规则（建议明确为整体字符匹配达50%或以上即为强相关）。

#### 2. 弱相关及弱相关补充

- 规则：`sug`未达到强相关标准，但仍与`q`明显存在关联（联想意义、语境联系），且与`sug`存在的关联匹配程度较低但未完全无关，则视为弱相关补充。
- **应用示例**：
  - 查询词：`Leaf Recall ID`
  - sug结果未完全命中核心词，但字面上存在弱关联意义，判为弱相关。

#### 2. query为用户名（常见普通用户）

- 规则：当`sug`结果为用户姓名或用户名，默认判定为强相关，不需要进行严格匹配。
- **补充说明**：
  - 查询普通用户昵称或ID时，直接返回同名或关联的用户信息均为无问题。

#### 3. 拼音兜底匹配（特殊情况）

- 规则：当存在前缀不匹配但出现拼音兜底召回时，该类召回必须排在低位，允许适当放宽相关性判定。
- **示例**：
  - 查询`q=旅行剑`，召回`sug=旅行简谱`
    - 虽然命中拼音相似的`lvxingjian`，但字面不同，应排在末位，否则视为排序错误。

### 三、排序规则说明

#### 1. 时效词排序

- 规则：时效性词语在产生后的48小时内视为有效，应排在前3位，超时则不计时效。
  - 当查询本身为单字或拼音且需求模糊时，不严格考虑时效。

#### 2. how to词排序

- 规则：明确需求为问题类型的查询词，首个`how to`类结果应排在top7之内。
  - 若查询为拼音或单字，需求模糊则不强制排序。
- **示例**：
  - 查询`q=锅包肉`，若`sug=锅包肉教程`，必须在top7内；否则视为排序问题。

#### 3. sug用户卡排序

- 规则：用户卡片类`sug`结果不参与评分和排序控制，默认视作正常，不影响其他排序。

### 四、占位符问题

- 规则：当查询词（`q`）自身存在错误或明显占位导致结果有误时，返回占位符（将q词复制一遍）视作无问题。
- **示例**：
  - 查询`q=出粗车`，返回`sug=出粗车`无问题，即使查询本身有误，也视作占位符合理。

### 五、规则优化与应用说明

- **优化建议**：
  - 增加详细示例，特别是在分词错误和拼音匹配兜底规则中，强调场景化解释。
  - 进一步明确关键字的比例计算方法，减少误判可能。
  - 确保各类规则逻辑清晰，避免模糊地带。
- **应用说明**：
  - 规则应用时需整体考虑用户查询意图与语义，避免机械判断。
  - 特殊情况应增加备注和复审机制。

### 1. 字面质量

#### 错别字

- 专有名词、网络热梗、明确通用写法无需管控。
  - √ 可爱修勾
  - × 怕恰狗壁纸（应为“帕恰狗壁纸”）
- 非专有名词大小写不算字面质量。
- 数字表述需统一。
  - √ 520礼物
  - × 5二零礼物
- IP及剧名注意歧义与截断问题。
  - √ 少年派2
  - × 少年派2集（存在歧义）

#### 冗余

- 无意义字词赘述、多字多余或重复影响理解。
  - √ 笑的肚子疼的搞笑视频
  - × 男鞋鞋推荐

#### 口语化

- 明显口语表达或无意义语气词管控。
  - × 播放西游记
  - × 国庆去哪玩好呀
  - √ 笑死人的笑话

#### 倒装

- 仅管控影响理解的倒装。
  - × 强烈格子衫推荐（应为“强烈推荐的格子衫”）

#### 语义不明

- 含义完全无法辨识需管控。
  - × 啊啊啊啊啊啊
  - √ 梦人间饭（可查小众内容）

### 2. 风控问题

#### 涉政问题

- 管控攻击、诋毁、负面评价政治、民族内容。
  - × 共产党害人、中国政府腐败
  - √ 中国最近打仗吗（正常疑问豁免）

#### 违法违规

- 严肃违法道德：杀人、自残、拐卖、未成年人怀孕生子等。
  - × 小猪佩奇怀孕
  - √ 盲山女大学生被拐片段（真实新闻豁免）
- 财产安全（偷窃技巧、非法赚钱）
  - × 提现到微信（存在诈骗风险）
  - √ 赚钱小游戏（正常需求豁免）
- 违禁物品（枪械改装、毒品、违法药品）
  - × AK47改装
  - √ 弹弓（正常需求豁免）

#### 明确色情

- 直白色情内容严格管控。
  - × 狗交配

#### 映射色情

- 有联想色情或软色情内容需管控。
  - × 打扑克、QQ弹弹
  - √ 韩剧推荐

#### 不良导向

- 封建迷信、负面价值观、赌博变体。
  - × 冥婚、法轮功
  - × 打麻将口诀
  - √ 僵尸、灵异事件真实视频

#### 谣言问题

- 虚假死亡、虚假消息、明星虚假事件。
  - × 李双江葬礼（未死亡假消息）
  - √ 周杰伦和昆凌离婚了吗（疑问语气豁免）

### 3. 重复问题

- 明显重复、没有其他含义需管控。
  - × 送男朋友的礼物-送男朋友什么礼物（重复）
  - √ 防水防汗粉底液-不脱妆粉底液（不重复）

### 4. 时效性

- 已过时的热点、热点事件时效性差的需管控（超出一周），周期性热点（节日、发布会）同理。
  - × 成都发生地震（超过一周）
  - √ 成都地震进展
  - × 2025出苹果15发布会（已过去）

以上规则适用于搜索Sug评测的具体判定，可作为Coze平台知识库建立及大模型训练的规范标准。
