#ifndef XGOEF_RENDERER_H_
#define XGOEF_RENDERER_H_

#include "include/cef_render_process_handler.h"
#include "include/wrapper/cef_message_router.h"

namespace xgoef {

// Renderer-side MessageRouter — injects window.cefQuery.
class RendererHandler : public CefRenderProcessHandler {
 public:
  RendererHandler() = default;

  void OnWebKitInitialized() override;
  void OnContextCreated(CefRefPtr<CefBrowser> browser,
                        CefRefPtr<CefFrame> frame,
                        CefRefPtr<CefV8Context> context) override;
  void OnContextReleased(CefRefPtr<CefBrowser> browser,
                         CefRefPtr<CefFrame> frame,
                         CefRefPtr<CefV8Context> context) override;
  bool OnProcessMessageReceived(CefRefPtr<CefBrowser> browser,
                                CefRefPtr<CefFrame> frame,
                                CefProcessId source_process,
                                CefRefPtr<CefProcessMessage> message) override;

 private:
  CefRefPtr<CefMessageRouterRendererSide> router_;
  IMPLEMENT_REFCOUNTING(RendererHandler);
};

}  // namespace xgoef

#endif  // XGOEF_RENDERER_H_
