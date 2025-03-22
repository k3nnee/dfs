FROM golang:1.24.1

WORKDIR /app

COPY ./backend/go.mod ./backend/go.sum /app/

RUN go mod tidy

COPY ./backend/ /app/

WORKDIR /app/master_node

EXPOSE 8080

CMD ["go", "run", "."]