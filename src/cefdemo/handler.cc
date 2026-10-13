#include "handler.h"

#include "include/cef_app.h"
#include "include/wrapper/cef_helpers.h"

DemoHandler::DemoHandler() = default;

void DemoHandler::OnAfterCreated(CefRefPtr<CefBrowser> browser) {
  CEF_REQUIRE_UI_THREAD();
  browser_count_++;
}

bool DemoHandler::DoClose(CefRefPtr<CefBrowser> browser) {
  CEF_REQUIRE_UI_THREAD();
  if (browser_count_ == 1) {
    closing_ = true;
  }
  return false;
}

void DemoHandler::OnBeforeClose(CefRefPtr<CefBrowser> browser) {
  CEF_REQUIRE_UI_THREAD();
  if (--browser_count_ == 0) {
    CefQuitMessageLoop();
  }
}

void DemoHandler::OnLoadEnd(CefRefPtr<CefBrowser> browser,
                            CefRefPtr<CefFrame> frame,
                            int httpStatusCode) {
  CEF_REQUIRE_UI_THREAD();
  if (frame->IsMain() && !closing_) {
    frame->ExecuteJavaScript(
        "document.getElementById('status').textContent = 'page loaded';",
        frame->GetURL(), 0);
  }
}
