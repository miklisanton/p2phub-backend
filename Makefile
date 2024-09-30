.PHONY: deps
deps:
	docker build -t dependencies -f ./dependencies.Dockerfile .
.PHONY: build
build:
	docker compose up --build -d

