
GYP(Generate Your Projects)
===========================

目录

1.	说明	1
2.	介绍	1
3.	示例骨架（Skeleton）	1
	- 3.1.	Chromium	2
	- 3.2.	Execute Target	3
	- 3.3.	Library Target	5
4.	用例	7
	- 4.1.	添加一个源文件（Source File）	7
	    - 4.1.1.	添加全平台源文件	7
		- 4.1.2.	添加特定平台源文件	7
	- 4.2.	添加一个可执行目标项（Execute Target）	9
		- 4.2.1.	全平台	9
		- 4.2.2.	特定平台	10
	- 4.3.	添加一个库目标项（Library Target）	10
		- 4.3.1.	全平台	11
		- 4.3.2.	特定平台	11
	- 4.4.	为目标项(Target)添加设置（Settings）	12
		- 4.4.1.	添加预处理宏定义	12
		- 4.4.2.	添加头文件包含目录	13
		- 4.4.3.	添加编译flags	13
		- 4.4.4.	添加链接flags	14
		- 4.4.5.	添加指定平台设置	14
	- 4.5.	配置目标项（Target）之间的依赖（Dependencies）	15
		- 4.5.1.	添加依赖的链接库	15
		- 4.5.2.	指定依赖者所需的设置	16
	- 4.6.	为Mac OS X bundles添加支持	17
	- 4.7.	在重构中正确的移除文件	17
	- 4.8.	定制编译步骤	17
5.	语言规格	18
	- 5.1.	目标	18
	- 5.2.	背景	18
	- 5.3.	预览	18
	- 5.4.	设计细节	19
		- 5.4.1.	顶级元素	20
		- 5.4.2.	目标项（Targets）	20
		- 5.4.3.	配置（Configurations）	21
		- 5.4.4.	条件（Conditionals）	22
		- 5.4.5.	行为（Actions）	22
		- 5.4.6.	规则（Rules）	23
		- 5.4.7.	拷贝（Copies）	24
		- 5.4.8.	生成XCode .pbxproj 文件	24
		- 5.4.9.	生成Visual Studio .vcxproj文件	24
		- 5.4.1.0	生成Visual Studio .sln文件	24
6.	输入格式参考	24

## 说明

本文档翻译GYP的相关文档，并做合并编辑工作。文档原始出处如下：


-	GYPUserDocumentation

https://code.google.com/p/gyp/wiki/GypUserDocumentation

GYP (Generate Your Projects) 
Status: Draft (as of 2009-05-19)

Mark Mentovai <mark@chromium.org>, Steven Knight <sgk@chromium.org> _et al._ 
Modified: 2009-05-19
-	GYPLanguageSpecification

https://code.google.com/p/gyp/wiki/GypLanguageSpecification

GYP (Generate Your Projects) 
Status: Draft (as of 2009-01-30)

Mark Mentovai <mark@chromium.org> _et al._ 
Modified: 2009-02-10
## 介绍

GYP是Google开发的跨平台项目工程生成系统。本文文档第2-4章节通过实例和用例介绍GYP的基本使用。这部分旨在提供用户级别的GYP使用指南。第5章则为GYP配置文件的语言规格做详细的说明。



GYP项目同Google的其他开源项目一样，托管在Google Code上：

https://code.google.com/p/gyp/
## 示例骨架（Skeleton）

本节介绍.gyp文件的基本骨架。
### Chromium

下面是一个Chromium项目的GYP配置骨架：
```
  {
    'variables': {
      .
      .
      .
    },
    'includes': [
      '../build/common.gypi',
    ],
    'target_defaults': {
      .
      .
      .
    },
    'targets': [
      {
        'target_name': 'target_1',
          .
          .
          .
      },
      {
        'target_name': 'target_2',
          .
          .
          .
      },
    ],
    'conditions': [
      ['OS=="linux"', {
        'targets': [
          {
            'target_name': 'linux_target_3',
              .
              .
              .
          },
        ],
      }],
      ['OS=="win"', {
        'targets': [
          {
            'target_name': 'windows_target_4',
              .
              .
              .
          },
        ],
      }, { # OS != "win"
        'targets': [
          {
            'target_name': 'non_windows_target_5',
              .
              .
              .
          },
      }],
    ],
  }
```


