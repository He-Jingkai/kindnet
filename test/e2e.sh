#!/usr/bin/env bash
# Copyright 2018 The Kubernetes Authors.
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

# hack script for running a kind e2e
# must be run with a kubernetes checkout in $PWD (IE from the checkout)
# TODO(bentheelder): replace this with kubetest integration
# Usage: SKIP="ginkgo skip regex" FOCUS="ginkgo focus regex" kind-e2e.sh

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# run e2es with kubetest
run_tests() {
    # export the KUBECONFIG
    KUBECONFIG="$(kind get kubeconfig-path)"
    export KUBECONFIG

    # base kubetest args
    KUBETEST_ARGS="--provider=skeleton --test --check-version-skew=false"

    # get the number of worker nodes
    # TODO(bentheelder): this is kinda gross
    NUM_NODES="$(kubectl get nodes \
        -o=jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.taints}{"\n"}{end}' \
        | grep -cv "node-role.kubernetes.io/master" \
    )"

    # ginkgo regexes
    SKIP="${SKIP:-"Slow|Alpha|Kubectl|\\[(Disruptive|Feature:[^\\]]+|Flaky)\\]"}"
    FOCUS="${FOCUS:-"\\[Conformance\\]"}"
    # if we set PARALLEL=true, skip serial tests set --ginkgo-parallel
    PARALLEL="${PARALLEL:-false}"
    if [[ "${PARALLEL}" == "true" ]]; then
        SKIP="\\[Serial\\]|${SKIP}"
        KUBETEST_ARGS="${KUBETEST_ARGS} --ginkgo-parallel"
    fi

    # add ginkgo args
    KUBETEST_ARGS="${KUBETEST_ARGS} --test_args=\"--ginkgo.focus=${FOCUS} --ginkgo.skip=${SKIP} --report-dir=${ARTIFACTS} --disable-log-dump=true --num-nodes=${NUM_NODES}\""

    # setting this env prevents ginkg e2e from trying to run provider setup
    export KUBERNETES_CONFORMANCE_TEST="y"

    # run kubetest, if it fails clean up and exit failure
    eval "kubetest ${KUBETEST_ARGS}"
}

# setup kind, build kubernetes, create a cluster, run the e2es
main() {
    # ensure artifacts exists when not in CI
    ARTIFACTS="${ARTIFACTS:-${PWD}/_artifacts}"
    mkdir -p "${ARTIFACTS}"
    export ARTIFACTS
    # now build an run the cluster and tests
    run_tests
}

main
