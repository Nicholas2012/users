FROM golang:latest
COPY . /app
WORKDIR /app
RUN go build -o bin/ ./cmd/...

FROM debian:12-slim
RUN apt-get update && apt-get --no-install-recommends -y install ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=0 /app/bin/api-server /api-server
CMD [ "/api-server" ]
