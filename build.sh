#!/bin/bash
#
# usage: ./build.bash /path/to/target.go

binary="DiscordScrobbler"

# argument handling and check
test "$1" && target="$1" # .go file to build
if ! test "$target"
then
  echo "target file required"
  exit 1
fi

# git check for uncommited changes.
if ! git diff-index --quiet HEAD --; then
    echo "Uncommited changes! Please commit before build."
    exit 1
fi

# git compress and then go generate.
git gc > /dev/null
trap "git reset HEAD --hard" INT
find . -name '*.go' -not -path "*/vendor/*" -not -path "*/pkg/*" -exec go generate {} \;
echo -e "\\nGo generate complete.\\n\\n"

# start build
if [[ ! -v binary ]]; then
    binary="$(basename "$(pwd)")" # default to default
    test "$2" && binary="$2" # binary output
fi

# Platforms to compile for
platforms="linux/386 linux/amd64 linux/arm windows/386 windows/amd64 darwin/386 darwin/amd64 freebsd/386 freebsd/amd64 freebsd/arm"

if ! test "$platforms"; then
  echo "no valid os/arch pairs were found to build"
  echo "- see: https://gist.github.com/jmervine/7d3f455e923cf2ac3c9e#file-golang-crosscompile-setup-bash"
  exit 1
fi

for platform in ${platforms}
do

    echo "$platform"
    IFS='/' read -ra splits <<< "$platform"

    goos="${splits[0]}"
    goarch="${splits[1]}"

    # ensure output file name
    output="$binary"
    test "$output" || output="$(basename "$target" | sed 's/\.go//')"

    # add exe to windows output
    [[ "windows" == "$goos" ]] && output="$output.exe"

    # set destination path for binary
    destination="./builds/$goos/$goarch/$output"

    echo "GOOS=$goos GOARCH=$goarch go build -x -o $destination $target"
    GOOS=$goos GOARCH=$goarch go build -o "$destination" "$target"
    echo -e "=-=-=-=-=-=-=-=-=-=-=-=-=-\\n\\n\\n\\n"

done

# Revert back to default.
git reset HEAD --hard
