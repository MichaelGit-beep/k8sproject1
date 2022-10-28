FROM golang:1.19.1-alpine3.16 AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o main

FROM golang:1.19.1-alpine3.16
WORKDIR /app 
COPY --from=builder /app/main ./
RUN ls
CMD [ "/app/main" ]