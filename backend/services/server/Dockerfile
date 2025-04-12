# Dockerfile for the server microservice
FROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go build -o server-service ./cmd
CMD [ "./server-service" ]
