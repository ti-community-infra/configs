#!/bin/bash

echo "Build merge and peribolos tools..."
go build -o tools/bin/merge cmd/merge/main.go
go build -o tools/bin/peribolos k8s.io/test-infra/prow/cmd/peribolos

echo "Generate configuration file:"
tools/bin/merge --org-part=tidb-community-bots=orgs/tidb-community-bots/org.yaml --merge-teams=true 2>&1 | tee config.yaml

echo "Start peribolos..."
tools/bin/peribolos "$@"
