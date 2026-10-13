---
name: cef-build
description: Configure and compile cefsimple from a CEF binary distribution or cef-project to prove the toolchain. Use when the user asks to compile CEF, build libcef_dll_wrapper, or run the sample browser. 也用于编译 CEF、编译 cefsimple、确认环境能编过。
---

# 编译 CEF 样本

环境用 `cef-setup-env`，包的版本用 `cef-download`。第一次编译的目标是跑通 `cefsimple`，不是编 Chromium。

## 只拿到了二进制包

在解压后的 `cef_binary_*` 目录里配置。学习构建按官方 Hands On Tutorial 关掉沙箱：

```bash
# Windows x64
cmake -B build -G "Visual Studio 17 2022" -A x64 -DUSE_SANDBOX=OFF
cmake --build build --config Release --target cefsimple

# Linux。学习构建关闭沙箱，并钉住能找到 libstdc++ 的 g++
cmake -B build -G "Unix Makefiles" -DCMAKE_BUILD_TYPE=Release \
  -DUSE_SANDBOX=OFF -DCMAKE_C_COMPILER=gcc -DCMAKE_CXX_COMPILER=g++
cmake --build build --target cefsimple

# macOS Apple 芯片；Intel 把 arm64 换成 x86_64，生成器仍是 Xcode
cmake -B build -G "Xcode" -DPROJECT_ARCH=arm64 -DUSE_SANDBOX=OFF
cmake --build build --config Release --target cefsimple
```

若这一版 CMake 报 `USE_SANDBOX` 未知，打开包内 `CMakeLists.txt` 顶部的选项说明，改用它实际提供的开关。不要自己写一个同名缓存变量应付报错。

产物：

- Windows：`build/tests/cefsimple/Release/cefsimple.exe`
- Linux：`build/tests/cefsimple/Release/cefsimple`。这一版的 `SET_CEF_TARGET_OUT_DIR` 在 Unix Makefiles 下同样使用 `Release/`。
- macOS：`build/tests/cefsimple/Release/cefsimple.app`

在有显示的环境里运行它，确认窗口里出现页面。Windows 也可以打开 `build/cef.sln`。

Chrome style 第一次使用某个 `root_cache_path` 时会弹出 “Chromium Additional Terms of Service”。非官方包里正文是空白占位。点 Accept 之后才出现浏览器窗口；选择写在该缓存目录里，换目录会再问一次。Linux 默认目录是 `~/.config/cef_user_data`。

cefsimple 的构建结束会打印一条手动命令：`sudo -- chown root:root` 和 `chmod 4755`，对象是 `chrome-sandbox`。`USE_SANDBOX=OFF` 的学习构建不要执行它。没有会话总线时，日志里的 dbus 和 UPower 错误不表示编译失败。

## 用 cef-project

```bash
git clone https://github.com/chromiumembedded/cef-project.git
cd cef-project
```

把顶层 `CMakeLists.txt` 的 `CEF_VERSION` 改成 `cef-download` 选定的字符串，再：

```bash
mkdir build && cd build
# Windows x64
cmake -G "Visual Studio 17 2022" -A x64 ..
cmake --build . --config Release --target cefsimple
# Linux
cmake -G "Unix Makefiles" -DCMAKE_BUILD_TYPE=Release \
  -DUSE_SANDBOX=OFF -DCMAKE_C_COMPILER=gcc -DCMAKE_CXX_COMPILER=g++ ..
cmake --build . --target cefsimple
# macOS arm64
cmake -G "Xcode" -DPROJECT_ARCH=arm64 ..
cmake --build . --config Release --target cefsimple
```

CMake 会把二进制包下到 `third_party/cef/`。Python 报错时设置 `PYTHON_EXECUTABLE` 后删掉 `build/` 再配置。换 `CEF_VERSION` 后同样删掉 `build/` 和已经解压的旧 `third_party/cef/cef_binary_*`。

## 沙箱

这次样本构建可以关掉沙箱，用来证明编译器和资源拷贝是对的。M138 及更新的 Windows 发布包不要把 `cef_sandbox.lib` 链进 exe。要出货时按 `docs/10-ship-and-upgrade.md` 使用 `bootstrap.exe` 和导出 `RunWinMain` 的 DLL。

## 失败时先看

- 找不到编译器：回到 `cef-setup-env`。
- 下载 404：`CEF_VERSION` 与 index 上的文件名不一致，`+` 必须保留。
- 程序秒退或白屏：看运行目录里有没有 `icudtl.dat`、pak、`locales/`。CMake 的 `cefsimple` 目标会拷贝它们；手拷时以包内 README.txt 为准。
- 子进程一闪：入口必须调用 `CefExecuteProcess`。样本已经这样做，不要在验证阶段改入口。
