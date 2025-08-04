这段代码的主要功能是计算大脑皮层和皮层下区域的各种几何属性，包括区域中心、顶点法线、区域方向和区域面积。这些计算结果将被用于生成连接性和结构性数据，可能用于神经科学研究。以下是对代码的详细解释：

### 主程序部分

这段代码是放在 `if __name__ == '__main__':` 块中的，这意味着它是作为脚本执行的主要程序。

#### 读取顶点和三角形数据

```
python
复制代码
verts = np.loadtxt(os.path.join(PRD, SUBJ_ID, 'surface', 'vertices.txt'))
tri = np.loadtxt(os.path.join(PRD, SUBJ_ID, 'surface', 'triangles.txt'))
tri = tri.astype(int)
region_mapping = np.loadtxt(os.path.join(PRD, SUBJ_ID, 'surface', 'region_mapping.txt')).astype(int)
```

- `verts`：加载顶点数据，这些数据定义了大脑表面网格的顶点坐标。
- `tri`：加载三角形数据，这些数据定义了每个三角形的顶点索引。
- `region_mapping`：加载区域映射数据，这些数据定义了每个顶点所属的脑区域。

#### 读取连接性数据

```
python
复制代码
weights = np.loadtxt(os.path.join(PRD, 'connectivity', 'weights.csv'))
tract_lengths = np.loadtxt(os.path.join(PRD, 'connectivity', 'tract_lengths.csv'))
weights = weights + weights.transpose() - np.diag(np.diag(weights))
weights = np.vstack([np.zeros((1, weights.shape[0])), weights])
weights = np.hstack([np.zeros((weights.shape[0], 1)), weights])
tract_lengths = tract_lengths + tract_lengths.transpose()
tract_lengths = np.vstack([np.zeros((1, tract_lengths.shape[0])), tract_lengths]) 
tract_lengths = np.hstack([np.zeros((tract_lengths.shape[0], 1)), tract_lengths])
```

- `weights`：加载连接权重矩阵，表示大脑区域之间的连接强度。
- `tract_lengths`：加载轨迹长度矩阵，表示大脑区域之间的连接长度。
- 对 `weights` 和 `tract_lengths` 矩阵进行对称化处理，并添加零填充行和列，以确保矩阵的大小和结构正确。

#### 保存处理后的连接性数据

```
python
复制代码
np.savetxt(os.path.join(PRD, SUBJ_ID, 'connectivity', 'weights_num.txt'), weights, fmt='%d')
np.savetxt(os.path.join(PRD, SUBJ_ID, 'connectivity', 'tract_lengths.txt'), tract_lengths, fmt='%.3f')
```

- 将处理后的连接权重矩阵和轨迹长度矩阵保存到指定路径。

#### 读取区域名称

```
python
复制代码
list_name = np.loadtxt(open(os.path.join('share', 'reference_table_' + PARCEL + ".csv"), "r"), delimiter=",", skiprows=1, usecols=(1, ), dtype='str')
```

- `list_name`：加载参考表中的区域名称列表。

#### 计算区域中心

```
python
复制代码
centers = compute_region_center_cortex(verts, region_mapping, list_name)
```

- 调用 `compute_region_center_cortex` 函数，计算每个区域的中心坐标。

#### 计算顶点法线和区域方向

```
python
复制代码
number_of_vertices = int(verts.shape[0])
number_of_triangles = int(tri.shape[0])
vertex_triangles = compute_vertex_triangles(number_of_vertices, number_of_triangles, tri)
triangle_normals = compute_triangle_normals(tri, verts)
triangle_angles = compute_triangle_angles(verts, number_of_triangles, tri)
vertex_normals = compute_vertex_normals(number_of_vertices, vertex_triangles, tri, triangle_angles, triangle_normals, verts)
orientations = compute_region_orientation_cortex(vertex_normals, region_mapping, list_name)
```

- `number_of_vertices` 和 `number_of_triangles`：获取顶点和三角形的数量。
- `vertex_triangles`：计算每个顶点关联的三角形列表。
- `triangle_normals`：计算每个三角形的法线向量。
- `triangle_angles`：计算每个三角形的内角。
- `vertex_normals`：计算每个顶点的法线向量。
- `orientations`：计算每个区域的平均方向。

#### 计算三角形面积和区域面积

```
python
复制代码
triangle_areas = compute_triangle_areas(verts, tri)
areas = compute_region_areas_cortex(triangle_areas, vertex_triangles, region_mapping, list_name)
```

- `triangle_areas`：计算每个三角形的面积。
- `areas`：计算每个区域的总表面积。

