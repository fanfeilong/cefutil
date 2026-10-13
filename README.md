# CEF 文档与技能

官方接口说明在 <https://chromiumembedded.github.io/cef/>。源码在 <https://github.com/chromiumembedded/cef>。预编译包在 <https://cef-builds.spotifycdn.com/index.html>。二次开发的起点是 <https://github.com/chromiumembedded/cef-project>。

文档和技能的核对日期是 **2026-10-13**。支持中的分支以 [Branches and Building](https://chromiumembedded.github.io/cef/branches_and_building.html) 为准。

## 有用的文档

`docs/` 说明现行 CEF 的结构，供人阅读，也供技能在动手前查阅。从 [文档地图](docs/README.md) 进入。

| 文档 | 内容 |
| --- | --- |
| [位置](docs/01-place.md) | 什么产品该用 CEF |
| [架构](docs/02-architecture.md) | Chrome bootstrap、Chrome style、Alloy style |
| [二进制边界](docs/03-binary-boundary.md) | C API、wrapper、`CEF_API_VERSION` |
| [进程和线程](docs/04-processes-and-threads.md) | 代码写在哪个进程 |
| [启动与关闭](docs/05-lifecycle.md) | 入口、消息循环、两段式关闭 |
| [应用怎么拆](docs/06-application-shape.md) | `CefApp`、`CefClient`、`CefBrowser` |
| [JavaScript](docs/07-javascript.md) | 消息路由 |
| [页面放在哪里](docs/08-hosting.md) | Views、原生父窗口、离屏渲染 |
| [集成顺序](docs/09-integrate.md) | 嵌进已有程序的检查表 |
| [发布与升级](docs/10-ship-and-upgrade.md) | 分支、沙箱、换包 |

## 技能

`skills/` 是给 coding agent 的操作技能。代理按这个顺序执行，不从归档文档里推断命令。

| 顺序 | 技能 | 作用 |
| --- | --- | --- |
| 1 | [cef-setup-env](skills/cef-setup-env/SKILL.md) | 配置 CMake、Python 和编译器 |
| 2 | [cef-download](skills/cef-download/SKILL.md) | 下载并校验当前分支的二进制包 |
| 3 | [cef-build](skills/cef-build/SKILL.md) | 编译并运行 `cefsimple` |
| 4 | [cef-create-project](skills/cef-create-project/SKILL.md) | 从 cefsimple 创建二次开发工程 |
| 5 | [cef-integrate](skills/cef-integrate/SKILL.md) | 集成到已经存在的程序 |

安装、本地打包和自动发布见 [skills/README.md](skills/README.md)。`master` 上 `skills/VERSION` 变化后，[Skills release](.github/workflows/skills-release.yml) 会发布 `skills-v<version>`，附件是 `cef-agent-skills-<version>.zip`。

全新程序走到第 4 步。已有程序在第 3 步确认能编过之后，使用第 5 步。

## 归档

[archive/](archive/README.md) 只保存旧的 PDF 和 Word。那些文件记录的是早期 GYP 和旧译文，不能用于现在的下载、编译或集成。
