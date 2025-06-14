# 03-实时渲染系统

 (Real-time Rendering System)

## 目录

- [03-实时渲染系统](#03-实时渲染系统)
  - [目录](#目录)
  - [1. 概述](#1-概述)
    - [1.1 核心功能](#11-核心功能)
    - [1.2 设计原则](#12-设计原则)
  - [2. 形式化定义](#2-形式化定义)
    - [2.1 渲染系统](#21-渲染系统)
    - [2.2 光照模型](#22-光照模型)
  - [3. 数学基础](#3-数学基础)
    - [3.1 线性代数](#31-线性代数)
    - [3.2 投影变换](#32-投影变换)
    - [3.3 光照计算](#33-光照计算)
  - [4. 渲染管线](#4-渲染管线)
    - [4.1 顶点着色器](#41-顶点着色器)
    - [4.2 片段着色器](#42-片段着色器)
    - [4.3 渲染器](#43-渲染器)
    - [4.4 材质系统](#44-材质系统)
    - [4.5 网格系统](#45-网格系统)
  - [5. Go语言实现](#5-go语言实现)
    - [5.1 渲染管理器](#51-渲染管理器)
    - [5.2 后处理系统](#52-后处理系统)
    - [5.3 阴影系统](#53-阴影系统)
  - [6. 性能优化](#6-性能优化)
    - [6.1 视锥体剔除](#61-视锥体剔除)
    - [6.2 LOD系统](#62-lod系统)
    - [6.3 实例化渲染](#63-实例化渲染)
  - [7. 总结](#7-总结)
    - [7.1 关键特性](#71-关键特性)
    - [7.2 扩展方向](#72-扩展方向)

## 1. 概述

实时渲染系统是游戏引擎的核心组件，负责将3D场景转换为2D图像并显示在屏幕上。现代实时渲染系统需要支持复杂的光照、材质、阴影等效果，同时保证60FPS以上的渲染性能。

### 1.1 核心功能

- **几何处理**: 顶点变换、裁剪、光栅化
- **光照计算**: 环境光、漫反射、镜面反射
- **材质渲染**: 纹理映射、法线贴图、PBR材质
- **阴影渲染**: 阴影映射、软阴影、实时阴影
- **后处理**: 抗锯齿、景深、运动模糊

### 1.2 设计原则

- **高性能**: 保证实时渲染性能
- **高质量**: 支持现代渲染效果
- **可扩展**: 支持自定义着色器
- **跨平台**: 支持多种图形API

## 2. 形式化定义

### 2.1 渲染系统

**定义 2.1.1** (渲染系统)
渲染系统是一个五元组 $R = (G, M, L, S, P)$，其中：

- $G$ 是几何数据 (Geometry Data)
- $M$ 是材质系统 (Material System)
- $L$ 是光照系统 (Lighting System)
- $S$ 是着色器系统 (Shader System)
- $P$ 是后处理系统 (Post-processing System)

**定义 2.1.2** (渲染管线)
渲染管线是一个函数序列 $P = (f_1, f_2, ..., f_n)$，其中每个 $f_i: V_i \rightarrow V_{i+1}$ 是渲染阶段函数。

### 2.2 光照模型

**定义 2.2.1** (Phong光照模型)
Phong光照模型定义为：
$$I = I_a + I_d + I_s$$

其中：

- $I_a = k_a \cdot A$ 是环境光
- $I_d = k_d \cdot (L \cdot N) \cdot D$ 是漫反射
- $I_s = k_s \cdot (R \cdot V)^n \cdot S$ 是镜面反射

**定义 2.2.2** (PBR光照模型)
基于物理的渲染模型定义为：
$$f_r = \frac{DFG}{4(N \cdot L)(N \cdot V)}$$

其中 $D$ 是法线分布函数，$F$ 是菲涅尔函数，$G$ 是几何函数。

## 3. 数学基础

### 3.1 线性代数

**定理 3.1.1** (矩阵变换的可逆性)
对于非奇异矩阵 $M$，存在逆矩阵 $M^{-1}$ 使得 $M \cdot M^{-1} = I$。

**证明**:
设 $M$ 是 $n \times n$ 非奇异矩阵，则 $\det(M) \neq 0$。
根据矩阵求逆公式：
$$M^{-1} = \frac{1}{\det(M)} \cdot \text{adj}(M)$$

其中 $\text{adj}(M)$ 是 $M$ 的伴随矩阵。

### 3.2 投影变换

**定义 3.2.1** (透视投影)
透视投影矩阵 $P$ 定义为：
$$P = \begin{bmatrix}
\frac{2n}{r-l} & 0 & \frac{r+l}{r-l} & 0 \\
0 & \frac{2n}{t-b} & \frac{t+b}{t-b} & 0 \\
0 & 0 & -\frac{f+n}{f-n} & -\frac{2fn}{f-n} \\
0 & 0 & -1 & 0
\end{bmatrix}$$

其中 $n, f$ 是近远平面，$l, r, t, b$ 是视锥体边界。

**定理 3.2.1** (投影变换保持直线)
投影变换将直线映射为直线。

### 3.3 光照计算

**定义 3.3.1** (法线分布函数)
GGX法线分布函数定义为：
$$D(h) = \frac{\alpha^2}{\pi((N \cdot H)^2(\alpha^2-1)+1)^2}$$

其中 $\alpha$ 是粗糙度参数。

**定义 3.3.2** (菲涅尔函数)
Schlick近似菲涅尔函数定义为：
$$F(v, h) = F_0 + (1-F_0)(1-(v \cdot h))^5$$

其中 $F_0$ 是基础反射率。

## 4. 渲染管线

### 4.1 顶点着色器

```go
// 顶点着色器
type VertexShader struct {
    program     uint32
    attributes  map[string]int32
    uniforms    map[string]int32
}

// 顶点数据
type Vertex struct {
    Position    mgl32.Vec3
    Normal      mgl32.Vec3
    TexCoord    mgl32.Vec2
    Tangent     mgl32.Vec3
    Bitangent   mgl32.Vec3
}

// 顶点着色器程序
const vertexShaderSource = `
# version 330 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aNormal;
layout (location = 2) in vec2 aTexCoord;
layout (location = 3) in vec3 aTangent;
layout (location = 4) in vec3 aBitangent;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

out vec3 FragPos;
out vec3 Normal;
out vec2 TexCoord;
out mat3 TBN;

void main() {
    FragPos = vec3(model * vec4(aPos, 1.0));
    Normal = mat3(transpose(inverse(model))) * aNormal;
    TexCoord = aTexCoord;

    vec3 T = normalize(vec3(model * vec4(aTangent, 0.0)));
    vec3 N = normalize(vec3(model * vec4(aNormal, 0.0)));
    T = normalize(T - dot(T, N) * N);
    vec3 B = cross(N, T);
    TBN = mat3(T, B, N);

    gl_Position = projection * view * vec4(FragPos, 1.0);
}
`

// 编译顶点着色器
func (vs *VertexShader) Compile() error {
    vs.program = gl.CreateProgram()

    vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
    gl.ShaderSource(vertexShader, vertexShaderSource)
    gl.CompileShader(vertexShader)

    // 检查编译错误
    var success int32
    gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &success)
    if success == 0 {
        var logLength int32
        gl.GetShaderiv(vertexShader, gl.INFO_LOG_LENGTH, &logLength)
        log := gl.GetShaderInfoLog(vertexShader)
        return fmt.Errorf("vertex shader compilation failed: %s", log)
    }

    gl.AttachShader(vs.program, vertexShader)
    gl.LinkProgram(vs.program)

    return nil
}
```

### 4.2 片段着色器

```go
// 片段着色器
type FragmentShader struct {
    program     uint32
    uniforms    map[string]int32
    textures    map[string]uint32
}

// PBR片段着色器
const fragmentShaderSource = `
# version 330 core
out vec4 FragColor;

in vec3 FragPos;
in vec3 Normal;
in vec2 TexCoord;
in mat3 TBN;

// 材质属性
uniform sampler2D albedoMap;
uniform sampler2D normalMap;
uniform sampler2D metallicMap;
uniform sampler2D roughnessMap;
uniform sampler2D aoMap;

// 光照属性
uniform vec3 lightPositions[4];
uniform vec3 lightColors[4];
uniform vec3 viewPos;

const float PI = 3.14159265359;

// PBR函数
float DistributionGGX(vec3 N, vec3 H, float roughness) {
    float a = roughness * roughness;
    float a2 = a * a;
    float NdotH = max(dot(N, H), 0.0);
    float NdotH2 = NdotH * NdotH;

    float nom   = a2;
    float denom = (NdotH2 * (a2 - 1.0) + 1.0);
    denom = PI * denom * denom;

    return nom / denom;
}

float GeometrySchlickGGX(float NdotV, float roughness) {
    float r = (roughness + 1.0);
    float k = (r * r) / 8.0;

    float nom   = NdotV;
    float denom = NdotV * (1.0 - k) + k;

    return nom / denom;
}

float GeometrySmith(vec3 N, vec3 V, vec3 L, float roughness) {
    float NdotV = max(dot(N, V), 0.0);
    float NdotL = max(dot(N, L), 0.0);
    float ggx2 = GeometrySchlickGGX(NdotV, roughness);
    float ggx1 = GeometrySchlickGGX(NdotL, roughness);

    return ggx1 * ggx2;
}

vec3 fresnelSchlick(float cosTheta, vec3 F0) {
    return F0 + (1.0 - F0) * pow(clamp(1.0 - cosTheta, 0.0, 1.0), 5.0);
}

void main() {
    // 获取材质属性
    vec3 albedo = pow(texture(albedoMap, TexCoord).rgb, vec3(2.2));
    float metallic = texture(metallicMap, TexCoord).r;
    float roughness = texture(roughnessMap, TexCoord).r;
    float ao = texture(aoMap, TexCoord).r;

    // 计算法线
    vec3 N = normalize(TBN * (texture(normalMap, TexCoord).rgb * 2.0 - 1.0));
    vec3 V = normalize(viewPos - FragPos);

    // 计算反射率
    vec3 F0 = vec3(0.04);
    F0 = mix(F0, albedo, metallic);

    // 反射方程
    vec3 Lo = vec3(0.0);
    for(int i = 0; i < 4; ++i) {
        vec3 L = normalize(lightPositions[i] - FragPos);
        vec3 H = normalize(V + L);
        float distance = length(lightPositions[i] - FragPos);
        float attenuation = 1.0 / (distance * distance);
        vec3 radiance = lightColors[i] * attenuation;

        // Cook-Torrance BRDF
        float NDF = DistributionGGX(N, H, roughness);
        float G   = GeometrySmith(N, V, L, roughness);
        vec3 F    = fresnelSchlick(max(dot(H, V), 0.0), F0);

        vec3 numerator    = NDF * G * F;
        float denominator = 4.0 * max(dot(N, V), 0.0) * max(dot(N, L), 0.0) + 0.0001;
        vec3 specular = numerator / denominator;

        vec3 kS = F;
        vec3 kD = vec3(1.0) - kS;
        kD *= 1.0 - metallic;

        float NdotL = max(dot(N, L), 0.0);

        Lo += (kD * albedo / PI + specular) * radiance * NdotL;
    }

    // 环境光照
    vec3 ambient = vec3(0.03) * albedo * ao;
    vec3 color = ambient + Lo;

    // HDR色调映射
    color = color / (color + vec3(1.0));
    // gamma校正
    color = pow(color, vec3(1.0/2.2));

    FragColor = vec4(color, 1.0);
}
`

// 编译片段着色器
func (fs *FragmentShader) Compile() error {
    fs.program = gl.CreateProgram()

    fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
    gl.ShaderSource(fragmentShader, fragmentShaderSource)
    gl.CompileShader(fragmentShader)

    // 检查编译错误
    var success int32
    gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &success)
    if success == 0 {
        var logLength int32
        gl.GetShaderiv(fragmentShader, gl.INFO_LOG_LENGTH, &logLength)
        log := gl.GetShaderInfoLog(fragmentShader)
        return fmt.Errorf("fragment shader compilation failed: %s", log)
    }

    gl.AttachShader(fs.program, fragmentShader)
    gl.LinkProgram(fs.program)

    return nil
}
```

### 4.3 渲染器

```go
// 渲染器
type Renderer struct {
    window      *glfw.Window
    vertexShader *VertexShader
    fragmentShader *FragmentShader
    camera      *Camera
    scene       *Scene
    lights      []*Light
    materials   map[string]*Material
    meshes      map[string]*Mesh
    textures    map[string]*Texture
}

// 相机
type Camera struct {
    position    mgl32.Vec3
    target      mgl32.Vec3
    up          mgl32.Vec3
    fov         float32
    aspect      float32
    near        float32
    far         float32
}

// 获取视图矩阵
func (c *Camera) GetViewMatrix() mgl32.Mat4 {
    return mgl32.LookAtV(c.position, c.target, c.up)
}

// 获取投影矩阵
func (c *Camera) GetProjectionMatrix() mgl32.Mat4 {
    return mgl32.Perspective(c.fov, c.aspect, c.near, c.far)
}

// 场景
type Scene struct {
    objects     []*GameObject
    lights      []*Light
    skybox      *Skybox
}

// 游戏对象
type GameObject struct {
    transform   *Transform
    mesh        *Mesh
    material    *Material
    visible     bool
}

// 渲染场景
func (r *Renderer) RenderScene() {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    // 设置相机矩阵
    viewMatrix := r.camera.GetViewMatrix()
    projectionMatrix := r.camera.GetProjectionMatrix()

    // 渲染天空盒
    if r.scene.skybox != nil {
        r.renderSkybox(r.scene.skybox, viewMatrix, projectionMatrix)
    }

    // 渲染所有对象
    for _, obj := range r.scene.objects {
        if obj.visible {
            r.renderObject(obj, viewMatrix, projectionMatrix)
        }
    }

    glfw.SwapBuffers(r.window)
}

// 渲染对象
func (r *Renderer) renderObject(obj *GameObject, viewMatrix, projectionMatrix mgl32.Mat4) {
    // 使用着色器程序
    gl.UseProgram(r.fragmentShader.program)

    // 设置变换矩阵
    modelMatrix := obj.transform.GetMatrix()
    gl.UniformMatrix4fv(gl.GetUniformLocation(r.fragmentShader.program, gl.Str("model\x00")), 1, false, &modelMatrix[0])
    gl.UniformMatrix4fv(gl.GetUniformLocation(r.fragmentShader.program, gl.Str("view\x00")), 1, false, &viewMatrix[0])
    gl.UniformMatrix4fv(gl.GetUniformLocation(r.fragmentShader.program, gl.Str("projection\x00")), 1, false, &projectionMatrix[0])

    // 设置相机位置
    gl.Uniform3fv(gl.GetUniformLocation(r.fragmentShader.program, gl.Str("viewPos\x00")), 1, &r.camera.position[0])

    // 设置光照
    r.setupLights()

    // 绑定材质
    obj.material.Bind(r.fragmentShader.program)

    // 渲染网格
    obj.mesh.Draw()
}

// 设置光照
func (r *Renderer) setupLights() {
    for i, light := range r.lights {
        if i >= 4 { // 最多4个光源
            break
        }

        posName := fmt.Sprintf("lightPositions[%d]", i)
        colorName := fmt.Sprintf("lightColors[%d]", i)

        gl.Uniform3fv(gl.GetUniformLocation(r.fragmentShader.program, gl.Str(posName+"\x00")), 1, &light.position[0])
        gl.Uniform3fv(gl.GetUniformLocation(r.fragmentShader.program, gl.Str(colorName+"\x00")), 1, &light.color[0])
    }
}
```

### 4.4 材质系统

```go
// 材质
type Material struct {
    albedoMap   *Texture
    normalMap   *Texture
    metallicMap *Texture
    roughnessMap *Texture
    aoMap       *Texture
    albedo      mgl32.Vec3
    metallic    float32
    roughness   float32
    ao          float32
}

// 绑定材质
func (m *Material) Bind(program uint32) {
    // 绑定纹理
    if m.albedoMap != nil {
        gl.ActiveTexture(gl.TEXTURE0)
        gl.BindTexture(gl.TEXTURE_2D, m.albedoMap.id)
        gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("albedoMap\x00")), 0)
    }

    if m.normalMap != nil {
        gl.ActiveTexture(gl.TEXTURE1)
        gl.BindTexture(gl.TEXTURE_2D, m.normalMap.id)
        gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("normalMap\x00")), 1)
    }

    if m.metallicMap != nil {
        gl.ActiveTexture(gl.TEXTURE2)
        gl.BindTexture(gl.TEXTURE_2D, m.metallicMap.id)
        gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("metallicMap\x00")), 2)
    }

    if m.roughnessMap != nil {
        gl.ActiveTexture(gl.TEXTURE3)
        gl.BindTexture(gl.TEXTURE_2D, m.roughnessMap.id)
        gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("roughnessMap\x00")), 3)
    }

    if m.aoMap != nil {
        gl.ActiveTexture(gl.TEXTURE4)
        gl.BindTexture(gl.TEXTURE_2D, m.aoMap.id)
        gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("aoMap\x00")), 4)
    }

    // 设置材质属性
    gl.Uniform3fv(gl.GetUniformLocation(program, gl.Str("albedo\x00")), 1, &m.albedo[0])
    gl.Uniform1f(gl.GetUniformLocation(program, gl.Str("metallic\x00")), m.metallic)
    gl.Uniform1f(gl.GetUniformLocation(program, gl.Str("roughness\x00")), m.roughness)
    gl.Uniform1f(gl.GetUniformLocation(program, gl.Str("ao\x00")), m.ao)
}
```

### 4.5 网格系统

```go
// 网格
type Mesh struct {
    VAO         uint32
    VBO         uint32
    EBO         uint32
    indices     []uint32
    vertexCount int32
}

// 创建网格
func NewMesh(vertices []Vertex, indices []uint32) *Mesh {
    mesh := &Mesh{
        indices:     indices,
        vertexCount: int32(len(indices)),
    }

    // 生成VAO和VBO
    gl.GenVertexArrays(1, &mesh.VAO)
    gl.GenBuffers(1, &mesh.VBO)
    gl.GenBuffers(1, &mesh.EBO)

    gl.BindVertexArray(mesh.VAO)

    // 绑定顶点数据
    gl.BindBuffer(gl.ARRAY_BUFFER, mesh.VBO)
    gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*int(unsafe.Sizeof(Vertex{})), gl.Ptr(vertices), gl.STATIC_DRAW)

    // 绑定索引数据
    gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.EBO)
    gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

    // 设置顶点属性
    // 位置
    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(unsafe.Sizeof(Vertex{})), gl.PtrOffset(0))
    gl.EnableVertexAttribArray(0)

    // 法线
    gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(unsafe.Sizeof(Vertex{})), gl.PtrOffset(unsafe.Sizeof(mgl32.Vec3{})))
    gl.EnableVertexAttribArray(1)

    // 纹理坐标
    gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(unsafe.Sizeof(Vertex{})), gl.PtrOffset(2*unsafe.Sizeof(mgl32.Vec3{})))
    gl.EnableVertexAttribArray(2)

    // 切线
    gl.VertexAttribPointer(3, 3, gl.FLOAT, false, int32(unsafe.Sizeof(Vertex{})), gl.PtrOffset(2*unsafe.Sizeof(mgl32.Vec3{})+unsafe.Sizeof(mgl32.Vec2{})))
    gl.EnableVertexAttribArray(3)

    // 副切线
    gl.VertexAttribPointer(4, 3, gl.FLOAT, false, int32(unsafe.Sizeof(Vertex{})), gl.PtrOffset(3*unsafe.Sizeof(mgl32.Vec3{})+unsafe.Sizeof(mgl32.Vec2{})))
    gl.EnableVertexAttribArray(4)

    gl.BindVertexArray(0)

    return mesh
}

// 渲染网格
func (m *Mesh) Draw() {
    gl.BindVertexArray(m.VAO)
    gl.DrawElements(gl.TRIANGLES, m.vertexCount, gl.UNSIGNED_INT, nil)
    gl.BindVertexArray(0)
}
```

## 5. Go语言实现

### 5.1 渲染管理器

```go
// 渲染管理器
type RenderManager struct {
    renderer    *Renderer
    renderQueue []RenderCommand
    postProcess *PostProcessor
    mutex       sync.Mutex
}

// 渲染命令
type RenderCommand struct {
    object      *GameObject
    priority    int
    shader      uint32
    material    *Material
}

// 添加渲染命令
func (rm *RenderManager) AddRenderCommand(cmd RenderCommand) {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    rm.renderQueue = append(rm.renderQueue, cmd)
}

// 排序渲染队列
func (rm *RenderManager) SortRenderQueue() {
    sort.Slice(rm.renderQueue, func(i, j int) bool {
        return rm.renderQueue[i].priority < rm.renderQueue[j].priority
    })
}

// 执行渲染
func (rm *RenderManager) Render() {
    rm.mutex.Lock()
    defer rm.mutex.Unlock()

    // 排序渲染队列
    rm.SortRenderQueue()

    // 执行渲染命令
    for _, cmd := range rm.renderQueue {
        rm.renderer.renderCommand(cmd)
    }

    // 后处理
    rm.postProcess.Process()

    // 清空渲染队列
    rm.renderQueue = rm.renderQueue[:0]
}
```

### 5.2 后处理系统

```go
// 后处理器
type PostProcessor struct {
    framebuffer *Framebuffer
    shaders     map[string]*Shader
    effects     []PostProcessEffect
}

// 后处理效果
type PostProcessEffect interface {
    Apply(input, output *Framebuffer)
    GetName() string
}

// 抗锯齿效果
type AntiAliasingEffect struct {
    shader *Shader
}

func (aa *AntiAliasingEffect) Apply(input, output *Framebuffer) {
    gl.BindFramebuffer(gl.FRAMEBUFFER, output.id)
    gl.Clear(gl.COLOR_BUFFER_BIT)

    aa.shader.Use()
    gl.BindTexture(gl.TEXTURE_2D, input.colorTexture)
    gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func (aa *AntiAliasingEffect) GetName() string {
    return "AntiAliasing"
}

// 景深效果
type DepthOfFieldEffect struct {
    shader *Shader
    focus  float32
    range  float32
}

func (dof *DepthOfFieldEffect) Apply(input, output *Framebuffer) {
    gl.BindFramebuffer(gl.FRAMEBUFFER, output.id)
    gl.Clear(gl.COLOR_BUFFER_BIT)

    dof.shader.Use()
    gl.Uniform1f(gl.GetUniformLocation(dof.shader.program, gl.Str("focus\x00")), dof.focus)
    gl.Uniform1f(gl.GetUniformLocation(dof.shader.program, gl.Str("range\x00")), dof.range)

    gl.BindTexture(gl.TEXTURE_2D, input.colorTexture)
    gl.BindTexture(gl.TEXTURE_2D, input.depthTexture)
    gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func (dof *DepthOfFieldEffect) GetName() string {
    return "DepthOfField"
}

// 应用后处理效果
func (pp *PostProcessor) Process() {
    input := pp.framebuffer
    output := NewFramebuffer(input.width, input.height)

    for _, effect := range pp.effects {
        effect.Apply(input, output)
        input, output = output, input
    }

    // 最终输出到屏幕
    gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
    gl.Clear(gl.COLOR_BUFFER_BIT)

    // 使用简单的着色器渲染最终结果
    finalShader := pp.shaders["final"]
    finalShader.Use()
    gl.BindTexture(gl.TEXTURE_2D, input.colorTexture)
    gl.DrawArrays(gl.TRIANGLES, 0, 6)
}
```

### 5.3 阴影系统

```go
// 阴影系统
type ShadowSystem struct {
    shadowMap   *Framebuffer
    lightSpaceMatrix mgl32.Mat4
    shader      *Shader
}

// 创建阴影映射
func (ss *ShadowSystem) CreateShadowMap(light *Light, scene *Scene) {
    // 设置阴影映射帧缓冲
    gl.BindFramebuffer(gl.FRAMEBUFFER, ss.shadowMap.id)
    gl.Clear(gl.DEPTH_BUFFER_BIT)

    // 计算光源空间矩阵
    ss.lightSpaceMatrix = ss.calculateLightSpaceMatrix(light)

    // 使用阴影着色器
    ss.shader.Use()
    gl.UniformMatrix4fv(gl.GetUniformLocation(ss.shader.program, gl.Str("lightSpaceMatrix\x00")), 1, false, &ss.lightSpaceMatrix[0])

    // 渲染场景到阴影映射
    for _, obj := range scene.objects {
        if obj.castShadow {
            ss.renderShadowObject(obj)
        }
    }

    gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

// 计算光源空间矩阵
func (ss *ShadowSystem) calculateLightSpaceMatrix(light *Light) mgl32.Mat4 {
    // 计算光源视图矩阵
    lightView := mgl32.LookAtV(light.position, light.position.Add(light.direction), mgl32.Vec3{0, 1, 0})

    // 计算光源投影矩阵
    lightProjection := mgl32.Ortho(-10, 10, -10, 10, 0.1, 100)

    return lightProjection.Mul4(lightView)
}

// 渲染阴影对象
func (ss *ShadowSystem) renderShadowObject(obj *GameObject) {
    modelMatrix := obj.transform.GetMatrix()
    gl.UniformMatrix4fv(gl.GetUniformLocation(ss.shader.program, gl.Str("model\x00")), 1, false, &modelMatrix[0])

    obj.mesh.Draw()
}
```

## 6. 性能优化

### 6.1 视锥体剔除

**定理 6.1.1** (视锥体剔除)
对于视锥体 $F$ 和包围盒 $B$，如果 $B \cap F = \emptyset$，则 $B$ 中的对象不需要渲染。

**实现**:
```go
// 视锥体
type Frustum struct {
    planes [6]mgl32.Vec4
}

// 包围盒
type BoundingBox struct {
    min mgl32.Vec3
    max mgl32.Vec3
}

// 视锥体剔除
func (f *Frustum) Contains(box *BoundingBox) bool {
    for _, plane := range f.planes {
        // 计算包围盒在平面法向量方向上的投影
        normal := mgl32.Vec3{plane.X, plane.Y, plane.Z}
        distance := plane.W

        // 计算包围盒的8个顶点
        vertices := []mgl32.Vec3{
            {box.min.X, box.min.Y, box.min.Z},
            {box.max.X, box.min.Y, box.min.Z},
            {box.min.X, box.max.Y, box.min.Z},
            {box.max.X, box.max.Y, box.min.Z},
            {box.min.X, box.min.Y, box.max.Z},
            {box.max.X, box.min.Y, box.max.Z},
            {box.min.X, box.max.Y, box.max.Z},
            {box.max.X, box.max.Y, box.max.Z},
        }

        // 检查所有顶点是否在平面的负半空间
        allOutside := true
        for _, vertex := range vertices {
            if mgl32.Dot(normal, vertex)+distance > 0 {
                allOutside = false
                break
            }
        }

        if allOutside {
            return false
        }
    }

    return true
}
```

### 6.2 LOD系统

```go
// LOD系统
type LODSystem struct {
    levels      []*Mesh
    distances   []float32
    currentLOD  int
}

// 选择LOD级别
func (lod *LODSystem) SelectLOD(cameraPos, objectPos mgl32.Vec3) *Mesh {
    distance := cameraPos.Sub(objectPos).Len()

    for i, dist := range lod.distances {
        if distance <= dist {
            lod.currentLOD = i
            return lod.levels[i]
        }
    }

    lod.currentLOD = len(lod.levels) - 1
    return lod.levels[lod.currentLOD]
}
```

### 6.3 实例化渲染

```go
// 实例化渲染
type InstancedRenderer struct {
    mesh        *Mesh
    instances   []*GameObject
    instanceBuffer uint32
    maxInstances  int
}

// 创建实例化渲染器
func NewInstancedRenderer(mesh *Mesh, maxInstances int) *InstancedRenderer {
    ir := &InstancedRenderer{
        mesh:         mesh,
        maxInstances: maxInstances,
    }

    // 创建实例缓冲区
    gl.GenBuffers(1, &ir.instanceBuffer)
    gl.BindBuffer(gl.ARRAY_BUFFER, ir.instanceBuffer)
    gl.BufferData(gl.ARRAY_BUFFER, maxInstances*int(unsafe.Sizeof(mgl32.Mat4{})), nil, gl.DYNAMIC_DRAW)

    // 设置实例属性
    gl.BindVertexArray(mesh.VAO)

    // 实例变换矩阵
    for i := 0; i < 4; i++ {
        gl.VertexAttribPointer(5+i, 4, gl.FLOAT, false, int32(unsafe.Sizeof(mgl32.Mat4{})), gl.PtrOffset(i*int(unsafe.Sizeof(mgl32.Vec4{}))))
        gl.VertexAttribDivisor(5+i, 1)
        gl.EnableVertexAttribArray(5 + i)
    }

    gl.BindVertexArray(0)

    return ir
}

// 渲染实例
func (ir *InstancedRenderer) Render() {
    if len(ir.instances) == 0 {
        return
    }

    // 更新实例数据
    transforms := make([]mgl32.Mat4, len(ir.instances))
    for i, instance := range ir.instances {
        transforms[i] = instance.transform.GetMatrix()
    }

    gl.BindBuffer(gl.ARRAY_BUFFER, ir.instanceBuffer)
    gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(transforms)*int(unsafe.Sizeof(mgl32.Mat4{})), gl.Ptr(transforms))

    // 渲染实例
    gl.BindVertexArray(ir.mesh.VAO)
    gl.DrawElementsInstanced(gl.TRIANGLES, ir.mesh.vertexCount, gl.UNSIGNED_INT, nil, int32(len(ir.instances)))
    gl.BindVertexArray(0)
}
```

## 7. 总结

实时渲染系统是现代游戏引擎的核心组件，需要综合考虑性能、质量和可扩展性。通过合理的架构设计和优化策略，可以构建出高性能的渲染系统。

### 7.1 关键特性

- **高性能渲染**: 支持60FPS以上的实时渲染
- **现代渲染效果**: 支持PBR、阴影、后处理等效果
- **可扩展架构**: 支持自定义着色器和渲染管线
- **跨平台支持**: 支持多种图形API和平台
- **优化技术**: 视锥体剔除、LOD、实例化等优化

### 7.2 扩展方向

- **光线追踪**: 支持实时光线追踪渲染
- **体积渲染**: 支持云、雾、烟等体积效果
- **程序化生成**: 支持程序化地形和植被
- **VR/AR**: 支持虚拟现实和增强现实
- **移动优化**: 针对移动平台的优化

---

**参考文献**:
1. Real-Time Rendering, Tomas Akenine-Möller
2. OpenGL Programming Guide, Dave Shreiner
3. Physically Based Rendering, Matt Pharr
4. Game Engine Architecture, Jason Gregory

**相关链接**:
- [01-游戏引擎架构](./01-Game-Engine-Architecture.md)
- [02-网络游戏服务器](./02-Network-Game-Server.md)
- [04-物理引擎](./04-Physics-Engine.md)
