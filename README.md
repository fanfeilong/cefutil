CEF全称是Chromium Embedded Framework，它是Chromium的Content API的封装库。
CEF官网地址：https://bitbucket.org/chromiumembedded/cef

CEF Util 货币 系统
------------------
CEF Util以文件+文件Task List的形式组织协调开发。每个文件`File` 对应一个`File Task List.md` 文件，上面列出了该文件的任务分割，每个开发者领取了任务时，自己在`FileName Task List.md`里 将自己的名字以 `@YourName` 添加到任务名词之前，完成第一次`Fork+Pull Request`后，会在该条目后面列出获得的cef货币，CEF货币是一种口头约定货币系统，根据任务的不同会发布不同的币种，比如翻译任务会发布`$(Translate)`币，该币可以用来在CEFUitl项目里发起悬赏翻译，或者可以用来在自己简历上列出，以表征自己在CEFUtil项目里的贡献。实际上，根据不同的任务，会有各种有趣的币被发布。

CEF 文档中文翻译任务
--------------------
- [x] [CEF General Usage](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF%20General%20Usage.md)
- [x] [CEF General Usage Task List](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF%20General%20Usage%20Task%20List.md)
- [x] [CEF General Usage中文版，欢迎查阅](https://github.com/fanfeilong/cefutil/blob/master/doc/CEF%20General%20Usage-zh-cn.md)

Chromium Documentations
-----------------
- [Chromium Build Guid](https://github.com/fanfeilong/cefutil/blob/master/doc/chromium_build_guid.md)
- [Chromium Generation Your Project](https://github.com/fanfeilong/cefutil/blob/master/doc/gyp.md)
- [Chromium Resoures](https://github.com/fanfeilong/cefutil/blob/master/doc/chromium_resources.md)
- [Chromium Content Register V8 Extension](https://github.com/fanfeilong/cefutil/blob/master/doc/content_register_v8_extension.md)

CEF FAQ
-------

#### 如何做进程间同步通信
参考CEF官方论坛的[这个帖子](http://www.magpcss.org/ceforum/viewtopic.php?f=6&t=10680)，尽量不要在CEF里做进程间同步通信，但如果你非要做。以下步骤是一个方案：

1. Start the call in the Browser process by "SendMessage".
2. In the BrowserProcess, configure a wait handle on an event until it gets signaled and wait.
3. The Render process receives the sent message from the Browser process.
4. The Render process calls the TryEval to execute the JavaScript.
5. The return code of the evaluation is being sent back to the Browser process, marshaled into a CefProcessMessage.
6. The Browser process receives the message, un-marshals the return code inside the message.
7. The event object is being signaled.
8. The Browser process returns to the caller the result of the JavaScript call.

如果是在Render进程首先发起调用，不要在CEFV8Handler的Execute里阻塞住执行流，参考CEF官方论坛的[这个使用Task的帖子](http://www.magpcss.org/ceforum/viewtopic.php?f=14&t=11132)。
