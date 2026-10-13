# 二进制边界：C API、包装层、版本

核对日期：2026-10-13。依据 [Architecture](https://chromiumembedded.github.io/cef/architecture.html)、[General Usage](https://chromiumembedded.github.io/cef/general_usage.html) 的 C++ Wrapper / Reference Counting / Strings，以及 [API Versioning](https://chromiumembedded.github.io/cef/api_versioning.html)。

## 为什么中间有一层 C

`libcef` 和宿主可能用不同的编译器运行库。字符串、对象必须由分配它的那一侧释放。CEF 因此把稳定出口做成 C 结构体和函数指针表。宿主不直接链 Chromium 的 C++ 类型。

C++ 用户链的是 `libcef_dll_wrapper`。翻译器读头文件里的 `/*--cef(...)--*/` 注释，生成 `cpptoc`（C++ 实现导出成 C）和 `ctocpp`（C 接口再包回 C++）。改包装层等于重新跑这套生成，不要手改生成文件。

纯 C 宿主可以不链包装层，直接使用 C API。官方有 `cefsimple_capi` 和 [Using the C API](https://chromiumembedded.github.io/cef/using_the_capi.html)。下面的生命周期和进程模型对两条路径相同。

## 引用计数

包装类继承 `CefBaseRefCounted`，用 `CefRefPtr` 持有。`IMPLEMENT_REFCOUNTING` 提供原子加减。创建浏览器、Handler、要跨回调活下去的对象时，交给 `CefRefPtr`，不要在栈上放一个 CEF 对象再把裸指针传进 `CreateBrowser` 之后不管。

`OnBeforeClose` 里丢掉你保存的 `CefRefPtr<CefBrowser>`。CEF 销毁浏览器对象的时机是这个回调，不是 `CloseBrowser` 返回的时候。

## 字符串

`cef_string_t` 带自己的析构函数，保证在分配它的堆上释放。C++ 侧用 `CefString`，它和 `std::string`（UTF-8）、`std::wstring` 互转。给 `CefSettings` 里的 `cef_string_t` 字段赋值时，用临时 `CefString` 包住那个字段再 `FromASCII` / `FromString`，让 CEF 自己拷贝到正确的堆。

不要把指向即将销毁的 `std::string` 缓冲区的指针写进 `cef_string_t` 又把 `dtor` 留空，除非你清楚这块内存在 CEF 读完之前一直有效。

## API 版本

M133 之前，CEF 的 ABI 没有版本。大版本一变，宿主就要重编。M133（2025 年 1 月）起引入 API 版本，让按某个稳定版本编出来的宿主，能跑在支持该版本的更新 `libcef` 上。

三档：

- **Legacy。** 2025 年之前就存在的 API，视作稳定，没有单独的版本号。
- **Versioned。** 2025 年起明确编号的稳定 API。
- **Experimental。** 标成实验的新 API。不定义 `CEF_API_VERSION` 时默认能看见它们。它们绑死当前 major/minor，没有前后 ABI 兼容，选了稳定版本后会被编译掉。

钉住稳定面的方法是在编译宿主和 `libcef_dll_wrapper` 时定义：

```text
CEF_API_VERSION=XXXYY
```

`XXX` 是 Chromium 主版本，`YY` 是该主版本内的递增号，主版本变化时从 `00` 重新计。例如 `13305` 表示 133 代上的第 5 个 API 版本。CMake 二进制包用 `-D api_version=XXXYY`。改这个宏之后，宿主和 wrapper 都要重编。

运行时，wrapper 第一次调用 `cef_api_hash` 时把这个版本交给 `libcef`。之后不能改。头文件里的结构体布局和库不一致时，进程会直接失败，而不是悄悄错位。

头文件用 `CEF_API_ADDED`、`CEF_API_REMOVED`、`CEF_API_RANGE` 把新增和删除包起来。稳定 ABI 的修改规则是：枚举在 `_NUM_VALUES` 前追加，结构体只在末尾加字段并自带 `size`，已有方法的参数和返回值不变。需要新签名时加新方法，旧槽位留着。

API 版本保证的是这一档 C/C++ 接口的二进制兼容。它不保证跨 Chromium 主版本的行为一致。换包之后要跑你自己的页面、关闭、JS 桥和渲染路径。较新的 CEF 会偶尔丢掉过旧的 API 版本，官方会在版本说明里写出来。

## 对集成的直接后果

- 业务代码只包含 `include/` 里的公共头，不包含 Chromium 的 `content/` 或 `chrome/`。
- 升级小版本、又希望少重编时，选定一个仍被该二进制支持的 `CEF_API_VERSION`，并让 wrapper 用同一个宏编译。
- 要用刚加上、还标着 experimental 的接口时，不要定义 `CEF_API_VERSION`，同时接受这份宿主只能配这一版 `libcef`。
