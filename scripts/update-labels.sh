#!/usr/bin/env bash

echo "Build label_sync..." >&2

label_sync=$PWD/tools/bin/label_sync
pushd $(mktemp -d)
  git clone https://github.com/kubernetes/test-infra.git
  cd test-infra/label_sync
  go build -o ${label_sync}
popd

# Generate labels.
pushd prow/config || exit
  "$label_sync" \
    --config=labels.yaml \
    --action=docs \
    --docs-template=labels.md.tmpl \
    --docs-output=labels.md
popd