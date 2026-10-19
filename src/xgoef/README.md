# xgoef

基于 **CEF** 的 **XGo** 客户端框架实验：用 [LLGo](https://github.com/xgo-dev/llgo) + [llcppg](https://github.com/goplus/llcppg) 对接 CEF **C API**，产品形态对标 [Tauri 2](https://v2.tauri.app)（Web UI + 原生后端 + IPC/插件）。

本目录是 `cefutil` 里与 `src/cefdemo` 并列的实验树，不是已可发布的发行版。

## 先读

1. [RESEARCH.md](RESEARCH.md) — 相关项目调研与选型
2. [DESIGN.md](DESIGN.md) — 分层与进程模型
3. 仓库文档 [docs/03-binary-boundary.md](../../docs/03-binary-boundary.md)、[docs/05-lifecycle.md](../../docs/05-lifecycle.md)、[docs/07-javascript.md](../../docs/07-javascript.md)

## 目录

```text
xgoef/
├── RESEARCH.md DESIGN.md README.md
├── go.mod
├── xgoef.conf.schema.json
├── cef/           # llcppg 配置 + 绑定说明（生成代码后续加入）
├── runtime/       # App / Window / 生命周期
├── ipc/           # Command + Event
├── plugin/        # 插件接口
├── conf/          # 配置
├── api/js/        # 前端桥
├── cmd/xgoef/     # CLI 占位
└── examples/hello # 最小应用形状
```

## 公式

```text
xgoef ≈ Tauri2 的应用模型 × CEF 的 Chromium × (XGo + LLGo + llcppg)
```

| Tauri 2 | xgoef |
| --- | --- |
| 系统 WebView | CEF（自带 Chromium） |
| Rust 核心 | XGo / Go（LLGo 编译） |
| `invoke` / events / plugins | `ipc` + `plugin`（同形状） |
| `tauri.conf.json` | `xgoef.conf.json` |

## 当前可跑通路径（Go + cgo + CEF C API）

```bash
# 需已下载 CEF，并编过 cefsimple_capi（作 xgoef-helper）
./scripts/build-xgoef-hello.sh
cd src/xgoef/build/hello && ./xgoef-hello
```

细节见 [docs/run.md](docs/run.md)。`cef/shim` 复用官方 `cefsimple_capi` 生命周期；子进程用原生 `xgoef-helper`，默认关 GPU（`XGOEF_ENABLE_GPU=1` 可改）。

## 下一步（llcppg / LLGo）

1. 安装 [XGo](https://xgo.dev) 与 [LLGo](https://github.com/xgo-dev/llgo)（llcppg 跟 LLGo main）
2. `cd cef && export CEF_ROOT=... && llcppg llcppg.cfg`
3. 用生成绑定逐步替换 shim；IPC 接 CEF MessageRouter

无 CEF 时 `go test ./...` 仍可通过（Host 接口隔离）。

## 与 cefdemo

`cefdemo` 继续作为 C++ / Views / Chrome style 的对照实现。xgoef 验证同一套生命周期能否用 C API + LLGo 表达，并在其上长出 Tauri 式应用 API。
