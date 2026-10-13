---
name: cef-host-surface
description: Choose and implement how a CEF browser is drawn, using Views, a native parent window, or off-screen rendering. Use when embedding into an existing UI, compositing a page into video or a game view, or switching Chrome style and Alloy style. 也用于离屏渲染、嵌进已有窗口、Alloy style、Views、OSR。
---

# 决定 CEF 页面画在哪里

读 `docs/02-architecture.md` 和 `docs/08-hosting.md`。改 `CefWindowInfo` 之前先确定 style。

## 选择

- 新窗口就是浏览器，要权限框、PDF 或扩展：Chrome style，Views。这是 cefsimple 的默认。
- 已有 Qt、Win32、Cocoa 或 GTK 窗口，只要客户区里一块网页：`CefWindowInfo::SetAsChild`，父窗口关闭走两段式，见 skill `cef-lifecycle`。
- 要把像素送进自己的合成器、视频或游戏画面：Alloy style（`--use-alloy-style`）加 `SetAsWindowless`，实现 `CefRenderHandler`。

官方架构页写明 Chrome runtime 要求 Views 且不支持离屏渲染。不要在 Chrome style 上调用 `SetAsWindowless`，除非当前这份头文件和一次运行证明可以。

## 离屏实现清单

1. `CefRenderHandler` 里实现头文件未标可选的方法，包括 `GetViewRect` 和 `OnPaint`。
2. `SetAsWindowless` 时尽量传入父窗口，否则上下文菜单等能力不可用。
3. 尺寸变化调用 `WasResized`，然后等待新的 `GetViewRect` 和 `OnPaint`。
4. 鼠标、键盘、焦点用 `CefBrowserHost::Send*` 送入。CEF 不会从你的窗口过程里自动看到这些事件。
5. 软件 `OnPaint` 路径不做加速合成，性能低于有窗口的浏览器。只有头文件里存在 `OnAcceleratedPaint` 时才接共享纹理，并在目标平台上验证。
6. 销毁仍用 `CloseBrowser` 和 `CefLifeSpanHandler`。

坐标：有窗口时用 `CefDisplay` 做 DIP 与设备像素的转换。离屏时使用 `GetScreenInfo` 提供的 `device_scale_factor`。

## 验证

Views 或原生窗口：页面可见，焦点能进输入框，关闭顺序正确。离屏：至少画出一帧，拖动窗口或改变尺寸后画面更新，点击能点到页面元素。对照样本是带 `--off-screen-rendering-enabled` 的 cefclient。
