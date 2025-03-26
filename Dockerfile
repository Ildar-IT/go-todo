FROM golang:latest

ENV GOPATH=/

COPY ./ ./

#install psql
RUN apt-get update
RUN apt-get -y install postgresql-client 

#execut
RUN chmod +x wait-for-postgres.sh


#build app
RUN go mod download
RUN go build -o main ./cmd/main.go
RUN go build -o ./migrations/main ./migrations/main.go

#execut
RUN chmod +x ./main
RUN chmod +x ./migrations/main

CMD ["./main"]

