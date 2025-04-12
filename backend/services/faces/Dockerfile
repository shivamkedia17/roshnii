# Dockerfile for the faces microservice
FROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go build -o faces-service ./cmd
CMD [ "./faces-service" ]
