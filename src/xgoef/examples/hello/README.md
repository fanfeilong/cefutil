# hello

最小 xgoef 应用形状：配置 + `greet` 命令 + 本地 `web/`。

当前无 CEF 绑定时：

```bash
cd src/xgoef
go test ./...
go run ./examples/hello   # 预期：cef: llcppg bindings not generated
```

绑定时机成熟后，用 `llgo run ./examples/hello`，并保证 `CEF_ROOT` 与 Resources 在可执行文件旁（流程同仓库 `cefdemo`）。

XGo 语法入口可在工具链就位后新增 `main.xgo`；现阶段以 `main.go` 为准（XGo 与 Go 可混编）。
