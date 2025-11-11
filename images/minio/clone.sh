clone_tag() {
  local repo_url="$1"
  local tag="$2"
  local dest_dir="$3"
  git init "$dest_dir"
  cd "$dest_dir" || exit 1
  git remote add origin "$repo_url"
  git fetch --depth 1 origin tag "$tag"
  git checkout FETCH_HEAD
}