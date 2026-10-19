# 路线图

对标 Tauri 2 的差距与设计见 **[tauri2-gap.md](tauri2-gap.md)**（P0/P1/P2 定义以该文为准）。

| 阶段 | 内容 | 状态 |
| --- | --- | --- |
| 0 | 调研 + 目录骨架 + IPC/配置可单测 | 完成 |
| 1a | LLGo + `cef/shim`（cefsimple_capi）Host；hello 出窗口（无 cgo） | 完成 |
| 1b | **XGo 一等**：`examples/hello/main.xgo` + `xgo build`（LLGo） | 完成 |
| 1c | llcppg 扩大 capi 绑定，收窄 shim | 待做 |
| 2 | `examples/hello` 弹出 CEF 窗口加载本地 `web/` | 完成 |
| 2b | 对标 Tauri 2 差距分析与设计文档 | 完成 |
| 3 | **P0** 接通 `api/js/xgoef.js` ↔ `ipc.Router` + 注入桥 | 待做 |
| 3b | **P0** Capabilities 强制执行 | 待做 |
| 4 | **P1** Plugin 契约 + 首批系统插件 | 待做 |
| 5 | **P1** `cmd/xgoef` create/dev/build + 资源协议 | 待做 |
| 6 | 框架 API（`runtime`/`ipc`/…）迁 `.xgo` | 待做 |
| 6b | **Classfile**：`xclass` + `*_app.xef`/`*_cmd.xef`（见 [classfile-design.md](classfile-design.md)） | 待做 |
| 7 | **P2+** Channels / 托盘菜单 / 打包更新器 | 待做 |

验收金样例：行为对齐 `src/cefdemo`（启动、本地页、关闭），API 对齐 Tauri 的 invoke 心智。
