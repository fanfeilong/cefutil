set(CEFDEMO_SRCS
  app.cc
  app.h
  handler.cc
  handler.h
  main_linux.cc
  )

set(CEF_TARGET "cefdemo")
ADD_LOGICAL_TARGET("libcef_lib" "${CEF_LIB_DEBUG}" "${CEF_LIB_RELEASE}")
SET_CEF_TARGET_OUT_DIR()

add_executable(${CEF_TARGET} ${CEFDEMO_SRCS})
SET_EXECUTABLE_TARGET_PROPERTIES(${CEF_TARGET})
add_dependencies(${CEF_TARGET} libcef_dll_wrapper)
target_link_libraries(${CEF_TARGET} libcef_lib libcef_dll_wrapper ${CEF_STANDARD_LIBS})
target_compile_definitions(${CEF_TARGET} PRIVATE
  "CEFDEMO_WEB_DIR=\"${CEF_TARGET_OUT_DIR}/web\""
  "CEFDEMO_ROOT_CACHE=\"${CEF_TARGET_OUT_DIR}/root_cache\"")

set_target_properties(${CEF_TARGET} PROPERTIES
  INSTALL_RPATH "$ORIGIN"
  BUILD_WITH_INSTALL_RPATH TRUE
  RUNTIME_OUTPUT_DIRECTORY "${CEF_TARGET_OUT_DIR}")

COPY_FILES("${CEF_TARGET}" "${CEF_BINARY_FILES}" "${CEF_BINARY_DIR}" "${CEF_TARGET_OUT_DIR}")
COPY_FILES("${CEF_TARGET}" "${CEF_RESOURCE_FILES}" "${CEF_RESOURCE_DIR}" "${CEF_TARGET_OUT_DIR}")

add_custom_command(
  TARGET ${CEF_TARGET}
  POST_BUILD
  COMMAND ${CMAKE_COMMAND} -E copy_directory
          "${CMAKE_CURRENT_SOURCE_DIR}/web"
          "${CEF_TARGET_OUT_DIR}/web"
  VERBATIM)
