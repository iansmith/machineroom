#!/bin/sh
echo starting to run the receive script
dir="/tmp/gitrecv_$$"
mkdir -p "$dir"
cat | tar -x -C "$dir"
echo "workingdir is $dir"

echo "----> Building alpha..."
cd "$dir/alpha"
make alpha

echo "----> Building beta..."
cd "$dir/beta"
make beta static/client.js

echo "---> Forcing down all docker containers"
cd "$dir"
docker ps -q | xargs docker kill

echo "---> Rebuilding docker images"
cd "$dir"
fig build && fig pull

echo "---> Restarting fig configuration"
cd "$dir"
fig up -d


echo "done"
exit 1
