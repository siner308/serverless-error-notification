#!/bin/bash

# build go package
GOOS=linux go build src/main.go

# build lambda package
zip main.zip main
