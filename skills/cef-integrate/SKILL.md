---
name: cef-integrate
description: Integrate Chromium Embedded Framework into an existing native C or C++ application using a CEF binary distribution. Use when adding CEF, embedding a browser window, wiring CefExecuteProcess, or the host builds but no page appears. 也用于把 CEF 嵌进已有原生程序、第一次创建浏览器窗口、子进程秒退或只有库没有页面。
---

# 把 CEF 嵌进已有程序

动手前读本仓库 `docs/02-architecture.md`、`docs/05-lifecycle.md`、`docs/09-integrate.md`。接口签名以用户下载的那份 `include/` 为准，不要沿用 2015 年的 GYP 或 `native function` 示例。

## 先写下来的决定

向用户确认，缺了就按默认并在回复里说明：

- 目标系统与架构。
- 分支：打开 https://chromiumembedded.github.io/cef/branches_and_building.html ，选仍受支持的分支。需要少升级用 Extended 或 LTC。
- style：默认 Chrome style。要离屏渲染则用 Alloy style（`--use-alloy-style`）。
- 放置：Views、已有原生父窗口，或离屏。Chrome style 按官方架构页要求 Views，且不支持离屏。
- 消息循环：已有 GUI 框架用外部泵；CEF 独占 UI 时用 `CefRunMessageLoop`。

## 步骤

1. 从 https://cef-builds.spotifycdn.com/index.html 取对应平台的 `cef_binary_*`。先让包内 `cefsimple` 编过并打开一个 URL。
2. 用与宿主相同的编译器和同一个 `CEF_API_VERSION`（若钉 ABI）编译 `libcef_dll_wrapper`。宿主链接 wrapper 与 `libcef`。
3. 按包内 README.txt 拷贝运行时文件（`libcef`、`icudtl.dat`、pak、snapshot、`locales/`）。macOS 使用 Framework + Helper 的 bundle 布局。
4. 在进程入口、创建业务窗口之前调用 `CefExecuteProcess`。返回值 `>= 0` 时立刻用该值退出。然后再 `CefInitialize`。
5. 实现 `CefApp`。浏览器进程 handler 的 `OnContextInitialized` 里 `CefBrowserHost::CreateBrowser`。
6. 实现 `CefClient` + `CefLifeSpanHandler`：`OnAfterCreated` 计数，`DoClose` 设标志并返回 `false`，`OnBeforeClose` 释放浏览器，计数归零再停止消息循环，然后 `CefShutdown`。
7. 原生父窗口的第一次系统关闭只调用 `CloseBrowser(false)` 并取消关闭。
8. 给出可写的 `CefSettings.cache_path` 和 `log_file`。开发期按沙箱文档传空的 sandbox info。发布包再按 `docs/10-ship-and-upgrade.md` 打开沙箱。

## 不要做

- 不要为了嵌一个窗口去编整份 Chromium，除非官方二进制没有该架构。
- 不要设置不存在的 `CefSettings.command_line_args`。开关放在 `CefApp::OnBeforeCommandLineProcessing`。
- 不要在子进程路径上跳过 `CefExecuteProcess`。
- 不要在 `CloseBrowser` 返回后立刻 `CefShutdown`。

## 做完怎么判断

冷启动打开指定 URL。关掉最后一个窗口后进程退出。日志里没有缺 pak。子进程不是启动后立刻消失。
