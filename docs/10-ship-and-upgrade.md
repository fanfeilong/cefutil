# 发布与升级

核对日期：2026-10-13。依据 [Branches and Building](https://chromiumembedded.github.io/cef/branches_and_building.html)、[Sandbox Setup](https://chromiumembedded.github.io/cef/sandbox_setup.html)、[API Versioning](https://chromiumembedded.github.io/cef/api_versioning.html)。

## 分支不是一个永远的“最新稳定版”

CEF 跟随 Chromium 的偶数里程碑，大约每四周一个发布分支，遇节假日会调。官方当前把支持分成：

- **Stable：** 进入稳定通道后大约两周的广泛安全修复。
- **Extended：** 自 M152 起每四个里程碑，含稳定期在内大约八周。
- **LTC/LTS：** 更长的平台无关安全修复。文档写的是自 M160 起每十二个里程碑；M138 到 M150 是每六个。

2026-10-13 官方表上的支持分支是：Beta 8078（Chromium 156）、Extended 7977（152）、LTC 7871（150）。这三行会过期。发版前重新打开 Branches and Building，不要把本文的版本号写进构建脚本。

二进制由 Spotify 的自动构建发布，Chromium 公告和 CEF 包之间通常有一段延迟。旧分支的包仍可下载，官方不再保证能编过。

## 运行时文件

Windows 和 Linux 上，默认把 `libcef`、`icudtl.dat`、pak、V8 snapshot 和 `locales/` 放在可执行文件旁边。可以用 `CefSettings.resources_dir_path` 和 `locales_dir_path` 改资源位置。macOS 的 bundle 结构是 Chromium 规定的，Framework 和 Helper 的相对位置不能随便挪。

包内 README.txt 会标哪些文件可以删。删 locales 以外的 pak 之前先在干净机器上跑一次冷启动。

## 沙箱

沙箱按进程类型限制渲染进程等子进程的系统权限。打开它，被攻破的渲染进程更难直接读磁盘或网络。

**Windows，M138 起。** `cef_sandbox` 静态库依赖 Chromium 自带的 Clang 和 libc++，不能再链进用 MSVC 编的普通 exe。官方给 `bootstrap.exe`（窗口子系统）和 `bootstrapc.exe`（控制台）。做法是：

1. 宿主编成 DLL。窗口程序导出 `RunWinMain`，控制台程序导出 `RunConsoleMain`。签名在 `cef_sandbox_win.h`。
2. 用 `bootstrap.exe --module=client` 启动，或把 bootstrap 改名为 `client.exe`，它会加载同目录的 `client.dll`。
3. 入口收到的 `sandbox_info` 传给 `CefExecuteProcess` 和 `CefInitialize`。
4. 改过图标或版本信息之后重新签名。

不走 bootstrap 的两条路：用 Chromium 的工具链把宿主编进 CEF 的 GN 构建；或者明确关闭沙箱，两个调用都传空的 `sandbox_info`。关闭沙箱只能作为你接受风险后的选择，不要在发布配置里无意留下。

启用沙箱时，Windows 上不能设置 `browser_subprocess_path`，主进程和子进程必须是同一个可执行文件（bootstrap 模式下是同一个 bootstrap）。

**macOS，M138 起。** 分发 `libcef_sandbox.dylib`，在 Helper 的 `main` 开头、加载 CEF 框架之前完成沙箱初始化。框架用 `CefScopedLibraryLoader` 动态加载。签名、权利项和公证见官方 [Code Signing](https://chromiumembedded.github.io/cef/code_signing.html)。

**Linux。** 按发行版的沙箱能力使用包里的 `chrome-sandbox`，具体以 Sandbox Setup 指向的说明为准。

## 升级时改什么

1. 在官方表上选仍受支持的分支，下载同平台同架构的二进制包。
2. 若钉了 `CEF_API_VERSION`，确认新包仍支持这个编号。不支持就升到新包仍接受的稳定编号，并重编 wrapper 和宿主。
3. 用同一套编译参数重编 `libcef_dll_wrapper`。不要混用上一版 wrapper 和这一版 `libcef`。
4. 替换运行时文件，保留你的 `cache_path` 数据目录，除非发行说明要求清缓存。
5. 对照头文件看你调用的实验 API 是否还在。M154 起不要再依赖 `CefRegisterExtension`。
6. 重测：冷启动、子进程、关闭最后一个窗口、`onbeforeunload`、JS 往返、离屏一帧（如果用了）、沙箱开着的发布包。

跨 Chromium 主版本时，页面行为、默认开关和 GPU 路径都会变。API 哈希通过只说明结构体布局匹配。

## 共享安装

Windows 上官方提供 bootstrap 安装器，可以按应用版本从 CDN 解析并校验 CEF 二进制，供多个程序共享。见 [Installer](https://chromiumembedded.github.io/cef/installer.html)。单机随包发布时不必上这套。
