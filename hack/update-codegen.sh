#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# corresponding to go mod init <module>
MODULE=cslab.ece.ntua.gr/actimanager
# api package
APIS_PKG=api
# generated output package
OUTPUT_PKG=internal/pkg/generated

GROUP_VERSION=cslab.ece.ntua.gr:v1alpha1

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.

${CODEGEN_PKG}/generate-groups.sh all \
 ${MODULE}/${OUTPUT_PKG} ${MODULE}/${APIS_PKG} \
 ${GROUP_VERSION} \
  --output-base "$(dirname "${BASH_SOURCE[0]}")/../../.." \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.go.txt"
