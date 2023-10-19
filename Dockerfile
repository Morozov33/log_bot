FROM golang:1.21-alpine
LABEL authors="Dmirtii Morozov"

WORKDIR /app

COPY . /app
RUN go build -o main .
ENTRYPOINT ["./main"]