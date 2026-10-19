# hello

最小 **XGo** xgoef 应用：`main.xgo` + 配置 + `greet` 命令 + 本地 `web/`。

```bash
# 工具链（一次）
./scripts/install-xgo-toolchain.sh

# 构建并运行（XGo → LLGo，无 cgo）
./scripts/build-xgoef-hello.sh
cd src/xgoef/build/hello && ./xgoef-hello
```

应用入口是 `main.xgo`（不是 `.go`）。框架底层的 `cef` / `runtime` 等包仍可为 Go，与 XGo 同模块混编；新产品代码默认写 XGo。

无 CEF 时可用 gc 测框架包：

```bash
cd src/xgoef
go test ./...
```
