# CEF 技能

这些技能按顺序驱动编程代理，在用户自己的机器上完成 CEF 的下载、环境、编译、二次开发工程和集成。每一步都使用官方二进制包和 CMake。代理不要阅读 `archive/`。

| 顺序 | Skill | 作用 |
| --- | --- | --- |
| 1 | `cef-setup-env` | 安装 CMake 和本平台编译器。cef-project 另外需要 Python 3.9–3.11 |
| 2 | `cef-download` | 按仍受支持的分支下载并校验 `cef_binary_*` |
| 3 | `cef-build` | 编译并运行 `cefsimple`，确认工具链 |
| 4 | `cef-create-project` | 在 cef-project 里从 cefsimple 复制出自己的目标 |
| 5 | `cef-integrate` | 把同一份二进制包嵌进已经存在的程序 |

用户要的是全新工程时走到第 4 步。用户已经有自己的程序时，第 4 步可以跳过，做完第 3 步后使用第 5 步。

## 安装

复制整个子目录，不要只复制正文。

```bash
mkdir -p .cursor/skills
cp -R skills/cef-setup-env skills/cef-download skills/cef-build \
      skills/cef-create-project skills/cef-integrate \
      /path/to/your-app/.cursor/skills/
```

本地打包：

```bash
scripts/pack-skills.sh
```

生成 `dist/cef-agent-skills-<version>.zip`。版本号在 `skills/VERSION`。

推到 `master` 且这次提交改了 `skills/VERSION` 时，GitHub Action 会用同一脚本打包，并发布 tag 为 `skills-v<version>` 的 Release。这个 tag 已经存在时，发布会失败，需要先把 `skills/VERSION` 改成新版本。只改技能正文、不改版本号，不会产生新的 Release。
