#!/bin/sh -x

##
## DOCKER KEYS
##
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9
sudo sh -c "echo deb https://get.docker.com/ubuntu docker main > /etc/apt/sources.list.d/docker.list"

##update everything
DEBIAN_FRONTEND=noninteractive sudo apt-get update -q
DEBIAN_FRONTEND=noninteractive sudo apt-get install -y -q dnsutils curl wget zip git mercurial build-essential ack-grep

 
# Setup sudo to allow no-password sudo for "admin"
# sudo groupadd -r admin
# sudo cp /etc/sudoers /etc/sudoers.orig
# sudo sed -i -e '/Defaults\s\+env_reset/a Defaults\texempt_group=admin' /etc/sudoers
# sudo sed -i -e 's/%admin ALL=(ALL) ALL/%admin ALL=NOPASSWD:ALL/g' /etc/sudoers
 
##
## PACKER
##

#sudo mkdir -p /usr/local/bin/packer
#cd /usr/local/bin/packer
#sudo wget https://dl.bintray.com/mitchellh/packer/packer_0.7.5_linux_amd64.zip
#sudo unzip packer_0.7.5_linux_amd64.zip
#sudo rm packer_0.7.5_linux_amd64.zip
#sudo chmod +x packer*

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


##
## install docker and config it to use dns from consul
##
sudo apt-get install -y -q lxc-docker
echo 'DOCKER_OPTS="--dns 10.0.2.15 --dns 8.8.8.8 --dns 8.8.4.4"' | sudo tee -a /etc/default/docker
sudo service docker restart
sudo usermod -a -G docker git
sudo usermod -a -G docker ubuntu
sudo usermod -a -G docker vagrant

##
## fig
##
sudo sh -c "curl -L https://github.com/docker/fig/releases/download/1.0.1/fig-`uname -s`-`uname -m` > /usr/local/bin/fig"
sudo chmod +x /usr/local/bin/fig

#
# FIX PATH
#
echo 'PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/bin/packer:/usr/local/bin/consul:/opt/tools/go/bin"' | sudo tee /etc/environment

#
# FIX SSH FOR OUR USERS
#
sudo groupadd -r igneous
echo "user igneous, password $1"
sudo useradd igneous -p "$1" -m -g docker -G igneous -s /bin/bash
echo 'Match user igneous' | sudo tee -a /etc/ssh/sshd_config
echo '     PasswordAuthentication yes' | sudo tee -a /etc/ssh/sshd_config

##
## pull some useful docker images, root is not in docker group but vagrant is
## do this last, because it mucks with current group
##
#newgrp docker
#echo "docker image: progrium/registrator"
#docker pull progrium/registrator
#echo "docker image: phusion/baseimage"
#docker pull phusion/baseimage:0.9.15


