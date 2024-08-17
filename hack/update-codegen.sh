#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DOMAIN=cslab.ece.ntua.gr
MODULE=${DOMAIN}/actimanager
APIS_PKG=api
OUTPUT_PKG=internal/pkg/generated

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

source "${CODEGEN_PKG}/kube_codegen.sh"

kube::codegen::gen_helpers \
    --input-pkg-root "${MODULE}/${APIS_PKG}" \
    --output-base "$(dirname "${BASH_SOURCE[0]}")/../../.." \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.go.txt"

kube::codegen::gen_client \
    --with-watch \
    --with-applyconfig \
    --input-pkg-root "${MODULE}/${APIS_PKG}" \
    --output-pkg-root "${MODULE}/${OUTPUT_PKG}" \
    --output-base "$(dirname "${BASH_SOURCE[0]}")/../../.." \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.go.txt"