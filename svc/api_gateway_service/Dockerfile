FROM golang:1.17 AS builder

WORKDIR /build

ADD ./main.go main.go
ADD ./service service
ADD ./go.mod go.mod

RUN apt-get update --allow-releaseinfo-change && apt upgrade -y
RUN go mod download all

ENV CGO_ENABLED=0 \
    GOOS=linux

RUN go build -o /api_gateway_service -a

FROM scratch
COPY --from=builder /api_gateway_service /app/api_gateway_service
EXPOSE 5000
ENTRYPOINT ["/app/api_gateway_service"]