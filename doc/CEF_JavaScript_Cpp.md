**⚠️ API 时效性警告**

本文档展示的 JavaScript 与 C++ 交互方法基于较早版本的 CEF API。随着 CEF 版本演进（特别是 CEF 100+ 版本），部分 API 可能已发生变化。

**🔥 推荐优先参考现代化指南：**
- 🚀 [现代开发指南 (2025)](modern_development_guide_2025.md) - 最新 API 使用方法
- ✨ [CEF 最新特性指南 (2025)](cef_features_2025.md) - 现代 JavaScript 绑定特性
- ⚡ [性能优化最佳实践 (2025)](performance_optimization_2025.md) - JavaScript 性能优化
- 🔒 [安全开发指南 (2025)](security_best_practices_2025.md) - JavaScript 安全最佳实践

**官方资源：**
- [最新 CEF API 文档](https://bitbucket.org/chromiumembedded/cef)
- [官方 JavaScript 绑定指南](https://bitbucket.org/chromiumembedded/cef/wiki/JavaScriptIntegration)
- [CEF 论坛相关讨论](https://magpcss.org/ceforum/viewforum.php?f=6)

**适用版本：** CEF3 早期到中期版本
**验证日期：** 2025年7月4日（建议使用现代化 API）

---

##### <a name="javascirpt-custome-handler"></a>JavaScript和Cpp交互示例(Custom Implementation)

一个CEF应用程序也可以提供自己的异步JavaScript绑定。

此处演示：
- JavaScript注册函数给Render进程，Render进程保存该JavaScript函数
- Render进程发消息通知Browser进程
- Browser进程处理后，回发消息给Render进程
- Render进程调用之前保存的JavaScript函数

#### 步骤
1. 首先在CefRenderProcessHandler的子类里覆写虚方法OnWebKitInitialized，并在该方法的实现里注册一个C++方法给JavaScript
	```
	//假设CefRenderProcessHandler的子类为CefRenderProcessHandlerImpl
	void CefRenderProcessHandlerImpl::OnWebKitInitialized(){
		std::string app_code =
		//-----------------------------------
		//声明JavaScript里要调用的Cpp方法
		"var app;"
		"if (!app)"
		"  app = {};"
		"(function() {"

		//  Send message 
		"  app.sendMessage = function(name, arguments) {"
		"    native function sendMessage();"
		"    return sendMessage(name, arguments);"
		"  };"

		// Registered Javascript Function, which will be called by Cpp
		"  app.registerJavascriptFunction = function(name,callback) {"
		"    native function registerJavascriptFunction();"
		"    return registerJavascriptFunction(name,callback);"
		"  };"

		"})();";
		//------------------------------------

		// Register app extension module
		
		// JavaScript里调用app.registerJavascriptFunction时，就会去通过CefRegisterExtension注册的CefV8Handler列表里查找
		// 找到"v8/app"对应的CefV8HandlerImpl，就调用它的Execute方法
		// 假设m_v8Handler是CefRenderProcessHandlerImpl的一个成员变量
		m_v8Handler = new CefV8HandlerImpl();
		CefRegisterExtension("v8/app", app_code,m_v8Handler);
	}
	```
2. 在CefV8Handler的子类的Execute方法里实现sendMessage和registerJavascriptFunction
	```
	// in CefV8HandlerImpl.h
	class CefV8HandlerImpl : publice CefV8Handler{
	public:
		CefV8HandlerImpl();
		~CefV8HandlerImpl();
	public:
		/**
		 *	CefV8Handler Methods, Which will be called when the V8 extension 
		 *  is called in the Javascript environment
		 */
		virtual bool Execute(const CefString& name
			,CefRefPtr<CefV8Value> object
			,const CefV8ValueList& arguments
			,CefRefPtr<CefV8Value>& retval
			,CefString& exception);
	private:
		// Map of message callbacks.
	typedef std::map<std::pair<std::string, int>,
					 std::pair<CefRefPtr<CefV8Context>, CefRefPtr<CefV8Value> > >
					 CallbackMap;
	CallbackMap callback_map_;
	}
	```
	```
	CefV8HandlerImpl::CefV8HandlerImpl()
	{

	}

	CefV8HandlerImpl::~CefV8HandlerImpl()
	{
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

	// in CefV8HandlerImpl.cpp
	bool CefV8HandlerImpl::Execute(const CefString& name  //JavaScript调用的C++方法名字
			,CefRefPtr<CefV8Value> object                 //JavaScript调用者对象
			,const CefV8ValueList& arguments              //JavaScript传递的参数
			,CefRefPtr<CefV8Value>& retval                //需要返回给JavaScript的值设置给这个对象
			,CefString& exception)                        //通知异常信息给JavaScript
	{
		// In the CefV8Handler::Execute implementation for “registerJavascriptFunction”.
		bool handle=false;
		if(name=="registerJavascriptFunction"){
			//保存JavaScript设置的回答函数
			if (arguments.size() == 2 && arguments[0]->IsString() &&
				arguments[1]->IsFunction()) {
			  std::string message_name = arguments[0]->GetStringValue();
			  CefRefPtr<CefV8Context> context = CefV8Context::GetCurrentContext();
			  int browser_id = context->GetBrowser()->GetIdentifier();
			  callback_map_.insert(
				  std::make_pair(std::make_pair(message_name, browser_id),
								 std::make_pair(context, arguments[1])));
			handle = true;
			
			// 此时可以发送一个信息给Browser进程，见第4个步骤
		}

		if(name=="sendMessage"){
			//处理SendMessage，
			handle = true;
		}

		// 如果没有处理，则发异常信息给js
		if(!handle){
			exception="not implement function";
		}

		return true;
	}
	```
3. 在HTML的JavaScript里，通过上面注册的方法向Render进程注册一个回调函数。
	
	```
	// In JavaScript register the callback function.
	app.setMessageCallback('binding_test', function(name, args) {
	  document.getElementById('result').value = "Response: "+args[0];
	});
	```
4. Render进程发送异步进程间通信到Browser进程。
5. Browser进程接收到进程间消息，并处理。
6. Browser进程处理完毕后，发送一个异步进程间消息给Render进程，返回结果。
7. Render进程接收到进程间消息，则调用最开始保存的JavaScript注册的回调函数处理之。
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
8. 在CefRenderProcessHandlerImpl::OnContextReleased()里释放JavaScript注册的回调函数以及其他V8资源。

	```
	void CefRenderProcessHandlerImpl::OnContextReleased(CefRefPtr<CefBrowser> browser,
									  CefRefPtr<CefFrame> frame,
									  CefRefPtr<CefV8Context> context) {
		if(m_v8Handler->Get()){
			m_v8Handler->Release();
		}
	}
	```