### 计算函数的详细解释

#### 计算区域中心

```
python
复制代码
def compute_region_center_cortex(vertices, region_mapping, list_name):
    regions = np.unique(region_mapping)
    region_center = np.zeros((list_name.shape[0], 3))
    for k in regions:
        vert = vertices[region_mapping == k, :]
        region_center[k, :] = np.mean(vert, axis=0)
    return region_center
```

- 计算每个区域的中心坐标，方法是找到所有属于该区域的顶点，并计算它们的平均值。

#### 计算顶点法线

```
python
复制代码
def compute_vertex_triangles(number_of_vertices, number_of_triangles, triangles):
    vertex_triangles = [[] for _ in range(number_of_vertices)]
    for k in range(number_of_triangles):
        vertex_triangles[triangles[k, 0]].append(k)
        vertex_triangles[triangles[k, 1]].append(k)
        vertex_triangles[triangles[k, 2]].append(k)
    return vertex_triangles
```

- 计算每个顶点关联的三角形列表。

```
python
复制代码
def compute_triangle_normals(triangles, vertices):
    tri_u = vertices[triangles[:, 1], :] - vertices[triangles[:, 0], :]
    tri_v = vertices[triangles[:, 2], :] - vertices[triangles[:, 0], :]
    tri_norm = np.cross(tri_u, tri_v)
    try:
        triangle_normals = tri_norm / np.sqrt(np.sum(tri_norm ** 2, axis=1))[:, np.newaxis]
    except FloatingPointError:
        triangle_normals = tri_norm
    return triangle_normals
```

- 计算每个三角形的法线向量。

```
python
复制代码
def compute_triangle_angles(vertices, number_of_triangles, triangles):
    verts = vertices
    angles = np.zeros((number_of_triangles, 3))
    for tt in range(number_of_triangles):
        triangle = triangles[tt, :]
        for ta in range(3):
            ang = np.roll(triangle, -ta)
            angles[tt, ta] = np.arccos(np.dot(
                (verts[ang[1], :] - verts[ang[0], :]) /
                np.sqrt(np.sum((verts[ang[1], :] - verts[ang[0], :]) ** 2, axis=0)),
                (verts[ang[2], :] - verts[ang[0], :]) /
                np.sqrt(np.sum((verts[ang[2], :] - verts[ang[0], :]) ** 2, axis=0))))
    return angles
```

- 计算每个三角形的内角。

```
python
复制代码
def compute_vertex_normals(number_of_vertices, vertex_triangles, triangles, triangle_angles, triangle_normals, vertices):
    vert_norms = np.zeros((number_of_vertices, 3))
    bad_normal_count = 0
    for k in range(number_of_vertices):
        try:
            tri_list = list(vertex_triangles[k])
            angle_mask = triangles[tri_list, :] == k
            angles = triangle_angles[tri_list, :]
            angles = angles[angle_mask][:, np.newaxis]
            angle_scaling = angles / np.sum(angles, axis=0)
            vert_norms[k, :] = np.mean(angle_scaling * triangle_normals[tri_list, :], axis=0)
            vert_norms[k, :] = vert_norms[k, :] / np.sqrt(np.sum(vert_norms[k, :] ** 2, axis=0))
        except (ValueError, FloatingPointError):
            vert_norms[k, :] = vertices[k] / np.sqrt(vertices[k].dot(vertices[k]))
            bad_normal_count += 1
    if bad_normal_count:
        print(" %d vertices have bad normals" % bad_normal_count)
    return vert_norms
```

- 计算每个顶点的法线向量。

```
python
复制代码
def compute_region_orientation_cortex(vertex_normals, region_mapping, list_name):
    regions = np.unique(region_mapping)
    average_orientation = np.zeros((list_name.shape[0], 3))
    for k in regions:
        orient = vertex_normals[region_mapping == k, :]
        avg_orient = np.mean(orient, axis=0)
        average_orientation[k, :] = avg_orient / np.sqrt(np.sum(avg_orient ** 2))
    return average_orientation
```

- 计算每个区域的平均法线方向。

#### 计算区域面积

```
python
复制代码
def compute_triangle_areas(vertices, triangles):
    tri_u = vertices[triangles[:, 1], :] - vertices[triangles[:, 0], :]
    tri_v = vertices[triangles[:, 2], :] - vertices[triangles[:, 0], :]
    tri_norm = np.cross(tri_u, tri_v)
    triangle_areas = np.sqrt(np.sum(tri_norm ** 2, axis=1)) / 2.0
    triangle_areas = triangle_areas[:, np.newaxis]
    return triangle_areas
```

