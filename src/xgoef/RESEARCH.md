# 相关项目调研

调研日期：2026-10-19。目标：用 **XGo + LLGo + llcppg** 基于 **CEF** 做一套类似 **Tauri 2** 的客户端框架（本目录 `xgoef`）。

## 结论（先看这个）

| 问题 | 结论 |
| --- | --- |
| 宿主语言 | **XGo 一等**（应用写 `.xgo`）；框架/绑定可混编 Go |
| 编译器 | **LLGo**（LLVM，C ABI 直连，无标准 Go 的 cgo 切换开销；`xgo build` 经 `// llgo 1.0`） |
| 绑定生成 | **llcppg**，目标是 **CEF C API**，不是 `libcef_dll_wrapper` |
| 产品定位 | 类似 Tauri 2：Web 前端 + 原生后端 + IPC/插件；渲染引擎换成自带 Chromium 的 CEF |
| 与本仓库关系 | 与 `src/cefdemo` 并列；共用 `CEF_VERSION` / `CEF_ROOT`；demo 继续当 C++ 基准 |

CEF 的稳定出口是 C（见仓库 `docs/03-binary-boundary.md`）。llcppg 对复杂 C++（模板、mangling、STL）仍弱，因此 **绑定面选 C API** 是硬约束，不是偏好。

---

## 1. XGo（语言）

- 仓库：<https://github.com/goplus/xgo>
- 站点：<https://xgo.dev>
- 公式：`XGo := C * Go * Python * JavaScript + Scratch`
- 与 Go 源码级兼容，可同包混编；通过 LLGo 接入 C 生态
- 对本项目：应用与框架 API 可用 XGo 书写；底层绑定包保持常规 Go，便于 `llgo`/`xgo` 共用

## 2. LLGo（编译器）

- 仓库：<https://github.com/xgo-dev/llgo>（亦见 goplus/llgo）
- 定位：基于 LLVM 的 Go 编译器，按 **C ABI** 与 C/C++/Python/JS 等互操作
- 相对标准 Go + cgo 的关键点：
  - 原生 goroutine **1:1 映射 OS 线程**、固定栈
  - 调 C **无需** Go↔C 栈/调度切换，适合 CEF 这种高频、强线程亲和的 C API
- 依赖：Go 工具链、LLVM/Clang/LLD、bdwgc、libffi 等（以官方 README 为准）
- 对本项目：应用默认 `xgo build`（LLGo 后端）；不依赖 cgo

## 3. llcppg（绑定生成）

- 仓库：<https://github.com/goplus/llcppg> / <https://github.com/xgo-dev/llcppg>
- 作用：从 C/C++ 头文件自动生成 **LLGo bindings**
- 用法概要：写 `llcppg.cfg` → `llcppg llcppg.cfg` → 得到可 `llgo run` 的包
- 配置字段（常用）：`name`, `cflags`, `include`, `libs`, `trimPrefixes`, `cplusplus`, `deps`, `mix`, `typeMap`, `symMap`, `headerOnly`, `staticLib`
- **C++ 模式**：`cplusplus: true` 可用，但模板、跨平台 mangling、STL 字符串等仍有缺口（见上游 issue #281 等）
- 对本项目：
  - `"cplusplus": false`
  - `include` 指向 `CEF_ROOT/include/capi/*.h` 中的核心头
  - `libs` 链 `$CEF_ROOT/Release/libcef.so`（及平台等价物）
  - 生成物落在 `cef/capi/`（生成目录，可 gitignore 或分阶段提交）

## 4. Tauri 2（产品对标）

- 文档：<https://v2.tauri.app>
- 结构：系统 WebView（WRY/TAO）+ Rust 核心 + 可选 JS API
- 要抄的能力（不是实现细节）：
  - **Commands**：前端 `invoke` → 后端处理 → JSON 往返（类 fetch）
  - **Events**：双向、fire-and-forget
  - **Plugins**：生命周期钩子 + 命令注册 + 可选前端 API
  - **配置驱动**：`tauri.conf.json` 一类清单（窗口、权限、资源）
  - **Capabilities / 权限**：命令白名单（二期）
