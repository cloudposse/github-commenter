FROM golang:1.11.2 as builder
RUN mkdir -p /go/src/github.com/cloudposse/github-commenter
WORKDIR /go/src/github.com/cloudposse/github-commenter
COPY . .
RUN go get && CGO_ENABLED=0 go build -v -o "./dist/bin/github-commenter" *.go


FROM alpine:3.8
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/cloudposse/github-commenter/dist/bin/github-commenter /usr/bin/github-commenter
ENV PATH $PATH:/usr/bin
ENTRYPOINT ["github-commenter"]
