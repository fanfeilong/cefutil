// Internal C++ helpers for the MessageRouter ↔ C ABI bridge.
#ifndef XGOEF_IPC_H_
#define XGOEF_IPC_H_

#include "include/wrapper/cef_message_router.h"

#include "xgoef_host.h"

namespace xgoef {

// MessageRouter Handler that forwards queries to the C callback.
class IpcHandler : public CefMessageRouterBrowserSide::Handler {
 public:
  IpcHandler() = default;

  bool OnQuery(CefRefPtr<CefBrowser> browser,
               CefRefPtr<CefFrame> frame,
               int64_t query_id,
               const CefString& request,
               bool persistent,
               CefRefPtr<Callback> callback) override;

  void OnQueryCanceled(CefRefPtr<CefBrowser> browser,
                       CefRefPtr<CefFrame> frame,
                       int64_t query_id) override;
};

void IpcSetHandler(xgoef_ipc_query_fn fn);
void IpcRespond(int64_t query_id, int ok, const char* result_utf8);
void IpcCancel(int64_t query_id);

}  // namespace xgoef

#endif  // XGOEF_IPC_H_