- 计算每个三角形的面积。

```
python
复制代码
def compute_region_areas_cortex(triangle_areas, vertex_triangles, region_mapping, list_name):
    regions = np.unique(region_mapping)
    region_surface_area = np.zeros((list_name.shape[0], 1))
    avt = np.array(vertex_triangles)
    for k in regions:
        regs = map(set, avt[region_mapping == k])
        region_triangles = set.union(*regs)
        region_surface_area[k] = triangle_areas[list(region_triangles)].sum()
    return region_surface_area
```

- 计算每个区域的总表面积。

### 总结

这段代码通过读取大脑皮层和皮层下区域的顶点和三角形数据，计算这些区域的各种几何属性（如区域中心、顶点法线、区域方向和区域面积），并将这些计算结果保存为文件。这些计算结果可以用于神经科学研究中的大脑结构和连接性的分析



这段代码的主要任务是处理和保存大脑的几何和连接性数据，包括区域表面积、平均方向、权重矩阵和区域中心坐标。以下是对这段代码的详细解释：

### 保存区域表面积和平均方向

```
python
复制代码
np.savetxt(os.path.join(PRD, SUBJ_ID, 'connectivity/areas.txt'), areas, fmt='%.2f')
np.savetxt(os.path.join(PRD, SUBJ_ID, 'connectivity/average_orientations.txt'), orientations, fmt='%.2f %.2f %.2f')
```

- `areas`：保存计算出的区域表面积，每个区域一个值。
- `orientations`：保存计算出的区域平均方向，每个区域一个三维向量。

### 加载区域表面积和权重矩阵

```
python
复制代码
areas = np.loadtxt(os.path.join(PRD, SUBJ_ID, 'connectivity/areas.txt'))
weights_num = np.loadtxt(os.path.join(PRD, SUBJ_ID, 'connectivity/weights_num.txt'))
weights_num = weights_num.astype(float)
```

- `areas`：重新加载之前保存的区域表面积。
- `weights_num`：加载连接权重矩阵，并将其转换为浮点型。

### 归一化权重矩阵

```
python
复制代码
weights_cor = np.zeros_like(weights_num)  # 初始化归一化的权重矩阵

for ta in range(len(weights_num)):
    for tb in range(len(weights_num[0])):
        if (areas[ta] + areas[tb] != 0):
            weights_cor[ta][tb] = weights_num[ta][tb] / (areas[ta] + areas[tb])
        else:
            weights_cor[ta][tb] = 0.0
```

- 通过将每个权重值除以对应区域的表面积之和，对权重矩阵进行归一化。
- `weights_cor` 是归一化后的权重矩阵。

### 保存归一化后的权重矩阵

```
python
复制代码
np.savetxt(os.path.join(PRD, SUBJ_ID, 'connectivity/weights_cor.txt'), weights_cor, fmt='%.4f')
```

- 将归一化后的权重矩阵保存到文件。

### 进一步归一化权重矩阵

```
python
复制代码
for ta in range(len(weights_cor)):
    for tb in range(len(weights_cor[0])):
        if weights_cor[ta].max() != 0:
            weights_num[ta][tb] = weights_cor[ta][tb] / weights_cor[ta].max()
```

- 将归一化后的权重矩阵的每一行的最大值作为分母，对每一行的权重值进行进一步归一化。

### 保存最终的权重矩阵

```
python
复制代码
np.savetxt(os.path.join(PRD, SUBJ_ID, 'connectivity/weights.txt'), weights_num, fmt='%.3f')
```

- 将最终归一化后的权重矩阵保存到文件。

### 保存区域中心坐标

```
python
复制代码
f = open(os.path.join(PRD, SUBJ_ID, 'connectivity/centres.txt'), 'w')
for i, name in enumerate(list_name):
    f.write(str(name) + ' ')
    for j in range(3):
        f.write('{:.4f} '.format(centers[i, j]))
    f.write('\n')
f.close()
```

- 将每个区域的中心坐标和名称保存到文件。
- 文件格式为：区域名称 后跟 3 个浮点数（中心坐标）。

### 保存皮层区域标签

```
python
复制代码
ref_table = np.loadtxt(open(os.path.join('share', 'reference_table_' + PARCEL + ".csv"), "rb"), delimiter=",", skiprows=1, usecols=(7,))
np.savetxt(os.path.join(PRD, SUBJ_ID, 'connectivity/cortical.txt'), ref_table, fmt='%d')
```

