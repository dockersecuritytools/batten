FROM golang:latest
MAINTAINER amir@scalock.com
RUN mkdir -p /go/src/github.com/dockersecuritytools/
WORKDIR /go/src/github.com/dockersecuritytools/
RUN git clone https://github.com/dockersecuritytools/batten.git
WORKDIR /go/src/github.com/dockersecuritytools/batten
RUN make all
CMD ["/go/src/github.com/dockersecuritytools/batten/bin/batten", "check"]