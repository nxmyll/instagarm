FROM golang:1.21.6 AS build
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app ./cmd/web

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/app .
COPY --from=build /app/static ./static
COPY --from=build /app/templates ./templates
EXPOSE 8080
CMD ["./app"]
