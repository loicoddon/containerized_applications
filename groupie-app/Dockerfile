FROM golang

WORKDIR /app

COPY . /app

EXPOSE 80

RUN go get github.com/go-sql-driver/mysql

CMD ["go","run","main.go"]