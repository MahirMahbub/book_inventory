FROM golang:1.18.0-alpine

RUN apk update && apk add --no-cache git \
    gcc \
    libc-dev \
    libffi-dev\
    make openssh-client\
    curl

RUN mkdir /app

WORKDIR /app

COPY go.mod go.sum ./
RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

RUN air -v

RUN go mod download
RUN go get github.com/cosmtrek/air
COPY . .

RUN go build -o tmpmain.exe .
EXPOSE 5002
CMD air
