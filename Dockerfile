FROM golang:alpine AS builder

# Git is required for fetching the dependencies.
RUN apk add --update --no-cache git

WORKDIR /work
ADD demo-parcer .

RUN go get -d -v
RUN go build -o ./demo-parcer

FROM scratch
COPY --from=builder /work/demo-parcer /go/bin/demo-parcer

ENTRYPOINT ["/go/bin/demo-parcer"]
