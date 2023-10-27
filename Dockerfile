FROM golang:latest
COPY . /app
WORKDIR /app
RUN go install ./...
CMD "api-server"