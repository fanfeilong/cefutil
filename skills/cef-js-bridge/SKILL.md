---
name: cef-js-bridge
description: Add or repair an asynchronous JavaScript-to-native bridge in a CEF host using CefMessageRouter. Use when the page must call C++ or the existing binding uses CefRegisterExtension, native function, or V8 calls from the browser process. 也用于页面调用 C++、JS 桥、cefQuery、替换已废弃的 V8 扩展。
---

# 接 CEF 的 JavaScript 桥

读 `docs/07-javascript.md`。V8 只存在于渲染进程的 `TID_RENDERER`。浏览器进程里构造 `CefV8Value` 是错误的。

## 选法

- 宿主只需要执行一段脚本：`CefFrame::ExecuteJavaScript`。它没有返回值。
- 要在页面上挂一个随文档创建的函数，且工作可以留在渲染进程：`CefRenderProcessHandler::OnContextCreated` 里写 window binding。
- 页面要调用宿主的文件、窗口或业务对象：用 `include/wrapper/cef_message_router.h` 的 `CefMessageRouter`。这是默认选择。

`CefRegisterExtension` 和 `native function` 自 M154 起废弃。看到旧代码时改成消息路由，不要修扩展字符串。

## 路由接线

两侧使用同一份 router 配置。

渲染进程的 `CefRenderProcessHandler`：

- 创建 renderer 侧 router。
- `OnContextCreated` 和 `OnContextReleased` 转给 router。
- `OnProcessMessageReceived` 先交给 router，已处理则返回 true。

浏览器进程：

- 创建 browser 侧 router，注册 `Handler`。
- 把头文件要求的生命周期、`OnBeforeBrowse`、`OnProcessMessageReceived` 全部转给 router。漏掉导航回调时，查询会落到已经卸载的文档上。

`Handler::OnQuery` 里不要做长时间阻塞。需要别的线程时，保存 callback，投递任务，完成后调用一次 `Success` 或 `Failure`。

参数用字符串，结构化数据用 JSON。在 `OnContextReleased` 之后不得再使用该 frame 的 V8 句柄。

## 最小验证

页面调用一次查询，原生侧返回固定字符串，页面把字符串显示出来。然后导航到另一个 URL，确认旧查询不会在新文档上触发。打开 `remote_debugging_port`，在 DevTools 里看渲染进程的异常。
