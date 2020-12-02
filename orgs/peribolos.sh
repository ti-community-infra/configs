#!/bin/bash

echo "Build merge and peribolos tools..."
go build -o tools/bin/merge cmd/merge/main.go
go build -o tools/bin/peribolos k8s.io/test-infra/prow/cmd/peribolos

echo "Generate configuration file:"
tools/bin/merge --org-part=tidb-community-bots=orgs/tidb-community-bots/org.yaml --merge-teams=true 2>&1 | tee config.yaml

echo " $(color-green 'done'), Start peribolos..."
for s in {5..1}; do
    echo -n $'\r'"in $s..."
    sleep 1s
done

tools/bin/peribolos "$@"
