CEF Util 货币 xi
- [ ] $(DEBUG)
- [ ] $(DEBUG)
- [ ] $(DEBUG)
- [ ] $(DEBUG)
- [ ] $(DEBUG)
- [ ] $(Translate)
- [ ] $(Translate)
- [ ] $(Translate)
- [ ] $(Translate)
- [ ] $(Translate)

CEF FAQ
========================

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