可以看到，一个gyp文件里面实际上就是一个Python字典（实际上除了引入Python风格注释，在字典和列表后面可以添加逗号外，跟Json并无区别）。从上面的例子我们看到，gyp配置的顶级元素有下面几种：



节点元素-说明-节点值类型

`variables`-定义可在本文件其他部分使用的变量-Python字典

`includes`-本文件所包含的其他.gypi文件列表-Python列表

`target_defaults`-对本文件所有targets都生效的设置-Python字典

`targets`-本文件所有编译项目设置，每个Target一个字典-Python列表

`conditions`-根据全局变量的值修改gyp其他定义节点-Python列表


### Execute Target

最简单的目标项（Target）就是一个简单的可执行程序。下面是一个可执行程序的.gyp配置。


```
{
    'targets': [
      {
        'target_name': 'foo',
        'type': 'executable',
        'msvs_guid': '5ECEC9E5-8F23-47B6-93E0-C3B328B3BE65',
        'dependencies': [
          'xyzzy',
          '../bar/bar.gyp:bar',
        ],
        'defines': [
          'DEFINE_FOO',
          'DEFINE_A_VALUE=value',
        ],
        'include_dirs': [
          '..',
        ],
        'sources': [
          'file1.cc',
          'file2.cc',
        ],
        'conditions': [
          ['OS=="linux"', {
            'defines': [
              'LINUX_DEFINE',
            ],
            'include_dirs': [
              'include/linux',
            ],
          }],
          ['OS=="win"', {
            'defines': [
              'WINDOWS_SPECIFIC_DEFINE',
            ],
          }, { # OS != "win",
            'defines': [
              'NON_WINDOWS_DEFINE',
            ],
          }]
        ],
      },
    ],
  }
```


从上面的配置可以看到，一个典型的target一般包含下表所示的子节点：

节点元素-说明-节点值类型

`target_name`-目标项名称，必须在所有.gyp文件中都唯一。该名字用来生成Visual Studio项目名；XCode 目标项名字；SCons目标项别名。-字符串

`type`-目标项类型，此例中为'executable'代表可执行程序。-`executeable`

`msvs_guid`-唯一标识目标项，生成Visual Studio sln文件所需-GUID

`dependencies`-列出依赖的其他Targets名称列表。
1、	所依赖的Targets会先被编译。
2、	如果所依赖的Target是一个库，会被链接到本项目。

所依赖的Targets的direct_dependent_settings节点的设置会被应用到本Target的编译和链接。参考下direct_dependent_settings的说明-Python列表

`defines`-C预处理定义列表，被加入到编译命令行选项（-D或/D）-Python列表

`include_dirs`-头文件包含路径，被加入到编译命令行选项（-I或/I）-Python列表

`sources`-源文件列表-Python列表

`conditions`-根据全局变量的值添加设置-Python列表


### Library Target

除了可执行程序目标项目，更多的是目标项目是类库。下面是一个类库目标项目的gyp设置示例，包含了大部分类库所需的gyp设置。


```
{
    'targets': [
      {
        'target_name': 'foo',
        'type': '<(library)'
        'msvs_guid': '5ECEC9E5-8F23-47B6-93E0-C3B328B3BE65',
        'dependencies': [
          'xyzzy',
          '../bar/bar.gyp:bar',
        ],
        'defines': [
          'DEFINE_FOO',
          'DEFINE_A_VALUE=value',
        ],
        'include_dirs': [
          '..',
        ],
        'direct_dependent_settings': {
          'defines': [
            'DEFINE_FOO',
            'DEFINE_ADDITIONAL',
          ],
          'linkflags': [
          ],
        },
        'export_dependent_settings': [
          '../bar/bar.gyp:bar',
        ],
        'sources': [
          'file1.cc',
          'file2.cc',
        ],
        'conditions': [
          ['OS=="linux"', {
            'defines': [
              'LINUX_DEFINE',
            ],
            'include_dirs': [
              'include/linux',
            ],
          ],
          ['OS=="win"', {
            'defines': [
              'WINDOWS_SPECIFIC_DEFINE',
            ],
          }, { # OS != "win",
            'defines': [
              'NON_WINDOWS_DEFINE',
            ],
          }]
        ],
    ],
  }
```


从上面的配置可以看到，和可执行目标项相比，类库目标项的gyp节点多出来几个：



节点类型-说明-节点值类型说明

`type`-库项目类型，一般定义为`<(library)`，如果需要指明静态或动态，则可以显示使用`static_library`或`shared_library`-`<(library)`

