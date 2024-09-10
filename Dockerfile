FROM golang:1.18 AS BUILDER

WORKDIR /app

COPY . .

RUN mkdir "bin"
RUN apt-get install postgresql-client

RUN go build -o bin/server ./cmd

EXPOSE 8080

CMD ["bin/server"]