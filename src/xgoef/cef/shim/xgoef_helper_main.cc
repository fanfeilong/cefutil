// Native CEF helper (renderer/GPU/etc). Must include MessageRouter renderer
// side so window.cefQuery is available to xgoef.js.

#include "include/cef_app.h"

#include "xgoef_app.h"

#if defined(CEF_X11)
#include <X11/Xlib.h>
#endif

#if defined(CEF_X11)
namespace {
int XErrorHandlerImpl(Display*, XErrorEvent*) {
  return 0;
}
int XIOErrorHandlerImpl(Display*) {
  return 0;
}
}  // namespace
#endif

NO_STACK_PROTECTOR
int main(int argc, char* argv[]) {
#if defined(CEF_X11)
  XSetErrorHandler(XErrorHandlerImpl);
  XSetIOErrorHandler(XIOErrorHandlerImpl);
#endif

  CefMainArgs main_args(argc, argv);
  CefRefPtr<xgoef::App> app(new xgoef::App());
  int code = CefExecuteProcess(main_args, app.get(), nullptr);
  return code >= 0 ? code : 0;
}