`static_library`

`shared_library`

`direct_dependent_settings`-该节点的设置会被应用到那些直接依赖于该Target的Targets。

参考`dependenties`节点说明。-Python字典

`export_dependent_settings`-导出Targets列表，这些Target的

`direct_dependent_settings`节点将会被

应用到那些直接依赖该Target的Targets-Python列表


## 用例

本章详细说明典型用例，这些用例并非完整的gyp文件，只列出每个用例相关的gyp片段。
### 添加一个源文件（Source File）

一个源文件可能是跨平台的，也可能是平台相关的，gyp分别为这两种需求提供了指定源文件的方法。
**添加全平台源文件**:
```
  {
    'targets': [
      {
        'target_name': 'my_target',
        'type': 'executable',
        'sources': [
          '../other/file_1.cc',
          'new_file.cc',
          'subdir/file3.cc',
        ],
      },
    ],
  },
```


跨全平台文件直接在`targets`/`sources`里列出。最好是以字母序排序。文件的路径是gyp文件所在地路径的相对路径。
#### 添加特定平台源文件
##### 指定后缀

为特定平台添加源文件最简单的方式是为源文件添加指定后缀，并在gyp文件里添加过滤规则。下表列出了特定平台源文件的约定后缀：



平台-后缀-示例

Linux-_linux-foo_linux.cc

Mac-_mac-foo_mac.cc

Posix-_posix-foo_posix.cc

Win-_win-foo_win.cc

gyp片段如下：
```
  {
    'targets': [
      {
        'target_name': 'foo',
        'type': 'executable',
        'sources': [
          'independent.cc',
          'specific_win.cc',
        ],
      },
    ],
  },
```


Chromium的gyp文件都有添加这些后缀源文件的过滤规则，比如本例的specific_win.cc文件会自动在非windows平台上被过滤掉。
##### 条件源文件

如果源文件没有这些特定的后缀，则可以通过`conditions`节点来特化平台相关的源文件，包括排除平台和指定平台两种办法。


-	排除平台（推荐做法）
```
  {
    'targets': [
      {
        'target_name': 'foo',
        'type': 'executable',
        'sources': [
          'linux_specific.cc',
        ],
        'conditions': [
          ['OS != "linux"', {
            'sources!': [ #注意这里是`sources!`
              # Linux-only; exclude on other platforms.
              'linux_specific.cc',
            ]
          }],
        ],
      },
    ],
  },
```

-	指定平台（不推荐）
```
  {
    'targets': [
      {
        'target_name': 'foo',
        'type': 'executable',
        'sources': [],
        ['OS == "linux"', {
          'sources': [
            # Only add to sources list on Linux.
            'linux_specific.cc',
          ]
        }],
      },
    ],
  },
```

### 添加一个可执行目标项（Execute Target）
**全平台**:

```
{
    'targets': [
      {
        'target_name': 'new_unit_tests',
        'type': 'executable',
        'defines': [
          'FOO',
        ],
        'include_dirs': [
          '..',
        ],
        'dependencies': [
          'other_target_in_this_file',
          'other_gyp2:target_in_other_gyp2',
        ],
        'sources': [
          'new_additional_source.cc',
          'new_unit_tests.cc',
        ],
      },
    ],
  }
```

**特定平台**

```
{
    'targets': [
    ],
    'conditions': [
      ['OS=="win"', {
        'targets': [
          {
            'target_name': 'new_unit_tests',
            'type': 'executable',
            'defines': [
              'FOO',
            ],
            'include_dirs': [
              '..',
            ],
            'dependencies': [
              'other_target_in_this_file',
              'other_gyp2:target_in_other_gyp2',
            ],
            'sources': [
              'new_additional_source.cc',
              'new_unit_tests.cc',
            ],
          },
        ],
      }],
    ],
  }
```

### 添加一个库目标项（Library Target）


**全平台**:

```
  {
    'targets': [
      {
        'target_name': 'new_library',
        'type': '<(library)',
        'defines': [
          'FOO',
          'BAR=some_value',
        ],
        'include_dirs': [
          '..',
        ],
        'dependencies': [
          'other_target_in_this_file',
          'other_gyp2:target_in_other_gyp2',
        ],
        'direct_dependent_settings': {
          'include_dirs': '.',
        },
        'export_dependent_settings': [
          'other_target_in_this_file',
        ],
        'sources': [
          'new_additional_source.cc',
          'new_library.cc',
        ],
      },
    ],
  }
```


