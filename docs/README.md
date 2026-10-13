# 文档地图

这些文章按依赖顺序写。前一篇的结论是后一篇的前提。

| 顺序 | 文章 | 读完能回答 |
| --- | --- | --- |
| 1 | [CEF 在产品里的位置](01-place.md) | 该不该用 CEF，而不是 Electron 或系统 WebView |
| 2 | [架构](02-architecture.md) | Chrome bootstrap、Chrome style 和 Alloy style、进程各管什么 |
| 3 | [二进制边界](03-binary-boundary.md) | 为什么有 C API，`CEF_API_VERSION` 钉住什么 |
| 4 | [进程和线程](04-processes-and-threads.md) | 代码该写在哪个进程、哪条线程 |
| 5 | [启动、消息循环、关闭](05-lifecycle.md) | 入口顺序，以及窗口为什么要关两次 |
| 6 | [应用怎么拆](06-application-shape.md) | `CefApp`、`CefClient`、`CefBrowser` 各放什么 |
| 7 | [JavaScript 和原生代码](07-javascript.md) | 页面怎么异步调用宿主 |
| 8 | [页面放在哪里](08-hosting.md) | Views、原生父窗口、离屏渲染 |
| 9 | [嵌进已有程序](09-integrate.md) | 从下载二进制包到第一次关掉窗口 |
| 10 | [发布与升级](10-ship-and-upgrade.md) | 分支、沙箱、换包时要重编什么 |

编程代理的操作技能在 [`skills/`](../skills/README.md)，顺序是 `cef-setup-env`、`cef-download`、`cef-build`、`cef-create-project`、`cef-integrate`。动手前读第 2、5、9 篇；做 JavaScript 桥再读第 7 篇；做离屏渲染再读第 8 篇。

官方说明书仍是接口细节的来源：<https://chromiumembedded.github.io/cef/>。本文不复制那套页面，只保留集成时会反复用到的结构。
