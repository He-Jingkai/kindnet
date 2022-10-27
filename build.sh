#!/usr/bin/env bash
# Copyright 2019 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o nounset
set -o errexit
set -o pipefail

# cd to the repo root
REPO_ROOT=$(git rev-parse --show-toplevel)
cd "${REPO_ROOT}"

# build the binary
export SOURCE_DIR="${REPO_ROOT}"

# TODO: verisoning
# build image
IMAGE="hejingkai/kindnetd"
TAG="arm64"
docker build \
  -t "${IMAGE}:${TAG}" \
  --build-arg="GOARCH=arm64" \
  -f Dockerfile \
  "${SOURCE_DIR}"

IMAGE="hejingkai/kindnetd"
TAG="amd64"
docker build \
  -t "${IMAGE}:${TAG}" \
  --build-arg="GOARCH=amd64" \
  -f Dockerfile \
  "${SOURCE_DIR}"

docker push hejingkai/kindnetd:amd64
docker push hejingkai/kindnetd:arm64