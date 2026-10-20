# xgoef × XGo Classfile：领域友好设计

> XGo 哲学：**Don't define a language for specific domain. Abstract domain knowledge for it.**  
> （SDF — Specific Domain Friendliness，不是再发明一门 DSL）

本文回答：在 xgoef（对标 Tauri 2 的 CEF 桌面壳）里，**如何用 classfile / class framework** 服务两类人：

1. **应用作者** — 写业务 command、窗口、生命周期，越少样板越好  
2. **插件作者** — 扩展系统能力（fs/dialog/…），可发布、可权限化  

依据：[goplus/xgo `doc/classfile.md`](https://github.com/goplus/xgo/blob/main/doc/classfile.md)、已有框架（spx / yap / cobra / mcp）、以及 `docs/tauri2-gap.md`。

---

## 1. 结论先看

| 问题 | 建议 |
| --- | --- |
| 要不要为 xgoef 发明 DSL？ | **不要**。做 **class framework**，后缀专用（如 `*.xef`） |
| classfile 主要给谁用？ | **应用作者**写 app / command / window；**官方体验优先** |
| 插件怎么开发？ | **双轨**：应用内轻量 `*_plugin.xef`；可发布插件用 **Go/XGo 包 + Plugin 契约**（对齐 Tauri crate） |
| 和现在的 `main.xgo` 什么关系？ | `main.xgo` 是过渡；目标态迁到 `hello_app.xef` + `greet_cmd.xef` |
| 何时做？ | **P0 IPC/ACL 之后**；classfile 依赖稳定的 `runtime.App` / `ipc.Router` |

一句话：**用 classfile 抽象「桌面应用装配」；用 Plugin 包契约抽象「可分发扩展」——两者互补，不是二选一。**

---

## 2. XGo Classfile 能给我们什么

Class framework 编译器自动完成：

1. 每个 classfile → 嵌入 `BaseType` 的 struct（变量=字段，函数=方法）  
2. 工程根 classfile → `MainEntry`；工作类 → `Main`  
3. 合成 `Main`：实例化所有 work class → 调用框架的 **Template Recv Method**（`XGot_<Base>_Main`）  
4. 生成 `func main() { new(App).Main() }`

框架作者写的是：**BaseType + Protocol + Template Recv Method**（领域知识）。  
应用作者写的是：**几乎只有业务逻辑的 `.xef` 文件**。

对标已有框架：

| 框架 | 领域知识 | 文件角色 |
| --- | --- | --- |
| spx | 游戏循环 / 精灵 | `main.spx` + `*.spx` |
| yap | HTTP 路由 | `*.yap`（文件名=路由） |
| cobra | CLI 子命令 | `*_app.gox` + `*_cmd.gox` |
| mcp | MCP resource/tool/prompt | 多种 `*_xxx.gox` |
| **xgoef（拟）** | CEF 桌面壳：窗、invoke、插件装配 | 见下 |

---

## 3. 推荐的 xgoef Class Framework 形状

### 3.1 专用后缀（硬约束）

`gox.mod` **禁止**抢占裸 `*.gox` / `*.xgo`。必须用框架独占后缀。

推荐：

```text
*.xef            — 可选：仅用于极简单文件 demo（慎用，易与多 kind 冲突）
*_app.xef        — 工程根（project）
*_cmd.xef        — IPC command（work）
*_win.xef        — 窗口定义（work，P1）
*_plugin.xef     — 应用内嵌插件（work，可选）
```

### 3.2 `gox.mod`（草案）

放在模块 `github.com/fanfeilong/cefutil/src/xgoef/xclass`（或独立模块 `github.com/fanfeilong/xgoef`）：

```go.mod
xgo 1.7

# project: 桌面应用根
project *_app.xef App github.com/fanfeilong/cefutil/src/xgoef/xclass

# work kinds（多 kind → 多 Protocol，对齐 mcp）
class -prefix=Cmd_    *_cmd.xef    Command   CommandProto
class -prefix=Win_    *_win.xef    Window    WindowProto
class -prefix=Plug_   *_plugin.xef Pluglet   PlugletProto
```

说明：

- `-prefix=` 避免 `run`/`fs` 等短名撞标识符（学 cobra）。  
- 三种 work class 分进 Template Recv 的三个 typed slice（学 mcp）。  
- 一期可只实现 `*_app.xef` + `*_cmd.xef`；window/plugin classfile 后补。

### 3.3 BaseType / Protocol（领域知识落点）

```go
// xclass/app.go — project base
type App struct {
    // 嵌入或持有 *runtime.App / conf / cef Host
    ConfPath string
}

func (a *App) MainEntry() {
    // 默认：读 xgoef.conf.json、准备 Host；用户代码可覆盖前后钩子
}

// xclass/command.go
type Command struct {
    Name       string   // 默认取自文件名 greet_cmd.xef → "greet"
    Permission string   // 默认同 Name；可覆盖为 "greet" / "fs:read"
}

type CommandProto interface {
    Main() // 或 Main(ctx *CmdContext)
    Spec() CommandSpec
}

// xclass/window.go（P1）
type Window struct {
    Label  string
    Title  string
    Width  int
    Height int
    URL    string
}

// xclass/pluglet.go — 应用内嵌插件（不是可发布包）
type Pluglet struct{}
type PlugletProto interface {
    Main()
    Name() string
}
```

### 3.4 Template Recv Method

```go
// XGot_App_Main — 编译器装配完成后交给框架
func XGot_App_Main(
    app     AppProto,
    cmds    []CommandProto,
    wins    []WindowProto,
    plugs   []PlugletProto,
) {
    // 1. 构造 runtime.App + CEF Host
    // 2. 注册每个 cmd → ipc.Router.Handle(spec.Name, ...)
    // 3. 合并 capabilities（文件声明 + conf）
    // 4. Setup 每个 pluglet / 外部 Plugin 包
    // 5. 按 wins（或 conf.windows）创建浏览器
    // 6. RunLoop；退出时 Teardown
}
```

这是 **CEF 生命周期 + IPC + ACL** 的唯一集中点；classfile 作者不应直接碰 `ExecuteProcess`。

---

## 4. 应用作者怎么写（SDF 体验）

### 4.1 目标目录（hello 演进）

```text
examples/hello/
  hello_app.xef          # project
  greet_cmd.xef          # command
  main_win.xef           # window（P1；一期可用 conf）
  xgoef.conf.json
  web/
```

### 4.2 `hello_app.xef`（草图）

```xgo
# 嵌入 xclass.App；顶层代码 → MainEntry
ConfPath = "xgoef.conf.json"

# 可选钩子：在框架跑 CEF 前/后
func OnReady() {
    echo "app ready"
}
```

### 4.3 `greet_cmd.xef`（草图）

```xgo
# 文件名 greet_cmd.xef → 默认 command 名 "greet"
# Permission 默认同名，可写入 capabilities

func Main(ctx *CmdContext) (any, error) {
    name := "world"
    if n, ok := ctx.Arg("name").(string); ok && n != "" {
        name = n
    }
    return "hello, ${name} (from xgoef)", nil
}
```

或更贴近当前 Router 签名：

```xgo
func Main(ctx *CmdContext, argsJSON []byte) (any, error) {
    ...
}
```

**体验对比：**

| 现在（`main.xgo`） | Classfile 目标 |
| --- | --- |
| 手写 CEF ExecuteProcess / InitAPIHash | 框架 Template Recv 代劳 |
| `app.Handle "greet", (ctx, args) => {...}` | 一个 `greet_cmd.xef` 文件 |
| conf 与代码分离但 command 挤在 main | command 按文件拆分，可多人并行 |

### 4.4 文件名即路由？（可选，学 yap）

一期：**文件名 → command 名**（`greet_cmd.xef` → `greet`）足够。  
二期若需要嵌套：`user_get_cmd.xef` → `user.get`，或显式 `Name = "user.get"`。

**不要**把 HTTP 路径编码进桌面 invoke 名；invoke 保持 Tauri 式扁平/点号命令空间。

---

## 5. 插件：两条轨道

Tauri 的插件 = **可发布的扩展包**（Rust crate + 可选 npm）。  
XGo classfile 更擅长 **工程内装配**。因此拆成：

### 轨道 A — 应用内嵌 Pluglet（classfile）

适用：只服务当前 app、不发布的胶水逻辑。

```text
analytics_plugin.xef
```

```xgo
func Name() string { return "analytics" }

func Main() {
    # 注册额外 command、挂生命周期
    handle "track", (ctx, args) => {
        return nil, nil
    }
}
```

`XGot_App_Main` 在 Setup 阶段调用各 Pluglet。

### 轨道 B — 可发布 Plugin 包（主路径，对标 Tauri）

适用：`fs` / `dialog` / `clipboard` 等官方或社区插件。

```text
module github.com/example/xgoef-plugin-fs

xgoef/
  plugin.go      # 实现 plugin.Plugin
  commands.go    # 或内部再用 *_cmd.xef（插件包自己的 gox.mod）
  api/js/fs.js   # 可选前端 API
```

**契约（相对现有 `plugin.Plugin` 扩展）：**

```go
type Plugin interface {
    Name() string
    Setup(app *runtime.App) error
    Teardown(app *runtime.App) error
    Commands() []CommandSpec    // name + permission id
    Permissions() []Permission  // 供 ACL 声明
}
```

应用侧装配：

```json
// xgoef.conf.json
"plugins": ["github.com/example/xgoef-plugin-fs"]
```

或在 `hello_app.xef`：

```xgo
use fsplugin.New()
```

**插件作者是否也用 classfile？**

- **可以**：插件模块自带 `gox.mod`，用 `*_cmd.xef` 声明该插件的 commands（内部 SDF）。  
- **不强制**：纯 Go/`*.xgo` 实现 `Plugin` 即可（更易被非 XGo 贡献者接受）。  
- **官方插件**建议：对外 Go API 稳定；对内可用 classfile 减样板。

```text
                    ┌─────────────────────────┐
                    │  应用仓库（*.xef）         │
                    │  app + cmd + win + pluglet│
                    └───────────┬─────────────┘
                                │ XGot_App_Main
                    ┌───────────▼─────────────┐
                    │  xgoef runtime + ACL     │
                    └───────────┬─────────────┘
              ┌─────────────────┼─────────────────┐
              ▼                 ▼                 ▼
        内嵌 Pluglet      官方 Plugin 包      社区 Plugin 包
        (*_plugin.xef)    (fs/dialog/…)      (Go module)
```

---

## 6. 与 Tauri 概念对照

| Tauri | xgoef Classfile / 包 |
| --- | --- |
| `src-tauri/src/lib.rs` + `#[tauri::command]` | `*_app.xef` + `*_cmd.xef` |
| `tauri.conf.json` windows | `*_win.xef` 和/或 conf（conf 保留给非代码配置） |
| `tauri-plugin-*` crate | 轨道 B Plugin 模块 |
| `app.manage` / 状态 | `App` 上的字段（classfile `var`）或显式 Store 插件 |
| Capabilities | conf + command `Permission` 字段；框架注册时写入 ACL |
| `invoke('plugin:fs\|read')` | `Permission` / 命令名约定 `fs:read` |

---

## 7. 刻意不做的「伪 DSL」

| 诱惑 | 为什么不做 |
| --- | --- |
| 新语言关键词 `window` / `command` 进 XGo 核心 | 违反 SDF；用 BaseType 方法即可 |
| 用 classfile 替代 `xgoef.conf.json` 全部配置 | 配置需无代码、CI、产品可读；代码只表达行为 |
| 一个 `*.xef` 表达整份 CEF C API | 绑定层保持 Go + llcppg；classfile 停在应用域 |
| 强制所有插件必须是 classfile | 抬高贡献门槛；包契约优先 |

允许的「领域友好」API（在 BaseType 上，不是新语法）：

```xgo
handle "greet", (ctx, args) => { ... }   # App/Command 提供的方法
listen "tick", payload => { ... }
emit "tick", {"n": 1}
createWindow label = "main", title = "Hi", width = 960, height = 640
```

这是 **普通 XGo 方法 + 命令式调用**，不是 DSL。

---

## 8. 实现分期

| 阶段 | 交付 | 依赖 |
| --- | --- | --- |
| C0 设计（本文） | `docs/classfile-design.md` | — |
| C1 骨架 | `xclass` 包：`App`/`Command`/`XGot_App_Main` 空实现 + `gox.mod` | P0 IPC 可后接 |
| C2 hello 迁移 | `hello_app.xef` + `greet_cmd.xef`；保留 conf | C1 + P0 桥 |
| C3 Window classfile | `*_win.xef` | P1 Window 模型 |
| C4 Pluglet | `*_plugin.xef` | C2 |
| C5 可发布插件模板 | `xgoef create-plugin` + 示例 `fs` | P1 Plugin 契约 |

**注意：** classfile 的 `main` 合成与现有「手动 `ExecuteProcess`」必须统一进 `XGot_App_Main`，避免双入口。

---

## 9. 最小可运行想象（C2）

```text
# greet_cmd.xef
func Main(ctx *CmdContext) (any, error) {
    return "hello, ${ctx.ArgString("name", "world")}", nil
}

# hello_app.xef
ConfPath = "xgoef.conf.json"
```

框架负责：CEF、注入 `xgoef.js`、注册 `greet`、执行 capabilities、开窗加载 `web/`。

应用作者 **看不到** `LockOSThread` / `InitAPIHash` / shim——那些是领域知识，藏在 `xclass` 里。

---

## 10. 开放问题

1. **后缀最终名**：`*.xef` vs `*_xgoef.gox`？（推荐短后缀 `*.xef`，与 yap/spx 一致）  
2. **Command.Main 签名**：固定 `(ctx, argsJSON)` 还是导出类型化 `Arg` helper？  
3. **与 LLGo**：classfile 生成 Go 后再 `llgo build`；`XGO_GOCMD` 包装器继续用。  
4. **测试**：官方插件与 app 用内置 `*_test.gox`；CEF 相关测放 integration。  
5. **是否允许纯 `main.xgo` 逃生舱**：是——高级用户可不用 classfile，直接调 `runtime.App`（文档标明「进阶」）。

---

## 11. 与其它文档的关系

| 文档 | 关系 |
| --- | --- |
| `DESIGN.md` | 分层总览；本文是应用/插件 **作者体验** 层 |
| `tauri2-gap.md` | 能力缺口；本文是用 XGo SDF **填缺口时的写法** |
| `roadmap.md` | 增加 classfile 分期条目 |

---

## 12. 一句话

**xgoef 不要做桌面 DSL；要做「桌面应用 class framework」：**  
用 `*_app.xef` / `*_cmd.xef` /（可选）`*_win.xef` 让应用作者只写业务；  
用 **Plugin 包契约** 让插件可发布、可 ACL；  
应用内嵌扩展再用 `*_plugin.xef`。  
这才是 XGo classfile 与 Tauri 插件模型在 xgoef 上的正确结合点。
