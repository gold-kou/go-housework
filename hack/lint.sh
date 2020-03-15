#!/usr/bin/env sh
set -Ceu

# go vet
go vet ./app/...

# golint
golint -set_exit_status $(go list ./app/... | grep -v /schemamodel)
