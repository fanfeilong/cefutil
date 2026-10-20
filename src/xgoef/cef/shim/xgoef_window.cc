#include "xgoef_window.h"

#include <map>
#include <mutex>
#include <string>

#include "include/base/cef_callback.h"
#include "include/cef_command_line.h"
#include "include/views/cef_browser_view.h"
#include "include/views/cef_window.h"
#include "include/wrapper/cef_closure_task.h"
#include "include/wrapper/cef_helpers.h"

#include "xgoef_handler.h"

namespace xgoef {
namespace {

std::mutex g_mu;
std::string g_pending_label = "main";
std::map<std::string, CefRefPtr<CefBrowser>> g_by_label;
std::map<int, std::string> g_label_by_id;

class SizedWindowDelegate : public CefWindowDelegate {
 public:
  SizedWindowDelegate(CefRefPtr<CefBrowserView> browser_view,
                      cef_runtime_style_t runtime_style,
                      int width,
                      int height,
                      const std::string& title)
      : browser_view_(browser_view),
        runtime_style_(runtime_style),
        width_(width > 0 ? width : 960),
        height_(height > 0 ? height : 640),
        title_(title) {}

  void OnWindowCreated(CefRefPtr<CefWindow> window) override {
    window->AddChildView(browser_view_);
    if (!title_.empty()) {
      window->SetTitle(title_);
    }
    window->Show();
  }

  void OnWindowDestroyed(CefRefPtr<CefWindow> window) override {
    (void)window;
    browser_view_ = nullptr;
  }

  bool CanClose(CefRefPtr<CefWindow> window) override {
    (void)window;
    CefRefPtr<CefBrowser> browser = browser_view_->GetBrowser();
    if (browser) {
      return browser->GetHost()->TryCloseBrowser();
    }
    return true;
  }

  CefSize GetPreferredSize(CefRefPtr<CefView> view) override {
    (void)view;
    return CefSize(width_, height_);
  }

  cef_runtime_style_t GetWindowRuntimeStyle() override {
    return runtime_style_;
  }

 private:
  CefRefPtr<CefBrowserView> browser_view_;
  const cef_runtime_style_t runtime_style_;
  const int width_;
  const int height_;
  const std::string title_;
  IMPLEMENT_REFCOUNTING(SizedWindowDelegate);
};

class StyleBrowserViewDelegate : public CefBrowserViewDelegate {
 public:
  explicit StyleBrowserViewDelegate(cef_runtime_style_t runtime_style)
      : runtime_style_(runtime_style) {}

  bool OnPopupBrowserViewCreated(CefRefPtr<CefBrowserView> browser_view,
                                 CefRefPtr<CefBrowserView> popup_browser_view,
                                 bool is_devtools) override {
    (void)browser_view;
    (void)is_devtools;
    CefWindow::CreateTopLevelWindow(new SizedWindowDelegate(
        popup_browser_view, runtime_style_, 960, 640, ""));
    return true;
  }

  cef_runtime_style_t GetBrowserRuntimeStyle() override {
    return runtime_style_;
  }

