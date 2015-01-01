# $script = <<SCRIPT
# export DEBIAN_FRONTEND=noninteractive
#     #docker
#     sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 36A1D7869245C8950F966E92D8576A8BA88D21E9
#     sudo sh -c "echo deb https://get.docker.com/ubuntu docker main > /etc/apt/sources.list.d/docker.list"

#     #update everything
#     sudo apt-get update
#     sudo apt-get install -y -q dnsutils curl wget zip git mercurial build-essential

#     sudo mkdir -p /usr/local/bin/packer
#     cd /usr/local/bin/packer
#     sudo wget https://dl.bintray.com/mitchellh/packer/packer_0.7.5_linux_amd64.zip
#     sudo unzip packer_0.7.5_linux_amd64.zip
#     sudo rm packer_0.7.5_linux_amd64.zip
#     sudo chmod +x packer*

#     sudo mkdir -p /usr/local/bin/consul
#     cd /usr/local/bin/consul
#     sudo wget https://dl.bintray.com/mitchellh/consul/0.4.1_linux_amd64.zip
#     sudo unzip 0.4.1_linux_amd64.zip
#     sudo rm 0.4.1_linux_amd64.zip 

#     #install docker
#     sudo apt-get install -y -q lxc-docker

#     sudo mkdir -p /opt/tools/go/src/github.com/hashicorp
#     export GOPATH=/opt/tools/go/src
#     cd /opt/tools/go/src/github.com/hashicorp
#     git clone https://github.com/hashicorp/consul-template.git
#     cd consul-template
#     make

#     sudo sh -c "curl -L https://github.com/docker/fig/releases/download/1.0.1/fig-`uname -s`-`uname -m` > /usr/local/bin/fig"; chmod +x /usr/local/bin/fig

#     docker pull progrium/registrator
#     docker pull progrium/consul

#     echo 'PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/bin/packer:/usr/local/bin/consul:/opt/tools/go/bin"' > /etc/environment
# SCRIPT

# Vagrant.configure("2") do |config|
#     config.vm.box = "larryli/utopic64"
#     config.vm.network "private_network", ip: "192.168.33.10"
#     config.vm.synced_folder "~", ENV['HOME'], type: "nfs"
#     config.vm.provision "shell", inline: $script.gsub("PWD",ENV['PWD'])
#end

$machineroom = <<MRSCRIPT
sudo gpasswd -a vagrant docker
sudo route -n add 10.0.0.0/16 192.168.33.10

docker build -t pg93:0.0.1 PWD/gotooling
docker build -t gotooling:0.0.1 PWD/database

cd PWD
fig pull && fig build

MRSCRIPT

Vagrant.configure("2") do |config|
    config.vm.box = "iansmith/machineroom-base"
    config.vm.network "private_network", ip: "192.168.33.10"
    config.vm.network "forwarded_port", guest: 2375, host: 2375
    config.vm.synced_folder "~", ENV['HOME'], type: "nfs"
    config.vm.provision "shell", inline: $machineroom.gsub("PWD",ENV['PWD'])
end
