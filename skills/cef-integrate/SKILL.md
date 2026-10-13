---
name: cef-integrate
description: Integrate an already downloaded CEF binary into an existing native C or C++ application. Use when the host project already exists and must show a browser view, share its message loop, and shut down cleanly. 也用于把 CEF 集成进已有工程、嵌浏览器窗口、接消息循环。
---

# 把 CEF 集成进已有工程

读 `docs/05-lifecycle.md`、`docs/06-application-shape.md` 和 `docs/09-integrate.md`。二进制包已经由 `cef-download` 选好，并且 `cef-build` 证明这套工具链能跑起 `cefsimple`。不要用 `archive/` 里的旧工程文件。

## 接到现有 CMake 或工程

1. 确定放置方式再写窗口代码。新窗口用 cefsimple 默认的 Chrome style 和 Views。已有父窗口用 `CefWindowInfo::SetAsChild`。要把像素送进自己的合成器时用 Alloy style（`--use-alloy-style`）和 `SetAsWindowless`。官方架构说明里 Chrome style 不支持离屏渲染；以当前包的 `CefWindowInfo` 头文件为准。
2. 把 `libcef_dll` 作为子目录加入现有工程，或链接 `cef-build` 已经产出的 `libcef_dll_wrapper`。包装层和宿主使用同一个编译器、同一个 `CEF_API_VERSION`（若要钉稳定 ABI）。宏的含义见 `docs/03-binary-boundary.md`。
3. 链接 `libcef`。按二进制包 README.txt 把 `libcef`、`icudtl.dat`、pak、snapshot 和 `locales/` 拷到可执行文件旁边。macOS 保持 Framework 与 Helper 的 bundle 相对位置。Linux 设置 rpath 为 `$ORIGIN`。
4. 在现有入口里、创建业务窗口之前调用 `CefExecuteProcess`。这行必须出现在主进程和子进程都会执行到的路径上。返回值大于等于 0 时立刻退出。
5. 填 `CefSettings` 的可写 `cache_path` 和 `log_file`，再 `CefInitialize`。
6. 已有 GUI 循环时不要调用 `CefRunMessageLoop`。保持 `multi_threaded_message_loop` 为 false，在主线程周期调用 `CefDoMessageLoopWork`。需要 CEF 独占 UI 时才用 `CefRunMessageLoop`。两种循环不能一起跑。
7. 在 `OnContextInitialized` 里 `CreateBrowser`。`CefLifeSpanHandler` 按 `docs/05-lifecycle.md` 计数；父窗口第一次系统关闭只调用 `CloseBrowser(false)` 并取消关闭。最后一个 `OnBeforeClose` 之后停止泵，再 `CefShutdown`。
8. 页面要调用原生代码时用 `include/wrapper/cef_message_router.h`，接线要求见 `docs/07-javascript.md`。不要新增 `CefRegisterExtension` 或 `native function`。

开发期按官方 Hands On Tutorial 关闭沙箱，把 sandbox 指针传空。准备发布时再按 `docs/10-ship-and-upgrade.md`：M138 及更新的 Windows 使用包里的 `bootstrap.exe` 加载导出 `RunWinMain` 的 DLL，不要把 `cef_sandbox.lib` 链进 MSVC 编的 exe。

## 做完怎么判断

现有程序能打开指定 URL。关闭最后一个浏览器窗口后进程退出。日志不报缺少 pak。渲染进程不是启动后立刻消失。若接了消息路由，页面一次查询能收到原生返回的字符串。
