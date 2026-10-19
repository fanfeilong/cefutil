#!/usr/bin/env bash
# Build examples/hello against CEF C API (cgo + xgoef_cef static lib).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VER="$(head -n1 "$ROOT/src/CEF_VERSION")"
PLAT="$(tail -n1 "$ROOT/src/CEF_VERSION")"
CEF_ROOT="${CEF_ROOT:-$ROOT/third_party/cef/cef_binary_${VER}_${PLAT}}"

if [[ ! -d "$CEF_ROOT/include" ]]; then
  echo "CEF_ROOT missing or incomplete: $CEF_ROOT" >&2
  echo "Run cef-download first (see README / skills/cef-download)." >&2
  exit 1
fi

OUT="$ROOT/src/xgoef/build/hello"
SHIM_BUILD="$ROOT/src/xgoef/build/shim"
mkdir -p "$OUT" "$SHIM_BUILD"

echo "==> configure xgoef_cef ($CEF_ROOT)"
cmake -S "$ROOT/src/xgoef/cef/shim" -B "$SHIM_BUILD" -G "Unix Makefiles" \
  -DCMAKE_BUILD_TYPE=Release -DUSE_SANDBOX=OFF \
  -DCMAKE_C_COMPILER=gcc -DCMAKE_CXX_COMPILER=g++ \
  -DCEF_ROOT="$CEF_ROOT"
cmake --build "$SHIM_BUILD" --target xgoef_cef -j"$(nproc)"

# Runtime files land in SHIM_BUILD/.../Release via CEF macros; locate archive.
LIB="$(find "$SHIM_BUILD" -name 'libxgoef_cef.a' | head -n1)"
RUNTIME_DIR="$(dirname "$LIB")"
if [[ -z "$LIB" || ! -f "$LIB" ]]; then
  echo "libxgoef_cef.a not found under $SHIM_BUILD" >&2
  exit 1
fi

echo "==> copy CEF runtime -> $OUT"
# Prefer files already copied next to the archive.
if [[ -f "$RUNTIME_DIR/libcef.so" ]]; then
  cp -a "$RUNTIME_DIR"/. "$OUT"/
else
  cp -a "$CEF_ROOT/Release"/. "$OUT"/
  cp -a "$CEF_ROOT/Resources"/. "$OUT"/
fi
cp -a "$ROOT/src/xgoef/examples/hello/web" "$OUT"/
cp -a "$ROOT/src/xgoef/examples/hello/xgoef.conf.json" "$OUT"/
mkdir -p "$OUT/api/js"
cp -a "$ROOT/src/xgoef/api/js/xgoef.js" "$OUT/api/js/"
# Fix html script path for the staged layout.
sed -i 's|src="../../../api/js/xgoef.js"|src="../api/js/xgoef.js"|' "$OUT/web/index.html"

# Native CEF helper for GPU/renderer (must not be the Go binary).
HELPER_SRC="$CEF_ROOT/build/tests/cefsimple_capi/Release/cefsimple_capi"
if [[ -x "$HELPER_SRC" ]]; then
  cp -f "$HELPER_SRC" "$OUT/xgoef-helper"
else
  echo "NOTE: $HELPER_SRC missing; build cefsimple_capi first for a stable helper." >&2
fi

echo "==> go build hello"
export CGO_ENABLED=1
export CC=gcc
export CXX=g++
# Keep CGO_CFLAGS minimal: they also apply to runtime/cgo. Package-specific
# flags live in cef/host_cgo.go; env only supplies CEF include/link paths.
export CGO_CFLAGS="-D_GNU_SOURCE -I$CEF_ROOT -I$ROOT/src/xgoef/cef/shim"
export CGO_LDFLAGS="-L$RUNTIME_DIR -L$OUT -lxgoef_cef -lcef -lX11 -lpthread -ldl -lrt -lstdc++ -Wl,-rpath,\$ORIGIN"

cd "$ROOT/src/xgoef"
go build -o "$OUT/xgoef-hello" ./examples/hello

echo "Built $OUT/xgoef-hello"
echo "Run:  cd $OUT && ./xgoef-hello"
