FROM golang:1.7

MAINTAINER Nicholas Page

RUN export PATH=$PATH:~/go/bin
RUN export GOROOT=~/go

RUN mkdir -p /go/src/gitlab.com/AndBobsYourUncle/server-utilities && mkdir /app
COPY . /go/src/gitlab.com/AndBobsYourUncle/server-utilities

WORKDIR /go/src/gitlab.com/AndBobsYourUncle/server-utilities

RUN go get -v
RUN go build -v -o server-utilities && cp server-utilities /app

WORKDIR /app
RUN rm -rf /go/src/gitlab.com/AndBobsYourUncle/server-utilities

EXPOSE 443

CMD ["/app/server-utilities"]