FROM golang:1.16

RUN mkdir -p /usr/src/app/
ADD . /usr/src/app/
WORKDIR /usr/src/app/

EXPOSE 50051

COPY . /usr/src/app/
CMD go run .server.go