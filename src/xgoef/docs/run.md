# 跑通 hello

当前可运行路径：**Go + cgo + CEF C API**（`cef/shim` 复用 `cefsimple_capi`）。  
llcppg / LLGo 直连是下一阶段；工具链可用 `scripts/install-xgo-toolchain.sh` 安装。

## 工具链（XGo / LLGo / llcppg）

```bash
./scripts/install-xgo-toolchain.sh
# 需要：Go 1.27、LLVM 22、LLGO_ROOT=/opt/xgo-toolchain/llgo
xgo version && llgo version && llcppg -h | head -1
```

注意：不要用 apt 里的旧 `xgo` 配 Go 1.27（export data 版本不匹配）；脚本从源码装 XGo。

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

产物目录含：`xgoef-hello`（Go 宿主）、`xgoef-helper`（纯 C 子进程）、`libcef.so`、Resources、`web/`。

## 关键行为

| 现象 | 原因 |
| --- | --- |
| 首次弹出 Chromium ToS | 新的 `root_cache`；点 Accept |
| 默认 `--disable-gpu --in-process-gpu` | Go 不宜作为 GPU 子进程；设 `XGOEF_ENABLE_GPU=1` 可关掉这些开关 |
| `browser_subprocess_path` → `xgoef-helper` | 渲染/网络等进程走原生 helper，不重入 Go |
| 点 invoke greet 报 host bridge | JS↔IPC 桥尚未接到 CEF MessageRouter（下一阶段） |

## 与 llcppg 的关系

- `cef/llcppg.cfg` 仍是生成绑定的目标配置。
- `cef/shim` 是在 llcppg/LLGo 工具链未齐前的 **C ABI 垫片**，供 cgo（及日后 LLGo）调用同一套 lifecycle。
