#include "xgoef_renderer.h"

#include "include/wrapper/cef_helpers.h"

namespace xgoef {

void RendererHandler::OnWebKitInitialized() {
  CefMessageRouterConfig config;
  router_ = CefMessageRouterRendererSide::Create(config);
}

void RendererHandler::OnContextCreated(CefRefPtr<CefBrowser> browser,
                                       CefRefPtr<CefFrame> frame,
                                       CefRefPtr<CefV8Context> context) {
  if (router_) {
    router_->OnContextCreated(browser, frame, context);
  }
}

void RendererHandler::OnContextReleased(CefRefPtr<CefBrowser> browser,
                                        CefRefPtr<CefFrame> frame,
                                        CefRefPtr<CefV8Context> context) {
  if (router_) {
    router_->OnContextReleased(browser, frame, context);
  }
}

bool RendererHandler::OnProcessMessageReceived(
    CefRefPtr<CefBrowser> browser,
    CefRefPtr<CefFrame> frame,
    CefProcessId source_process,
    CefRefPtr<CefProcessMessage> message) {
  if (router_) {
    return router_->OnProcessMessageReceived(browser, frame, source_process,
                                             message);
  }
  return false;
}

}  // namespace xgoef
