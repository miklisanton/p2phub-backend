FROM golang:latest

COPY . .

CMD ["go", "test", "./internal/tasks/...", "-v"]
