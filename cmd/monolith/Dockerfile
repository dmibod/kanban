FROM golang:1.12.4-alpine3.9 as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -mod vendor -o main ./monolith/main.go
FROM scratch
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]