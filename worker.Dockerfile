FROM golang:1.24.1

WORKDIR /app

COPY ./backend/go.mod ./backend/
COPY ./backend/go.sum ./backend/

WORKDIR /app/backend

RUN go mod tidy

COPY ./backend/worker_node/ ./worker_node/

WORKDIR ./worker_node

EXPOSE 8080

# Run the application
CMD ["go", "run", "."]