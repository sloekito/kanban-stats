FROM nordstrom/awscli:1.8.1
MAINTAINER Innovation Platform Team "invcldtm@nordstrom.com"

ADD build/kanban-stats /bin/
ADD scripts/createdb.sh /bin/
ADD scripts/influxdb-backup.sh /bin/

ENV INFLUXDB_RELEASE 0.9.4.1
ADD dist/influxdb_${INFLUXDB_RELEASE}_amd64.deb /
RUN dpkg -i /influxdb_${INFLUXDB_RELEASE}_amd64.deb
    
RUN apt-get update -qy \
 && apt-get install -qy --no-install-recommends --no-install-suggests \
      curl \
      jq \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
 
ENTRYPOINT ["bin/kanban-stats"]