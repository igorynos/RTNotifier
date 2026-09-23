FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /app ./cmd/rtnotifier
FROM scratch
COPY --from=build /app /app
ENTRYPOINT ["/app"]
