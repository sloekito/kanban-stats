FROM nordstrom/baseimage-ubuntu:14.04.2
MAINTAINER Innovation Platform Team "invcldtm@nordstrom.com"

ADD build/daily_trello /bin/daily_trello

ENTRYPOINT ["bin/daily_trello"]