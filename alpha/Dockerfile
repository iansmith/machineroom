FROM phusion/baseimage:0.9.15
MAINTAINER eng@igneous.io


RUN mkdir /etc/service/alpha
ADD run-alpha /etc/service/alpha/run
RUN chmod +x /etc/service/alpha/run

ADD alpha /alpha
RUN chmod +x /alpha


CMD ["/sbin/my_init"]