#!/bin/sh
# Pack agent skills into a release zip.
# Output: dist/cef-agent-skills-<version>.zip
set -eu
cd "$(dirname "$0")/.."

version=$(tr -d '[:space:]' < skills/VERSION)
out_dir=dist
stage=$(mktemp -d)
trap 'rm -rf "$stage"' EXIT

mkdir -p "$out_dir" "$stage/skills"
cp skills/README.md skills/VERSION "$stage/skills/"
for name in cef-integrate cef-lifecycle cef-js-bridge cef-host-surface cef-upgrade; do
  mkdir -p "$stage/skills/$name"
  cp "skills/$name/SKILL.md" "$stage/skills/$name/SKILL.md"
done

out="$out_dir/cef-agent-skills-$version.zip"
rm -f "$out"
(
  cd "$stage"
  zip -r -q "$OLDPWD/$out" skills
)
echo "$out"
