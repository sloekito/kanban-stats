FROM nordstrom/baseimage-alpine:3.2

MAINTAINER Innovation Platform Team "invcldtm@nordstrom.com"

ADD hello-world /hello-world

ENTRYPOINT /hello-world