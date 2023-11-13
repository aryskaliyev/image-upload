FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o web ./cmd/web

RUN apt-get update && apt-get install -y sqlite3

RUN sqlite3 mydb.db < forum.sql

EXPOSE 4000

CMD ["./web"] 
