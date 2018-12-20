#!/bin/bash -e
# Run from directory above via ./scripts/cov.sh

rm -rf ./cov.out
go test ./... -coverprofile=cov.out
go tool cover -html=cov.out
