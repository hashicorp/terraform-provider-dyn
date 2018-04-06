#!/usr/bin/env bash

artifacts_dir=${1}
project_dir=$PWD
project_name=$(basename ${project_dir})
git_branch=${CI_BRANCH:-$(git rev-parse --abbrev-ref HEAD)}
short_git_commit=$(git rev-parse --short HEAD)

# For every possible build combination (e.g lambda/service) loop over and build

    # find what type of service it is
      # makes sure the build are created in
    # /artifacts/{lambda,service}/{service-name} format
build_dest="${artifacts_dir}/${project_name}"
    # ensure the directory exists
mkdir -p $build_dest
    # build the binary

make vet fmtcheck
echo "Building ${project_name} in ${build_dest}"

go build -o ${build_dest}/${project_name}

echo "Compressing $build_dest/$git_branch-$short_git_commit.zip"
cd $build_dest
zip $build_dest/$git_branch-$short_git_commit.zip $project_name
