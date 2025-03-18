FROM golang:1.24.1

WORKDIR /app

COPY ./backend/go.mod ./

RUN go mod tidy

COPY ./backend .

EXPOSE 8080

CMD ["go", "run", "."]

