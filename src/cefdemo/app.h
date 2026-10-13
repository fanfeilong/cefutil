#ifndef CEFUTIL_SRC_APP_H_
#define CEFUTIL_SRC_APP_H_

#include "include/cef_app.h"

class DemoApp : public CefApp, public CefBrowserProcessHandler {
 public:
  DemoApp();

  CefRefPtr<CefBrowserProcessHandler> GetBrowserProcessHandler() override {
    return this;
  }

  void OnContextInitialized() override;

 private:
  IMPLEMENT_REFCOUNTING(DemoApp);
};

#endif
