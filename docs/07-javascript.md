# JavaScript 和原生代码

核对日期：2026-10-13。依据 [JavaScript Integration](https://chromiumembedded.github.io/cef/javascript_integration.html) 和 [General Usage](https://chromiumembedded.github.io/cef/general_usage.html) 的 Asynchronous JavaScript Bindings。

V8 只在渲染进程的 `TID_RENDERER` 上跑。每个 frame 有自己的 `CefV8Context`。`OnContextCreated` 到 `OnContextReleased` 是这个上下文的寿命。上下文释放之后不能再碰它，也不能把 V8 对象一直留在成员里。

## 三条由浅到深的路

**宿主执行一段脚本。** `CefFrame::ExecuteJavaScript` 在浏览器进程和渲染进程都可以调用，不要求当前已经进入 V8 上下文。适合改 DOM、调用页面上已经存在的函数。它没有返回值。需要结果时，让页面再发一条消息回来。

**往 `window` 上挂值。** 在 `CefRenderProcessHandler::OnContextCreated` 里取 `context->GetGlobal()`，用 `CefV8Value` 挂属性或函数。每次文档创建都会再来一次，所以可以按 frame 给不同的能力。函数的实现是 `CefV8Handler::Execute`，它仍然跑在渲染进程里。

**异步桥到浏览器进程。** 页面要读文件、调你的业务对象、改宿主窗口时，渲染进程只负责收调用、发消息、在回来时执行 JS 回调。官方提供的实现是 `CefMessageRouter`，头文件在 `include/wrapper/cef_message_router.h`，独立样本是二进制包里的 message_router 示例。

## 不要再用的扩展写法

`CefRegisterExtension` 和 `native function` 这套 V8 扩展，官方在 **M154** 标为废弃。扩展在 DOM 存在之前注入，加载时访问 DOM 会崩溃，而且一旦注册就不能按 frame 修改。新代码用 `OnContextCreated` 的 window binding，或用消息路由。

## 消息路由怎么接

路由分两侧，配置（JS 函数名等）必须一致。

渲染进程：

- 在 `OnWebKitInitialized` 或路由要求的初始化点创建 renderer 侧 router。
- `OnContextCreated` 里 `OnContextCreated`，让路由注入 `window.cefQuery` 一类函数（具体名字由配置决定）。
- `OnContextReleased` 里通知路由释放。
- `OnProcessMessageReceived` 里把消息交给路由，返回值表示已处理。

浏览器进程：

- 创建 browser 侧 router，并注册一个或多个 `Handler`。
- `CefLifeSpanHandler`、`CefRequestHandler::OnBeforeBrowse`、`CefClient::OnProcessMessageReceived` 里按头文件要求的那几个回调转发给路由。漏掉 `OnBeforeBrowse` 时，跨文档的查询会串到旧页面上。

`Handler::OnQuery` 在浏览器进程执行。处理完调用 `callback->Success` 或 `Failure`。可以先返回 true 把查询留住，等工作完成再调用 callback。callback 只能调用一次。

查询的参数是字符串。结构化数据在两侧自己约定 JSON。不要把 `CefV8Value` 或裸函数指针塞进进程消息。

## 自己写消息时要守的规则

不用路由、直接 `SendProcessMessage` 时：

- 渲染侧用 map 记住 `(消息名, browser id)` 到 `(CefV8Context, 函数)` 的对应。
- 执行回调前 `context->Enter()`，之后 `Exit()`。
- `OnContextReleased` 里把这个 context 的条目删掉。留着的 V8 句柄会在文档销毁后被用到。
- 浏览器侧的结果用另一条消息送回，在渲染侧再执行函数。

`CefV8Value::SetRethrowExceptions(true)` 只在调用栈上面还有 JS 时使用。从原生菜单直接进 `ExecuteFunction` 时不要重抛，否则异常没人接，进程会坏掉。
