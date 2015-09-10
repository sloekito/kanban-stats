FROM nordstrom/awscli:1.8.1
MAINTAINER Innovation Platform Team "invcldtm@nordstrom.com"

ADD build/daily_trello /bin/
ADD scripts/influxdb-backup.sh /bin/

RUN apt-get update -qy \
 && apt-get install -qy --no-install-recommends --no-install-suggests \
      curl \
      jq \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
 
ENTRYPOINT ["bin/daily_trello"]