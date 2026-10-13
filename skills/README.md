# CEF 集成用的 agent skills

这些 skill 给编程代理用，用来把 CEF 嵌进已有的原生程序，或在已有集成上改生命周期、JS 桥、呈现方式和版本。它们假设代理能读本仓库的 `docs/`，并且会打开用户选定的那份 CEF 头文件核对签名。

## 包含哪些

| Skill | 什么时候用 |
| --- | --- |
| `cef-integrate` | 工程里还没有 CEF，或只拷了库却跑不起来 |
| `cef-lifecycle` | 启动、消息循环、关窗口、退出崩溃或进程残留 |
| `cef-js-bridge` | 页面要调用原生能力，或现有绑定是 V8 扩展 |
| `cef-host-surface` | 选择或修改 Views、原生父窗口、离屏渲染 |
| `cef-upgrade` | 更换 `libcef` 版本、打开沙箱、处理 API 哈希失败 |

## 安装

把本目录下的五个子目录复制到代理会加载的 skill 路径。Cursor 会读取项目里的 `.cursor/skills/<name>/SKILL.md`，也会读取用户目录下的技能目录。复制整个子目录，不要只复制正文。

```bash
# 在你的程序仓库里
mkdir -p .cursor/skills
cp -R skills/cef-integrate skills/cef-lifecycle skills/cef-js-bridge \
      skills/cef-host-surface skills/cef-upgrade \
      /path/to/your-app/.cursor/skills/
```

## 发布包

`scripts/pack-skills.sh` 把上述目录打成 zip，作为本仓库的 release 附件。zip 内是 `skills/` 目录，解压后按上面的方式复制。版本号来自 `skills/VERSION`。

```bash
scripts/pack-skills.sh
# 生成 dist/cef-agent-skills-<version>.zip
```

文档和 skill 的核对日期写在各文件开头。分支号会过期，以官方 [Branches and Building](https://chromiumembedded.github.io/cef/branches_and_building.html) 为准。
