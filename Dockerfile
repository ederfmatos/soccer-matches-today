FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY . .
RUN PATH="/go/bin:${PATH}" GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=amd64 && \
    go build -tags musl -ldflags '-s -w -extldflags "-static"' -o app .

FROM scratch
COPY --from=builder /build/app ./app
EXPOSE 8080
ENV TZ=America/Sao_Paulo
ENTRYPOINT ["./app"]