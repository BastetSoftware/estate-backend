FROM golang:1.20-alpine
WORKDIR /app

COPY *.go ./
COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN go build -o /estate

EXPOSE 8080

CMD [ "/estate" ]

