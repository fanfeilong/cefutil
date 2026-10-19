# cef — CEF C API 绑定层

## 策略

- **只绑 C API**（`CEF_ROOT/include/capi/`），不绑 `libcef_dll_wrapper`。
- 用 **llcppg** 生成 LLGo 包；手写层只做 string / refcount / 生命周期糖衣。
- 参考 [bnema/purego-cef](https://github.com/bnema/purego-cef) 的「生成 + 小手写层」分层，实现换成 llcppg/LLGo。

## 生成（需本机已装 llgo、llcppg，并解压 CEF）

```bash
export CEF_ROOT=/absolute/path/to/cef_binary_*_linux64
cd "$(dirname "$0")"
# 按需编辑 llcppg.cfg 的 include 列表后：
llcppg llcppg.cfg
```

生成目录名由 `llcppg.cfg` 的 `name` 决定（默认 `capi`）。可将 `capi/` 加入版本库，或在 CI 中生成。

## 最小头文件集合（一期）

优先纳入启动闭环，避免一次吃下全部 capi：

- `cef_base_capi.h`
- `cef_app_capi.h`
- `cef_client_capi.h`
- `cef_browser_capi.h`
- `cef_browser_process_handler_capi.h`
- `cef_life_span_handler_capi.h`
- `cef_frame_capi.h`
- `cef_string_capi.h` / `cef_string_types.h`（若拆分）
- `cef_api_hash.h`（若在 include 根而非 capi）

具体文件名以当前 `CEF_ROOT/include` 为准；`llcppg.cfg` 里用 `mix: true` 并显式列出 `include`。

## 本目录

- `llcppg.cfg` — 生成配置模板（`cflags`/`libs` 依赖 `CEF_ROOT`）
- `host.go` / `host_llgo.go` / `host_stub.go` — Host 接口；LLGo `go:linkname` 实现（`!llgo` stub 供 gc 测试）
- `shim/` — 当前可跑通的 C 宿主（复用 `cefsimple_capi`）；llcppg 就位后逐步变薄
