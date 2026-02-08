#!/usr/bin/env bash
set -euo pipefail

OWNER="${OWNER:-odysseia-greek}"
GROUP="${GROUP:-dev}"
ROOT="${ROOT:-.}"

PATCH_SCRIPT="bump-patch.sh"
MINOR_SCRIPT="bump-minor.sh"

find_container_dirs() {
  (cd "$ROOT" && find . -maxdepth 2 -type f \( -name Containerfile -o -name Dockerfile \) -print) \
    | awk -F/ '{print $2}' \
    | sort -u
}

latest_tag() {
  local pkg="$1"

  gh api "/orgs/${OWNER}/packages/container/${pkg}/versions" \
    --paginate \
    | jq -r '
        .[]
        | .metadata.container.tags[]?
      ' \
    | grep -E '^v?[0-9]+\.[0-9]+\.[0-9]+$' \
    | head -n1
}

bump_patch() {
  local v="${1#v}"
  IFS='.' read -r maj min pat <<<"$v"
  echo "v${maj}.${min}.$((pat+1))"
}

bump_minor() {
  local v="${1#v}"
  IFS='.' read -r maj min _ <<<"$v"
  echo "v${maj}.$((min+1)).0"
}

# Initialise scripts
cat >"$PATCH_SCRIPT" <<EOF
#!/usr/bin/env bash
set -euo pipefail
# Auto-generated — delete lines you don't want
EOF

cat >"$MINOR_SCRIPT" <<EOF
#!/usr/bin/env bash
set -euo pipefail
# Auto-generated — delete lines you don't want
EOF

chmod +x "$PATCH_SCRIPT" "$MINOR_SCRIPT"

echo "# OWNER=$OWNER GROUP=$GROUP"
echo

while IFS= read -r dir; do
  [[ -z "$dir" ]] && continue

  latest="$(latest_tag "$dir" || true)"

  if [[ -z "$latest" ]]; then
    echo "# $dir: no semver tag found"
    echo
    continue
  fi

  patch="$(bump_patch "$latest")"
  minor="$(bump_minor "$latest")"

  # Pretty output
  echo "# $dir: latest=$latest"
  echo "# patch"
  echo "( cd \"$dir\" && archimedes images single -t \"$patch\" -m=false -g=\"$GROUP\" )"
  echo "# minor"
  echo "( cd \"$dir\" && archimedes images single -t \"$minor\" -m=false -g=\"$GROUP\" )"
  echo

  # Append to scripts
  echo "( cd \"$dir\" && archimedes images single -t \"$patch\" -m=false -g=\"$GROUP\" )" >>"$PATCH_SCRIPT"
  echo "( cd \"$dir\" && archimedes images single -t \"$minor\" -m=false -g=\"$GROUP\" )" >>"$MINOR_SCRIPT"

done < <(find_container_dirs)

echo "# Generated:"
echo "#  - $PATCH_SCRIPT"
echo "#  - $MINOR_SCRIPT"