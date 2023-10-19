FROM golang:alpine
LABEL authors="Dmirtii Morozov"

WORKDIR /app

COPY ./src /app
RUN go build -o main .
ENTRYPOINT ["./main"]