- 与 Tauri 的刻意差异：
  - 渲染：CEF 全量 Chromium（安装包更大，行为跨平台更齐）
  - 后端语言：XGo/Go + LLGo，不是 Rust
  - 不依赖系统 WebView（见 `docs/01-place.md` 对 Tauri 默认路径的对比）

## 5. CEF 与本仓库

- 官方：<https://chromiumembedded.github.io/cef/>，C API 说明 [Using the C API](https://chromiumembedded.github.io/cef/using_the_capi.html)
- 本仓库：`docs/` 架构说明 + `skills/` 下载/构建流程 + `src/cefdemo` C++ Views 样例
- 钉版本：`src/CEF_VERSION`（当前调研时为 152.x linux64）
- 不可违背的约束（与语言无关）：
  1. 先 `cef_execute_process` / `CefExecuteProcess`
  2. 再 `cef_api_hash`（C）或 wrapper 等价路径，再其它 API
  3. 同一可执行文件多进程；沙箱、cache 路径按 docs/skills
  4. JS 桥优先 `CefMessageRouter` 思路（`docs/07-javascript.md`），避免已废弃的 V8 extension

## 6. 同类 / 可借鉴项目

| 项目 | 路径 | 启发 | 不照搬的原因 |
| --- | --- | --- | --- |
| **energye/energy** | Go + LCL + CEF/WebView | 完整桌面框架、IPC、资源协议、跨平台 | 绑定与运行时路径不同；目标是 XGo/LLGo，不是 Energy 运行时 |
| **bnema/purego-cef** | Go + purego，CEF C API，无 cgo | **生成绑定 + 手写 ergonomic 层**；生命周期/refcount 分层 | 用 purego 动态加载；我们用 llcppg + LLGo 静态/链接式 C ABI |
| **turutcrane/cefingo** | cgo + CEF C API | capi 用法与线程注意点 | cgo 抢占问题；明确提醒未关抢占 |
| **Gurkengewuerz/cefgo** | 较旧 CEF C API 绑定 | JS↔Go 回调雏形 | CEF 版本过旧 |
| **cztomczak/cef2go** / **richardwilkes/cef** | 早期 Go CEF | 历史参考 | 基本停更 |
| **cefsimple_capi**（官方） | 纯 C 样例 | 正确的 C API 生命周期与 refcount | 非 Go；作绑定验收金样例 |

## 7. 技术选型矩阵

```text
                    标准 Go+cgo     purego-cef      LLGo+llcppg(本方案)
C 调用开销          高（栈切换）    中（动态符号）   低（C ABI 直连）
线程模型            M:N 抢占风险    同左             1:1 OS 线程，贴合 CEF
绑定生成            手工/脚本       自研生成器       llcppg（XGo 生态）
语言生态            Go              Go               XGo + Go
与 cefutil 对齐     弱              弱               强（同 CEF_ROOT/文档）
```

## 8. 风险与开放问题

1. **llcppg 对 CEF 头文件规模**：`include/capi` 头文件多、结构体/函数指针表密集，可能需分批 `include` + `mix`，或先做最小子集（app / browser / client / life_span / frame）。
2. **回调方向（C→XGo）**：Handler 是 C 函数指针表；需验证 LLGo 导出回调与 CEF 线程模型是否匹配，必要时加极薄 C shim（`cef/shim`）。
3. **Views / Chrome style**：C++ demo 用 Views；C API 是否完整暴露 Views 需对照当前 `CEF_ROOT/include/capi`。不行则一期用 `cef_browser_host_create_browser` + 平台窗口。
4. **打包与多进程**：辅助进程、rpath、Resources 拷贝仍按 `cef-create-project` / cefdemo 的 CMake 经验，由 `cmd/xgoef` 或脚本复用。
5. **许可**：本仓库 GPL-3.0；CEF/Chromium 有各自许可，发行时需分开说明。

## 9. 建议实施顺序

1. 本骨架（目录、接口、配置、调研文档）← 当前
2. 最小 `llcppg.cfg`：只绑 `cef_app_capi.h` / `cef_browser_capi.h` 等启动闭环
3. `examples/hello`：`execute_process` → `initialize` → 空白页 → `shutdown`
4. IPC：MessageRouter 风格 Commands + 前端 `api/js/xgoef.js`
5. Plugin + 配置权限模型（对标 Tauri 2 capabilities）
6. CLI：`xgoef create|dev|build`
