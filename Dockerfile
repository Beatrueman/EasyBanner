# 编译环境
FROM reg.redrock.team/library/golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o EasyBanner

# 运行环境
FROM debian:stable-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates\
    && update-ca-certificates

WORKDIR /app/EasyBanner

COPY --from=builder /app/EasyBanner .
COPY templates/ ./templates

EXPOSE 8080

CMD ["./EasyBanner"]
