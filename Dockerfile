FROM golang:1.13.3-buster as builder
ENV GO111MODULE=on
ENV CGO_ENABLED=0
WORKDIR /usr/src/
COPY . /usr/src
RUN go build -v -o "bin/github-commenter" *.go

FROM alpine:3.12
RUN apk add --no-cache ca-certificates
COPY --from=builder /usr/src/bin/* /usr/bin/
ENV PATH $PATH:/usr/bin
ENTRYPOINT ["github-commenter"]
