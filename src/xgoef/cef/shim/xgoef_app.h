#ifndef XGOEF_APP_H_
#define XGOEF_APP_H_

#include "include/cef_app.h"

#include "xgoef_renderer.h"

namespace xgoef {

class App : public CefApp, public CefBrowserProcessHandler {
 public:
  App();

  CefRefPtr<CefBrowserProcessHandler> GetBrowserProcessHandler() override {
    return this;
  }
  CefRefPtr<CefRenderProcessHandler> GetRenderProcessHandler() override {
    return renderer_;
  }

  void OnContextInitialized() override;
  CefRefPtr<CefClient> GetDefaultClient() override;

 private:
  CefRefPtr<RendererHandler> renderer_;
  IMPLEMENT_REFCOUNTING(App);
};

}  // namespace xgoef

#endif  // XGOEF_APP_H_
