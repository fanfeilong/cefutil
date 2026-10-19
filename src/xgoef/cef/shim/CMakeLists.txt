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

set(CEFSIMPLE_CAPI_DIR "${CEF_ROOT}/tests/cefsimple_capi")

set(XGOEF_CEF_SRCS
  xgoef_host.c
  ${CEFSIMPLE_CAPI_DIR}/simple_app.c
  ${CEFSIMPLE_CAPI_DIR}/simple_browser_list.c
  ${CEFSIMPLE_CAPI_DIR}/simple_display_handler.c
  ${CEFSIMPLE_CAPI_DIR}/simple_handler.c
  ${CEFSIMPLE_CAPI_DIR}/simple_handler_linux.c
  ${CEFSIMPLE_CAPI_DIR}/simple_life_span_handler.c
  ${CEFSIMPLE_CAPI_DIR}/simple_load_handler.c
  ${CEFSIMPLE_CAPI_DIR}/simple_views.c
)

ADD_LOGICAL_TARGET("libcef_lib" "${CEF_LIB_DEBUG}" "${CEF_LIB_RELEASE}")
SET_CEF_TARGET_OUT_DIR()

list(REMOVE_ITEM CEF_COMPILER_FLAGS "-std=c99")
list(APPEND CEF_COMPILER_FLAGS "-std=c11")

add_library(xgoef_cef STATIC ${XGOEF_CEF_SRCS})
SET_LIBRARY_TARGET_PROPERTIES(xgoef_cef)
target_include_directories(xgoef_cef PRIVATE
  "${CMAKE_CURRENT_SOURCE_DIR}"
  "${CEF_ROOT}"
)
target_compile_definitions(xgoef_cef PRIVATE CEF_X11)
set_target_properties(xgoef_cef PROPERTIES
  C_STANDARD 11
  C_STANDARD_REQUIRED ON
  POSITION_INDEPENDENT_CODE ON
  ARCHIVE_OUTPUT_DIRECTORY "${CEF_TARGET_OUT_DIR}"
)
add_dependencies(xgoef_cef libcef_dll_wrapper)
target_link_libraries(xgoef_cef libcef_lib libcef_dll_wrapper ${CEF_STANDARD_LIBS})

COPY_FILES("xgoef_cef" "${CEF_BINARY_FILES}" "${CEF_BINARY_DIR}" "${CEF_TARGET_OUT_DIR}")
COPY_FILES("xgoef_cef" "${CEF_RESOURCE_FILES}" "${CEF_RESOURCE_DIR}" "${CEF_TARGET_OUT_DIR}")
