# 架构：两层 Chromium，一种启动方式

核对日期：2026-10-13。依据 [Architecture](https://chromiumembedded.github.io/cef/architecture.html)。

## Chromium 里和嵌入有关的两层

现代 Chromium 里，嵌入者会碰到两层：

- **content。** 核心对象是 `WebContents`。它包含多进程、网络、渲染和大部分通用 Web 平台。
- **chrome。** 核心对象是 `Browser`。权限对话框、PDF、扩展、自动填充、设备选择这些能力在这一层，而不是 content 里。

CEF3 从 2013 年到 2020 年前后，浏览器窗口做在 content 的 `WebContents` 上。这套实现后来被命名为 **Alloy runtime**。2020 年起，CEF 把同一套对外 API 又做在 chrome 的 `Browser` 上，称为 **Chrome runtime**。这样权限、PDF、扩展等不必再从 chrome 层逐个摘进 CEF。

## 现在的启动方式

自 **M128**（2024 年中）起，Alloy 的独立启动路径（bootstrap）已经删除。进程启动统一走 **Chrome bootstrap**。窗口外观仍有两种 style，都建在这条启动路径上：

| Style | 建立在 | 你得到什么 | 典型开关 |
| --- | --- | --- | --- |
| Chrome style | chrome 层 `Browser` | Chrome 的工具栏、权限和设备对话框、PDF、扩展等 | cefsimple 的默认 |
| Alloy style | content 层 `WebContents` | 传统 CEF 窗口，没有整套 Chrome UI | `--use-alloy-style` |

官方在架构页写明：撰写当时，Chrome runtime 要求使用 Views，并且不支持离屏渲染。Alloy style 仍是离屏渲染和“只嵌一块页面、不要 Chrome 壳”的路径。发行前打开你这份二进制的 `include/internal/cef_types_win.h`（以及对应平台头）和 `CefWindowInfo`，确认 `SetAsWindowless` 在选定 style 下是否可用。

内部类名可以帮助读源码，不必在业务代码里依赖它们：

- 启动与各进程的 Content 回调：`ChromeMainDelegateCef`、`ChromeContentClientCef`、`ChromeContentBrowserClientCef`、`ChromeContentRendererClientCef`。
- 浏览器进程里的全局上下文：`CefContext`，由 `CefInitialize` 创建、`CefShutdown` 销毁。
- 浏览器进程里对外的 `CefBrowser` / `CefBrowserHost`：Chrome style 是 `ChromeBrowserHostImpl`，Alloy style 是 `AlloyBrowserHostImpl`。
- 渲染进程里的 `CefBrowser`：`CefBrowserImpl`。

## 进程

浏览器进程就是宿主进程。窗口、大多数应用逻辑、网络入口在这里。Blink 和 V8 在独立的渲染进程里，默认按源（scheme + 域名）起进程。GPU、网络服务等进程按需出现。

Windows 和 Linux 上，同一个可执行文件会被多次拉起，靠命令行区分进程类型，入口里要先调用 `CefExecuteProcess`。macOS 上子进程必须是单独的 Helper app，并且按系统要求放进 app bundle。

进程之间只有异步消息。浏览器进程里的 C++ 对象不能拿去给渲染进程里的 JavaScript 直接调用。`CefBrowser` 和 `CefFrame` 在两侧都有对象，两侧的标识符一致，用来把消息送回同一个 frame。

## 库的边界

`libcef` 导出 C API。这层把宿主和 Chromium 的编译器、堆、内部类型隔开。`libcef_dll_wrapper` 是生成出来的 C++ 包装，以源码形式放在二进制包里，编进你的程序。业务代码调用的是包装后的 `CefBrowser`、`CefClient` 这些类。包装层的生成方式和版本宏见 [二进制边界](03-binary-boundary.md)。

## 和样本程序的对应关系

[教程](https://chromiumembedded.github.io/cef/tutorial.html) 用 `cefsimple` 把结构压成四块：

1. 入口：区分主进程和子进程，初始化，跑消息循环。
2. `CefApp`（样本里的 `SimpleApp`）：进程级回调。浏览器进程的逻辑在 `CefBrowserProcessHandler`，渲染进程的逻辑在 `CefRenderProcessHandler`。
3. `CefClient`（样本里的 `SimpleHandler`）：某个浏览器实例的回调。
4. `CefBrowserHost::CreateBrowser`：在 `OnContextInitialized` 里创建窗口。最后一个窗口的 `OnBeforeClose` 里退出消息循环。

`cefclient` 是功能全集，包含 Chrome style、Alloy style、离屏渲染、自定义 scheme 和请求拦截。新集成先跑通 `cefsimple` 的路径，再按缺的能力从 `cefclient` 里摘。
