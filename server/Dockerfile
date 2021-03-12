FROM golang:1.16 as builder

WORKDIR /go/app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/grpc ./server/

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /bin/grpc /bin/grpc

EXPOSE 2015
ENTRYPOINT ["/bin/grpc"]
