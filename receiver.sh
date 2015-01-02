#!/bin/sh -x 
dir=`mktemp -d /tmp/gitrecvXXXXXX`
cat | tar -x -C "$dir"
source /etc/environment
echo "workdir is $dir"
cd "$dir/alpha"
make alpha
cd "$dir/beta"
make beta static/client.js

echo "done"
exit 1
