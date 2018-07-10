# build stage
FROM golang:1.10.2-alpine AS build
ARG dir=/go/src/github.com/cohhei/go-to-the-handson/04
ADD . ${dir}
RUN apk update && \
    apk add --virtual build-dependencies build-base git && \
    cd ${dir} && \
    go get -u github.com/lib/pq && \
    go build -o todo-api

# final stage
FROM alpine:3.7
ARG dir=/go/src/github.com/cohhei/go-to-the-handson/04
WORKDIR /app
COPY --from=build ${dir}/todo-api /app/
EXPOSE 8080
CMD ./todo-api