#!/bin/bash

set -e -o pipefail


sudo systemctl disable --now systemd-timesyncd
sudo apt-add-repository -yu ppa:maas/3.4
sudo apt install maas -y
apt update -y
apt install -y postgresql

export MAAS_DBUSER=maasuser
export MAAS_DBPASS=maaspass
export MAAS_DBNAME=maas
export HOSTNAME=localhost

sudo -i -u postgres psql -c "CREATE USER \"$MAAS_DBUSER\" WITH ENCRYPTED PASSWORD '$MAAS_DBPASS'"

sudo -i -u postgres createdb -O "$MAAS_DBUSER" "$MAAS_DBNAME"

echo "host    $MAAS_DBNAME    $MAAS_DBUSER    0/0     md5" /etc/postgresql/14/main/pg_hba.conf

sudo maas init region+rack --database-uri "postgres://$MAAS_DBUSER:$MAAS_DBPASS@$HOSTNAME/$MAAS_DBNAME"
