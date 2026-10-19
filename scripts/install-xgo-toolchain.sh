#!/usr/bin/env bash
# Install XGo + LLGo + llcppg for cefutil / xgoef development (Ubuntu 24.04).
# Idempotent. Puts tools on /usr/local/bin and sets LLGO_ROOT.
set -euo pipefail

ROOT="${XGO_TOOLCHAIN_ROOT:-/opt/xgo-toolchain}"
GO_VER="${XGO_GO_VERSION:-1.27.2}"
CODENAME="$(lsb_release -cs)"

run_root() {
  if [[ "$(id -u)" -eq 0 ]]; then
    "$@"
  else
    sudo "$@"
  fi
}

echo "==> Go ${GO_VER}"
if ! /usr/local/go-1.27/bin/go version 2>/dev/null | grep -q "go${GO_VER}"; then
  tmp="$(mktemp -d)"
  curl -fsSL -o "$tmp/go.tgz" "https://go.dev/dl/go${GO_VER}.linux-amd64.tar.gz"
  run_root rm -rf /usr/local/go-1.27
  run_root mkdir -p /usr/local/go-1.27
  run_root tar -C /usr/local/go-1.27 --strip-components=1 -xzf "$tmp/go.tgz"
  run_root ln -sfn /usr/local/go-1.27/bin/go /usr/local/bin/go
  run_root ln -sfn /usr/local/go-1.27/bin/gofmt /usr/local/bin/gofmt
fi
export PATH="/usr/local/go-1.27/bin:/usr/local/bin:${PATH}"
go version

echo "==> LLVM 22 + native deps"
if ! llvm-config-22 --version 2>/dev/null | grep -q '^22\.'; then
  run_root mkdir -p /etc/apt/keyrings
  curl -fsSL https://apt.llvm.org/llvm-snapshot.gpg.key \
    | run_root gpg --dearmor -o /etc/apt/keyrings/llvm-snapshot.gpg
  echo "deb [signed-by=/etc/apt/keyrings/llvm-snapshot.gpg] http://apt.llvm.org/${CODENAME}/ llvm-toolchain-${CODENAME}-22 main" \
    | run_root tee /etc/apt/sources.list.d/llvm-22.list >/dev/null
  run_root apt-get update -qq
  run_root apt-get install -y -qq \
    llvm-22-dev clang-22 libclang-22-dev lld-22 lldb-22 libunwind-22-dev libc++-22-dev \
    pkg-config libgc-dev libssl-dev zlib1g-dev libffi-dev libcjson-dev libsqlite3-dev libuv1-dev \
    g++ g++-14 libstdc++-14-dev make git curl
fi
for t in llvm-config clang clang++ ld.lld lld lldb; do
  if [[ -x "/usr/lib/llvm-22/bin/$t" ]]; then
    run_root ln -sfn "/usr/lib/llvm-22/bin/$t" "/usr/local/bin/$t"
  elif [[ -x "/usr/bin/${t}-22" ]]; then
    run_root ln -sfn "/usr/bin/${t}-22" "/usr/local/bin/$t"
  fi
done
export PATH="/usr/local/bin:/usr/lib/llvm-22/bin:${PATH}"
export LIBRARY_PATH="/usr/lib/gcc/x86_64-linux-gnu/14:/usr/lib/x86_64-linux-gnu${LIBRARY_PATH:+:$LIBRARY_PATH}"

echo "==> LLGo (source)"
run_root mkdir -p "$ROOT"
run_root chown "$(id -u):$(id -g)" "$ROOT" 2>/dev/null || true
if [[ ! -d "$ROOT/llgo/.git" ]]; then
  git clone --depth 1 https://github.com/xgo-dev/llgo.git "$ROOT/llgo"
fi
(
  cd "$ROOT/llgo"
  ./install.sh
)
export LLGO_ROOT="$ROOT/llgo"
export GOPATH="${GOPATH:-$HOME/go}"
export PATH="${GOPATH}/bin:${PATH}"
run_root ln -sfn "${GOPATH}/bin/llgo" /usr/local/bin/llgo
llgo version

echo "==> XGo (source; apt package lags Go 1.27 export data)"
if [[ ! -d "$ROOT/xgo/.git" ]]; then
  git clone --depth 1 https://github.com/goplus/xgo.git "$ROOT/xgo"
fi
(
  cd "$ROOT/xgo"
  ./all.bash
)
run_root ln -sfn "$ROOT/xgo/bin/xgo" /usr/local/bin/xgo
xgo version

echo "==> llcppg (llgo install)"
if [[ ! -d "$ROOT/llcppg/.git" ]]; then
  git clone --depth 1 https://github.com/goplus/llcppg.git "$ROOT/llcppg"
fi
(
  cd "$ROOT/llcppg"
  go mod download
  llgo install ./cmd/llcppg
)
run_root ln -sfn "${GOPATH}/bin/llcppg" /usr/local/bin/llcppg

run_root tee /etc/profile.d/xgo-toolchain.sh >/dev/null <<EOF
export LLGO_ROOT=$ROOT/llgo
export GOPATH="\${GOPATH:-\$HOME/go}"
export PATH="/usr/local/bin:/usr/local/go-1.27/bin:/usr/lib/llvm-22/bin:\${GOPATH}/bin:\${PATH}"
export LIBRARY_PATH="/usr/lib/gcc/x86_64-linux-gnu/14:/usr/lib/x86_64-linux-gnu\${LIBRARY_PATH:+:\$LIBRARY_PATH}"
EOF

echo "==> OK"
echo "  go:     $(go version)"
echo "  xgo:    $(xgo version)"
echo "  llgo:   $(llgo version)"
echo "  llcppg: $(command -v llcppg)"
echo "  LLGO_ROOT=$LLGO_ROOT"
echo "  LLVM:   $(llgo env LLGO_LLVM_VERSION)"
