FROM golang:1.22-alpine

WORKDIR /app

#using air for live reloading
RUN go install github.com/cosmtrek/air@latest

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

EXPOSE 3000

CMD ["air", "-c", ".air.toml"]