需要注意的是，`type`默认使用`<(library)`，允许用户在gyp全局配置里设置具体使用静态还算共享库。或者可以固定使用静态或动态连接库：`static_library`或`shared_library`
**特定平台**:

```
  {
    'targets': [
    ],
    'conditions': [
      ['OS=="win"', {
        'targets': [
          {
            'target_name': 'new_library',
            'type': '<(library)',
            'defines': [
              'FOO',
              'BAR=some_value',
            ],
            'include_dirs': [
              '..',
            ],
            'dependencies': [
              'other_target_in_this_file',
              'other_gyp2:target_in_other_gyp2',
            ],
            'direct_dependent_settings': {
              'include_dirs': '.',
            },
            'export_dependent_settings': [
              'other_target_in_this_file',
            ],
            'sources': [
              'new_additional_source.cc',
              'new_library.cc',
            ],
          },
        ],
      }],
    ],
  }
```

### 为目标项(Target)添加设置（Settings）

本节详细说明常用目标项设置的配置，本节gyp代码都是片段代码。
#### 添加预处理宏定义

预处理宏定义添加到defines节点：

```
  {
    'targets': [
      {
        'target_name': 'existing_target',
        'defines': [
          'FOO',
          'BAR=some_value',
        ],
      },
    ],
  },
```

#### 添加头文件包含目录

头文件包含目录添加到`include_dirs`节点，该节点可能直接在在target下也可以在`conditions`节点下。

```
  {
    'targets': [
      {
        'target_name': 'existing_target',
        'include_dirs': [
          '..',
          'include',
        ],
      },
    ],
  },
```

#### 添加编译flags

编译flags添加到`cflags`节点，由于编译flags一般都算平台相关的，所以cflags节点一般在conditions节点下。

```
  {
    'targets': [
      {
        'target_name': 'existing_target',
        'conditions': [
          ['OS=="win"', {
            'cflags': [
              '/WX',
            ],
          }, { # OS != "win"
            'cflags': [
              '-Werror',
            ],
          }],
        ],
      },
    ],
  },
```

#### 添加链接flags

链接flags是链接器相关的。



在linux和非mac的posix系统上，添加到`ldflags`节点上。

```
  {
    'targets': [
      {
        'target_name': 'existing_target',
        'conditions': [
          ['OS=="linux"', {
            'ldflags': [
              '-pthread',
            ],
          }],
        ],
      },
    ],
  },
```


在OS X系统上，添加到`xcode_settings`节点上。

在Windows系统上，添加到`msvs_settings`节点上。


#### 添加指定平台设置

每个设置关键字都可以添加感叹号后缀用以移除某些设置，比如chromium的common.gypi里将Linux平台的 -Werror编译flags移除：

```
  {
    'targets': [
      {
        'target_name': 'third_party_target',
        'conditions': [
          ['OS=="linux"', {
            'cflags!': [
              '-Werror',
            ],
          }],
        ],
      },
    ],
  },
```



### 配置目标项（Target）之间的依赖（Dependencies）

gyp系统提供了指定目标项之间互相依赖的配置方案。下面分别说明需要配置的步骤。
#### 添加依赖的链接库

在`dependencies`节点添加依赖的链接库：

```
  {
    'targets': [
      {
        'target_name': 'foo',
        'dependencies': [
          'libbar',
        ],
      },
      {
        'target_name': 'libbar',
        'type': '<(library)',
        'sources': [
        ],
      },
    ],
  }
```

如果所依赖的目标项在另外一个gyp文件，则需要指明gyp文件的相对路径：

```
  {
    'targets': [
      {
        'target_name': 'foo',
        'dependencies': [
          '../bar/bar.gyp:libbar',
        ],
      },
    ],
  }
```

如果更新了bar.gyp的名字，则需要更新所有依赖于bar.gyp的地方。
#### 指定依赖者所需的设置

目标项可以在`direct_dependent_settings`节点添加依赖者所需设置，这些设置会被链接到依赖者，形成最终的配置文件。

