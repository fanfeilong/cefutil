我来说windows下CEF3.2171的关闭流程，里面会引用一部分官方库的文档和个人的伪代码，为了辅助理解——
以下是截取自cef_life_span_handler.h的头文件文档，所以一部分文档他还是写在头文件里的，根据他的流程，能很快的去梳理相关逻辑
// The CefLifeSpanHandler::OnBeforeClose() method will be called immediately
  // before the browser object is destroyed. The application should only exit
  // after OnBeforeClose() has been called for all existing browsers.
  //
  // If the browser represents a modal window and a custom modal loop
  // implementation was provided in CefLifeSpanHandler::RunModal() this callback
  // should be used to restore the opener window to a usable state.
  //
  // By way of example consider what should happen during window close when the
  // browser is parented to an application-provided top-level OS window.
  // 1.  User clicks the window close button which sends an OS close
  //     notification (e.g. WM_CLOSE on Windows, performClose: on OS-X and
  //     "delete_event" on Linux).
  // 2.  Application's top-level window receives the close notification and:
  //     A. Calls CefBrowserHost::CloseBrowser(false).
  //     B. Cancels the window close.
  // 3.  JavaScript 'onbeforeunload' handler executes and shows the close
  //     confirmation dialog (which can be overridden via
  //     CefJSDialogHandler::OnBeforeUnloadDialog()).
  // 4.  User approves the close.
  // 5.  JavaScript 'onunload' handler executes.
  // 6.  Application's DoClose() handler is called. Application will:
  //     A. Set a flag to indicate that the next close attempt will be allowed.
  //     B. Return false.
  // 7.  CEF sends an OS close notification.
  // 8.  Application's top-level window receives the OS close notification and
  //     allows the window to close based on the flag from #6B.
  // 9.  Browser OS window is destroyed.
  // 10. Application's CefLifeSpanHandler::OnBeforeClose() handler is called and
  //     the browser object is destroyed.
  // 11. Application exits by calling CefQuitMessageLoop() if no other browsers
  //     exist.

pre1. 如果你要管理CefBrowser的生命周期，意味者你必须实现相关 CefLifeSpanHandler接口，在OnAfterCreated里管理和获取CefBrowser的每一个browser，在DoClose和OnBeforeClose里管理关闭
pre2. 这里要注意整个流程对应开发者来说不是线性代码，都是基于消息、事件投递、接口切面的一种编程习惯，大家要对异步接口的开发有心理准备，你实现的一个接口都是等待框架或者别人调用。

