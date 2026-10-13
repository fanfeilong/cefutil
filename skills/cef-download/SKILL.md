---
name: cef-download
description: Download a current CEF binary distribution for the host platform from the Spotify automated builder. Use when the user needs libcef, a cef_binary tarball, or the CEF_VERSION string for cef-project. 也用于下载 CEF、选择分支、获取 cef_binary 压缩包。
---

# 下载 CEF 二进制包

先读 `docs/10-ship-and-upgrade.md`。不要从 `archive/` 里的文档抄版本号或下载地址。

## 选定版本

1. 打开 https://chromiumembedded.github.io/cef/branches_and_building.html ，抄下仍在支持窗口内的分支。少升级选 Extended 或 LTC，要跟新的 Web 平台选 Beta。
2. 打开 https://cef-builds.spotifycdn.com/index.html ，在该 Chromium 主版本下找与本机匹配的平台包。文件名形式是：

```text
cef_binary_<CEF_VERSION>_<platform>.tar.bz2
```

`<CEF_VERSION>` 含加号，例如 `152.0.6+g708dc14+chromium-152.0.7977.83`。`<platform>` 只能是这些值之一：`windows64`、`windowsarm64`、`windows32`、`macosarm64`、`macosx64`、`linux64`、`linuxarm64`、`linuxarm`。Beta 频道的包名在平台后多一段 `_beta`。

3. 把选中的完整 `CEF_VERSION` 写进回复和工程笔记。cef-project 仓库里默认的 `CEF_VERSION` 会过期，不能原样使用。

## 下载和校验

URL 把版本字符串里的 `+` 写成 `%2B`：

```text
https://cef-builds.spotifycdn.com/cef_binary_<CEF_VERSION>_<platform>.tar.bz2
https://cef-builds.spotifycdn.com/cef_binary_<CEF_VERSION>_<platform>.tar.bz2.sha1
```

先取 `.sha1`，再下载 tar.bz2，用 SHA1 核对。解压后目录名与去掉 `.tar.bz2` 的文件名相同。包内必须有 `CMakeLists.txt`、`include/`、`Release/` 或 `Debug/`、`Resources/`、`libcef_dll/`。

用 cef-project 时不必手写 curl。把同一个 `CEF_VERSION` 写进它顶层 `CMakeLists.txt` 的 `set(CEF_VERSION ...)`，配置阶段会调用 `cmake/DownloadCEF.cmake` 下载到 `third_party/cef/`。下载失败时用上面的 URL 对照文件名，不要改成别的镜像。

## 做完怎么判断

解压目录里能看到 `libcef` 库、`icudtl.dat`、pak 和 `locales/`。记录平台、`CEF_VERSION` 和解压路径，交给 `cef-build` 或 `cef-create-project`。
