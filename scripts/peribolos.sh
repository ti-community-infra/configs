#!/bin/bash

echo "Build merge and peribolos tools..."
go build -o tools/bin/merge cmd/merge/main.go
# Add -mod=mod flag for go build because the default behavior of go build have changed in go 1.16.
# Refer: https://github.com/golang/go/issues/44212#issuecomment-776937327
go build -mod=mod -o tools/bin/peribolos k8s.io/test-infra/prow/cmd/peribolos

echo "Generate configuration file:"
tools/bin/merge --org-part=ti-community-infra=orgs/ti-community-infra/org.yaml --merge-teams=true 2>&1 | tee orgs/config.yaml

echo "Test the generated file"
go test ./...

echo "Start peribolos..."
tools/bin/peribolos "$@"
