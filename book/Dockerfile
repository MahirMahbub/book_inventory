FROM golang:1.18.0-alpine

RUN apk update && apk add --no-cache git \
    gcc \
    libc-dev \
    libffi-dev\
    openssh

RUN mkdir /app

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download golang.org/x/tools
RUN go mod download

COPY . .
#ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main
RUN go build -o main .

CMD ["./main"]
