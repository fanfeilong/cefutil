⚠️ **重要提示**: 本仓库的部分内容可能已过时。CEF 项目仍在活跃开发中，请访问 [CEF 官方网站](https://bitbucket.org/chromiumembedded/cef) 获取最新信息。

**最后更新**: 2025年7月4日  
**CEF 当前版本**: 138.x  
**状态**: 部分文档需要更新 - 查看 [更新计划](update/README.md)

---

CEF全称是Chromium Embedded Framework，它是Chromium的Content API的封装库。

- CEF官网地址：https://bitbucket.org/chromiumembedded/cef
- CEF官方论坛：https://magpcss.org/ceforum
- CEFSharp：https://github.com/cefsharp/CefSharp
- ⚠️ ChromiumFX（已不活跃）: https://bitbucket.org/chromiumfx/chromiumfx

## CEF 文档中文翻译任务
--------------------
- [x] [CEF General Usage](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF%20General%20Usage.md)
- [x] [CEF General Usage Task List](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF%20General%20Usage%20Task%20List.md)
- [x] [CEF General Usage中文版，欢迎查阅](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF%20General%20Usage-zh-cn.md)

## CEF FAQ
-------
- [CEF 的关闭流程](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF_Close.md)
- [CEF 的C和C++封装](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF_cpp2c_annotation.md)
- [CEF 的JavaScript和C++跨语言交互](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF_JavaScript_Cpp.md)
- [Cef经典N大问题 by 扫地僧](http://blog.csdn.net/weolar/article/details/51994895)

## Chromium Documentations
-----------------
- [Chromium Resoures](https://github.com/fanfeilong/cefutil/blob/master/doc/chromium_resources.md)
- [Chromium Content Register V8 Extension](https://github.com/fanfeilong/cefutil/blob/master/doc/content_register_v8_extension.md)

### 构建系统文档（已过时）
**⚠️ 注意：以下 GYP 构建系统文档已过时，Chromium 已迁移到 GN 构建系统**
- [Chromium Generation Your Project (GYP)](https://github.com/fanfeilong/cefutil/blob/master/doc/gyp.md) - 仅供历史参考
- [Chromium GYP 中文翻译](https://github.com/fanfeilong/cefutil/blob/master/doc/gyp.pdf) - 仅供历史参考

**推荐使用最新的 GN 构建系统：**
- [GN 官方文档](https://gn.googlesource.com/gn/)
- [Chromium GN 快速开始](https://chromium.googlesource.com/chromium/src/+/main/tools/gn/docs/quick_start.md)

### 🔧 现代开发环境配置要求 (2025年)

**⚠️ 重要：以下信息基于 2025年7月的最新要求，替代了本仓库中的过时指南**

#### 操作系统要求
- **Windows**: Windows 10 64位 (版本 1903 或更新) 或 Windows 11
- **macOS**: macOS 12.0 (Monterey) 或更新版本
- **Linux**: Ubuntu 20.04 LTS 或更新版本

#### 编译器要求
- **推荐**: Visual Studio 2022 Community/Professional/Enterprise
- **最低**: Visual Studio 2019 (版本 16.11 或更新)
- **必须组件**: 
  - MSVC v143 工具集 (VS 2022) 或 MSVC v142 (VS 2019)
  - Windows SDK 10.0.22000 或更新版本
  - CMake 3.21 或更新版本

#### 依赖工具
- **depot_tools**: [官方版本](https://chromium.googlesource.com/chromium/tools/depot_tools.git)
- **Python**: 3.8 - 3.11 (depot_tools 会自动管理)
- **Git**: 2.30 或更新版本
- **Ninja**: 包含在 depot_tools 中

#### Visual C++ 运行时要求
- **CEF 126.x+**: Visual C++ 2019 Redistributable (x64)
- **CefSharp**: Visual C++ 2019 Redistributable 或更新版本

#### 构建系统现状
- **当前主流**: GN (Generate Ninja) - Chromium 官方构建系统
- **CEF**: 使用 GN + Ninja 构建管道
- **CefSharp**: 支持 .NET Framework 4.6.2+ 和 .NET 6+

#### 快速环境检查
```cmd
# 检查 Python 版本
python --version

# 检查 Git 版本  
git --version

# 检查 Visual Studio 版本
where cl.exe
cl.exe
```

#### 获取最新信息
- [CEF 官方构建指南](https://bitbucket.org/chromiumembedded/cef/wiki/BranchesAndBuilding)
- [Chromium 开发者文档](https://chromium.googlesource.com/chromium/src/+/main/docs/)
- [CefSharp 构建指南](https://github.com/cefsharp/CefSharp/wiki)

---

### 🆕 CEF 最新版本特性亮点 (2025年)

**⚠️ 重要：CEF 版本跟随 Chromium 发布周期，以下特性基于 CEF 137.x/138.x 版本**

#### Chromium 137/138 核心新特性

##### 🎨 CSS 增强
- **CSS `if()` 函数**: 简洁表达条件值，提升样式逻辑表达能力
  ```css
  div {
    color: var(--color);
    background-color: if(style(--color: white): black; else: white);
  }
  ```
- **CSS `reading-flow` 和 `reading-order`**: 控制元素的键盘焦点导航顺序
- **CSS 锚点定位增强**: `anchor-size()` 函数支持更多属性
- **CSS 横向书写模式**: 支持 `sideways-rl` 和 `sideways-lr`

##### 🔐 安全与隐私
- **Document-Isolation-Policy**: 无需 COOP/COEP 即可启用跨源隔离
- **Ed25519 加密算法**: Web Cryptography API 支持现代椭圆曲线签名
- **HSTS 跟踪防护**: 防止第三方利用 HSTS 缓存进行用户跟踪
- **:visited 链接历史分区**: 增强用户浏览历史隐私保护

##### 📱 移动与设备支持
- **Device Posture API**: 检测可折叠设备的物理姿态
- **File System Access for Android**: Android 平台文件系统访问
- **Capture All Screens**: 企业环境下捕获所有连接的屏幕
- **Element Capture**: 精确捕获 DOM 子树而非整个标签页

##### 🎮 媒体与图形
- **H265 (HEVC) 编解码器**: WebRTC 和 MediaRecorder 支持
- **Float16Array**: 高精度浮点数组支持
- **Canvas 浮点颜色**: 支持医疗可视化等高精度应用
- **WebGPU 增强**: 32位浮点纹理混合、纹理视图用法控制

##### 🔗 API 现代化
- **Observable API**: 响应式编程范式，处理异步事件流
- **fetchLater API**: 延迟请求，页面销毁时自动发送
- **WebAuthn 条件创建**: 密码升级为密钥支持
- **FedCM 多身份提供商**: 单一对话框显示多个 IdP

##### ⚡ 性能优化
- **JavaScript Promise Integration (JSPI)**: WebAssembly 与 JS Promise 集成
- **WebAssembly Memory64**: 支持大于 4GB 的线性内存
- **Energy Saver 冻结**: 节能模式下冻结后台标签页
- **编译提示魔法注释**: 优化 JavaScript 解析和编译

#### 开发者体验改进

##### 🛠️ 调试与监控
- **CSP 哈希报告**: 脚本资源 URL 和哈希值报告
- **响应时间置信度**: 识别双峰性能分布
- **崩溃报告调用栈**: 无响应页面的 JavaScript 调用栈

##### 📝 新增 Web API
- **Rewriter API**: AI 支持的文本重写和改写
- **Writer API**: AI 支持的文本生成
- **highlightsFromPoint API**: 精确的高亮检测交互

#### CEF 特定优势

##### 🔧 嵌入式应用优势
- **改进的进程隔离**: 更安全的多进程架构
- **现代 GPU 加速**: WebGPU 全面支持和优化
- **移动平台支持**: Android WebView 功能对等
- **企业级功能**: 屏幕捕获、设备控制等企业需求

##### 📊 兼容性与稳定性
- **长期支持版本**: 企业环境稳定性保证
- **向后兼容**: 渐进式 API 演进，减少破坏性变更
- **多平台一致性**: Windows、macOS、Linux 行为一致

#### 获取最新信息

- **CEF 发布页面**: [Spotify CEF Builds](https://cef-builds.spotifycdn.com/index.html)
- **官方更新日志**: [CEF Forum Announcements](https://magpcss.org/ceforum/viewforum.php?f=10)
- **Chromium 发布说明**: [Chrome Release Notes](https://developer.chrome.com/release-notes/)
- **版本跟踪**: [CEF 版本分支信息](https://bitbucket.org/chromiumembedded/cef/wiki/BranchesAndBuilding)

**💡 提示**: CEF 通常在 Chromium 稳定版发布后 2-4 周内跟进相应版本，建议关注官方论坛获取最新构建信息。

---

### ⚡ CEF 现代性能优化建议 (2025年)

**基于 CEF 137.x/138.x 和现代 Web 技术的性能优化最佳实践**

#### 🎮 WebGPU 渲染优化

##### 内存管理优化
- **32位浮点纹理混合**: 适用于 HDR 渲染和科学可视化
  ```javascript
  // 请求 float32-blendable 特性
  const device = await adapter.requestDevice({
    requiredFeatures: ["float32-blendable"]
  });
  ```

- **纹理视图用法控制**: 精确控制纹理使用权限
  ```javascript
  const view = texture.createView({
    format: 'rgba8unorm-srgb',
    usage: GPUTextureUsage.RENDER_ATTACHMENT // 限制使用权限
  });
  ```

##### 渲染管线优化
- **集群式渲染**: 1024+ 光源场景性能提升 10-15 倍
  - 视锥体切片：32x18x48 集群配置
  - GPU 驱动的光源剔除和分组
  - 寄存器压力优化（目标 >90% 占用率）

- **子组（Subgroups）优化**: 机器学习工作负载加速
  ```wgsl
  // WGSL 子组函数示例
  let sum = subgroupAdd(localValue);
  let broadcast = subgroupBroadcast(value, 0);
  ```

##### Apple 平台特定优化
- **Metal 直接映射**: 近乎 1:1 的 API 映射，消除转换开销
- **半精度浮点数**: iOS/visionOS 内存减少 50%
  ```javascript
  // 请求 shader-f16 支持
  const device = await adapter.requestDevice({
    requiredFeatures: ["shader-f16"]
  });
  ```

- **瓦片式渲染优化**: 
  - 合并相关操作，最小化渲染通道边界
  - 利用片上内存效率

#### 🧠 WebAssembly 性能优化

##### 内存分配器升级
```bash
# 使用 mimalloc 替代 dlmalloc
emcc -sMALLOC=mimalloc your_app.c
```

**性能提升数据**:
- 单线程: 2x 性能提升（dlmalloc 2660ms vs mimalloc 1466ms）
- 多线程: 近乎完美的扩展性能
- 内存密集型应用: 10x+ 性能提升

##### 文件系统现代化
```bash
# 使用 WasmFS 替代 JS FS
emcc -sWASMFS your_app.c
```

**WasmFS 优势**:
- pthread 文件操作: 32x 性能提升
- 主线程操作: 2x 性能提升
- OPFS 支持: 现代持久化存储

#### 🔧 进程架构优化

##### 多进程配置调优
- **主进程职责**: UI 线程专注用户交互
- **渲染进程优化**: 独立 GPU 上下文，避免阻塞
- **沙盒效率**: 现代化沙盒减少 IPC 开销

##### GPU 进程现代化
- **命令缓冲池**: 减少创建/销毁开销
- **纹理共享优化**: DMA-BUF 零拷贝传输
- **同步点优化**: 显式围栏支持

#### 💾 内存管理最佳实践

##### 堆内存优化
```cpp
// CEF 应用中的内存池模式
class TexturePool {
  std::vector<CefRefPtr<CefTexture>> pool_;
public:
  CefRefPtr<CefTexture> Acquire(int width, int height);
  void Release(CefRefPtr<CefTexture> texture);
};
```

##### V8 优化配置
```cpp
// V8 堆限制优化
CefSettings settings;
settings.javascript_flags = 
  "--max-old-space-size=4096 "      // 4GB 老生代内存
  "--max-semi-space-size=256 "      // 256MB 新生代内存
  "--optimize-for-size";            // 优化代码大小
```

#### 🌐 网络性能优化

##### HTTP/3 和 QUIC
- **自动协议升级**: 减少延迟 20-30%
- **多路复用改进**: 消除队头阻塞
- **连接迁移**: 移动网络无缝切换

##### 资源加载策略
```cpp
// 预加载关键资源
CefRequestContext::LoadExtension(
  CefString("resource-hint-extension"),
  CefLoadHandler::Create([](bool success) {
    // 预连接到关键域名
    // 预加载关键样式表
  })
);
```

#### 📊 性能监控和调试

##### 现代调试工具
- **PIX on Windows**: GPU 性能分析
  ```bash
  # Chrome 启动参数
  --disable-gpu-sandbox 
  --disable-gpu-watchdog 
  --enable-dawn-features=emit_hlsl_debug_symbols
  ```

- **Chrome DevTools WebGPU**: 实时性能监控
- **内存使用追踪**: Process Explorer 集成

##### 性能指标监控
```javascript
// GPU 性能监控
const adapter = await navigator.gpu.requestAdapter();
const adapterInfo = adapter.info;
console.log(`GPU: ${adapterInfo.vendor} ${adapterInfo.architecture}`);

// 渲染时间测量
const querySet = device.createQuerySet({
  type: 'timestamp',
  count: 2,
});
```

#### 🚀 部署优化建议

##### 二进制大小优化
- **Link Time Optimization (LTO)**: 减少二进制大小 15-25%
- **死代码消除**: 移除未使用的 WebGPU 后端
- **压缩配置**: Brotli/Gzip 优化静态资源

##### 启动性能
```cpp
// 延迟初始化非关键组件
CefSettings settings;
settings.background_color = CefColorSetARGB(255, 255, 255, 255);
settings.no_sandbox = true;  // 开发环境可选
settings.single_process = false;  // 保持多进程架构
```

#### 🎯 平台特定优化

##### Windows 平台
- **D3D12 后端**: 优先使用 D3D12 而非 D3D11
- **MSAA 优化**: 4x MSAA 平衡质量和性能
- **HDR 支持**: 10-bit 颜色输出

##### macOS/iOS 平台
- **Metal Performance Shaders**: 机器学习工作负载加速
- **TrueDepth 相机**: WebXR/AR 应用优化
- **Neural Engine**: AI 推理加速

##### Linux 平台
- **Vulkan 优先**: 现代驱动程序支持
- **Wayland 显示**: 合成器直接集成
- **Mesa 优化**: 开源驱动程序调优

#### 📈 性能基准和目标

##### 目标性能指标
- **启动时间**: < 2 秒（冷启动）
- **首次绘制**: < 100ms
- **GPU 占用率**: > 85%（图形密集型应用）
- **内存使用**: < 设备内存的 30%

##### 基准测试工具
- **MotionMark**: WebGPU 渲染性能
- **WebAssembly 基准**: 计算性能
- **Lighthouse**: 整体 Web 性能

**💡 关键建议**: 现代 CEF 应用应充分利用 WebGPU、现代 WebAssembly 特性和平台特定优化来实现近原生的性能表现。

---

### 🔒 CEF 安全最佳实践 (2025年)

**基于现代威胁模型和 CEF 最新安全特性的全面安全指南**

#### 🛡️ 进程隔离和沙盒

##### 多进程架构安全
```cpp
// 强制启用多进程模式
CefSettings settings;
settings.single_process = false;  // 永远不要使用单进程模式
settings.no_sandbox = false;     // 确保沙盒启用

// 配置进程间通信安全
settings.command_line_args_disabled = true;  // 禁用命令行参数传递
```

##### 沙盒强化配置
```cpp
// Windows 平台沙盒配置
settings.sandbox_info = my_sandbox_info;  // 自定义沙盒信息

// 限制文件系统访问
settings.user_data_path = CefString("safe_app_data");
settings.cache_path = CefString("safe_cache");
```

#### 🔐 内容安全策略 (CSP)

##### 严格的 CSP 配置
```html
<!-- 推荐的 CSP 头部配置 -->
<meta http-equiv="Content-Security-Policy" content="
  default-src 'self';
  script-src 'self' 'unsafe-inline' 'wasm-unsafe-eval';
  style-src 'self' 'unsafe-inline';
  img-src 'self' data: https:;
  font-src 'self' data:;
  connect-src 'self' https:;
  media-src 'self' https:;
  object-src 'none';
  frame-src 'none';
  base-uri 'self';
  form-action 'self';
">
```

##### CEF 应用 CSP 实施
```cpp
// 在资源处理器中强制 CSP
class SecureResourceHandler : public CefResourceHandler {
public:
  void GetResponseHeaders(CefRefPtr<CefResponse> response, 
                         int64& response_length,
                         CefString& redirectUrl) override {
    CefResponse::HeaderMap headers;
    response->GetHeaderMap(headers);
    
    // 强制 CSP
    headers["Content-Security-Policy"] = 
      "default-src 'self'; script-src 'self' 'unsafe-inline' 'wasm-unsafe-eval'";
    
    // 安全头部
    headers["X-Frame-Options"] = "DENY";
    headers["X-Content-Type-Options"] = "nosniff";
    headers["Referrer-Policy"] = "strict-origin-when-cross-origin";
    
    response->SetHeaderMap(headers);
  }
};
```

#### 🌐 网络安全

##### HTTPS 强制和证书固定
```cpp
// 配置 HTTPS 强制
class SecureRequestHandler : public CefRequestHandler {
public:
  bool OnBeforeResourceLoad(CefRefPtr<CefBrowser> browser,
                           CefRefPtr<CefFrame> frame,
                           CefRefPtr<CefRequest> request,
                           CefRefPtr<CefRequestCallback> callback) override {
    
    CefString url = request->GetURL();
    
    // 强制 HTTPS
    if (url.ToString().substr(0, 7) == "http://") {
      std::string https_url = "https://" + url.ToString().substr(7);
      request->SetURL(https_url);
    }
    
    // 验证允许的域名
    if (!IsAllowedDomain(url.ToString())) {
      callback->Cancel();
      return true;
    }
    
    return false;
  }
  
private:
  bool IsAllowedDomain(const std::string& url) {
    // 实现域名白名单验证
    static const std::vector<std::string> allowed_domains = {
      "yourapp.com", "api.yourapp.com", "cdn.yourapp.com"
    };
    // ... 验证逻辑
    return true;
  }
};
```

##### 现代 TLS 配置
```cpp
// TLS 安全配置
CefSettings settings;
settings.command_line_args = {
  "--ssl-version-min=tls1.2",           // 最低 TLS 1.2
  "--disable-features=VizDisplayCompositor", // 禁用不必要特性
  "--disable-web-security",             // 仅开发环境使用
  "--cipher-suite-blacklist=0x0004,0x0005"  // 禁用弱密码套件
};
```

#### 🔒 数据保护

##### 敏感数据处理
```cpp
// 安全的用户数据存储
class SecureDataManager {
private:
  std::string EncryptData(const std::string& data) {
    // 使用 AES-256-GCM 加密
    // Windows: CryptProtectData
    // macOS: Keychain Services
    // Linux: libsecret
    return encrypted_data;
  }
  
  std::string DecryptData(const std::string& encrypted_data) {
    // 相应的解密逻辑
    return decrypted_data;
  }
  
public:
  void StoreSecurely(const std::string& key, const std::string& value) {
    std::string encrypted = EncryptData(value);
    // 存储到安全位置
  }
};
```

##### 内存安全措施
```cpp
// 安全内存清理
class SecureString {
private:
  std::vector<char> data_;
  
public:
  ~SecureString() {
    // 显式清零内存
    explicit_bzero(data_.data(), data_.size());
  }
};

// JavaScript 堆隔离
CefSettings settings;
settings.javascript_flags = 
  "--enable-heap-profiling "          // 启用堆分析
  "--max-old-space-size=2048 "        // 限制内存使用
  "--enable-pointer-compression";     // 启用指针压缩
```

#### 🎯 输入验证和 XSS 防护

##### JavaScript Bridge 安全
```cpp
// 安全的 JavaScript 绑定
class SecureV8Handler : public CefV8Handler {
public:
  bool Execute(const CefString& name,
               CefRefPtr<CefV8Value> object,
               const CefV8ValueList& arguments,
               CefRefPtr<CefV8Value>& retval,
               CefString& exception) override {
    
    if (name == "secureAPI") {
      // 严格的参数验证
      if (arguments.size() != 2 || 
          !arguments[0]->IsString() || 
          !arguments[1]->IsObject()) {
        exception = "Invalid arguments";
        return true;
      }
      
      // 输入清理和验证
      std::string input = arguments[0]->GetStringValue();
      if (!ValidateInput(input)) {
        exception = "Invalid input format";
        return true;
      }
      
      // 安全处理
      ProcessSecurely(input);
      retval = CefV8Value::CreateBool(true);
      return true;
    }
    
    return false;
  }
  
private:
  bool ValidateInput(const std::string& input) {
    // 实施严格的输入验证
    // 检查长度、字符集、格式等
    return true;
  }
};
```

##### DOM XSS 防护
```javascript
// 客户端 XSS 防护函数
function sanitizeInput(input) {
  const div = document.createElement('div');
  div.textContent = input;
  return div.innerHTML;
}

// 安全的 DOM 操作
function updateContent(userInput) {
  const sanitized = sanitizeInput(userInput);
  document.getElementById('content').textContent = sanitized;  // 使用 textContent
}
```

#### 🔍 安全监控和日志

##### 安全事件监控
```cpp
// 安全事件日志记录
class SecurityMonitor : public CefLoadHandler {
public:
  void OnLoadError(CefRefPtr<CefBrowser> browser,
                   CefRefPtr<CefFrame> frame,
                   ErrorCode errorCode,
                   const CefString& errorText,
                   const CefString& failedUrl) override {
    
    // 记录可疑的网络错误
    if (errorCode == ERR_CERT_INVALID || 
        errorCode == ERR_SSL_PROTOCOL_ERROR) {
      LogSecurityEvent("SSL Error", failedUrl.ToString(), errorCode);
    }
  }
  
private:
  void LogSecurityEvent(const std::string& event_type,
                       const std::string& url,
                       int error_code) {
    // 发送到安全日志系统
    // 包括时间戳、用户ID、设备指纹等
  }
};
```

##### 运行时安全检查
```cpp
// 定期安全检查
class SecurityChecker {
public:
  void PerformSecurityCheck() {
    // 检查进程完整性
    if (!VerifyProcessIntegrity()) {
      HandleSecurityBreach();
    }
    
    // 检查内存使用异常
    if (GetMemoryUsage() > MAX_ALLOWED_MEMORY) {
      HandleMemoryLeak();
    }
    
    // 检查网络连接
    ValidateNetworkConnections();
  }
  
private:
  bool VerifyProcessIntegrity() {
    // 验证进程签名和权限
    return true;
  }
  
  void HandleSecurityBreach() {
    // 安全事件响应
    // 1. 记录详细日志
    // 2. 通知管理员
    // 3. 采取防护措施
  }
};
```

#### 🔄 更新和补丁管理

##### 自动更新安全
```cpp
// 安全的自动更新机制
class SecureUpdater {
public:
  bool CheckForUpdates() {
    // 1. 使用 HTTPS 检查更新
    // 2. 验证数字签名
    // 3. 检查版本回滚攻击
    
    std::string update_url = "https://updates.yourapp.com/check";
    CefRefPtr<CefRequest> request = CefRequest::Create();
    request->SetURL(update_url);
    
    // 添加认证头
    CefRequest::HeaderMap headers;
    headers["Authorization"] = "Bearer " + GetSecureToken();
    headers["User-Agent"] = GetSecureUserAgent();
    request->SetHeaderMap(headers);
    
    return ProcessUpdateResponse(request);
  }
  
private:
  bool VerifyUpdateSignature(const std::string& update_data,
                           const std::string& signature) {
    // 使用公钥验证数字签名
    return true;
  }
};
```

#### 📱 移动平台安全

##### Android 特定安全措施
```cpp
// Android WebView 安全配置
#ifdef ANDROID
void ConfigureAndroidSecurity(CefSettings& settings) {
  // 禁用文件系统访问
  settings.command_line_args.push_back("--disable-file-system");
  
  // 启用严格站点隔离
  settings.command_line_args.push_back("--site-per-process");
  
  // 禁用不安全的协议
  settings.command_line_args.push_back("--disable-features=InsecureOriginWarnings");
}
#endif
```

#### 🎯 威胁建模和风险评估

##### 常见威胁向量
1. **代码注入攻击**: XSS、SQL 注入、命令注入
2. **中间人攻击**: TLS 降级、证书伪造
3. **本地文件包含**: 路径遍历、文件读取
4. **内存安全**: 缓冲区溢出、UAF 漏洞
5. **供应链攻击**: 第三方库、依赖污染

##### 防护策略矩阵
| 威胁类型 | 防护措施 | 实施优先级 |
|---------|---------|-----------|
| XSS | CSP + 输入验证 | 高 |
| MITM | 证书固定 + HSTS | 高 |
| 文件包含 | 沙盒 + 路径验证 | 中 |
| 内存安全 | ASLR + CFI | 中 |
| 供应链 | 签名验证 + SCA | 高 |

**🚨 关键提醒**: 安全是持续过程，需要定期更新威胁模型、进行安全审计和渗透测试。务必关注 CEF 官方安全公告和 CVE 数据库。

---

## CEF Build
-----------------
如果要编译CEF，则需要同步编译Chromium的最新源码
- [Chromium Build Guid](https://github.com/fanfeilong/cefutil/blob/master/doc/chromium_build_guid.md)
- [Chromium for developer](http://www.chromium.org/developers)
- [windows build instructions](https://chromium.googlesource.com/chromium/src/+/master/docs/windows_build_instructions.md)
  - **官方 depot_tools 仓库**: https://chromium.googlesource.com/chromium/tools/depot_tools.git
  - **系统要求**: Python 3.8+ (depot_tools 会自动管理 Python 环境)
  - **自动更新**: 运行 `gclient` 时会自动更新 depot_tools
  - **包含工具**: GN、autoninja、siso、ninja 等最新构建工具
  - 如遇网络问题，可搜索 GitHub 上的镜像仓库：[depot_tools mirrors](https://github.com/search?q=depot_tools)
  - google 官方开发的vs插件，专门为chromium源码提供的。[vs-chromium](https://github.com/chromium/vs-chromium)
  - 同样的原因，depot_tool里内置的chromium的源码下载地址放在googleapi.com上，也是被墙的。下面这个地址里有chromium的所有项目的git地址：
    - https://chromium.googlesource.com/?format=JSON
    - https://chromium.googlesource.com/?format=TEXT

相关链接
---------
- [servo, the embeddable browser engine](http://blogs.s-osg.org/servo-the-embeddable-browser-engine/)

- [geckoview](https://mozilla.github.io/geckoview/)
Android offers a built-in WebView, which applications can hook into in order to display web pages within the context of their app. However, Android's WebView is not really intended for building browsers, and hence, many advanced Web APIs are disabled. Furthermore, it is also a moving target: different phones might have different versions of WebView, all of which your app has to support.

That is where GeckoView comes in. GeckoView is:

* Full-featured: GeckoView is designed to expose the entire power of the Web to applications, and all that through a straightforward API. Think of it as harnessing the full power of Gecko (the engine that powers Firefox), while its API is WebView-like and easy to use.
* Suited for apps and browsers: GeckoView is particularly suited for building mobile browsers, but it can be embedded as a web engine component in any kind of app.
* Self-Contained: Because GeckoView is a standalone library that you bundle with your application, you can be confident that the code you test is the code that will actually run.
* Standards Compliant: Like Firefox, GeckoView offers excellent support for modern Web standards.

---

### 🌍 CEF 多平台支持指南 (2025年)

**全面覆盖 Windows、macOS、Linux、Android 和 iOS 平台的现代开发指南**

#### 🪟 Windows 平台

##### 开发环境配置
```cmd
# 推荐的 Windows 开发环境
Visual Studio 2022 Community/Professional
Windows 11 SDK (10.0.22621.0 或更新)
Windows 10 version 1903+ 或 Windows 11
```

##### Windows 特定优化
```cpp
// Windows 平台配置
#ifdef _WIN32
void ConfigureWindowsPlatform(CefSettings& settings) {
  // 启用 DirectComposition
  settings.command_line_args.push_back("--enable-direct-composition");
  
  // D3D12 渲染器
  settings.command_line_args.push_back("--use-d3d12");
  
  // 高 DPI 支持
  settings.command_line_args.push_back("--enable-high-dpi-support");
  settings.command_line_args.push_back("--force-device-scale-factor=1.0");
  
  // Windows 11 特性
  settings.command_line_args.push_back("--enable-features=WebGPU");
}
#endif
```

##### Windows WebView2 兼容性
```cpp
// WebView2 迁移兼容层
class CEFWebView2Bridge {
public:
  // 提供类似 WebView2 的 API
  void NavigateToUrl(const std::wstring& url) {
    if (cef_browser_) {
      cef_browser_->GetMainFrame()->LoadURL(url);
    }
  }
  
  void SetPermissionRequestHandler(
    std::function<void(const std::string&)> handler) {
    permission_handler_ = handler;
  }
};
```

#### 🍎 macOS 平台

##### 开发环境要求
```bash
# macOS 开发环境
macOS 12.0 (Monterey) 或更新版本
Xcode 14.0 或更新版本
Command Line Tools for Xcode
```

##### macOS 特定功能
```objc
// macOS 平台配置
#ifdef __APPLE__
void ConfigureMacOSPlatform(CefSettings& settings) {
  // Metal 渲染器
  settings.command_line_args.push_back("--use-metal");
  
  // Core Animation 支持
  settings.command_line_args.push_back("--enable-core-animation-api");
  
  // macOS 沙盒兼容
  settings.command_line_args.push_back("--enable-sandbox-logging");
  
  // Retina 显示支持
  settings.device_scale_factor = 2.0;
}

// NSView 集成
@interface CEFNSView : NSView
- (void)embedCEFBrowser:(CefRefPtr<CefBrowser>)browser;
@end
#endif
```

##### App Store 部署
```objc
// App Store 兼容配置
void ConfigureAppStoreCompatibility(CefSettings& settings) {
  // 禁用沙盒外功能
  settings.no_sandbox = false;
  settings.single_process = false;
  
  // 使用应用程序容器
  NSString* containerPath = [[NSHomeDirectory() 
    stringByAppendingPathComponent:@"Library"]
    stringByAppendingPathComponent:@"Containers"];
  
  settings.user_data_path = [containerPath UTF8String];
}
```

#### 🐧 Linux 平台

##### 发行版支持
```bash
# 支持的 Linux 发行版
Ubuntu 20.04 LTS, 22.04 LTS, 24.04 LTS
Debian 11, 12
CentOS 8, 9
Fedora 36+
Arch Linux (rolling)
```

##### Linux 构建配置
```cpp
// Linux 平台配置
#ifdef __linux__
void ConfigureLinuxPlatform(CefSettings& settings) {
  // Wayland 支持
  settings.command_line_args.push_back("--enable-features=UseOzonePlatform");
  settings.command_line_args.push_back("--ozone-platform=wayland");
  
  // Vulkan 渲染器
  settings.command_line_args.push_back("--enable-vulkan");
  
  // 字体配置
  settings.command_line_args.push_back("--force-color-profile=srgb");
}
#endif
```

##### 依赖包管理
```bash
# Ubuntu/Debian 依赖
apt-get install libnss3-dev libgconf-2-4 libxrandr2 \
  libasound2-dev libpangocairo-1.0-0 libatk1.0-0 \
  libcairo-gobject2 libgtk-3-0 libgdk-pixbuf2.0-0

# Fedora 依赖
dnf install nss-devel gtk3-devel alsa-lib-devel \
  libXrandr-devel libXcomposite-devel libXcursor-devel

# Arch Linux 依赖
pacman -S gtk3 nss alsa-lib libxrandr libxcomposite libxcursor
```

#### 📱 Android 平台

##### Android Studio 配置
```gradle
// app/build.gradle
android {
    compileSdk 34
    
    defaultConfig {
        minSdk 21  // CEF Android 最低要求
        targetSdk 34
        ndk {
            abiFilters 'arm64-v8a', 'armeabi-v7a', 'x86_64'
        }
    }
    
    buildFeatures {
        prefab true
    }
}

dependencies {
    implementation 'org.chromium.net:cronet-api:119.6045.31'
    implementation 'org.chromium.net:cronet-common:119.6045.31'
}
```

##### JNI 集成
```cpp
// Android JNI 桥接
#ifdef ANDROID
extern "C" JNIEXPORT void JNICALL
Java_com_yourapp_cef_CefManager_initializeCef(
    JNIEnv* env,
    jobject thiz,
    jstring app_data_path) {
    
  CefSettings settings;
  
  // Android 特定配置
  settings.command_line_args.push_back("--use-vulkan");
  settings.command_line_args.push_back("--enable-gpu-rasterization");
  
  // 获取应用数据路径
  const char* path = env->GetStringUTFChars(app_data_path, nullptr);
  settings.user_data_path = CefString(path);
  env->ReleaseStringUTFChars(app_data_path, path);
  
  CefInitialize(settings, nullptr, nullptr);
}
#endif
```

##### Kotlin 封装
```kotlin
// CEF Android Kotlin 封装
class CefWebView @JvmOverloads constructor(
    context: Context,
    attrs: AttributeSet? = null,
    defStyleAttr: Int = 0
) : FrameLayout(context, attrs, defStyleAttr) {
    
    private var cefBrowser: Long = 0
    
    fun loadUrl(url: String) {
        CefManager.loadUrl(cefBrowser, url)
    }
    
    fun evaluateJavaScript(script: String, callback: (String) -> Unit) {
        CefManager.evaluateJavaScript(cefBrowser, script) { result ->
            Handler(Looper.getMainLooper()).post {
                callback(result)
            }
        }
    }
}
```

#### 📱 iOS 平台 (实验性支持)

##### iOS 集成配置
```swift
// iOS CEF 集成 (基于 WKWebView 桥接)
import WebKit
import UIKit

class CEFWebView: UIView {
    private var webView: WKWebView!
    
    override func awakeFromNib() {
        super.awakeFromNib()
        setupWebView()
    }
    
    private func setupWebView() {
        let config = WKWebViewConfiguration()
        
        // CEF 兼容层
        config.preferences.setValue(true, forKey: "allowFileAccessFromFileURLs")
        config.setValue(true, forKey: "allowUniversalAccessFromFileURLs")
        
        webView = WKWebView(frame: bounds, configuration: config)
        addSubview(webView)
        
        // 自动布局
        webView.translatesAutoresizingMaskIntoConstraints = false
        NSLayoutConstraint.activate([
            webView.topAnchor.constraint(equalTo: topAnchor),
            webView.leadingAnchor.constraint(equalTo: leadingAnchor),
            webView.trailingAnchor.constraint(equalTo: trailingAnchor),
            webView.bottomAnchor.constraint(equalTo: bottomAnchor)
        ])
    }
}
```

#### 🔧 跨平台构建系统

##### CMake 配置
```cmake
# 跨平台 CMakeLists.txt
cmake_minimum_required(VERSION 3.21)
project(CEFMultiPlatform)

# CEF 查找
find_package(CEF REQUIRED)

# 平台特定源文件
if(WIN32)
    set(PLATFORM_SOURCES src/windows/platform_win.cpp)
    set(PLATFORM_LIBS d3d11 dxgi)
elseif(APPLE)
    if(IOS)
        set(PLATFORM_SOURCES src/ios/platform_ios.mm)
        find_library(WEBKIT_FRAMEWORK WebKit)
        set(PLATFORM_LIBS ${WEBKIT_FRAMEWORK})
    else()
        set(PLATFORM_SOURCES src/macos/platform_mac.mm)
        find_library(METAL_FRAMEWORK Metal)
        set(PLATFORM_LIBS ${METAL_FRAMEWORK})
    endif()
elseif(ANDROID)
    set(PLATFORM_SOURCES src/android/platform_android.cpp)
    set(PLATFORM_LIBS log android)
else()
    set(PLATFORM_SOURCES src/linux/platform_linux.cpp)
    find_package(PkgConfig REQUIRED)
    pkg_check_modules(GTK3 REQUIRED gtk+-3.0)
    set(PLATFORM_LIBS ${GTK3_LIBRARIES})
endif()

# 目标配置
add_executable(${PROJECT_NAME} 
    src/main.cpp 
    src/common/app.cpp
    ${PLATFORM_SOURCES}
)

target_link_libraries(${PROJECT_NAME} 
    ${CEF_LIBRARIES} 
    ${PLATFORM_LIBS}
)
```

##### 条件编译管理
```cpp
// platform_config.h - 平台配置统一管理
#pragma once

// 平台检测
#if defined(_WIN32)
    #define PLATFORM_WINDOWS 1
    #include "windows/platform_specific.h"
#elif defined(__APPLE__)
    #include <TargetConditionals.h>
    #if TARGET_OS_IPHONE
        #define PLATFORM_IOS 1
        #include "ios/platform_specific.h"
    #else
        #define PLATFORM_MACOS 1
        #include "macos/platform_specific.h"
    #endif
#elif defined(ANDROID)
    #define PLATFORM_ANDROID 1
    #include "android/platform_specific.h"
#elif defined(__linux__)
    #define PLATFORM_LINUX 1
    #include "linux/platform_specific.h"
#endif

// 跨平台抽象接口
class PlatformInterface {
public:
    virtual ~PlatformInterface() = default;
    virtual void InitializePlatform() = 0;
    virtual void ShutdownPlatform() = 0;
    virtual void* GetNativeWindowHandle() = 0;
    
    static std::unique_ptr<PlatformInterface> Create();
};
```

#### 📦 包管理和分发

##### Windows 分发
```xml
<!-- Package.appxmanifest for MSIX -->
<Package xmlns="http://schemas.microsoft.com/appx/manifest/foundation/windows10">
  <Identity Name="YourApp.CEF" 
            Publisher="CN=YourCompany" 
            Version="1.0.0.0" />
  
  <Properties>
    <DisplayName>Your CEF Application</DisplayName>
    <PublisherDisplayName>Your Company</PublisherDisplayName>
    <Description>CEF-based application</Description>
  </Properties>
  
  <Applications>
    <Application Id="YourApp" 
                 Executable="YourApp.exe" 
                 EntryPoint="YourApp.App">
      <uap:VisualElements DisplayName="Your CEF App"
                         Square150x150Logo="Assets/Square150x150Logo.png"
                         Square44x44Logo="Assets/Square44x44Logo.png"
                         BackgroundColor="transparent" />
    </Application>
  </Applications>
</Package>
```

##### macOS 分发
```bash
# macOS 应用程序包结构
YourApp.app/
├── Contents/
│   ├── Info.plist
│   ├── MacOS/
│   │   └── YourApp
│   ├── Resources/
│   │   └── icon.icns
│   └── Frameworks/
│       └── Chromium Embedded Framework.framework/

# 代码签名
codesign --force --deep --sign "Developer ID Application" YourApp.app
```

##### Linux 分发
```bash
# AppImage 打包
linuxdeploy --appdir AppDir --executable YourApp --output appimage

# Flatpak 配置
flatpak-builder build com.yourcompany.YourApp.yml
flatpak build-export repo build

# Snap 配置
snapcraft
```

#### 🔄 持续集成/持续部署

##### GitHub Actions 多平台构建
```yaml
# .github/workflows/multi-platform.yml
name: Multi-Platform Build

on: [push, pull_request]

jobs:
  build:
    strategy:
      matrix:
        include:
          - os: windows-latest
            target: windows-x64
          - os: macos-latest
            target: macos-universal
          - os: ubuntu-latest
            target: linux-x64
          - os: ubuntu-latest
            target: android-arm64
    
    runs-on: ${{ matrix.os }}
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Build Environment
      run: |
        if [ "${{ matrix.target }}" == "android-arm64" ]; then
          echo "Setup Android NDK"
        fi
    
    - name: Build CEF Application
      run: |
        mkdir build && cd build
        cmake .. -DTARGET_PLATFORM=${{ matrix.target }}
        cmake --build . --config Release
    
    - name: Package Application
      run: |
        # 平台特定的打包逻辑
        if [ "${{ runner.os }}" == "Windows" ]; then
          echo "Create MSIX package"
        elif [ "${{ runner.os }}" == "macOS" ]; then
          echo "Create DMG package"
        else
          echo "Create AppImage package"
        fi
```

#### 📊 平台兼容性矩阵

| 特性 | Windows | macOS | Linux | Android | iOS |
|------|---------|--------|-------|---------|-----|
| CEF 核心 | ✅ | ✅ | ✅ | ✅ | ⚠️ |
| WebGPU | ✅ | ✅ | ✅ | ✅ | ❌ |
| 硬件加速 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 文件系统访问 | ✅ | ⚠️ | ✅ | ⚠️ | ❌ |
| 多进程 | ✅ | ✅ | ✅ | ⚠️ | ❌ |
| 自定义协议 | ✅ | ✅ | ✅ | ✅ | ⚠️ |

**图例**: ✅ 完全支持, ⚠️ 部分支持/有限制, ❌ 不支持

#### 🚀 最佳实践建议

##### 代码组织
```cpp
// 推荐的项目结构
project/
├── src/
│   ├── common/          // 跨平台公共代码
│   ├── platform/        // 平台特定代码
│   │   ├── windows/
│   │   ├── macos/
│   │   ├── linux/
│   │   └── android/
│   └── app/             // 应用程序逻辑
├── resources/           // 跨平台资源
├── cmake/               // 构建配置
└── docs/                // 平台特定文档
```

##### 性能优化建议
1. **渲染后端选择**: Windows (D3D12) > macOS (Metal) > Linux (Vulkan)
2. **内存管理**: 各平台使用原生内存分配器
3. **文件 I/O**: 利用平台特定的异步 I/O API
4. **网络优化**: HTTP/3 和 QUIC 在所有平台启用

**🌟 重要提示**: 现代 CEF 提供了出色的跨平台一致性，但充分利用平台特定优化能显著提升用户体验。建议采用"通用核心 + 平台优化"的架构模式。

---
