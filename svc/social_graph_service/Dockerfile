FROM golang:1.17 AS builder

WORKDIR /build

ADD ./main.go main.go
ADD ./service service
ADD ./go.mod go.mod

RUN apt-get update --allow-releaseinfo-change && apt upgrade -y
RUN go mod download all

ENV CGO_ENABLED=0 \
    GOOS=linux

RUN go build -o /social_graph_service -a

FROM scratch
COPY --from=builder /social_graph_service /app/social_graph_service
EXPOSE 9090
ENTRYPOINT ["/app/social_graph_service"]