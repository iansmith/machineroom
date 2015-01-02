#!/bin/sh 

PROG="alpha"

## this horrible hack is necessary because of the fact that we are doing
## tricks with the gopath that confuse godeps.  it wants there to be a
## a version control system repository at the level of the project
## that you are running godeps at. 

cd /go/src/github.com/igneous-systems/$PROG/
git init .
git config --global user.name "Godep NeedsWork"
git config --global user.email "brokenhack@example.com"
git commit --allow-empty -m "no message"

cd /go/src/github.com/igneous-systems/lib/
git init .
git config --global user.name "Godep NeedsWork"
git config --global user.email "brokenhack@example.com"
git commit --allow-empty -m "no message"


cd /go/src/github.com/igneous-systems/$PROG/
godep save .

GODEPOK="y"
if [ "$?" != "0" ]; then
		echo "******* "
		echo "******* godep failed!"
		echo "******* "
		GODEPOK="n"
fi

#scary
rm -rf .git
cd /go/src/github.com/igneous-systems/lib/
rm -rf .git

##show us what happened
if [ "$GODEPOK" == "n" ]; then
	exit 1
fi

cat /go/src/github.com/igneous-systems/$PROG/Godeps/Godeps.json
