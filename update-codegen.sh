#!/bin/bash

ROOT_PACKAGE="github.com/dichque/grafana-operator"
CUSTOM_RESOURCE_NAME="grafana"
CUSTOM_RESOURCE_VERSION="v1"
GO111MODULE=off

# go get -u k8s.io/code-generator/...
cd $GOPATH/src/k8s.io/code-generator

./generate-groups.sh all "$ROOT_PACKAGE/pkg/client" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION" --output-base "${GOPATH}/src" --go-header-file "hack/boilerplate.go.txt"
