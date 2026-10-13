# CEF 在产品里的位置

核对日期：2026-10-13。分支号以 [Branches and Building](https://chromiumembedded.github.io/cef/branches_and_building.html) 为准。

CEF（Chromium Embedded Framework）是一份嵌进现有原生程序的浏览器库。宿主进程仍然是你的程序。CEF 提供 C 接口的 `libcef`，以及随二进制包发布的 `libcef_dll_wrapper` 源码，把这份 C 接口包成 C++ 类。页面跑在 Chromium 上，窗口、生命周期和安全更新由宿主负责。

官方仓库在 GitHub：<https://github.com/chromiumembedded/cef>。文档入口：<https://chromiumembedded.github.io/cef/>。预编译包由 Spotify 的自动构建提供：<https://cef-builds.spotifycdn.com/index.html>。

## 谁在用

留下 CEF 的产品，主体已经是原生程序，网页只是其中一块表面：

- Spotify 桌面客户端把界面嵌在自己的进程里，并托管官方二进制构建。
- OBS 的浏览器源要把网页合成进视频帧。
- Steam 客户端长期自带一份 `libcef`，版本往往落后于官方当前分支。

这些程序要的是：三个桌面系统上同一份 Chromium、由自己决定何时升级、以及把页面画进自己的窗口或像素缓冲。

整份界面都是网页、还要和 Node 或 VS Code 插件共用代码的程序，走的是 Electron。只做 Windows、愿意让系统更新引擎的程序，走 WebView2。界面可以接受各系统网页组件差异的程序，走系统 WebView（Tauri 默认就是这条）。CEF 的代价是安装包里一份完整 Chromium，以及你必须自己发布安全更新。

## 用它之前先定的三件事

1. **页面要画在哪里。** 跨平台的 Views 窗口、嵌进已有原生父窗口，还是不建窗口、由你自己合成像素（离屏渲染）。Chrome style 目前要求 Views，并且官方架构说明写明它不支持离屏渲染。离屏渲染走 Alloy style。动手前用你选定的那份头文件核对 `CefWindowInfo`。
2. **谁拥有消息循环。** 程序可以让 CEF 独占线程（`CefRunMessageLoop`），也可以在自己的 UI 循环里泵 CEF，或在 Windows/Linux 上让 CEF 的 UI 线程单独跑。这决定了关闭顺序和哪些调用必须回到 UI 线程。
3. **二进制怎么升级。** 自 M133（2025 年 1 月）起，可以用 `CEF_API_VERSION` 钉住一套稳定 ABI，使宿主不必每次 Chromium 小版本都重编。行为兼容不在保证范围内，跨大版本仍要自己测。自 M138 起，Windows 沙箱不再把 `cef_sandbox` 静态库链进你的 exe，而是用官方的 `bootstrap.exe` 加载你的 DLL。

## 这份仓库讲什么

官方 General Usage 是接口说明书。这里按开发者真正卡住的层次写：进程边界、启动和关闭、Handler 怎么分工、JavaScript 怎么跨进程、页面怎么画出来、怎样把库放进已有工程。文中的调用顺序来自官方教程、General Usage、Architecture、JavaScript Integration、API Versioning 和 Sandbox Setup，代码是说明顺序用的骨架，落地时以你下载的那份 `include/` 为准。
