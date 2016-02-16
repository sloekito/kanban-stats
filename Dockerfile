FROM nordstrom/baseimage-alpine:3.2

MAINTAINER Innovation Platform Team "invcldtm@nordstrom.com"

ARG APP_NAME
ADD $APP_NAME /$APP_NAME

ENV APP_NAME $APP_NAME
ENTRYPOINT /$APP_NAME