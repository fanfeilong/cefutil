# cef/shim

可运行的 CEF C API 宿主垫片（复用 `CEF_ROOT/tests/cefsimple_capi` 源码）。

- `xgoef_host.h` / `xgoef_host.c`：供 Go cgo（及日后 LLGo）调用的 lifecycle API
- `CMakeLists.txt`：编出 `libxgoef_cef.a`，并拷贝 CEF 运行时文件

构建入口：仓库根目录 `scripts/build-xgoef-hello.sh`。

原则：shim **不写业务**；业务在 XGo/Go 的 `runtime` / `ipc`。llcppg 生成绑定后，此目录应收窄为仅保留回调/Views 难点。
