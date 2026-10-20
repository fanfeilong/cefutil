#include "xgoef_ipc.h"

#include <map>
#include <mutex>
#include <string>

#include "include/base/cef_callback.h"
#include "include/wrapper/cef_closure_task.h"
#include "include/wrapper/cef_helpers.h"

#include "xgoef_window.h"

namespace xgoef {
namespace {

std::mutex g_mu;
xgoef_ipc_query_fn g_handler = nullptr;
std::map<int64_t, CefRefPtr<CefMessageRouterBrowserSide::Callback>> g_pending;

void RespondOnUI(int64_t query_id, int ok, std::string result) {
  CEF_REQUIRE_UI_THREAD();
  CefRefPtr<CefMessageRouterBrowserSide::Callback> cb;
  {
    std::lock_guard<std::mutex> lock(g_mu);
    auto it = g_pending.find(query_id);
    if (it == g_pending.end()) {
      return;
    }
    cb = it->second;
    g_pending.erase(it);
  }
  if (!cb) {
    return;
  }
  if (ok) {
    cb->Success(result);
  } else {
    cb->Failure(-2, result);
  }
}

}  // namespace

bool IpcHandler::OnQuery(CefRefPtr<CefBrowser> browser,
                         CefRefPtr<CefFrame> frame,
                         int64_t query_id,
                         const CefString& request,
                         bool persistent,
                         CefRefPtr<Callback> callback) {
  CEF_REQUIRE_UI_THREAD();
  (void)frame;
  (void)persistent;

  xgoef_ipc_query_fn handler = nullptr;
  {
    std::lock_guard<std::mutex> lock(g_mu);
    handler = g_handler;
    g_pending[query_id] = callback;
  }

  if (!handler) {
    {
      std::lock_guard<std::mutex> lock(g_mu);
      g_pending.erase(query_id);
    }
    callback->Failure(-1, "xgoef host bridge not connected");
    return true;
  }

  const std::string label = WindowLabelFor(browser);
  const std::string req = request.ToString();
  handler(query_id, label.c_str(), req.c_str());
  return true;
}

void IpcHandler::OnQueryCanceled(CefRefPtr<CefBrowser> browser,
                                 CefRefPtr<CefFrame> frame,
                                 int64_t query_id) {
  CEF_REQUIRE_UI_THREAD();
  (void)browser;
  (void)frame;
  IpcCancel(query_id);
}

void IpcSetHandler(xgoef_ipc_query_fn fn) {
  std::lock_guard<std::mutex> lock(g_mu);
  g_handler = fn;
}

void IpcRespond(int64_t query_id, int ok, const char* result_utf8) {
  std::string result = result_utf8 ? result_utf8 : "";
  if (CefCurrentlyOn(TID_UI)) {
    RespondOnUI(query_id, ok, std::move(result));
    return;
  }
  CefPostTask(TID_UI, base::BindOnce(&RespondOnUI, query_id, ok, result));
}

void IpcCancel(int64_t query_id) {
  std::lock_guard<std::mutex> lock(g_mu);
  g_pending.erase(query_id);
}

}  // namespace xgoef

extern "C" void xgoef_ipc_set_handler(xgoef_ipc_query_fn fn) {
  xgoef::IpcSetHandler(fn);
}

extern "C" void xgoef_ipc_respond(int64_t query_id,
                                  int ok,
                                  const char* result_utf8) {
  xgoef::IpcRespond(query_id, ok, result_utf8);
}
