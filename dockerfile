FROM golang:alpine

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o cmd/main.go

CMD [ "./cmd/main.go" ]