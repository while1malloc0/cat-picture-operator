#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

vendor/k8s.io/code-generator/generate-groups.sh \
deepcopy \
github.com/while1malloc0/cat-picture-operator/operator/pkg/generated \
github.com/while1malloc0/cat-picture-operator/operator/pkg/apis \
cat:v1 \
--go-header-file "./tmp/codegen/boilerplate.go.txt"
