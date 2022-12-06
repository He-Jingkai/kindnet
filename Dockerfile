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

# STEP 1: Build kindnetd binary
FROM golang:1.18 AS builder
# golang envs
# ENV GO111MODULE="on"
ENV GOPROXY=https://proxy.golang.org
# copy in sources
WORKDIR /src
COPY . .
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o /go/bin/kindnetd        ./cmd/kindnetd
RUN go build -o /opt/cni/bin/host-local ./cmd/host-local
RUN go build -o /opt/cni/bin/ptp        ./cmd/ptp
RUN go build -o /opt/cni/bin/bridge     ./cmd/bridge
RUN go build -o /opt/cni/bin/portmap    ./cmd/portmap
RUN go build -o /opt/cni/bin/loopback   ./cmd/loopback

# STEP 2: Build small image
FROM gcr.io/istio-release/base:latest
COPY --from=builder --chown=root:root /go/bin/kindnetd /bin/kindnetd
COPY --from=builder --chown=root:root /opt/cni/bin     /opt/cni/bin
CMD ["/bin/kindnetd"]
