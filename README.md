⚠️ **重要提示**: 本仓库的部分内容可能已过时。CEF 项目仍在活跃开发中，请访问 [CEF 官方网站](https://bitbucket.org/chromiumembedded/cef) 获取最新信息。

**最后更新**: 2025年7月4日  
**CEF 当前版本**: 138.x  
**状态**: 部分文档需要更新 - 查看 [更新计划](update/README.md)

---

CEF全称是Chromium Embedded Framework，它是Chromium的Content API的封装库。

- CEF官网地址：https://bitbucket.org/chromiumembedded/cef
- CEF官方论坛：http://www.magpcss.org/ceforum
- CEFSharp：https://github.com/cefsharp/CefSharp
- ⚠️ ChromiumFX（已不活跃）: https://bitbucket.org/chromiumfx/chromiumfx

## CEF 文档中文翻译任务
--------------------
- [x] [CEF General Usage](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF%20General%20Usage.md)
- [x] [CEF General Usage Task List](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF%20General%20Usage%20Task%20List.md)
- [x] [CEF General Usage中文版，欢迎查阅](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF%20General%20Usage-zh-cn.md)

## CEF FAQ
-------
- [CEF 的关闭流程](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF_Close.md)
- [CEF 的C和C++封装](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF_cpp2c_annotation.md)
- [CEF 的JavaScript和C++跨语言交互](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF_JavaScript_Cpp.md)
- [Cef经典N大问题 by 扫地僧](http://blog.csdn.net/weolar/article/details/51994895)

## Chromium Documentations
-----------------
- [Chromium Generation Your Project](https://github.com/fanfeilong/cefutil/blob/master/doc/gyp.md)
- [Chromium Resoures](https://github.com/fanfeilong/cefutil/blob/master/doc/chromium_resources.md)
- [Chromium Content Register V8 Extension](https://github.com/fanfeilong/cefutil/blob/master/doc/content_register_v8_extension.md)
- [Chromium GYP 中文翻译](https://github.com/fanfeilong/cefutil/blob/master/doc/gyp.pdf)

## CEF Build
-----------------
如果要编译CEF，则需要同步编译Chromium的最新源码
- [Chromium Build Guid](https://github.com/fanfeilong/cefutil/blob/master/doc/chromium_build_guid.md)
- [Chromium for developer](http://www.chromium.org/developers)
- [windows build instructions](https://chromium.googlesource.com/chromium/src/+/master/docs/windows_build_instructions.md)
  - depot_tool的官方地址被墙(一种深深的怨恨)，可以在github 里搜索depot tool 关键字搜索相关的克隆项目，例如：
    - [depot tool](https://github.com/cybertk/depot_tools)
  - google 官方开发的vs插件，专门为chromium源码提供的。[vs-chromium](https://github.com/chromium/vs-chromium)
  - 同样的原因，depot_tool里内置的chromium的源码下载地址放在googleapi.com上，也是被墙的。下面这个地址里有chromium的所有项目的git地址：
    - https://chromium.googlesource.com/?format=JSON
    - https://chromium.googlesource.com/?format=TEXT

相关链接
---------
- [servo, the embeddable browser engine](http://blogs.s-osg.org/servo-the-embeddable-browser-engine/)

- [geckoview](https://mozilla.github.io/geckoview/)
Android offers a built-in WebView, which applications can hook into in order to display web pages within the context of their app. However, Android’s WebView is not really intended for building browsers, and hence, many advanced Web APIs are disabled. Furthermore, it is also a moving target: different phones might have different versions of WebView, all of which your app has to support.

That is where GeckoView comes in. GeckoView is:

* Full-featured: GeckoView is designed to expose the entire power of the Web to applications, and all that through a straightforward API. Think of it as harnessing the full power of Gecko (the engine that powers Firefox), while its API is WebView-like and easy to use.
* Suited for apps and browsers: GeckoView is particularly suited for building mobile browsers, but it can be embedded as a web engine component in any kind of app.
* Self-Contained: Because GeckoView is a standalone library that you bundle with your application, you can be confident that the code you test is the code that will actually run.
* Standards Compliant: Like Firefox, GeckoView offers excellent support for modern Web standards.
