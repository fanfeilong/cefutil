---
name: cef-lifecycle
description: Fix CEF startup, message-loop integration, and browser shutdown in a host application. Use when the process stays alive after close, exits before onbeforeunload, deadlocks, or sub-processes flash and quit. 也用于关窗口进程不退、onbeforeunload 来不及弹出、消息循环卡死、CefShutdown 崩溃。
---

# 修正 CEF 的启动和关闭

读 `docs/04-processes-and-threads.md` 和 `docs/05-lifecycle.md`。先看宿主入口和 `CefLifeSpanHandler`，再改。

## 判断卡在哪一段

- 没有任何页面，辅助进程一闪退出：入口没有在每个进程里调用 `CefExecuteProcess`，或 macOS Helper 不在 bundle 里，或命令行被宿主吃掉了。
- 窗口在，页面白：资源路径。核对 `locales`、pak、`resources_dir_path`。
- 关不掉或进程残留：没到 `OnBeforeClose`，或计数未归零，或外部泵仍在跑却没有 `CefQuitMessageLoop` 的等价物。
- 一点关闭就退出，页面来不及问用户：第一次系统关闭被放行，没有 `CloseBrowser(false)`。
- 和 Qt、Win32、GTK 一起卡死：同时跑了 `CefRunMessageLoop` 和宿主循环，或在 UI 线程上等待其他线程。

## 消息循环只留一种

- 独占：`CefRunMessageLoop`，结束用 `CefQuitMessageLoop`。
- 外部泵：`multi_threaded_message_loop = false`，主线程周期调用 `CefDoMessageLoopWork`。打开 `external_message_pump` 时仍要保留最大间隔。`OnScheduleMessagePumpWork` 可能在任意线程，投回主线程。不要在 `CefDoMessageLoopWork` 里面重入调用它。
- 独立 UI 线程：仅 Windows 和 Linux，`multi_threaded_message_loop = true`。不要再调用上面两个泵函数。`CefInitialize` 和 `CefShutdown` 留在主线程。

## 关闭骨架

父窗口第一次收到系统关闭：若 `IsClosing` 为假，`CloseBrowser(false)` 并取消这次系统关闭。`DoClose` 里对主浏览器设 `IsClosing` 并返回 `false`。第二次系统关闭放行。`OnBeforeClose` 丢掉 `CefRefPtr<CefBrowser>`。最后一个浏览器关闭后再停循环。循环退出之后才 `CefShutdown`。

`CloseBrowser(true)` 只用于已经决定不再询问页面的退出。

## 线程

生命周期回调在 UI 线程。不要在 `TID_IO` 上碰窗口。需要回到 UI 时用 `CefPostTask(TID_UI, ...)`。

macOS 的应用退出按 `cefsimple_mac.mm` 的顺序，不要把 `WM_CLOSE` 翻译成 Cocoa 调用了事。
