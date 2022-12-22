#!/bin/bash

export PATH=$PATH:$(go env GOPATH)/bin

# Add the command to the CLI

# The first argument is the name of the command

cobra-cli add -a "Julien CAGNIART" -l "MIT" $1 $2 $3
