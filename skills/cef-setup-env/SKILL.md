---
name: cef-setup-env
description: Install the toolchain required to configure and compile a CEF binary distribution or cef-project. Use before the first CEF build on Windows, macOS, or Linux. 也用于配置 CEF 编译环境、安装 CMake、Visual Studio、Xcode 或 GTK 开发包。
---

# 配置编译 CEF 二进制包的环境

这条路径使用官方二进制包和 CMake，不克隆 Chromium，不安装 depot_tools，不用 GYP。只有 `cef-download` 确认该架构没有二进制包时，才改走官方 Master Build Quick Start。

读 `docs/09-integrate.md` 的前两节。版本和包由 `cef-download` 决定。

## 各平台要有的工具

共同要求：

- CMake 3.21 或更新。用 `cmake --version` 确认。
- Python 3.9、3.10 或 3.11。cef-project 的 CMake 会调用 Python 下载 clang-format。3.12 及以上先不要假设能过；配置报 Python 错误时改用 3.11，并用 `PYTHON_EXECUTABLE` 指到那个解释器。
- Git，用来克隆 https://github.com/chromiumembedded/cef-project 。

Windows：

- Visual Studio 2022，含“使用 C++ 的桌面开发”和与目标架构匹配的 MSVC、Windows 10/11 SDK。
- 生成器使用 `"Visual Studio 17 2022"`，架构用 `-A x64`、`-A Win32` 或 `-A ARM64`。
- 在装有 VS 的开发者命令行里运行 CMake，保证能找到编译器。

macOS：

- 安装 Xcode 和命令行工具（`xcode-select --install`）。
- 生成器使用 `"Xcode"`。Apple 芯片用 `-DPROJECT_ARCH=arm64`，Intel 用 `-DPROJECT_ARCH=x86_64`。
- cef-project 的 README 把已测试范围写到 Xcode 16.4。本机更新时仍先用这套生成器；只有 CMake 明确报生成器或 SDK 错误时再按报错调整，不要降级整个系统来猜。

Linux：

- 发行版至少能满足官方分支页上的 Ubuntu 20.04 或 Debian 10 要求。cef-project 推荐 Ubuntu 22.04 和 GCC 11 及以上。
- 安装 `build-essential`、`libgtk-3-dev`、cmake 和所选 Python。
- 生成器可用 `"Unix Makefiles"` 或 `Ninja`，并设置 `-DCMAKE_BUILD_TYPE=Release`。

## 做完怎么判断

`cmake --version` 不低于 3.21，`python3 --version` 落在 3.9–3.11，本机编译器能被 CMake 找到（Windows 上是 cl，macOS 上是 Xcode，Linux 上是 g++ 或 clang++）。然后运行 `cef-download` 和 `cef-build`。
