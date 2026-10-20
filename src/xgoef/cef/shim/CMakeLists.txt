cmake_minimum_required(VERSION 3.21)
project(xgoef_cef C CXX)

if(NOT CEF_ROOT)
  message(FATAL_ERROR "Pass -DCEF_ROOT=/absolute/path/to/cef_binary_*")
endif()

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${CEF_ROOT}/cmake")
find_package(CEF REQUIRED)

# Prefer sharing an already-built wrapper from the CEF tree when present.
set(CEF_WRAPPER_BUILD "${CEF_ROOT}/build/libcef_dll_wrapper")
if(EXISTS "${CEF_WRAPPER_BUILD}/libcef_dll_wrapper.a")
  add_library(libcef_dll_wrapper STATIC IMPORTED)
  set_target_properties(libcef_dll_wrapper PROPERTIES
    IMPORTED_LOCATION "${CEF_WRAPPER_BUILD}/libcef_dll_wrapper.a")
else()
  add_subdirectory(${CEF_LIBCEF_DLL_WRAPPER_PATH} libcef_dll_wrapper)
endif()

set(XGOEF_CEF_SRCS
  xgoef_host.cc
  xgoef_app.cc
  xgoef_handler.cc
  xgoef_handler_linux.cc
  xgoef_renderer.cc
  xgoef_ipc.cc
)

ADD_LOGICAL_TARGET("libcef_lib" "${CEF_LIB_DEBUG}" "${CEF_LIB_RELEASE}")
SET_CEF_TARGET_OUT_DIR()

add_library(xgoef_cef STATIC ${XGOEF_CEF_SRCS})
SET_LIBRARY_TARGET_PROPERTIES(xgoef_cef)
target_include_directories(xgoef_cef PRIVATE
  "${CMAKE_CURRENT_SOURCE_DIR}"
  "${CEF_ROOT}"
)
target_compile_definitions(xgoef_cef PRIVATE CEF_X11)
set_target_properties(xgoef_cef PROPERTIES
  CXX_STANDARD 17
  CXX_STANDARD_REQUIRED ON
  POSITION_INDEPENDENT_CODE ON
  ARCHIVE_OUTPUT_DIRECTORY "${CEF_TARGET_OUT_DIR}"
)
add_dependencies(xgoef_cef libcef_dll_wrapper)
target_link_libraries(xgoef_cef libcef_lib libcef_dll_wrapper ${CEF_STANDARD_LIBS})

# Native helper with renderer-side MessageRouter (not the LLGo host binary).
add_executable(xgoef_helper xgoef_helper_main.cc)
SET_EXECUTABLE_TARGET_PROPERTIES(xgoef_helper)
target_include_directories(xgoef_helper PRIVATE
  "${CMAKE_CURRENT_SOURCE_DIR}"
  "${CEF_ROOT}"
)
target_compile_definitions(xgoef_helper PRIVATE CEF_X11)
set_target_properties(xgoef_helper PROPERTIES
  CXX_STANDARD 17
  CXX_STANDARD_REQUIRED ON
  OUTPUT_NAME "xgoef-helper"
  RUNTIME_OUTPUT_DIRECTORY "${CEF_TARGET_OUT_DIR}"
)
add_dependencies(xgoef_helper libcef_dll_wrapper xgoef_cef)
target_link_libraries(xgoef_helper xgoef_cef libcef_lib libcef_dll_wrapper ${CEF_STANDARD_LIBS})

COPY_FILES("xgoef_cef" "${CEF_BINARY_FILES}" "${CEF_BINARY_DIR}" "${CEF_TARGET_OUT_DIR}")
COPY_FILES("xgoef_cef" "${CEF_RESOURCE_FILES}" "${CEF_RESOURCE_DIR}" "${CEF_TARGET_OUT_DIR}")
COPY_FILES("xgoef_helper" "${CEF_BINARY_FILES}" "${CEF_BINARY_DIR}" "${CEF_TARGET_OUT_DIR}")
COPY_FILES("xgoef_helper" "${CEF_RESOURCE_FILES}" "${CEF_RESOURCE_DIR}" "${CEF_TARGET_OUT_DIR}")
