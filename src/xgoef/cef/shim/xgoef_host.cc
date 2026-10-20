/* Browser-process entry helpers — C++ CefApp with MessageRouter. */

#include "xgoef_host.h"

#include "include/base/cef_logging.h"
#include "include/cef_app.h"
#include "include/cef_command_line.h"

#include "xgoef_app.h"

#if defined(CEF_X11)
#include <X11/Xlib.h>
#endif

namespace {

CefRefPtr<xgoef::App> g_app;

#if defined(CEF_X11)
int XErrorHandlerImpl(Display* display, XErrorEvent* event) {
  (void)display;
  (void)event;
  return 0;
}

int XIOErrorHandlerImpl(Display* display) {
  (void)display;
  return 0;
}
#endif

}  // namespace

extern "C" void xgoef_cef_api_hash(void) {
  // C++ API uses CefInitialize which checks the API version via the wrapper.
  // Keep a no-op entry so existing Go call sites stay valid.
}

extern "C" int xgoef_cef_execute_process(int argc, char** argv) {
  CefMainArgs main_args(argc, argv);
  g_app = new xgoef::App();

  int exit_code = CefExecuteProcess(main_args, g_app.get(), nullptr);
  if (exit_code >= 0) {
    g_app = nullptr;
    return exit_code;
  }
  return -1;
}

extern "C" int xgoef_cef_initialize(int argc,
                                    char** argv,
                                    const char* root_cache_path,
                                    int no_sandbox,
                                    const char* subprocess_path) {
  CHECK(g_app);

#if defined(CEF_X11)
  XSetErrorHandler(XErrorHandlerImpl);
  XSetIOErrorHandler(XIOErrorHandlerImpl);
#endif

  CefMainArgs main_args(argc, argv);

  CefSettings settings;
  if (no_sandbox) {
    settings.no_sandbox = true;
  }
  if (root_cache_path && root_cache_path[0] != '\0') {
    CefString(&settings.root_cache_path).FromASCII(root_cache_path);
  }
  if (subprocess_path && subprocess_path[0] != '\0') {
    CefString(&settings.browser_subprocess_path).FromASCII(subprocess_path);
  }

  if (!CefInitialize(main_args, settings, g_app.get(), nullptr)) {
    return CefGetExitCode();
  }

  // initialize retains a reference; clear our holder.
  g_app = nullptr;
  return 0;
}

extern "C" void xgoef_cef_run_loop(void) {
  CefRunMessageLoop();
}

extern "C" void xgoef_cef_shutdown(void) {
  CefShutdown();
}
