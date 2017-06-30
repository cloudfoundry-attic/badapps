#!/bin/bash
set -e

GOPATH=$PWD

cd src/code.cloudfoundry.org/badapps
go get -t ./...

ginkgo -r -p -race $@

