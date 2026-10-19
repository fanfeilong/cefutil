/* xgoef_host.h — C ABI surface for the Go/LLGo Host implementation.
 *
 * Lifecycle mirrors tests/cefsimple_capi (and docs/05-lifecycle.md).
 * Startup URL is taken from argv `--url=` (same as cefsimple_capi).
 */
#ifndef XGOEF_HOST_H_
#define XGOEF_HOST_H_

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

#ifdef __cplusplus
}
#endif

#endif /* XGOEF_HOST_H_ */
