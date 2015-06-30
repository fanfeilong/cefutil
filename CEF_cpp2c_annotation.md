CEF Cpp2c Anotation
===================

CEF 使用工具生成的Cpp和C接口互相转换，基本流程是这样的
1. Cef源码使用C++开发
2. 通过libcef_dll/cpptoc导出C接口
3. 通过libcef_dll/ctocpp对导出的C接口反向做一层C++封装，让C++用户直接使用C++接口

CEF cpptoc
----------
使用cpptoc导出c接口，该部分位于libcef_dll/cpptoc子目录下
首先看下基类，涉及到两个文件
1. libcef_dll/cpptoc/cpptoc.h
2. libcef_dll/cpptoc/base_cpptoc.h

base_cpptoc.h只定义了一个继承自CppToC的类：
```
class CefBaseCppToC
    : public CefCppToC<CefBaseCppToC, CefBase, cef_base_t> {
 public:
  CefBaseCppToC();
};
```

关键看cpptoc.h，该类是一个模板类
```
template <class ClassName, class BaseName, class StructName>
class CefCppToC : public CefBase {
public:
	static StructName* Wrap(CefRefPtr<BaseName> c){
		if (!c.get())
		  return NULL;

		// Wrap our object with the CefCppToC class.
		ClassName* wrapper = new ClassName();
		wrapper->wrapper_struct_.object_ = c.get();
		// Add a reference to our wrapper object that will be released once our
		// structure arrives on the other side.
		wrapper->AddRef();
		// Return the structure pointer that can now be passed to the other side.
		return wrapper->GetStruct();
	}
	static CefRefPtr<BaseName> Unwrap(StructName* s){
	    if (!s)
		  return NULL;

		// Cast our structure to the wrapper structure type.
		WrapperStruct* wrapperStruct = GetWrapperStruct(s);

		// If the type does not match this object then we need to unwrap as the
		// derived type.
		if (wrapperStruct->type_ != kWrapperType)
		  return UnwrapDerived(wrapperStruct->type_, s);

		// Add the underlying object instance to a smart pointer.
		CefRefPtr<BaseName> objectPtr(wrapperStruct->object_);
		// Release the reference to our wrapper object that was added before the
		// structure was passed back to us.
		wrapperStruct->wrapper_->Release();
		// Return the underlying object instance.
		return objectPtr
	}
	static CefRefPtr<BaseName> Get(StructName* s){
		DCHECK(s);
		WrapperStruct* wrapperStruct = GetWrapperStruct(s);
		// Verify that the wrapper offset was calculated correctly.
		DCHECK_EQ(kWrapperType, wrapperStruct->type_);
		return wrapperStruct->object_
	}

public:
	StructName* GetStruct(){
		return &wrapper_struct_.struct_;
	}
	void AddRef();
	bool Release();
	bool HasOneRef();

protected:
	CefCppToC(){
		wrapper_struct_.type_ = kWrapperType;
		wrapper_struct_.wrapper_ = this;
		memset(GetStruct(), 0, sizeof(StructName));

		cef_base_t* base = reinterpret_cast<cef_base_t*>(GetStruct());
		base->size = sizeof(StructName);
		base->add_ref = struct_add_ref;
		base->release = struct_release;
		base->has_one_ref = struct_has_one_ref;
	}
	virtual ~CefCppToC();

private:
  	struct WrapperStruct {
    	CefWrapperType type_;
    	BaseName* object_;
    	CefCppToC<ClassName, BaseName, StructName>* wrapper_;
    	StructName struct_;
  	};
	static WrapperStruct* GetWrapperStruct(StructName* s) {
		// Offset using the WrapperStruct size instead of individual member sizes
		// to avoid problems due to platform/compiler differences in structure
		// padding.
		return reinterpret_cast<WrapperStruct*>(
			reinterpret_cast<char*>(s) -
			(sizeof(WrapperStruct) - sizeof(StructName)));
  	}
	static CefRefPtr<BaseName> UnwrapDerived(CefWrapperType type, StructName* s);
	void UnderlyingAddRef() const{
		wrapper_struct_.object_->AddRef();
	}
	bool UnderlyingRelease() const;
	bool UnderlyingHasOneRef() const;
	static void CEF_CALLBACK struct_add_ref(cef_base_t* base){
		DCHECK(base);
		if (!base)
		  return;

		WrapperStruct* wrapperStruct =
			GetWrapperStruct(reinterpret_cast<StructName*>(base));
		// Verify that the wrapper offset was calculated correctly.
		DCHECK_EQ(kWrapperType, wrapperStruct->type_);

		wrapperStruct->wrapper_->AddRef();
	}
	static int CEF_CALLBACK struct_release(cef_base_t* base);
	static int CEF_CALLBACK struct_has_one_ref(cef_base_t* base);

private:
	WrapperStruct wrapper_struct_;
  	CefRefCount ref_count_;
	static CefWrapperType kWrapperType;
}
```
1. 可以看到，CppToC通过模板参数指定要封装的类及其父类的名字ClassName、BaseName，最后一个模板参数是要转换的C结构体的名字
2. CppToC内部持有一个WrapperStruct结构体成员变量wrapper_struct_,该结构体包含了双向转化所需的信息： 
    - WrapperStruct则持有了CEF C++对象BaseName* object_
	- WrapperStruct则持有了CEF C对象StructName struct_
3. CppToC在构造函数里初始化C结构体变量wrapper_struct_.struct_，把该C对象的函数指针都赋值好
    ```
	cef_base_t* base = reinterpret_cast<cef_base_t*>(GetStruct());
	base->size = sizeof(StructName);
	base->add_ref = struct_add_ref;
	base->release = struct_release;
	base->has_one_ref = struct_has_one_ref;
	```
	其中struct_add_ref、struct_release、struct_has_one_ref等静态函数分别转调用UnderlyingAddRef、UnderlyingRelease、UnderlyingHasOneRef，这三个成员函数则实际上调用的是
	wrapper_struct_.object的AddRef、Release、HasOneRef方法
4. 实际上wrapper_struct_.struct_指向的是include/capi/子目录里定义的各种c结构体，他们都有共同的基类结构体cef_base_t，所以CppToC这里把wrapper_struct_.object转成cef_base_t*来初始化，而子类结构体的函数指针则交给
CppToC的子类去设置

所以CppToC的原理就三条
1. 每个CEF类有一个对应的C结构体，C结构体里为C++类的public方法定义一堆对应的函数指针
2. 每个CEF类都有一个对应的CppToC的子类，而CppToC内部会创建一个WrapperStruct，同时持有了C++类和C结构体,通过GetStruct()方法获取C对象wrapper_struct_.struct_，通过GetWrappperStruct(StructName*)从C对象获取WrappperStruct，进而获取C++对象wrapper_struct_.object_
3. CppToC的构造函数以及子类的构造函数负责将wrapper_struct_.struct_的函数指针赋值为内部的静态成员函数，这些静态成员函数内部通过GetWrapperStruct得到wrapper_struct_,进而得到C++对象wrapper_struct_.object_，从而调用C++的接口

在理解了以上内容后，再去看cpptoc目录下的其他CppToC子类的代码就没很清晰了。

CEF ctocpp
----------
