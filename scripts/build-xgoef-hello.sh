#!/usr/bin/env bash
# Build examples/hello with LLGo (no cgo) against the CEF C API shim.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VER="$(head -n1 "$ROOT/src/CEF_VERSION")"
PLAT="$(tail -n1 "$ROOT/src/CEF_VERSION")"
CEF_ROOT="${CEF_ROOT:-$ROOT/third_party/cef/cef_binary_${VER}_${PLAT}}"

export PATH="/usr/local/bin:/usr/local/go-1.27/bin:/usr/lib/llvm-22/bin:${HOME}/go/bin:${PATH}"
export LLGO_ROOT="${LLGO_ROOT:-/opt/xgo-toolchain/llgo}"
export LIBRARY_PATH="/usr/lib/gcc/x86_64-linux-gnu/14:/usr/lib/x86_64-linux-gnu${LIBRARY_PATH:+:$LIBRARY_PATH}"
unset CGO_CFLAGS CGO_LDFLAGS CGO_CXXFLAGS CGO_CPPFLAGS CGO_FFLAGS
export CGO_ENABLED=0

if ! command -v llgo >/dev/null; then
  echo "llgo not found; run ./scripts/install-xgo-toolchain.sh first" >&2
  exit 1
fi
if [[ ! -d "$CEF_ROOT/include" ]]; then
  echo "CEF_ROOT missing or incomplete: $CEF_ROOT" >&2
  exit 1
fi
if [[ -z "${LLGO_ROOT:-}" || ! -d "$LLGO_ROOT/runtime" ]]; then
  echo "LLGO_ROOT invalid: ${LLGO_ROOT:-<empty>}" >&2
  exit 1
fi

OUT="$ROOT/src/xgoef/build/hello"
SHIM_BUILD="$ROOT/src/xgoef/build/shim"
PC_DIR="$ROOT/src/xgoef/build/pkgconfig"
mkdir -p "$OUT" "$SHIM_BUILD" "$PC_DIR"

echo "==> configure xgoef_cef ($CEF_ROOT)"
cmake -S "$ROOT/src/xgoef/cef/shim" -B "$SHIM_BUILD" -G "Unix Makefiles" \
  -DCMAKE_BUILD_TYPE=Release -DUSE_SANDBOX=OFF \
  -DCMAKE_C_COMPILER=gcc -DCMAKE_CXX_COMPILER=g++ \
  -DCEF_ROOT="$CEF_ROOT"
cmake --build "$SHIM_BUILD" --target xgoef_cef -j"$(nproc)"

LIB="$(find "$SHIM_BUILD" -name 'libxgoef_cef.a' | head -n1)"
RUNTIME_DIR="$(dirname "$LIB")"
if [[ -z "$LIB" || ! -f "$LIB" ]]; then
  echo "libxgoef_cef.a not found under $SHIM_BUILD" >&2
  exit 1
fi

echo "==> write pkg-config xgoef-cef.pc"
# Do NOT put -Wl,-rpath,$ORIGIN (or any $VAR) in Libs: LLGo runs os.Expand on
# pkg-config output, so $ORIGIN becomes empty and "-Wl,-rpath," steals the
# following -lffi into RUNPATH (undefined ffi_* at link time).
# LLGo already maps every -L to a runtime -rpath; keep Libs comma-free.
cat > "$PC_DIR/xgoef-cef.pc" <<EOF
Name: xgoef-cef
Description: xgoef CEF C API host shim
Version: 0.0.1
Libs: -L${RUNTIME_DIR} -lxgoef_cef -L${OUT} -lcef -lX11 -lpthread -ldl -lrt -lstdc++
Cflags: -I${ROOT}/src/xgoef/cef/shim -I${CEF_ROOT}
EOF
export PKG_CONFIG_PATH="$PC_DIR:${PKG_CONFIG_PATH:-}"
pkg-config --libs xgoef-cef >/dev/null

echo "==> copy CEF runtime -> $OUT"
if [[ -f "$RUNTIME_DIR/libcef.so" ]]; then
  cp -a "$RUNTIME_DIR"/. "$OUT"/
else
  cp -a "$CEF_ROOT/Release"/. "$OUT"/
  cp -a "$CEF_ROOT/Resources"/. "$OUT"/
fi
# Keep the static archive available for the linker next to OUT as well.
cp -f "$LIB" "$OUT"/
cp -a "$ROOT/src/xgoef/examples/hello/web" "$OUT"/
cp -a "$ROOT/src/xgoef/examples/hello/xgoef.conf.json" "$OUT"/
mkdir -p "$OUT/api/js"
cp -a "$ROOT/src/xgoef/api/js/xgoef.js" "$OUT/api/js/"
sed -i 's|src="../../../api/js/xgoef.js"|src="../api/js/xgoef.js"|' "$OUT/web/index.html"

HELPER_SRC="$CEF_ROOT/build/tests/cefsimple_capi/Release/cefsimple_capi"
if [[ -x "$HELPER_SRC" ]]; then
  cp -f "$HELPER_SRC" "$OUT/xgoef-helper"
else
  echo "NOTE: $HELPER_SRC missing; build cefsimple_capi first for a stable helper." >&2
fi

echo "==> llgo build hello (CGO_ENABLED=0)"
cd "$ROOT/src/xgoef"
go get github.com/goplus/lib@v0.6.3
llgo build -o "$OUT/xgoef-hello" ./examples/hello

echo "Built $OUT/xgoef-hello"
echo "Run:  cd $OUT && ./xgoef-hello"
echo "Note: LLGo host — no cgo"
