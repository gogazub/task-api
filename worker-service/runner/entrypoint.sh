#!/usr/bin/env bash
set -euo pipefail

SRC="${SRC:-/work/main.cpp}"
OUT="${OUT:-/work/a.out}"
STD="${STD:-c++20}"
CXX="${CXX:-g++}"
OPT="${OPT:--O2}"

mkdir -p "$(dirname "$SRC")"
mkdir -p "$(dirname "$OUT")"

if [[ ! -s "$SRC" ]]; then
  echo ">> No $SRC found, reading source from STDIN..."
  cat > "$SRC"
fi

echo ">> Compiling: $CXX -std=$STD $OPT $SRC -o $OUT"
$CXX -std="$STD" $OPT "$SRC" -o "$OUT"

echo ">> Running: $OUT"
exec "$OUT" "$@"
