#!/bin/sh -x
GOPATH=/go:`godep path` gopherjs $@
