# 跑通 hello

当前可运行路径：**XGo 应用** + **LLGo** 编译 + CEF C API（`cef/host_llgo.go` → `cef/shim`，**不用 cgo**）。  
一等语言是 **XGo**（`examples/hello/main.xgo`）；工具链：`scripts/install-xgo-toolchain.sh`。

## 工具链（XGo / LLGo / llcppg）

```bash
./scripts/install-xgo-toolchain.sh
# 需要：Go 1.27、LLVM 22、LLGO_ROOT=/opt/xgo-toolchain/llgo
xgo version && llgo version && llcppg -h | head -1
```

注意：不要用 apt 里的旧 `xgo` 配 Go 1.27（export data 版本不匹配）；脚本从源码装 XGo。

`src/xgoef/go.mod` 带 `// llgo 1.0`，因此 `xgo build` 会把编译器指到 LLGo。构建脚本通过 `scripts/xgo-llgo-gocmd.sh` 去掉 xgo 注入的、LLGo 尚不支持的 `-ldflags -X ...`。

## 前提

1. 已按仓库 skills 下载 CEF（`src/CEF_VERSION`），解压在 `third_party/cef/...`
2. 已编译过官方 `cefsimple_capi`（提供 `xgoef-helper`）：

```bash
CEF_ROOT=/path/to/cef_binary_*_linux64
cmake -S "$CEF_ROOT" -B "$CEF_ROOT/build" -G "Unix Makefiles" \
  -DCMAKE_BUILD_TYPE=Release -DUSE_SANDBOX=OFF \
  -DCMAKE_C_COMPILER=gcc -DCMAKE_CXX_COMPILER=g++
cmake --build "$CEF_ROOT/build" --target cefsimple_capi -j"$(nproc)"
```

## 构建与运行

```bash
./scripts/build-xgoef-hello.sh
cd src/xgoef/build/hello
./xgoef-hello
```

产物目录含：`xgoef-hello`（由 `main.xgo` 经 xgo+llgo 生成）、`xgoef-helper`（纯 C 子进程）、`libcef.so`、Resources、`web/`。

## 开发约定

| 层 | 语言 | 说明 |
| --- | --- | --- |
| 应用（`examples/*`） | **XGo**（`.xgo`） | 一等开发面；`xgo build` |
| 框架 API（`runtime`/`ipc`/`plugin`/`conf`） | Go 或 XGo | 可混编；新 API 优先 XGo |
| CEF 绑定（`cef/`） | Go + LLGo `linkname` | 贴近 C ABI；llcppg 生成物同左 |

## 关键行为

| 现象 | 原因 |
| --- | --- |
| 首次弹出 Chromium ToS | 新的 `root_cache`；点 Accept |
| 默认 `--disable-gpu --in-process-gpu` | 宿主不宜作 GPU 子进程；设 `XGOEF_ENABLE_GPU=1` 可关掉这些开关 |
| `browser_subprocess_path` → `xgoef-helper` | 渲染/网络等进程走原生 helper，不重入 LLGo 宿主 |
| 点 invoke greet 报 host bridge | JS↔IPC 桥尚未接到 CEF MessageRouter（下一阶段） |

## 与 llcppg 的关系

- 宿主已用 LLGo `go:linkname` 调 shim（无 cgo）。
- `cef/llcppg.cfg` 用于后续扩大到完整 `include/capi`；shim 应收窄为仅保留回调难点。
