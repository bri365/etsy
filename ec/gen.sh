#!/usr/bin/env bash

PATH=$(pwd):$PATH protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative ec.proto
