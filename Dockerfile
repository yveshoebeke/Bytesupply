FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go get ./...
RUN GOOS=linux go build -ldflags="-s -w" -o /go/bin/bytesupply-app ./main.go
RUN mkdir /go/bin/log
RUN ["touch", "/go/bin/log/bytesupply.log"]
RUN ["chmod", "a+w", "/go/bin/log/bytesupply.log"]
RUN mkdir -p /go/bin/data/qTurHm/
RUN mkdir -p /go/bin/data/messages/
COPY ./static/ /go/bin/static/
COPY ./sitemap.xml /go/bin/sitemap.xml
COPY ./robots.txt /go/bin/robots.txt
RUN ["chmod", "+x", "/go/bin"]
FROM alpine:3.9
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
FROM build 
WORKDIR /go/bin
EXPOSE 80
ENTRYPOINT bytesupply-app --port 80