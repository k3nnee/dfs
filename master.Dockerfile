FROM golang:1.24.1

WORKDIR /app

COPY ./backend/go.mod ./backend/

WORKDIR /app/backend

RUN go mod tidy

COPY ./backend/master_node/ ./master_node/

WORKDIR ./master_node

EXPOSE 8080

# Run the application
CMD ["go", "run", "."]