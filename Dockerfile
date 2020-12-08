FROM golang:latest

RUN mkdir -p /ganjapp

WORKDIR /ganjapp

ADD . /ganjapp

RUN mkdir -p /ganjapp/data

VOLUME /ganjapp/data

RUN mkdir -p /ganjapp/temp

VOLUME /ganjapp/temp

RUN rm -rf /ganjapp/.git*

RUN go build ./ganjapp.go

RUN rm -rf /ganjapp/app

RUN rm -rf /ganjapp/artwork

RUN rm -rf /ganjapp/controllers

RUN rm -rf /ganjapp/middleware

RUN rm -rf /ganjapp/models

RUN rm -rf /ganjapp/utilities

RUN rm -rf /ganjapp/Dockerfile

RUN rm -rf /ganjapp/go.*

RUN rm -rf /ganjapp/*.go

ENV GIN_MODE release

ENV PORT 8080

ENV GANJAPP_ROOT /ganjapp

ENV GANJAPP_DATABASE_TYPE sqlite

ENV GANJAPP_DATABASE_DSN ""

ENV GANJAPP_COOKIE_KEY setme

ENV GANJAPP_JWT_KEY setme

ENV GANJAPP_JWT_AUDIENCE servers.ganj.app

ENV GANJAPP_S3_ENDPOINT localhost:9001

ENV GANJAPP_S3_ACCESSKEY setme

ENV GANJAPP_S3_SECRET_KEY setme

ENV GANJAPP_S3_HTTPS false

EXPOSE ${PORT:+8080}/tcp

CMD ["./ganjapp"]