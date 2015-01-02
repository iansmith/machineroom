#!/bin/sh -x
GOPATH=/go:`godep path` go $@
