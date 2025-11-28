#!/usr/bin/env just --justfile

update:
  go get -u
  go mod tidy -v

lint:
  golangci-lint run