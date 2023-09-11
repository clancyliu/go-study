FROM golang:1.18 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY=https://goproxy.cn,direct \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY . ./
RUN go build ./src/main/k8s/main.go


FROM alpine
COPY --from=builder /app/main /app/main
EXPOSE 9080
CMD ["/app/main"]