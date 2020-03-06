#!/usr/bin/env sh
set -Ceu

# goimports (openapiにより作成されたファイルはformat対象にしない)
test -z "$(goimports -l $(find ./app -type f -name '*.go' | grep -v /openapi-generated/) | grep -v /schemamodel/ | tee /dev/stderr)"

# go vet
go vet ./app/...

# golint (openapiにより作成されたファイルはformat対象にしない)
golint -set_exit_status $(go list ./... | grep -v /vendor/ | grep -v /openapi-generated/ | grep -v /schemamodel/)