FROM golang

MAINTAINER Abhi Yerra <abhi@berkeley.edu>

ADD . /go/src/github.com/forestly/treely

RUN cd /go/src/github.com/forestly/treely && go get ./...
RUN go install github.com/forestly/treely

WORKDIR /go/src/github.com/forestly/treely

ENTRYPOINT /go/bin/treely serve

EXPOSE 3001
