**⚠️ 文档时效性声明**

本文档基于较早版本的 CEF 编写，部分内容可能已过时。当前 CEF 最新版本为 138.x（基于 Chromium 138）。

**🔥 强烈推荐优先参考 2025年现代化指南：**
- 🚀 [现代开发指南 (2025)](modern_development_guide_2025.md) - 最新工具链和环境配置
- ✨ [CEF 最新特性指南 (2025)](cef_features_2025.md) - CEF 137.x/138.x 核心特性
- ⚡ [性能优化最佳实践 (2025)](performance_optimization_2025.md) - WebGPU、WebAssembly 等
- 🔒 [安全开发指南 (2025)](security_best_practices_2025.md) - 现代安全威胁防护
- 🌐 [多平台支持指南 (2025)](multi_platform_guide_2025.md) - 跨平台开发完整指南

**官方资源：**
- [CEF 官方文档](https://bitbucket.org/chromiumembedded/cef)
- [CEF 官方论坛](https://magpcss.org/ceforum)
- [最新构建指南](https://bitbucket.org/chromiumembedded/cef/wiki/BranchesAndBuilding)

**适用版本范围：** CEF3 早期版本（可能不适用于 CEF 100+ 版本）
**最后验证时间：** 2025年7月4日

---

**📖 关于本文档**

此文档作为历史参考保留，介绍了 CEF3 的基础概念。对于现代 CEF 开发，建议使用上述 2025年现代化指南。

---

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

Each binary distribution also contains a README.txt file that describes the platform-specific distribution in greater detail and a LICENSE.txt file that contains CEF’s BSD license. When distributing an application based on CEF you should include the license text somewhere in your application’s distribution. For example, you can list it on an “About” or “Credits” page in your application’s UI, or in the documentation bundled with your application. License and credit information is also available inside of a CEF3 browser window by loading “about:license” and “about:credits” respectively.

每个二进制包包含一个README.txt文件和一个LICENSE.txt文件，README.txt用以描述平台相关的细节，而LICENSE.txt包含CEF的BSD版权说明。如果你发布了基于CEF的应用，则应该应用程序的某个地方包含该版权声明。例如，你可以在"关于”和“授权"页面列出该版权声明，或者单独一个文档包含该版权声明。“关于”和“授权”信息也可以分别在CEF浏览器的"about:license"和"about:credits"页面查看。

Applications based on CEF binary distributions can be built using standard platform build tools. This includes Visual Studio on Windows, Xcode on Mac OS X and gcc/make on Linux. The project Downloads page contains information about the OS and build tool versions required for specific binary releases. When building on Linux also pay careful attention to the listed package dependencies.

基于CEF二进制包的应用程序可以使用每个平台上的经典编译工具。包括Windows平台上的Visual Studio，Mac OSX平台上的Xcode，以及Linux平台上的gcc/make编译工具链。CEF项目的下载页面包含了这些平台上编译特定版本CEF所需的编译工具的版本信息。在Linux上编译CEF时需要特别注意依赖工具链。

See the Tutorial Wiki page for detailed instructions on how to create a simple application using the CEF3 binary distribution.

Tutorial Wiki页面有更多关于如何使用CEF3二进制包创建简单应用程序的细节。

#### Building from Source Code
#### 从源码编译

CEF can be built from source code either locally or using automated build systems like TeamCity. This requires the download of Chromium and CEF source code via either Subversion (SVN) or Git. The Chromium code base is quite large and building Chromium from source code is only recommended on moderately powerful machines with more than 4GB of RAM. Detailed instructions for building Chromium and CEF from source code are available on the BranchesAndBuilding page.

CEF可以从源码编译，用户可以使用本地编译系统或者像TeamCity这样的自动化编译系统编译。首先你需要使用svn或者git下载Chromium和CEF的源码。由于Chromium源码很大，只建议在内存大于4GB的现代机器上编译。编译Chromium和CEF的细节请参考 [BranchesAndBuilding]()页面。

#### Sample Application

#### 示例应用程序

The cefclient sample application is a complete working example of CEF integration and is included in source code form with each binary distribution. The easiest way to create a new application using CEF is to start with the cefclient application and remove the parts that you don’t need. Many of the examples in this document originate from the cefclient application.

cefclient是一个完整的CEF客户端应用程序示例，并且它的源码包含在CEF每个二进制发布包中。使用CEF创建一个新的应用程序，最简单的方法是先从cefclient应用程序开始，删除你不需要的部分。本文档中许多示例都是来源于cefclient应用程序。

#### Important Concepts

#### 重要概念

There are some important underlying concepts to developing CEF3-based applications that should be understood before proceeding.

在开发基于CEF3的应用程序前，有一些重要的基础概念应该被理解。

##### C++ Wrapper

####  C++ 封装

The libcef shared library exports a C API that isolates the user from the CEF runtime and code base. The libcef_dll_wrapper project, which is distributed in source code form as part of the binary release, wraps this exported C API in a C++ API that is then linked into the client application. The code for this C/C++ API translation layer is automatically generated by the translator tool. Direct usage of the C API is described on the UsingTheCAPI page.

libcef 动态链接库导出 C API 使得使用者不用关心CEF运行库和基础代码。libcef_dll_wrapper 工程把 C API 封装成 C++ API同时包含在客户端应用程序工程中，与cefclient一样，源代码作为CEF二进制发布包的一部份共同发布。C/C++ API的转换层代码是由转换工具自动生成。UsingTheCAPI 页面描述了如何使用C API。

##### Processes
##### 进程

CEF3 runs using multiple processes. The main process which handles window creation, painting and network access is called the “browser” process. This is generally the same process as the host application and the majority of the application logic will run in the browser process. Blink rendering and JavaScript execution occur in a separate “render” process. Some application logic, such as JavaScript bindings and DOM access, will also run in the render process. The default process model will spawn a new render process for each unique origin (scheme + domain). Other processes will be spawned as needed, such as “plugin” processes to handle plugins like Flash and “gpu” processes to handle accelerated compositing.

CEF3是多进程架构的。"browser"被定义为主进程，负责窗口管理，界面绘制和网络交互。Blink的渲染和Js的执行被放在一个独立的"render"
进程中；除此之外，render进程还负责Js Binding和对Dom节点的访问。
默认的进程模型中，会为每个标签页创建一个新的"render"进程。其他进程按需创建，象管理插件的进程和处理合成加速的进程。

By default the main application executable will be spawned multiple times to represent separate processes. This is handled via command-line flags that are passed into the CefExecuteProcess function. If the main application executable is large, takes a long time to load, or is otherwise unsuitable for non-browser processes the host can use a separate executable for those other processes. This can be configured via the CefSettings.browser_subprocess_path variable. See the “Application Structure” section for more information.

默认情况下，主应用程序会被多次启动运行各自独立的进程。这是通过传递不同的命令行参数给CefExecuteProcess函数做到的。如果主应用程序很大，加载时间比较长，或者不能在非浏览器进程里使用，则宿主程序可使用独立的可执行文件去运行这些进程。这可以通过配置CefSettings.browser_subprocess_path变量做到。更多细节请参考`Application Structure`一节。

The separate processes spawned by CEF3 communicate using Inter-Process Communication (IPC). Application logic implemented in the browser and render processes can communicate by sending asynchronous messages back and forth. JavaScriptIntegration in the render process can expose asynchronous APIs that are handled in the browser process. See the “Inter-Process Communication” section for more information.

CEF3的进程之间可以通过IPC进行通信。"Browser"和"Render"进程可以通过发送异步消息进行双向通信。甚至在Render进程可以注册在Browser进程响应的异步JavaScript API。更多细节，请参考`Inter-Process Communication`一节。

CEF3 supports a single-process run mode for debugging purposes via the "--single-process" command-line flag. Platform-specific debugging tips are also available for Windows, Mac OS X and Linux.

通过设置命令行的"--single-process"，CEF3就可以支持用于调试目的的单进程运行模型。支持的平台为：Windows，Mac OS X 和Linux。

##### Threads

##### 线程

Each process in CEF3 runs multiple threads. For a complete list of threads see the cef_thread_id_t enumeration. The browser process for example contains the following commonly-referenced threads:

在CEF3中，每个进程都会运行多个线程。完整的线程类型表请参照cef_thread_id_t。例如，在browser进程中包含如下主要的线程：

- **TID_UI** thread is the main thread in the browser. This will be the same as the main application thread if CefInitialize() is called with a CefSettings.multi_threaded_message_loop value of false.
- **TID_IO** thread is used to process IPC and network messages.
- **TID_FILE** thread is used to interact with the file system. 

- **TID_UI** 线程是浏览器的主线程。如果应用程序在调用调用CefInitialize()时，传递CefSettings.multi_threaded_message_loop=false，这个线程也是应用程序的主线程。
- **TID_IO** 线程主要负责处理IPC消息以及网络通信。
- **TID_FILE** 线程负责与文件系统交互。 

Due to the multi-threaded nature of CEF it’s important to use locking or message passing to protect data members from access on multiple threads. The IMPLEMENT_LOCKING macro provides Lock() and Unlock() methods and an AutoLock scoped object for synchronizing access to blocks of code. The CefPostTask family of functions support easy asynchronous message passing between threads. See the “Posting Tasks” section for more information.

由于CEF采用多线程架构，有必要使用锁和闭包来保证在多不同线程安全的传递数据。IMPLEMENT_LOCKING定义提供了Lock()和Unlock()方法以及AutoLock对象来保证不同代码块同步访问数据。CefPostTask函数组支持简易的线程间异步消息传递。更多信息，请参考"Posting Tasks"章节。

The current thread can be verified using the CefCurrentlyOn() function. The cefclient application uses the following defines to verify that methods are executed on the expected thread:

判断当前工作线程可以通过使用CefCurrentlyOn()方法，cefclient工程使用下面的定义来确保方法在期望的线程中被执行。

```
#define REQUIRE_UI_THREAD()   ASSERT(CefCurrentlyOn(TID_UI));
#define REQUIRE_IO_THREAD()   ASSERT(CefCurrentlyOn(TID_IO));
#define REQUIRE_FILE_THREAD() ASSERT(CefCurrentlyOn(TID_FILE));
```

##### Reference Counting
##### 引用计数

All framework classes implement the CefBase interface and all instance pointers are handled using the CefRefPtr smart pointer implementation that automatically handles reference counting via calls to AddRef() and Release().The easiest way to implement these classes is as follows:

所有的框架类从CefBase继承，实例指针由CefRefPtr管理，CefRefPtr通过调用AddRef()和Release()方法自动管理引用计数。框架类的实现方式如下：

```
class MyClass : public CefBase {
 public:
  // Various class methods here...

 private:
  // Various class members here...

  IMPLEMENT_REFCOUNTING(MyClass);  // Provides atomic refcounting implementation.
};

// References a MyClass instance
CefRefPtr<MyClass> my_class = new MyClass();
```

##### Strings

##### 字符串

CEF defines its own data structure for representing strings. This is for a few different reasons:

CEF为字符串定义了自己的数据结构。下面是这样做的理由：

- The libcef library and the host application may use different runtimes for managing heap memory. All objects, including strings, need to be freed using the same runtime that allocated the memory.
- The libcef library can be compiled to support different underlying string types (UTF8, UTF16 or wide). The default is UTF16 but it can be changed by modifying the defines in cef_string.h and recompiling CEF. When choosing the wide string type keep in mind that the size will vary depending on the platform. 

- libcef包和宿主程序可能使用不同的运行时，对堆管理的方式也不同。所有的对象，包括字符串，需要确保和申请堆内存使用相同的运行时环境。
- libcef包可以编译为支持不同的字符串类型(UTF8，UTF16以及WIDE)。默认采用的是UTF16，默认字符集可以通过更改cef_string.h文件中的定义，然后重新编译来修改。当使用宽字节集的时候，切记字符的长度由当前使用的平台决定。

For UTF16 the string structure looks like this:

UTF16字符串结构体示例如下：

```
typedef struct _cef_string_utf16_t {
  char16* str;  // Pointer to the string
  size_t length;  // String length
  void (*dtor)(char16* str);  // Destructor for freeing the string on the correct heap
} cef_string_utf16_t;
```

The selected string type is then typedef’d to the generic type:

通过typedef来设置常用的字符编码。

```
typedef char16 cef_char_t;
typedef cef_string_utf16_t cef_string_t;
```

CEF provides a number of C API functions for operating on the CEF string types (mapped via #defines to the type-specific functions). For example:

CEF提供了一批C语言的方法来操作字符串(通过#define的方式来适应不同的字符编码)

- **cef_string_set** will assign a string value to the structure with or without copying the value.
- **cef_string_clear** will clear the string value.
- **cef_string_cmp** will compare two string values. 

- **cef_string_set** 对制定的字符串变量赋值(支持深拷贝或浅拷贝)。
- **cef_string_clear** 清空字符串。
- **cef_string_cmp** 比较两个字符串. 

CEF also provides functions for converting between all supported string types (ASCII, UTF8, UTF16 and wide). See the cef_string.h and cef_string_types.h headers for the complete list of functions.

CEF也提供了所有的字符编码的字符串格式之间相互转换的方法。具体函数列表请查阅cef_string.h和cef_string_types.h文件。

Usage of CEF strings in C++ is simplified by the CefString class. CefString provides automatic conversion to and from std::string (UTF8) and std::wstring (wide) types. It can also be used to wrap an existing cef_string_t structure for assignment purposes.

在C++中，通常使用CefString类来管理CEF的字符串。CefString支持与std::string(UTF8)、std::wstring(wide)类型的相互转换。也可以用来包裹一个cef_string_t结构来对其进行赋值。

Assignment to and from std::string:

和std::string的相互转换：

```
std::string str = “Some UTF8 string”;

// Equivalent ways of assigning |str| to |cef_str|. Conversion from UTF8 will occur if necessary.
CefString cef_str(str);
cef_str = str;
cef_str.FromString(str);

// Equivalent ways of assigning |cef_str| to |str|. Conversion to UTF8 will occur if necessary.
str = cef_str;
str = cef_str.ToString();
```

Assignment to and from std::wstring:

和std::wstring的相互转换：

```
std::wstring str = “Some wide string”;

// Equivalent ways of assigning |str| to |cef_str|. Conversion from wide will occur if necessary.
CefString cef_str(str);
cef_str = str;
cef_str.FromWString(str);

// Equivalent ways of assigning |cef_str| to |str|. Conversion to wide will occur if necessary.
str = cef_str;
str = cef_str.ToWString();
```

Use the FromASCII() method if you know that the format is ASCII:

如果是ASCII编码，使用FromASCII进行赋值：

```
const char* cstr = “Some ASCII string”;
CefString cef_str;
cef_str.FromASCII(cstr);
```

Some structures like CefSettings have cef_string_t members. CefString can be used to simplify the assignment to those members:

一些结构体(比如CefSettings)含有cef_string_t类型的成员，CefString支持直接赋值给这些成员。

```
CefSettings settings;
const char* path = “/path/to/log.txt”;

// Equivalent assignments.
CefString(&settings.log_file).FromASCII(path);
cef_string_from_ascii(path, strlen(path), &settings.log_file);
```

##### Command Line Arguments
##### 命令行参数
Many features in CEF3 and Chromium can be configured using command line arguments. These arguments take the form "--some-argument[=optional-param]" and are passed into CEF via CefExecuteProcess() and the CefMainArgs structure (see the “Application Structure” section below). To disable processing of arguments from the command line set CefSettings.command_line_args_disabled to true before passing the CefSettings structure into CefInitialize(). To specify command line arguments inside the host application implement the CefApp::OnBeforeCommandLineProcessing() method. See comments in client_switches.cpp for more information on how to discover supported command line switches.

在CEF3和Chromium中许多特性可以使用命令行参数进行配置。这些参数采用"--some-argument[=optional-param]"形式，并通过CefExecuteProcess()和CefMainArgs结构（参考下面的"应用资源布局"章节）传递给CEF。在传递CefSettings结构给CefInitialize()之前，我们可以设置CefSettings.command_line_args_disabled为false来禁用对命令行参数的处理。如果想指定命令行参数传入主应用程序，实现CefApp::OnBeforeCommandLineProcessing()方法。更多关于如何查找已支持的命令行选项的信息，请查看client_switches.cpp文件的注释。

#### Application Layout
#### 应用程序文件结构

Application layout can differ significantly depending on the platform. For example, on Mac OS X your application layout must follow a specific app bundle structure. Windows and Linux are more flexible, allowing you to customize the location where CEF libraries and resources are stored. For a complete working example of the required layout you can download a client archive from the project Downloads page. Some files are optional and some are required as detailed in the README.txt file for each platform.

应用资源布局依赖于平台，有很大的不同。比如，在Mac OS X上，你的资源布局必须遵循特定的app bundles结构；Window与Linux则更灵活，允许你定制CEF库文件与资源文件所在的位置。为了获取到特定地、可以正常工作的示例，你可以从工程的下载页面，下载到一个client压缩包，每个平台对应的README.txt文件详细说明了哪些文件是可选的，哪些文件是必须的。

##### Windows
##### Windows操作系统

On Windows the default layout places the libcef library and related resources next to the application executable. The directory structure looks like this:

在Windows平台上，默认的资源布局将libcef库文件、相关资源与可执行文件放置在同级目录，文件夹结构如下类似：

```
Application/
    cefclient.exe  <= cefclient application executable 
    libcef.dll <= main CEF library 
    icudt.dll <= ICU unicode support library 
    ffmpegsumo.dll <= HTML5 audio/video support library 
    libEGL.dll, libGLESv2.dll, … <= accelerated compositing support libraries 
    cef.pak, devtools_resources.pak <= non-localized resources and strings 
    locales/
        en-US.pak, … <= locale-specific resources and strings 
```

The location of the CEF libraries and resource files can be customized using the CefSettings structure (see the README.txt file or “CefSettings” section for details). The cefclient application on Windows compiles in resources via the BINARY resource type in cefclient.rc but an application could just as easily load resources from the local file system.

使用结构体CefSettings可以定制CEF库文件、资源文件的位置（查看README.txt文件或者本文中CefSettings部分获取更详细的信息）。虽然在Windows平台上，cefclient项目将资源文件以二进制形式编译进cefclient.rc文件，但是改为从文件系统加载资源也很容易。


##### Linux
##### Linux操作系统

On Linux the default layout places the libcef library and related resources next to the application executable. Note however that there’s a discrepancy between where libcef.so is located in the client distribution and where it’s located in the binary distribution that you build yourself. The location depends on how the linker rpath value is set when building the application executable. For example, a value of “-Wl,-rpath,.” (“.” meaning the current directory) will allow you to place libcef.so next to the application executable. The path to libcef.so can also be specified using the LD_LIBRARY_PATH environment variable.

在Linux平台上，默认的资源布局将libcef库文件、相关资源与可执行文件放置在同级目录。注意：在你编译的版本与发行版本应用程序中，libcef.so的位置是有差异的，此文件的位置取决于编译可执行程序时，编译器rpath的值。比如，编译选项为“-Wl,-rpath,.”（“.”意思是当前文件夹），这样libcef.so与可执行文件处于同级目录。libcef.so文件的路径可以通过环境变量中的“`LD_LIBRARY_PATH`”指定。

```
Application/
    cefclient  <= cefclient application executable 
    libcef.so <= main CEF library 
    ffmpegsumo.so <-- HTML5 audio/video support library 
    cef.pak, devtools_resources.pak <= non-localized resources and strings 
    locales/
        en-US.pak, … <= locale-specific resources and strings 
    files/
        binding.html, … <= cefclient application resources 
```

The location of the CEF libraries and resource files can be customized using the CefSettings structure (see the README.txt file of “CefSettings” section for details).

使用结构体CefSettings可以定制CEF库文件、资源文件（查看README.txt文件或者本文中CefSettings部分获取更详细的信息）。

##### Mac OS X
##### Mac X平台

The application (app bundle) layout on Mac OS X is mandated by the Chromium implementation and consequently is not very flexible. The directory structure looks like this:

在Mac X平台上，app bundles委托给了Chromium实现，因此不是很灵活。文件夹结构如下类似：

```
cefclient.app/
    Contents/
        Frameworks/
            Chromium Embedded Framework.framework/
                Libraries/
                    ffmpegsumo.so <= HTML5 audio/video support library 
                    libcef.dylib <= main CEF library
                Resources/
                    cef.pak, devtools_resources.pak <= non-localized resources and strings
                    *.png, *.tiff <= Blink image and cursor resources 
                    en.lproj/, … <= locale-specific resources and strings 
            libplugin_carbon_interpose.dylib <= plugin support library
            cefclient Helper.app/
                Contents/
                    Info.plist
                    MacOS/
                        cefclient Helper <= helper executable 
                    Pkginfo
            cefclient Helper EH.app/
                Contents/
                    Info.plist
                    MacOS/
                        cefclient Helper EH <= helper executable 
                    Pkginfo
            cefclient Helper NP.app/
                Contents/
                    Info.plist
                    MacOS/
                        cefclient Helper NP <= helper executable 
                    Pkginfo
        Info.plist
        MacOS/
            cefclient <= cefclient application executable 
        Pkginfo
        Resources/
            binding.html, … <= cefclient application resources 
```

The "Chromium Embedded Framework.framework" is an unversioned framework that contains all CEF binaries and resources. Executables (cefclient, cefclient Helper, etc) are linked to libcef.dylib using install_name_tool and a path relative to @executable_path.

列表中的“Chromium Embedded Framework.framework”，这个未受版本管控的框架包含了所有的CEF库文件、资源文件。使用install_name_tool与@executable_path，将cefclient，cefclient helper等可执行文件，连接到了libcef.dylib上。

The "cefclient Helper" apps are used for executing separate processes (renderer, plugin, etc) with different characteristics. They need to have separate app bundles and Info.plist files so that, among other things, they don't show dock icons. The "EH" helper, which is used when launching plugin processes, has the MH_NO_HEAP_EXECUTION bit cleared to allow an executable heap. The "NP" helper, which is used when launching NaCl plugin processes only, has the MH_PIE bit cleared to disable ASLR. This is set up as part of the build process using scripts from the tools/ directory. Examine the Xcode project included with the binary distribution or the originating cefclient.gyp file for a better idea of the script dependencies.

应用程序cefclient helper用来执行不同特点、独立的进程（renderer，plugin等），这些进程需要独立的资源布局与Info.plist等文件，它们没有显示停靠图标。用来启动插件进程的EH Helper清除了MH_NO_HEAP_EXECUTION标志位，这样就允许一个可执行堆。只能用来启动NaCL插件进程的NP Helper，清除了MH_PIE标志位，这样就禁用了ASLR。这些都是tools文件夹下面，用来构建进程脚本的一部分。为了理清脚本的依赖关系，更好的做法是检查发行版本中的Xcode工程或者原始文件cefclient.gyp。

#### Application Structure
#### 应用程序结构

Every CEF3 application has the same general structure.

- Provide an entry-point function that initializes CEF and runs either sub-process executable logic or the CEF message loop.
- Provide an implementation of CefApp to handle process-specific callbacks.
- Provide an implementation of CefClient to handle browser-instance-specific callbacks.
- Call CefBrowserHost::CreateBrowser() to create a browser instance and manage the browser life span using CefLifeSpanHandler. 

每个CEF3应用程序都是相同的结构

- 提供入口函数，用于初始化CEF、运行子进程执行逻辑或者CEF消息循环。
- 提供CefApp实现，用于处理进程相关的回调。
- 提供CefClient实现，用于处理browser实例相关的回调
- 执行CefBrowserHost::CreateBrowser()创建一个browser实例，使用CefLifeSpanHandler管理browser对象生命周期。

##### Entry-Point Function
##### 入口函数

As described in the “Processes” section a CEF3 application will run multiple processes. The processes can all use the same executable or a separate executable can be specified for the sub-processes. Execution of the process begins in the entry-point function. Complete platform-specific examples for Windows, Linux and Mac OS-X are available in cefclient_win.cc, cefclient_gtk.cc and cefclient_mac.mm respectively.

像本文中进程章节描述的那样，一个CEF3应用程序会运行多个进程，这些进程能够使用同一个执行器或者为子进程定制的、单独的执行器。进程的执行从入口函数开始，示例cefclient_win.cc、cefclient_gtk.cc、cefclient_mac.mm分别对应Windows、Linux和Mac OS-X平台下的实现。

When launching sub-processes CEF will specify configuration information using the command-line that must be passed into the CefExecuteProcess function via the CefMainArgs structure. The definition of CefMainArgs is platform-specific. On Linux and Mac OS X it accepts the argc and argv values which are passed into the main() function.

当执行子进程是，CEF将使用命令行参数指定配置信息，这些命令行参数必须通过CefMainArgs结构体传入到CefExecuteProcess函数。CefMainArgs的定义与平台相关，在Linux、Mac OS X平台下，它接收main函数传入的argc和argv参数值。

```
CefMainArgs main_args(argc, argv);
```

On Windows it accepts the instance handle (HINSTANCE) which is passed into the wWinMain() function. The instance handle is also retrievable via GetModuleHandle(NULL).

在Windows平台下，它接收wWinMain函数传入的参数：实例句柄（HINSTANCE），这个实例能够通过函数GetModuleHandle(NULL)获取。

```
CefMainArgs main_args(hInstance);
```

##### Single Executable
##### 单一执行体

When running as a single executable the entry-point function is required to differentiate between the different process types. The single executable structure is supported on Windows and Linux but not on Mac OS X.

当以单一执行体运行时，根据不同的进程类型，入口函数有差异。Windows、Linux平台支持单一执行体架构，Mac OS X平台则不行。


```
int main(int argc, char* argv[]) {
  // Structure for passing command-line arguments.
  // The definition of this structure is platform-specific.
  CefMainArgs main_args(argc, argv);

  // Optional implementation of the CefApp interface.
  CefRefPtr<MyApp> app(new MyApp);

  // Execute the sub-process logic, if any. This will either return immediately for the browser
  // process or block until the sub-process should exit.
  int exit_code = CefExecuteProcess(main_args, app.get());
  if (exit_code >= 0) {
    // The sub-process terminated, exit now.
    return exit_code;
  }

  // Populate this structure to customize CEF behavior.
  CefSettings settings;

  // Initialize CEF in the main process.
  CefInitialize(main_args, settings, app.get());

  // Run the CEF message loop. This will block until CefQuitMessageLoop() is called.
  CefRunMessageLoop();

  // Shut down CEF.
  CefShutdown();

  return 0;
}
```

##### Separate Sub-Process Executable
##### 独立的子进程执行体

When using a separate sub-process executable you need two separate executable projects and two separate entry-point functions.

当使用独立的子进程执行体时，你需要2个分开的可执行工程和2个分开的入口函数。

Main application entry-point function:

主程序的入口函数：

```
// Program entry-point function.
// 程序入口函数
int main(int argc, char* argv[]) {
  // Structure for passing command-line arguments.
  // The definition of this structure is platform-specific.
  // 传递命令行参数的结构体。
  // 这个结构体的定义与平台相关。
  CefMainArgs main_args(argc, argv);

  // Optional implementation of the CefApp interface.
  // 可选择性地实现CefApp接口
  CefRefPtr<MyApp> app(new MyApp);

  // Populate this structure to customize CEF behavior.
  // 填充这个结构体，用于定制CEF的行为。
  CefSettings settings;

  // Specify the path for the sub-process executable.
  // 指定子进程的执行路径
  CefString(&settings.browser_subprocess_path).FromASCII(“/path/to/subprocess”);

  // Initialize CEF in the main process.
  // 在主进程中初始化CEF 
  CefInitialize(main_args, settings, app.get());

  // Run the CEF message loop. This will block until CefQuitMessageLoop() is called.
  // 执行消息循环，此时会堵塞，直到CefQuitMessageLoop()函数被调用。
  CefRunMessageLoop();

  // Shut down CEF.
  // 关闭CEF
  CefShutdown();

  return 0;
}
```

Sub-process application entry-point function:
子进程程序的入口函数：

```
// Program entry-point function.
// 程序入口函数
int main(int argc, char* argv[]) {
  // Structure for passing command-line arguments.
  // The definition of this structure is platform-specific.
  // 传递命令行参数的结构体。
  // 这个结构体的定义与平台相关。
  CefMainArgs main_args(argc, argv);

  // Optional implementation of the CefApp interface.
  // 可选择性地实现CefApp接口
  CefRefPtr<MyApp> app(new MyApp);

  // Execute the sub-process logic. This will block until the sub-process should exit.
  // 执行子进程逻辑，此时会堵塞直到子进程退出。
  return CefExecuteProcess(main_args, app.get());
}
```

##### Message Loop Integration
##### 集成消息循环

CEF can also integrate with an existing application message loop instead of running its own message loop. There are two ways to do this.

CEF可以不用它自己提供的消息循环，而与已经存在的程序中消息环境集成在一起，有两种方式可以做到：

1. Call CefDoMessageLoopWork() on a regular basis instead of calling CefRunMessageLoop(). Each call to CefDoMessageLoopWork() will perform a single iteration of the CEF message loop. Caution should be used with this approach. Calling the method too infrequently will starve the CEF message loop and negatively impact browser performance. Calling the method too frequently will negatively impact CPU usage.
2. Set CefSettings.multi_threaded_message_loop = true (Windows only). This will cause CEF to run the browser UI thread on a separate thread from the main application thread. With this approach neither CefDoMessageLoopWork() nor CefRunMessageLoop() need to be called. CefInitialize() and CefShutdown() should still be called on the main application thread. You will need to provide your own mechanism for communicating with the main application thread (see for example the message window usage in cefclient_win.cpp). You can test this mode in cefclient on Windows by running with the “--multi-threaded-message-loop” command-line flag. 

1. 周期性执行CefDoMessageLoopWork()函数，替代调用CefRunMessageLoop()。CefDoMessageLoopWork()的每一次调用，都将执行一次CEF消息循环的单次迭代。需要注意的是，此方法调用次数太少时，CEF消息循环会饿死，将极大的影响browser的性能，调用次数太频繁又将影响CPU使用率。

2. 设置CefSettings.multi_threaded_message_loop=true（Windows平台下有效），这个设置项将导致CEF运行browser UI运行在单独的线程上，而不是在主线程上，这种场景下CefDoMessageLoopWork()或者CefRunMessageLoop()都不需要调用，CefInitialze()、CefShutdown()仍然在主线程中调用。你需要提供主程序线程通信的机制（查看cefclient_win.cpp中提供的消息窗口实例）。在Windows平台下，你可以通过命令行参数“--multi-threaded-message-loop”测试上述消息模型。

##### CefSettings
##### CefSettings

The CefSettings structure allows configuration of application-wide CEF settings. Some commonly configured members include:

CefSettings结构体允许定义全局的CEF配置，经常用到的配置项如下：

- **single_process** Set to true to use a single process for the browser and renderer. Also configurable using the "single-process" command-line switch. See the “Processes” section for more information.
- **browser_subprocess_path** The path to a separate executable that will be launched for sub-processes. See the “Separate Sub-Process Executable” section for more information.
- **multi_threaded_message_loop** Set to true to have the browser process message loop run in a separate thread. See the “Message Loop Integration” section for more information.
- **command_line_args_disabled** Set to true to disable configuration of browser process features using standard CEF and Chromium command-line arguments. See the “Command Line Arguments” section for more information.
- **cache_path** The location where cache data will be stored on disk. If empty an in-memory cache will be used for some features and a temporary disk cache will be used for others. HTML5 databases such as localStorage will only persist across sessions if a cache path is specified.
- **locale** The locale string that will be passed to Blink. If empty the default locale of "en-US" will be used. This value is ignored on Linux where locale is determined using environment variable parsing with the precedence order: LANGUAGE, LC_ALL, LC_MESSAGES and LANG. Also configurable using the "lang" command-line switch.
- **log_file** The directory and file name to use for the debug log. If empty, the default name of "debug.log" will be used and the file will be written to the application directory. Also configurable using the "log-file" command-line switch.
- **log_severity** The log severity. Only messages of this severity level or higher will be logged. Also configurable using the "log-severity" command-line switch with a value of "verbose", "info", "warning", "error", "error-report" or "disable".
- **resources_dir_path** The fully qualified path for the resources directory. If this value is empty the cef.pak and/or devtools_resources.pak files must be located in the module directory on Windows/Linux or the app bundle Resources directory on Mac OS X. Also configurable using the "resources-dir-path" command-line switch.
- **locales_dir_path** The fully qualified path for the locales directory. If this value is empty the locales directory must be located in the module directory. This value is ignored on Mac OS X where pack files are always loaded from the app bundle Resources directory. Also configurable using the "locales-dir-path" command-line switch.
- **remote_debugging_port** Set to a value between 1024 and 65535 to enable remote debugging on the specified port. For example, if 8080 is specified the remote debugging URL will be http://localhost:8080. CEF can be remotely debugged from any CEF or Chrome browser window. Also configurable using the "remote-debugging-port" command-line switch. 

- **single_process** 设置为true时，browser和renderer使用一个进程。此项也可以通过命令行参数“single-process”配置。查看本文中“进程”章节获取更多的信息。
- **browser_subprocess_path** 设置用于启动子进程单独执行器的路径。参考本文中“单独子进程执行器”章节获取更多的信息。
- **cache_path** 设置磁盘上用于存放缓存数据的位置。如果此项为空，某些功能将使用内存缓存，多数功能将使用临时的磁盘缓存。形如本地存储的HTML5数据库只能在设置了缓存路径才能跨session存储。
- **locale** 此设置项将传递给Blink。如果此项为空，将使用默认值“en-US”。在Linux平台下此项被忽略，使用环境变量中的值，解析的依次顺序为：LANGUAE，LC_ALL，LC_MESSAGES和LANG。此项也可以通过命令行参数“lang”配置。
- **log_file** 此项设置的文件夹和文件名将用于输出debug日志。如果此项为空，默认的日志文件名为debug.log，位于应用程序所在的目录。此项也可以通过命令参数“log-file”配置。
- **log_severity** 此项设置日志级别。只有此等级、或者比此等级高的日志的才会被记录。此项可以通过命令行参数“log-severity”配置，可以设置的值为“verbose”，“info”，“warning”，“error”，“error-report”，“disable”
- **resources_dir_path** 此项设置资源文件夹的位置。如果此项为空，Windows平台下cef.pak、Linux平台下devtools_resourcs.pak、Mac OS X下的app bundle Resources目录必须位于组件目录。此项也可以通过命令行参数“resource-dir-path”配置。
- **locales_dir_path** 此项设置locale文件夹位置。如果此项为空，locale文件夹必须位于组件目录，在Mac OS X平台下此项被忽略，pak文件从app bundle Resources目录。此项也可以通过命令行参数“locales-dir-path”配置。
- **remote_debugging_port** 此项可以设置1024-65535之间的值，用于在指定端口开启远程调试。例如，如果设置的值为8080，远程调试的URL为http://localhost:8080。CEF或者Chrome浏览器能够调试CEF。此项也可以通过命令行参数“remote-debugging-port”配置。

##### CefBrowser and CefFrame
##### CefBrowser和CefFrame

The CefBrowser and CefFrame objects are used for sending commands to the browser and for retrieving state information in callback methods. Each CefBrowser object will have a single main CefFrame object representing the top-level frame and zero or more CefFrame objects representing sub-frames. For example, a browser that loads two iframes will have three CefFrame objects (the top-level frame and the two iframes).

CefBrowser和CefFrame对象被用来发送命令给浏览器以及在回调函数里获取状态信息。每个CefBrowser对象包含一个主CefFrame对象，主CefFrame对象代表页面的顶层frame；同时每个CefBrowser对象可以包含零个或多个的CefFrame对象，分别代表不同的子Frame。例如，一个浏览器加载了两个iframe，则该CefBrowser对象拥有三个CefFrame对象（顶层frame和两个iframe）。

To load a URL in the browser main frame:
下面的代码在浏览器的主frame里加载一个URL：
```
browser->GetMainFrame()->LoadURL(some_url);
```

To navigate the browser backwards:
下面的代码执行浏览器的回退操作：
```
browser->GoBack();
```

To retrieve the main frame HTML contents:
下面的代码从主frame里获取HTML内容：
```
// Implementation of the CefStringVisitor interface.
class Visitor : public CefStringVisitor {
 public:
  Visitor() {}

  // Called asynchronously when the HTML contents are available.
  virtual void Visit(const CefString& string) OVERRIDE {
    // Do something with |string|...
  }

  IMPLEMENT_REFCOUNTING(Visitor);
};

browser->GetMainFrame()->GetSource(new Visitor());
```

CefBrowser and CefFrame objects exist in both the browser process and the render process. Host behavior can be controlled in the browser process via the CefBrowser::GetHost() method. For example, the native handle for a windowed browser can be retrieved as follows:

CefBrowser和CefFrame对象在Browser进程和Render进程都有对等的代理对象。在Browser进程里，Host（宿主）行为控制可以通过CefBrowser::GetHost()方法控制。例如，浏览器窗口的原生句柄可以用下面的代码获取：

```
// CefWindowHandle is defined as HWND on Windows, NSView* on Mac OS X
// and GtkWidget* on Linux.
CefWindowHandle window_handle = browser->GetHost()->GetWindowHandle();
```

Other methods are available for history navigation, loading of strings and requests, sending edit commands, retrieving text/html contents, and more. See the documentation for the complete list of supported methods.

其他方法包括历史导航，加载字符串和请求，发送编辑命令，提取text/html内容等。请参考支持函数相关的文档或者CefBrowser的头文件注释。

##### CefApp

The CefApp interface provides access to process-specific callbacks. Important callbacks include:

CefApp接口提供了不同进程的可定制回调函数。毕竟重要的回调函数如下：

- **OnBeforeCommandLineProcessing** which provides the opportunity to programmatically set command-line arguments. See the “Command Line Arguments” section for more information.
- **OnRegisterCustomSchemes** which provides an opportunity to register custom schemes. See the “”Request Handling” section for more information.
- **GetBrowserProcessHandler** which returns the handler for functionality specific to the browser process including the OnContextInitialized() method.
- **GetRenderProcessHandler** which returns the handler for functionality specific to the render process. This includes JavaScript-related callbacks and process messages. See JavaScriptIntegration and the “Inter-Process Communication” section for more information. 

- **OnBeforeCommandLineProcessing**  提供了以编程方式设置命令行参数的机会，更多细节，请参考`Command Line Arguments`一节。
- **OnRegisterCustomSchemes**  提供了注册自定义schemes的机会，更多细节，请参考`Request Handling`一节。
- **GetBrowserProcessHandler** 返回定制Browser进程的Handler，该Handler包括了诸如OnContextInitialized的回调。
- **GetRenderProcessHandler** 返回定制Render进程的Handler，该Handler包含了JavaScript相关的一些回调以及消息处理的回调。更多细节，请参考`JavascriptIntegration`和`Inter-Process Communication`两节。

Example CefApp implementation:

CefApp子类的例子：

```
// MyApp implements CefApp and the process-specific interfaces.
class MyApp : public CefApp,
              public CefBrowserProcessHandler,
              public CefRenderProcessHandler {
 public:
  MyApp() {}

  // CefApp methods. Important to return |this| for the handler callbacks.
  virtual void OnBeforeCommandLineProcessing(
      const CefString& process_type,
      CefRefPtr<CefCommandLine> command_line) {
    // Programmatically configure command-line arguments...
  }
  virtual void OnRegisterCustomSchemes(
      CefRefPtr<CefSchemeRegistrar> registrar) OVERRIDE {
    // Register custom schemes...
  }
  virtual CefRefPtr<CefBrowserProcessHandler> GetBrowserProcessHandler()
      OVERRIDE { return this; }
  virtual CefRefPtr<CefRenderProcessHandler> GetRenderProcessHandler()
      OVERRIDE { return this; }

  // CefBrowserProcessHandler methods.
  virtual void OnContextInitialized() OVERRIDE {
    // The browser process UI thread has been initialized...
  }
  virtual void OnRenderProcessThreadCreated(CefRefPtr<CefListValue> extra_info)
                                            OVERRIDE {
    // Send startup information to a new render process...
  }

  // CefRenderProcessHandler methods.
  virtual void OnRenderThreadCreated(CefRefPtr<CefListValue> extra_info)
                                     OVERRIDE {
    // The render process main thread has been initialized...
    // Receive startup information in the new render process...
  }
  virtual void OnWebKitInitialized(CefRefPtr<ClientApp> app) OVERRIDE {
    // WebKit has been initialized, register V8 extensions...
  }
  virtual void OnBrowserCreated(CefRefPtr<CefBrowser> browser) OVERRIDE {
    // Browser created in this render process...
  }
  virtual void OnBrowserDestroyed(CefRefPtr<CefBrowser> browser) OVERRIDE {
    // Browser destroyed in this render process...
  }
  virtual bool OnBeforeNavigation(CefRefPtr<CefBrowser> browser,
                                  CefRefPtr<CefFrame> frame,
                                  CefRefPtr<CefRequest> request,
                                  NavigationType navigation_type,
                                  bool is_redirect) OVERRIDE {
    // Allow or block different types of navigation...
  }
  virtual void OnContextCreated(CefRefPtr<CefBrowser> browser,
                                CefRefPtr<CefFrame> frame,
                                CefRefPtr<CefV8Context> context) OVERRIDE {
    // JavaScript context created, add V8 bindings here...
  }
  virtual void OnContextReleased(CefRefPtr<CefBrowser> browser,
                                 CefRefPtr<CefFrame> frame,
                                 CefRefPtr<CefV8Context> context) OVERRIDE {
    // JavaScript context released, release V8 references here...
  }
  virtual bool OnProcessMessageReceived(
      CefRefPtr<CefBrowser> browser,
      CefProcessId source_process,
      CefRefPtr<CefProcessMessage> message) OVERRIDE {
    // Handle IPC messages from the browser process...
  }

  IMPLEMENT_REFCOUNTING(MyApp);
};
```

##### CefClient

The CefClient interface provides access to browser-instance-specific callbacks. A single CefClient instance can be shared among any number of browsers. Important callbacks include:

CefClient提供访问browser-instance-specific的回调接口。单实例CefClient可以共数任意数量的浏览器进程。以下为几个重要的回调：

- Handlers for things like browser life span, context menus, dialogs, display notifications, drag events, focus events, keyboard events and more. The majority of handlers are optional. See the documentation in cef_client.h for the side effects, if any, of not implementing a specific handler.
- **OnProcessMessageReceived** which is called when an IPC message is received from the render process. See the “Inter-Process Communication” section for more information. 

- 比如处理browser的生命周期，右键菜单，对话框，通知显示， 拖曳事件，焦点事件，键盘事件等等。如果没有对某个特定的处理接口进行实现会造成什么影响，请查看cef_client.h文件中相关说明。

Example CefClient implementation:

CefClient子类的例子：

```
// MyHandler implements CefClient and a number of other interfaces.
class MyHandler : public CefClient,
                  public CefContextMenuHandler,
                  public CefDisplayHandler,
                  public CefDownloadHandler,
                  public CefDragHandler,
                  public CefGeolocationHandler,
                  public CefKeyboardHandler,
                  public CefLifeSpanHandler,
                  public CefLoadHandler,
                  public CefRequestHandler {
 public:
  MyHandler();

  // CefClient methods. Important to return |this| for the handler callbacks.
  virtual CefRefPtr<CefContextMenuHandler> GetContextMenuHandler() OVERRIDE {
    return this;
  }
  virtual CefRefPtr<CefDisplayHandler> GetDisplayHandler() OVERRIDE {
    return this;
  }
  virtual CefRefPtr<CefDownloadHandler> GetDownloadHandler() OVERRIDE {
    return this;
  }
  virtual CefRefPtr<CefDragHandler> GetDragHandler() OVERRIDE {
    return this;
  }
  virtual CefRefPtr<CefGeolocationHandler> GetGeolocationHandler() OVERRIDE {
    return this;
  }
  virtual CefRefPtr<CefKeyboardHandler> GetKeyboardHandler() OVERRIDE {
    return this;
  }
  virtual CefRefPtr<CefLifeSpanHandler> GetLifeSpanHandler() OVERRIDE {
    return this;
  }
  virtual CefRefPtr<CefLoadHandler> GetLoadHandler() OVERRIDE {
    return this;
  }
  virtual CefRefPtr<CefRequestHandler> GetRequestHandler() OVERRIDE {
    return this;
  }
  virtual bool OnProcessMessageReceived(CefRefPtr<CefBrowser> browser,
                                        CefProcessId source_process,
                                        CefRefPtr<CefProcessMessage> message)
                                        OVERRIDE {
    // Handle IPC messages from the render process...
  }

  // CefContextMenuHandler methods
  virtual void OnBeforeContextMenu(CefRefPtr<CefBrowser> browser,
                                   CefRefPtr<CefFrame> frame,
                                   CefRefPtr<CefContextMenuParams> params,
                                   CefRefPtr<CefMenuModel> model) OVERRIDE {
    // Customize the context menu...
  }
  virtual bool OnContextMenuCommand(CefRefPtr<CefBrowser> browser,
                                    CefRefPtr<CefFrame> frame,
                                    CefRefPtr<CefContextMenuParams> params,
                                    int command_id,
                                    EventFlags event_flags) OVERRIDE {
    // Handle a context menu command...
  }

  // CefDisplayHandler methods
  virtual void OnLoadingStateChange(CefRefPtr<CefBrowser> browser,
                                    bool isLoading,
                                    bool canGoBack,
                                    bool canGoForward) OVERRIDE {
    // Update UI for browser state...
  }
  virtual void OnAddressChange(CefRefPtr<CefBrowser> browser,
                               CefRefPtr<CefFrame> frame,
                               const CefString& url) OVERRIDE {
    // Update the URL in the address bar...
  }
  virtual void OnTitleChange(CefRefPtr<CefBrowser> browser,
                             const CefString& title) OVERRIDE {
    // Update the browser window title...
  }
  virtual bool OnConsoleMessage(CefRefPtr<CefBrowser> browser,
                                const CefString& message,
                                const CefString& source,
                                int line) OVERRIDE {
    // Log a console message...
  }

  // CefDownloadHandler methods
  virtual void OnBeforeDownload(
      CefRefPtr<CefBrowser> browser,
      CefRefPtr<CefDownloadItem> download_item,
      const CefString& suggested_name,
      CefRefPtr<CefBeforeDownloadCallback> callback) OVERRIDE {
    // Specify a file path or cancel the download...
  }
  virtual void OnDownloadUpdated(
      CefRefPtr<CefBrowser> browser,
      CefRefPtr<CefDownloadItem> download_item,
      CefRefPtr<CefDownloadItemCallback> callback) OVERRIDE {
    // Update the download status...
  }

  // CefDragHandler methods
  virtual bool OnDragEnter(CefRefPtr<CefBrowser> browser,
                           CefRefPtr<CefDragData> dragData,
                           DragOperationsMask mask) OVERRIDE {
    // Allow or deny drag events...
  }

  // CefGeolocationHandler methods
  virtual void OnRequestGeolocationPermission(
      CefRefPtr<CefBrowser> browser,
      const CefString& requesting_url,
      int request_id,
      CefRefPtr<CefGeolocationCallback> callback) OVERRIDE {
    // Allow or deny geolocation API access...
  }

  // CefKeyboardHandler methods
  virtual bool OnPreKeyEvent(CefRefPtr<CefBrowser> browser,
                             const CefKeyEvent& event,
                             CefEventHandle os_event,
                             bool* is_keyboard_shortcut) OVERRIDE {
    // Perform custom handling of key events...
  }

  // CefLifeSpanHandler methods
  virtual bool OnBeforePopup(CefRefPtr<CefBrowser> browser,
                             CefRefPtr<CefFrame> frame,
                             const CefString& target_url,
                             const CefString& target_frame_name,
                             const CefPopupFeatures& popupFeatures,
                             CefWindowInfo& windowInfo,
                             CefRefPtr<CefClient>& client,
                             CefBrowserSettings& settings,
                             bool* no_javascript_access) OVERRIDE {
    // Allow or block popup windows, customize popup window creation...
  }
  virtual void OnAfterCreated(CefRefPtr<CefBrowser> browser) OVERRIDE {
    // Browser window created successfully...
  }
  virtual bool DoClose(CefRefPtr<CefBrowser> browser) OVERRIDE {
    // Allow or block browser window close...
  }
  virtual void OnBeforeClose(CefRefPtr<CefBrowser> browser) OVERRIDE {
    // Browser window is closed, perform cleanup...
  }

  // CefLoadHandler methods
  virtual void OnLoadStart(CefRefPtr<CefBrowser> browser,
                           CefRefPtr<CefFrame> frame) OVERRIDE {
    // A frame has started loading content...
  }
  virtual void OnLoadEnd(CefRefPtr<CefBrowser> browser,
                         CefRefPtr<CefFrame> frame,
                         int httpStatusCode) OVERRIDE {
    // A frame has finished loading content...
  }
  virtual void OnLoadError(CefRefPtr<CefBrowser> browser,
                           CefRefPtr<CefFrame> frame,
                           ErrorCode errorCode,
                           const CefString& errorText,
                           const CefString& failedUrl) OVERRIDE {
    // A frame has failed to load content...
  }
  virtual void OnRenderProcessTerminated(CefRefPtr<CefBrowser> browser,
                                         TerminationStatus status) OVERRIDE {
    // A render process has crashed...
  }

  // CefRequestHandler methods
  virtual CefRefPtr<CefResourceHandler> GetResourceHandler(
      CefRefPtr<CefBrowser> browser,
      CefRefPtr<CefFrame> frame,
      CefRefPtr<CefRequest> request) OVERRIDE {
    // Optionally intercept resource requests...
  }
  virtual bool OnQuotaRequest(CefRefPtr<CefBrowser> browser,
                              const CefString& origin_url,
                              int64 new_size,
                              CefRefPtr<CefQuotaCallback> callback) OVERRIDE {
    // Allow or block quota requests...
  }
  virtual void OnProtocolExecution(CefRefPtr<CefBrowser> browser,
                                   const CefString& url,
                                   bool& allow_os_execution) OVERRIDE {
    // Handle execution of external protocols...
  }

  IMPLEMENT_REFCOUNTING(MyHandler);
};
```

##### Browser Life Span

Browser生命周期

Browser life span begins with a call to CefBrowserHost::CreateBrowser() or CefBrowserHost::CreateBrowserSync(). Convenient places to execute this logic include the CefBrowserProcessHandler::OnContextInitialized() callback or platform-specific message handlers like WM_CREATE on Windows.

Browser生命周期从执行 CefBrowserHost::CreateBrowser() 或者 CefBrowserHost::CreateBrowserSync() 开始。可以在CefBrowserProcessHandler::OnContextInitialized() 回调或者特殊平台例如windows的WM_CREATE 中方便的执行业务逻辑。

```
// Information about the window that will be created including parenting, size, etc.
// The definition of this structure is platform-specific.

// 定义的结构体与平台相关

CefWindowInfo info;
// On Windows for example...
info.SetAsChild(parent_hwnd, client_rect);

// Customize this structure to control browser behavior.
CefBrowserSettings settings;

// CefClient implementation.
CefRefPtr<MyClient> client(new MyClient);

// Create the browser asynchronously. Initially loads the Google URL.
CefBrowserHost::CreateBrowser(info, client.get(), “http://www.google.com”, settings);

The CefLifeSpanHandler class provides the callbacks necessary for managing browser life span. Below is an extract of the relevant methods and members.

CefLifeSpanHandler 类提供管理 browser生合周期必需的回调。以下为相关方法和成员。

class MyClient : public CefClient,
                 public CefLifeSpanHandler,
                 ... {
  // CefClient methods.
  virtual CefRefPtr<CefLifeSpanHandler> GetLifeSpanHandler() OVERRIDE {
    return this;
  }

  // CefLifeSpanHandler methods.
  void OnAfterCreated(CefRefPtr<CefBrowser> browser) OVERRIDE;
  bool DoClose(CefRefPtr<CefBrowser> browser) OVERRIDE;
  void OnBeforeClose(CefRefPtr<CefBrowser> browser) OVERRIDE;

  // Member accessors.
  CefRefPtr<CefBrowser> GetBrower() { return m_Browser; }
  bool IsClosing() { return m_bIsClosing; }

 private:
  CefRefPtr<CefBrowser> m_Browser;
  int m_BrowserId;
  int m_BrowserCount;
  bool m_bIsClosing;

  IMPLEMENT_REFCOUNTING(MyHandler);
  IMPLEMENT_LOCKING(MyHandler);
};
```

The OnAfterCreated() method will be called immediately after the browser object is created. The host application can use this method to keep a reference to the main browser object.

当browser对象创建后OnAfterCreated() 方法立即执行。宿主程序可以用这个方法来保持对browser对象的引用。

```
void MyClient::OnAfterCreated(CefRefPtr<CefBrowser> browser) {
  // Must be executed on the UI thread.
  REQUIRE_UI_THREAD();
  // Protect data members from access on multiple threads.
  AutoLock lock_scope(this);

  if (!m_Browser.get())   {
    // Keep a reference to the main browser.
    m_Browser = browser;
    m_BrowserId = browser->GetIdentifier();
  }

  // Keep track of how many browsers currently exist.
  m_BrowserCount++;
}
```

To destroy the browser call CefBrowserHost::CloseBrowser().

执行CefBrowserHost::CloseBrowser()销毁browser对象。

```
// Notify the browser window that we would like to close it. This will result in a call to 
// MyHandler::DoClose() if the JavaScript 'onbeforeunload' event handler allows it.
browser->GetHost()->CloseBrowser(false);
```

If the browser is parented to another window then the close event may originate in the OS function for that parent window (for example, by clicking the X on the parent window). The parent window then needs to call CloseBrowser(false) and wait for a second OS close event to indicate that the browser has allowed the close. The second OS close event will not be sent if the close is canceled by a JavaScript ‘onbeforeunload’ event handler or by the DoClose() callback. Notice the IsClosing() check in the below examples -- it will return false for the first OS close event and true for the second (after DoClose is called).

browser对象的关闭事件来源于他的父窗口的关闭方法（比如，在父窗口上点击X控钮。）。父窗口需要调用  CloseBrowser(false) 并且等待操作系统的第二个关闭事件来决定是否允许关闭。如果在JavaScript 'onbeforeunload'事件处理或者 DoClose()回调中取消了关闭操作，则操作系统的第二个关闭事件可能不会发送。注意一下面示例中对IsCloseing()的判断-它在第一个关闭事件中返回false，在第二个关闭事件中返回true(当 DoCloase 被调用后)。

Handling in the parent window WndProc on Windows:
Windows平台下，在父窗口的WndProc里处理WM_ClOSE消息：
```
case WM_CLOSE:
  if (g_handler.get() && !g_handler->IsClosing()) {
    CefRefPtr<CefBrowser> browser = g_handler->GetBrowser();
    if (browser.get()) {
      // Notify the browser window that we would like to close it. This will result in a call to 
      // MyHandler::DoClose() if the JavaScript 'onbeforeunload' event handler allows it.
      browser->GetHost()->CloseBrowser(false);

      // Cancel the close.
      return 0;
    }
  }

  // Allow the close.
  break;

case WM_DESTROY:
  // Quitting CEF is handled in MyHandler::OnBeforeClose().
  return 0;
}
```

Handling the “delete_event” signal on Linux:
Linux平台下，处理`delete_event`信号:
```
gboolean delete_event(GtkWidget* widget, GdkEvent* event,
                      GtkWindow* window) {
  if (g_handler.get() && !g_handler->IsClosing()) {
    CefRefPtr<CefBrowser> browser = g_handler->GetBrowser();
    if (browser.get()) {
      // Notify the browser window that we would like to close it. This will result in a call to 
      // MyHandler::DoClose() if the JavaScript 'onbeforeunload' event handler allows it.
      browser->GetHost()->CloseBrowser(false);

      // Cancel the close.
      return TRUE;
    }
  }

  // Allow the close.
  return FALSE;
}
```

Handling the windowShouldClose: selector on Mac OS X:
MacOS X平台下，处理windowShouldClose选择器:
```
// Called when the window is about to close. Perform the self-destruction
// sequence by getting rid of the window. By returning YES, we allow the window
// to be removed from the screen.
- (BOOL)windowShouldClose:(id)window {
  if (g_handler.get() && !g_handler->IsClosing()) {
    CefRefPtr<CefBrowser> browser = g_handler->GetBrowser();
    if (browser.get()) {
      // Notify the browser window that we would like to close it. This will result in a call to 
      // MyHandler::DoClose() if the JavaScript 'onbeforeunload' event handler allows it.
      browser->GetHost()->CloseBrowser(false);

      // Cancel the close.
      return NO;
    }
  }

  // Try to make the window go away.
  [window autorelease];

  // Clean ourselves up after clearing the stack of anything that might have the
  // window on it.
  [self performSelectorOnMainThread:@selector(cleanup:)
                         withObject:window
                      waitUntilDone:NO];

  // Allow the close.
  return YES;
}
```

The DoClose() method sets the m_bIsClosing flag and returns false to send the second OS close event.

DoClose方法设置m_blsClosing 标志位为true，并返回false以再次发送操作系统的关闭事件。

```
bool MyClient::DoClose(CefRefPtr<CefBrowser> browser) {
  // Must be executed on the UI thread.
  REQUIRE_UI_THREAD();
  // Protect data members from access on multiple threads.
  AutoLock lock_scope(this);

  // Closing the main window requires special handling. See the DoClose()
  // documentation in the CEF header for a detailed description of this
  // process.
  if (m_BrowserId == browser->GetIdentifier()) {
    // Notify the browser that the parent window is about to close.
    browser->GetHost()->ParentWindowWillClose();

    // Set a flag to indicate that the window close should be allowed.
    m_bIsClosing = true;
  }

  // Allow the close. For windowed browsers this will result in the OS close
  // event being sent.
  return false;
}
```

When the OS function receives the second OS close event it allows the parent window to actually close. This then results in a call to OnBeforeClose(). Make sure to release any references to the browser object in the OnBeforeClose() callback.

当操作系统捕捉到第二次关闭事件，它才会允许父窗口真正关闭。该动作会先触发OnBeforeClose()回调，请在该回调里释放所有对浏览器对象的引用。

```
void MyHandler::OnBeforeClose(CefRefPtr<CefBrowser> browser) {
  // Must be executed on the UI thread.
  REQUIRE_UI_THREAD();
  // Protect data members from access on multiple threads.
  AutoLock lock_scope(this);

  if (m_BrowserId == browser->GetIdentifier()) {
    // Free the browser pointer so that the browser can be destroyed.
    m_Browser = NULL;
  }

  if (--m_BrowserCount == 0) {
    // All browser windows have closed. Quit the application message loop.
    CefQuitMessageLoop();
  }
}
```

See the cefclient application for complete working examples on each platform.

完整的流程，请参考cefclient例子里对不同平台的处理。

##### 离屏渲染(Off-Screen Rendering)

With off-screen rendering CEF does not create a native browser window. Instead, CEF provides the host application with invalidated regions and a pixel buffer and the host application notifies CEF of mouse, keyboard and focus events. Off-screen rendering does not currently support accelerated compositing so performance may suffer as compared to a windowed browser. Off-screen browsers will receive the same notifications as windowed browsers including the life span notifications described in the previous section. To use off-screen rendering:

在离屏渲染模式下，CEF不会创建原生浏览器窗口。CEF为宿主程序提供无效的区域和像素缓存区，而宿主程序负责通知鼠标键盘以及焦点事件给CEF。离屏渲染目前不支持混合加速，所以性能上可能无法和非离屏渲染相比。离屏浏览器将收到和窗口浏览器同样的事件通知，例如前一节介绍的生命周期事件。下面介绍如何使用离屏渲染：

- Implement the CefRenderHandler interface. All methods are required unless otherwise indicated.
- Call CefWindowInfo::SetAsOffScreen() and optionally CefWindowInfo::SetTransparentPainting() before passing the CefWindowInfo structure to CefBrowserHost::CreateBrowser(). If no parent window is passed to SetAsOffScreen some functionality like context menus may not be available.
- The CefRenderHandler::GetViewRect() method will be called to retrieve the desired view rectangle.
- The CefRenderHandler::OnPaint() method will be called to provide invalid regions and the updated pixel buffer. The cefclient application draws the buffer using OpenGL but your application can use whatever technique you prefer.
- To resize the browser call CefBrowserHost::WasResized(). This will result in a call to GetViewRect() to retrieve the new size followed by a call to OnPaint().
- Call the CefBrowserHost::SendXXX() methods to notify the browser of mouse, keyboard and focus events.
- Call CefBrowserHost::CloseBrowser() to destroy browser. 

-  实现CefRenderHandler接口。除非特别说明，所有的方法都需要覆写。
-  调用CefWindowInfo::SetAsOffScreen()，将CefWindowInfo传递给CefBrowserHost::CreateBrowser()之前还可以选择设置CefWindowInfo::SetTransparentPainting()。如果没有父窗口被传递给SetAsOffScreen,则有些类似上下文菜单这样的功能将不可用。
-  CefRenderHandler::GetViewRect方法将被调用以获得所需要的可视区域。
-  CefRenderHandler::OnPaint() 方法将被调用以提供无效区域（脏区域）以及更新过的像素缓存。cefclient程序里使用OpenGL绘制缓存，但你可以使用任何别的绘制技术。
-  可以调用CefBrowserHost::WasResized()方法改变浏览器大小。这将导致对GetViewRect()方法的调用，以获取新的浏览器大小，然后调用OnPaint()重新绘制。
-  调用CefBrowserHost::SendXXX()方法通知浏览器的鼠标、键盘和焦点事件。
-  调用CefBrowserHost::CloseBrowser()销毁浏览器。

Run cefclient with the “--off-screen-rendering-enabled” command-line flag for a working example.
使用命令行参数`--off-screen-rendering-enabled`运行cefclient，可以测试离屏渲染的效果。

##### 投递任务(Posting Tasks)

Tasks can be posted between the various threads in a single process using the CefPostTask family of methods (see the cef_task.h header file for the complete list). The task will execute asynchronously on the message loop of the target thread. For example, to execute the MyObject::MyMethod method on the UI thread and pass it two parameters:

任务(Task)可以通过CefPostTask在一个进程内的不同的线程之间投递。CefPostTask有一系列的重载方法，详细内容请参考cef_task.h头文件。任务将会在被投递线程的消息循环里异步执行。例如，为了在UI线程上执行MyObject::MyMethod方法，并传递两个参数，代码如下：

```
CefPostTask(TID_UI, NewCefRunnableMethod(object, &MyObject::MyMethod, param1, param2));
```

To execute the MyFunction function on the IO thread and pass it two parameters:
为了在IO线程在执行MyFunction方法，同时传递两个参数，代码如下：

```
CefPostTask(TID_IO, NewCefRunnableFunction(MyFunction, param1, param2));
```

See the cef_runnable.h header file for more information on the NewCefRunnable templates.
参考cef_runnable.h头文件以了解更多关于NewCefRunnable模板方法的细节。

If the host application needs to keep a reference to a run loop it can use the CefTaskRunner class. For example, to get the task runner for the UI thread:

如果宿主程序需要保留一个运行循环的引用，则可以使用CefTaskRunner类。例如，获取UI线程的任务运行器(task runner)，代码如下：

```
CefRefPtr<CefTaskRunner> task_runner = CefTaskRunner::GetForThread(TID_UI);
```

##### 进程间通信(Inter-Process Communication (IPC))

Since CEF3 runs in multiple processes it is necessary to provide mechanisms for communicating between those processes. CefBrowser and CefFrame objects exist in both the browser and render processes which helps to facilitate this process. Each CefBrowser and CefFrame object also has a unique ID value associated with it that will match on both sides of the process boundary.

由于CEF3运行在多进程环境下，所以需要提供一个进程间通信机制。CefBrowser和CefFrame对象在Borwser和Render进程里都有代理对象。CefBrowser和CefFrame对象都有一个唯一ID值绑定，便于在两个进程间定位匹配的代理对象。

###### Process Startup Messages
###### 处理启动消息

To provide all render processes with the same information on startup implement CefBrowserProcessHandler::OnRenderProcessThreadCreated() in the browser process. This will pass information to CefRenderProcessHandler::OnRenderThreadCreated() in the render process.
为了给所有的Render进程提供一样的启动信息，请在Browser进程实现CefBrowserProcessHander::OnRenderProcessThreadCreated()方法。在这里传入的信息会在Render进程的CefRenderProcessHandler::OnRenderThreadCreated()方法里接受。

###### Process Runtime Messages
###### 处理运行时消息

To pass information at any time during the process lifespan use process messages via the CefProcessMessage class. These messages are associated with a specific CefBrowser instance and are sent using the CefBrowser::SendProcessMessage() method. The process message should contain whatever state information is required via CefProcessMessage::GetArgumentList().
在进程声明周期内，任何时候你都可以通过CefProcessMessage类传递进程间消息。这些信息和特定的CefBrowser实例绑定在一起，用户可以通过CefBrowser::SendProcessMessage()方法发送。进程间消息可以包含任意的状态信息，用户可以通过CefProcessMessage::GetArgumentList()获取。

```
// Create the message object.
CefRefPtr<CefProcessMessage> msg= CefProcessMessage::Create(“my_message”);

// Retrieve the argument list object.
CefRefPtr<CefListValue> args = msg>GetArgumentList();

// Populate the argument values.
args->SetString(0, “my string”);
args->SetInt(0, 10);

// Send the process message to the render process.
// Use PID_BROWSER instead when sending a message to the browser process.
browser->SendProcessMessage(PID_RENDERER, msg);
```

A message sent from the browser process to the render process will arrive in CefRenderProcessHandler::OnProcessMessageReceived(). A message sent from the render process to the browser process will arrive in CefClient::OnProcessMessageReceived().
一个从Browser进程发送到Render进程的消息将会在CefRenderProcessHandler::OnProcessMessageReceived()方法里被接收。一个从Render进程发送到Browser进程的消息将会在CefClient::OnProcessMessageReceived()方法里被接收。

```
bool MyHandler::OnProcessMessageReceived(
    CefRefPtr<CefBrowser> browser,
    CefProcessId source_process,
    CefRefPtr<CefProcessMessage> message) {
  // Check the message name.
  const std::string& message_name = message->GetName();
  if (message_name == “my_message”) {
    // Handle the message here...
    return true;
  }
  return false;
}
```

To associate the message with a particular CefFrame pass the frame ID (retrievable via CefFrame::GetIdentifier()) as an argument and retrieve the associated CefFrame in the receiving process via the CefBrowser::GetFrame() method.

我们可以调用CefFrame::GerIdentifier()获取CefFrame的ID，并通过进程间消息发送给另一个进程，然后在接收端通过CefBrowser::GetFrame()找到对应的CefFrame。通过这种方式可以将进程间消息和特定的CefFrame联系在一起。

```
// Helper macros for splitting and combining the int64 frame ID value.
#define MAKE_INT64(int_low, int_high) \
    ((int64) (((int) (int_low)) | ((int64) ((int) (int_high))) << 32))
#define LOW_INT(int64_val) ((int) (int64_val))
#define HIGH_INT(int64_val) ((int) (((int64) (int64_val) >> 32) & 0xFFFFFFFFL))

// Sending the frame ID.
const int64 frame_id = frame->GetIdentifier();
args->SetInt(0, LOW_INT(frame_id));
args->SetInt(1, HIGH_INT(frame_id));

// Receiving the frame ID.
const int64 frame_id = MAKE_INT64(args->GetInt(0), args->GetInt(1));
CefRefPtr<CefFrame> frame = browser->GetFrame(frame_id);
```

##### Asynchronous JavaScript Bindings
##### 异步JavaScript绑定

JavaScriptIntegration is implemented in the render process but frequently need to communicate with the browser process. The JavaScript APIs themselves should be designed to work asynchronously using closures and promises.

JavaScritp被集成在Render进程，但是需要频繁和Browser进程交互。 JavaScript API应该被设计成可使用闭包异步执行。

Generic Message Router

###### 通用消息转发

Starting with trunk revision 1574 CEF provides a generic implementation for routing asynchronous messages between JavaScript running in the renderer process and C++ running in the browser process. An application interacts with the router by passing it data from standard CEF C++ callbacks (OnBeforeBrowse, OnProcessMessageRecieved, OnContextCreated, etc). The renderer-side router supports generic JavaScript callback registration and execution while the browser-side router supports application-specific logic via one or more application-provided Handler instances.

从1574版本开始，CEF提供了在Render进程执行的JavaScript和在Browser进程执行的C++代码之间同步通信的转发器。应用程序通过C++回调函数（OnBeforeBrowse, OnProcessMessageRecieved, OnContextCreated等）传递数据。Render进程支持通用的JavaScript回调函数注册机制，Browser进程则支持应用程序注册特定的Handler进行处理。

The JavaScript bindings look like this:
下面的代码示例在JavaScript端扩展window对象，添加cefQuery函数：

```
// Create and send a new query.
var request_id = window.cefQuery({
    request: 'my_request',
    persistent: false,
    onSuccess: function(response) {},
    onFailure: function(error_code, error_message) {}
});

// Optionally cancel the query.
window.cefQueryCancel(request_id);
```

The C++ handler looks like this:
对应的C++ Handler代码如下：

```
class Callback : public CefBase {
 public:
  ///
  // Notify the associated JavaScript onSuccess callback that the query has
  // completed successfully with the specified |response|.
  ///
  virtual void Success(const CefString& response) =0;

  ///
  // Notify the associated JavaScript onFailure callback that the query has
  // failed with the specified |error_code| and |error_message|.
  ///
  virtual void Failure(int error_code, const CefString& error_message) =0;
};

class Handler {
 public:
  ///
  // Executed when a new query is received. |query_id| uniquely identifies the
  // query for the life span of the router. Return true to handle the query
  // or false to propagate the query to other registered handlers, if any. If
  // no handlers return true from this method then the query will be
  // automatically canceled with an error code of -1 delivered to the
  // JavaScript onFailure callback. If this method returns true then a
  // Callback method must be executed either in this method or asynchronously
  // to complete the query.
  ///
  virtual bool OnQuery(CefRefPtr<CefBrowser> browser,
                       CefRefPtr<CefFrame> frame,
                       int64 query_id,
                       const CefString& request,
                       bool persistent,
                       CefRefPtr<Callback> callback) {
    return false;
  }

  ///
  // Executed when a query has been canceled either explicitly using the
  // JavaScript cancel function or implicitly due to browser destruction,
  // navigation or renderer process termination. It will only be called for
  // the single handler that returned true from OnQuery for the same
  // |query_id|. No references to the associated Callback object should be
  // kept after this method is called, nor should any Callback methods be
  // executed.
  ///
  virtual void OnQueryCanceled(CefRefPtr<CefBrowser> browser,
                               CefRefPtr<CefFrame> frame,
                               int64 query_id) {}
};
```

See wrapper/cef_message_router.h for complete usage documentation.
完整的用法请参考wrapper/cef_message_router.h

###### Custom Implementation
###### 自定义实现

A CEF-based application can also provide its own custom implementation of asynchronous JavaScript bindings. A simplistic implementation could work as follows:

一个CEF应用程序也可以提供自己的异步JavaScript绑定。典型的实现如下：

1. The JavaScript binding in the render process is passed a callback function.
1. Render进程的JavaScript传递一个回调函数。
```
// In JavaScript register the callback function.
app.setMessageCallback('binding_test', function(name, args) {
  document.getElementById('result').value = "Response: "+args[0];
});
```

2. The render process keeps a reference to the callback function.
2. Render进程的C++端通过一个map持有JavaScript端注册的回调函数。

```
// Map of message callbacks.
typedef std::map<std::pair<std::string, int>,
                 std::pair<CefRefPtr<CefV8Context>, CefRefPtr<CefV8Value> > >
                 CallbackMap;
CallbackMap callback_map_;

// In the CefV8Handler::Execute implementation for “setMessageCallback”.
if (arguments.size() == 2 && arguments[0]->IsString() &&
    arguments[1]->IsFunction()) {
  std::string message_name = arguments[0]->GetStringValue();
  CefRefPtr<CefV8Context> context = CefV8Context::GetCurrentContext();
  int browser_id = context->GetBrowser()->GetIdentifier();
  callback_map_.insert(
      std::make_pair(std::make_pair(message_name, browser_id),
                     std::make_pair(context, arguments[1])));
}
```

3. The render process sends an asynchronous IPC message to the browser process requesting that work be performed.
3. Render进程发送异步进程间通信到Browser进程。

4. The browser process receives the IPC message and performs the work.
4. Browser进程接收到进程间消息，并处理。

5. Upon completion of the work the browser process sends an asynchronous IPC message back to the render process with the result.
5. Browser进程处理完毕后，发送一个异步进程间消息给Render进程，返回结果。

6. The render process receives the IPC message and executes the callback function with the result.
6. Render进程接收到进程间消息，则调用最开始保存的JavaScript注册的回调函数处理之。

```
// Execute the registered JavaScript callback if any.
if (!callback_map_.empty()) {
  const CefString& message_name = message->GetName();
  CallbackMap::const_iterator it = callback_map_.find(
      std::make_pair(message_name.ToString(),
                     browser->GetIdentifier()));
  if (it != callback_map_.end()) {
    // Keep a local reference to the objects. The callback may remove itself
    // from the callback map.
    CefRefPtr<CefV8Context> context = it->second.first;
    CefRefPtr<CefV8Value> callback = it->second.second;

    // Enter the context.
    context->Enter();

    CefV8ValueList arguments;

    // First argument is the message name.
    arguments.push_back(CefV8Value::CreateString(message_name));

    // Second argument is the list of message arguments.
    CefRefPtr<CefListValue> list = message->GetArgumentList();
    CefRefPtr<CefV8Value> args = CefV8Value::CreateArray(list->GetSize());
    SetList(list, args);  // Helper function to convert CefListValue to CefV8Value.
    arguments.push_back(args);

    // Execute the callback.
    CefRefPtr<CefV8Value> retval = callback->ExecuteFunction(NULL, arguments);
    if (retval.get()) {
      if (retval->IsBool())
        handled = retval->GetBoolValue();
    }

    // Exit the context.
    context->Exit();
  }
}
```

7. Release any V8 references associated with the context in CefRenderProcessHandler::OnContextReleased().
7. 在CefRenderProcessHandler::OnContextReleased()里释放JavaScript注册的回调函数以及其他V8资源。

```
void MyHandler::OnContextReleased(CefRefPtr<CefBrowser> browser,
                                  CefRefPtr<CefFrame> frame,
                                  CefRefPtr<CefV8Context> context) {
  // Remove any JavaScript callbacks registered for the context that has been released.
  if (!callback_map_.empty()) {
    CallbackMap::iterator it = callback_map_.begin();
    for (; it != callback_map_.end();) {
      if (it->second.first->IsSame(context))
        callback_map_.erase(it++);
      else
        ++it;
    }
  }
}
```

##### Synchronous Requests
##### 同步请求

In rare cases it may be necessary to implement synchronous communication between the browser and render processes. This should be avoided whenever possible because it will negatively impact performance in the render process. However, if you must have synchronous communication consider using synchronous XMLHttpRequests which will block the render process while awaiting handling in the browser process network layer. The browser process can handle the requests using a custom scheme handler or network interception. See the “Network Layer” section for more information.

某些特殊场景下，也许会需要在Browser进程和Render进程做进程间同步通信。这应该被尽可能避免，因为这会对Render进程的性能造成负面影响。然而如果你一定要做进程间同步通信，可以考虑使用XMLHttpRequest，XMLHttpRequest在等待Browser进程的网络响应的时候会等待。Browser进程可以通过自定义scheme Handler或者网络交互处理XMLHttpRequest。更多细节，请参考`Network Layer`一节。

###### 网络层(Network Layer)

By default network requests in CEF3 will be handled in a manner transparent to the host application. For applications wishing for a closer relationship with the network layer CEF3 exposes a range of network-related functionalities.

默认情况下，CEF3的网络请求会被宿主程序手工处理。然而CEF3也暴露了一系列网络相关的函数用以处理网络请求。


Network-related callbacks can occur on different threads so make sure to pay attention to the documentation and properly protect your data members.

网络相关的回调函数可在不同线程被调用，因此要注意相关文档的说明，并对自己的数据进行线程安全保护。

###### 自定义请求(Custom Requests)


The simplest way to load a URL in a browser frame is via the CefFrame::LoadURL() method.
通过CefFrame::LoadURL()方法可简单加载一个url:

```
browser->GetMainFrame()->LoadURL(some_url);
```

Applications wishing to send more complex requests containing custom request headers or upload data can use the CefFrame::LoadRequest() method. This method accepts a CefRequest object as the single argument.

如果希望发送更复杂的请求，则可以调用CefFrame::LoadRequest()方法。该方法接受一个CefRequest对象作为唯一的参数。


```
// Create a CefRequest object.
CefRefPtr<CefRequest> request = CefRequest::Create();

// Set the request URL.
request->SetURL(some_url);

// Set the request method. Supported methods include GET, POST, HEAD, DELETE and PUT.
request->SetMethod(“POST”);

// Optionally specify custom headers.
CefRequest::HeaderMap headerMap;
headerMap.insert(
    std::make_pair("X-My-Header", "My Header Value"));
request->SetHeaderMap(headerMap);

// Optionally specify upload content.
// The default “Content-Type” header value is "application/x-www-form-urlencoded".
// Set “Content-Type” via the HeaderMap if a different value is desired.
const std::string& upload_data = “arg1=val1&arg2=val2”;
CefRefPtr<CefPostData> postData = CefPostData::Create();
CefRefPtr<CefPostDataElement> element = CefPostDataElement::Create();
element->SetToBytes(upload_data.size(), upload_data.c_str());
postData->AddElement(element);
request->SetPostData(postData);
```

####### 浏览器无关请求(Browser-Independent Requests)

Applications can send network requests not associated with a particular browser via the CefURLRequest class. Implement the CefURLRequestClient interface to handle the resulting response. CefURLRequest can be used in both the browser and render processes.

应用程序可以通过CefURLRequest类发送和浏览器无关的网络请求。并实现CefURLRequestClient接口处理响应。CefURLRequest可以在Browser和Render进程被使用。

```
class MyRequestClient : public CefURLRequestClient {
 public:
  MyRequestClient()
    : upload_total_(0),
      download_total_(0) {}

  virtual void OnRequestComplete(CefRefPtr<CefURLRequest> request) OVERRIDE {
    CefURLRequest::Status status = request->GetRequestStatus();
    CefURLRequest::ErrorCode error_code = request->GetRequestError();
    CefRefPtr<CefResponse> response = request->GetResponse();

    // Do something with the response...
  }

  virtual void OnUploadProgress(CefRefPtr<CefURLRequest> request,
                                uint64 current,
                                uint64 total) OVERRIDE {
    upload_total_ = total;
  }

  virtual void OnDownloadProgress(CefRefPtr<CefURLRequest> request,
                                  uint64 current,
                                  uint64 total) OVERRIDE {
    download_total_ = total;
  }

  virtual void OnDownloadData(CefRefPtr<CefURLRequest> request,
                              const void* data,
                              size_t data_length) OVERRIDE {
    download_data_ += std::string(static_cast<const char*>(data), data_length);
  }

 private:
  uint64 upload_total_;
  uint64 download_total_;
  std::string download_data_;

 private:
  IMPLEMENT_REFCOUNTING(MyRequestClient);
};
```

To send the request:
下面的代码发送一个请求:

```
// Set up the CefRequest object.
CefRefPtr<CefRequest> request = CefRequest::Create();
// Populate |request| as shown above...

// Create the client instance.
CefRefPtr<MyRequestClient> client = new MyRequestClient();

// Start the request. MyRequestClient callbacks will be executed asynchronously.
CefRefPtr<CefURLRequest> url_request = CefURLRequest::Create(request, client.get());
// To cancel the request: url_request->Cancel();
```

Requests made with CefURLRequest can also specify custom behaviors via the CefRequest::SetFlags() method. Supported bit flags include:
可以通过CefRequest::SetFlags定制请求的行为，这些标志位包括：

- **UR_FLAG_SKIP_CACHE** If set the cache will be skipped when handling the request.
- **UR_FLAG_ALLOW_CACHED_CREDENTIALS** If set cookies may be sent with the request and saved from the response. UR_FLAG_ALLOW_CACHED_CREDENTIALS must also be set.
- **UR_FLAG_REPORT_UPLOAD_PROGRESS** If set upload progress events will be generated when a request has a body.
- **UR_FLAG_REPORT_LOAD_TIMING** If set load timing info will be collected for the request.
- **UR_FLAG_REPORT_RAW_HEADERS** If set the headers sent and received for the request will be recorded.
- **UR_FLAG_NO_DOWNLOAD_DATA** If set the CefURLRequestClient::OnDownloadData method will not be called.
- **UR_FLAG_NO_RETRY_ON_5XX** If set 5XX redirect errors will be propagated to the observer instead of automatically re-tried. This currently only applies for requests originated in the browser process. 

- **UR_FLAG_SKIP_CACHE** If set the cache will be skipped when handling the request.
如果设置了该标志位，则处理请求响应时，缓存将被跳过。
- **UR_FLAG_ALLOW_CACHED_CREDENTIALS** If set cookies may be sent with the request and saved from the response. UR_FLAG_ALLOW_CACHED_CREDENTIALS must also be set.
如果设置了该标志位，则可能会发送cookie并在响应端被保存。同时UR_FLAG_ALLOW_CACHED_CREDENTIALS标志位也必须被设置。
- **UR_FLAG_REPORT_UPLOAD_PROGRESS** If set upload progress events will be generated when a request has a body.
如果设置了该标志位，则当请求拥有请求体时，上载进度事件将会被触发。

- **UR_FLAG_REPORT_LOAD_TIMING** If set load timing info will be collected for the request.
如果设置了该标志位，则时间信息会被收集。

- **UR_FLAG_REPORT_RAW_HEADERS** If set the headers sent and received for the request will be recorded.
如果设置了该标志位，则头部会被发送，并且接收端会被记录。

- **UR_FLAG_NO_DOWNLOAD_DATA** If set the CefURLRequestClient::OnDownloadData method will not be called.
如果设置了该标志位，则CefURLRequestClient::OnDownloadData方法不会被调用。

- **UR_FLAG_NO_RETRY_ON_5XX** If set 5XX redirect errors will be propagated to the observer instead of automatically re-tried. This currently only applies for requests originated in the browser process. 
如果设置了该标志位，则5xx重定向错误会被交给相关Observer去处理，而不是自动重试。这个功能目前只能在Browser进程的请求端使用。

For example, to skip the cache and not report download data:
例如，为了跳过缓存并不报告下载数据，代码如下：

```
request->SetFlags(UR_FLAG_SKIP_CACHE | UR_FLAG_NO_DOWNLOAD_DATA);
```

####### 请求响应(Request Handling)

CEF3 supports two approaches for handling network requests inside of an application. The scheme handler approach allows registration of a handler for requests targeting a particular origin (scheme + domain). The request interception approach allows handling of arbitrary requests at the application’s discretion.

CEF3 支持两种方式处理网络请求。一种是实现scheme Handler，这种方式允许为特定的(sheme+domain)请求注册特定的请求响应。另一种是请求拦截，允许处理任意的网络请求。

It’s important to register custom schemes (anything other than “HTTP”, “HTTPS”, etc) with CEF so that they’ll behave as expected. For example, if you wish your scheme to behave the same as HTTP (support POST requests and enforce HTTP access control (CORS) restrictions) it should be registered as a “standard” scheme. If your custom scheme will be performing cross-origin requests to other schemes consider using the HTTP scheme instead of a custom scheme to avoid potential issues. If you wish to use custom schemes the attributes are registered via the CefApp::OnRegisterCustomSchemes() callback.
注册自定义scheme（有别于`HTTP`,`HTTPS`等）可以让CEF按希望的方式处理请求。例如，如果你希望特定的shceme被当策划那个HTTP一样处理，则应该注册一个`standard`的scheme。如果你的自定义shceme可被跨域执行，则应该考虑使用使用HTTP scheme代替自定义scheme以避免潜在问题。如果你希望使用自定义scheme，实现CefApp::OnRegisterCustomSchemes回调。

```
void MyApp::OnRegisterCustomSchemes(CefRefPtr<CefSchemeRegistrar> registrar) {
  // Register "client" as a standard scheme.
  registrar->AddCustomScheme("client", true, false, false);
}
```

##### Scheme响应(Scheme Handler)

A scheme handler is registered via the CefRegisterSchemeHandlerFactory() function. A good place to call this function is from CefBrowserProcessHandler::OnContextInitialized(). For example, you can register a handler for “client://myapp/” requests:

通过CefRegisterSchemeHandlerFactory方法注册一个scheme响应，最好在CefBrowserProcessHandler::OnContextInitialized()方法里调用。例如，你可以注册一个"client://myapp/"的请求：

```
CefRegisterSchemeHandlerFactory("client", “myapp”, new MySchemeHandlerFactory());
```

Handlers can be used with both built-in schemes (HTTP, HTTPS, etc) and custom schemes. When using a built-in scheme choose a domain name unique to your application (like “myapp” or “internal”). Implement the CefSchemeHandlerFactory and CefResourceHandler classes to handle the request and provide response data. See the “Scheme Handler” test in cefclient (implemented in scheme_test.[cpp|h]) for a working example.

scheme Handler类可以被用在内置shcme(HTTP,HTTPS等），也可以被用在自定义scheme上。当使用内置shceme，选择一个对你的应用程序来说唯一的域名。实现CefSchemeHandlerFactory和CefResoureHandler类去处理请求并返回响应数据。可以参考cefclient/sheme_test.h的例子。

```
// Implementation of the factory for creating client request handlers.
class MySchemeHandlerFactory : public CefSchemeHandlerFactory {
 public:
  virtual CefRefPtr<CefResourceHandler> Create(CefRefPtr<CefBrowser> browser,
                                               CefRefPtr<CefFrame> frame,
                                               const CefString& scheme_name,
                                               CefRefPtr<CefRequest> request)
                                               OVERRIDE {
    // Return a new resource handler instance to handle the request.
    return new MyResourceHandler();
  }

  IMPLEMENT_REFCOUNTING(MySchemeHandlerFactory);
};

// Implementation of the resource handler for client requests.
class MyResourceHandler : public CefResourceHandler {
 public:
  MyResourceHandler() {}

  virtual bool ProcessRequest(CefRefPtr<CefRequest> request,
                              CefRefPtr<CefCallback> callback)
                              OVERRIDE {
    // Evaluate |request| to determine proper handling...
    // Execute |callback| once header information is available.
    // Return true to handle the request.
    return true;
  }

  virtual void GetResponseHeaders(CefRefPtr<CefResponse> response,
                                  int64& response_length,
                                  CefString& redirectUrl) OVERRIDE {
    // Populate the response headers.
    response->SetMimeType("text/html");
    response->SetStatus(200);

    // Specify the resulting response length.
    response_length = ...;
  }

  virtual void Cancel() OVERRIDE {
    // Cancel the response...
  }

  virtual bool ReadResponse(void* data_out,
                            int bytes_to_read,
                            int& bytes_read,
                            CefRefPtr<CefCallback> callback)
                            OVERRIDE {
    // Read up to |bytes_to_read| data into |data_out| and set |bytes_read|.
    // If data isn't immediately available set bytes_read=0 and execute
    // |callback| asynchronously.
    // Return true to continue the request or false to complete the request.
    return …;
  }

 private:
  IMPLEMENT_REFCOUNTING(MyResourceHandler);
};
```

If the response data is known at request time the CefStreamResourceHandler class provides a convenient default implementation of CefResourceHandler.

如果响应数据类型是已知的，则CefStreamResourceHandler类提供了CefResourceHandler类的默认实现。

```
// CefStreamResourceHandler is part of the libcef_dll_wrapper project.
#include “include/wrapper/cef_stream_resource_handler.h”

const std::string& html_content = “<html><body>Hello!</body></html>”;

// Create a stream reader for |html_content|.
CefRefPtr<CefStreamReader> stream =
    CefStreamReader::CreateForData(
        static_cast<void*>(const_cast<char*>(html_content.c_str())),
        html_content.size());

// Constructor for HTTP status code 200 and no custom response headers.
// There’s also a version of the constructor for custom status code and response headers.
return new CefStreamResourceHandler("text/html", stream);
```

##### 请求拦截(Request Interception)

The CefRequestHandler::GetResourceHandler() method supports the interception of arbitrary requests. It uses the same CefResourceHandler class as the scheme handler approach. See ClientHandler::GetResourceHandler (implemented in client_handler.cpp) for a working example.

CefRequestHandler::GetResourceHandler()方法支持拦截任意请求。参考client_handler.cpp。

```
CefRefPtr<CefResourceHandler> MyHandler::GetResourceHandler(
      CefRefPtr<CefBrowser> browser,
      CefRefPtr<CefFrame> frame,
      CefRefPtr<CefRequest> request) {
  // Evaluate |request| to determine proper handling...
  if (...)
    return new MyResourceHandler();

  // Return NULL for default handling of the request.
  return NULL;
}
```

####### 其他回调(Other Callbacks)

The CefRequestHandler interface provides callbacks for various network-related events incuding authentication, cookie handling, external protocol handling, certificate errors and so on.

CefRequestHander接口还提供了其他回调函数以定制其他网络相关事件。包括授权、cookei处理、外部协议处理、证书错误等。

##### Proxy Resolution

Proxy settings are configured in CEF3 using the same command-line flags as Google Chrome.
CEF3使用类似Google Chrome一样的方式，通过命令行参数传递代理配置。

```
--proxy-server=host:port
      Specify the HTTP/SOCKS4/SOCKS5 proxy server to use for requests. An individual proxy
      server is specified using the format:

        [<proxy-scheme>://]<proxy-host>[:<proxy-port>]

      Where <proxy-scheme> is the protocol of the proxy server, and is one of:

        "http", "socks", "socks4", "socks5".

      If the <proxy-scheme> is omitted, it defaults to "http". Also note that "socks" is equivalent to
      "socks5".

      Examples:

        --proxy-server="foopy:99"
            Use the HTTP proxy "foopy:99" to load all URLs.

        --proxy-server="socks://foobar:1080"
            Use the SOCKS v5 proxy "foobar:1080" to load all URLs.

        --proxy-server="sock4://foobar:1080"
            Use the SOCKS v4 proxy "foobar:1080" to load all URLs.

        --proxy-server="socks5://foobar:66"
            Use the SOCKS v5 proxy "foobar:66" to load all URLs.

      It is also possible to specify a separate proxy server for different URL types, by prefixing
      the proxy server specifier with a URL specifier:

      Example:

        --proxy-server="https=proxy1:80;http=socks4://baz:1080"
            Load https://* URLs using the HTTP proxy "proxy1:80". And load http://*
            URLs using the SOCKS v4 proxy "baz:1080".

--no-proxy-server
      Disables the proxy server.

--proxy-auto-detect
      Autodetect  proxy  configuration.

--proxy-pac-url=URL
      Specify proxy autoconfiguration URL.
```

If the proxy requires authentication the CefRequestHandler::GetAuthCredentials() callback will be executed with an |isProxy| value of true to retrieve the username and password.

如果代理请求授权，CefRequestHandler::GetAuthCredentials()回调会被调用。如果isProxy参数为true，则需要返回用户名和密码。

```
bool MyHandler::GetAuthCredentials(
    CefRefPtr<CefBrowser> browser,
    CefRefPtr<CefFrame> frame,
    bool isProxy,
    const CefString& host,
    int port,
    const CefString& realm,
    const CefString& scheme,
    CefRefPtr<CefAuthCallback> callback) {
  if (isProxy) {
    // Provide credentials for the proxy server connection.
    callback->Continue("myuser", "mypass");
    return true;
  }
  return false;
}
```

Web content loading during application startup can be delayed due to network proxy resolution (for example, if "Automatically detect proxy settings" is checked on Windows). For best user experience consider designing your application to first show a static splash page and then redirect to the actual website using meta refresh -- the redirect will be blocked until proxy resolution completes. For testing purposes proxy resolution can be disabled using the "--no-proxy-server" command-line flag. Proxy resolution delays can also be duplicated in Google Chrome by running "chrome --url=..." from the command line. 
网络内容加载可能会因为代理而有延迟。为了更好的用户体验，可以考虑让你的应用程序先显示一个闪屏，等内容加载好了再通过meta refresh显示真实网页。可以指定`--no-proxy-server`禁用代理并做相关测试。代理延迟也可以通过chrome浏览器重现，方式是使用命令行传参：`chrome -url=...`.


