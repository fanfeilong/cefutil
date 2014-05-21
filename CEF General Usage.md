**这是一个翻译文档，持续更新中**

CEF General Usage(CEF3预览)
===========================

介绍
----

CEF全称Chromium Embedded Framework,是一个基于Google Chromium 的开源项目。Google Chromium项目主要是为Google Chrome应用开发的，而CEF的目标则是为第三方应用提供可嵌入浏览器支持。CEF隔离底层Chromium和Blink的复杂代码，并提供一套产品级稳定的API，发布跟踪具体Chromium版本的分支，以及二进制包。CEF的大部分特性都提供了丰富的默认实现，让使用者做尽量少的定制即可满足需求。在本文发布的时候，世界上已经有很多公司和机构采用CEF，CEF的安装量超过了100万。[CEF wikipedia]页面上有使用CEF的公司和机构的不完全的列表。CEF的典型应用场景包括：

- 嵌入一个兼容HTML5的浏览器控件到一个已经存在的本地应用。
- 创建一个轻量化的壳浏览器，用以托管主要用Web技术开发的应用。
- 有些应用有独立的绘制框架，使用CEF对Web内容做离线渲染。
- 使用CEF做自动化Web测试。

CEF3是基于Chomuim Content API多进程构架的下一代CEF，拥有下列优势：

- 改进的性能和稳定性（JavaScript和插件在一个独立的进程内执行）。
- 支持Retina显示器
- 支持WebGL和3D CSS的GPU加速
- 类似WebRTC和语音输入这样的前卫特性。
- 通过DevTools远程调试协议以及ChromeDriver2提供更好的自动化UI测试
- 更快获得当前以及未来的Web特性和标准的能力

本文档介绍CEF3开发中涉及到的一般概念。

开始
----

#### 使用二进制包

CEF3的二进制包可以在[这个页面下载](http://www.magpcss.net/cef_downloads/)。其中包含了在特定平台（Windows，Mac OS X 以及 Linux）编译特定版本CEF3所需的全部文件。不同平台拥有共同的结构：

- **cefclient** 
- **Debug**
- **include**
- **libcef_dll**
- **Release**
- **Resources**
- **tools**