- 加载参考表中的第 7 列数据（皮层区域标签），并保存到文件。

### 总结

这段代码主要执行以下任务：

1. **保存计算结果**：包括区域表面积、平均方向、和权重矩阵。
2. **归一化权重矩阵**：通过两步归一化过程，使权重矩阵标准化。
3. **保存区域中心坐标**：以文本格式保存每个区域的中心坐标。
4. **保存皮层区域标签**：从参考表中加载并保存皮层区域标签。

这些计算和保存的结果可以用于后续的神经科学分析，如研究大脑结构和功能之间的关系。





```
import nibabel as nib
import numpy as np

# 文件路径
atlas_file = '246.nii'  # 246 图谱的 NIfTI 文件
vertices_file = 'vertices.txt'
triangles_file = 'triangles.txt'
output_file = 'areas.txt'

# 加载 246 图谱的 NIfTI 文件
atlas_nii = nib.load(atlas_file)
atlas_data = atlas_nii.get_fdata()
atlas_shape = atlas_data.shape

# 加载顶点和三角形数据
vertices = np.loadtxt(vertices_file)
triangles = np.loadtxt(triangles_file, dtype=int)

# 获取 NIfTI 图像的仿射矩阵，用于将顶点坐标转换为体素坐标
affine = atlas_nii.affine

# 将顶点坐标转换为体素坐标
def get_voxel_coords(vertices, affine, atlas_shape):
    # 添加一个维度，进行坐标变换
    vertices_h = np.hstack((vertices, np.ones((vertices.shape[0], 1))))
    voxel_coords = np.dot(np.linalg.inv(affine), vertices_h.T).T[:, :3]
    # 裁剪坐标，确保在合法范围内
    voxel_coords = np.clip(voxel_coords, [0, 0, 0], np.array(atlas_shape) - 1)
    return np.round(voxel_coords).astype(int)

voxel_coords = get_voxel_coords(vertices, affine, atlas_shape)

def compute_triangle_areas(vertices, triangles):
    tri_areas = np.zeros(triangles.shape[0])
    for i, tri in enumerate(triangles):
        v0, v1, v2 = vertices[tri[0]], vertices[tri[1]], vertices[tri[2]]
        tri_areas[i] = 0.5 * np.linalg.norm(np.cross(v1 - v0, v2 - v0))
    return tri_areas

# 计算每个三角形的面积
triangle_areas = compute_triangle_areas(vertices, triangles)

# 获取顶点的标签
vertex_labels = np.zeros(vertices.shape[0], dtype=int)
for i, tri in enumerate(triangles):
    v0, v1, v2 = voxel_coords[tri[0]], voxel_coords[tri[1]], voxel_coords[tri[2]]
    labels = [atlas_data[tuple(v0)], atlas_data[tuple(v1)], atlas_data[tuple(v2)]]
    label = int(np.median(labels))
    vertex_labels[tri[0]] = label
    vertex_labels[tri[1]] = label
    vertex_labels[tri[2]] = label

# 计算每个区域的总面积
areas = np.zeros(int(np.max(vertex_labels)) + 1)
for i, area in enumerate(triangle_areas):
    v0, v1, v2 = triangles[i]
    label = vertex_labels[v0]
    areas[label] += area

# 跳过背景区域
valid_areas = areas[1:]

# 保存区域面积到文件
with open(output_file, 'w') as f:
    for label, area in enumerate(valid_areas, start=1):
        f.write(f'{label} {area:.6f}\n')

print(f"Saved areas to {output_file}")
```

### 代码解释

1. **获取 NIfTI 图像的仿射矩阵**：
   - 使用 `atlas_nii.affine` 获取仿射矩阵，用于将顶点坐标转换为体素坐标。
2. **将顶点坐标转换为体素坐标**：
   - 定义 `get_voxel_coords` 函数，将顶点坐标转换为体素坐标。
   - 使用 `np.clip` 确保转换后的坐标在合法范围内。
3. **计算每个三角形的面积**：
   - 定义 `compute_triangle_areas` 函数，计算每个三角形的面积。
4. **获取顶点的标签**：
   - 根据转换后的体素坐标从 `atlas_data` 中获取标签。
5. **计算每个区域的总面积**：
   - 将每个顶点的三角形面积累加到相应的区域，得到每个区域的总面积。
6. **保存结果**：
   - 将计算得到的每个区域的总面积保存到 `areas.txt` 文件中。

通过这些步骤，你可以从 246 图谱的 NIfTI 文件生成每个区域的总面积，并保存结果到 `areas.txt` 