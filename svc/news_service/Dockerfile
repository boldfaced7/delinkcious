FROM golang:1.17 AS builder

WORKDIR /build

ADD ./main.go main.go
ADD ./service service
ADD ./go.mod go.mod

RUN apt-get update --allow-releaseinfo-change && apt upgrade -y
RUN go mod download all

ENV CGO_ENABLED=0 \
    GOOS=linux

RUN go build -o /news_service -a

FROM scratch
COPY --from=builder /news_service /app/news_service
EXPOSE 7070
ENTRYPOINT ["/app/news_service"]