FROM dependencies AS builder
# Copy the application source code.
COPY ./main.go /go/src/app/main.go
# Build the application.
RUN CGO_ENABLED=0 \
go build -o /go/bin/tg-bot /go/src/app/main.go
ENTRYPOINT [ "/go/bin/tg-bot" ]

FROM alpine:latest
COPY --from=builder /go/bin/tg-bot /bin/tg-bot
WORKDIR /app
ENTRYPOINT [ "/bin/tg-bot" ]
