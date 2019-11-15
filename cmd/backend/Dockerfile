FROM golang:1.13 AS golang-build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/backend/main.go cmd/backend/main.go
COPY pkg ./pkg
RUN CGO_ENABLED=0 go build -o backend cmd/backend/main.go

FROM alpine:3.8
RUN apk add --no-cache ca-certificates
COPY --from=golang-build /build/backend /
CMD ["/backend"]
