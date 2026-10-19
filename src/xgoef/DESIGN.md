# xgoef 架构设计

对标 Tauri 2 的分层，后端换成 XGo/LLGo，WebView 换成 CEF。

## 总览

```text
┌──────────────────────────────────────────────────────────┐
│  Web UI (任意前端栈)                                      │
│  api/js → window.__XGOEF__.invoke / listen / emit         │
└────────────────────────────┬─────────────────────────────┘
                             │ JSON 消息（异步）
┌────────────────────────────▼─────────────────────────────┐
│  xgoef 应用层（XGo 一等；可混编 Go）                        │
│  runtime.App · ipc.Router · plugin.Registry · conf.File  │
└────────────────────────────┬─────────────────────────────┘
                             │ 手写 ergonomic API
┌────────────────────────────▼─────────────────────────────┐
│  cef 包（llcppg 生成 + 少量手写）                          │
│  C API 类型、refcount、string、生命周期封装                │
└────────────────────────────┬─────────────────────────────┘
                             │ LLGo C ABI
┌────────────────────────────▼─────────────────────────────┐
│  libcef（Spotify cef_binary_*，与 src/CEF_VERSION 对齐）   │
└──────────────────────────────────────────────────────────┘
```

## 包职责

| 包 | 职责 |
| --- | --- |
| `cef/` | llcppg 配置与生成绑定；禁止业务逻辑 |
| `cef/shim/` | 仅当回调/Views 无法直接从 Go 表达时使用的薄 C 层 |
| `runtime/` | 进程入口、初始化、消息循环、窗口/浏览器 |
| `ipc/` | Command 注册与调用、Event 总线 |
| `plugin/` | 插件生命周期与命令挂载 |
| `conf/` | `xgoef.conf.json` 解析 |
| `api/js/` | 注入或随资源分发的前端桥 |
| `cmd/xgoef/` | 脚手架与构建辅助（后期） |
| `examples/hello/` | 最小可运行目标 |

## 对标 Tauri 2 的 API 形状

### Commands

```xgo
# 后端（XGo：命令式 + lambda + ${expr}）
app.Handle "greet", (ctx, argsJSON) => {
    return "hello, xgo", nil
}

# 前端
const msg = await window.__XGOEF__.invoke('greet', { name: 'xgo' })
```

### Events

```xgo
app.Emit "tick", {"n": 1}
app.Listen "ui-ready", payload => {
    echo "ready:", payload
}
```

### Plugins

```xgo
app.Use fsplugin.New()
# 插件可注册命令、挂钩 OnReady / OnExit
```

### 配置

见 `xgoef.conf.schema.json` / `examples/hello/xgoef.conf.json`：窗口标题尺寸、启动 URL、CEF 路径提示、允许的 command 列表。

## 进程与线程（必须遵守）

1. 每个进程入口最先：`ExecuteProcess`；子进程直接退出码返回。
2. 主进程：`InitAPIHash` → 填 `Settings`（含绝对 `root_cache_path`）→ `Initialize` → 建浏览器 → `RunLoop` → `Shutdown`。
3. UI 相关 CEF 调用在 UI 线程；LLGo 1:1 线程模型下仍应用 `runtime.LockOSThread`（或框架保证主 goroutine 钉在主线程）。
4. 关闭：两段式（先请求关浏览器，最后一次 `OnBeforeClose` 再 `QuitLoop`），与 `docs/05-lifecycle.md` 一致。

## IPC 实现路线

一期不做完整 `CefMessageRouter` C++ wrapper 绑定，采用：

1. 渲染侧注入 `api/js/xgoef.js`（或 `OnContextCreated` 挂 `window.__XGOEF__`）。
2. 查询走 CEF process message / 官方 message_router 的 C 等价物（绑定就绪后切换）。
3. 载荷约定为 JSON；命令名字符串路由到 `ipc.Router`。

二期对齐 Tauri capabilities：配置里声明哪些窗口可调哪些 command。

## 与 cefdemo 的分工

| | `src/cefdemo` | `src/xgoef` |
| --- | --- | --- |
| 语言 | C++ | XGo / Go + LLGo |
| API | C++ wrapper + Views | C API（+ 必要时 shim） |
| 角色 | 技能验证过的黄金样例 | 框架实验与产品化方向 |
| CEF 包 | 同一 `CEF_ROOT` | 同一 `CEF_ROOT` |

## 非目标（当前阶段）

- 移动端
- 替代本仓库 skills 的下载逻辑（继续用 `cef-download`）
- 完整复刻 Tauri 插件市场与签名更新器
- 在未安装 LLGo/llcppg/CEF 的环境下保证 `llgo build` 通过（骨架以接口与文档为主）

## 对标 Tauri 2 的差距

分层与 API **形状**见上文；与 Tauri 2 的能力对照、优先级（P0–P3）和补齐设计见：

→ **[docs/tauri2-gap.md](docs/tauri2-gap.md)**

## XGo Classfile / 插件作者体验

XGo 的 classfile（SDF，不是新 DSL）如何服务「写应用 / 写插件」见：

→ **[docs/classfile-design.md](docs/classfile-design.md)**

要点：应用用 `*_app.xef` + `*_cmd.xef` 装配；可发布插件走 Go/XGo **Plugin 包契约**；应用内嵌扩展可用 `*_plugin.xef`。