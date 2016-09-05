# Starts a workspace with latest Go
# GOPATH set to /go
FROM golang

COPY . /go/src/github.com/rutigs/berkeley 
COPY slaves.json /go/src/github.com/rutigs/berkeley

RUN go install github.com/rutigs/berkeley
ENTRYPOINT ["/go/bin/berkeley"]

# Todo make this work
# CMD ["-m", "-addr=127.0.0.1:8080", "-slaves=/go/src/github.com/slaves.json"]

EXPOSE 8080
