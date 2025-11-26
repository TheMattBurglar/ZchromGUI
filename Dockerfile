# Dockerfile
# Build stage
FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /server ./web

# Runtime stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /server /server
# Copy static web assets so FileServer("./web/static") can find them at runtime
COPY --from=build /src/web/static /web/static
WORKDIR /
ENV PORT=8080
EXPOSE 8080
CMD ["/server"]
