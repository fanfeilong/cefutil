---
name: cef-create-project
description: Create a new CEF secondary-development target from cef-project by copying cefsimple into an application directory. Use when the user wants a new CEF application project rather than editing Chromium. 也用于创建 CEF 二次开发工程、从 cefsimple 拉起自己的程序。
---

# 创建 CEF 二次开发工程

在 `cef-build` 已经把官方 `cefsimple` 编过并打开页面之后再做。新工程是 cef-project 里的一个目标，不是 Chromium 源码树里的补丁。

官方步骤在 https://github.com/chromiumembedded/cef-project 的 README，“Add your own project source code”。

## 步骤

1. 克隆 cef-project，按 `cef-download` 改好 `CEF_VERSION`，按 `cef-build` 先编过一次 `cefsimple`。
2. 在仓库根目录建 `myproject/`（名字用用户的产品名，全文一致）。
3. 把 `third_party/cef/cef_binary_*/tests/cefsimple/` 的内容复制到 `myproject/`。这是最小可运行结构：入口、`CefApp`、`CefClient`、CMake 已经带资源拷贝和 macOS bundle。
4. 在文件名和文件内容里把 `cefsimple` 换成项目名，把路径里的 `tests` 换成 `examples`。平台入口文件（`cefsimple_win.cc`、`cefsimple_linux.cc`、`cefsimple_mac.mm`、`process_helper_mac.cc`）一起改名。
5. 在顶层 `CMakeLists.txt` 末尾、已有的 `add_subdirectory` 之后加上：

```cmake
add_subdirectory(myproject)
```

6. 删掉 `build/` 后按 `cef-build` 的生成器重新配置，然后：

```bash
cmake --build build --config Release --target myproject
```

Linux 的单配置生成器去掉 `--config Release`。

7. 运行生成的可执行文件或 `.app`，确认仍能打开默认页面，并且关掉窗口后进程退出。

## 随后改什么

保持入口顺序不变：`CefExecuteProcess`，返回值大于等于 0 就退出，然后 `CefInitialize`，再跑消息循环。第一个业务改动只换初始 URL。生命周期、父窗口和 JavaScript 桥分别按 `docs/05-lifecycle.md`、`docs/08-hosting.md`、`docs/07-javascript.md` 改，不要在这一步换成 `CefRegisterExtension`。

Bazel 路径是实验性质。用户没有指定 Bazel 时用 CMake。

## 做完怎么判断

新目标单独编过，不依赖改 `tests/cefsimple` 的源文件。运行时资源在输出目录或 app bundle 里。顶层 CMake 仍从 Spotify 拉取与 `CEF_VERSION` 相同的包。
