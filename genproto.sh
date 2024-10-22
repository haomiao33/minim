#!/bin/sh

protoc --go_out=. --go-grpc_out=. proto/RouterService.proto