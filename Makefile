.PHONY: deps build
build: deps
	docker compose up --build -d
deps:
	docker build -t dependencies -f ./dependencies.Dockerfile .