```
{
    'targets': [
      {
        'target_name': 'foo',
        'type': 'executable',
        'dependencies': [
          'libbar',
        ],
      },
      {
        'target_name': 'libbar',
        'type': '<(library)',
        'defines': [
          'LOCAL_DEFINE_FOR_LIBBAR',
          'DEFINE_TO_USE_LIBBAR',
        ],
        'include_dirs': [
          '..',
          'include/libbar',
        ],
        'direct_dependent_settings': {
          'defines': [
            'DEFINE_TO_USE_LIBBAR',
          ],
          'include_dirs': [
            'include/libbar',
          ],
        },
      },
    ],
  }
```

如上所示，目标项foo依赖于同文件的libbar，而后者在`direct_dependent_settings` 里指定了`defines`和`include_dirs`项，所以foo的编译命令最终会加下面选项：

-D DEFINE_TO_USE_LIBBAR –I include/libbar

需要注意的是，`direct_dependent_settings`里的设置并不影响libbar本身的设置，比如libbar在自己的`defines`节点有`DEFINE_TO_USE_LIBBAR`宏定义，在自己的`include_dirs`里有

`include/libbar`目录。
**为Mac OS X bundles添加支持**:

```
{
      'target_name': 'test_app',
      'product_name': 'Test App Gyp',
      'type': 'executable',
      'mac_bundle': 1,
      'sources': [
        'main.m',
        'TestAppAppDelegate.h',
        'TestAppAppDelegate.m',
      ],
      'mac_bundle_resources': [
        'TestApp/English.lproj/InfoPlist.strings',
        'TestApp/English.lproj/MainMenu.xib',
      ],
      'link_settings': {
        'libraries': [
          '$(SDKROOT)/System/Library/Frameworks/Cocoa.framework',
        ],
      },
      'xcode_settings': {
        'INFOPLIST_FILE': 'TestApp/TestApp-Info.plist',
      },
    },
```


### 在重构中正确的移除文件

TODO（GYP官方文档还没提供此部分）
### 定制编译步骤

TODO（GYP官方文档还没提供此部分）
## 语言规格
### 目标

GYP工具的目标是为Chromium提供一个平台无关的配置系统，GYP工具通过gyp配置系统生成各平台的编译系统。包括Visual Studio, Xcode, SCons以及make files等。GYP的目标是让输入格式尽量通用，但不会浪费时间去做到完全正确。
### 背景

无论是Google内部还算外部，都有很多项目试图创建简单、通用的跨平台编译系统，允许每个平台有一定的灵活性以解决不同平台的差异。事实上，没有哪个系统明显能满足Chrojmium的需求，看上去这是一个棘手的问题。GYP的目标是成功创建一套工具，为Chromium高度定制，而不是做一套能处理所有跨平台编译问题的工具。

Mac系统下的编译环境最复杂。所以我们以Xcode为模型做初始配置，然后将配置适配到其他平台。（译注：聪明的做法!）.
### 预览

GYP系统拥有以下几个特点：
-	输入配置写在.gyp后后缀的文件里。
-	每个.gyp文件指定了如何编译文件里定义的「组件（Component）」
-	在Mac系统，一个.gyp文件生成一个Xcode .xcodeproj bundle，用以编译Target。
-	在Windows系统，一个.gyp文件生成了一个Visual Studio .sln文件，并为每个Target生成一个.vcproj工程文件。、
-	在Linux系统，一个.gyp文件生成一个SCons文件，为每个Target生成一个makefile文件。
-	.gyp文件的语法是Python数据结构（字典、列表、字符串）。
-	禁止任意使用Python代码
-	对全局和本地.gyp文件里使用eval()执行，从而.gyp里只允许有表达式而非任意Python语句。
-	所有的输入应该与Json格式兼容，除了两种情况：1、使用#写单行注释2、允许在列表和字典后添加逗号。
-	输入是一个key-value构成的字典。
-	无效的关键字是非法的，直接被忽略掉。



?
### 设计细节

GYP的设计原则：
-	一致的重用关键字
-	平台相关的关键字带上平台相关的概念前缀
-	例如：msvs-disabled_warnings, xcode_framework_dirs
-	输入语法是声明性的，数据驱动的
-	通过使用Python的eval()而非exec来强制这点，前者只解析表达式，后者可以解析任意Python语句。
-	关键字的具体语义是递归延迟的，直到所有的文件被读取，并且根据配置选择适当的文件。
-	源文件列表相关
-	源文件列表是平坦化的，源文件列表的排序遵循通常的开发者规则，当源文件被链接后，源文件列表会保持它们被列出的顺序。
-	源文件列表不包含逻辑文件夹机制，例如Visual Studio的Filter元素以及Xcode的Groups。
-	文件夹嵌套自动映射到文件夹系统。
-	占位符



