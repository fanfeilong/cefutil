/* xgoef_host.h — C ABI surface for the Go/LLGo Host implementation.
 *
 * Lifecycle mirrors cefsimple (C++) and docs/05-lifecycle.md.
 * Startup URL is taken from argv `--url=` (same as cefsimple).
 * IPC uses CefMessageRouter; see xgoef_ipc_*.
 */
#ifndef XGOEF_HOST_H_
#define XGOEF_HOST_H_

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/* Call before any other CEF function (Linux/Windows). */
void xgoef_cef_api_hash(void);

/* Returns >=0 when this process is a CEF helper and should exit with that code.
 * Returns -1 when this is the browser process. */
int xgoef_cef_execute_process(int argc, char** argv);

/* Initialize browser process. root_cache_path / subprocess_path may be NULL.
 * no_sandbox non-zero disables the sandbox (learning builds).
 * When subprocess_path is set, CEF helpers use that native binary instead of
 * re-entering the Go host executable (required for reliable GPU/renderer).
 * Returns 0 on success, non-zero on failure (see cef_get_exit_code). */
int xgoef_cef_initialize(int argc,
                         char** argv,
                         const char* root_cache_path,
                         int no_sandbox,
                         const char* subprocess_path);

void xgoef_cef_run_loop(void);
void xgoef_cef_shutdown(void);

/* ---- IPC (CefMessageRouter bridge) ---- */

/* Called on the CEF UI thread when the page issues cefQuery with a JSON
 * request body. The host must eventually call xgoef_ipc_respond. */
typedef void (*xgoef_ipc_query_fn)(int64_t query_id, const char* request_utf8);

/* Register the Go/LLGo invoke handler. Pass NULL to clear. */
void xgoef_ipc_set_handler(xgoef_ipc_query_fn fn);

/* Complete a query. ok!=0 → Success(result_utf8); else Failure(-2, result_utf8).
 * Safe to call from any thread (posted to UI if needed). */
void xgoef_ipc_respond(int64_t query_id, int ok, const char* result_utf8);

#ifdef __cplusplus
}
#endif

#endif /* XGOEF_HOST_H_ */
