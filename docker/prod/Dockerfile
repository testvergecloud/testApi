FROM golang:1.21-alpine as build-stage

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Release stage
FROM alpine:latest as release-stage
COPY --from=build-stage /app/main /app/main
COPY --from=build-stage /app/app.env /app/app.env
WORKDIR /app

EXPOSE  3000

RUN chmod +x main

CMD ["/app/main"]
