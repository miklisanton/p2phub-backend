FROM dependencies AS builder
# Copy the application source code.
COPY ./main.go /go/src/app/main.go
# Build the application.
RUN CGO_ENABLED=0 \
go build -o /go/bin/web-server /go/src/app/main.go
ENTRYPOINT [ "/go/bin/web-server" ]

FROM alpine:latest
COPY --from=builder /go/bin/web-server /bin/web-server
WORKDIR /app
ENTRYPOINT [ "/bin/web-server" ]
