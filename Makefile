.PHONY: deps build
build: network
	docker compose up --build -d
network: deps
	docker network create app-network || true
deps:
	docker build -t dependencies -f ./dependencies.Dockerfile .

