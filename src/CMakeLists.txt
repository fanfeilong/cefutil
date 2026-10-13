cmake_minimum_required(VERSION 3.21)
project(cefutil_src)

# CEF_ROOT is the extracted cef_binary_* directory from cef-download.
if(NOT CEF_ROOT)
  message(FATAL_ERROR "Pass -DCEF_ROOT=/absolute/path/to/cef_binary_<version>_<platform>")
endif()

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${CEF_ROOT}/cmake")
find_package(CEF REQUIRED)

add_subdirectory(${CEF_LIBCEF_DLL_WRAPPER_PATH} libcef_dll_wrapper)
add_subdirectory(cefdemo)
