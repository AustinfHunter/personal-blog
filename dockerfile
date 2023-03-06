ARG PORT 
FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go mod tidy
RUN go build -o ./bin/blog-backend

EXPOSE $PORT

CMD ["./bin/blog-backend "]
