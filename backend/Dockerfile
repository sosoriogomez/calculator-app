FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN go build -o /calculator-api ./cmd/server

FROM alpine:3.20
RUN adduser -D -H appuser
COPY --from=build /calculator-api /calculator-api
USER appuser
EXPOSE 8080
ENTRYPOINT ["/calculator-api"]
