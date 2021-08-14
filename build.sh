#!/bin/bash

# build go package
go build src/main.go

# build lambda package
zip function.zip main