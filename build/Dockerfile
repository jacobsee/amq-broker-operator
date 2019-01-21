FROM alpine:3.8

RUN apk upgrade --update --no-cache

USER nobody

ADD build/_output/bin/amq-broker-operator /usr/local/bin/amq-broker-operator
