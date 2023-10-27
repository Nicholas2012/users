build:
	go build -o bin/ ./cmd/...

docker-build:
	docker-compose build

up: docker-build
	docker-compose up -d

down:
	docker-compose down
