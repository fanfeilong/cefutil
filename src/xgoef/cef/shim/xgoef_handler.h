#ifndef XGOEF_HANDLER_H_
#define XGOEF_HANDLER_H_

#include <list>

#include "include/cef_client.h"
#include "include/cef_display_handler.h"
#include "include/cef_life_span_handler.h"
#include "include/cef_load_handler.h"
#include "include/cef_request_handler.h"
#include "include/wrapper/cef_message_router.h"

#include "xgoef_ipc.h"

namespace xgoef {

class ClientHandler : public CefClient,
                      public CefDisplayHandler,
                      public CefLifeSpanHandler,
                      public CefLoadHandler,
                      public CefRequestHandler {
 public:
  explicit ClientHandler(bool is_alloy_style);
  ~ClientHandler() override;

  static ClientHandler* GetInstance();

  CefRefPtr<CefDisplayHandler> GetDisplayHandler() override { return this; }
  CefRefPtr<CefLifeSpanHandler> GetLifeSpanHandler() override { return this; }
  CefRefPtr<CefLoadHandler> GetLoadHandler() override { return this; }
  CefRefPtr<CefRequestHandler> GetRequestHandler() override { return this; }

  bool OnProcessMessageReceived(CefRefPtr<CefBrowser> browser,
                                CefRefPtr<CefFrame> frame,
                                CefProcessId source_process,
                                CefRefPtr<CefProcessMessage> message) override;

  void OnTitleChange(CefRefPtr<CefBrowser> browser,
                     const CefString& title) override;

  void OnAfterCreated(CefRefPtr<CefBrowser> browser) override;
  bool DoClose(CefRefPtr<CefBrowser> browser) override;
  void OnBeforeClose(CefRefPtr<CefBrowser> browser) override;

  void OnLoadError(CefRefPtr<CefBrowser> browser,
                   CefRefPtr<CefFrame> frame,
                   ErrorCode errorCode,
                   const CefString& errorText,
                   const CefString& failedUrl) override;

  bool OnBeforeBrowse(CefRefPtr<CefBrowser> browser,
                      CefRefPtr<CefFrame> frame,
                      CefRefPtr<CefRequest> request,
                      bool user_gesture,
                      bool is_redirect) override;
  void OnRenderProcessTerminated(CefRefPtr<CefBrowser> browser,
                                 TerminationStatus status,
                                 int error_code,
                                 const CefString& error_string) override;

  void CloseAllBrowsers(bool force_close);

 private:
  void PlatformTitleChange(CefRefPtr<CefBrowser> browser,
                           const CefString& title);

  const bool is_alloy_style_;
  typedef std::list<CefRefPtr<CefBrowser>> BrowserList;
  BrowserList browser_list_;
  bool is_closing_ = false;

  CefRefPtr<CefMessageRouterBrowserSide> message_router_;
  IpcHandler ipc_handler_;

  IMPLEMENT_REFCOUNTING(ClientHandler);
};

}  // namespace xgoef

#endif  // XGOEF_HANDLER_H_
