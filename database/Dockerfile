FROM phusion/baseimage:0.9.15
MAINTAINER eng@igneous.io

RUN DEBIAN_FRONTEND=noninteractive apt-get update -q
RUN DEBIAN_FRONTEND=noninteractive apt-get install -q -y curl postgresql-9.3 wget zip
RUN usermod -a -G sudo postgres

ADD run-postgres /etc/service/postgres/run
RUN chmod +x /etc/service/postgres/run

ADD provision.sh /tmp/provision.sh
RUN chmod +x /tmp/provision.sh

ADD provision.sql /tmp/provision.sql
RUN chmod +x /tmp/provision.sh

RUN echo 'host all  all    0.0.0.0/0  trust' >> /etc/postgresql/9.3/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/9.3/main/postgresql.conf

RUN /tmp/provision.sh
CMD ["/sbin/my_init"]