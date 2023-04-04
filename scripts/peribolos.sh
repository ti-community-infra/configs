#!/bin/bash

echo "Generate configuration file:"
go run cmd/merge/main.go --org-part=ti-community-infra=orgs/ti-community-infra/org.yaml --merge-teams=true 2>&1 | tee orgs/config.yaml

echo "Test the generated file"
go test ./...

peribolos_bin=$PWD/tools/bin/peribolos
if [ -f $peribolos_bin ]; then
    echo 'tool existed.'
else
    pushd $(mktemp -d)
    git clone --depth 1 https://github.com/kubernetes/test-infra.git
    cd test-infra/prow/cmd/peribolos
    go build -o ${peribolos_bin}
    popd
fi
echo "Start peribolos..."
tools/bin/peribolos "$@"
