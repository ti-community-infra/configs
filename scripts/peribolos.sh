#!/bin/bash

echo "Build merge and peribolos tools..."
go build -o tools/bin/merge cmd/merge/main.go
go build -o tools/bin/peribolos k8s.io/test-infra/prow/cmd/peribolos

echo "Generate configuration file:"
tools/bin/merge --org-part=ti-community-infra=orgs/ti-community-infra/org.yaml --merge-teams=true 2>&1 | tee orgs/config.yaml

echo "Test the generated file"
go test ./...

echo "Start peribolos..."
tools/bin/peribolos "$@"
