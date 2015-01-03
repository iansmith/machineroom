#!/bin/sh -x

##
## DOCKER KEYS
##
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9
sudo sh -c "echo deb https://get.docker.com/ubuntu docker main > /etc/apt/sources.list.d/docker.list"

##update everything
DEBIAN_FRONTEND=noninteractive sudo apt-get update -q
DEBIAN_FRONTEND=noninteractive sudo apt-get install -y -q dnsutils curl wget zip git mercurial build-essential ack-grep

##
## CONSUL
##

sudo mkdir -p /usr/local/bin/consul
cd /usr/local/bin/consul
sudo wget https://dl.bintray.com/mitchellh/consul/0.4.1_linux_amd64.zip
sudo unzip 0.4.1_linux_amd64.zip
sudo rm 0.4.1_linux_amd64.zip 
sudo chmod +x /usr/local/bin/consul/consul
sudo mkdir -p /etc/consul.d
sudo mkdir -p /var/consul
sudo cp /tmp/bootstrap.conf /etc/consul.d/bootstrap.conf
sudo cp /tmp/consul.conf /etc/init/consul.conf
sudo service consul start

##
## GIT RECEIVE
##
cd /tmp
wget https://raw.github.com/progrium/gitreceive/master/gitreceive
sudo mv gitreceive /usr/local/bin/gitreceive
sudo chmod +x /usr/local/bin/gitreceive
sudo /usr/local/bin/gitreceive init
sudo cp /tmp/receiver /home/git/receiver
chmod +x /home/git/receiver

##
## install docker and config it to use dns from consul
##
sudo apt-get install -y -q lxc-docker
echo 'DOCKER_OPTS="--dns 172.31.27.37 --dns 8.8.8.8 --dns 8.8.4.4"' | sudo tee -a /etc/default/docker
sudo service docker restart
sudo usermod -a -G docker git
sudo usermod -a -G docker ubuntu

##
## fig
##
sudo sh -c "curl -L https://github.com/docker/fig/releases/download/1.0.1/fig-`uname -s`-`uname -m` > /usr/local/bin/fig"
sudo chmod +x /usr/local/bin/fig

#
# FIX PATH
#
echo 'PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/bin/packer:/usr/local/bin/consul:/opt/tools/go/bin"' | sudo tee /etc/environment




