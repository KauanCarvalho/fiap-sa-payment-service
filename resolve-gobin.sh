#!/bin/sh

set -e

GOBIN=${GOBIN:-$(go env GOBIN)}

if [ -z "$GOBIN" ]; then
  GOBIN="$(go env GOPATH)"/bin
fi

echo "$GOBIN"
