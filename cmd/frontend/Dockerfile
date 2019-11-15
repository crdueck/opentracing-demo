FROM golang:1.13 AS golang-build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/frontend/main.go cmd/frontend/main.go
COPY pkg ./pkg
RUN CGO_ENABLED=0 go build -o frontend cmd/frontend/main.go

FROM alpine:3.8
RUN apk add --no-cache ca-certificates
COPY --from=golang-build /build/frontend /
CMD ["/frontend"]
