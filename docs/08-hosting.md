# 页面放在哪里

核对日期：2026-10-13。依据 [Tutorial](https://chromiumembedded.github.io/cef/tutorial.html) 的 Architecture in 60 Seconds，[General Usage](https://chromiumembedded.github.io/cef/general_usage.html) 的 Off-Screen Rendering，以及 [Architecture](https://chromiumembedded.github.io/cef/architecture.html)。

CEF 不决定你的产品窗口长什么样。它提供三种放置方式。cefsimple 默认是 Chrome style 加 Views。

## Views

Views 是 CEF 自带的跨平台控件树。Chrome style 依赖它：工具栏、权限气泡和浏览器视图都是 Views 控件。cefsimple 的默认路径就是这种。

适合新窗口本身就是一块浏览器、希望三个系统行为接近的程序。窗口对象、焦点和弹层由 CEF 管。你仍通过 `CefClient` 拦截导航、下载和 JS 对话框。

## 原生父窗口

`CefWindowInfo::SetAsChild` 把浏览器放进已经存在的 `HWND`、`NSView` 或 X 窗口。宿主继续用自己的 UI 框架画菜单和工具栏，CEF 只占客户区里的一块。

这种窗口的关闭必须走 [启动、消息循环、关闭](05-lifecycle.md) 里的两段式：第一次系统关闭只调用 `CloseBrowser(false)`，`DoClose` 之后的第二次才销毁父窗口。

坐标有两套。Views 和多数浏览器回调用 DIP（左上角原点）。Windows API 要设备像素，用 `CefDisplay` 转换。macOS 的 Cocoa 也用点，但原点在左下。离屏浏览器要用 `CefRenderHandler::GetScreenInfo` 里最近一次的 `device_scale_factor`。浏览器缩放为 100% 时，DIP 和 CSS 像素的比例才是 1。

## 离屏渲染

离屏渲染不创建浏览器原生窗口。CEF 把脏矩形和像素缓冲交给你，你把鼠标、键盘和焦点用 `CefBrowserHost::Send*` 送回去。

启用步骤：

1. 实现 `CefRenderHandler`。头文件里没标可选的方法都要实现。
2. `CefWindowInfo::SetAsWindowless` 之后再 `CreateBrowser`。不传父窗口时，上下文菜单等依赖父窗口的功能会缺失。
3. `GetViewRect` 返回你要的视图大小。
4. `OnPaint` 收到像素缓冲和脏区。样本用 OpenGL 画，你可以用自己的合成器。
5. 尺寸变化时调用 `WasResized`。CEF 会再要一次 `GetViewRect`，然后 `OnPaint`。
6. 输入用 `SendMouse*`、`SendKeyEvent`、`SendFocusEvent` 等。
7. 销毁仍是 `CloseBrowser`，生命周期回调和有窗口时一样。

官方 General Usage 写明：这条软件像素路径目前不支持加速合成，性能会低于有窗口的浏览器。较新的头文件里有 `OnAcceleratedPaint`，用共享纹理把 GPU 帧交出来。这是追加的回调，不是所有分支、所有 style 都有。以你下载的 `cef_render_handler.h` 为准。

架构页的约束要一起看：Chrome runtime 在文档撰写时要求 Views，并且不支持离屏渲染。要离屏时用 Alloy style（Chrome bootstrap 加上 `--use-alloy-style`），并在目标二进制上做一次最小窗口验证。cefclient 用 `--off-screen-rendering-enabled` 可以对照官方样本。

OBS 这一类把网页合成进视频的程序，走的就是离屏而不是 Views。它们接受像素路径的复杂度和版本跟进成本，换取页面能进入自己的画面。
