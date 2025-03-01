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

docker buildx create --use --name kindnet-builder
docker buildx inspect kindnet-builder --bootstrap

#IMAGE="hejingkai/kindnetd"
IMAGE="registry.cn-hangzhou.aliyuncs.com/jkhe/kindnetd"
TAG="1.0.109"
docker buildx build \
  -t "${IMAGE}:${TAG}" \
  --platform=linux/arm64,linux/amd64 \
  -f Dockerfile \
  "${SOURCE_DIR}" \
  --push

docker buildx rm kindnet-builder