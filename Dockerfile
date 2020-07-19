FROM golang:1.14-stretch

RUN mkdir /app

WORKDIR /app

COPY . /app

ENV PORT=8080

RUN apt-get update && apt-get install poppler-utils -y

RUN go build -mod=readonly -o server

CMD [ "./server" ]