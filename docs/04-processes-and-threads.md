# 进程和线程

核对日期：2026-10-13。依据 [General Usage](https://chromiumembedded.github.io/cef/general_usage.html) 的 Processes、Threads、Posting Tasks、IPC。

## 进程各管什么

| 进程 | 里面有什么 | 你的代码通常放哪 |
| --- | --- | --- |
| Browser | 窗口、浏览器对象、网络入口、大多数 Handler | 宿主的主逻辑、`CefClient` |
| Renderer | Blink、V8、该源的 frame | `CefRenderProcessHandler`、V8 绑定 |
| GPU | 合成与加速 | 一般不写业务 |
| 其他服务进程 | 网络服务等 | 一般不写业务 |

渲染进程默认按源拆分。同一个页面里不同源的 iframe 会进不同进程。因此“在 OnContextCreated 里拿到的对象”只在当前渲染进程、当前 frame 的上下文里有效。

子进程由 CEF 拉起，命令行里带有进程类型。入口必须把原始命令行交进去：

- Linux / macOS：`CefMainArgs(argc, argv)`。
- Windows：`CefMainArgs(hInstance)`。

`CefExecuteProcess` 在子进程里会阻塞到该进程该退出，并返回一个大于等于 0 的退出码。主进程里它立刻返回负数，随后才是 `CefInitialize`。子进程如果没走到 `CefExecuteProcess`，页面不会起来，日志里通常是子进程秒退。

主程序很大、或不适合当子进程时，设置 `CefSettings.browser_subprocess_path` 指向专用可执行文件。macOS 必须这么做，并且 Helper 要先进沙箱初始化、再加载框架。Windows 上若启用沙箱，M138 之后不能再用这条设置把子进程指到另一个 exe，见 [发布与升级](10-ship-and-upgrade.md)。

## 线程

每个进程里有多条线程。常用的几个：

- `TID_UI`：浏览器进程的主线程。`multi_threaded_message_loop` 为 false 时，它就是调用 `CefInitialize` 的那条线程。
- `TID_IO`：浏览器进程里处理 IPC 和网络。
- `TID_FILE_*`：碰文件系统。阻塞操作放这里，或放你自己创建的线程。
- `TID_RENDERER`：渲染进程的主线程。所有 Blink 和 V8 调用都在这里。

头文件 `include/wrapper/cef_helpers.h` 里的 `CEF_REQUIRE_UI_THREAD()` 等宏，只在调试构建里检查当前线程。回调写错线程时，发布构建不会帮你挡。

跨线程用 `CefPostTask`，不要在 UI 线程上锁等待渲染进程。`include/wrapper/cef_closure_task.h` 把 `base::BindOnce` 收成可以投递的任务。需要保存某条线程的投递器时，用 `CefTaskRunner::GetForThread`。

多线程碰到的成员用 `include/base/cef_lock.h` 里的 `base::AutoLock`。更常见的做法是：状态只放在拥有它的那条线程上，别的线程只投递任务。

## 进程之间怎么说话

`CefFrame::SendProcessMessage` 把 `CefProcessMessage` 送到对侧。渲染进程发往浏览器进程用 `PID_BROWSER`，反方向用 `PID_RENDERER`。对侧入口是：

- 渲染进程收到：`CefRenderProcessHandler::OnProcessMessageReceived`。
- 浏览器进程收到：`CefClient::OnProcessMessageReceived`。

消息是拷贝，不是共享内存对象。参数放进 `CefListValue`。返回值也是另一条消息，或者走下一节的消息路由。

创建浏览器时可以带 `extra_info`（`CefDictionaryValue`）。它会在渲染进程的 `CefRenderProcessHandler::OnBrowserCreated` 里交给和这个浏览器相关的每个渲染进程。适合放只读的启动配置，不适合放之后还要改的状态。

JavaScript 要调用宿主时，绑定发生在渲染进程，真正做事的代码多半在浏览器进程。这条链路必须是异步的。同步等待会卡住渲染进程。官方把“用同步 XHR 打到自定义 scheme”标成少见的退路，页面会在等待期间失去响应。产品代码用 [JavaScript 桥](07-javascript.md) 里的消息路由。
