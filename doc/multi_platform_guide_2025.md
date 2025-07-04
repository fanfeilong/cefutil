# CEF 多平台支持指南 (2025年)

**文档版本**: 2025年7月4日  
**覆盖平台**: Windows, macOS, Linux, Android, iOS  
**构建系统**: 现代化跨平台解决方案

## 🪟 Windows 平台

### 开发环境配置
```cmd
# 推荐的 Windows 开发环境
Visual Studio 2022 Community/Professional
Windows 11 SDK (10.0.22621.0 或更新)
Windows 10 version 1903+ 或 Windows 11
```

### Windows 特定优化
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

## 🍎 macOS 平台

### 开发环境要求
```bash
# macOS 开发环境
macOS 12.0 (Monterey) 或更新版本
Xcode 14.0 或更新版本
Command Line Tools for Xcode
```

### macOS 特定功能
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

## 🐧 Linux 平台

### 发行版支持
```bash
# 支持的 Linux 发行版
Ubuntu 20.04 LTS, 22.04 LTS, 24.04 LTS
Debian 11, 12
CentOS 8, 9
Fedora 36+
Arch Linux (rolling)
```

### Linux 构建配置
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

### 依赖包管理
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

## 📱 Android 平台

### Android Studio 配置
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

### JNI 集成
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

## 📱 iOS 平台 (实验性支持)

### iOS 集成配置
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

## 🔧 跨平台构建系统

### CMake 配置
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

### 条件编译管理
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

## 📊 平台兼容性矩阵

| 特性 | Windows | macOS | Linux | Android | iOS |
|------|---------|--------|-------|---------|-----|
| CEF 核心 | ✅ | ✅ | ✅ | ✅ | ⚠️ |
| WebGPU | ✅ | ✅ | ✅ | ✅ | ❌ |
| 硬件加速 | ✅ | ✅ | ✅ | ✅ | ✅ |
| 文件系统访问 | ✅ | ⚠️ | ✅ | ⚠️ | ❌ |
| 多进程 | ✅ | ✅ | ✅ | ⚠️ | ❌ |
| 自定义协议 | ✅ | ✅ | ✅ | ✅ | ⚠️ |

**图例**: ✅ 完全支持, ⚠️ 部分支持/有限制, ❌ 不支持

## 🚀 最佳实践建议

### 代码组织
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

### 性能优化建议
1. **渲染后端选择**: Windows (D3D12) > macOS (Metal) > Linux (Vulkan)
2. **内存管理**: 各平台使用原生内存分配器
3. **文件 I/O**: 利用平台特定的异步 I/O API
4. **网络优化**: HTTP/3 和 QUIC 在所有平台启用

## 📚 相关文档

- [现代开发指南](modern_development_guide_2025.md)
- [CEF 最新特性指南](cef_features_2025.md)
- [性能优化最佳实践](performance_optimization_2025.md)
- [安全开发指南](security_best_practices_2025.md)

**🌟 重要提示**: 现代 CEF 提供了出色的跨平台一致性，但充分利用平台特定优化能显著提升用户体验。建议采用"通用核心 + 平台优化"的架构模式。 