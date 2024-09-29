FROM dependencies AS builder
# Copy the application source code.
COPY ./main.go /go/src/app/main.go
# Build the application.
RUN CGO_ENABLED=0 \
go build -o /go/bin/observer /go/src/app/main.go
ENTRYPOINT [ "/go/bin/observer" ]

FROM alpine:latest
COPY --from=builder /go/bin/observer /bin/observer
WORKDIR /app
ENTRYPOINT [ "/bin/observer" ]
