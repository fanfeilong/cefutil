# 启动、消息循环、关闭

核对日期：2026-10-13。依据 [General Usage](https://chromiumembedded.github.io/cef/general_usage.html) 的 Application Structure、Message Loop Integration、Browser Life Span，以及 [Tutorial](https://chromiumembedded.github.io/cef/tutorial.html)。

## 入口顺序

Windows 和 Linux 的单可执行文件入口是同一条路径。子进程和主进程共用它。

```cpp
int Main(int argc, char* argv[]) {
  CefMainArgs main_args(argc, argv);
  CefRefPtr<App> app(new App);

  const int exit_code = CefExecuteProcess(main_args, app, /*sandbox_info=*/nullptr);
  if (exit_code >= 0) {
    return exit_code;
  }

  CefSettings settings;
  CefString(&settings.cache_path).FromASCII("/path/to/cache");
  if (!CefInitialize(main_args, settings, app, /*sandbox_info=*/nullptr)) {
    return 1;
  }

  CefRunMessageLoop();
  CefShutdown();
  return 0;
}
```

Windows 的 `CefMainArgs` 收的是 `HINSTANCE`。`sandbox_info` 在启用沙箱时不是空指针，见 [发布与升级](10-ship-and-upgrade.md)。`CefSettings` 要在 `CefInitialize` 之前填完。缓存目录为空时，localStorage 一类存储不会跨进程生命周期保留。

`CefApp::GetBrowserProcessHandler` 返回的对象在 `OnContextInitialized` 里创建第一个浏览器。此时 CEF 的 UI 循环已经在跑，可以调用 `CefBrowserHost::CreateBrowser`。

## 三种消息循环

**CEF 自己的循环。** `CefRunMessageLoop` 阻塞到 `CefQuitMessageLoop`。适合 CEF 就是这个进程的 UI，或你可以把别的工作都投递到 `TID_UI`。

**接到已有循环上。** `multi_threaded_message_loop = false`，在主线程周期性调用 `CefDoMessageLoopWork`，同时继续派发你自己的窗口消息。调用太稀，页面饿死；调用太密，空转 CPU。

把 `CefSettings.external_message_pump` 打开后，CEF 会通过 `CefBrowserProcessHandler::OnScheduleMessagePumpWork` 要求你提前再泵一次。这个回调是补充，不是完整的工作清单。即使没收到回调，也要继续按你的最大间隔泵。回调可能出现在任意线程，先投回主线程。延迟小于等于 0 时尽快泵；延迟大于 0 时，用“请求的延迟”和“你的最大间隔”里较短的那个。官方样本把最大间隔放在大约 33 毫秒（`1000/30`），这是调度策略，不是帧率承诺。泵函数不可重入：如果一次 `CefDoMessageLoopWork` 里又进入了嵌套的原生消息循环，样本会把后续工作推迟到这次调用返回。你的程序如果在模态对话框期间仍要让页面动，需要在这条循环里继续泵。

**CEF 单独一条 UI 线程。** 仅 Windows 和 Linux。`multi_threaded_message_loop = true` 时不要再调用 `CefRunMessageLoop` 或 `CefDoMessageLoopWork`。`CefInitialize` 和 `CefShutdown` 仍在主线程。CEF 的浏览器 UI 调用在另一条线程上，和你的窗口线程通信要自己做（cefclient 在 Windows 上用隐藏消息窗口）。macOS 没有这个模式。

## 创建一个浏览器

`CefBrowserHost::CreateBrowser` 是异步的。`CreateBrowserSync` 在 UI 线程上同步返回浏览器对象。窗口信息、`CefClient`、初始 URL、`CefBrowserSettings` 在这次调用里给定。`OnAfterCreated` 才是“对象已经存在”的时刻，在这里保存 `CefRefPtr<CefBrowser>` 并增加计数。

父窗口是你自己的原生窗口时，`CefWindowInfo::SetAsChild` 把浏览器放进客户区。Views 宿主走另一套窗口创建，见 [页面放在哪里](08-hosting.md)。

## 关闭顺序

关闭不是 `CloseBrowser` 返回之后浏览器就没了。有原生父窗口时，顺序是：

1. 用户关闭父窗口。窗口过程收到系统关闭（Windows 上是 `WM_CLOSE`）。
2. 若还没进入关闭，调用 `GetHost()->CloseBrowser(false)`，并取消这次系统关闭。`false` 表示允许页面的 `onbeforeunload` 询问用户。
3. 页面同意后，CEF 调用 `CefLifeSpanHandler::DoClose`。这里把“下次系统关闭可以放行”的标志设上，并返回 `false`。返回 `false` 会让 CEF 再发一次系统关闭。
4. 父窗口再次收到系统关闭。标志已设，放行，窗口真正销毁。
5. `OnBeforeClose` 里释放浏览器指针。计数到 0 时调用 `CefQuitMessageLoop`（或停止你自己的泵）。
6. 消息循环退出后，主线程调用 `CefShutdown`。

`DoClose` 返回 `true` 表示你已经自己处理了关闭、不需要 CEF 再发系统关闭。窗口化浏览器若在这里返回 `true` 却没有真正关掉父窗口，进程会停在“浏览器还在、窗口关不掉”。

`CloseBrowser(true)` 跳过 `onbeforeunload`。用在你已经决定退出、不再问页面的时候。

最后一个浏览器没走到 `OnBeforeClose` 就调用 `CefShutdown`，会在退出时崩溃或挂住。多浏览器时用计数，不要在关第一个窗口时退出循环。

macOS 的退出还要处理应用激活和 Helper 进程，照 `cefsimple_mac.mm` 里的注释做，不要把 Windows 的 `WM_CLOSE` 直接译过去。
