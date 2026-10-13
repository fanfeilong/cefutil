#include "app.h"

#include <filesystem>

#include <X11/Xlib.h>

#include "include/cef_app.h"
#include "include/cef_command_line.h"

namespace {

int XErrorHandlerImpl(Display* display, XErrorEvent* event) {
  return 0;
}

int XIOErrorHandlerImpl(Display* display) {
  return 0;
}

}  // namespace

int main(int argc, char* argv[]) {
  CefMainArgs main_args(argc, argv);

  const int exit_code = CefExecuteProcess(main_args, nullptr, nullptr);
  if (exit_code >= 0) {
    return exit_code;
  }

  XSetErrorHandler(XErrorHandlerImpl);
  XSetIOErrorHandler(XIOErrorHandlerImpl);

  CefSettings settings;
#if !defined(CEF_USE_SANDBOX)
  settings.no_sandbox = true;
#endif
  std::filesystem::create_directories(CEFDEMO_ROOT_CACHE);
  CefString(&settings.root_cache_path).FromASCII(CEFDEMO_ROOT_CACHE);

  CefRefPtr<DemoApp> app(new DemoApp);
  if (!CefInitialize(main_args, settings, app.get(), nullptr)) {
    return CefGetExitCode();
  }

  CefRunMessageLoop();
  CefShutdown();
  return 0;
}
