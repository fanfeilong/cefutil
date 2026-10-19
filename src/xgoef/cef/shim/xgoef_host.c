/* Browser-process entry helpers adapted from cefsimple_capi/cefsimple_linux.c */

#include "xgoef_host.h"

#include <stdlib.h>
#include <string.h>

#if defined(CEF_X11)
#include <X11/Xlib.h>
#endif

#include "include/capi/cef_app_capi.h"
#include "include/cef_api_hash.h"
#include "tests/cefsimple_capi/simple_app.h"
#include "tests/cefsimple_capi/simple_utils.h"

static simple_app_t* g_app = NULL;

#if defined(CEF_X11)
static int XErrorHandlerImpl(Display* display, XErrorEvent* event) {
  (void)display;
  (void)event;
  return 0;
}

static int XIOErrorHandlerImpl(Display* display) {
  (void)display;
  return 0;
}
#endif

void xgoef_cef_api_hash(void) {
  cef_api_hash(CEF_API_VERSION, 0);
}

int xgoef_cef_execute_process(int argc, char** argv) {
  cef_main_args_t main_args = {};
  main_args.argc = argc;
  main_args.argv = argv;

  g_app = simple_app_create();
  CHECK(g_app);

  /* Both execute_process (helper path) and initialize take one reference.
   * Hold an extra so the browser path still has a ref for initialize. */
  g_app->app.base.add_ref(&g_app->app.base);

  int exit_code = cef_execute_process(&main_args, &g_app->app, NULL);
  if (exit_code >= 0) {
    /* Helper process: execute_process consumed one ref; drop the extra. */
    g_app->app.base.release(&g_app->app.base);
    g_app = NULL;
    return exit_code;
  }

  /* Browser process: execute_process did not consume a ref. */
  return -1;
}

int xgoef_cef_initialize(int argc,
                         char** argv,
                         const char* root_cache_path,
                         int no_sandbox,
                         const char* subprocess_path) {
  CHECK(g_app);

#if defined(CEF_X11)
  XSetErrorHandler(XErrorHandlerImpl);
  XSetIOErrorHandler(XIOErrorHandlerImpl);
#endif

  cef_main_args_t main_args = {};
  main_args.argc = argc;
  main_args.argv = argv;

  cef_settings_t settings = {};
  settings.size = sizeof(cef_settings_t);
  if (no_sandbox) {
    settings.no_sandbox = 1;
  }
  if (root_cache_path && root_cache_path[0] != '\0') {
    cef_string_from_ascii(root_cache_path, (size_t)strlen(root_cache_path),
                          &settings.root_cache_path);
  }
  if (subprocess_path && subprocess_path[0] != '\0') {
    cef_string_from_ascii(subprocess_path, (size_t)strlen(subprocess_path),
                          &settings.browser_subprocess_path);
  }

  if (!cef_initialize(&main_args, &settings, &g_app->app, NULL)) {
    return cef_get_exit_code();
  }

  /* initialize consumed one ref; the remaining ref lives until shutdown. */
  g_app = NULL;
  return 0;
}

void xgoef_cef_run_loop(void) { cef_run_message_loop(); }

void xgoef_cef_shutdown(void) { cef_shutdown(); }