1. 用户发起了一个系统级别的关闭操作，比如CLOSE，window下会被WM_CLOSE和相关封装触发
2. 系统的顶层窗口消息处理会获取到CLOSE消息
&ensp;&ensp;A. 调用  CefBrowserHost::CloseBrowser(false), 通知HOST开始关闭流程，参数force_close=false，我们不需要交互式关闭。
&ensp;&ensp;B. 取消系统级关闭，比如以下代码我"return 0; "拒接了系统窗口的关闭流程，进入了个人和CEF的流程
```c++
     case WM_CLOSE:
      if (client && !client->IsClosing()) {
        CefRefPtr<CefBrowser> browser = client->GetBrowser();
        if (browser.get()) {
          // Notify the browser window that we would like to close it. This
          // will result in a call to ClientHandler::DoClose() if the
          // JavaScript 'onbeforeunload' event handler allows it.
          browser->GetHost()->CloseBrowser(false);

          // Cancel the close.
          return 0;
        }
      }

      // Allow the close.
      break;
```
这是单browser关闭
```C++
browser->GetHost()->CloseBrowser(false);
```
同样可以替换成多browser关闭
```C++
void ClientHandler::CloseAllBrowsers(bool force_close) {
  if (!CefCurrentlyOn(TID_UI)) {
    // Execute on the UI thread.
    CefPostTask(TID_UI,
        base::Bind(&ClientHandler::CloseAllBrowsers, this, force_close));
    return;
  }

  if (!popup_browsers_.empty()) {
    // Request that any popup browsers close.
    BrowserList::const_iterator it = popup_browsers_.begin();
    for (; it != popup_browsers_.end(); ++it)
      (*it)->GetHost()->CloseBrowser(force_close);
  }

  if (browser_.get()) {
    // Request that the main browser close.
    browser_->GetHost()->CloseBrowser(force_close);
  }
}
```
intermittent.
最终会触发一个或多个CloseBrowser => CefBrowserHostImpl::CloseBrowser => CefBrowserHostImpl::CloseContents
``` C++
void CefBrowserHostImpl::CloseBrowser(bool force_close) {
  if (CEF_CURRENTLY_ON_UIT()) {
    // Exit early if a close attempt is already pending and this method is
    // called again from somewhere other than WindowDestroyed().
    if (destruction_state_ >= DESTRUCTION_STATE_PENDING &&
        (IsWindowless() || !window_destroyed_)) {
      if (force_close && destruction_state_ == DESTRUCTION_STATE_PENDING) {
        // Upgrade the destruction state.
        destruction_state_ = DESTRUCTION_STATE_ACCEPTED;
      }
      return;
    }

    if (destruction_state_ < DESTRUCTION_STATE_ACCEPTED) {
      destruction_state_ = (force_close ? DESTRUCTION_STATE_ACCEPTED :
                                          DESTRUCTION_STATE_PENDING);
    }

    content::WebContents* contents = web_contents();
    if (contents && contents->NeedToFireBeforeUnload()) {
      // Will result in a call to BeforeUnloadFired() and, if the close isn't
      // canceled, CloseContents().
      contents->DispatchBeforeUnload(false);
    } else {
      CloseContents(contents);
    }
  } else {
    CEF_POST_TASK(CEF_UIT,
        base::Bind(&CefBrowserHostImpl::CloseBrowser, this, force_close));
  }
}
```
我们跳过其他细节（ 我也没关注:)  ），和相关CEF_POST_TASK把函数投递到UI线程上去处理的线程类型,最后会进入 CloseContents
```C++
void CefBrowserHostImpl::CloseContents(content::WebContents* source) {
  if (destruction_state_ == DESTRUCTION_STATE_COMPLETED)
    return;

  bool close_browser = true;

  // If this method is called in response to something other than
  // WindowDestroyed() ask the user if the browser should close.
  if (client_.get() && (IsWindowless() || !window_destroyed_)) {
    CefRefPtr<CefLifeSpanHandler> handler =
        client_->GetLifeSpanHandler();
    if (handler.get()) {
      close_browser = !handler->DoClose(this);   // CefLifeSpanHandler::DoClose
    }
  }

  if (close_browser) {
    if (destruction_state_ != DESTRUCTION_STATE_ACCEPTED)
      destruction_state_ = DESTRUCTION_STATE_ACCEPTED;

    if (!IsWindowless() && !window_destroyed_) {
      // A window exists so try to close it using the platform method. Will
      // result in a call to WindowDestroyed() if/when the window is destroyed
      // via the platform window destruction mechanism.
      PlatformCloseWindow();
    } else {
      // Keep a reference to the browser while it's in the process of being
      // destroyed.
      CefRefPtr<CefBrowserHostImpl> browser(this);

      // No window exists. Destroy the browser immediately.
      DestroyBrowser();       // CefLifeSpanHandler::OnBeforeClose
      if (!IsWindowless()) {
        // Release the reference added in PlatformCreateWindow().
        Release();
      }
    }
  } else if (destruction_state_ != DESTRUCTION_STATE_NONE) {
    destruction_state_ = DESTRUCTION_STATE_NONE;
  }
}
```
在CloseContents里发现了我们的DoClose接口，其实这个函数里还有我们要处理的OnBeforeClose接口，DestroyBrowser包含了他
```C++
void CefBrowserHostImpl::DestroyBrowser() {
  CEF_REQUIRE_UIT();

  destruction_state_ = DESTRUCTION_STATE_COMPLETED;

  if (client_.get()) {
    CefRefPtr<CefLifeSpanHandler> handler = client_->GetLifeSpanHandler();
    if (handler.get()) {
      // Notify the handler that the window is about to be closed.
      handler->OnBeforeClose(this);        //`CefLifeSpanHandler::OnBeforeClose
    }
  }
  // 省略
}
```
到这里我们要有这样的调用栈初步概念（之后会修改..)
```C++
=> 表示同步调用 ->异步调用
[os]window close =>
    [个人]browser(s) close =>
        [cef]CefBrowserHostImpl::CloseBrowser =>
            [cef]CefBrowserHostImpl::CloseContents =>
                [实现cef接口]CefLifeSpanHandler::DoClose =>
                [cef]CefLifeSpanHandler::DestroyBrowser =>
                    [实现cef接口]CefLifeSpanHandler::OnBeforeClose =>
