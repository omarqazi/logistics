# Start with golang base image
FROM golang
MAINTAINER Omar Qazi (omar@bqe.com)

# Compile latest source
ADD . /go/src/github.com/omarqazi/logistics
RUN go get github.com/omarqazi/logistics
RUN go get bitbucket.org/liamstask/goose/cmd/goose
RUN go install github.com/omarqazi/logistics

WORKDIR /go/src/github.com/omarqazi/logistics

CMD /go/bin/logistics
EXPOSE 8080
