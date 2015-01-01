START HERE
==========

Do this first because it will take a long time.  
```
git clone git@github.com:iansmith/machineroom
cd machineroom
vagrant up
```

This will pull a 700MB vagrant box called 
[machineroom-base](https://atlas.hashicorp.com/iansmith/boxes/machineroom-base) 
from  [atlas](http://atlas.hashicorp.com/).  Further, there are various
images that need to "pulled" into the docker cache and that takes a while
too.  

Read on while your downloads roll.

MACHINEROOM
============

The intent of the machineroom VM is to provide a simulation of the deploy 
environment for development.  In other words it provides the same 
"infrastructure" as our web app deployment environment but runs on a single
workstation.  

The quality of the "sim" is currently pretty poor.  Worse, the current sim
makes no effort at security of the internal containers.

The machineroom VM uses docker containers to simulate both _machines_ and
_processes_ in the real deployment environment.  Each docker container 
represents both a machine and a process running on that machine.  Thus, if
the real infrastructure has two machines, one for hot failover, with two
processes running on each one, the machineroom simulator must have four
containers.

The machineroom VMs "externally visible" interface should be a couple of
"front door" servers that expect to talk to users via a UI and mantle 
instances via an https port.  This should mimic the production environment,
but does not right now.

Because it's for development, the machineroom VM has a couple of "holes"
punched into it to allow easy development.  Primarily, it allows the
service bus to be "routed" to the developer's web browser through routing/DNS
chicanery with a name like "alpha.service.consul".  Further, there are the
usual tricks with vagrant VMs and the host operating system to allow you
to "live edit" files that are being served.


OBJECTIVES OF THIS DEMO
=======================

* Illustrate what a programming model that includes a service bus looks like.
* Show how we can simulate a production environment on a single dev workstation.

We would like to illustrate, but haven't yet:

* How you can deploy to a staging/production server from this same "configuration" that you use for development.
* Show how automated tools could be used as part of the deployment workflow.


DEMO
====


OS X SETUP OF ROUTES TO CONTAINERS
===================================

There is a line in the `Vagrantfile` that is an attempt to add a route 
that goes via the virtualbox vm to the containers.  It may not work on 
your system if you have your VM in a different place networkologically. 

The line that might be wrong is:
```
sudo route -n add 10.0.2.0/24 192.168.33.10
```

After doing this route operation correctly you should be able to do:
`ping 10.0.2.15` (assuming there is a container running there) and it should
respond with a standard response, not `Request timeout for icmp_seq`.

OS X SETUP CONSUL NAME RESOLUTION
===================================
You are going to want "easy" debugging via the web browser.  This can be done
by creating a new "resolver" for the consul-based DNS world.  

```
sudo mkdir /etc/resolver
sudo tee /etc/resolver/consul >/dev/null <<EOF
nameserver 192.168.33.10
port 8600
EOF
```

The docs on resolver indicate that you can add the port number (after a dot)
to the initial nameserver line, but Mac OSX's internal tools don't appear
to respect that configuration option so separate lines are recommended.  
Again, if your virtualbox VM isn't on the same IP address as mine, you may 
have to configure that line.

You can test it with your web browser (only) to a consul service 
`http://app.service.consul`.  You can look at what the resolver is sending
with `dig @192.168.33.10 -p 8600 app.service.consul` but note that the 
operating system's internal tools do not go through the same path as `dig`
so that is only helpful to see what the consul service is responding with. 

Sadly, OSX caches pretty aggressively on DNS servers so you may want
to induce a refresh like this:

```
sudo launchctl unload -w /System/Library/LaunchDaemons/com.apple.discoveryd.plist
sudo launchctl load -w /System/Library/LaunchDaemons/com.apple.discoveryd.plist
```

Apple's docs indicate that `sudo discoveryutil udnsflushcaches` should work,
but it doesn't for me.  Conversation with APilloud indicated that the problem
may be that we are manipulating the set of resolvers, not the content of the
resolvers ("the unicast dns cache") and this might be why the brute-force 
approach is needed.

Later experiments revealed that the command above to flush the udns (unicast
dns) cache _is_ important if you have previously goofed up the config for your
`.consul` routes, and end up with a cached value that originated from your
"normal" nameserver.  Look at the last section of the dig report where it 
should say ``;;SERVER: 192.168.33.10#8600(192.168.33.10)``.  
Clearing your udns cache is thus still advised.


FIG
===

[Fig](http://www.fig.sh/index.html) is a tool for orchestrating a 
bunch of docker containers from a single
config file. It does nothing more than docker commands, it's just a convenient
way of doing the same commands to a lot of containers at once.  The config
file that all fig commands work from is [fig.yml](http://www.fig.sh/yml.html)
and is just a thin wrapper over docker's command line parameters.

Note: yaml is the  work of satan.

`fig up` brings up the entire network of containers (7 in our case here).  I have noticed that control-c doesn't _always_ seem to kill all the nodes in the  network.  To reliably bring everything down, use `fig kill`.  
You can see the running nodes in the network with `fig ps` from a 
different shell. 

`fig pull` pulls all the containers for parts of the network that are using
a "fixed" container. You shouldn't have to do this as it is done once when
the vagrant box starts up.  `fig build` rebuilds all the containers that are
built locally.  Currently, this is needed after *any* go code change 
because we don't have a `makefile` that under stands how to do this
automatically. This operation is generally cheap due to docker caching so it
does not hurt to do: `fig kill; fig build; fig up`.


BETA
====


TODO
=====
* The DNS configuration is still bogus in that hosts like `alpha.service.consul` seems to resolve to the same "host" on the bridge as `beta.service.consul`.

* Need a way to signal to the registrator that we want a port exposed other than publishing it to the host (which leads to complications because of port
collisions, see previous entry).

* How do we bootstrap the initial consul server? Beta and registrator do this
through a docker link currently.

* Nginx understands graceful restart, can we version the newly pushed code and
do a rolling update?

