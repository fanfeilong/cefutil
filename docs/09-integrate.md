# 嵌进已有程序

核对日期：2026-10-13。步骤对应 [Tutorial](https://chromiumembedded.github.io/cef/tutorial.html)、[General Usage](https://chromiumembedded.github.io/cef/general_usage.html) 和 [Sandbox Setup](https://chromiumembedded.github.io/cef/sandbox_setup.html)。先读 [位置](01-place.md) 到 [页面放在哪里](08-hosting.md)，再按下面的顺序做。不要从编译整份 Chromium 开始。

## 1. 选定分支和放置方式

打开 [Branches and Building](https://chromiumembedded.github.io/cef/branches_and_building.html)，选一个仍在支持窗口内的分支。需要少升级的产品用 Extended 或 LTC；要跟新的 Web 平台用当前 Beta，并接受它很快变成不受支持的旧分支。

同时写下三行决定，后面的代码都由它们推出：

- style：Chrome style，还是 Alloy style。
- 放置：Views、原生父窗口，还是离屏。
- 循环：CEF 独占、外部泵，还是（Windows/Linux）独立 UI 线程。

离屏渲染不要选 Chrome style，除非你在这份头文件里证实 `SetAsWindowless` 可用。

## 2. 拿二进制包，而不是自己编 Chromium

从 <https://cef-builds.spotifycdn.com/index.html> 下载与目标系统、架构匹配的 `cef_binary_*`。包内的 `CMakeLists.txt` 和 `cefsimple` 能先在本机编过，再把同一套文件放进你的工程。

自己编 Chromium 只在你要改 CEF 补丁或官方包没有的架构时才做。那条路径用 GN，见 Master Build Quick Start。GYP 已经退出 Chromium，旧笔记里的 Visual Studio 2010 步骤不能用。

## 3. 编 wrapper，再编宿主

1. 用和宿主相同的编译器、相同的 `CEF_API_VERSION`（如果要钉 ABI）编译 `libcef_dll_wrapper`。
2. 宿主链接 wrapper 和 `libcef`。
3. 把 Resources 和 Release/Debug 里 README.txt 标成必需的文件放到运行目录。缺 `icudtl.dat`、pak 或 `locales` 时，进程能起来但页面是白的或立刻退出。

各平台的目录形状以包内 README.txt 为准。macOS 必须是 bundle：`Chromium Embedded Framework.framework` 和 Helper.app 放在 `Contents/Frameworks`。Linux 要让 rpath 或 `LD_LIBRARY_PATH` 找到 `libcef.so`。

## 4. 把入口接到现有主函数

在创建任何业务窗口之前：

1. 构造 `CefMainArgs` 和 `CefApp`。
2. 调用 `CefExecuteProcess`。返回值大于等于 0 就退出进程。
3. 填 `CefSettings`：至少给出可写的 `cache_path`、日志路径和严重级别。
4. `CefInitialize`。
5. 按你选的方式跑消息循环。
6. 在 `OnContextInitialized` 里 `CreateBrowser`。

已有 GUI 框架时，用外部泵，不要再调用 `CefRunMessageLoop`，否则两条循环抢线程。`CefInitialize` 所在的线程就是 `TID_UI`（独立 UI 线程模式除外）。

## 5. 放上最小的 Client

先实现 `CefLifeSpanHandler`：`OnAfterCreated` 计数并保存浏览器，`DoClose` 设关闭标志并返回 `false`，`OnBeforeClose` 释放指针，计数归零时停止消息循环。父窗口的关闭处理按 [启动、消息循环、关闭](05-lifecycle.md) 取消第一次系统关闭。

这一步能打开 URL、能关掉进程，再加加载回调和 JS 桥。

## 6. 需要页面调原生时再加路由

按 [JavaScript 和原生代码](07-javascript.md) 接 `CefMessageRouter`。先做一个往返：页面发起查询，浏览器进程 `Success` 一个固定字符串。通了之后再把 `Handler` 接到你的业务对象。业务对象若有线程亲和，在 `OnQuery` 里投递过去，完成后再调用 callback。

## 7. 沙箱放到最后开

开发期可以按沙箱文档把 `windows_sandbox_info` 传空、明确关闭沙箱，先把功能跑通。发布构建再打开。M138 及更新的 Windows 包使用 `bootstrap.exe` 加载导出了 `RunWinMain` 的客户端 DLL，而不是把 `cef_sandbox.lib` 链进 exe。macOS 在 Helper 的 `main` 最前面调用沙箱初始化，并用 `CefScopedLibraryLoader` 加载框架。细节见 [发布与升级](10-ship-and-upgrade.md)。

## 8. 对照样本查问题

| 现象 | 先看 |
| --- | --- |
| 子进程窗口一闪就没，没有页面 | 入口是否在所有进程里都调用了 `CefExecuteProcess`；macOS Helper 是否在 bundle 里 |
| 白屏，日志报找不到 pak | 资源目录、`resources_dir_path`、`locales_dir_path` |
| 关掉窗口进程还在 | `OnBeforeClose` 是否调用了 `CefQuitMessageLoop`；外部泵是否还在转 |
| 关窗口直接退出，没有 `onbeforeunload` | 第一次 `WM_CLOSE` 被放行了，没有 `CloseBrowser(false)` |
| JS 函数不存在 | 绑定写在了浏览器进程，或用了已废弃的 `CefRegisterExtension` |
| 查询没有回调 | 路由的 `OnBeforeBrowse` / `OnProcessMessageReceived` 没转发 |
| 换一个 libcef 就崩溃在启动 | `CEF_API_VERSION` 和这份二进制支持的版本不一致 |

远程调试端口打开后，用 DevTools 看渲染进程。浏览器进程的日志看 `CefSettings.log_file`。
