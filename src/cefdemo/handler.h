#ifndef CEFUTIL_SRC_HANDLER_H_
#define CEFUTIL_SRC_HANDLER_H_

#include "include/cef_client.h"

class DemoHandler : public CefClient, public CefLifeSpanHandler, public CefLoadHandler {
 public:
  DemoHandler();

  CefRefPtr<CefLifeSpanHandler> GetLifeSpanHandler() override { return this; }
  CefRefPtr<CefLoadHandler> GetLoadHandler() override { return this; }

  void OnAfterCreated(CefRefPtr<CefBrowser> browser) override;
  bool DoClose(CefRefPtr<CefBrowser> browser) override;
  void OnBeforeClose(CefRefPtr<CefBrowser> browser) override;
  void OnLoadEnd(CefRefPtr<CefBrowser> browser,
                 CefRefPtr<CefFrame> frame,
                 int httpStatusCode) override;

 private:
  int browser_count_ = 0;
  bool closing_ = false;

  IMPLEMENT_REFCOUNTING(DemoHandler);
};

#endif
