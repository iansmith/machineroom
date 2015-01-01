FROM phusion/baseimage:0.9.15
MAINTAINER eng@igneous.io

ADD beta /beta
RUN chmod +x /beta

RUN mkdir /etc/service/beta
ADD run-beta /etc/service/beta/run
RUN chmod +x /etc/service/beta/run

CMD ["/sbin/my_init"]