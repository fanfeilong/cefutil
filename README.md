# CEF 架构与集成

这个仓库说明 Chromium Embedded Framework 怎么嵌进已有的原生程序：进程和线程怎么分，窗口怎么创建和关闭，JavaScript 怎么跨进程调用宿主，二进制包怎么升级。

官方接口说明在 <https://chromiumembedded.github.io/cef/>。源码在 <https://github.com/chromiumembedded/cef>。预编译包在 <https://cef-builds.spotifycdn.com/index.html>。

本文核对日期是 **2026-10-13**。当天仍受支持的发布分支以官方 [Branches and Building](https://chromiumembedded.github.io/cef/branches_and_building.html) 为准（当时为 Beta 156、Extended 152、LTC 150）。不要把旧笔记里的 CEF 138 或 GYP 构建步骤当成现行步骤。

## 怎么读

从 [文档地图](docs/README.md) 按顺序读。已经决定要用 CEF、准备动手时，直接看 [嵌进已有程序](docs/09-integrate.md)。

1. [CEF 在产品里的位置](docs/01-place.md)
2. [架构](docs/02-architecture.md)
3. [二进制边界](docs/03-binary-boundary.md)
4. [进程和线程](docs/04-processes-and-threads.md)
5. [启动、消息循环、关闭](docs/05-lifecycle.md)
6. [应用怎么拆](docs/06-application-shape.md)
7. [JavaScript 和原生代码](docs/07-javascript.md)
8. [页面放在哪里](docs/08-hosting.md)
9. [嵌进已有程序](docs/09-integrate.md)
10. [发布与升级](docs/10-ship-and-upgrade.md)

## 给编程代理的 skill

[`skills/`](skills/README.md) 是一组给 coding agent 用的技能，用来在用户自己的程序里集成或修改 CEF。五个技能分别覆盖首次集成、生命周期、JavaScript 桥、窗口与离屏渲染、升级和沙箱。

安装时把子目录复制到程序仓库的 `.cursor/skills/`。发布包用下面的命令生成，作为本仓库 release 的附件：

```bash
scripts/pack-skills.sh
```

产物是 `dist/cef-agent-skills-<version>.zip`。

## 官方样本

集成时以二进制包里的 `cefsimple` 为最小结构，以 `cefclient` 为功能对照。离屏渲染、自定义 scheme、请求拦截都先在 `cefclient` 里看到可运行的做法，再摘进自己的进程。
