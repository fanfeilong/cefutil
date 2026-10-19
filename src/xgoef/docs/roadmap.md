# 路线图

| 阶段 | 内容 | 状态 |
| --- | --- | --- |
| 0 | 调研 + 目录骨架 + IPC/配置可单测 | 完成 |
| 1a | cgo + `cef/shim`（cefsimple_capi）Host；hello 出窗口 | 完成 |
| 1b | `CEF_ROOT` + llcppg 生成最小 capi，替换 shim | 待做 |
| 2 | `examples/hello` 弹出 CEF 窗口加载本地 `web/` | 完成 |
| 3 | 接通 `api/js/xgoef.js` ↔ `ipc.Router` | 待做 |
| 4 | Plugin / capabilities | 待做 |
| 5 | `cmd/xgoef` create/dev/build + 资源拷贝复用 cefdemo | 待做 |

验收金样例：行为对齐 `src/cefdemo`（启动、本地页、关闭），API 对齐 Tauri 的 invoke 心智。
