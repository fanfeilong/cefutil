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
- 只编译官方二进制包时，系统自带的 Python 即可。要走 cef-project 时，它的 CMake 会用 Python 下载 clang-format，README 要求 3.9、3.10 或 3.11。3.12 配置失败时，把 `PYTHON_EXECUTABLE` 指到 3.11；本机没有 3.11 时改用 `cef-create-project` 里对着已解压目录建工程的做法，不要为了 clang-format 更换系统 Python。
- Git，用来克隆 https://github.com/chromiumembedded/cef-project 。只有走 cef-project 时才需要。

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
- 安装 `build-essential`、`cmake`、`libgtk-3-dev`。运行窗口还需要 `libnss3`、`libasound2t64`（旧发行版是 `libasound2`）、`libgbm1`、`libxcomposite1`、`libxdamage1`、`libxrandr2`、`libxss1`、`libxkbcommon0`、`libatk1.0-0`、`libatk-bridge2.0-0`、`libcups2`、`libpango-1.0-0`、`libdrm2`。
- 生成器可用 `"Unix Makefiles"` 或 `Ninja`，并设置 `-DCMAKE_BUILD_TYPE=Release`。
- 显式指定编译器：`-DCMAKE_C_COMPILER=gcc -DCMAKE_CXX_COMPILER=g++`。`/usr/bin/c++` 若是 clang，而 clang 选中的 GCC 版本没有 `libstdc++.so`，链接会报 `cannot find -lstdc++`。Ubuntu 24.04 上曾出现 clang 18 选择 GCC 14、库却装在 GCC 13 的情况，改用 g++ 13 后通过。

## 做完怎么判断

`cmake --version` 不低于 3.21，本机编译器能被 CMake 找到（Windows 上是 cl，macOS 上是 Xcode，Linux 上是 g++）。走 cef-project 时再确认 `python3 --version` 落在 3.9–3.11。然后运行 `cef-download` 和 `cef-build`。
