#include "app.h"

#include <string>

#include "handler.h"

#include "include/cef_browser.h"
#include "include/views/cef_browser_view.h"
#include "include/views/cef_window.h"
#include "include/wrapper/cef_helpers.h"

namespace {

class DemoWindowDelegate : public CefWindowDelegate {
 public:
  explicit DemoWindowDelegate(CefRefPtr<CefBrowserView> browser_view)
      : browser_view_(browser_view) {}

  void OnWindowCreated(CefRefPtr<CefWindow> window) override {
    window->SetTitle("cefutil src");
    window->AddChildView(browser_view_);
    window->Show();
  }

  void OnWindowDestroyed(CefRefPtr<CefWindow> window) override {
    browser_view_ = nullptr;
  }

  bool CanClose(CefRefPtr<CefWindow> window) override {
    CefRefPtr<CefBrowser> browser = browser_view_->GetBrowser();
    if (browser) {
      return browser->GetHost()->TryCloseBrowser();
    }
    return true;
  }

  CefSize GetPreferredSize(CefRefPtr<CefView> view) override {
    return CefSize(960, 640);
  }

 private:
  CefRefPtr<CefBrowserView> browser_view_;
  IMPLEMENT_REFCOUNTING(DemoWindowDelegate);
};

std::string PageUrl() {
  return "file://" + std::string(CEFDEMO_WEB_DIR) + "/index.html";
}

}  // namespace

DemoApp::DemoApp() = default;

void DemoApp::OnContextInitialized() {
  CEF_REQUIRE_UI_THREAD();

  CefRefPtr<DemoHandler> handler(new DemoHandler);
  CefBrowserSettings browser_settings;
  CefRefPtr<CefBrowserView> browser_view = CefBrowserView::CreateBrowserView(
      handler, PageUrl(), browser_settings, nullptr, nullptr, nullptr);
  CefWindow::CreateTopLevelWindow(new DemoWindowDelegate(browser_view));
}
