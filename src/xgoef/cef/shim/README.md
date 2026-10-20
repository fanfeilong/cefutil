# cef/shim

CEF C++ host pad（MessageRouter + Views），供 **LLGo**（`go:linkname`，无 cgo）链接。

| 文件 | 作用 |
| --- | --- |
| `xgoef_host.h` / `xgoef_host.cc` | lifecycle C ABI：`xgoef_cef_*` |
| `xgoef_ipc.*` | `cefQuery` → C 回调 → `xgoef_ipc_respond` |
| `xgoef_app.*` / `xgoef_handler.*` / `xgoef_renderer.*` | 浏览器 + 渲染进程 MessageRouter |
| `xgoef_helper_main.cc` | 原生 `xgoef-helper`（渲染侧注入 `window.cefQuery`） |

构建入口：仓库根目录 `scripts/build-xgoef-hello.sh`（编出 `libxgoef_cef.a` + `xgoef-helper`，再 `xgo build`）。

原则：shim **不写业务**；业务在 XGo/Go 的 `runtime` / `ipc` / `acl`。llcppg 扩大 capi 后，此目录应收窄为仅保留回调/Views/MessageRouter 难点。
