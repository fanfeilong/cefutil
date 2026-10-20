// Window label registry + host→renderer event delivery (ExecuteJavaScript).
#ifndef XGOEF_WINDOW_H_
#define XGOEF_WINDOW_H_

#include <string>

#include "include/cef_browser.h"

#include "xgoef_host.h"

namespace xgoef {

// Next browser created will be bound to this label (cleared after take/bind).
void WindowSetPendingLabel(const std::string& label);
std::string WindowTakePendingLabel();

// Bind / unbind browsers for EmitTo / Close / Focus.
void WindowBind(const std::string& label, CefRefPtr<CefBrowser> browser);
void WindowUnbind(CefRefPtr<CefBrowser> browser);
std::string WindowLabelFor(CefRefPtr<CefBrowser> browser);

// Create another Views window (must run on UI thread or is posted).
// Returns 0 on success.
int WindowCreate(const char* label,
                 const char* title,
                 const char* url,
                 int width,
                 int height);

int WindowClose(const char* label);
int WindowFocus(const char* label);

// Deliver JSON event envelope to one / all window main frames.
int WindowEmitTo(const char* label, const char* json_utf8);
int WindowEmitAll(const char* json_utf8);

}  // namespace xgoef

#endif  // XGOEF_WINDOW_H_
