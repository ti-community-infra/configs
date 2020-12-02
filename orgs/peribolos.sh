#!/bin/bash

GO111MODULE=on go build -o tools/bin/merge cmd/merge/main.go
GO111MODULE=on go build -o tools/bin/peribolos k8s.io/test-infra/prow/cmd/peribolos

tools/bin/merge --org-part=tidb-community-bots=orgs/tidb-community-bots/org.yaml --merge-teams=true 2>&1 | tee config.yaml

tools/bin/peribolos "$@"
