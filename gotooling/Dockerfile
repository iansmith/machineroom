FROM phusion/baseimage:0.9.15
MAINTAINER eng@igneous.io

RUN DEBIAN_FRONTEND=noninteractive apt-get update -q
RUN DEBIAN_FRONTEND=noninteractive apt-get install -q -y build-essential curl git mercurial ack-grep

WORKDIR /tmp
RUN curl -s -O https://storage.googleapis.com/golang/go1.4.linux-amd64.tar.gz
RUN gunzip go1.4.linux-amd64.tar.gz

RUN mkdir -p /opt
RUN mkdir -p /go/src/github.com/igneous-systems
RUN mkdir -p /go/bin
RUN sudo mkdir -p /go/pkg
     
WORKDIR /opt
RUN tar xf /tmp/go1.4.linux-amd64.tar

RUN chmod 666 /etc/profile
RUN echo 'export GOPATH=/go' >> /etc/profile
RUN echo 'export GOROOT=/opt/go' >> /etc/profile
RUN echo 'export PATH=/opt/go/bin:/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin' >> /etc/profile
RUN sudo chmod 644 /etc/profile

          
RUN curl -sL https://deb.nodesource.com/setup | sudo bash -
RUN DEBIAN_FRONTEND=noninteractive apt-get install -q -y nodejs
RUN npm install -g grunt-cli

ENV GOPATH=/go
ENV GOROOT=/opt/go
ENV PATH=/opt/go/bin:/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
     
RUN go get github.com/gopherjs/gopherjs
RUN go get github.com/tools/godep
      