例子：

```
{
  'target_defaults': {
    'defines': [
      'U_STATIC_IMPLEMENTATION',
      ['LOGFILE', 'foo.log',],
    ],
    'include_dirs': [
      '..',
    ],
  },
  'targets': [
    {
      'target_name': 'foo',
      'type': 'static_library',
      'sources': [
        'foo/src/foo.cc',
        'foo/src/foo_main.cc',
      ],
      'include_dirs': [
         'foo',
         'foo/include',
      ],
      'conditions': [
         [ 'OS==mac', { 'sources': [ 'platform_test_mac.mm' ] } ]
      ],
      'direct_dependent_settings': {
        'defines': [
          'UNIT_TEST',
        ],
        'include_dirs': [
          'foo',
          'foo/include',
        ],
      },
    },
  ],
}
```



#### 顶级元素

gyp字典的顶级元素如下：

节点关键字-说明

`conditions`-条件节可以包含任意其他顶级节点

`includes`-所包含的.gypi文件列表

`target_defaults`-默认目标项配置，被所有顶级targets继承

`targets`-target列表

`variables`-变量列表，变量定义是字符串键值对


**目标项（Targets）
**
节点关键字-说明

`actions`-自定义行为列表，每个行为接受输入文件，产生输出文件。

`all_dependent_settings`-为所有直接或间接依赖的目标项提供统一设置。注意与

`direct_dependent_settings`和`link_settings`区别开来。该节点可以包含几乎所有target的子元素，除了`configurations`,

`target_name`以及`type`元素。

`configurations`-编译配置列表

`copies`-拷贝动作列表

`defines`-预处理宏定义列表

`dependencies`-依赖target列表，在同一个文件夹内的直接写target名字；在其他.gyp文件里的，需要指明相对路径，例如：

../path/to/other.gyp:target_we_want

`direct_dependent_settings`-所有直接依赖该target的其他target所需的设置。该节点可以包含几乎所有的target子节点，除了`configurations`,

`target_name`以及`type`节点。注意与`all_dependent_settings`

和`link_settings`的区别。

`include_dirs`-头文件包含目录列表，被加到C/C++编译选项-I或/I

`libraries`-依赖库列表

`link_settings`-所有链接了该target的其他target所需的设置。`type`为

`executable`和`shared_library`的target是可链接的，所以如果它们依赖于某个不可链接的目标，比如`static_library`，它们将采用不可链接目标的`link_settings`。该节点可以包含几乎所有的target子节点，除了`configurations`,`target_name`,`type`等。注意和`all_dependent_settings`以及

`direct_dependent_settins`区别。

`rules`-一个特殊的自定义行为，接受输入文件列表，产生输出文件列表。

`sources`-源文件列表

`target_conditions`-类似`conditions`，但是直到他们被合并到实际target时才解析。

`target_name`-目标项名字

`type`-目标项类型，`executable`,`<(library)`,`static_library`或

`shared_library`，如果设置为`none`，则不会被链接。这在资源项目或文档项目时很有用。

`msvs_props`-Visual Studio 属性文件(.vsprops)列表

`xcode_config_file`-Xcode配置文件(.xcconfig)列表

`xcode_framework_dirs`-Xcode Framework 目录列表


#### 配置（Configurations）

`congiguration` 节点可以被`targets`或`target_defaults`节点包含。`configurations`节点是一个列表，因此当前无法覆盖在`target_defaults`里定义的`configurations`.



特别需要注意到是，应该为同一个工程的每个target定义同样的`configurations`集合。即使某个工程跨越了多个.gyp文件。



一个Configuration可以包含几乎所有的`target`子节点，除了

`actions`, `all_dependent_settins`, `configurations`, `dependencies`, `direct_dependent_settings`, 

`libraries`, `link_settings`, `sources`, `target_name`, `type`.



最后，每个Configuration必须包含`configuration_name`节点。
#### 条件（Conditionals）

Conditionals节点可以是.gyp文件的任何字典的直接子节点。有两种类型，

`conditions`节点和`target_conditions`节点，前者在加载.gyp文件时被处理，后者在所有依赖项合并完后处理。



