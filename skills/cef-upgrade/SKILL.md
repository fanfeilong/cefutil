---
name: cef-upgrade
description: Upgrade a CEF binary distribution, pin CEF_API_VERSION, and enable the platform sandbox including the Windows bootstrap executable. Use when bumping libcef, fixing an API hash crash, preparing a release build, or replacing cef_sandbox.lib. 也用于升级 CEF、打开沙箱、bootstrap.exe、API 哈希启动失败。
---

# 升级 CEF 并打开沙箱

读 `docs/03-binary-boundary.md` 和 `docs/10-ship-and-upgrade.md`。版本号不要从本仓库抄。打开 https://chromiumembedded.github.io/cef/branches_and_building.html 看当天仍支持的分支。

## 换包

1. 选定仍受支持的分支，下载同一系统、同一架构的二进制包。
2. 若工程定义了 `CEF_API_VERSION`，确认新包的支持列表包含这个 `XXXYY`。不包含就改到新包仍接受的稳定编号。
3. 用相同编译器和相同宏重编 `libcef_dll_wrapper`，再编宿主。禁止旧 wrapper 配新 `libcef`。
4. 替换 `libcef`、pak、snapshot、`icudtl.dat`、`locales`。保留应用自己的缓存目录，除非发行说明要求清空。
5. 搜索并删除对 `CefRegisterExtension`、`OnWebKitInitialized` 里注册扩展、`native function` 的依赖。这些在 M154 起废弃。
6. 重测冷启动、子进程、最后一个窗口关闭、`onbeforeunload`、JS 往返，以及离屏的一帧（若使用）。

启动时因 API 哈希失败退出：宿主的 `CEF_API_VERSION` 与库不一致，或 wrapper 用了另一个编号编译。把两边调成同一个仍受支持的编号后重编。

不定义 `CEF_API_VERSION` 时宿主使用实验 API，只能配这一版 `libcef`。要让后续小版本少重编，就钉一个稳定编号，并接受实验接口会被编译掉。

## Windows 沙箱（M138 及更新）

不要把 `cef_sandbox.lib` 链进 MSVC 编的 exe。按 `cef_sandbox_win.h`：

- 宿主编成 DLL，导出 `RunWinMain` 或 `RunConsoleMain`。
- 用包里的 `bootstrap.exe` 或 `bootstrapc.exe` 加载它。
- 把收到的 `sandbox_info` 传入 `CefExecuteProcess` 和 `CefInitialize`。
- 启用沙箱时不要设置 `CefSettings.browser_subprocess_path`。

开发期可以传空 sandbox info 并明确关闭沙箱。发布配置不要留下这个状态，除非用户书面接受渲染进程无沙箱。

## macOS 沙箱（M138 及更新）

Helper 的 `main` 在加载框架之前完成沙箱初始化。用 `CefScopedLibraryLoader` 加载 Framework。分发 `libcef_sandbox.dylib`。签名和公证按 https://chromiumembedded.github.io/cef/code_signing.html 。

## 做完怎么判断

发布配置的程序在干净机器上启动，渲染进程处于沙箱中（平台工具能看到受限进程），页面能打开，关闭后进程退出。API 版本宏、wrapper 和 `libcef` 来自同一次选定的包。
