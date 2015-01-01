#!/bin/sh

sudo -u postgres /usr/lib/postgresql/9.3/bin/pg_ctl -w -o "--config-file=/etc/postgresql/9.3/main/postgresql.conf" -D /var/lib/postgresql/9.3/main start    
sudo -u postgres /usr/lib/postgresql/9.3/bin/createdb alpha
sudo -u postgres psql < /tmp/provision.sql
sudo -u postgres /usr/lib/postgresql/9.3/bin/pg_ctl -o "--config-file=/etc/postgresql/9.3/main/postgresql.conf" -D /var/lib/postgresql/9.3/main stop

