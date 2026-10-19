# cef/shim

极薄的 C（或 Objective-C）垫片，仅在下列情况启用：

1. CEF Handler 回调表无法用当前 llcppg/LLGo 直接从 XGo 填充
2. Views / 平台窗口创建必须走少量 C++ 才能接上 C API
3. 需要把 `cef_string_t` 与 Go string 的所有权规则钉死在一处

原则：shim **不写业务**；每个导出符号都有对应的 Go/XGo 声明，并在 RESEARCH 风险项里记一笔。
