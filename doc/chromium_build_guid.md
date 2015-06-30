Windows7 64位系统编译Chromium指南
=================================

1. 代理相关问题，以ssh隧道代理为例：
	- 配置bitvise，设置socket5代理。指定socket代理端口号127.0.0.1:xxxx;
    - 配置privoxy，设置http代理。
		- 转发socket5：forward-socket5 127.0.0.1:xxxx  .
		- 监听端口：listen-addrenss 127.0.0.1:yyyy .
	- 命令行下netsh winhttp set proxy 127.0.0.1:yyyy “<local>” 设置winhttp代理。供gclient下载工具链时使用。
    - 环境变量添加http_proxy，值为127.0.0.1:yyyy;这是因为chromium的Python下载脚本会用到这个环境变量设置代理。
    - 使用git前可以通过git config --global http.proxy=127.0.0.1:yyyy配置代理。
    - 如果有vps，一定直接在vsp上下载chromium源码，打包成tar，在公司网络下直接拉取chromium源码简直是找死。

2. 安装64位windows7，本指南可能不适用于Windows7 32位系统，请自重。
3. 安装VisualStudio2010，设置环境变量GYP_MSVS_VERSION=2010
4. 安装VisualStudio2010 SP1。
5. 安装windows 8.0 sdk（不要安装8.1），并修改Windows Kits\8.0\Include\WinRT\asyncinfo.h
   注释第66行的class关键字：

   ```
   65 namespace ABI (namespace Windows { namespace Foundation { 
   66 enum /*class*/ AsyncStatus {                
   ```
   
6. 安装June 2010 DirectX SDK。如果有遇到”Error Code:S1023”，则通过控制面板将
   Microsoft Visual C++ 2010 x86 Redistributable
   Microsoft Visual C++ 2010 x64 Redistributable
   都卸载掉，重启后重新安装June 2010 DirectX SDK。
7. 下载depot_tools，解压，并将depot_tools根目录添加到Path环境变量。
   为了让depot_tools正常可用，可能需要配置好代理，这个尤为重要，很多人都死在这。
8. 在cmd命令行下执行命令： gclient 
   该命令会优先下载git、svn、python并解压到depot_tools目录下，并设置环境变量。
9. 下载chromium源码，可以直接下载chromium的源码压缩包也可以通过gclient直接下载。
- 下载最新的chromium源码压缩包，并且最好用7z解压：
```
http://chromium-browser-source.commondatastorage.googleapis.com/chromium_tarball.html
```
- 下载指定版本的chromium源码压缩包：
```
http://chromium-browser-source.commondatastorage.googleapis.com/chromium.rXXXXX.tgz
其中rXXXXX表示版本号，比如r197479表示Revision197479。
```
- 所有可用的压缩包版本号列表页面是：
```
http://chromium-browser-source.commondatastorage.googleapis.com/
```
- 通过gclient下载源码。
创建存放chromium源码的路径，路径不要带空格，比如d:\dev\chromium。
打开cmd命令行，cd到chromium源码路径。
执行gclient config --git-deps https://chromium.googlesource.com/chromium/src.git
执行gclient sync 下载chromium源码。
