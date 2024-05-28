FROM golang:1.16-alpine as builder

WORKDIR /app
COPY . .

RUN go build -o iptc ./cmd

FROM chainguard/go:latest
WORKDIR /root/
COPY --from=builder /app/iptc .
COPY .env .

EXPOSE 8080
CMD ["./iptc"]
