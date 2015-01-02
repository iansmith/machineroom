#!/bin/sh -x

## this horrible hack is necessary because of the fact that we are doing
## tricks with the gopath that confuse godeps.  it wants there to be a
## a version control system repository at the level of the project
## that you are running godeps at. 

cd /go/src/github.com/igneous-systems/beta/
git init .
git config --global user.name "Godep NeedsWork"
git config --global user.email "brokenhack@example.com"
git commit --allow-empty -m "no message"

cd /go/src/github.com/igneous-systems/beta/
godep save ./...

#scary
rm -rf .git

##show us what happened
cat /go/src/github.com/igneous-systems/beta/Godeps/Godeps.json