[os]阻止window close
```
3-5. js和交互级的操作，略，其中有一点js onunload让你有机会在js层面去通知关闭，比如通知后台的native 线程去做自我清理动作，特别是cefQuery持久化以后的和JS交互的native线程
6.DoClose接口被调用
&ensp;&ensp;A. 允许设置一些状态表示browser正在关闭
&ensp;&ensp;B. "return false;" 我们不需要紧急关闭（好像是预留接口，或者其他平台用)
7.cef 发送一个系统级别的关闭动作发现之前的调用栈已经需要修改了，因为这里异步操作)
CloseContents => PlatformCloseWindow =>  CefBrowserHostImpl::PlatformCloseWindow =>PostMessage WM_CLOSE
```C++
void CefBrowserHostImpl::PlatformCloseWindow() {
  if (window_info_.window != NULL) {
    HWND frameWnd = GetAncestor(window_info_.window, GA_ROOT);
    PostMessage(frameWnd, WM_CLOSE, 0, 0);
  }
}
```
PS: host公共代码在这个文件browser_host_impl.cc里，平台特殊接口在browser_host_impl_xxx.cc里
PlatformCloseWindow在browser_host_impl_win.cc里

允许关闭"if (client && !client->IsClosing()) " 之前我代码里IsClosing已经允许（
```C++
 case WM_CLOSE:
      if (client && !client->IsClosing()) {
// 这次不会进入IF分支
```
8.系统窗口已销毁
```C++
LRESULT CALLBACK CefBrowserHostImpl::WndProc(HWND hwnd, UINT message,
                                             WPARAM wParam, LPARAM lParam) {
// 略
  case WM_DESTROY:
    if (browser) {
      // Clear the user data pointer.
      gfx::SetWindowUserData(hwnd, NULL);

      // Force the browser to be destroyed and release the reference added in
      // PlatformCreateWindow().
      browser->WindowDestroyed();
    }
    return 0;
// 略
  }
```
cef竟然在等系统群发销毁消息然后处理的代码，看到这里才发现CefBrowserHostImpl这个类是CefBrowser和系统整合的桥梁，不同平台对 CefBrowserHostImpl 的实现也不一样

9.框架发起销毁操作
10. 调用OnBeforeClose
WM_DESTROY => CefBrowserHostImpl::WindowDestroyed => ClientHandler::OnBeforeClose
```C++
void ClientHandler::OnBeforeClose(CefRefPtr<CefBrowser> browser) {
  CEF_REQUIRE_UI_THREAD();

  message_router_->OnBeforeClose(browser);

  if (GetBrowserId() == browser->GetIdentifier()) {
    {
      base::AutoLock lock_scope(lock_);
      // Free the browser pointer so that the browser can be destroyed
      browser_ = NULL;
    }

    if (osr_handler_.get()) {
      osr_handler_->OnBeforeClose(browser);
      osr_handler_ = NULL;
    }
  } else if (browser->IsPopup()) {
    // Remove from the browser popup list.
    BrowserList::iterator bit = popup_browsers_.begin();
    for (; bit != popup_browsers_.end(); ++bit) {
      if ((*bit)->IsSame(browser)) {
        popup_browsers_.erase(bit);
        break;
      }
    }
  }

  if (--browser_count_ == 0) {
    // All browser windows have closed.
    // Remove and delete message router handlers.
    MessageHandlerSet::const_iterator it = message_handler_set_.begin();
    for (; it != message_handler_set_.end(); ++it) {
      message_router_->RemoveHandler(*(it));
      delete *(it);
    }
    message_handler_set_.clear();
    message_router_ = NULL;

    // Quit the application message loop.
    //AppQuitMessageLoop();
  XCefAppManage::Instance()->QuitMessageLoop();
  }
}
```
我判断到 "--browser_count_ == 0" ,  开始发起真实关闭。

11. 发起的真实关闭逻辑
个人在OnBeforeClose里销毁所有browser和相关操作cefquery组件销毁以后，发起了总的关闭流程（关闭是异步的）
```C++
void            XCefAppManage::QuitMessageLoop() {
  if (GetCefSettings().multi_threaded_message_loop)
  {
    // Running in multi-threaded message loop mode. Need to execute
    // PostQuitMessage on the main application thread.
    if (NULL == message_wnd__)
    {
      message_wnd__ = ::FindWindow(XWinUtil::GetMessageWindowClassName(XWinUtil::GetParentProcessID()), NULL);
    }
    DCHECK(message_wnd__);
    PostMessage(message_wnd__, WM_COMMAND, ID_QUIT, 0);
  }
  else {
    CefQuitMessageLoop();
  }
}
```
如果CefRunMessageLoop创建的消息循环对应CefQuitMessageLoop去终结，其他方式的处理就千奇百怪了，我根据cefclient的例子用一个隐藏窗口来结束系统的消息循环
```C++
PostMessage(message_wnd__, WM_COMMAND, ID_QUIT, 0);
```
处理非常简单，调用了原生API PostQuitMessage退出循环
```C++
    case ID_QUIT:
      PostQuitMessage(0);
      return 0;
```
PS: 如果需要把CEF混入其他框架的话，不同框架的处理都需要适配，切记。
12. CefShutDown() 销毁Cef资源，如果有资源没关闭，比如browser的计数不为0，debug下会check assert。
最后我再来修改一下之前的调用栈
```C++
=> 表示同步调用 ->异步调用
[os]on window close WM_CLOSE =>
    [个人]browser(s) close =>
        [cef]CefBrowserHostImpl::CloseBrowser =>
            [cef]CefBrowserHostImpl::CloseContents =>
                [实现cef接口]CefLifeSpanHandler::DoClose =>
                [cef]CefBrowserHostImpl::PlatformCloseWindow =>
                    [os] post window close -> WM_CLOSE
[os]阻止window close

post被接收
[os] 第二次on window close WM_CLOSE =>
     [个人] 允许 window close ->WM_DESTROY

[os] 系统群发WM_DESTROY
    [cef] CefBrowserHostImpl::WndProc =>
        [cef & os] 每个browser proc内的 WM_CLOSE 处理
        [cef]CefLifeSpanHandler::DestroyBrowser =>
            [实现cef接口]CefLifeSpanHandler::OnBeforeClose =>
                发起总关闭(个人使用XCefAppManage::QuitMessageLoop)->

退出消息循环
CefShutDown
```

这就一些对CEF CLOSE的个人理解的心得