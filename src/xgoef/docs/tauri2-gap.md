# 对标 Tauri 2：能力差距与设计

调研基准：[Tauri 2](https://v2.tauri.app/)（文档 / Architecture / IPC / ACL / Plugins，2026-10）。  
对照对象：本仓库 `src/xgoef` 当前实现（XGo 应用 + LLGo + CEF C API）。

本文回答三件事：

1. **我们已有什么**（相对 Tauri 的心智模型）
2. **还差什么**（按能力域拆开，区分「未做」与「刻意不做」）
3. **怎么补**（分优先级的设计落点，不估日历工期）

相关文档：`DESIGN.md`（分层）、`RESEARCH.md`（选型）、`docs/roadmap.md`（阶段）、仓库根 `docs/07-javascript.md`（CEF JS 桥约束）。

---

## 1. 对标边界（先对齐预期）

| 维度 | Tauri 2 | xgoef | 含义 |
| --- | --- | --- | --- |
| 渲染 | 系统 WebView（WRY/TAO） | **自带 Chromium 的 CEF** | 包体更大、行为更齐、跨发行版一致；**不追求 Tauri 的小体积** |
| 后端语言 | Rust（+ 插件侧 Swift/Kotlin） | **XGo 一等**，底层可混编 Go；LLGo 链 CEF | API 形状学 Tauri，实现栈不同 |
| 平台 | Desktop + **Mobile** | **Desktop-first**（现阶段 Linux 可跑） | 移动端标为非目标 |
| 安全模型 | Capabilities / ACL 一等公民 | 配置里有占位，**未执行** | 必须补，否则 invoke 裸奔 |
| 插件生态 | 官方插件矩阵 + 社区 | 仅接口骨架 | 官方插件不照搬，先定契约再挑子集 |

**刻意不对齐的点（非缺口）：**

- 不换成系统 WebView；继续 CEF（本仓库定位）。
- 不追求 &lt;1MB 安装包；接受 `libcef` 量级。
- 不复制 Tauri 移动端与全部官方插件市场。
- 更新器 / 签名分发可后置；本仓库已有 `docs/10-ship-and-upgrade.md` 可复用思路。

---

## 2. 当前 xgoef 快照（「有」到什么程度）

| 能力 | 状态 | 证据 |
| --- | --- | --- |
| 应用用 XGo 书写 | ✅ | `examples/hello/main.xgo`，`xgo build` → LLGo |
| 窗口 + 本地页 | ✅（薄） | CEF Views/`cefsimple_capi` 路径；URL 经 `--url=` |
| 配置文件 | 🟡 骨架 | `xgoef.conf.json` + schema；窗口字段极少 |
| Command 路由（宿主侧） | 🟡 内存态 | `ipc.Router.Handle` / `Invoke` 可单测；**未接到 CEF** |
| Event 总线（宿主侧） | 🟡 进程内 | `ipc.Bus`；**不跨渲染进程** |
| 前端 JS 桥 | 🟡 占位 | `api/js/xgoef.js`（`invoke`/`listen`/`emit`）；无 `__xgoefPost` → 报 *host bridge not connected* |
| Plugin 生命周期 | 🟡 接口 | `Setup`/`Teardown`；无命令注入、无 JS API 包 |
| Capabilities | ⬜ 配置占位 | schema 有 `capabilities`；**运行时不校验** |
| CLI（create/dev/build） | ⬜ | `cmd/xgoef` 仅 `info` |
| 多窗口 / 托盘 / 菜单 | ⬜ | — |
| 资源协议 / 资产嵌入 | ⬜ | 现用 `file://` |
| 官方类插件（fs/dialog/…） | ⬜ | — |

结论：我们有 **Tauri 式分层的轮廓** 和 **能出窗的 CEF 宿主**；缺的是把「前端 ↔ 宿主」真正打通，以及围绕安全与窗口/资源的产品能力。

---

## 3. Tauri 2 能力地图 → xgoef 差距矩阵

图例：✅ 已有 · 🟡 部分 · ⬜ 缺失 · ✕ 刻意不做（或远期）

### 3.1 核心运行时与窗口

| Tauri 能力 | xgoef | 差距说明 |
| --- | --- | --- |
| App / Builder 生命周期 | 🟡 | 有 `runtime.App.Run`；缺精细 hook（ready / exit / page-load） |
| Window / Webview 多实例 | ⬜ | 仅第一窗；无 label、无 `createWindow` |
| Window 状态（位置/尺寸/最大化） | ⬜ | conf 仅 title/width/height/url |
| 菜单 / 托盘 / 任务栏 | ⬜ | CEF/Views 可做，未封装 |
| 主题 / 装饰 / 透明 | ⬜ | — |
| 多显示器 / DPI | ⬜ | 依赖 CEF 默认 |
| 打印 / 开发者工具 API | ⬜ | CEF 支持，未暴露 XGo API |

### 3.2 IPC（Commands / Events / Channels）

| Tauri 能力 | xgoef | 差距说明 |
| --- | --- | --- |
| `invoke` → 后端 command | 🟡 | JS + Router 齐，**传输层未接**（MessageRouter / process message） |
| 命令参数/返回值 JSON | 🟡 | 有；缺 TypeScript 类型生成、缺结构化错误码 |
| 异步 command | ⬜ | Handler 同步；无取消、无超时 |
| Events 双向 | 🟡 | JS/宿主各有本地订阅；**跨进程未通** |
| 按 webview/window 定向 emit | ⬜ | — |
| Channels（有序大数据流） | ⬜ | 下载进度、子进程 stdout 等场景需要 |
| `eval` / 注入脚本 | ⬜ | CEF 可做；未封装 |
| IPC 权限门禁 | ⬜ | 见 ACL |

**P0 设计落点（IPC）：** 见 §5.1。

### 3.3 安全：Capabilities / ACL / CSP

| Tauri 能力 | xgoef | 差距说明 |
| --- | --- | --- |
| Capability 绑定窗口/webview | ⬜ | conf 占位未执行 |
| Permission 标识（`plugin:perm`） | ⬜ | — |
| 路径/命令 scope（allow/deny） | ⬜ | — |
| remote URL 能力边界 | ⬜ | 当前可加载任意 `file://`/`http(s)` |
| CSP / 资产完整性 | ⬜ | Tauri codegen 侧有；我们未做 |
| 隔离：未匹配 capability → 无 IPC | ⬜ | 必须有，否则 greet 模型不可产品化 |

**P0 设计落点（ACL）：** 见 §5.2。

### 3.4 前端工程与资产

| Tauri 能力 | xgoef | 差距说明 |
| --- | --- | --- |
| `@tauri-apps/api` 风格包 | 🟡 | 单文件 `xgoef.js`；无 npm 包、无 TS |
| `withGlobalTauri` / 自动注入 | ⬜ | 需在 `OnContextCreated` 注入桥 |
| `frontendDist` / `devUrl` 热更新 | 🟡 | conf 有字段；无 `xgoef dev` |
| 自定义协议（`asset://` / `tauri://`） | ⬜ | 现 `file://`；缺沙箱与范围限制 |
| 图标 / 资源嵌入打包 | ⬜ | 手工拷贝 `web/` + CEF Resources |

### 3.5 插件系统

| Tauri 能力 | xgoef | 差距说明 |
| --- | --- | --- |
| Plugin = 原生逻辑 + 可选 JS API | 🟡 | 仅 Setup/Teardown |
| 插件注册 command / 权限声明 | ⬜ | — |
| 官方插件矩阵（fs/dialog/shell/…） | ⬜ | 按需移植子集，不 mirror 全表 |
| 移动端插件（Swift/Kotlin） | ✕ | 非目标 |

**建议首批官方插件（桌面）：** `fs`（scoped）、`dialog`、`clipboard`、`shell`（受限 open）、`os`、`process`、`notification`、`global-shortcut`、`window-state`。其余（sql/updater/websocket…）后置。

### 3.6 工具链与 DX

| Tauri 能力 | xgoef | 差距说明 |
| --- | --- | --- |
| `create-tauri-app` | ⬜ | `xgoef create` 占位 |
| `tauri dev` / `tauri build` | ⬜ | 现有 `scripts/build-xgoef-hello.sh` |
| 多前端模板（Vite/React/…） | ⬜ | — |
| 平台打包（msi/dmg/deb/AppImage） | ⬜ | 依赖 CEF 运行时布局；需打包器设计 |
| IDE / 代码生成（command 宏） | ⬜ | XGo 侧可用 codegen 或约定式注册 |

### 3.7 发行与运维

| Tauri 能力 | xgoef | 差距说明 |
| --- | --- | --- |
| Updater 插件 | ⬜ / 远期 | 可对齐 `docs/10-ship-and-upgrade.md` |
| 单实例 | ⬜ | — |
| Autostart / Deep link | ⬜ | — |
| 崩溃报告 / 日志插件 | ⬜ | CEF crashpad 可接 |

### 3.8 移动端

| Tauri 能力 | xgoef |
| --- | --- |
| iOS / Android | ✕ 非目标（CEF 移动路径与 Tauri 不同） |

---

## 4. 优先级（按「产品可用」排序）

### P0 — 没有就不叫 Tauri-like 应用

1. **IPC 传输接通**：`xgoef.js` ↔ `ipc.Router`（CEF MessageRouter 或 process message + C shim）。
2. **桥注入**：浏览器上下文创建时注入 `__XGOEF__` / `__xgoefPost`。
3. **Capabilities 强制执行**：未授权 command → 拒绝；默认拒绝未声明窗口。
4. **hello 端到端**：按钮 `invoke('greet')` 返回宿主字符串（不再 *host bridge not connected*）。

### P1 — 可做正经桌面壳

5. Window API：label、多窗、关闭/聚焦、基础状态。
6. 自定义资源协议（或严格 `asset` 根目录）替代裸 `file://`。
7. Event 跨进程（含 `emit_to(label)`）。
8. Plugin 契约：注册 commands + 可选 `api/js/plugins/<name>.js`。
9. `xgoef dev` / `xgoef build`（封装 CEF 拷贝、llgo/xgo、前端 dist）。

### P2 — 体验对齐

10. Channels、异步 command、超时/取消。
11. 首批系统插件：fs（scoped）、dialog、clipboard、os。
12. 菜单 / 托盘、window-state 持久化。
13. TS 类型与 `invoke` 声明生成。

### P3 — 发行与生态

14. 打包器、updater、单实例、deep link。
15. 更多插件；文档站点与模板。

---

## 5. 关键缺口的设计草案

### 5.1 IPC 传输（P0）

**目标：** 前端 `invoke` 进入 `ipc.Router.Invoke`，响应回 Promise。

**推荐路径（与仓库 `docs/07-javascript.md` 一致）：**

1. **短期（shim 扩展）**  
   - 在 `cef/shim` 增加 MessageRouter（或简化的 `cef_process_message`）C API。  
   - 浏览器进程：收到 query → 调 `App.Router.Invoke` → 回写。  
   - 渲染进程：`xgoef.js` 的 `__xgoefPost` 接到 `cefQuery` / 等价物。

2. **中期**  
   - llcppg 生成更多 `capi`；shim 只保留回调难写部分。

**消息约定（已与 `xgoef.js` 对齐，保持稳定）：**

```json
{ "type": "invoke", "id": 1, "cmd": "greet", "args": { "name": "xgo" } }
{ "type": "response", "id": 1, "ok": true, "payload": "hello, xgo (from xgoef)" }
{ "type": "event", "event": "tick", "payload": { "n": 1 } }
```

**线程：** UI/浏览器线程投递；重活丢到工作线程再 `post` 回 UI（CEF 规则不变）。

### 5.2 Capabilities（P0）

**配置（演进现有 schema，对齐 Tauri 语义但简化）：**

```json
"capabilities": [
  {
    "identifier": "main",
    "windows": ["main"],
    "permissions": ["greet", "core:default"]
  }
]
```

**运行时：**

- 每个 WebView 绑定 window label（首版默认 `"main"`）。
- `Invoke` 前：`cap.Allows(label, cmd)`；失败返回明确错误（勿静默）。
- 插件权限形如 `fs:read`；应用命令用短名 `greet`。
- **默认拒绝**：无匹配 capability 的窗口不能 IPC。

一期不做 remote URL ACL、不做路径 scope 引擎；fs 插件落地时再加 allow 列表。

### 5.3 Window 模型（P1）

```text
App
 └─ Window (label, title, bounds, url)
      └─ WebView (CEF browser) 
```

- conf `windows[].label`（缺省 `main`）。
- XGo API 草图：`app.CreateWindow {label, title, width, height, url}`、`app.Window(label).Close()`。
- 映射 CEF Views / Alloy：优先继续 `cefsimple_capi` Views 路径，避免两套窗口栈。

### 5.4 资源与协议（P1）

- 开发：`devUrl`（Vite）+ 可选 CORS/unsafe 仅 dev。
- 生产：`frontendDist` 映射为 `xgoef://app/...`（CEF 自定义 scheme）或受限目录的 `file://`。
- 禁止默认放开任意本地路径读取（与 capabilities / fs scope 联动）。

### 5.5 Plugin 契约（P1）

```xgo
type Plugin interface {
    Name() string
    Setup(app *runtime.App) error
    Teardown(app *runtime.App) error
    # 新增：
    # Commands() []CommandSpec   // name + permission id
    # Permissions() []string
}
```

- JS：约定 `api/js/plugins/<name>.js` 或打包进单一 `xgoef.js` 命名空间 `__XGOEF__.plugins.fs`。
- 配置 `plugins: ["fs", "dialog"]` → 构建期或启动期装载。

### 5.6 CLI（P1）

| 命令 | 行为 |
| --- | --- |
| `xgoef create` | 模板：`main.xgo` + `xgoef.conf.json` + `web/` |
| `xgoef dev` | 可选拉起前端 dev server；`xgo build`/`llgo`；拷贝 CEF runtime；热重载策略二期 |
| `xgoef build` | 正式产物目录布局（hello 脚本 upstream） |

复用现有 `scripts/build-xgoef-hello.sh` / `install-xgo-toolchain.sh`，避免两套工具链逻辑。

---

## 6. 架构演进图（目标态）

```text
┌────────────────────────────── Web ──────────────────────────────┐
│  @xgoef/api (TS)  ·  __XGOEF__.invoke / listen / channel         │
└───────────────────────────────┬─────────────────────────────────┘
                                │ CEF MessageRouter / scheme
┌───────────────────────────────▼─────────────────────────────────┐
│  runtime.App (XGo/Go)                                             │
│   ├─ acl.Guard(capabilities)                                      │
│   ├─ ipc.Router + EventHub + Channel                              │
│   ├─ window.Registry                                              │
│   └─ plugin.Registry → fs/dialog/...                              │
└───────────────────────────────┬─────────────────────────────────┘
                                │ LLGo C ABI
┌───────────────────────────────▼─────────────────────────────────┐
│  cef/ (llcppg + thin shim)  →  libcef                             │
└─────────────────────────────────────────────────────────────────┘
```

---

## 7. 验收标准（相对 Tauri 心智，而非功能清单）

**里程碑 A（P0 完成）：**

- hello 页点击 *invoke greet* 显示宿主返回值。
- 未在 capabilities 声明的 command 调用失败。
- 文档：`docs/run.md` 描述桥接与 ACL。

**里程碑 B（P1 核心）：**

- 至少两个 window label；从宿主 `emit` 到指定窗。
- 一个 scoped `fs` 或 `dialog` 插件可演示。
- `xgoef build` 打出可拷贝运行的目录。

**里程碑 C：**

- Channels 或异步 command 之一落地；打包/更新器设计评审通过。

---

## 8. 风险与依赖

| 风险 | 影响 | 缓解 |
| --- | --- | --- |
| MessageRouter C API 绑定成本 | 卡住 P0 | 先 shim 手写最小 query/response |
| CEF 多进程 vs LLGo 宿主 | 辅助进程不稳 | 保持 `xgoef-helper`；GPU 默认 disable |
| ACL 过严/过松 | 安全或 DX | 默认拒绝 + hello 明示 permission |
| 插件范围膨胀 | 失焦 | 严格按 §3.5 首批列表 |
| 与 Tauri API 名一一对应强迫症 | 无意义工时 | **对齐心智，不保证符号同名** |

---

## 9. 文档与代码后续动作

1. 实现 P0：shim 桥 + ACL + hello E2E（更新 `roadmap.md` 状态）。
2. 扩展 `xgoef.conf.schema.json`：`windows[].label`、正式 `capabilities[]`。
3. 新增 `docs/ipc.md`、`docs/acl.md`（实现时写；本稿为总设计）。
4. `DESIGN.md` 仅保留分层总览，差距与演进以本文为准。

---

## 10. 一句话结论

相对 Tauri 2，xgoef **已具备「XGo 后端 + CEF 窗 + 配置/IPC/插件的外形」**；真正的产品差距集中在：

1. **IPC 未接通**  
2. **无强制 ACL**  
3. **窗口/资源/插件/CLI 仍停留在骨架**

先打通 P0，再按 P1 做壳，系统插件与发行能力放在其后。刻意保留 CEF 与桌面范围，不在体积和移动端上与 Tauri 竞赛。
