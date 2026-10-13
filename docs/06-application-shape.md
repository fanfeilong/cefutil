# 应用怎么拆：App、Client、Browser

核对日期：2026-10-13。依据 [General Usage](https://chromiumembedded.github.io/cef/general_usage.html) 的 CefApp、CefClient、CefBrowser，以及 [Tutorial](https://chromiumembedded.github.io/cef/tutorial.html)。

一个进程里通常只有一个 `CefApp`。一个 `CefClient` 可以被多个浏览器共享，也可以一个窗口一个 Client。浏览器实例是 `CefBrowser`，里面至少有一个主 `CefFrame`，每个 iframe 再有自己的 frame。

## CefApp：进程级

`CefApp` 在主进程和子进程里都会被用到，所以它必须能在两类进程里构造。按进程拆开：

- `OnBeforeCommandLineProcessing`：在 CEF 读命令行之后追加开关。关停“从进程命令行配置 CEF”时，把 `CefSettings.command_line_args_disabled` 设为 true，改由这个回调提供参数。
- `OnRegisterCustomSchemes`：注册自定义 scheme。要在所有进程一致，并且早于使用该 scheme 的浏览器。
- `GetBrowserProcessHandler`：只在浏览器进程被调用。`OnContextInitialized` 放在这里。
- `GetRenderProcessHandler`：只在渲染进程被调用。V8 上下文、进程消息的渲染侧在这里。

`GetBrowserProcessHandler` 若返回 `this`，`CefApp` 的实现类就要同时继承 `CefBrowserProcessHandler`。渲染侧同理。不要在浏览器进程里做 V8 操作，那个接口不会在这里被调用。

## CefClient：浏览器实例级

`CefClient` 的 getter 返回各个 Handler。没实现的 getter 返回空，CEF 用默认行为。和集成最相关的几个：

| Handler | 作用 |
| --- | --- |
| `CefLifeSpanHandler` | 创建之后、是否允许关闭、关闭前释放 |
| `CefLoadHandler` | 开始加载、结束、错误 |
| `CefDisplayHandler` | 标题、地址、状态、控制台消息 |
| `CefRequestHandler` | 导航、资源请求、证书错误 |
| `CefRenderHandler` | 离屏渲染时的尺寸和像素 |
| `CefDialogHandler` / `CefJSDialogHandler` | 文件框、`alert`、`onbeforeunload` |
| `CefKeyboardHandler` / `CefContextMenuHandler` | 快捷键和菜单要不要交给宿主 |

`OnProcessMessageReceived` 在 `CefClient` 自己身上，是浏览器进程收到渲染进程消息的入口。

回调在哪条线程上，以头文件注释为准。生命周期回调在 UI 线程。网络相关回调经常在 IO 线程。在 IO 线程里弹原生对话框，或在 UI 线程里做同步网络，都会把对应的进程卡住。

## CefBrowser 和 CefFrame

两侧进程都有这两个类型，标识符一致。只有浏览器进程里 `GetHost()` 是可用的：取原生窗口句柄、关闭、发送输入、调整离屏尺寸、打开 DevTools，都走 Host。

加载地址用 `GetMainFrame()->LoadURL`。取 HTML 是异步的，`GetSource` 接受一个 `CefStringVisitor`，在回调里才能读到字符串。历史导航用 `GoBack` / `GoForward`。

`CefBrowserSettings` 是单个浏览器的设置，`CefSettings` 是整个进程的设置。背景色、默认字体、是否允许本地文件访问页面里的通用设置，放在浏览器设置里；缓存目录、日志、子进程路径、消息循环模式，放在进程设置里。

## 命令行从哪来

Chromium 的大量行为是 `--开关`。三条路：

1. 进程启动参数，经 `CefMainArgs` 进入。子进程的参数由 CEF 自己追加，宿主要原样保留。
2. `CefApp::OnBeforeCommandLineProcessing` 里追加。
3. `CefBrowserProcessHandler::OnBeforeChildProcessLaunch` 里给子进程追加你自己的参数。

调试时给 cefclient 或你的程序加上 `--remote-debugging-port=9222`，用 Chrome 打开 `http://localhost:9222`。这个开关也可以写成 `CefSettings.remote_debugging_port`。