`conditions`节点和`target_conditions`节点的值是一个列表，每个列表都包含三部分，第一部分是条件表达式；第二部分是条件表达所为true时会被融合的设置；第三部分是可选的，表示条件表达式为false时的设置。其中，第一部分和第二部分是必须的。

```
'conditions': [
          ['OS=="linux"', {
            'cflags!': [
              '-Werror',
            ],
          },{

  #可选的

	}],
]
```

-如左边示例，

黄色是条件表达式

灰色是条件表达式成立时的设置

蓝色是条件表达所不成立时的设置

蓝色部分是可选的。




#### 行为（Actions）

一个行为提供了对输入的一个或多个文件进行处理并产生输出文件的操作。`actions`节点的值是一个列表，列表的每个元素都是一个字典，这个字典的子元素如下：

节点关键字-说明-类型

`action_name`-行为名称，有些编译系统需要唯一的行为名称，有些则完全忽略行为名称。-string

`inputs`-输入文件列表-list

`outputs`-输出文件列表-list

`action`-行为命令行列表。为了最大化跨平台兼容，如果使用Python解释器，列表的第一个元素是`python`。-list

`message`-编译过程中需要显示的信息。-string



编译环境会比较输入和输出，如果任意输出文件丢失或比输入文件生成时间早，则会执行该行为。如果所有的输出文件都存在并且生成时间比输入文件早，则输出文件被认为是最新的，行为不会被执行。

Actions在Xcode里被实现为shell脚本，并且在编译前执行。Actions在Visual Studio里实现为编译文件的FileConfiguration元素的VCCustomBuildTool，并分别将

`action`指定为CommandLine，`outputs`指定为Outputs，`inputs`指定为`Additional Dependencies`.

下面是一个例子：


```
      'sources': [
        # libraries.cc is generated by the js2c action below.
        '<(INTERMEDIATE_DIR)/libraries.cc',
      ],
      'actions': [
        {
          'variables': {
            'core_library_files': [
              'src/runtime.js',
              'src/v8natives.js',
              'src/macros.py',
            ],
          },
          'action_name': 'js2c',
          'inputs': [
            'tools/js2c.py',
            '<@(core_library_files)', #gyp的变量名引用方式
          ],
          'outputs': [
            '<(INTERMEDIATE_DIR)/libraries.cc',
            '<(INTERMEDIATE_DIR)/libraries-empty.cc',
          ],
          'action': ['python', 'tools/js2c.py', '<@(_outputs)', 'CORE', '<@(core_library_files)'],
        },
      ],
```

#### 规则（Rules）

`rules`节点提供了自定义编译行为，对给定输入文件列表，产生输出文件列表。rules节点的值是一个列表，列表的每个元素是一个字典，字典的子节点如下：

节点关键字-说明-类型

`rule_name`-规则名字，有些编译系统要求唯一的规则名字，有些则完全忽略规则名字。-sting

`extension`-target的所有源文件如果标记了该关键字，则都会被当作rule的inputs-string

`inputs`-输入文件列表-list

`outputs`-输出文件列表-list

`action`-命令行列表。为了尽量兼容不同的平台，执行python脚本应该在列表的第一个元素指定`python`。-list

`message`-编译时显示的消息，可以使用RULE_INPUT_XXX变量。-string





下面的变量名可以被`output`, `action`, `message`使用。

变量名-说明

RULE_INPUT_PATH-当前文件的绝对路径

RULE_INPUT_DIRNAME-当前文件所在地目录

RULE_INPUT_NAME-当前文件的文件名

RULE_INPUT_ROOT-当前文件的不带扩展名的文件名

RULE_INPUT_EXT-当前文件的扩展名



Rules可以看做是Action生成器。所有标记了同一个extension的源文件都会被同一个Rules执行，并指定了同样的`inputs`, `outputs`, `action`, `message`，源文件被添加到`inputs`列表。
#### 拷贝（Copies）

`copies`节点提供了拷贝动作。`copies`节点是一个列表，列表的每个元素是一个字典，字典的子节点如下：

节点关键字-说明-类型

`destination`-目标拷贝路径-string

`files`-待拷贝文件列表-list


#### 生成XCode .pbxproj 文件
#### 生成Visual Studio .vcxproj文件
#### 生成Visual Studio .sln文件
## 输入格式参考

https://code.google.com/p/gyp/wiki/InputFormatReference



