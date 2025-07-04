# CEF 现代开发指南 (2025年)

**文档版本**: 2025年7月4日  
**适用版本**: CEF 137.x/138.x+  
**维护状态**: 活跃更新中

## 🔧 开发环境配置要求

### 操作系统要求
- **Windows**: Windows 10 64位 (版本 1903 或更新) 或 Windows 11
- **macOS**: macOS 12.0 (Monterey) 或更新版本
- **Linux**: Ubuntu 20.04 LTS 或更新版本

### 编译器要求
- **推荐**: Visual Studio 2022 Community/Professional/Enterprise
- **最低**: Visual Studio 2019 (版本 16.11 或更新)
- **必须组件**: 
  - MSVC v143 工具集 (VS 2022) 或 MSVC v142 (VS 2019)
  - Windows SDK 10.0.22000 或更新版本
  - CMake 3.21 或更新版本

### 依赖工具
- **depot_tools**: [官方版本](https://chromium.googlesource.com/chromium/tools/depot_tools.git)
- **Python**: 3.8 - 3.11 (depot_tools 会自动管理)
- **Git**: 2.30 或更新版本
- **Ninja**: 包含在 depot_tools 中

### Visual C++ 运行时要求
- **CEF 126.x+**: Visual C++ 2019 Redistributable (x64)
- **CefSharp**: Visual C++ 2019 Redistributable 或更新版本

### 构建系统现状
- **当前主流**: GN (Generate Ninja) - Chromium 官方构建系统
- **CEF**: 使用 GN + Ninja 构建管道
- **CefSharp**: 支持 .NET Framework 4.6.2+ 和 .NET 6+

### 快速环境检查
```cmd
# 检查 Python 版本
python --version

# 检查 Git 版本  
git --version

# 检查 Visual Studio 版本
where cl.exe
cl.exe
```

## 📚 相关文档

- [CEF 最新特性指南](cef_features_2025.md)
- [性能优化最佳实践](performance_optimization_2025.md)
- [安全开发指南](security_best_practices_2025.md)
- [多平台开发指南](multi_platform_guide_2025.md)

## 🔗 官方资源

- [CEF 官方构建指南](https://bitbucket.org/chromiumembedded/cef/wiki/BranchesAndBuilding)
- [Chromium 开发者文档](https://chromium.googlesource.com/chromium/src/+/main/docs/)
- [CefSharp 构建指南](https://github.com/cefsharp/CefSharp/wiki) 