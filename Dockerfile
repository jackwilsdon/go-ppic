# Build the binary.
FROM golang:alpine AS builder
COPY . /go/src/github.com/jackwilsdon/go-ppic
WORKDIR /go/src/github.com/jackwilsdon/go-ppic
RUN apk add --no-cache git
RUN go get -d -v github.com/jackwilsdon/go-ppic/cmd/ppicd/... && \
    CGO_ENABLED=0 go install -ldflags "-extldflags -static" github.com/jackwilsdon/go-ppic/cmd/ppicd

# Create a blank image containing the server.
FROM scratch
COPY --from=builder /go/bin/ppicd /ppicd

ENV PORT=3000
EXPOSE 3000

ENTRYPOINT ["/ppicd"]