 private:
  const cef_runtime_style_t runtime_style_;
  IMPLEMENT_REFCOUNTING(StyleBrowserViewDelegate);
};

void EmitOnUI(CefRefPtr<CefBrowser> browser, std::string json) {
  CEF_REQUIRE_UI_THREAD();
  if (!browser) {
    return;
  }
  CefRefPtr<CefFrame> frame = browser->GetMainFrame();
  if (!frame) {
    return;
  }
  // json is a JSON object; embed directly into the dispatch call.
  const std::string code =
      "window.__xgoefDispatch&&window.__xgoefDispatch(" + json + ");";
  frame->ExecuteJavaScript(code, frame->GetURL(), 0);
}

CefRefPtr<CefBrowser> FindLocked(const std::string& label) {
  auto it = g_by_label.find(label);
  if (it == g_by_label.end()) {
    return nullptr;
  }
  return it->second;
}

}  // namespace

void WindowSetPendingLabel(const std::string& label) {
  std::lock_guard<std::mutex> lock(g_mu);
  g_pending_label = label.empty() ? "main" : label;
}

std::string WindowTakePendingLabel() {
  std::lock_guard<std::mutex> lock(g_mu);
  std::string out = g_pending_label.empty() ? "main" : g_pending_label;
  g_pending_label = "main";
  return out;
}

void WindowBind(const std::string& label, CefRefPtr<CefBrowser> browser) {
  if (!browser) {
    return;
  }
  std::string use = label.empty() ? "main" : label;
  std::lock_guard<std::mutex> lock(g_mu);
  // If label already taken, keep unique by appending browser id.
  if (g_by_label.count(use) && !g_by_label[use]->IsSame(browser)) {
    use = use + "-" + std::to_string(browser->GetIdentifier());
  }
  g_by_label[use] = browser;
  g_label_by_id[browser->GetIdentifier()] = use;
}

void WindowUnbind(CefRefPtr<CefBrowser> browser) {
  if (!browser) {
    return;
  }
  std::lock_guard<std::mutex> lock(g_mu);
  auto it = g_label_by_id.find(browser->GetIdentifier());
  if (it == g_label_by_id.end()) {
    return;
  }
  g_by_label.erase(it->second);
  g_label_by_id.erase(it);
}

std::string WindowLabelFor(CefRefPtr<CefBrowser> browser) {
  if (!browser) {
    return "main";
  }
  std::lock_guard<std::mutex> lock(g_mu);
  auto it = g_label_by_id.find(browser->GetIdentifier());
  if (it == g_label_by_id.end()) {
    return "main";
  }
  return it->second;
}

void DoWindowCreate(std::string lab,
                    std::string tit,
                    std::string u,
                    int width,
                    int height) {
  CEF_REQUIRE_UI_THREAD();
  ClientHandler* raw = ClientHandler::GetInstance();
  if (!raw) {
    return;
  }
  CefRefPtr<ClientHandler> handler = raw;
  WindowSetPendingLabel(lab);

  CefRefPtr<CefCommandLine> command_line =
      CefCommandLine::GetGlobalCommandLine();
  cef_runtime_style_t runtime_style = CEF_RUNTIME_STYLE_DEFAULT;
  if (command_line && command_line->HasSwitch("use-alloy-style")) {
    runtime_style = CEF_RUNTIME_STYLE_ALLOY;
  }

  CefBrowserSettings browser_settings;
  CefRefPtr<CefBrowserView> browser_view = CefBrowserView::CreateBrowserView(
      handler, u, browser_settings, nullptr, nullptr,
      new StyleBrowserViewDelegate(runtime_style));
  CefWindow::CreateTopLevelWindow(
      new SizedWindowDelegate(browser_view, runtime_style, width, height, tit));
}

void DoWindowClose(std::string lab) {
  CEF_REQUIRE_UI_THREAD();
  CefRefPtr<CefBrowser> browser;
  {
    std::lock_guard<std::mutex> lock(g_mu);
    browser = FindLocked(lab);
  }
  if (browser) {
    browser->GetHost()->CloseBrowser(false);
  }
}

void DoWindowFocus(std::string lab) {
  CEF_REQUIRE_UI_THREAD();
  CefRefPtr<CefBrowser> browser;
  {
    std::lock_guard<std::mutex> lock(g_mu);
    browser = FindLocked(lab);
  }
  if (!browser) {
    return;
  }
  if (auto view = CefBrowserView::GetForBrowser(browser)) {
    if (auto window = view->GetWindow()) {
      window->Show();
      window->Activate();
    }
  }
  browser->GetHost()->SetFocus(true);
}

void DoEmitTo(std::string lab, std::string json) {
  CEF_REQUIRE_UI_THREAD();
  CefRefPtr<CefBrowser> browser;
  {
    std::lock_guard<std::mutex> lock(g_mu);
    browser = FindLocked(lab);
  }
  EmitOnUI(browser, json);
}

void DoEmitAll(std::string json) {
  CEF_REQUIRE_UI_THREAD();
  std::map<std::string, CefRefPtr<CefBrowser>> copy;
  {
    std::lock_guard<std::mutex> lock(g_mu);
    copy = g_by_label;
  }
  for (const auto& kv : copy) {
    EmitOnUI(kv.second, json);
  }
}

int WindowCreate(const char* label,
                 const char* title,
                 const char* url,
                 int width,
                 int height) {
  std::string lab = label ? label : "main";
  std::string tit = title ? title : "";
  std::string u = url && url[0] ? url : "about:blank";
  if (CefCurrentlyOn(TID_UI)) {
    DoWindowCreate(lab, tit, u, width, height);
  } else {
    CefPostTask(TID_UI, base::BindOnce(&DoWindowCreate, lab, tit, u, width,
                                       height));
  }
  return 0;
}

int WindowClose(const char* label) {
  std::string lab = label && label[0] ? label : "main";
  if (CefCurrentlyOn(TID_UI)) {
    DoWindowClose(lab);
  } else {
    CefPostTask(TID_UI, base::BindOnce(&DoWindowClose, lab));
  }
  return 0;
}

int WindowFocus(const char* label) {
  std::string lab = label && label[0] ? label : "main";
  if (CefCurrentlyOn(TID_UI)) {
    DoWindowFocus(lab);
  } else {
    CefPostTask(TID_UI, base::BindOnce(&DoWindowFocus, lab));
  }
  return 0;
}

int WindowEmitTo(const char* label, const char* json_utf8) {
  if (!json_utf8) {
    return -1;
  }
  std::string lab = label && label[0] ? label : "main";
  std::string json = json_utf8;
  if (CefCurrentlyOn(TID_UI)) {
    DoEmitTo(lab, json);
  } else {
    CefPostTask(TID_UI, base::BindOnce(&DoEmitTo, lab, json));
  }
  return 0;
}

int WindowEmitAll(const char* json_utf8) {
  if (!json_utf8) {
    return -1;
  }
  std::string json = json_utf8;
  if (CefCurrentlyOn(TID_UI)) {
    DoEmitAll(json);
  } else {
    CefPostTask(TID_UI, base::BindOnce(&DoEmitAll, json));
  }
  return 0;
}

}  // namespace xgoef

extern "C" void xgoef_window_set_pending_label(const char* label) {
  xgoef::WindowSetPendingLabel(label ? label : "main");
}

extern "C" int xgoef_window_create(const char* label,
                                   const char* title,
                                   const char* url,
                                   int width,
                                   int height) {
  return xgoef::WindowCreate(label, title, url, width, height);
}

extern "C" int xgoef_window_close(const char* label) {
  return xgoef::WindowClose(label);
}

extern "C" int xgoef_window_focus(const char* label) {
  return xgoef::WindowFocus(label);
}

extern "C" int xgoef_ipc_emit_to(const char* label, const char* json_utf8) {
  return xgoef::WindowEmitTo(label, json_utf8);
}

extern "C" int xgoef_ipc_emit_all(const char* json_utf8) {
  return xgoef::WindowEmitAll(json_utf8);
}
